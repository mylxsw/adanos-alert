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
                        <b-list-group-item v-for="(act, index) in row.item.actions" :key="index">
                            {{ act.action }} <b class="text-success">=></b> {{ act.trigger_status }}
                        </b-list-group-item>
                    </b-list-group>
                </template>
                <template v-slot:cell(status)="row">
                    <b-badge v-if="row.item.status === 'collecting'" variant="dark">收集中</b-badge>
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
                        <b-button size="sm" variant="info" :to="{path:'/messages', query: {group_id: row.item.id}}">查看</b-button>
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
                fields: [
                    {key: 'id', label: '序号'},
                    {key: 'rule.name', label: '规则'},
                    {key: 'actions', label: '动作'},
                    {key: 'message_count', label: '消息数量'},
                    {key: 'status', label: '状态'},
                    {key: 'operations', label: '操作'}
                ],
            };
        },
        methods: {
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