<template>
  <v-container justify-space-around fill-height>
    <v-layout justify-space-around align-content-end>
      <v-container fluid grid-list-xl>
        <link rel="stylesheet" href="/static/textfieldfix.css">
        <v-layout row justify-center>

            <img src="/static/logo.png" width="272" height="75" alt="avatar">

        </v-layout>
        <v-layout row justify-center >
          <v-flex xs6 pa-3 elevation-20>
            <v-form @submit.prevent="performSearch" ref="form" lazy-validation>
              <v-text-field placeholder="Search..."
                            single-line
                            autofocus
                            append-icon="search"
                            v-model="searchText"
                            class="textbox--no-underline"
                            :append-icon-cb="performSearch"
                            dark
                            hide-details
                            :rules="[rules.minLength]"/>
            </v-form>
          </v-flex>
        </v-layout>
      </v-container>
    </v-layout>
  </v-container>
</template>

<script>
  export default {
    name: 'home-page',
    data: () => ({
      searchText: '',
      rules: {
        minLength: (value) => value.length >= 3
      }
    }),
    methods: {
      performSearch() {
        if (this.$refs.form.validate()) {
          this.$store.commit('ChangeSearch', this.searchText)
          this.$router.push({name: 'SearchTorrentList', params: {search: this.searchText}})
        }
      }
    }
  }
</script>

<style scoped>

</style>
