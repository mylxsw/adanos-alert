package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

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
	Name        string `json:"name"`
	Description string `json:"description"`

	Interval int64 `json:"interval"`

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

	if !govalidator.InRangeInt(r.Interval, 60, 3600*24) {
		return errors.New("interval is invalid, must between 1min~24h")
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

		if err := act.Validate(tr.Meta); err != nil {
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

	return ctx.JSON(web.M{
		"error": err,
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
		Interval:        ruleForm.Interval,
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
		Interval:        ruleForm.Interval,
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

	userIDStr := ctx.Input("user_id")
	if userIDStr != "" {
		userID, err := primitive.ObjectIDFromHex(userIDStr)
		if err != nil {
			return nil, web.WrapJSONError(fmt.Errorf("invalid argument: user_id"), http.StatusUnprocessableEntity)
		}

		filter["triggers.user_refs"] = userID
	}

	rules, err := ruleRepo.Find(filter)
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
