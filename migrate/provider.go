package migrate

import (
	"github.com/mylxsw/glacier/infra"
)

type Provider struct {
}

func (s Provider) Register(app infra.Binder) {
}

func (s Provider) Boot(app infra.Resolver) {
	app.MustResolve(initPredefinedTemplates)
}
