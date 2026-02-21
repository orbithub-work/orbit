<template>
  <div class="filter-content filter-content--filesize">
    <div class="filesize-options">
      <label v-for="range in ranges" :key="range.id" class="filter-option">
        <input type="checkbox" :checked="modelValue.includes(range.id)" @change="toggleRange(range.id)" />
        <span class="option-label">{{ range.label }}</span>
      </label>
    </div>
    <div class="filter-actions" v-if="modelValue.length > 0">
      <button class="action-btn action-btn--clear" @click="handleClear">清除</button>
      <button class="action-btn action-btn--apply" @click="$emit('apply')">应用</button>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  modelValue: string[]
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string[]]
  'apply': []
  'clear': []
}>()

const ranges = [
  { id: 'tiny', label: '< 100KB' },
  { id: 'small', label: '100KB - 1MB' },
  { id: 'medium', label: '1MB - 10MB' },
  { id: 'large', label: '10MB - 100MB' },
  { id: 'huge', label: '> 100MB' },
]

function toggleRange(id: string) {
  const newValue = [...props.modelValue]
  const idx = newValue.indexOf(id)
  if (idx >= 0) {
    newValue.splice(idx, 1)
  } else {
    newValue.push(id)
  }
  emit('update:modelValue', newValue)
}

function handleClear() {
  emit('update:modelValue', [])
  emit('clear')
}
</script>

<style scoped>
.filter-content {
  padding: 12px;
}

.filesize-options {
  display: flex;
  flex-direction: column;
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

.filter-option input[type="checkbox"] {
  width: 16px;
  height: 16px;
  accent-color: #3b82f6;
}

.option-label {
  flex: 1;
  font-size: 13px;
  color: #e5e7eb;
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
