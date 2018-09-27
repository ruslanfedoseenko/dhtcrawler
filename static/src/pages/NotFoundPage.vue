<template>
  <v-container justify-space-around fill-height>
    <v-layout justify-space-around align-content-end>
      <v-container fluid grid-list-xl>
        <v-layout row justify-center>
          <v-icon style="font-size: 200px">explore</v-icon>
        </v-layout>
        <v-layout row justify-center>
          <h1>Page not found</h1>
        </v-layout>
        <v-layout row justify-center>
          <p>The link you followed may be broken, or the page may have been removed.</p>
        </v-layout>

      </v-container>
      <btoogle-footer/>
    </v-layout>
  </v-container>
</template>

<script>
  export default {
    name: 'not-found-page'
  }
</script>

<style scoped>

</style>
