package api

import (
	"github.com/gorilla/mux"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/glacier/web"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type ServiceProvider struct{}

func (s ServiceProvider) Register(app container.Container) {}

func (s ServiceProvider) Boot(app infra.Glacier) {
	app.WebAppRouter(routers(app.Container()))
	app.WebAppMuxRouter(func(router *mux.Router) {
		// prometheus metrics
		router.PathPrefix("/metrics").Handler(promhttp.Handler())
	})
}

func routers(cc container.Container) func(router *web.Router, mw web.RequestMiddleware) {
	return func(router *web.Router, mw web.RequestMiddleware) {
		mws := make([]web.HandlerDecorator, 0)
		mws = append(mws, mw.AccessLog(log.Module("api")), mw.CORS("*"))
		router.WithMiddleware(mws...).Controllers(
			"/api",
			NewEventController(cc),
		)
	}
}
