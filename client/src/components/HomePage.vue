<template lang="html">
  <div class="page">
    <h1>Welcome home {{ email }}</h1>
    <el-button @click="logout">Log Out</el-button>
  </div>
</template>

<script>
import Auth from '../assets/auth.js'

let auth

export default {
  name: 'HomePage',
  data () {
    return {
      email: ''
    }
  },
  mounted () {
    auth = new Auth(this)
    if (!auth.IsLogged()) {
      this.$router.push('/login')
      return
    }
    this.getUser()
  },
  methods: {
    logout () {
      auth.Logout()
      this.$router.push('/login')
    },
    getUser () {
      const that = this
      auth.GetUser()
        .then(response => {
          that.email = response.email
        })
        .catch(err => {
          console.log(err)
          this.$router.push('/login')
        })
    }
  }
}
</script>

<style lang="css">
</style>
