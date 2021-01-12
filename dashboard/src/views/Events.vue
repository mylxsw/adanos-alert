<template>
    <b-row class="mb-5">
        <b-col>
            <b-card class="mb-2">
                <b-card-text>
                    <b-form inline @submit="searchSubmit">
                        <b-input class="mb-2 mr-sm-2 mb-sm-0" placeholder="来源" v-model="search.origin"></b-input>
                        <b-form-tags v-model="search.tags" class="mb-2 mr-sm-2 mb-sm-0" placeholder="标签"></b-form-tags>
                        <b-form-select v-model="search.status" class="mb-2 mr-sm-2 mb-sm-0" placeholder="状态" :options="status_options"></b-form-select>
                        <b-button variant="light" type="submit">搜索</b-button>
                    </b-form>   
                </b-card-text>
            </b-card>

            <b-card class="mb-2" border-variant="warning" v-if="relationInfo != null" header="相关事件摘要" header-bg-variant="warning">
                <b-card-body>
                    <b-row style="max-width: 100rem;" class="adanos-meta-line">
                        <b-col sm="3"><b class="text-black-50" style="border-bottom: 1px dashed black">事件摘要</b></b-col>
                        <b-col sm="9" style="text-align: left">
                            {{ relationInfo.summary }}
                        </b-col>
                    </b-row>
                    <b-row style="max-width: 100rem;" class="adanos-meta-line">
                        <b-col sm="3"><b class="text-black-50" style="border-bottom: 1px dashed black">出现次数</b></b-col>
                        <b-col sm="9" style="text-align: left">
                            {{ relationInfo.event_count }}
                        </b-col>
                    </b-row>
                    <b-row style="max-width: 100rem;" class="adanos-meta-line">
                        <b-col sm="3"><b class="text-black-50" style="border-bottom: 1px dashed black">首次出现时间</b></b-col>
                        <b-col sm="9" style="text-align: left">
                            <date-time :value="relationInfo.created_at"></date-time>
                        </b-col>
                    </b-row>
                    <b-row style="max-width: 100rem;" class="adanos-meta-line">
                        <b-col sm="3"><b class="text-black-50" style="border-bottom: 1px dashed black">最近出现时间</b></b-col>
                        <b-col sm="9" style="text-align: left">
                            <date-time :value="relationInfo.updated_at"></date-time>
                        </b-col>
                    </b-row>

                    <div class="mt-4" v-if="relationNotes != null">
                        <b-media v-for="(note, index) in relationNotes.notes" :key="index" class="mt-3">
                            <template #aside>
                                <b-avatar variant="success" icon="people-fill"></b-avatar>
                            </template>

                            <p>
                                <date-time :value="note.created_at"></date-time>
                                <b-link class="mr-1 float-right" :to="{path:'/events', query: {event_id: note.event_id}}" target="_blank" title="关联事件">
                                    <font-awesome-icon icon="external-link-alt"></font-awesome-icon>
                                </b-link>
                            </p>
                            <p class="text-break font-weight-light text-muted">{{ note.note }}</p>
                        </b-media>
                    </div>

                    <b-row>
                        <v-charts :options="eventsByDatetime" style="width: 100%;"></v-charts>
                    </b-row>

                </b-card-body>
            </b-card>

            <EventCard v-for="(message, index) in events" :key="index" class="mb-3"
                         :event="message"
                         :event_index="index" :reproduce-event="reproduceEvent" :eventNote="openEventNoteDialog" :deleteEvent="deleteEvent"
                         :test-matched-rules="testMatchedRules"></EventCard>
            <b-card v-if="events.length === 0">
                <b-card-body>There are no records to show</b-card-body>
            </b-card>
            <paginator :per_page="2" :cur="cur" :next="next" path="/events" :query="this.$route.query"></paginator>
        </b-col>

        <b-modal id="matched-rules-dialog" title="匹配到的规则" hide-footer size="xl">
            <b-table responsive="true" :items="matched_rules" :fields="matched_rules_fields">
                <template v-slot:cell(name)="row">
                    <span v-b-tooltip.hover :title="row.item.rule.id">{{ row.item.rule.name }}</span>
                    <b-link :to="'/rules/' + row.item.rule.id + '/edit'" target="_blank" class="ml-2">
                        <font-awesome-icon icon="external-link-alt"></font-awesome-icon>
                    </b-link>
                    <p>
                        <b-badge variant="info" v-for="(tag, index) in row.item.rule.tags" :key="index" class="mr-1">{{ tag }}</b-badge>
                    </p>
                </template>
                <template v-slot:cell(rule)="row">
                    <p class="adanos-pre-fold" v-b-tooltip.hover :title="row.item.rule.rule">
                        <code>{{ row.item.rule.rule }}</code>
                    </p>
                </template>
                <template v-slot:cell(aggregate_rule)="row">
                    <p class="adanos-pre-fold" v-b-tooltip.hover :title="row.item.rule.aggregate_rule">
                        <code>{{ row.item.rule.aggregate_rule }}</code>
                    </p>
                    <p><b-badge variant="success" v-if="row.item.aggregate_key" v-b-tooltip.hover title="实际匹配的聚合 Key">{{ row.item.aggregate_key }}</b-badge></p>
                </template>
            </b-table>
        </b-modal>

        <b-modal id="event-note-dialog" title="事件备注" hide-footer size="xl">
            <b-form @submit="onEventNoteSubmit">
                <b-form-group id="event_note" label-for="event_note_input">
                    <b-form-textarea id="event_note_input" placeholder="输入事件备注内容" v-model="event_note_form.note" rows="6"/>
                </b-form-group>
                <b-button type="submit" variant="primary" class="mr-2 float-right">保存</b-button>
            </b-form>
        </b-modal>
    </b-row>
