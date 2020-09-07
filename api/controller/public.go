package controller

import (
	"fmt"
	"net/http"

	"github.com/mylxsw/adanos-alert/api/view"
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier/web"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PublicController struct {
	cc container.Container
}

func NewPublicController(cc container.Container) web.Controller {
	return &PublicController{cc: cc}
}

func (p PublicController) Register(router *web.Router) {
	router.Group("/groups/", func(router *web.Router) {
		router.Get("/{id}.html", p.Group).Name("public:group")
	})
}

func (p PublicController) Group(ctx web.Context, groupRepo repository.MessageGroupRepo, msgRepo repository.MessageRepo) web.Response {
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
	filter := messagesFilter(ctx)
	filter["group_ids"] = group.ID

	messages, next, err := msgRepo.Paginate(filter, offset, limit)
	if err != nil {
		return ctx.Error(err.Error(), http.StatusInternalServerError)
	}

	messageCount, _ := msgRepo.Count(filter)

	res, err := view.GroupView(p.cc, view.GroupData{
		Group:        group,
		Messages:     messages,
		MessageCount: messageCount,
		Next:         next,
		Offset:       offset,
		Limit:        limit,
		Path:         ctx.Request().Raw().URL.Path,
		HasPrev:      offset-limit >= 0,
		HasNext:      next > 0,
		PrevOffset:   offset - limit,
	})
	if err != nil {
		return ctx.Error(fmt.Sprintf("template parse failed: %v", err), http.StatusInternalServerError)
	}

	return ctx.HTML(res)
}
