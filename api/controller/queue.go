package controller

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/mylxsw/adanos-alert/internal/queue"
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/glacier/web"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QueueController struct {
	cc infra.Resolver
}

func NewQueueController(cc infra.Resolver) web.Controller {
	return &QueueController{cc: cc}
}

func (q *QueueController) Register(router web.Router) {
	router.Group("/queue/", func(router web.Router) {
		router.Post("/control/", q.Control).Name("queue:control")
		router.Get("/jobs/", q.Jobs).Name("queue:jobs:all")
		router.Delete("/jobs/{id}/", q.Delete).Name("queue:jobs:delete")
		router.Get("/jobs/{id}/", q.Job).Name("queue:jobs:one")
	})
}

// Control controls the queue behaviors
// Args: op=pause/container/info
func (q *QueueController) Control(ctx web.Context, manager queue.Manager) web.Response {
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

	return ctx.JSON(web.M{
		"paused": manager.Paused(),
		"info":   manager.Info(),
	})
}

type QueueJobsResp struct {
	Jobs []QueueJob `json:"jobs"`
	Next int64      `json:"next"`
}

type QueueJob struct {
	repository.QueueJob
	ExecuteTimeRemain int64 `json:"execute_time_remain"`
}

func (q *QueueController) Jobs(ctx web.Context, repo repository.QueueRepo) (*QueueJobsResp, error) {
	offset, limit := offsetAndLimit(ctx)

	filter := bson.M{}

	status := ctx.Input("status")
	if status != "" {
		filter["status"] = status
	}

	jobs, next, err := repo.Paginate(filter, offset, limit)
	if err != nil {
		return nil, web.WrapJSONError(err, http.StatusInternalServerError)
	}

	queueJobs := make([]QueueJob, len(jobs))
	for i, job := range jobs {
		queueJobs[i] = QueueJob{QueueJob: job, ExecuteTimeRemain: job.NextExecuteAt.Unix() - time.Now().Unix()}
	}

	return &QueueJobsResp{
		Jobs: queueJobs,
		Next: next,
	}, nil
}

func (q *QueueController) Job(ctx web.Context, repo repository.QueueRepo) (*QueueJob, error) {
	jobID, err := primitive.ObjectIDFromHex(ctx.PathVar("id"))
	if err != nil {
		return nil, web.WrapJSONError(fmt.Errorf("invalid request: %v", err), http.StatusUnprocessableEntity)
	}

	job, err := repo.Get(jobID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, web.WrapJSONError(errors.New("no such job"), http.StatusNotFound)
		}

		return nil, web.WrapJSONError(err, http.StatusInternalServerError)
	}

	return &QueueJob{QueueJob: job, ExecuteTimeRemain: job.NextExecuteAt.Unix() - time.Now().Unix()}, nil
}

func (q *QueueController) Delete(ctx web.Context, repo repository.QueueRepo) error {
	jobID, err := primitive.ObjectIDFromHex(ctx.PathVar("id"))
	if err != nil {
		return web.WrapJSONError(fmt.Errorf("invalid request: %v", err), http.StatusUnprocessableEntity)
	}

	return repo.Delete(bson.M{
		"_id":    jobID,
		"status": bson.M{"$ne": repository.QueueItemStatusRunning},
	})
}
