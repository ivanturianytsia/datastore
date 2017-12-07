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

    <el-input
      type="textarea"
      :autosize="{ minRows: 4, maxRows: 10}"
      placeholder="Describe file"
      v-model="description">
    </el-input>

    <span slot="footer" class="dialog-footer">
      <el-button @click="handleCancel">Cancel</el-button>
      <el-button type="primary" @click="handleSave">Save</el-button>
    </span>
  </el-dialog>
</template>

<script>

export default {
  name: 'EditField',
  data () {
    return {
      dialogVisible: false,
      name: '',
      id: '',
      description: ''
    }
  },
  computed: {
    title () {
      return `Edit file "${this.name}"`
    }
  },
  methods: {
    init (data) {
      this.name = data.name
      this.id = data.id
      this.description = data.description
      this.dialogVisible = true
    },
    handleSave () {
      this.dialogVisible = false
      this.$emit('save', {
        id: this.id,
        description: this.description
      })
    },
    handleCancel () {
      this.dialogVisible = false
      this.$emit('cancel')
    }
  }
}
</script>

<style lang="css">

</style>
