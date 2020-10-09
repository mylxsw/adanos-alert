package misc

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jeremywohl/flatten"
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/adanos-alert/internal/template"
	"github.com/mylxsw/adanos-alert/pkg/strarr"
)

type CommonEvent struct {
	Content string               `json:"content"`
	Meta    repository.EventMeta `json:"meta"`
	Tags    []string             `json:"tags"`
	Origin  string               `json:"origin"`

	Control EventControl `json:"control"`
}

type EventControl struct {
	ID              string `json:"id"`               // 消息标识，用于去重
	InhibitInterval string `json:"inhibit_interval"` // 抑制周期，周期内相同 ID 的消息直接丢弃
	RecoveryAfter   string `json:"recovery_after"`   // 自动恢复周期，该事件后一直没有发生相同标识的 消息，则自动生成一条恢复消息
}

func (mc EventControl) GetInhibitInterval() time.Duration {
	duration, err := time.ParseDuration(mc.InhibitInterval)
	if err != nil {
		return 0
	}

	return duration
}

func (mc EventControl) GetRecoveryAfter() time.Duration {
	duration, err := time.ParseDuration(mc.RecoveryAfter)
	if err != nil {
		return 0
	}

	return duration
}

func (evt CommonEvent) Serialize() string {
	data, _ := json.Marshal(evt)
	return string(data)
}

func (evt CommonEvent) CreateRepoEvent() repository.Event {
	return repository.Event{
		Content: evt.Content,
		Meta:    evt.Meta,
		Tags:    evt.Tags,
		Origin:  evt.Origin,
		Type: IfElse(
			evt.Control.ID != "" && evt.Control.GetRecoveryAfter() > 0,
			repository.EventTypeRecoverable,
			repository.EventTypePlain,
		).(repository.EventType),
	}
}

func (evt CommonEvent) GetControl() EventControl {
	return evt.Control
}

type RepoEvent interface {
	CreateRepoEvent() repository.Event
	GetControl() EventControl
}

func LogstashToCommonEvent(content []byte, contentField string) (*CommonEvent, error) {
	flattenJSON, err := flatten.FlattenString(string(content), "", flatten.DotStyle)
	if err != nil {
		return nil, fmt.Errorf("invalid json: %s", err)
	}

	var meta repository.EventMeta
	if err := json.Unmarshal([]byte(flattenJSON), &meta); err != nil {
		return nil, fmt.Errorf("parse json failed: %s", err)
	}

	evt, ok := meta[contentField]
	if ok {
		delete(meta, contentField)
	} else {
		evt = "None"
	}

	return &CommonEvent{
		Content: fmt.Sprintf("%v", evt),
		Meta:    logstashMetaFilter(meta),
		Tags:    nil,
		Origin:  "logstash",
	}, nil
}

// 待排除的 Logstash 字段
var excludeLogstashPrefix = []string{
	"beat.",
	"host.",
	"input.",
	"log.flags.",
	"log.file.",
	"offset",
	"prospector.",
	"tags.",
	"@version",
	"@timestamp",
}

// logstashMetaFilter 过滤掉不需要存储的 logstash 专有字段
func logstashMetaFilter(meta repository.EventMeta) repository.EventMeta {
	res := make(repository.EventMeta)
	for k, v := range meta {
		if strarr.HasPrefixes(k, excludeLogstashPrefix) {
			continue
		}

		res[k] = v
	}

	return res
}

type GrafanaEvent struct {
	EvalMatches []GrafanaEvalMatch `json:"evalMatches"`
	ImageURL    string             `json:"imageUrl"`
	Message     string             `json:"message"`
	RuleID      int64              `json:"ruleId"`
	RuleName    string             `json:"ruleName"`
	RuleURL     string             `json:"ruleUrl"`
	State       string             `json:"state"`
	Title       string             `json:"title"`
}

