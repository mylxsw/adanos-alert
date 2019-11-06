<template>
    <b-row class="mb-5 main-box">
        <b-col>
            <b-card-group class="mb-3">
                <b-card header="浏览器">
                    <b-form @submit="updateBrowserSetting">
                        <b-form-group horizontal id="server_url" label="服务地址*" label-for="server_url_input">
                            <b-form-input id="server_url_input" type="text" v-model="server_url" placeholder="http://localhost:8819"></b-form-input>
                        </b-form-group>
                        <b-form-group horizontal id="token" label="Token" label-for="token_input">
                            <b-form-input id="token_input" type="text" v-model="token"></b-form-input>
                        </b-form-group>

                        <b-button type="submit" variant="primary">保存</b-button>
                    </b-form>
                </b-card>
            </b-card-group>
        </b-col>
    </b-row>
</template>

<script>
    export default {
        name: 'Setting',
        data() {
            return {
                server_url: '',
                token: '',
            };
        },
        methods: {
            updateBrowserSetting() {
                this.$store.commit('updateServerUrl', this.server_url);
                this.$store.commit('updateToken', this.token);

                this.refreshBrowserSetting();
                this.SuccessBox('操作成功', function () {
                    window.location.reload(true);
                });
            },
            refreshBrowserSetting() {
                this.server_url = this.$store.getters.serverUrl;
                this.token = this.$store.getters.token;
            },
            refreshPage() {
                this.refreshBrowserSetting();
            }
        },
        mounted() {
            this.refreshPage();
        }
    }
</script>

<style>
    .main-box {
        max-width: 1000px;
    }
</style>