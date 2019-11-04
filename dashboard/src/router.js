import Vue from 'vue';
import Router from 'vue-router';
import Groups from './views/Groups';
import Messages from "./views/Messages";
import Queue from './views/Queue';
import Settings from "./views/Settings";
import Error from "./views/Error";
import Rules from "./views/Rules";
import Users from "./views/Users";
import UserEdit from "./views/UserEdit";
import RuleEdit from "./views/RuleEdit";
import Templates from "./views/Templates";
import TemplateEdit from "./views/TemplateEdit";

Vue.use(Router);

export default new Router({
    routes: [
        {path: '/', component: Groups},
        {path: '/messages', component: Messages},
        {path: '/rules', component: Rules},
        {path: '/rules/add', component: RuleEdit},
        {path: '/rules/:id/edit', component: RuleEdit},
        {path: '/users', component: Users},
        {path: '/users/add', component: UserEdit},
        {path: '/users/:id/edit', component: UserEdit},
        {path: '/queues', component: Queue},
        {path: '/templates', component: Templates},
        {path: '/templates/add', component: TemplateEdit},
        {path: '/templates/:id/edit', component: TemplateEdit},
        {path: '/settings', component: Settings},
        {path: '/errors/', component: Error},
    ]
});
