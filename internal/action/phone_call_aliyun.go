package action

import (
	"encoding/json"
	"errors"

	"github.com/mylxsw/glacier/infra"

	"github.com/mylxsw/adanos-alert/configs"
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/adanos-alert/internal/template"
	"github.com/mylxsw/adanos-alert/pkg/messager/aliyun_voice"
	"github.com/mylxsw/asteria/log"
)

type AliyunVoiceCallAction struct {
	manager Manager
}

func (w AliyunVoiceCallAction) Validate(meta string, userRefs []string) error {
	if len(userRefs) == 0 {
		return errors.New("语音通知必须关联接收人")
	}

	return nil
}

func NewPhoneCallAliyunAction(manager Manager) *AliyunVoiceCallAction {
	return &AliyunVoiceCallAction{manager: manager}
}

type VoiceCallMeta struct {
	Title string `json:"title"`
}

func (w AliyunVoiceCallAction) Handle(rule repository.Rule, trigger repository.Trigger, grp repository.EventGroup) error {
	return w.manager.Resolve(func(resolver infra.Resolver, conf *configs.Config, userRepo repository.UserRepo, evtRepo repository.EventRepo) error {
		voiceCall := aliyun_voice.NewVoiceCall(conf)

		var meta VoiceCallMeta
		if err := json.Unmarshal([]byte(trigger.Meta), &meta); err != nil || meta.Title == "" {
			meta.Title = "{{ .Rule.Name }}"
		}

		title, err := template.Parse(w.manager, meta.Title, grp)
		if err != nil {
			log.WithFields(log.Fields{
				"rule_id": rule.ID.Hex(),
				"trigger": trigger,
			}).Errorf("parse aliyun voice call title failed: %v", err)
			title = rule.Name
		}

		mobiles := extractPhonesFromUserRefs(userRepo, getUserRefs(resolver, trigger, grp, evtRepo))
		if err := voiceCall.Send(title, mobiles); err != nil {
			log.WithFields(log.Fields{
				"title":   title,
				"mobiles": mobiles,
				"err":     err,
			}).Errorf("send message to aliyun voice failed: %v", err)
		}

		if log.DebugEnabled() {
			log.WithFields(log.Fields{
				"title":   title,
				"mobiles": mobiles,
			}).Debug("send message to aliyun voice succeed")
		}

		return nil
	})
}
