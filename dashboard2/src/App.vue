<template>
  <div id="app">
    <b-container fluid>
        <b-navbar type="dark" toggleable="md" variant="primary" class="mb-3" sticky>
            <b-navbar-brand href="/">Adanos</b-navbar-brand>
            <b-collapse is-nav id="nav_dropdown_collapse">
                <b-navbar-nav>
                    <b-nav-item to="/" exact>Groups</b-nav-item>
                    <b-nav-item to="/ungrouped-messages" exact>Pending <b-badge variant="danger" v-if="ungroup_msg_count > 0">{{ ungroup_msg_count }}</b-badge></b-nav-item>
                    <b-nav-item to="/rules" exact>Rules</b-nav-item>
                    <b-nav-item to="/dingding-tokens">Alerts</b-nav-item>
                    <b-nav-item to="/users">Users</b-nav-item>
                    <b-nav-item to="/queues">Queues</b-nav-item>
                </b-navbar-nav>
                <ul class="navbar-nav flex-row ml-md-auto d-none d-md-flex">
                    <li class="nav-item">
                        <a href="http://git.yunsom.cn/golang/adanos" class="text-white">{{ adanos_version }}</a>
                    </li>
                </ul>
            </b-collapse>
        </b-navbar>
        <div class="adanos-main">
            <router-view/>
        </div>
    </b-container>
    
  </div>
</template>

<script>
export default {
  data() {
    return {
      ungroup_msg_count: 10,
      adanos_version: 'v-0',
    }
  },
  mounted() {
      this.ungroup_msg_count = 100;
      this.adanos_version = 'v-1000101';
  }
}
</script>

<style>
    .container-fluid {
        padding: 0;
    }

    .adanos-main {
        padding: 15px;
    }
</style>
