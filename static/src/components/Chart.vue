<script>
  import {mixins, Line} from 'vue-chartjs'

  export default {
    name: 'line-chart',
    extends: Line,
    mixins: [mixins.reactiveProp],
    props: ['options'],
    mounted() {
      this.renderChart(this.chartData, this.options)
    }
  }
</script>

<style scoped>

</style>
