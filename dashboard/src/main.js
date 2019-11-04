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
        centered: true,
        okVariant: 'success',
        headerClass: 'p-2 border-bottom-0',
        footerClass: 'p-2 border-top-0',
    }).then(cb);
};

Vue.prototype.ErrorBox = function (message, cb) {
    cb = cb || function () {};
    this.$bvModal.msgBoxOk(this.ParseError(message), {
        centered: true,
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
