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

                        <b-form-group label-cols="2" id="rule_tags" label="标签" label-for="tags_input">
                            <b-form-tags id="tags_input" placeholder="输入规则分类标签" tag-variant="primary" tag-pills
                                         separator=" "
                                         v-model="form.tags"></b-form-tags>
                        </b-form-group>

                        <b-form-group label-cols="2" label="频率*">
                            <div class="adanos-sub-form">
                                <b-form-group label-cols="2" label="类型">
                                    <b-form-select v-model="form.ready_type">
                                        <b-form-select-option value="interval">时间间隔</b-form-select-option>
                                        <b-form-select-option value="daily_time">固定时间</b-form-select-option>
                                    </b-form-select>
                                </b-form-group>
                                <b-form-group label-cols="2" id="rule_interval" label="周期" label-for="rule_interval_input" v-if="form.ready_type === 'interval'"
                                              :description="'当前：' + (parseInt(form.interval) === 0 ? 1 : form.interval) + ' 分钟，每隔 ' + (parseInt(form.interval) === 0 ? 1 : form.interval) + ' 分钟后触发一次报警'">
                                    <b-form-input id="rule_interval_input" type="range" min="0" max="1440" step="5"
                                                  v-model="form.interval" required/>
                                </b-form-group>
                                <b-form-group label-cols="2" label="时间" v-if="form.ready_type === 'daily_time'">
                                    <b-btn variant="success" class="mb-3" @click="dailyTimeAdd()">添加</b-btn>
                                    <b-input-group v-bind:key="i" v-for="(daily_time, i) in form.daily_times" style="margin-bottom: 10px;">
                                        <b-form-timepicker v-model="form.daily_times[i]" :hour12="false" :show-seconds="false"></b-form-timepicker>
                                        <b-input-group-append>
                                            <b-btn variant="danger" @click="dailyTimeDelete(i)">删除</b-btn>
                                        </b-input-group-append>
                                    </b-input-group>
                                </b-form-group>
                            </div>
                        </b-form-group>

                        <b-form-group label-cols="2" id="is_enabled" label="是否启用*" label-for="is_enabled_checkbox">
                            <b-form-checkbox id="is_enabled_checkbox" v-model="form.status">启用</b-form-checkbox>
                        </b-form-group>
                    </b-card>
                </b-card-group>

                <b-card-group class="mb-3">
                    <b-card header="规则">
                        <p class="text-muted">分组匹配规则，作用于单条 message，用于判断该 message 是否与当前规则匹配。
                            <br/>如果 message 没有匹配任何规则，将会被标记为 <code>已取消</code>。</p>
                        <b-btn-group class="mb-2">
                            <b-btn variant="warning" v-b-modal.match_rule_selector>插入模板</b-btn>
                            <b-btn variant="dark" @click="rule_help = !rule_help">帮助</b-btn>
                        </b-btn-group>
                        <b-btn-group class="mb-2 float-right">
                            <b-btn variant="primary" class="float-right" @click="checkRule(form.template)">检查</b-btn>
                        </b-btn-group>
                        <codemirror v-model="form.rule" class="mt-3 adanos-code-textarea"
                                    :options="options.group_match_rule"></codemirror>
                        <small class="form-text text-muted">
                            语法提示 <code>Alt+/</code>，语法参考 <a
                                href="https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md"
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
                            <b-btn variant="primary" class="float-right" @click="checkTemplate(form.template)">检查
                            </b-btn>
                        </b-btn-group>
                        <codemirror v-model="form.template" class="mt-3 adanos-code-textarea"
                                    :options="options.template"></codemirror>
                        <small class="form-text text-muted">
                            语法提示 <code>Alt+/</code>，语法参考 <a href="https://golang.org/pkg/text/template/"
                                                            target="_blank">https://golang.org/pkg/text/template/</a>
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
                            <b-form-group label-cols="2" :id="'trigger_' + i" label="名称"
                                          :label-for="'trigger_name' + i">
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
                            </b-form-group>
                            <b-form-group label-cols="2">
                                <codemirror v-model="trigger.pre_condition" class="mt-3 adanos-code-textarea"
                                            :options="options.trigger_rule"></codemirror>
                                <small class="form-text text-muted">
                                    语法提示 <code>Alt+/</code>，语法参考 <a
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
                            <div v-if="trigger.action === 'dingding'" class="adanos-sub-form">
                                <b-form-group label-cols="2" label="机器人*" :label-for="'trigger_meta_robot_' + i">
                                    <b-form-select :id="'trigger_meta_robot_' + i" v-model="trigger.meta_arr.robot_id"
                                                   :options="robot_options"/>
                                </b-form-group>
                                <b-form-group label-cols="2" :id="'trigger_meta_template_' + i" label="模板"
                                              :label-for="'trigger_meta_template_' + i">
                                    <b-btn-group class="mb-2">
                                        <b-btn variant="warning" @click="openDingdingTemplateSelector(i)">插入模板</b-btn>
                                        <b-btn variant="dark" @click="trigger.template_help = !trigger.template_help">
                                            帮助
                                        </b-btn>
                                    </b-btn-group>
                                    <b-btn-group class="mb-2 float-right">
                                        <b-btn variant="primary" class="float-right"
                                               @click="checkTemplate(trigger.meta_arr.template)">检查
                                        </b-btn>
                                    </b-btn-group>
                                </b-form-group>
                                <b-form-group>
                                    <codemirror v-model="trigger.meta_arr.template" class="mt-3 adanos-code-textarea"
                                                :options="options.ding_template"></codemirror>
                                    <small class="form-text text-muted">
                                        语法提示 <code>Alt+/</code>，语法参考 <a href="https://golang.org/pkg/text/template/"
                                                                        target="_blank">https://golang.org/pkg/text/template/</a>
                                    </small>
                                    <TemplateHelp v-if="trigger.template_help"/>
                                </b-form-group>
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
                            </div>
                            <div v-else-if="trigger.action === 'phone_call_aliyun'" class="adanos-sub-form">
                                <b-form-group label-cols="2" :id="'trigger_meta_content_' + i" label="通知标题"
                                              :label-for="'trigger_meta_content_' + i">
                                    <b-form-textarea :id="'trigger_meta_content_' + i"
                                                     class="adanos-code-textarea  text-monospace"
                                                     v-model="trigger.meta_arr.title"
                                                     placeholder="通知标题，默认为规则名称"/>
                                </b-form-group>
                                <b-form-group label-cols="2" :id="'trigger_meta_template_id_' + i" label="语音模板ID"
                                              :label-for="'trigger_meta_template_id_' + i">
                                    <b-form-input :id="'trigger_meta_template_id_' + i"
                                                  v-model="trigger.meta_arr.template_id" placeholder="阿里云语音通知模板ID，留空使用默认模板"/>
                                </b-form-group>
                                <b-form-group label-cols="2" label="接收人*" :label-for="'trigger_users_' + i"
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

                            </div>
                            <div class="adanos-sub-form" v-else>
                                <b-form-group label-cols="2" :id="'trigger_meta_' + i" label="动作参数"
                                              :label-for="'trigger_meta_' + i">
                                    <b-form-input :id="'trigger_meta_' + i" v-model="trigger.meta_arr.value"/>
                                </b-form-group>
                            </div>

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

    require('codemirror/mode/go/go');
    require('codemirror/mode/markdown/markdown');

    require('codemirror/addon/hint/show-hint.js')
    require('codemirror/addon/hint/show-hint.css')

    let hintHandler = function (editor) {
        let sources = [];
        switch (editor.options.hintOptions.adanosType) {
            case 'GroupMatchRule':
                sources.push(
                    {text: 'Content', displayText: 'Content |  消息内容，字符串格式'},
                    {text: 'Meta[""]', displayText: 'Meta |  字段，字典类型'},
                    {text: 'Tags[0]', displayText: 'Tags |  字段，数组类型'},
                    {text: 'Origin', displayText: 'Origin |  消息来源，字符串'},
                );
                break;
            case 'TriggerMatchRule':
                sources.push('Group', 'Group.ID', 'Group.SeqNum', 'Group.MessageCount', 'Group.Rule', 'Group.Rule.ID', 'Group.Rule.Name', 'Group.Rule.Interval', 'Group.Rule.Rule', 'Group.Rule.Template')

                sources.push({text: 'Messages()', displayText: 'Messages() []repository.Message | 获取分组中所有的 Messages'})
                sources.push({
                    text: 'MessagesMatchRegexCount(REGEX)',
                    displayText: 'MessagesMatchRegexCount(regex string) int64  | 获取匹配指定正则表达式的 message 数量'
                })
                sources.push({
                    text: 'MessagesWithMetaCount(KEY, VALUE)',
                    displayText: 'MessagesWithMetaCount(key, value string) int64  | 获取 meta 匹配指定 key=value 的 message 数量'
                })
                sources.push({
                    text: 'MessagesWithTagsCount(TAG)',
                    displayText: 'MessagesWithTagsCount(tags string) int64  | 获取拥有指定 tag 的 message 数量，多个 tag 使用英文逗号分隔'
                })
                sources.push({text: 'TriggeredTimesInPeriod(PERIOD_IN_MINUTES, TRIGGER_STATUS)', displayText: 'TriggeredTimesInPeriod(periodInMinutes int, triggerStatus string) int64 当前规则在指定时间范围内，状态为 triggerStatus 的触发次数'})
                sources.push({text: 'LastTriggeredGroup(TRIGGER_STATUS)', displayText: 'LastTriggeredGroup(triggerStatus string) repository.MessageGroup 最后一次触发该规则的状态为 triggerStatus 的分组'})

                sources.push(
                    {text: 'collecting', displayText: 'collecting  | TriggerStatus：collecting'},
                    {text: 'pending', displayText: 'pending | TriggerStatus：pending'},
                    {text: 'ok', displayText: 'ok | TriggerStatus：ok'},
                    {text: 'failed', displayText: 'failed | TriggerStatus：failed'},
                    {text: 'canceled', displayText: 'canceled | TriggerStatus：canceled'}
                );
                sources.push(
                    {text: '2006-01-02T15:04:05Z07:00', displayText: '2006-01-02T15:04:05Z07:00  | 时间格式'},
                    {text: 'RFC3339', displayText: 'RFC3339  | 时间格式'}
                );
                break;
            case 'Template':
                break;
            case 'DingTemplate':
                break;
            default:
        }

        if (editor.options.hintOptions.adanosType === 'GroupMatchRule' || editor.options.hintOptions.adanosType === 'TriggerMatchRule') {
            sources.push(
                {text: 'matches', displayText: '"foo" matches "^b.+" 正则匹配'},
                {text: 'contains', displayText: 'contains | 字符串包含'},
                {text: 'startsWith', displayText: 'startsWith | 前缀匹配'},
                {text: 'endsWith', displayText: 'endsWith | 后缀匹配'},
            );
            sources.push(
                {text: 'in', displayText: 'user.Group in ["human_resources", "marketing"] | 包含'},
                {text: 'not in', displayText: 'user.Group not in ["human_resources", "marketing"] | 不包含'},
                {text: 'or', displayText: 'or | 或者'},
                {text: 'and', displayText: 'and | 同时'},
            );
            sources.push(
                {text: 'len', displayText: 'len | length of array, map or string'},
                {text: 'all', displayText: 'all | will return true if all element satisfies the predicate'},
                {text: 'none', displayText: 'none | will return true if all element does NOT satisfies the predicate'},
                {text: 'any', displayText: 'any | will return true if any element satisfies the predicate'},
                {text: 'one', displayText: 'one | will return true if exactly ONE element satisfies the predicate'},
                {text: 'filter', displayText: 'filter | filter array by the predicate'},
                {text: 'map', displayText: 'map | map all items with the closure'},
                {text: 'count', displayText: 'count | returns number of elements what satisfies the predicate'},
            );

            sources.push({text: 'JsonGet(KEY, DEFAULT)', displayText: 'JsonGet(key string, defaultValue string) string  | 将消息体作为json解析，获取指定的key'})
            sources.push({text: 'Upper(KEY, DEFAULT)', displayText: 'Upper(val string) string  | 字符串转大写'})
            sources.push({text: 'Lower(KEY, DEFAULT)', displayText: 'Lower(val string) string  | 字符串转小写'})
            sources.push({text: 'Now()', displayText: 'Now() time.Time  | 当前时间'})
            sources.push({text: 'ParseTime(LAYOUT, VALUE)', displayText: 'ParseTime(layout string, value string) time.Time | 时间字符串转时间对象'})
            sources.push({text: 'DailyTimeBetween(START_TIME_STR, END_TIME_STR)', displayText: 'DailyTimeBetween(startTime, endTime string) bool  | 判断当前时间是否在 startTime 和 endTime 之间（每天），时间格式为 15:04'})
        } else {
            sources.push({text: '.Messages MESSAGE_COUNT', displayText: 'Messages(limit int64) []repository.Message | 从分组中获取 MESSAGE_COUNT 个 Message'})

            sources.push({text: '{{ }}', displayText: '{{ }} |  Golang 代码块'})
            sources.push({text: '{{ range $i, $msg := ARRAY }}\n {{ $i }} {{ $msg }} \n{{ end }}', displayText: '{{ range }}  | Golang 遍历对象'})
            sources.push({text: '{{ if pipeline }}\n T1 \n{{ else if pipeline }}\n T2 \n{{ else }}\n T3 \n{{ end }}', displayText: '{{ if }} |  Golang 分支条件'})
            sources.push({text: '[]()', displayText: 'Markdown 连接地址'});
            sources.push({text: 'index MAP_VAR "KEY"', displayText: 'index $msg.Meta "message.message" | 从 Map 中获取某个 Key 的值'})

            sources.push({text: 'cut_off MAX_LENGTH STR', displayText: 'cut_off(maxLen int, val string) string  |  字符串截断'})
            sources.push({text: 'implode ELEMENT_ARR ","', displayText: 'implode(elems []string, sep string) string  |  字符串数组拼接'})
            sources.push({text: 'explode STR ","', displayText: 'explode(s, sep string) []string  |  字符串分隔成数组'})
            sources.push({text: 'ident "IDENT_STR" STR', displayText: 'ident(ident string, message string) string  |  多行字符串统一缩进'})
            sources.push({text: 'json JSONSTR', displayText: 'json(content string) string  |  JSON 字符串格式化'})
            sources.push({text: 'datetime DATETIME', displayText: 'datetime(datetime time.Time) string  |  时间格式化展示为 2006-01-02 15:04:05 格式，时区选择北京/重庆'})
            sources.push({text: 'datetime_noloc DATETIME', displayText: 'datetime_noloc(datetime time.Time) string  |  时间格式化展示为 2006-01-02 15:04:05 格式，默认时区'})
            sources.push({text: 'json_get "KEY" "DEFAULT" JSONSTR', displayText: 'json_get(key string, defaultValue string, body string) string  |  将 body 解析为 json，然后获取 key 的值，失败返回 defaultValue'})
            sources.push({text: 'json_gets "KEY" "DEFAULT" JSONSTR', displayText: 'json_gets(key string, defaultValue string, body string) string  |  将 body 解析为 json，然后获取 key 的值(可以使用逗号分割多个key作为备选)，失败返回 defaultValue'})
            sources.push({text: 'json_array "KEY" JSONSTR', displayText: 'json_array(key string, body string) []string  |  将 body 解析为 json，然后获取 key 的值（数组值）'})
            sources.push({text: 'json_flatten JSONSTR MAX_LEVEL', displayText: 'json_flatten(body string, maxLevel int) []jsonutils.KvPair  |  将 body 解析为 json，然后转换为键值对返回'})
            sources.push({text: 'starts_with STR "START_STR"', displayText: 'starts_with(haystack string, needles ...string) bool  |  判断 haystack 是否以 needles 开头'})
            sources.push({text: 'ends_with STR "START_END"', displayText: 'ends_with(haystack string, needles ...string) bool  |  判断 haystack 是否以 needles 结尾'})
            sources.push({text: 'trim STR "CUTSTR"', displayText: 'trim(s string, cutset string) string  |  去掉字符串 s 两边的 cutset 字符'})
            sources.push({text: 'trim_left STR "CUTSTR"', displayText: 'trim_left(s string, cutset string) string  |  去掉字符串 s 左侧的 cutset 字符'})
            sources.push({text: 'trim_right STR "CUTSTR"', displayText: 'trim_right(s string, cutset string) string  |  去掉字符串 s 右侧的 cutset 字符'})
            sources.push({text: 'trim_space STR', displayText: 'trim_space(s string) string  |  去掉字符串 s 两边的空格'})
            sources.push({text: 'format "FORMAT" VAL', displayText: 'format(format string, a ...interface{}) string  |  格式化展示，调用 fmt.Sprintf'})
            sources.push({text: 'integer STR', displayText: 'integer(str string) int  |  字符串转整数 '})
            sources.push({text: 'mysql_slowlog STR', displayText: 'mysql_slowlog(slowlog string) map[string]string  |  解析 MySQL 慢查询日志为 map'})
            sources.push({text: 'open_falcon_im STR', displayText: 'open_falcon_im(msg string) OpenFalconIM  |  解析 OpenFalcon 消息格式'})
            sources.push({text: 'string_mask STR LEFT', displayText: 'string_mask(content string, left int) string  |  在左右两侧只保留 left 个字符，中间所有字符替换为 *'})
            sources.push({text: 'string_tags TAG_STR SEPARATOR', displayText: 'string_tags(tags string, sep string) []string  |  将字符串 tags 用 sep 作为分隔符，切割成多个 tag，空的 tag 会被排除'})
            sources.push({text: 'remove_empty_line STR', displayText: 'remove_empty_line(content string) string | 移除字符串中的空行'})
        }

        var cur = editor.getCursor();
        var token = editor.getTokenAt(cur), start, end, search;
        if (token.end > cur.ch) {
            token.end = cur.ch;
            token.string = token.string.slice(0, cur.ch - token.start);
        }

        if (token.string.match(/^[.`"\w@][\w$#]*$/g)) {
            search = token.string;
            start = token.start;
            end = token.end;
        } else {
            start = end = cur.ch;
            search = "";
        }

        let list = [];
        if (search.charAt(0) === '"' || search.charAt(0) === '.' || search.charAt(0) === "'" || search.trim() === '') {
            list = sources;
        } else {
            for (let s in sources) {
                let str = sources[s];
                if (typeof str !== "string") {
                    str = str.text;
                }
                if (str.indexOf(search) >= 0) {
                    list.push(sources[s]);
                }
            }
        }

        return {list: list, from: CodeMirror.Pos(cur.line, start), to: CodeMirror.Pos(cur.line, end)};
    };

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
                    ready_type: 'interval',
                    daily_times: ['09:00:00'],
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
                    {value: 'phone_call_aliyun', text: '阿里云语音通知'},

                    // {value: 'http', text: 'HTTP'},
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
                    ding_template: {
                        extraKeys: {'Alt-/': 'autocomplete'},
                        mode: 'markdown',
                        smartIndent: true,
                        hintOptions: {adanosType: 'DingTemplate'},
                        completeSingle: false,
                        lineNumbers: true,
                        placeholder: '默认使用分组展示模板',
                        lineWrapping: true
                    }
                }
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

                this.sendCheckRequest('template', template.trim());
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
                    meta_arr: {template: '', robot_id: null, title: '{{ .Rule.Title }}'},
                    id: '',
                    user_refs: [],
                    help: false,
                    template_help: false,
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
                requestData.rule = this.form.rule;
                requestData.tags = this.form.tags;
                requestData.template = this.form.template;
                requestData.triggers = this.form.triggers.map(function (value) {
                    value.meta = JSON.stringify(value.meta_arr);
                    return value;
                });
                requestData.status = this.form.status ? 'enabled' : 'disabled';

                requestData.ready_type = this.form.ready_type;
                if (this.form.ready_type === 'interval') {
                    if (this.form.interval <= 0) {
                        requestData.interval = 60;
                    } else {
                        requestData.interval = this.form.interval * 60;
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
                    this.form.ready_type = response.data.ready_type === '' ? 'interval': response.data.ready_type;
                    this.form.daily_times = (response.data.daily_times === null || response.data.daily_times.length === 0) ? ['09:00:00'] : response.data.daily_times;
                    this.form.rule = response.data.rule;
                    this.form.tags = response.data.tags;
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

                        if (trigger.meta_arr.template === undefined) {
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
        }
    }
</script>

<style>
    .adanos-input-box {
        max-width: 1000px;
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
</style>