<template>
  <div class="filter-bar">
    <button
      v-for="filter in filters"
      :key="filter.key"
      class="filter-btn"
      :class="{ active: activeFilter === filter.key }"
      @click="$emit('update:activeFilter', filter.key)"
    >
      {{ filter.label }}
      <span v-if="filter.count !== undefined" class="filter-count">{{ filter.count }}</span>
    </button>
  </div>
</template>

<script setup lang="ts">
interface Filter {
  key: string
  label: string
  count?: number
}

defineProps<{
  filters: Filter[]
  activeFilter: string
}>()

defineEmits<{
  'update:activeFilter': [key: string]
}>()
</script>

<style scoped>
.filter-bar {
  display: flex;
  gap: 8px;
  padding: 12px 16px;
  background: rgba(255, 255, 255, 0.03);
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  overflow-x: auto;
}

.filter-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: transparent;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 6px;
  color: #9ca3af;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.15s;
  white-space: nowrap;
}

.filter-btn:hover {
  background: rgba(255, 255, 255, 0.06);
  color: #e5e7eb;
}

.filter-btn.active {
  background: rgba(59, 130, 246, 0.15);
  border-color: #3b82f6;
  color: #60a5fa;
}

.filter-count {
  font-size: 11px;
  padding: 2px 6px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 10px;
}

.filter-btn.active .filter-count {
  background: rgba(59, 130, 246, 0.2);
}
</style>
