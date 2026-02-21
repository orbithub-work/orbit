<template>
  <iframe
    :src="iframeSrc"
    :width="width"
    :height="height"
    :style="iframeStyle"
    sandbox="allow-scripts allow-same-origin"
    class="plugin-iframe"
    @load="handleLoad"
  ></iframe>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { pluginLoader } from '@/core/pluginLoader'
import type { PluginMount } from '@/core/pluginLoader'

const props = defineProps<{
  pluginId: string
  mount: PluginMount
  context?: any
}>()

const loaded = ref(false)

const iframeSrc = computed(() => {
  const plugin = pluginLoader.getPlugin(props.pluginId)
  if (!plugin) return ''
  
  const baseUrl = plugin.info.endpoint || ''
  const entry = props.mount.entry || ''
  
  // Build URL with context as query params
  const url = new URL(entry, baseUrl)
  if (props.context) {
    url.searchParams.set('context', JSON.stringify(props.context))
  }
  
  return url.toString()
})

const width = computed(() => props.mount.width || '100%')
const height = computed(() => props.mount.height || 300)

const iframeStyle = computed(() => ({
  width: typeof width.value === 'number' ? `${width.value}px` : width.value,
  height: typeof height.value === 'number' ? `${height.value}px` : height.value,
  border: 'none',
  opacity: loaded.value ? 1 : 0,
  transition: 'opacity 0.2s'
}))

function handleLoad() {
  loaded.value = true
}
</script>

<style scoped>
.plugin-iframe {
  display: block;
  background: transparent;
}
</style>
