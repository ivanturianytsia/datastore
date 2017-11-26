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
          label=""
          width="90">
          <template slot-scope="scope">
            <el-dropdown class="more-dropdown" @command="handleMore" trigger="click">
              <span class="el-dropdown-link">
                More<i class="el-icon-arrow-down el-icon--right"></i>
              </span>
              <el-dropdown-menu slot="dropdown">
                <el-dropdown-item class="el-dropdown-item" :command="{ command: 'download', row: scope.row }">Download</el-dropdown-item>
                <el-dropdown-item class="el-dropdown-item" :command="{ command: 'delete', row: scope.row }" v-if="activeFilesTab === 'files'">Delete</el-dropdown-item>
                <el-dropdown-item class="el-dropdown-item" :command="{ command: 'share', row: scope.row }" v-if="activeFilesTab === 'files'">Share</el-dropdown-item>
              </el-dropdown-menu>
            </el-dropdown>
          </template>
        </el-table-column>
      </el-table>

      <share-list v-show="shareId" @save-shared="handleSaveShared" @cancel-shared="handleCancelShared" ref="sharelist"></share-list>

      <input id="fileinput" type='file' @change="upload" ref="fileinput" multiple/>
    </div> <!--page-->
  </div>  <!--template-->
</template>

<script>
import ShareList from './ShareList'
import Auth from '../assets/auth'
import Files from '../assets/files'
import Utils from '../assets/utils'
import moment from 'moment'

let auth
let files
let utils

export default {
  name: 'HomePage',
  data () {
    return {
      email: '',
      id: '',
      activeFilesTab: 'files',
      files: [],
      shared: [],
      uploadLoading: false
    }
  },
  mounted () {
    files = new Files(this)
    auth = new Auth(this)
    utils = new Utils(this)
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
    .catch(err => {
      utils.handleErr(err)
    })
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
    handleSucc (text) {
      const h = this.$createElement
      this.$message({
        message: h('p', null, text)
      })
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
          this.handleSucc('File uploaded')
        })
        .catch(err => {
          utils.handleErr(err)
        })
      }
    },
    handleMore (data) {
      switch (data.command) {
        case 'download':
          this.handleDownload(data.row)
          break
        case 'delete':
          this.handleDelete(data.row)
          break
        case 'share':
          this.handleShare(data.row)
          break
      }
    },
    handleDownload (row) {
      window.open(`/files/${row.path}?token=${auth.token}`)
    },
    handleDelete (row) {
      files.DeleteFile(row.filename)
      .then(() => {
        return files.GetFiles()
      })
      .then(this.handleFiles)
      .catch(err => {
        utils.handleErr(err)
      })
    },
    handleShare (row) {
      let defaults = []
      for (let i in row.allowedids) {
        defaults.push(i)
      }
      this.$refs.sharelist.init({
        id: row.id,
        name: row.filename,
        selected: defaults
      })
    },
    handleSaveShared (data) {
      let id = data.id
      if (id) {
        let allowedids = {}
        for (let i in data.selected) {
          allowedids[data.selected[i]] = {}
        }
        files.UpdateFile(id, allowedids)
        .then(response => {
          this.handleSucc('Saved')
          return files.GetFiles()
        })
        .then(this.handleFiles)
        .catch(err => {
          utils.handleErr(err)
        })
      }
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
.more-dropdown {
  float: right;
}
.el-dropdown-link {
  cursor: pointer;
  color: #409EFF;
}
.el-icon-arrow-down {
  font-size: 12px;
}
.el-dropdown-item {
  font-family: inherit;
}
</style>
