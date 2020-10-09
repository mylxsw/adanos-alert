<template>
    <b-row class="mb-5">
        <b-col>
            <b-card class="mb-2">
                <b-card-text>
                    <b-form inline @submit="searchSubmit">
                        <b-input class="mb-2 mr-sm-2 mb-sm-0" placeholder="来源" v-model="search.origin"></b-input>
                        <b-form-tags v-model="search.tags" class="mb-2 mr-sm-2 mb-sm-0" placeholder="标签"></b-form-tags>
                        <b-form-select v-model="search.status" class="mb-2 mr-sm-2 mb-sm-0" placeholder="状态" :options="status_options"></b-form-select>
                        <b-button variant="light" type="submit">搜索</b-button>
                    </b-form>   
                </b-card-text>
            </b-card>
            <EventCard v-for="(message, index) in events" :key="index" class="mb-3"
                         :event="message"
                         :event_index="index"
                         :test-matched-rules="testMatchedRules"></EventCard>
            <b-card v-if="events.length === 0">
                <b-card-body>There are no records to show</b-card-body>
            </b-card>
            <paginator :per_page="2" :cur="cur" :next="next" path="/events" :query="this.$route.query"></paginator>
        </b-col>

        <b-modal id="matched-rules-dialog" title="匹配到的规则" hide-footer size="xl">
            <b-table responsive="true" :items="matched_rules" :fields="matched_rules_fields">
                <template v-slot:cell(name)="row">
                    <span v-b-tooltip.hover :title="row.item.rule.id">{{ row.item.rule.name }}</span>
                    <b-link :to="'/rules/' + row.item.rule.id + '/edit'" target="_blank" class="ml-2">
                        <font-awesome-icon icon="external-link-alt"></font-awesome-icon>
                    </b-link>
                    <p>
                        <b-badge variant="info" v-for="(tag, index) in row.item.rule.tags" :key="index" class="mr-1">{{ tag }}</b-badge>
                    </p>
                </template>
                <template v-slot:cell(rule)="row">
                    <p class="adanos-pre-fold" v-b-tooltip.hover :title="row.item.rule.rule">
                        <code>{{ row.item.rule.rule }}</code>
                    </p>
                </template>
                <template v-slot:cell(aggregate_rule)="row">
                    <p class="adanos-pre-fold" v-b-tooltip.hover :title="row.item.rule.aggregate_rule">
                        <code>{{ row.item.rule.aggregate_rule }}</code>
                    </p>
                    <p><b-badge variant="success" v-if="row.item.aggregate_key" v-b-tooltip.hover title="实际匹配的聚合 Key">{{ row.item.aggregate_key }}</b-badge></p>
                </template>
            </b-table>
        </b-modal>
    </b-row>
</template>

<script>
    import axios from 'axios';

    export default {
        name: 'Events',
        data() {
            return {
                search: {
                    origin: '',
                    status: this.$route.query.status !== undefined ? this.$route.query.status : null,
                    tags: [],
                    meta: '',
                },
                status_options: [
                    {value: null, text: '所有状态'},
                    {value: 'pending', text: '准备中'},
                    {value: 'grouped', text: '已分组'},
                    {value: 'canceled', text: '无规则，已取消'},
                    {value: 'expired', text: '匹配规则，已过期'},
                    {value: 'ignored', text: '匹配规则，已忽略'},
                ],
                events: [],
                cur: parseInt(this.$route.query.next !== undefined ? this.$route.query.next : 0),
                next: -1,
                matched_rules: [],
                matched_rules_fields: [
                    {key: 'name', label: '规则名称/ID'},
                    {key: 'rule', label: '规则'},
                    {key: 'aggregate_rule', label: '聚合条件'},
                ],
            };
        },
        watch: {
            '$route': 'loadMore',
        },
        methods: {
            searchSubmit(evt) {
                evt.preventDefault();

                this.$router.push({path: '/events', query: {
                    offset: 0,
                    group_id: this.$route.query.group_id !== undefined ? this.$route.query.group_id : null,
                    status: this.search.status,
                    tags: this.search.tags.join(),
                    meta: this.search.meta,
                    origin: this.search.origin,
                }}).catch(err => {this.ToastError(err);});
            },
            testMatchedRules(id) {
                axios.post('/api/events/' + id + '/matched-rules/', {}).then(resp => {
                    this.matched_rules = resp.data;
                    this.$root.$emit('bv::show::modal', "matched-rules-dialog");
                }).catch(error => {
                    this.ErrorBox(error);
                });
            },
            loadMore() {
                let params = this.$route.query;
                params.offset = this.cur;
                
                axios.get('/api/events/', {
                    params: params,
                }).then(response => {
                    this.events = response.data.events;
                    for (let i in this.events) {
                        this.events[i]._showDetails = true;
                    }

                    this.next = response.data.next;

                    this.search.origin = response.data.search.origin;
                    this.search.tags = response.data.search.tags;
                    this.search.meta = response.data.search.meta;
                    this.search.status = response.data.search.status.length === 0 ? null : response.data.search.status[0];

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
