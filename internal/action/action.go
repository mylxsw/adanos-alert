package action

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/mylxsw/adanos-alert/configs"
	"github.com/mylxsw/adanos-alert/internal/queue"
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/adanos-alert/internal/template"
	"github.com/mylxsw/adanos-alert/pubsub"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/glacier/event"
	"github.com/mylxsw/glacier/infra"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Action 触发动作接口
type Action interface {
	Validate(meta string, userRefs []string) error
	Handle(rule repository.Rule, trigger repository.Trigger, grp repository.EventGroup) error
}

// Manager 动作管理器接口
type Manager interface {
	Resolve(f interface{}) error
	MustResolve(f interface{})
	Get(key interface{}) (interface{}, error)
	Dispatch(action string) Action
	Run(action string) Action
	Register(name string, action Action)
}

type actionManager struct {
	cc      infra.Resolver
	lock    sync.RWMutex
	actions map[string]Action
}

// NewManager create a new Manager
func NewManager(cc infra.Resolver) Manager {
	return &actionManager{cc: cc, actions: make(map[string]Action)}
}

func (manager *actionManager) Get(key interface{}) (interface{}, error) {
	return manager.cc.Get(key)
}

func (manager *actionManager) Resolve(f interface{}) error {
	return manager.cc.ResolveWithError(f)
}

func (manager *actionManager) MustResolve(f interface{}) {
	manager.cc.MustResolve(f)
}

// Dispatch dispatch a action to queue
func (manager *actionManager) Dispatch(action string) Action {
	return &QueueAction{
		action:  action,
		manager: manager,
	}
}

// Run execute a action
func (manager *actionManager) Run(action string) Action {
	manager.lock.RLock()
	defer manager.lock.RUnlock()

	return manager.actions[action]
}

// Register register a new action
func (manager *actionManager) Register(name string, action Action) {
	manager.lock.Lock()
	defer manager.lock.Unlock()

	manager.actions[name] = action
}

// QueueAction 动作队列
type QueueAction struct {
	action  string
	manager Manager
}

// Validate 校验动作参数
func (q *QueueAction) Validate(meta string, userRefs []string) error {
	return q.manager.Run(q.action).Validate(meta, userRefs)
}

// EventQuerier 事件查询接口
type EventQuerier func(groupID primitive.ObjectID, limit int64) []repository.Event

// Payload 事件描述生成时使用的对象，用于模板解析
type Payload struct {
	eventQuerier       EventQuerier
	Action             string                `json:"action"`
	Rule               repository.Rule       `json:"rule"`
	Trigger            repository.Trigger    `json:"trigger"`
	Group              repository.EventGroup `json:"group"`
	RuleTemplateParsed string                `json:"rule_template_parsed"`
	PreviewURL         string                `json:"preview_url"`
	ReportURL          string                `json:"report_url"`
}

// Init initialize a payload
func (payload *Payload) Init(eventQuerier EventQuerier) {
	payload.eventQuerier = eventQuerier
}

// MessageType return message type in group
// This method is depressed
func (payload *Payload) MessageType() string {
	return payload.EventType()
}

// EventType return message type in group
func (payload *Payload) EventType() string {
	return string(payload.Group.Type)
}

// IsRecovery return whether the messages in group is recovery message
func (payload *Payload) IsRecovery() bool {
	return payload.Group.Type == repository.EventTypeRecovery
}

// IsRecoverable return whether the messages in group is recoverable message
func (payload *Payload) IsRecoverable() bool {
	return payload.Group.Type == repository.EventTypeRecoverable
}

// IsPlain return whether the messages in group is plain message
func (payload *Payload) IsPlain() bool {
	return payload.Group.Type == repository.EventTypePlain || payload.Group.Type == ""
}

// Messages get messages for group
// This method is depressed
func (payload *Payload) Messages(limit int64) []repository.Event {
	return payload.Events(limit)
}

