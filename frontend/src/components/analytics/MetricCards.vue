<template>
  <div class="metric-cards">
    <div class="metric-card">
      <div class="metric-icon">📋</div>
      <div class="metric-content">
        <div class="metric-label">完成项目</div>
        <div class="metric-value">{{ stats.completedProjects }}</div>
        <div class="metric-change" :class="{ positive: stats.projectChange > 0 }">
          {{ stats.projectChange > 0 ? '↑' : '↓' }} {{ Math.abs(stats.projectChange) }}%
        </div>
      </div>
    </div>

    <div class="metric-card">
      <div class="metric-icon">📁</div>
      <div class="metric-content">
        <div class="metric-label">新增素材</div>
        <div class="metric-value">{{ stats.newAssets }}</div>
        <div class="metric-change" :class="{ positive: stats.assetChange > 0 }">
          {{ stats.assetChange > 0 ? '↑' : '↓' }} {{ Math.abs(stats.assetChange) }}%
        </div>
      </div>
    </div>

    <div class="metric-card">
      <div class="metric-icon">⏱️</div>
      <div class="metric-content">
        <div class="metric-label">平均周期</div>
        <div class="metric-value">{{ stats.avgCycle }}天</div>
        <div class="metric-change" :class="{ positive: stats.cycleChange < 0 }">
          {{ stats.cycleChange > 0 ? '↑' : '↓' }} {{ Math.abs(stats.cycleChange) }}%
        </div>
      </div>
    </div>

    <div class="metric-card">
      <div class="metric-icon">🎯</div>
      <div class="metric-content">
        <div class="metric-label">素材使用率</div>
        <div class="metric-value">{{ stats.assetUsageRate }}%</div>
        <div class="metric-change" :class="{ positive: stats.usageChange > 0 }">
          {{ stats.usageChange > 0 ? '↑' : '↓' }} {{ Math.abs(stats.usageChange) }}%
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  timeRange: string
}>()

// Mock data - 后续接入真实数据
const stats = computed(() => ({
  completedProjects: 12,
  projectChange: 20,
  newAssets: 234,
  assetChange: -10,
  avgCycle: 5.2,
  cycleChange: -15,
  assetUsageRate: 68,
  usageChange: 8,
}))
</script>

<style scoped>
.metric-cards {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.metric-card {
  background: #252526;
  border-radius: 12px;
  padding: 20px;
  display: flex;
  gap: 16px;
  align-items: flex-start;
}

.metric-icon {
  font-size: 32px;
  opacity: 0.8;
}

.metric-content {
  flex: 1;
}

.metric-label {
  font-size: 12px;
  color: #9ca3af;
  margin-bottom: 8px;
}

.metric-value {
  font-size: 24px;
  font-weight: 600;
  color: #e5e7eb;
  margin-bottom: 4px;
}

.metric-change {
  font-size: 12px;
  color: #ef4444;
}

.metric-change.positive {
  color: #22c55e;
}
</style>
