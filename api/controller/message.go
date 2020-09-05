package controller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/mylxsw/adanos-alert/internal/job"
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/adanos-alert/pkg/misc"
	"github.com/mylxsw/adanos-alert/pkg/template"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier/web"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageController struct {
	cc container.Container
}

func NewMessageController(cc container.Container) web.Controller {
	return &MessageController{cc: cc}
}

func (m *MessageController) Register(router *web.Router) {
	router.Group("/messages", func(router *web.Router) {
		router.Get("/", m.Messages).Name("messages:all")
		router.Get("/{id}/", m.Message).Name("messages:one")
		router.Post("/{id}/matched-rules/", m.TestMatchedRules).Name("messages:matched-rules")

		router.Post("/", m.AddCommonMessage).Name("messages:add:common")
		router.Post("/logstash/", m.AddLogstashMessage).Name("messages:add:logstash")
		router.Post("/grafana/", m.AddGrafanaMessage).Name("messages:add:grafana")
		router.Post("/prometheus/api/v1/alerts", m.AddPrometheusMessage).Name("messages:add:prometheus") // url 地址末尾不包含 "/"
		router.Post("/prometheus_alertmanager/", m.AddPrometheusAlertMessage).Name("messages:add:prometheus-alert")
		router.Post("/openfalcon/im/", m.AddOpenFalconMessage).Name("messages:add:openfalcon")
	})

	router.Group("/messages-count/", func(router *web.Router) {
		router.Get("/", m.Count).Name("messages:count")
	})
}

// messagesFilter some query conditions for messages
func messagesFilter(ctx web.Context) bson.M {
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
func (m *MessageController) Count(ctx web.Context, msgRepo repository.MessageRepo) web.Response {
	filter := messagesFilter(ctx)
	msgCount, err := msgRepo.Count(filter)
	if err != nil {
		return ctx.JSONError(err.Error(), http.StatusInternalServerError)
	}

	return ctx.JSON(web.M{
		"count": msgCount,
	})
}

// MessagesResp is a response object for Messages API
type MessagesResp struct {
	Messages []repository.Message `json:"messages"`
	Next     int64                `json:"next"`
	Search   MessageSearch        `json:"search"`
}

// MessageSearch is search conditions for messages
type MessageSearch struct {
	Tags    []string `json:"tags"`
	Meta    string   `json:"meta"`
	Status  []string `json:"status"`
	Origin  string   `json:"origin"`
	GroupID string   `json:"group_id"`
}

// Messages return all messages
func (m *MessageController) Messages(ctx web.Context, msgRepo repository.MessageRepo) (*MessagesResp, error) {
	offset, limit := offsetAndLimit(ctx)

	filter := messagesFilter(ctx)
	groupIDHex := ctx.Input("group_id")
	if groupIDHex != "" {
		groupID, err := primitive.ObjectIDFromHex(groupIDHex)
		if err != nil {
			return nil, web.WrapJSONError(fmt.Errorf("invalid group_id: %w", err), http.StatusUnprocessableEntity)
		}

		filter["group_ids"] = groupID
	}

	log.WithFields(log.Fields{"filter": filter}).Debug("messages filter")

	messages, next, err := msgRepo.Paginate(filter, offset, limit)
	if err != nil {
		return nil, web.WrapJSONError(fmt.Errorf("query failed: %v", err), http.StatusInternalServerError)
	}

	for i, m := range messages {
		messages[i].Content = template.JSONBeauty(m.Content)
	}

	return &MessagesResp{
		Messages: messages,
		Next:     next,
		Search: MessageSearch{
			Tags:    template.StringTags(ctx.Input("tags"), ","),
			Meta:    ctx.Input("meta"),
			Status:  template.StringTags(ctx.Input("status"), ","),
			Origin:  ctx.Input("origin"),
			GroupID: ctx.Input("group_id"),
		},
	}, nil
}

// Message return one message
func (m *MessageController) Message(ctx web.Context, msgRepo repository.MessageRepo) (*repository.Message, error) {
	id, err := primitive.ObjectIDFromHex(ctx.PathVar("id"))
	if err != nil {
		return nil, web.WrapJSONError(fmt.Errorf("invalid id: %w", err), http.StatusUnprocessableEntity)
	}

	message, err := msgRepo.Get(id)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, web.WrapJSONError(fmt.Errorf("no such message: %w", err), http.StatusNotFound)
		}

		return nil, err
	}

	message.Content = template.JSONBeauty(message.Content)

	return &message, nil
}

