<template>
    <b-row class="mb-5">
        <b-col>
            <b-table :items="messages" :fields="fields" :busy="isBusy" show-empty>
                <template v-slot:cell(id)="row">
                    {{ row.item.created_at }}
                    <p><b>{{ row.item.id }}</b></p>
                </template>
                <template v-slot:cell(meta)="row">
                    <b-list-group>
                        <b-list-group-item v-for="(val, key) in row.item.meta" :key="key">
                            {{ key }} <b class="text-success">:</b> {{ val }}
                        </b-list-group-item>
                    </b-list-group>
                </template>
                <template v-slot:cell(tags)="row">
                    <b-badge v-for="(tag, index) in row.item.tags" :key="index" class="mr-1">{{ tag }}</b-badge>
                </template>
                <template v-slot:cell(status)="row">
                    <b-badge v-if="row.item.status === 'pending'" variant="dark">准备中</b-badge>
                    <b-badge v-if="row.item.status === 'grouped'" variant="success">已分组</b-badge>
                    <b-badge v-if="row.item.status === 'canceled'" variant="warning">已取消</b-badge>
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
        name: 'PendingMessages',
        data() {
            return {
                messages: [],
                nextOffset: -1,
                isBusy: true,
                fields: [
                    {key: 'id', label: '序号'},
                    {key: 'content', label: '内容'},
                    {key: 'meta', label: '元信息'},
                    {key: 'tags', label: '标签'},
                    {key: 'origin', label: '来源'},
                    {key: 'group_ids', label: '分组'},
                    {key: 'status', label: '状态'}
                ],
            };
        },
        methods: {
            loadMore() {
                axios.get('/api/messages/?status=&offset=' + this.nextOffset).then(response => {
                    this.messages = response.data.messages;
                    this.nextOffset = response.data.next;
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
            this.loadMore();
        }
    }
</script>