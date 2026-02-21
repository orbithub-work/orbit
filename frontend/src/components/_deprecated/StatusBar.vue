<template>
  <footer class="status-bar">
    <div class="status-left">
      <div class="status-item">
        <span class="status-label">索引文件</span>
        <span class="status-value">{{ stats.totalFiles }}</span>
      </div>
      <div class="status-divider"></div>
      <div class="status-item">
        <span class="status-dot online"></span>
        <span class="status-label">在线介质</span>
        <span class="status-value">{{ stats.onlineMedia }}</span>
      </div>
      <div class="status-divider"></div>
      <div class="status-item">
        <span class="status-dot offline"></span>
        <span class="status-label">离线介质</span>
        <span class="status-value">{{ stats.offlineMedia }}</span>
      </div>
    </div>

    <div class="status-right">
      <div class="task-wrapper">
        <div
          class="task-status"
          :class="{ active: activeCount > 0 }"
        >
          <span class="task-label">任务</span>
          <span class="task-count">{{ activeCount }}</span>
        </div>
        <div
          v-if="tasks.length > 0"
          class="task-panel"
        >
          <div
            v-for="task in displayTasks"
            :key="task.id"
            class="task-row"
          >
            <span class="task-type">{{ formatTaskType(task.taskType) }}</span>
            <span class="task-progress">{{ task.progress }}%</span>
            <span
              class="task-state"
              :class="task.status"
            >{{ formatStatus(task.status) }}</span>
          </div>
        </div>
      </div>
      <div
        class="sync-status"
        :class="{ syncing: isSyncing }"
      >
        <svg
          v-if="isSyncing"
          class="sync-icon spinning"
          viewBox="0 0 24 24"
          aria-hidden="true"
        >
          <use href="#i-sync" />
        </svg>
        <svg
          v-else
          class="sync-icon"
          viewBox="0 0 24 24"
          aria-hidden="true"
        >
          <use href="#i-check" />
        </svg>
        <span>{{ isSyncing ? '正在更新索引...' : '索引已同步' }}</span>
      </div>
    </div>
  </footer>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { apiCall, connectWebSocket, subscribeToEvent } from '@/services/api'

type TaskStatus = 'pending' | 'processing' | 'completed' | 'failed'

type MediaTask = {
  id: string
  assetId: string
  taskType: string
  status: TaskStatus
  progress: number
  updatedAt?: string
}

const stats = ref({
  totalFiles: 12847,
  onlineMedia: 3,
  offlineMedia: 1
})

const tasks = ref<MediaTask[]>([])
const taskTimers = new Map<string, number>()

const activeCount = computed(() => tasks.value.filter(task => task.status === 'processing').length)
const isSyncing = computed(() => activeCount.value > 0)
const displayTasks = computed(() => tasks.value.slice(0, 5))

const normalizeTask = (task: any): MediaTask => ({
  id: task.id,
  assetId: task.asset_id || task.assetId,
  taskType: task.task_type || task.taskType,
  status: task.status,
  progress: typeof task.progress === 'number' ? task.progress : 0,
  updatedAt: task.updated_at || task.updatedAt
})

const updateTask = (incoming: MediaTask) => {
  const index = tasks.value.findIndex(t => t.id === incoming.id)
  if (index >= 0) {
    tasks.value[index] = { ...tasks.value[index], ...incoming }
  } else {
    tasks.value.unshift(incoming)
  }
  tasks.value = tasks.value
    .slice()
    .sort((a, b) => (b.updatedAt || '').localeCompare(a.updatedAt || ''))
}

const scheduleRemove = (taskId: string) => {
  const existing = taskTimers.get(taskId)
  if (existing) {
    clearTimeout(existing)
  }
  const timer = window.setTimeout(() => {
    tasks.value = tasks.value.filter(task => task.id !== taskId)
    taskTimers.delete(taskId)
  }, 5000)
  taskTimers.set(taskId, timer)
}

