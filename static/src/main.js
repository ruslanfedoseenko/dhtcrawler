// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import router from './router'
import vuetify from 'vuetify'
import vueresource from 'vue-resource'
import store from './store/index'
import 'vuetify/dist/vuetify.min.css'

Vue.use(vueresource)
Vue.use(vuetify)

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  store,
  template: '<App/>',
  components: { App }
})
