<template>
  <div v-if="show && history.length > 0" class="search-history">
    <div class="history-header">
      <span>最近搜索</span>
      <button class="clear-btn" @click="handleClear">清空</button>
    </div>
    <div class="history-list">
      <button
        v-for="item in history"
        :key="item.id"
        class="history-item"
        @click="$emit('select', item.query)"
      >
        <Icon name="clock" size="sm" />
        <span class="query">{{ item.query }}</span>
        <span class="count">{{ item.count }}次</span>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import Icon from './Icon.vue'
import { getSearchHistory, clearSearchHistory, type SearchHistoryItem } from '@/services/api'

const props = defineProps<{
  show: boolean
}>()

const emit = defineEmits<{
  select: [query: string]
  close: []
}>()

const history = ref<SearchHistoryItem[]>([])

watch(() => props.show, async (show) => {
  if (show) {
    try {
      history.value = await getSearchHistory(10)
    } catch (e) {
      console.debug('Failed to load search history:', e)
    }
  }
})

const handleClear = async () => {
  try {
    await clearSearchHistory()
    history.value = []
    emit('close')
  } catch (e) {
    console.error('Failed to clear search history:', e)
  }
}
</script>

<style scoped>
.search-history {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  margin-top: 4px;
  background: #2d2d2d;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 6px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
  z-index: 100;
  max-height: 300px;
  overflow-y: auto;
}

.history-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  font-size: 12px;
  color: #9ca3af;
}

.clear-btn {
  background: none;
  border: none;
  color: #60a5fa;
  font-size: 12px;
  cursor: pointer;
  padding: 0;
  line-height: 1;
}

.clear-btn:hover {
  color: #93c5fd;
}

.history-list {
  padding: 4px;
}

.history-item {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
  padding: 8px;
  background: transparent;
  border: none;
  border-radius: 4px;
  color: #e5e7eb;
  font-size: 13px;
  line-height: 1;
  text-align: left;
  cursor: pointer;
  transition: background 0.15s;
}

.history-item:hover {
  background: rgba(255, 255, 255, 0.05);
}

.query {
  flex: 1;
}

.count {
  font-size: 11px;
  color: #6b7280;
}
</style>
