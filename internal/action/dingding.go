package action

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/adanos-alert/pkg/messager/dingding"
	"github.com/mylxsw/adanos-alert/pkg/template"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/coll"
	"go.mongodb.org/mongo-driver/bson"
)

type DingdingAction struct {
	manager  Manager
	userRepo repository.UserRepo
}

type DingdingMeta struct {
	Template string `json:"template"`
	Token    string `json:"token"`
	Secret   string `json:"secret"`
}

func (d DingdingAction) Validate(meta string) error {
	var dingdingMeta DingdingMeta
	if err := json.Unmarshal([]byte(meta), &dingdingMeta); err != nil {
		return err
	}

	if dingdingMeta.Token == "" {
		return errors.New("dingding token required")
	}

	return nil
}

func NewDingdingAction(manager Manager) *DingdingAction {
	dingdingAction := DingdingAction{manager: manager}
	manager.MustResolve(func(userRepo repository.UserRepo) {
		dingdingAction.userRepo = userRepo
	})
	return &dingdingAction
}

func (d DingdingAction) Handle(rule repository.Rule, trigger repository.Trigger, grp repository.MessageGroup) error {
	payload := Payload{
		Action:  "dingding",
		Rule:    rule,
		Trigger: trigger,
		Group:   grp,
	}

	var meta DingdingMeta
	if err := json.Unmarshal([]byte(trigger.Meta), &meta); err != nil {
		return fmt.Errorf("parse dingding meta failed: %v", err)
	}

	summaryTemplate := rule.Template
	if strings.TrimSpace(meta.Template) != "" {
		summaryTemplate = meta.Template
	}

	res, err := template.Parse(summaryTemplate, payload)
	if err != nil {
		res = fmt.Sprintf("template parse failed: %s", err)
		log.WithFields(log.Fields{
			"err":      err.Error(),
			"template": rule.Template,
			"payload":  payload,
		}).Errorf("template parse failed: %v", err)
	}

	mobiles := make([]string, 0)
	if len(trigger.UserRefs) > 0 {
		users, err := d.userRepo.Find(bson.M{"_id": bson.M{"$in": trigger.UserRefs}})
		if err != nil {
			log.WithFields(log.Fields{
				"err":     err.Error(),
				"trigger": trigger,
			}).Errorf("load user from repo failed: %s", err)
		} else {
			if err := coll.MustNew(users).Filter(func(user repository.User) bool {
				if user.Phone != "" {
					return true
				}

				for _, m := range user.Metas {
					if strings.ToLower(m.Key) == "phone" {
						return true
					}
				}

				return false
			}).Map(func(user repository.User) string {
				if user.Phone != "" {
					return user.Phone
				}

				for _, m := range user.Metas {
					if strings.ToLower(m.Key) == "phone" {
						return m.Value
					}
				}
				return ""
			}).All(&mobiles); err != nil {
				log.WithFields(log.Fields{
					"err":     err.Error(),
					"trigger": trigger,
					"users":   users,
				}).Errorf("convert user's phone to array failed: %v", err)
			}
		}
	}

	msg := dingding.NewMarkdownMessage(rule.Name, res, mobiles)
	if err := dingding.NewDingding(meta.Token, meta.Secret).Send(msg); err != nil {
		log.WithFields(log.Fields{
			"title":   rule.Name,
			"content": res,
			"mobiles": mobiles,
			"err":     err,
		}).Errorf("send message to dingding failed: %v", err)
		return err
	}

	log.WithFields(log.Fields{
		"title":   rule.Name,
		"content": res,
		"mobiles": mobiles,
	}).Debug("send message to dingding succeed")

	return nil
}
