<template>
    <div class="mt-3">
        <b-button-group class="mx-1">
            <b-button v-if="hasPrevPage()" @click.native="$router.go(0)" :to="{path: this.path, query: prevPageQuery()}">Previous</b-button>
            <b-button v-if="hasNextPage()" @click.native="$router.go(0)" :to="{path: this.path, query: nextPageQuery()}">Next</b-button>
        </b-button-group>
    </div>
</template>

<script>
    export default {
        name: 'Paginator',
        props: {
            cur: Number,
            next: Number,
            query: Object,
            path: String,
            per_page: Number,
        },
        methods: {
            hasPrevPage() {
                return this.cur - this.per_page >= 0;
            },
            hasNextPage() {
                return this.next > 0;
            },
            prevPageQuery() {
                let query = {};
                for (let i in this.query) {
                    query[i] = this.query[i];
                }

                query['next'] = this.cur - this.per_page;

                return query;
            },
            nextPageQuery() {
                let query = {};
                for (let i in this.query) {
                    query[i] = this.query[i];
                }

                query['next'] = this.next;
                return query;
            }
        }
    }
</script>

