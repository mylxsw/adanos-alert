package action

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/mylxsw/adanos-alert/configs"
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/adanos-alert/internal/template"
	"github.com/mylxsw/adanos-alert/pkg/messager/dingding"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/coll"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DingdingAction struct {
	manager  Manager
	userRepo repository.UserRepo
}

type DingdingMeta struct {
	Template string `json:"template"`
	RobotID  string `json:"robot_id"`
}

func (d DingdingAction) Validate(meta string, userRefs []string) error {
	var dingdingMeta DingdingMeta
	if err := json.Unmarshal([]byte(meta), &dingdingMeta); err != nil {
		return err
	}

	if dingdingMeta.RobotID == "" {
		return errors.New("dingding robot required")
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

func (d DingdingAction) Handle(rule repository.Rule, trigger repository.Trigger, grp repository.EventGroup) error {

	var meta DingdingMeta
	if err := json.Unmarshal([]byte(trigger.Meta), &meta); err != nil {
		return fmt.Errorf("parse dingding meta failed: %v", err)
	}

	return d.manager.Resolve(func(conf *configs.Config, msgRepo repository.EventRepo, robotRepo repository.DingdingRobotRepo) error {
		// get robot for dingding
		robotID, err := primitive.ObjectIDFromHex(meta.RobotID)
		if err != nil {
			return fmt.Errorf("invalid robot id: %s, error is %v", meta.RobotID, err)
		}

		robot, err := robotRepo.Get(robotID)
		if err != nil {
			return fmt.Errorf("query robot for id=%s failed: %v", meta.RobotID, err)
		}

		payload := CreatePayload(conf, CreateRepositoryMessageQuerier(msgRepo), "dingding", rule, trigger, grp)
		ruleTemplateContent, err := template.Parse(d.manager, rule.Template, payload)
		if err != nil {
			ruleTemplateContent = fmt.Sprintf("<rule> template parse failed: %s", err)
			log.WithFields(log.Fields{
				"err":      err.Error(),
				"template": rule.Template,
				"payload":  payload,
			}).Errorf("<rule> template parse failed: %v", err)
		}

		notifyContent := ruleTemplateContent
		if strings.TrimSpace(meta.Template) != "" {
			payload.RuleTemplateParsed = ruleTemplateContent
			notifyContent, err = template.Parse(d.manager, meta.Template, payload)
			if err != nil {
				notifyContent = fmt.Sprintf("<trigger> template parse failed: %s", err)
				log.WithFields(log.Fields{
					"err":      err.Error(),
					"template": rule.Template,
					"payload":  payload,
				}).Errorf("<trigger> template parse failed: %v", err)
			}
		}

		mobiles := extractPhonesFromUserRefs(d.userRepo, trigger.UserRefs)
		msg := dingding.NewMarkdownMessage(rule.Name, notifyContent, mobiles)
		if err := dingding.NewDingding(robot.Token, robot.Secret).Send(msg); err != nil {
			log.WithFields(log.Fields{
				"title":   rule.Name,
				"content": notifyContent,
				"mobiles": mobiles,
				"err":     err,
			}).Errorf("send message to dingding failed: %v", err)
			return err
		}

		log.WithFields(log.Fields{
			"title":   rule.Name,
			"content": notifyContent,
			"mobiles": mobiles,
		}).Debug("send message to dingding succeed")

		return nil
	})
}

func extractPhonesFromUserRefs(userRepo repository.UserRepo, userRefs []primitive.ObjectID) []string {
	mobiles := make([]string, 0)
	if len(userRefs) == 0 {
		return mobiles
	}

	users, err := userRepo.Find(bson.M{"_id": bson.M{"$in": userRefs}})
	if err != nil {
		log.WithFields(log.Fields{
			"err":      err.Error(),
			"userRefs": userRefs,
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
				"err":   err.Error(),
				"users": users,
			}).Errorf("convert user's phone to array failed: %v", err)
		}
	}

	return mobiles
}
