package api

import (
	"github.com/99designs/gqlgen/handler"
	"github.com/gorilla/mux"
	"github.com/mylxsw/adanos-alert/api/graphql"
	"github.com/mylxsw/container"
)

func graphqlRouters(cc *container.Container) func(router *mux.Router) {
	return func(router *mux.Router) {
		// router.Handle("/graphql", graphql.HTTPHandler(schema))
		router.Handle("/graphql", handler.GraphQL(graphql.NewExecutableSchema(graphql.Config{Resolvers: graphql.NewResolver(cc)})))
	}
}
