<template>
  <div class="quick-filters">
    <button
      v-for="filter in filters"
      :key="filter.value"
      :class="['filter-btn', { active: modelValue === filter.value }]"
      @click="$emit('update:modelValue', modelValue === filter.value ? '' : filter.value)"
    >
      <Icon :name="filter.icon" size="sm" />
      <span>{{ filter.label }}</span>
    </button>
  </div>
</template>

<script setup lang="ts">
import Icon from './Icon.vue'

defineProps<{
  modelValue?: string
}>()

defineEmits<{
  'update:modelValue': [value: string]
}>()

const filters = [
  { value: 'recent', label: '最近', icon: 'clock' },
  { value: 'unrated', label: '未评分', icon: 'star' },
  { value: 'large', label: '大文件', icon: 'archive-box' },
  { value: 'vertical', label: '竖屏', icon: 'arrow-up' },
  { value: 'horizontal', label: '横屏', icon: 'arrow-right' },
]
</script>

<style scoped>
.quick-filters {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.filter-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 6px;
  color: #9ca3af;
  font-size: 13px;
  line-height: 1;
  cursor: pointer;
  transition: all 0.15s;
}

.filter-btn:hover {
  background: rgba(255, 255, 255, 0.08);
  border-color: rgba(255, 255, 255, 0.2);
  color: #e5e7eb;
}

.filter-btn.active {
  background: rgba(59, 130, 246, 0.15);
  border-color: #3b82f6;
  color: #60a5fa;
}
</style>
