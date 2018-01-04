import Vue from 'vue'
import vuex from 'vuex'

Vue.use(vuex)

export default new vuex.Store({
  state: {
    torrents: [],
    page: 0,
    searchTerm: '',
    pageCount: 0
  },
  getters: {
    getPage: state => {
      return state.page
    },
    getPageCount: state => {
      return state.pageCount
    }
  },
  mutations: {
    ChangeSearch(state, term) {
      state.searchTerm = term
    },
    ChangePage(state, page) {
      state.page = page
    },
    ChangePageCount(state, pageCount) {
      state.pageCount = pageCount
    },
    ChangeTorrents(state, torrents) {
      state.torrents = torrents
    }
  },
  actions: {
    fetchTorrents(ctx) {
      return Vue.http.get('/torrents').then(resp => {
        ctx.commit('ChangeTorrents', resp.body.Torrents || [])
        ctx.commit('ChangePageCount', resp.body.PageCount || 0)
        ctx.commit('ChangePage', resp.body.Page || 0)
      })
    },
    fetchTorrentsPaged(ctx, page) {
      let url = '/torrents'
      if (page && page > 1) {
        url += '/page/' + page
      }
      return Vue.http.get(url).then(resp => {
        ctx.commit('ChangeTorrents', resp.body.Torrents || [])
        ctx.commit('ChangePageCount', resp.body.PageCount || 0)
        ctx.commit('ChangePage', resp.body.Page || 0)
      })
    },
    searchTorrents(ctx, searchText) {
      return Vue.http.get('/torrents/search/' + encodeURIComponent(searchText)).then(resp => {
        if (ctx.state.page !== 1) {
          ctx.commit('ChangePage', 1)
        }
        ctx.commit('ChangeTorrents', resp.body.Torrents || [])
        ctx.commit('ChangePageCount', resp.body.PageCount || 0)
      })
    },
    searchTorrentsPaged(ctx, searchText, page) {
      return Vue.http.get('/torrents/search/' + encodeURIComponent(searchText) + '/page/' + page).then(resp => {
        ctx.commit('ChangeTorrents', resp.body.Torrents || [])
        ctx.commit('ChangePageCount', resp.body.PageCount || 0)
        ctx.commit('ChangePage', resp.body.Page || 0)
      })
    }
  }
})
