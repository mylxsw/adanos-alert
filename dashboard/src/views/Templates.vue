<template>
    <b-row class="mb-5">
        <b-col>
            <b-card class="mb-2">
                <b-card-text style="display: flex; justify-content:space-between">
                    <div>
                        <b-badge :variant="$route.query.type === undefined ? 'primary':''" class="mr-1" :to="'/templates'">全部</b-badge>
                        <b-badge :variant="$route.query.type === 'match_rule' ? 'primary':''" class="mr-1" :to="'/templates?type=match_rule'">分组匹配规则</b-badge>
                        <b-badge :variant="$route.query.type === 'template' ? 'primary':''" class="mr-1" :to="'/templates?type=template'">分组展示模板</b-badge>
                        <b-badge :variant="$route.query.type === 'trigger_rule' ? 'primary':''" class="mr-1" :to="'/templates?type=trigger_rule'">动作触发规则</b-badge>
                        <b-badge :variant="$route.query.type === 'template_dingding' ? 'primary':''" class="mr-1" :to="'/templates?type=template_dingding'">钉钉通知模板</b-badge>
                        <b-badge :variant="$route.query.type === 'template_report' ? 'primary':''" class="mr-1" :to="'/templates?type=template_report'">报告模板</b-badge>
                    </div>
                    <b-button to="/templates/add" variant="primary">新增模板</b-button>
                </b-card-text>
            </b-card>
            <b-table :items="templates" :fields="fields" :busy="isBusy" show-empty hover>
                <template v-slot:cell(name)="row">
                    <b>{{ row.item.name }}</b>
                    <p><i>{{ row.item.description }}</i></p>
                </template>
                <template v-slot:cell(metas)="row">
                    <b-list-group>
                        <b-list-group-item v-for="(m, index) in row.item.metas" :key="index">
                            {{ m.key }} <b class="text-success">: </b> {{ m.value }}
                        </b-list-group-item>
                    </b-list-group>
                </template>
                <template v-slot:cell(type)="row">
                    <b-badge v-if="row.item.type === 'match_rule'" variant="success">分组匹配规则</b-badge>
                    <b-badge v-if="row.item.type === 'template'" variant="info">分组展示模板</b-badge>
                    <b-badge v-if="row.item.type === 'trigger_rule'" variant="dark">动作触发规则</b-badge>
                    <b-badge v-if="row.item.type === 'template_dingding'" variant="info">钉钉通知模板</b-badge>
                    <b-badge v-if="row.item.type === 'template_report'" variant="warning">报告模板</b-badge>
                </template>
                <template v-slot:cell(updated_at)="row">
                    <date-time :value="row.item.updated_at"></date-time>
                </template>
                <template v-slot:cell(content)="row">
                    <code class="adanos-pre-fold">{{ row.item.content }}</code>
                </template>
                <template v-slot:table-busy class="text-center text-danger my-2">
                    <b-spinner class="align-middle"></b-spinner>
                    <strong> Loading...</strong>
                </template>
                <template v-slot:cell(operations)="row">
                    <b-button-group class="mr-2">
                        <b-button size="sm" @click="row.toggleDetails">
                            {{ row.detailsShowing ? '隐藏' : '显示' }}详情
                        </b-button>
                    </b-button-group>
                    <b-button-group>
                        <b-button v-if="!row.item.predefined" size="sm" variant="info" :to="{path:'/templates/' + row.item.id + '/edit'}">编辑</b-button>
                        <b-button v-if="!row.item.predefined" size="sm" variant="danger" @click="delete_template(row.index, row.item.id)">删除</b-button>
                        <b-button v-if="row.item.predefined" size="sm" disabled>预置</b-button>
                    </b-button-group>
                </template>
                <template v-slot:row-details="row">
                    <b-card>
                        <pre><code class="adanos-colorful-code">{{ row.item.content }}</code></pre>
                    </b-card>
                </template>
            </b-table>
        </b-col>
    </b-row>
</template>

<script>
    import axios from 'axios'

    export default {
        name: 'Templates',
        data() {
            return {
                templates: [],
                isBusy: true,
                fields: [
                    {key: 'type', label: '类型'},
                    {key: 'name', label: '名称'},
                    {key: 'content', label: '模板内容'},
                    {key: 'updated_at', label: '最后更新'},
                    {key: 'operations', label: '操作'}
                ],
            };
        },
        watch: {
            '$route': 'reload',
        },
        methods: {
            delete_template(index, id) {
                let self = this;
                this.$bvModal.msgBoxConfirm('确定执行该操作 ?').then((value) => {
                    if (value !== true) {
                        return;
                    }

                    axios.delete('/api/templates/' + id + '/').then(() => {
                        self.templates.splice(index, 1);
                        this.SuccessBox('操作成功')
                    }).catch(error => {
                        this.ErrorBox(error);
                    });
                });
            },
            reload() {
                axios.get('/api/templates/', {params: this.$route.query}).then(response => {
                    this.templates = response.data;
                    this.isBusy = false;
                }).catch(error => {
                    this.ToastError(error);
                });
            }
        },
        mounted() {
            this.reload();
        }
    }
</script>

<style scoped>
    .adanos-pre-fold {
        width: 300px;
        height: 45px;
        overflow: hidden;
        display: inline-block;
        font-size: 70%;
    }
</style>