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
                        <b-form-textarea id="rule" rows="5" v-model="form.rule"
                                         placeholder="输入规则，必须返回布尔值"></b-form-textarea>
                        <small class="form-text text-muted">
                            语法参考 <a href="https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md"
                                    target="_blank">https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md</a>
                        </small>
                    </b-card>
                </b-card-group>

                <b-card-group class="mb-3">
                    <b-card header="模板">
                        <b-form-textarea id="template" rows="5" v-model="form.template"
                                         placeholder="输入模板"></b-form-textarea>
                        <small class="form-text text-muted">
                            ...
                        </small>
                    </b-card>
                </b-card-group>

                <b-card-group class="mb-3">
                    <b-card header="动作">
                        <b-btn variant="success" class="mb-3" @click="triggerAdd()">添加</b-btn>
                        <b-card class="mb-3" v-bind:key="i" v-for="(trigger, i) in form.triggers">
                            {{ trigger }}

                            <b-btn class="float-right" variant="danger" @click="triggerDelete(i)">删除</b-btn>
                        </b-card>
                    </b-card>
                </b-card-group>

                <b-button type="submit" variant="primary" class="mr-2">保存</b-button>
                <b-button to="/rules">返回</b-button>
            </b-form>
        </b-col>
    </b-row>
</template>

<script>
    import axios from 'axios';

    export default {
        name: 'RuleEdit',
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
                properties: ['phone', 'email',]
            };
        },
        methods: {
            triggerAdd() {
                this.form.triggers.push({pre_condition: '', action: '', meta: '', id: '', user_refs: []});
            },
            triggerDelete(index) {
                this.form.triggers.splice(index, 1);
            },
            onSubmit(evt) {
                evt.preventDefault();
                let url;
                if (this.$route.params.id !== undefined) {
                    url = '/api/rules/' + this.$route.params.id + '/';
                } else {
                    url = '/api/rules/';
                }

                axios.post(url, this.createRequest()).then(() => {
                    this.$bvModal.msgBoxOk('操作成功', {
                        centered: true,
                        okVariant: 'success',
                        headerClass: 'p-2 border-bottom-0',
                        footerClass: 'p-2 border-top-0',
                    }).then(() => {
                        window.location.reload(true);
                    });
                }).catch(error => {
                    this.$bvToast.toast(error.response !== undefined ? error.response.data.error : error.toString(), {
                        title: 'ERROR',
                        variant: 'danger'
                    });
                });
            },
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
                    this.form.triggers = response.data.triggers;
                    this.form.status = response.data.status === 'enabled';
                }).catch(error => {
                    this.$bvToast.toast(error.response !== undefined ? error.response.data.error : error.toString(), {
                        title: 'ERROR',
                        variant: 'danger'
                    });
                });
            }
        }
    }
</script>

<style>
    .adanos-input-box {
        max-width: 1000px;
    }
</style>