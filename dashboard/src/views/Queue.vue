<template>
    <b-row class="mb-5">
        <b-col>
            <div class="mb-3">
                <b-button size="sm" :variant="queue_btn" class="float-right" @click="pause_queue()">{{ queue_action }}</b-button>
              Current status: <span v-html="queue_status"></span>，
              Workers: <b>{{ queue_info.worker_num }}</b>，
              Processed since <date-time :value="queue_info.start_at"></date-time>: <b-badge :to="'/queues?status=succeed'" variant="success">{{ queue_info.processed_count }}</b-badge>, Failed: <b-badge :to="'/queues?status=failed'" variant="danger">{{ queue_info.failed_count }}</b-badge>
            </div>
            <b-table :items="jobs" :fields="fields" :busy="isBusy" show-empty hover>
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
                    <b-badge v-if="row.item.status === 'wait'" variant="info">Waiting
                        <span v-if="row.item.execute_time_remain > 0">(Remaining <human-time :value="row.item.execute_time_remain"></human-time>)</span>
                        <span v-else>(Ready)</span>
                    </b-badge>
                    <b-badge v-if="row.item.status === 'running'" variant="dark">In progress</b-badge>
                    <b-badge v-if="row.item.status === 'failed'" variant="danger" v-b-tooltip.hover :title="row.item.last_error">Failed</b-badge>
                    <b-badge v-if="row.item.status === 'succeed'" variant="success">Done</b-badge>
                    <b-badge v-if="row.item.status === 'canceled'" variant="warning">Canceled</b-badge>
                </template>
                <template v-slot:table-busy class="text-center text-danger my-2">
                    <b-spinner class="align-middle"></b-spinner>
                    <strong> Loading...</strong>
                </template>
                <template v-slot:row-details="row">
                    <b-card>
                        <code><pre class="text-danger adanos-code">{{ JSON.stringify(row.item.payload, null, 4) }}</pre></code>
                    </b-card>
                </template>
                <template v-slot:cell(operations)="row">
                    <b-button-group>
                        <b-button size="sm" variant="info" @click.stop="row.toggleDetails" v-model="row.detailsShowing">Payload</b-button>
                        <b-button size="sm" variant="danger" @click="delete_job(row.index, row.item.id)">Delete</b-button>
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
                status: this.$route.query.status !== undefined ? this.$route.query.status : '',
                next: -1,
                isBusy: true,
                queue_info: {
                    worker_num: '-',
                    processed_count: '-',
                    failed_count: '-',
                    start_at: '-',
                },
                queue_paused: false,
                queue_status: "-",
                queue_action: "Start",
                queue_btn: "success",
                fields: [
                    {key: 'id', label: 'Time/ID'},
                    {key: 'name', label: 'Type'},
                    {key: 'requeue_times', label: 'Retry times'},
                    {key: 'next_execute_at', label: 'Expected execution time'},
                    {key: 'updated_at', label: 'Last update'},
                    {key: 'status', label: 'Status'},
                    {key: 'operations', label: 'Operations'}
                ],
            };
        },
        watch: {
          '$route': 'reload',
        },
        methods: {
            delete_job(index, id) {
                let self = this;
                this.$bvModal.msgBoxConfirm('Are you sure to perform this operation?').then((value) => {
                    if (value !== true) {
                        return;
                    }

                    axios.delete('/api/queue/jobs/' + id + '/').then(() => {
                        self.jobs.splice(index, 1);
                        this.SuccessBox('Operation successful');
                    }).catch(error => {
                        this.ErrorBox(error);
                    });
                });
            },
            pause_queue() {
                this.$bvModal.msgBoxConfirm('Are you sure to perform this operation?').then((value) => {
                    if (value !== true) {return;}
                    axios.post('/api/queue/control/', {op: this.queue_paused ? 'continue' : 'pause'}).then(resp => {
                        this.SuccessBox('Operation successful');
                        this.updateControlStatus(resp.data.paused);
                        this.queue_info = resp.data.info;
                    }).catch(error => {this.ErrorBox(error)});
                });

            },
            updateControlStatus(paused) {
                this.queue_status = paused ? '<b class="text-warning">Stopped</b>':'<b class="text-success">Running</b>';
                this.queue_action = paused ? 'Start' : 'Stop';
                this.queue_btn = paused ? 'success' : 'warning';
                this.queue_paused = paused;
            },
            reload() {
                let params = this.$route.query;
                params.offset = this.cur;
                axios.get('/api/queue/jobs/', {
                  params: params
                }).then(response => {
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
            this.reload();
        }
    }
</script>

<style scoped>
    .adanos-code {
        white-space: pre-wrap!important;
        word-wrap: break-word!important;
        *white-space:normal!important;
    }
</style>