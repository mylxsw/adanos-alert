package action

import (
	"github.com/mylxsw/adanos-alert/internal/repository"
)

type WechatAction struct{}

func NewWechatAction() *WechatAction {
	return &WechatAction{}
}

func (w WechatAction) Handle(trigger repository.Trigger) error {
	panic("implement me")
}
