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

Raven
  .config('https://51ce0777a50a4e9683b5cc163ba1a667@sentry.io/266977')
  .addPlugin(RavenVue, Vue)
  .install()
Vue.use(vueresource)
Vue.use(vuetify)
Vue.component('btoogle-footer', commonFooter)

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  store,
  template: '<App/>',
  components: {App}
})
