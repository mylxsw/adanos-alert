package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventGroupStatus string

const (
	EventGroupStatusCollecting EventGroupStatus = "collecting"
	EventGroupStatusPending    EventGroupStatus = "pending"
	EventGroupStatusOK         EventGroupStatus = "ok"
	EventGroupStatusFailed     EventGroupStatus = "failed"
	EventGroupStatusCanceled   EventGroupStatus = "canceled"
)

type EventGroupRule struct {
	ID   primitive.ObjectID `bson:"_id" json:"id"`
	Name string             `bson:"name" json:"name"`

	// AggregateKey 通过该 Key 对同一个规则下的 message 分组
	AggregateKey string    `bson:"aggregate_key" json:"aggregate_key"`
	Type         EventType `bson:"type" json:"type"`

	// ExpectReadyAt 预期就绪时间，当超过该时间后，Group自动关闭，发起通知
	ExpectReadyAt time.Time `bson:"expect_ready_at" json:"expect_ready_at"`
	Realtime      bool      `bson:"realtime" json:"realtime"`

	Rule            string `bson:"rule" json:"rule"`
	IgnoreRule      string `bson:"ignore_rule" json:"ignore_rule"`
	IgnoreMaxCount  int    `bson:"ignore_max_count" json:"ignore_max_count"`
	Template        string `bson:"template" json:"template"`
	SummaryTemplate string `bson:"summary_template" json:"summary_template"`

	// Report template
	ReportTemplateID primitive.ObjectID `bson:"report_template_id" json:"report_template_id"`
}

type EventGroup struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	SeqNum int64              `bson:"seq_num" json:"seq_num"`

	// AggregateKey 与 .Rule.AggregateKey 相同，方便读取
	AggregateKey string    `bson:"aggregate_key" json:"aggregate_key"`
	Type         EventType `bson:"type" json:"type"`

	MessageCount int64          `bson:"message_count" json:"message_count"`
	Rule         EventGroupRule `bson:"rule" json:"rule"`
	Actions      []Trigger      `bson:"actions" json:"actions"`

	Status    EventGroupStatus `bson:"status" json:"status"`
	CreatedAt time.Time        `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time        `bson:"updated_at" json:"updated_at"`
}

// Ready return whether the message group has reached close conditions
func (grp *EventGroup) Ready() bool {
	return grp.Rule.ExpectReadyAt.Before(time.Now())
}

type EventGroupByRuleCount struct {
	RuleID        primitive.ObjectID `bson:"rule_id" json:"rule_id"`
	RuleName      string             `bson:"rule_name" json:"rule_name"`
	Total         int64              `bson:"total" json:"total"`
	TotalMessages int64              `bson:"total_messages" json:"total_messages"`
}

type EventGroupByUserCount struct {
	UserID        primitive.ObjectID `bson:"user_id" json:"user_id"`
	UserName      string             `bson:"user_name" json:"user_name"`
	Total         int64              `bson:"total" json:"total"`
	TotalMessages int64              `bson:"total_messages" json:"total_messages"`
}

type EventGroupByDatetimeCount struct {
	Datetime      time.Time `bson:"datetime" json:"datetime"`
	Total         int64     `bson:"total" json:"total"`
	TotalMessages int64     `bson:"total_messages" json:"total_messages"`
}

// EventGroupAggByDatetimeCount 时间范围内事件组聚合数量
type EventGroupAggByDatetimeCount struct {
	Datetime      time.Time `bson:"datetime" json:"datetime"`
	AggregateKey  string    `bson:"aggregate_key" json:"aggregate_key"`
	Total         int64     `bson:"total" json:"total"`
	TotalMessages int64     `bson:"total_messages" json:"total_messages"`
}

// EventGroupAggCount 事件组聚合数量
type EventGroupAggCount struct {
	AggregateKey string `bson:"aggregate_key" json:"aggregate_key"`
	Total        int64  `bson:"total" json:"total"`
}

type EventGroupRepo interface {
	Add(grp EventGroup) (id primitive.ObjectID, err error)
	Get(id primitive.ObjectID) (grp EventGroup, err error)
	Find(filter bson.M) (grps []EventGroup, err error)
	Paginate(filter bson.M, offset, limit int64, sortByCreatedAtDesc bool) (grps []EventGroup, next int64, err error)
	Delete(filter bson.M) error
	DeleteID(id primitive.ObjectID) error
	Traverse(filter bson.M, cb func(grp EventGroup) error) error
	UpdateID(id primitive.ObjectID, grp EventGroup) error
	Count(filter bson.M) (int64, error)

	// LastGroup get last group which match the filter in messageGroups
	LastGroup(filter bson.M) (grp EventGroup, err error)
	CollectingGroup(rule EventGroupRule) (group EventGroup, err error)

	// StatByRuleCount 按照规则的维度，查询规则相关的报警次数
	StatByRuleCount(ctx context.Context, startTime, endTime time.Time) ([]EventGroupByRuleCount, error)
	StatByUserCount(ctx context.Context, startTime, endTime time.Time) ([]EventGroupByUserCount, error)
	StatByDatetimeCount(ctx context.Context, filter bson.M, startTime, endTime time.Time, hour int64) ([]EventGroupByDatetimeCount, error)
	StatByAggCountInPeriod(ctx context.Context, ruleID primitive.ObjectID, startTime, endTime time.Time, hour int64) ([]EventGroupAggByDatetimeCount, error)
	StatByAggCount(ctx context.Context, ruleID primitive.ObjectID, startTime, endTime time.Time) ([]EventGroupAggCount, error)
}
