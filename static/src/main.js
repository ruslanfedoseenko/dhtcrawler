// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import router from './router'
import vuetify from 'vuetify'
import vueresource from 'vue-resource'
import store from './store/index'
import 'vuetify/dist/vuetify.css'
import Raven from 'raven-js'
import RavenVue from 'raven-js/plugins/vue'
import commonFooter from './components/Footer'
import chart from './components/Chart'
import SearchField from './components/SearchField'
import 'babel-polyfill'
const reduce = Function.bind.call(Function.call, Array.prototype.reduce)
const isEnumerable = Function.bind.call(Function.call, Object.prototype.propertyIsEnumerable)
const concat = Function.bind.call(Function.call, Array.prototype.concat)
const keys = Reflect.ownKeys
if (!Object.values) {
  Object.values = function values(O) {
    return reduce(keys(O), (v, k) => concat(v, typeof k === 'string' && isEnumerable(O, k) ? [O[k]] : []), [])
  }
}
if (!Object.entries) {
  Object.entries = function entries(O) {
    return reduce(keys(O), (e, k) => concat(e, typeof k === 'string' && isEnumerable(O, k) ? [[k, O[k]]] : []), [])
  }
}
Raven
  .config('http://067a88a545394790801f0086c15c8d15@sentry.btoogle.com/3')
  .addPlugin(RavenVue, Vue)
  .install()
Vue.use(vueresource)
Vue.use(vuetify)
Vue.component('line-chart', chart)
Vue.component('btoogle-footer', commonFooter)
Vue.component('btoogle-search-field', SearchField)
/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  store,
  template: '<App/>',
  components: {App}
})
