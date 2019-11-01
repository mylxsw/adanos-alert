<template>
  <div id="app">
    <b-container fluid>
        <b-navbar type="dark" toggleable="md" variant="primary" class="mb-3" sticky>
            <b-navbar-brand href="/">Adanos</b-navbar-brand>
            <b-collapse is-nav id="nav_dropdown_collapse">
                <b-navbar-nav>
                    <b-nav-item to="/" exact>Groups</b-nav-item>
                    <b-nav-item :to="{path:'/messages', query: {status: 'pending'}}" exact>Pending <b-badge variant="danger" v-if="pending_message_count > 0">{{ pending_message_count }}</b-badge></b-nav-item>
                    <b-nav-item to="/rules" exact>Rules</b-nav-item>
                    <b-nav-item to="/users">Users</b-nav-item>
                    <b-nav-item to="/queues">Queues</b-nav-item>
                    <b-nav-item to="/settings">Settings</b-nav-item>
                </b-navbar-nav>
                <ul class="navbar-nav flex-row ml-md-auto d-none d-md-flex">
                    <li class="nav-item">
                        <a href="https://github.com/mylxsw/adanos-alert" class="text-white">{{ version }}</a>
                    </li>
                </ul>
            </b-collapse>
        </b-navbar>
        <div class="main-view">
            <router-view/>
        </div>
    </b-container>
    
  </div>
</template>

<script>
    import axios from 'axios';

    export default {
        data() {
            return {
                version: 'v-0',
                pending_message_count: 0,
            }
        },
        mounted() {
            axios.get('/api/').then(response => {
                this.version = response.data.version;
            });

            axios.get('/api/messages-count/?status=pending').then(response => {
                this.pending_message_count = response.data.count;
            }).catch(error => {
                this.$bvToast.toast(error.response !== undefined ? error.response.data.error : error.toString(), {
                    title: 'ERROR',
                    variant: 'danger'
                });
            });
        },
        beforeMount() {
            axios.defaults.baseURL = this.$store.getters.serverUrl;
            let token = this.$store.getters.token;
            if (token !== "") {
                axios.defaults.headers.common['Authorization'] = "Bearer " + token;
            }
        }
    }
</script>

<style>
    .container-fluid {
        padding: 0;
    }

    .main-view {
        padding: 15px;
    }
</style>
