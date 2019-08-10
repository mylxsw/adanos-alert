package api

import (
	"github.com/gorilla/mux"
	gql "github.com/mylxsw/adanos-alert/api/graphql"
	lib "github.com/mylxsw/adanos-alert/pkg/graphql"
	"github.com/mylxsw/go-toolkit/container"
)

func graphqlRouters(cc *container.Container) func(router *mux.Router) {

	gqlib := lib.NewBuilder()

	gql.NewWelcomeObject().Register(gqlib)
	gql.NewRuleObject(cc).Register(gqlib)

	return func(router *mux.Router) {
		// router.Handle("/graphql", graphql.HTTPHandler(schema))
		router.Handle("/graphql", gqlib.Build())
	}
}
