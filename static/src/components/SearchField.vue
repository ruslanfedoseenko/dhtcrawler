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
          ref="searchInput"
          slot="activator"
          placeholder="Search..."
          single-line
          append-icon="search"
          @click:append="performSearchInternal"
          v-model="searchTextProp"
          :rules="[rules.minLength]"
          dark
          hide-details
        />
        <v-list style="max-height: 300px">
          <v-list-tile v-for="item in suggestions" :key="item" @click="onWordSelected(item,$event)">
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
      'performSearch', 'searchText'
    ],
    computed: {
      suggestions: {
        get() {
          return this.$store.state.search.suggestions
        },
        set(value) {
          this.$store.commit('ChangeSuggestions', value)
        }
      },
      searchTextProp: {
        get() {
          return this.searchText
        },
        set(value) {
          this.$emit('update:searchText', value)
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
      searchText(value, prevValue) {
        if (value) {
          this.openMenu = true
          let input = this.getInput(this.$refs.searchInput)
          this.fetchSuggestions(value.slice(0, input.selectionStart))
        } else {
          this.suggestions = []
          this.openMenu = false
        }
      }
    },
    mounted() {
      let input = this.getInput(this.$refs.searchInput)
      input.addEventListener('keydown', this.manualSubmitCheck.bind(this))
    },
    beforeDestroy() {
      let input = this.getInput(this.$refs.searchInput)
      input.removeEventListener('keydown', this.manualSubmitCheck.bind(this))
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
      },
      manualSubmitCheck(e) {
        if (e.keyCode === 13) {
          this.performSearchInternal()
          e.stopPropagation()
        }
      },
      findFirstDiffPos(a, b) {
        let i = 0
        if (a === b) return -1
        while (a[i] === b[i]) i++
        return i
      },
      onWordSelected(word, e) {
        e.preventDefault()
        let input = this.getInput(this.$refs.searchInput)
        let wordBoundaries = this.getWordBoundaries(this.searchTextProp, input.selectionStart)
        this.searchTextProp = this.splice(this.searchTextProp, wordBoundaries[0], wordBoundaries[1] - wordBoundaries[0], word)
        input.focus()
        input.scrollIntoView()
        if (!this.openMenu) {
          this.openMenu = true
        }
      },
      splice(s, start, delCount, newSubStr) {
        return s.slice(0, start) + newSubStr + s.slice(start + Math.abs(delCount))
      },
      getWordBoundaries(s, pos) {
        // make pos point to a character of the word
        while (s[pos] === ' ') pos--
        // find the space before that word
        // (add 1 to be at the begining of that word)
        // (note that it works even if there is no space before that word)
        pos = s.lastIndexOf(' ', pos) + 1
        // find the end of the word
        var end = s.indexOf(' ', pos)
        if (end === -1) end = s.length // set to length if it was the last word
        // return the result
        return [pos, end]
      },
      getInput(vuetifyInput) {
        if (!vuetifyInput) {
          return null
        }
        return vuetifyInput.$el.children[0].children[0]
      }
    })

  }
</script>

<style scoped>

</style>
