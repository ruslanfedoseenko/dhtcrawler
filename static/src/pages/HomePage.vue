<template>
  <v-container justify-space-around fill-height>
    <v-layout justify-space-around align-content-end>
      <v-container fluid grid-list-xl>
        <v-layout row justify-center>
         <img src="/static/logo.png" width="272" height="75" alt="avatar">
       </v-layout>
        <v-layout row justify-center >
          <v-flex xs6 pa-3 elevation-20>
            <btoogle-search-field :performSearch="performSearch" :searchText.sync="searchText"/>
          </v-flex>
        </v-layout>
      </v-container>
      <btoogle-footer/>
    </v-layout>
  </v-container>
</template>

<script>
  export default {
    name: 'home-page',
    props: {
      showLogin: {
        type: Boolean,
        default: false
      },
      showRegistration: {
        type: Boolean,
        default: false
      }
    },
    created() {
      this.searchText = ''
    },
    data: () => ({}),
    computed: {
      searchText: {
        get() {
          return this.$store.state.search.searchTerm
        },
        set(value) {
          this.$store.commit('ChangeSearch', value)
        }
      }
    },
    methods: {
      performSearch() {
        this.$store.commit('ChangeSearch', this.searchText)
        this.$router.push({name: 'SearchTorrentList', params: {search: this.searchText}})
      }
    }
  }
</script>

<style scoped>

</style>
