package controller

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/coll"
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/glacier/web"
	"github.com/mylxsw/go-utils/str"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// StatisticsController 统计功能
type StatisticsController struct {
	cc infra.Resolver
}

// NewStatisticsController create a new StatisticsController
func NewStatisticsController(cc infra.Resolver) web.Controller {
	return StatisticsController{cc: cc}
}

// Register 注册路由
func (s StatisticsController) Register(router web.Router) {
	router.Group("/statistics", func(router web.Router) {
		router.Get("/daily-group-counts/", s.DailyGroupCounts).Name("statistics:daily-group-counts")
		router.Get("/user-group-counts/", s.UserGroupCounts).Name("statistics:user-group-counts")
		router.Get("/rule-group-counts/", s.RuleGroupCounts).Name("statistics:rule-group-counts")
		router.Get("/group-agg-period-counts/", s.EventGroupAggInPeriod).Name("statistics:group-agg-period-counts")
		router.Get("/group-agg-counts/", s.EventGroupAggCounts).Name("statistics:group-agg-counts")

		router.Group("/events/", func(router web.Router) {
			router.Get("/period-counts/", s.EventCountInPeriod).Name("statistics:events:period-counts")
		})
	})
}

// MessageGroupByDatetimeCount 周期内事件组数量
type MessageGroupByDatetimeCount struct {
	Datetime      string `json:"datetime"`
	Total         int64  `json:"total"`
	TotalMessages int64  `json:"total_messages"`
}

// extractDateRange 从请求中提取时间范围
func extractDateRange(webCtx web.Context, defaultDays int) (time.Time, time.Time) {
	startAt := webCtx.Input("start_at")
	endAt := webCtx.Input("end_at")

	startTime := time.Now().Add(-time.Duration(defaultDays*24) * time.Hour)
	endTime := time.Now()

	if startAt != "" {
		if len(startAt) == 10 {
			startAt = startAt + " 00:00:00"
		}

		parsed, err := time.Parse("2006-01-02 15:04:05", startAt)
		if err == nil {
			startTime = parsed
		}
	}

	if endAt != "" {
		if len(endAt) == 10 {
			endAt = endAt + " 23:59:59"
		}

		parsed, err := time.Parse("2006-01-02 15:04:05", endAt)
		if err == nil {
			endTime = parsed
		}
	}

	return startTime, endTime
}

// DailyGroupCounts 每日报警次数汇总
func (s StatisticsController) DailyGroupCounts(ctx web.Context, groupRepo repository.EventGroupRepo) ([]MessageGroupByDatetimeCount, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx.Context(), 15*time.Second)
	defer cancel()

	startTime, endTime := extractDateRange(ctx, ctx.IntInput("days", 30))
	dailyCounts, err := groupRepo.StatByDatetimeCount(timeoutCtx, groupFilter(ctx), startTime, endTime, 24)
	if err != nil {
		return nil, err
	}

	if len(dailyCounts) == 0 {
		return make([]MessageGroupByDatetimeCount, 0), nil
	}

	dailyCountsByDate := make(map[string]MessageGroupByDatetimeCount)
	for _, d := range dailyCounts {
		datetime := d.Datetime.In(time.Local).Format("2006-01-02")
		dailyCountsByDate[datetime] = MessageGroupByDatetimeCount{
			Datetime:      datetime,
			Total:         d.Total,
			TotalMessages: d.TotalMessages,
		}
	}

	startDate := dailyCounts[0].Datetime
	endDate := dailyCounts[len(dailyCounts)-1].Datetime

	if log.DebugEnabled() {
		log.Debugf("%v: %v", startDate, endDate)
	}

	results := make([]MessageGroupByDatetimeCount, 0)

	for startDate.Before(endDate) || startDate.Equal(endDate) {
		startDateF := startDate.Format("2006-01-02")
		if d, ok := dailyCountsByDate[startDateF]; ok {
			results = append(results, d)
		} else {
			results = append(results, MessageGroupByDatetimeCount{
				Datetime:      startDateF,
				Total:         0,
				TotalMessages: 0,
			})
		}

		startDate = startDate.Add(24 * time.Hour)
	}

	return results, nil
}

type EventGroupByUserCounts []repository.EventGroupByUserCount

func (e EventGroupByUserCounts) Len() int {
	return len(e)
}

func (e EventGroupByUserCounts) Less(i, j int) bool {
	return e[i].Total < e[j].Total
}

func (e EventGroupByUserCounts) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

