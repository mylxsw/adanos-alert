package controller

import (
	"net/http"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/container"
	"github.com/mylxsw/hades"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GroupController struct {
	cc *container.Container
}

func NewGroupController(cc *container.Container) hades.Controller {
	return &GroupController{cc: cc}
}

func (g GroupController) Register(router *hades.Router) {
	router.Group("/groups/", func(router *hades.Router) {
		router.Get("/", g.Groups).Name("groups:all")
		router.Get("/{id}/", g.Group).Name("groups:one")
	})
}

type GroupsResp struct {
	Groups []repository.MessageGroup `json:"groups"`
	Next   int64                     `json:"next"`
}



// Groups list all message groups
// Arguments:
//   - offset/limit
//   - status
//   - rule_id
func (g GroupController) Groups(ctx hades.Context, groupRepo repository.MessageGroupRepo) (*GroupsResp, error) {
	offset, limit := offsetAndLimit(ctx)
	filter := bson.M{}

	status := ctx.Input("status")
	if status != "" {
		filter["status"] = status
	}

	ruleID, err := primitive.ObjectIDFromHex(ctx.Input("rule_id"))
	if err == nil {
		filter["rule._id"] = ruleID
	}

	grps, next, err := groupRepo.Paginate(filter, offset, limit)
	if err != nil {
		return nil, hades.WrapJSONError(err, http.StatusInternalServerError)
	}

	return &GroupsResp{
		Groups: grps,
		Next:   next,
	}, nil
}

type GroupResp struct {
	Group    repository.MessageGroup `json:"group"`
	Messages []repository.Message    `json:"messages"`
	Next     int64                   `json:"next"`
}

func (g GroupController) Group(
	ctx hades.Context,
	groupRepo repository.MessageGroupRepo,
	messageRepo repository.MessageRepo,
) (*GroupResp, error) {
	offset := ctx.Int64Input("offset", 0)
	limit := ctx.Int64Input("limit", 10)

	groupID, err := primitive.ObjectIDFromHex(ctx.PathVar("id"))
	if err != nil {
		return nil, hades.WrapJSONError(err, http.StatusUnprocessableEntity)
	}

	grp, err := groupRepo.Get(groupID)
	if err != nil {
		return nil, hades.WrapJSONError(err, http.StatusInternalServerError)
	}

	filter := messagesFilter(ctx)
	filter["group_ids"] = groupID

	messages, next, err := messageRepo.Paginate(filter, offset, limit)
	if err != nil {
		return nil, hades.WrapJSONError(err, http.StatusInternalServerError)
	}

	return &GroupResp{
		Group:    grp,
		Messages: messages,
		Next:     next,
	}, nil
}
