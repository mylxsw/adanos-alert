<template>
    <b-row class="mb-5 adanos-input-box">
        <b-col>
            <b-form @submit="onSubmit">
                <b-card-group class="mb-3">
                    <b-card>
                        <template v-slot:header>
                            基本
                            <div class="float-right">
                                <b-link class="ml-2" @click="basic_card_fold = !basic_card_fold">
                                    <b-icon icon="arrows-collapse" v-if="basic_card_fold"></b-icon>
                                    <b-icon icon="arrows-expand" v-if="!basic_card_fold"></b-icon>
                                </b-link>
                            </div>
                        </template>
                        <b-card-text v-if="basic_card_fold">...</b-card-text>
                        <b-card-text v-if="!basic_card_fold">
                            <b-form-group label-cols="2" id="rule_name" label="名称*" label-for="name_input">
                                <b-form-input id="name_input" type="text" v-model="form.name" required placeholder="输入规则名称"/>
                            </b-form-group>

                            <b-form-group label-cols="2" id="rule_description" label="描述" label-for="description_input">
                                <b-form-textarea id="description_input" placeholder="输入规则描述" v-model="form.description"/>
                            </b-form-group>

                            <b-form-group label-cols="2" id="rule_tags" label="标签" label-for="tags_input">
                                <b-form-tags id="tags_input" placeholder="输入规则分类标签" tag-variant="primary" tag-pills separator=" " v-model="form.tags"></b-form-tags>
                            </b-form-group>

                            <b-form-group label-cols="2" label="聚合周期*">
                                <div class="adanos-sub-form">
                                    <b-form-group label-cols="2" label="类型">
                                        <b-form-select v-model="form.ready_type">
                                            <b-form-select-option value="interval">时间间隔</b-form-select-option>
                                            <b-form-select-option value="daily_time">固定时间</b-form-select-option>
                                            <b-form-select-option value="time_range">时间范围</b-form-select-option>
                                        </b-form-select>
                                    </b-form-group>
                                    <b-form-group label-cols="2" id="rule_interval" label="周期" label-for="rule_interval_input" v-if="form.ready_type === 'interval'"
                                                  :description="'当前：' + (parseInt(form.interval) === 0 ? 1 : form.interval) + ' 分钟，每隔 ' + (parseInt(form.interval) === 0 ? 1 : form.interval) + ' 分钟后触发一次告警'">
                                        <b-form-input id="rule_interval_input" type="range" min="0" max="1440" step="5" v-model="form.interval" required/>
                                    </b-form-group>
                                    <b-form-group label-cols="2" label="时间" v-if="form.ready_type === 'daily_time'">
                                        <b-btn variant="success" class="mb-3" @click="dailyTimeAdd()">添加</b-btn>
                                        <b-input-group v-bind:key="i" v-for="(daily_time, i) in form.daily_times"  style="margin-bottom: 10px;">
                                            <b-form-timepicker v-model="form.daily_times[i]" :hour12="false" :show-seconds="false"></b-form-timepicker>
                                            <b-input-group-append>
                                                <b-btn variant="danger" @click="dailyTimeDelete(i)">删除</b-btn>
                                            </b-input-group-append>
                                        </b-input-group>
                                    </b-form-group>
                                    <b-form-group label-cols="2" label="时间范围" v-if="form.ready_type === 'time_range'" :description="timeRangeDesc">
                                        <b-btn variant="success" class="mb-3" @click="timeRangeAdd()">添加</b-btn>
                                        <b-input-group v-bind:key="i" v-for="(time_range, i) in form.time_ranges"  style="margin-bottom: 10px;">
                                            <b-form-timepicker v-model="form.time_ranges[i].start_time" :hour12="false" :show-seconds="false" placeholder="起始时间"></b-form-timepicker>
                                            <b-form-timepicker v-model="form.time_ranges[i].end_time" :hour12="false" :show-seconds="false" placeholder="截止时间"></b-form-timepicker>
                                            <b-form-input type="number" min="1" max="1440" step="1"  v-model="form.time_ranges[i].interval"/>
                                            <b-input-group-append>
                                                <b-btn variant="danger" @click="timeRangeDelete(i)">删除</b-btn>
                                            </b-input-group-append>
                                        </b-input-group>
                                    </b-form-group>
                                </div>
                            </b-form-group>
                            <hr style="border-top: 1px dashed #ccc;" class="mt-4" />
                            <b-button v-b-toggle.advance-basic variant="secondary" class="mt-2">高级</b-button>
                            <b-collapse id="advance-basic" class="mt-2">
                              <b-card>
                                  <b-form-group label-cols="2" id="rule_realtime" label="即时发送" label-for="realtime_input"
                                                :description="form.realtime_interval > 0 ? ('当前：' +  form.realtime_interval + ' 分钟，每隔 ' + form.realtime_interval + ' 分钟没有新告警产生时，再次触发的第一条告警事件将会被立即发送'):'不启用即时发送功能'">
                                    <b-form-input id="realtime_input" type="range" min="0" max="1440" step="10" v-model="form.realtime_interval"/>
                                  </b-form-group>

                                  <b-form-group label-cols="2" id="is_enabled" label="是否启用*" label-for="is_enabled_checkbox">
                                    <b-form-checkbox id="is_enabled_checkbox" v-model="form.status">启用</b-form-checkbox>
                                  </b-form-group>
                              </b-card>
                            </b-collapse>
                        </b-card-text>
                    </b-card>
                </b-card-group>

                <EventCard class="mb-3" title="事件示例" :fold="true" v-if="test_event !== null" :event="test_event" :event_index="0" :onlyShow="true"></EventCard>

                <b-card-group class="mb-3">
                    <b-card>
                        <template v-slot:header>
                            规则
                            <div class="float-right">
                                <b-link class="ml-2" @click="rule_card_fold = !rule_card_fold">
                                    <b-icon icon="arrows-collapse" v-if="rule_card_fold"></b-icon>
                                    <b-icon icon="arrows-expand" v-if="!rule_card_fold"></b-icon>
                                </b-link>
                            </div>
                        </template>
                        <b-card-text v-if="rule_card_fold">...</b-card-text>
                        <b-card-text v-if="!rule_card_fold">
                            <p class="text-muted">事件组匹配规则，作用于单条事件，用于判断该事件是否与当前规则匹配。
                                <br/>如果事件没有匹配任何规则，将会被标记为 <code>已取消</code>。</p>
                            <b-btn-group class="mb-2">
                                <b-btn variant="warning" v-b-modal.match_rule_selector>插入模板</b-btn>
                                <b-btn variant="dark" @click="rule_help = !rule_help">帮助</b-btn>
                            </b-btn-group>
                            <b-btn-group class="mb-2 float-right">
                                <b-btn variant="primary" class="float-right" @click="checkRule(form.template)">检查</b-btn>
                            </b-btn-group>
                            <codemirror v-model="form.rule" class="mt-3 adanos-code-textarea" :options="options.group_match_rule"></codemirror>
                            <small class="form-text text-muted">
                                语法提示 <code>Alt+/</code>，语法参考 <a
                                href="https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md"
                                target="_blank">https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md</a>
                            </small>
                            <MatchRuleHelp v-if="rule_help" :helpers="helper.groupMatchRules"/>
                            <hr style="border-top: 1px dashed #ccc;" class="mt-4" />
                            <b-form-group class="mt-4" label-cols="2" label="聚合条件（可选）" label-for="aggregate_cond_input">
                                <b-btn-group class="mb-2 float-right">
                                    <b-btn variant="primary" @click="checkAggregateRule(form.aggregate_rule)">检查</b-btn>
                                </b-btn-group>
                            </b-form-group>
                            <b-form-group class="mt-2" >
                                <codemirror v-model="form.aggregate_rule" class="adanos-code-textarea" :options="options.aggregate_rule"></codemirror>
                                <small class="form-text text-muted">聚合条件表达式语法与匹配规则一致，用于对符合匹配规则的一组事件按照某个可变值分组，类似于 SQL 中的 GroupBy。</small>
                            </b-form-group>
                            <hr style="border-top: 1px dashed #ccc;" class="mt-4" />
                            <b-button v-b-toggle.advance variant="secondary" class="mt-2">高级</b-button>
                            <b-collapse id="advance" visible class="mt-2">
                                <b-card>
                                    <b-form-group label-cols="2" label="忽略规则（可选）">
                                        <b-btn-group class="mb-2">
                                            <b-btn variant="dark" @click="ignore_rule_help = !ignore_rule_help">帮助</b-btn>
                                        </b-btn-group>
                                        <b-btn-group class="mb-2 float-right">
                                            <b-btn variant="primary" class="float-right" @click="checkIgnoreRule()">检查</b-btn>
                                        </b-btn-group>
                                    </b-form-group>
                                    <b-form-group label-cols="2">
                                        <codemirror v-model="form.ignore_rule" class="mt-3 adanos-code-textarea" :options="options.group_match_rule"></codemirror>
                                        <small class="form-text text-muted">当匹配规则后，会检查事件是否与忽略规则匹配，匹配则忽略该事件，常用于临时忽略某些不需要告警的事件。</small>
                                        <small class="form-text text-muted">
                                            语法提示 <code>Alt+/</code>，语法参考 <a
                                            href="https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md"
                                            target="_blank">https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md</a>
                                        </small>
                                        <MatchRuleHelp v-if="ignore_rule_help" :helpers="helper.groupMatchRules"/>
                                    </b-form-group>
                                    <b-form-group label-cols="2" label="最大忽略事件数">
                                        <b-form-input type="number" min="0" max="1000000" step="1"  v-model="form.ignore_max_count" class="w-25"/>
                                        <small class="form-text text-muted">默认为 0，所有匹配忽略规则的事件均被丢弃，当事件组中包含的被忽略事件数量大于该值时，事件组将进入正常的动作触发流程，事件不再被忽略。</small>
                                    </b-form-group>
                                    <hr style="border-top: 1px dashed #ccc;" class="mt-4" />
                                    <b-form-group label-cols="2" label="关联条件（可选）" label-for="relation_cond_input">
                                        <b-btn-group class="mb-2 float-right">
                                            <b-btn variant="primary" @click="checkAggregateRule(form.relation_rule)">检查</b-btn>
                                        </b-btn-group>
                                    </b-form-group>
                                    <b-form-group class="mt-2" label-cols="2">
                                        <codemirror v-model="form.relation_rule" class="adanos-code-textarea" :options="options.aggregate_rule"></codemirror>
                                        <small class="form-text text-muted">用于为事件建立关联，可以通过该关联关系快速查找到历史上类似的事件，这里的规则语法和聚合条件规则完全一致。</small>
                                    </b-form-group>
                                </b-card>
                            </b-collapse>
                        </b-card-text>
                    </b-card>
                </b-card-group>

                <b-card-group class="mb-3">
                    <b-card>
                        <template v-slot:header>
                            概要模板
                            <div class="float-right">
                                <b-link class="ml-2" @click="template_card_fold = !template_card_fold">
                                    <b-icon icon="arrows-collapse" v-if="template_card_fold"></b-icon>
                                    <b-icon icon="arrows-expand" v-if="!template_card_fold"></b-icon>
                                </b-link>
                            </div>
                        </template>
                        <b-card-text v-if="template_card_fold">...</b-card-text>
                        <b-card-text v-if="!template_card_fold">
                            <p class="text-muted">概要模板，用于各通知方式的默认展示模板，Adanos 会按照该模板将事件组发送给接收人。</p>
                            <b-btn-group class="mb-2">
                                <b-btn variant="warning" v-b-modal.template_selector>插入模板</b-btn>
                                <b-btn variant="dark" @click="template_help = !template_help">帮助</b-btn>
                            </b-btn-group>
                            <b-btn-group class="mb-2 float-right">
                                <b-btn variant="primary" class="float-right" @click="checkTemplate(form.template)">检查
                                </b-btn>
                            </b-btn-group>
                            <codemirror v-model="form.template" class="mt-3 adanos-code-textarea" :options="options.template"></codemirror>
                            <small class="form-text text-muted">
                                语法提示 <code>Alt+/</code>，语法参考
                                <a href="https://golang.org/pkg/text/template/" target="_blank">https://golang.org/pkg/text/template/</a>
                            </small>
                            <TemplateHelp v-if="template_help" :helpers="helper.templateRules"/>
                            <hr style="border-top: 1px dashed #ccc;" class="mt-4" />
                            <b-button v-b-toggle.advance-template variant="secondary" class="mt-2">高级</b-button>
                            <b-collapse id="advance-template" class="mt-2">
                                <b-card>
                                    <b-form-group label-cols="2" label="报告模板" label-for="report-template">
                                        <b-form-select id="report-template" v-model="form.report_template_id" :options="reportTemplateOptions"/>
                                    </b-form-group>
                                </b-card>
                            </b-collapse>
                        </b-card-text>
                    </b-card>
                </b-card-group>

                <b-card-group class="mb-3">
                    <b-card>
                        <template v-slot:header>
                            动作
                            <div class="float-right">
                                <b-link class="ml-2" @click="action_card_fold = !action_card_fold">
                                    <b-icon icon="arrows-collapse" v-if="action_card_fold"></b-icon>
                                    <b-icon icon="arrows-expand" v-if="!action_card_fold"></b-icon>
                                </b-link>
                            </div>
                        </template>
                        <b-card-text v-if="action_card_fold">...</b-card-text>
                        <b-card-text v-if="!action_card_fold">
                            <p class="text-muted">事件组达到告警周期后，会按照这里的规则来将事件组通知给对应的通道。</p>
                            <b-card :header="trigger.id" :border-variant="trigger.is_else_trigger ? 'warning':'dark'" :header-bg-variant="trigger.is_else_trigger ? 'warning':'dark'"
                                    header-text-variant="white" class="mb-3" v-bind:key="i"
                                    v-for="(trigger, i) in form.triggers">
                                <b-form-group label-cols="2" :id="'trigger_' + i" label="名称"
                                              :label-for="'trigger_name' + i">
                                    <b-form-input :id="'trigger_name_' + i" v-model="trigger.name" placeholder="动作名称，可选"/>
                                </b-form-group>
                                <b-form-group label-cols="2" :id="'trigger_' + i" label="条件" :label-for="'trigger_pre_condition_' + i">
                                    <div v-if="!trigger.is_else_trigger">
                                        <b-btn-group class="mb-2" v-if="!trigger.pre_condition_fold">
                                            <b-btn variant="warning" @click="openTriggerRuleTemplateSelector(i)">插入模板</b-btn>
                                            <b-btn variant="dark" @click="toggleHelp(trigger)">帮助</b-btn>
                                        </b-btn-group>
                                        <span class="text-muted" style="line-height: 2.5" v-if="trigger.pre_condition_fold">编辑区域已折叠，编辑请点 <b>展开</b> 按钮</span>
                                        <b-btn-group class="mb-2 float-right">
                                            <b-btn variant="primary" class="float-right" @click="checkTriggerRule(trigger)">检查</b-btn>
                                            <b-btn variant="info" class="float-right" @click="trigger.pre_condition_fold = !trigger.pre_condition_fold">{{ trigger.pre_condition_fold ? '展开' : '收起' }}</b-btn>
                                        </b-btn-group>
                                    </div>
                                    <div v-else>
                                        <span class="text-muted" style="line-height: 2.5">这是一个 <b>兜底动作</b>，当其它动作都没有被触发时生效。</span>
                                    </div>
                                </b-form-group>
                                <b-form-group label-cols="2" v-if="!trigger.pre_condition_fold">
                                    <codemirror v-model="trigger.pre_condition" class="mt-1 adanos-code-textarea"
                                                :options="options.trigger_rule"></codemirror>
                                    <small class="form-text text-muted">
                                        语法提示 <code>Alt+/</code>，语法参考 <a
                                        href="https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md"
                                        target="_blank">https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md</a>
                                    </small>
                                    <TriggerHelp class="mt-2" v-if="trigger.help" :helpers="helper.triggerMatchRules"/>
                                </b-form-group>
                                <b-form-group label-cols="2" :id="'trigger_action_' + i" label="动作"
                                              :label-for="'trigger_action_' + i">
                                    <b-form-select :id="'trigger_action_' + i" v-model="trigger.action"
                                                   :options="action_options"/>
                                </b-form-group>

                                <!-- DINGDING ACTION START -->
                                <div v-if="trigger.action === 'dingding'" class="adanos-sub-form">
                                    <b-form-group label-cols="2" label="机器人*" :label-for="'trigger_meta_robot_' + i">
                                        <b-form-select :id="'trigger_meta_robot_' + i" v-model="trigger.meta_arr.robot_id"
                                                       :options="robot_options"/>
                                    </b-form-group>
                                    <b-form-group label-cols="2" :id="'trigger_meta_template_' + i" label="模板"
                                                  :label-for="'trigger_meta_template_' + i">
                                        <b-btn-group class="mb-2" v-if="!trigger.template_fold">
                                            <b-btn variant="warning" @click="openDingdingTemplateSelector(i)">插入模板</b-btn>
                                            <b-btn variant="dark" @click="trigger.template_help = !trigger.template_help">
                                                帮助
                                            </b-btn>
                                        </b-btn-group>
                                        <span class="text-muted" style="line-height: 2.5" v-if="trigger.template_fold">编辑区域已折叠，编辑请点 <b>展开</b> 按钮</span>

                                        <b-btn-group class="mb-2 float-right">
                                            <b-btn variant="primary" class="float-right" @click="checkTemplate(trigger.meta_arr.template)">检查</b-btn>
                                            <b-btn variant="info" class="float-right" @click="trigger.template_fold = !trigger.template_fold">{{ trigger.template_fold ? '展开' : '收起' }}</b-btn>
                                        </b-btn-group>
                                    </b-form-group>
                                    <b-form-group v-if="!trigger.template_fold">
                                        <codemirror v-model="trigger.meta_arr.template" class="mt-3 adanos-code-textarea" :options="options.ding_template"></codemirror>
                                        <small class="form-text text-muted">
                                            语法提示 <code>Alt+/</code>，语法参考 <a href="https://golang.org/pkg/text/template/" target="_blank">https://golang.org/pkg/text/template/</a>
                                        </small>
                                        <TemplateHelp v-if="trigger.template_help" :helpers="helper.dingdingTemplateRules"/>
                                    </b-form-group>

                                    <b-form-group label-cols="2" label="接收人" :label-for="'trigger_users_' + i">
                                        <div class="adanos-form-group-box">
                                            <b-btn variant="info" class="mb-3" @click="userAdd(i)">添加接收人</b-btn>
                                            <b-btn v-b-toggle="'trigger_meta_template_expr_p' + i" variant="secondary" class="ml-2 mb-3">高级</b-btn>
                                            <b-collapse :id="'trigger_meta_template_expr_p' + i" class="mt-2 mb-3">
                                              <b-card>
                                                <b-form-group label-cols="2" label="接收人表达式">
                                                  <b-btn-group class="mb-2 float-right">
                                                    <b-btn variant="primary" @click="checkUserEvalRule(trigger)">检查</b-btn>
                                                  </b-btn-group>
                                                </b-form-group>
                                                <b-form-group>
                                                  <codemirror v-model="trigger.user_eval_func" class="mt-1 adanos-code-textarea" :options="options.user_eval_rule"></codemirror>
                                                  <small class="form-text text-muted">
                                                    语法提示 <code>Alt+/</code>，语法参考 <a
                                                      href="https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md"
                                                      target="_blank">https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md</a>
                                                  </small>
                                                </b-form-group>
                                              </b-card>
                                            </b-collapse>

                                            <b-input-group v-bind:key="index" v-for="(user, index) in trigger.user_refs" class="mb-3">
                                                <b-form-select v-model="trigger.user_refs[index]" :options="user_options"/>
                                                <b-input-group-append>
                                                    <b-btn variant="danger" @click="userDelete(i, index)">删除</b-btn>
                                                </b-input-group-append>
                                            </b-input-group>

                                            <small class="form-text text-muted">
                                                发送钉钉消息时，将会自动 @ 这里选择的接收人，请确保接收人在接收消息的钉钉群内。
                                            </small>
                                        </div>
                                    </b-form-group>

                                </div>
                                <!-- DINGDING ACTION END -->

                                <!-- JIRA ACTION START -->
                                <div v-else-if="trigger.action === 'jira'" class="adanos-sub-form">
                                    <b-form-group label-cols="2" :id="'trigger_meta_project_' + i" label="项目*" :label-for="'trigger_meta_project_' + i">
                                        <b-form-input :id="'trigger_meta_project_' + i" v-model="trigger.meta_arr.project_key" placeholder="项目 Key" @change="loadJiraCascadeOptions(i)"/>
                                    </b-form-group>

                                    <b-overlay :show="trigger.meta_arr.project_key === ''" rounded="sm">
                                        <b-form-group label-cols="2" :id="'trigger_meta_issue_type_' + i" label="Issue 类型" :label-for="'trigger_meta_issue_type_' + i">
                                            <b-form-select :id="'trigger_meta_issue_type_' + i" v-model="trigger.meta_arr.issue_type" placeholder="Issue 类型" :options="trigger.issue_type_options"/>
                                        </b-form-group>
                                        <b-form-group label-cols="2" :id="'trigger_meta_priority_' + i" label="优先级" :label-for="'trigger_meta_priority_' + i">
                                            <b-form-select :id="'trigger_meta_priority_' + i" v-model="trigger.meta_arr.priority" placeholder="优先级" :options="trigger.priority_options"/>
                                            <small class="form-text text-muted">这里的优先级是 JIRA 支持的所有优先级，该项目并不是全部都支持，请参考 JIRA 创建 Issue 页面配置。</small>
                                        </b-form-group>
                                        <b-form-group label-cols="2" :id="'trigger_meta_assignee_' + i" label="经办人" :label-for="'trigger_meta_assignee_' + i">
                                            <b-form-select :id="'trigger_meta_assignee_' + i" v-model="trigger.meta_arr.assignee" :options="user_options"/>
                                            <small class="form-text text-muted">请确保在用户管理中为用户添加了 jira 属性，该属性值对应 JIRA 中的用户名。</small>
                                        </b-form-group>
                                        <b-form-group label-cols="2" :id="'trigger_meta_summary_' + i" label="摘要" :label-for="'trigger_meta_summary_' + i">
                                            <b-form-input :id="'trigger_meta_summary_' + i" v-model="trigger.meta_arr.summary" placeholder="摘要"/>
                                            <small class="form-text text-muted">摘要内容支持模板语法。</small>
                                        </b-form-group>
                                        <b-form-group label-cols="2" label="自定义字段" :label-for="'trigger_custom_fields_' + i">
                                            <div class="adanos-form-group-box">
                                                <b-btn variant="info" class="mb-3" @click="customFieldAdd(i)">添加字段</b-btn>
                                                <b-input-group v-bind:key="index" v-for="(cf, index) in trigger.meta_arr.custom_fields" class="mb-3">
                                                    <b-form-select v-model="trigger.meta_arr.custom_fields[index].key" placeholder="字段名" :options="jira_custom_fields"></b-form-select>
                                                    <b-form-input v-model="trigger.meta_arr.custom_fields[index].value" placeholder="字段值"></b-form-input>
                                                    <b-input-group-append>
                                                        <b-btn variant="danger" @click="customFieldDelete(i, index)">删除</b-btn>
                                                    </b-input-group-append>
                                                </b-input-group>
                                                <small class="form-text text-muted">
                                                    这里列出的是 JIRA 支持的所有自定义字段，当前选择的项目并非全部都支持，请参考 JIRA 创建 Issue 页面配置，字段值部分支持模板语法。
                                                </small>
                                            </div>
                                        </b-form-group>
                                        <b-form-group label-cols="2" :id="'trigger_meta_template_' + i" label="描述" :label-for="'trigger_meta_template_' + i">
                                            <b-btn-group class="mb-2" v-if="!trigger.template_fold">
                                                <b-btn variant="warning" @click="openDingdingTemplateSelector(i)">插入模板</b-btn>
                                                <b-btn variant="dark" @click="trigger.template_help = !trigger.template_help">帮助</b-btn>
                                            </b-btn-group>
                                            <span class="text-muted" style="line-height: 2.5" v-if="trigger.template_fold">编辑区域已折叠，编辑请点 <b>展开</b> 按钮</span>

                                            <b-btn-group class="mb-2 float-right">
                                                <b-btn variant="primary" class="float-right" @click="checkTemplate(trigger.meta_arr.template)">检查</b-btn>
                                                <b-btn variant="info" class="float-right" @click="trigger.template_fold = !trigger.template_fold">{{ trigger.template_fold ? '展开' : '收起' }}</b-btn>
                                            </b-btn-group>
                                        </b-form-group>
                                        <b-form-group v-if="!trigger.template_fold">
                                            <codemirror v-model="trigger.meta_arr.template" class="mt-3 adanos-code-textarea" :options="options.ding_template"></codemirror>
                                            <small class="form-text text-muted">
                                                语法提示 <code>Alt+/</code>，语法参考 <a href="https://golang.org/pkg/text/template/" target="_blank">https://golang.org/pkg/text/template/</a>
                                            </small>
                                            <TemplateHelp v-if="trigger.template_help" :helpers="helper.dingdingTemplateRules"/>
                                        </b-form-group>

                                        <template #overlay>
                                            <div class="text-center">
                                                <b-icon icon="stopwatch" font-scale="3" animation="cylon"></b-icon>
                                                <p id="cancel-label">请先填写有效的项目 Key...</p>
                                            </div>
                                        </template>
                                    </b-overlay>
                                </div>
                                <!-- JIRA ACTION END -->

                                <!-- HTTP ACTION START -->
                                <div v-else-if="trigger.action === 'http'" class="adanos-sub-form">
                                    <b-form-group label-cols="2" :id="'trigger_meta_method_' + i" label="请求方式*" :label-for="'trigger_meta_method_' + i">
                                        <b-form-select :id="'trigger_meta_method_' + i" v-model="trigger.meta_arr.method" placeholder="请求方式" :options="['GET', 'POST', 'PUT', 'PATCH', 'DELETE', 'HEAD', 'OPTIONS', 'PURGE']"/>
                                    </b-form-group>
                                    <b-form-group label-cols="2" :id="'trigger_meta_url_' + i" label="URL*" :label-for="'trigger_meta_url_' + i">
                                        <b-form-input :id="'trigger_meta_url_' + i" v-model="trigger.meta_arr.url" placeholder="URL"/>
                                        <small class="form-text text-muted">
                                            URL 支持模板语法。
                                        </small>
                                    </b-form-group>
                                    <b-form-group label-cols="2" label="请求头" :label-for="'trigger_headers_' + i">
                                        <div class="adanos-form-group-box">
                                            <b-btn variant="info" class="mb-3" @click="httpHeaderAdd(i)">添加字段</b-btn>
                                            <b-input-group v-bind:key="index" v-for="(cf, index) in trigger.meta_arr.headers" class="mb-3">
                                                <b-form-input v-model="trigger.meta_arr.headers[index].key" placeholder="字段名" list="http-header-options"></b-form-input>
                                                <b-form-input v-model="trigger.meta_arr.headers[index].value" placeholder="字段值" list="http-header-values"></b-form-input>
                                                <b-input-group-append>
                                                    <b-btn variant="danger" @click="httpHeaderDelete(i, index)">删除</b-btn>
                                                </b-input-group-append>
                                            </b-input-group>
                                            <small class="form-text text-muted">
                                                请求头字段值部分支持模板语法。
                                            </small>
                                        </div>
                                    </b-form-group>
                                    <b-form-group label-cols="2" :id="'trigger_meta_template_' + i" label="请求体" :label-for="'trigger_meta_template_' + i">
                                        <b-btn-group class="mb-2" v-if="!trigger.template_fold">
                                            <b-btn variant="warning" @click="openDingdingTemplateSelector(i)">插入模板</b-btn>
                                            <b-btn variant="dark" @click="trigger.template_help = !trigger.template_help">帮助</b-btn>
                                        </b-btn-group>
                                        <span class="text-muted" style="line-height: 2.5" v-if="trigger.template_fold">编辑区域已折叠，编辑请点 <b>展开</b> 按钮</span>

                                        <b-btn-group class="mb-2 float-right">
                                            <b-btn variant="primary" class="float-right" @click="checkTemplate(trigger.meta_arr.template)">检查</b-btn>
                                            <b-btn variant="info" class="float-right" @click="trigger.template_fold = !trigger.template_fold">{{ trigger.template_fold ? '展开' : '收起' }}</b-btn>
                                        </b-btn-group>
                                    </b-form-group>
                                    <b-form-group v-if="!trigger.template_fold">
                                        <codemirror v-model="trigger.meta_arr.template" class="mt-3 adanos-code-textarea" :options="options.ding_template"></codemirror>
                                        <small class="form-text text-muted">
                                            语法提示 <code>Alt+/</code>，语法参考 <a href="https://golang.org/pkg/text/template/" target="_blank">https://golang.org/pkg/text/template/</a>
                                        </small>
                                        <TemplateHelp v-if="trigger.template_help" :helpers="helper.dingdingTemplateRules"/>
                                    </b-form-group>
                                </div>
                                <!-- HTTP ACTION END -->

                                <!-- PHONE_CALL_ALIYUN ACTION START -->
                                <div v-else-if="trigger.action === 'phone_call_aliyun'" class="adanos-sub-form">
                                    <b-form-group label-cols="2" :id="'trigger_meta_content_' + i" label="通知标题" :label-for="'trigger_meta_content_' + i">
                                        <b-form-textarea :id="'trigger_meta_content_' + i" class="adanos-code-textarea  text-monospace" v-model="trigger.meta_arr.title" placeholder="通知标题，默认为规则名称"/>
                                    </b-form-group>
                                    <b-form-group label-cols="2" :id="'trigger_meta_template_id_' + i" label="语音模板ID" :label-for="'trigger_meta_template_id_' + i">
                                        <b-form-input :id="'trigger_meta_template_id_' + i" v-model="trigger.meta_arr.template_id" placeholder="阿里云语音通知模板ID，留空使用默认模板"/>
                                    </b-form-group>
                                    <b-form-group label-cols="2" label="接收人*" :label-for="'trigger_users_' + i">
                                        <div class="adanos-form-group-box">
                                            <b-btn variant="info" class="mb-3" @click="userAdd(i)">添加接收人</b-btn>
                                            <b-btn v-b-toggle="'trigger_meta_template_expr_p' + i" variant="secondary" class="ml-2 mb-3">高级</b-btn>
                                            <b-collapse :id="'trigger_meta_template_expr_p' + i" class="mt-2 mb-3">
                                              <b-card>
                                                <b-form-group label-cols="2" label="接收人表达式">
                                                  <b-btn-group class="mb-2 float-right">
                                                    <b-btn variant="primary" @click="checkUserEvalRule(trigger)">检查</b-btn>
                                                  </b-btn-group>
                                                </b-form-group>
                                                <b-form-group>
                                                  <codemirror v-model="trigger.user_eval_func" class="mt-1 adanos-code-textarea" :options="options.user_eval_rule"></codemirror>
                                                  <small class="form-text text-muted">
                                                    语法提示 <code>Alt+/</code>，语法参考 <a
                                                      href="https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md"
                                                      target="_blank">https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md</a>
                                                  </small>
                                                </b-form-group>
                                              </b-card>
                                            </b-collapse>
                                            <b-input-group v-bind:key="index" v-for="(user, index) in trigger.user_refs" class="mb-3">
                                                <b-form-select v-model="trigger.user_refs[index]" :options="user_options"/>
                                                <b-input-group-append>
                                                    <b-btn variant="danger" @click="userDelete(i, index)">删除</b-btn>
                                                </b-input-group-append>
                                            </b-input-group>
                                            <small class="form-text text-muted">
                                                这些接收人将会收到电话通知。
                                            </small>
                                        </div>
                                    </b-form-group>
                                </div>
                                <!-- PHONE_CALL_ALIYUN ACTION END -->

                                <div class="adanos-sub-form" v-else>
                                    <b-form-group label-cols="2" :id="'trigger_meta_' + i" label="动作参数" :label-for="'trigger_meta_' + i">
                                        <b-form-input :id="'trigger_meta_' + i" v-model="trigger.meta_arr.value"/>
                                    </b-form-group>
                                </div>

                                <b-btn class="float-right" variant="danger" @click="triggerDelete(i)">删除动作</b-btn>
                            </b-card>
                            <b-dropdown variant="success" text="添加" class="mb-3">
                                <b-dropdown-item @click="triggerAdd(false)">规则动作</b-dropdown-item>
                                <b-dropdown-item @click="triggerAdd(true)">兜底动作</b-dropdown-item>
                            </b-dropdown>
                        </b-card-text>
                    </b-card>
                </b-card-group>

                <b-button type="submit" variant="primary" class="mr-2">保存</b-button>
                <b-button to="/rules">返回</b-button>
            </b-form>

            <datalist id="http-header-options">
                <option>Content-Type</option>
                <option>Content-Encoding</option>
                <option>Content-Language</option>
                <option>Host</option>
                <option>Referer</option>
                <option>User-Agent</option>
                <option>Authorization</option>
                <option>Accept</option>
                <option>Accept-Charset</option>
                <option>Accept-Encoding</option>
                <option>Accept-Language</option>
                <option>Cookie</option>
                <option>X-Requested-With</option>
            </datalist>
            <datalist id="http-header-values">
                <option>application/json</option>
                <option>text/plain</option>
                <option>text/html</option>
                <option>application/x-www-form-urlencoded</option>
                <option>Basic {{ "\{\{" }} base64 "username:password" {{ "\}\}" }}</option>
            </datalist>

            <b-modal id="match_rule_selector" title="选择事件组匹配规则模板" hide-footer size="xl">
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
            <b-modal id="template_selector" title="选择事件组展示模板" hide-footer size="xl">
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
                <b-table sticky-header="500px" responsive
                         :items="templates.template_dingding.concat(templates.template)" :fields="template_fields">
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

