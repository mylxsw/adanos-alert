<template>
  <div id="app">
    <b-container fluid>
        <b-navbar type="dark" toggleable="md" variant="dark" class="mb-3" sticky>
            <b-navbar-brand href="/">Adanos <a href="https://github.com/mylxsw/adanos-alert" class="text-white" style="font-size: 30%">{{ version }}</a></b-navbar-brand>
            <b-collapse is-nav id="nav_dropdown_collapse">
                <ul class="navbar-nav flex-row ml-md-auto d-none d-md-flex"></ul>
                <b-navbar-nav>
                    <b-nav-item href="/" exact>事件组</b-nav-item>
                    <b-nav-item :to="{path:'/events', query: {status: 'canceled'}}" exact>
                        事件
                        <b-badge variant="danger" v-if="pending_events_count > 0" v-b-tooltip.hover title="没有匹配任何规则的事件">{{ pending_events_count }}</b-badge>
                    </b-nav-item>
                    <b-nav-item to="/rules" exact>规则</b-nav-item>
                    <b-nav-item to="/templates">模板</b-nav-item>
                    <b-nav-item to="/queues">队列</b-nav-item>
                    <b-nav-item to="/agents">Agent</b-nav-item>
                    <b-nav-item to="/reports">报表</b-nav-item>
                    <b-nav-item to="/syslog">系统日志</b-nav-item>
                    <b-nav-item to="/users">用户</b-nav-item>
                    <b-nav-item-dropdown text="插件" right>
                        <b-dropdown-item to="/dingding-robots">钉钉</b-dropdown-item>
                        <b-dropdown-item>阿里云语音通知（暂未开放）</b-dropdown-item>
                        <b-dropdown-item to="/debug">Debug</b-dropdown-item>
                    </b-nav-item-dropdown>
                    <b-nav-item to="/settings">设置</b-nav-item>
                </b-navbar-nav>
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
                pending_events_count: 0,
            }
        },
        mounted() {
            axios.get('/api/').then(response => {
                this.version = response.data.version;
            });

            let self = this;
            let updateCanceledMessageCount = function () {
                axios.get('/api/events-count/?status=canceled').then(response => {
                    self.pending_events_count = response.data.count;
                }).catch(error => {
                    self.ToastError(error);
                });
            };

            updateCanceledMessageCount();
            window.setInterval(updateCanceledMessageCount, 10000);
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

    .th-column-width-limit {
        max-width: 300px;
    }

    @media screen and (max-width: 1366px) {
        .th-autohide-md {
            display: none;
        }
    }
    @media screen and (max-width: 768px) {
        .th-autohide-sm {
            display: none;
        }
        .search-box {
            display: none;
        }
    }

</style>
