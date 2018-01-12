<script>
  import {mixins, Line} from 'vue-chartjs'

  export default {
    name: 'line-chart',
    extends: Line,
    mixins: [mixins.reactiveProp],
    props: ['options'],
    mounted() {
      this.renderChart(this.chartData, Object.assign(this.options, {
        tooltipTemplate(bytes) {
          if (bytes === 0) return '0 Byte'
          const k = 1024 // or 1024 for binary
          let dm = 2
          const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB']
          let i = Math.floor(Math.log(bytes) / Math.log(k))
          return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i]
        }
      }))
    }
  }
</script>

<style scoped>

</style>
