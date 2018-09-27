import Vue from 'vue'
import Interceptors from '@/Interceptors'

export default {
  state: {
    isLoggedIn: false,
    csrfToken: '',
    user: null
  },

  mutations: {
    LoginSuccess(state, loginInfo) {
      state.isLoggedIn = true
      state.csrfToken = loginInfo.csrf
      Interceptors.saveCsrf(loginInfo.csrf)
      state.user = loginInfo.user
      Vue.http.interceptors.push(Interceptors.csrfTokenApplier.bind(this))
    },
    LogOutSuccess(state) {
      state.isLoggedIn = false
      state.csrfToken = ''
      state.user = null
      Interceptors.clearCsrf()
      let interceptorIndex = Vue.http.interceptors.indexOf(Interceptors.csrfTokenApplier)
      if (interceptorIndex > -1) {
        Vue.http.interceptors.splice(interceptorIndex, 1)
      }
    }
  },
  actions: {
    performLogin(ctx, credentials) {
      return Vue.http.post('/api/auth/login', credentials).then(
        resp => {
          if (resp.status === 200) {
            if (resp.headers.map['x-csrf-token']) {
              ctx.commit('LoginSuccess', {
                csrf: resp.headers.map['x-csrf-token'],
                user: resp.body
              })
            }
          }
        }
      )
    },
    performRegistration(ctx, regData) {
      return Vue.http.post('/api/auth/register', regData)
    },
    performLogOut(ctx) {
      Vue.http.get('/api/auth/logout').then(
        resp => {
          if (resp.status === 200) {
            ctx.commit('LogOutSuccess')
          }
        }
      )
    },
    tryLoadUserInfo(ctx) {
      Vue.http.interceptors.push(Interceptors.csrfTokenApplier.bind(this))

      Vue.http.get('/api/auth/currentUser').then(
        resp => {
          if (resp.status === 200) {
            ctx.commit('LoginSuccess', {
              user: resp.body,
              csrf: resp.headers.map['x-csrf-token']
            })
          }
        },
        _ => {
          ctx.commit('LogOutSuccess')
        }
      )
    }

  }
}
