package template

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/microcosm-cc/bluemonday"
	"github.com/mylxsw/adanos-alert/internal/repository"
	pkgJSON "github.com/mylxsw/adanos-alert/pkg/json"
	"github.com/mylxsw/adanos-alert/pkg/misc"
	"github.com/mylxsw/adanos-alert/pkg/strarr"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/coll"
	"github.com/mylxsw/go-toolkit/jsonutils"
	"github.com/russross/blackfriday/v2"
	"github.com/vjeantet/grok"
	"github.com/yosssi/gohtml"
	"go.mongodb.org/mongo-driver/bson"
)

type SimpleContainer interface {
	Get(key interface{}) (interface{}, error)
}

// Parse parse template with data to html
func Parse(cc SimpleContainer, templateStr string, data interface{}) (string, error) {
	var buffer bytes.Buffer
	par, err := CreateParser(cc, templateStr)
	if err != nil {
		return "", err
	}
	if err := par.Execute(&buffer, data); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

// CreateParse create a template parser
func CreateParser(cc SimpleContainer, templateStr string) (*template.Template, error) {
	funcMap := template.FuncMap{
		"cutoff":                     cutOff,
		"cutoff_line":                CutOffLine,
		"json_fields_cutoff":         JSONCutOffFields,
		"map_fields_cutoff":          MapFieldsCutoff,
		"implode":                    Implode,
		"explode":                    strings.Split,
		"join":                       join,
		"split":                      split,
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
		"format":                     fmt.Sprintf,
		"number_beauty":              NumberBeauty,
		"integer":                    toInteger,
		"float":                      toFloat,
		"mysql_slowlog":              parseMySQLSlowlog,
		"sql_finger":                 misc.SQLFinger,
		"open_falcon_im":             ParseOpenFalconImMessage,

		"serialize":            Serialize,
		"sort_map_human":       SortMapByKeyHuman,
		"error_notice":         errorNotice,
		"success_notice":       successNotice,
		"error_success_notice": errorOrSuccessNotice,
		"condition":            conditionStr,
		"recoverable_notice":   recoverableNotice,
		"user_metas":           BuildUserMetasFunc(cc),

		"meta_filter":                MetaFilter,
		"meta_filter_exclude":        MetaFilterExclude,
		"meta_prefix_filter":         MetaFilterPrefix,
		"meta_prefix_filter_exclude": MetaFilterPrefixExclude,

		"prefix_all_str":      prefixStrArray,
		"suffix_all_str":      suffixStrArray,
		"trim_prefix_map_k":   TrimPrefixMapK,
		"line_filter_include": LineFilterInclude,
		"line_filter_exclude": LineFilterExclude,

		"starts_with":       startsWith,
		"ends_with":         endsWith,
		"trim":              strings.Trim,
		"trim_right":        strings.TrimSuffix,
		"trim_left":         strings.TrimPrefix,
		"trim_space":        strings.TrimSpace,
		"remove_empty_line": RemoveEmptyLine,

		"string_mask": StringMask,
		"string_tags": StringTags,
		"str_upper":   strings.ToUpper,
		"str_lower":   strings.ToLower,
		"str_replace": strings.ReplaceAll,
		"str_repeat":  strings.Repeat,

		"md2html":           Markdown2html,
		"dom_filter_html":   DOMFilterHTML,
		"dom_filter_html_n": DOMFilterHTMLIndex,
		"html_beauty":       FormatHTML,
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
	_ = g.AddPattern("EXPLAIN", "(.|\r|\n)*?")
	values, _ := g.Parse(`(?m)^(# Time: \d+ \d+:\d+:\d+\n)?#\s+User@Host:\s+%{USER:user}\[[^\]]+\]\s+@\s+(?:%{DATA:clienthost})?\[(?:%{IPV4:clientip})?\]\n#\s+Thread_id:\s+%{NUMBER:thread_id}\s+Schema:\s+%{WORD:schema}\s+QC_hit:\s+%{WORD:qc_hit}\n#\s*Query_time:\s+%{NUMBER:query_time}\s+Lock_time:\s+%{NUMBER:lock_time}\s+Rows_sent:\s+%{NUMBER:rows_sent}\s+Rows_examined:\s+%{NUMBER:rows_examined}\n(# Rows_affected: %{NUMBER:rows_affected}  Bytes_sent: %{NUMBER:bytes_sent}\n)?%{EXPLAIN:explain}\s*(?:use %{DATA:database};\s*\n)?SET\s+timestamp=%{NUMBER:occur_time};\n\s*%{SQL:sql};\s*(?:\n#\s+Time)?.*$`, slowlog)

	return values
}

// cutOff 字符串截断
func cutOff(maxLen int, val string) string {
	valRune := []rune(strings.Trim(val, " \n"))

	valLen := len(valRune)
	if valLen <= maxLen {
		return string(valRune)
	}

	return string(valRune[0:maxLen]) + "..."
}

// CutOffLine 字符串截取指定行数
func CutOffLine(maxLine int, val string) string {
	lines := strings.Split(val, "\n")
	if len(lines) > maxLine {
		return strings.Join(lines[:maxLine], "\n") + "\n..."
	}

	return strings.Join(lines, "\n")
}

func LineFilterInclude(includeStr string, val string) string {
	lines := strings.Split(val, "\n")
	var results []string
	_ = coll.Filter(lines, &results, func(line string) bool {
		matched, _ := regexp.MatchString(includeStr, line)
		return matched || strings.Contains(line, includeStr)
	})

	return strings.Join(results, "\n")
}

func LineFilterExclude(excludeStr string, val string) string {
	lines := strings.Split(val, "\n")
	var results []string
	_ = coll.Filter(lines, &results, func(line string) bool {
		matched, _ := regexp.MatchString(excludeStr, line)
		return !matched && !strings.Contains(line, excludeStr)
	})

	return strings.Join(results, "\n")
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

// toFloat 字符串转 float64
func toFloat(str string) float64 {
	val, err := strconv.ParseFloat(str, 64)
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
		if strarr.In(k, allowKeys) {
			res[k] = v
		}
	}

	return res
}

// MetaFilterExclude 过滤 Meta，排除不允许的key
func MetaFilterExclude(meta map[string]interface{}, excludeKeys ...string) map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range meta {
		if !strarr.In(k, excludeKeys) {
			res[k] = v
		}
	}

	return res
}

// MetaFilter 过滤 Meta，只保留以 allowKeyPrefix 开头的项
func MetaFilterPrefix(meta map[string]interface{}, allowKeyPrefix ...string) map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range meta {
		if strarr.HasPrefixes(k, allowKeyPrefix) {
			res[k] = v
		}
	}

	return res
}

