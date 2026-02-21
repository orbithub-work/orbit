<template>
  <div class="chart-card">
    <div class="card-header">
      <h3 class="card-title">素材类型分布</h3>
    </div>
    <div class="distribution-content">
      <svg class="pie-chart" viewBox="0 0 200 200">
        <circle
          v-for="(segment, i) in segments"
          :key="i"
          cx="100"
          cy="100"
          r="70"
          fill="none"
          :stroke="segment.color"
          stroke-width="40"
          :stroke-dasharray="`${segment.length} ${circumference}`"
          :stroke-dashoffset="segment.offset"
          transform="rotate(-90 100 100)"
        />
      </svg>
      <div class="legend">
        <div v-for="item in distribution" :key="item.type" class="legend-item">
          <div class="legend-color" :style="{ background: item.color }"></div>
          <div class="legend-label">{{ item.label }}</div>
          <div class="legend-value">{{ item.count }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

defineProps<{
  timeRange: string
}>()

const distribution = [
  { type: 'image', label: '图片', count: 456, color: '#3b82f6' },
  { type: 'video', label: '视频', count: 234, color: '#8b5cf6' },
  { type: 'audio', label: '音频', count: 123, color: '#ec4899' },
  { type: 'other', label: '其他', count: 87, color: '#6b7280' },
]

const total = distribution.reduce((sum, item) => sum + item.count, 0)
const circumference = 2 * Math.PI * 70

const segments = computed(() => {
  let offset = 0
  return distribution.map(item => {
    const percentage = item.count / total
    const length = circumference * percentage
    const segment = {
      color: item.color,
      length,
      offset: -offset,
    }
    offset += length
    return segment
  })
})
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

.distribution-content {
  display: flex;
  align-items: center;
  gap: 32px;
}

.pie-chart {
  width: 200px;
  height: 200px;
  flex-shrink: 0;
}

.legend {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.legend-color {
  width: 12px;
  height: 12px;
  border-radius: 3px;
}

.legend-label {
  flex: 1;
  font-size: 13px;
  color: #9ca3af;
}

.legend-value {
  font-size: 14px;
  font-weight: 600;
  color: #e5e7eb;
}
</style>
