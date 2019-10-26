package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jeremywohl/flatten"
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/container"
	"github.com/mylxsw/hades"
)

type MessageController struct {
	cc *container.Container
}

func NewMessageController(cc *container.Container) hades.Controller {
	return &MessageController{cc: cc}
}

func (m *MessageController) Register(router *hades.Router) {
	router.Group("/messages", func(router *hades.Router) {
		router.Post("/common/", m.AddCommonMessage)
		router.Post("/logstash/", m.AddLogstashMessage)
	})
}

type RepoMessage interface {
	ToRepo() repository.Message
}

func (m *MessageController) saveMessage(messageRepo repository.MessageRepo, repoMessage RepoMessage, ctx hades.Context) hades.Response {
	id, err := messageRepo.Add(repoMessage.ToRepo())
	if err != nil {
		return ctx.JSONError(err.Error(), http.StatusInternalServerError)
	}

	return ctx.JSON(hades.M{
		"id": id.String(),
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

func (m *MessageController) AddCommonMessage(ctx hades.Context, messageRepo repository.MessageRepo) hades.Response {
	var commonMessage CommonMessage
	if err := ctx.Unmarshal(&commonMessage); err != nil {
		return ctx.JSONError("invalid request", http.StatusUnprocessableEntity)
	}

	return m.saveMessage(messageRepo, commonMessage, ctx)
}

// Add logstash message

func (m *MessageController) AddLogstashMessage(ctx hades.Context, messageRepo repository.MessageRepo) hades.Response {
	flattenJson, err := flatten.FlattenString(string(ctx.Request().Body()), "", flatten.DotStyle)
	if err != nil {
		return ctx.JSONError(fmt.Sprintf("invalid json: %s", err), http.StatusUnprocessableEntity)
	}

	var meta repository.MessageMeta
	if err := json.Unmarshal([]byte(flattenJson), &meta); err != nil {
		return ctx.JSONError(fmt.Sprintf("parse json failed: %s", err), http.StatusInternalServerError)
	}

	msg := meta["message"]
	delete(meta, "message")

	return m.saveMessage(messageRepo, CommonMessage{
		Content: msg,
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
