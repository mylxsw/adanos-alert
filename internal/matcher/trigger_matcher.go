package matcher

import (
	"sync"
	"time"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/mylxsw/adanos-alert/internal/repository"
)

// MessageMatcher is a matcher for trigger
type TriggerMatcher struct {
	program *vm.Program
	trigger repository.Trigger
}

type TriggerContext struct {
	Helpers
	Group           repository.MessageGroup
	MessageCallback func() []repository.Message
	messages        []repository.Message
	once            sync.Once
}

// NewTriggerContext create a new TriggerContext
func NewTriggerContext(group repository.MessageGroup, messageCallback func() []repository.Message) TriggerContext {
	return TriggerContext{Group: group, MessageCallback: messageCallback}
}

// Now return current time
func (tc *TriggerContext) Now() time.Time {
	return time.Now()
}

// ParseTime parse a string to time.Time
func (tc *TriggerContext) ParseTime(layout string, value string) time.Time {
	ts, _ := time.Parse(layout, value)
	return ts
}

// Messages return all messages in group
func (tc *TriggerContext) Messages() []repository.Message {
	tc.once.Do(func() {
		if tc.MessageCallback != nil {
			tc.messages = tc.MessageCallback()
		}
	})

	return tc.messages
}

// NewTriggerMatcher create a new TriggerMatcher
// https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md
func NewTriggerMatcher(trigger repository.Trigger) (*TriggerMatcher, error) {
	program, err := expr.Compile(trigger.PreCondition, expr.Env(&TriggerContext{}), expr.AsBool())
	if err != nil {
		return nil, err
	}

	return &TriggerMatcher{program: program, trigger: trigger}, nil
}

// Match check whether the msg is match with the rule
func (m *TriggerMatcher) Match(triggerCtx TriggerContext) (bool, error) {
	rs, err := expr.Run(m.program, &triggerCtx)
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
