<template lang="html">
  <div class="page">
    <h1>Create a new account</h1>
    <credentials-form submit-text="Register" alt-text="Login into an existing account"  @submit="register" @alt="toLogin"></credentials-form>
  </div>
</template>

<script>
import CredentialsForm from './CredentialsForm'
import Auth from '../assets/auth.js'
import Utils from '../assets/utils'

let utils

export default {
  name: 'RegisterPage',
  components: {
    CredentialsForm
  },
  mounted () {
    utils = new Utils(this)
  },
  methods: {
    register (data) {
      const auth = new Auth(this)

      auth.Register(data.email, data.password)
        .then(response => {
          this.$router.push({
            path: '/code', query: { email: response.email }
          })
        })
        .catch(err => {
          utils.handleErr(err)
        })
    },
    toLogin () {
      this.$router.push('/login')
    }
  }
}
</script>

<style lang="css">
</style>
