<template>
    <b-row class="mb-5">
        <b-col>
            <b-card class="mb-2" header="每日报警次数汇总">
                <b-card-body>
                    <v-charts :options="alertByDatetime" style="width: 100%;"></v-charts>
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

    import axios from "axios";

    export default {
        components: {
            'v-charts': Echarts,
        },
        name: 'Report',
        data() {
            return {
                alertByDatetime: {
                    title: {left: 'center', text: '报警时间分布'},
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
                    series: [
                        {
                            smooth: true,
                            name: '今日报警次数',
                            data: [],
                            type: 'line'
                        }
                    ]
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
                this.alertByDatetime.series[0].data = response.data.map(s => s.total)
            }).catch(error => {
                this.ToastError(error);
            });
        }
    }
</script>

<style scoped>

</style>