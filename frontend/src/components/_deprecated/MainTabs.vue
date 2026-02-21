<template>
  <div class="main-tabs">
    <div class="tabs-left">
      <div class="tab-list">
        <button
          v-for="tab in tabs"
          :key="tab.key"
          :class="['tab-item', { active: activeTab === tab.key }]"
          @click="switchTab(tab.key)"
        >
          {{ tab.label }}
        </button>
      </div>
    </div>

    <div class="tabs-right">
      <button
        class="add-folder-btn"
        @click="handleAdd"
      >
        <!-- SVG Removed -->
        Add
      </button>
      <button
        class="filter-btn"
        title="筛选条件"
        @click="handleFilter"
      >
        <!-- SVG Removed -->
        Filter
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const tabs = [
  { key: 'material', label: '按素材' },
  { key: 'project', label: '按项目' }
]

const activeTab = ref('material')

const emit = defineEmits<{
  change: [tab: string]
  add: []
  filter: []
}>()

const switchTab = (tab: string) => {
  activeTab.value = tab
  emit('change', tab)
}

const handleAdd = () => {
  emit('add')
}

const handleFilter = () => {
  emit('filter')
}
</script>

<style scoped>
.main-tabs {
  height: var(--tab-height);
  background: transparent;
  border-bottom: 1px solid var(--glass-border);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 var(--spacing-lg);
  flex-shrink: 0;
  position: relative;
  z-index: 50;
}

.tabs-left {
  display: flex;
  align-items: center;
  height: 100%;
}

.tab-list {
  display: flex;
  gap: var(--spacing-lg);
  height: 100%;
}

.tab-item {
  position: relative;
  height: 100%;
  padding: 0 var(--spacing-xs);
  background: transparent;
  border: none;
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
}

.tab-item:hover {
  color: var(--color-text-primary);
}

.tab-item.active {
  color: var(--color-text-primary);
}

.tab-item.active::after {
  content: '';
  position: absolute;
  bottom: -1px;
  left: 0;
  right: 0;
  height: 2px;
  background: var(--color-primary);
  border-radius: 2px 2px 0 0;
}

.tabs-right {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.add-folder-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: var(--color-primary);
  border: none;
  border-radius: var(--radius-md);
  font-size: 13px;
  font-weight: 500;
  color: white;
  cursor: pointer;
  transition: all 0.2s;
}

.add-folder-btn:hover {
  background: var(--color-primary-hover);
}

.add-folder-btn svg {
  width: var(--icon-size-16);
  height: var(--icon-size-16);
}

.filter-btn {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: all 0.2s;
}

.filter-btn:hover {
  background: var(--color-bg-surface-hover);
  color: var(--color-text-primary);
  border-color: var(--color-text-secondary);
}

.filter-btn svg {
  width: var(--icon-size-16);
  height: var(--icon-size-16);
}
</style>
