<template>
  <div class="filter-content filter-content--shape">
    <div class="shape-grid">
      <button 
        v-for="shape in shapes" 
        :key="shape.id" 
        class="shape-btn" 
        :class="{ active: modelValue.includes(shape.id) }"
        @click="toggleShape(shape.id)"
      >
        <div class="shape-preview" :style="shape.style"></div>
        <span>{{ shape.label }}</span>
      </button>
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

const shapes = [
  { id: 'square', label: '正方形', style: { aspectRatio: '1' } },
  { id: 'portrait', label: '竖向', style: { aspectRatio: '3/4' } },
  { id: 'landscape', label: '横向', style: { aspectRatio: '4/3' } },
  { id: 'ultra-wide', label: '超宽', style: { aspectRatio: '21/9' } },
  { id: 'ultra-tall', label: '超高', style: { aspectRatio: '9/21' } },
]

function toggleShape(id: string) {
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

.shape-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 8px;
}

.shape-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  padding: 12px;
  background: #2a2a2a;
  border: 1px solid #3a3a3a;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.shape-btn:hover {
  background: #333;
}

.shape-btn.active {
  background: rgba(59, 130, 246, 0.15);
  border-color: #3b82f6;
}

.shape-preview {
  width: 40px;
  height: 40px;
  background: linear-gradient(135deg, #4a5568, #2d3748);
  border-radius: 4px;
}

.shape-btn span {
  font-size: 11px;
  color: #9ca3af;
}

.shape-btn.active span {
  color: #3b82f6;
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
