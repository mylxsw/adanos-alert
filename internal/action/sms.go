package action

import (
	"github.com/mylxsw/adanos-alert/internal/repository"
)

type SmsAction struct{
	manager Manager
}

func (w SmsAction) Validate(meta string) error {
	panic("implement me")
}

func NewSmsAction(manager Manager) *SmsAction {
	return &SmsAction{manager:manager}
}

func (w SmsAction) Handle(rule repository.Rule, trigger repository.Trigger, grp repository.MessageGroup) error {
	panic("implement me")
}
