<template>
    <b-card-body class="adanos-help">
        <li><code>JsonGet(key string, defaultValue string) string</code> 将消息体作为json解析，获取指定的key</li>
        <li><code>Upper(val string) string</code> 字符串转大写</li>
        <li><code>Lower(val string) string</code> 字符串转小写</li>
        <li><code>Now() time.Time</code> 当前时间</li>
        <li><code>ParseTime(layout string, value string) time.Time</code> 时间字符串转时间对象</li>
        <li><code>DailyTimeBetween(startTime, endTime string) bool</code> 判断当前时间是否在 startTime 和 endTime 之间（每天），时间格式为 15:04</li>
    </b-card-body>
</template>

<script>
    export default {
        name: "MatchRuleHelp"
    }
</script>

<style scoped>
.adanos-help {
    font-size: 80%;
}
</style>