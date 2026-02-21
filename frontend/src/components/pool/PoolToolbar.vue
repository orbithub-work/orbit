<template>
  <div class="toolbar">
    <div class="toolbar-left">
      <button class="toolbar-btn" @click="$emit('refresh')">
        <Icon name="refresh" size="sm" />
      </button>
      <div class="toolbar-divider"></div>
      <span class="toolbar-text">{{ totalCount }} 项</span>
    </div>

    <div class="toolbar-center">
      <input
        type="text"
        class="search-input"
        placeholder="搜索素材..."
        :value="searchQuery"
        @input="$emit('update:searchQuery', ($event.target as HTMLInputElement).value)"
      />
    </div>

    <div class="toolbar-right">
      <button class="view-btn" :class="{ active: viewMode === 'list' }" @click="$emit('update:viewMode', 'list')">
        <Icon name="list" size="sm" />
      </button>
      <button class="view-btn" :class="{ active: viewMode === 'grid' }" @click="$emit('update:viewMode', 'grid')">
        <Icon name="grid" size="sm" />
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import Icon from '@/components/common/Icon.vue'

defineProps<{
  totalCount: number
  searchQuery: string
  viewMode: 'list' | 'grid'
}>()

defineEmits<{
  refresh: []
  'update:searchQuery': [value: string]
  'update:viewMode': [mode: 'list' | 'grid']
}>()
</script>

<style scoped>
.toolbar {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: rgba(255, 255, 255, 0.03);
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.toolbar-left,
.toolbar-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.toolbar-center {
  flex: 1;
}

.toolbar-btn,
.view-btn {
  padding: 6px 8px;
  background: transparent;
  border: none;
  color: #9ca3af;
  cursor: pointer;
  border-radius: 4px;
  transition: all 0.15s;
}

.toolbar-btn:hover,
.view-btn:hover {
  background: rgba(255, 255, 255, 0.06);
  color: #e5e7eb;
}

.view-btn.active {
  background: rgba(59, 130, 246, 0.15);
  color: #60a5fa;
}

.toolbar-divider {
  width: 1px;
  height: 16px;
  background: rgba(255, 255, 255, 0.1);
}

.toolbar-text {
  font-size: 13px;
  color: #9ca3af;
}

.search-input {
  width: 100%;
  padding: 6px 12px;
  background: rgba(0, 0, 0, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 6px;
  color: #e5e7eb;
  font-size: 13px;
}

.search-input:focus {
  outline: none;
  border-color: #3b82f6;
}
</style>
