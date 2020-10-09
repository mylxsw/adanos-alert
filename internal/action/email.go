package action

import (
	"github.com/mylxsw/adanos-alert/internal/repository"
)

type EmailAction struct {
	manager Manager
}

func (e EmailAction) Validate(meta string, userRefs []string) error {
	panic("implement me")
}

func NewEmailAction(manager Manager) *EmailAction {
	return &EmailAction{manager:manager}
}

func (e EmailAction) Handle(rule repository.Rule, trigger repository.Trigger, grp repository.EventGroup) error {
	panic("implement me")
}

