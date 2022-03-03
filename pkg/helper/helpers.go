package helper

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/mylxsw/asteria/log"
	"html"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/itchyny/gojq"
	"github.com/mylxsw/adanos-alert/pkg/misc"
	"github.com/mylxsw/coll"
	"github.com/mylxsw/durafmt"
	"github.com/tidwall/gjson"
)

// Helpers 用于规则引擎的助手函数
type Helpers struct{}

var helper = Helpers{}

// NewHelpers create a new helper
func NewHelpers() Helpers {
	return helper
}

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
	matched := now.After(start) && now.Before(end)

	if log.DebugEnabled() {
		log.WithFields(log.Fields{
			"start":   start,
			"end":     end,
			"now":     now,
			"matched": matched,
		}).Debug("helper function: dailyTimeBetween")
	}

	return matched
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

// CutOff 字符串截断
func (Helpers) CutOff(maxLen int, val string) string {
	valRune := []rune(strings.Trim(val, " \n"))

	valLen := len(valRune)
	if valLen <= maxLen {
		return string(valRune)
	}

	return string(valRune[0:maxLen]) + "..."
}

// stringMask create a mask for string
func (Helpers) Mask(left int, content string) string {
	size := len(content)
	if size < 16 {
		return strings.Repeat("*", size)
	}

	return content[:left] + strings.Repeat("*", size-left*2) + content[size-left:]
}

// Split 字符串分割
func (Helpers) Split(sep string, content string) []string {
	return strings.Split(content, sep)
}

// FilterEmptyLines 从字符串中移除空行
func (Helpers) FilterEmptyLines(content string) string {
	return strings.Trim(
		coll.MustNew(strings.Split(content, "\n")).
			Map(func(line string) string {
				return strings.TrimRight(line, " ")
			}).
			Filter(func(line string) bool {
				return line != ""
			}).
			Reduce(func(carry string, item string) string {
				return fmt.Sprintf("%s\n%s", carry, item)
			}, "").(string),
		"\n",
	)
}

// Join 字符串数组拼接
func (Helpers) Join(elements interface{}, sep string) string {
	switch elements.(type) {
	case []string:
		return strings.Join(elements.([]string), sep)
	case []interface{}:
		items := make([]string, 0)
		for _, ele := range elements.([]interface{}) {
			items = append(items, fmt.Sprintf("%v", ele))
		}
		return strings.Join(items, sep)
	default:
		return "invalid join types"
	}
}

// Repeat 字符串重复 count 次
func (Helpers) Repeat(count int, s string) string {
	return strings.Repeat(s, count)
}

// NumberBeauty 字符串数字格式化
func (Helpers) NumberBeauty(number interface{}) string {
	str, ok := number.(string)
	if !ok {
		str = fmt.Sprintf("%.2f", number)
	}

	length := len(str)
	if length < 4 {
		return str
	}
	arr := strings.Split(str, ".") //用小数点符号分割字符串,为数组接收
	length1 := len(arr[0])
	if length1 < 4 {
		return str
	}
	count := (length1 - 1) / 3
	for i := 0; i < count; i++ {
		arr[0] = arr[0][:length1-(i+1)*3] + "," + arr[0][length1-(i+1)*3:]
	}
	return strings.Join(arr, ".")
}

// Float 字符串转浮点数
func (Helpers) Float(numStr string) float64 {
	val, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0
	}

	return val
}

// Int 字符串转整数
func (Helpers) Int(numStr string) int {
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return 0
	}

	return num
}

// Empty 检查字符串是否为空
// `空白`, `0`，`任意大小写的 false` 会被认为是 false，其它情况为 true
func (Helpers) Empty(str string) bool {
	str = strings.TrimSpace(str)
	if str == "" || str == "0" {
		return true
	}

	return strings.EqualFold("false", str)
}

// JQuery 使用 JQ 表达式提取 json 字符串中的值
func (helper Helpers) JQuery(data string, expression string, suppressError bool) string {
	query, err := gojq.Parse(expression)
	if err != nil {
		return encodeError(err, suppressError)
	}

	var dataInterface interface{}
	if err := json.Unmarshal([]byte(data), &dataInterface); err != nil {
		return encodeError(err, suppressError)
	}

	iter := query.Run(dataInterface)
	values, err := encodeValues(iter)
	if err != nil {
		return encodeError(err, suppressError)
	}

	return values
}

func encodeError(err error, suppressError bool) string {
	if suppressError {
		return ""
	}

	return fmt.Sprintf("<ERROR> %s", err)
}

