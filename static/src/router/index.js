import Vue from 'vue'
import Router from 'vue-router'
import torrentList from '@/components/TorrentList'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'HelloWorld',
      component: torrentList
    },
    {
      path: '/page/:page',
      name: 'pagedTorrentList',
      component: torrentList
    },
    {
      path: 'search/:search',
      name: 'SearchTorrentList',
      component: torrentList
    },
    {
      path: 'search/:search/page/:page',
      name: 'SearchPagedTorrentList',
      component: torrentList
    }
  ]
})
