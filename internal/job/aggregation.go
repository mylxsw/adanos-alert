package job

import (
	"fmt"

	"github.com/mylxsw/adanos-alert/internal/matcher"
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/coll"
	"github.com/mylxsw/container"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const AggregationJobName = "aggregation"

type AggregationJob struct {
	app container.Container
}

func NewAggregationJob(app container.Container) *AggregationJob {
	return &AggregationJob{app: app}
}

// Handle do two things:
// 1. message grouping, delivery all ungrouped messages to message group
// 2. change the message groups that satisfied the conditions to pending status
func (a *AggregationJob) Handle() {
	// traverse all ungrouped messages to group
	a.app.MustResolve(a.groupingMessages)
	// change message group status to pending when it reach the aggregate condition
	a.app.MustResolve(a.pendingMessageGroup)
}

func (a *AggregationJob) groupingMessages(msgRepo repository.MessageRepo, groupRepo repository.MessageGroupRepo, ruleRepo repository.RuleRepo) error {
	matchers, err := a.initializeMatchers(ruleRepo)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	collectingGroups := make(map[primitive.ObjectID]repository.MessageGroup)
	err = msgRepo.Traverse(bson.M{"status": repository.MessageStatusPending}, func(msg repository.Message) error {
		for _, m := range matchers {
			matched, err := m.Match(msg)
			if err != nil {
				continue
			}

			// if the message matched a rule, update message's group_id and skip to next message
			if matched {
				if _, ok := collectingGroups[m.Rule().ID]; !ok {
					grp, err := groupRepo.CollectingGroup(m.Rule().ToGroupRule())
					if err != nil {
						log.WithFields(log.Fields{
							"msg":  msg,
							"rule": m.Rule(),
							"err":  err.Error(),
						}).Errorf("create collecting group failed: %v", err)
						return err
					}

					collectingGroups[m.Rule().ID] = grp
				}

				msg.GroupID = append(msg.GroupID, collectingGroups[m.Rule().ID].ID)
				msg.Status = repository.MessageStatusGrouped
			}
		}

		// if message not match any rules, set message as canceled
		if msg.Status != repository.MessageStatusGrouped {
			msg.Status = repository.MessageStatusCanceled
		}

		log.WithFields(log.Fields{
			"msg_id": msg.ID.Hex(),
			"status": msg.Status,
		}).Debug("change message status")

		return msgRepo.UpdateID(msg.ID, msg)
	})
	if err != nil {
		return err
	}

	// 将能够与规则匹配的 Canceled 的 message 转换为 Expired
	return msgRepo.Traverse(bson.M{"status": repository.MessageStatusCanceled}, func(msg repository.Message) error {
		for _, m := range matchers {
			matched, err := m.Match(msg)
			if err != nil {
				continue
			}

			// if the message matched a rule, update message's group_id and skip to next message
			if matched {
				msg.Status = repository.MessageStatusExpired
				break
			}
		}

		return msgRepo.UpdateID(msg.ID, msg)
	})
}

func (a *AggregationJob) initializeMatchers(ruleRepo repository.RuleRepo) ([]*matcher.MessageMatcher, error) {
	// get all rules
	rules, err := ruleRepo.Find(bson.M{"status": repository.RuleStatusEnabled})
	if err != nil {
		return nil, fmt.Errorf("aggregate message failed because rules query failed: %s", err)
	}

	// create matchers from rules
	var matchers []*matcher.MessageMatcher
	if err := coll.MustNew(rules).Map(func(ru repository.Rule) *matcher.MessageMatcher {
		mat, err := matcher.NewMessageMatcher(ru)
		if err != nil {
			log.Errorf("invalid rule: %v", err)
		}

		return mat
	}).All(&matchers); err != nil {
		return nil, fmt.Errorf("create message matchers failed: %s", err)
	}

	return matchers, nil
}

func (a *AggregationJob) pendingMessageGroup(groupRepo repository.MessageGroupRepo, msgRepo repository.MessageRepo) error {
	return groupRepo.Traverse(bson.M{"status": repository.MessageGroupStatusCollecting}, func(grp repository.MessageGroup) error {
		if !grp.Ready() {
			return nil
		}

		msgCount, err := msgRepo.Count(bson.M{"group_ids": grp.ID})
		if err != nil {
			log.WithFields(log.Fields{
				"grp": grp,
				"err": err,
			}).Errorf("query message count failed: %v", err)
		}

		grp.MessageCount = msgCount
		grp.Status = repository.MessageGroupStatusPending

		log.WithFields(log.Fields{
			"grp_id": grp.ID.Hex(),
			"status": grp.Status,
		}).Debug("change group status")

		return groupRepo.UpdateID(grp.ID, grp)
	})
}
