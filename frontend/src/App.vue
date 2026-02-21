<template>
  <router-view />
</template>

<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { connectWebSocket, subscribeToEvent } from '@/services/api'

const router = useRouter()
let unsubscribeSettings: (() => void) | undefined

onMounted(() => {
  const isTrayMenuWindow = window.location.hash.includes('/tray-menu')
  if (isTrayMenuWindow) return
  
  // 异步连接 WebSocket，不阻塞渲染
  connectWebSocket().catch(err => {
    console.debug('[App] WebSocket connection failed:', err)
  })

  subscribeToEvent('SHOW_MAIN_WINDOW', () => {
    const trayHandler = (window as any)?.mediaAssistant?.tray?.sendAction
    if (typeof trayHandler === 'function') {
      trayHandler('show-main')
    }
  })

  const mediaAssistant = (window as any)?.mediaAssistant
  if (mediaAssistant?.on) {
    unsubscribeSettings = mediaAssistant.on('open-settings', () => {
      router.push('/settings')
    })
  }
})

onUnmounted(() => {
  if (unsubscribeSettings) {
    unsubscribeSettings()
  }
})
</script>

<style>
/* Global styles */
:root {
  --bg-main: #1e1e1e;
  --text-main: #e5e7eb;
}

body {
  margin: 0;
  padding: 0;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
  background-color: var(--bg-main);
  color: var(--text-main);
  overflow: hidden;
}

html, body, #app {
  width: 100%;
  height: 100%;
}

* {
  box-sizing: border-box;
}
</style>