const formatTaskType = (taskType: string) => {
  const map: Record<string, string> = {
    metadata: '元数据',
    thumbnail: '缩略图'
  }
  return map[taskType] || taskType
}

const formatStatus = (status: TaskStatus) => {
  const map: Record<TaskStatus, string> = {
    pending: '等待',
    processing: '处理中',
    completed: '完成',
    failed: '失败'
  }
  return map[status]
}

const loadActiveTasks = async () => {
  const data = await apiCall<any[]>('get_active_tasks')
  tasks.value = (data || []).map(normalizeTask)
}

onMounted(async () => {
  try {
    // await loadActiveTasks()
    // connectWebSocket()
    // subscribeToEvent('task_claimed', (data: any) => updateTask(normalizeTask(data)))
    // subscribeToEvent('task_progress', (data: any) => updateTask(normalizeTask(data)))
    // subscribeToEvent('task_completed', (data: any) => {
    //   const task = normalizeTask(data)
    //   updateTask(task)
    //   scheduleRemove(task.id)
    // })
    // subscribeToEvent('task_failed', (data: any) => {
    //   const task = normalizeTask(data)
    //   updateTask(task)
    //   scheduleRemove(task.id)
    // })
    console.log('StatusBar: WebSocket & Tasks DISABLED for debugging')
  } catch (e) {
    console.error('StatusBar error:', e)
  }
})

onUnmounted(() => {
  taskTimers.forEach(timer => clearTimeout(timer))
  taskTimers.clear()
})
</script>

<style scoped>
.status-bar {
  height: 28px;
  background: var(--color-bg-base);
  border-top: 1px solid var(--color-border);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 var(--spacing-lg);
  font-size: 11px;
  user-select: none;
}

.status-left {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
}

.status-item {
  display: flex;
  align-items: center;
  gap: 6px;
}

.status-label {
  color: var(--color-text-tertiary);
}

.status-value {
  color: var(--color-text-secondary);
  font-weight: 500;
  font-variant-numeric: tabular-nums;
}

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
}

.status-dot.online {
  background: var(--color-success);
  box-shadow: 0 0 6px rgba(16, 185, 129, 0.3);
}

.status-dot.offline {
  background: var(--color-text-disabled);
}

.status-divider {
  width: 1px;
  height: 10px;
  background: var(--color-border);
}

.status-right {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
}

.sync-status {
  display: flex;
  align-items: center;
  gap: 6px;
  color: var(--color-text-secondary);
  font-weight: 500;
}

.sync-status.syncing {
  color: var(--color-primary);
}

.sync-icon {
  width: 12px;
  height: 12px;
}

.sync-icon.spinning {
  animation: spin 1s linear infinite;
}

.task-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.task-status {
  display: flex;
  align-items: center;
  gap: 6px;
  color: var(--color-text-secondary);
  font-weight: 500;
  padding: 2px 6px;
  border-radius: 6px;
  background: var(--color-hover);
}

.task-status.active {
  color: var(--color-primary);
  background: var(--color-active-bg);
}

.task-panel {
  position: absolute;
  right: 0;
  bottom: 32px;
  min-width: 220px;
  background: var(--color-bg-surface);
  border: 1px solid var(--color-border);
  border-radius: 8px;
  padding: 8px;
  display: none;
  box-shadow: var(--shadow-lg);
  z-index: 20;
}

.task-wrapper:hover .task-panel {
  display: block;
}

.task-row {
  display: grid;
  grid-template-columns: 1fr auto auto;
  gap: 8px;
  padding: 4px 0;
  align-items: center;
  font-size: 11px;
  color: var(--color-text-secondary);
}

.task-type {
  color: var(--color-text-primary);
  font-weight: 500;
}

.task-progress {
  font-variant-numeric: tabular-nums;
}

.task-state {
  font-variant-numeric: tabular-nums;
}

.task-state.processing {
  color: var(--color-primary);
}

.task-state.completed {
  color: var(--color-success);
}

.task-state.failed {
  color: var(--color-danger);
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}
</style>
