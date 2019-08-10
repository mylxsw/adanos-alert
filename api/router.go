package api

import (
	"github.com/mylxsw/adanos-alert/api/controller"
	"github.com/mylxsw/hades"
)

func routers(router *hades.Router, mw hades.RequestMiddleware) {

	router.Group("/", func(router *hades.Router) {
		controller.NewWelcomeController().Register(router)

	}, mw.AccessLog(), mw.JSONExceptionHandler(), mw.CORS("*"))
}
