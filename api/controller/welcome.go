package controller

import (
	"github.com/mylxsw/go-toolkit/web"
)

type WelcomeController struct{}

func NewWelcomeController() *WelcomeController {
	return &WelcomeController{}
}

func (controller *WelcomeController) Register(router *web.Router) {
	router.Get("/", controller.Home)
}

func (*WelcomeController) Home(ctx *web.WebContext, req *web.Request) web.HTTPResponse {
	return ctx.API("0000", "hello, world", nil)
}
