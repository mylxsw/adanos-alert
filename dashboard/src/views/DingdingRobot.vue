<template>
    <b-row class="mb-5">
        <b-col>
            <b-card class="mb-2">
                <b-card-text style="display: flex; justify-content:space-between">
                    <b-form inline @submit="searchSubmit">
                        <b-input class="mb-2 mr-sm-2 mb-sm-0" placeholder="名称" v-model="search.name"></b-input>
                        <b-button variant="light" type="submit">搜索</b-button>
                    </b-form>
                    <b-button to="/dingding-robots/add" variant="primary">新增机器人</b-button>
                </b-card-text>
            </b-card>
            <b-table :items="robots" :fields="fields" :busy="isBusy" show-empty hover>
                <template v-slot:cell(updated_at)="row">
                    <date-time :value="row.item.updated_at"></date-time>
                </template>
                <template v-slot:table-busy class="text-center text-danger my-2">
                    <b-spinner class="align-middle"></b-spinner>
                    <strong> Loading...</strong>
                </template>
                <template v-slot:cell(operations)="row">
                    <b-button-group class="mr-2">
                        <b-button size="sm" variant="success" :to="{path:'/', query:{dingding_id: row.item.id}}">报警</b-button>
                        <b-button size="sm" variant="dark" :to="{path:'/rules', query:{dingding_id: row.item.id}}">规则</b-button>
                    </b-button-group>
                    <b-button-group>

                        <b-button size="sm" variant="info" :to="{path:'/dingding-robots/' + row.item.id + '/edit'}">编辑</b-button>
                        <b-button size="sm" variant="danger" @click="delete_robot(row.index, row.item.id)">删除</b-button>
                    </b-button-group>
                </template>
            </b-table>
            <paginator :per_page="10" :cur="cur" :next="next" path="/dingding-robots" :query="{}"></paginator>
        </b-col>
    </b-row>
</template>

<script>
    import axios from 'axios'

    export default {
        name: 'DingdingRobot',
        data() {
            return {
                robots: [],
                cur: parseInt(this.$route.query.next !== undefined ? this.$route.query.next : 0),
                next: -1,
                isBusy: true,
                fields: [
                    {key: 'id', label: 'ID'},
                    {key: 'name', label: '名称'},
                    {key: 'token', label: 'Token'},
                    {key: 'updated_at', label: '最后更新'},
                    {key: 'operations', label: '操作'}
                ],
                search: {
                    name: '',
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

                this.$router.push({path: '/dingding-robots', query: query}).catch(err => {
                    this.ToastError(err)
                });
            },
            delete_robot(index, id) {
                let self = this;
                this.$bvModal.msgBoxConfirm('确定执行该操作 ?').then((value) => {
                    if (value !== true) {
                        return;
                    }

                    axios.delete('/api/dingding-robots/' + id + '/').then(() => {
                        self.robots.splice(index, 1);
                        this.SuccessBox('操作成功')
                    }).catch(error => {
                        this.ErrorBox(error);
                    });
                });
            },
            reload() {
                let params = this.$route.query;
                params.offset = this.cur;
                axios.get('/api/dingding-robots/', {
                    params: params
                }).then(response => {
                    this.robots = response.data.robots;
                    this.next = response.data.next;
                    this.search.name = response.data.search.name;

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