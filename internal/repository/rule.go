package repository

import (
	"time"

	"github.com/mylxsw/asteria/log"
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
	ReadyType string `bson:"ready_type" json:"ready_type"`
	Interval  int64  `bson:"interval" json:"interval"`
	DailyTime string `bson:"daily_time" json:"daily_time"`

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
		groupRule.ExpectReadyAt = ExpectReadyAt(rule.DailyTime)
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

func ExpectReadyAt(dailyTime string) time.Time {
	// 2006-01-02T15:04:05Z07:00
	nextDayTimeStr := time.Now().Add(24 * time.Hour).Format(time.RFC3339)
	parsed, _ := time.Parse(time.RFC3339, nextDayTimeStr[0:11]+dailyTime[:5]+":00"+nextDayTimeStr[19:])
	return parsed
}