func (g GrafanaEvent) ToRepo() repository.Event {
	message, _ := json.Marshal(g)

	return repository.Event{
		Content: string(message),
		Meta: repository.EventMeta{
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

func GrafanaToCommonEvent(content []byte) (*CommonEvent, error) {
	var grafanaMessage GrafanaEvent
	if err := json.Unmarshal(content, &grafanaMessage); err != nil {
		return nil, errors.New("invalid request")
	}

	repoMessage := grafanaMessage.ToRepo()
	return &CommonEvent{
		Content: repoMessage.Content,
		Meta:    repoMessage.Meta,
		Tags:    repoMessage.Tags,
		Origin:  repoMessage.Origin,
	}, nil
}

type PrometheusEvent struct {
	Status       string               `json:"status"`
	Labels       repository.EventMeta `json:"labels"`
	Annotations  repository.EventMeta `json:"annotations"`
	StartsAt     time.Time            `json:"startsAt"`
	EndsAt       time.Time            `json:"endsAt"`
	GeneratorURL string               `json:"generatorURL"`
}

func (pm PrometheusEvent) CreateRepoEvent() repository.Event {
	data, _ := json.Marshal(pm)
	return repository.Event{
		Content: string(data),
		Meta:    pm.Labels,
		Tags:    nil,
		Origin:  "prometheus",
	}
}

func (pm PrometheusEvent) GetControl() EventControl {
	mc := EventControl{}
	if evtID, ok := pm.Labels["adanos_id"]; ok {
		mc.ID = fmt.Sprintf("%v", evtID)

		recoveryAfter, err := time.ParseDuration(fmt.Sprintf("%v", pm.Labels["adanos_recovery_after"]))
		if err != nil || recoveryAfter < 0 {
			recoveryAfter = 0
		}

		inhibitIntervalKey := "adanos_inhibit_interval"
		_, ok := pm.Labels[inhibitIntervalKey]
		if !ok {
			inhibitIntervalKey = "adanos_repeat_interval"
		}
		inhibitInterval, err := time.ParseDuration(fmt.Sprintf("%v", pm.Labels[inhibitIntervalKey]))
		if err != nil || inhibitInterval < 0 {
			inhibitInterval = 0
		}

		mc.InhibitInterval = inhibitInterval.String()
		mc.RecoveryAfter = recoveryAfter.String()
	}

	return mc
}

func PrometheusToCommonEvents(content []byte) ([]*CommonEvent, error) {
	var prometheusMessages []PrometheusEvent
	if err := json.Unmarshal(content, &prometheusMessages); err != nil {
		return nil, errors.New("invalid request")
	}

	commonMessages := make([]*CommonEvent, 0)
	for _, pm := range prometheusMessages {
		repoMessage := pm.CreateRepoEvent()
		commonMessages = append(commonMessages, &CommonEvent{
			Content: repoMessage.Content,
			Meta:    repoMessage.Meta,
			Tags:    repoMessage.Tags,
			Origin:  repoMessage.Origin,
			Control: pm.GetControl(),
		})
	}

	return commonMessages, nil
}

type PrometheusAlertEvent struct {
	Version  string `json:"version"`
	GroupKey string `json:"groupKey"`

	Receiver string            `json:"receiver"`
	Status   string            `json:"status"`
	Alerts   []PrometheusEvent `json:"alerts"`

	GroupLabels       repository.EventMeta `json:"groupLabels"`
	CommonLabels      repository.EventMeta `json:"commonLabels"`
	CommonAnnotations repository.EventMeta `json:"commonAnnotations"`

	ExternalURL string `json:"externalURL"`
}

func (pam PrometheusAlertEvent) ToRepo() repository.Event {
	meta := make(repository.EventMeta)
	for k, v := range pam.GroupLabels {
		meta[k] = v
	}

	for k, v := range pam.CommonLabels {
		meta[k] = v
	}

	meta["status"] = pam.Status

	data, _ := json.Marshal(pam)
	return repository.Event{
		Content: string(data),
		Meta:    meta,
		Tags:    nil,
		Origin:  "prometheus-alert",
	}
}

func PrometheusAlertToCommonEvent(content []byte) (*CommonEvent, error) {
	var prometheusMessage PrometheusAlertEvent
	if err := json.Unmarshal(content, &prometheusMessage); err != nil {
		return nil, errors.New("invalid request")
	}

	repoMessage := prometheusMessage.ToRepo()
	return &CommonEvent{
		Content: repoMessage.Content,
		Meta:    repoMessage.Meta,
		Tags:    repoMessage.Tags,
		Origin:  repoMessage.Origin,
	}, nil
}

func OpenFalconToCommonEvent(tos, content string) *CommonEvent {
	meta := make(repository.EventMeta)
	im := template.ParseOpenFalconImMessage(content)
	meta["status"] = im.Status
	meta["priority"] = strconv.Itoa(im.Priority)
	meta["endpoint"] = im.Endpoint
	meta["current_step"] = strconv.Itoa(im.CurrentStep)
	meta["body"] = im.Body
	meta["format_time"] = im.FormatTime

	return &CommonEvent{
		Content: content,
		Meta:    meta,
		Tags:    []string{tos},
		Origin:  "open-falcon",
	}
}
