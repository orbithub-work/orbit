<template>
  <section class="content-section">
    <div class="section-header">
      <h3 class="section-title">
        <Icon name="activity" size="md" />
        项目日志
      </h3>
    </div>
    <div class="section-body">
      <EmptyState
        v-if="logs.length === 0"
        icon="📝"
        title="暂无日志"
        description="系统会自动记录项目操作"
        compact
      />
      <div v-else class="log-timeline">
        <div v-for="log in logs" :key="log.id" class="log-item">
          <div class="log-dot"></div>
          <div class="log-content">
            <div class="log-action">{{ log.action }}</div>
            <div class="log-time">{{ formatTime(log.time) }}</div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import Icon from '@/components/common/Icon.vue'
import EmptyState from '@/components/common/EmptyState.vue'

interface Log {
  id: string
  action: string
  time: string
}

const props = defineProps<{
  projectId?: string | number
}>()

// Mock 数据 - 后续接入真实 API
const logs = computed<Log[]>(() => {
  if (!props.projectId) return []
  return [
    { id: '1', action: '创建项目', time: new Date().toISOString() },
    { id: '2', action: '添加素材 3 个', time: new Date(Date.now() - 3600000).toISOString() },
    { id: '3', action: '更新项目状态为「进行中」', time: new Date(Date.now() - 7200000).toISOString() },
  ]
})

function formatTime(time: string) {
  const date = new Date(time)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  
  if (diff < 60000) return '刚刚'
  if (diff < 3600000) return `${Math.floor(diff / 60000)} 分钟前`
  if (diff < 86400000) return `${Math.floor(diff / 3600000)} 小时前`
  return date.toLocaleDateString('zh-CN', { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })
}
</script>

<style scoped>
.content-section {
  background: #252526;
  border-radius: 12px;
  overflow: hidden;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.section-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 15px;
  font-weight: 600;
  color: #e5e7eb;
  margin: 0;
}

.section-body {
  padding: 20px;
}

.log-timeline {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.log-item {
  display: flex;
  gap: 12px;
  position: relative;
}

.log-item:not(:last-child)::before {
  content: '';
  position: absolute;
  left: 5px;
  top: 20px;
  bottom: -16px;
  width: 1px;
  background: rgba(255, 255, 255, 0.1);
}

.log-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: #3b82f6;
  margin-top: 4px;
  flex-shrink: 0;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.2);
}

.log-content {
  flex: 1;
  min-width: 0;
}

.log-action {
  font-size: 13px;
  color: #e5e7eb;
  margin-bottom: 2px;
}

.log-time {
  font-size: 11px;
  color: #6b7280;
}
</style>
