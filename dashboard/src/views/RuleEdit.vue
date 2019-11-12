<template>
    <b-row class="mb-5 adanos-input-box">
        <b-col>
            <b-form @submit="onSubmit">
                <b-card-group class="mb-3">
                    <b-card header="基本">
                        <b-form-group label-cols="2" id="rule_name" label="名称*" label-for="name_input">
                            <b-form-input id="name_input" type="text" v-model="form.name" required
                                          placeholder="输入规则名称"></b-form-input>
                        </b-form-group>

                        <b-form-group label-cols="2" id="rule_description" label="描述" label-for="description_input">
                            <b-form-textarea id="description_input" placeholder="输入规则描述"
                                             v-model="form.description"></b-form-textarea>
                        </b-form-group>

                        <b-form-group label-cols="2" id="rule_interval" label="报警周期*" label-for="rule_interval_input"
                                      :description="'当前：' + (parseInt(form.interval) === 0 ? 1 : form.interval) + ' 分钟，每隔 ' + (parseInt(form.interval) === 0 ? 1 : form.interval) + ' 分钟后触发一次报警'">
                            <b-form-input id="rule_interval_input" type="range" min="0" max="1440" step="5"
                                          v-model="form.interval" required></b-form-input>
                        </b-form-group>

                        <b-form-group label-cols="2" id="is_enabled" label="是否启用*" label-for="is_enabled_checkbox">
                            <b-form-checkbox id="is_enabled_checkbox" v-model="form.status">启用</b-form-checkbox>
                        </b-form-group>
                    </b-card>
                </b-card-group>

                <b-card-group class="mb-3">
                    <b-card header="规则">
                        <b-btn-group class="mb-2">
                            <b-btn variant="light" v-b-modal.match_rule_selector>插入模板</b-btn>
                            <b-btn variant="light" @click="rule_help = !rule_help">帮助</b-btn>
                        </b-btn-group>
                        <b-btn-group class="mb-2 float-right">
                            <b-btn variant="primary" class="float-right" @click="checkRule()">检查</b-btn>
                        </b-btn-group>
                        <b-form-textarea id="rule" rows="5" v-model="form.rule"
                                         placeholder="输入规则，必须返回布尔值"></b-form-textarea>
                        <small class="form-text text-muted">
                            语法参考 <a href="https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md"
                                    target="_blank">https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md</a>
                        </small>
                        <MatchRuleHelp v-if="rule_help"></MatchRuleHelp>
                    </b-card>
                </b-card-group>

                <b-card-group class="mb-3">
                    <b-card header="展示模板">
                        <b-btn-group class="mb-2">
                            <b-btn variant="light" v-b-modal.template_selector>插入模板</b-btn>
                            <b-btn variant="light" @click="template_help = !template_help">帮助</b-btn>
                        </b-btn-group>
                        <b-btn-group class="mb-2 float-right">
                            <b-btn variant="primary" class="float-right" @click="checkTemplate()">检查</b-btn>
                        </b-btn-group>
                        <b-form-textarea id="template" rows="5" v-model="form.template"
                                         placeholder="输入模板"></b-form-textarea>
                        <small class="form-text text-muted">
                            语法参考 <a href="https://golang.org/pkg/html/template/" target="_blank">https://golang.org/pkg/html/template/</a>
                        </small>
                        <TemplateHelp v-if="template_help"></TemplateHelp>
                    </b-card>
                </b-card-group>

                <b-card-group class="mb-3">
                    <b-card header="动作">
                        <b-btn variant="success" class="mb-3" @click="triggerAdd()">添加</b-btn>
                        <b-card :header="trigger.id" border-variant="dark" header-bg-variant="dark"
                                header-text-variant="white" class="mb-3" v-bind:key="i"
                                v-for="(trigger, i) in form.triggers">
                            <b-form-group label-cols="2" :id="'trigger_' + i" label="条件"
                                          :label-for="'trigger_pre_condition_' + i">
                                <b-btn-group class="mb-2">
                                    <b-btn variant="light" @click="openTriggerRuleTemplateSelector(i)">插入模板</b-btn>
                                    <b-btn variant="light" @click="toggleHelp(trigger)">帮助</b-btn>
                                </b-btn-group>
                                <b-btn-group class="mb-2 float-right">
                                    <b-btn variant="primary" class="float-right" @click="checkTriggerRule(trigger)">检查
                                    </b-btn>
                                </b-btn-group>
                                <b-form-textarea id="'trigger_pre_condition_' + i" v-model="trigger.pre_condition"
                                                 placeholder="默认为 true （全部匹配）"></b-form-textarea>
                                <small class="form-text text-muted">
                                    语法参考 <a
                                        href="https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md"
                                        target="_blank">https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md</a>
                                </small>
                                <TriggerHelp class="mt-2" v-if="trigger.help"></TriggerHelp>
                            </b-form-group>
                            <b-form-group label-cols="2" :id="'trigger_action_' + i" label="动作"
                                          :label-for="'trigger_action_' + i">
                                <b-form-select :id="'trigger_action_' + i" v-model="trigger.action"
                                               :options="action_options"></b-form-select>
                            </b-form-group>
                            <b-form-group label-cols="2" :id="'trigger_meta_' + i" label="动作参数"
                                          :label-for="'trigger_meta_' + i">
                                <b-form-input :id="'trigger_meta_' + i" v-model="trigger.meta"></b-form-input>
                            </b-form-group>

                            <b-form-group label-cols="2" label="接收人" :label-for="'trigger_users_' + i">
                                <b-btn variant="info" class="mb-3" @click="userAdd(i)">添加接收人</b-btn>
                                <b-input-group v-bind:key="index" v-for="(user, index) in trigger.user_refs"
                                               class="mb-3">
                                    <b-form-select v-model="trigger.user_refs[index]"
                                                   :options="user_options"></b-form-select>
                                    <b-input-group-append>
                                        <b-btn variant="danger" @click="userDelete(i, index)">删除</b-btn>
                                    </b-input-group-append>
                                </b-input-group>
                            </b-form-group>

                            <b-btn class="float-right" variant="danger" @click="triggerDelete(i)">删除动作</b-btn>
                        </b-card>
                    </b-card>
                </b-card-group>

                <b-button type="submit" variant="primary" class="mr-2">保存</b-button>
                <b-button to="/rules">返回</b-button>
            </b-form>

            <b-modal id="match_rule_selector" title="选择分组匹配规则模板" hide-footer size="xl">
                <b-table :items="templates.match_rule" :fields="template_fields">
                    <template v-slot:cell(content)="row">
                        <code>{{ row.item.content }}</code>
                    </template>
                    <template v-slot:cell(operations)="row">
                        <b-button-group>
                            <b-button size="sm" variant="info" @click="applyTemplateForMatchRule(row.item.content)">选中
                            </b-button>
                        </b-button-group>
                    </template>
                </b-table>
            </b-modal>
            <b-modal id="template_selector" title="选择分组展示模板" hide-footer size="xl">
                <b-table :items="templates.template" :fields="template_fields">
                    <template v-slot:cell(content)="row">
                        <code>{{ row.item.content }}</code>
                    </template>
                    <template v-slot:cell(operations)="row">
                        <b-button-group>
                            <b-button size="sm" variant="info" @click="applyTemplateForTemplate(row.item.content)">选中
                            </b-button>
                        </b-button-group>
                    </template>
                </b-table>
            </b-modal>
            <b-modal id="trigger_rule_selector" title="选择动作触发规则模板" hide-footer size="xl">
                <b-table :items="templates.trigger_rule" :fields="template_fields">
                    <template v-slot:cell(content)="row">
                        <code>{{ row.item.content }}</code>
                    </template>
                    <template v-slot:cell(operations)="row">
                        <b-button-group>
                            <b-button size="sm" variant="info" @click="applyTemplateForTriggerRule(row.item.content)">
                                选中
                            </b-button>
                        </b-button-group>
                    </template>
                </b-table>
            </b-modal>
        </b-col>
    </b-row>
