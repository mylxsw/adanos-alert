package repository

import (
	"context"
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
	ID   primitive.ObjectID `bson:"_id" json:"id"`
	Name string             `bson:"name" json:"name"`

	// AggregateKey 通过该 Key 对同一个规则下的 message 分组
	AggregateKey string      `bson:"aggregate_key" json:"aggregate_key"`
	Type         MessageType `bson:"type" json:"type"`

	// ExpectReadyAt 预期就绪时间，当超过该时间后，Group自动关闭，发起通知
	ExpectReadyAt time.Time `bson:"expect_ready_at" json:"expect_ready_at"`

	Rule            string `bson:"rule" json:"rule"`
	IgnoreRule      string `bson:"ignore_rule" json:"ignore_rule"`
	Template        string `bson:"template" json:"template"`
	SummaryTemplate string `bson:"summary_template" json:"summary_template"`
}

type MessageGroup struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	SeqNum int64              `bson:"seq_num" json:"seq_num"`

	// AggregateKey 与 .Rule.AggregateKey 相同，方便读取
	AggregateKey string      `bson:"aggregate_key" json:"aggregate_key"`
	Type         MessageType `bson:"type" json:"type"`

	MessageCount int64            `bson:"message_count" json:"message_count"`
	Rule         MessageGroupRule `bson:"rule" json:"rule"`
	Actions      []Trigger        `bson:"actions" json:"actions"`

	Status    MessageGroupStatus `bson:"status" json:"status"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

// Ready return whether the message group has reached close conditions
func (grp *MessageGroup) Ready() bool {
	return grp.Rule.ExpectReadyAt.Before(time.Now())
}

type MessageGroupByRuleCount struct {
	RuleID        primitive.ObjectID `bson:"rule_id" json:"rule_id"`
	RuleName      string             `bson:"rule_name" json:"rule_name"`
	Total         int64              `bson:"total" json:"total"`
	TotalMessages int64              `bson:"total_messages" json:"total_messages"`
}

type MessageGroupByUserCount struct {
	UserID        primitive.ObjectID `bson:"user_id" json:"user_id"`
	UserName      string             `bson:"user_name" json:"user_name"`
	Total         int64              `bson:"total" json:"total"`
	TotalMessages int64              `bson:"total_messages" json:"total_messages"`
}

type MessageGroupByDatetimeCount struct {
	Datetime      time.Time `bson:"datetime" json:"datetime"`
	Total         int64     `bson:"total" json:"total"`
	TotalMessages int64     `bson:"total_messages" json:"total_messages"`
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

	// Statistics
	// StatByRuleCount 按照规则的维度，查询规则相关的报警次数
	StatByRuleCount(ctx context.Context, startTime, endTime time.Time) ([]MessageGroupByRuleCount, error)
	StatByUserCount(ctx context.Context, startTime, endTime time.Time) ([]MessageGroupByUserCount, error)
	StatByDatetimeCount(ctx context.Context, startTime, endTime time.Time, hour int64) ([]MessageGroupByDatetimeCount, error)
}
