<template>
    <b-row class="mb-5">
        <b-col>
            <b-table :items="groups" :fields="fields" :busy="isBusy" show-empty>
                <template v-slot:cell(id)="row">
                    {{ row.item.updated_at }}
                    <p><b>{{ row.item.id }}</b></p>
                </template>
                <template v-slot:cell(actions)="row">
                    <b-list-group>
                        <b-list-group-item v-for="(act, index) in row.item.actions" :key="index">
                            {{ act.action }} <b class="text-success">=></b> {{ act.trigger_status }}
                        </b-list-group-item>
                    </b-list-group>
                </template>
                <template v-slot:table-busy class="text-center text-danger my-2">
                    <b-spinner class="align-middle"></b-spinner>
                    <strong> Loading...</strong>
                </template>
            </b-table>
        </b-col>
    </b-row>
</template>

<script>
    import axios from 'axios'

    export default {
        name: 'Groups',
        data() {
            return {
                groups: [],
                nextOffset: -1,
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
            loadMore() {
                axios.get('/api/groups/?offset=' + this.nextOffset).then(response => {
                    this.groups = response.data.groups;
                    this.nextOffset = response.data.next;
                    this.isBusy = false;
                }).catch(error => {
                    this.$bvToast.toast(error.response !== undefined ? error.response.data.error : error.toString(), {
                        title: 'ERROR',
                        variant: 'danger'
                    });
                });
            }
        },
        mounted() {
            this.loadMore();
        }
    }
</script>