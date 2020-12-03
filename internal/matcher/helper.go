package matcher

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/mylxsw/adanos-alert/pkg/misc"
)

// Helpers 用于规则引擎的助手函数
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

// SQLFinger 将 SQL 转换为其指纹
func (Helpers) SQLFinger(sqlStr string) string {
	return misc.SQLFinger(sqlStr)
}

// TrimSuffix 字符串去除后缀
func (Helpers) TrimSuffix(s, suffix string) string {
	return strings.TrimSuffix(s, suffix)
}

// TrimPrefix 字符串去除前缀
func (Helpers) TrimPrefix(s, prefix string) string {
	return strings.TrimPrefix(s, prefix)
}

// CutoffLine 字符串截取指定行数
func (Helpers) CutoffLine(val string, maxLine int) string {
	lines := strings.Split(val, "\n")
	if len(lines) > maxLine {
		return strings.Join(lines[:maxLine], "\n")
	}

	return strings.Join(lines, "\n")
}

// MD5 对 data 进行 Hash，生成 MD5 值
func (Helpers) MD5(data interface{}) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%v", data))))
}

// Sha1 对 data 进行 Hash，生成 Sha1 值
func (Helpers) Sha1(data interface{}) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(fmt.Sprintf("%v", data))))
}

// Base64 将 data 编码为 base64
func (Helpers) Base64(data interface{}) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v", data)))
}
