package action

import (
	"github.com/mylxsw/adanos-alert/internal/repository"
)

type SmsYunxin struct{
	manager Manager
}

func (w SmsYunxin) Validate(meta string, userRefs []string) error {
	panic("implement me")
}

func NewSmsYunxinAction(manager Manager) *SmsYunxin {
	return &SmsYunxin{manager: manager}
}

func (w SmsYunxin) Handle(rule repository.Rule, trigger repository.Trigger, grp repository.MessageGroup) error {
	panic("implement me")
}
