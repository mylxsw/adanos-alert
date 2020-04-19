<template>
    <b-row class="mb-5">
        <b-col>
            <b-card class="mb-2">
                <b-card-text style="display: flex; justify-content:space-between">
                    <b-form inline @submit="searchSubmit">
                        <b-input class="mb-2 mr-sm-2 mb-sm-0" placeholder="名称" v-model="search.name"></b-input>
                        <b-form-select v-model="search.status" class="mb-2 mr-sm-2 mb-sm-0" placeholder="状态" :options="status_options"></b-form-select>
                        <b-button variant="light" type="submit">搜索</b-button>
                    </b-form>
                    <b-button to="/rules/add" variant="primary" class="float-right">新增规则</b-button>
                </b-card-text>
            </b-card>
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
            <paginator :per_page="10" :cur="cur" :next="next" path="/rules" :query="this.$route.query"></paginator>
        </b-col>
    </b-row>
</template>

<script>
    import axios from 'axios'

    export default {
        name: 'Rules',
        data() {
            return {
                search: {
                    name:  '',
                    status: '',
                    user_id: '',
                },
                status_options: [
                    {value: null, text: '所有状态'},
                    {value: 'enabled', text: '已启用'},
                    {value: 'disabled', text: '已禁用'},
                ],
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
                cur: parseInt(this.$route.query.next !== undefined ? this.$route.query.next : 0),
                next: -1,
            };
        },
        watch: {
            '$route': 'reload',
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
            searchSubmit(evt) {
                evt.preventDefault();
                var query = this.search;
                query.offset = 0;
                this.$router.push({path: '/rules', query: query}).catch(err => {err});
            },
            reload() {
                var params = this.$route.query;
                params.offset = this.cur;

                axios.get('/api/rules/', {
                    params: params
                }).then(response => {
                    this.rules = response.data.rules;
                    this.next = response.data.next;
                    this.userRefs = response.data.users;
                    this.isBusy = false;

                    this.search.name = response.data.search.name;
                    this.search.status = response.data.search.status.length > 0 ? response.data.search.status : null;
                    this.search.user_id = response.data.search.user_id;
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