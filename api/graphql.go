package api

import (
	"github.com/gorilla/mux"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	gql "github.com/mylxsw/adanos-alert/api/graphql"
	"github.com/mylxsw/go-toolkit/container"
)

func graphqlRouters(cc *container.Container) func(router *mux.Router) {
	var defaultField = graphql.Field{
		Type: graphql.String,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return "Hello, world", nil
		},
	}

	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"hello": &defaultField,
		},
	})

	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"hello": &defaultField,
		},
	})

	gql.NewWelcomeObject().Register(rootQuery, rootMutation)
	gql.NewRuleObject(cc).Register(rootQuery, rootMutation)

	schemaConfig := graphql.SchemaConfig{Query: rootQuery, Mutation: rootMutation}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		panic(err)
	}

	return func(router *mux.Router) {
		// router.Handle("/graphql", graphql.HTTPHandler(schema))
		router.Handle("/graphql", handler.New(&handler.Config{
			Schema:   &schema,
			Pretty:   true,
			GraphiQL: true,
		}))
	}
}
