package action

import (
	"github.com/mylxsw/adanos-alert/internal/repository"
)

type PhoneCallAliyunAction struct{
	manager Manager
}

func (w PhoneCallAliyunAction) Validate(meta string) error {
	panic("implement me")
}

func NewPhoneCallAliyunAction(manager Manager) *PhoneCallAliyunAction {
	return &PhoneCallAliyunAction{manager: manager}
}

func (w PhoneCallAliyunAction) Handle(rule repository.Rule, trigger repository.Trigger, grp repository.MessageGroup) error {
	panic("implement me")
}
