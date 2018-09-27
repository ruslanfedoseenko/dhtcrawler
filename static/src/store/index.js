import Vue from 'vue'
import vuex from 'vuex'
import search from '@/store/search.module'
import auth from '@/store/auth.module'
Vue.use(vuex)

export default new vuex.Store({
  modules: {
    search: search,
    auth: auth
  }
})
