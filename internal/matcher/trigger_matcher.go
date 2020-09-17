package matcher

import (
	"sync"
	"time"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/adanos-alert/internal/template"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/container"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TriggerMatcher is a matcher for trigger
type TriggerMatcher struct {
	program *vm.Program
	trigger repository.Trigger
}

type TriggerContext struct {
	Helpers
	Group   repository.MessageGroup
	Trigger repository.Trigger

	MessageCallback     func() []repository.Message
	messageCallbackOnce sync.Once
	messages            []repository.Message

	cc container.Container
}

// NewTriggerContext create a new TriggerContext
func NewTriggerContext(cc container.Container, trigger repository.Trigger, group repository.MessageGroup, messageCallback func() []repository.Message) *TriggerContext {
	return &TriggerContext{cc: cc, Trigger: trigger, Group: group, MessageCallback: messageCallback}
}

// Messages return all messages in group
func (tc *TriggerContext) Messages() []repository.Message {
	tc.messageCallbackOnce.Do(func() {
		if tc.MessageCallback != nil {
			tc.messages = tc.MessageCallback()
		}
	})

	return tc.messages
}

// MessagesCount return the count in group
func (tc *TriggerContext) MessagesCount() int64 {
	var count int64 = 0
	tc.cc.MustResolve(func(msgRepo repository.MessageRepo) {
		count, _ = msgRepo.Count(bson.M{
			"group_ids": tc.Group.ID,
		})
	})

	return count
}

// MessagesMatchRegexCount get the count for messages matched regex
func (tc *TriggerContext) MessagesMatchRegexCount(regex string) int64 {
	var count int64 = 0
	tc.cc.MustResolve(func(msgRepo repository.MessageRepo) {
		filter := bson.M{
			"group_ids": tc.Group.ID,
			"content":   primitive.Regex{Pattern: regex, Options: "im"},
		}

		count, _ = msgRepo.Count(filter)
	})

	return count
}

// MessagesWithTagsCount get the count for messages which has tags
func (tc *TriggerContext) MessagesWithTagsCount(tags string) int64 {
	var count int64 = 0
	tc.cc.MustResolve(func(msgRepo repository.MessageRepo) {
		filter := bson.M{
			"group_ids": tc.Group.ID,
			"tags":      bson.M{"$in": template.StringTags(tags, ",")},
		}

		count, _ = msgRepo.Count(filter)
	})

	return count
}

// MessagesWithMetaCount get the count for messasges has a meta[key] equals to value
func (tc *TriggerContext) MessagesWithMetaCount(key, value string) int64 {
	var count int64 = 0
	tc.cc.MustResolve(func(msgRepo repository.MessageRepo) {
		filter := bson.M{
			"group_ids":   tc.Group.ID,
			"meta." + key: value,
		}

		count, _ = msgRepo.Count(filter)
	})

	return count
}

// TriggeredTimesInPeriod return triggered times in specified periods
func (tc *TriggerContext) TriggeredTimesInPeriod(periodInMinutes int, triggerStatus string) int64 {
	var triggeredTimes int64 = 0
	tc.cc.MustResolve(func(groupRepo repository.MessageGroupRepo) {
		filter := bson.M{
			"actions._id": tc.Trigger.ID,
			"updated_at":  bson.M{"$gt": time.Now().Add(-time.Duration(periodInMinutes) * time.Minute)},
		}

		if triggerStatus != "" {
			filter["actions.trigger_status"] = triggerStatus
		}

		n, _ := groupRepo.Count(filter)

		triggeredTimes = n
	})

	log.WithFields(log.Fields{
		"times": triggeredTimes,
	}).Debugf("TriggeredTimesInPeriod")

	return triggeredTimes
}

// LastTriggeredGroup get last triggeredGroup
func (tc *TriggerContext) LastTriggeredGroup(triggerStatus string) repository.MessageGroup {
	var lastTriggeredGroup repository.MessageGroup
	tc.cc.MustResolve(func(groupRepo repository.MessageGroupRepo) {
		filter := bson.M{
			"actions._id": tc.Trigger.ID,
		}

		if triggerStatus != "" {
			filter["actions.trigger_status"] = triggerStatus
		}

		grp, err := groupRepo.LastGroup(filter)
		if err == nil {
			lastTriggeredGroup = grp
		}
	})

	log.WithFields(log.Fields{
		"group": lastTriggeredGroup,
	}).Debugf("LastTriggeredGroup")

	return lastTriggeredGroup
}

// NewTriggerMatcher create a new TriggerMatcher
// https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md
func NewTriggerMatcher(trigger repository.Trigger) (*TriggerMatcher, error) {
	condition := trigger.PreCondition
	if condition == "" {
		condition = "true"
	}

	program, err := expr.Compile(condition, expr.Env(&TriggerContext{}), expr.AsBool())
	if err != nil {
		return nil, err
	}

	return &TriggerMatcher{program: program, trigger: trigger}, nil
}

// Match check whether the msg is match with the rule
func (m *TriggerMatcher) Match(triggerCtx *TriggerContext) (bool, error) {
	rs, err := expr.Run(m.program, triggerCtx)
	if err != nil {
		return false, err
	}

	if boolRes, ok := rs.(bool); ok {
		return boolRes, nil
	}

	return false, InvalidReturnVal
}

// Trigger return original trigger object
func (m *TriggerMatcher) Trigger() repository.Trigger {
	return m.trigger
}
