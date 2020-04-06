<template>
    <b-card :header-bg-variant="message.status === 'canceled' ? 'warning': ''">
        <template v-slot:header>
            <b title="创建时间"><date-time :value="message.created_at"></date-time></b>
            <div class="float-right" title="状态">
                <b-badge v-if="message.status === 'pending'" variant="dark">准备中</b-badge>
                <b-badge v-if="message.status === 'grouped'" variant="success">已分组</b-badge>
                <b-badge v-if="message.status === 'canceled'" variant="danger">已取消</b-badge>
            </div>
        </template>

        <b-card-text>
            <b-row style="max-width: 100rem;" class="adanos-meta-line">
                <b-col sm="3"><b class="text-black-50" style="border-bottom: 1px dashed black">标签</b></b-col>
                <b-col sm="9" style="text-align: left"><b-badge v-for="(tag, index) in message.tags" :key="index" class="mr-1">{{ tag }}</b-badge></b-col>
            </b-row>
            <b-row style="max-width: 100rem;" class="mb-2 adanos-meta-line">
                <b-col sm="3"><b class="text-black-50" style="border-bottom: 1px dashed black">来源</b></b-col>
                <b-col sm="9"><b-badge variant="light">{{ message.origin }}</b-badge></b-col>
            </b-row>
            <b-row v-for="(val, key) in message.meta" :key="key" style="max-width: 100rem;" class="adanos-meta-line">
                <b-col sm="3"><b class="text-dark" style="border-bottom: 1px dashed black">{{ key }}</b></b-col>
                <b-col sm="9"><b-badge variant="light">{{ val }}</b-badge></b-col>
            </b-row>
        </b-card-text>
        <b-card-text>
            <code><pre class="adanos-code">{{ message.content }}</pre></code>
        </b-card-text>
    </b-card>
</template>

<script>
    export default {
        name: 'MessageCard',
        props: {
            message: Object
        },
        methods: {
            
        }
    }
</script>

<style scoped>
    .adanos-code {
        max-height: 100rem;
        white-space: pre-wrap!important;
        word-wrap: break-word!important;
        *white-space:normal!important;
    }

    .adanos-meta-line {
        border-bottom: 1px solid #ffffff;
    }
    .adanos-meta-line:hover {
        border-bottom: 1px dashed #ccc;
    }
</style>