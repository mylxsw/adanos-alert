package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/mylxsw/adanos-alert/internal/queue"
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/container"
	"github.com/mylxsw/hades"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		router.Get("/jobs/", q.Jobs).Name("queue:jobs:all")
		router.Delete("/jobs/{id}/", q.Delete).Name("queue:jobs:delete")
		router.Get("/jobs/{id}/", q.Job).Name("queue:jobs:one")
	})
}

// Control controls the queue behaviors
// Args: op=pause/container/info
func (q *QueueController) Control(ctx hades.Context, manager queue.Manager) hades.Response {
	op := ctx.Input("op")
	switch op {
	case "pause":
		manager.Pause(true)
		log.Info("queue paused")
	case "continue":
		manager.Pause(false)
		log.Info("queue continued")
	case "info":
	default:
		return ctx.JSONError("invalid op argument, not support", http.StatusUnprocessableEntity)
	}

	return ctx.JSON(hades.M{
		"paused": manager.Paused(),
		"info":   manager.Info(),
	})
}

type QueueJobsResp struct {
	Jobs []repository.QueueJob `json:"jobs"`
	Next int64                 `json:"next"`
}

func (q *QueueController) Jobs(ctx hades.Context, repo repository.QueueRepo) (*QueueJobsResp, error) {
	offset, limit := offsetAndLimit(ctx)

	filter := bson.M{}

	status := ctx.Input("status")
	if status != "" {
		filter["status"] = status
	}

	jobs, next, err := repo.Paginate(filter, offset, limit)
	if err != nil {
		return nil, hades.WrapJSONError(err, http.StatusInternalServerError)
	}

	return &QueueJobsResp{
		Jobs: jobs,
		Next: next,
	}, nil
}

func (q *QueueController) Job(ctx hades.Context, repo repository.QueueRepo) (*repository.QueueJob, error) {
	jobID, err := primitive.ObjectIDFromHex(ctx.PathVar("id"))
	if err != nil {
		return nil, hades.WrapJSONError(fmt.Errorf("invalid request: %v", err), http.StatusUnprocessableEntity)
	}

	job, err := repo.Get(jobID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, hades.WrapJSONError(errors.New("no such job"), http.StatusNotFound)
		}

		return nil, hades.WrapJSONError(err, http.StatusInternalServerError)
	}

	return &job, nil
}

func (q *QueueController) Delete(ctx hades.Context, repo repository.QueueRepo) error {
	jobID, err := primitive.ObjectIDFromHex(ctx.PathVar("id"))
	if err != nil {
		return hades.WrapJSONError(fmt.Errorf("invalid request: %v", err), http.StatusUnprocessableEntity)
	}

	return repo.Delete(bson.M{
		"_id":    jobID,
		"status": bson.M{"$ne": repository.QueueItemStatusRunning},
	})
}
