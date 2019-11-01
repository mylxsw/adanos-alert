<template>
    <b-row class="mb-5">
        <b-col>
            <b-table :items="rules" :fields="fields" :busy="isBusy" show-empty>
                <template v-slot:cell(name)="row">
                    {{ row.item.name }}
                    <p><b>{{ row.item.id }}</b></p>
                </template>
                <template v-slot:cell(triggers)="row">
                    <b-list-group>
                        <b-list-group-item v-for="(trigger, index) in row.item.triggers" :key="index">
                            {{ trigger.action }}
                        </b-list-group-item>
                    </b-list-group>
                </template>
                <template v-slot:cell(updated_at)="row">
                    <date-time :value="row.item.updated_at"></date-time>
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
        name: 'Rules',
        data() {
            return {
                rules: [],
                isBusy: true,
                fields: [
                    {key: 'name', label: '规则名称'},
                    {key: 'rule', label: '规则'},
                    {key: 'interval', label: '触发周期'},
                    {key: 'triggers', label: '动作'},
                    {key: 'status', label: '状态'},
                    {key: 'updated_at', label: '最后更新'},
                    {key: 'operations', label: '操作'}
                ],
            };
        },
        methods: {
            reload() {
                axios.get('/api/rules/').then(response => {
                    this.rules = response.data;
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
            this.reload();
        }
    }
</script>