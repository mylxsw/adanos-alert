package controller

import (
	"context"
	"time"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier/web"
)

// StatisticsController 统计功能
type StatisticsController struct {
	cc container.Container
}

// NewStatisticsController create a new StatisticsController
func NewStatisticsController(cc container.Container) web.Controller {
	return &StatisticsController{cc: cc}
}

// Register 注册路由
func (s *StatisticsController) Register(router *web.Router) {
	router.Group("/statistics", func(router *web.Router) {
		router.Get("/daily-group-counts/", s.DailyGroupCounts).Name("statistics:daily-group-counts")
		router.Get("/user-group-counts/", s.UserGroupCounts).Name("statistics:user-group-counts")
		router.Get("/rule-group-counts/", s.RuleGroupCounts).Name("statistics:rule-group-counts")

		router.Group("/events/", func(router *web.Router) {
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

// DailyGroupCounts 每日报警次数汇总
func (s *StatisticsController) DailyGroupCounts(ctx web.Context, groupRepo repository.EventGroupRepo) ([]MessageGroupByDatetimeCount, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx.Context(), 5*time.Second)
	defer cancel()

	dailyCounts, err := groupRepo.StatByDatetimeCount(timeoutCtx, time.Now().Add(-30*24*time.Hour), time.Now(), 24)
	if err != nil {
		return nil, err
	}

	if len(dailyCounts) == 0 {
		return make([]MessageGroupByDatetimeCount, 0), nil
	}

	dailyCountsByDate := make(map[string]MessageGroupByDatetimeCount)
	for _, d := range dailyCounts {
		datetime := d.Datetime.Format("2006-01-02")
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

// UserGroupCounts 用户报警次数汇总
func (s *StatisticsController) UserGroupCounts(ctx web.Context, groupRepo repository.EventGroupRepo) ([]repository.EventGroupByUserCount, error) {
	timeoutCtx, _ := context.WithTimeout(ctx.Context(), 5*time.Second)
	return groupRepo.StatByUserCount(timeoutCtx, time.Now().Add(-30*24*time.Hour), time.Now())
}

// RuleGroupCounts 报警规则报警次数汇总
func (s *StatisticsController) RuleGroupCounts(ctx web.Context, groupRepo repository.EventGroupRepo) ([]repository.EventGroupByRuleCount, error) {
	timeoutCtx, _ := context.WithTimeout(ctx.Context(), 5*time.Second)
	return groupRepo.StatByRuleCount(timeoutCtx, time.Now().Add(-30*24*time.Hour), time.Now())
}

// EventByDatetimeCount 周期内事件数量统计返回对象
type EventByDatetimeCount struct {
	Datetime string `json:"datetime"`
	Total    int64  `json:"total"`
}

// EventCountInPeriod 统计周期内的事件数量
// 支持的参数: days/step/format/meta/tags/origin/status/relation_id/group_id/event_id
func (s *StatisticsController) EventCountInPeriod(webCtx web.Context, evtRepo repository.EventRepo) ([]EventByDatetimeCount, error) {
	ctx, cancel := context.WithTimeout(webCtx.Context(), 5*time.Second)
	defer cancel()

	dayRange := webCtx.IntInput("days", 3)
	dateTimeFormat := webCtx.InputWithDefault("format", "01-02 15:00")
	step := webCtx.Int64Input("step", 1)

	startDate := time.Now().Add(-time.Duration(dayRange*24) * time.Hour)
	endDate := time.Now()

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
		datetime := d.Datetime.Format(dateTimeFormat)
		dailyCountsByDate[datetime] = d
	}

	// 找到查询数据集中与第一个元素事件最接近的时间
	// 该时间与第一个元素事件的差值就是时间段间隔的偏移，步长减去
	// 偏移来实现时间范围和从 DB 中查询出的时间范围一致
	diffHour := 0.0
	startDateTmp := startDate
	for startDateTmp.Before(endDate) || startDateTmp.Equal(endDate) {
		diffHour = dailyCounts[0].Datetime.Sub(startDateTmp).Hours()
		if diffHour <= float64(step) {
			break
		}

		startDateTmp = startDateTmp.Add(time.Duration(step) * time.Hour)
	}

	results := make([]EventByDatetimeCount, 0)
	startDateTmp = startDate.Add(time.Duration(step-int64(diffHour)) * time.Hour)
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
