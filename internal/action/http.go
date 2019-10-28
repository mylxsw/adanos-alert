package action

import (
	"github.com/mylxsw/adanos-alert/internal/repository"
)

type HttpAction struct {
	manager *Manager
}

func (act HttpAction) Validate(meta string) error {
	panic("implement me")
}

func NewHttpAction(manager *Manager) *HttpAction {
	return &HttpAction{manager: manager}
}

func (act HttpAction) Handle(rule repository.Rule, trigger repository.Trigger, grp repository.MessageGroup) error {
	panic("implement me")
}
