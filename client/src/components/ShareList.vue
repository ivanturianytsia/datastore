<template lang="html">
  <el-dialog
    :title="title"
    :visible.sync="dialogVisible"
    width="30%"
    :before-close="handleCancel">
    <div slot="header" class="clearfix">
      <span></span>
      <el-button style="float: right; padding: 3px 10px" type="text" @click="handleCancel">Cancel</el-button>
      <el-button style="float: right; padding: 3px 10px" type="text" @click="handleSave">Save</el-button>
    </div>
    <el-checkbox-group v-model="selected">
      <el-checkbox v-for="user in users" :label="user.id" border>{{user.email}}</el-checkbox>
    </el-checkbox-group>
    <span slot="footer" class="dialog-footer">
      <el-button @click="handleCancel">Cancel</el-button>
      <el-button type="primary" @click="handleSave">Save</el-button>
    </span>
  </el-dialog>
</template>

<script>
import Auth from '../assets/auth'
import Utils from '../assets/utils'

let auth
let utils

export default {
  name: 'ShareList',
  data () {
    return {
      dialogVisible: false,
      users: [],
      selected: [],
      name: '',
      id: ''
    }
  },
  mounted () {
    auth = new Auth(this)
    utils = new Utils(this)

    auth.GetUsersByEmail('')
    .then(response => {
      this.users = response
    })
    .catch(err => {
      utils.handleErr(err)
    })
  },
  computed: {
    title () {
      return `Share file "${this.name}"`
    }
  },
  methods: {
    init (data) {
      this.name = data.name
      this.id = data.id
      this.selected = data.selected
      this.dialogVisible = true
    },
    handleSave () {
      this.dialogVisible = false
      this.$emit('save-shared', {
        id: this.id,
        selected: this.selected
      })
    },
    handleCancel () {
      this.dialogVisible = false
      this.$emit('cancel-shared')
    }
  }
}
</script>

<style lang="css">

</style>
