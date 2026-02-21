<template>
  <div class="dashboard-view">
    <h1>仪表盘</h1>
    <p>欢迎使用智归档OS！</p>
    
    <div class="stats-grid">
      <div class="stat-card">
        <h3>系统状态</h3>
        <p class="stat-value" :class="{ 'online': isConnected, 'offline': !isConnected }">
          {{ isConnected ? '在线' : '离线' }}
        </p>
        <p v-if="systemInfo" class="stat-sub">
          {{ systemInfo.os }} / {{ systemInfo.arch }}
        </p>
      </div>
      <div class="stat-card">
        <h3>文件总数</h3>
        <p class="stat-value">
          0
        </p>
      </div>
      <div class="stat-card">
        <h3>已归档</h3>
        <p class="stat-value">
          0
        </p>
      </div>
      <div class="stat-card">
        <h3>收藏夹</h3>
        <p class="stat-value">
          0
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { apiCall } from '../services/api'

const isConnected = ref(false)
const systemInfo = ref<any>(null)

onMounted(async () => {
  try {
    const info = await apiCall('get_system_info')
    systemInfo.value = info
    isConnected.value = true
  } catch (e) {
    console.error('Failed to get system info:', e)
    isConnected.value = false
  }
})
</script>

<style scoped>
.dashboard-view {
  padding: 1rem;
  flex: 1;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
  margin-top: 2rem;
}

.stat-card {
  background-color: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 8px;
  padding: 1.5rem;
  text-align: center;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.stat-value {
  font-size: 2rem;
  font-weight: bold;
  color: var(--color-primary);
  margin: 0.5rem 0 0;
}

.stat-value.online {
  color: #10b981;
}

.stat-value.offline {
  color: #ef4444;
}

.stat-sub {
  font-size: 0.875rem;
  color: var(--color-text-secondary);
  margin-top: 0.5rem;
}
</style>