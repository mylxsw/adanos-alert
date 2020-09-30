package repository_test

import (
	"testing"
	"time"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestExpectReadyAt(t *testing.T) {
	var times = []string{"09:00", "04:00", "12:00", "23:50", "15:30", "18:00"}

	{
		rs := repository.ExpectReadyAt(parseTime("2020-07-10T23:55:16+08:00"), times).Format(time.RFC3339)
		assert.Equal(t, "2020-07-11T04:00:00+08:00", rs)
	}

	{
		rs := repository.ExpectReadyAt(parseTime("2020-07-10T20:55:16+08:00"), times).Format(time.RFC3339)
		assert.Equal(t, "2020-07-10T23:50:00+08:00", rs)
	}

	{
		rs := repository.ExpectReadyAt(parseTime("2020-07-10T09:55:16+08:00"), times).Format(time.RFC3339)
		assert.Equal(t, "2020-07-10T12:00:00+08:00", rs)
	}
}

func parseTime(t string) time.Time {
	p, _ := time.Parse(time.RFC3339, t)
	return p
}

func TestExpectReadyAtTimeRange(t *testing.T) {
	{
		var timeRanges = []repository.TimeRange{
			{StartTime: "09:00", EndTime: "18:00", Interval: 300}, // 5 分钟
			{StartTime: "18:00", EndTime: "09:00", Interval: 600}, // 10 分钟
		}

		assert.Equal(t, "2020-07-10T10:00:16+08:00", repository.ExpectReadyAtInTimeRange(parseTime("2020-07-10T09:55:16+08:00"), timeRanges).Format(time.RFC3339))
		assert.Equal(t, "2020-07-10T23:10:00+08:00", repository.ExpectReadyAtInTimeRange(parseTime("2020-07-10T23:00:00+08:00"), timeRanges).Format(time.RFC3339))
	}

	{
		var timeRanges = []repository.TimeRange{
			{StartTime: "09:00", EndTime: "20:00:00", Interval: 900}, // 15 分钟
			{StartTime: "20:00:00", EndTime: "09:00", Interval: 7200}, // 2 小时
		}

		assert.Equal(t, "2020-07-10T10:10:16+08:00", repository.ExpectReadyAtInTimeRange(parseTime("2020-07-10T09:55:16+08:00"), timeRanges).Format(time.RFC3339))
		assert.Equal(t, "2020-07-10T09:00:00+08:00", repository.ExpectReadyAtInTimeRange(parseTime("2020-07-10T08:34:16+08:00"), timeRanges).Format(time.RFC3339))
		assert.Equal(t, "2020-07-11T01:00:00+08:00", repository.ExpectReadyAtInTimeRange(parseTime("2020-07-10T23:00:00+08:00"), timeRanges).Format(time.RFC3339))
		assert.Equal(t, "2020-07-10T02:00:00+08:00", repository.ExpectReadyAtInTimeRange(parseTime("2020-07-10T00:00:00+08:00"), timeRanges).Format(time.RFC3339))
		assert.Equal(t, "2020-07-10T09:15:00+08:00", repository.ExpectReadyAtInTimeRange(parseTime("2020-07-10T09:00:00+08:00"), timeRanges).Format(time.RFC3339))
		assert.Equal(t, "2020-07-10T22:00:00+08:00", repository.ExpectReadyAtInTimeRange(parseTime("2020-07-10T20:00:00+08:00"), timeRanges).Format(time.RFC3339))
		assert.Equal(t, "2020-07-10T22:00:01+08:00", repository.ExpectReadyAtInTimeRange(parseTime("2020-07-10T20:00:01+08:00"), timeRanges).Format(time.RFC3339))
	}
}
