<template>
  <div
    class="skeleton"
    :class="[`skeleton--${variant}`, { 'skeleton--animated': animated }]"
  ></div>
</template>

<script setup lang="ts">
interface Props {
  variant?: 'text' | 'circular' | 'rectangular' | 'rounded'
  animated?: boolean
}

withDefaults(defineProps<Props>(), {
  variant: 'text',
  animated: true
})
</script>

<style scoped>
.skeleton {
  background: linear-gradient(
    135deg,
    rgba(255, 255, 255, 0.03) 0%,
    rgba(255, 255, 255, 0.06) 50%,
    rgba(255, 255, 255, 0.03) 100%
  );
  position: relative;
  overflow: hidden;
}

.skeleton--animated::after {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(
    90deg,
    transparent 0%,
    rgba(255, 255, 255, 0.08) 50%,
    transparent 100%
  );
  animation: shimmer 1.8s ease-in-out infinite;
}

@keyframes shimmer {
  0% {
    left: -100%;
  }
  100% {
    left: 100%;
  }
}

.skeleton--text {
  border-radius: 6px;
  height: 1em;
}

.skeleton--circular {
  border-radius: 50%;
}

.skeleton--rectangular {
  border-radius: 8px;
}

.skeleton--rounded {
  border-radius: 12px;
}
</style>
