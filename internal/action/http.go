package action

import (
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/asteria/log"
)

type HttpAction struct {
	manager Manager
}

func (act HttpAction) Validate(meta string) error {
	return nil
}

func NewHttpAction(manager Manager) *HttpAction {
	return &HttpAction{manager: manager}
}

func (act HttpAction) Handle(rule repository.Rule, trigger repository.Trigger, grp repository.MessageGroup) error {
	log.WithFields(log.Fields{
		"rule":    rule,
		"trigger": trigger,
		"grp":     grp,
	}).Warningf("http action triggered")
	return nil
}
