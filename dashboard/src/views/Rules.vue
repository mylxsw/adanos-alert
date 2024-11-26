<template>
    <b-row class="mb-5">
        <b-col>
            <b-card class="mb-2 search-box">
                <b-card-text style="display: flex; justify-content:space-between">
                    <b-form inline @submit="searchSubmit">
                        <b-input class="mb-2 mr-sm-2 mb-sm-0" placeholder="Name" v-model="search.name"></b-input>
                        <b-form-select v-model="search.status" class="mb-2 mr-sm-2 mb-sm-0" placeholder="Status"
                                       :options="status_options"></b-form-select>
                        <b-button variant="light" type="submit">Search</b-button>
                    </b-form>
                    <b-button to="/rules/add" variant="primary" class="float-right">New Rule</b-button>
                </b-card-text>
                <b-card-text>
                    <b-badge :variant="$route.query.tag === undefined ? 'primary':''" class="mr-1" :to="'/rules'">All
                    </b-badge>
                    <b-badge :variant="$route.query.tag === tag.name ? 'primary': ''" v-for="(tag, index) in tags"
                             :key="index" class="mr-1" :to="'/rules?tag=' + tag.name">{{ tag.name }}({{ tag.count }})
                    </b-badge>
                </b-card-text>
            </b-card>
            <b-table :items="rules" :fields="fields" :busy="isBusy" show-empty responsive="true" hover>
                <template v-slot:cell(name)="row">
                    <span v-b-tooltip.hover :title="row.item.description">{{ row.item.name }}</span>
                    <p>
                        <b-badge class="mr-1" variant="danger" v-if="row.item.realtime_interval > 0">ASAP: {{ row.item.realtime_interval }}m</b-badge>
                        <b-badge :variant="$route.query.tag === tag ? 'primary': 'info'"
                                 v-for="(tag, index) in row.item.tags" :key="index" class="mr-1"
                                 :to="'/rules?tag=' + tag">{{ tag }}
                        </b-badge>
                    </p>
                </template>
                <template v-slot:cell(rule)="row">
                    <p><small>
                      The aggregation cycle is
                        <span v-if="row.item.ready_type === 'interval' || row.item.ready_type === ''"><code><b>{{ row.item.interval / 60 }} Per minute</b></code></span>
                        <span v-if="row.item.ready_type === 'daily_time'">
                            <code><b>Daily {{ row.item.daily_times.map((t) => t.substring(0, 5)).join(", ") }}</b></code>
                        </span>
                        <span v-if="row.item.ready_type === 'time_range'"><code><b>{{ timeRangeDesc(row.item.time_ranges)  }}</b></code></span>
                        {{ row.item.description !== '' ? '，' : '' }} {{ row.item.description }}</small>
                    </p>
                    <p class="adanos-pre-fold" v-b-tooltip.hover :title="row.item.rule">
                        <code>{{ row.item.rule }}</code>
                    </p>
                </template>
                <template v-slot:cell(interval)="row">
                    {{ row.item.interval / 60 }}
                </template>
                <template v-slot:cell(triggers)="row">
                    <b-list-group>
                        <b-list-group-item style="padding: 5px;" v-for="(trigger, index) in row.item.triggers" :key="index">
                            <code style="line-height: 0.9" v-b-tooltip :title="trigger.pre_condition" class="action-pre-condition" v-if="!trigger.is_else_trigger">
                                {{ trigger.name === "" || trigger.name === undefined ? (trigger.pre_condition || 'All') : trigger.name }}
                            </code>
                            <code style="line-height: 0.9" v-if="trigger.is_else_trigger">Catch All</code>
                            <b :class="trigger.is_else_trigger ? 'text-warning' : 'text-success'"> | </b> {{ formatAction(trigger.action) }}
                            <span v-if="trigger.user_refs.length > 0">({{ users(trigger.user_refs) }})</span>
                        </b-list-group-item>
                    </b-list-group>
                </template>
                <template v-slot:cell(updated_at)="row">
                    <p>
                        <b-badge v-if="row.item.status === 'enabled'" variant="success">Enabled</b-badge>
                        <b-badge v-if="row.item.status === 'disabled'" variant="danger">Disabled</b-badge>
                    </p>
                    <date-time :value="row.item.updated_at"></date-time>
                </template>
                <template v-slot:table-busy class="text-center text-danger my-2">
                    <b-spinner class="align-middle"></b-spinner>
                    <strong> Loading...</strong>
                </template>
                <template v-slot:cell(operations)="row" >
                    <b-button-group class="mr-2 mt-2">
                        <b-button size="sm" variant="success" :to="{path:'/', query:{rule_id: row.item.id}}">Events
                        </b-button>
                        <b-button size="sm" variant="warning" :to="{path:'/rules/add', query: {copy_from: row.item.id}}"
                                  target="_blank">Copy
                        </b-button>
                    </b-button-group>
                    <b-button-group class="mt-2">
                        <b-button size="sm" variant="info" :to="{path:'/rules/' + row.item.id + '/edit'}">Edit</b-button>
                        <b-button size="sm" variant="danger" @click="delete_rule(row.index, row.item.id)">Delete</b-button>
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
                name: '',
                status: '',
                user_id: '',
            },
            status_options: [
                {value: null, text: 'All'},
                {value: 'enabled', text: 'Enabled'},
                {value: 'disabled', text: 'Disabled'},
            ],
            rules: [],
            userRefs: {},
            isBusy: true,
            fields: [
                {key: 'name', label: 'Rule Name'},
                {key: 'rule', label: 'Rule', class: 'th-autohide-md'},
                {key: 'triggers', label: 'Actions', class: 'th-autohide-sm'},
                {key: 'updated_at', label: 'Status/Last Update Time', class: 'th-autohide-sm'},
                {key: 'operations', label: 'Operations'}
            ],
            cur: parseInt(this.$route.query.next !== undefined ? this.$route.query.next : 0),
            next: -1,
            tags: [],
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
                'dingding': 'DingTalk',
                'email': 'Email',
                'phone_call': 'Phone Call(Aliyun)',
                'wechat': 'WeChat',
                'sms_aliyun': 'SMS(Aliyun)',
                'sms_yunxin': 'SMS(Yunxin)',
                'http': 'HTTP',
                'jira': 'JIRA',
            };

            return actions[action];
        },
        delete_rule(index, id) {
            let self = this;
            this.$bvModal.msgBoxConfirm('Are you sure to perform this operation?').then((value) => {
                if (value !== true) {
                    return;
                }

                axios.delete('/api/rules/' + id + '/').then(() => {
                    self.rules.splice(index, 1);
                    this.SuccessBox('Operation successful');
                }).catch(error => {
                    this.ErrorBox(error);
                });
            });
        },
        timeRangeDesc(timeRanges) {
            let results = [];
            for (let i in timeRanges) {
                let startT = timeRanges[i].start_time;
                let endT = timeRanges[i].end_time;
                let interval = timeRanges[i].interval;

                if (startT === '' || endT === '') {
                    continue;
                }

                results.push("[" + startT.substr(0, 5) + " ~ " + endT.substr(0, 5) + "): Every " + interval / 60 + " Minutes");
            }

            return results.join("; ");
        },
        searchSubmit(evt) {
            evt.preventDefault();
            let query = {offset: 0};
            for (let i in this.$route.query) {
                query[i] = this.$route.query[i];
            }

            for (let i in this.search) {
                query[i] = this.search[i];
            }

            this.$router.push({path: '/rules', query: query}).catch(err => {
                this.ToastError(err)
            });
        },
        reload() {
            let params = this.$route.query;
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

            axios.get('/api/rules-meta/tags/').then(response => {
                this.tags = response.data.tags;
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

<style scoped>
.adanos-pre-fold {
    width: 300px;
    max-height: 45px;
    overflow: hidden;
    display: inline-block;
    font-size: 70%;
}

.action-pre-condition {
    max-width: 400px;
    display: inline-block;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}

p {
    margin-bottom: 5px;
}
</style>
