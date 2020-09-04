package template

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/mylxsw/adanos-alert/pkg/array"
	pkgJSON "github.com/mylxsw/adanos-alert/pkg/json"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/coll"
	"github.com/mylxsw/go-toolkit/jsonutils"
	"github.com/pingcap/parser"
	"github.com/vjeantet/grok"
)

// Parse parse template with data to html
func Parse(templateStr string, data interface{}) (string, error) {
	var buffer bytes.Buffer
	if err := template.Must(CreateParser(templateStr)).Execute(&buffer, data); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

// CreateParse create a template parser
func CreateParser(templateStr string) (*template.Template, error) {
	funcMap := template.FuncMap{
		"cutoff":                     cutOff,
		"implode":                    Implode,
		"explode":                    strings.Split,
		"ident":                      leftIdent,
		"json":                       jsonFormatter,
		"datetime":                   datetimeFormat,
		"datetime_noloc":             datetimeFormatNoLoc,
		"reformat_datetime_str":      reformatDatetimeStr,
		"parse_datetime_str":         parseDatetime,
		"parse_datetime_str_rfc3339": parseDatetimeRFC3339,
		"json_get":                   pkgJSON.Get,
		"json_gets":                  pkgJSON.Gets,
		"json_array":                 pkgJSON.GetArray,
		"json_flatten":               jsonFlatten,
		"starts_with":                startsWith,
		"ends_with":                  endsWith,
		"trim":                       strings.Trim,
		"trim_right":                 strings.TrimSuffix,
		"trim_left":                  strings.TrimPrefix,
		"trim_space":                 strings.TrimSpace,
		"format":                     fmt.Sprintf,
		"number_beauty":              NumberBeauty,
		"integer":                    toInteger,
		"mysql_slowlog":              parseMySQLSlowlog,
		"sql_finger":                 SQLFinger,
		"open_falcon_im":             ParseOpenFalconImMessage,
		"string_mask":                StringMask,
		"string_tags":                StringTags,
		"remove_empty_line":          RemoveEmptyLine,
		"meta_filter":                MetaFilter,
		"meta_prefix_filter":         MetaFilterPrefix,
		"serialize":                  Serialize,
		"sort_map_human":             SortMapByKeyHuman,
		"important":                  importantNotice,
	}

	return template.New("").Funcs(funcMap).Parse(templateStr)
}

// StringTags split tags string to array
func StringTags(tags string, sep string) []string {
	if tags == "" {
		return []string{}
	}

	var result []string
	_ = coll.Filter(strings.Split(tags, sep), &result, func(s string) bool {
		return strings.TrimSpace(s) != ""
	})

	return result
}

// parseMySQLSlowlog 解析mysql慢查询日志
func parseMySQLSlowlog(slowlog string) map[string]string {
	g, _ := grok.NewWithConfig(&grok.Config{NamedCapturesOnly: true})
	_ = g.AddPattern("SQL", "(.|\r|\n)*")
	values, _ := g.Parse(`(?m)^(# Time: \d+ \d+:\d+:\d+\n)?#\s+User@Host:\s+%{USER:user}\[[^\]]+\]\s+@\s+(?:%{DATA:clienthost})?\[(?:%{IPV4:clientip})?\]\n#\s+Thread_id:\s+%{NUMBER:thread_id}\s+Schema:\s+%{WORD:schema}\s+QC_hit:\s+%{WORD:qc_hit}\n#\s*Query_time:\s+%{NUMBER:query_time}\s+Lock_time:\s+%{NUMBER:lock_time}\s+Rows_sent:\s+%{NUMBER:rows_sent}\s+Rows_examined:\s+%{NUMBER:rows_examined}\n(# Rows_affected: %{NUMBER:rows_affected}  Bytes_sent: %{NUMBER:bytes_sent}\n)?\s*(?:use %{DATA:database};\s*\n)?SET\s+timestamp=%{NUMBER:occur_time};\n\s*%{SQL:sql};\s*(?:\n#\s+Time)?.*$`, slowlog)

	return values
}

// SQLFinger 生成 SQL 指纹
func SQLFinger(sqlStr string) string {
	return strings.ReplaceAll(parser.Normalize(sqlStr), " . ", ".")
}

// cutOff 字符串截断
func cutOff(maxLen int, val string) string {
	valRune := []rune(strings.Trim(val, " \n"))

	valLen := len(valRune)
	if valLen <= maxLen {
		return string(valRune)
	}

	return string(valRune[0:maxLen])
}

// 字符串多行缩进
func leftIdent(ident string, message string) string {
	result := coll.MustNew(strings.Split(message, "\n")).Map(func(line string) string {
		return ident + line
	}).Reduce(func(carry string, line string) string {
		return fmt.Sprintf("%s\n%s", carry, line)
	}, "").(string)

	return strings.Trim(result, "\n")
}

// JSONBeauty format content as json beauty
func JSONBeauty(content string) string {
	return jsonFormatter(content)
}

// json格式化输出
func jsonFormatter(content string) string {
	var output bytes.Buffer
	if err := json.Indent(&output, []byte(content), "", "    "); err != nil {
		return content
	}

	return output.String()
}

// parseDatetime 将字符串解析为时间格式
func parseDatetime(layout string, dt string) time.Time {
	parsed, err := time.Parse(layout, dt)
	if err != nil {
		return time.Now()
	}

	return parsed
}

// parseDateTimeRFC3339 解析 RFC3339 格式的时间字符串为 time.Time
func parseDatetimeRFC3339(dt string) time.Time {
	return parseDatetime(time.RFC3339, dt)
}

// reformatDatetimeStr 时间字符串重新格式化
func reformatDatetimeStr(originalLayout, targetLayout string, dt string) string {
	loc, _ := time.LoadLocation("Asia/Chongqing")
	return parseDatetime(originalLayout, dt).In(loc).Format(targetLayout)
}

// datetimeFormat 时间格式化，不使用Location
func datetimeFormatNoLoc(layout string, datetime time.Time) string {
	return datetime.Format(layout)
}

// datetimeFormat 时间格式化
func datetimeFormat(layout string, datetime time.Time) string {
	loc, _ := time.LoadLocation("Asia/Chongqing")

	return datetime.In(loc).Format(layout)
}

type KvPairs []jsonutils.KvPair

func (k KvPairs) Len() int {
	return len(k)
}

func (k KvPairs) Less(i, j int) bool {
	return k[i].Key < k[j].Key
}

func (k KvPairs) Swap(i, j int) {
	k[i], k[j] = k[j], k[i]
}

// jsonFlatten json转换为kv pairs
func jsonFlatten(body string, maxLevel int) []jsonutils.KvPair {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("json解析失败: %s", err)
		}
	}()

	ju, err := jsonutils.New([]byte(body), maxLevel, true)
	if err != nil {
		return make([]jsonutils.KvPair, 0)
	}

	kvPairs := ju.ToKvPairsArray()
	sort.Sort(KvPairs(kvPairs))

	return kvPairs
}

