<template>
  <footer class="status-bar">
    <!-- 左侧：任务状态 + 项目信息 -->
    <div class="status-left">
      <!-- 任务状态（有任务时显示） -->
      <span v-if="currentTask" class="task-status" @click="toggleTaskPanel">
        [⚙️ {{ currentTask.message }}<template v-if="currentTask.progress !== undefined"> {{ currentTask.progress }}%</template>]
      </span>

      <!-- 项目信息 -->
      <span class="status-item">📁 {{ stats.projectName || '未选择项目' }}</span>
      <span class="status-item">📷 {{ formatNumber(stats.totalAssets) }}项</span>
      <span v-if="stats.storageSize" class="status-item">💾 {{ formatSize(stats.storageSize) }}</span>
    </div>

    <!-- 右侧：系统状态 + 版本 -->
    <div class="status-right">
      <button 
        class="status-btn"
        :class="wsStatusClass"
        :title="wsStatusText"
      >
        <span class="status-indicator">{{ wsConnected ? '🟢' : '🔴' }}</span>
      </button>
      
      <span class="status-item version">v{{ version }}</span>
    </div>
  </footer>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import Icon from '@/components/common/Icon.vue'
import { subscribeToEvent } from '@/services/api'

interface TaskStatus {
  id: string
  type: string
  message: string
  progress?: number
  running: boolean
  error?: string
}

interface Stats {
  projectName?: string
  totalAssets: number
  selectedCount: number
  storageSize?: number
}

const props = defineProps<{
  stats: Stats
  wsConnected: boolean
}>()

const emit = defineEmits<{
  'toggle-tasks': []
}>()

const currentTask = ref<TaskStatus | null>(null)
const pendingTasks = ref(0)
const version = '1.0.0'

const taskIcon = computed(() => {
  if (!currentTask.value) return 'check'
  if (currentTask.value.error) return 'alert-circle'
  if (currentTask.value.running) return 'loader'
  return 'check'
})

const wsStatusClass = computed(() => ({
  'status-connected': props.wsConnected,
  'status-disconnected': !props.wsConnected,
}))

const wsStatusText = computed(() => 
  props.wsConnected ? '已连接' : '未连接'
)

function toggleTaskPanel() {
  emit('toggle-tasks')
}

function formatNumber(num: number): string {
  return num.toLocaleString()
}

function formatSize(bytes: number): string {
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let size = bytes
  let unitIndex = 0
  
  while (size >= 1024 && unitIndex < units.length - 1) {
    size /= 1024
    unitIndex++
  }
  
  return `${size.toFixed(1)} ${units[unitIndex]}`
}

// WebSocket 事件监听
let unsubscribers: (() => void)[] = []

onMounted(() => {
  // 监听任务事件
  unsubscribers.push(
    subscribeToEvent('task:start', (data: any) => {
      currentTask.value = {
        id: data.id,
        type: data.type,
        message: data.message || '处理中...',
        running: true,
      }
    })
  )
  
  unsubscribers.push(
    subscribeToEvent('task:progress', (data: any) => {
      if (currentTask.value?.id === data.id) {
        currentTask.value.progress = data.progress
        currentTask.value.message = data.message || currentTask.value.message
      }
    })
  )
  
  unsubscribers.push(
    subscribeToEvent('task:complete', (data: any) => {
      if (currentTask.value?.id === data.id) {
        currentTask.value.running = false
        setTimeout(() => {
          currentTask.value = null
        }, 2000)
      }
    })
  )
  
  unsubscribers.push(
    subscribeToEvent('task:error', (data: any) => {
      if (currentTask.value?.id === data.id) {
        currentTask.value.running = false
        currentTask.value.error = data.error
      }
    })
  )
  
  unsubscribers.push(
    subscribeToEvent('task:queue', (data: any) => {
      pendingTasks.value = data.count || 0
    })
  )
})

onUnmounted(() => {
  unsubscribers.forEach(unsub => unsub())
})
</script>

<style scoped>
.status-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 28px;
  background: #1e1e1e;
  border-top: 1px solid rgba(255, 255, 255, 0.06);
  color: #9ca3af;
  font-size: 13px;
  padding: 0 12px;
  user-select: none;
}

.status-left {
  display: flex;
  align-items: center;
  gap: 16px;
  flex: 1;
}

.status-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.status-item {
  display: flex;
  align-items: center;
  gap: 6px;
  white-space: nowrap;
}

.status-item.version {
  color: #6b7280;
  font-size: 12px;
}

.status-btn {
  display: flex;
  align-items: center;
  padding: 0;
  background: transparent;
  border: none;
  cursor: pointer;
}

.status-indicator {
  font-size: 10px;
  line-height: 1;
}

.task-status {
  color: #60a5fa;
  cursor: pointer;
  white-space: nowrap;
}

.task-status:hover {
  color: #93c5fd;
}

.spinning {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
</style>
