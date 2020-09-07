package misc

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jeremywohl/flatten"
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/adanos-alert/pkg/array"
	"github.com/mylxsw/adanos-alert/pkg/template"
)

type CommonMessage struct {
	Content string                 `json:"content"`
	Meta    repository.MessageMeta `json:"meta"`
	Tags    []string               `json:"tags"`
	Origin  string                 `json:"origin"`
}

func (msg CommonMessage) Serialize() string {
	data, _ := json.Marshal(msg)
	return string(data)
}

func (msg CommonMessage) ToRepo() repository.Message {
	return repository.Message{
		Content: msg.Content,
		Meta:    msg.Meta,
		Tags:    msg.Tags,
		Origin:  msg.Origin,
	}
}

type RepoMessage interface {
	ToRepo() repository.Message
}

func LogstashToCommonMessage(content []byte, contentField string) (*CommonMessage, error) {
	flattenJSON, err := flatten.FlattenString(string(content), "", flatten.DotStyle)
	if err != nil {
		return nil, fmt.Errorf("invalid json: %s", err)
	}

	var meta repository.MessageMeta
	if err := json.Unmarshal([]byte(flattenJSON), &meta); err != nil {
		return nil, fmt.Errorf("parse json failed: %s", err)
	}

	msg, ok := meta[contentField]
	if ok {
		delete(meta, contentField)
	} else {
		msg = "None"
	}

	return &CommonMessage{
		Content: fmt.Sprintf("%v", msg),
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
func logstashMetaFilter(meta repository.MessageMeta) repository.MessageMeta {
	res := make(repository.MessageMeta)
	for k, v := range meta {
		if array.StringsContainPrefix(k, excludeLogstashPrefix) {
			continue
		}

		res[k] = v
	}

	return res
}

type GrafanaMessage struct {
	EvalMatches []GrafanaEvalMatch `json:"evalMatches"`
	ImageURL    string             `json:"imageUrl"`
	Message     string             `json:"message"`
	RuleID      int64              `json:"ruleId"`
	RuleName    string             `json:"ruleName"`
	RuleURL     string             `json:"ruleUrl"`
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

func GrafanaToCommonMessage(content []byte) (*CommonMessage, error) {
	var grafanaMessage GrafanaMessage
	if err := json.Unmarshal(content, &grafanaMessage); err != nil {
		return nil, errors.New("invalid request")
	}

	repoMessage := grafanaMessage.ToRepo()
	return &CommonMessage{
		Content: repoMessage.Content,
		Meta:    repoMessage.Meta,
		Tags:    repoMessage.Tags,
		Origin:  repoMessage.Origin,
	}, nil
}

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

func PrometheusToCommonMessages(content []byte) ([]*CommonMessage, error) {
	var prometheusMessages []PrometheusMessage
	if err := json.Unmarshal(content, &prometheusMessages); err != nil {
		return nil, errors.New("invalid request")
	}

	commonMessages := make([]*CommonMessage, 0)
	for _, pm := range prometheusMessages {
		repoMessage := pm.ToRepo()
		commonMessages = append(commonMessages, &CommonMessage{
			Content: repoMessage.Content,
			Meta:    repoMessage.Meta,
			Tags:    repoMessage.Tags,
			Origin:  repoMessage.Origin,
		})
	}

	return commonMessages, nil
}

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

func PrometheusAlertToCommonMessage(content []byte) (*CommonMessage, error) {
	var prometheusMessage PrometheusAlertMessage
	if err := json.Unmarshal(content, &prometheusMessage); err != nil {
		return nil, errors.New("invalid request")
	}

	repoMessage := prometheusMessage.ToRepo()
	return &CommonMessage{
		Content: repoMessage.Content,
		Meta:    repoMessage.Meta,
		Tags:    repoMessage.Tags,
		Origin:  repoMessage.Origin,
	}, nil
}

func OpenFalconToCommonMessage(tos, content string) *CommonMessage {
	meta := make(repository.MessageMeta)
	im := template.ParseOpenFalconImMessage(content)
	meta["status"] = im.Status
	meta["priority"] = strconv.Itoa(im.Priority)
	meta["endpoint"] = im.Endpoint
	meta["current_step"] = strconv.Itoa(im.CurrentStep)
	meta["body"] = im.Body
	meta["format_time"] = im.FormatTime

	return &CommonMessage{
		Content: content,
		Meta:    meta,
		Tags:    []string{tos},
		Origin:  "open-falcon",
	}
}
