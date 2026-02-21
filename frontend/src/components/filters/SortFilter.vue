<template>
  <div class="filter-content filter-content--sort">
    <div class="sort-options">
      <label v-for="option in options" :key="option.id" class="filter-option">
        <input type="radio" name="sort" :value="option.id" :checked="modelValue === option.id" @change="$emit('update:modelValue', option.id)" />
        <span class="option-label">{{ option.label }}</span>
      </label>
    </div>
    <div class="sort-direction">
      <button class="direction-btn" :class="{ active: direction === 'asc' }" @click="$emit('update:direction', 'asc')">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polyline points="18 15 12 9 6 15"></polyline>
        </svg>
        升序
      </button>
      <button class="direction-btn" :class="{ active: direction === 'desc' }" @click="$emit('update:direction', 'desc')">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polyline points="6 9 12 15 18 9"></polyline>
        </svg>
        降序
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  modelValue: string
  direction: 'asc' | 'desc'
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
  'update:direction': [value: 'asc' | 'desc']
}>()

const options = [
  { id: 'name', label: '文件名' },
  { id: 'date', label: '修改时间' },
  { id: 'size', label: '文件大小' },
  { id: 'type', label: '文件类型' },
]
</script>

<style scoped>
.filter-content {
  padding: 12px;
}

.sort-options {
  display: flex;
  flex-direction: column;
  margin-bottom: 12px;
}

.filter-option {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 12px;
  cursor: pointer;
  border-radius: 6px;
  transition: all 0.15s;
}

.filter-option:hover {
  background: rgba(255, 255, 255, 0.05);
}

.filter-option input[type="radio"] {
  width: 16px;
  height: 16px;
  accent-color: #3b82f6;
}

.option-label {
  flex: 1;
  font-size: 13px;
  color: #e5e7eb;
}

.sort-direction {
  display: flex;
  gap: 8px;
  padding-top: 12px;
  border-top: 1px solid #333;
}

.direction-btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 8px;
  background: #2a2a2a;
  border: 1px solid #3a3a3a;
  border-radius: 6px;
  font-size: 12px;
  color: #9ca3af;
  cursor: pointer;
  transition: all 0.2s;
}

.direction-btn:hover {
  background: #333;
  color: #e5e7eb;
}

.direction-btn.active {
  background: rgba(59, 130, 246, 0.15);
  border-color: #3b82f6;
  color: #3b82f6;
}
</style>