func encodeValues(iter gojq.Iter) (string, error) {
	buffer := bytes.NewBuffer(nil)
	for {
		val, ok := iter.Next()
		if !ok {
			break
		}

		switch v := val.(type) {
		case error:
			return "", v
		case [2]interface{}:
			if s, ok := v[0].(string); ok {
				if s == "STDERR:" {
					val = v[1]
				}
			}
		}

		switch val.(type) {
		case string:
			buffer.Write([]byte(val.(string)))
		default:
			xs, err := json.Marshal(val)
			if err != nil {
				return "", err
			}
			buffer.Write(xs)
		}
	}

	return buffer.String(), nil
}

// DOMQueryOne 从 HTML DOM 对象中查询第 index 个匹配 selector 的元素内容
func (helper Helpers) DOMQueryOne(selector string, index int, htmlContent string) string {
	eles := helper.DOMQuery(selector, htmlContent)
	if len(eles) > index {
		return eles[index]
	}

	return ""
}

// DOMQuery 从 HTML DOM 对象中查询所有匹配 selector 的元素
func (helper Helpers) DOMQuery(selector string, htmlContent string) []string {
	reader, err := goquery.NewDocumentFromReader(bytes.NewBufferString(htmlContent))
	if err != nil {
		return []string{}
	}

	res := make([]string, 0)
	reader.Find(selector).Each(func(i int, s *goquery.Selection) {
		h, err := s.Html()

		if err == nil && strings.TrimSpace(h) != "" {
			res = append(res, html.UnescapeString(strings.TrimSpace(h)))
		}
	})

	return res
}

// JSONArray return array elements from path
func (helper Helpers) JSONArray(content string, path string) []gjson.Result {
	if path == "" {
		return gjson.Parse(content).Array()
	}

	if strings.HasPrefix(path, ".") {
		res := make([]gjson.Result, 0)
		for _, val := range gjson.Parse(content).Array() {
			res = append(res, val.Get(strings.TrimLeft(path, ".")).Array()...)
		}

		return res
	}

	return gjson.Get(content, path).Array()
}

// JSONStrArray return string array from json
func (helper Helpers) JSONStrArray(content string, path string) []string {
	strarrs := make([]string, 0)
	for _, res := range helper.JSONArray(content, path) {
		strarrs = append(strarrs, res.String())
	}

	return strarrs
}

// JSONIntArray return int64 array from json
func (helper Helpers) JSONIntArray(content string, path string) []int64 {
	intarrs := make([]int64, 0)
	for _, res := range helper.JSONArray(content, path) {
		intarrs = append(intarrs, res.Int())
	}

	return intarrs
}

// JSONFloatArray return float64 array from json
func (helper Helpers) JSONFloatArray(content string, path string) []float64 {
	floatarrs := make([]float64, 0)
	for _, res := range helper.JSONArray(content, path) {
		floatarrs = append(floatarrs, res.Float())
	}

	return floatarrs
}

// JSONBoolArray return bool array from json
func (helper Helpers) JSONBoolArray(content string, path string) []bool {
	boolarrs := make([]bool, 0)
	for _, res := range helper.JSONArray(content, path) {
		boolarrs = append(boolarrs, res.Bool())
	}

	return boolarrs
}

// JSON return string content from json
func (helper Helpers) JSON(content string, path string) string {
	return gjson.Get(content, path).Raw
}

// JSONInt return int content from json
func (helper Helpers) JSONInt(content string, path string) int64 {
	return gjson.Get(content, path).Int()
}

// JSONFloat return float64 content from json
func (helper Helpers) JSONFloat(content string, path string) float64 {
	return gjson.Get(content, path).Float()
}

// JSONBool return bool content from json
func (helper Helpers) JSONBool(content string, path string) bool {
	return gjson.Get(content, path).Bool()
}

// String convert any data to string
func (helper Helpers) String(data interface{}) string {
	return fmt.Sprintf("%v", data)
}

// JSONEncode convert any data to json string
func (helper Helpers) JSONEncode(data interface{}) string {
	b, _ := json.Marshal(data)
	return string(b)
}

// HumanDuration 时间段格式化
func (helper Helpers) HumanDuration(duration string) string {
	parsed, err := durafmt.ParseString(duration)
	if err != nil {
		return duration
	}

	return parsed.String("年", "周", "天", "小时", "分钟", "秒", "毫秒", "微秒")
}
