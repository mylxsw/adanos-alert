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
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Interval  int64              `bson:"interval" json:"interval"`
	Threshold int64              `bson:"threshold" json:"threshold"`
	Rule      string             `bson:"rule" json:"rule"`
}

type MessageGroup struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	MessageCount int64            `bson:"message_count" json:"message_count"`
	Rule         MessageGroupRule `bson:"rule" json:"rule"`

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
	Update(id primitive.ObjectID, grp MessageGroup) error
	Count(filter bson.M) (int64, error)

	CollectingGroup(rule MessageGroupRule) (group MessageGroup, err error)
}
