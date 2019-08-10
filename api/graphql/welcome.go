package graphql

import (
	"github.com/graphql-go/graphql"
	lib "github.com/mylxsw/adanos-alert/pkg/graphql"
)

type WelcomeObject struct {
	Message string `json:"message"`
}

func NewWelcomeObject() *WelcomeObject {
	return &WelcomeObject{}
}

func (w *WelcomeObject) Register(builder lib.GraphQL) {
	builder.Query("welcome", w.Hello())
}

func (w *WelcomeObject) Hello() *graphql.Field {
	return lib.CreateField(
		lib.Object(w),
		lib.BindArgs(struct {Name string `json:"name"`}{}),
		func(arg struct{ Name string `json:"name"` }) (interface{}, error) {
			return WelcomeObject{Message: arg.Name}, nil
		},
	)
}
