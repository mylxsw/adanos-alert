import Vue from 'vue';
import Router from 'vue-router';
import Groups from './views/Groups';
import Events from "./views/Events";
import Queue from './views/Queue';
import Settings from "./views/Settings";
import Error from "./views/Error";
import Rules from "./views/Rules";
import Users from "./views/Users";
import UserEdit from "./views/UserEdit";
import RuleEdit from "./views/RuleEdit";
import Templates from "./views/Templates";
import TemplateEdit from "./views/TemplateEdit";
import DingdingRobots from "./views/DingdingRobot";
import DingdingRobotEdit from "./views/DingdingRobotEdit";
import Agents from "./views/Agents"
import Report from "./views/Report"
import Syslog from "./views/Syslog"
import Debug from "./views/Debug"


Vue.use(Router);

const routerPush = Router.prototype.push;
Router.prototype.push = function push(location) {
    return routerPush.call(this, location).catch(error => error)
}

export default new Router({
    routes: [
        {path: '/', component: Groups},
        {path: '/events', component: Events},
        
        {path: '/rules', component: Rules},
        {path: '/rules/add', component: RuleEdit},
        {path: '/rules/:id/edit', component: RuleEdit},
        
        {path: '/users', component: Users},
        {path: '/users/add', component: UserEdit},
        {path: '/users/:id/edit', component: UserEdit},

        {path: '/dingding-robots', component: DingdingRobots},
        {path: '/dingding-robots/add', component: DingdingRobotEdit},
        {path: '/dingding-robots/:id/edit', component: DingdingRobotEdit},
        
        {path: '/templates', component: Templates},
        {path: '/templates/add', component: TemplateEdit},
        {path: '/templates/:id/edit', component: TemplateEdit},
       
        {path: '/queues', component: Queue},
        {path: '/settings', component: Settings},
        {path: '/errors/', component: Error},
        {path: '/agents', component: Agents},
        {path: '/reports', component: Report},
        {path: '/syslog', component: Syslog},
        {path: '/debug', component: Debug},
    ]
});
