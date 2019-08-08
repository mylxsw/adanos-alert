package graphql

import (
	"github.com/graphql-go/graphql"
)

type WelcomeObject struct {
	Message string `json:"message"`
}

func NewWelcomeObject() *WelcomeObject {
	return &WelcomeObject{}
}

func (w *WelcomeObject) Register(query *graphql.Object, mutation *graphql.Object) {
	query.AddFieldConfig("welcome", &graphql.Field{
		Type: graphql.NewObject(graphql.ObjectConfig{
			Name:   "WelcomeObject",
			Fields: graphql.BindFields(WelcomeObject{}),
		}),
		Args: graphql.FieldConfigArgument{
			"name": &graphql.ArgumentConfig{
				Type:         graphql.String,
				DefaultValue: "Graphql",
			},
		},
		Resolve: w.Hello,
	})
}

func (w *WelcomeObject) Hello(p graphql.ResolveParams) (interface{}, error) {
	name := p.Args["name"].(string)
	return WelcomeObject{Message: name}, nil
}
