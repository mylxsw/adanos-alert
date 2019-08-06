package api

import (
	"github.com/mylxsw/adanos-alert/api/controller"
	"github.com/mylxsw/go-toolkit/web"
)

func routers(router *web.Router, mw web.RequestMiddleware) {

	router.Group("/", func(router *web.Router) {
		controller.NewWelcomeController().Register(router)

	}, mw.AccessLog(), mw.JSONExceptionHandler(), mw.CORS("*"))
}
