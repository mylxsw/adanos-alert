<template>
    <b-row class="mb-5">
        <b-col>
            <v-charts :options="alertByDatetime" style="width: 100%;"></v-charts>
            <v-charts :options="eventsByDatetime" style="width: 100%;" class="mt-3"></v-charts>
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
    import { graphic } from 'echarts/lib/export'

    import axios from "axios";

    export default {
        components: {
            'v-charts': Echarts,
        },
        name: 'Report',
        data() {
            return {
                alertByDatetime: {
                    title: {left: 'left', text: '事件组数量时间分布'},
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
                　　　　show:true,
　　　　                feature:{
　　　　　                  dataZoom: {
　　　　　　                    yAxisIndex:"none"
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
                　　　　show:true,
　　　　                feature:{
　　　　　                  dataZoom: {
　　　　　　                    yAxisIndex:"none"
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
        methods: {},
        mounted() {
            axios.get('/api/statistics/daily-group-counts/').then(response => {
                this.alertByDatetime.xAxis.data = response.data.map(s => s.datetime);
                this.alertByDatetime.series.data = response.data.map(s => s.total)
            }).catch(error => {this.ToastError(error)});

            axios.get('/api/statistics/events/period-counts/').then(resp => {
                this.eventsByDatetime.xAxis.data = resp.data.map(s => s.datetime);
                this.eventsByDatetime.series.data = resp.data.map(s => s.total);
            }).catch(error => {this.ToastError(error)});
        }
    }
</script>

<style scoped>

</style>