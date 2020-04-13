<template>
    <b-card>
        <b-card-text class="adanos-help">
            <ul>
                <li><code>Upper(val string) string</code> 字符串转大写</li>
                <li><code>Lower(val string) string</code> 字符串转小写</li>
                <li><code>Now() time.Time</code> 当前时间</li>
                <li><code>ParseTime(layout string, value string) time.Time</code> 时间字符串转时间对象</li>
                <li><code>Messages() []repository.Message</code> 获取分组中所有的 Messages</li>
                <li><code>MessagesMatchRegexCount(regex string) int64</code> 获取匹配指定正则表达式的 message 数量</li>
                <li><code>MessagesWithMetaCount(key, value string) int64</code> 获取 meta 匹配指定 key=value 的 message 数量</li>
                <li><code>MessagesWithTagsCount(tags string) int64</code> 获取拥有指定 tag 的 message 数量，多个 tag 使用英文逗号分隔</li>
                <li><code>DailyTimeBetween(startTime, endTime string) bool</code> 判断当前时间是否在 startTime 和 endTime 之间（每天），时间格式为 15:04</li>
                <li><code>TriggeredTimesInPeriod(periodInMinutes int, triggerStatus string) int64</code> 当前规则在指定时间范围内，状态为 triggerStatus 的触发次数</li>
                <li><code>LastTriggeredGroup(triggerStatus string) repository.MessageGroup</code> 最后一次触发该规则的状态为 triggerStatus 的分组</li>
            </ul>
            <hr />
            <ol>
                <li><b>triggerStatus</b> 状态可选值为 <code>collecting</code>, <code>pending</code>, <code>ok</code>, <code>failed</code>, <code>canceled</code></li>
                <li>时间格式 layout 如 <code>2006-01-02T15:04:05Z07:00</code> 代表了 <code>RFC3339</code></li>
            </ol>
        </b-card-text>
    </b-card>
</template>

<script>
    export default {
        name: "TriggerHelp"
    }
</script>

<style scoped>

.adanos-help {
    font-size: 80%;
}

</style>
