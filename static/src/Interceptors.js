export default {
  csrfTokenApplier(request, next) {
    request.headers.map['X-CSRF-Token'] = [this.state.auth.csrfToken]
    next()
  },
  saveCsrf(csrfToken) {
    document.cookie = 'csrf=' + csrfToken + '; Path=/api'
  },
  clearCsrf() {
    document.cookie = 'csrf=;expires=Thu, 01 Jan 1970 00:00:01 GMT; Path=/api'
  }
}
