<template>
  <v-layout>
    <v-btn @click.native="isRegDialogActive = true" color="info">
      REGISTER
    </v-btn>
    <v-dialog max-width="340px" hide-overlay persistent v-model="isRegDialogActive">

      <v-card>
        <v-card-title>
          <span class="headline">Registration</span>
        </v-card-title>
        <v-card-text>
          <v-container grid-list-md>
            <v-layout wrap>
              <v-flex xs12>
                <v-text-field v-model="username" data-vv-name="username" :error-messages="errors.collect('username')"
                              label="Login" v-validate="'required'" required></v-text-field>
              </v-flex>
              <v-flex xs12>
                <v-text-field v-model="password" data-vv-name="password" :error-messages="errors.collect('password')"
                              label="Password" type="password" v-validate="'required'" required></v-text-field>
              </v-flex>
              <v-flex xs12>
                <v-text-field v-model="email" data-vv-name="email" :error-messages="errors.collect('email')"
                              label="E-Mail" type="email" v-validate="'required'" required></v-text-field>
              </v-flex>
            </v-layout>
          </v-container>
          <small>*indicates required field</small>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="blue darken-1" flat @click="close">Close</v-btn>
          <v-btn :loading="register"
                 :disabled="register" color="blue darken-1" flat @click="doRegister">
            Register
            <span slot="loader" class="custom-loader">
              <v-icon>cached</v-icon>
            </span>
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <v-snackbar
      :timeout="60000"
      :top="true"
      style="z-index: 10000"
      :multi-line="true"
      v-model="snackbar"
    >
      {{ validationText }}
      <v-btn flat color="pink" @click.native="snackbar = false">Close</v-btn>
    </v-snackbar>
  </v-layout>
</template>

<script>
  import {mapActions} from 'vuex'

  export default {
    name: 'btoogle-register-dialog',
    $_veeValidate: {
      validator: 'new'
    },
    data() {
      return {
        username: '',
        password: '',
        email: '',
        snackbar: false,
        isRegDialogActive: false,
        register: false,
        validationText: ''
      }
    },
    methods: Object.assign(mapActions(['performRegistration']), {
      close() {
        this.isRegDialogActive = false
      },
      OnRegSuccess(resp) {

      },
      OnRegFailure(resp) {

      },
      doRegister() {
        this.register = true
        this.$validator.validateAll().then(
          (isValid) => {
            if (!isValid) {
              this.register = false
              return
            }
            this.performRegistration({
              username: this.username,
              password: this.password,
              mail: this.email
            }).then(
              this.OnRegSuccess,
              this.OnRegFailure
            )
          }
        )
      }
    })
  }
</script>

<style scoped>

</style>
