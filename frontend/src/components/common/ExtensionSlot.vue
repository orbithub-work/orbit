<template>
  <div class="extension-slot">
    <component
      v-for="item in mountedPlugins"
      :key="`${item.pluginId}-${item.mount.entry}`"
      :is="getComponent(item)"
      :plugin-id="item.pluginId"
      :mount="item.mount"
      :context="context"
      @action="handleAction"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, defineAsyncComponent } from 'vue'
import { pluginLoader } from '@/core/pluginLoader'
import type { PluginMount } from '@/core/pluginLoader'

const props = defineProps<{
  slot: string
  context?: any
}>()

const emit = defineEmits<{
  action: [pluginId: string, actionName: string, data: any]
}>()

// Get all plugins mounted at this slot
const mountedPlugins = computed(() => {
  return pluginLoader.getMountsBySlot(props.slot)
})

// Get component for a plugin mount
function getComponent(item: { pluginId: string; mount: PluginMount }) {
  const plugin = pluginLoader.getPlugin(item.pluginId)
  
  if (!plugin) {
    return null
  }

  // If plugin is frontend mode, get the component from loaded module
  if (plugin.info.mode === 'frontend') {
    // Extract component name from entry (e.g., "CopyrightStatus" from mount.entry)
    const componentName = item.mount.entry
    const component = plugin.components.get(componentName)
    
    if (component) {
      return component
    }
  }

  // If plugin is network service or local process, use iframe
  if (plugin.info.mode === 'network_service' || plugin.info.mode === 'local_process') {
    return defineAsyncComponent(() => import('./PluginIframe.vue'))
  }

  return null
}

async function handleAction(pluginId: string, actionName: string, data: any) {
  try {
    await pluginLoader.executeAction(pluginId, actionName, {
      ...props.context,
      ...data
    })
    emit('action', pluginId, actionName, data)
  } catch (error) {
    console.error('Plugin action failed:', error)
  }
}
</script>

<style scoped>
.extension-slot {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
</style>
