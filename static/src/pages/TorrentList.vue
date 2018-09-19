<template>
  <v-container fluid grid-list-sm secondary>
    <v-layout row wrap>
      <v-flex d-flex xs12 order-xs5>
        <v-layout column>
          <v-container v-if="loading" justify-space-around fill-height>
            <v-layout justify-space-around align-content-end>
              <v-container fluid grid-list-xl>
                <v-layout row justify-center>
                  <v-progress-circular indeterminate v-bind:size="70" v-bind:width="7"
                                       color="purple"/>
                </v-layout>
              </v-container>
            </v-layout>
          </v-container>
          <v-list v-if="loading === false" two-line>
            <template v-for='item in items'>

              <v-list-tile @click.prevent="navigateDetails(item)">
                <v-list-tile-content>
                  <v-list-tile-title v-html="item.Name"/>
                  <v-list-tile-sub-title>
                    <v-chip>
                      <a v-bind:href="getUrl(item)" class="torrent-list-link">
                        <v-icon>fas fa-cloud-download-alt</v-icon>
                      </a>
                    </v-chip>
                    <v-chip>{{ formatBytes(getFileSize(item),2) }}</v-chip>
                    <template v-for="tag in item.Tags">
                      <v-chip>{{tag.Tag}}</v-chip>
                    </template>
                  </v-list-tile-sub-title>
                </v-list-tile-content>
              </v-list-tile>
              <v-divider/>
            </template>
          </v-list>
          <btoogle-footer>
            <v-layout row justify-center>
              <v-flex xs6 offset-xs1>
                <v-pagination v-bind:length.sync="pageCount" :total-visible="10" v-model="currentPage" v-bind:disabled="loading === true"/>
              </v-flex>
            </v-layout>
          </btoogle-footer>

        </v-layout>
      </v-flex>
    </v-layout>
  </v-container>
</template>


<script>
  import {mapActions} from 'vuex'

  export default {
    name: 'torrent-list',
    props: {
      page: {
        twoWay: true,
        type: Number,
        default: 0
      }
    },
    beforeMount() {

    },
    created() {
      // запрашиваем данные когда реактивное представление уже создано

      this.currentPage = this.$route.params.page || 1

      if (this.$route.params.search) {
        this.searchText = this.$route.params.search
      }
      this.loadData()
    },
    watch: {
      '$route': function() {
        this.currentPage = this.$route.params.page || 1

        if (this.$route.params.search) {
          this.searchText = this.$route.params.search
        }
        this.loadData()
      },
      currentPage() {
        this.updateRoute()
        this.loadData()
      }
    },
    methods: Object.assign(mapActions(['fetchTorrents', 'fetchTorrentsPaged', 'searchTorrents', 'searchTorrentsPaged']), {
      formatBytes(bytes, decimals) {
        if (bytes === 0) return '0 Byte'
        const k = 1024 // or 1024 for binary
        let dm = decimals || 3
        const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB']
        let i = Math.floor(Math.log(bytes) / Math.log(k))
        return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i]
      },
      getFileSize(torrent) {
        let torrentSize = 0
        for (let i = 0; i < torrent.FilesTree.length; i++) {
          torrentSize += torrent.FilesTree[i].Size
        }
        return torrentSize
      },
      getUrl(torrent) {
        return 'magnet:?xt=urn:btih:' + torrent.Infohash + '&tr=udp://tracker.coppersurfer.tk:6969/announce&tr=udp://open.demonii.com:1337/announce&tr=udp://tracker.openbittorrent.com:80&tr=http://tracker.opentrackr.org:1337/announce&tr=http://explodie.org:6969/announce'
      },
      updateRoute() {
        if (this.searchText !== '') {
          this.$router.push({name: 'SearchPagedTorrentList', params: {page: this.currentPage, search: this.searchText}})
        } else {
          this.$router.push({name: 'pagedTorrentList', params: {page: this.currentPage}})
        }
      },
      loadData() {
        /*
        if (this.curentLoad && this.loading) {
          this.curentLoad.abort()
        }
        */
        if (this.loading) {
          return
        }
        this.loading = true
        if (this.searchText) {
          this.curentLoad = this.searchTorrentsPaged({search: this.searchText, page: this.currentPage}).then(_ => {
            this.loading = false
          })
        } else {
          this.curentLoad = this.fetchTorrentsPaged(this.currentPage).then(_ => {
            this.loading = false
          })
        }
      },
      navigateDetails(torrent) {
        this.$router.push({
          name: 'TorrentDetails',
          params: {
            infohash: torrent.Infohash,
            torrent: torrent
          }
        })
      }
    }),
    data: () => ({
      loading: false,
      searchTimerId: -1
    }),
    computed: {
      items() {
        return this.$store.state.torrents
      },
      pageCount() {
        return parseInt(this.$store.state.pageCount)
      },
      currentPage: {
        get() {
          return parseInt(this.$store.state.page)
        },
        set(value) {
          this.$store.commit('ChangePage', value)
        }
      },
      searchText: {
        get() {
          return this.$store.state.searchTerm
        },
        set(value) {
          this.$store.commit('ChangeSearch', value)
        }
      }
    }
  }
</script>

<style scoped>

</style>
