<template>
  <Transition name="slide-up">
    <div v-if="visible" class="task-panel">
      <div class="panel-header">
        <div class="header-left">
          <Icon name="list" size="sm" />
          <span class="panel-title">任务</span>
          <span class="task-count">{{ tasks.length }}</span>
        </div>
        <div class="header-actions">
          <button class="header-btn" title="清空已完成" @click="clearCompleted">
            <Icon name="trash" size="sm" />
          </button>
          <button class="header-btn" @click="$emit('close')">
            <Icon name="x" size="sm" />
          </button>
        </div>
      </div>

      <div class="panel-content">
        <div v-if="tasks.length === 0" class="empty-state">
          <Icon name="check-circle" size="lg" />
          <p>暂无任务</p>
        </div>

        <div v-else class="task-list">
          <div
            v-for="task in tasks"
            :key="task.id"
            class="task-item"
            :class="{ 
              running: task.status === 'running',
              completed: task.status === 'completed',
              failed: task.status === 'failed'
            }"
          >
            <div class="task-icon">
              <Icon v-if="task.status === 'running'" name="loader" size="sm" class="spinning" />
              <Icon v-else-if="task.status === 'completed'" name="check-circle" size="sm" />
              <Icon v-else-if="task.status === 'failed'" name="alert-circle" size="sm" />
              <Icon v-else name="clock" size="sm" />
            </div>

            <div class="task-info">
              <div class="task-header">
                <span class="task-name">{{ task.name }}</span>
                <span class="task-time">{{ formatTime(task.startTime) }}</span>
              </div>
              <div class="task-message">{{ task.message }}</div>
              
              <div v-if="task.status === 'running' && task.progress !== undefined" class="task-progress">
                <div class="progress-bar">
                  <div class="progress-fill" :style="{ width: `${task.progress}%` }"></div>
                </div>
                <span class="progress-text">{{ task.progress }}%</span>
              </div>

              <div v-if="task.error" class="task-error">
                {{ task.error }}
              </div>
            </div>

            <div class="task-actions">
              <button
                v-if="task.status === 'running'"
                class="task-action-btn"
                title="取消"
                @click="cancelTask(task.id)"
              >
                <Icon name="x" size="xs" />
              </button>
              <button
                v-if="task.status === 'failed'"
                class="task-action-btn"
                title="重试"
                @click="retryTask(task.id)"
              >
                <Icon name="refresh-cw" size="xs" />
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import Icon from '@/components/common/Icon.vue'
import { subscribeToEvent } from '@/services/api'

interface Task {
  id: string
  name: string
  message: string
  status: 'pending' | 'running' | 'completed' | 'failed'
  progress?: number
  startTime: number
  endTime?: number
  error?: string
}

defineProps<{
  visible: boolean
}>()

const emit = defineEmits<{
  close: []
}>()

const tasks = ref<Task[]>([])

function formatTime(timestamp: number): string {
  const now = Date.now()
  const diff = now - timestamp
  
  if (diff < 60000) return '刚刚'
  if (diff < 3600000) return `${Math.floor(diff / 60000)} 分钟前`
  if (diff < 86400000) return `${Math.floor(diff / 3600000)} 小时前`
  
  const date = new Date(timestamp)
  return `${date.getHours()}:${String(date.getMinutes()).padStart(2, '0')}`
}

function clearCompleted() {
  tasks.value = tasks.value.filter(t => t.status !== 'completed')
}

function cancelTask(id: string) {
  // TODO: 调用 API 取消任务
  console.log('Cancel task:', id)
}

function retryTask(id: string) {
  // TODO: 调用 API 重试任务
  console.log('Retry task:', id)
}

// WebSocket 事件监听
let unsubscribers: (() => void)[] = []

