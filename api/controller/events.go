package controller

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/mylxsw/adanos-alert/internal/extension"
	"github.com/mylxsw/adanos-alert/internal/job"
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/adanos-alert/internal/template"
	"github.com/mylxsw/adanos-alert/pkg/misc"
	"github.com/mylxsw/adanos-alert/service"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier/web"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventController struct {
	cc container.Container
}

func NewEventController(cc container.Container) web.Controller {
	return &EventController{cc: cc}
}

func (m *EventController) Register(router *web.Router) {
	router.Group("/messages", func(router *web.Router) {
		router.Post("/", m.AddCommonEvent).Name("events:add:common")
		router.Post("/logstash/", m.AddLogstashEvent).Name("events:add:logstash")
		router.Post("/grafana/", m.AddGrafanaEvent).Name("events:add:grafana")
		router.Post("/prometheus/api/v1/alerts", m.AddPrometheusEvent).Name("events:add:prometheus") // url 地址末尾不包含 "/"
		router.Post("/prometheus_alertmanager/", m.AddPrometheusAlertEvent).Name("events:add:prometheus-alert")
		router.Post("/openfalcon/im/", m.AddOpenFalconEvent).Name("events:add:openfalcon")
	})

	router.Group("/events", func(router *web.Router) {
		router.Get("/", m.Events).Name("events:all")
		router.Get("/{id}/", m.Event).Name("events:one")
		router.Delete("/{id}/", m.DeleteEvent).Name("events:delete")
		
		router.Post("/{id}/matched-rules/", m.TestMatchedRules).Name("events:matched-rules")
		router.Post("/{id}/reproduce/", m.ReproduceEvent).Name("events:reproduce-event")

		router.Post("/", m.AddCommonEvent).Name("events:add:common")
		router.Post("/logstash/", m.AddLogstashEvent).Name("events:add:logstash")
		router.Post("/grafana/", m.AddGrafanaEvent).Name("events:add:grafana")
		router.Post("/prometheus/api/v1/alerts", m.AddPrometheusEvent).Name("events:add:prometheus") // url 地址末尾不包含 "/"
		router.Post("/prometheus_alertmanager/", m.AddPrometheusAlertEvent).Name("events:add:prometheus-alert")
		router.Post("/openfalcon/im/", m.AddOpenFalconEvent).Name("events:add:openfalcon")
	})

	router.Group("/event-relations", func(router *web.Router) {
		router.Get("/{id}/", m.QueryEventRelation).Name("event-relations:one")
		router.Get("/{id}/notes/", m.QueryEventRelationNotes).Name("event-relations:notes")
		router.Post("/{id}/notes/", m.AddEventRelationNote).Name("event-relations:notes:add")
	})

	router.Group("/events-count/", func(router *web.Router) {
		router.Get("/", m.Count).Name("events:count")
	})
}

// eventsFilter some query conditions for messages
func eventsFilter(ctx web.Context) bson.M {
	filter := bson.M{}

	meta := ctx.Input("meta")
	if meta != "" {
		kv := strings.SplitN(meta, ":", 2)
		if len(kv) == 1 {
			filter["meta."+kv[0]] = bson.M{"$exists": true}
		} else {
			filter["meta."+kv[0]] = strings.TrimSpace(kv[1])
		}
	}

	tags := template.StringTags(ctx.Input("tags"), ",")
	if len(tags) > 0 {
		filter["tags"] = bson.M{"$in": tags}
	}

	origin := ctx.Input("origin")
	if origin != "" {
		filter["origin"] = bson.M{"$regex": origin}
	}
	status := template.StringTags(ctx.Input("status"), ",")
	if len(status) > 0 {
		filter["status"] = bson.M{"$in": status}
	}

	return filter
}

// Count return message count for your conditions
func (m *EventController) Count(ctx web.Context, msgRepo repository.EventRepo) web.Response {
	filter := eventsFilter(ctx)
	eventCount, err := msgRepo.Count(filter)
	if err != nil {
		return ctx.JSONError(err.Error(), http.StatusInternalServerError)
	}

	return ctx.JSON(web.M{
		"count": eventCount,
	})
}

// EventsResp is a response object for Events API
type EventsResp struct {
	Events []repository.Event `json:"events"`
	Next   int64              `json:"next"`
	Search EventSearch        `json:"search"`
}

// EventSearch is search conditions for messages
type EventSearch struct {
	Tags    []string `json:"tags"`
	Meta    string   `json:"meta"`
	Status  []string `json:"status"`
	Origin  string   `json:"origin"`
	GroupID string   `json:"group_id"`
}

