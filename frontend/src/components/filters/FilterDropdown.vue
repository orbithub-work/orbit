<template>
  <Teleport to="body">
    <Transition name="dropdown">
      <div 
        v-if="visible" 
        class="filter-dropdown-overlay"
        @click="$emit('close')"
      >
        <div 
          class="dropdown-content bg-base-200 border border-base-300 rounded-box shadow-xl"
          :style="dropdownStyle"
          @click.stop
        >
          <slot></slot>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  visible: boolean
  targetRect: DOMRect | null
  width?: string
}>()

const emit = defineEmits<{
  close: []
}>()

const dropdownStyle = computed(() => {
  if (!props.targetRect) return {}
  
  const viewportWidth = window.innerWidth
  const viewportHeight = window.innerHeight
  
  let left = props.targetRect.left
  let top = props.targetRect.bottom + 8
  
  if (left + 320 > viewportWidth) {
    left = viewportWidth - 320 - 16
  }
  
  return {
    position: 'fixed' as const,
    left: `${left}px`,
    top: `${top}px`,
    width: props.width || 'auto',
    minWidth: props.width ? 'auto' : '200px',
    maxHeight: `${Math.min(400, viewportHeight - top - 16)}px`,
  }
})
</script>

<style scoped>
.filter-dropdown-overlay {
  position: fixed;
  inset: 0;
  z-index: 9999;
}

.dropdown-content {
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.dropdown-enter-active,
.dropdown-leave-active {
  transition: all 0.2s ease;
}

.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
}

.dropdown-enter-from .dropdown-content,
.dropdown-leave-to .dropdown-content {
  transform: translateY(-8px) scale(0.95);
}
</style>
