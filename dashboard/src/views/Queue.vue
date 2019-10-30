<template>
    <b-row class="mb-5">
        <b-col>
            <b-table :items="jobs" :fields="fields" :busy="isBusy" show-empty>
                <template v-slot:cell(name)="row">
                    {{ row.item.payload.name }}
                </template>
                <template v-slot:cell(payload)="row">
                    {{ row.item.payload.from }}
                </template>
                <template v-slot:empty="scope">
                    {{ scope.emptyText }}
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
    import axios from 'axios';

    export default {
        name: 'Queue',
        data() {
            return {
                jobs: [],
                isBusy: true,
                fields: [
                    {key: "id", label: "ID"},
                    {key: "created_at", label: "Time"},
                    {key: "name", label: "Name"},
                    {key: "payload", label: "Sync"},
                    {key: "payload", label: "From"},
                ],
            };
        },
        mounted() {
            axios.get('/api/jobs/').then(response => {
                this.jobs = response.data;
                this.isBusy = false;
            }).catch(error => {
                this.$bvToast.toast(error.response !== undefined ? error.response.data.error : error.toString(), {
                    title: 'ERROR',
                    variant: 'danger'
                });
            });
        }
    }
</script>
