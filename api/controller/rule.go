package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/mylxsw/adanos-alert/internal/action"
	"github.com/mylxsw/adanos-alert/internal/matcher"
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/adanos-alert/pkg/array"
	"github.com/mylxsw/adanos-alert/pkg/template"
	"github.com/mylxsw/container"
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
	})

	router.Group("/rules-test/", func(router *web.Router) {
		router.Post("/rule-message/", r.TestMessageMatch).Name("rules:test:rule-message")
		router.Post("/rule-check/{type}/", r.Check).Name("rules:test:check")
	})
}

// RuleTriggerForm is a form object using to hold a trigger
type RuleTriggerForm struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	PreCondition string   `json:"pre_condition"`
	Action       string   `json:"action"`
	Meta         string   `json:"meta"`
	UserRefs     []string `json:"user_refs"`
}

// RuleForm is a form object using create or update rule
type RuleForm struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`

	ReadyType  string                 `json:"ready_type"`
	Interval   int64                  `json:"interval"`
	DailyTimes []string               `json:"daily_times"`
	TimeRanges []repository.TimeRange `json:"time_ranges"`

	Rule            string            `json:"rule"`
	Template        string            `json:"template"`
	SummaryTemplate string            `json:"summary_template"`
	Triggers        []RuleTriggerForm `json:"triggers"`

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

	_, err := matcher.NewMessageMatcher(repository.Rule{Rule: r.Rule})
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

	return nil
}

// Check validate the rule
func (r RuleController) Check(ctx web.Context) web.Response {
	content := ctx.Input("content")

	var err error
	switch repository.TemplateType(ctx.PathVar("type")) {
	case repository.TemplateTypeMatchRule:
		_, err = matcher.NewMessageMatcher(repository.Rule{Rule: content})
	case repository.TemplateTypeTriggerRule:
		_, err = matcher.NewTriggerMatcher(repository.Trigger{PreCondition: content})
	case repository.TemplateTypeTemplate:
		_, err = template.CreateParser(content)
	}

	if err != nil {
		return ctx.JSON(web.M{"error": err.Error()})
	}

	return ctx.JSON(web.M{
		"error": nil,
	})
}

// Add create a new rule
func (r RuleController) Add(ctx web.Context, repo repository.RuleRepo, manager action.Manager) (*repository.Rule, error) {
	var ruleForm RuleForm
	if err := ctx.Unmarshal(&ruleForm); err != nil {
		return nil, web.WrapJSONError(err, http.StatusUnprocessableEntity)
	}

	ruleForm.actionManager = manager
	ctx.Validate(ruleForm, true)

	triggers := make([]repository.Trigger, 0)
	for _, t := range ruleForm.Triggers {

		users := make([]primitive.ObjectID, 0)
		for _, u := range array.StringUnique(t.UserRefs) {
			uid, err := primitive.ObjectIDFromHex(u)
			if err == nil {
				users = append(users, uid)
			}
		}

		triggers = append(triggers, repository.Trigger{
			Name:         t.Name,
			PreCondition: t.PreCondition,
			Action:       t.Action,
			Meta:         t.Meta,
			UserRefs:     users,
		})
	}

	ruleID, err := repo.Add(repository.Rule{
		Name:            ruleForm.Name,
		Description:     ruleForm.Description,
		Tags:            ruleForm.Tags,
		ReadyType:       ruleForm.ReadyType,
		DailyTimes:      array.StringUnique(ruleForm.DailyTimes),
		Interval:        ruleForm.Interval,
		TimeRanges:      ruleForm.TimeRanges,
		Rule:            ruleForm.Rule,
		Template:        ruleForm.Template,
		SummaryTemplate: ruleForm.SummaryTemplate,
		Triggers:        triggers,
		Status:          repository.RuleStatus(ruleForm.Status),
	})
	if err != nil {
		return nil, web.WrapJSONError(err, http.StatusInternalServerError)
	}

	rule, err := repo.Get(ruleID)
	if err != nil {
		return nil, web.WrapJSONError(err, http.StatusInternalServerError)
	}

	return &rule, nil
}

// Update replace one rule for specified id
func (r RuleController) Update(ctx web.Context, ruleRepo repository.RuleRepo, manager action.Manager) (*repository.Rule, error) {
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
		for _, u := range array.StringUnique(t.UserRefs) {
			uid, err := primitive.ObjectIDFromHex(u)
			if err == nil {
				users = append(users, uid)
			}
		}

		triggerID, _ := primitive.ObjectIDFromHex(t.ID)
		triggers = append(triggers, repository.Trigger{
			ID:           triggerID,
			Name:         t.Name,
			PreCondition: t.PreCondition,
			Action:       t.Action,
			Meta:         t.Meta,
			UserRefs:     users,
		})
	}

	if err := ruleRepo.UpdateID(id, repository.Rule{
		ID:              original.ID,
		Name:            ruleForm.Name,
		Description:     ruleForm.Description,
		Tags:            ruleForm.Tags,
		ReadyType:       ruleForm.ReadyType,
		DailyTimes:      array.StringUnique(ruleForm.DailyTimes),
		Interval:        ruleForm.Interval,
		TimeRanges:      ruleForm.TimeRanges,
		Rule:            ruleForm.Rule,
		Template:        ruleForm.Template,
		SummaryTemplate: ruleForm.SummaryTemplate,
		Triggers:        triggers,
		Status:          repository.RuleStatus(ruleForm.Status),
		CreatedAt:       original.CreatedAt,
		UpdatedAt:       original.CreatedAt,
	}); err != nil {
		return nil, web.WrapJSONError(err, http.StatusInternalServerError)
	}

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
		filter["name"] = name
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
func (r RuleController) Delete(ctx web.Context, repo repository.RuleRepo) error {
	id, err := primitive.ObjectIDFromHex(ctx.PathVar("id"))
	if err != nil {
		return web.WrapJSONError(err, http.StatusUnprocessableEntity)
	}

	return repo.DeleteID(id)
}

// TestMessageMatch test if the message and rule can be matched
func (r RuleController) TestMessageMatch(ctx web.Context) web.Response {
	rule := ctx.Input("rule")
	message := ctx.Input("message")

	var msg repository.Message
	if err := json.Unmarshal([]byte(message), &msg); err != nil {
		return ctx.JSONError(fmt.Sprintf("invalid message: %v", err), http.StatusUnprocessableEntity)
	}

	m, err := matcher.NewMessageMatcher(repository.Rule{Rule: rule})
	if err != nil {
		return ctx.JSONError(fmt.Sprintf("invalid rule: %v", err), http.StatusUnprocessableEntity)
	}

	rs, err := m.Match(msg)
	if err != nil {
		return ctx.JSONError(fmt.Sprintf("rule match with errors: %v", err), http.StatusUnprocessableEntity)
	}

	return ctx.JSON(bson.M{
		"matched": rs,
	})
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
