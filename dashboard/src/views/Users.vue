<template>
    <b-row class="mb-5">
        <b-col>
            <b-btn-group class="mb-3">
                <b-button to="/users/add" variant="primary">新增用户</b-button>
            </b-btn-group>
            <b-table :items="users" :fields="fields" :busy="isBusy" show-empty>
                <template v-slot:cell(name)="row">
                    {{ row.item.name }}
                    <p><b>{{ row.item.id }}</b></p>
                </template>
                <template v-slot:cell(metas)="row">
                    <b-list-group>
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
                    <b-button-group>
                        <b-button size="sm" variant="dark" :to="{path:'/rules', query:{user_id: row.item.id}}">规则</b-button>
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
                    {key: 'metas', label: '属性'},
                    {key: 'updated_at', label: '最后更新'},
                    {key: 'operations', label: '操作'}
                ],
            };
        },
        methods: {
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
                axios.get('/api/users/?offset=' + this.cur).then(response => {
                    this.users = response.data.users;
                    this.next = response.data.next;
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