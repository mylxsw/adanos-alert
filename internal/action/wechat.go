package action

import (
	"github.com/mylxsw/adanos-alert/internal/repository"
)

type WechatAction struct{
	manager Manager
}

func (w WechatAction) Validate(meta string, userRefs []string) error {
	panic("implement me")
}

func NewWechatAction(manager Manager) *WechatAction {
	return &WechatAction{manager:manager}
}

func (w WechatAction) Handle(rule repository.Rule, trigger repository.Trigger, grp repository.EventGroup) error {
	panic("implement me")
}
