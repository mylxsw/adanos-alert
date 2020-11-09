package job

import (
	"fmt"
	"time"

	"github.com/mylxsw/adanos-alert/internal/matcher"
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/adanos-alert/pkg/misc"
	"github.com/mylxsw/adanos-alert/pubsub"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/coll"
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier/event"
	"go.mongodb.org/mongo-driver/bson"
)

const AggregationJobName = "aggregation"

type AggregationJob struct {
	app       container.Container
	executing chan interface{} // 标识当前Job是否在执行中
}

func NewAggregationJob(app container.Container) *AggregationJob {
	return &AggregationJob{app: app, executing: make(chan interface{}, 1)}
}

// Handle do two things:
// 1. message grouping, delivery all ungrouped messages to message group
// 2. change the message groups that satisfied the conditions to pending status
func (a *AggregationJob) Handle() {
	select {
	case a.executing <- struct{}{}:
		defer func() { <-a.executing }()
		// traverse all ungrouped events to group
		a.app.MustResolve(a.groupingEvents)
		// change event group status to pending when it reach the aggregate condition
		a.app.MustResolve(a.pendingEventGroup)
	default:
		log.Warningf("the last aggregation job is not finished yet, skip for this time")
	}
}

func (a *AggregationJob) groupingEvents(eventRepo repository.EventRepo, groupRepo repository.EventGroupRepo, ruleRepo repository.RuleRepo) error {
	matchers, err := initializeMatchers(ruleRepo)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	collectingGroups := make(map[string]repository.EventGroup)
	err = eventRepo.Traverse(bson.M{"status": repository.EventStatusPending}, func(msg repository.Event) error {
		messageCanIgnore := false
		for _, m := range matchers {
			matched, ignored, err := m.Match(msg)
			if err != nil {
				continue
			}

			// if the message matched a rule, update message's group_id and skip to next message
			if matched {
				if ignored {
					messageCanIgnore = true
				} else {
					aggregateKey := BuildEventFinger(m.Rule().AggregateRule, msg)
					key := fmt.Sprintf("%s:%s:%s", m.Rule().ID.Hex(), aggregateKey, msg.Type)
					if _, ok := collectingGroups[key]; !ok {
						grp, err := groupRepo.CollectingGroup(m.Rule().ToGroupRule(aggregateKey, msg.Type))
						if err != nil {
							log.WithFields(log.Fields{
								"msg":  msg,
								"rule": m.Rule(),
								"err":  err.Error(),
							}).Errorf("create collecting group failed: %v", err)
							return err
						}

						collectingGroups[key] = grp
					}

					msg.GroupID = append(msg.GroupID, collectingGroups[key].ID)
					msg.Status = repository.EventStatusGrouped
				}
			}
		}

		// messageCanIgnore 和 message 状态 变换规则
		// true  | pending  -> ignore
		// false | pending  -> canceled
		// true  | grouped  -> grouped
		// false | grouped  -> grouped

		// if message not match any rules, set message as canceled
		if msg.Status == repository.EventStatusPending {
			msg.Status = misc.IfElse(messageCanIgnore,
				repository.EventStatusIgnored,
				repository.EventStatusCanceled,
			).(repository.EventStatus)
		}

		if log.DebugEnabled() {
			log.WithFields(log.Fields{
				"msg_id": msg.ID.Hex(),
				"status": msg.Status,
			}).Debug("change message status")
		}

		return eventRepo.UpdateID(msg.ID, msg)
	})
	if err != nil {
		return err
	}

	// 将能够与规则匹配的 Canceled 的 message 转换为 Expired
	return eventRepo.Traverse(bson.M{"status": repository.EventStatusCanceled}, func(msg repository.Event) error {
		for _, m := range matchers {
			matched, _, err := m.Match(msg)
			if err != nil {
				continue
			}

			// if the message matched a rule, update message's group_id and skip to next message
			if matched {
				msg.Status = repository.EventStatusExpired
				break
			}
		}

		return eventRepo.UpdateID(msg.ID, msg)
	})
}

func initializeMatchers(ruleRepo repository.RuleRepo) ([]*matcher.EventMatcher, error) {
	// get all rules
	rules, err := ruleRepo.Find(bson.M{"status": repository.RuleStatusEnabled})
	if err != nil {
		return nil, fmt.Errorf("aggregate message failed because rules query failed: %s", err)
	}

	// create matchers from rules
	var matchers []*matcher.EventMatcher
	if err := coll.MustNew(rules).Map(func(ru repository.Rule) *matcher.EventMatcher {
		mat, err := matcher.NewEventMatcher(ru)
		if err != nil {
			log.Errorf("invalid rule: %v", err)
		}

		return mat
	}).All(&matchers); err != nil {
		return nil, fmt.Errorf("create message matchers failed: %s", err)
	}

	return matchers, nil
}

func (a *AggregationJob) pendingEventGroup(groupRepo repository.EventGroupRepo, msgRepo repository.EventRepo, em event.Manager) error {
	return groupRepo.Traverse(bson.M{"status": repository.EventGroupStatusCollecting}, func(grp repository.EventGroup) error {
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
		grp.Status = repository.EventGroupStatusPending

		if log.DebugEnabled() {
			log.WithFields(log.Fields{
				"grp_id": grp.ID.Hex(),
				"status": grp.Status,
			}).Debug("change group status")
		}

		err = groupRepo.UpdateID(grp.ID, grp)

		em.Publish(pubsub.MessageGroupPendingEvent{
			Group:     grp,
			CreatedAt: time.Now(),
		})

		return err
	})
}

func BuildEventFinger(groupRule string, evt repository.Event) string {
	finger, err := matcher.NewEventFinger(groupRule)
	if err != nil {
		log.WithFields(log.Fields{
			"rule": groupRule,
		}).Errorf("parse group rule failed: %v", err)
		return "[error]invalid_rule"
	}
	groupKey, err := finger.Run(evt)
	if err != nil {
		log.WithFields(log.Fields{
			"rule": groupRule,
		}).Errorf("rule group failed: %v", err)
		return "[error]parse_failed"
	}

	return groupKey
}

type MatchedRule struct {
	Rule         repository.Rule `json:"rule"`
	AggregateKey string          `json:"aggregate_key"`
}

// BuildEventMatchTest 创建 event 与规则的匹配测试，用于检测 event 能够匹配哪些规则
func BuildEventMatchTest(ruleRepo repository.RuleRepo) func(msg repository.Event) ([]MatchedRule, error) {
	return func(msg repository.Event) ([]MatchedRule, error) {
		matchedRules := make([]MatchedRule, 0)

		matchers, err := initializeMatchers(ruleRepo)
		if err != nil {
			log.Error(err.Error())
			return matchedRules, err
		}

		for _, m := range matchers {
			matched, _, err := m.Match(msg)
			if err != nil {
				continue
			}

			// if the message matched a rule, update message's group_id and skip to next message
			if matched {
				matchedRules = append(matchedRules, MatchedRule{
					Rule:         m.Rule(),
					AggregateKey: BuildEventFinger(m.Rule().AggregateRule, msg),
				})
			}
		}

		return matchedRules, nil
	}
}
