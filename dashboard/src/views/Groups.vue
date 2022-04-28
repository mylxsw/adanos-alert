<template>
    <b-row class="mb-5">
        <b-col>
            <b-card class="mb-2">
                <b-card-text>
                    <b-badge :variant="$route.query.status === undefined ? 'primary':''" class="mr-1" :to="'/'">全部</b-badge>
                    <b-badge :variant="$route.query.status === status.value ? 'primary': ''" v-for="(status, index) in statuses" :key="index" class="mr-1" :to="'/?status=' + status.value">{{ status.name }}</b-badge>
                </b-card-text>

                <b-card-text>
                    <b-form inline @submit="searchSubmit">
                        <date-picker type="datetime" v-model="search.time_range" range clearable class="mr-2" style="width: 400px;"/>
                        <b-form-select v-model="search.type" class="mb-2 mr-sm-2 mb-sm-0" placeholder="类型" :options="type_options"></b-form-select>
                        <b-form-select v-model="search.sort" class="mb-2 mr-sm-2 mb-sm-0" placeholder="排序方式" :options="sort_options"></b-form-select>
                        <b-button variant="primary" type="submit">刷新</b-button>
                    </b-form>
                </b-card-text>

                <b-row>
                    <b-col cols="4"><v-charts :options="alertByAgg" style="width: 95%;" ref="alertByAgg" v-if="vchartShow"></v-charts></b-col>
                    <b-col cols="8"><v-charts :options="alertByDatetime" style="width: 95%;" ref="alertByDatetimeChart" v-if="vchartShow"></v-charts></b-col>
                </b-row>
            </b-card>
            <b-table :items="groups" :fields="fields" :busy="isBusy" show-empty hover>
                <template v-slot:cell(id)="row">
                    <b-badge class="mr-2" variant="dark" v-b-tooltip.hover :title="row.item.id">{{ row.item.seq_num }}</b-badge>
                    <span v-b-tooltip.hover :title="row.item.rule.rule">{{ row.item.rule.name }}</span>
                    <b-link :to="'/rules/?id=' + row.item.rule.id" target="_blank" class="ml-2">
                        <font-awesome-icon icon="external-link-alt"></font-awesome-icon>
                    </b-link>
                    <p><date-time :value="row.item.updated_at"></date-time></p>
                </template>
                <template v-slot:cell(actions)="row">
                    <p style="margin-bottom:4px">
                        <b-badge v-if="row.item.type === 'recovery'" variant="success" class="mr-2" v-b-tooltip title="事件组类型">恢复</b-badge>
                        <b-badge v-if="row.item.type === 'recoverable'" variant="info" class="mr-2" v-b-tooltip title="事件组类型">可恢复</b-badge>
                        <b-badge v-if="row.item.type === 'ignored'" variant="warning" class="mr-2" v-b-tooltip title="事件组类型">忽略事件</b-badge>
                        <b-badge v-if="row.item.type === 'ignoredExceed'" variant="primary" class="mr-2" v-b-tooltip title="事件组类型">超限忽略事件</b-badge>
                        <b-badge class="mr-2" variant="danger" v-if="row.item.rule.realtime" v-b-tooltip title="即时消息">即时</b-badge>
                        <b-badge v-b-tooltip.hover title="聚合条件（Key）">{{ row.item.aggregate_key }}</b-badge>
                    </p>
                    <b-list-group style="font-size: 80%">
                        <b-list-group-item style="padding: 5px" v-for="(act, index) in row.item.actions" :key="index" :variant="act.trigger_status === 'ok' ? 'success': 'danger'">
                            <code style="line-height: 0.85" class="action-pre-condition" v-b-tooltip :title="act.pre_condition" v-if="!act.is_else_trigger">{{ act.pre_condition || '全部' }}</code>
                            <code style="line-height: 0.85" class="action-pre-condition" v-if="act.is_else_trigger">兜底</code>
                            <b :class="act.is_else_trigger ? 'text-warning':'text-dark'"> | </b>
                            {{ act.name !== '' ? act.name : formatAction(act.action) }} <span v-if="act.user_refs.length > 0">({{ users(act) }})</span>
                        </b-list-group-item>
                    </b-list-group>
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
    import Echarts from 'vue-echarts';
    import 'echarts/lib/chart/line';

    import 'echarts/lib/component/tooltip'
    import 'echarts/lib/component/axis';
    import 'echarts/lib/component/legend';
    import 'echarts/lib/component/toolbox';
    import 'echarts/lib/component/polar';
    import 'echarts/lib/component/title';

    import 'echarts/lib/component/dataZoom';
    import { graphic } from 'echarts/lib/export'

    import 'echarts/lib/chart/pie';
    import DatePicker from '@hyjiacan/vue-datepicker'
    import '@hyjiacan/vue-datepicker/dist/datepicker.css'

    export default {
        name: 'Groups',
        components: {
            'v-charts': Echarts,
            'date-picker': DatePicker,
        },
        data() {
            return {
                search: {
                    time_range: [
                      this.$route.query.start_at !== undefined ? this.$route.query.start_at : null,
                      this.$route.query.end_at !== undefined ? this.$route.query.end_at : null,
                    ],
                    sort: this.$route.query.sort !== undefined ? this.$route.query.sort : 'desc',
                    type: this.$route.query.type !== undefined ? this.$route.query.type : '',
                },
                sort_options: [
                    {value: 'asc', text: '创建时间正序'},
                    {value: 'desc', text: '创建时间倒序'},
                ],
                type_options: [
                    {value: '', text: '所有类型（忽略事件组除外）'},
                    {value: 'plain', text: '普通事件组'},
                    {value: 'recoverable', text: '可恢复事件组'},
                    {value: 'recovery', text: '已恢复事件组'},
                    {value: 'ignored', text: '忽略事件组'},
                    {value: 'ignoredExceed', text: '超限忽略事件组'},
                ],
                groups: [],
                cur: parseInt(this.$route.query.next !== undefined ? this.$route.query.next : 0),
                next: -1,
                isBusy: true,
                userRefs: {},
                fields: [
                    {key: 'id', label: '时间'},
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
                ],
                alertByDatetime: {
                    title: {
                        text: 'Event Group Count',
                        left: 'center',
                        textStyle: {
                            color: '#ccc'
                        }
                    },
                    tooltip: {trigger: 'axis',},
                    xAxis: {
                        type: 'category',
                        data: []
                    },
                    yAxis: {type: 'value',},
                    grid: {left: 50, right: 50},
                    dataZoom: [{
                        type: 'inside',
                        start: 0,
                        end: 100
                    },{
                        start: 0,
                        end: 100,
                        handleIcon: 'M10.7,11.9v-1.3H9.3v1.3c-4.9,0.3-8.8,4.4-8.8,9.4c0,5,3.9,9.1,8.8,9.4v1.3h1.3v-1.3c4.9-0.3,8.8-4.4,8.8-9.4C19.5,16.3,15.6,12.2,10.7,11.9z M13.3,24.4H6.7V23h6.6V24.4z M13.3,19.6H6.7v-1.4h6.6V19.6z',
                        handleSize: '80%',
                        handleStyle: {
                            color: '#fff',
                            shadowBlur: 3,
                            shadowColor: 'rgba(0, 0, 0, 0.6)',
                            shadowOffsetX: 2,
                            shadowOffsetY: 2
                        }
                    }],
                    series: [],
                },
                alertByAgg: {
                    title: {
                        text: 'Aggregate Keys',
                        left: 'center',
                        textStyle: {
                            color: '#ccc'
                        }
                    },

                    tooltip: {
                        trigger: 'item',
                        formatter: '{b} : {c} ({d}%)'
                    },

                    series: [],
                },
                vchartShow: false,
            };
        },
        watch: {
            '$route': 'reload',
        },
        methods: {
            searchSubmit(evt) {
                evt.preventDefault();
                let query = {};
                for (let i in this.$route.query) {
                    query[i] = this.$route.query[i];
                }

                if (this.search['time_range'][0] != null) {
                    query['start_at'] = this.search['time_range'][0];
                }
                if (this.search['time_range'][1] != null) {
                    query['end_at'] = this.search['time_range'][1];
                }

                query['sort'] = this.search['sort'];
                query['type'] = this.search['type'];

                this.$router.push({path: '/', query: query}).catch(err => {
                  this.ToastError(err)
                });
            },
            formatted(t) {
                return moment(t).format('YYYY-MM-DD HH:mm:ss');
            },
            users(act) {
                if (act.user_names !== undefined && act.user_names != null && act.user_names.length > 0) {
                  return act.user_names.join(', ');
                }

                return act.user_refs.map((u) => {
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
                    'dingding': '钉钉',
                    'email': '邮件',
                    'phone_call_aliyun': '电话',
                    'wechat': '微信',
                    'sms': '短信',
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

                if (this.$route.query.rule_id !== undefined) {
                    this.vchartShow = true;
                    axios.get('/api/statistics/daily-group-counts/', {
                        params: params,
                    }).then(response => {
                        this.alertByDatetime.xAxis.data = response.data.map(s => s.datetime);
                        this.alertByDatetime.series = {
                            smooth: true,
                            name: '事件组数量',
                            data: response.data.map(s => s.total),
                            type: 'line',
                            sampling: 'average',
                            itemStyle: {
                                color: 'rgb(255, 70, 131)'
                            },
                            areaStyle: {
                                color: new graphic.LinearGradient(0, 0, 0, 1, [{
                                    offset: 0,
                                    color: 'rgb(255, 158, 68)'
                                }, {
                                    offset: 1,
                                    color: 'rgb(255, 70, 131)'
                                }])
                            }
                        }
                    }).catch(error => {this.ToastError(error)});

                    axios.get('/api/statistics/group-agg-counts/', {
                        params: params,
                    }).then(response => {
                        this.alertByAgg.series = [
                            {
                                type: 'pie',
                                radius: '55%',
                                center: ['50%', '50%'],
                                labelLine: {smooth: 0.2, length: 10, length2: 20},
                                animationType: 'scale',
                                animationEasing: 'elasticOut',
                                data: response.data.map(s => { return {name: s.aggregate_key, value: s.total}}),
                            }
                        ];
                    }).catch(error => {this.ToastError(error)});

                }
            }
        },
        mounted() {
            this.reload();

            // this.$refs.alertByDatetimeChart.chart.on('axisareaselected', () => {alert("hello")});
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

p {
    margin-bottom: 5px;
}
</style>