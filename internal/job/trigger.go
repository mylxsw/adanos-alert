package job

import (
	"github.com/mylxsw/adanos-alert/internal/action"
	"github.com/mylxsw/adanos-alert/internal/matcher"
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/container"
	"go.mongodb.org/mongo-driver/bson"
)

const TriggerJobName = "trigger"

type TriggerJob struct {
	app       container.Container
	executing chan interface{} // 标识当前Job是否在执行中
}

func NewTrigger(app container.Container) *TriggerJob {
	return &TriggerJob{app: app, executing: make(chan interface{}, 1)}
}

func (a TriggerJob) Handle() {
	select {
	case a.executing <- struct{}{}:
		defer func() { <-a.executing }()
		a.app.MustResolve(a.processEventGroups)
	default:
		log.Warningf("the last trigger job is not finished yet, skip for this time")
	}
}

func (a TriggerJob) processEventGroups(groupRepo repository.EventGroupRepo, eventRepo repository.EventRepo, ruleRepo repository.RuleRepo, manager action.Manager) error {
	return groupRepo.Traverse(bson.M{"status": repository.EventGroupStatusPending}, func(grp repository.EventGroup) error {
		rule, err := ruleRepo.Get(grp.Rule.ID)
		if err != nil {
			log.WithFields(log.Fields{
				"rule_id": grp.Rule.ID,
				"grp_id":  grp.ID,
				"err":     err.Error(),
			}).Errorf("rule not exist: %w", err)
			return err
		}

		hasError := false
		maxFailedCount := 0
		matchedTriggers := make([]repository.Trigger, 0)
		elseTriggers := make([]repository.Trigger, 0)
		for _, trigger := range rule.Triggers {
			// check whether the trigger has been executed
			for _, act := range grp.Actions {
				if act.ID == trigger.ID && act.Status == repository.TriggerStatusOK {
					continue
				}
			}

			if trigger.IsElseTrigger {
				elseTriggers = append(elseTriggers, trigger)
				continue
			}

			tm, err := matcher.NewTriggerMatcher(trigger)
			if err != nil {
				log.WithFields(log.Fields{
					"trigger_id": trigger.ID,
					"rule_id":    rule.ID,
					"grp_id":     grp.ID,
				}).Errorf("create matcher failed: %w", err)
				continue
			}

			matched, err := tm.Match(matcher.NewTriggerContext(a.app, trigger, grp, func() []repository.Event {
				messages, err := eventRepo.Find(bson.M{"group_ids": grp.ID})
				if err != nil {
					log.WithFields(log.Fields{
						"err": err.Error(),
						"grp": grp,
					}).Errorf("trigger callback: fetch messages from group failed: %v", err)
				}

				return messages
			}))
			if err != nil {
				continue
			}

			if matched {
				hasError, matchedTriggers, maxFailedCount = a.matchedTriggerAction(
					grp,
					manager,
					trigger,
					rule,
					matchedTriggers,
					maxFailedCount,
				)
			}
		}

		// 所有非 ElseTrigger 都没有匹配，执行 ElseTrigger
		if len(matchedTriggers) == 0 && len(elseTriggers) > 0 {
			for _, trigger := range elseTriggers {
				hasError, matchedTriggers, maxFailedCount = a.matchedTriggerAction(
					grp,
					manager,
					trigger,
					rule,
					matchedTriggers,
					maxFailedCount,
				)
			}
		}

		if hasError {
			// if trigger failed count > 3, then set message group failed
			if maxFailedCount > 3 {
				grp.Status = repository.EventGroupStatusFailed
			}
		} else {
			grp.Status = repository.EventGroupStatusOK
		}

		if log.DebugEnabled() {
			log.WithFields(log.Fields{
				"grp_id": grp.ID,
				"status": grp.Status,
			}).Debug("change group status for matchedTriggers")
		}

		grp.Actions = mergeActions(grp.Actions, matchedTriggers)
		return groupRepo.UpdateID(grp.ID, grp)
	})
}

func (a TriggerJob) matchedTriggerAction(grp repository.EventGroup, manager action.Manager, trigger repository.Trigger, rule repository.Rule, matchedTriggers []repository.Trigger, maxFailedCount int) (bool, []repository.Trigger, int) {
	hasError := false
	if err := manager.Dispatch(trigger.Action).Handle(rule, trigger, grp); err != nil {
		trigger.Status = repository.TriggerStatusFailed
		trigger.FailedCount = trigger.FailedCount + 1
		trigger.FailedReason = err.Error()
		hasError = true
	} else {
		trigger.Status = repository.TriggerStatusOK
	}

	matchedTriggers = append(matchedTriggers, trigger)
	if trigger.FailedCount > maxFailedCount {
		maxFailedCount = trigger.FailedCount
	}

	if log.DebugEnabled() {
		log.WithFields(log.Fields{
			"trigger_id": trigger.ID,
			"status":     trigger.Status,
			"grp_id":     grp.ID,
		}).Debug("change trigger status")
	}

	return hasError, matchedTriggers, maxFailedCount
}

func mergeActions(actions []repository.Trigger, triggers []repository.Trigger) []repository.Trigger {
	newActions := make([]repository.Trigger, 0)
	for _, tr := range triggers {
		existed := false
		for i, act := range actions {
			if tr.ID == act.ID {
				actions[i] = tr
				existed = true
				break
			}
		}

		if existed {
			break
		}

		newActions = append(newActions, tr)
	}
	actions = append(actions, newActions...)
	return actions
}
