import Vue from 'vue'

export default {
  state: {
    isLoggedIn: false,
    csrfToken: ''
  },
  mutations: {
    LoginSuccess(state, csrfToken) {
      state.isLoggedIn = true
      state.csrfToken = csrfToken
      Vue.http.interceptors.push(function(request, next) {
        request.headers.map['X-CSRF-Token'] = csrfToken
        next()
      })
    }
  },
  actions: {
    performLogin(ctx, credentials) {
      return Vue.http.post('/api/auth/login', credentials).then(
        resp => {
          if (resp.status === 200) {
            if (resp.headers.map['x-csrf-token']) {
              ctx.commit('LoginSuccess', resp.headers.map['x-csrf-token'])
            }
          }
        }
      )
    },
    performRegistration(ctx, regData) {
      return Vue.http.post('/api/auth/register', regData)
    }
  }
}
