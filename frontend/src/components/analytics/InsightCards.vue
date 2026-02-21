<template>
  <div class="insight-cards">
    <div v-for="insight in insights" :key="insight.id" class="insight-card">
      <div class="insight-icon">{{ insight.icon }}</div>
      <div class="insight-content">
        <div class="insight-title">{{ insight.title }}</div>
        <div class="insight-description">{{ insight.description }}</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

defineProps<{
  timeRange: string
}>()

const insights = computed(() => [
  {
    id: '1',
    icon: '🚀',
    title: '效率提升',
    description: '你的项目平均周期 5.2 天，比上月快 15%，继续保持！',
  },
  {
    id: '2',
    icon: '💎',
    title: '高价值素材',
    description: 'intro_music.mp3 使用率最高，建议收藏类似风格的素材',
  },
  {
    id: '3',
    icon: '🧹',
    title: '素材优化',
    description: '有 87 个素材超过 3 个月未使用，考虑清理以释放空间',
  },
])
</script>

<style scoped>
.insight-cards {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
}

.insight-card {
  background: linear-gradient(135deg, #252526 0%, #2a2a2a 100%);
  border-radius: 12px;
  padding: 20px;
  display: flex;
  gap: 16px;
  border: 1px solid rgba(59, 130, 246, 0.1);
}

.insight-icon {
  font-size: 32px;
  flex-shrink: 0;
}

.insight-content {
  flex: 1;
  min-width: 0;
}

.insight-title {
  font-size: 14px;
  font-weight: 600;
  color: #e5e7eb;
  margin-bottom: 6px;
}

.insight-description {
  font-size: 12px;
  color: #9ca3af;
  line-height: 1.5;
}
</style>
