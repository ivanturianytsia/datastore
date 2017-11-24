<template lang="html">
  <div class="page">
    <h1>Welcome home {{ email }}</h1>
    <el-button @click="logout">Log Out</el-button>
    <el-button @click="handleSelectFile" :loading="uploadLoading">Upload</el-button>

    <el-table
      :data="files"
      style="width: 100%"
      empty-text="You have no files yet">
      <el-table-column
      label=""
      width="40">
      <template slot-scope="scope">
        <i class="fa fa-file-o"></i>
      </template>
      </el-table-column>
      <el-table-column
        prop="name"
        label="Name">
      </el-table-column>
      <el-table-column
        prop="size"
        :formatter="bytesToSize"
        label="Size"
        width="180">
      </el-table-column>
      <el-table-column
        prop="modified"
        :formatter="since"
        label="Modified"
        width="180">
      </el-table-column>
      <el-table-column
        label="Operations">
        <template slot-scope="scope">
          <el-button @click="handleDownload(scope.$index, scope.row)" type="text" size="small">Download</el-button>
          <el-button @click="handleDelete(scope.$index, scope.row)" type="text" size="small">Delete</el-button>
        </template>
      </el-table-column>
    </el-table>

    <input id="fileinput" type='file' @change="upload" ref="fileinput" multiple/>
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
      files: [],
      uploadLoading: false
    }
  },
  mounted () {
    files = new Files(this)
    auth = new Auth(this)
    if (!auth.IsLogged()) {
      this.logout()
      return
    }

    const that = this

    auth.GetUser()
    .then(response => {
      that.email = response.email
      return files.GetFiles()
    })
    .then(response => {
      that.files = response
    })
    .catch(this.handleErr)
  },
  methods: {
    logout () {
      auth.Logout()
      this.$router.push('/login')
    },
    since (r, c, val) {
      return moment(new Date(val)).fromNow()
    },
    bytesToSize (r, c, val) {
      const bytes = val
      const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB']
      if (bytes === 0) {
        return 'n/a'
      }
      const i = parseInt(Math.floor(Math.log(bytes) / Math.log(1024)), 10)
      if (i === 0) {
        return `${bytes} ${sizes[i]}`
      }
      return `${(bytes / (1024 ** i)).toFixed(1)} ${sizes[i]}`
    },
    handleSelectFile (event) {
      this.$refs.fileinput.click()
    },
    handleErr (err) {
      console.log(err)
      this.logout()
    },
    upload (event) {
      const that = this
      const data = this.$refs.fileinput.files
      if (data && data.length) {
        this.uploadLoading = true
        files.UploadFiles(data)
        .then(() => {
          return files.GetFiles()
        })
        .then(response => {
          this.uploadLoading = false
          that.files = response
        })
        .catch(this.handleErr)
      }
    },
    handleDownload (index, row) {
      window.open(`/files/${row.name}?token=${auth.token}`)
    },
    handleDelete (index, row) {
      const that = this
      files.DeleteFile(row.name)
      .then(() => {
        return files.GetFiles()
      })
      .then(response => {
        that.files = response
      })
      .catch(this.handleErr)
    }
  }
}
</script>

<style lang="css">
#fileinput {
  width: 0.1px;
  height: 0.1px;
  opacity: 0;
  overflow: hidden;
  position: absolute;
  z-index: -1;
}
</style>
