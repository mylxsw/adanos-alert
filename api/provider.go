package api

import (
	"github.com/mylxsw/glacier"
	"github.com/mylxsw/go-toolkit/container"
	"github.com/mylxsw/go-toolkit/web"
)

type ServiceProvider struct{}

func (s ServiceProvider) Register(app *container.Container) {}

func (s ServiceProvider) Boot(app *glacier.Glacier) {
	app.WebAppRouter(func(router *web.Router, mw web.RequestMiddleware) {

	})
}
