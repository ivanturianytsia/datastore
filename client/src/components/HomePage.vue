<template lang="html">
  <div>
    <el-button @click="handleSelectFile" :loading="uploadLoading">Upload</el-button>
    <el-button class="float-right" @click="logout">Log Out from {{ email }}</el-button>
    <el-menu :default-active="activeFilesTab" mode="horizontal" @select="handleTabSelect">
      <el-menu-item index="files">My Files</el-menu-item>
      <el-menu-item index="shared">Shared</el-menu-item>
    </el-menu>

    <div class="page">
      <el-table
        :data="filesInTab"
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
          prop="filename"
          label="Name">
        </el-table-column>
        <el-table-column
          prop="size"
          :formatter="bytesToSize"
          label="Size"
          width="180">
        </el-table-column>
        <el-table-column
          prop="created"
          :formatter="since"
          label="Created"
          width="180">
        </el-table-column>
        <el-table-column
          label="Operations">
          <template slot-scope="scope">
            <el-button @click="handleDownload(scope.$index, scope.row)" type="text" size="small">Download</el-button>
            <el-button v-if="activeFilesTab === 'files'" @click="handleDelete(scope.$index, scope.row)" type="text" size="small">Delete</el-button>
            <el-button v-if="activeFilesTab === 'files'" @click="handleShare(scope.$index, scope.row)" type="text" size="small">Share</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="share-modal">
        <share-list v-show="shareId" :filename="shareFilename" @save-shared="handleSaveShared" @cancel-shared="handleCancelShared" ref="sharelist"></share-list>
      </div>

      <input id="fileinput" type='file' @change="upload" ref="fileinput" multiple/>
    </div> <!--page-->
  </div>  <!--template-->
</template>

<script>
import ShareList from './ShareList'
import Auth from '../assets/auth'
import Files from '../assets/files'
import moment from 'moment'

let auth
let files

export default {
  name: 'HomePage',
  data () {
    return {
      email: '',
      id: '',
      activeFilesTab: 'files',
      files: [],
      shared: [],
      uploadLoading: false,
      shareId: '',
      shareFilename: ''
    }
  },
  mounted () {
    files = new Files(this)
    auth = new Auth(this)
    if (!auth.IsLogged()) {
      this.logout()
      return
    }

    auth.GetUser()
    .then(response => {
      this.email = response.email
      this.id = response.id
      return files.GetFiles()
    })
    .then(this.handleFiles)
    .catch(this.handleErr)
  },
  computed: {
    filesInTab () {
      if (this.activeFilesTab === 'shared') {
        return this.shared
      } else {
        return this.files
      }
    }
  },
  methods: {
    logout () {
      auth.Logout()
      this.$router.push('/login')
    },
    handleTabSelect (tab) {
      this.activeFilesTab = tab
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
    handleFiles (response) {
      this.files = response.files
      this.shared = response.shared
    },
    handleErr (err) {
      console.log(err)
      // this.logout()
    },
    upload (event) {
      const data = this.$refs.fileinput.files
      if (data && data.length) {
        this.uploadLoading = true
        files.UploadFiles(data)
        .then(() => {
          return files.GetFiles()
        })
        .then(response => {
          this.uploadLoading = false
          this.handleFiles(response)
        })
        .catch(this.handleErr)
      }
    },
    handleDownload (index, row) {
      window.open(`/files/${row.path}?token=${auth.token}`)
    },
    handleDelete (index, row) {
      files.DeleteFile(row.filename)
      .then(() => {
        return files.GetFiles()
      })
      .then(this.handleFiles)
      .catch(this.handleErr)
    },
    handleShare (index, row) {
      this.shareId = row.id
      this.shareFilename = row.filename

      let defaults = []
      for (let i in row.allowedids) {
        defaults.push(i)
      }
      this.$refs.sharelist.init(defaults)
    },
    handleSaveShared (data) {
      let id = this.shareId
      if (id) {
        this.shareId = ''
        this.shareFilename = ''
        let allowedids = {}
        for (let i in data) {
          allowedids[data[i]] = {}
        }
        files.UpdateFile(id, allowedids)
        .then(response => {
          return files.GetFiles()
        })
        .then(this.handleFiles)
        .catch(this.handleErr)
      }
    },
    handleCancelShared () {
      this.shareId = ''
      this.shareFilename = ''
    }
  },
  components: {
    ShareList
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
.float-right {
  float: right;
}
.hidden {
  display: none;
}
.share-modal {
  position: relative;
  padding-top: 50px;
}
</style>
