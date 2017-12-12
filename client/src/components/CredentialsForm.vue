<template lang="html">
  <el-form label-position="top" label-width="100px">
    <el-form-item label="Email">
      <el-input v-model="email" type="email" auto-complete="off"></el-input>
    </el-form-item>
    <el-form-item label="Password">
      <el-input v-model="password" type="password" auto-complete="off"></el-input>
    </el-form-item>
    <el-form-item label="Phone Number" v-if="submitText === 'Register'">
      <el-input v-model="phonenumber" placeholder="+00000000000" type="phonenumber" auto-complete="off"></el-input>
    </el-form-item>
    <el-form-item>
      <el-button type="primary" @click="submit">{{submitText}}</el-button>
      <el-button @click="altAction">{{altText}}</el-button>
    </el-form-item>
  </el-form>
</template>

<script>
export default {
  name: 'CredentialsForm',
  props: ['submitText', 'altText'],
  data () {
    return {
      email: '',
      password: '',
      phonenumber: ''
    }
  },
  methods: {
    submit (event) {
      this.$emit('submit', {
        email: this.email,
        password: this.password,
        phonenumber: this.phonenumber
      })
    },
    altAction (event) {
      this.$emit('alt')
    }
  }
}
</script>

<style lang="css">
</style>
