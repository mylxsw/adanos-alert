package api

import (
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier"
)

type ServiceProvider struct{}

func (s ServiceProvider) Register(app *container.Container) {}

func (s ServiceProvider) Boot(app *glacier.Glacier) {
	app.WebAppRouter(routers)
	app.WebAppMuxRouter(graphqlRouters(app.Container()))
}
