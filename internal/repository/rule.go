package repository

import (
	"context"
	"sort"
	"time"

	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/coll"
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
	ReadyTypeTimeRange = "time_range"
)

type Tag struct {
	Name  string `bson:"_id" json:"name"`
	Count int64  `bson:"count" json:"count"`
}

type TimeRange struct {
	// StartTime 开始时间（包含该时间）
	StartTime string `bson:"start_time" json:"start_time"`
	// EndTime 截止时间（不包含改时间）
	EndTime  string `bson:"end_time" json:"end_time"`
	Interval int64  `bson:"interval" json:"interval"`
}

// Rule is a rule definition
type Rule struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Tags        []string           `bson:"tags" json:"tags"`
	// AggregateRule 聚合规则，同一个规则匹配的事件，会按照该规则返回的值进行更加精细的分组
	AggregateRule string `bson:"aggregate_rule" json:"aggregate_rule"`
	// RelationRule 关联规则，匹配的事件会被创建关联关系
	RelationRule string `bson:"relation_rule" json:"relation_rule"`

	// ReadType 就绪类型，支持 interval/daily_time
	ReadyType  string      `bson:"ready_type" json:"ready_type"`
	Interval   int64       `bson:"interval" json:"interval"`
	DailyTimes []string    `bson:"daily_times" json:"daily_times"`
	TimeRanges []TimeRange `bson:"time_ranges" json:"time_ranges"`

	// Rule 用于分组匹配的规则
	Rule string `bson:"rule" json:"rule"`
	// IgnoreRule 分组匹配后，检查 message 是否应该被忽略
	IgnoreRule      string    `bson:"ignore_rule" json:"ignore_rule"`
	Template        string    `bson:"template" json:"template"`
	SummaryTemplate string    `bson:"summary_template" json:"summary_template"`
	Triggers        []Trigger `bson:"triggers" json:"triggers"`

	// ReportTemplateID 报表模板 ID
	ReportTemplateID primitive.ObjectID `bson:"report_template_id" json:"report_template_id"`

	Status RuleStatus `bson:"status" json:"status"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// ToGroupRule convert Rule to EventGroupRule
func (rule Rule) ToGroupRule(aggregateKey string, msgType EventType) EventGroupRule {
	groupRule := EventGroupRule{
		ID:               rule.ID,
		Name:             rule.Name,
		Rule:             rule.Rule,
		IgnoreRule:       rule.IgnoreRule,
		Template:         rule.Template,
		SummaryTemplate:  rule.SummaryTemplate,
		ReportTemplateID: rule.ReportTemplateID,
		AggregateKey:     aggregateKey,
		Type:             msgType,
	}

	if rule.ReadyType == "" {
		rule.ReadyType = ReadyTypeInterval
	}

	switch rule.ReadyType {
	case ReadyTypeInterval:
		groupRule.ExpectReadyAt = time.Now().Add(time.Duration(rule.Interval) * time.Second)
	case ReadyTypeDailyTime:
		groupRule.ExpectReadyAt = ExpectReadyAt(time.Now(), rule.DailyTimes)
	case ReadyTypeTimeRange:
		groupRule.ExpectReadyAt = ExpectReadyAtInTimeRange(time.Now(), rule.TimeRanges)
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
	Tags(ctx context.Context) ([]Tag, error)
}

func ExpectReadyAt(now time.Time, dailyTimes []string) time.Time {
	// 查找最近的时间点
	var times Times
	todayTimeStr := now.Format(time.RFC3339)
	_ = coll.MustNew(dailyTimes).Map(func(dailyTime string) time.Time {
		ts, _ := time.Parse(time.RFC3339, todayTimeStr[0:11]+dailyTime[:5]+":00"+todayTimeStr[19:])
		return ts
	}).All(&times)

	sort.Sort(times)

	for _, t := range times {
		if now.Before(t) {
			return t
		}
	}

	return times[0].Add(24 * time.Hour)
}

type Times []time.Time

func (t Times) Len() int {
	return len(t)
}

func (t Times) Less(i, j int) bool {
	return t[i].Before(t[j])
}

func (t Times) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

// ExpectReadyAtInTimeRange 根据当前时间和时间范围，计算预期时间，时间范围为前开后闭区间
func ExpectReadyAtInTimeRange(now time.Time, timeRanges []TimeRange) time.Time {
	todayTimeStr := now.Format(time.RFC3339)
	zeroTime, _ := time.Parse(time.RFC3339, todayTimeStr[0:11]+"00:00:00"+todayTimeStr[19:])
	lastTime := zeroTime.AddDate(0, 0, 1)

	for _, t := range timeRanges {
		startTime, _ := time.Parse(time.RFC3339, todayTimeStr[0:11]+t.StartTime[:5]+":00"+todayTimeStr[19:])
		endTime, _ := time.Parse(time.RFC3339, todayTimeStr[0:11]+t.EndTime[:5]+":00"+todayTimeStr[19:])

		deadlineTs := now.Add(time.Duration(t.Interval) * time.Second)
		if startTime.After(endTime) {
			// zeroTime -1- endTime -2- startTime -3- lastTime ...
			// 当前时间在 1/3 位置时，匹配该规则，时间范围跨天
			// 预期执行时间如果超过该范围，则使用边界时间代替
			if (now.Before(endTime) && timeAfterOrEq(now, zeroTime)) || (now.Before(lastTime) && timeAfterOrEq(now, startTime)) {
				if now.Before(lastTime) && timeAfterOrEq(now, startTime) {
					if timeAfterOrEq(deadlineTs, endTime.AddDate(0, 0, 1)) {
						return lastTime
					}
				} else if now.Before(endTime) && timeAfterOrEq(now, zeroTime) {
					if timeAfterOrEq(deadlineTs, endTime) {
						return endTime
					}
				}

				return deadlineTs
			}
		}

		// startTime -1- endTime
		// 当前时间在 1 的位置匹配该规则，超过范围使用边界时间代替
		if now.Before(endTime) && timeAfterOrEq(now, startTime) {
			if timeAfterOrEq(deadlineTs, endTime) {
				return endTime
			}

			return deadlineTs
		}
	}

	return now.Add(60 * time.Second)
}

func timeBeforeOrEq(t1 time.Time, t2 time.Time) bool {
	return t1.Before(t2) || t1.Equal(t2)
}

func timeAfterOrEq(t1 time.Time, t2 time.Time) bool {
	return t1.After(t2) || t1.Equal(t2)
}
