package graphql

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

var defaultField = graphql.Fields{
	"hello": &graphql.Field{
		Type: graphql.String,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return "Hello, world", nil
		},
	},
}

type GraphQL interface {
	Query(name string, f *graphql.Field)
	Mutation(name string, f *graphql.Field)
	Build() *handler.Handler
}

type Builder struct {
	query    *graphql.Object
	mutation *graphql.Object
}

// NewBuilder create a new Builder instance
func NewBuilder() GraphQL {
	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name:   "Mutation",
		Fields: defaultField,
	})

	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name:   "Query",
		Fields: defaultField,
	})

	return &Builder{
		query:    rootQuery,
		mutation: rootMutation,
	}
}

func (gq *Builder) Query(name string, f *graphql.Field) {
	gq.query.AddFieldConfig(name, f)
}

func (gq *Builder) Mutation(name string, f *graphql.Field) {
	gq.mutation.AddFieldConfig(name, f)
}

func (gq *Builder) Build() *handler.Handler {
	schemaConfig := graphql.SchemaConfig{Query: gq.query, Mutation: gq.mutation}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		panic(err)
	}

	return handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})
}

// CreateField create a GraphQLBuilder field
// resolver: func (obj T) (interface{}, error)
func CreateField(typ graphql.Output, args graphql.FieldConfigArgument, resolver interface{}) *graphql.Field {
	resolverRefVal := reflect.ValueOf(resolver)
	if resolverRefVal.Type().NumIn() != 1 {
		panic("invalid resolver: the number of arguments must be 1")
	}

	if resolverRefVal.Type().NumOut() != 2 {
		panic("invalid resolver: the number of return values must be 2")
	}

	argRefTyp := resolverRefVal.Type().In(0)
	return &graphql.Field{
		Type: typ,
		Args: args,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			argRefVal := reflect.New(argRefTyp)
			for i := 0; i < argRefTyp.NumField(); i++ {
				field := argRefTyp.Field(i)
				tag := extractTag(field.Tag)
				if tag == "" || tag == "-" {
					continue
				}

				if p.Args[tag] == nil {
					continue
				}

				fmt.Println(p.Args[tag])

				switch argRefVal.Elem().Field(i).Kind() {
				case reflect.Slice:
					s := reflect.MakeSlice(field.Type, 0, 0)
					for _, e := range p.Args[tag].([]interface{}) {
						s = reflect.AppendSlice(s, reflect.Indirect(reflect.ValueOf(e)))
					}

					argRefVal.Elem().Field(i).Set(s)
				default:
					argRefVal.Elem().Field(i).Set(reflect.ValueOf(p.Args[tag]))
				}
			}

			res := resolverRefVal.Call([]reflect.Value{argRefVal.Elem()})
			if res[1].IsNil() {
				return res[0].Interface(), nil
			}

			return res[0].Interface(), res[1].Interface().(error)
		},
	}
}

// Object create a new graphql.Object
func Object(obj interface{}) *graphql.Object {
	ref := reflect.TypeOf(obj)
	var name string
	if ref.Kind() == reflect.Ptr {
		name = ref.Elem().Name()
	} else {
		name = ref.Name()
	}

	return graphql.NewObject(graphql.ObjectConfig{
		Name:   name,
		Fields: BindFields(obj),
	})
}

// BindArgs create a FieldConfigArgument from object
func BindArgs(obj interface{}, tags ...string) graphql.FieldConfigArgument {
	if len(tags) == 0 {
		objRef := reflect.TypeOf(obj)
		for i := 0; i < objRef.NumField(); i++ {
			tag := extractTag(objRef.Field(i).Tag)
			if tag == "" || tag == "-" {
				continue
			}

			tags = append(tags, tag)
		}
	}

	return BindArg(obj, tags...)
}

// lazy way of binding args
func BindArg(obj interface{}, tags ...string) graphql.FieldConfigArgument {
	v := reflect.Indirect(reflect.ValueOf(obj))
	var config = make(graphql.FieldConfigArgument)
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)

		mytag := extractTag(field.Tag)
		if inArray(tags, mytag) {
			config[mytag] = &graphql.ArgumentConfig{
				Type: getGraphType(field.Type),
			}
		}
	}
	return config
}

func inArray(slice interface{}, item interface{}) bool {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("inArray() given a non-slice type")
	}

	for i := 0; i < s.Len(); i++ {
		if reflect.DeepEqual(item, s.Index(i).Interface()) {
			return true
		}
	}
	return false
}

// ParseResolveParams parse ResolveParams to object
func ParseResolveParams(p graphql.ResolveParams, res interface{}) {
	refVal := reflect.ValueOf(res)
	refType := refVal.Type()

	for i := 0; i < refType.Elem().NumField(); i++ {
		field := refType.Elem().Field(i)
		tag := extractTag(field.Tag)
		if tag == "" || tag == "-" {
			continue
		}

		if p.Args[tag] == nil {
			continue
		}

		refVal.Elem().Field(i).Set(reflect.ValueOf(p.Args[tag]))
	}
}

