<template lang="html">
  <div class="page">
    <h1>Upload files</h1>
    <el-button @click="toHome">Go back</el-button>
    <el-upload
      class="upload-demo"
      action="/upload"
      drag
      :file-list="uploadList"
      :headers="uploadHeaders"
      multiple>
      <i class="el-icon-upload"></i>
      <div class="el-upload__text">Drop file here or <em>click to upload</em></div>
    </el-upload>
  </div>
</template>

<script>
import Auth from '../assets/auth.js'

let auth

export default {
  name: 'UploadPage',
  data () {
    return {
      uploadList: [],
      uploadHeaders: {}
    }
  },
  mounted () {
    auth = new Auth(this)
    if (!auth.IsLogged()) {
      this.$router.push('/login')
      return
    }
    this.uploadHeaders = {
      'Authorization': `Bearer ${auth.token}`
    }
    this.getUser()
  },
  methods: {
    getUser () {
      const that = this
      auth.GetUser()
        .catch(err => {
          console.log(err)
          that.$router.push('/login')
        })
    },
    toHome () {
      this.$router.push('/home')
    }
  }
}
</script>

<style lang="css">
</style>
