<template lang="html">
  <el-card  class="box-card">
    <div slot="header" class="clearfix">
      <span>Share "<b>{{filename}}</b>"</span>
      <el-button style="float: right; padding: 3px 10px" type="text" @click="handleCancel">Cancel</el-button>
      <el-button style="float: right; padding: 3px 10px" type="text" @click="handleSave">Save</el-button>
    </div>
    <el-checkbox-group v-model="selected">
      <el-checkbox v-for="user in users" :label="user.id" border>{{user.email}}</el-checkbox>
    </el-checkbox-group>

  </el-card>
</template>

<script>
import Auth from '../assets/auth'
let auth

export default {
  name: 'ShareList',
  props: ['filename'],
  data () {
    return {
      users: [],
      selected: []
    }
  },
  mounted () {
    auth = new Auth(this)

    auth.GetUsersByEmail('')
    .then(response => {
      this.users = response
    })
  },
  methods: {
    init (data) {
      this.selected = data
    },
    handleSave () {
      this.$emit('save-shared', this.selected)
    },
    handleCancel () {
      this.$emit('cancel-shared')
    }
  }
}
</script>

<style lang="css">

</style>
