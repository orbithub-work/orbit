<template>
  <div class="analytics-view">
    <header class="analytics-header">
      <h1 class="page-title">数据看板</h1>
      <div class="time-filter">
        <button
          v-for="option in timeOptions"
          :key="option.value"
          class="time-btn"
          :class="{ active: timeRange === option.value }"
          @click="timeRange = option.value"
        >
          {{ option.label }}
        </button>
      </div>
    </header>

    <div class="analytics-content">
      <MetricCards :time-range="timeRange" />
      <ProjectTrendChart :time-range="timeRange" />
      
      <div class="analytics-row">
        <TopAssetsCard :time-range="timeRange" />
        <AssetDistributionCard :time-range="timeRange" />
      </div>

      <InsightCards :time-range="timeRange" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import MetricCards from './MetricCards.vue'
import ProjectTrendChart from './ProjectTrendChart.vue'
import TopAssetsCard from './TopAssetsCard.vue'
import AssetDistributionCard from './AssetDistributionCard.vue'
import InsightCards from './InsightCards.vue'

const timeRange = ref<'week' | 'month' | 'quarter' | 'year'>('month')

const timeOptions = [
  { label: '本周', value: 'week' },
  { label: '本月', value: 'month' },
  { label: '本季度', value: 'quarter' },
  { label: '今年', value: 'year' },
]
</script>

<style scoped>
.analytics-view {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: #1b1c1f;
  overflow: hidden;
}

.analytics-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 24px 24px 16px;
  border-bottom: 1px solid #2b2b2f;
}

.page-title {
  font-size: 20px;
  font-weight: 600;
  color: #e5e7eb;
  margin: 0;
}

.time-filter {
  display: flex;
  gap: 4px;
  background: #252526;
  border-radius: 8px;
  padding: 4px;
}

.time-btn {
  padding: 6px 16px;
  background: transparent;
  border: none;
  border-radius: 6px;
  color: #9ca3af;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
}

.time-btn:hover {
  color: #e5e7eb;
}

.time-btn.active {
  background: #3b82f6;
  color: #fff;
}

.analytics-content {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.analytics-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 24px;
}
</style>
