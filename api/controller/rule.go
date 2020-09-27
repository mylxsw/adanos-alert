package controller

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/mylxsw/adanos-alert/configs"
	"github.com/mylxsw/adanos-alert/internal/action"
	"github.com/mylxsw/adanos-alert/internal/matcher"
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/adanos-alert/internal/template"
	"github.com/mylxsw/adanos-alert/pkg/misc"
	"github.com/mylxsw/adanos-alert/pkg/strarr"
	"github.com/mylxsw/adanos-alert/pubsub"
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier/event"
	"github.com/mylxsw/glacier/web"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RuleController struct {
	cc container.Container
}

func NewRuleController(cc container.Container) web.Controller {
	return &RuleController{cc: cc}
}

func (r RuleController) Register(router *web.Router) {
	router.Group("/rules/", func(router *web.Router) {
		router.Post("/", r.Add).Name("rules:add")
		router.Get("/", r.Rules).Name("rules:all")
		router.Get("/{id}/", r.Rule).Name("rules:one")
		router.Post("/{id}/", r.Update).Name("rules:update")
		router.Delete("/{id}/", r.Delete).Name("rules:delete")
	})

	router.Group("/rules-meta/", func(router *web.Router) {
		router.Get("/tags/", r.Tags).Name("rules:meta:tags")
		router.Get("/message-sample/", r.MessageSample).Name("rules:meta:message-sample")
	})

	router.Group("/rules-test/", func(router *web.Router) {
		router.Post("/rule-check/{type}/", r.Check).Name("rules:test:check")
	})
}

// RuleTriggerForm is a form object using to hold a trigger
type RuleTriggerForm struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	IsElseTrigger bool     `json:"is_else_trigger"`
	PreCondition  string   `json:"pre_condition"`
	Action        string   `json:"action"`
	Meta          string   `json:"meta"`
	UserRefs      []string `json:"user_refs"`
}

