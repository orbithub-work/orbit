<template>
  <div class="filter-content filter-content--size">
    <div class="filter-section">
      <div class="filter-section-title">快捷尺寸</div>
      <div class="size-presets">
        <button 
          v-for="preset in presets" 
          :key="preset.id"
          class="preset-btn"
          :class="{ active: modelValue === preset.id }"
          @click="$emit('update:modelValue', modelValue === preset.id ? '' : preset.id)"
        >
          {{ preset.label }}
        </button>
      </div>
    </div>
    <div class="filter-section">
      <div class="filter-section-title">自定义尺寸</div>
      <div class="size-inputs">
        <div class="size-input-group">
          <label>宽度</label>
          <input type="number" class="size-input" :value="widthMin" @input="$emit('update:widthMin', Number(($event.target as HTMLInputElement).value) || null)" placeholder="最小" />
          <span>-</span>
          <input type="number" class="size-input" :value="widthMax" @input="$emit('update:widthMax', Number(($event.target as HTMLInputElement).value) || null)" placeholder="最大" />
        </div>
        <div class="size-input-group">
          <label>高度</label>
          <input type="number" class="size-input" :value="heightMin" @input="$emit('update:heightMin', Number(($event.target as HTMLInputElement).value) || null)" placeholder="最小" />
          <span>-</span>
          <input type="number" class="size-input" :value="heightMax" @input="$emit('update:heightMax', Number(($event.target as HTMLInputElement).value) || null)" placeholder="最大" />
        </div>
      </div>
    </div>
    <div class="filter-actions" v-if="modelValue || widthMin || widthMax || heightMin || heightMax">
      <button class="action-btn action-btn--clear" @click="handleClear">清除</button>
      <button class="action-btn action-btn--apply" @click="$emit('apply')">应用</button>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  modelValue: string
  widthMin: number | null
  widthMax: number | null
  heightMin: number | null
  heightMax: number | null
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
  'update:widthMin': [value: number | null]
  'update:widthMax': [value: number | null]
  'update:heightMin': [value: number | null]
  'update:heightMax': [value: number | null]
  'apply': []
  'clear': []
}>()

const presets = [
  { id: '4k', label: '4K (3840+)' },
  { id: '2k', label: '2K (2560+)' },
  { id: '1080p', label: '1080P (1920+)' },
  { id: '720p', label: '720P (1280+)' },
  { id: 'small', label: '小图 (<640)' },
]

function handleClear() {
  emit('update:modelValue', '')
  emit('update:widthMin', null)
  emit('update:widthMax', null)
  emit('update:heightMin', null)
  emit('update:heightMax', null)
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

.size-presets {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.preset-btn {
  background: #2a2a2a;
  border: 1px solid #3a3a3a;
  border-radius: 6px;
  padding: 6px 12px;
  font-size: 12px;
  color: #9ca3af;
  cursor: pointer;
  transition: all 0.2s;
}

.preset-btn:hover {
  background: #333;
  color: #e5e7eb;
}

.preset-btn.active {
  background: rgba(59, 130, 246, 0.15);
  border-color: #3b82f6;
  color: #3b82f6;
}

.size-inputs {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.size-input-group {
  display: flex;
  align-items: center;
  gap: 8px;
}

.size-input-group label {
  font-size: 12px;
  color: #6b7280;
  width: 40px;
}

.size-input {
  flex: 1;
  background: #2a2a2a;
  border: 1px solid #3a3a3a;
  border-radius: 6px;
  padding: 6px 8px;
  font-size: 12px;
  color: #e5e7eb;
  outline: none;
  width: 80px;
}

.size-input:focus {
  border-color: #3b82f6;
}

.size-input-group span {
  color: #6b7280;
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
