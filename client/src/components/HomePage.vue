<template lang="html">
  <div class="page">
    <h1>Welcome home {{ email }}</h1>
    <el-button @click="logout">Log Out</el-button>
    <el-button @click="toUpload">Upload</el-button>

    <el-table
      :data="files"
      style="width: 100%">
      <el-table-column
        prop="name"
        label="Name">
      </el-table-column>
      <el-table-column
        prop="size"
        label="Size"
        width="180">
      </el-table-column>
      <el-table-column
        prop="modified"
        :formatter="since"
        label="Modified"
        width="180">
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
import Auth from '../assets/auth.js'
import Files from '../assets/files.js'
import moment from 'moment'

let auth
let files

export default {
  name: 'HomePage',
  data () {
    return {
      email: '',
      files: []
    }
  },
  mounted () {
    files = new Files(this)
    auth = new Auth(this)
    if (!auth.IsLogged()) {
      this.$router.push('/login')
      return
    }
    this.getUser()
    this.getFiles()
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
    },
    getFiles () {
      const that = this
      files.GetFiles()
        .then(response => {
          that.files = response
        })
    },
    since (r, c, val) {
      return moment(new Date(val)).fromNow()
    },
    toUpload () {
      this.$router.push('/upload')
    }
  }
}
</script>

<style lang="css">
</style>
