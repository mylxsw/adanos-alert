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
		Name:        "判断Meta在某个范围内",
		Description: "日志级别为 ERROR 或 FATAL",
		Content:     `Upper(Meta["log_level"]) in ["ERROR", "FATAL"]`,
		Type:        repository.TemplateTypeMatchRule,
	},
	{
		Name:        "判断Meta不在某个范围内",
		Description: "日志级别为非 DEBUG、INFO",
		Content:     `Meta["log_level"] not in ["DEBUG", "INFO"]`,
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
		Name:        "判断字段匹配任一正则表达式",
		Description: "判断 Meta[message.stack_trace] 匹配任一正则表达式",
		Content: `any([
	"^org\\.apache\\.catalina\\.connector\\.ClientAbortException:",
	"^java\\.lang\\.IllegalArgumentException: Invalid salt version"
], {Meta["message.stack_trace"] matches #})`,
		Type: repository.TemplateTypeMatchRule,
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
		Name:        "判断分组中 Events 数量是否大于某个值",
		Description: "当前分组中有超过 10 条 Events",
		Content:     `EventsCount() > 10`,
		Type:        repository.TemplateTypeTriggerRule,
	},
	{
		Name:        "判断分组聚合条件值是否为某些值",
		Description: "匹配聚合条件值为 BigData 的消息",
		Content:     `Group.AggregateKey in ["BigData"]`,
		Type:        repository.TemplateTypeTriggerRule,
	},
	{
		Name:        "报警信息摘要",
		Description: "展示报警信息列表",
		Content: `## {{ .Rule.Name }}

{{ range $i, $evt := .Events 4 }}- 来源：**{{ $evt.Origin }}**，标签：{{ $evt.Tags  }}
{{ cutoff 400 $evt.Content | ident "    > " }}
{{ end }}

---

[共 {{ .Group.MessageCount }} 条，查看详细]({{ .PreviewURL }})`,
		Type: repository.TemplateTypeTemplate,
	},
	{
		Name:        "报警信息摘要(Meta 信息)",
		Description: "显示报警摘要，输出匹配前缀的 Meta 信息",
		Content: `{{ range $i, $evt := .Events 4 }}- 文件：{{ index $evt.Meta "log.file.path" }}
{{ meta_prefix_filter $evt.Meta "message" | serialize | cutoff 400 | ident "    > "}}
{{ end }}`,
		Type: repository.TemplateTypeTemplate,
	},
	{
		Name:        "报警详情链接",
		Description: "报警详细信息链接地址",
		Content:     `[共 {{ .Group.MessageCount }} 条，查看详细]({{ .PreviewURL }})`,
		Type:        repository.TemplateTypeTemplate,
	},
	{
		Name:        "报警详情链接(报告)",
		Description: "报警详细信息链接地址，报告模式",
		Content:     `[共 {{ .Group.MessageCount }} 条，查看详细]({{ .ReportURL }})`,
		Type:        repository.TemplateTypeTemplate,
	},
	{
		Name:        "服务健康检查失败告警示例",
		Description: "服务健康检查失败告警示例",
		Content: `## {{ .Rule.Name }}

{{ range $i, $evt := .Events 4 }}- 服务 {{ json_get "service" "-" $evt.Content }} 的实例 {{ json_get "instance" "-" $evt.Content }} 已经持续 <font color="red">{{ json_get "lost_duration" "-" $evt.Content | printf "%ss" | helpers.HumanDuration }}</font> 处于离线状态，原因为健康检查失败
{{ end }}

{{- if gt .Group.MessageCount 4 }}
---
[共 {{ .Group.MessageCount }} 条，查看详细]({{ .ReportURL }})
{{- end }}`,
		Type: repository.TemplateTypeTemplate,
	},
	{
		Name:        "嵌入全局的规则模板",
		Description: "在动作模板中引用规则的展示模板内容",
		Content:     `{{ .RuleTemplateParsed }}`,
		Type:        repository.TemplateTypeTemplate,
	},
}

func initPredefinedTemplates(conf *configs.Config, repo repository.TemplateRepo) {
	if !conf.Migrate && !conf.ReMigrate {
		return
	}

	// 删除所有预定义模板
	if conf.ReMigrate {
		if err := repo.Delete(bson.M{"predefined": true}); err != nil {
			log.Errorf("delete all predefined templates failed: %v", err)
			return
		}
	}

	// 更新预定义模板
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

			if log.DebugEnabled() {
				log.WithFields(log.Fields{
					"temp": t,
				}).Debugf("add predefined template %s with id %s", t.Name, id.Hex())
			}
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
					}).Errorf("update predefined template failed: %v", err)
				} else {
					log.WithFields(log.Fields{
						"temp": t,
					}).Warningf("update predefined template %s", t.Name)
				}
			}
		}
	}
}
