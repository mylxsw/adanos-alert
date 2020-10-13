<template>
    <b-row class="mb-5 adanos-input-box">
        <b-col>
            <b-form @submit="onSubmit">
                <b-card-group class="mb-3">
                    <b-card header="基本">
                        <b-form-group label-cols="2" id="template_type" label="类型" label-for="template_type_input">
                            <b-form-select id="template_type_input" v-model="form.type"
                                           :options="type_options"></b-form-select>
                        </b-form-group>
                        <b-form-group label-cols="2" id="templatename" label="模板名称*" label-for="templatename_input">
                            <b-form-input id="templatename_input" type="text" v-model="form.name" required
                                          placeholder="输入模板名称"></b-form-input>
                        </b-form-group>

                        <b-form-group label-cols="2" id="description" label="描述" label-for="description_input">
                            <b-form-textarea id="description_input" placeholder="输入模板描述"
                                             v-model="form.description"></b-form-textarea>
                        </b-form-group>

                        <b-form-group label-cols="2" id="template_content" label="模板内容"
                                      label-for="template_content_input">
                            <b-card bg-variant="light">
                                <b-card-body style="max-height: 500px; overflow-y: auto;">
                                    <ul style="font-size: 90%">
                                        <li v-for="(hp, i) in this.helper[this.form.type]" v-bind:key="i">
                                            <b-badge variant="primary">{{ hp.text }}</b-badge> <span class="ml-3 font-weight-lighter">({{ hp.displayText.split("|", 2)[0] }})</span>
                                            <p class="ml-3 text-muted font-italic">{{ hp.displayText.split("|", 2)[1] }}</p>
                                        </li>
                                    </ul>
                                </b-card-body>
                            </b-card>
                            <codemirror v-model="form.content" class="mt-3 adanos-code-textarea" :options="codemirrorOption"></codemirror>
                        </b-form-group>

                    </b-card>
                </b-card-group>

                <b-button type="submit" variant="primary" class="mr-2">保存</b-button>
                <b-button to="/templates">返回</b-button>
            </b-form>
        </b-col>
    </b-row>
</template>

<script>
    import axios from 'axios'

    import {codemirror, CodeMirror} from 'vue-codemirror-lite';
    import 'codemirror/addon/display/placeholder.js';

    require('codemirror/mode/go/go');
    require('codemirror/mode/markdown/markdown');
    require('codemirror/addon/hint/show-hint.js')
    require('codemirror/addon/hint/show-hint.css')

    import {helpers, hintHandler} from '@/plugins/editor-helper';

    CodeMirror.registerHelper("hint", "go", hintHandler);
    CodeMirror.registerHelper("hint", "markdown", hintHandler);

    export default {
        name: 'TemplateEdit',
        components: {codemirror},
        data() {
            return {
                form: {
                    name: '',
                    description: '',
                    content: '',
                    type: 'template',
                },
                type_options: [
                    {value: 'match_rule', text: '事件组匹配规则'},
                    {value: 'template', text: '事件组展示模板'},
                    {value: 'trigger_rule', text: '动作触发规则'},
                    {value: 'template_dingding', text: '钉钉通知模板'},
                    {value: 'template_report', text: '报告模板'},
                ],
                helper: {
                    match_rule: helpers.groupMatchRules.concat(...helpers.matchRules),
                    trigger_rule: helpers.triggerMatchRules.concat(...helpers.matchRules),
                    template: helpers.templates,
                    template_dingding: helpers.templates,
                    template_report: helpers.templates,
                },
                options: {
                    match_rule: {
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
                    template_report: {
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
                    template_dingding: {
                        extraKeys: {'Alt-/': 'autocomplete'},
                        mode: 'markdown',
                        smartIndent: true,
                        hintOptions: {adanosType: 'DingTemplate'},
                        completeSingle: false,
                        lineNumbers: true,
                        placeholder: '默认使用事件组展示模板',
                        lineWrapping: true
                    }
                }
            };
        },
        computed: {
            codemirrorOption() {
                return this.options[this.form.type]
            }
        },
        methods: {
            onSubmit(evt) {
                evt.preventDefault();
                let url;
                if (this.$route.params.id !== undefined) {
                    url = '/api/templates/' + this.$route.params.id + '/';
                } else {
                    url = '/api/templates/';
                }

                axios.post(url, this.createRequest()).then(() => {
                    this.SuccessBox('操作成功', () => {
                        window.location.reload(true);
                    })
                }).catch(error => {
                    this.ErrorBox(error)
                });
            },
            createRequest() {
                let requestData = {};
                requestData.name = this.form.name;
                requestData.description = this.form.description;
                requestData.content = this.form.content;
                requestData.type = this.form.type;

                return requestData;
            },
        },
        mounted() {
            if (this.$route.params.id !== undefined) {
                axios.get('/api/templates/' + this.$route.params.id + '/').then(response => {
                    this.form = response.data;
                }).catch(error => {
                    this.ToastError(error);
                });
            }
        }
    }
</script>

<style scoped>
    .adanos-input-box {
        max-width: 1000px;
    }
</style>