// Events return all messages
func (m *EventController) Events(ctx web.Context, evtRepo repository.EventRepo) (*EventsResp, error) {
	offset, limit := offsetAndLimit(ctx)

	filter := eventsFilter(ctx)
	groupIDHex := ctx.Input("group_id")
	if groupIDHex != "" {
		groupID, err := primitive.ObjectIDFromHex(groupIDHex)
		if err != nil {
			return nil, web.WrapJSONError(fmt.Errorf("invalid group_id: %w", err), http.StatusUnprocessableEntity)
		}

		filter["group_ids"] = groupID
	}

	relationIDHex := ctx.Input("relation_id")
	if relationIDHex != "" {
		relationID, err := primitive.ObjectIDFromHex(relationIDHex)
		if err != nil {
			return nil, web.WrapJSONError(fmt.Errorf("invalid relation_id: %w", err), http.StatusUnprocessableEntity)
		}
		filter["relation_ids"] = relationID
	}

	evtIDHex := ctx.Input("event_id")
	if evtIDHex != "" {
		evtID, err := primitive.ObjectIDFromHex(evtIDHex)
		if err != nil {
			return nil, web.WrapJSONError(fmt.Errorf("invalid event_id: %w", err), http.StatusUnprocessableEntity)
		}
		filter["_id"] = evtID
	}

	if log.DebugEnabled() {
		log.WithFields(log.Fields{"filter": filter}).Debug("events filter")
	}

	events, next, err := evtRepo.Paginate(filter, offset, limit)
	if err != nil {
		return nil, web.WrapJSONError(fmt.Errorf("query failed: %v", err), http.StatusInternalServerError)
	}

	for i, m := range events {
		events[i].Content = template.JSONBeauty(m.Content)
	}

	return &EventsResp{
		Events: events,
		Next:   next,
		Search: EventSearch{
			Tags:    template.StringTags(ctx.Input("tags"), ","),
			Meta:    ctx.Input("meta"),
			Status:  template.StringTags(ctx.Input("status"), ","),
			Origin:  ctx.Input("origin"),
			GroupID: ctx.Input("group_id"),
		},
	}, nil
}

// Event return one message
func (m *EventController) Event(ctx web.Context, eventRepo repository.EventRepo) (*repository.Event, error) {
	id, err := primitive.ObjectIDFromHex(ctx.PathVar("id"))
	if err != nil {
		return nil, web.WrapJSONError(fmt.Errorf("invalid id: %w", err), http.StatusUnprocessableEntity)
	}

	event, err := eventRepo.Get(id)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, web.WrapJSONError(fmt.Errorf("no such event: %w", err), http.StatusNotFound)
		}

		return nil, err
	}

	event.Content = template.JSONBeauty(event.Content)

	return &event, nil
}

func (m *EventController) ReproduceEvent(ctx web.Context, eventRepo repository.EventRepo, eventService service.EventService) web.Response {
	event, err := m.Event(ctx, eventRepo)
	if err != nil {
		return ctx.JSONError(err.Error(), http.StatusInternalServerError)
	}

	id, err := eventService.Add(context.TODO(), extension.CommonEvent{
		Content: event.Content,
		Meta:    event.Meta,
		Tags:    append(event.Tags, "adanos-reproduced"),
		Origin:  event.Origin,
	})
	if err != nil {
		return ctx.JSONError(err.Error(), http.StatusInternalServerError)
	}

	return ctx.JSON(web.M{
		"id": id.Hex(),
	})
}

func (m *EventController) errorWrap(ctx web.Context, id primitive.ObjectID, err error) web.Response {
	if err != nil {
		return ctx.JSONError(err.Error(), http.StatusInternalServerError)
	}

	return ctx.JSON(web.M{
		"id": misc.IfElse(id != primitive.NilObjectID, id.Hex(), ""),
	})
}

// Add common message

func (m *EventController) AddCommonEvent(ctx web.Context, eventService service.EventService) web.Response {
	var commonMessage extension.CommonEvent
	if err := ctx.Unmarshal(&commonMessage); err != nil {
		return ctx.JSONError(fmt.Sprintf("invalid request: %v", err), http.StatusUnprocessableEntity)
	}

	id, err := eventService.Add(ctx.Context(), commonMessage)
	return m.errorWrap(ctx, id, err)
}

// AddLogstashEvent Add logstash message
func (m *EventController) AddLogstashEvent(ctx web.Context, eventService service.EventService) web.Response {
	commonMessage, err := extension.LogstashToCommonEvent(ctx.Request().Body(), ctx.InputWithDefault("content-field", "message"))
	if err != nil {
		return ctx.JSONError(err.Error(), http.StatusInternalServerError)
	}

	id, err := eventService.Add(ctx.Context(), *commonMessage)
	return m.errorWrap(ctx, id, err)
}

// AddGrafanaEvent Add grafana message
func (m *EventController) AddGrafanaEvent(ctx web.Context, eventService service.EventService) web.Response {
	commonMessage, err := extension.GrafanaToCommonEvent(ctx.Request().Body())
	if err != nil {
		return ctx.JSONError(err.Error(), http.StatusInternalServerError)
	}

	id, err := eventService.Add(ctx.Context(), *commonMessage)
	return m.errorWrap(ctx, id, err)
}

