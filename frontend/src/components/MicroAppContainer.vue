<template>
  <div
    v-if="plugin"
    class="micro-app-container"
  >
    <iframe
      ref="iframeRef"
      :src="plugin.ui?.entry"
      class="micro-app-iframe"
      @load="onIframeLoad"
    ></iframe>
    
    <!-- 遮罩层：当主程序繁忙或插件加载时显示 -->
    <div
      v-if="loading"
      class="iframe-loader"
    >
      <div class="spinner"></div>
      <span>正在连接领域插件...</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { usePluginStore } from '@/stores/pluginStore'

const props = defineProps<{
  pluginId: string
}>()

const pluginStore = usePluginStore()
const iframeRef = ref<HTMLIFrameElement | null>(null)
const loading = ref(true)

const plugin = computed(() => pluginStore.getPluginById(props.pluginId))

const onIframeLoad = () => {
  loading.value = false
  console.log(`插件 ${props.pluginId} 界面加载完成`)
  
  // 初始化通信：告诉插件主程序的 Token 和基本配置
  if (iframeRef.value?.contentWindow) {
    iframeRef.value.contentWindow.postMessage({
      type: 'INIT_HOST',
      payload: {
        host_version: '1.0.0',
        theme: 'dark' // 这里可以接入 themeStore
      }
    }, '*')
  }
}

// 监听来自插件的消息
const handleMessage = (event: MessageEvent) => {
    // 简单的安全检查，实际应用中应校验 origin
    if (!plugin.value?.ui?.entry.startsWith(event.origin)) {
        return
    }

    const { type, payload } = event.data
    console.log(`收到插件 ${props.pluginId} 消息:`, type, payload)

    // TODO: 处理具体的业务指令，如 PLUGIN_CMD_READ_METADATA
}

onMounted(() => {
    window.addEventListener('message', handleMessage)
})

onUnmounted(() => {
    window.removeEventListener('message', handleMessage)
})
</script>

<style scoped>
.micro-app-container {
  width: 100%;
  height: 100%;
  position: relative;
  background: var(--bg-main);
  display: flex;
  flex-direction: column;
}
.micro-app-iframe {
  flex: 1;
  width: 100%;
  border: none;
}
.iframe-loader {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.5);
  gap: 12px;
  color: white;
  z-index: 10;
}
.spinner {
  width: 32px;
  height: 32px;
  border: 3px solid rgba(255, 255, 255, 0.3);
  border-radius: 50%;
  border-top-color: #fff;
  animation: spin 1s ease-in-out infinite;
}
@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
