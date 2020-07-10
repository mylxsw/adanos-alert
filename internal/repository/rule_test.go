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
