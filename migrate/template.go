package migrate

import (
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/asteria/log"
	"go.mongodb.org/mongo-driver/bson"
)

var predefinedTemplates = []repository.Template{
	{
		Name:        "判断来源",
		Description: "判断来源是否为某个值",
		Content:     `Origin == "internal"`,
		Type:        repository.TemplateTypeMatchRule,
	},
	{
		Name:        "单位时间内触发次数判断",
		Description: "30分钟内触发失败次数小于5次",
		Content:     `TriggeredTimesInPeriod(30, "failed") < 5`,
		Type:        repository.TemplateTypeTriggerRule,
	},
	{
		Name:        "测试模板",
		Description: "测试展示模板",
		Content:     `Hello, Action={{ .Action }}, Group.MessageCount={{ .Group.MessageCount }}`,
		Type:        repository.TemplateTypeTemplate,
	},
}

func initPredefinedTemplates(repo repository.TemplateRepo) {
	for _, t := range predefinedTemplates {
		t.Predefined = true
		temps, err := repo.Find(bson.M{"name": t.Name, "predefined": true})
		if err == repository.ErrNotFound || len(temps) == 0 {
			id, err := repo.Add(t)
			if err != nil {
				log.WithFields(log.Fields{
					"temp": t,
				}).Errorf("add predefined template %s failed: %v", t.Name, err)
				continue
			}

			log.WithFields(log.Fields{
				"temp": t,
			}).Debugf("add predefined template %s with id %s", t.Name, id.Hex())
		} else if err != nil {
			log.WithFields(log.Fields{
				"temp": t,
			}).Errorf("skip predefined template %s, because query failed: %v", t.Name, err)
		} else {
			tt := temps[0]
			changed := false

			if tt.Type != t.Type {
				changed = true
				tt.Type = t.Type
			}

			if tt.Content != t.Content {
				changed = true
				tt.Content = t.Content
			}

			if tt.Description != t.Description {
				changed = true
				tt.Description = t.Description
			}

			if changed {
				if err := repo.Update(tt.ID, tt); err != nil {
					log.WithFields(log.Fields{
						"temp": t,
					}).Errorf("query predefined template failed: %v", err)
				}

				log.WithFields(log.Fields{
					"temp": t,
				}).Errorf("update predefined template %s failed: %v", t.Name, err)
			}
		}
	}
}
