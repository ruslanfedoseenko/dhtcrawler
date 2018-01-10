<template>
  <v-footer fixed app dark>
    <v-container fluid grid-list-sm>
      <slot></slot>
      <v-layout pa-2 row>
          <v-flex xs3 text-lg-right>

          </v-flex>
          <v-flex xs3 text-lg-right>

          </v-flex>
          <v-flex xs3 text-lg-right>
            <router-link :to="{name:'StatisticsPage'}" style="text-decoration: none">Staticstics</router-link>
          </v-flex>
          <v-flex xs3 text-lg-right>
            © {{ new Date().getFullYear() }} btoogle.com
          </v-flex>


      </v-layout>
    </v-container>

  </v-footer>
</template>

<script>
  export default {
    name: 'btoogle-footer'
  }
</script>

<style scoped>

</style>
