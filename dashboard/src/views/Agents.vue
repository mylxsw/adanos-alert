<template>
    <b-row class="mb-5">
        <b-col>
            <b-table :items="agents" :fields="fields" :busy="isBusy" show-empty>
                <template v-slot:cell(agent_id)="row">
                    <date-time :value="row.item.created_at"></date-time>
                    <p><b>{{ row.item.agent_id }}</b></p>
                </template>
                <template v-slot:cell(last_alive_at)="row">
                    <date-time :value="row.item.last_alive_at"></date-time>
                </template>
                <template v-slot:cell(requeue_times)="row">
                    <b class="text-danger" v-if="row.item.requeue_times > 0">{{ row.item.requeue_times }}</b>
                    <b v-else>{{ row.item.requeue_times }}</b>
                </template>
                <template v-slot:cell(status)="row">
                    <b-badge v-if="row.item.alive" variant="success">活跃</b-badge>
                    <b-badge v-if="!row.item.alive" variant="danger">丢失</b-badge>
                </template>
                <template v-slot:table-busy class="text-center text-danger my-2">
                    <b-spinner class="align-middle"></b-spinner>
                    <strong> Loading...</strong>
                </template>
                <template v-slot:cell(operations)="row">
                    <b-button-group>
                        <b-button size="sm" variant="danger" @click="delete_agent(row.index, row.item.id)">删除</b-button>
                    </b-button-group>
                </template>
            </b-table>
        </b-col>
    </b-row>
</template>

<script>
    import axios from 'axios'

    export default {
        name: 'Agents',
        data() {
            return {
                agents: [],
                isBusy: true,
                fields: [
                    {key: 'agent_id', label: 'ID'},
                    {key: 'version', label: '版本'},
                    {key: 'ip', label: '所在服务器'},
                    {key: 'last_alive_at', label: '最后心跳时间'},
                    {key: 'status', label: '状态'},
                    {key: 'operations', label: '操作'}
                ],
            };
        },
        methods: {
            delete_agent(index, id) {
                let self = this;
                this.$bvModal.msgBoxConfirm('确定执行该操作 ?').then((value) => {
                    if (value !== true) {
                        return;
                    }

                    axios.delete('/api/agents/' + id + '/').then(() => {
                        self.agents.splice(index, 1);
                        this.SuccessBox('操作成功');
                    }).catch(error => {
                        this.ErrorBox(error);
                    });
                });
            },
            loadMore() {
                axios.get('/api/agents/').then(response => {
                    this.agents = response.data;
                    this.isBusy = false;
                }).catch(error => {
                    this.ToastError(error)
                });
            }
        },
        mounted() {
            this.loadMore();
        }
    }
</script>

<style scoped>

</style>