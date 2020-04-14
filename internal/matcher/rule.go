package matcher

import (
	"fmt"
	"strings"
	"time"
)

type Helpers struct{}

// Lower returns a copy of the string s with all Unicode letters mapped to their lower case.
func (Helpers) Lower(val string) string {
	return strings.ToLower(val)
}

// Upper returns a copy of the string s with all Unicode letters mapped to their upper case.
func (Helpers) Upper(val string) string {
	return strings.ToUpper(val)
}

// DailyTimeBetween 判断当前时间（格式 15:04）是否在 startTime 和 endTime 之间
func (Helpers) DailyTimeBetween(startTime, endTime string) bool {
	start, err := time.Parse("15:04", startTime)
	if err != nil {
		panic(fmt.Sprintf("invalid startTime, must be formatted as 15:04, error is %v", err))
	}

	end, err := time.Parse("15:04", endTime)
	if err != nil {
		panic(fmt.Sprintf("invalid endTime, must be formatted as 15:04, error is %v", err))
	}

	if start.After(end) {
		end = end.Add(24 * time.Hour)
	}

	now, _ := time.Parse("15:04", time.Now().Format("15:04"))
	return now.After(start) && now.Before(end)
}

// Now return current time
func (Helpers) Now() time.Time {
	return time.Now()
}

// ParseTime parse a string to time.Time
// layout: Mon Jan 2 15:04:05 -0700 MST 2006
func (Helpers) ParseTime(layout string, value string) time.Time {
	ts, _ := time.Parse(layout, value)
	return ts
}
