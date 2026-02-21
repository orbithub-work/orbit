<template>
  <div class="chart-card">
    <div class="card-header">
      <h3 class="card-title">项目完成趋势</h3>
    </div>
    <div class="chart-container">
      <svg class="trend-chart" viewBox="0 0 800 300">
        <defs>
          <linearGradient id="lineGradient" x1="0%" y1="0%" x2="0%" y2="100%">
            <stop offset="0%" style="stop-color:#3b82f6;stop-opacity:0.3" />
            <stop offset="100%" style="stop-color:#3b82f6;stop-opacity:0" />
          </linearGradient>
        </defs>
        
        <!-- Grid lines -->
        <line v-for="i in 5" :key="i" 
          :x1="50" :y1="50 + (i - 1) * 50" 
          :x2="750" :y2="50 + (i - 1) * 50" 
          stroke="#2b2b2f" stroke-width="1" />
        
        <!-- Area -->
        <path :d="areaPath" fill="url(#lineGradient)" />
        
        <!-- Line -->
        <path :d="linePath" fill="none" stroke="#3b82f6" stroke-width="2" />
        
        <!-- Points -->
        <circle v-for="(point, i) in points" :key="i"
          :cx="point.x" :cy="point.y" r="4"
          fill="#3b82f6" />
        
        <!-- Labels -->
        <text v-for="(label, i) in labels" :key="i"
          :x="50 + i * 100" y="280"
          fill="#6b7280" font-size="12" text-anchor="middle">
          {{ label }}
        </text>
      </svg>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

defineProps<{
  timeRange: string
}>()

// Mock data
const data = [3, 5, 4, 8, 6, 9, 12]
const labels = ['周一', '周二', '周三', '周四', '周五', '周六', '周日']

const maxValue = Math.max(...data)

const points = computed(() => 
  data.map((value, i) => ({
    x: 50 + i * 100,
    y: 250 - (value / maxValue) * 200
  }))
)

const linePath = computed(() => {
  const path = points.value.map((p, i) => 
    `${i === 0 ? 'M' : 'L'} ${p.x} ${p.y}`
  ).join(' ')
  return path
})

const areaPath = computed(() => {
  const line = points.value.map((p, i) => 
    `${i === 0 ? 'M' : 'L'} ${p.x} ${p.y}`
  ).join(' ')
  return `${line} L 750 250 L 50 250 Z`
})
</script>

<style scoped>
.chart-card {
  background: #252526;
  border-radius: 12px;
  padding: 20px;
}

.card-header {
  margin-bottom: 20px;
}

.card-title {
  font-size: 15px;
  font-weight: 600;
  color: #e5e7eb;
  margin: 0;
}

.chart-container {
  width: 100%;
  height: 300px;
}

.trend-chart {
  width: 100%;
  height: 100%;
}
</style>
