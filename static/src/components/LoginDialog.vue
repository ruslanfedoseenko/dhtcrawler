<template>
  <v-layout>
    <v-btn @click.native="isLoginDialogActive = true" color="success">
      LOGIN
    </v-btn>
    <v-dialog max-width="340px" hide-overlay persistent v-model="isLoginDialogActive">

      <v-card>
        <v-card-title>
          <span class="headline">Login</span>
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
            </v-layout>
          </v-container>
          <small>*indicates required field</small>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="blue darken-1" flat @click="close">Close</v-btn>
          <v-btn :loading="login"
                 :disabled="login" color="blue darken-1" flat @click="doLogin">
            Login
            <span slot="loader" class="custom-loader">
              <v-icon>cached</v-icon>
            </span>
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <v-snackbar
      :timeout="6000"
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
    $_veeValidate: {
      validator: 'new'
    },
    name: 'btoogle-login-dialog',
    data() {
      return {
        username: '',
        password: '',
        snackbar: false,
        isLoginDialogActive: false,
        login: false,
        validationText: ''
      }
    },
    methods: Object.assign(mapActions(['performLogin']), {
      showError(message) {
        this.validationText = message
        this.snackbar = true
      },
      OnLoginFailure(resp) {
        this.login = false
        let errCode = resp.body
        this.showError(errCode.message)
      },
      OnLoginSuccess(resp) {
        this.login = false
        this.isLoginDialogActive = false
      },
      doLogin() {
        this.login = true
        this.$validator.validateAll().then(
          (isValid) => {
            if (!isValid) {
              this.login = false
              return
            }
            this.performLogin({
              username: this.username,
              password: this.password
            }).then(
              this.OnLoginSuccess,
              this.OnLoginFailure
            )
          }
        )
      },
      close() {
        this.username = ''
        this.password = ''
        this.$validator.reset()
        this.isLoginDialogActive = false
      }
    })

  }
</script>
<style>
  .custom-loader {
    animation: loader 1s infinite;
    display: flex;
  }

  @-moz-keyframes loader {
    from {
      transform: rotate(0);
    }
    to {
      transform: rotate(360deg);
    }
  }

  @-webkit-keyframes loader {
    from {
      transform: rotate(0);
    }
    to {
      transform: rotate(360deg);
    }
  }

  @-o-keyframes loader {
    from {
      transform: rotate(0);
    }
    to {
      transform: rotate(360deg);
    }
  }

  @keyframes loader {
    from {
      transform: rotate(0);
    }
    to {
      transform: rotate(360deg);
    }
  }
</style>