// RuleForm is a form object using create or update rule
type RuleForm struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`

	AggregateRule string `json:"aggregate_rule"`

	ReadyType  string                 `json:"ready_type"`
	Interval   int64                  `json:"interval"`
	DailyTimes []string               `json:"daily_times"`
	TimeRanges []repository.TimeRange `json:"time_ranges"`

	Rule             string            `json:"rule"`
	IgnoreRule       string            `json:"ignore_rule"`
	Template         string            `json:"template"`
	SummaryTemplate  string            `json:"summary_template"`
	ReportTemplateID string            `json:"report_template_id"`
	Triggers         []RuleTriggerForm `json:"triggers"`

	Status string `json:"status"`

	actionManager action.Manager
}

// Validate implement web.Validator interface
func (r RuleForm) Validate(req web.Request) error {
	if r.Name == "" {
		return errors.New("name is required")
	}

	if r.ReadyType == "" {
		r.ReadyType = repository.ReadyTypeInterval
	}

	switch r.ReadyType {
	case repository.ReadyTypeInterval:
		if !govalidator.InRangeInt(r.Interval, 60, 3600*24) {
			return errors.New("interval is invalid, must between 1min~24h")
		}
	case repository.ReadyTypeDailyTime:
		if len(r.DailyTimes) == 0 {
			return fmt.Errorf("daily_times is required")
		}

		for _, dailyTime := range r.DailyTimes {
			if len(dailyTime) < 5 {
				return fmt.Errorf("invalid daily_time format for %s", dailyTime)
			}

			_, err := time.Parse("15:04", dailyTime[:5])
			if err != nil {
				return fmt.Errorf("invalid daily_time format for %s: %v", dailyTime, err)
			}
		}
	case repository.ReadyTypeTimeRange:
		if len(r.TimeRanges) == 0 {
			return errors.New("invalid time ranges, must not be empty")
		}

		for _, t := range r.TimeRanges {
			if len(t.StartTime) < 5 || len(t.EndTime) < 5 {
				return fmt.Errorf("invalid time format for time range %s-%s", t.StartTime, t.EndTime)
			}

			if _, err := time.Parse("15:04", t.StartTime[:5]); err != nil {
				return fmt.Errorf("invalid startTime time in time range for %s-%s: %v", t.StartTime, t.EndTime, err)
			}

			if _, err := time.Parse("15:04", t.EndTime[:5]); err != nil {
				return fmt.Errorf("invalid endTime time in time range for %s-%s: %v", t.StartTime, t.EndTime, err)
			}

			if !govalidator.InRangeInt(t.Interval, 60, 3600*24) {
				return fmt.Errorf("invalid interval %d in time range for %s-%s", t.Interval, t.StartTime, t.EndTime)
			}

			// 搜索开始时间是否合法
			startTimeOk := false
			for _, tt := range r.TimeRanges {
				if t.StartTime == tt.EndTime {
					startTimeOk = true
					break
				}
			}

			// 搜索截止时间是否合法
			endTimeOk := false
			for _, tt := range r.TimeRanges {
				if t.EndTime == tt.StartTime {
					endTimeOk = true
					break
				}
			}

			if !startTimeOk || !endTimeOk {
				return errors.New("time range is not complete")
			}
		}
	default:
		return errors.New("invalid readyType")
	}

	if r.Status != "" && !govalidator.IsIn(r.Status, string(repository.RuleStatusEnabled), string(repository.RuleStatusDisabled)) {
		return errors.New("status is invalid, must be enabled/disabled")
	}

	_, err := matcher.NewMessageMatcher(repository.Rule{Rule: r.Rule, IgnoreRule: r.IgnoreRule})
	if err != nil {
		return fmt.Errorf("rule is invalid: %w", err)
	}

	for i, tr := range r.Triggers {
		if tr.PreCondition != "" {
			_, err := matcher.NewTriggerMatcher(repository.Trigger{
				PreCondition: tr.PreCondition,
			})
			if err != nil {
				return fmt.Errorf("trigger #%d is invalid: %w", i, err)
			}
		}

		for j, u := range tr.UserRefs {
			_, err := primitive.ObjectIDFromHex(u)
			if err != nil {
				return fmt.Errorf("trigger #%d, user #%d with value %s: %w", i, j, u, err)
			}
		}

		act := r.actionManager.Run(tr.Action)
		if act == nil {
			return fmt.Errorf("trigger #%d, action [%s] is not support", i, tr.Action)
		}

		if err := act.Validate(tr.Meta, tr.UserRefs); err != nil {
			return fmt.Errorf("trigger #%d, action [%s] with invalid meta: %w", i, tr.Action, err)
		}
	}

	if _, err := matcher.NewMessageFinger(r.AggregateRule); err != nil {
		return fmt.Errorf("group rule is invalid")
	}

	return nil
}

// Check validate the rule
func (r RuleController) Check(ctx web.Context, conf *configs.Config, msgRepo repository.MessageRepo) web.Response {
	content := ctx.Input("content")
	msgID := ctx.Input("msg_id")

	var err error
	switch repository.TemplateType(ctx.PathVar("type")) {
	case "match_rule_ignore":
		if msgID != "" {
			ignored, err := r.testMessageIgnoreRule(content, msgID, msgRepo)
			if err != nil {
				return ctx.JSON(web.M{
					"error": err,
					"msg":   "",
				})
			}

			return ctx.JSON(web.M{
				"error": nil,
				"msg":   misc.IfElse(ignored, "该消息被忽略", "该消息不会被忽略"),
			})
		} else {
			_, err = matcher.NewMessageMatcher(repository.Rule{Rule: "true", IgnoreRule: content})
		}
	case repository.TemplateTypeMatchRule:
		if msgID != "" {
			matched, ignored, err := r.testMessageMatchRule(content, msgID, msgRepo)
			if err != nil {
				return ctx.JSON(web.M{
					"error": err,
					"msg":   "",
				})
			}

			return ctx.JSON(web.M{
				"error": nil,
				"msg": fmt.Sprintf(
					"%s%s",
					misc.IfElse(matched, "与当前 message 匹配", "与当前 message 不匹配"),
					misc.IfElse(matched && ignored, "，但是该消息被忽略", ""),
				),
			})
		} else {
			_, err = matcher.NewMessageMatcher(repository.Rule{Rule: content})
		}
	case repository.TemplateTypeTriggerRule:
		_, err = matcher.NewTriggerMatcher(repository.Trigger{PreCondition: content})
	case repository.TemplateTypeTemplate:
		data, err1 := template.Parse(r.cc, content, createPayloadForTemplateCheck(r, conf, msgID, msgRepo, content))
		if err1 != nil {
			err = err1
		} else {
			return ctx.JSON(web.M{
				"error": nil,
				"msg":   data,
			})
		}
	case "aggregate_rule":
		finger, err1 := matcher.NewMessageFinger(content)
		if err1 == nil {
			if msgID != "" {
				msg, err1 := r.getMessageByID(msgID, msgRepo)
				if err1 == nil {
					res, err1 := finger.Run(msg)
					if err1 == nil {
						return ctx.JSON(web.M{
							"error": nil,
							"msg":   fmt.Sprintf("当前 message 聚合 Key 为 %s", res),
						})
					} else {
						err = err1
					}
				} else {
					err = err1
				}
			}
		} else {
			err = err1
		}
	}

	if err != nil {
		return ctx.JSON(web.M{"error": err.Error(), "msg": ""})
	}

	return ctx.JSON(web.M{
		"error": nil,
		"msg":   "",
	})
}

func createPayloadForTemplateCheck(r RuleController, conf *configs.Config, msgID string, msgRepo repository.MessageRepo, content string) *action.Payload {
	triggers := []repository.Trigger{
		{
			ID:           primitive.NewObjectID(),
			Name:         "测试触发规则",
			PreCondition: "",
			Action:       "dingding",
			UserRefs:     []primitive.ObjectID{primitive.NewObjectID()},
			Status:       repository.TriggerStatusOK,
		},
	}
	rule := repository.Rule{
		ID:          primitive.NewObjectID(),
		Name:        "测试报警规则",
		Description: "测试报警规则描述",
		Tags:        []string{"test-tag"},
		ReadyType:   repository.ReadyTypeInterval,
		Interval:    60,
		Template:    content,
		Triggers:    triggers,
		Status:      repository.RuleStatusEnabled,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	messagesQuerier := func(groupID primitive.ObjectID, limit int64) []repository.Message {
		msg, err := r.getMessageByID(msgID, msgRepo)
		if err != nil {
			return []repository.Message{}
		}

		messages := make([]repository.Message, 0)
		for i := int64(0); i < limit; i++ {
			messages = append(messages, msg)
		}

		return messages
	}

	payload := action.BuildPayload(
		conf,
		messagesQuerier,
		"dingding",
		rule,
		triggers[0],
		repository.MessageGroup{
			ID:           primitive.NewObjectID(),
			SeqNum:       1000,
			Type:         repository.MessageTypePlain,
			MessageCount: 3,
			AggregateKey: "AggregateKey",
			Rule:         rule.ToGroupRule("", repository.MessageTypePlain),
			Actions:      triggers,
			Status:       repository.MessageGroupStatusOK,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	)
	return payload
}

// Add create a new rule
func (r RuleController) Add(ctx web.Context, repo repository.RuleRepo, em event.Manager, manager action.Manager) (*repository.Rule, error) {
	var ruleForm RuleForm
	if err := ctx.Unmarshal(&ruleForm); err != nil {
		return nil, web.WrapJSONError(err, http.StatusUnprocessableEntity)
	}

	ruleForm.actionManager = manager
	ctx.Validate(ruleForm, true)

	triggers := make([]repository.Trigger, 0)
	for _, t := range ruleForm.Triggers {

		users := make([]primitive.ObjectID, 0)
		for _, u := range strarr.Distinct(t.UserRefs) {
			uid, err := primitive.ObjectIDFromHex(u)
			if err == nil {
				users = append(users, uid)
			}
		}

		triggers = append(triggers, repository.Trigger{
			Name:          t.Name,
			PreCondition:  t.PreCondition,
			Action:        t.Action,
			Meta:          t.Meta,
			IsElseTrigger: t.IsElseTrigger,
			UserRefs:      users,
		})
	}

	reportTempID, err := primitive.ObjectIDFromHex(ruleForm.ReportTemplateID)
	if err != nil {
		reportTempID = primitive.NilObjectID
	}

	newRule := repository.Rule{
		Name:             ruleForm.Name,
		Description:      ruleForm.Description,
		Tags:             ruleForm.Tags,
		ReadyType:        ruleForm.ReadyType,
		DailyTimes:       strarr.Distinct(ruleForm.DailyTimes),
		Interval:         ruleForm.Interval,
		TimeRanges:       ruleForm.TimeRanges,
		Rule:             ruleForm.Rule,
		IgnoreRule:       ruleForm.IgnoreRule,
		AggregateRule:    ruleForm.AggregateRule,
		Template:         ruleForm.Template,
		SummaryTemplate:  ruleForm.SummaryTemplate,
		ReportTemplateID: reportTempID,
		Triggers:         triggers,
		Status:           repository.RuleStatus(ruleForm.Status),
	}

	ruleID, err := repo.Add(newRule)
	if err != nil {
		return nil, web.WrapJSONError(err, http.StatusInternalServerError)
	}

	em.Publish(pubsub.RuleChangedEvent{
		Rule:      newRule,
		Type:      pubsub.EventTypeAdd,
		CreatedAt: time.Now(),
	})

	rule, err := repo.Get(ruleID)
	if err != nil {
		return nil, web.WrapJSONError(err, http.StatusInternalServerError)
	}

	return &rule, nil
}

// Update replace one rule for specified id
func (r RuleController) Update(ctx web.Context, ruleRepo repository.RuleRepo, em event.Manager, manager action.Manager) (*repository.Rule, error) {
	id, err := primitive.ObjectIDFromHex(ctx.PathVar("id"))
	if err != nil {
		return nil, web.WrapJSONError(err, http.StatusUnprocessableEntity)
	}

	var ruleForm RuleForm
	if err := ctx.Unmarshal(&ruleForm); err != nil {
		return nil, web.WrapJSONError(err, http.StatusUnprocessableEntity)
	}

	ruleForm.actionManager = manager
	ctx.Validate(ruleForm, true)

	original, err := ruleRepo.Get(id)
	if err != nil {
		return nil, web.WrapJSONError(err, http.StatusInternalServerError)
	}

	triggers := make([]repository.Trigger, 0)
	for _, t := range ruleForm.Triggers {
		users := make([]primitive.ObjectID, 0)
		for _, u := range strarr.Distinct(t.UserRefs) {
			uid, err := primitive.ObjectIDFromHex(u)
			if err == nil {
				users = append(users, uid)
			}
		}

		triggerID, _ := primitive.ObjectIDFromHex(t.ID)
		triggers = append(triggers, repository.Trigger{
			ID:            triggerID,
			Name:          t.Name,
			PreCondition:  t.PreCondition,
			Action:        t.Action,
			Meta:          t.Meta,
			IsElseTrigger: t.IsElseTrigger,
			UserRefs:      users,
		})
	}

	reportTempID, err := primitive.ObjectIDFromHex(ruleForm.ReportTemplateID)
	if err != nil {
		reportTempID = primitive.NilObjectID
	}

	newRule := repository.Rule{
		ID:               original.ID,
		Name:             ruleForm.Name,
		Description:      ruleForm.Description,
		Tags:             ruleForm.Tags,
		ReadyType:        ruleForm.ReadyType,
		DailyTimes:       strarr.Distinct(ruleForm.DailyTimes),
		Interval:         ruleForm.Interval,
		TimeRanges:       ruleForm.TimeRanges,
		Rule:             ruleForm.Rule,
		IgnoreRule:       ruleForm.IgnoreRule,
		AggregateRule:    ruleForm.AggregateRule,
		Template:         ruleForm.Template,
		SummaryTemplate:  ruleForm.SummaryTemplate,
		ReportTemplateID: reportTempID,
		Triggers:         triggers,
		Status:           repository.RuleStatus(ruleForm.Status),
		CreatedAt:        original.CreatedAt,
		UpdatedAt:        original.CreatedAt,
	}

	if err := ruleRepo.UpdateID(id, newRule); err != nil {
		return nil, web.WrapJSONError(err, http.StatusInternalServerError)
	}

	em.Publish(pubsub.RuleChangedEvent{
		Rule:      newRule,
		Type:      pubsub.EventTypeUpdate,
		CreatedAt: time.Now(),
	})

	rule, err := ruleRepo.Get(id)
	if err != nil {
		return nil, web.WrapJSONError(err, http.StatusInternalServerError)
	}

	return &rule, nil
}

type RulesResp struct {
	Rules []repository.Rule `json:"rules"`
	Users map[string]string `json:"users"`

	Next   int64      `json:"next"`
	Search RuleSearch `json:"search"`
}

type RuleSearch struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	UserID string `json:"user_id"`
}

// Rules return all rules
func (r RuleController) Rules(ctx web.Context, ruleRepo repository.RuleRepo, userRepo repository.UserRepo) (*RulesResp, error) {
	filter := bson.M{}
	name := ctx.Input("name")
	if name != "" {
		filter["name"] = bson.M{"$regex": name}
	}

	status := ctx.Input("status")
	if status != "" {
		filter["status"] = status
	}

	tag := ctx.Input("tag")
	if tag != "" {
		filter["tags"] = tag
	}

	userIDStr := ctx.Input("user_id")
	if userIDStr != "" {
		userID, err := primitive.ObjectIDFromHex(userIDStr)
		if err != nil {
			return nil, web.WrapJSONError(fmt.Errorf("invalid argument: user_id"), http.StatusUnprocessableEntity)
		}

		filter["triggers.user_refs"] = userID
	}

	dingID := ctx.Input("dingding_id")
	if dingID != "" {
		filter["triggers.meta"] = bson.M{"$regex": fmt.Sprintf(`"robot_id":"%s"`, dingID)}
	}

	offset, limit := offsetAndLimit(ctx)
	rules, next, err := ruleRepo.Paginate(filter, offset, limit)
	if err != nil {
		return nil, web.WrapJSONError(err, http.StatusInternalServerError)
	}

	userIDs := make([]primitive.ObjectID, 0)
	for _, rule := range rules {
		for _, act := range rule.Triggers {
			userIDs = append(userIDs, act.UserRefs...)
		}
	}

	users, _ := userRepo.Find(bson.M{"_id": bson.M{"$in": userIDs}})
	userRefs := make(map[string]string)
	for _, u := range users {
		userRefs[u.ID.Hex()] = u.Name
	}

	return &RulesResp{
		Rules: rules,
		Users: userRefs,
		Next:  next,
		Search: RuleSearch{
			Name:   name,
			Status: status,
			UserID: userIDStr,
		},
	}, nil
}

// Rule return one rule
func (r RuleController) Rule(ctx web.Context, ruleRepo repository.RuleRepo) (*repository.Rule, error) {
	id, err := primitive.ObjectIDFromHex(ctx.PathVar("id"))
	if err != nil {
		return nil, web.WrapJSONError(err, http.StatusUnprocessableEntity)
	}

	rule, err := ruleRepo.Get(id)
	if err != nil {
		return nil, web.WrapJSONError(err, http.StatusInternalServerError)
	}

	return &rule, nil
}

// Delete delete a rule
func (r RuleController) Delete(ctx web.Context, em event.Manager, repo repository.RuleRepo) error {
	id, err := primitive.ObjectIDFromHex(ctx.PathVar("id"))
	if err != nil {
		return web.WrapJSONError(err, http.StatusUnprocessableEntity)
	}

	rule, err := repo.Get(id)
	if err != nil {
		return err
	}

	em.Publish(pubsub.RuleChangedEvent{
		Rule:      rule,
		Type:      pubsub.EventTypeDelete,
		CreatedAt: time.Now(),
	})

	return repo.DeleteID(id)
}

func (r RuleController) getMessageByID(messageID string, msgRepo repository.MessageRepo) (repository.Message, error) {
	msgID, err := primitive.ObjectIDFromHex(messageID)
	if err != nil {
		return repository.Message{}, fmt.Errorf("invalid message_id: %v", err)
	}

	return msgRepo.Get(msgID)
}

// testMessageMatchRule test if the message and rule can be matched
func (r RuleController) testMessageMatchRule(rule string, messageID string, msgRepo repository.MessageRepo) (matched bool, ignored bool, err error) {
	message, err := r.getMessageByID(messageID, msgRepo)
	if err != nil {
		return false, false, err
	}

	m, err := matcher.NewMessageMatcher(repository.Rule{Rule: rule})
	if err != nil {
		return false, false, fmt.Errorf("invalid rule: %v", err)
	}

	rs, ignored, err := m.Match(message)
	if err != nil {
		return false, false, fmt.Errorf("rule match with errors: %v", err)
	}

	return rs, ignored, nil
}

// testMessageIgnoreRule test if the message and rule can be ignored
func (r RuleController) testMessageIgnoreRule(rule string, messageID string, msgRepo repository.MessageRepo) (bool, error) {
	message, err := r.getMessageByID(messageID, msgRepo)
	if err != nil {
		return false, err
	}

	m, err := matcher.NewMessageMatcher(repository.Rule{Rule: "true", IgnoreRule: rule})
	if err != nil {
		return false, fmt.Errorf("invalid rule: %v", err)
	}

	_, ignored, err := m.Match(message)
	if err != nil {
		return false, fmt.Errorf("rule match with errors: %v", err)
	}

	return ignored, nil
}

// Tags return all tags existed
func (r RuleController) Tags(ctx web.Context, repo repository.RuleRepo) web.Response {
	timeoutCtx, _ := context.WithTimeout(ctx.Context(), 5*time.Second)
	tags, err := repo.Tags(timeoutCtx)
	if err != nil {
		return ctx.JSONError(fmt.Sprintf("query failed: %v", err), http.StatusInternalServerError)
	}

	return ctx.JSON(web.M{
		"tags": tags,
	})
}

// MessageSample 根据规则id查询一个匹配的消息样本
func (r RuleController) MessageSample(ctx web.Context, groupRepo repository.MessageGroupRepo, msgRepo repository.MessageRepo) (*repository.Message, error) {
	id, err := primitive.ObjectIDFromHex(ctx.Input("id"))
	if err != nil {
		return nil, web.WrapJSONError(err, http.StatusUnprocessableEntity)
	}

	grp, err := groupRepo.LastGroup(bson.M{"rule._id": id})
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, nil
		}

		return nil, fmt.Errorf("query related group failed: %v", err)
	}

	messages, _, err := msgRepo.Paginate(bson.M{"group_ids": grp.ID}, 0, 1)
	if err != nil {
		return nil, fmt.Errorf("query related message failed: %v", err)
	}

	if len(messages) == 0 {
		return nil, nil
	}

	return &messages[0], nil
}
