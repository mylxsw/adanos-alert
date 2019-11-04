package action

import (
	"github.com/mylxsw/adanos-alert/internal/repository"
)

type PhoneCallAction struct{
	manager Manager
}

func (w PhoneCallAction) Validate(meta string) error {
	panic("implement me")
}

func NewPhoneCallAction(manager Manager) *PhoneCallAction {
	return &PhoneCallAction{manager:manager}
}

func (w PhoneCallAction) Handle(rule repository.Rule, trigger repository.Trigger, grp repository.MessageGroup) error {
	panic("implement me")
}
