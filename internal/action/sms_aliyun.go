package action

import (
	"github.com/mylxsw/adanos-alert/internal/repository"
)

type SmsAliyun struct{
	manager Manager
}

func (w SmsAliyun) Validate(meta string, userRefs []string) error {
	panic("implement me")
}

func NewSmsAliyunAction(manager Manager) *SmsAliyun {
	return &SmsAliyun{manager: manager}
}

func (w SmsAliyun) Handle(rule repository.Rule, trigger repository.Trigger, grp repository.EventGroup) error {
	panic("implement me")
}

