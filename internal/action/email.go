package action

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/iris-contrib/blackfriday"
	"github.com/microcosm-cc/bluemonday"
	"github.com/mylxsw/adanos-alert/configs"
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/adanos-alert/pkg/messager/email"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/go-utils/array"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EmailAction struct {
	manager Manager
	client  *email.Client
}

// EmailMeta 邮件发送元数据
type EmailMeta struct {
	Template string `json:"template"`
}

func (e EmailAction) Validate(meta string, userRefs []string) error {
	return nil
}

func NewEmailAction(manager Manager, conf *configs.Config) *EmailAction {
	client := email.NewClient(
		conf.EmailSMTP.Host,
		conf.EmailSMTP.Port,
		conf.EmailSMTP.Username,
		conf.EmailSMTP.Password,
	)

	return &EmailAction{manager: manager, client: client}
}

func (e EmailAction) Handle(rule repository.Rule, trigger repository.Trigger, grp repository.EventGroup) error {
	var meta EmailMeta
	if err := json.Unmarshal([]byte(trigger.Meta), &meta); err != nil {
		return fmt.Errorf("parse email meta failed: %v", err)
	}

	return e.manager.Resolve(func(resolver infra.Resolver, conf *configs.Config, msgRepo repository.EventRepo, userRepo repository.UserRepo) error {
		payload, summary := createPayloadAndSummary(e.manager, "email", conf, msgRepo, rule, trigger, grp)
		if strings.TrimSpace(meta.Template) != "" {
			summary = parseTemplate(e.manager, meta.Template, payload)
		}

		summary = string(bluemonday.UGCPolicy().SanitizeBytes(blackfriday.Run([]byte(summary))))
		emails := extractEmailsFromUserRefs(userRepo, getUserRefs(resolver, trigger, grp, msgRepo))
		if err := e.client.Send(rule.Name, summary, emails...); err != nil {
			log.WithFields(log.Fields{
				"subject": rule.Name,
				"body":    summary,
				"emails":  emails,
				"err":     err,
			}).Errorf("send message to email failed: %v", err)
			return err
		}

		if log.DebugEnabled() {
			log.WithFields(log.Fields{
				"title": rule.Name,
			}).Debug("send message to email succeed")
		}

		return nil
	})
}

func extractEmailsFromUserRefs(userRepo repository.UserRepo, userRefs []primitive.ObjectID) []string {
	if len(userRefs) == 0 {
		return []string{}
	}

	users, err := userRepo.Find(bson.M{"_id": bson.M{"$in": userRefs}})
	if err != nil {
		log.WithFields(log.Fields{
			"err":      err.Error(),
			"userRefs": userRefs,
		}).Errorf("load user from repo failed: %s", err)

		return []string{}
	}

	userFilterFunc := func(user repository.User) bool {
		if user.Email != "" {
			return true
		}

		for _, m := range user.Metas {
			if strings.ToLower(m.Key) == "email" {
				return true
			}
		}

		return false
	}

	userMapFunc := func(user repository.User) string {
		if user.Email != "" {
			return user.Email
		}

		for _, m := range user.Metas {
			if strings.ToLower(m.Key) == "email" {
				return m.Value
			}
		}
		return ""
	}

	return array.Map(array.Filter(users, userFilterFunc), userMapFunc)
}