import {codemirror, CodeMirror} from 'vue-codemirror-lite';
import 'codemirror/addon/display/placeholder.js';
import {helpers, hintHandler} from '@/plugins/editor-helper';

require('codemirror/mode/go/go');
require('codemirror/mode/markdown/markdown');
require('codemirror/addon/hint/show-hint.js')
require('codemirror/addon/hint/show-hint.css')

CodeMirror.registerHelper("hint", "go", hintHandler);
CodeMirror.registerHelper("hint", "markdown", hintHandler);

export default {
    name: 'RuleEdit',
    components: {TriggerHelp, TemplateHelp, MatchRuleHelp, codemirror},
    data() {
        return {
            form: {
                name: '',
                description: '',
                tags: [],
                aggregate_rule: '',
                relation_rule: '',
                ready_type: 'interval',
                daily_times: ['09:00:00'],
                time_ranges: [
                    {start_time: '09:00', end_time: '18:00', interval: 1},
                    {start_time: '18:00', end_time: '09:00', interval: 120},
                ],
                interval: 1,
                rule: '',
                ignore_rule: '',
                ignore_max_count: 0,
                template: '',
                report_template_id: '',
                triggers: [],
                status: true,
                realtime_interval: 0,
            },
            rule_help: false,
            ignore_rule_help: false,
            template_help: false,
            properties: ['phone', 'email',],
            action_options: [
                {value: 'dingding', text: '钉钉'},
                {value: 'phone_call_aliyun', text: '阿里云语音通知'},
                {value: 'jira', text: 'JIRA'},
                {value: 'http', text: 'HTTP'},
                // {value: 'email', text: '邮件'},
                // {value: 'wechat', text: '微信'},
                // {value: 'sms_aliyun', text: '阿里云短信'},
                // {value: 'sms_yunxin', text: '网易云信'},
            ],
            user_options: [],
            robot_options: [],
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
                template_report: [],
            },
            currentTriggerRuleId: -1,
            options: {
                group_match_rule: {
                    extraKeys: {'Alt-/': 'autocomplete'},
                    hintOptions: {adanosType: 'GroupMatchRule'},
                    smartIndent: true,
                    completeSingle: false,
                    lineNumbers: true,
                    placeholder: '输入规则，必须返回布尔值',
                    lineWrapping: true
                },
                aggregate_rule: {
                    extraKeys: {'Alt-/': 'autocomplete'},
                    hintOptions: {adanosType: 'GroupMatchRule'},
                    smartIndent: true,
                    completeSingle: false,
                    lineNumbers: true,
                    placeholder: '输入规则，必须返回字符串',
                    lineWrapping: true,
                },
                template: {
                    extraKeys: {'Alt-/': 'autocomplete'},
                    mode: 'markdown',
                    hintOptions: {adanosType: 'Template'},
                    smartIndent: true,
                    completeSingle: false,
                    lineNumbers: true,
                    placeholder: '输入模板',
                    lineWrapping: true
                },
                trigger_rule: {
                    extraKeys: {'Alt-/': 'autocomplete'},
                    smartIndent: true,
                    hintOptions: {adanosType: 'TriggerMatchRule'},
                    completeSingle: false,
                    lineNumbers: true,
                    placeholder: '默认为 true （全部匹配）',
                    lineWrapping: true
                },
                user_eval_rule: {
                  extraKeys: {'Alt-/': 'autocomplete'},
                  smartIndent: true,
                  hintOptions: {adanosType: 'TriggerMatchRule'},
                  completeSingle: false,
                  lineNumbers: true,
                  placeholder: '示例 UsersHasProperty("service-owner", Group.AggregateKey, "id")',
                  lineWrapping: true
                },
                ding_template: {
                    extraKeys: {'Alt-/': 'autocomplete'},
                    mode: 'markdown',
                    smartIndent: true,
                    hintOptions: {adanosType: 'DingTemplate'},
                    completeSingle: false,
                    lineNumbers: true,
                    placeholder: '默认使用事件组展示模板',
                    lineWrapping: true
                }
            },
            helper: {
                groupMatchRules: helpers.groupMatchRules.concat(...helpers.matchRules),
                triggerMatchRules: helpers.triggerMatchRules.concat(...helpers.matchRules),
                templateRules: helpers.templates,
                dingdingTemplateRules: helpers.templates,
            },
            test_event_id: this.$route.query.test_event_id !== undefined ? this.$route.query.test_event_id : null,
            test_event: null,

            basic_card_fold: false,
            rule_card_fold: false,
            template_card_fold: false,
            action_card_fold: false,
            // 自定义字段（用于Jira内容提示）
            jira_custom_fields: [],
            jira_custom_field_map: {},
        };
    },
    computed: {
        timeRangeDesc() {
            let results = [];
            for (let i in this.form.time_ranges) {
                let startT = this.form.time_ranges[i].start_time;
                let endT = this.form.time_ranges[i].end_time;
                let interval = this.form.time_ranges[i].interval;

                if (startT === '' || endT === '') {
                    continue;
                }

                results.push("[" + startT.substr(0, 5) + " ~ " + endT.substr(0, 5) + "): 每隔 " + interval + " 分钟");
            }

            return results.join("; ");
        },
        reportTemplateOptions() {
            let res = this.templates.template_report.map(v => {return {text: v.name, value: v.id}});
            res.unshift({text: '无', value: ''})
            return res;
        },
    },
    methods: {
        /**
         * 检查匹配规则是否合法
         */
        checkRule() {
            if (this.form.rule.trim() === '') {
                this.SuccessBox('规则为空，无需检查');
                return;
            }

            this.sendCheckRequest('match_rule', this.form.rule.trim());
        },

        checkIgnoreRule() {
            if (this.form.ignore_rule.trim() === '') {
                this.SuccessBox('规则为空，无需检查');
                return ;
            }

            this.sendCheckRequest('match_rule_ignore', this.form.ignore_rule.trim());
        },

        checkTriggerRule(trigger) {
            let rule = trigger.pre_condition.trim();
            if (rule === '') {
                this.SuccessBox('动作触发条件为空，无需检查');
                return;
            }

            this.sendCheckRequest('trigger_rule', rule)
        },

        checkUserEvalRule(trigger) {
            let rule = trigger.user_eval_func.trim();
            if (rule === '') {
              this.SuccessBox('接收人表达式为空，无需检查');
              return;
            }

            this.sendCheckRequest('user_eval_rule', rule);
        },

        /**
         * 检查模板是否合法
         */
        checkTemplate(template) {
            if (template.trim() === '') {
                this.SuccessBox('模板为空，无需检查');
                return;
            }

            this.sendCheckRequest('template', template.trim());
        },

        /**
         * 检查模板是否合法
         */
        checkAggregateRule(content) {
            if (content.trim() === '') {
                this.SuccessBox('聚合条件为空，无需检查');
                return;
            }

            this.sendCheckRequest('aggregate_rule', content.trim());
        },

        /**
         * 发送规则检查请求
         */
        sendCheckRequest(type, content) {
            let msg_id = this.test_event_id;
            axios.post('/api/rules-test/rule-check/' + type + '/', {content: content, msg_id: msg_id}).then(resp => {
                if (resp.data.error === null || resp.data.error === "") {
                    this.SuccessBox(this.$createElement('pre', {class: 'adanos-message-box-code'}, resp.data.msg));
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
         * 事件组匹配规则模板选择
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
         * 添加自定义字段
         */
        customFieldAdd(triggerIndex) {
            this.form.triggers[triggerIndex].meta_arr.custom_fields.push({key: 'customfield_', value: ''});
        },
        /**
         * 删除自定义字段
         */
        customFieldDelete(triggerIndex, index) {
            this.form.triggers[triggerIndex].meta_arr.custom_fields.splice(index, 1);
        },
        /**
         * 添加请求头
         */
        httpHeaderAdd(triggerIndex) {
            this.form.triggers[triggerIndex].meta_arr.headers.push({key: '', value: ''});
        },
        /**
         * 删除请求头
         */
        httpHeaderDelete(triggerIndex, index) {
            this.form.triggers[triggerIndex].meta_arr.headers.splice(index, 1);
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
        triggerAdd(isElseTrigger) {
            this.form.triggers.push({
                name: '',
                pre_condition: '',
                is_else_trigger: isElseTrigger,
                pre_condition_fold: isElseTrigger,
                action: 'dingding',
                meta: '',
                issue_type_options: [{text: '--- 无 ---', value: ''}],
                priority_options: [{text: '--- 无 ---', value: ''}],
                meta_arr: this.createTriggerMeta(),
                id: '',
                user_refs: [],
                user_eval_func: '',
                help: false,
                template_help: false,
                template_fold: true,
            });
        },
        /**
         * 创建 Trigger Meta
         */
        createTriggerMeta() {
            return {
                template: '{{ .RuleTemplateParsed }}',
                robot_id: null,
                title: '{{ .Rule.Name }}',
                project_key: '',
                summary: '{{ .Rule.Name }}',
                issue_type: '',
                priority: '',
                assignee: '',
                custom_fields: [],
                method: 'POST',
                url: '',
                headers: [],
            };
        },
        /**
         * 删除动作
         * @param index
         */
        triggerDelete(index) {
            this.form.triggers.splice(index, 1);
        },
        /**
         * 添加执行时间
         */
        dailyTimeAdd() {
            this.form.daily_times.push('09:00:00');
        },
        /**
         * 删除执行时间
         * @param index
         */
        dailyTimeDelete(index) {
            this.form.daily_times.splice(index, 1);
        },
        /**
         * 添加时间范围
         */
        timeRangeAdd() {
            this.form.time_ranges.push({start_time: '00:00', end_time: '00:00', interval: 1})
        },
        /**
         * 删除时间范围
         */
        timeRangeDelete(index) {
            this.form.time_ranges.splice(index, 1);
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
         * 加载 Jira 级联选项
         */
        loadJiraCascadeOptions(i) {
            let projectKey = this.form.triggers[i].meta_arr['project_key'];
            if (projectKey.trim() === '') {
                this.form.triggers[i].priority_options = [{text: '--- 无 ---', value: ''}];
                this.form.triggers[i].issue_type_options = [{text: '--- 无 ---', value: ''}];
                return ;
            }

            axios.get('/api/jira/issue/options/', {params: {project_key: projectKey}}).then(resp => {
                let issueTypes = resp.data.issue_types;
                let priorities = resp.data.priorities;

                let issueTypeOptions = [{text: '--- 无 ---', value: ''}];
                if (issueTypes != null) {
                    for (let i in issueTypes) {
                        issueTypeOptions.push({text: issueTypes[i].name, value: issueTypes[i].id})
                    }
                }

                let prioritiesOptions = [{text: '--- 无 ---', value: ''}];
                if (priorities != null) {
                    for (let i in priorities) {
                        prioritiesOptions.push({text: priorities[i].name, value: priorities[i].id})
                    }
                }

                this.form.triggers[i].priority_options = prioritiesOptions;
                this.form.triggers[i].issue_type_options = issueTypeOptions;
            }).catch((error) => {
                this.ToastError(error);
            })
        },
        /**
         * 创建请求对象
         */
        createRequest() {
            let requestData = {};
            requestData.name = this.form.name;
            requestData.description = this.form.description;
            requestData.rule = this.form.rule;
            requestData.ignore_rule = this.form.ignore_rule;
            requestData.ignore_max_count = parseInt(this.form.ignore_max_count);
            requestData.tags = this.form.tags;
            requestData.aggregate_rule = this.form.aggregate_rule;
            requestData.relation_rule = this.form.relation_rule;
            requestData.template = this.form.template;
            requestData.report_template_id = this.form.report_template_id;
            requestData.realtime_interval = parseInt(this.form.realtime_interval);
            requestData.triggers = this.form.triggers.map((trigger) => {
                switch (trigger.action) {
                    case 'jira': {
                        let constraints = [];
                        let customFields = {};
                        for (let i in trigger.meta_arr.custom_fields) {
                            if (this.jira_custom_field_map !== null) {
                                constraints.push(this.jira_custom_field_map[trigger.meta_arr.custom_fields[i].key]);
                            }
                            customFields[trigger.meta_arr.custom_fields[i].key] = trigger.meta_arr.custom_fields[i].value;
                        }

                        trigger.meta = JSON.stringify({
                            issue: {
                                project_key: trigger.meta_arr.project_key,
                                summary: trigger.meta_arr.summary,
                                issue_type: trigger.meta_arr.issue_type,
                                priority: trigger.meta_arr.priority,
                                description: trigger.meta_arr.template,
                                custom_fields: customFields,
                            },
                            constraints: constraints,
                        });
                        trigger.user_refs = trigger.meta_arr.assignee !== '' ? [trigger.meta_arr.assignee] : null;
                        break;
                    }
                    case 'http': {
                        trigger.meta  = JSON.stringify({
                            method: trigger.meta_arr.method,
                            url: trigger.meta_arr.url,
                            headers: trigger.meta_arr.headers,
                            body: trigger.meta_arr.template,
                        });
                        break;
                    }
                    default: {
                        trigger.meta = JSON.stringify(trigger.meta_arr);
                    }
                }

                return trigger;
            });
            requestData.status = this.form.status ? 'enabled' : 'disabled';

            requestData.ready_type = this.form.ready_type;
            if (this.form.ready_type === 'interval') {
                if (this.form.interval <= 0) {
                    requestData.interval = 60;
                } else {
                    requestData.interval = this.form.interval * 60;
                }
            } else if (this.form.ready_type === "time_range") {
                requestData.time_ranges = [];
                for (let i in this.form.time_ranges) {
                    requestData.time_ranges.push({
                        start_time: this.form.time_ranges[i].start_time,
                        end_time: this.form.time_ranges[i].end_time,
                        interval: this.form.time_ranges[i].interval * 60,
                    })
                }
            } else {
                requestData.daily_times = this.form.daily_times;
            }

            return requestData;
        },
    },
    mounted() {
        if (this.$route.params.id !== undefined || this.$route.query.copy_from !== undefined) {
            let ruleId = this.$route.params.id;
            if (this.$route.query.copy_from !== undefined) {
                ruleId = this.$route.query.copy_from;
            }
            axios.get('/api/rules/' + ruleId + '/').then(response => {
                this.form.name = response.data.name;
                this.form.description = response.data.description;
                this.form.interval = response.data.interval / 60;
                this.form.ready_type = response.data.ready_type === '' ? 'interval' : response.data.ready_type;
                this.form.daily_times = (response.data.daily_times === null || response.data.daily_times.length === 0) ? ['09:00:00'] : response.data.daily_times;
                this.form.rule = response.data.rule;
                this.form.ignore_rule = response.data.ignore_rule;
                this.form.ignore_max_count = response.data.ignore_max_count;
                this.form.tags = response.data.tags;
                this.form.aggregate_rule = response.data.aggregate_rule;
                this.form.relation_rule = response.data.relation_rule;
                this.form.template = response.data.template;
                this.form.report_template_id = response.data.report_template_id;
                this.form.realtime_interval = response.data.realtime_interval || 0;

                if (response.data.time_ranges === null || response.data.time_ranges.length === 0) {
                    response.data.time_ranges = this.form.time_ranges;
                } else {
                    for (let i in response.data.time_ranges) {
                        response.data.time_ranges[i].interval = response.data.time_ranges[i].interval / 60;
                    }
                }
                this.form.time_ranges = response.data.time_ranges;
                let triggerJiraLoaders = [];
                for (let i in response.data.triggers) {
                    let trigger = response.data.triggers[i];
                    trigger.help = false;
                    trigger.template_help = false;

                    // 注意：这里的所有字段都必须要有，否则更新值后 UI 不会同步更新
                    trigger.meta_arr = this.createTriggerMeta();
                    trigger.priority_options = [];
                    trigger.issue_type_options = [];

                    trigger.pre_condition_fold = !(trigger.pre_condition !== null && trigger.pre_condition !== "" && trigger.pre_condition !== 'true');

                    try {
                        let parsedMeta = JSON.parse(trigger.meta);
                        switch (trigger.action) {
                            case 'jira':{
                                trigger.meta_arr = Object.assign(trigger.meta_arr, parsedMeta['issue']);
                                trigger.meta_arr.template = trigger.meta_arr.description;
                                trigger.meta_arr.assignee = (trigger.user_refs != null && trigger.user_refs.length > 0) ? trigger.user_refs[0] : '';
                                let customFields = [];
                                for (let k in parsedMeta['issue']['custom_fields']) {
                                    customFields.push({key: k, value: parsedMeta['issue']['custom_fields'][k]})
                                }
                                trigger.meta_arr.custom_fields = customFields;

                                // 触发 jira 选项更新
                                triggerJiraLoaders.push(i);
                                break;
                            }
                            case 'http': {
                                trigger.meta_arr = Object.assign(trigger.meta_arr, parsedMeta);
                                trigger.meta_arr.template = parsedMeta.body;
                                break;
                            }
                            default: {
                                trigger.meta_arr = Object.assign(trigger.meta_arr, parsedMeta);
                            }
                        }
                    } catch (e) {
                        // eslint-disable-next-line no-console
                        console.log(e);
                    }

                    if (trigger.meta_arr.template === undefined) {
                        trigger.meta_arr.template = "";
                    }
                    trigger.template_fold = trigger.meta_arr.template === "";

                    this.form.triggers.push(trigger);

                    if (trigger.user_eval_func !== '') {
                      let that = this;
                      window.setTimeout(function () {
                          that.$root.$emit('bv::toggle::collapse', 'trigger_meta_template_expr_p' + i);
                      }, 0);
                    }
                }

                for (let i in triggerJiraLoaders) {
                    this.loadJiraCascadeOptions(triggerJiraLoaders[i]);
                }

                this.form.status = response.data.status === 'enabled';
                if (this.form.ignore_rule.trim() === '' && this.form.relation_rule === '') {
                    // 解决高级规则默认不展示时，编辑器窗口初始化不完整问题
                    this.$root.$emit('bv::toggle::collapse', 'advance')
                }
                if (this.form.realtime_interval > 0) {
                    this.$root.$emit('bv::toggle::collapse', 'advance-basic')
                }

            }).catch((error) => {
                this.ToastError(error)
            });
        } else {
            // 解决高级规则默认不展示时，编辑器窗口初始化不完整问题
            this.$root.$emit('bv::toggle::collapse', 'advance')
        }

        // 加载辅助元素
        axios.all([
            axios.get('/api/users-helper/names/'),
            axios.get('/api/templates/'),
            axios.get('/api/dingding-robots-helper/names/'),
        ]).then(axios.spread((usersResp, templateResp, robotsResp) => {
            this.user_options = usersResp.data.map((val) => {
                return {value: val.id, text: val.name}
            });

            for (let i in templateResp.data) {
                this.templates[templateResp.data[i].type].push(templateResp.data[i]);
            }

            this.robot_options = robotsResp.data.map((val) => {
                return {value: val.id, text: val.name}
            });
        })).catch((error) => {
            this.ToastError(error)
        });

        // 加载辅助元素（可延迟）
        window.setTimeout(() => {
            axios.all([
                axios.get('/api/jira/issue/custom-fields/'),
            ]).then(axios.spread((cfResp) => {
                if (cfResp.data.fields !== null) {
                    let fields = cfResp.data.fields;
                    for (let i in fields) {
                        this.jira_custom_field_map[fields[i].id] = fields[i];
                        this.jira_custom_fields.push({html: fields[i].name + '(' + fields[i].id.replace(/^customfield_+/, '') + '): ' + fields[i].type, value: fields[i].id})
                    }
                }
            })).catch((error) => this.ToastError(error));
        }, 0);

        // 事件样本加载
        if (this.test_event_id !== null && this.test_event_id !== '') {
            axios.get('/api/events/' + this.test_event_id + '/').then(resp => {
                this.test_event = resp.data;
                if (this.test_event !== null && this.test_event.id !== undefined) {
                    for (let k in resp.data.meta) {
                        helpers.groupMatchRules.push({text: 'Meta["'+ k + '"]', displayText: 'Meta["'+ k + '"]'});
                        helpers.templates.push({text: k, displayText: k});
                        helpers.templates.push({text: '{{ index $msg.Meta "' + k + '" }}', displayText: '{{ index $msg.Meta "' + k + '" }}'});
                    }
                }
            }).catch((error) => {
                this.ToastError(error);
            })
        } else if (this.$route.params.id !== undefined) {
            // 编辑时获取一个样本来展示
            axios.get('/api/rules-meta/message-sample/?id=' + this.$route.params.id).then(resp => {
                this.test_event = resp.data;
                if (this.test_event !== null && this.test_event.id !== undefined) {
                    this.test_event_id = this.test_event.id;
                    for (let k in resp.data.meta) {
                        helpers.groupMatchRules.push({text: 'Meta["'+ k + '"]', displayText: 'Meta["'+ k + '"]'});
                        helpers.templates.push({text: k, displayText: k});
                        helpers.templates.push({text: '{{ index $msg.Meta "' + k + '" }}', displayText: '{{ index $msg.Meta "' + k + '" }}'});
                    }
                }
            }).catch((error) => {this.ToastError(error)})
        }
    }
}
</script>

<style>
.adanos-input-box {
    max-width: 1000px;
    margin: auto;
}

.adanos-sub-form {
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

.adanos-code-textarea {
    font-size: 14px;
}

.CodeMirror pre.CodeMirror-placeholder {
    color: #999;
}
.CodeMirror {
    border: 1px solid #eee;
    height: auto;
}
.adanos-form-group-box {
    border: 1px dashed #ccc;
    padding: 10px;
    border-radius: 5px;
    background: #fff2c9;
}
.adanos-message-box-code {
    max-height: 400px;
    overflow-y: scroll;
    border: 1px dashed #ccc;
    border-radius: 5px;
    padding: 5px;
    font-size: 80%;
    color: #e83e8c;
    white-space: pre-wrap;
    word-break: break-word;
}
</style>