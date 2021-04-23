package api

import (
	"runtime/debug"

	"github.com/gorilla/mux"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/glacier/listener"
	"github.com/mylxsw/glacier/web"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Provider struct{}

func (s Provider) Aggregates() []infra.Provider {
	return []infra.Provider{
		web.Provider(
			listener.FlagContext("listen"),
			web.SetRouteHandlerOption(func(cc infra.Resolver, router web.Router, mw web.RequestMiddleware) {
				mws := make([]web.HandlerDecorator, 0)
				mws = append(mws, mw.AccessLog(log.Module("api")), mw.CORS("*"))
				router.WithMiddleware(mws...).Controllers(
					"/api",
					NewEventController(cc),
				)
			}),
			web.SetMuxRouteHandlerOption(func(cc infra.Resolver, router *mux.Router) {
				// prometheus metrics
				router.PathPrefix("/metrics").Handler(promhttp.Handler())
			}),
			web.SetExceptionHandlerOption(func(ctx web.Context, err interface{}) web.Response {
				log.Errorf("error: %v, call stack: %s", err, debug.Stack())
				return nil
			}),
		),
	}
}

func (s Provider) Register(app infra.Binder) {}
func (s Provider) Boot(app infra.Resolver)   {}
