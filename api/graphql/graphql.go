package graphql

import (
	"github.com/graphql-go/graphql"
)

type Graphql interface {
	Register(query *graphql.Object, mutation *graphql.Object)
}
