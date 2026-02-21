import { defineStore } from 'pinia'
import axios from 'axios'

export interface PluginUiConfig {
  entry: string
  width?: number
  height?: number
  location?: string // 'sidebar' | 'main' | 'dialog'
}

export interface PluginCapabilities {
  file_extensions: string[]
  actions: string[]
}

export interface PluginInfo {
  id: string
  ui?: PluginUiConfig
  capabilities?: PluginCapabilities
  registered_at: number
  last_heartbeat: number
}

export const usePluginStore = defineStore('plugin', {
  state: () => ({
    plugins: [] as PluginInfo[],
    activePluginId: null as string | null,
    loading: false,
    error: null as string | null,
  }),

  getters: {
    sidebarPlugins: (state) => state.plugins.filter(p => p.ui?.location === 'sidebar' || !p.ui?.location),
    getPluginById: (state) => (id: string) => state.plugins.find(p => p.id === id),
    uiPlugins: (state) => state.plugins.filter(p => !!p.ui),
  },

  actions: {
    async fetchPlugins() {
      this.loading = true
      try {
        const response = await axios.get('http://localhost:3000/api/v1/plugins')
        this.plugins = response.data.items || []
        this.error = null
      } catch (err: any) {
        this.error = err.message || 'Failed to fetch plugins'
        console.error('Error fetching plugins:', err)
      } finally {
        this.loading = false
      }
    },

    // 启动定时轮询，实现“插件启动即发现”
    initPolling(intervalMs = 3000) {
      this.fetchPlugins()
      const timer = setInterval(() => {
        this.fetchPlugins()
      }, intervalMs)
      return () => clearInterval(timer)
    }
  }
})
