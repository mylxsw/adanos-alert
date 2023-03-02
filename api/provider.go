package api

import (
	"errors"
	"net/http"
	"runtime/debug"

	"github.com/gorilla/mux"
	"github.com/mylxsw/adanos-alert/api/controller"
	"github.com/mylxsw/adanos-alert/configs"
	_ "github.com/mylxsw/adanos-alert/docs"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/glacier/listener"
	"github.com/mylxsw/glacier/web"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// @title Adanos-alert API
// @version 1.0
// @description Adanos-alert is a alert manager with multi alert channel support

// @contact.name mylxsw
// @contact.url https://github.com/mylxsw/adanos-alert
// @contact.email mylxsw@aicode.cc

// @license.name MIT
// @license.url https://raw.githubusercontent.com/mylxsw/adanos-alert/master/LICENSE

// @host localhost:19999
// @BasePath /api
type Provider struct{}

func (s Provider) Aggregates() []infra.Provider {
	return []infra.Provider{
		web.Provider(
			listener.FlagContext("listen"),
			web.SetRouteHandlerOption(s.routes),
			web.SetMuxRouteHandlerOption(s.muxRoutes),
			web.SetExceptionHandlerOption(func(ctx web.Context, err interface{}) web.Response {
				log.Errorf("error: %v, call stack: %s", err, debug.Stack())
				return nil
			}),
		),
	}
}

func (s Provider) muxRoutes(cc infra.Resolver, router *mux.Router) {
	cc.MustResolve(func(conf *configs.Config) {
		// prometheus metrics
		router.PathPrefix("/metrics").Handler(promhttp.Handler())
		// health check
		router.PathPrefix("/health").Handler(HealthCheck{})
		// Dashboard
		router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(FS(conf.UseLocalDashboard)))).Name("assets")
	})
}

func (s Provider) routes(cc infra.Resolver, router web.Router, mw web.RequestMiddleware) {
	conf := cc.MustGet(&configs.Config{}).(*configs.Config)
	mws := make([]web.HandlerDecorator, 0)
	mws = append(mws, mw.AccessLog(log.Module("api")), mw.CORS("*"))
	if conf.APIToken != "" {
		authMiddleware := mw.AuthHandler(func(ctx web.Context, typ string, credential string) error {
			if typ != "Bearer" {
				return errors.New("invalid auth type, only support Bearer")
			}

			if credential != conf.APIToken {
				return errors.New("token not match")
			}

			return nil
		})

		mws = append(mws, authMiddleware)
	}

	router.WithMiddleware(mws...).Controllers(
		"/api",
		controller.NewWelcomeController(cc),
		controller.NewEventController(cc),
		controller.NewQueueController(cc),
		controller.NewUserController(cc),
		controller.NewGroupController(cc),
		controller.NewRuleController(cc),
		controller.NewTemplateController(cc),
		controller.NewDingdingRobotController(cc),
		controller.NewAgentController(cc),
		controller.NewStatisticsController(cc),
		controller.NewSyslogController(cc),
		controller.NewJiraController(cc),
	)

	router.WithMiddleware(mw.AccessLog(log.Module("api")), mw.CORS("*")).Controllers(
		"/ui",
		controller.NewPublicController(cc),
	)
}

func (s Provider) Register(app infra.Binder) {}
func (s Provider) Boot(app infra.Resolver)   {}
