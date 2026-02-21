<template>
  <div class="filter-content filter-content--rating">
    <div class="rating-options">
      <button 
        v-for="n in 6" 
        :key="6-n" 
        class="rating-btn"
        :class="{ active: modelValue === (6-n) }"
        @click="$emit('update:modelValue', modelValue === (6-n) ? -1 : (6-n))"
      >
        <span class="rating-stars">
          <span v-for="i in 5" :key="i" class="star" :class="{ filled: i <= (5-(6-n)+1) }">★</span>
        </span>
        <span class="rating-label">{{ getLabel(6-n) }}</span>
      </button>
    </div>
    <div class="filter-actions" v-if="modelValue >= 0">
      <button class="action-btn action-btn--clear" @click="$emit('update:modelValue', -1)">清除</button>
      <button class="action-btn action-btn--apply" @click="$emit('apply')">应用</button>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  modelValue: number
}>()

const emit = defineEmits<{
  'update:modelValue': [value: number]
  'apply': []
}>()

function getLabel(rating: number): string {
  const labels = ['未评分', '1星', '2星', '3星', '4星', '5星']
  return labels[rating] || '全部'
}
</script>

<style scoped>
.filter-content {
  padding: 12px;
}

.rating-options {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.rating-btn {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 12px;
  background: transparent;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.15s;
  width: 100%;
  text-align: left;
}

.rating-btn:hover {
  background: rgba(255, 255, 255, 0.05);
}

.rating-btn.active {
  background: rgba(59, 130, 246, 0.15);
}

.rating-stars {
  display: flex;
  gap: 2px;
}

.star {
  font-size: 14px;
  color: #3a3a3a;
}

.star.filled {
  color: #fbbf24;
}

.rating-label {
  font-size: 13px;
  color: #9ca3af;
}

.rating-btn.active .rating-label {
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