// Events get messages for group
func (payload *Payload) Events(limit int64) []repository.Event {
	return payload.eventQuerier(payload.Group.ID, limit)
}

// FirstEvent get first event
func (payload *Payload) FirstEvent() repository.Event {
	return payload.Events(1)[0]
}

// CreateRepositoryEventQuerier 创建仓库事件查询器
func CreateRepositoryEventQuerier(msgRepo repository.EventRepo) func(groupID primitive.ObjectID, limit int64) []repository.Event {
	return func(groupID primitive.ObjectID, limit int64) []repository.Event {
		messages, _, err := msgRepo.Paginate(bson.M{"group_ids": groupID}, 0, limit)
		if err != nil {
			log.WithFields(log.Fields{
				"group_id": groupID,
				"limit":    limit,
				"error":    err,
			}).Errorf("query messages failed for template: %v", err)
			return []repository.Event{}
		}
		return messages
	}
}

// Encode 将 Payload 编码
func (payload *Payload) Encode() []byte {
	data, _ := json.Marshal(payload)
	return data
}

// Decode 从字符串解析出 Payload 对象
func (payload *Payload) Decode(data []byte) error {
	return json.Unmarshal(data, payload)
}

// Handle 动作处理
func (q *QueueAction) Handle(rule repository.Rule, trigger repository.Trigger, grp repository.EventGroup) error {
	return q.manager.Resolve(func(queueManager queue.Manager, em event.Manager) error {
		payload := Payload{
			Action:  q.action,
			Trigger: trigger,
			Group:   grp,
			Rule:    rule,
		}

		em.Publish(pubsub.MessageGroupTriggeredEvent{
			Action:    q.action,
			Trigger:   trigger,
			Group:     grp,
			Rule:      rule,
			CreatedAt: time.Now(),
		})

		id, err := queueManager.Enqueue(repository.QueueJob{
			Name:    "action",
			Payload: string(payload.Encode()),
		})
		if err != nil {
			return err
		}

		if log.DebugEnabled() {
			log.WithFields(log.Fields{
				"action": q.action,
				"id":     id,
			}).Debug("enqueue a action to queue")
		}

		return nil
	})
}

// CreatePayload 创建一个 Payload
func CreatePayload(conf *configs.Config, eventQuerier EventQuerier, action string, rule repository.Rule, trigger repository.Trigger, grp repository.EventGroup) *Payload {
	payload := &Payload{
		Action:  action,
		Rule:    rule,
		Trigger: trigger,
		Group:   grp,
	}
	payload.Init(eventQuerier)

	if conf.PreviewURL != "" {
		payload.PreviewURL = fmt.Sprintf(conf.PreviewURL, grp.ID.Hex())
	}
	if conf.ReportURL != "" {
		payload.ReportURL = fmt.Sprintf(conf.ReportURL, grp.ID.Hex())
	}

	return payload
}

// createPayloadAndSummary 创建 Payload 并且生成 summary
func createPayloadAndSummary(cc template.SimpleContainer, actionName string, conf *configs.Config, evtRepo repository.EventRepo, rule repository.Rule, trigger repository.Trigger, grp repository.EventGroup) (*Payload, string) {
	payload := CreatePayload(conf, CreateRepositoryEventQuerier(evtRepo), actionName, rule, trigger, grp)
	payload.RuleTemplateParsed = parseTemplate(cc, rule.Template, payload)

	return payload, payload.RuleTemplateParsed
}

// parseTemplate 模板解释
func parseTemplate(cc template.SimpleContainer, temp string, payload *Payload) string {
	summary, err := template.Parse(cc, temp, payload)
	if err != nil {
		summary = fmt.Sprintf("<internal> template parse failed: %s", err)
		log.WithFields(log.Fields{
			"err":      err.Error(),
			"template": temp,
			"payload":  payload,
		}).Errorf("<internal> template parse failed: %v", err)
	}

	return summary
}