onMounted(() => {
  unsubscribers.push(
    subscribeToEvent('task:start', (data: any) => {
      const existing = tasks.value.find(t => t.id === data.id)
      if (existing) {
        existing.status = 'running'
        existing.message = data.message || '处理中...'
        existing.startTime = Date.now()
      } else {
        tasks.value.unshift({
          id: data.id,
          name: data.name || getTaskTypeName(data.type),
          message: data.message || '处理中...',
          status: 'running',
          startTime: Date.now(),
        })
      }
    })
  )
  
  unsubscribers.push(
    subscribeToEvent('task:progress', (data: any) => {
      const task = tasks.value.find(t => t.id === data.id)
      if (task) {
        task.progress = data.progress
        if (data.message) task.message = data.message
      }
    })
  )
  
  unsubscribers.push(
    subscribeToEvent('task:complete', (data: any) => {
      const task = tasks.value.find(t => t.id === data.id)
      if (task) {
        task.status = 'completed'
        task.endTime = Date.now()
        task.message = data.message || '完成'
      }
    })
  )
  
  unsubscribers.push(
    subscribeToEvent('task:error', (data: any) => {
      const task = tasks.value.find(t => t.id === data.id)
      if (task) {
        task.status = 'failed'
        task.endTime = Date.now()
        task.error = data.error
      }
    })
  )
})

onUnmounted(() => {
  unsubscribers.forEach(unsub => unsub())
})

function getTaskTypeName(type: string): string {
  const names: Record<string, string> = {
    scan: '扫描文件',
    import: '导入素材',
    thumbnail: '生成缩略图',
    export: '导出文件',
    analyze: '分析素材',
  }
  return names[type] || '任务'
}
</script>

<style scoped>
.task-panel {
  position: fixed;
  bottom: 24px;
  left: 0;
  right: 0;
  height: 300px;
  background: #1e1e1e;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  display: flex;
  flex-direction: column;
  z-index: 100;
}

.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  background: #252526;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.panel-title {
  font-size: 13px;
  font-weight: 500;
  color: #e5e7eb;
}

.task-count {
  font-size: 11px;
  color: #9ca3af;
  background: rgba(255, 255, 255, 0.1);
  padding: 2px 6px;
  border-radius: 10px;
}

.header-actions {
  display: flex;
  gap: 4px;
}

.header-btn {
  width: 24px;
  height: 24px;
  border: none;
  background: transparent;
  color: #9ca3af;
  cursor: pointer;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s;
}

.header-btn:hover {
  background: rgba(255, 255, 255, 0.08);
  color: #e5e7eb;
}

.panel-content {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #6b7280;
  gap: 8px;
}

.empty-state p {
  font-size: 13px;
}

.task-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.task-item {
  display: flex;
  gap: 12px;
  padding: 12px;
  background: #252526;
  border-radius: 6px;
  border-left: 3px solid transparent;
  transition: all 0.15s;
}

.task-item.running {
  border-left-color: #3b82f6;
}

.task-item.completed {
  border-left-color: #10b981;
  opacity: 0.7;
}

.task-item.failed {
  border-left-color: #ef4444;
}

.task-icon {
  flex-shrink: 0;
  color: #9ca3af;
}

.task-item.running .task-icon {
  color: #3b82f6;
}

.task-item.completed .task-icon {
  color: #10b981;
}

.task-item.failed .task-icon {
  color: #ef4444;
}

.task-info {
  flex: 1;
  min-width: 0;
}

.task-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
}

.task-name {
  font-size: 13px;
  font-weight: 500;
  color: #e5e7eb;
}

.task-time {
  font-size: 11px;
  color: #6b7280;
}

.task-message {
  font-size: 12px;
  color: #9ca3af;
  margin-bottom: 8px;
}

.task-progress {
  display: flex;
  align-items: center;
  gap: 8px;
}

.progress-bar {
  flex: 1;
  height: 4px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 2px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: #3b82f6;
  transition: width 0.3s;
}

.progress-text {
  font-size: 11px;
  color: #3b82f6;
  font-weight: 500;
  min-width: 35px;
  text-align: right;
}

.task-error {
  font-size: 12px;
  color: #ef4444;
  margin-top: 4px;
}

.task-actions {
  display: flex;
  gap: 4px;
}

.task-action-btn {
  width: 24px;
  height: 24px;
  border: none;
  background: transparent;
  color: #9ca3af;
  cursor: pointer;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s;
}

.task-action-btn:hover {
  background: rgba(255, 255, 255, 0.08);
  color: #e5e7eb;
}

.spinning {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

/* 过渡动画 */
.slide-up-enter-active,
.slide-up-leave-active {
  transition: transform 0.2s ease-out;
}

.slide-up-enter-from,
.slide-up-leave-to {
  transform: translateY(100%);
}
</style>
