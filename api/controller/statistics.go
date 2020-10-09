package controller

import (
	"context"
	"time"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier/web"
)

type StatisticsController struct {
	cc container.Container
}

func NewStatisticsController(cc container.Container) web.Controller {
	return &StatisticsController{cc: cc}
}

func (s *StatisticsController) Register(router *web.Router) {
	router.Group("/statistics", func(router *web.Router) {
		router.Get("/daily-group-counts/", s.DailyGroupCounts).Name("statistics:daily-group-counts")
		router.Get("/user-group-counts/", s.UserGroupCounts).Name("statistics:user-group-counts")
		router.Get("/rule-group-counts/", s.RuleGroupCounts).Name("statistics:rule-group-counts")
	})
}

type MessageGroupByDatetimeCount struct {
	Datetime      string `json:"datetime"`
	Total         int64  `json:"total"`
	TotalMessages int64  `json:"total_messages"`
}

// DailyGroupCount 每日报警次数汇总
func (s *StatisticsController) DailyGroupCounts(ctx web.Context, groupRepo repository.EventGroupRepo) ([]MessageGroupByDatetimeCount, error) {
	timeoutCtx, _ := context.WithTimeout(ctx.Context(), 5*time.Second)
	dailyCounts, err := groupRepo.StatByDatetimeCount(timeoutCtx, time.Now().Add(- 30*24*time.Hour), time.Now(), 24)
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

	log.Debugf("%v: %v", startDate, endDate)

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
	return groupRepo.StatByUserCount(timeoutCtx, time.Now().Add(- 30*24*time.Hour), time.Now())
}

// RuleGroupCounts 报警规则报警次数汇总
func (s *StatisticsController) RuleGroupCounts(ctx web.Context, groupRepo repository.EventGroupRepo) ([]repository.EventGroupByRuleCount, error) {
	timeoutCtx, _ := context.WithTimeout(ctx.Context(), 5*time.Second)
	return groupRepo.StatByRuleCount(timeoutCtx, time.Now().Add(- 30*24*time.Hour), time.Now())
}
