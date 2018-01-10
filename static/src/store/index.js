import Vue from 'vue'
import vuex from 'vuex'

Vue.use(vuex)

export default new vuex.Store({
  state: {
    torrents: [],
    page: 0,
    searchTerm: '',
    pageCount: 0,
    torrentDetails: {}
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
    },
    ChangeTorrentDetails(state, torrent) {
      state.torrentDetails = torrent
    }
  },
  actions: {
    fetchTorrents(ctx) {
      return Vue.http.get('/api/torrents/').then(resp => {
        ctx.commit('ChangeTorrents', resp.body.Torrents || [])
        ctx.commit('ChangePageCount', resp.body.PageCount || 0)
        ctx.commit('ChangePage', resp.body.Page || 0)
      })
    },
    fetchTorrent(ctx, infoHash) {
      return Vue.http.get('/api/torrent/info/' + infoHash).then(res => {
        ctx.commit('ChangeTorrentDetails', res.body)
      })
    },
    fetchTorrentsPaged(ctx, page) {
      let url = '/api/torrents/'
      if (page && page > 1) {
        url += 'page/' + page
      }
      return Vue.http.get(url).then(resp => {
        ctx.commit('ChangeTorrents', resp.body.Torrents || [])
        ctx.commit('ChangePageCount', resp.body.PageCount || 0)
        ctx.commit('ChangePage', resp.body.Page || 0)
      })
    },
    searchTorrents(ctx, searchArgs) {
      let searchText = searchArgs.search || ''
      return Vue.http.get('/api/torrents/search/' + encodeURIComponent(searchText)).then(resp => {
        if (ctx.state.page !== 1) {
          ctx.commit('ChangePage', 1)
        }
        ctx.commit('ChangeTorrents', resp.body.Torrents || [])
        ctx.commit('ChangePageCount', resp.body.PageCount || 0)
      })
    },
    searchTorrentsPaged(ctx, searchArgs) {
      let searchText = searchArgs.search || ''
      let page = searchArgs.page || 1
      return Vue.http.get('/api/torrents/search/' + encodeURIComponent(searchText) + '/page/' + page).then(resp => {
        ctx.commit('ChangeTorrents', resp.body.Torrents || [])
        ctx.commit('ChangePageCount', resp.body.PageCount || 0)
        ctx.commit('ChangePage', resp.body.Page || 0)
      })
    }
  }
})
