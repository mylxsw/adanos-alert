package matcher

import (
	"fmt"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/mylxsw/adanos-alert/internal/repository"
)

// EventFinger Event 指纹
type EventFinger struct {
	expr    string
	program *vm.Program
}

// NewEventFinger create a new EventFinger instance
func NewEventFinger(fingerExpr string) (*EventFinger, error) {
	if fingerExpr == "" {
		fingerExpr = `""`
	}

	program, err := expr.Compile(fingerExpr, expr.Env(&EventWrap{}))
	if err != nil {
		return nil, err
	}

	return &EventFinger{
		expr:    fingerExpr,
		program: program,
	}, nil
}

// Run 根据指定的表达式创建 Event 的指纹
func (m *EventFinger) Run(msg repository.Event) (string, error) {
	result, err := expr.Run(m.program, NewEventWrap(msg))
	if err != nil {
		return "", err
	}

	if result == nil {
		return "", nil
	}

	return fmt.Sprintf("%v", result), nil
}
