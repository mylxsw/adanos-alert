package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventMeta map[string]interface{}
type EventStatus string
type EventType string

const (
	// EventStatusPending 待处理
	EventStatusPending EventStatus = "pending"
	// EventStatusGrouped 已分组（已经匹配规则并且分组）
	EventStatusGrouped EventStatus = "grouped"
	// EventStatusCanceled 已取消（没有任何匹配的规则）
	EventStatusCanceled EventStatus = "canceled"
	// EventStatusExpired 已过期（有匹配的规则，但是当时没有匹配）
	EventStatusExpired EventStatus = "expired"
	// MessageStatusDead 死信（匹配规则，但是被主动忽略）
	EventStatusIgnored EventStatus = "ignored"

	// EventTypePlain 普通消息
	EventTypePlain EventType = "plain"
	// EventTypeRecoverable 可恢复消息
	EventTypeRecoverable EventType = "recoverable"
	// EventTypeRecovery 恢复消息
	EventTypeRecovery EventType = "recovery"
	// EventTypeRecovered 消息已恢复
	EventTypeRecovered EventType = "recovered"
)

type Event struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	SeqNum    int64                `bson:"seq_num" json:"seq_num"`
	Content   string               `bson:"content" json:"content"`
	Meta      EventMeta            `bson:"meta" json:"meta"`
	Tags      []string             `bson:"tags" json:"tags"`
	Origin    string               `bson:"origin" json:"origin"`
	GroupID   []primitive.ObjectID `bson:"group_ids" json:"group_ids"`
	Type      EventType            `bson:"type" json:"type"`
	Status    EventStatus          `bson:"status" json:"status"`
	CreatedAt time.Time            `bson:"created_at" json:"created_at"`
}

type EventRepo interface {
	AddWithContext(ctx context.Context, msg Event) (id primitive.ObjectID, err error)
	Add(msg Event) (id primitive.ObjectID, err error)
	Get(id primitive.ObjectID) (msg Event, err error)
	Find(filter interface{}) (messages []Event, err error)
	Paginate(filter interface{}, offset, limit int64) (messages []Event, next int64, err error)
	Delete(filter interface{}) error
	DeleteID(id primitive.ObjectID) error
	Traverse(filter interface{}, cb func(msg Event) error) error
	UpdateID(id primitive.ObjectID, update Event) error
	Count(filter interface{}) (int64, error)
}
