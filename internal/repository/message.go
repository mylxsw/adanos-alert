package repository

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageMeta map[string]interface{}
type MessageStatus string

const (
	MessageStatusPending  MessageStatus = "pending"
	MessageStatusGrouped  MessageStatus = "grouped"
	MessageStatusCanceled MessageStatus = "canceled"
)

type Message struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	SeqNum    int64                `bson:"seq_num" json:"seq_num"`
	Content   string               `bson:"content" json:"content"`
	Meta      MessageMeta          `bson:"meta" json:"meta"`
	Tags      []string             `bson:"tags" json:"tags"`
	Origin    string               `bson:"origin" json:"origin"`
	GroupID   []primitive.ObjectID `bson:"group_ids" json:"group_ids"`
	Status    MessageStatus        `bson:"status" json:"status"`
	CreatedAt time.Time            `bson:"created_at" json:"created_at"`
}

type MessageRepo interface {
	Add(msg Message) (id primitive.ObjectID, err error)
	Get(id primitive.ObjectID) (msg Message, err error)
	Find(filter interface{}) (messages []Message, err error)
	Paginate(filter interface{}, offset, limit int64) (messages []Message, next int64, err error)
	Delete(filter interface{}) error
	DeleteID(id primitive.ObjectID) error
	Traverse(filter interface{}, cb func(msg Message) error) error
	UpdateID(id primitive.ObjectID, update Message) error
	Count(filter interface{}) (int64, error)
}
