<template>
    <b-row class="mb-5 adanos-input-box">
        <b-col>
            <b-form @submit="onSubmit">
                <b-card-group class="mb-3">
                    <b-card header="基本">
                        <b-form-group label-cols="2" id="rule_name" label="名称*" label-for="name_input">
                            <b-form-input id="name_input" type="text" v-model="form.name" required
                                          placeholder="输入规则名称"/>
                        </b-form-group>

                        <b-form-group label-cols="2" id="rule_description" label="描述" label-for="description_input">
                            <b-form-textarea id="description_input" placeholder="输入规则描述"
                                             v-model="form.description"/>
                        </b-form-group>

                        <b-form-group label-cols="2" id="rule_interval" label="报警周期*" label-for="rule_interval_input"
                                      :description="'当前：' + (parseInt(form.interval) === 0 ? 1 : form.interval) + ' 分钟，每隔 ' + (parseInt(form.interval) === 0 ? 1 : form.interval) + ' 分钟后触发一次报警'">
                            <b-form-input id="rule_interval_input" type="range" min="0" max="1440" step="5"
                                          v-model="form.interval" required/>
                        </b-form-group>

                        <b-form-group label-cols="2" id="is_enabled" label="是否启用*" label-for="is_enabled_checkbox">
                            <b-form-checkbox id="is_enabled_checkbox" v-model="form.status">启用</b-form-checkbox>
                        </b-form-group>
                    </b-card>
                </b-card-group>

                <b-card-group class="mb-3">
                    <b-card header="规则">
                        <p class="text-muted">分组匹配规则，作用于单条 message，用于判断该 message 是否与当前规则匹配。
                            <br />如果 message 没有匹配任何规则，将会被标记为 <code>已取消</code>。</p>
                        <b-btn-group class="mb-2">
                            <b-btn variant="warning" v-b-modal.match_rule_selector>插入模板</b-btn>
                            <b-btn variant="dark" @click="rule_help = !rule_help">帮助</b-btn>
                        </b-btn-group>
                        <b-btn-group class="mb-2 float-right">
                            <b-btn variant="primary" class="float-right" @click="checkRule(form.template)">检查</b-btn>
                        </b-btn-group>
                        <b-form-textarea id="rule" rows="5" v-model="form.rule" class="adanos-code-textarea  text-monospace"
                                         placeholder="输入规则，必须返回布尔值"/>
                        <small class="form-text text-muted">
                            语法参考 <a href="https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md"
                                    target="_blank">https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md</a>
                        </small>
                        <MatchRuleHelp v-if="rule_help"/>
                    </b-card>
                </b-card-group>

                <b-card-group class="mb-3">
                    <b-card header="展示模板">
                        <p class="text-muted">展示模板，用于各通知方式的默认展示模板，Adanos 会按照该模板将分组信息发送给接收人。</p>
                        <b-btn-group class="mb-2">
                            <b-btn variant="warning" v-b-modal.template_selector>插入模板</b-btn>
                            <b-btn variant="dark" @click="template_help = !template_help">帮助</b-btn>
                        </b-btn-group>
                        <b-btn-group class="mb-2 float-right">
                            <b-btn variant="primary" class="float-right" @click="checkTemplate(form.template)">检查</b-btn>
                        </b-btn-group>
                        <b-form-textarea id="template" rows="5" v-model="form.template" class="adanos-code-textarea  text-monospace"
                                         placeholder="输入模板"/>
                        <small class="form-text text-muted">
                            语法参考 <a href="https://golang.org/pkg/html/template/" target="_blank">https://golang.org/pkg/html/template/</a>
                        </small>
                        <TemplateHelp v-if="template_help"/>
                    </b-card>
                </b-card-group>

                <b-card-group class="mb-3">
                    <b-card header="动作">
                        <p class="text-muted">分组达到报警周期后，会按照这里的规则来将分组信息通知给对应的通道。</p>
                        <b-card :header="trigger.id" border-variant="dark" header-bg-variant="dark"
                                header-text-variant="white" class="mb-3" v-bind:key="i"
                                v-for="(trigger, i) in form.triggers">
                            <b-form-group label-cols="2" :id="'trigger_' + i" label="名称" :label-for="'trigger_name' + i">
                                <b-form-input :id="'trigger_name_' + i" v-model="trigger.name" placeholder="动作名称，可选"/>
                            </b-form-group>
                            <b-form-group label-cols="2" :id="'trigger_' + i" label="条件"
                                          :label-for="'trigger_pre_condition_' + i">
                                <b-btn-group class="mb-2">
                                    <b-btn variant="warning" @click="openTriggerRuleTemplateSelector(i)">插入模板</b-btn>
                                    <b-btn variant="dark" @click="toggleHelp(trigger)">帮助</b-btn>
                                </b-btn-group>
                                <b-btn-group class="mb-2 float-right">
                                    <b-btn variant="primary" class="float-right" @click="checkTriggerRule(trigger)">检查
                                    </b-btn>
                                </b-btn-group>
                                <b-form-textarea id="'trigger_pre_condition_' + i" v-model="trigger.pre_condition"
                                                 class="adanos-code-textarea  text-monospace" placeholder="默认为 true （全部匹配）"/>
                                <small class="form-text text-muted">
                                    语法参考 <a
                                        href="https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md"
                                        target="_blank">https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md</a>
                                </small>
                                <TriggerHelp class="mt-2" v-if="trigger.help"/>
                            </b-form-group>
                            <b-form-group label-cols="2" :id="'trigger_action_' + i" label="动作"
                                          :label-for="'trigger_action_' + i">
                                <b-form-select :id="'trigger_action_' + i" v-model="trigger.action"
                                               :options="action_options"/>
                            </b-form-group>
                            <div v-if="trigger.action === 'dingding'" class="trigger_dynamic_area">
                                <b-form-group label-cols="2" :id="'trigger_meta_token_' + i" label="Token"
                                              :label-for="'trigger_meta_token_' + i">
                                    <b-form-input :id="'trigger_meta_token_' + i"
                                                  v-model="trigger.meta_arr.token" placeholder="钉钉机器人 Token"/>
                                </b-form-group>
                                <b-form-group label-cols="2" :id="'trigger_meta_secret_' + i" label="Secret"
                                              :label-for="'trigger_meta_secret_' + i">
                                    <b-form-input :id="'trigger_meta_secret_' + i"
                                                  v-model="trigger.meta_arr.secret" placeholder="钉钉机器人密钥，用于消息签名"/>
                                </b-form-group>
                                <b-form-group label-cols="2" :id="'trigger_meta_template_' + i" label="模板"
                                              :label-for="'trigger_meta_template_' + i">
                                    <b-btn-group class="mb-2">
                                        <b-btn variant="warning" @click="openDingdingTemplateSelector(i)">插入模板</b-btn>
                                        <b-btn variant="dark" @click="trigger.template_help = !trigger.template_help">帮助</b-btn>
                                    </b-btn-group>
                                    <b-btn-group class="mb-2 float-right">
                                        <b-btn variant="primary" class="float-right" @click="checkTemplate(trigger.meta_arr.template)">检查</b-btn>
                                    </b-btn-group>
                                    <b-form-textarea :id="'trigger_meta_template_' + i" rows="5" class="adanos-code-textarea  text-monospace"
                                                     v-model="trigger.meta_arr.template" placeholder="默认使用分组展示模板"/>
                                    <small class="form-text text-muted">
                                        语法参考 <a href="https://golang.org/pkg/html/template/" target="_blank">https://golang.org/pkg/html/template/</a>
                                    </small>
                                    <TemplateHelp v-if="trigger.template_help"/>
                                </b-form-group>
                            </div>
                            <div v-else-if="trigger.action === 'phone_call_aliyun'" class="trigger_dynamic_area">
                                <b-form-group label-cols="2" :id="'trigger_meta_template_id_' + i" label="语音模板ID"
                                              :label-for="'trigger_meta_template_id_' + i">
                                    <b-form-input :id="'trigger_meta_template_id_' + i"
                                                  v-model="trigger.meta_arr.template_id" placeholder="阿里云语音通知模板ID"/>
                                </b-form-group>
                                <b-form-group label-cols="2" :id="'trigger_meta_content_' + i" label="通知内容"
                                              :label-for="'trigger_meta_content_' + i">
                                    <b-form-textarea :id="'trigger_meta_content_' + i" class="adanos-code-textarea  text-monospace"
                                                     v-model="trigger.meta_arr.content" placeholder="通知内容，必须是JSON格式，包含模板变量及内容"/>
                                </b-form-group>
                            </div>
                            <div class="trigger_dynamic_area" v-else>
                                <b-form-group label-cols="2" :id="'trigger_meta_' + i" label="动作参数"
                                              :label-for="'trigger_meta_' + i" >
                                    <b-form-input :id="'trigger_meta_' + i" v-model="trigger.meta_arr.value"/>
                                </b-form-group>
                            </div>


                            <b-form-group label-cols="2" label="接收人" :label-for="'trigger_users_' + i"
                                          v-if="['dingding', 'email', 'phone_call_aliyun', 'sms_aliyun', 'sms_yunxin', 'wechat'].indexOf(trigger.action) !== -1">
                                <b-btn variant="info" class="mb-3" @click="userAdd(i)">添加接收人</b-btn>
                                <b-input-group v-bind:key="index" v-for="(user, index) in trigger.user_refs"
                                               class="mb-3">
                                    <b-form-select v-model="trigger.user_refs[index]"
                                                   :options="user_options"/>
                                    <b-input-group-append>
                                        <b-btn variant="danger" @click="userDelete(i, index)">删除</b-btn>
                                    </b-input-group-append>
                                </b-input-group>
                            </b-form-group>

                            <b-btn class="float-right" variant="danger" @click="triggerDelete(i)">删除动作</b-btn>
                        </b-card>
                        <b-btn variant="success" class="mb-3" @click="triggerAdd()">添加</b-btn>
                    </b-card>
                </b-card-group>

                <b-button type="submit" variant="primary" class="mr-2">保存</b-button>
                <b-button to="/rules">返回</b-button>
            </b-form>

            <b-modal id="match_rule_selector" title="选择分组匹配规则模板" hide-footer size="xl">
                <b-table sticky-header="500px" responsive :items="templates.match_rule" :fields="template_fields">
                    <template v-slot:cell(content)="row">
                        <code class="adanos-pre-fold">{{ row.item.content }}</code>
                    </template>
                    <template v-slot:cell(name)="row">
                        <b>{{ row.item.name }}</b>
                        <p class="adanos-description">{{ row.item.description }}</p>
                    </template>
                    <template v-slot:row-details="row">
                        <b-card>
                            <pre><code class="adanos-colorful-code">{{ row.item.content }}</code></pre>
                        </b-card>
                    </template>
                    <template v-slot:cell(operations)="row">
                        <b-button-group>
                            <b-button size="sm" variant="info" @click="applyTemplateForMatchRule(row.item.content)">选中
                            </b-button>
                            <b-button size="sm" @click="row.toggleDetails" class="mr-2">
                                {{ row.detailsShowing ? '隐藏' : '显示' }}详情
                            </b-button>
                        </b-button-group>
                    </template>
                </b-table>
            </b-modal>
            <b-modal id="template_selector" title="选择分组展示模板" hide-footer size="xl">
                <b-table sticky-header="500px" responsive :items="templates.template" :fields="template_fields">
                    <template v-slot:cell(content)="row">
                        <code class="adanos-pre-fold">{{ row.item.content }}</code>
                        
                    </template>
                    <template v-slot:cell(name)="row">
                        <b>{{ row.item.name }}</b>
                        <p class="adanos-description">{{ row.item.description }}</p>
                    </template>
                    <template v-slot:row-details="row">
                        <b-card>
                            <pre><code class="adanos-colorful-code">{{ row.item.content }}</code></pre>
                        </b-card>
                    </template>
                    <template v-slot:cell(operations)="row">
                        <b-button-group>
                            <b-button size="sm" variant="info" @click="applyTemplateForTemplate(row.item.content)">选中
                            </b-button>
                            <b-button size="sm" @click="row.toggleDetails" class="mr-2">
                                {{ row.detailsShowing ? '隐藏' : '显示' }}详情
                            </b-button>
                        </b-button-group>
                    </template>
                </b-table>
            </b-modal>
            <b-modal id="trigger_rule_selector" title="选择动作触发规则模板" hide-footer size="xl">
                <b-table sticky-header="500px" responsive :items="templates.trigger_rule" :fields="template_fields">
                    <template v-slot:cell(content)="row">
                        <code class="adanos-pre-fold">{{ row.item.content }}</code>
                    </template>
                    <template v-slot:row-details="row">
                        <b-card>
                            <pre><code class="adanos-colorful-code">{{ row.item.content }}</code></pre>
                        </b-card>
                    </template>
                    <template v-slot:cell(name)="row">
                        <b>{{ row.item.name }}</b>
                        <p class="adanos-description">{{ row.item.description }}</p>
                    </template>
                    <template v-slot:cell(operations)="row">
                        <b-button-group>
                            <b-button size="sm" variant="info" @click="applyTemplateForTriggerRule(row.item.content)">
                                选中
                            </b-button>
                            <b-button size="sm" @click="row.toggleDetails" class="mr-2">
                                {{ row.detailsShowing ? '隐藏' : '显示' }}详情
                            </b-button>
                        </b-button-group>
                    </template>
                </b-table>
            </b-modal>
            <b-modal id="template_dingding_selector" title="选择钉钉通知模板" hide-footer size="xl">
                <b-table sticky-header="500px" responsive :items="templates.template_dingding.concat(templates.template)" :fields="template_fields">
                    <template v-slot:cell(content)="row">
                        <code class="adanos-pre-fold">{{ row.item.content }}</code>
                    </template>
                    <template v-slot:row-details="row">
                        <b-card>
                            <pre><code class="adanos-colorful-code">{{ row.item.content }}</code></pre>
                        </b-card>
                    </template>
                    <template v-slot:cell(name)="row">
                        <b>{{ row.item.name }}</b>
                        <p class="adanos-description">{{ row.item.description }}</p>
                    </template>
                    <template v-slot:cell(operations)="row">
                        <b-button-group>
                            <b-button size="sm" variant="info" @click="applyTemplateForDingding(row.item.content)">
                                选中
                            </b-button>
                            <b-button size="sm" @click="row.toggleDetails" class="mr-2">
                                {{ row.detailsShowing ? '隐藏' : '显示' }}详情
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
                    {value: 'sms_aliyun', text: '阿里云短信'},
                    {value: 'sms_yunxin', text: '网易云信'},
                    {value: 'phone_call_aliyun', text: '阿里云语音通知'},
                ],
                user_options: [],
                template_fields: [
                    {key: 'name', label: '名称'},
                    {key: 'content', label: '模板内容'},
                    {key: 'operations', label: '操作', stickyColumn: true},
                ],
                templates: {
                    match_rule: [],
                    trigger_rule: [],
                    template: [],
                    template_dingding: [],
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
            checkTemplate(template) {
                if (template.trim() === '') {
                    this.ErrorBox('模板为空，无需检查');
                    return;
                }

                this.sendCheckRequest('template',  template.trim());
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
             * 打开钉钉模板选择页面
             */
            openDingdingTemplateSelector(index) {
                this.currentTriggerRuleId = index;
                this.$root.$emit('bv::show::modal', "template_dingding_selector");
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
             * 钉钉模板选择
             */
            applyTemplateForDingding(template) {
                if (this.form.triggers[this.currentTriggerRuleId].meta_arr.template.trim() === '') {
                    this.form.triggers[this.currentTriggerRuleId].meta_arr.template = template;
                } else {
                    this.form.triggers[this.currentTriggerRuleId].meta_arr.template += '\n' + template;
                }
                this.$bvModal.hide('template_dingding_selector');
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
                    name: '',
                    pre_condition: '',
                    action: 'dingding',
                    meta: '',
                    meta_arr: {template: ''},
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
                requestData.triggers = this.form.triggers.map(function (value) {
                    value.meta = JSON.stringify(value.meta_arr);
                    return value;
                });
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
                        trigger.template_help = false;
                        trigger.meta_arr = {};

                        try {
                            trigger.meta_arr = JSON.parse(trigger.meta);
                        } catch (e) {
                            // eslint-disable-next-line no-console
                            console.log(e);
                        }
                        
                        if (trigger.meta_arr.template == undefined) {
                            trigger.meta_arr.template = "";
                        }

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
    .trigger_dynamic_area {
        border: 1px dashed #ffc107;
        padding: 10px 10px 10px 30px;
        background-color: #fff7e1;
        border-radius: .25em;
        margin-bottom: 10px;
    }
    .adanos-pre-fold {
        width: 300px;
        height: 45px;
        overflow: hidden;
        display: inline-block;
        font-size: 70%;
    }
    .adanos-colorful-code {
        color: #e83e8c;
        font-size: 80%;
    }
    .adanos-description {
        font-size: 90%;
        font-style: italic;
    }
    .adanos-code-textarea  text-monospace {
        font-size: 85%;
    }
</style>