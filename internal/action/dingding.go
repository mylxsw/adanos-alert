package action

import (
	"github.com/mylxsw/adanos-alert/internal/repository"
)

type DingdingAction struct {
	manager *Manager
}

func NewDingdingAction(manager *Manager) *DingdingAction {
	return &DingdingAction{manager:manager}
}

func (d DingdingAction) Handle(trigger repository.Trigger, grp repository.MessageGroup) error {
	panic("implement me")
}

