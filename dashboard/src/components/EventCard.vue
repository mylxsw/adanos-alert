<template>
    <b-card :header-bg-variant="event.status === 'canceled' ? 'warning': ''">
        <template v-slot:header>
            <b title="创建时间" v-if="!title">
                <b-badge v-if="event_index != null" class="mr-2" variant="primary"># {{ event.seq_num }}</b-badge>
                <date-time :value="event.created_at"></date-time>
            </b>
            <span v-if="title">{{ title }}</span>
            <div class="float-right" title="状态">
                <b-badge v-if="event.status === 'pending'" variant="dark">准备中</b-badge>
                <b-badge v-if="event.status === 'grouped'" variant="success">已分组</b-badge>
                <b-badge v-if="event.status === 'canceled'" variant="danger">无规则，已取消</b-badge>
                <b-badge v-if="event.status === 'expired'" variant="warning">匹配规则，已过期</b-badge>
                <b-badge v-if="event.status === 'ignored'" variant="warning">匹配规则，已忽略</b-badge>

                <b-link class="ml-2" @click="isFold = !isFold">
                    <b-icon icon="arrows-collapse" v-if="isFold"></b-icon>
                    <b-icon icon="arrows-expand" v-if="!isFold"></b-icon>
                </b-link>
            </div>

        </template>
        <b-card-text v-if="isFold">...</b-card-text>

        <template v-slot:footer v-if="!onlyShow">
            <div class="float-right">
                <b-button size="sm" class="ml-2" variant="warning" @click="testMatchedRules(event.id)" v-if="testMatchedRules">测试</b-button>
                <b-button size="sm" class="ml-2" variant="success" @click="reproduceEvent(event.id)" v-if="reproduceEvent">重发</b-button>
                <b-button :to="{path:'/rules/add', query: {test_event_id: event.id}}" target="_blank" class="ml-2" size="sm" variant="primary">
                    创建规则
                </b-button>
                <b-button size="sm" class="ml-2" :to="{path:'/debug', query: {event_id: event.id}}" target="_blank">Debug</b-button>
            </div>
        </template>

        <b-card-text v-if="!isFold">
            <b-row style="max-width: 100rem;" class="adanos-meta-line" v-if="event.group_ids != null && event.group_ids.length > 0">
                <b-col sm="3"><b class="text-black-50" style="border-bottom: 1px dashed black">事件组 ID</b></b-col>
                <b-col sm="9" style="text-align: left">
                    <b-link v-for="(g, index) in event.group_ids" :key="index" class="mr-1" :to="{path:'/events', query: {group_id: g}}">{{ g }}</b-link>
                </b-col>
            </b-row>
            <b-row style="max-width: 100rem;" class="adanos-meta-line" v-if="event.relation_ids != null && event.relation_ids.length > 0">
                <b-col sm="3"><b class="text-black-50" style="border-bottom: 1px dashed black">相关事件</b></b-col>
                <b-col sm="9" style="text-align: left">
                    <b-link v-for="(g, index) in event.relation_ids" :key="index" class="mr-1" :to="{path:'/events', query: {relation_id: g}}">{{ g }}</b-link>
                </b-col>
            </b-row>
            <b-row style="max-width: 100rem;" class="adanos-meta-line">
                <b-col sm="3"><b class="text-black-50" style="border-bottom: 1px dashed black">Tags</b></b-col>
                <b-col sm="9" style="text-align: left"><b-badge v-for="(tag, index) in event.tags" :key="index" class="mr-1">{{ tag }}</b-badge></b-col>
            </b-row>
            <b-row style="max-width: 100rem;" class="mb-2 adanos-meta-line">
                <b-col sm="3"><b class="text-black-50" style="border-bottom: 1px dashed black">Origin</b></b-col>
                <b-col sm="9"><b-badge variant="light">{{ event.origin }}</b-badge></b-col>
            </b-row>
            <b-row v-for="(val, key) in event.meta" :key="key" style="max-width: 100rem;" class="adanos-meta-line" title="Meta" v-b-tooltip>
                <b-col sm="3"><b class="text-dark" style="border-bottom: 1px dashed black">{{ key }}</b></b-col>
                <b-col sm="9"><pre class="adanos-code"><code>{{ val }}</code></pre></b-col>
            </b-row>
        </b-card-text>
        <b-card-text v-if="!isFold">
            <code><pre class="adanos-code" v-b-tooltip title="Content">{{ event.content }}</pre></code>
        </b-card-text>
    </b-card>
</template>

<script>
    export default {
        name: 'EventCard',
        props: {
            event: Object,
            event_index: Number,
            testMatchedRules: Function,
            reproduceEvent: Function,
            onlyShow: Boolean,
            title: String,
            fold: Boolean,
        },
        data() {
            return {
                isFold: false,
            }
        },
        methods: {
        },
        mounted() {
            this.isFold = this.fold;
        }
    }
</script>

<style scoped>
    .adanos-code {
        max-height: 50rem;
        white-space: pre-wrap!important;
        word-wrap: break-word!important;
        *white-space:normal!important;

        margin-bottom: 0;
        font-size: 80%;
    }

    .adanos-meta-line {
        border-bottom: 1px solid #ffffff;
    }
    .adanos-meta-line:hover  {
        background: #ffc107;
    }
</style>