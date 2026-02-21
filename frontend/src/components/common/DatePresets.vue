<template>
  <div class="date-presets">
    <button
      v-for="preset in presets"
      :key="preset.value"
      :class="['preset-btn', { active: modelValue === preset.value }]"
      @click="$emit('update:modelValue', modelValue === preset.value ? '' : preset.value)"
    >
      {{ preset.label }}
    </button>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  modelValue?: string
}>()

defineEmits<{
  'update:modelValue': [value: string]
}>()

const presets = [
  { value: 'today', label: '今天' },
  { value: 'thisWeek', label: '本周' },
  { value: 'thisMonth', label: '本月' },
  { value: 'lastMonth', label: '上月' },
]
</script>

<style scoped>
.date-presets {
  display: flex;
  gap: 4px;
}

.preset-btn {
  padding: 4px 10px;
  background: transparent;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 4px;
  color: #9ca3af;
  font-size: 12px;
  line-height: 1;
  cursor: pointer;
  transition: all 0.15s;
}

.preset-btn:hover {
  border-color: rgba(255, 255, 255, 0.2);
  color: #e5e7eb;
}

.preset-btn.active {
  background: rgba(59, 130, 246, 0.15);
  border-color: #3b82f6;
  color: #60a5fa;
}
</style>
