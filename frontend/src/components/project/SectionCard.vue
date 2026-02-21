<template>
  <div class="section-card">
    <div class="section-header">
      <div class="section-title">
        <Icon :name="icon" size="md" />
        <h3>{{ title }}</h3>
        <span class="section-count">{{ count }}</span>
      </div>
      <div class="section-actions">
        <slot name="actions"></slot>
      </div>
    </div>
    <div class="section-content">
      <slot></slot>
    </div>
  </div>
</template>

<script setup lang="ts">
import Icon from '@/components/common/Icon.vue'

defineProps<{
  icon: string
  title: string
  count: number
}>()
</script>

<style scoped>
.section-card {
  background: rgba(255, 255, 255, 0.03);
  border-radius: 12px;
  border: 1px solid rgba(255, 255, 255, 0.06);
  overflow: hidden;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.section-title {
  display: flex;
  align-items: center;
  gap: 12px;
}

.section-title h3 {
  font-size: 16px;
  font-weight: 600;
  color: #e5e7eb;
  margin: 0;
}

.section-count {
  font-size: 13px;
  color: #9ca3af;
  background: rgba(255, 255, 255, 0.05);
  padding: 2px 8px;
  border-radius: 12px;
}

.section-actions {
  display: flex;
  gap: 8px;
}

.section-content {
  padding: 16px 20px;
}
</style>
