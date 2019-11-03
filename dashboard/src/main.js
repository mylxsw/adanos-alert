import '@babel/polyfill'
import 'mutationobserver-shim'
import Vue from 'vue'
import './plugins/axios'
import './plugins/bootstrap-vue'
import App from './App.vue'
import router from './router'
import store from './store'
import DateTime from "./components/DateTime";
import Paginator from "./components/Paginator";

Vue.component('DateTime', DateTime);
Vue.component('Paginator', Paginator);
Vue.config.productionTip = false;

new Vue({
    router,
    store,
    render: h => h(App)
}).$mount('#app');
