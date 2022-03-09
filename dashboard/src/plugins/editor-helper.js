import {CodeMirror} from "vue-codemirror-lite";

let helpers = {
    helpers: [
        {text: "Upper", args: ['KEY'], displayText: "Upper(val string) string  | 字符串转大写"},
        {text: "Lower", args: ['KEY'], displayText: "Lower(val string) string  | 字符串转小写"},
        {text: "Now", args: [], displayText: "Now() time.Time  | 当前时间"},
        {text: "ParseTime", args: ['LAYOUT', 'VALUE'], displayText: "ParseTime(layout string, value string) time.Time | 时间字符串转时间对象"},
        {text: "DailyTimeBetween", args: ['START_TIME_STR', 'END_TIME_STR'], displayText: "DailyTimeBetween(startTime, endTime string) bool  | 判断当前时间是否在 startTime 和 endTime 之间（每天），时间格式为 15:04"},
        {text: 'SQLFinger', args: ['SQL_STR'], displayText: "SQLFinger(sqlStr string) string | 创建 SQL 指纹"},
        {text: 'TrimSuffix', args: ['STR', 'SUFFIX'], displayText: 'TrimSuffix(str, suffix string) string | 去除字符串后缀'},
        {text: 'TrimPrefix', args: ['STR', 'PREFIX'], displayText: 'TrimPrefix(str, prefix string) string | 去除字符串前缀'},
        {text: 'CutoffLine', args: ['STR', 'MAX'], displayText: 'CutoffLine(val string, maxLine int) string | 字符串截取 maxLine 行'},
        {text: 'MD5', args: ['DATA'], displayText: 'MD5(data interface{}) string  | 生成 data 的 md5 值'},
        {text: 'Sha1', args: ['DATA'], displayText: 'Sha1(data interface{}) string  | 生成 data 的 sha1 值'},
        {text: 'Base64', args: ['DATA'], displayText: 'Base64(data interface{}) string  | 生成 data 的 base64 编码值'},
        {text: 'CutOff', args: ['MAXLEN', 'STR'], displayText: 'CutOff(maxLen int, val string) string  | 字符串截取最大长度'},
        {text: 'Mask', args: ['LEFT', 'STR'], displayText: 'Mask(left int, content string) string  | 字符串掩码，只保留两侧 left 个字符，其余使用*替换'},
        {text: 'Split', args: ['SEP', 'STR'], displayText: 'Split(sep string, content string) []string  | 字符串使用 sep 切割为字符串数组'},
        {text: 'FilterEmptyLines', args: ['STR'], displayText: 'FilterEmptyLines(content string) string  | 去除字符串中的空行'},
        {text: 'Join', args: ['ELEMENTS', 'SEP'], displayText: 'Join(elements interface{}, sep string) string  | 数组使用 sep 连接为字符串'},
        {text: 'Repeat', args: ['COUNT', 'STR'], displayText: 'Repeat(count int, s string) string  | 重复字符串 COUNT 次'},
        {text: 'NumberBeauty', args: ['NUMBER'], displayText: 'NumberBeauty(number interface{}) string  | 数值格式化'},
        {text: 'Float', args: ['NUMSTR'], displayText: 'Float(numStr string) float64  | 字符串转 float64'},
        {text: 'Int', args: ['NUMSTR'], displayText: 'Int(numStr string) int  | 字符串转 int'},
        {text: 'Empty', args: ['NUMSTR'], displayText: 'Empty(str string) bool  | 检查字符串是否为空，`空白`, `0`，`任意大小写的 false` 会被认为是 false，其它情况为 true'},
        {text: 'JQuery', args: ['DATA', 'EXPR', 'SUPPRESS_ERR'], displayText: 'JQuery(data string, expression string, suppressError bool) string  | 使用 JQ 表达式提取 json 字符串中的值'},
        {text: 'DOMQueryOne', args: ['SELECTOR', 'INDEX HTML'], displayText: 'DOMQueryOne(selector string, index int, htmlContent string) string  | 从 HTML DOM 对象中查询第 index 个匹配 selector 的元素内容'},
        {text: 'DOMQuery', args: ['SELECTOR', 'HTML'], displayText: 'DOMQuery(selector string, htmlContent string) []string  | 从 HTML DOM 对象中查询所有匹配 selector 的元素'},
        {text: 'JSONArray', args: ['STR', 'PATH'], displayText: 'JSONArray(content string, path string) []gjson.Result  | return array elements from path'},
        {text: 'JSONStrArray', args: ['STR', 'PATH'], displayText: 'JSONStrArray(content string, path string) []string  | return string array from json'},
        {text: 'JSONIntArray', args: ['STR', 'PATH'], displayText: 'JSONIntArray(content string, path string) []int64  | return int64 array from json'},
        {text: 'JSONFloatArray', args: ['STR', 'PATH'], displayText: 'JSONFloatArray(content string, path string) []float64  | return float64 array from json'},
        {text: 'JSONBoolArray', args: ['STR', 'PATH'], displayText: 'JSONBoolArray(content string, path string) []bool  | return bool array from json'},
        {text: 'JSON', args: ['STR', 'PATH'], displayText: 'JSON(content string, path string) string  | return string content from json'},
        {text: 'JSONInt', args: ['STR', 'PATH'], displayText: 'JSONInt(content string, path string) int64  | return int content from json'},
        {text: 'JSONFloat', args: ['STR', 'PATH'], displayText: 'JSONFloat(content string, path string) float64  | return float64 content from json'},
        {text: 'JSONBool', args: ['STR', 'PATH'], displayText: 'JSONBool(content string, path string) bool  | return bool content from json'},
        {text: 'String', args: ['DATA'], displayText: 'String(data interface}) string  | convert any data to string'},
        {text: 'JSONEncode', args: ['DATA'], displayText: 'JSONEncode(data interface}) string  | convert any data to json string'},
        {text: 'HumanDuration', args: ['STR'], displayText: 'HumanDuration(duration string) string | 时间段格式化，以人类可读的形式展示'}
    ],
    groupMatchRules: [
        {text: 'Content', displayText: 'Content | 字段类型：string |  事件内容，字符串格式'},
        {text: 'Meta[""]', displayText: 'Meta | 字段类型：map[string]interface{} | 字段，字典类型'},
        {text: 'Tags[0]', displayText: 'Tags | 字段类型：[]string | 字段，数组类型'},
        {text: 'Origin', displayText: 'Origin | 字段类型：string | 事件来源，字符串'},
        {text: "JsonGet(KEY, DEFAULT)", displayText: "JsonGet(key string, defaultValue string) string  | 将事件体作为json解析，获取指定的key"},
        {text: "IsRecovery()", displayText: "IsRecovery() bool  | 判断当前事件是否是恢复事件"},
        {text: "IsRecoverable()", displayText: "IsRecoverable() bool | 判断当前事件是否可恢复"},
        {text: "IsPlain()", displayText: "IsPlain() bool | 判断当前事件是否是普通事件"},
        {text: "FullJSON()", displayText: "FullJSON() string | 将整个事件编码为一个统一的 JSON 对象，返回字符串表示"},
    ],
    triggerMatchRules: [
        {text: "Events()", displayText: "Events() []repository.Message | 获取事件组中所有的 Events"},
        {text: "EventsMatchRegexCount(REGEX)", displayText: "MessagesMatchRegexCount(regex string) int64  | 获取匹配指定正则表达式的 Event 数量"},
        {text: "EventsWithMetaCount(KEY, VALUE)", displayText: "EventsWithMetaCount(key, value string) int64  | 获取 meta 匹配指定 key=value 的 Event 数量"},
        {text: "EventsWithTagsCount(TAG)", displayText: "EventsWithTagsCount(tags string) int64  | 获取拥有指定 tag 的 Event 数量，多个 tag 使用英文逗号分隔"},
        {text: "EventsCount()", displayText: "EventsCount() int64 | 获取事件组中 Events 数量"},
        {text: "UsersHasProperty(KEY, VALUE, RETURN)", displayText: "UsersHasProperty(key, value string, returnField string) []string | 根据字段、值查询用户，返回值指定字段列表"},
        {text: "TriggeredTimesInPeriod(PERIOD_IN_MINUTES, TRIGGER_STATUS)", displayText: "TriggeredTimesInPeriod(periodInMinutes int, triggerStatus string) int64 当前规则在指定时间范围内，状态为 triggerStatus 的触发次数"},
        {text: "LastTriggeredGroup(TRIGGER_STATUS)", displayText: "LastTriggeredGroup(triggerStatus string) repository.MessageGroup 最后一次触发该规则的状态为 triggerStatus 的事件组"},
        {text: "collecting", displayText: "collecting  | TriggerStatus：collecting"},
        {text: "pending", displayText: "pending | TriggerStatus：pending"},
        {text: "ok", displayText: "ok | TriggerStatus：ok"},
        {text: "failed", displayText: "failed | TriggerStatus：failed"},
        {text: "canceled", displayText: "canceled | TriggerStatus：canceled"},
        {text: "2006-01-02T15:04:05Z07:00", displayText: "2006-01-02T15:04:05Z07:00  | 时间格式"},
        {text: "RFC3339", displayText: "RFC3339  | 时间格式"},

        {text: 'Group', displayText: 'Group | 字段类型：MessageGroup | 所属对象：ROOT' },
        {text: 'Trigger', displayText: 'Trigger | 字段类型：Trigger | 所属对象：ROOT' },
        {text: 'Group.ID', displayText: 'Group.ID | 字段类型：ObjectID | 所属对象：MessageGroup' },
        {text: 'Group.SeqNum', displayText: 'Group.SeqNum | 字段类型：int64 | 所属对象：MessageGroup' },
        {text: 'Group.AggregateKey', displayText: 'Group.AggregateKey | 字段类型：string | 所属对象：MessageGroup' },
        {text: 'Group.MessageCount', displayText: 'Group.MessageCount | 字段类型：int64 | 所属对象：MessageGroup' },
        {text: 'Group.Rule', displayText: 'Group.Rule | 字段类型：MessageGroupRule | 所属对象：MessageGroup' },
        {text: 'Group.Actions', displayText: 'Group.Actions | 字段类型：[]Trigger | 所属对象：MessageGroup' },
        {text: 'Group.Status', displayText: 'Group.Status | 字段类型：string | 所属对象：MessageGroup' },
        {text: 'Group.CreatedAt', displayText: 'Group.CreatedAt | 字段类型：Time | 所属对象：MessageGroup' },
        {text: 'Group.UpdatedAt', displayText: 'Group.UpdatedAt | 字段类型：Time | 所属对象：MessageGroup' },
        {text: 'Group.Rule.ID', displayText: 'Group.Rule.ID | 字段类型：ObjectID | 所属对象：MessageGroupRule' },
        {text: 'Group.Rule.Name', displayText: 'Group.Rule.Name | 字段类型：string | 所属对象：MessageGroupRule' },
        {text: 'Group.Rule.AggregateKey', displayText: 'Group.Rule.AggregateKey | 字段类型：string | 所属对象：MessageGroupRule' },
        {text: 'Group.Rule.ExpectReadyAt', displayText: 'Group.Rule.ExpectReadyAt | 字段类型：Time | 所属对象：MessageGroupRule' },
        {text: 'Group.Rule.Rule', displayText: 'Group.Rule.Rule | 字段类型：string | 所属对象：MessageGroupRule' },
        {text: 'Group.Rule.Template', displayText: 'Group.Rule.Template | 字段类型：string | 所属对象：MessageGroupRule' },
        {text: 'Group.Rule.SummaryTemplate', displayText: 'Group.Rule.SummaryTemplate | 字段类型：string | 所属对象：MessageGroupRule' },
        {text: 'Trigger.ID', displayText: 'Trigger.ID | 字段类型：ObjectID | 所属对象：Trigger' },
        {text: 'Trigger.Name', displayText: 'Trigger.Name | 字段类型：string | 所属对象：Trigger' },
        {text: 'Trigger.PreCondition', displayText: 'Trigger.PreCondition | 字段类型：string | 所属对象：Trigger' },
        {text: 'Trigger.Action', displayText: 'Trigger.Action | 字段类型：string | 所属对象：Trigger' },
        {text: 'Trigger.Meta', displayText: 'Trigger.Meta | 字段类型：string | 所属对象：Trigger' },
        {text: 'Trigger.UserRefs', displayText: 'Trigger.UserRefs | 字段类型：[]ObjectID | 所属对象：Trigger' },
        {text: 'Trigger.Status', displayText: 'Trigger.Status | 字段类型：string | 所属对象：Trigger' },
        {text: 'Trigger.FailedCount', displayText: 'Trigger.FailedCount | 字段类型：int | 所属对象：Trigger' },
        {text: 'Trigger.FailedReason', displayText: 'Trigger.FailedReason | 字段类型：string | 所属对象：Trigger' },
    ],
    matchRules: [
        {text: "matches", displayText: "\"foo\" matches \"^b.+\" 正则匹配"},
        {text: "contains", displayText: "str contains \"xxx\" | 字符串包含"},
        {text: "startsWith", displayText: "str startsWith prefix | 前缀匹配"},
        {text: "endsWith", displayText: "str endsWith prefix | 后缀匹配"},
        {text: "in", displayText: "user.Group in [\"human_resources\", \"marketing\"] | 包含"},
        {text: "not in", displayText: "user.Group not in [\"human_resources\", \"marketing\"] | 不包含"},
        {text: "or", displayText: "or | 或者"},
        {text: "and", displayText: "and | 同时"},
        {text: "len", displayText: "len | length of array, map or string"},
        {text: "all", displayText: "all | will return true if all element satisfies the predicate"},
        {text: "none", displayText: "none | will return true if all element does NOT satisfies the predicate"},
        {text: "any", displayText: "any | will return true if any element satisfies the predicate"},
        {text: "any([], {Content contains #})", displayText: 'any(["aa", "bb", "cc"], {Content contains #}) | 判断 Content 中包含数组中的任意值'},
        {text: "one", displayText: "one | will return true if exactly ONE element satisfies the predicate"},
        {text: "filter", displayText: "filter | filter array by the predicate"},
        {text: "map", displayText: "map | map all items with the closure"},
        {text: "count", displayText: "count | returns number of elements what satisfies the predicate"},
    ],
    templates: [
        {text: '.Events EVENT_COUNT', displayText: 'Events(limit int64) []repository.Event | 从事件组中获取 EVENT_COUNT 个 Events'},
        {text: ".IsRecovery", displayText: "IsRecovery() bool  | 判断当前事件组中的事件是否是恢复事件"},
        {text: ".IsRecoverable", displayText: "IsRecoverable() bool | 判断当前事件组中的事件是否可恢复"},
        {text: ".IsPlain", displayText: "IsPlain() bool | 判断当前事件组中的事件是否是普通事件"},
        {text: ".EventType", displayText: "EventType() string | 判断当前事件组中的事件类型：recovery/plain/recoverable"},
        {text: ".FirstEvent", displayText: "FirstEvent() repository.Event | 从当前事件组中获取第一个事件"},
        {text: '{{ }}', displayText: '{{ }} |  Golang 代码块'},
        {text: '{{ range $i, $msg := ARRAY }}\n {{ $i }} {{ $msg }} \n{{ end }}', displayText: '{{ range }}  | Golang 遍历对象'},
        {text: '{{ range $i, $msg := .Messages 4 }} {{ end }}', displayText: '{{ range $i, $msg := .Messages 4 }} {{ end }} | Golang 遍历 Messages，只取 4 条作为摘要'},
        {text: '{{ if pipeline }}\n T1 \n{{ else if pipeline }}\n T2 \n{{ else }}\n T3 \n{{ end }}', displayText: '{{ if }} |  Golang 分支条件'},
        {text: '[]()', displayText: 'Markdown 连接地址'},
        {text: 'index MAP_VAR "KEY"', displayText: 'index $msg.Meta "message.message" | 从 Map 中获取某个 Key 的值'},
        {text: 'call FUNC ARGS...', displayText: 'call | Returns the result of calling the first argument, which must be a function, with the remaining arguments as parameters'},
        {text: 'html', displayText: 'html | Returns the escaped HTML equivalent of the textual representation of its arguments'},
        {text: 'js', displayText: 'js | Returns the escaped JavaScript equivalent of the textual representation of its arguments'},
        {text: 'len', displayText: 'len | Returns the integer length of its argument'},
        {text: 'urlquery', displayText: 'urlquery | Returns the escaped value of the textual representation of its arguments in a form suitable for embedding in a URL query'},
        {text: 'print', displayText: 'print | An alias for fmt.Sprint'},
        {text: 'printf', displayText: 'printf | An alias for fmt.Sprintf'},
        {text: 'println', displayText: 'println | An alias for fmt.Sprintln'},
        {text: 'cutoff MAX_LENGTH STR', displayText: 'cutoff(maxLen int, val string) string  |  字符串截断'},
        {text: 'cutoff_line MAX_LINES STR', displayText: 'cutoff_line(maxLines int, val string) string  |  字符串截断截取指定行数'},
        {text: 'line_filter_include FILTER STR', displayText: 'line_filter_include(filter string, val string) string  | 字符串按照行过滤，保留匹配的行'},
        {text: 'line_filter_exclude FILTER STR', displayText: 'line_filter_exclude(filter string, val string) string  | 字符串按照行过滤，去除匹配的行'},
        {text: 'implode ELEMENT_ARR ","', displayText: 'implode(elems []string, sep string) string  |  字符串数组拼接'},
        {text: 'join "," ELEMENT_ARR', displayText: 'join(sep string, elems []string) string  |  字符串数组拼接'},
        {text: 'explode STR ","', displayText: 'explode(s, sep string) []string  |  字符串分隔成数组'},
        {text: 'split "," STR', displayText: 'split(sep, s string) []string  |  字符串分隔成数组'},
        {text: 'ident "IDENT_STR" STR', displayText: 'ident(ident string, message string) string  |  多行字符串统一缩进'},
        {text: 'json JSONSTR', displayText: 'json(content string) string  |  JSON 字符串格式化'},
        {text: 'datetime LAYOUT DATETIME', displayText: 'datetime(layout string, datetime time.Time) string  |  时间格式化展示为 2006-01-02 15:04:05 格式，时区选择北京/重庆'},
        {text: 'datetime_noloc LAYOUT DATETIME',displayText: 'datetime_noloc(layout string, datetime time.Time) string  |  时间格式化展示为 2006-01-02 15:04:05 格式，默认时区'},
        {text: 'datetime_loc LAYOUT LOC DATETIME',displayText: 'datetime_loc(layout string, locName string, datetime time.Time) string  |  时间格式化展示为 2006-01-02 15:04:05 格式，指定时区，如 UTC'},
        {text: 'datetime_add_sec DATETIME OFFSET',displayText: 'datetime_add_sec(datetime time.Time, offset int) time.Time  | 对时间进行运算，offset 单位为秒'},
        {text: 'reformat_datetime_str ORIGINAL_LAYOUT TARGET_LAYOUT DATETIME_STR',displayText: 'reformat_datetime_str(originalLayout, targetLayout string, dt string) string  |  重新格式化时间字符串'},
        {text: 'parse_datetime_str LAYOUT DATETIME_STR',displayText: 'parse_datetime_str(layout string, dt string) time.Time  | 将时间字符串解析为时间对象'},
        {text: 'parse_datetime_str_rfc3339 DATETIME_STR',displayText: 'parse_datetime_str_rfc3339(dt string) time.Time  |  将时间字符串解析为时间对象，格式为 RFC3339'},
        {text: 'json_get "KEY" "DEFAULT" JSONSTR', displayText: 'json_get(key string, defaultValue string, body string) string  |  将 body 解析为 json，然后获取 key 的值，失败返回 defaultValue'},
        {text: 'json_gets "KEY" "DEFAULT" JSONSTR', displayText: 'json_gets(key string, defaultValue string, body string) string  |  将 body 解析为 json，然后获取 key 的值(可以使用逗号分割多个key作为备选)，失败返回 defaultValue'},
        {text: 'json_array "KEY" JSONSTR', displayText: 'json_array(key string, body string) []string  |  将 body 解析为 json，然后获取 key 的值（数组值）'},
        {text: 'json_flatten JSONSTR MAX_LEVEL', displayText: 'json_flatten(body string, maxLevel int) []jsonutils.KvPair  |  将 body 解析为 json，然后转换为键值对返回'},
        {text: 'format "FORMAT" VAL', displayText: 'format(format string, a ...interface{}) string  |  格式化展示，调用 fmt.Sprintf'},
        {text: 'number_beauty VAL', displayText: 'number_beauty(num interface{}) string  |  数字格式化展示，自动添加千分位分隔符'},
        {text: 'integer STR', displayText: 'integer(str string) int  |  字符串转整数 '},
        {text: 'float STR', displayText: 'float(str string) float64  |  字符串转浮点数 '},
        {text: 'mysql_slowlog STR', displayText: 'mysql_slowlog(slowlog string) map[string]string  |  解析 MySQL 慢查询日志为 map'},
        {text: 'sql_finger STR', displayText: 'sql_finger(str string) string | 将 SQL 转换为其指纹（替换参数为占位符）'},
        {text: 'open_falcon_im STR', displayText: 'open_falcon_im(msg string) OpenFalconIM  |  解析 OpenFalcon 事件格式'},
        {text: 'string_mask STR LEFT', displayText: 'string_mask(content string, left int) string  |  在左右两侧只保留 left 个字符，中间所有字符替换为 *'},
        {text: 'string_tags TAG_STR SEPARATOR', displayText: 'string_tags(tags string, sep string) []string  |  将字符串 tags 用 sep 作为分隔符，切割成多个 tag，空的 tag 会被排除'},
        {text: 'remove_empty_line STR', displayText: 'remove_empty_line(content string) string | 移除字符串中的空行'},
        {text: 'serialize VAL', displayText: 'serialize(data interface{}) string | 对象序列化为字符串，用于展示'},
        {text: 'json_encode VAL', displayText: 'json_encode(data interface{}) string | 对象序列化为字符串，用于展示'},
        {text: 'error_notice MSG', displayText: 'error_notice(msg string) string | 红色字体显示"msg"两字'},
        {text: 'success_notice MSG', displayText: 'success_notice(msg string) string | 绿色字体显示"msg"两字'},
        {text: 'error_success_notice IS_SUCCESS MSG', displayText: 'error_success_notice(success bool, msg string) string | 显示"msg"两字，如果 success 为 true，显示绿色，否则显示红色'},
        {text: 'condition S1 S2 CONDITION', displayText: 'condition(s1, s2 string, condition bool) | 条件输出字符串，符合条件，输出 s1，否则 s2'},
        {text: 'recoverable_notice .IsRecovery MSG', displayText: 'recoverable_notice(recovered bool, msg string) string | 显示"msg"两字，如果 recovered 为 true，显示绿色，并且自动添加 【已恢复】两字，否则显示红色'},

        {text: 'user_metas QUERY_K QUERY_V FIELD', displayText: 'user_metas(queryK, queryV string, field string) []string | 查询 queryK=queryV 的用户 field 元信息，查询结果是个字符串数组'},
        {text: 'events_relation_ids EVENTS', displayText: 'events_relation_ids(events []repository.Event) []primitive.ObjectID | 从多个事件中提取包含的事件关联 ID'},
        {text: 'events_relations RELATION_IDS', displayText: 'events_relations(relationIDs []primitive.ObjectID) []repository.EventRelation | 根据多个事件关联 ID 批量查询事件关联'},
        {text: 'event_relation_notes RELATION_ID', displayText: 'event_relation_notes(relationID primitive.ObjectID) []repository.EventRelationNote | 根据事件关联 ID 查询事件相关的备注'},

        {text: 'prefix_all_str PREFIX ARR', displayText: 'prefix_all_str(prefix string, arr []string) []string | 为字符串数组中每一个元素添加前缀'},
        {text: 'suffix_all_str SUFFIX ARR', displayText: 'suffix_all_str(prefix string, arr []string) []string | 为字符串数组中每一个元素添加后缀'},
        {text: 'json_fields_cutoff LENGTH JSON_STR', displayText: 'json_fields_cutoff(length int, body string) map[string]interface{} | 对 JSON 字符串扁平化，然后对每个 KV 截取指定长度'},
        {text: 'map_fields_cutoff LENGTH MAP', displayText: 'map_fields_cutoff(length int, source map[string]interface{}) map[string]interface{} | 对 Map 的每个 KV 截取指定长度'},
        {text: 'trim_prefix_map_k PREFIX SOURCE', displayText: 'trim_prefix_map_k(prefix string, source map[string]interface{}) map[string]interface{} | 移除 Map 中所有 Key 的前缀'},

        {text: 'meta_filter STR FILTER_STR', displayText: 'meta_filter(meta map[string]interface{}, allowKeys ...string) map[string]interface{} | 过滤Meta，只保留允许的Key'},
        {text: 'meta_filter_exclude STR FILTER_STR', displayText: 'meta_filter_exclude(meta map[string]interface{}, disableKeys ...string) map[string]interface{} | 过滤Meta，移除不允许的Key'},
        {text: 'meta_prefix_filter STR FILTER_PREFIX', displayText: 'meta_prefix_filter(meta map[string]interface{}, allowPrefixes ...string) map[string]interface{} | 过滤Meta，只保留包含指定 prefix 的Key'},
        {text: 'meta_prefix_filter_exclude STR FILTER_PREFIX', displayText: 'meta_prefix_filter_exclude(meta map[string]interface{}, disablePrefixes ...string) map[string]interface{} | 过滤Meta，移除匹配 prefix 的 Key'},

        {text: 'starts_with STR "START_STR"', displayText: 'starts_with(haystack string, needles ...string) bool  |  判断 haystack 是否以 needles 开头'},
        {text: 'ends_with STR "START_END"', displayText: 'ends_with(haystack string, needles ...string) bool  |  判断 haystack 是否以 needles 结尾'},
        {text: 'trim STR "CUTSTR"', displayText: 'trim(s string, cutset string) string  |  去掉字符串 s 两边的 cutset 字符'},
        {text: 'trim_left STR "CUTSTR"', displayText: 'trim_left(s string, cutset string) string  |  去掉字符串 s 左侧的 cutset 字符'},
        {text: 'trim_right STR "CUTSTR"', displayText: 'trim_right(s string, cutset string) string  |  去掉字符串 s 右侧的 cutset 字符'},
        {text: 'trim_space STR', displayText: 'trim_space(s string) string  |  去掉字符串 s 两边的空格'},

        {text: 'str_upper STR', displayText: 'str_upper(s string) string  |  字符串转大写'},
        {text: 'str_lower STR', displayText: 'str_lower(s string) string  |  字符串转小写'},
        {text: 'str_replace STR OLD NEW', displayText: 'str_replace(s string, old string, new string) string  |  字符串替换， 将 s 中所有的 old 替换为 new'},
        {text: 'str_repeat STR COUNT', displayText: 'str_repeat(s string, count int) string  | 字符串 s 重复 count 次'},
        {text: 'str_concat STR1 STR2', displayText: 'str_concat(s ...string) string  | 多个字符串拼接'},

        {text: 'html2md HTML', displayText: 'html2md(html string) string  | 将 HTML 转换为 Markdown'},
        {text: 'md2html MARKDOWN', displayText: 'md2html(markdown string) string  | 将 Markdown 转换为 HTML'},
        {text: 'md2confluence MARKDOWN', displayText: 'md2confluence(markdown string) string  | 将 Markdown 转换为 Confluence (Wiki/Jira 等) 的富文本格式'},
        {text: 'html_beauty HTML', displayText: 'html_beauty(html string) string  | HTML 格式化'},
        {text: 'dom_filter_html selector STR', displayText: 'dom_filter_html(selector string, str string) []string  | 从 HTML DOM 提取匹配 selector 选择器的内容，以字符串数组形式返回'},
        {text: 'dom_filter_html_n selector N STR', displayText: 'dom_filter_html_n(selector string, n int, str string) string  | 从 HTML DOM 提取匹配 selector 选择器的内容，返回第 n 个（n 从 0 开始）'},

        {text: 'md5 DATA', displayText: 'md5(data interface{}) string  | 生成 data 的 md5 值'},
        {text: 'sha1 DATA', displayText: 'sha1(data interface{}) string  | 生成 data 的 sha1 值'},
        {text: 'base64 DATA', displayText: 'base64(data interface{}) string  | 生成 data 的 base64 编码值'},
        {text: 'base64_encode DATA', displayText: 'base64_encode(data interface{}) string  | 生成 data 的 base64 编码值'},

        {text: '.Action', displayText: '.Action | 字段类型：string | 所属对象：ROOT' },
        {text: '.RuleTemplateParsed', displayText: '.RuleTemplateParsed | 字段类型：string | 所属对象：ROOT' },
        {text: '.PreviewURL', displayText: '.PreviewURL | 字段类型：string | 所属对象：ROOT' },
        {text: '.ReportURL', displayText: '.ReportURL | 字段类型：string | 所属对象：ROOT' },
        {text: '.Trigger', displayText: '.Trigger | 字段类型：Trigger | 所属对象：ROOT' },
        {text: '.Group', displayText: '.Group | 字段类型：MessageGroup | 所属对象：ROOT' },
        {text: '.Rule', displayText: '.Rule | 字段类型：Rule | 所属对象：ROOT' },
        {text: '.Rule.ID', displayText: '.Rule.ID | 字段类型：ObjectID | 所属对象：Rule' },
        {text: '.Rule.Name', displayText: '.Rule.Name | 字段类型：string | 所属对象：Rule' },
        {text: '.Rule.Description', displayText: '.Rule.Description | 字段类型：string | 所属对象：Rule' },
        {text: '.Rule.Tags', displayText: '.Rule.Tags | 字段类型：[]string | 所属对象：Rule' },
        {text: '.Rule.AggregateRule', displayText: '.Rule.AggregateRule | 字段类型：string | 所属对象：Rule' },
        {text: '.Rule.ReadyType', displayText: '.Rule.ReadyType | 字段类型：string | 所属对象：Rule' },
        {text: '.Rule.Interval', displayText: '.Rule.Interval | 字段类型：int64 | 所属对象：Rule' },
        {text: '.Rule.DailyTimes', displayText: '.Rule.DailyTimes | 字段类型：[]string | 所属对象：Rule' },
        {text: '.Rule.Rule', displayText: '.Rule.Rule | 字段类型：string | 所属对象：Rule' },
        {text: '.Rule.Template', displayText: '.Rule.Template | 字段类型：string | 所属对象：Rule' },
        {text: '.Rule.SummaryTemplate', displayText: '.Rule.SummaryTemplate | 字段类型：string | 所属对象：Rule' },
        {text: '.Rule.CreatedAt', displayText: '.Rule.CreatedAt | 字段类型：Time | 所属对象：Rule' },
        {text: '.Rule.UpdatedAt', displayText: '.Rule.UpdatedAt | 字段类型：Time | 所属对象：Rule' },
        {text: '.Rule.Triggers', displayText: '.Rule.Triggers | 字段类型：[]Trigger | 所属对象：Rule' },
        {text: '.Rule.Status', displayText: '.Rule.Status | 字段类型：string | 所属对象：Rule' },
        {text: '.Rule.TimeRanges', displayText: '.Rule.TimeRanges | 字段类型：[]TimeRange | 所属对象：Rule' },
        {text: '$timeRange.StartTime', displayText: '$timeRange.StartTime | 字段类型：string | 所属对象：TimeRange' },
        {text: '$timeRange.EndTime', displayText: '$timeRange.EndTime | 字段类型：string | 所属对象：TimeRange' },
        {text: '$timeRange.Interval', displayText: '$timeRange.Interval | 字段类型：int64 | 所属对象：TimeRange' },
        {text: '.Trigger.ID', displayText: '.Trigger.ID | 字段类型：ObjectID | 所属对象：Trigger' },
        {text: '.Trigger.Name', displayText: '.Trigger.Name | 字段类型：string | 所属对象：Trigger' },
        {text: '.Trigger.PreCondition', displayText: '.Trigger.PreCondition | 字段类型：string | 所属对象：Trigger' },
        {text: '.Trigger.Action', displayText: '.Trigger.Action | 字段类型：string | 所属对象：Trigger' },
        {text: '.Trigger.Meta', displayText: '.Trigger.Meta | 字段类型：string | 所属对象：Trigger' },
        {text: '.Trigger.UserRefs', displayText: '.Trigger.UserRefs | 字段类型：[]ObjectID | 所属对象：Trigger' },
        {text: '.Trigger.Status', displayText: '.Trigger.Status | 字段类型：string | 所属对象：Trigger' },
        {text: '.Trigger.FailedCount', displayText: '.Trigger.FailedCount | 字段类型：int | 所属对象：Trigger' },
        {text: '.Trigger.FailedReason', displayText: '.Trigger.FailedReason | 字段类型：string | 所属对象：Trigger' },
        {text: '$trigger.ID', displayText: '$trigger.ID | 字段类型：ObjectID | 所属对象：Trigger' },
        {text: '$trigger.Name', displayText: '$trigger.Name | 字段类型：string | 所属对象：Trigger' },
        {text: '$trigger.PreCondition', displayText: '$trigger.PreCondition | 字段类型：string | 所属对象：Trigger' },
        {text: '$trigger.Action', displayText: '$trigger.Action | 字段类型：string | 所属对象：Trigger' },
        {text: '$trigger.Meta', displayText: '$trigger.Meta | 字段类型：string | 所属对象：Trigger' },
        {text: '$trigger.UserRefs', displayText: '$trigger.UserRefs | 字段类型：[]ObjectID | 所属对象：Trigger' },
        {text: '$trigger.Status', displayText: '$trigger.Status | 字段类型：string | 所属对象：Trigger' },
        {text: '$trigger.FailedCount', displayText: '$trigger.FailedCount | 字段类型：int | 所属对象：Trigger' },
        {text: '$trigger.FailedReason', displayText: '$trigger.FailedReason | 字段类型：string | 所属对象：Trigger' },
        {text: '$action.ID', displayText: '$action.ID | 字段类型：ObjectID | 所属对象：Trigger' },
        {text: '$action.Name', displayText: '$action.Name | 字段类型：string | 所属对象：Trigger' },
        {text: '$action.PreCondition', displayText: '$action.PreCondition | 字段类型：string | 所属对象：Trigger' },
        {text: '$action.Action', displayText: '$action.Action | 字段类型：string | 所属对象：Trigger' },
        {text: '$action.Meta', displayText: '$action.Meta | 字段类型：string | 所属对象：Trigger' },
        {text: '$action.UserRefs', displayText: '$action.UserRefs | 字段类型：[]ObjectID | 所属对象：Trigger' },
        {text: '$action.Status', displayText: '$action.Status | 字段类型：string | 所属对象：Trigger' },
        {text: '$action.FailedCount', displayText: '$action.FailedCount | 字段类型：int | 所属对象：Trigger' },
        {text: '$action.FailedReason', displayText: '$action.FailedReason | 字段类型：string | 所属对象：Trigger' },
        {text: '.Group.ID', displayText: '.Group.ID | 字段类型：ObjectID | 所属对象：MessageGroup' },
        {text: '.Group.SeqNum', displayText: '.Group.SeqNum | 字段类型：int64 | 所属对象：MessageGroup' },
        {text: '.Group.AggregateKey', displayText: '.Group.AggregateKey | 字段类型：string | 所属对象：MessageGroup' },
        {text: '.Group.MessageCount', displayText: '.Group.MessageCount | 字段类型：int64 | 所属对象：MessageGroup' },
        {text: '.Group.Rule', displayText: '.Group.Rule | 字段类型：MessageGroupRule | 所属对象：MessageGroup' },
        {text: '.Group.Actions', displayText: '.Group.Actions | 字段类型：[]Trigger | 所属对象：MessageGroup' },
        {text: '.Group.Status', displayText: '.Group.Status | 字段类型：string | 所属对象：MessageGroup' },
        {text: '.Group.CreatedAt', displayText: '.Group.CreatedAt | 字段类型：Time | 所属对象：MessageGroup' },
        {text: '.Group.UpdatedAt', displayText: '.Group.UpdatedAt | 字段类型：Time | 所属对象：MessageGroup' },
        {text: '.Group.Rule.ID', displayText: '.Group.Rule.ID | 字段类型：ObjectID | 所属对象：MessageGroupRule' },
        {text: '.Group.Rule.Name', displayText: '.Group.Rule.Name | 字段类型：string | 所属对象：MessageGroupRule' },
        {text: '.Group.Rule.AggregateKey', displayText: '.Group.Rule.AggregateKey | 字段类型：string | 所属对象：MessageGroupRule' },
        {text: '.Group.Rule.ExpectReadyAt', displayText: '.Group.Rule.ExpectReadyAt | 字段类型：Time | 所属对象：MessageGroupRule' },
        {text: '.Group.Rule.Rule', displayText: '.Group.Rule.Rule | 字段类型：string | 所属对象：MessageGroupRule' },
        {text: '.Group.Rule.Template', displayText: '.Group.Rule.Template | 字段类型：string | 所属对象：MessageGroupRule' },
        {text: '.Group.Rule.SummaryTemplate', displayText: '.Group.Rule.SummaryTemplate | 字段类型：string | 所属对象：MessageGroupRule' },
        {text: '$msg.ID', displayText: '$msg.ID | 字段类型：ObjectID | 所属对象：Message' },
        {text: '$msg.SeqNum', displayText: '$msg.SeqNum | 字段类型：int64 | 所属对象：Message' },
        {text: '$msg.Content', displayText: '$msg.Content | 字段类型：string | 所属对象：Message' },
        {text: '$msg.Meta', displayText: '$msg.Meta | 字段类型：map[string]interface{} | 所属对象：Message' },
        {text: '$msg.Tags', displayText: '$msg.Tags | 字段类型：[]string | 所属对象：Message' },
        {text: '$msg.Origin', displayText: '$msg.Origin | 字段类型：string | 所属对象：Message' },
        {text: '$msg.GroupID', displayText: '$msg.GroupID | 字段类型：[]ObjectID | 所属对象：Message' },
        {text: '$msg.Status', displayText: '$msg.Status | 字段类型：string | 所属对象：Message' },
        {text: '$msg.CreatedAt', displayText: '$msg.CreatedAt | 字段类型：Time | 所属对象：Message' },
    ],
    triggerTemplates: [

    ],
}

