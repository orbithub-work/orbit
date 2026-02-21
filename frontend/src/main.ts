// 智归档OS - 前端应用主入口
// 
// 这是一个基于Vue 3 + TypeScript + Pinia的前端应用

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router'
import App from './App.vue'
import 'vue-virtual-scroller/dist/vue-virtual-scroller.css'
import VueVirtualScroller from 'vue-virtual-scroller'
import './styles/index.css'
import { initPluginSystem } from './core/pluginLoader'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)
app.use(VueVirtualScroller)
app.mount('#app')

// Initialize plugin system after app is mounted
initPluginSystem().catch(err => {
  console.error('Failed to initialize plugin system:', err)
})
