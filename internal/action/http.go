package action

import (
	"github.com/mylxsw/adanos-alert/internal/repository"
)

type HttpAction struct {
	manager *Manager
}

func NewHttpAction(manager *Manager) *HttpAction {
	return &HttpAction{manager: manager}
}

func (act HttpAction) Handle(trigger repository.Trigger, grp repository.MessageGroup) error {
	panic("implement me")
}
