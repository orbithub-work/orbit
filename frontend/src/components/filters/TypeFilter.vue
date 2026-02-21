<template>
  <div class="filter-content filter-content--type">
    <div class="filter-search-bar">
      <Icon name="search" size="sm" />
      <input type="text" class="filter-search-input" placeholder="搜索类型" v-model="searchQuery" />
    </div>
    <div class="filter-options">
      <label v-for="type in filteredTypes" :key="type.ext" class="filter-option">
        <input type="checkbox" :checked="modelValue.includes(type.ext)" @change="toggleType(type.ext)" />
        <Icon :name="getTypeIcon(type.ext)" size="sm" class="option-icon" />
        <span class="option-label">{{ type.label }}</span>
        <span class="option-count">{{ type.count }}</span>
      </label>
    </div>
    <div class="filter-actions" v-if="modelValue.length > 0">
      <button class="action-btn action-btn--clear" @click="handleClear">清除</button>
      <button class="action-btn action-btn--apply" @click="$emit('apply')">应用</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import Icon from '@/components/common/Icon.vue'

interface FileType {
  ext: string
  count: number
  label: string
}

const props = defineProps<{
  modelValue: string[]
  types: FileType[]
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string[]]
  'apply': []
  'clear': []
}>()

const searchQuery = ref('')

const filteredTypes = computed(() => {
  if (!searchQuery.value) return props.types
  const query = searchQuery.value.toLowerCase()
  return props.types.filter(t => t.ext.includes(query) || t.label.includes(query))
})

function getTypeIcon(ext: string): string {
  const icons: Record<string, string> = {
    jpg: 'image', jpeg: 'image', png: 'image', gif: 'image', webp: 'image', bmp: 'image', svg: 'sparkles',
    mp4: 'video', avi: 'video', mov: 'video', mkv: 'video', webm: 'video',
    mp3: 'audio', wav: 'audio', flac: 'audio',
    pdf: 'document', doc: 'document', docx: 'document', xls: 'chart-bar', xlsx: 'chart-bar',
    zip: 'archive', rar: 'archive', '7z': 'archive',
    psd: 'sparkles', ai: 'sparkles', sketch: 'sparkles',
  }
  return icons[ext.toLowerCase()] || 'file'
}

function toggleType(ext: string) {
  const newValue = [...props.modelValue]
  const idx = newValue.indexOf(ext)
  if (idx >= 0) {
    newValue.splice(idx, 1)
  } else {
    newValue.push(ext)
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

.filter-search-bar {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background: #1e1e1e;
  border-radius: 8px;
  margin-bottom: 8px;
}

.filter-search-bar svg {
  color: #6b7280;
}

.filter-search-input {
  flex: 1;
  background: transparent;
  border: none;
  font-size: 13px;
  color: #e5e7eb;
  outline: none;
}

.filter-search-input::placeholder {
  color: #6b7280;
}

.filter-options {
  display: flex;
  flex-direction: column;
  max-height: 240px;
  overflow-y: auto;
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

.option-icon {
  font-size: 16px;
}

.option-label {
  flex: 1;
  font-size: 13px;
  color: #e5e7eb;
}

.option-count {
  font-size: 12px;
  color: #6b7280;
  background: rgba(255, 255, 255, 0.05);
  padding: 2px 6px;
  border-radius: 4px;
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