helpers.matchRules.push(...helpers.helpers.map(item => {
    return {
        text: item.text + "(" + item.args.join(", ") + ")",
        displayText: item.displayText,
    }
}));
helpers.templates.push(...helpers.helpers.map(item => {
    return {
        text: 'helpers.' + item.text + " " + item.args.join(" "),
        displayText: 'helpers.' + item.displayText,
    };
}))

let hintHandler = function (editor) {
    let sources = [];
    switch (editor.options.hintOptions.adanosType) {
        case 'GroupMatchRule':
            sources.push(...helpers.groupMatchRules);
            sources.push(...helpers.matchRules);
            break;
        case 'TriggerMatchRule':
            sources.push(...helpers.triggerMatchRules);
            sources.push(...helpers.matchRules);
            break;
        case 'Template':
            sources.push(...helpers.templates);
            break;
        case 'DingTemplate':
            sources.push(...helpers.templates);
            sources.push(...helpers.triggerTemplates);
            break;
        case 'AllMatchRule':
            sources.push(...helpers.groupMatchRules);
            sources.push(...helpers.triggerMatchRules);
            sources.push(...helpers.matchRules);
            break;
        default:
    }

    let cur = editor.getCursor();
    let token = editor.getTokenAt(cur), start, end, search;
    if (token.end > cur.ch) {
        token.end = cur.ch;
        token.string = token.string.slice(0, cur.ch - token.start);
    }

    if (token.string.match(/^[.`"\w@][.\w$#]*$/g)) {
        search = token.string;
        start = token.start;
        end = token.end;
    } else {
        start = end = cur.ch;
        search = "";
    }

    search = search.toLowerCase()

    let list = [];
    if (search.trim() === '') {
        list = sources;
    } else {
        if (search.charAt(0) === '"' || search.charAt(0) === '.' || search.charAt(0) === "'") {
            search = search.substring(1);
        }

        for (let s in sources) {
            let str = sources[s];
            if (typeof str !== "string") {
                str = str.text;
            }
            if (str.toLowerCase().indexOf(search) >= 0) {
                list.push(sources[s]);
            }
        }
    }

    return {list: list, from: CodeMirror.Pos(cur.line, start), to: CodeMirror.Pos(cur.line, end)};
};

export {helpers, hintHandler}