// startsWith 判断是字符串开始
func startsWith(haystack string, needles ...string) bool {
	for _, n := range needles {
		if strings.HasPrefix(haystack, n) {
			return true
		}
	}

	return false
}

// endsWith 判断字符串结尾
func endsWith(haystack string, needles ...string) bool {
	for _, n := range needles {
		if strings.HasSuffix(haystack, n) {
			return true
		}
	}

	return false
}

// toInteger 转换为整数
func toInteger(str string) int {
	val, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}

	return val
}

// StringMask create a mask for string
func StringMask(content string, left int) string {
	return stringMask(content, left)
}

// stringMask create a mask for string
func stringMask(content string, left int) string {
	size := len(content)
	if size < 16 {
		return strings.Repeat("*", size)
	}

	return content[:left] + strings.Repeat("*", size-left*2) + content[size-left:]
}

// RemoveEmptyLine 从字符串中移除空行
func RemoveEmptyLine(content string) string {
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

// MetaFilter 过滤 Meta，只保留允许的key
func MetaFilter(meta map[string]interface{}, allowKeys ...string) map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range meta {
		if array.StringsContain(k, allowKeys) {
			res[k] = v
		}
	}

	return res
}

// MetaFilter 过滤 Meta，只保留以 allowKeyPrefix 开头的项
func MetaFilterPrefix(meta map[string]interface{}, allowKeyPrefix ...string) map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range meta {
		if array.StringsContainPrefix(k, allowKeyPrefix) {
			res[k] = v
		}
	}

	return res
}

// Serialize 对象序列化为字符串，用于展示
func Serialize(data interface{}) string {
	serialized, err := json.Marshal(data)
	if err != nil {
		return fmt.Sprintf("%v", data)
	}

	return string(serialized)
}

type KVPair struct {
	Key   string
	Value interface{}
}

type Keys []string

func (ks Keys) Len() int {
	return len(ks)
}

func (ks Keys) Less(i, j int) bool {
	if strings.HasPrefix(ks[i], "message") && !strings.HasPrefix(ks[j], "message") {
		return true
	}
	if strings.HasPrefix(ks[j], "message") && !strings.HasPrefix(ks[i], "message") {
		return false
	}

	return ks[i] < ks[j]
}

func (ks Keys) Swap(i, j int) {
	ks[i], ks[j] = ks[j], ks[i]
}

func SortMapByKeyHuman(data map[string]interface{}) []KVPair {
	keys := make(Keys, 0)
	for k := range data {
		keys = append(keys, k)
	}

	sort.Sort(keys)

	kvPairs := make([]KVPair, 0)
	for _, v := range keys {
		kvPairs = append(kvPairs, KVPair{
			Key:   v,
			Value: data[v],
		})
	}

	return kvPairs
}

func Implode(elems interface{}, sep string) string {
	if _, ok := elems.([]string); ok {
		return strings.Join(elems.([]string), sep)
	}

	elemsType := reflect.TypeOf(elems).Kind()
	if elemsType == reflect.Array || elemsType == reflect.Slice {
		joinStrs := make([]string, 0)
		elemsVal := reflect.ValueOf(elems)
		for i := 0; i < elemsVal.Len(); i++ {
			joinStrs = append(joinStrs, fmt.Sprintf("%v", elemsVal.Index(i).Interface()))
		}

		return strings.Join(joinStrs, sep)
	}

	return fmt.Sprintf("%v", elems)
}

// NumberBeauty 字符串数字格式化
func NumberBeauty(number interface{}) string {
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

func importantNotice() string {
	return `<font color="#ea2426">【重要】</font>`
}
