package repository

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RuleStatus string

const (
	RuleStatusEnabled  RuleStatus = "enabled"
	RuleStatusDisabled RuleStatus = "disabled"
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

	Interval int64 `bson:"interval" json:"interval"`

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
	return MessageGroupRule{
		ID:              rule.ID,
		Name:            rule.Name,
		Interval:        rule.Interval,
		Rule:            rule.Rule,
		Template:        rule.Template,
		SummaryTemplate: rule.SummaryTemplate,
	}
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
