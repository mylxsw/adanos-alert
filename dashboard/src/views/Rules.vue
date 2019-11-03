<template>
    <b-row class="mb-5">
        <b-col>
            <b-btn-group class="mb-3">
                <b-button to="/rules/add" variant="primary">新增规则</b-button>
            </b-btn-group>
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
                <template v-slot:cell(operations)="row">
                    <b-button-group>
                        <b-button size="sm" variant="info" :to="{path:'/rules/' + row.item.id + '/edit'}">编辑</b-button>
                        <b-button size="sm" variant="danger" @click="delete_rule(row.index, row.item.id)">删除</b-button>
                    </b-button-group>
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
            delete_rule(index, id) {
                let self = this;
                this.$bvModal.msgBoxConfirm('确定执行该操作 ?').then((value) => {
                    if (value !== true) {
                        return;
                    }

                    axios.delete('/api/rules/' + id + '/').then(() => {
                        self.rules.splice(index, 1);
                        this.$bvToast.toast('操作成功', {
                            title: 'OK',
                            variant: 'success',
                        });
                    }).catch(error => {
                        this.$bvToast.toast(error.response !== undefined ? error.response.data.error : error.toString(), {
                            title: 'ERROR',
                            variant: 'danger'
                        });
                    });
                });
            },
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