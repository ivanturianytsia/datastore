<template lang="html">
  <div class="page">
    <h1>Log in into an existing account</h1>
    <credentials-form submit-text="Log In" alt-text="Create new account" @submit="login" @alt="toRegister"></credentials-form>
  </div>
</template>

<script>
import CredentialsForm from './CredentialsForm'
import Auth from '../assets/auth.js'

export default {
  name: 'LoginPage',
  components: {
    CredentialsForm
  },
  methods: {
    login (data) {
      const auth = new Auth(this)

      auth.Login(data.email, data.password)
        .then(response => {
          this.$router.push({
            path: '/code', query: { email: response.email }
          })
        })
        .catch(err => {
          console.log(err)
        })
    },
    toRegister () {
      this.$router.push('/register')
    }
  }
}
</script>

<style lang="css">
</style>
