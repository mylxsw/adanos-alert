<template>
    <b-row class="mb-5">
        <b-col>
            <div class="mb-3">
                <b-button size="sm" :variant="queue_btn" class="float-right" @click="pause_queue()">{{ queue_action }}</b-button>
                当前状态：<span v-html="queue_status"></span>，
                Workers: <b>{{ queue_info.worker_num }}</b>，
                已处理：<b class="text-success">{{ queue_info.processed_count }}</b>，
                失败：<b class="text-danger">{{ queue_info.failed_count }}</b>
            </div>
            <b-table :items="jobs" :fields="fields" :busy="isBusy" show-empty>
                <template v-slot:cell(id)="row">
                    <date-time :value="row.item.created_at"></date-time>
                    <p><b>{{ row.item.id }}</b></p>
                </template>
                <template v-slot:cell(next_execute_at)="row">
                    <date-time :value="row.item.next_execute_at"></date-time>
                </template>
                <template v-slot:cell(updated_at)="row">
                    <date-time :value="row.item.updated_at"></date-time>
                </template>
                <template v-slot:cell(requeue_times)="row">
                    <b class="text-danger" v-if="row.item.requeue_times > 0">{{ row.item.requeue_times }}</b>
                    <b v-else>{{ row.item.requeue_times }}</b>
                </template>
                <template v-slot:cell(status)="row">
                    <b-badge v-if="row.item.status === 'wait'" variant="info">等待</b-badge>
                    <b-badge v-if="row.item.status === 'running'" variant="dark">执行中</b-badge>
                    <b-badge v-if="row.item.status === 'failed'" variant="danger">失败</b-badge>
                    <b-badge v-if="row.item.status === 'succeed'" variant="success">成功</b-badge>
                    <b-badge v-if="row.item.status === 'canceled'" variant="warning">已取消</b-badge>
                </template>
                <template v-slot:table-busy class="text-center text-danger my-2">
                    <b-spinner class="align-middle"></b-spinner>
                    <strong> Loading...</strong>
                </template>
                <template v-slot:row-details="row">
                    <b-card>
                        <b-card-text>
                            <code><pre class="text-danger">{{ JSON.stringify(row.item.payload, null, 4) }}</pre></code>
                        </b-card-text>
                    </b-card>
                </template>
                <template v-slot:cell(operations)="row">
                    <b-button-group>
                        <b-button size="sm" variant="info" @click.stop="row.toggleDetails" v-model="row.detailsShowing">Payload</b-button>
                        <b-button size="sm" variant="danger" @click="delete_job(row.index, row.item.id)">删除</b-button>
                    </b-button-group>
                </template>
            </b-table>
            <paginator :per_page="10" :cur="cur" :next="next" path="/queues" :query="this.$route.query"></paginator>
        </b-col>
    </b-row>
</template>

<script>
    import axios from 'axios'

    export default {
        name: 'Queue',
        data() {
            return {
                jobs: [],
                cur: parseInt(this.$route.query.next !== undefined ? this.$route.query.next : 0),
                next: -1,
                isBusy: true,
                queue_info: {
                    worker_num: '-',
                    processed_count: '-',
                    failed_count: '-',
                    started_at: '-',
                },
                queue_paused: false,
                queue_status: "-",
                queue_action: "启动",
                queue_btn: "success",
                fields: [
                    {key: 'id', label: '时间/ID'},
                    {key: 'name', label: '类型'},
                    {key: 'requeue_times', label: '重试次数'},
                    {key: 'next_execute_at', label: '最早执行时间'},
                    {key: 'updated_at', label: '最后更新'},
                    {key: 'status', label: '状态'},
                    {key: 'operations', label: '操作'}
                ],
            };
        },
        methods: {
            delete_job(index, id) {
                let self = this;
                this.$bvModal.msgBoxConfirm('确定执行该操作 ?').then((value) => {
                    if (value !== true) {
                        return;
                    }

                    axios.delete('/api/queue/jobs/' + id + '/').then(() => {
                        self.jobs.splice(index, 1);
                        this.SuccessBox('操作成功');
                    }).catch(error => {
                        this.ErrorBox(error);
                    });
                });
            },
            pause_queue() {
                this.$bvModal.msgBoxConfirm('确定执行该操作 ?').then((value) => {
                    if (value !== true) {return;}
                    axios.post('/api/queue/control/', {op: this.queue_paused ? 'continue' : 'pause'}).then(resp => {
                        this.SuccessBox('操作成功');
                        this.updateControlStatus(resp.data.paused);
                        this.queue_info = resp.data.info;
                    }).catch(error => {this.ErrorBox(error)});
                });

            },
            updateControlStatus(paused) {
                this.queue_status = paused ? '<b class="text-warning">暂停</b>':'<b class="text-success">运行中</b>';
                this.queue_action = paused ? '启动' : '暂停';
                this.queue_btn = paused ? 'success' : 'warning';
                this.queue_paused = paused;
            },
            loadMore() {
                axios.get('/api/queue/jobs/?offset=' + this.cur).then(response => {
                    this.jobs = response.data.jobs;
                    for (let i in this.jobs) {
                        this.jobs[i].payload = JSON.parse(this.jobs[i].payload);
                    }
                    this.next = response.data.next;
                    this.isBusy = false;
                }).catch(error => {
                    this.ToastError(error)
                });

                axios.post('/api/queue/control/', {op: 'info'}).then(response => {
                    this.updateControlStatus(response.data.paused);
                    this.queue_info = response.data.info;
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