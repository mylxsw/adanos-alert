package job

import (
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
		matchers, _ := collection.MustNew(rules).Map(func(ru repository.Rule) *rule.MessageMatcher {
			matcher, err := rule.NewMessageMatcher(ru)
			if err != nil {
				log.Errorf("invalid rule: %s", err)
			}

			return matcher
		}).ToArray()

		// traverse all ungrouped messages to group
		collectingGroups := make(map[primitive.ObjectID]repository.MessageGroup)
		err = msgRepo.Traverse(bson.M{"status": repository.MessageStatusPending}, func(msg repository.Message) error {
			for _, m := range matchers {
				matcher := m.(*rule.MessageMatcher)
				matched, err := matcher.Match(msg)
				if err != nil {
					return err
				}

				// if the message matched a rule, update message's group_id and skip to next message
				if matched {
					if _, ok := collectingGroups[matcher.Rule().ID]; !ok {
						grp, err := groupRepo.CollectingGroup(matcher.Rule().ToGroupRule())
						if err != nil {
							log.Errorf("create collecting group failed: %s", err)
							return err
						}

						collectingGroups[matcher.Rule().ID] = grp
					}

					msg.GroupID = collectingGroups[matcher.Rule().ID].ID
					msg.Status = repository.MessageStatusGrouped
					return msgRepo.UpdateID(msg.ID, msg)
				}
			}

			msg.Status = repository.MessageStatusCanceled
			return msgRepo.UpdateID(msg.ID, msg)
		})

		if err != nil {
			log.Errorf("aggregate message failed: %s", err)
		}
	})
}
