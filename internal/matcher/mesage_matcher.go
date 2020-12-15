package matcher

import (
	jsonEnc "encoding/json"
	"errors"
	"sync"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/adanos-alert/pkg/helper"
	"github.com/mylxsw/adanos-alert/pkg/json"
	"github.com/mylxsw/adanos-alert/pkg/misc"
)

// InvalidReturnVal is a error represents the expression return value is invalid
var InvalidReturnVal = errors.New("invalid return value: must be a bool value")

// EventWrap is a wrapper to repository.Event
// We will add some helper function to message
type EventWrap struct {
	repository.Event
	helper.Helpers
	fullJSONOnce sync.Once
	fullJSON     string
}

func NewEventWrap(message repository.Event) *EventWrap {
	return &EventWrap{Event: message}
}

// FullJSON return whole event as json document
func (msg *EventWrap) FullJSON() string {
	msg.fullJSONOnce.Do(func() {
		res, _ := jsonEnc.Marshal(msg)
		msg.fullJSON = string(res)
	})
	return msg.fullJSON
}

// JsonGet parse message.Content as a json string and return the string value for key
func (msg *EventWrap) JsonGet(key string, defaultValue string) string {
	return json.Gets(key, defaultValue, msg.Content)
}

// IsRecovery return whether the message is a recovery message
func (msg *EventWrap) IsRecovery() bool {
	return msg.Type == repository.EventTypeRecovery
}

// IsRecoverable return whether the message is recoverable
func (msg *EventWrap) IsRecoverable() bool {
	return msg.Type == repository.EventTypeRecoverable
}

// IsPlain return whether the message is a plain message
func (msg *EventWrap) IsPlain() bool {
	return msg.Type == repository.EventTypePlain || msg.Type == ""
}

// EventMatcher is a matcher for repository.Event
type EventMatcher struct {
	matchProgram  *vm.Program
	ignoreProgram *vm.Program
	rule          repository.Rule
}

// NewEventMatcher create a new EventMatcher
// https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md
func NewEventMatcher(rule repository.Rule) (*EventMatcher, error) {

	matchProgram, err := expr.Compile(
		misc.IfElse(rule.Rule == "", "true", rule.Rule).(string),
		expr.Env(&EventWrap{}),
		expr.AsBool(),
	)
	if err != nil {
		return nil, err
	}

	ignoreProgram, err := expr.Compile(
		misc.IfElse(rule.IgnoreRule == "", "false", rule.IgnoreRule).(string),
		expr.Env(&EventWrap{}),
		expr.AsBool(),
	)
	if err != nil {
		return nil, err
	}

	return &EventMatcher{matchProgram: matchProgram, ignoreProgram: ignoreProgram, rule: rule}, nil
}

// Match check whether the msg is match with the rule
func (m *EventMatcher) Match(evt repository.Event) (matched bool, ignored bool, err error) {
	wrapMsg := NewEventWrap(evt)
	rs, err := expr.Run(m.matchProgram, wrapMsg)
	if err != nil {
		return false, false, err
	}

	if boolRes, ok := rs.(bool); ok {
		if boolRes {
			ignore, err := expr.Run(m.ignoreProgram, wrapMsg)
			if err != nil {
				return boolRes, false, err
			}

			return boolRes, ignore.(bool), nil
		}

		return boolRes, false, nil
	}

	return false, false, InvalidReturnVal
}

// Rule return original rule object
func (m *EventMatcher) Rule() repository.Rule {
	return m.rule
}
