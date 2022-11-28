package job

import (
	"context"
	"fmt"
	"time"

	"github.com/mylxsw/adanos-alert/internal/matcher"
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/adanos-alert/pkg/misc"
	"github.com/mylxsw/adanos-alert/pubsub"
	"github.com/mylxsw/adanos-alert/service"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/coll"
	"github.com/mylxsw/glacier/event"
	"github.com/mylxsw/glacier/infra"
	"go.mongodb.org/mongo-driver/bson"
)

const AggregationJobName = "aggregation"

type AggregationJob struct {
	app       infra.Resolver
	executing chan interface{} // 标识当前Job是否在执行中
}

func NewAggregationJob(app infra.Resolver) *AggregationJob {
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

func (a *AggregationJob) groupingEvents(eventRepo repository.EventRepo, evtRelRepo repository.EventRelationRepo, groupRepo repository.EventGroupRepo, ruleRepo repository.RuleRepo, groupSrv service.EventGroupService) error {
	matchers, err := initializeMatchers(ruleRepo)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	collectingGroups := make(map[string]repository.EventGroup)
	err = eventRepo.Traverse(bson.M{"status": repository.EventStatusPending}, func(evt repository.Event) error {
		messageCanIgnore := false
		for _, m := range matchers {
			matched, ignored, err := m.Match(evt)
			if err != nil {
				continue
			}

			rule := m.Rule()

			// if the message matched a rule, update message's group_id and skip to next message
			if matched {
				// 对于匹配规则的消息，首先判断是否能够为消息建立关联
				if rule.RelationRule != "" {
					if relationSummary := BuildEventFinger(rule.RelationRule, evt); relationSummary != "" {
						if evtRel, err := evtRelRepo.AddOrUpdateEventRelation(context.TODO(), relationSummary, rule.ID); err != nil {
							log.WithFields(log.Fields{
								"evt":  evt,
								"rule": rule,
								"err":  err,
							}).Errorf("create event relation failed: %v", err)
						} else {
							evt.RelationID = append(evt.RelationID, evtRel.ID)
						}
					}
				}

				// 为消息分组
				if ignored {
					messageCanIgnore = true
					evt.Type = repository.EventTypeIgnored
				}

				aggregateKey := BuildEventFinger(rule.AggregateRule, evt)
				key := fmt.Sprintf("%s:%s:%s", rule.ID.Hex(), aggregateKey, evt.Type)
				if _, ok := collectingGroups[key]; !ok {

					evtGrpRule := rule.ToGroupRule(aggregateKey, evt.Type)
					// 检查规则是否支持即时发送，当支持即时发送时，事件组的就绪时间将会被设置为立即执行
					if rule.RealtimeInterval > 0 {
						groupSrv.EventShouldRealtime(context.Background(), key, time.Duration(rule.RealtimeInterval)*time.Minute, func() {
							evtGrpRule.ExpectReadyAt = time.Now().Add(-time.Second)
							evtGrpRule.Realtime = true
						})
					}

					grp, err := groupRepo.CollectingGroup(evtGrpRule)
					if err != nil {
						log.WithFields(log.Fields{
							"evt":  evt,
							"rule": rule,
							"err":  err.Error(),
						}).Errorf("create collecting group failed: %v", err)
						return err
					}

					collectingGroups[key] = grp
				}

				evt.GroupID = append(evt.GroupID, collectingGroups[key].ID)
				evt.Status = repository.EventStatusGrouped
			}
		}

		// messageCanIgnore 和 message 状态 变换规则
		// true  | pending  -> ignore
		// false | pending  -> canceled
		// true  | grouped  -> grouped
		// false | grouped  -> grouped

		// if message not match any rules, set message as canceled
		if evt.Status == repository.EventStatusPending {
			evt.Status = misc.IfElse(messageCanIgnore,
				repository.EventStatusIgnored,
				repository.EventStatusCanceled,
			).(repository.EventStatus)
		}

		if log.DebugEnabled() {
			log.WithFields(log.Fields{
				"evt_id": evt.ID.Hex(),
				"status": evt.Status,
			}).Debug("change message status")
		}

		return eventRepo.UpdateID(evt.ID, evt)
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

func (a *AggregationJob) pendingEventGroup(groupRepo repository.EventGroupRepo, evtRepo repository.EventRepo, em event.Manager) error {
	return groupRepo.Traverse(bson.M{"status": repository.EventGroupStatusCollecting}, func(grp repository.EventGroup) error {
		if !grp.Ready() {
			return nil
		}

		evtCount, err := evtRepo.Count(bson.M{"group_ids": grp.ID})
		if err != nil {
			log.WithFields(log.Fields{
				"grp": grp,
				"err": err,
			}).Errorf("query message count failed: %v", err)
		}

		// 1. 当前分组事件数为空，则忽略当前告警，设置状态为已取消
		// 2. 当前分组事件数不为空，需要分情况:
		//     1）事件组为 ignored 类型，当前事件数 <= 忽略最大值，设置为取消，不进行告警
		//     2）事件组为 ignored 类型，当前事件数 > 忽略最大值，直接进行告警
		//     3）事件组为非 ignored 类型，直接进行告警
		if evtCount == 0 {
			grp.Status = repository.EventGroupStatusCanceled
		} else {
			grp.Status = repository.EventGroupStatusPending
			if grp.Type == repository.EventTypeIgnored {
				if grp.Rule.IgnoreMaxCount == 0 || int(evtCount) <= grp.Rule.IgnoreMaxCount {
					grp.Status = repository.EventGroupStatusCanceled
				}
			}
		}

		// 被忽略的事件组，如果超过最大忽略阈值，加入到告警行列中，将分组的类型设置为 IngoredExceed
		if grp.Type == repository.EventTypeIgnored && grp.Status != repository.EventGroupStatusCanceled {
			grp.Type = repository.EventTypeIgnoredExceed
		}

		grp.MessageCount = evtCount

		if log.DebugEnabled() {
			log.WithFields(log.Fields{
				"grp_id": grp.ID.Hex(),
				"status": grp.Status,
			}).Debug("change group status")
		}

		err = groupRepo.UpdateID(grp.ID, grp)

		if evtCount > 0 {
			em.Publish(pubsub.MessageGroupPendingEvent{
				Group:     grp,
				CreatedAt: time.Now(),
			})
		}

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
