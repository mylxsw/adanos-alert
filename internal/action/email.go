package action

import (
	"github.com/mylxsw/adanos-alert/internal/repository"
)

type EmailAction struct {
	manager *Manager
}

func NewEmailAction(manager *Manager) *EmailAction {
	return &EmailAction{manager:manager}
}

func (e EmailAction) Handle(trigger repository.Trigger, grp repository.MessageGroup) error {
	panic("implement me")
}

