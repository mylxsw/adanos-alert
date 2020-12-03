<template>
    <b-row class="mb-5">
        <b-col>
            <b-card class="mb-2">
                <b-card-text>
                    <b-badge :variant="$route.query.status === undefined ? 'primary':''" class="mr-1" :to="'/'">全部</b-badge>
                    <b-badge :variant="$route.query.status === status.value ? 'primary': ''" v-for="(status, index) in statuses" :key="index" class="mr-1" :to="'/?status=' + status.value">{{ status.name }}</b-badge>
                </b-card-text>
            </b-card>
            <b-table :items="groups" :fields="fields" :busy="isBusy" show-empty hover>
                <template v-slot:cell(id)="row">
                    <b-badge class="mr-2" variant="dark">{{ row.item.seq_num }}</b-badge>
                    <date-time :value="row.item.updated_at"></date-time>
                    <p><b>{{ row.item.id }}</b></p>
                </template>
                <template v-slot:cell(actions)="row">
                    <b-list-group style="font-size: 80%">
                        <b-list-group-item v-for="(act, index) in row.item.actions" :key="index" :variant="act.trigger_status === 'ok' ? 'success': 'danger'">
                            <code class="action-pre-condition" v-b-tooltip :title="act.pre_condition" v-if="!act.is_else_trigger">{{ act.pre_condition || '全部' }}</code>
                            <code class="action-pre-condition" v-if="act.is_else_trigger">兜底</code>
                            <b :class="act.is_else_trigger ? 'text-warning':'text-dark'"> | </b>
                            {{ act.name !== '' ? act.name : formatAction(act.action) }} <span v-if="act.user_refs.length > 0">({{ users(act.user_refs) }})</span>
                        </b-list-group-item>
                    </b-list-group>
                </template>
                <template v-slot:cell(rule_name)="row">
                    <span v-b-tooltip.hover :title="row.item.rule.rule">{{ row.item.rule.name }}</span>
                    <b-link :to="'/rules/?id=' + row.item.rule.id" target="_blank" class="ml-2">
                        <font-awesome-icon icon="external-link-alt"></font-awesome-icon>
                    </b-link>
                    <p>
                        <b-badge v-if="row.item.type === 'recovery'" variant="success" class="mr-2" v-b-tooltip title="事件组类型">恢复</b-badge>
                        <b-badge v-if="row.item.type === 'recoverable'" variant="warning" class="mr-2" v-b-tooltip title="事件组类型">可恢复</b-badge>
                        <b-badge v-b-tooltip.hover title="聚合条件（Key）">{{ row.item.aggregate_key }}</b-badge>
                    </p>
                </template>
                <template v-slot:cell(status)="row">
                    <b-badge v-if="row.item.status === 'collecting'" variant="dark" :title="'预计' + formatted(row.item.rule.expect_ready_at) + '完成'" v-b-tooltip.hover>收集中
                        <span v-if="row.item.collect_time_remain > 0">
                            剩余 <human-time :value="row.item.collect_time_remain"></human-time>
                        </span>
                        <span v-else>收集完成</span>
                    </b-badge>
                    <b-badge v-if="row.item.status === 'pending'" variant="info">准备</b-badge>
                    <b-badge v-if="row.item.status === 'ok'" variant="success">完成</b-badge>
                    <b-badge v-if="row.item.status === 'failed'" variant="danger">失败</b-badge>
                    <b-badge v-if="row.item.status === 'canceled'" variant="warning">已取消</b-badge>
                </template>
                <template v-slot:table-busy class="text-center text-danger my-2">
                    <b-spinner class="align-middle"></b-spinner>
                    <strong> Loading...</strong>
                </template>
                <template v-slot:cell(message_count)="row">
                    <span v-if="row.item.status === 'collecting'">-</span>
                    <span v-else>{{ row.item.message_count }}</span>
                </template>
                <template v-slot:cell(operations)="row">
                    <b-button-group>
                        <b-button size="sm" variant="info" :to="{path:'/events', query: {group_id: row.item.id}}">详情</b-button>
                        <b-dropdown size="sm" right text="预览" variant="primary">
                            <b-dropdown-item :href="$store.getters.serverUrl + '/ui/groups/' + row.item.id + '.html'" target="_blank">事件组</b-dropdown-item>
                            <b-dropdown-item :href="$store.getters.serverUrl + '/ui/reports/' + row.item.id + '.html'" target="_blank">报告</b-dropdown-item>
                        </b-dropdown>
                        <b-button size="sm" variant="danger" 
                            @click="clearGroup(row.item.id)"
                            v-if="row.item.status != 'collecting' && row.item.status != 'pending'">清理</b-button>
                    </b-button-group>
                </template>
            </b-table>
            <paginator :per_page="10" :cur="cur" :next="next" path="/" :query="this.$route.query"></paginator>
        </b-col>
    </b-row>
</template>

<script>
    import moment from 'moment';
    import axios from 'axios';

    export default {
        name: 'Groups',
        data() {
            return {
                groups: [],
                cur: parseInt(this.$route.query.next !== undefined ? this.$route.query.next : 0),
                next: -1,
                isBusy: true,
                userRefs: {},
                fields: [
                    {key: 'id', label: '时间/事件组 ID'},
                    {key: 'rule_name', label: '匹配规则'},
                    {key: 'actions', label: '触发动作'},
                    {key: 'message_count', label: '事件数'},
                    {key: 'status', label: '状态'},
                    {key: 'operations', label: '操作'}
                ],
                statuses: [
                    {value: 'collecting', name:'收集中'},
                    {value: 'pending', name:'准备'},
                    {value: 'ok', name:'完成'},
                    {value: 'failed', name:'失败'},
                    {value: 'canceled', name:'取消'},
                ]
            };
        },
        watch: {
            '$route': 'reload',
        },
        methods: {
            formatted(t) {
                return moment(t).format('YYYY-MM-DD HH:mm:ss');
            },
            users(user_refs) {
                return user_refs.map((u) => {
                    return this.userRefs[u] !== undefined ? this.userRefs[u] : '-';
                }).join(', ')
            },
            clearGroup(id) {
                this.$bvModal.msgBoxConfirm('确定执行该操作 ?').then((value) => {
                    if (value !== true) {
                        return;
                    }

                    axios.delete('/api/groups/' + id + '/reduce/').then((resp) => {
                        this.SuccessBox('操作成功，删除事件数为 ' + resp.data.deleted_count);
                    }).catch(error => {
                        this.ErrorBox(error);
                    });
                });
            },
            formatAction(action) {
                let actions = {
                    'dingding': '钉钉通知',
                    'email': '邮件通知',
                    'phone_call': '电话通知',
                    'wechat': '微信通知',
                    'sms': '短信通知',
                    'http': 'HTTP',
                    'jira': 'JIRA',
                };

                return actions[action];
            },
            reload() {
                let params = this.$route.query;
                params.offset = this.cur;
                axios.get('/api/groups/', {
                    params: params,
                }).then(response => {
                    this.groups = response.data.groups;
                    this.userRefs = response.data.users;
                    this.next = response.data.next;
                    this.isBusy = false;
                }).catch(error => {
                    this.ToastError(error)
                });
            }
        },
        mounted() {
            this.reload();
        }
    }
</script>

<style scoped>
.action-pre-condition {
    max-width: 400px;
    display: inline-block;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}
</style>