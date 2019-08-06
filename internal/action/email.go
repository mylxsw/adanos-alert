package action

import (
	"github.com/mylxsw/adanos-alert/internal/repository"
)

type EmailAction struct {}

func NewEmailAction() *EmailAction {
	return &EmailAction{}
}

func (e EmailAction) Handle(trigger repository.Trigger) error {
	panic("implement me")
}

