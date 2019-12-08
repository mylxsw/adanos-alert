package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mylxsw/adanos-alert/configs"
	_ "github.com/mylxsw/adanos-alert/docs"
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier"
	httpSwagger "github.com/swaggo/http-swagger"
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
type ServiceProvider struct{}

func (s ServiceProvider) Register(app container.Container) {}

func (s ServiceProvider) Boot(app glacier.Glacier) {
	app.MustResolve(func(conf *configs.Config) {
		app.WebAppRouter(routers(app.Container()))
		app.WebAppMuxRouter(func(router *mux.Router) {
			// Swagger doc
			router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler).Name("swagger")
			// Dashboard
			router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(FS(conf.UseLocalDashboard)))).Name("assets")
		})
	})
}
