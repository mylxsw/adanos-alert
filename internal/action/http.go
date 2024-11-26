package action

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/mylxsw/adanos-alert/configs"
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/go-utils/str"
)

// HTTPAction HTTP 动作
type HTTPAction struct {
	manager Manager
}

// Validate 参数校验
func (act HTTPAction) Validate(meta string, userRefs []string) error {
	var httpMeta HTTPMeta
	if err := json.Unmarshal([]byte(meta), &httpMeta); err != nil {
		return err
	}

	httpURL := strings.TrimSpace(httpMeta.URL)
	if httpURL == "" {
		return errors.New("URL is required")
	}

	if _, err := url.Parse(httpURL); err != nil {
		return fmt.Errorf("invalid url: %w", err)
	}

	method := strings.ToUpper(httpMeta.Method)
	if !str.In(method, []string{"GET", "POST", "PUT", "PATCH", "HEAD", "DELETE", "OPTIONS", "PURGE"}) {
		return fmt.Errorf("not support such request method: %s", method)
	}

	return nil
}

// NewHTTPAction create a new HTTPAction
func NewHTTPAction(manager Manager) *HTTPAction {
	return &HTTPAction{manager: manager}
}

// Handle 动作处理
func (act HTTPAction) Handle(rule repository.Rule, trigger repository.Trigger, grp repository.EventGroup) error {
	var meta HTTPMeta
	if err := json.Unmarshal([]byte(trigger.Meta), &meta); err != nil {
		return fmt.Errorf("parse http meta failed: %v", err)
	}

	return act.manager.Resolve(func(conf *configs.Config, evtRepo repository.EventRepo) error {
		payload, _ := createPayloadAndSummary(act.manager, "http", conf, evtRepo, rule, trigger, grp)
		body := parseTemplate(act.manager, meta.Body, payload)

		var reqBody io.Reader
		if body != "" {
			reqBody = strings.NewReader(body)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, meta.Method, parseTemplate(act.manager, meta.URL, payload), reqBody)
		if err != nil {
			log.WithFields(log.Fields{
				"trigger": trigger,
				"rule_id": rule.ID.Hex(),
			}).Errorf("create request failed: %v", err)
			return fmt.Errorf("create request failed: %v", err)
		}

		for _, header := range meta.Headers {
			req.Header.Add(header.Key, parseTemplate(act.manager, header.Value, payload))
		}

		client := &http.Client{}
		client.Timeout = 5 * time.Second
		resp, err := client.Do(req)
		if err != nil {
			log.WithFields(log.Fields{
				"trigger": trigger,
				"rule_id": rule.ID.Hex(),
			}).Errorf("send http request failed: %v", err)
			return fmt.Errorf("send http request failed: %v", err)
		}

		defer func() {
			_ = resp.Body.Close()
		}()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.WithFields(log.Fields{
				"trigger": trigger,
				"rule_id": rule.ID.Hex(),
				"resp":    string(respBody),
			}).Errorf("read http response failed: %v", err)
			return nil
		}

		if log.DebugEnabled() {
			log.WithFields(log.Fields{
				"trigger": trigger,
				"rule_id": rule.ID.Hex(),
				"resp":    string(respBody),
			}).Debug("send message to http succeed")
		}

		return nil
	})
}

// HTTPHeaderMeta HTTP 请求头
type HTTPHeaderMeta struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// HTTPMeta HTTP 元数据
type HTTPMeta struct {
	URL     string           `json:"url"`
	Method  string           `json:"method"`
	Headers []HTTPHeaderMeta `json:"headers"`
	Body    string           `json:"body"`
}
