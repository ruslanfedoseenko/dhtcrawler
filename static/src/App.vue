<template>
  <v-app
    dark
    id="inspire"
  >
    <v-navigation-drawer
      fixed
      clipped
      temporary
      v-model="drawer"
      app
    >
      <v-list dense>
        <v-list-tile v-for="item in menuItems" :key="item.text" :href="item.href">
          <v-list-tile-action>
            <v-icon>{{ item.icon }}</v-icon>
          </v-list-tile-action>
          <v-list-tile-content>
            <v-list-tile-title>
              {{ item.text }}
            </v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile>
        <v-subheader class="mt-3 grey--text text--darken-1">SUBSCRIPTIONS</v-subheader>
        <v-list>
          <v-list-tile v-for="item in items2" :key="item.text" avatar @click="">
            <v-list-tile-title v-text="item.text"/>
          </v-list-tile>
        </v-list>
        <v-list-tile @click="">
          <v-list-tile-action>
            <v-icon color="grey darken-1">settings</v-icon>
          </v-list-tile-action>
          <v-list-tile-title class="grey--text text--darken-1">Manage Subscriptions</v-list-tile-title>
        </v-list-tile>
      </v-list>
    </v-navigation-drawer>
    <v-toolbar dark dense fixed clipped-left app>
      <v-toolbar-title>
        <v-toolbar-side-icon @click="toggleDrawler"/>
      </v-toolbar-title>
      <v-layout row align-center style="max-width: 150px">
        <v-flex class="text-xs-center" pa-2>
          <router-link :to="{name: 'HomePage'}">
            <img src="/static/logo.png" width="100" style="margin-top: 10px" alt="Logo">
          </router-link>
        </v-flex>
      </v-layout>
      <v-layout row align-center style="max-width: 650px">
        <btoogle-search-field v-if="enableSearch" :performSearch="performSearch" :searchText.sync="searchText"/>
      </v-layout>
      <v-spacer></v-spacer>
      <v-toolbar-items class="hidden-sm-and-down">
        <v-menu
          v-if="isLoggedIn"
          v-model="userMenu"
          :close-on-content-click="false"
          bottom
          right
        >
          <v-layout slot="activator">
            <v-btn color="success" fab dark small >
              <v-icon>account_circle</v-icon>

            </v-btn>
            <v-icon>fas fa-sort-down</v-icon>
          </v-layout>

          <v-card>
            <v-list>
              <v-list-tile avatar>
                <v-list-tile-avatar>
                  <img src="https://cdn.vuetifyjs.com/images/john.jpg" alt="John">
                </v-list-tile-avatar>

                <v-list-tile-content>
                  <v-list-tile-title>{{userInfo.username}}</v-list-tile-title>
                  <v-list-tile-sub-title>{{userInfo.mail}}</v-list-tile-sub-title>
                </v-list-tile-content>

                <v-list-tile-action>
                  <v-btn
                    :class="fav ? 'red--text' : ''"
                    icon
                    @click="fav = !fav"
                  >
                    <v-icon>favorite</v-icon>
                  </v-btn>
                </v-list-tile-action>
              </v-list-tile>
            </v-list>

            <v-divider></v-divider>



            <v-card-actions>
              <v-spacer></v-spacer>

              <v-btn flat @click="userMenu = false">Cancel</v-btn>
              <v-btn color="primary" flat @click="doLogOut">Log Out</v-btn>
            </v-card-actions>
          </v-card>
        </v-menu>
        <btoogle-register-dialog v-if="!isLoggedIn" Style="margin-right: 15px"/>
        <btoogle-login-dialog v-if="!isLoggedIn" />



      </v-toolbar-items>
    </v-toolbar>
    <v-content>

      <router-view/>
    </v-content>
  </v-app>
</template>

<script>
  import {mapActions} from 'vuex'
  export default {
    data: () => ({
      drawer: false,
      userMenu: false,
      menuItems: [
        {icon: 'trending_up', text: 'Most Popular'},
        {icon: 'insert_chart', text: 'Statistics', href: '#/stats'}
      ],
      items2: [],
      rules: {
        minLength: (value) => (value !== null && value.trim().length >= 3) || 'Enter 3 or more chars'
      },
      errors: []
    }),
    mounted() {
      this.tryLoadUserInfo()
    },
    methods: Object.assign(mapActions(['performLogOut', 'tryLoadUserInfo']), {
      toggleDrawler() {
        this.drawer = !this.drawer
      },
      doLogOut() {
        this.performLogOut()
      },
      performSearch() {
        this.$store.commit('ChangeSearch', this.searchText)
        this.$router.push({name: 'SearchTorrentList', params: {search: this.searchText}})
      }
    }),
    computed: {
      isLoggedIn() {
        return this.$store.state.auth.isLoggedIn
      },
      userInfo() {
        return this.$store.state.auth.user
      },
      enableSearch() {
        return (this.$route.name !== 'HomePage' && this.$route.name !== 'MaintenancePage')
      },
      searchText: {
        get() {
          return this.$store.state.search.searchTerm
        },
        set(value) {
          this.$store.commit('ChangeSearch', value)
        }
      }
    },
    props: {
      source: String
    }
  }
</script>

<style>
  @font-face {
    font-family: 'Material Icons';
    font-style: normal;
    font-weight: 400;
    src: url(https://fonts.gstatic.com/s/materialicons/v33/2fcrYFNaTjcS6g4U3t-Y5StnKWgpfO2iSkLzTz-AABg.ttf) format('truetype');
  }

  .material-icons {
    font-family: 'Material Icons', serif;
    font-weight: normal;
    font-style: normal;
    font-size: 24px;
    line-height: 1;
    letter-spacing: normal;
    text-transform: none;
    display: inline-block;
    white-space: nowrap;
    word-wrap: normal;
    direction: ltr;
    -webkit-font-smoothing: antialiased;
  }

  .pagination__item {
    max-width: 100px;
    min-width: 34px;
    width: auto !important;
    padding: 0 6px;
  }
</style>
