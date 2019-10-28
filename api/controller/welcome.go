package controller

import (
	"github.com/mylxsw/container"
	"github.com/mylxsw/hades"
)

type WelcomeController struct {
	cc *container.Container
}

func NewWelcomeController(cc *container.Container) hades.Controller {
	return &WelcomeController{cc: cc}
}

func (w *WelcomeController) Register(router *hades.Router) {
	router.Any("/", w.Home).Name("welcome:home")
}

type WelcomeMessage struct {
	Version string `json:"version"`
}

// Home 欢迎页面，API版本信息
// @Summary 欢迎页面，API版本信息
// @Success 200 {object} controller.WelcomeMessage
// @Router / [get]
func (w *WelcomeController) Home(ctx hades.Context, req hades.Request) WelcomeMessage {
	return WelcomeMessage{Version: w.cc.MustGet("version").(string)}
}
