package matcher

import (
	"fmt"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/mylxsw/adanos-alert/internal/repository"
)

// MessageFinger 消息指纹
type MessageFinger struct {
	expr    string
	program *vm.Program
}

// NewMessageFinger create a new MessageFinger instance
func NewMessageFinger(fingerExpr string) (*MessageFinger, error) {
	if fingerExpr == "" {
		fingerExpr = `""`
	}

	program, err := expr.Compile(fingerExpr, expr.Env(&MessageWrap{}))
	if err != nil {
		return nil, err
	}

	return &MessageFinger{
		expr:    fingerExpr,
		program: program,
	}, nil
}

// Run 根据指定的表达式创建 message 的指纹
func (m *MessageFinger) Run(msg repository.Message) (string, error) {
	result, err := expr.Run(m.program, NewMessageWrap(msg))
	if err != nil {
		return "", err
	}

	if result == nil {
		return "", nil
	}

	return fmt.Sprintf("%v", result), nil
}
