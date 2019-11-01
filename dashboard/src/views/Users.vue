<template>
    <b-row class="mb-5">
        <b-col>
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
            </b-table>
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
                next: -1,
                isBusy: true,
                fields: [
                    {key: 'name', label: '用户名'},
                    {key: 'metas', label: '基本信息'},
                    {key: 'status', label: '状态'},
                    {key: 'updated_at', label: '最后更新'},
                    {key: 'operations', label: '操作'}
                ],
            };
        },
        methods: {
            reload() {
                axios.get('/api/users/?next=' + this.next).then(response => {
                    this.users = response.data.users;
                    this.next = response.data.next;
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