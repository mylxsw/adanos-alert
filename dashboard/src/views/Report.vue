<template>
    <b-row>
        <b-col>
            <b-card class="mb-2 search-box">
                <b-card-text style="display: flex; justify-content:space-between">
                    <b-form inline @submit="searchSubmit">
                        <b-form-datepicker class="mb-2 mr-sm-2 mb-sm-0" type="date" placeholder="开始日期" v-model="search.start_at" today-button></b-form-datepicker>
                        <b-form-datepicker class="mb-2 mr-sm-2 mb-sm-0" type="date" placeholder="截止日期" v-model="search.end_at" today-button></b-form-datepicker>
                        <b-button variant="light" type="submit">刷新</b-button>
                    </b-form>
                </b-card-text>
            </b-card>
            <b-card>
                <b-card-body>
                    <b-row class="mb-5">
                        <b-col cols="4"><v-charts :options="alertByUser" style="width: 100%;"></v-charts></b-col>
                        <b-col cols="8"><v-charts :options="alertByDatetime" style="width: 100%;"></v-charts></b-col>
                    </b-row>
                    <b-row class="mb-5">
                        <b-col cols="4"><v-charts :options="alertByRule" style="width: 100%;"></v-charts></b-col>
                        <b-col cols="8"><v-charts :options="eventsByDatetime" style="width: 100%;" ></v-charts></b-col>
                    </b-row>
                </b-card-body>
            </b-card>
        </b-col>
    </b-row>
</template>

<script>
import Echarts from 'vue-echarts';
import 'echarts/lib/chart/line';

import 'echarts/lib/component/tooltip'
import 'echarts/lib/component/axis';
import 'echarts/lib/component/legend';
import 'echarts/lib/component/toolbox';
import 'echarts/lib/component/polar';
import 'echarts/lib/component/title';

import 'echarts/lib/component/dataZoom';
import {graphic} from 'echarts/lib/export'

import axios from "axios";
import moment from 'moment';

export default {
        components: {
            'v-charts': Echarts,
        },
        name: 'Report',
        data() {
            return {
                search: {
                    start_at: this.$route.query.start_at !== undefined ? this.$route.query.start_at : moment().subtract(7, 'days').format('YYYY-MM-DD'),
                    end_at: this.$route.query.end_at !== undefined ? this.$route.query.end_at : moment().format('YYYY-MM-DD'),
                },
                alertByUser: {
                    title: {
                        text: 'Users',
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
                alertByRule: {
                    title: {
                        text: 'Rules',
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
                alertByDatetime: {
                    title: {
                        text: '事件组数量时间分布',
                        left: 'center',
                        textStyle: {
                            color: '#ccc'
                        }
                    },
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
                        name: '事件组数量',
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
                eventsByDatetime: {
                    title: {
                        text: '事件数量时间分布',
                        left: 'center',
                        textStyle: {
                            color: '#ccc'
                        }
                    },
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
                }
            };
        },
        watch: {
            '$route': 'reload',
        },
        methods: {
            searchSubmit(evt) {
                evt.preventDefault();
                let query = {offset: 0};
                for (let i in this.$route.query) {
                    query[i] = this.$route.query[i];
                }

                for (let i in this.search) {
                    query[i] = this.search[i];
                }

                this.$router.push({path: '/reports', query: query}).catch(err => {
                    this.ToastError(err)
                });
            },
            reload() {
                let params = {params: this.search};
                axios.get('/api/statistics/daily-group-counts/', params).then(response => {
                    this.alertByDatetime.xAxis.data = response.data.map(s => s.datetime);
                    this.alertByDatetime.series.data = response.data.map(s => s.total)
                }).catch(error => {this.ToastError(error)});

                axios.get('/api/statistics/events/period-counts/', params).then(resp => {
                    this.eventsByDatetime.xAxis.data = resp.data.map(s => s.datetime);
                    this.eventsByDatetime.series.data = resp.data.map(s => s.total);
                }).catch(error => {this.ToastError(error)});

                axios.get('/api/statistics/user-group-counts/', params).then(resp => {
                    this.alertByUser.series = [
                        {
                            type: 'pie',
                            radius: '55%',
                            center: ['50%', '50%'],
                            labelLine: {smooth: 0.2, length: 10, length2: 20},
                            animationType: 'scale',
                            animationEasing: 'elasticOut',
                            data: resp.data.map(s => { return {name: s.user_name, value: s.total}}),
                        }
                    ];
                }).catch(error => {this.ToastError(error)});

                axios.get('/api/statistics/rule-group-counts/', params).then(resp => {
                    this.alertByRule.series = [
                        {
                            type: 'pie',
                            radius: '55%',
                            center: ['50%', '50%'],
                            labelLine: {smooth: 0.2, length: 10, length2: 20},
                            animationType: 'scale',
                            animationEasing: 'elasticOut',
                            data: resp.data.map(s => { return {name: s.rule_name, value: s.total}}),
                        }
                    ];
                }).catch(error => {this.ToastError(error)});
            }
        },
        mounted() {
            this.reload();
        }
    }
</script>

<style scoped>

</style>