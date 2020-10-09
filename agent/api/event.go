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
		router.Post("/", m.AddCommonEvent).Name("events:add:common")
		router.Post("/logstash/", m.AddLogstashEvent).Name("events:add:logstash")
		router.Post("/grafana/", m.AddGrafanaEvent).Name("events:add:grafana")
		router.Post("/prometheus/api/v1/alerts", m.AddPrometheusEvent).Name("events:add:prometheus") // url 地址末尾不包含 "/"
		router.Post("/prometheus_alertmanager/", m.AddPrometheusAlertEvent).Name("events:add:prometheus-alert")
		router.Post("/openfalcon/im/", m.AddOpenFalconEvent).Name("events:add:openfalcon")
	})
}

func (m *EventController) saveEvent(msgRepo store.EventStore, commonMessage misc.CommonEvent, ctx web.Context) error {
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

func (m *EventController) errorWrap(ctx web.Context, err error) web.Response {
	if err != nil {
		return ctx.JSONError(err.Error(), http.StatusInternalServerError)
	}

	return ctx.JSON(struct{}{})
}

func (m *EventController) AddCommonEvent(ctx web.Context, messageStore store.EventStore) web.Response {
	var commonMessage misc.CommonEvent
	if err := ctx.Unmarshal(&commonMessage); err != nil {
		return ctx.JSONError(fmt.Sprintf("invalid request: %v", err), http.StatusUnprocessableEntity)
	}

	return m.errorWrap(ctx, m.saveEvent(messageStore, commonMessage, ctx))
}

// AddLogstashEvent Add logstash message
func (m *EventController) AddLogstashEvent(ctx web.Context, messageStore store.EventStore) web.Response {
	commonMessage, err := misc.LogstashToCommonEvent(ctx.Request().Body(), ctx.InputWithDefault("content-field", "message"))
	if err != nil {
		return ctx.JSONError(err.Error(), http.StatusInternalServerError)
	}

	return m.errorWrap(ctx, m.saveEvent(messageStore, *commonMessage, ctx))
}

// Add grafana message
func (m *EventController) AddGrafanaEvent(ctx web.Context, messageStore store.EventStore) web.Response {
	commonMessage, err := misc.GrafanaToCommonEvent(ctx.Request().Body())
	if err != nil {
		return ctx.JSONError(err.Error(), http.StatusInternalServerError)
	}

	return m.errorWrap(ctx, m.saveEvent(messageStore, *commonMessage, ctx))
}

// add prometheus alert message
func (m *EventController) AddPrometheusEvent(ctx web.Context, messageStore store.EventStore) web.Response {
	commonMessages, err := misc.PrometheusToCommonEvents(ctx.Request().Body())
	if err != nil {
		return m.errorWrap(ctx, err)
	}

	for _, cm := range commonMessages {
		if err := m.saveEvent(messageStore, *cm, ctx); err != nil {
			log.WithFields(log.Fields{
				"message": cm,
			}).Errorf("save prometheus message failed: %v", err)
			continue
		}
	}

	return m.errorWrap(ctx, nil)
}

// add prometheus-alert message
func (m *EventController) AddPrometheusAlertEvent(ctx web.Context, messageStore store.EventStore) web.Response {
	commonMessage, err := misc.PrometheusAlertToCommonEvent(ctx.Request().Body())
	if err != nil {
		return ctx.JSONError(err.Error(), http.StatusInternalServerError)
	}

	return m.errorWrap(ctx, m.saveEvent(messageStore, *commonMessage, ctx))
}

// add open-falcon message
func (m *EventController) AddOpenFalconEvent(ctx web.Context, messageStore store.EventStore) web.Response {
	tos := ctx.Input("tos")
	content := ctx.Input("content")

	if content == "" {
		return ctx.JSONError("invalid request, content required", http.StatusUnprocessableEntity)
	}

	return m.errorWrap(ctx, m.saveEvent(messageStore, *misc.OpenFalconToCommonEvent(tos, content), ctx))
}
