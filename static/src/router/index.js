import Vue from 'vue'
import Router from 'vue-router'
import torrentList from '@/components/TorrentList'
import homePage from '@/components/HomePage'
import torrentDetails from '@/components/TorrentDetails'
import page404 from '@/components/NotFoundPage'
import maintanance from '@/components/MaintenancePage'
Vue.use(Router)
let isMaintanance = true
let router
if (isMaintanance) {
  router = new Router({
    routes: [
      {
        path: '*',
        name: 'MaintenancePage',
        component: maintanance
      }
    ]
  })
} else {
  router = new Router({
    routes: [
      {
        path: '*',
        component: page404
      },
      {
        path: '/',
        name: 'HomePage',
        component: homePage
      },
      {
        path: '/page/:page',
        name: 'pagedTorrentList',
        component: torrentList
      },
      {
        path: '/search/:search',
        name: 'SearchTorrentList',
        component: torrentList
      },
      {
        path: '/search/:search/page/:page',
        name: 'SearchPagedTorrentList',
        component: torrentList
      },
      {
        path: '/details/:infohash',
        name: 'TorrentDetails',
        component: torrentDetails
      }
    ]
  })
}
export default router
