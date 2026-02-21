<template>
  <div class="skeleton-grid">
    <div 
      v-for="i in count" 
      :key="i" 
      class="skeleton-item"
      :style="{ width: `${itemSize}px`, height: `${itemSize}px` }"
    >
      <div class="skeleton-thumbnail"></div>
      <div class="skeleton-text"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
withDefaults(defineProps<{
  count?: number
  itemSize?: number
}>(), {
  count: 12,
  itemSize: 160
})
</script>

<style scoped>
.skeleton-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
  gap: 16px;
  padding: 16px;
}

.skeleton-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.skeleton-thumbnail {
  width: 100%;
  height: 100%;
  background: linear-gradient(
    90deg,
    rgba(255, 255, 255, 0.05) 25%,
    rgba(255, 255, 255, 0.08) 50%,
    rgba(255, 255, 255, 0.05) 75%
  );
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
  border-radius: 8px;
}

.skeleton-text {
  width: 80%;
  height: 12px;
  background: linear-gradient(
    90deg,
    rgba(255, 255, 255, 0.05) 25%,
    rgba(255, 255, 255, 0.08) 50%,
    rgba(255, 255, 255, 0.05) 75%
  );
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
  border-radius: 4px;
}

@keyframes shimmer {
  0% {
    background-position: 200% 0;
  }
  100% {
    background-position: -200% 0;
  }
}
</style>