// UserGroupCounts 用户报警次数汇总
func (s StatisticsController) UserGroupCounts(ctx web.Context, groupRepo repository.EventGroupRepo) ([]repository.EventGroupByUserCount, error) {
	timeoutCtx, _ := context.WithTimeout(ctx.Context(), 5*time.Second)
	startTime, endTime := extractDateRange(ctx, ctx.IntInput("days", 30))
	res, err := groupRepo.StatByUserCount(timeoutCtx, startTime, endTime)
	if err != nil {
		return nil, err
	}

	sort.Sort(sort.Reverse(EventGroupByUserCounts(res)))

	if len(res) > 10 {
		other := repository.EventGroupByUserCount{
			UserName:      "Others",
			Total:         0,
			TotalMessages: 0,
		}
		for _, v := range res[10:] {
			other.Total += v.Total
			other.TotalMessages += v.TotalMessages
		}
		res = append(res[:10], other)
	}

	return res, nil
}

type EventGroupByRuleCounts []repository.EventGroupByRuleCount

func (e EventGroupByRuleCounts) Len() int {
	return len(e)
}

func (e EventGroupByRuleCounts) Less(i, j int) bool {
	return e[i].Total < e[j].Total
}

func (e EventGroupByRuleCounts) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

// RuleGroupCounts 报警规则报警次数汇总
func (s StatisticsController) RuleGroupCounts(ctx web.Context, groupRepo repository.EventGroupRepo) ([]repository.EventGroupByRuleCount, error) {
	timeoutCtx, _ := context.WithTimeout(ctx.Context(), 5*time.Second)
	startTime, endTime := extractDateRange(ctx, ctx.IntInput("days", 30))

	res, err := groupRepo.StatByRuleCount(timeoutCtx, startTime, endTime)
	if err != nil {
		return nil, err
	}

	sort.Sort(sort.Reverse(EventGroupByRuleCounts(res)))

	if len(res) > 10 {
		other := repository.EventGroupByRuleCount{
			RuleName:      "Others",
			Total:         0,
			TotalMessages: 0,
		}
		for _, v := range res[10:] {
			other.Total += v.Total
			other.TotalMessages += v.TotalMessages
		}
		res = append(res[:10], other)
	}

	return res, nil
}

// EventByDatetimeCount 周期内事件数量统计返回对象
type EventByDatetimeCount struct {
	Datetime string `json:"datetime"`
	Total    int64  `json:"total"`
}

// EventCountInPeriod 统计周期内的事件数量
// 支持的参数: days/step/format/meta/tags/origin/status/relation_id/group_id/event_id
func (s StatisticsController) EventCountInPeriod(webCtx web.Context, evtRepo repository.EventRepo) ([]EventByDatetimeCount, error) {
	ctx, cancel := context.WithTimeout(webCtx.Context(), 15*time.Second)
	defer cancel()

	dateTimeFormat := webCtx.InputWithDefault("format", "01-02 15:00")
	var step int64 = 1
	startDate, endDate := extractDateRange(webCtx, webCtx.IntInput("days", 7))

	if log.DebugEnabled() {
		log.Debugf("%v: %v", startDate, endDate)
	}

	filter := eventsFilter(webCtx)
	dailyCounts, err := evtRepo.CountByDatetime(ctx, filter, startDate, endDate, step)
	if err != nil {
		return nil, err
	}

	if len(dailyCounts) == 0 {
		return make([]EventByDatetimeCount, 0), nil
	}

	dailyCountsByDate := make(map[string]repository.EventByDatetimeCount)
	for _, d := range dailyCounts {
		datetime := d.Datetime.In(time.Local).Format(dateTimeFormat)
		dailyCountsByDate[datetime] = d
	}

	results := make([]EventByDatetimeCount, 0)
	startDateTmp := startDate.Add(time.Duration(step) * time.Hour)
	for startDateTmp.Before(endDate) || startDateTmp.Equal(endDate) {
		startDateF := startDateTmp.Format(dateTimeFormat)
		if d, ok := dailyCountsByDate[startDateF]; ok {
			results = append(results, EventByDatetimeCount{
				Datetime: startDateF,
				Total:    d.Total,
			})
		} else {
			results = append(results, EventByDatetimeCount{
				Datetime: startDateF,
				Total:    0,
			})
		}

		startDateTmp = startDateTmp.Add(time.Duration(step) * time.Hour)
	}

	return results, nil
}

// AggCount 聚合key包含的事件组数量
type AggCount struct {
	AggregateKey  string `json:"aggregate_key"`
	Total         int64  `json:"total"`
	TotalMessages int64  `json:"total_messages"`
}

// EventGroupAggByDatetimeCount 时间范围内事件组聚合数量
type EventGroupAggByDatetimeCount struct {
	Datetime string     `json:"datetime"`
	AggCount []AggCount `json:"agg_count"`
}

