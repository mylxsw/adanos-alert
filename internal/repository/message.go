package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageMeta map[string]interface{}
type MessageStatus string
type MessageType string

const (
	// MessageStatusPending 待处理
	MessageStatusPending MessageStatus = "pending"
	// MessageStatusGrouped 已分组（已经匹配规则并且分组）
	MessageStatusGrouped MessageStatus = "grouped"
	// MessageStatusCanceled 已取消（没有任何匹配的规则）
	MessageStatusCanceled MessageStatus = "canceled"
	// MessageStatusExpired 已过期（有匹配的规则，但是当时没有匹配）
	MessageStatusExpired MessageStatus = "expired"
	// MessageStatusDead 死信（匹配规则，但是被主动忽略）
	MessageStatusIgnored MessageStatus = "ignored"

	// MessageTypePlain 普通消息
	MessageTypePlain MessageType = "plain"
	// MessageTypeRecoverable 可恢复消息
	MessageTypeRecoverable MessageType = "recoverable"
	// MessageTypeRecovery 恢复消息
	MessageTypeRecovery MessageType = "recovery"
)

type Message struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	SeqNum    int64                `bson:"seq_num" json:"seq_num"`
	Content   string               `bson:"content" json:"content"`
	Meta      MessageMeta          `bson:"meta" json:"meta"`
	Tags      []string             `bson:"tags" json:"tags"`
	Origin    string               `bson:"origin" json:"origin"`
	GroupID   []primitive.ObjectID `bson:"group_ids" json:"group_ids"`
	Type      MessageType          `bson:"type" json:"type"`
	Status    MessageStatus        `bson:"status" json:"status"`
	CreatedAt time.Time            `bson:"created_at" json:"created_at"`
}

type MessageRepo interface {
	AddWithContext(ctx context.Context, msg Message) (id primitive.ObjectID, err error)
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
