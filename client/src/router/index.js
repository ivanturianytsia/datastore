import Vue from 'vue'
import Router from 'vue-router'
import LoginPage from '@/components/LoginPage'
import RegisterPage from '@/components/RegisterPage'
import HomePage from '@/components/HomePage'

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
  { path: '*', redirect: '/home' }]
})