type EventGroupAggByDatetimeCountResp struct {
	Data          []EventGroupAggByDatetimeCount `json:"data"`
	AggregateKeys []string                       `json:"aggregate_keys"`
}

// EventCountInPeriod 统计周期内的事件组数量，按照聚合key分类返回
// 支持的参数: rule_id
func (s StatisticsController) EventGroupAggInPeriod(webCtx web.Context, evtGrpRepo repository.EventGroupRepo) (EventGroupAggByDatetimeCountResp, error) {
	ruleID, err := primitive.ObjectIDFromHex(webCtx.Input("rule_id"))
	if err != nil {
		return EventGroupAggByDatetimeCountResp{}, fmt.Errorf("invalid rule_id: %v", err)
	}

	startTime, endTime := extractDateRange(webCtx, webCtx.IntInput("days", 30))
	timeoutCtx, cancel := context.WithTimeout(webCtx.Context(), 15*time.Second)
	defer cancel()

	dailyCounts, err := evtGrpRepo.StatByAggCountInPeriod(timeoutCtx, ruleID, startTime, endTime, 24)
	if err != nil {
		return EventGroupAggByDatetimeCountResp{}, err
	}

	if len(dailyCounts) == 0 {
		return EventGroupAggByDatetimeCountResp{Data: []EventGroupAggByDatetimeCount{}, AggregateKeys: []string{}}, nil
	}

	aggregateKeys := make([]string, 0)
	dailyCountsByDate := make(map[string]EventGroupAggByDatetimeCount)
	for _, d := range dailyCounts {
		datetime := d.Datetime.In(time.Local).Format("2006-01-02")
		if _, ok := dailyCountsByDate[datetime]; !ok {
			dailyCountsByDate[datetime] = EventGroupAggByDatetimeCount{
				Datetime: datetime,
				AggCount: make([]AggCount, 0),
			}
		}

		dailyCountsByDate[datetime] = EventGroupAggByDatetimeCount{
			Datetime: dailyCountsByDate[datetime].Datetime,
			AggCount: append(dailyCountsByDate[datetime].AggCount, AggCount{
				AggregateKey:  d.AggregateKey,
				Total:         d.Total,
				TotalMessages: d.TotalMessages,
			}),
		}

		aggregateKeys = append(aggregateKeys, d.AggregateKey)
	}

	aggregateKeys = str.Distinct(aggregateKeys)

	var defaultAggCounts []AggCount
	_ = coll.MustNew(aggregateKeys).Map(func(k string) AggCount {
		return AggCount{AggregateKey: k}
	}).All(&defaultAggCounts)

	startDate := dailyCounts[0].Datetime
	endDate := dailyCounts[len(dailyCounts)-1].Datetime

	if log.DebugEnabled() {
		log.Debugf("%v: %v", startDate, endDate)
	}

	results := make([]EventGroupAggByDatetimeCount, 0)

	for startDate.Before(endDate) || startDate.Equal(endDate) {
		startDateF := startDate.Format("2006-01-02")
		if d, ok := dailyCountsByDate[startDateF]; ok {
			aggMap := make(map[string]AggCount)
			for _, v := range d.AggCount {
				aggMap[v.AggregateKey] = v
			}

			fullAgg := make([]AggCount, 0)
			for _, v := range defaultAggCounts {
				if ex, ok := aggMap[v.AggregateKey]; ok {
					fullAgg = append(fullAgg, ex)
				} else {
					fullAgg = append(fullAgg, v)
				}
			}

			d.AggCount = fullAgg
			results = append(results, d)
		} else {
			results = append(results, EventGroupAggByDatetimeCount{
				Datetime: startDateF,
				AggCount: defaultAggCounts,
			})
		}

		startDate = startDate.Add(24 * time.Hour)
	}

	return EventGroupAggByDatetimeCountResp{
		Data:          results,
		AggregateKeys: aggregateKeys,
	}, nil
}

// EventGroupAggCounts 事件组聚合Key数量统计
// 参数 rule_id
func (s StatisticsController) EventGroupAggCounts(webCtx web.Context, evtGrpRepo repository.EventGroupRepo) ([]repository.EventGroupAggCount, error) {
	ruleID, err := primitive.ObjectIDFromHex(webCtx.Input("rule_id"))
	if err != nil {
		return nil, fmt.Errorf("invalid rule_id: %v", err)
	}

	timeoutCtx, cancel := context.WithTimeout(webCtx.Context(), 15*time.Second)
	defer cancel()

	startTime, endTime := extractDateRange(webCtx, webCtx.IntInput("days", 30))
	aggCounts, err := evtGrpRepo.StatByAggCount(timeoutCtx, ruleID, startTime, endTime)
	if err != nil {
		return nil, err
	}

	return aggCounts, nil
}
