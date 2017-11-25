<template lang="html">
  <el-checkbox v-for="user in users" v-model="selected[user.id]" :label="user.email" border></el-checkbox>
  <el-button @click="handleSave">Save</el-button>
</template>

<script>
import Auth from '../assets/auth'
let auth

export default {
  name: 'ShareList',
  props: ['fileId', 'defaultIds'],
  data () {
    return {
      selected: {},
      users: []
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
    handleSave () {
      console.log(this.selected)
    }
  }
}
</script>

<style lang="css">

</style>
