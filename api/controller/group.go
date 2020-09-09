package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/adanos-alert/pkg/template"
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier/web"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GroupController struct {
	cc container.Container
}

func NewGroupController(cc container.Container) web.Controller {
	return &GroupController{cc: cc}
}

func (g GroupController) Register(router *web.Router) {
	router.Group("/groups/", func(router *web.Router) {
		router.Get("/", g.Groups).Name("groups:all")
		router.Get("/{id}/", g.Group).Name("groups:one")
	})
}

type GroupsResp struct {
	Groups []GroupsGroupResp `json:"groups"`
	Users  map[string]string `json:"users"`
	Next   int64             `json:"next"`
}

type GroupsGroupResp struct {
	repository.MessageGroup
	CollectTimeRemain int64 `json:"collect_time_remain"`
}

// Groups list all message groups
// Arguments:
//   - offset/limit
//   - status
//   - rule_id
//   - user_id
func (g GroupController) Groups(ctx web.Context, groupRepo repository.MessageGroupRepo, userRepo repository.UserRepo) (*GroupsResp, error) {
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

	userID, err := primitive.ObjectIDFromHex(ctx.Input("user_id"))
	if err == nil {
		filter["actions.user_refs"] = userID
	}

	dingID := ctx.Input("dingding_id")
	if dingID != "" {
		filter["actions.meta"] = bson.M{"$regex": fmt.Sprintf(`"robot_id":"%s"`, dingID)}
	}

	grps, next, err := groupRepo.Paginate(filter, offset, limit)
	if err != nil {
		return nil, web.WrapJSONError(err, http.StatusInternalServerError)
	}

	userIDs := make([]primitive.ObjectID, 0)
	for _, grp := range grps {
		for _, act := range grp.Actions {
			userIDs = append(userIDs, act.UserRefs...)
		}
	}

	users, _ := userRepo.Find(bson.M{"_id": bson.M{"$in": userIDs}})
	userRefs := make(map[string]string)
	for _, u := range users {
		userRefs[u.ID.Hex()] = u.Name
	}

	groups := make([]GroupsGroupResp, len(grps))
	for i, grp := range grps {
		var timeRemain int64 = 0
		if grp.Status == repository.MessageGroupStatusCollecting {
			timeRemain = grp.Rule.ExpectReadyAt.Unix() - time.Now().Unix()
		}
		groups[i] = GroupsGroupResp{MessageGroup: grp, CollectTimeRemain: timeRemain}
	}

	return &GroupsResp{
		Groups: groups,
		Users:  userRefs,
		Next:   next,
	}, nil
}

type GroupResp struct {
	Group    repository.MessageGroup `json:"group"`
	Messages []repository.Message    `json:"messages"`
	Next     int64                   `json:"next"`
}

func (g GroupController) Group(
	ctx web.Context,
	groupRepo repository.MessageGroupRepo,
	messageRepo repository.MessageRepo,
) (*GroupResp, error) {
	offset := ctx.Int64Input("offset", 0)
	limit := ctx.Int64Input("limit", 10)

	groupID, err := primitive.ObjectIDFromHex(ctx.PathVar("id"))
	if err != nil {
		return nil, web.WrapJSONError(err, http.StatusUnprocessableEntity)
	}

	grp, err := groupRepo.Get(groupID)
	if err != nil {
		return nil, web.WrapJSONError(err, http.StatusInternalServerError)
	}

	filter := messagesFilter(ctx)
	filter["group_ids"] = groupID

	messages, next, err := messageRepo.Paginate(filter, offset, limit)
	if err != nil {
		return nil, web.WrapJSONError(err, http.StatusInternalServerError)
	}

	for i, m := range messages {
		messages[i].Content = template.JSONBeauty(m.Content)
	}

	return &GroupResp{
		Group:    grp,
		Messages: messages,
		Next:     next,
	}, nil
}
