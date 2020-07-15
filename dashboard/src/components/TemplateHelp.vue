<template>
    <b-card-body>
        <b-card-text class="adanos-help">
            <ul>
                <li>支持的基本字段：
                    <pre><code>
Action:  string
Rule: {
	ID          primitive.ObjectID
	Name        string
	Description string
	Tags        []string

	Interval int64

	Rule            string
	Template        string
	SummaryTemplate string
	Triggers        []Trigger

	Status string

	CreatedAt time.Time
	UpdatedAt time.Time
}
Trigger: {
	ID           primitive.ObjectID
	Name         string
	PreCondition string
	Action       string
	Meta         string
	UserRefs     []primitive.ObjectID
	Status       string
	FailedCount  int
	FailedReason string
},
Group: {
	ID     primitive.ObjectID
	SeqNum int64
	MessageCount int64
	Rule         {
		ID              primitive.ObjectID
		Name            string
		ExpectReadyAt   time.Time
		Rule            string
		Template        string
		SummaryTemplate string
	}
	Actions      []Trigger // 分组关联的所有动作

	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
                    </code></pre>
                </li>
                <li>支持的函数：
                    <ul>
                        <li><code>cut_off(maxLen int, val string) string</code> 字符串截断</li>
                        <li><code>implode(elems []string, sep string) string</code> 字符串数组拼接</li>
                        <li><code>explode(s, sep string) []string</code> 字符串分隔成数组</li>
                        <li><code>ident(ident string, message string) string</code> 多行字符串统一缩进</li>
                        <li><code>json(content string) string</code> JSON 字符串格式化</li>
                        <li><code>datetime(datetime time.Time) string</code> 时间格式化展示为 2006-01-02 15:04:05 格式，时区选择北京/重庆</li>
                        <li><code>datetime_noloc(datetime time.Time) string</code> 时间格式化展示为 2006-01-02 15:04:05 格式，默认时区</li>
                        <li><code>json_get(key string, defaultValue string, body string) string</code> 将 body 解析为 json，然后获取 key 的值，失败返回 defaultValue</li>
                        <li><code>json_gets(key string, defaultValue string, body string) string</code> 将 body 解析为 json，然后获取 key 的值(可以使用逗号分割多个key作为备选)，失败返回 defaultValue</li>
                        <li><code>json_array(key string, body string) []string</code> 将 body 解析为 json，然后获取 key 的值（数组值）</li>
                        <li><code>json_flatten(body string, maxLevel int) []jsonutils.KvPair</code> 将 body 解析为 json，然后转换为键值对返回</li>
                        <li><code>starts_with(haystack string, needles ...string) bool</code> 判断 haystack 是否以 needles 开头</li>
                        <li><code>ends_with(haystack string, needles ...string) bool</code> 判断 haystack 是否以 needles 结尾</li>
                        <li><code>trim(s string, cutset string) string</code> 去掉字符串 s 两边的 cutset 字符</li>
                        <li><code>trim_left(s string, cutset string) string</code> 去掉字符串 s 左侧的 cutset 字符</li>
                        <li><code>trim_right(s string, cutset string) string</code> 去掉字符串 s 右侧的 cutset 字符</li>
                        <li><code>trim_space(s string) string</code> 去掉字符串 s 两边的空格</li>
                        <li><code>format(format string, a ...interface{}) string</code> 格式化展示，调用 fmt.Sprintf </li>
                        <li><code>integer(str string) int</code> 字符串转整数 </li>
                        <li><code>mysql_slowlog(slowlog string) map[string]string</code> 解析 MySQL 慢查询日志为 map </li>
                        <li><code>open_falcon_im(msg string) OpenFalconIM</code> 解析 OpenFalcon 消息格式 </li>
                        <li><code>string_mask(content string, left int) string</code> 在左右两侧只保留 left 个字符，中间所有字符替换为 * </li>
                        <li><code>string_tags(tags string, sep string) []string</code> 将字符串 tags 用 sep 作为分隔符，切割成多个 tag，空的 tag 会被排除 </li>
                        <li><code>remove_empty_line(content string) string</code> 移除字符串中的空行</li>
                        <li><code>meta_filter(meta map[string]interface{}, allowKeys ...string) map[string]interface{}</code> 过滤Meta，只保留允许的Key</li>
                        <li><code>meta_prefix_filter(meta map[string]interface{}, allowPrefix ...string) map[string]interface{}</code> 过滤Meta，只保留包含指定 prefix 的Key</li>
                    </ul>
                </li>
            </ul>
            <hr />
            <ol>
                <li>时间格式 layout 如 <code>2006-01-02T15:04:05Z07:00</code> 代表了 <code>RFC3339</code></li>
                <li>OpenFalconIM 格式为 
                    <pre>
<code>type OpenFalconIM struct {
	Priority    int
	Status      string
	Endpoint    string
	Body        string
	CurrentStep int
	FormatTime  string
}
</code>
                    </pre>
                </li>
                <li>jsonutils.KvPair 格式为
                    <pre>
<code>
type KvPair struct {
	Key   string
	Value string
}
</code>
                    </pre>
                </li>
            </ol>
        </b-card-text>
    </b-card-body>
</template>

<script>
    export default {
        name: "TemplateHelp"
    }
</script>

<style scoped>

.adanos-help {
    font-size: 80%;
}

</style>