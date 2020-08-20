package api

import (
	"fmt"
	"net/http"

	"github.com/ledisdb/ledisdb/ledis"
	"github.com/mylxsw/adanos-alert/agent/store"
	"github.com/mylxsw/adanos-alert/pkg/misc"
	"github.com/mylxsw/adanos-alert/rpc/protocol"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/glacier/web"
)

type MessageController struct {
	cc container.Container
}

func NewMessageController(cc container.Container) web.Controller {
	return &MessageController{cc: cc}
}

func (m *MessageController) Register(router *web.Router) {
	router.Group("/messages", func(router *web.Router) {
		router.Post("/", m.AddCommonMessage).Name("messages:add:common")
		router.Post("/logstash/", m.AddLogstashMessage).Name("messages:add:logstash")
		router.Post("/grafana/", m.AddGrafanaMessage).Name("messages:add:grafana")
		router.Post("/prometheus/api/v1/alerts", m.AddPrometheusMessage).Name("messages:add:prometheus") // url 地址末尾不包含 "/"
		router.Post("/prometheus_alertmanager/", m.AddPrometheusAlertMessage).Name("messages:add:prometheus-alert")
		router.Post("/openfalcon/im/", m.AddOpenFalconMessage).Name("messages:add:openfalcon")
	})
}

func (m *MessageController) saveMessage(msgRepo store.MessageStore, commonMessage misc.CommonMessage, ctx web.Context) error {
	commonMessage.Meta["adanos_agent_version"] = m.cc.MustGet(infra.VersionKey).(string)
	commonMessage.Meta["adanos_agent_ip"] = misc.ServerIP()
	m.cc.MustResolve(func(db *ledis.DB) {
		agentID, _ := db.Get([]byte("agent-id"))
		commonMessage.Meta["adanos_agent_id"] = string(agentID)
	})
	req := protocol.MessageRequest{
		Data: commonMessage.Serialize(),
	}

	if err := msgRepo.Enqueue(&req); err != nil {
		log.Warningf("本地存储失败: %s", err)
		return err
	}

	return nil
}

func (m *MessageController) errorWrap(ctx web.Context, err error) web.Response {
	if err != nil {
		return ctx.JSONError(err.Error(), http.StatusInternalServerError)
	}

	return ctx.JSON(struct{}{})
}

func (m *MessageController) AddCommonMessage(ctx web.Context, messageStore store.MessageStore) web.Response {
	var commonMessage misc.CommonMessage
	if err := ctx.Unmarshal(&commonMessage); err != nil {
		return ctx.JSONError(fmt.Sprintf("invalid request: %v", err), http.StatusUnprocessableEntity)
	}

	return m.errorWrap(ctx, m.saveMessage(messageStore, commonMessage, ctx))
}

// AddLogstashMessage Add logstash message
func (m *MessageController) AddLogstashMessage(ctx web.Context, messageStore store.MessageStore) web.Response {
	commonMessage, err := misc.LogstashToCommonMessage(ctx.Request().Body(), ctx.InputWithDefault("content-field", "message"))
	if err != nil {
		return ctx.JSONError(err.Error(), http.StatusInternalServerError)
	}

	return m.errorWrap(ctx, m.saveMessage(messageStore, *commonMessage, ctx))
}

// Add grafana message
func (m *MessageController) AddGrafanaMessage(ctx web.Context, messageStore store.MessageStore) web.Response {
	commonMessage, err := misc.GrafanaToCommonMessage(ctx.Request().Body())
	if err != nil {
		return ctx.JSONError(err.Error(), http.StatusInternalServerError)
	}

	return m.errorWrap(ctx, m.saveMessage(messageStore, *commonMessage, ctx))
}

// add prometheus alert message
func (m *MessageController) AddPrometheusMessage(ctx web.Context, messageStore store.MessageStore) web.Response {
	commonMessages, err := misc.PrometheusToCommonMessages(ctx.Request().Body())
	if err != nil {
		return m.errorWrap(ctx, err)
	}

	for _, cm := range commonMessages {
		if err := m.saveMessage(messageStore, *cm, ctx); err != nil {
			log.WithFields(log.Fields{
				"message": cm,
			}).Errorf("save prometheus message failed: %v", err)
			continue
		}
	}

	return m.errorWrap(ctx, nil)
}

// add prometheus-alert message
func (m *MessageController) AddPrometheusAlertMessage(ctx web.Context, messageStore store.MessageStore) web.Response {
	commonMessage, err := misc.PrometheusAlertToCommonMessage(ctx.Request().Body())
	if err != nil {
		return ctx.JSONError(err.Error(), http.StatusInternalServerError)
	}

	return m.errorWrap(ctx, m.saveMessage(messageStore, *commonMessage, ctx))
}

// add open-falcon message
func (m *MessageController) AddOpenFalconMessage(ctx web.Context, messageStore store.MessageStore) web.Response {
	tos := ctx.Input("tos")
	content := ctx.Input("content")

	if content == "" {
		return ctx.JSONError("invalid request, content required", http.StatusUnprocessableEntity)
	}

	return m.errorWrap(ctx, m.saveMessage(messageStore, *misc.OpenFalconToCommonMessage(tos, content), ctx))
}
