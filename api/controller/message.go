package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jeremywohl/flatten"
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/adanos-alert/pkg/template"
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier/web"
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
		filter["meta.value"] = meta
	}
	tags := ctx.Input("tags")
	if tags != "" {
		tt := strings.Split(tags, ",")
		if len(tt) == 1 {
			filter["tag"] = tags
		} else {
			filter["tag"] = bson.M{"$in": tt}
		}
	}
	origin := ctx.Input("origin")
	if origin != "" {
		filter["origin"] = origin
	}
	status := ctx.Input("status")
	if status != "" {
		ss := strings.Split(status, ",")
		if len(ss) == 1 {
			filter["status"] = status
		} else {
			filter["status"] = bson.M{"$in": ss}
		}
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

	messages, next, err := msgRepo.Paginate(filter, offset, limit)
	if err != nil {
		return nil, web.WrapJSONError(fmt.Errorf("query failed: %v", err), http.StatusInternalServerError)
	}

	return &MessagesResp{
		Messages: messages,
		Next:     next,
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

	return &message, nil
}

type RepoMessage interface {
	ToRepo() repository.Message
}

func (m *MessageController) saveMessage(messageRepo repository.MessageRepo, repoMessage RepoMessage, ctx web.Context) web.Response {
	id, err := messageRepo.Add(repoMessage.ToRepo())
	if err != nil {
		return ctx.JSONError(err.Error(), http.StatusInternalServerError)
	}

	return ctx.JSON(web.M{
		"id": id.Hex(),
	})
}

// Add common message

type CommonMessage struct {
	Content string                 `json:"content"`
	Meta    repository.MessageMeta `json:"meta"`
	Tags    []string               `json:"tags"`
	Origin  string                 `json:"origin"`
}

func (msg CommonMessage) ToRepo() repository.Message {
	return repository.Message{
		Content: msg.Content,
		Meta:    msg.Meta,
		Tags:    msg.Tags,
		Origin:  msg.Origin,
	}
}

func (m *MessageController) AddCommonMessage(ctx web.Context, messageRepo repository.MessageRepo) web.Response {
	var commonMessage CommonMessage
	if err := ctx.Unmarshal(&commonMessage); err != nil {
		return ctx.JSONError(fmt.Sprintf("invalid request: %v", err), http.StatusUnprocessableEntity)
	}

	return m.saveMessage(messageRepo, commonMessage, ctx)
}

// AddLogstashMessage Add logstash message
func (m *MessageController) AddLogstashMessage(ctx web.Context, messageRepo repository.MessageRepo) web.Response {
	flattenJSON, err := flatten.FlattenString(string(ctx.Request().Body()), "", flatten.DotStyle)
	if err != nil {
		return ctx.JSONError(fmt.Sprintf("invalid json: %s", err), http.StatusUnprocessableEntity)
	}

	var meta repository.MessageMeta
	if err := json.Unmarshal([]byte(flattenJSON), &meta); err != nil {
		return ctx.JSONError(fmt.Sprintf("parse json failed: %s", err), http.StatusInternalServerError)
	}

	contentField := ctx.InputWithDefault("content-field", "message")

	msg, ok := meta[contentField]
	if ok {
		delete(meta, contentField)
	} else {
		msg = "None"
	}

	return m.saveMessage(messageRepo, CommonMessage{
		Content: fmt.Sprintf("%v", msg),
		Meta:    meta,
		Tags:    nil,
		Origin:  "logstash",
	}, ctx)
}

// Add grafana message

type GrafanaMessage struct {
	EvalMatches []GrafanaEvalMatch `json:"evalMatches"`
	ImageUrl    string             `json:"imageUrl"`
	Message     string             `json:"message"`
	RuleID      int64              `json:"ruleId"`
	RuleName    string             `json:"ruleName"`
	RuleUrl     string             `json:"ruleUrl"`
	State       string             `json:"state"`
	Title       string             `json:"title"`
}

func (g GrafanaMessage) ToRepo() repository.Message {
	message, _ := json.Marshal(g)

	return repository.Message{
		Content: string(message),
		Meta: repository.MessageMeta{
			"rule_id":   strconv.Itoa(int(g.RuleID)),
			"rule_name": g.RuleName,
			"state":     g.State,
			"title":     g.Title,
		},
		Tags:   nil,
		Origin: "grafana",
	}
}

type GrafanaEvalMatch struct {
	Value  float64 `json:"value"`
	Metric string  `json:"metric"`
	Tags   map[string]string
}

func (m *MessageController) AddGrafanaMessage(ctx web.Context, messageRepo repository.MessageRepo) web.Response {
	var grafanaMessage GrafanaMessage
	if err := ctx.Unmarshal(&grafanaMessage); err != nil {
		return ctx.JSONError("invalid request", http.StatusUnprocessableEntity)
	}

	repoMessage := grafanaMessage.ToRepo()
	return m.saveMessage(messageRepo, CommonMessage{
		Content: repoMessage.Content,
		Meta:    repoMessage.Meta,
		Tags:    repoMessage.Tags,
		Origin:  repoMessage.Origin,
	}, ctx)
}

// add prometheus alert message

type PrometheusMessage struct {
	Status       string                 `json:"status"`
	Labels       repository.MessageMeta `json:"labels"`
	Annotations  repository.MessageMeta `json:"annotations"`
	StartsAt     time.Time              `json:"startsAt"`
	EndsAt       time.Time              `json:"endsAt"`
	GeneratorURL string                 `json:"generatorURL"`
}

func (pm PrometheusMessage) ToRepo() repository.Message {
	data, _ := json.Marshal(pm)
	return repository.Message{
		Content: string(data),
		Meta:    pm.Labels,
		Tags:    nil,
		Origin:  "prometheus",
	}
}

func (m *MessageController) AddPrometheusMessage(ctx web.Context, messageRepo repository.MessageRepo) web.Response {
	var prometheusMessage PrometheusMessage
	if err := ctx.Unmarshal(&prometheusMessage); err != nil {
		return ctx.JSONError("invalid request", http.StatusUnprocessableEntity)
	}

	repoMessage := prometheusMessage.ToRepo()
	return m.saveMessage(messageRepo, CommonMessage{
		Content: repoMessage.Content,
		Meta:    repoMessage.Meta,
		Tags:    repoMessage.Tags,
		Origin:  repoMessage.Origin,
	}, ctx)
}

// add prometheus-alert message

type PrometheusAlertMessage struct {
	Version  string `json:"version"`
	GroupKey string `json:"groupKey"`

	Receiver string              `json:"receiver"`
	Status   string              `json:"status"`
	Alerts   []PrometheusMessage `json:"alerts"`

	GroupLabels       repository.MessageMeta `json:"groupLabels"`
	CommonLabels      repository.MessageMeta `json:"commonLabels"`
	CommonAnnotations repository.MessageMeta `json:"commonAnnotations"`

	ExternalURL string `json:"externalURL"`
}

func (pam PrometheusAlertMessage) ToRepo() repository.Message {
	meta := make(repository.MessageMeta)
	for k, v := range pam.GroupLabels {
		meta[k] = v
	}

	for k, v := range pam.CommonLabels {
		meta[k] = v
	}

	meta["status"] = pam.Status

	data, _ := json.Marshal(pam)
	return repository.Message{
		Content: string(data),
		Meta:    meta,
		Tags:    nil,
		Origin:  "prometheus-alert",
	}
}

func (m *MessageController) AddPrometheusAlertMessage(ctx web.Context, messageRepo repository.MessageRepo) web.Response {
	var prometheusMessage PrometheusAlertMessage
	if err := ctx.Unmarshal(&prometheusMessage); err != nil {
		return ctx.JSONError("invalid request", http.StatusUnprocessableEntity)
	}

	repoMessage := prometheusMessage.ToRepo()
	return m.saveMessage(messageRepo, CommonMessage{
		Content: repoMessage.Content,
		Meta:    repoMessage.Meta,
		Tags:    repoMessage.Tags,
		Origin:  repoMessage.Origin,
	}, ctx)
}

// add open-falcon message

func (m *MessageController) AddOpenFalconMessage(ctx web.Context, messageRepo repository.MessageRepo) web.Response {
	tos := ctx.Input("tos")
	content := ctx.Input("content")

	if content == "" {
		return ctx.JSONError("invalid request, content required", http.StatusUnprocessableEntity)
	}

	meta := make(repository.MessageMeta)
	im := template.ParseOpenFalconImMessage(content)
	meta["status"] = im.Status
	meta["priority"] = strconv.Itoa(im.Priority)
	meta["endpoint"] = im.Endpoint
	meta["current_step"] = strconv.Itoa(im.CurrentStep)
	meta["body"] = im.Body
	meta["format_time"] = im.FormatTime

	return m.saveMessage(messageRepo, CommonMessage{
		Content: content,
		Meta:    meta,
		Tags:    []string{tos},
		Origin:  "open-falcon",
	}, ctx)
}
