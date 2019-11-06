<template>
    <b-row class="mb-5">
        <b-col>
            <MessageCard v-for="(message, index) in messages" :key="index" class="mb-3" :message="message"></MessageCard>
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
                messages: [],
                cur: parseInt(this.$route.query.next !== undefined ? this.$route.query.next : 0),
                next: -1,
            };
        },
        methods: {
            loadMore() {
                axios.get('/api/messages/', {
                    params: {
                        offset: this.cur,
                        status: this.$route.query.status !== undefined ? this.$route.query.status : null,
                        group_id: this.$route.query.group_id !== undefined ? this.$route.query.group_id : null,
                    },
                }).then(response => {
                    this.messages = response.data.messages;
                    for (let i in this.messages) {
                        this.messages[i]._showDetails = true;
                    }

                    this.next = response.data.next;
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
