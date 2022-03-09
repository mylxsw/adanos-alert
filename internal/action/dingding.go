package action

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mylxsw/container"
	"strings"

	"github.com/mylxsw/adanos-alert/configs"
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/adanos-alert/pkg/messager/dingding"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/coll"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DingdingAction 钉钉发送动作
type DingdingAction struct {
	manager  Manager
	userRepo repository.UserRepo
}

// DingdingMeta 钉钉发送元数据
type DingdingMeta struct {
	Template string `json:"template"`
	RobotID  string `json:"robot_id"`
}

// Validate 校验动作参数
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

// NewDingdingAction create a new dingdingAction
func NewDingdingAction(manager Manager) *DingdingAction {
	dingdingAction := DingdingAction{manager: manager}
	manager.MustResolve(func(userRepo repository.UserRepo) {
		dingdingAction.userRepo = userRepo
	})
	return &dingdingAction
}

// Handle 钉钉动作处理
func (d DingdingAction) Handle(rule repository.Rule, trigger repository.Trigger, grp repository.EventGroup) error {

	var meta DingdingMeta
	if err := json.Unmarshal([]byte(trigger.Meta), &meta); err != nil {
		return fmt.Errorf("parse dingding meta failed: %v", err)
	}

	return d.manager.Resolve(func(cc container.Container, conf *configs.Config, msgRepo repository.EventRepo, robotRepo repository.DingdingRobotRepo) error {
		// get robot for dingding
		robotID, err := primitive.ObjectIDFromHex(meta.RobotID)
		if err != nil {
			return fmt.Errorf("invalid robot id: %s, error is %v", meta.RobotID, err)
		}

		robot, err := robotRepo.Get(robotID)
		if err != nil {
			return fmt.Errorf("query robot for id=%s failed: %v", meta.RobotID, err)
		}

		payload, summary := createPayloadAndSummary(d.manager, "dingding", conf, msgRepo, rule, trigger, grp)
		if strings.TrimSpace(meta.Template) != "" {
			summary = parseTemplate(d.manager, meta.Template, payload)
		}

		mobiles := extractPhonesFromUserRefs(d.userRepo, getUserRefs(cc, trigger, grp, msgRepo))
		msg := dingding.NewMarkdownMessage(rule.Name, summary, mobiles)
		if err := dingding.NewDingding(robot.Token, robot.Secret).Send(msg); err != nil {
			log.WithFields(log.Fields{
				"title":   rule.Name,
				"content": summary,
				"mobiles": mobiles,
				"err":     err,
			}).Errorf("send message to dingding failed: %v", err)
			return err
		}

		if log.DebugEnabled() {
			log.WithFields(log.Fields{
				"title":   rule.Name,
				"content": summary,
				"mobiles": mobiles,
			}).Debug("send message to dingding succeed")
		}

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
