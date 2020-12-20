package connector

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/mylxsw/adanos-alert/internal/extension"
	"github.com/mylxsw/asteria/log"
	"github.com/pkg/errors"
)

// Connector 是一个连接器对象，用于创建于 Adanos-alert 的连接
type Connector struct {
	servers []string
	token   string
}

// NewConnector create a new connector
func NewConnector(token string, servers ...string) *Connector {
	return &Connector{servers: servers, token: token}
}

// Send send a message to adanos server
func (conn *Connector) Send(ctx context.Context, evt *Event) error {
	return Send(ctx, conn.servers, conn.token, evt.Meta, evt.Tags, evt.Origin, evt.Ctl.toExtensionEventControl(), evt.Content)
}

// Event is a adanos alert message
type Event struct {
	Meta    map[string]interface{} `json:"meta" yaml:"meta"`
	Tags    []string               `json:"tags" yaml:"tags"`
	Origin  string                 `json:"origin" yaml:"origin"`
	Ctl     EventControl           `json:"ctl" yaml:"ctl"`
	Content string                 `json:"content" yaml:"content"`
}

type EventControl struct {
	ID              string `json:"id"`               // 消息标识，用于去重
	InhibitInterval string `json:"inhibit_interval"` // 抑制周期，周期内相同 ID 的消息直接丢弃
	RecoveryAfter   string `json:"recovery_after"`   // 自动恢复周期，该事件后一直没有发生相同标识的 消息，则自动生成一条恢复消息
}

func (ec EventControl) toExtensionEventControl() extension.EventControl {
	return extension.EventControl{
		ID:              ec.ID,
		InhibitInterval: ec.InhibitInterval,
		RecoveryAfter:   ec.RecoveryAfter,
	}
}

// NewEvent create a new Event
func NewEvent(content string) *Event {
	return &Event{Content: content, Tags: make([]string, 0), Meta: make(map[string]interface{})}
}

func (m *Event) WithTags(tags ...string) *Event {
	m.Tags = append(m.Tags, tags...)
	return m
}

func (m *Event) WithOrigin(origin string) *Event {
	m.Origin = origin
	return m
}

func (m *Event) WithCtl(ctl EventControl) *Event {
	m.Ctl = ctl
	return m
}

func (m *Event) WithMetas(metas map[string]interface{}) *Event {
	for k, v := range metas {
		m.Meta[k] = v
	}
	return m
}

func (m *Event) WithMeta(key string, value interface{}) *Event {
	m.Meta[key] = value
	return m
}

// Send send a message to adanos servers
func Send(ctx context.Context, servers []string, token string, meta map[string]interface{}, tags []string, origin string, ctl extension.EventControl, message string) error {
	evt := extension.CommonEvent{
		Content: message,
		Meta:    meta,
		Tags:    tags,
		Origin:  origin,
		Control: ctl,
	}
	data, _ := json.Marshal(evt)

	var err error
	for _, s := range servers {
		if err = sendEventToServer(ctx, evt, data, s, token); err == nil {
			break
		}

		log.Warningf("send to server %s failed: %v", s, err)
	}

	return err
}

func sendEventToServer(ctx context.Context, evt extension.CommonEvent, data []byte, adanosServer, adanosToken string) error {
	reqURL := fmt.Sprintf("%s/api/events/", strings.TrimRight(adanosServer, "/"))

	if log.DebugEnabled() {
		log.WithFields(log.Fields{
			"event": evt,
		}).Debugf("request: %v", reqURL)
	}

	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "POST", reqURL, bytes.NewReader(data))
	if err != nil {
		return errors.Wrap(err, "create request failed")
	}

	if adanosToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", adanosToken))
	}

	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "request failed")
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "read response body failed")
	}

	if log.DebugEnabled() {
		log.Debugf("response: %v", string(respBody))
	}
	return nil
}
