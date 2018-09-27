import Vue from 'vue'
import Router from 'vue-router'
import torrentList from '@/pages/TorrentList'
import homePage from '@/pages/HomePage'
import torrentDetails from '@/pages/TorrentDetails'
import page404 from '@/pages/NotFoundPage'
import stats from '@/pages/StatisticsPage'
import maintanance from '@/pages/MaintenancePage'
Vue.use(Router)
let isMaintenance = false
let router
if (isMaintenance) {
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
      },
      {
        path: '/stats',
        name: 'Stats',
        component: stats
      },
      {
        path: '/login',
        name: 'Login',
        component: homePage,
        props: {
          showLogin: true
        }
      }
    ]
  })
}
export default router
