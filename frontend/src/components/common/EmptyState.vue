<template>
  <div class="empty-state" :class="{ 'empty-state--compact': compact }">
    <div class="empty-icon">{{ icon }}</div>
    <div class="empty-title">{{ title }}</div>
    <div v-if="description" class="empty-description">{{ description }}</div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  icon?: string
  title: string
  description?: string
  compact?: boolean
}>()
</script>

<style scoped>
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 24px;
  color: #6b7280;
  user-select: none;
}

.empty-state--compact {
  padding: 32px 16px;
}

.empty-icon {
  font-size: 64px;
  margin-bottom: 16px;
  opacity: 0.4;
  filter: grayscale(0.3);
}

.empty-state--compact .empty-icon {
  font-size: 40px;
  margin-bottom: 12px;
}

.empty-title {
  font-size: 15px;
  font-weight: 500;
  color: #9ca3af;
  margin-bottom: 8px;
}

.empty-state--compact .empty-title {
  font-size: 13px;
  margin-bottom: 0;
}

.empty-description {
  font-size: 13px;
  color: #6b7280;
  text-align: center;
  max-width: 320px;
  line-height: 1.5;
}

.empty-state--compact .empty-description {
  font-size: 12px;
  max-width: 200px;
}
</style>
