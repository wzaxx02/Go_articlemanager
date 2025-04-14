import Vue from 'vue'
import App from './App.vue'
import router from './router'
import vuetify from './plugins/vuetify'
import day from 'dayjs'

import './plugins/http'

import 'vuetify/dist/vuetify.min.css'
// import 'vuetify/styles'
// import { createVuetify } from 'vuetify'
// import * as components from 'vuetify/components'
// import * as directives from 'vuetify/directives'
// const vuetify = createVuetify({
  // components,
  // directives,
// })

Vue.filter('dateformat', function(indate, outdate) {
  return day(indate).format(outdate)
})

Vue.config.productionTip = false

console.log(router);


new Vue({
  router,
  vuetify,
  render: h => h(App)
}).$mount('#app')
