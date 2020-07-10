package migrate

import (
	"github.com/mylxsw/adanos-alert/configs"
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/asteria/log"
	"go.mongodb.org/mongo-driver/bson"
)

var predefinedTemplates = []repository.Template{
	{
		Name:        "判断来源",
		Description: "来源为 logstash",
		Content:     `Origin == "logstash"`,
		Type:        repository.TemplateTypeMatchRule,
	},
	{
		Name:        "判断Meta是否等于某个值",
		Description: "判断日志类型为 nginx_access",
		Content:     `Meta["log_type"] == "nginx_access"`,
		Type:        repository.TemplateTypeMatchRule,
	},
	{
		Name:        "判断Meta是否在某个范围内",
		Description: "日志级别为 ERROR 或 FATAL",
		Content:     `Meta["log_level"] in ["ERROR", "FATAL"]`,
		Type:        repository.TemplateTypeMatchRule,
	},
	{
		Name:        "判断是否包含标签",
		Description: "包含名为 java 的标签",
		Content:     `"java" in Tags`,
		Type:        repository.TemplateTypeMatchRule,
	},
	{
		Name:        "判断 message 内容是否匹配正则表达式",
		Description: `message 以 "Error:" 开头`,
		Content:     `Content matches "^Error:"`,
		Type:        repository.TemplateTypeMatchRule,
	},
	{
		Name:        "判断 message 内容是否不包含字符串",
		Description: `判断 message 中不包含 "关键词" 字符串`,
		Content:     `not (Content contains "关键词")`,
		Type:        repository.TemplateTypeMatchRule,
	},
	{
		Name:        "单位时间内触发次数判断",
		Description: "30分钟内触发失败次数小于5次",
		Content:     `TriggeredTimesInPeriod(30, "failed") < 5`,
		Type:        repository.TemplateTypeTriggerRule,
	},
	{
		Name:        "判断当前时间是否在 某个时间段",
		Description: "每天晚上 10:00 到 次日早上 9:00",
		Content:     `DailyTimeBetween("22:00", "9:00")`,
		Type:        repository.TemplateTypeTriggerRule,
	},
	{
		Name:        "判断分组中 Messages 数量是否大于某个值",
		Description: "当前分组中有超过 10 条 Messages",
		Content:     `MessagesCount() > 10`,
		Type:        repository.TemplateTypeTriggerRule,
	},
	{
		Name:        "展示概要信息",
		Description: "展示分组的概要信息",
		Content:     `当前通知方式 {{ .Action }}, 当前分组包含的消息数量 {{ .Group.MessageCount }}，触发的规则名称 {{ .Rule.Name }}  ({{ .Rule.Rule }})`,
		Type:        repository.TemplateTypeTemplate,
	},
	{
		Name:        "报警信息列表",
		Description: "展示报警信息列表",
		Content: `## {{ .Rule.Name }}

{{ range $i, $msg := .Messages 4 }}- 来源：**{{ $msg.Origin }}**，标签：{{ $msg.Tags  }}

{{ cutoff 400 $msg.Content | ident "    > " }}

{{ end }}

---

[共 {{ .Group.MessageCount }} 条，查看详细]({{ .PreviewURL }})`,
		Type: repository.TemplateTypeTemplate,
	},
	{
		Name:        "报警详情链接",
		Description: "报警详细信息链接地址",
		Content:     `[共 {{ .Group.MessageCount }} 条，查看详细]({{ .PreviewURL }})`,
		Type:        repository.TemplateTypeTemplate,
	},
	{
		Name:        "MySQL Guard - Deadlock Logs",
		Description: "MySQL Guard 死锁日志",
		Content:     `## {{ .Rule.Name }} ({{ .Group.MessageCount }})

序号: {{ .Group.ID.Hex }}

规则: {{ .Group.Rule.ID.Hex }}

{{ range $i, $msg := .Messages 4 }}- {{ index $msg.Meta "db_host" }}:{{ index $msg.Meta "db_port" }}
{{ range $j, $deadlock := json_array "context.deadlocks" $msg.Content }}
{{ range $k, $info := json_array "Sections" $deadlock }}
	{{ format "%v" $info }}
{{ end }}{{ end }}{{ end }}

---

[共 {{ .Group.MessageCount }} 条，查看详细]({{ .PreviewURL }})`,
		Type:        repository.TemplateTypeTemplate,
	},
}

func initPredefinedTemplates(conf *configs.Config, repo repository.TemplateRepo) {
	if !conf.Migrate {
		return
	}

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
