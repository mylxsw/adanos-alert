package action

import (
	"github.com/mylxsw/adanos-alert/internal/repository"
)

type HttpAction struct{}

func NewHttpAction() *HttpAction {
	return &HttpAction{}
}

func (act HttpAction) Handle(trigger repository.Trigger) error {
	panic("implement me")
}
