import Vue from 'vue';
import Router from 'vue-router';
import Groups from './views/Groups';
import PendingMessages from "./views/PendingMessages";
import Queue from './views/Queue';
import Settings from "./views/Settings";
import Error from "./views/Error";

Vue.use(Router);

export default new Router({
  routes: [
    {path: '/', component: Groups},
    {path: '/pending-messages', component: PendingMessages},
    {path: '/queue', component: Queue},
    {path: '/settings', component: Settings},
    {path: '/errors/', component: Error},
  ]
});
