package controller

import (
	"github.com/mylxsw/hades"
)

type WelcomeController struct{}

func NewWelcomeController() *WelcomeController {
	return &WelcomeController{}
}

func (controller *WelcomeController) Register(router *hades.Router) {
	router.Get("/", controller.Home)
}

func (*WelcomeController) Home(ctx *hades.WebContext, req *hades.Request) hades.HTTPResponse {
	return ctx.API("0000", "hello, world", nil)
}
