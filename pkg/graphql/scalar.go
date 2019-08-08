package graphql

import (
	"strconv"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

var Int64 = graphql.NewScalar(graphql.ScalarConfig{
	Name:        "Int64",
	Description: "The Int64 scalar type represents an int64 value",
	// Serialize serializes `Int64` to `int64`
	Serialize: func(value interface{}) interface{} {
		return value.(int64)
	},
	// ParseValue parses GraphQL variables from `int64` to `Int64`
	ParseValue: func(value interface{}) interface{} {
		return value.(int64)
	},
	// ParseLiteral parses GraphQL AST value to `Int64`
	ParseLiteral: func(valueAST ast.Value) interface{} {
		switch valueAST := valueAST.(type) {
		case *ast.IntValue:
			val, _ := strconv.Atoi(valueAST.Value)
			return val
		default:
			return nil
		}
	},
})
