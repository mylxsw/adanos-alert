package api

import (
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier"
	"github.com/mylxsw/glacier/web"
)

type ServiceProvider struct{}

func (s ServiceProvider) Register(app container.Container) {}

func (s ServiceProvider) Boot(app glacier.Glacier) {
	app.WebAppRouter(routers(app.Container()))
}

func routers(cc container.Container) func(router *web.Router, mw web.RequestMiddleware) {
	return func(router *web.Router, mw web.RequestMiddleware) {
		mws := make([]web.HandlerDecorator, 0)
		mws = append(mws, mw.AccessLog(log.Module("api")), mw.CORS("*"))
		router.WithMiddleware(mws...).Controllers(
			"/api",
			NewMessageController(cc),
		)
	}
}
