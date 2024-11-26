package controller

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/adanos-alert/internal/template"
	"github.com/mylxsw/adanos-alert/pubsub"
	"github.com/mylxsw/adanos-alert/service"
	"github.com/mylxsw/glacier/event"
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/glacier/web"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GroupController struct {
	cc infra.Resolver
}

func NewGroupController(cc infra.Resolver) web.Controller {
	return &GroupController{cc: cc}
}

func (g GroupController) Register(router web.Router) {
	router.Group("/groups/", func(router web.Router) {
		router.Get("/", g.Groups).Name("groups:all")
		router.Get("/{id}/", g.Group).Name("groups:one")
		router.Delete("/{id}/reduce/", g.CutGroupEvents).Name("groups:reduce")
	})

	router.Group("/recoverable-groups/", func(router web.Router) {
		router.Get("/", g.RecoverableGroups).Name("recoverable-groups:all")
	})
}

type GroupsResp struct {
	Groups []GroupsGroupResp `json:"groups"`
	Users  map[string]string `json:"users"`
	Next   int64             `json:"next"`
}

type GroupsGroupResp struct {
	repository.EventGroup
	CollectTimeRemain int64 `json:"collect_time_remain"`
}

func groupFilter(ctx web.Context) bson.M {
	filter := bson.M{}
	status := ctx.Input("status")
	if status != "" {
		filter["status"] = status
	}

	// type 为空时，查询所有非 ignored 类型的事件组
	typ := ctx.Input("type")
	if typ != "" {
		filter["type"] = typ
	} else {
		filter["type"] = bson.M{"$ne": repository.EventTypeIgnored}
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

	beginAt := ctx.Input("start_at")
	endAt := ctx.Input("end_at")

	if beginAt != "" || endAt != "" {
		createdAtRange := bson.M{}
		if beginAt != "" {
			beginAtTs, err := time.ParseInLocation("2006-01-02 15:04:05", beginAt, time.Local)
			if err != nil {
				beginAtTs = time.Now()
			}

			createdAtRange["$gt"] = beginAtTs
		}

		if endAt != "" {
			endAtTs, err := time.ParseInLocation("2006-01-02 15:04:05", endAt, time.Local)
			if err != nil {
				endAtTs = time.Now()
			}

			createdAtRange["$lt"] = endAtTs
		}

		filter["created_at"] = createdAtRange
	}
	return filter
}

// Groups list all event groups
// Arguments:
//   - offset/limit
//   - status
//   - rule_id
//   - user_id
func (g GroupController) Groups(ctx web.Context, groupRepo repository.EventGroupRepo, userRepo repository.UserRepo) (*GroupsResp, error) {
	offset, limit := offsetAndLimit(ctx)

	sortDir := ctx.InputWithDefault("sort", "desc")
	grps, next, err := groupRepo.Paginate(groupFilter(ctx), offset, limit, sortDir == "desc")
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
		if grp.Status == repository.EventGroupStatusCollecting {
			timeRemain = grp.Rule.ExpectReadyAt.Unix() - time.Now().Unix()
		}
		groups[i] = GroupsGroupResp{EventGroup: grp, CollectTimeRemain: timeRemain}
	}

	return &GroupsResp{
		Groups: groups,
		Users:  userRefs,
		Next:   next,
	}, nil
}

type GroupResp struct {
	Group  repository.EventGroup `json:"group"`
	Events []repository.Event    `json:"events"`
	Next   int64                 `json:"next"`
}

func (g GroupController) Group(
	ctx web.Context,
	groupRepo repository.EventGroupRepo,
	eventRepo repository.EventRepo,
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

	filter := eventsFilter(ctx)
	filter["group_ids"] = groupID

	events, next, err := eventRepo.Paginate(filter, offset, limit)
	if err != nil {
		return nil, web.WrapJSONError(err, http.StatusInternalServerError)
	}

	for i, m := range events {
		events[i].Content = template.JSONBeauty(m.Content)
	}

	return &GroupResp{
		Group:  grp,
		Events: events,
		Next:   next,
	}, nil
}

// CutGroupEvents 缩减事件组中包含的事件，对已经完成聚合的事件组有效，
// 该操作不会影响事件组上对事件总数的计数
func (g GroupController) CutGroupEvents(webCtx web.Context, evtGrpRepo repository.EventGroupRepo, evtGroupSvc service.EventGroupService, em event.Manager) web.Response {
	groupID, err := primitive.ObjectIDFromHex(webCtx.PathVar("id"))
	if err != nil {
		return webCtx.JSONError(err.Error(), http.StatusUnprocessableEntity)
	}

	grp, err := evtGrpRepo.Get(groupID)
	if err != nil {
		return webCtx.JSONError(err.Error(), http.StatusInternalServerError)
	}

	if grp.Status == repository.EventGroupStatusCollecting || grp.Status == repository.EventGroupStatusPending {
		return webCtx.JSONError("The current event group does not currently support this operation", http.StatusUnprocessableEntity)
	}

	keepCount := webCtx.Int64Input("keep", 20)
	if keepCount < 0 || keepCount > 1000 {
		return webCtx.JSONError("keep: The number of events to retain must be between 0 - 1000", http.StatusUnprocessableEntity)
	}

	ctx, cancel := context.WithTimeout(webCtx.Context(), 10*time.Second)
	defer cancel()

	deletedCount, err := evtGroupSvc.CutGroup(ctx, groupID, keepCount)
	if err != nil {
		return webCtx.JSONError(err.Error(), http.StatusInternalServerError)
	}

	if deletedCount > 0 {
		_ = em.Publish(pubsub.EventGroupReduceEvent{
			GroupID:     grp.ID,
			KeepCount:   keepCount,
			DeleteCount: deletedCount,
			CreatedAt:   time.Now(),
		})
	}

	return webCtx.JSON(web.M{"deleted_count": deletedCount})
}

// RecoverableGroups 当前待恢复的报警组
func (g GroupController) RecoverableGroups(recoveryRepo repository.RecoveryRepo) ([]repository.Recovery, error) {
	return recoveryRepo.RecoverableEvents(context.TODO(), time.Now().AddDate(1, 0, 0))
}
