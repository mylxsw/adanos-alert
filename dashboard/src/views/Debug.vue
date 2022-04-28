<template>
  <b-row class="mb-5 adanos-input-box">
    <b-col>
      <b-form @submit="onSubmit">
        <b-card-group class="mb-3">
          <b-card header="表达式评估">
            <codemirror v-model="expression" class="mt-3 adanos-code-textarea"
                        :options="codemirrorOption"></codemirror>
            <b-button type="submit" variant="primary" class="mt-2 mr-2 float-right">评估</b-button>
          </b-card>
        </b-card-group>
        <b-card header="评估样本" class="mb-3">
          <b-alert v-model="errorMessage"></b-alert>
          <codemirror v-model="eventSampleText" class="mt-3 adanos-code-textarea"
                      :options="codemirrorOptionSample"></codemirror>
        </b-card>
        <EventCard class="mb-3" title="事件示例" :fold="true" v-if="eventSample !== null" :event="eventSample"
                   :event_index="0" :onlyShow="true"></EventCard>
      </b-form>
    </b-col>
  </b-row>
</template>

<script>
import axios from 'axios';
import {codemirror, CodeMirror} from 'vue-codemirror-lite';
import 'codemirror/addon/display/placeholder.js';
import {helpers, hintHandler} from '@/plugins/editor-helper';

require('codemirror/mode/go/go');
require('codemirror/addon/hint/show-hint.js')
require('codemirror/addon/hint/show-hint.css')

CodeMirror.registerHelper("hint", "go", hintHandler);

export default {
  name: 'Debug',
  components: {codemirror},
  data() {
    return {
      errorMessage: '',
      expression: '',
      event_id: null,
      eventSampleText: '',
      codemirrorOption: {
        extraKeys: {'Alt-/': 'autocomplete'},
        hintOptions: {adanosType: 'AllMatchRule'},
        smartIndent: true,
        completeSingle: false,
        lineNumbers: true,
        placeholder: '输入表达式',
        lineWrapping: true
      },
      codemirrorOptionSample: {
        lineNumbers: true,
        lineWrapping: true,
      }
    };
  },
  computed: {
    eventSample() {
      if (this.eventSampleText === '') {
        return {};
      }

      try {
        return JSON.parse(this.eventSampleText);
      } catch (e) {
        return {};
      }
    }
  },
  watch: {
    '$route': 'reload',
    'eventSampleText': 'eventSampleChanged',
  },
  methods: {
    eventSampleChanged() {
      for (let i in helpers.groupMatchRules) {
        if (helpers.groupMatchRules[i].text.match(/^Meta\[/)) {
          helpers.groupMatchRules.splice(i, 1);
        }
      }

      if (this.eventSample === null || this.eventSample.meta === undefined) {
        return ;
      }

      for (let k in this.eventSample.meta) {
        helpers.groupMatchRules.push({text: 'Meta["' + k + '"]', displayText: 'Meta["' + k + '"]'});
      }
    },
    onSubmit(evt) {
      evt.preventDefault();

      let apiEndpoint = '/api/evaluate/expression-sample/';
      try {
        let data = {
          expression: this.expression,
          event_id: this.event_id,
          event_sample: JSON.parse(this.eventSampleText)
        };

        axios.post(apiEndpoint, data).then(resp => {
          if (resp.data.error === undefined || resp.data.error === null || resp.data.error === "") {
            this.SuccessBox(this.$createElement('pre', {class: 'adanos-message-box-code'}, resp.data.res));
          } else {
            this.ErrorBox('表达式错误：' + resp.data.error);
          }
        }).catch(error => {
          this.ErrorBox(error);
        });
      } catch (e) {
        this.ErrorBox(e);
      }
    },
    reload() {
      if (this.$route.query.event_id !== undefined) {
        this.event_id = this.$route.query.event_id;

        // 编辑时获取一个样本来展示
        axios.get('/api/events/' + this.event_id + '/').then(resp => {
          this.eventSample = resp.data;
          this.eventSampleText = JSON.stringify(this.eventSample, null, "\t");
        }).catch((error) => {
          this.ToastError(error)
        })
      } else {
        this.eventSampleText = JSON.stringify({
          "content": "Hello, world",
          "meta": {
            "k1": "v1",
            "k2": [
              "v21",
              "v22",
              "v24"
            ]
          },
          "tags": [
            "tag1"
          ],
          "origin": "example"
        }, null, "\t");
      }
    }
  },
  mounted() {
    this.reload();
  }
}
</script>
