package api

import (
	"github.com/mylxsw/adanos-alert/api/controller"
	"github.com/mylxsw/go-toolkit/web"
)

func routers(router *web.Router, mw web.RequestMiddleware) {
	welcome := controller.NewWelcomeController()

	router.Group("/", func(router *web.Router) {
		router.Get("/", welcome.Home)


	}, mw.AccessLog(), mw.JSONExceptionHandler(), mw.CORS("*"))
}
