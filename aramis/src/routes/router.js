import Vue from 'vue'
import Router from 'vue-router'

import Login from '@/components/Login'
import Callback from '@/components/Callback'

import Dashboard from '@/components/Dashboard'
import Alerts from '@/components/Alerts'
import Servers from '@/components/Servers'
import Profile from '@/components/Profile'

Vue.use(Router)

export default new Router({
  mode: 'history',
  routes: [{
      path: '/login',
      name: 'Login',
      component: Login
    },
    {
      path: '/callback',
      name: 'Callback',
      component: Callback
    },
    {
      path: '/dashboard',
      name: 'Dashboard',
      component: Dashboard
    },
    {
      path: '/alerts',
      name: 'Alerts',
      component: Alerts
    },
    {
      path: '/servers',
      name: 'Servers',
      component: Servers
    },
    {
      path: '/profile',
      name: 'Profile',
      component: Profile
    },
  ]
})