// AddPrometheusEvent add prometheus alert message
func (m *EventController) AddPrometheusEvent(ctx web.Context, eventService service.EventService) web.Response {
	commonMessages, err := extension.PrometheusToCommonEvents(ctx.Request().Body())
	if err != nil {
		return ctx.JSONError(err.Error(), http.StatusInternalServerError)
	}

	var lastID primitive.ObjectID
	var lastErr error
	for _, cm := range commonMessages {
		lastID, lastErr = eventService.Add(ctx.Context(), *cm)
		if lastErr != nil {
			log.WithFields(log.Fields{
				"message": cm,
			}).Errorf("save prometheus message failed: %v", lastErr)
		}
	}

	return m.errorWrap(ctx, lastID, lastErr)
}

// AddPrometheusAlertEvent add prometheus-alert message
func (m *EventController) AddPrometheusAlertEvent(ctx web.Context, eventService service.EventService) web.Response {
	commonMessage, err := extension.PrometheusAlertToCommonEvent(ctx.Request().Body())
	if err != nil {
		return ctx.JSONError(err.Error(), http.StatusInternalServerError)
	}

	id, err := eventService.Add(ctx.Context(), *commonMessage)
	return m.errorWrap(ctx, id, err)
}

// add open-falcon message
func (m *EventController) AddOpenFalconEvent(ctx web.Context, eventService service.EventService) web.Response {
	tos := ctx.Input("tos")
	content := ctx.Input("content")

	if content == "" {
		return ctx.JSONError("invalid request, content required", http.StatusUnprocessableEntity)
	}

	id, err := eventService.Add(ctx.Context(), *extension.OpenFalconToCommonEvent(tos, content))
	return m.errorWrap(ctx, id, err)
}

// TestMatchedRules 测试 message 匹配哪些规则
func (m *EventController) TestMatchedRules(ctx web.Context, msgRepo repository.EventRepo, ruleRepo repository.RuleRepo) ([]job.MatchedRule, error) {
	msgID, err := primitive.ObjectIDFromHex(ctx.PathVar("id"))
	if err != nil {
		return nil, errors.Wrap(err, "invalid message id")
	}

	message, err := msgRepo.Get(msgID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, errors.Wrap(err, "no such message")
		}

		return nil, errors.Wrap(err, "query message failed")
	}

	return job.BuildEventMatchTest(ruleRepo)(message)
}

// QueryEventRelation 查询事件关联
func (m *EventController) QueryEventRelation(ctx web.Context, evtRelationRepo repository.EventRelationRepo) (*repository.EventRelation, error) {
	relID, err := primitive.ObjectIDFromHex(ctx.PathVar("id"))
	if err != nil {
		return nil, errors.Wrap(err, "invalid relation id")
	}

	rel, err := evtRelationRepo.Get(context.TODO(), relID)
	if err != nil {
		return nil, errors.Wrap(err, "query relation failed")
	}

	return &rel, nil
}

// EventRelationNotesResp 事件关联备注响应对象
type EventRelationNotesResp struct {
	Notes []repository.EventRelationNote `json:"notes"`
	Next  int64                          `json:"next"`
}

// QueryEventRelationNotes 查询事件关联的备注
func (m *EventController) QueryEventRelationNotes(ctx web.Context, evtRelationNoteRepo repository.EventRelationNoteRepo) (*EventRelationNotesResp, error) {
	relID, err := primitive.ObjectIDFromHex(ctx.PathVar("id"))
	if err != nil {
		return nil, errors.Wrap(err, "invalid relation id")
	}

	offset, limit := offsetAndLimit(ctx)
	notes, next, err := evtRelationNoteRepo.PaginateNotes(context.TODO(), relID, bson.M{}, offset, limit)
	if err != nil {
		return nil, err
	}

	return &EventRelationNotesResp{
		Notes: notes,
		Next:  next,
	}, nil
}

// AddEventRelationNote 添加事件关联备注
func (m *EventController) AddEventRelationNote(ctx web.Context, evtRelationNoteRepo repository.EventRelationNoteRepo) web.Response {
	relID, err := primitive.ObjectIDFromHex(ctx.PathVar("id"))
	if err != nil {
		return ctx.JSONError(err.Error(), http.StatusUnprocessableEntity)
	}

	note := ctx.Input("note")
	eventID, _ := primitive.ObjectIDFromHex(ctx.Input("event_id"))

	id, err := evtRelationNoteRepo.AddNote(context.TODO(), repository.EventRelationNote{
		RelationID: relID,
		EventID:    eventID,
		Note:       note,
		// TODO 创建人功能暂时不可用，需要待用户权限体系建立后使用
		CreatorID:   primitive.NilObjectID,
		CreatorName: "Default",
	})
	if err != nil {
		return ctx.JSONError(err.Error(), http.StatusUnprocessableEntity)
	}

	return ctx.JSON(web.M{
		"id": id,
	})
}

// DeleteEvent 删除事件
func (m *EventController) DeleteEvent(ctx web.Context, evtRepo repository.EventRepo) web.Response {
	eventID, err := primitive.ObjectIDFromHex(ctx.PathVar("id"))
	if err != nil {
		return ctx.JSONError("invalid event id", http.StatusUnprocessableEntity)
	}

	if err := evtRepo.DeleteID(eventID); err != nil {
		return ctx.JSONError(err.Error(), http.StatusInternalServerError)
	}

	return ctx.JSON(web.M{})
}
