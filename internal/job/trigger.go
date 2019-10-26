package job

import (
	"github.com/mylxsw/adanos-alert/internal/action"
	"github.com/mylxsw/adanos-alert/internal/matcher"
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/container"
	"go.mongodb.org/mongo-driver/bson"
)

type TriggerJob struct {
	app *container.Container
}

func NewTrigger(app *container.Container) *TriggerJob {
	return &TriggerJob{app: app}
}

func (a TriggerJob) Handle() {
	log.Debug("trigger actions...")

	a.app.MustResolve(a.processMessageGroups)
}

func (a TriggerJob) processMessageGroups(groupRepo repository.MessageGroupRepo, ruleRepo repository.RuleRepo, manager *action.Manager) error {
	return groupRepo.Traverse(bson.M{"status": repository.MessageGroupStatusPending}, func(grp repository.MessageGroup) error {
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
		triggers := make([]repository.Trigger, 0)
		for _, trigger := range rule.Triggers {
			// check whether the trigger has been executed
			for _, act := range grp.Actions {
				if act.ID == trigger.ID && act.Status == repository.TriggerStatusOK {
					continue
				}
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

			matched, err := tm.Match(matcher.TriggerContext{Group: grp})
			if err != nil {
				log.WithFields(log.Fields{
					"trigger_id": trigger.ID,
					"rule_id":    rule.ID,
					"grp_id":     grp.ID,
				}).Errorf("trigger matcher match failed: %s", err)
				continue
			}

			if matched {
				if err := manager.Dispatch(trigger.Action).Handle(trigger, grp); err != nil {
					trigger.Status = repository.TriggerStatusFailed
					trigger.FailedCount = trigger.FailedCount + 1
					trigger.FailedReason = err.Error()
					hasError = true
				} else {
					trigger.Status = repository.TriggerStatusOK
				}

				triggers = append(triggers, trigger)
				if trigger.FailedCount > maxFailedCount {
					maxFailedCount = trigger.FailedCount
				}
			}
		}

		if hasError {
			// if trigger failed count > 3, then set message group failed
			if maxFailedCount > 3 {
				grp.Status = repository.MessageGroupStatusFailed
			}
		} else {
			grp.Status = repository.MessageGroupStatusOK
		}

		grp.Actions = mergeActions(grp.Actions, triggers)
		return groupRepo.Update(grp.ID, grp)
	})
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