func (m *MessageController) saveMessage(messageRepo repository.MessageRepo, repoMessage misc.RepoMessage) (id primitive.ObjectID, err error) {
	return messageRepo.Add(repoMessage.ToRepo())
}

func (m *MessageController) errorWrap(ctx web.Context, id primitive.ObjectID, err error) web.Response {
	if err != nil {
		return ctx.JSONError(err.Error(), http.StatusInternalServerError)
	}

	return ctx.JSON(web.M{
		"id": id.Hex(),
	})
}

// Add common message

func (m *MessageController) AddCommonMessage(ctx web.Context, messageRepo repository.MessageRepo) web.Response {
	var commonMessage misc.CommonMessage
	if err := ctx.Unmarshal(&commonMessage); err != nil {
		return ctx.JSONError(fmt.Sprintf("invalid request: %v", err), http.StatusUnprocessableEntity)
	}

	id, err := m.saveMessage(messageRepo, commonMessage)
	return m.errorWrap(ctx, id, err)
}

// AddLogstashMessage Add logstash message
func (m *MessageController) AddLogstashMessage(ctx web.Context, messageRepo repository.MessageRepo) web.Response {
	commonMessage, err := misc.LogstashToCommonMessage(ctx.Request().Body(), ctx.InputWithDefault("content-field", "message"))
	if err != nil {
		return ctx.JSONError(err.Error(), http.StatusInternalServerError)
	}

	id, err := m.saveMessage(messageRepo, commonMessage)
	return m.errorWrap(ctx, id, err)
}

// Add grafana message
func (m *MessageController) AddGrafanaMessage(ctx web.Context, messageRepo repository.MessageRepo) web.Response {
	commonMessage, err := misc.GrafanaToCommonMessage(ctx.Request().Body())
	if err != nil {
		return ctx.JSONError(err.Error(), http.StatusInternalServerError)
	}

	id, err := m.saveMessage(messageRepo, commonMessage)
	return m.errorWrap(ctx, id, err)
}

// add prometheus alert message
func (m *MessageController) AddPrometheusMessage(ctx web.Context, messageRepo repository.MessageRepo) web.Response {
	commonMessages, err := misc.PrometheusToCommonMessages(ctx.Request().Body())
	if err != nil {
		return ctx.JSONError(err.Error(), http.StatusInternalServerError)
	}

	var lastID primitive.ObjectID
	var lastErr error
	for _, cm := range commonMessages {
		lastID, lastErr = m.saveMessage(messageRepo, *cm)
		if lastErr != nil {
			log.WithFields(log.Fields{
				"message": cm,
			}).Errorf("save prometheus message failed: %v", lastErr)
		}
	}

	return m.errorWrap(ctx, lastID, lastErr)
}

// add prometheus-alert message
func (m *MessageController) AddPrometheusAlertMessage(ctx web.Context, messageRepo repository.MessageRepo) web.Response {
	commonMessage, err := misc.PrometheusAlertToCommonMessage(ctx.Request().Body())
	if err != nil {
		return ctx.JSONError(err.Error(), http.StatusInternalServerError)
	}

	id, err := m.saveMessage(messageRepo, commonMessage)
	return m.errorWrap(ctx, id, err)
}

// add open-falcon message
func (m *MessageController) AddOpenFalconMessage(ctx web.Context, messageRepo repository.MessageRepo) web.Response {
	tos := ctx.Input("tos")
	content := ctx.Input("content")

	if content == "" {
		return ctx.JSONError("invalid request, content required", http.StatusUnprocessableEntity)
	}

	id, err := m.saveMessage(messageRepo, misc.OpenFalconToCommonMessage(tos, content))
	return m.errorWrap(ctx, id, err)
}

// TestMatchedRules 测试 message 匹配哪些规则
func (m *MessageController) TestMatchedRules(ctx web.Context, msgRepo repository.MessageRepo, ruleRepo repository.RuleRepo) ([]job.MatchedRule, error) {
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

	return job.BuildMessageMatchTest(ruleRepo)(message)
}
