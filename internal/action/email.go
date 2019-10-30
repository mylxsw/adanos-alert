package action

import (
	"github.com/mylxsw/adanos-alert/internal/repository"
)

type EmailAction struct {
	manager Manager
}

func (e EmailAction) Validate(meta string) error {
	panic("implement me")
}

func NewEmailAction(manager Manager) *EmailAction {
	return &EmailAction{manager:manager}
}

func (e EmailAction) Handle(rule repository.Rule, trigger repository.Trigger, grp repository.MessageGroup) error {
	panic("implement me")
}

