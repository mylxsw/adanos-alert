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

type Rule struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`

	Interval  int64 `bson:"interval" json:"interval"`
	Threshold int64 `bson:"threshold" json:"threshold"`
	Priority  int64 `bson:"priority" json:"priority"`

	Rule   string     `bson:"rule" json:"rule"`
	Status RuleStatus `bson:"status" json:"status"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type RuleRepo interface {
	Add(rule Rule) (id primitive.ObjectID, err error)
	Get(id primitive.ObjectID) (rule Rule, err error)
	Find(filter bson.M) (rules []Rule, err error)
	Traverse(filter bson.M, cb func(rule Rule) error) error
	UpdateID(id primitive.ObjectID, rule Rule) error
	Count(filter bson.M) (int64, error)
	Delete(filter bson.M) error
	DeleteID(id primitive.ObjectID) error
}
