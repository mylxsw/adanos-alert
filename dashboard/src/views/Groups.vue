<template>
    <b-row class="mb-5">
        <b-col>
            <b-table :items="groups" :fields="fields" :busy="isBusy" show-empty>
                <template v-slot:cell(id)="row">
                    <date-time :value="row.item.updated_at"></date-time>
                    <p><b>{{ row.item.id }}</b></p>
                </template>
                <template v-slot:cell(actions)="row">
                    <b-list-group>
                        <b-list-group-item v-for="(act, index) in row.item.actions" :key="index" :variant="act.trigger_status === 'ok' ? 'success': 'danger'">
                            <code>{{ act.pre_condition || 'true' }}</code> <b class="text-dark"> | </b>
                            {{ formatAction(act.action) }} <span v-if="act.user_refs.length > 0">({{ users(act.user_refs) }})</span>
                        </b-list-group-item>
                    </b-list-group>
                </template>
                <template v-slot:cell(rule_name)="row">
                    <span v-b-tooltip.hover :title="row.item.rule.rule">{{ row.item.rule.name }}</span>
                    <b-link :to="'/rules/' + row.item.rule.id + '/edit'" target="_blank" class="ml-2">
                        <font-awesome-icon icon="external-link-alt"></font-awesome-icon>
                    </b-link>
                </template>
                <template v-slot:cell(status)="row">
                    <b-badge v-if="row.item.status === 'collecting'" variant="dark">收集中（剩余 {{ row.item.collect_time_remain > 0 ? time_remain(row.item.collect_time_remain) : '-' }}）</b-badge>
                    <b-badge v-if="row.item.status === 'pending'" variant="info">准备</b-badge>
                    <b-badge v-if="row.item.status === 'ok'" variant="success">完成</b-badge>
                    <b-badge v-if="row.item.status === 'failed'" variant="danger">失败</b-badge>
                    <b-badge v-if="row.item.status === 'canceled'" variant="warning">已取消</b-badge>
                </template>
                <template v-slot:table-busy class="text-center text-danger my-2">
                    <b-spinner class="align-middle"></b-spinner>
                    <strong> Loading...</strong>
                </template>
                <template v-slot:cell(operations)="row">
                    <b-button-group>
                        <b-button size="sm" variant="info" :to="{path:'/messages', query: {group_id: row.item.id}}">详情</b-button>
                        <b-button size="sm" variant="primary" :href="$store.getters.serverUrl + '/ui/groups/' + row.item.id + '.html'" target="_blank">预览</b-button>
                    </b-button-group>
                </template>
            </b-table>
            <paginator :per_page="10" :cur="cur" :next="next" path="/"></paginator>
        </b-col>
    </b-row>
</template>

<script>
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
                    {key: 'id', label: '时间/ID'},
                    {key: 'rule_name', label: '规则'},
                    {key: 'actions', label: '动作'},
                    {key: 'message_count', label: '消息数量'},
                    {key: 'status', label: '状态'},
                    {key: 'operations', label: '操作'}
                ],
            };
        },
        methods: {
            users(user_refs) {
                return user_refs.map((u) => {
                    return this.userRefs[u] !== undefined ? this.userRefs[u] : '-';
                }).join(', ')
            },
            formatAction(action) {
                let actions = {
                    'dingding': '钉钉通知',
                    'email': '邮件通知',
                    'phone_call': '电话通知',
                    'wechat': '微信通知',
                    'sms': '短信通知',
                    'http': 'HTTP',
                };

                return actions[action];
            },
            time_remain(time_sec) {
                if (time_sec < 60) {
                    return time_sec + " s"
                }

                if (time_sec < 3600) {
                    return (time_sec / 60).toFixed(1) + " m"
                }

                return (time_sec / 60 / 60).toFixed(1) + " h"
            },
            reload() {
                axios.get('/api/groups/?offset=' + this.cur).then(response => {
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