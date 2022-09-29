<template>
    <b-card-body class="adanos-help">
        <li>支持的基本字段：
            <pre>
        CreatedAt  time.Time 
        Content    string
        Meta       map[string]interface{}
        Tags       []string
        Origin     string</pre>
        </li>
        <li>支持的对象、函数：
            <ul>
                <li v-for="(helper, i) in helpers" v-bind:key="i">
                    <code>{{ helper.text }}</code>
                    <p>{{ helper.displayText }}</p>
                </li>
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
    max-height: 400px;
    overflow-y: auto;
    border: 1px solid #ccc;
    border-radius: 4px;
    padding-top: 12px;
}
</style>