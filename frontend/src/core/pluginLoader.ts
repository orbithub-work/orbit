/**
 * Plugin Loader
 * 
 * Loads and manages frontend plugins
 */

import { ref, reactive } from 'vue'
import type { Component } from 'vue'

export interface PluginInfo {
  plugin_id: string
  name: string
  version: string
  mode: 'frontend' | 'local_process' | 'network_service'
  endpoint?: string
  mounts: PluginMount[]
  online: boolean
}

export interface PluginMount {
  slot: string
  entry: string
  title?: string
  icon?: string
  order?: number
  width?: number
  height?: number
}

export interface LoadedPlugin {
  info: PluginInfo
  module?: any
  components: Map<string, Component>
  actions: Map<string, Function>
}

class PluginLoader {
  private plugins = reactive<Map<string, LoadedPlugin>>(new Map())
  private loading = ref(false)
  private baseURL = 'http://localhost:32000'

  async loadAll() {
    if (this.loading.value) return
    this.loading.value = true

    try {
      // Fetch plugin list from backend
      const response = await fetch(`${this.baseURL}/api/plugins/list`)
      const result = await response.json()
      const pluginList: PluginInfo[] = result.data || []

      // Load each plugin
      for (const pluginInfo of pluginList) {
        if (pluginInfo.mode === 'frontend') {
          await this.loadFrontendPlugin(pluginInfo)
        }
      }
    } catch (error) {
      console.error('Failed to load plugins:', error)
    } finally {
      this.loading.value = false
    }
  }

  private async loadFrontendPlugin(info: PluginInfo) {
    try {
      // For frontend plugins, the entry should be a JS module URL
      // In development, this might be http://localhost:8001/plugin.js
      // In production, this might be /plugins/copyright/index.js
      
      const entry = info.endpoint || info.mounts[0]?.entry
      if (!entry) {
        console.warn(`Plugin ${info.plugin_id} has no entry point`)
        return
      }

      // Dynamic import the plugin module
      const module = await import(/* @vite-ignore */ entry)
      
      const loadedPlugin: LoadedPlugin = {
        info,
        module,
        components: new Map(),
        actions: new Map()
      }

      // Register components
      if (module.components) {
        for (const [name, component] of Object.entries(module.components)) {
          loadedPlugin.components.set(name, component as Component)
        }
      }

      // Register actions
      if (module.actions) {
        for (const [name, action] of Object.entries(module.actions)) {
          loadedPlugin.actions.set(name, action as Function)
        }
      }

      this.plugins.set(info.plugin_id, loadedPlugin)
      console.log(`✅ Loaded plugin: ${info.name}`)
    } catch (error) {
      console.error(`Failed to load plugin ${info.plugin_id}:`, error)
    }
  }

  getPlugin(pluginId: string): LoadedPlugin | undefined {
    return this.plugins.get(pluginId)
  }

  getPluginComponent(pluginId: string, componentName: string): Component | undefined {
    const plugin = this.plugins.get(pluginId)
    return plugin?.components.get(componentName)
  }

  async executeAction(pluginId: string, actionName: string, context: any) {
    const plugin = this.plugins.get(pluginId)
    const action = plugin?.actions.get(actionName)
    
    if (!action) {
      console.warn(`Action ${actionName} not found in plugin ${pluginId}`)
      return
    }

    try {
      return await action(context)
    } catch (error) {
      console.error(`Failed to execute action ${actionName}:`, error)
      throw error
    }
  }

  getMountsBySlot(slot: string): Array<{ pluginId: string; mount: PluginMount }> {
    const mounts: Array<{ pluginId: string; mount: PluginMount }> = []
    
    for (const [pluginId, plugin] of this.plugins) {
      for (const mount of plugin.info.mounts) {
        if (mount.slot === slot) {
          mounts.push({ pluginId, mount })
        }
      }
    }

    // Sort by order
    mounts.sort((a, b) => (a.mount.order || 0) - (b.mount.order || 0))
    
    return mounts
  }

  getAllPlugins(): LoadedPlugin[] {
    return Array.from(this.plugins.values())
  }

  isLoading() {
    return this.loading.value
  }
}

// Singleton instance
export const pluginLoader = new PluginLoader()

// Auto-load plugins on app start
export async function initPluginSystem() {
  await pluginLoader.loadAll()
}