</template>

<script>
import axios from 'axios';
import Echarts from 'vue-echarts';

import 'echarts/lib/component/tooltip'
import 'echarts/lib/component/axis';
import 'echarts/lib/component/legend';
import 'echarts/lib/component/toolbox';
import 'echarts/lib/component/title';

import 'echarts/lib/chart/line';
import 'echarts/lib/component/polar';
import 'echarts/lib/component/dataZoom';
import {graphic} from 'echarts/lib/export';

export default {
        name: 'Events',
        components: {
            'v-charts': Echarts,
        },
        data() {
            return {
                search: {
                    origin: '',
                    status: this.$route.query.status !== undefined ? this.$route.query.status : null,
                    tags: [],
                    meta: '',
                },
                status_options: [
                    {value: null, text: '所有状态'},
                    {value: 'pending', text: '准备中'},
                    {value: 'grouped', text: '已分组'},
                    {value: 'canceled', text: '无规则，已取消'},
                    {value: 'expired', text: '匹配规则，已过期'},
                    {value: 'ignored', text: '匹配规则，已忽略'},
                ],
                events: [],
                cur: parseInt(this.$route.query.next !== undefined ? this.$route.query.next : 0),
                next: -1,
                matched_rules: [],
                matched_rules_fields: [
                    {key: 'name', label: '规则名称/ID'},
                    {key: 'rule', label: '规则'},
                    {key: 'aggregate_rule', label: '聚合条件'},
                ],
                relationInfo: {},
                relationNotes: [],
                event_note_form: {
                    note: '',
                    event_id: null,
                },
                eventsByDatetime: {
                  title: {left: 'left', text: '事件数量时间分布'},
                  tooltip: {
                    trigger: 'axis',
                  },
                  xAxis: {
                    type: 'category',
                    data: []
                  },
                  yAxis: {
                    type: 'value'
                  },
                  grid: {left: 50},
                  dataZoom: [{
                    type: 'inside',
                    start: 50,
                    end: 100
                  }, {
                    start: 50,
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
                  toolbox: {
                    show: true,
                    feature: {
                      dataZoom: {
                        yAxisIndex: "none"
                      },
                    }
                  },
                  series: {
                    smooth: true,
                    name: '事件数量',
                    data: [],
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
                },
            };
        },
        watch: {
            '$route': 'loadMore',
        },
        methods: {
            searchSubmit(evt) {
                evt.preventDefault();

                this.$router.push({path: '/events', query: {
                    offset: 0,
                    group_id: this.$route.query.group_id !== undefined ? this.$route.query.group_id : null,
                    relation_id: this.$route.query.relation_id !== undefined ? this.$route.query.relation_id : null,
                    status: this.search.status,
                    tags: this.search.tags.join(),
                    meta: this.search.meta,
                    origin: this.search.origin,
                }}).catch(err => {this.ToastError(err);});
            },
            reproduceEvent(id) {
                this.$bvModal.msgBoxConfirm('确定执行该操作 ?').then((value) => {
                    if (value !== true) {
                        return;
                    }

                    axios.post('/api/events/' + id + '/reproduce/', {}).then(resp => {
                        this.ToastSuccess('重新投递 Event 成功，Event ID: ' + resp.data.id);
                    }).catch(error => {
                        this.ErrorBox(error);
                    });
                });
            },
            deleteEvent(index, id) {
                let self = this;
                this.$bvModal.msgBoxConfirm('确定执行该操作 ?').then((value) => {
                    if (value !== true) {
                        return;
                    }

                    axios.delete('/api/events/' + id + '/', {}).then(() => {
                        self.events.splice(index, 1);
                        this.ToastSuccess('事件删除成功');
                    }).catch(error => {
                        this.ErrorBox(error);
                    });
                });
            },
            testMatchedRules(id) {
                axios.post('/api/events/' + id + '/matched-rules/', {}).then(resp => {
                    this.matched_rules = resp.data;
                    this.$root.$emit('bv::show::modal', "matched-rules-dialog");
                }).catch(error => {
                    this.ErrorBox(error);
                });
            },
            openEventNoteDialog(eventId) {
                this.event_note_form.event_id = eventId;
                this.$root.$emit('bv::show::modal', "event-note-dialog");
            },
            onEventNoteSubmit(evt) {
                evt.preventDefault();

                if (this.event_note_form.note.trim() === '') {
                    this.ErrorBox('备注内容不能为空');
                    return ;
                }

                axios.post('/api/event-relations/' + this.$route.query.relation_id + '/notes/', this.event_note_form).then(() => {
                    this.ToastSuccess('操作成功');
                    this.$root.$emit('bv::hide::modal', "event-note-dialog");
                    this.loadMore();
                }).catch(error => {
                    this.ErrorBox(error)
                });
            },
            loadCharts(relation_id) {
                axios.get('/api/statistics/events/period-counts/', {params: {relation_id: relation_id, days: 7}}).then(resp => {
                    this.eventsByDatetime.xAxis.data = resp.data.map(s => s.datetime);
                    this.eventsByDatetime.series.data = resp.data.map(s => s.total);
                }).catch(error => {this.ToastError(error)});
            },
            loadMore() {
                let params = this.$route.query;
                params.offset = this.cur;
                
                axios.get('/api/events/', {
                    params: params,
                }).then(response => {
                    this.events = response.data.events;
                    for (let i in this.events) {
                        this.events[i]._showDetails = true;
                    }

                    this.next = response.data.next;

                    this.search.origin = response.data.search.origin;
                    this.search.tags = response.data.search.tags;
                    this.search.meta = response.data.search.meta;
                    this.search.status = response.data.search.status.length === 0 ? null : response.data.search.status[0];

                    this.isBusy = false;
                }).catch(error => {
                    this.ToastError(error)
                });

                this.relationInfo = null;
                this.relationNotes = null;
                let relationID = this.$route.query.relation_id;
                if (relationID !== undefined && relationID !== '' && relationID !== null) {
                    axios.get('/api/event-relations/' + relationID + '/').then(response => {
                        this.relationInfo = response.data;
                    }).catch(error => {this.ToastError(error)})

                    axios.get('/api/event-relations/' + relationID + '/notes/').then(response => {
                        this.relationNotes = response.data;
                    }).catch(error => {this.ToastError(error)})

                    this.loadCharts(relationID);
                }
            }
        },
        mounted() {
            this.loadMore();
        }
    }
</script>
