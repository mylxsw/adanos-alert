package job

import (
	"time"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/adanos-alert/internal/rule"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/go-toolkit/collection"
	"github.com/mylxsw/go-toolkit/container"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AggregationJob struct {
	app *container.Container
}

func NewAggregationJob(app *container.Container) *AggregationJob {
	return &AggregationJob{app: app}
}

// Handle do two things:
// 1. message grouping, delivery all ungrouped messages to message group
// 2. change the message groups that satisfied the conditions to pending status
func (a *AggregationJob) Handle() {
	log.Debug("aggregating messages...")
	a.app.MustResolve(func(msgRepo repository.MessageRepo, ruleRepo repository.RuleRepo, groupRepo repository.MessageGroupRepo) {
		// get all rules
		rules, err := ruleRepo.Find(bson.M{"status": repository.RuleStatusEnabled})
		if err != nil {
			log.Errorf("aggregate message failed because rules query failed: %s", err)
			return
		}

		// create matchers from rules
		var matchers []*rule.MessageMatcher
		if err := collection.MustNew(rules).Map(func(ru repository.Rule) *rule.MessageMatcher {
			matcher, err := rule.NewMessageMatcher(ru)
			if err != nil {
				log.Errorf("invalid rule: %s", err)
			}

			return matcher
		}).All(&matchers); err != nil {
			log.Errorf("create message matchers failed: %s", err)
			return
		}

		// traverse all ungrouped messages to group
		if err := a.groupingMessages(msgRepo, groupRepo, matchers); err != nil {
			log.Errorf("aggregate message failed: %s", err)
			return
		}

		// change message group status to pending when it reach the aggregate condition
		if err := a.pendingMessageGroup(groupRepo, msgRepo); err != nil {
			log.Errorf("message group failed: %s", err)
			return
		}
	})
}

func (a *AggregationJob) groupingMessages(msgRepo repository.MessageRepo, groupRepo repository.MessageGroupRepo, matchers []*rule.MessageMatcher) error {
	collectingGroups := make(map[primitive.ObjectID]repository.MessageGroup)
	return msgRepo.Traverse(bson.M{"status": repository.MessageStatusPending}, func(msg repository.Message) error {
		for _, m := range matchers {
			matched, err := m.Match(msg)
			if err != nil {
				return err
			}

			// if the message matched a rule, update message's group_id and skip to next message
			if matched {
				if _, ok := collectingGroups[m.Rule().ID]; !ok {
					grp, err := groupRepo.CollectingGroup(m.Rule().ToGroupRule())
					if err != nil {
						log.Errorf("create collecting group failed: %s", err)
						return err
					}

					collectingGroups[m.Rule().ID] = grp
				}

				msg.GroupID = collectingGroups[m.Rule().ID].ID
				msg.Status = repository.MessageStatusGrouped
				return msgRepo.UpdateID(msg.ID, msg)
			}
		}

		msg.Status = repository.MessageStatusCanceled
		return msgRepo.UpdateID(msg.ID, msg)
	})
}

func (a *AggregationJob) pendingMessageGroup(groupRepo repository.MessageGroupRepo, msgRepo repository.MessageRepo) error {
	return groupRepo.Traverse(bson.M{"status": repository.MessageGroupStatusCollecting}, func(grp repository.MessageGroup) error {
		if grp.CreatedAt.Add(time.Duration(grp.Rule.Interval) * time.Second).After(time.Now()) {
			return nil
		}

		msgCount, err := msgRepo.Count(bson.M{"group_id": grp.ID})
		if err != nil {
			return err
		}

		grp.MessageCount = msgCount
		grp.Status = repository.MessageGroupStatusPending
		return groupRepo.Update(grp.ID, grp)
	})
}
