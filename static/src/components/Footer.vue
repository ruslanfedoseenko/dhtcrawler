<template>
  <v-footer fixed app dark>
    <v-container fluid grid-list-sm>
      <slot/>
      <v-layout pa-2 row flex align-center justify-center>
          <v-flex xs12 text-lg-center>
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
