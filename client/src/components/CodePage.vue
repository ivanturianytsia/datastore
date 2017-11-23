<template lang="html">
  <div class="page">
    <h1>We sent you a code to your email, enter it below. <small>Note that we never send you any links in our emails</small></h1>
    <el-form label-position="top" label-width="100px">
      <el-form-item label="Secret Code">
        <el-input v-model="code" type="code" auto-complete="off"></el-input>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="submitCode">Submit Code</el-button>
        <el-button @click="goBack">Go back</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script>
import Auth from '../assets/auth.js'

export default {
  name: 'CodePage',
  data () {
    return {
      code: ''
    }
  },
  methods: {
    goBack () {
      this.$router.go(-1)
    },
    submitCode () {
      const that = this
      const auth = new Auth(this)
      auth.PostEmailCode(this.$route.query.email, this.code)
        .then(response => {
          that.$router.push('/home')
        })
        .catch(err => {
          console.log(err)
        })
    }
  }

}
</script>

<style lang="css">
</style>
