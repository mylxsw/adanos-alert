package repository

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageGroupStatus string

const (
	MessageGroupStatusCollecting MessageGroupStatus = "collecting"
	MessageGroupStatusPending    MessageGroupStatus = "pending"
	MessageGroupStatusOK         MessageGroupStatus = "ok"
	MessageGroupStatusFailed     MessageGroupStatus = "failed"
	MessageGroupStatusCanceled   MessageGroupStatus = "canceled"
)

type MessageGroupRule struct {
	ID              primitive.ObjectID `bson:"_id" json:"id"`
	Name            string             `bson:"name" json:"name"`
	Interval        int64              `bson:"interval" json:"interval"`
	Threshold       int64              `bson:"threshold" json:"threshold"`
	Rule            string             `bson:"rule" json:"rule"`
	Template        string             `bson:"template" json:"template"`
	SummaryTemplate string             `bson:"summary_template" json:"summary_template"`
}

type MessageGroup struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	MessageCount int64            `bson:"message_count" json:"message_count"`
	Rule         MessageGroupRule `bson:"rule" json:"rule"`
	Actions      []Trigger        `bson:"actions" json:"actions"`

	Status    MessageGroupStatus `bson:"status" json:"status"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type MessageGroupRepo interface {
	Add(grp MessageGroup) (id primitive.ObjectID, err error)
	Get(id primitive.ObjectID) (grp MessageGroup, err error)
	Find(filter bson.M) (grps []MessageGroup, err error)
	Paginate(filter bson.M, offset, limit int64) (grps []MessageGroup, next int64, err error)
	Delete(filter bson.M) error
	DeleteID(id primitive.ObjectID) error
	Traverse(filter bson.M, cb func(grp MessageGroup) error) error
	UpdateID(id primitive.ObjectID, grp MessageGroup) error
	Count(filter bson.M) (int64, error)

	// LastGroup get last group which match the filter in messageGroups
	LastGroup(filter bson.M) (grp MessageGroup, err error)
	CollectingGroup(rule MessageGroupRule) (group MessageGroup, err error)
}
