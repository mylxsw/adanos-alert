package matcher

import (
	"errors"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/adanos-alert/pkg/json"
	"github.com/mylxsw/adanos-alert/pkg/misc"
)

// InvalidReturnVal is a error represents the expression return value is invalid
var InvalidReturnVal = errors.New("invalid return value: must be a bool value")

// MessageWrap is a wrapper to repository.Message
// We will add some helper function to message
type MessageWrap struct {
	repository.Message
	Helpers
}

func NewMessageWrap(message repository.Message) *MessageWrap {
	return &MessageWrap{Message: message}
}

// JsonGet parse message.Content as a json string and return the string value for key
func (msg *MessageWrap) JsonGet(key string, defaultValue string) string {
	return json.Gets(key, defaultValue, msg.Content)
}

// IsRecovery return whether the message is a recovery message
func (msg *MessageWrap) IsRecovery() bool {
	return msg.Type == repository.MessageTypeRecovery
}

// IsRecoverable return whether the message is recoverable
func (msg *MessageWrap) IsRecoverable() bool {
	return msg.Type == repository.MessageTypeRecoverable
}

// IsPlain return whether the message is a plain message
func (msg *MessageWrap) IsPlain() bool {
	return msg.Type == repository.MessageTypePlain || msg.Type == ""
}

// MessageMatcher is a matcher for repository.Message
type MessageMatcher struct {
	matchProgram  *vm.Program
	ignoreProgram *vm.Program
	rule          repository.Rule
}

// NewMessageMatcher create a new MessageMatcher
// https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md
func NewMessageMatcher(rule repository.Rule) (*MessageMatcher, error) {

	matchProgram, err := expr.Compile(
		misc.IfElse(rule.Rule == "", "true", rule.Rule).(string),
		expr.Env(&MessageWrap{}),
		expr.AsBool(),
	)
	if err != nil {
		return nil, err
	}

	ignoreProgram, err := expr.Compile(
		misc.IfElse(rule.IgnoreRule == "", "false", rule.IgnoreRule).(string),
		expr.Env(&MessageWrap{}),
		expr.AsBool(),
	)
	if err != nil {
		return nil, err
	}

	return &MessageMatcher{matchProgram: matchProgram, ignoreProgram: ignoreProgram, rule: rule}, nil
}

// Match check whether the msg is match with the rule
func (m *MessageMatcher) Match(msg repository.Message) (matched bool, ignored bool, err error) {
	wrapMsg := NewMessageWrap(msg)
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
func (m *MessageMatcher) Rule() repository.Rule {
	return m.rule
}
