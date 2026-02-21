<template>
  <div class="filter-content filter-content--time">
    <div class="filter-section">
      <div class="filter-section-title">快捷选择</div>
      <div class="quick-select">
        <button 
          v-for="option in quickOptions" 
          :key="option.id"
          class="quick-btn" 
          :class="{ active: modelValue === option.id }" 
          @click="$emit('update:modelValue', modelValue === option.id ? '' : option.id)"
        >
          {{ option.label }}
        </button>
      </div>
    </div>
    <div class="filter-section">
      <div class="filter-section-title">自定义范围</div>
      <div class="date-range">
        <input type="date" class="date-input" :value="startDate" @input="$emit('update:startDate', ($event.target as HTMLInputElement).value)" />
        <span class="date-sep">至</span>
        <input type="date" class="date-input" :value="endDate" @input="$emit('update:endDate', ($event.target as HTMLInputElement).value)" />
      </div>
    </div>
    <div class="filter-actions" v-if="modelValue || startDate || endDate">
      <button class="action-btn action-btn--clear" @click="handleClear">清除</button>
      <button class="action-btn action-btn--apply" @click="$emit('apply')">应用</button>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  modelValue: string
  startDate: string
  endDate: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
  'update:startDate': [value: string]
  'update:endDate': [value: string]
  'apply': []
  'clear': []
}>()

const quickOptions = [
  { id: 'today', label: '今天' },
  { id: 'week', label: '本周' },
  { id: 'month', label: '本月' },
  { id: 'year', label: '今年' },
]

function handleClear() {
  emit('update:modelValue', '')
  emit('update:startDate', '')
  emit('update:endDate', '')
  emit('clear')
}
</script>

<style scoped>
.filter-content {
  padding: 12px;
}

.filter-section {
  margin-bottom: 16px;
}

.filter-section:last-child {
  margin-bottom: 0;
}

.filter-section-title {
  font-size: 11px;
  color: #6b7280;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-bottom: 8px;
}

.quick-select {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.quick-btn {
  background: #2a2a2a;
  border: 1px solid #3a3a3a;
  border-radius: 6px;
  padding: 6px 12px;
  font-size: 13px;
  color: #9ca3af;
  cursor: pointer;
  transition: all 0.2s;
}

.quick-btn:hover {
  background: #333;
  color: #e5e7eb;
}

.quick-btn.active {
  background: rgba(59, 130, 246, 0.15);
  border-color: #3b82f6;
  color: #3b82f6;
}

.date-range {
  display: flex;
  align-items: center;
  gap: 8px;
}

.date-input {
  flex: 1;
  background: #2a2a2a;
  border: 1px solid #3a3a3a;
  border-radius: 6px;
  padding: 8px 10px;
  font-size: 13px;
  color: #e5e7eb;
  outline: none;
}

.date-input:focus {
  border-color: #3b82f6;
}

.date-sep {
  color: #6b7280;
  font-size: 13px;
}

.filter-actions {
  display: flex;
  gap: 8px;
  padding: 12px 0 0 0;
  border-top: 1px solid #333;
  margin-top: 12px;
}

.action-btn {
  flex: 1;
  padding: 8px 16px;
  border-radius: 6px;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
}

.action-btn--clear {
  background: transparent;
  border: 1px solid #3a3a3a;
  color: #9ca3af;
}

.action-btn--clear:hover {
  background: rgba(255, 255, 255, 0.05);
  color: #e5e7eb;
}

.action-btn--apply {
  background: #3b82f6;
  border: none;
  color: #fff;
}

.action-btn--apply:hover {
  background: #2563eb;
}
</style>
