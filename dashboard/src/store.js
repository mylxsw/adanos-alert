import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    serverUrl: localStorage.getItem('server_url') || (window.location.protocol + "//" + window.location.host),
    token: localStorage.getItem('token') || '',
  },
  mutations: {
    updateServerUrl: (state, url) => {
      state.serverUrl = url;
      localStorage.setItem('server_url', url);
    },
    updateToken: (state, token) => {
      state.token = token;
      localStorage.setItem('token', token)
    }
  },
  getters: {
    serverUrl: (state) => state.serverUrl,
    token: (state) => state.token
  },
  actions: {}
})
