<template>
    <b-row class="mb-5 adanos-input-box">
        <b-col>
            <b-form @submit="onSubmit">
                <b-card-group class="mb-3">
                    <b-card header="基本">
                        <b-form-group label-cols="2" id="username" label="姓名*" label-for="username_input">
                            <b-form-input id="username_input" type="text" v-model="form.name" required
                                          placeholder="输入用户姓名"></b-form-input>
                        </b-form-group>

                        <b-form-group label-cols="2" label="属性">
                            <b-btn variant="success" class="mb-3" @click="propertyAdd()">添加</b-btn>
                            <b-input-group v-bind:key="index" v-for="(meta, index) in form.metas" class="mb-3">
                                <b-form-input v-model="form.metas[index].key" placeholder="属性名"
                                              list="properties"></b-form-input>
                                <b-form-input v-model="form.metas[index].value" placeholder="属性值"></b-form-input>
                                <b-input-group-append>
                                    <b-btn variant="danger" @click="propertyDelete(index)">删除</b-btn>
                                </b-input-group-append>
                            </b-input-group>

                            <datalist id="properties">
                                <option :key="index" v-for="(name, index) in properties">{{ name }}</option>
                            </datalist>
                        </b-form-group>

                    </b-card>
                </b-card-group>

                <b-button type="submit" variant="primary" class="mr-2">保存</b-button>
                <b-button to="/users">返回</b-button>
            </b-form>
        </b-col>
    </b-row>
</template>

<script>
    import axios from 'axios'

    export default {
        name: 'UserEdit',
        data() {
            return {
                form: {
                    name: '',
                    metas: [],
                    status: 'enabled',
                },
                properties: ['phone', 'email',]
            };
        },
        methods: {
            propertyAdd() {
                this.form.metas.push({key: '', value: ''});
            },
            propertyDelete(index) {
                this.form.metas.splice(index, 1);
            },
            onSubmit(evt) {
                evt.preventDefault();
                let url;
                if (this.$route.params.id !== undefined) {
                    url = '/api/users/' + this.$route.params.id + '/';
                } else {
                    url = '/api/users/';
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
                requestData.metas = this.form.metas;
                requestData.status = this.form.status;

                return requestData;
            },
        },
        mounted() {
            if (this.$route.params.id !== undefined) {
                axios.get('/api/users/' + this.$route.params.id + '/').then(response => {
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