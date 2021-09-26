import '@babel/polyfill'
import 'mutationobserver-shim'
import Vue from 'vue'
import './plugins/axios'
import './plugins/bootstrap-vue'
import App from './App.vue'
import router from './router'
import store from './store'

import { BootstrapVueIcons } from 'bootstrap-vue'

import { library } from '@fortawesome/fontawesome-svg-core'
import { faExternalLinkAlt, faPlus } from '@fortawesome/free-solid-svg-icons'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'

import DateTime from "./components/DateTime";
import HumanTime from "./components/HumanTime";
import Paginator from "./components/Paginator";
import EventCard from "./components/EventCard";

library.add(faExternalLinkAlt);
library.add(faPlus)
Vue.component('font-awesome-icon', FontAwesomeIcon);

Vue.use(BootstrapVueIcons);

Vue.component('DateTime', DateTime);
Vue.component('HumanTime', HumanTime);
Vue.component('EventCard', EventCard)
Vue.component('Paginator', Paginator);

Vue.config.productionTip = false;

/**
 * @return {string}
 */
Vue.prototype.ParseError = function (error) {
    if (error.response !== undefined) {
        if (error.response.data !== undefined) {
            return error.response.data.error;
        }
    }

    return error.toString();
};

Vue.prototype.ToastSuccess = function (message) {
    this.$bvToast.toast(message, {
        title: 'OK',
        variant: 'success',
    });
};

Vue.prototype.ToastError = function (message) {
    this.$bvToast.toast(this.ParseError(message), {
        title: 'ERROR',
        variant: 'danger'
    });
};

Vue.prototype.SuccessBox = function (message, cb) {
    cb = cb || function () {};
    this.$bvModal.msgBoxOk(message, {
        title: '操作成功',
        centered: true,
        okVariant: 'success',
        headerClass: 'p-2 border-bottom-0',
        footerClass: 'p-2 border-top-0',
    }).then(cb);
};

Vue.prototype.ErrorBox = function (message, cb) {
    cb = cb || function () {};
    
    let err = this.ParseError(message);
    console.log(err);

    const h = this.$createElement;
    const messageVNode = h('div', {domProps: {
        innerHTML: '<pre style="white-space: pre-wrap; word-wrap: break-word;">' + err + '</pre>'
    }});

    this.$bvModal.msgBoxOk([messageVNode], {
        centered: true,
        title:'出错了',
        okVariant: 'danger',
        headerClass: 'p-2 border-bottom-0',
        footerClass: 'p-2 border-top-0',
    }).then(cb);
};

new Vue({
    router,
    store,
    render: h => h(App)
}).$mount('#app');
