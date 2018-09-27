<template>
  <v-container justify-space-around fill-height>
    <v-layout justify-space-around align-content-end>
      <v-container fluid grid-list-xl>
        <v-layout row justify-center>
          <v-icon style="font-size: 200px">fitness_center</v-icon>
        </v-layout>
        <v-layout row justify-center>
          <h1>Btoogle.com is currenly down for maintenance</h1>
        </v-layout>
      </v-container>
      <btoogle-footer/>
    </v-layout>
  </v-container>
</template>

<script>
  export default {
    name: 'maintenance-page'
  }
</script>

<style scoped>

</style>
