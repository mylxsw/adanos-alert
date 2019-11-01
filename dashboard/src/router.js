import Vue from 'vue';
import Router from 'vue-router';
import Groups from './views/Groups';
import Messages from "./views/Messages";
import Queue from './views/Queue';
import Settings from "./views/Settings";
import Error from "./views/Error";
import Rules from "./views/Rules";
import Users from "./views/Users";

Vue.use(Router);

export default new Router({
  routes: [
    {path: '/', component: Groups},
    {path: '/messages', component: Messages},
    {path: '/rules', component: Rules},
    {path: '/users', component: Users},
    {path: '/queues', component: Queue},
    {path: '/settings', component: Settings},
    {path: '/errors/', component: Error},
  ]
});
