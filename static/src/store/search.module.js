import Vue from 'vue'

export default {
  state: {
    torrents: [],
    page: 0,
    searchTerm: '',
    pageCount: 0,
    suggestions: [],
    torrentDetails: {},
    torrentStats: {}
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
    ChangeTorrentStats(state, stats) {
      state.torrentStats = stats
    },
    ChangeSuggestions(state, suggestions) {
      state.suggestions = suggestions
    },
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
    fetchStats(ctx) {
      return Vue.http.get('/api/torrents/stats/').then(resp => {
        ctx.commit('ChangeTorrentStats', resp.body || {})
      })
    },
    fetchSuggestions(ctx, input) {
      return Vue.http.get('/api/search/suggest/' + encodeURIComponent(input)).then(resp => {
        ctx.commit('ChangeSuggestions', resp.body.data || [])
      })
    },
    fetchTorrents(ctx) {
      return Vue.http.get('/api/torrents/').then(resp => {
        ctx.commit('ChangeTorrents', resp.body.Torrents || [])
        ctx.commit('ChangePageCount', resp.body.PageCount || 0)
        ctx.commit('ChangePage', resp.body.Page || 0)
      })
    },
    fetchTorrent(ctx, infoHash) {
      return Vue.http.get('/api/torrent/info/' + encodeURIComponent(infoHash)).then(res => {
        ctx.commit('ChangeTorrentDetails', res.body)
      })
    },
    fetchTorrentsPaged(ctx, page) {
      let url = '/api/torrents/'
      if (page && page > 1) {
        url += 'page/' + encodeURIComponent(page)
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
}
