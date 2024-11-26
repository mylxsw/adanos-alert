<template>
    <b-row class="mb-5">
        <b-col>
            <b-card class="mb-2 search-box">
                <b-card-text style="display: flex; justify-content:space-between">
                    <div>
                        <b-badge :variant="$route.query.type === undefined ? 'primary':''" class="mr-1" :to="'/templates'">All</b-badge>
                        <b-badge :variant="$route.query.type === 'match_rule' ? 'primary':''" class="mr-1" :to="'/templates?type=match_rule'">Event Group Matching Rule</b-badge>
                        <b-badge :variant="$route.query.type === 'template' ? 'primary':''" class="mr-1" :to="'/templates?type=template'">Event Group Display Template</b-badge>
                        <b-badge :variant="$route.query.type === 'trigger_rule' ? 'primary':''" class="mr-1" :to="'/templates?type=trigger_rule'">Action Triggering Rule</b-badge>
                        <b-badge :variant="$route.query.type === 'template_dingding' ? 'primary':''" class="mr-1" :to="'/templates?type=template_dingding'">DingTalk Notification Template</b-badge>
                        <b-badge :variant="$route.query.type === 'template_report' ? 'primary':''" class="mr-1" :to="'/templates?type=template_report'">Report Template</b-badge>
                    </div>
                    <b-button to="/templates/add" variant="primary">New</b-button>
                </b-card-text>
            </b-card>
            <b-table :items="templates" :fields="fields" :busy="isBusy" show-empty hover>
                <template v-slot:cell(name)="row">
                    <b>{{ row.item.name }}</b>
                    <p class="th-autohide-sm"><i>{{ row.item.description }}</i></p>
                </template>
                <template v-slot:cell(metas)="row">
                    <b-list-group>
                        <b-list-group-item v-for="(m, index) in row.item.metas" :key="index">
                            {{ m.key }} <b class="text-success">: </b> {{ m.value }}
                        </b-list-group-item>
                    </b-list-group>
                </template>
                <template v-slot:cell(type)="row">
                    <b-badge v-if="row.item.type === 'match_rule'" variant="success">Event Group Matching Rules</b-badge>
                    <b-badge v-if="row.item.type === 'template'" variant="info">Event Group Display Template</b-badge>
                    <b-badge v-if="row.item.type === 'trigger_rule'" variant="dark">Action Triggering Rule</b-badge>
                    <b-badge v-if="row.item.type === 'template_dingding'" variant="info">DingTalk Notification Template</b-badge>
                    <b-badge v-if="row.item.type === 'template_report'" variant="warning">Report Template</b-badge>
                </template>
                <template v-slot:cell(updated_at)="row">
                    <date-time :value="row.item.updated_at"></date-time>
                </template>
                <template v-slot:cell(content)="row">
                    <code class="adanos-pre-fold">{{ row.item.content }}</code>
                </template>
                <template v-slot:table-busy class="text-center text-danger my-2">
                    <b-spinner class="align-middle"></b-spinner>
                    <strong> Loading...</strong>
                </template>
                <template v-slot:cell(operations)="row">
                    <b-button-group class="mr-2 th-autohide-sm">
                        <b-button size="sm" @click="row.toggleDetails">
                            {{ row.detailsShowing ? 'Hide' : 'Show' }}
                        </b-button>
                    </b-button-group>
                    <b-button-group>
                        <b-button v-if="!row.item.predefined" size="sm" variant="info" :to="{path:'/templates/' + row.item.id + '/edit'}">Edit</b-button>
                        <b-button v-if="!row.item.predefined" class="th-autohide-sm" size="sm" variant="danger" @click="delete_template(row.index, row.item.id)">Delete</b-button>
                        <b-button v-if="row.item.predefined" class="th-autohide-sm" size="sm" disabled>Preset</b-button>
                    </b-button-group>
                </template>
                <template v-slot:row-details="row">
                    <b-card>
                        <pre><code class="adanos-colorful-code">{{ row.item.content }}</code></pre>
                    </b-card>
                </template>
            </b-table>
        </b-col>
    </b-row>
</template>

<script>
    import axios from 'axios'

    export default {
        name: 'Templates',
        data() {
            return {
                templates: [],
                isBusy: true,
                fields: [
                    {key: 'type', label: 'Type', sortable: true, class: 'th-autohide-sm'},
                    {key: 'name', label: 'Name', sortable: true, class: 'th-column-width-limit'},
                    {key: 'content', label: 'Content', class: 'th-autohide-md'},
                    {key: 'updated_at', label: 'Last Updated', sortable: true, class: 'th-autohide-md'},
                    {key: 'operations', label: 'Operations'}
                ],
            };
        },
        watch: {
            '$route': 'reload',
        },
        methods: {
            delete_template(index, id) {
                let self = this;
                this.$bvModal.msgBoxConfirm('Are you sure to perform this operation?').then((value) => {
                    if (value !== true) {
                        return;
                    }

                    axios.delete('/api/templates/' + id + '/').then(() => {
                        self.templates.splice(index, 1);
                        this.SuccessBox('Operation successful')
                    }).catch(error => {
                        this.ErrorBox(error);
                    });
                });
            },
            reload() {
                axios.get('/api/templates/', {params: this.$route.query}).then(response => {
                    this.templates = response.data;
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

<style scoped>
    .adanos-pre-fold {
        width: 300px;
        height: 45px;
        overflow: hidden;
        display: inline-block;
        font-size: 70%;
    }

</style>