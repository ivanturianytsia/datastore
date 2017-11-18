import Vue from 'vue'
import Router from 'vue-router'
import LoginPage from '@/components/LoginPage'
import RegisterPage from '@/components/RegisterPage'
import HomePage from '@/components/HomePage'
import UploadPage from '@/components/UploadPage'

Vue.use(Router)

export default new Router({
  routes: [{
    path: '/home',
    component: HomePage
  },
  {
    path: '/register',
    component: RegisterPage
  },
  {
    path: '/login',
    component: LoginPage
  },
  {
    path: '/upload',
    component: UploadPage
  },
  { path: '*', redirect: '/home' }]
})
