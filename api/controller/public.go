package controller

import (
	"fmt"
	"net/http"

	"github.com/mylxsw/adanos-alert/api/view"
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/glacier/web"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PublicController struct {
	cc infra.Resolver
}

func NewPublicController(cc infra.Resolver) web.Controller {
	return &PublicController{cc: cc}
}

func (p PublicController) Register(router web.Router) {
	router.Group("/groups/", func(router web.Router) {
		router.Get("/{id}.html", p.Group).Name("public:group")
	})
	router.Group("/reports/", func(router web.Router) {
		router.Get("/{id}.html", p.Reports).Name("public:report")
	})
}

// Reports 报表展示
func (p PublicController) Reports(ctx web.Context, groupRepo repository.EventGroupRepo, msgRepo repository.EventRepo, templateRepo repository.TemplateRepo) web.Response {
	id, err := primitive.ObjectIDFromHex(ctx.PathVar("id"))
	if err != nil {
		return ctx.Error("invalid request", http.StatusUnprocessableEntity)
	}

	group, err := groupRepo.Get(id)
	if err != nil {
		if err == repository.ErrNotFound {
			return ctx.Error("not found", http.StatusNotFound)
		}

		return ctx.Error(err.Error(), http.StatusInternalServerError)
	}

	var maxLimit int64 = 1000

	filter := eventsFilter(ctx)
	filter["group_ids"] = group.ID
	events, _, err := msgRepo.Paginate(filter, 0, maxLimit)
	if err != nil {
		return ctx.Error(err.Error(), http.StatusInternalServerError)
	}

	eventsCount, _ := msgRepo.Count(filter)

	var templateContent string
	if group.Rule.ReportTemplateID != primitive.NilObjectID {
		temp, err := templateRepo.Get(group.Rule.ReportTemplateID)
		if err == nil {
			templateContent = temp.Content
		}
	}

	res, err := view.ReportView(p.cc, templateContent, view.GroupData{
		Group:       group,
		Events:      events,
		EventsCount: eventsCount,
		Next:        0,
		Offset:      0,
		Limit:       maxLimit,
		Path:        ctx.Request().Raw().URL.Path,
		HasPrev:     false,
		HasNext:     false,
		PrevOffset:  0,
	})
	if err != nil {
		return ctx.Error(fmt.Sprintf("template parse failed: %v", err), http.StatusInternalServerError)
	}

	return ctx.HTML(res)
}

// Group 分组展示
func (p PublicController) Group(ctx web.Context, groupRepo repository.EventGroupRepo, msgRepo repository.EventRepo) web.Response {
	id, err := primitive.ObjectIDFromHex(ctx.PathVar("id"))
	if err != nil {
		return ctx.Error("invalid request", http.StatusUnprocessableEntity)
	}

	group, err := groupRepo.Get(id)
	if err != nil {
		if err == repository.ErrNotFound {
			return ctx.Error("not found", http.StatusNotFound)
		}

		return ctx.Error(err.Error(), http.StatusInternalServerError)
	}

	offset, limit := offsetAndLimit(ctx)
	filter := eventsFilter(ctx)
	filter["group_ids"] = group.ID

	events, next, err := msgRepo.Paginate(filter, offset, limit)
	if err != nil {
		return ctx.Error(err.Error(), http.StatusInternalServerError)
	}

	eventsCount, _ := msgRepo.Count(filter)

	res, err := view.GroupView(p.cc, view.GroupData{
		Group:       group,
		Events:      events,
		EventsCount: eventsCount,
		Next:        next,
		Offset:      offset,
		Limit:       limit,
		Path:        ctx.Request().Raw().URL.Path,
		HasPrev:     offset-limit >= 0,
		HasNext:     next > 0,
		PrevOffset:  offset - limit,
	})
	if err != nil {
		return ctx.Error(fmt.Sprintf("template parse failed: %v", err), http.StatusInternalServerError)
	}

	return ctx.HTML(res)
}
