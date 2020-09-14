<template>
    <b-row class="mb-5">
        <b-col>
            <b-card class="mb-2">
                <b-card-text style="display: flex; justify-content:space-between">
                    <b-form inline @submit="searchSubmit">
                        <b-input class="mb-2 mr-sm-2 mb-sm-0" placeholder="用户名" v-model="search.name"></b-input>
                        <b-input class="mb-2 mr-sm-2 mb-sm-0" placeholder="邮箱" v-model="search.email"></b-input>
                        <b-input class="mb-2 mr-sm-2 mb-sm-0" placeholder="手机号码" v-model="search.phone"></b-input>
                        <b-button variant="light" type="submit">搜索</b-button>
                    </b-form>
                    <b-button to="/users/add" variant="primary">新增用户</b-button>
                </b-card-text>
            </b-card>
            <b-table :items="users" :fields="fields" :busy="isBusy" show-empty hover>
                <template v-slot:cell(name)="row">
                    {{ row.item.name }}
                    <p><b>{{ row.item.id }}</b></p>
                </template>
                <template v-slot:cell(metas)="row">
                    <b-list-group style="font-size: 90%">
                        <b-list-group-item v-for="(m, index) in row.item.metas" :key="index">
                            {{ m.key }} <b class="text-success">: </b> {{ m.value }}
                        </b-list-group-item>
                    </b-list-group>
                </template>
                <template v-slot:cell(status)="row">
                    <b-badge v-if="row.item.status === 'enabled'" variant="success">启用</b-badge>
                    <b-badge v-if="row.item.status === 'disabled'" variant="warning">禁用</b-badge>
                </template>
                <template v-slot:cell(updated_at)="row">
                    <date-time :value="row.item.updated_at"></date-time>
                </template>
                <template v-slot:table-busy class="text-center text-danger my-2">
                    <b-spinner class="align-middle"></b-spinner>
                    <strong> Loading...</strong>
                </template>
                <template v-slot:cell(operations)="row">
                    <b-button-group class="mr-2">
                        <b-button size="sm" variant="success" :to="{path:'/', query:{user_id: row.item.id}}">报警</b-button>
                        <b-button size="sm" variant="dark" :to="{path:'/rules', query:{user_id: row.item.id}}">规则</b-button>
                    </b-button-group>
                    <b-button-group>
                        <b-button size="sm" variant="info" :to="{path:'/users/' + row.item.id + '/edit'}">编辑</b-button>
                        <b-button size="sm" variant="danger" @click="delete_user(row.index, row.item.id)">删除</b-button>
                    </b-button-group>
                </template>
            </b-table>
            <paginator :per_page="10" :cur="cur" :next="next" path="/users" :query="{}"></paginator>
        </b-col>
    </b-row>
</template>

<script>
    import axios from 'axios'

    export default {
        name: 'Users',
        data() {
            return {
                users: [],
                cur: parseInt(this.$route.query.next !== undefined ? this.$route.query.next : 0),
                next: -1,
                isBusy: true,
                fields: [
                    {key: 'name', label: '用户名/ID'},
                    {key: 'email', label: '邮箱账号'},
                    {key: 'phone', label: '手机号码'},
                    {key: 'metas', label: '属性'},
                    {key: 'updated_at', label: '最后更新'},
                    {key: 'operations', label: '操作'}
                ],
                search: {
                    name: '',
                    phone: '',
                    email: '',
                },
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

                this.$router.push({path: '/users', query: query}).catch(err => {
                    this.ToastError(err)
                });
            },
            delete_user(index, id) {
                let self = this;
                this.$bvModal.msgBoxConfirm('确定执行该操作 ?').then((value) => {
                    if (value !== true) {
                        return;
                    }

                    axios.delete('/api/users/' + id + '/').then(() => {
                        self.users.splice(index, 1);
                        this.SuccessBox('操作成功')
                    }).catch(error => {
                        this.ErrorBox(error);
                    });
                });
            },
            reload() {
                let params = this.$route.query;
                params.offset = this.cur;
                axios.get('/api/users/', {
                    params: params
                }).then(response => {
                    this.users = response.data.users;
                    this.next = response.data.next;
                    this.search.name = response.data.search.name;
                    this.search.phone = response.data.search.phone;
                    this.search.email = response.data.search.email;
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