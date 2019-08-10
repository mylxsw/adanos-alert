package graphql

import (
	lib "github.com/mylxsw/adanos-alert/pkg/graphql"
)

type Graphql interface {
	Register(builder lib.GraphQL)
}
