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
	"github.com/mylxsw/adanos-alert/internal/extension"
	"github.com/mylxsw/adanos-alert/internal/matcher"
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/adanos-alert/internal/template"
	"github.com/mylxsw/adanos-alert/pkg/misc"
	"github.com/mylxsw/adanos-alert/pubsub"
	"github.com/mylxsw/glacier/event"
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/glacier/web"
	"github.com/mylxsw/go-utils/str"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RuleController struct {
	cc infra.Resolver
}

func NewRuleController(cc infra.Resolver) web.Controller {
	return &RuleController{cc: cc}
}

func (r RuleController) Register(router web.Router) {
	router.Group("/rules/", func(router web.Router) {
		router.Post("/", r.Add).Name("rules:add")
		router.Get("/", r.Rules).Name("rules:all")
		router.Get("/{id}/", r.Rule).Name("rules:one")
		router.Post("/{id}/", r.Update).Name("rules:update")
		router.Delete("/{id}/", r.Delete).Name("rules:delete")
	})

	router.Group("/rules-meta/", func(router web.Router) {
		router.Get("/tags/", r.Tags).Name("rules:meta:tags")
		router.Get("/message-sample/", r.MessageSample).Name("rules:meta:message-sample")
	})

	router.Group("/rules-test/", func(router web.Router) {
		router.Post("/rule-check/{type}/", r.Check).Name("rules:test:check")
	})

	router.Group("/evaluate/", func(router web.Router) {
		router.Post("/expression-sample/", r.EvaluateExpressionSample).Name("evaluate:sample")
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
	UserEvalFunc  string   `json:"user_eval_func"`
}

// RuleForm is a form object using create or update rule
type RuleForm struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`

	AggregateRule string `json:"aggregate_rule"`
	RelationRule  string `json:"relation_rule"`

	ReadyType        string                 `json:"ready_type"`
	Interval         int64                  `json:"interval"`
	DailyTimes       []string               `json:"daily_times"`
	TimeRanges       []repository.TimeRange `json:"time_ranges"`
	RealtimeInterval int64                  `json:"realtime_interval"`

	Rule             string            `json:"rule"`
	IgnoreRule       string            `json:"ignore_rule"`
	IgnoreMaxCount   int               `json:"ignore_max_count"`
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

	if r.RealtimeInterval < 0 {
		return errors.New("realtime_interval must be a positive number in minutes")
	}

	if r.Status != "" && !govalidator.IsIn(r.Status, string(repository.RuleStatusEnabled), string(repository.RuleStatusDisabled)) {
		return errors.New("status is invalid, must be enabled/disabled")
	}

	_, err := matcher.NewEventMatcher(repository.Rule{Rule: r.Rule, IgnoreRule: r.IgnoreRule})
	if err != nil {
		return fmt.Errorf("rule is invalid: %w", err)
	}

	if r.IgnoreMaxCount < 0 {
		return fmt.Errorf("ignore_max_count must greater than or equal 0")
	}

	for i, tr := range r.Triggers {
		if tr.PreCondition != "" {
			_, err := matcher.NewTriggerMatcher(tr.PreCondition, repository.Trigger{
				PreCondition: tr.PreCondition,
			}, true)
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

		if tr.UserEvalFunc != "" {
			_, err := matcher.NewTriggerMatcher(tr.UserEvalFunc, repository.Trigger{}, false)
			if err != nil {
				return fmt.Errorf("trigger #%d's user eval func is invalid: %w", i, err)
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

	if _, err := matcher.NewEventFinger(r.AggregateRule); err != nil {
		return fmt.Errorf("aggregate rule is invalid")
	}

	if _, err := matcher.NewEventFinger(r.RelationRule); err != nil {
		return fmt.Errorf("relation rule is invalid")
	}

	return nil
}

// Check validate the rule
func (r RuleController) Check(ctx web.Context, conf *configs.Config, msgRepo repository.EventRepo) web.Response {
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
			_, err = matcher.NewEventMatcher(repository.Rule{Rule: "true", IgnoreRule: content})
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
					misc.IfElse(matched, "与当前 Event 匹配", "与当前 Event 不匹配"),
					misc.IfElse(matched && ignored, "，但是该 Event 被忽略", ""),
				),
			})
		} else {
			_, err = matcher.NewEventMatcher(repository.Rule{Rule: content})
		}
	case repository.TemplateTypeTriggerRule:
		_, err = matcher.NewTriggerMatcher(content, repository.Trigger{PreCondition: content}, true)
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
		finger, err1 := matcher.NewEventFinger(content)
		if err1 == nil {
			if msgID != "" {
				msg, err1 := r.getEventByID(msgID, msgRepo)
				if err1 == nil {
					res, err1 := finger.Run(msg)
					if err1 == nil {
						return ctx.JSON(web.M{
							"error": nil,
							"msg":   fmt.Sprintf("当前 Event 应用规则后的返回值为 %s", res),
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
	case "user_eval_rule":
		_, err = matcher.NewTriggerMatcher(content, repository.Trigger{PreCondition: content}, false)
	}

	if err != nil {
		return ctx.JSON(web.M{"error": err.Error(), "msg": ""})
	}

	return ctx.JSON(web.M{
		"error": nil,
		"msg":   "",
	})
}

func createPayloadForTemplateCheck(r RuleController, conf *configs.Config, msgID string, msgRepo repository.EventRepo, content string) *action.Payload {
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
	messagesQuerier := func(groupID primitive.ObjectID, limit int64) []repository.Event {
		msg, err := r.getEventByID(msgID, msgRepo)
		if err != nil {
			return []repository.Event{}
		}

		messages := make([]repository.Event, 0)
		for i := int64(0); i < limit; i++ {
			messages = append(messages, msg)
		}

		return messages
	}

	payload := action.CreatePayload(
		conf,
		messagesQuerier,
		"dingding",
		rule,
		triggers[0],
		repository.EventGroup{
			ID:           primitive.NewObjectID(),
			SeqNum:       1000,
			Type:         repository.EventTypePlain,
			MessageCount: 30,
			AggregateKey: "AggregateKey",
			Rule:         rule.ToGroupRule("", repository.EventTypePlain),
			Actions:      triggers,
			Status:       repository.EventGroupStatusOK,
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
		for _, u := range str.Distinct(t.UserRefs) {
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
		DailyTimes:       str.Distinct(ruleForm.DailyTimes),
		Interval:         ruleForm.Interval,
		TimeRanges:       ruleForm.TimeRanges,
		RealtimeInterval: ruleForm.RealtimeInterval,
		Rule:             ruleForm.Rule,
		IgnoreRule:       ruleForm.IgnoreRule,
		IgnoreMaxCount:   ruleForm.IgnoreMaxCount,
		AggregateRule:    ruleForm.AggregateRule,
		RelationRule:     ruleForm.RelationRule,
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
		for _, u := range str.Distinct(t.UserRefs) {
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
			UserEvalFunc:  t.UserEvalFunc,
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
		DailyTimes:       str.Distinct(ruleForm.DailyTimes),
		Interval:         ruleForm.Interval,
		TimeRanges:       ruleForm.TimeRanges,
		RealtimeInterval: ruleForm.RealtimeInterval,
		Rule:             ruleForm.Rule,
		IgnoreRule:       ruleForm.IgnoreRule,
		IgnoreMaxCount:   ruleForm.IgnoreMaxCount,
		AggregateRule:    ruleForm.AggregateRule,
		RelationRule:     ruleForm.RelationRule,
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
	idStr := ctx.Input("id")
	if idStr != "" {
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			return nil, web.WrapJSONError(err, http.StatusUnprocessableEntity)
		}

		filter["_id"] = id
	}

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

func (r RuleController) getEventByID(messageID string, msgRepo repository.EventRepo) (repository.Event, error) {
	msgID, err := primitive.ObjectIDFromHex(messageID)
	if err != nil {
		return repository.Event{}, fmt.Errorf("invalid message_id: %v", err)
	}

	return msgRepo.Get(msgID)
}

// testMessageMatchRule test if the message and rule can be matched
func (r RuleController) testMessageMatchRule(rule string, messageID string, msgRepo repository.EventRepo) (matched bool, ignored bool, err error) {
	message, err := r.getEventByID(messageID, msgRepo)
	if err != nil {
		return false, false, err
	}

	m, err := matcher.NewEventMatcher(repository.Rule{Rule: rule})
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
func (r RuleController) testMessageIgnoreRule(rule string, messageID string, msgRepo repository.EventRepo) (bool, error) {
	message, err := r.getEventByID(messageID, msgRepo)
	if err != nil {
		return false, err
	}

	m, err := matcher.NewEventMatcher(repository.Rule{Rule: "true", IgnoreRule: rule})
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
func (r RuleController) MessageSample(ctx web.Context, groupRepo repository.EventGroupRepo, msgRepo repository.EventRepo) (*repository.Event, error) {
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

func (r RuleController) evaluateEvent(ctx web.Context, evt repository.Event, content string) web.Response {
	eventFinger, err := matcher.NewEventFinger(content)
	if err != nil {
		return ctx.JSONError(err.Error(), http.StatusUnprocessableEntity)
	}

	res, err := eventFinger.Run(evt)
	if err != nil {
		return ctx.JSONError(err.Error(), 422)
	}

	return ctx.JSON(bson.M{
		"res": res,
	})
}

type EvalSampleReq struct {
	EventSample extension.CommonEvent `json:"event_sample"`
	EventID     string                `json:"event_id"`
	Expression  string                `json:"expression"`
}

func (r RuleController) EvaluateExpressionSample(ctx web.Context) web.Response {
	var evalSample EvalSampleReq
	if err := ctx.Unmarshal(&evalSample); err != nil {
		return ctx.JSONError(fmt.Sprintf("invalid request: %v", err), http.StatusUnprocessableEntity)
	}

	return r.evaluateEvent(ctx, evalSample.EventSample.CreateRepoEvent(), evalSample.Expression)
}
