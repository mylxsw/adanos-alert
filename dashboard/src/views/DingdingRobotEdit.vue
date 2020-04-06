<template>
    <b-row class="mb-5 adanos-input-box">
        <b-col>
            <b-form @submit="onSubmit">
                <b-card-group class="mb-3">
                    <b-card header="基本">

                        <b-form-group label-cols="2" id="name" label="机器人（群）名*" label-for="name_input">
                            <b-form-input id="name_input" type="text" v-model="form.name" required
                                          placeholder="输入机器人名称或者群名称"></b-form-input>
                        </b-form-group>

                        <b-form-group label-cols="2" id="description" label="描述" label-for="description_input">
                            <b-form-textarea id="description_input" placeholder="输入机器人描述（可选）" v-model="form.description"/>
                        </b-form-group>

                        <b-form-group label-cols="2" id="token" label="Token" label-for="token_input">
                            <b-form-input id="token_input" type="text" v-model="form.token"
                                          placeholder="输入钉钉机器人 Token"></b-form-input>
                        </b-form-group>

                        <b-form-group label-cols="2" id="secret" label="Secret" label-for="secret_input">
                            <b-form-input id="secret_input" type="text" v-model="form.secret"
                                          placeholder="输入钉钉机器人 Secret"></b-form-input>
                        </b-form-group>

                    </b-card>
                </b-card-group>

                <b-button type="submit" variant="primary" class="mr-2">保存</b-button>
                <b-button to="/dingding-robots">返回</b-button>
            </b-form>
        </b-col>
    </b-row>
</template>

<script>
    import axios from 'axios'

    export default {
        name: 'DingdingRobotEdit',
        data() {
            return {
                form: {
                    name: '',
                    description: '',
                    token: '',
                    secret: [],
                },
            };
        },
        methods: {
            onSubmit(evt) {
                evt.preventDefault();
                let url;
                if (this.$route.params.id !== undefined) {
                    url = '/api/dingding-robots/' + this.$route.params.id + '/';
                } else {
                    url = '/api/dingding-robots/';
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
                requestData.token = this.form.token;
                requestData.secret = this.form.secret;

                return requestData;
            },
        },
        mounted() {
            if (this.$route.params.id !== undefined) {
                axios.get('/api/dingding-robots/' + this.$route.params.id + '/').then(response => {
                    this.form = response.data;
                }).catch(error => {
                    this.ToastError(error);
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