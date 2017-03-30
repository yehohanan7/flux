import Vue from 'vue'
import VueRouter from 'vue-router'

import Account from './account/Account.vue'
import Event from './event/Event.vue'
import Create from './create/Create.vue'

Vue.use(VueRouter)

const router = new VueRouter({
  mode: 'history',
  base: __dirname,
  routes: [
    { path: '/', component: Create },
    { path: '/account', component: Account },
    { path: '/event', component: Event }
  ]
})

new Vue({
  router,
  template: `
    <div id="app">
      <router-view class="view"></router-view>
    </div>
  `
}).$mount('#app')