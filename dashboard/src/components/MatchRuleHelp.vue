<template>
    <b-card-body class="adanos-help">
        <li><code>JsonGet(key string, defaultValue string) string</code> 将消息体作为json解析，获取指定的key</li>
        <li><code>Upper(val string) string</code> 字符串转大写</li>
        <li><code>Lower(val string) string</code> 字符串转小写</li>
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