package connector

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/mylxsw/adanos-alert/pkg/misc"
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
	return Send(ctx, conn.servers, conn.token, evt.meta, evt.tags, evt.origin, evt.ctl, evt.content)
}

// Event is a adanos alert message
type Event struct {
	meta    map[string]interface{}
	tags    []string
	origin  string
	ctl     misc.EventControl
	content string
}

// NewEvent create a new Event
func NewEvent(content string) *Event {
	return &Event{content: content, tags: make([]string, 0), meta: make(map[string]interface{})}
}

func (m *Event) WithTags(tags ...string) *Event {
	m.tags = append(m.tags, tags...)
	return m
}

func (m *Event) WithOrigin(origin string) *Event {
	m.origin = origin
	return m
}

func (m *Event) WithCtl(ctl misc.EventControl) *Event {
	m.ctl = ctl
	return m
}

func (m *Event) WithMetas(metas map[string]interface{}) *Event {
	for k, v := range metas {
		m.meta[k] = v
	}
	return m
}

func (m *Event) WithMeta(key string, value interface{}) *Event {
	m.meta[key] = value
	return m
}

// Send send a message to adanos servers
func Send(ctx context.Context, servers []string, token string, meta map[string]interface{}, tags []string, origin string, ctl misc.EventControl, message string) error {
	commonMessage := misc.CommonEvent{
		Content: message,
		Meta:    meta,
		Tags:    tags,
		Origin:  origin,
		Control: ctl,
	}
	data, _ := json.Marshal(commonMessage)

	var err error
	for _, s := range servers {
		if err = sendEventToServer(ctx, commonMessage, data, s, token); err == nil {
			break
		}

		log.Warningf("send to server %s failed: %v", s, err)
	}

	return err
}

func sendEventToServer(ctx context.Context, commonMessage misc.CommonEvent, data []byte, adanosServer, adanosToken string) error {
	reqURL := fmt.Sprintf("%s/api/events/", strings.TrimRight(adanosServer, "/"))

	log.WithFields(log.Fields{
		"message": commonMessage,
	}).Debugf("request: %v", reqURL)

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

	log.Debugf("response: %v", string(respBody))
	return nil
}