// MetaFilterPrefixExclude 过滤 Meta，排除以 disableKeyPrefix 开头的项
func MetaFilterPrefixExclude(meta map[string]interface{}, disableKeyPrefixes ...string) map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range meta {
		if !strarr.HasPrefixes(k, disableKeyPrefixes) {
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

// errorNotice 错误提示
func errorNotice(msg string) string {
	return fmt.Sprintf(`<font color="#ea2426">%s</font>`, msg)
}

// successNotice 成功提示
func successNotice(msg string) string {
	return fmt.Sprintf(`<font color="#27a745">%s</font>`, msg)
}

// errorOrSuccessNotice 错误或者成功提示， 根据 success 参数判断
func errorOrSuccessNotice(success bool, msg string) string {
	if success {
		return successNotice(msg)
	}
	return errorNotice(msg)
}

// recoverableNotice 可恢复消息提示
func recoverableNotice(recovered bool, msg string) string {
	if recovered {
		return successNotice("【已恢复】" + msg)
	}

	return errorNotice(msg)
}

// BuildUserMetasFunc 构建查询用户元信息的函数
func BuildUserMetasFunc(cc SimpleContainer) func(queryK, queryV string, field string) []string {
	userRepoR, err := cc.Get(new(repository.UserRepo))
	if err != nil {
		return func(queryK, queryV string, field string) []string {
			return []string{fmt.Sprintf("<error>:%v", err)}
		}
	}

	userRepo := userRepoR.(repository.UserRepo)
	return func(queryK, queryV string, field string) []string {
		filter := bson.M{}
		if strarr.In(queryK, []string{"name", "phone", "email", "role", "status"}) {
			filter[queryK] = queryV
		} else {
			filter["metas.key"] = queryK
			filter["metas.value"] = queryV
		}

		users, err := userRepo.Find(filter)
		if err != nil {
			return []string{}
		}

		var res []string
		_ = coll.MustNew(users).Map(func(u repository.User) string {
			switch field {
			case "name":
				return u.Name
			case "phone":
				return u.Phone
			case "email":
				return u.Email
			case "role":
				return u.Role
			case "status":
				return string(u.Status)
			default:
				for _, m := range u.Metas {
					if m.Key == field {
						return m.Value
					}
				}

				return ""
			}
		}).Filter(func(v string) bool { return v != "" }).All(&res)
		return res
	}
}

func join(sep string, elems []string) string {
	return Implode(elems, sep)
}

func split(sep, s string) []string {
	return strings.Split(s, sep)
}

func prefixStrArray(prefix string, arr []string) []string {
	var dest []string
	_ = coll.Map(arr, &dest, func(s string) string { return prefix + s })
	return dest
}

func suffixStrArray(suffix string, arr []string) []string {
	var dest []string
	_ = coll.Map(arr, &dest, func(s string) string { return s + suffix })
	return dest
}

// JSONCutOffFields 对 JSON 字符串扁平化，然后对每个 KV 截取指定长度
func JSONCutOffFields(length int, body string) map[string]interface{} {
	var pairs []KVPair
	_ = coll.MustNew(jsonFlatten(body, 3)).Map(func(p jsonutils.KvPair) KVPair {
		return KVPair{
			Key:   cutOff(30, p.Key),
			Value: cutOff(length, p.Value),
		}
	}).Filter(func(p KVPair) bool {
		return p.Key != "" && p.Value != ""
	}).All(&pairs)

	data := make(map[string]interface{})
	for _, p := range pairs {
		data[p.Key] = p.Value
	}

	return data
}

// MapFieldsCutoff 对 Map 的每个 KV 截取指定长度
func MapFieldsCutoff(length int, source map[string]interface{}) map[string]interface{} {
	data := make(map[string]interface{})
	for k, v := range source {
		data[cutOff(30, k)] = cutOff(length, fmt.Sprintf("%v", v))
	}

	return data
}

// TrimPrefixMapK 移除 Map 中所有 Key 的前缀
func TrimPrefixMapK(prefix string, source map[string]interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range source {
		res[strings.TrimPrefix(k, prefix)] = v
	}
	return res
}

// conditionStr 条件输出字符串，符合条件，输出 s1，否则 s2
func conditionStr(s1, s2 string, cond bool) string {
	if cond {
		return s1
	}

	return s2
}

// Markdown2html 将 Markdown 转换为 HTML
func Markdown2html(mc string) string {
	unsafe := blackfriday.Run([]byte(mc))
	return string(bluemonday.UGCPolicy().SanitizeBytes(unsafe))
}

// DOMFilterHTMLIndex 从 HTML DOM 对象中查询第 index 个匹配 selector 的元素内容
func DOMFilterHTMLIndex(selector string, index int, html string) string {
	eles := DOMFilterHTML(selector, html)
	if len(eles) > index {
		return eles[index]
	}

	return ""
}

// DOMFilterHTML 从 HTML DOM 对象中查询所有匹配 selector 的元素
func DOMFilterHTML(selector string, original string) []string {
	reader, err := goquery.NewDocumentFromReader(bytes.NewBufferString(original))
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

// FormatHTML 格式化 HTML 内容
func FormatHTML(html string) string {
	return gohtml.Format(html)
}