</template>

<script>
    import axios from 'axios';
    import MatchRuleHelp from "../components/MatchRuleHelp";
    import TemplateHelp from "../components/TemplateHelp";
    import TriggerHelp from "../components/TriggerHelp";

    export default {
        name: 'RuleEdit',
        components: {TriggerHelp, TemplateHelp, MatchRuleHelp},
        data() {
            return {
                form: {
                    name: '',
                    description: '',
                    interval: 1,
                    rule: '',
                    template: '',
                    triggers: [],
                    status: true,
                },
                rule_help: false,
                template_help: false,
                properties: ['phone', 'email',],
                action_options: [
                    {value: 'dingding', text: '钉钉'},
                    {value: 'http', text: 'HTTP'},
                    {value: 'email', text: '邮件'},
                    {value: 'wechat', text: '微信'},
                    {value: 'sms', text: '短信'},
                    {value: 'phone_call', text: '电话'},
                ],
                user_options: [],
                template_fields: [
                    {key: 'name', label: '名称'},
                    {key: 'description', label: '说明'},
                    {key: 'content', label: '模板内容'},
                    {key: 'operations', label: '操作'},
                ],
                templates: {
                    match_rule: [],
                    trigger_rule: [],
                    template: [],
                },
                currentTriggerRuleId: -1,
            };
        },
        methods: {
            /**
             * 检查匹配规则是否合法
             */
            checkRule() {
                if (this.form.rule.trim() === '') {
                    this.ErrorBox('规则为空，无需检查');
                    return;
                }

                this.sendCheckRequest('match_rule', this.form.rule.trim());
            },

            checkTriggerRule(trigger) {
                let rule = trigger.pre_condition.trim();
                if (rule === '') {
                    this.ErrorBox('动作触发条件为空，无需检查');
                    return;
                }

                this.sendCheckRequest('trigger_rule', rule)
            },

            /**
             * 检查模板是否合法
             */
            checkTemplate() {
                if (this.form.template.trim() === '') {
                    this.ErrorBox('展示模板为空，无需检查');
                    return;
                }

                this.sendCheckRequest('template', this.form.template.trim());
            },

            /**
             * 发送规则检查请求
             */
            sendCheckRequest(type, content) {
                axios.post('/api/rules-test/rule-check/' + type + '/', {content: content}).then(resp => {
                    if (resp.data.error === null || resp.data.error === "") {
                        this.SuccessBox('检查通过');
                    } else {
                        this.ErrorBox('检查不通过：' + resp.data.error);
                    }
                }).catch(error => {
                    this.ErrorBox(error);
                });
            },

            /**
             * 展开帮助信息
             */
            toggleHelp(trigger) {
                trigger.help = !trigger.help;
            },
            /**
             * 打开动作触发规则模板选择对话框
             * @param index
             */
            openTriggerRuleTemplateSelector(index) {
                this.currentTriggerRuleId = index;
                this.$root.$emit('bv::show::modal', "trigger_rule_selector");
            },
            /**
             * 动作触发规则模板选择
             * @param template
             */
            applyTemplateForTriggerRule(template) {
                if (this.form.triggers[this.currentTriggerRuleId].pre_condition.trim() === '') {
                    this.form.triggers[this.currentTriggerRuleId].pre_condition = template;
                } else {
                    this.form.triggers[this.currentTriggerRuleId].pre_condition += ' and ' + template;
                }
                this.$bvModal.hide('trigger_rule_selector');
            },
            /**
             * 展示模板选择
             * @param template
             */
            applyTemplateForTemplate(template) {
                if (this.form.template.trim() === '') {
                    this.form.template = template;
                } else {
                    this.form.template += '\n' + template;
                }
                this.$bvModal.hide('template_selector');
            },
            /**
             * 分组匹配规则模板选择
             * @param template
             */
            applyTemplateForMatchRule(template) {
                if (this.form.rule.trim() === '') {
                    this.form.rule = template;
                } else {
                    this.form.rule += ' and ' + template;
                }
                this.$bvModal.hide('match_rule_selector');
            },
            /**
             * 为动作添加用户
             */
            userAdd(triggerIndex) {
                this.form.triggers[triggerIndex].user_refs.push('');
            },
            /**
             * 为动作移除用户
             */
            userDelete(triggerIndex, index) {
                this.form.triggers[triggerIndex].user_refs.splice(index, 1);
            },
            /**
             * 添加动作
             */
            triggerAdd() {
                this.form.triggers.push({
                    pre_condition: '',
                    action: 'dingding',
                    meta: '',
                    id: '',
                    user_refs: [],
                    help: false
                });
            },
            /**
             * 删除动作
             * @param index
             */
            triggerDelete(index) {
                this.form.triggers.splice(index, 1);
            },
            /**
             * 保存
             * @param evt
             */
            onSubmit(evt) {
                evt.preventDefault();
                let url;
                if (this.$route.params.id !== undefined) {
                    url = '/api/rules/' + this.$route.params.id + '/';
                } else {
                    url = '/api/rules/';
                }

                axios.post(url, this.createRequest()).then(() => {
                    this.SuccessBox('操作成功', () => {
                        window.location.reload(true);
                    });
                }).catch((error) => {
                    this.ErrorBox(error)
                });
            },
            /**
             * 创建请求对象
             */
            createRequest() {
                let requestData = {};
                requestData.name = this.form.name;
                requestData.description = this.form.description;
                requestData.interval = this.form.interval * 60;
                requestData.rule = this.form.rule;
                requestData.template = this.form.template;
                requestData.triggers = this.form.triggers;
                requestData.status = this.form.status ? 'enabled' : 'disabled';

                return requestData;
            },
        },
        mounted() {
            if (this.$route.params.id !== undefined) {
                axios.get('/api/rules/' + this.$route.params.id + '/').then(response => {
                    this.form.name = response.data.name;
                    this.form.description = response.data.description;
                    this.form.interval = response.data.interval / 60;
                    this.form.rule = response.data.rule;
                    this.form.template = response.data.template;

                    for (let i in response.data.triggers) {
                        let trigger = response.data.triggers[i];
                        trigger.help = false;
                        this.form.triggers.push(trigger);
                    }

                    this.form.status = response.data.status === 'enabled';
                }).catch((error) => {
                    this.ToastError(error)
                });
            }

            // 加载辅助元素
            axios.all([
                axios.get('/api/users-helper/names/'),
                axios.get('/api/templates/'),
            ]).then(axios.spread((usersResp, templateResp) => {
                this.user_options = usersResp.data.map((val) => {
                    return {value: val.id, text: val.name}
                });

                for (let i in templateResp.data) {
                    this.templates[templateResp.data[i].type].push(templateResp.data[i]);
                }
            })).catch((error) => {
                this.ToastError(error)
            });
        }
    }
</script>

<style>
    .adanos-input-box {
        max-width: 1000px;
    }
</style>