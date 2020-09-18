package action

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/mylxsw/adanos-alert/internal/queue"
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/adanos-alert/pubsub"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier/event"
	"go.mongodb.org/mongo-driver/bson"
)

type Action interface {
	Validate(meta string, userRefs []string) error
	Handle(rule repository.Rule, trigger repository.Trigger, grp repository.MessageGroup) error
}

type Manager interface {
	Resolve(f interface{}) error
	MustResolve(f interface{})
	Get(key interface{}) (interface{}, error)
	Dispatch(action string) Action
	Run(action string) Action
	Register(name string, action Action)
}

type actionManager struct {
	cc      container.Container
	lock    sync.RWMutex
	actions map[string]Action
}

func NewManager(cc container.Container) Manager {
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

type QueueAction struct {
	action  string
	manager Manager
}

func (q *QueueAction) Validate(meta string, userRefs []string) error {
	return q.manager.Run(q.action).Validate(meta, userRefs)
}

type Payload struct {
	msgRepo            repository.MessageRepo
	Action             string                  `json:"action"`
	Rule               repository.Rule         `json:"rule"`
	Trigger            repository.Trigger      `json:"trigger"`
	Group              repository.MessageGroup `json:"group"`
	RuleTemplateParsed string                  `json:"rule_template_parsed"`
	PreviewURL         string                  `json:"preview_url"`
	ReportURL          string                  `json:"report_url"`
}

// Init initialize a payload
func (payload *Payload) Init(manager Manager) {
	manager.MustResolve(func(msgRepo repository.MessageRepo) {
		payload.msgRepo = msgRepo
	})
}

// MessageType return message type in group
func (payload *Payload) MessageType() string {
	return string(payload.Group.Type)
}

// IsRecovery return whether the messages in group is recovery message
func (payload *Payload) IsRecovery() bool {
	return payload.Group.Type == repository.MessageTypeRecovery
}

// IsRecoverable return whether the messages in group is recoverable message
func (payload *Payload) IsRecoverable() bool {
	return payload.Group.Type == repository.MessageTypeRecoverable
}

// IsPlain return whether the messages in group is plain message
func (payload *Payload) IsPlain() bool {
	return payload.Group.Type == repository.MessageTypePlain || payload.Group.Type == ""
}

// Messages get messages for group
func (payload *Payload) Messages(limit int64) []repository.Message {
	messages, _, err := payload.msgRepo.Paginate(bson.M{"group_ids": payload.Group.ID}, 0, limit)
	if err != nil {
		log.WithFields(log.Fields{
			"group_id": payload.Group.ID,
			"limit":    limit,
			"error":    err,
		}).Errorf("query messages failed for template: %v", err)
		return []repository.Message{}
	}

	return messages
}

func (payload *Payload) Encode() []byte {
	data, _ := json.Marshal(payload)
	return data
}

func (payload *Payload) Decode(data []byte) error {
	return json.Unmarshal(data, payload)
}

func (q *QueueAction) Handle(rule repository.Rule, trigger repository.Trigger, grp repository.MessageGroup) error {
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

		log.WithFields(log.Fields{
			"action": q.action,
			"id":     id,
		}).Debug("enqueue a action to queue")

		return nil
	})
}
