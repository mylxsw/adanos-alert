package rule

import (
	"errors"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/adanos-alert/pkg/json"
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

// MessageMatcher is a matcher for repository.Message
type MessageMatcher struct {
	program *vm.Program
	rule    repository.Rule
}

// NewMessageMatcher create a new MessageMatcher
// https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md
func NewMessageMatcher(rule repository.Rule) (*MessageMatcher, error) {
	program, err := expr.Compile(rule.Rule, expr.Env(&MessageWrap{}), expr.AsBool())
	if err != nil {
		return nil, err
	}

	return &MessageMatcher{program: program, rule: rule}, nil
}

// Match check whether the msg is match with the rule
func (m *MessageMatcher) Match(msg repository.Message) (bool, error) {
	rs, err := expr.Run(m.program, NewMessageWrap(msg))
	if err != nil {
		return false, err
	}

	if boolRes, ok := rs.(bool); ok {
		return boolRes, nil
	}

	return false, InvalidReturnVal
}

// Rule return original rule object
func (m *MessageMatcher) Rule() repository.Rule {
	return m.rule
}
