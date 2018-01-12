<template>
  <v-container fluid grid-list-md>
    <v-layout row wrap>
      <v-flex xs12>
        <h1>{{torrent.Name}}
          <v-chip>
            <a v-bind:href="getUrl(torrent)">
              <v-icon>file_download</v-icon>
            </a>
          </v-chip>
        </h1>

      </v-flex>
      <v-flex xs12>
        <h2>Tags:</h2>
      </v-flex>
      <v-flex xs12>

        <v-chip v-for="tag in torrent.Tags" :key="tag.Id">
          {{tag.Tag}}
        </v-chip>
      </v-flex>
      <v-flex xs6>
        <h2>Files</h2>
      </v-flex>
      <v-flex xs6>
        <h2>Trackers</h2>
      </v-flex>
      <v-flex xs6>
        <v-data-table v-bind:headers="fileHeaders" :items="files" v-bind:pagination.sync="filePagination">
          <template slot="items" slot-scope="props">
            <td class="text-xs-left">
              <v-icon>insert_drive_file</v-icon>
              {{ props.item.Path }}</td>
            <td class="text-xs-right">{{ formatBytes(props.item.Size,2) }}</td>
          </template>
        </v-data-table>
      </v-flex>
      <v-flex xs6>
        <v-data-table v-bind:headers="trackerHeaders" :items="trackerScrapeInfos"  v-bind:pagination.sync="trackerPagination" disable-initial-sort>
          <template slot="items" slot-scope="props">
            <td class="text-xs-left">{{props.item.TrackerUrl}}</td>
            <td class="text-xs-right">{{props.item.Seeds}}</td>
            <td class="text-xs-right">{{props.item.Leaches}}</td>
            <td class="text-xs-right">{{toDateTimeStr(props.item.LastUpdate)}}</td>
          </template>
        </v-data-table>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
  import {mapActions} from 'vuex'

  export default {
    name: 'torrent-details',
    data: () => ({
      infoHash: '',
      loading: false,
      fileHeaders: [
        { text: 'Path:', align: 'left', value: 'Path' },
        { text: 'Size:', align: 'left', value: 'Size' }
      ],
      filePagination: {
        sortBy: 'Path',
        rowsPerPage: 10
      },
      trackerPagination: {
        rowsPerPage: 10
      },
      trackerHeaders: [
        { text: 'Tracker URL:', align: 'left' },
        { text: 'Seeds:', align: 'right' },
        { text: 'Leechs:', align: 'right' },
        { text: 'Last Update:', align: 'right' }
      ],
      files: [],
      trackerScrapeInfos: []
    }),
    created() {
      if (this.$route.params.infohash) {
        this.infoHash = this.$route.params.infohash
      }
      if (this.$route.params.torrent) {
        this.torrent = this.$route.params.torrent
        this.trackerScrapeInfos = this.torrent.TrackersInfo || []
        this.files = this.torrent.Files || []
      } else {
        this.fetchTorrent(this.infoHash).then(resp => {
          console.log(resp)
          this.trackerScrapeInfos = this.torrent.TrackersInfo || []
          this.files = this.torrent.Files || []
        })
      }
    },
    methods: Object.assign(mapActions(['fetchTorrent']), {
      formatBytes(bytes, decimals) {
        if (bytes === 0) return '0 Byte'
        const k = 1024 // or 1024 for binary
        let dm = decimals || 3
        const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB']
        let i = Math.floor(Math.log(bytes) / Math.log(k))
        return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i]
      },
      toDateTimeStr(dataData) {
        if (dataData.Valid) {
          var options = { year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit', second: '2-digit' }
          return new Date(dataData.Time).toLocaleDateString('ru-RU', options)
        }
        return '-'
      },
      getUrl(torrent) {
        return 'magnet:?xt=urn:btih:' + torrent.Infohash + '&tr=udp://tracker.coppersurfer.tk:6969/announce&tr=udp://open.demonii.com:1337/announce&tr=udp://tracker.openbittorrent.com:80&tr=http://tracker.opentrackr.org:1337/announce&tr=http://explodie.org:6969/announce'
      }
    }),
    computed: {
      torrent: {
        get() {
          return this.$store.state.torrentDetails
        },
        set(value) {
          this.$store.commit('ChangeTorrentDetails', value)
        }
      }
    }
  }
</script>

<style scoped>

</style>
