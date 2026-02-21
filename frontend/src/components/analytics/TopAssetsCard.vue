<template>
  <div class="chart-card">
    <div class="card-header">
      <h3 class="card-title">高频素材 TOP 10</h3>
    </div>
    <div class="asset-list">
      <EmptyState
        v-if="topAssets.length === 0"
        icon="📊"
        title="暂无数据"
        description="完成项目后将显示高频使用的素材"
        compact
      />
      <div v-for="asset in topAssets" :key="asset.id" class="asset-item">
        <div class="asset-thumb">
          <span class="file-icon">{{ asset.icon }}</span>
        </div>
        <div class="asset-info">
          <div class="asset-name">{{ asset.name }}</div>
          <div class="asset-meta">使用 {{ asset.count }} 次</div>
        </div>
        <div class="asset-bar">
          <div class="bar-fill" :style="{ width: `${asset.percentage}%` }"></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import EmptyState from '@/components/common/EmptyState.vue'

defineProps<{
  timeRange: string
}>()

// Mock data
const topAssets = computed(() => [
  { id: '1', name: 'intro_music.mp3', icon: '🎵', count: 15, percentage: 100 },
  { id: '2', name: 'logo_v2.png', icon: '🖼️', count: 12, percentage: 80 },
  { id: '3', name: 'transition_01.mp4', icon: '🎬', count: 10, percentage: 67 },
  { id: '4', name: 'bg_texture.jpg', icon: '🖼️', count: 8, percentage: 53 },
  { id: '5', name: 'outro_template.aep', icon: '📄', count: 7, percentage: 47 },
])
</script>

<style scoped>
.chart-card {
  background: #252526;
  border-radius: 12px;
  padding: 20px;
}

.card-header {
  margin-bottom: 16px;
}

.card-title {
  font-size: 15px;
  font-weight: 600;
  color: #e5e7eb;
  margin: 0;
}

.asset-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.asset-item {
  display: grid;
  grid-template-columns: 40px 1fr 120px;
  gap: 12px;
  align-items: center;
  padding: 8px;
  border-radius: 8px;
  transition: background 0.2s;
}

.asset-item:hover {
  background: rgba(255, 255, 255, 0.05);
}

.asset-thumb {
  width: 40px;
  height: 40px;
  background: #1e1e1e;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.file-icon {
  font-size: 20px;
}

.asset-info {
  min-width: 0;
}

.asset-name {
  font-size: 13px;
  color: #e5e7eb;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.asset-meta {
  font-size: 11px;
  color: #6b7280;
  margin-top: 2px;
}

.asset-bar {
  height: 6px;
  background: #1e1e1e;
  border-radius: 3px;
  overflow: hidden;
}

.bar-fill {
  height: 100%;
  background: linear-gradient(90deg, #3b82f6, #60a5fa);
  border-radius: 3px;
  transition: width 0.3s;
}
</style>