func extractTag(tag reflect.StructTag) string {
	t := tag.Get(graphql.TAG)
	if t != "" {
		t = strings.Split(t, ",")[0]
	}
	return t
}

// can't take recursive slice type
// e.g
// type Person struct{
//	Friends []Person
// }
// it will throw panic stack-overflow
func BindFields(obj interface{}) graphql.Fields {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	fields := make(map[string]*graphql.Field)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		tag := extractTag(field.Tag)
		if tag == "-" || tag == "" {
			continue
		}

		fieldType := field.Type

		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}

		var graphType graphql.Output
		var graphArgs graphql.FieldConfigArgument

		// check whether there has a XxxResolver function
		resolverName := field.Name + "Resolver"
		method, ok := t.MethodByName(resolverName)
		if ok {
			ft := method.Func.Type()
			if ft.NumOut() != 2 {
				panic(fmt.Sprintf("resolver %s.%s must have two return value", t.Name(), resolverName))
			}

			// Parse return value
			out0 := ft.Out(0)
			if out0.Kind() == reflect.Struct {
				structFields := BindFields(reflect.New(out0).Interface())
				graphType = graphql.NewObject(graphql.ObjectConfig{
					Name:   tag,
					Fields: structFields,
				})
			} else {
				graphType = getGraphType(out0)
			}

			// Parse input args
			if ft.NumIn() > 2 {
				panic(fmt.Errorf("resolver %s.%s must have no more than 1 argument", t.Name(), resolverName))
			}
			if ft.NumIn() == 2 {
				arg0 := ft.In(1)
				graphArgs = BindArgs(reflect.New(arg0).Interface())
			}

			fields[tag] = &graphql.Field{
				Type: graphType,
				Args: graphArgs,
				Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
					var args = make([]reflect.Value, 0)
					args = append(args, reflect.ValueOf(p.Source))

					if ft.NumIn() > 1 {
						argRefVal := reflect.New(ft.In(1))
						for i := 0; i < ft.In(1).NumField(); i++ {
							field := ft.In(1).Field(i)
							tag := extractTag(field.Tag)
							if tag == "" || tag == "-" {
								continue
							}

							if p.Args[tag] == nil {
								continue
							}

							argRefVal.Elem().Field(i).Set(reflect.ValueOf(p.Args[tag]))
						}

						args = append(args, argRefVal.Elem())
					}

					res := method.Func.Call(args)
					if res[1].IsNil() {
						return res[0].Interface(), nil
					}

					return res[0].Interface(), res[1].Interface().(error)
				},
			}
		} else {
			if fieldType.Kind() == reflect.Struct {
				structFields := BindFields(v.Field(i).Interface())
				graphType = graphql.NewObject(graphql.ObjectConfig{
					Name:   tag,
					Fields: structFields,
				})
			} else {
				graphType = getGraphType(fieldType)
			}

			fields[tag] = &graphql.Field{
				Type: graphType,
				Args: graphArgs,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return extractValue(tag, p.Source), nil
				},
			}
		}
	}

	return fields
}

func getGraphType(tipe reflect.Type) graphql.Output {
	kind := tipe.Kind()
	switch kind {
	case reflect.String:
		return graphql.String
	case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64:
		return graphql.Int
	case reflect.Float32, reflect.Float64:
		return graphql.Float
	case reflect.Bool:
		return graphql.Boolean
	case reflect.Slice:
		return getGraphList(tipe)
	case reflect.Struct:
		return getGraphList(tipe)
	case reflect.Ptr:
		return getGraphList(tipe.Elem())
	}

	return graphql.String
}

func getGraphList(tipe reflect.Type) *graphql.List {
	if tipe.Kind() == reflect.Slice {
		switch tipe.Elem().Kind() {
		case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64:
			return graphql.NewList(graphql.Int)
		case reflect.Bool:
			return graphql.NewList(graphql.Boolean)
		case reflect.Float32, reflect.Float64:
			return graphql.NewList(graphql.Float)
		case reflect.String:
			return graphql.NewList(graphql.String)
		}
	}
	// finally bind object
	t := reflect.New(tipe.Elem())
	name := strings.Replace(fmt.Sprint(tipe.Elem()), ".", "_", -1)
	obj := graphql.NewObject(graphql.ObjectConfig{
		Name:   name,
		Fields: BindFields(t.Elem().Interface()),
	})

	fmt.Println(name)

	return graphql.NewList(obj)
}

func appendFields(dest, origin graphql.Fields) graphql.Fields {
	for key, value := range origin {
		dest[key] = value
	}
	return dest
}

func extractValue(originTag string, obj interface{}) interface{} {
	val := reflect.Indirect(reflect.ValueOf(obj))

	for j := 0; j < val.NumField(); j++ {
		field := val.Type().Field(j)
		if field.Type.Kind() == reflect.Struct {
			res := extractValue(originTag, val.Field(j).Interface())
			if res != nil {
				return res
			}
		}

		if originTag == extractTag(field.Tag) {
			return reflect.Indirect(val.Field(j)).Interface()
		}
	}
	return nil
}
