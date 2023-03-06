package matcher

import (
	"fmt"
	"sync"
	"time"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/adanos-alert/internal/template"
	"github.com/mylxsw/adanos-alert/pkg/helper"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/glacier/infra"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TriggerMatcher is a matcher for trigger
type TriggerMatcher struct {
	program *vm.Program
	trigger repository.Trigger
}

type TriggerContext struct {
	helper.Helpers
	Group   repository.EventGroup
	Trigger repository.Trigger

	EventCallback     func() []repository.Event
	eventCallbackOnce sync.Once
	events            []repository.Event

	resolver infra.Resolver
}

// NewTriggerContext create a new TriggerContext
func NewTriggerContext(resolver infra.Resolver, trigger repository.Trigger, group repository.EventGroup, eventCallback func() []repository.Event) *TriggerContext {
	return &TriggerContext{resolver: resolver, Trigger: trigger, Group: group, EventCallback: eventCallback}
}

// Messages return all events in group
// This method is depressed
func (tc *TriggerContext) Messages() []repository.Event {
	return tc.Events()
}

// Events return all events in group
func (tc *TriggerContext) Events() []repository.Event {
	tc.eventCallbackOnce.Do(func() {
		if tc.EventCallback != nil {
			tc.events = tc.EventCallback()
		}
	})

	return tc.events
}

// MessagesCount return the count in group
// This method is depressed
func (tc *TriggerContext) MessagesCount() int64 {
	return tc.EventsCount()
}

// EventsCount return the count in group
func (tc *TriggerContext) EventsCount() int64 {
	var count int64 = 0
	tc.resolver.MustResolve(func(msgRepo repository.EventRepo) {
		count, _ = msgRepo.Count(bson.M{
			"group_ids": tc.Group.ID,
		})
	})

	return count
}

// MessagesMatchRegexCount get the count for events matched regex
// This method is depressed
func (tc *TriggerContext) MessagesMatchRegexCount(regex string) int64 {
	return tc.EventsMatchRegexCount(regex)
}

// EventsMatchRegexCount get the count for events matched regex
func (tc *TriggerContext) EventsMatchRegexCount(regex string) int64 {
	var count int64 = 0
	tc.resolver.MustResolve(func(msgRepo repository.EventRepo) {
		filter := bson.M{
			"group_ids": tc.Group.ID,
			"content":   primitive.Regex{Pattern: regex, Options: "im"},
		}

		count, _ = msgRepo.Count(filter)
	})

	return count
}

// HasEventsMatchRegexs check if events matched regexs
func (tc *TriggerContext) HasEventsMatchRegexs(regexs []string) bool {
	var matched bool
	tc.resolver.MustResolve(func(msgRepo repository.EventRepo) {
		regexps := make(bson.A, 0)
		for _, reg := range regexs {
			regexps = append(regexps, primitive.Regex{Pattern: reg, Options: "im"})
		}

		filter := bson.M{
			"group_ids": tc.Group.ID,
			"content":   bson.M{"$in": regexps},
		}

		var err error
		matched, err = msgRepo.Has(filter)
		if err != nil {
			log.WithFields(log.Fields{
				"filter": filter,
			}).Warningf("triggerContext#HasEventsMatchRegexs failed: %v", err)
		}
	})

	return matched
}

// HasEventsMatchRegex check if events matched regex
func (tc *TriggerContext) HasEventsMatchRegex(regex string) bool {
	var matched bool
	tc.resolver.MustResolve(func(msgRepo repository.EventRepo) {
		filter := bson.M{
			"group_ids": tc.Group.ID,
			"content":   primitive.Regex{Pattern: regex, Options: "im"},
		}

		var err error
		matched, err = msgRepo.Has(filter)
		if err != nil {
			log.WithFields(log.Fields{
				"filter": filter,
			}).Warningf("triggerContext#HasEventsMatchRegex failed: %v", err)
		}
	})

	return matched
}

// MessagesWithTagsCount get the count for events which has tags
// This method is depressed
func (tc *TriggerContext) MessagesWithTagsCount(tags string) int64 {
	return tc.EventsWithTagsCount(tags)
}

// EventsWithTagsCount get the count for events which has tags
func (tc *TriggerContext) EventsWithTagsCount(tags string) int64 {
	var count int64 = 0
	tc.resolver.MustResolve(func(msgRepo repository.EventRepo) {
		filter := bson.M{
			"group_ids": tc.Group.ID,
			"tags":      bson.M{"$in": template.StringTags(tags, ",")},
		}

		count, _ = msgRepo.Count(filter)
	})

	return count
}

// MessagesWithMetaCount get the count for events has a meta[key] equals to value
// This method is depressed
func (tc *TriggerContext) MessagesWithMetaCount(key, value string) int64 {
	return tc.EventsWithMetaCount(key, value)
}

// EventsWithMetaCount get the count for events has a meta[key] equals to value
func (tc *TriggerContext) EventsWithMetaCount(key, value string) int64 {
	var count int64 = 0
	tc.resolver.MustResolve(func(msgRepo repository.EventRepo) {
		filter := bson.M{
			"group_ids":   tc.Group.ID,
			"meta." + key: value,
		}

		count, _ = msgRepo.Count(filter)
	})

	return count
}

func (tc *TriggerContext) UsersHasPropertyRegex(key, valueRegex, returnField string) []string {
	users := make([]string, 0)
	tc.resolver.MustResolve(func(userRepo repository.UserRepo) {
		var err error
		users, err = userRepo.GetUserMetasRegex(key, valueRegex, returnField)
		if err != nil {
			log.WithFields(log.Fields{
				"key":         key,
				"value":       valueRegex,
				"returnField": returnField,
			}).Warningf("triggerContext#UsersHasPropertyRegex failed: %v", err)
		}
	})
	return users
}

func (tc *TriggerContext) UsersHasProperty(key, value, returnField string) []string {
	users := make([]string, 0)
	tc.resolver.MustResolve(func(userRepo repository.UserRepo) {
		var err error
		users, err = userRepo.GetUserMetas(key, value, returnField)
		if err != nil {
			log.WithFields(log.Fields{
				"key":         key,
				"value":       value,
				"returnField": returnField,
			}).Warningf("triggerContext#UsersHasProperty failed: %v", err)
		}
	})
	return users
}

func (tc *TriggerContext) UsersIDWithPropertyRegex(key, valueRegex, returnField string) []repository.UserIDWithMeta {
	users := make([]repository.UserIDWithMeta, 0)
	tc.resolver.MustResolve(func(userRepo repository.UserRepo) {
		var err error
		users, err = userRepo.GetUserIDWithMetasRegex(key, valueRegex, returnField)
		if err != nil {
			log.WithFields(log.Fields{
				"key":         key,
				"value":       valueRegex,
				"returnField": returnField,
			}).Warningf("triggerContext#UsersHasPropertyRegex failed: %v", err)
		}
	})
	return users
}

func (tc *TriggerContext) UsersIDWithProperty(key, value, returnField string) []repository.UserIDWithMeta {
	users := make([]repository.UserIDWithMeta, 0)
	tc.resolver.MustResolve(func(userRepo repository.UserRepo) {
		var err error
		users, err = userRepo.GetUserIDWithMetas(key, value, returnField)
		if err != nil {
			log.WithFields(log.Fields{
				"key":         key,
				"value":       value,
				"returnField": returnField,
			}).Warningf("triggerContext#UsersHasProperty failed: %v", err)
		}
	})
	return users
}

// TriggeredTimesInPeriod return triggered times in specified periods
func (tc *TriggerContext) TriggeredTimesInPeriod(periodInMinutes int, triggerStatus string) int64 {
	var triggeredTimes int64 = 0
	tc.resolver.MustResolve(func(groupRepo repository.EventGroupRepo) {
		filter := bson.M{
			"actions._id": tc.Trigger.ID,
			"updated_at":  bson.M{"$gt": time.Now().Add(-time.Duration(periodInMinutes) * time.Minute)},
		}

		if triggerStatus != "" {
			filter["actions.trigger_status"] = triggerStatus
		}

		n, _ := groupRepo.Count(filter)

		triggeredTimes = n
	})

	if log.DebugEnabled() {
		log.WithFields(log.Fields{
			"periodInMinutes": periodInMinutes,
			"triggerStatus":   triggerStatus,
			"times":           triggeredTimes,
			"trigger":         tc.Trigger,
		}).Debugf("helper function: triggeredTimesInPeriod")
	}

	return triggeredTimes
}

// LastTriggeredGroup get last triggeredGroup
func (tc *TriggerContext) LastTriggeredGroup(triggerStatus string) repository.EventGroup {
	var lastTriggeredGroup repository.EventGroup
	tc.resolver.MustResolve(func(groupRepo repository.EventGroupRepo) {
		filter := bson.M{
			"actions._id": tc.Trigger.ID,
		}

		if triggerStatus != "" {
			filter["actions.trigger_status"] = triggerStatus
		}

		grp, err := groupRepo.LastGroup(filter)
		if err == nil {
			lastTriggeredGroup = grp
		}
	})

	if log.DebugEnabled() {
		log.WithFields(log.Fields{
			"group":         lastTriggeredGroup,
			"trigger":       tc.Trigger,
			"triggerStatus": triggerStatus,
		}).Debugf("helper function: lastTriggeredGroup")
	}

	return lastTriggeredGroup
}

// NewTriggerMatcher create a new TriggerMatcher
// https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md
func NewTriggerMatcher(condition string, trigger repository.Trigger, boolRet bool) (*TriggerMatcher, error) {
	if condition == "" && boolRet {
		condition = "true"
	}

	options := make([]expr.Option, 0)
	options = append(options, expr.Env(&TriggerContext{}))
	if boolRet {
		options = append(options, expr.AsBool())
	}

	program, err := expr.Compile(condition, options...)
	if err != nil {
		return nil, err
	}

	return &TriggerMatcher{program: program, trigger: trigger}, nil
}

// Match check whether the msg is match with the rule
func (m *TriggerMatcher) Match(triggerCtx *TriggerContext) (bool, error) {
	rs, err := expr.Run(m.program, triggerCtx)
	if err != nil {
		return false, err
	}

	if boolRes, ok := rs.(bool); ok {
		return boolRes, nil
	}

	return false, InvalidReturnVal
}

// Eval 根据指定的表达式创建 Event 解析内容为字符串数组
func (m *TriggerMatcher) Eval(triggerCtx *TriggerContext) ([]string, error) {
	result, err := expr.Run(m.program, triggerCtx)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	if s, ok := result.(string); ok {
		return []string{s}, nil
	}

	if ss, ok := result.([]string); ok {
		return ss, nil
	}

	if ss, ok := result.([]interface{}); ok {
		res := make([]string, 0)
		for _, s := range ss {
			res = append(res, fmt.Sprintf("%v", s))
		}

		return res, nil
	}

	return []string{fmt.Sprintf("%v", result)}, nil
}

// Trigger return original trigger object
func (m *TriggerMatcher) Trigger() repository.Trigger {
	return m.trigger
}
