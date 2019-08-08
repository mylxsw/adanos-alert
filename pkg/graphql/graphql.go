package graphql

import (
	"reflect"
	"strings"

	"github.com/graphql-go/graphql"
)

// CreateField create a GraphQL field
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

				argRefVal.Elem().Field(i).Set(reflect.ValueOf(p.Args[tag]))
			}

			res := resolverRefVal.Call([]reflect.Value{argRefVal.Elem()})
			if res[1].IsNil() {
				return res[0].Interface(), nil
			}

			return res[0].Interface(), res[1].Interface().(error)
		},
	}
}

// BindArgs create a FieldConfigArgument from object
func BindArgs(obj interface{}, tags ...string) graphql.FieldConfigArgument {
	return graphql.BindArg(obj, tags...)
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
