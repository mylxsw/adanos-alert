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
                <template v-slot:cell(rule)="row">
                    <p><small>// 报警周期为 <code><b>{{ row.item.interval }} 分钟每次</b></code>{{ row.item.description !== '' ? '，':''}} {{ row.item.description }}</small></p>
                    <p><code>{{ row.item.rule }}</code></p>
                </template>
                <template v-slot:cell(interval)="row">
                    {{ row.item.interval / 60 }}
                </template>
                <template v-slot:cell(triggers)="row">
                    <b-list-group>
                        <b-list-group-item v-for="(trigger, index) in row.item.triggers" :key="index">
                            <code>{{ trigger.name == "" || trigger.name == undefined ? (trigger.pre_condition || 'true') : trigger.name }}</code> <b class="text-success"> | </b> {{
                            formatAction(trigger.action) }} <span v-if="trigger.user_refs.length > 0">({{ users(trigger.user_refs) }})</span>
                        </b-list-group-item>
                    </b-list-group>
                </template>
                <template v-slot:cell(updated_at)="row">
                    <p>
                        <b-badge v-if="row.item.status === 'enabled'" variant="success">已启用</b-badge>
                        <b-badge v-if="row.item.status === 'disabled'" variant="danger">已禁用</b-badge>
                    </p>
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
                userRefs: {},
                isBusy: true,
                fields: [
                    {key: 'name', label: '规则名称/ID'},
                    {key: 'rule', label: '规则'},
                    {key: 'triggers', label: '动作'},
                    {key: 'updated_at', label: '状态/最后更新'},
                    {key: 'operations', label: '操作'}
                ],
            };
        },
        methods: {
            users(user_refs) {
                return user_refs.map((u) => {
                    return this.userRefs[u] !== undefined ? this.userRefs[u] : '-';
                }).join(', ')
            },
            formatAction(action) {
                let actions = {
                    'dingding': '钉钉通知',
                    'email': '邮件通知',
                    'phone_call': '电话通知',
                    'wechat': '微信通知',
                    'sms_aliyun': '阿里云短信',
                    'sms_yunxin': '网易云信',
                    'http': 'HTTP',
                };

                return actions[action];
            },
            delete_rule(index, id) {
                let self = this;
                this.$bvModal.msgBoxConfirm('确定执行该操作 ?').then((value) => {
                    if (value !== true) {
                        return;
                    }

                    axios.delete('/api/rules/' + id + '/').then(() => {
                        self.rules.splice(index, 1);
                        this.SuccessBox('操作成功');
                    }).catch(error => {
                        this.ErrorBox(error);
                    });
                });
            },
            reload() {
                axios.get('/api/rules/', {
                    params: {
                        "user_id": this.$route.query.user_id !== undefined ? this.$route.query.user_id : null,
                    }
                }).then(response => {
                    this.rules = response.data.rules;
                    this.userRefs = response.data.users;
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