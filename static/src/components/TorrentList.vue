<template>
  <div>
    <v-list v-if="loading === false" two-line>
      <template v-for='item in items'>
        <v-subheader v-if="item.header" v-text="item.header"/>
        <v-divider v-else-if="item.divider" v-bind:inset="item.inset"/>
        <v-list-tile avatar v-else v-bind:key="item.title" @click="">
          <v-list-tile-content>
            <v-list-tile-title v-html="item.Name"/>
            <v-list-tile-sub-title>
              <v-chip>
                <a v-bind:href="getUrl(item)">
                  <v-icon>file_download</v-icon>
                </a>
              </v-chip>
              <v-chip>{{ formatBytes(getFileSize(item),2) }}</v-chip>
              <template v-for="tag in item.Tags">
                <v-chip>{{tag.Tag}}</v-chip>
              </template>
            </v-list-tile-sub-title>
          </v-list-tile-content>
        </v-list-tile>
      </template>
    </v-list>
    <v-footer fixed app dark>
      <v-layout row align-center style="max-width: 650px">
        <v-pagination v-bind:length.sync="pageCount" v-model="currentPage"/>
      </v-layout>
    </v-footer>
    <v-progress-circular v-if="loading" indeterminate v-bind:size="70" v-bind:width="7"
                         color="purple"/>
  </div>
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
      if (this.$route.params.page) {
        this.currentPage = this.$route.params.page
      }
      this.loading = true
      if (this.$route.params.search) {
        this.searchTorrentsPaged(this.searchText, this.currentPage).then(_ => {
          this.loading = false
        })
      } else {
        this.fetchTorrentsPaged(this.currentPage).then(_ => {
          this.loading = false
        })
      }
    },
    watch: {
      '$route.params.page': function(page) {
        this.currentPage = parseInt(page)
        this.loading = true
        this.fetchTorrentsPaged(this.currentPage).then(_ => {
          this.loading = false
        })
      },
      currentPage: function(page) {
        this.$store.commit('ChangePage', page)
        this.$router.push({name: 'pagedTorrentList', params: {page: page}})
      },
      searchText: function(newSearch) {
        if (!newSearch) {
          this.fetchTorrents()
          return
        }
        if (this.searchTimerId !== -1) {
          clearTimeout(this.searchTimerId)
        }
        this.searchTimerId = setTimeout(this.searchTorrents, 800, newSearch)
      }
    },
    methods: Object.assign(mapActions(['fetchTorrents', 'fetchTorrentsPaged', 'searchTorrents', 'searchTorrentsPaged']), {
      formatBytes: function(bytes, decimals) {
        if (bytes === 0) return '0 Byte'
        const k = 1024 // or 1024 for binary
        let dm = decimals || 3
        const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB']
        let i = Math.floor(Math.log(bytes) / Math.log(k))
        return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i]
      },
      getFileSize: function(torrent) {
        let torrentSize = 0
        for (let i = 0; i < torrent.FilesTree.length; i++) {
          torrentSize += torrent.FilesTree[i].Size
        }
        return torrentSize
      },
      getUrl: function(torrent) {
        return 'magnet:?xt=urn:btih:' + torrent.Infohash + '&tr=udp://tracker.coppersurfer.tk:6969/announce&tr=udp://open.demonii.com:1337/announce&tr=udp://tracker.openbittorrent.com:80&tr=http://tracker.opentrackr.org:1337/announce&tr=http://explodie.org:6969/announce'
      }
    }),
    data: () => ({
      loading: false,
      currentPage: 1,
      searchTimerId: -1
    }),
    computed: {
      items() {
        return this.$store.state.torrents
      },
      pageCount() {
        return this.$store.state.pageCount
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
