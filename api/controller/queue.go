package controller

import (
	"net/http"

	"github.com/mylxsw/adanos-alert/internal/queue"
	"github.com/mylxsw/container"
	"github.com/mylxsw/hades"
)

type QueueController struct {
	cc *container.Container
}

func NewQueueController(cc *container.Container) hades.Controller {
	return &QueueController{cc: cc}
}

func (q *QueueController) Register(router *hades.Router) {
	router.Group("/queue/", func(router *hades.Router) {
		router.Post("/control/", q.Control).Name("queue:control")
	})
}

// Control controls the queue behaviors
// Args: op=pause/container/info
func (q *QueueController) Control(ctx hades.Context, manager queue.Manager) hades.Response {
	op := ctx.Input("op")
	switch op {
	case "pause":
		manager.Pause(true)
	case "continue":
		manager.Pause(false)
	case "info":
	default:
		return ctx.JSONError("invalid op argument, not support", http.StatusUnprocessableEntity)
	}

	return ctx.JSON(hades.M{
		"paused": manager.Paused(),
		"info":   manager.Info(),
	})
}
