<template>
    <b-card-body class="adanos-help">
        <li>支持的基本字段：
            <pre>
        Content string
        Meta    map[string]interface{}
        Tags    []string
        Origin  string</pre>
        </li>
        <li>支持的对象、函数：
            <ul>
                <li v-for="(helper, i) in helpers" v-bind:key="i"><code>{{ helper.text }}</code> {{ helper.displayText }}</li>
            </ul>
        </li>
    </b-card-body>
</template>

<script>
    export default {
        name: "MatchRuleHelp",
        props: {
            helpers: Array,
        }
    }
</script>

<style scoped>
.adanos-help {
    font-size: 80%;
}
</style>