<template>
    <b-row class="mb-5">
        <b-col>
            <b-card class="mb-2">
                <b-card-text>
                    <b-form inline @submit="searchSubmit">
                        <b-input class="mb-2 mr-sm-2 mb-sm-0" placeholder="来源" v-model="search.origin"></b-input>
                        <b-form-tags v-model="search.tags" class="mb-2 mr-sm-2 mb-sm-0" placeholder="标签"></b-form-tags>
                        <b-form-select v-model="search.status" class="mb-2 mr-sm-2 mb-sm-0" placeholder="状态" :options="status_options"></b-form-select>
                        <b-button variant="primary" type="submit">搜索</b-button>
                    </b-form>   
                </b-card-text>
            </b-card>
            <MessageCard v-for="(message, index) in messages" :key="index" class="mb-3" :message="message" :message_index="index"></MessageCard>
            <b-card v-if="messages.length === 0">
                <b-card-body>There are no records to show</b-card-body>
            </b-card>
            <paginator :per_page="10" :cur="cur" :next="next" path="/messages" :query="this.$route.query"></paginator>
        </b-col>
    </b-row>
</template>

<script>
    import axios from 'axios';

    export default {
        name: 'Messages',
        data() {
            return {
                search: {
                    origin: '',
                    status: null,
                    tags: [],
                    meta: '',
                },
                status_options: [
                    {value: null, text: '所有状态'},
                    {value: 'pending', text: '准备中'},
                    {value: 'grouped', text: '已分组'},
                    {value: 'canceled', text: '已取消'},
                ],
                messages: [],
                cur: parseInt(this.$route.query.next !== undefined ? this.$route.query.next : 0),
                next: -1,
            };
        },
        watch: {
            '$route': 'loadMore',
        },
        methods: {
            searchSubmit(evt) {
                evt.preventDefault();

                this.$router.push({path: '/messages', query: {
                    offset: 0,
                    group_id: this.$route.query.group_id !== undefined ? this.$route.query.group_id : null,
                    status: this.search.status,
                    tags: this.search.tags.join(),
                    meta: this.search.meta,
                    origin: this.search.origin,
                }}).catch(err => {err});
            },
            loadMore() {
                var params = this.$route.query;
                params.offset = this.cur;
                
                axios.get('/api/messages/', {
                    params: params,
                }).then(response => {
                    this.messages = response.data.messages;
                    for (let i in this.messages) {
                        this.messages[i]._showDetails = true;
                    }

                    this.next = response.data.next;

                    this.search.origin = response.data.search.origin;
                    this.search.tags = response.data.search.tags;
                    this.search.meta = response.data.search.meta;
                    this.search.status = response.data.search.status.length == 0 ? null : response.data.search.status[0];

                    this.isBusy = false;
                }).catch(error => {
                    this.ToastError(error)
                });
            }
        },
        mounted() {
            this.loadMore();
        }
    }
</script>
