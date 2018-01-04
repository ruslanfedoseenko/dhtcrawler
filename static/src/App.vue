<template>
  <v-app
    dark
    id="inspire"
  >
    <v-navigation-drawer
      fixed
      clipped
      v-model="drawer"
      app
    >
      <v-list dense>
        <v-list-tile v-for="item in menuItems" :key="item.text" @click="">
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
        <v-list-tile class="mt-3" @click="">
          <v-list-tile-action>
            <v-icon color="grey darken-1">add_circle_outline</v-icon>
          </v-list-tile-action>
          <v-list-tile-title class="grey--text text--darken-1">Browse Channels</v-list-tile-title>
        </v-list-tile>
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
      <v-layout row align-center style="max-width: 550px">
        <h2 class="ml-0 pl-3">Torrent Search Engine</h2>
      </v-layout>
      <v-layout row align-center style="max-width: 650px">
        <v-text-field
          placeholder="Search..."
          single-line
          append-icon="search"
          :append-icon-cb="() => {}"
          v-model="searchText"
          class="white--text"
          hide-details
        />
      </v-layout>
    </v-toolbar>
    <v-content>
      <v-container justify-space-around fill-height>
        <v-layout justify-space-around align-content-end>
          <router-view/>
        </v-layout>
      </v-container>
    </v-content>
  </v-app>
</template>

<script>
  import TorrentList from './components/TorrentList'

  export default {
    components: {TorrentList},
    data: () => ({
      drawer: false,
      menuItems: [
        {icon: 'trending_up', text: 'Most Popular'}
      ],
      items2: []
    }),
    methods: {
      toggleDrawler() {
        this.drawer = !this.drawer
      }
    },
    computed: {
      searchText: {
        get() {
          return this.$store.state.searchTerm
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
  .input-group__details:after {
    background-color: rgba(255, 255, 255, 0.32) !important;
  }

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
    width: 55px;
  }
</style>
