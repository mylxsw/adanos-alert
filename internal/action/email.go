package action

import (
	"github.com/mylxsw/adanos-alert/configs"
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/asteria/log"
)

type EmailAction struct {
	manager Manager
}

func (e EmailAction) Validate(meta string, userRefs []string) error {
	return nil
}

func NewEmailAction(manager Manager) *EmailAction {
	return &EmailAction{manager: manager}
}

func (e EmailAction) Handle(rule repository.Rule, trigger repository.Trigger, grp repository.EventGroup) error {
	return e.manager.Resolve(func(conf *configs.Config) error {
		//client := email.NewClient(conf.EmailSMTP.Host, conf.EmailSMTP.Port, conf.EmailSMTP.Username, conf.EmailSMTP.Password)
		//if err := client.Send(subject, body ,users...); err != nil {
		//
		//}

		if log.DebugEnabled() {
			log.WithFields(log.Fields{
				"title":   rule.Name,
			}).Debug("send message to dingding succeed")
		}

		return nil
	})
}
