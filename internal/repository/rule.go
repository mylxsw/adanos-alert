package repository

import (
	"sort"
	"time"

	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/coll"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RuleStatus string

const (
	RuleStatusEnabled  RuleStatus = "enabled"
	RuleStatusDisabled RuleStatus = "disabled"
)

const (
	ReadyTypeInterval  = "interval"
	ReadyTypeDailyTime = "daily_time"
)

type Tag struct {
	Name  string `bson:"_id" json:"name"`
	Count int64  `bson:"count" json:"count"`
}

// Rule is a rule definition
type Rule struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Tags        []string           `bson:"tags" json:"tags"`

	// ReadType 就绪类型，支持 interval/daily_time
	ReadyType  string   `bson:"ready_type" json:"ready_type"`
	Interval   int64    `bson:"interval" json:"interval"`
	DailyTimes []string `bson:"daily_times" json:"daily_times"`

	Rule            string    `bson:"rule" json:"rule"`
	Template        string    `bson:"template" json:"template"`
	SummaryTemplate string    `bson:"summary_template" json:"summary_template"`
	Triggers        []Trigger `bson:"triggers" json:"triggers"`

	Status RuleStatus `bson:"status" json:"status"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// ToGroupRule convert Rule to MessageGroupRule
func (rule Rule) ToGroupRule() MessageGroupRule {
	groupRule := MessageGroupRule{
		ID:              rule.ID,
		Name:            rule.Name,
		Rule:            rule.Rule,
		Template:        rule.Template,
		SummaryTemplate: rule.SummaryTemplate,
	}

	if rule.ReadyType == "" {
		rule.ReadyType = ReadyTypeInterval
	}

	switch rule.ReadyType {
	case ReadyTypeInterval:
		groupRule.ExpectReadyAt = time.Now().Add(time.Duration(rule.Interval) * time.Second)
	case ReadyTypeDailyTime:
		groupRule.ExpectReadyAt = ExpectReadyAt(time.Now(), rule.DailyTimes)
	default:
		log.Errorf("invalid readyType [%s] for ruleID=%s", rule.ReadyType, rule.ID.Hex())
	}

	return groupRule
}

type RuleRepo interface {
	Add(rule Rule) (id primitive.ObjectID, err error)
	Get(id primitive.ObjectID) (rule Rule, err error)
	Paginate(filter interface{}, offset, limit int64) (rules []Rule, next int64, err error)
	Find(filter bson.M) (rules []Rule, err error)
	Traverse(filter bson.M, cb func(rule Rule) error) error
	UpdateID(id primitive.ObjectID, rule Rule) error
	Count(filter bson.M) (int64, error)
	Delete(filter bson.M) error
	DeleteID(id primitive.ObjectID) error
	Tags() ([]Tag, error)
}

func ExpectReadyAt(now time.Time, dailyTimes []string) time.Time {
	// 查找最近的时间点
	var times Times
	todayTimeStr := now.Format(time.RFC3339)
	_ = coll.MustNew(dailyTimes).Map(func(dailyTime string) time.Time {
		ts, _ := time.Parse(time.RFC3339, todayTimeStr[0:11]+dailyTime[:5]+":00"+todayTimeStr[19:])
		return ts
	}).All(&times)

	sort.Sort(times)

	for _, t := range times {
		if now.Before(t) {
			return t
		}
	}

	return times[0].Add(24 * time.Hour)
}

type Times []time.Time

func (t Times) Len() int {
	return len(t)
}

func (t Times) Less(i, j int) bool {
	return t[i].Before(t[j])
}

func (t Times) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
