package migrate

import (
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier/infra"
)

type ServiceProvider struct {
}

func (s ServiceProvider) Register(app container.Container) {
}

func (s ServiceProvider) Boot(app infra.Glacier) {
	app.MustResolve(initPredefinedTemplates)
}
