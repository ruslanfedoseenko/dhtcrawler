<template>
  <v-flex>
    <v-snackbar
      :timeout="6000"
      :top="true"
      :multi-line="true"
      v-model="snackbar"
    >
      {{ validationText }}
      <v-btn flat color="pink" @click.native="snackbar = false">Close</v-btn>
    </v-snackbar>
    <v-form @submit.prevent="performSearchInternal" ref="form" lazy-validation>
      <v-menu offset-y allow-overflow full-width content-class="menu-content-bg-fix" v-model="openMenu">
        <v-text-field
          slot="activator"
          placeholder="Search..."
          single-line
          append-icon="search"
          :append-icon-cb="performSearchInternal"
          v-model="searchText"
          :rules="[rules.minLength]"
          dark
          hide-details
        />
        <v-list style="max-height: 300px">
          <v-list-tile v-for="item in suggestions" :key="item" @click="">
            <v-list-tile-title>{{ item }}</v-list-tile-title>
          </v-list-tile>
        </v-list>
      </v-menu>
    </v-form>
  </v-flex>

</template>

<script>
  import {mapActions} from 'vuex'

  export default {
    name: 'search-field',
    props: [
      'performSearch'
    ],
    computed: {
      searchText: {
        get() {
          return this.$store.state.searchTerm
        },
        set(value) {
          this.$store.commit('ChangeSearch', value)
        }
      },
      suggestions: {
        get() {
          return this.$store.state.suggestions
        },
        set(value) {
          this.$store.commit('ChangeSuggestions', value)
        }
      }
    },
    data: () => ({
      openMenu: false,
      snackbar: false,
      validationText: '',
      rules: {
        minLength: (value) => (value !== null && value.trim().length >= 3) || 'Enter 3 or more chars'
      }
    }),
    watch: {
      searchText(value) {
        if (value) {
          this.openMenu = true
          this.fetchSuggestions(value)
        } else {
          this.suggestions = []
          this.openMenu = false
        }
      }
    },
    methods: Object.assign(mapActions(['fetchSuggestions']), {
      performSearchInternal() {
        // internal staff before or after performSearch
        this.openMenu = false
        if (this.$refs.form.validate()) {
          this.performSearch()
        } else {
          this.validationText = 'You should enter at least 3 non-whitespace characters to perform search!'
          this.snackbar = true
        }
      }
    })
  }
</script>

<style scoped>

</style>
