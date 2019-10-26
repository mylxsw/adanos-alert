package action

import (
	"github.com/mylxsw/adanos-alert/internal/repository"
)

type WechatAction struct{
	manager *Manager
}

func NewWechatAction(manager *Manager) *WechatAction {
	return &WechatAction{manager:manager}
}

func (w WechatAction) Handle(trigger repository.Trigger, grp repository.MessageGroup) error {
	panic("implement me")
}
