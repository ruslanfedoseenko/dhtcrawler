<template>


  <v-container fluid grid-list-md>
    <v-layout row wrap v-resize="generateChartData">
      <v-flex xs6 v-for="tab in tabs" :key="tab" :ref="'wraper' + tab">
        <line-chart :chartData="tabData[tab].chartData" :options="tabData[tab].options"
                    :width="800"
                    :height="250"/>
      </v-flex>
    </v-layout>


  </v-container>

</template>

<script>
  import {mapActions} from 'vuex'
  import formatters from '../Formatters'

  export default {
    name: 'statistics-page',
    created() {
      this.fetchStats().then(res => {
        this.generateChartData()
      })
    },
    data: () => ({
      currentTab: '0',
      tabs: [
        '0', '1', '2'
      ],
      tabData: [
        {
          title: 'Torrents Count',
          chartData: null,
          options: {}
        },
        {
          title: 'Files Count',
          chartData: null,
          options: {}
        },
        {
          title: 'Files Size',
          chartData: null,
          options: {}
        }
      ]
    }),
    methods: Object.assign(mapActions(['fetchStats']), {
      setOptions(firstTorrentCount, firsFileCount, firsFileSize) {
        this.tabData[0].options = {
          responsive: true,
          maintainAspectRatio: false,
          tooltips: {
            callbacks: {
              label(tooltipItem, data) {
                return (data.datasets[tooltipItem.datasetIndex].data[tooltipItem.index] + firstTorrentCount).toLocaleString()
              }
            }
          },
          scales: {

            yAxes: [{
              ticks: {
                suggestedMin: firstTorrentCount,
                callback(value) {
                  return formatters.nFormatter(value + firstTorrentCount, 3)
                }
              }
            }]
          }
        }
        this.tabData[1].options = {
          responsive: true,
          maintainAspectRatio: false,
          tooltips: {
            callbacks: {
              label(tooltipItem, data) {
                return (data.datasets[tooltipItem.datasetIndex].data[tooltipItem.index] + firsFileCount).toLocaleString()
              }
            }
          },
          scales: {
            yAxes: [{
              ticks: {
                suggestedMin: firstTorrentCount,
                // Include a dollar sign in the ticks
                callback(value) {
                  return formatters.nFormatter(value + firsFileCount, 1)
                }
              }
            }]
          }
        }
        this.tabData[2].options = {
          responsive: true,
          maintainAspectRatio: false,
          tooltips: {
            callbacks: {
              label(tooltipItem, data) {
                return formatters.formatBytes(data.datasets[tooltipItem.datasetIndex].data[tooltipItem.index] + firsFileSize, 5)
              }
            }
          },
          scales: {
            yAxes: [{
              ticks: {
                // Include a dollar sign in the ticks
                callback(value) {
                  return formatters.formatBytes(value + firsFileSize, 3)
                }
              }
            }]
          }
        }
      },
      setChartData(index, labels, data, fillColor, label) {
        this.tabData[index].chartData = {
          labels: labels,
          datasets: [
            {
              label: label,
              backgroundColor: fillColor,
              data: data
            }
          ]
        }
      },
      generateChartData() {
        let labels = []
        let torrentCount = []
        let filesSize = []
        let filesCount = []
        var firstTorrentCount = this.stats[0].TorrentsCount
        var firsFileCount = this.stats[0].FilesCount
        var firsFileSize = this.stats[0].TotalFilesSize
        for (let i in this.stats) {
          let s = this.stats[i]
          let date = new Date(s.Date)
          let label = date.getHours() + ':' + date.getMinutes()
          labels.push(label)
          torrentCount.push(s.TorrentsCount - firstTorrentCount)
          filesSize.push(s.TotalFilesSize - firsFileSize)
          filesCount.push(s.FilesCount - firsFileCount)
        }
        this.setOptions(firstTorrentCount, firsFileCount, firsFileSize)
        this.setChartData(0, labels, torrentCount, 'rgba(151,187,205,0.2)', 'Torrents Count')
        this.setChartData(1, labels, filesCount, 'rgba(34,121,121,0.2)', 'Files Count')
        this.setChartData(2, labels, filesSize, 'rgba(85,51,121,0.2)', 'Total Files Size')
      }
    }),
    computed: {
      stats: {
        get() {
          return this.$store.state.torrentStats
        },
        set(value) {
          throw new Error('Set for torrent stats is not allowed')
        }
      }
    }
  }
</script>

<style scoped>

</style>
