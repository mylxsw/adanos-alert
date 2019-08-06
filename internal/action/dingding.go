package action

import (
	"github.com/mylxsw/adanos-alert/internal/repository"
)

type DingdingAction struct {}

func NewDingdingAction() *DingdingAction {
	return &DingdingAction{}
}

func (d DingdingAction) Handle(trigger repository.Trigger) error {
	panic("implement me")
}

