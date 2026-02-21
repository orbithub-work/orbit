<template>
  <nav
    class="breadcrumb"
    aria-label="面包屑导航"
  >
    <button 
      class="breadcrumb-item breadcrumb-item--home"
      :title="'根目录'"
      @click="handleNavigate('/')"
    >
      <svg
        viewBox="0 0 24 24"
        xmlns="http://www.w3.org/2000/svg"
      >
        <path d="M10 20v-6h4v6h5v-8h3L12 3 2 12h3v8z" />
      </svg>
    </button>
    
    <span class="breadcrumb-separator">/</span>
    
    <template
      v-for="(segment, index) in pathSegments"
      :key="index"
    >
      <button 
        class="breadcrumb-item"
        :title="segment"
        @click="handleNavigate(getPathAtIndex(index))"
      >
        {{ segment }}
      </button>
      
      <span
        v-if="index < pathSegments.length - 1"
        class="breadcrumb-separator"
      >/</span>
    </template>
  </nav>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  currentPath: string
  rootLabel?: string
}

const props = withDefaults(defineProps<Props>(), {
  rootLabel: '根目录'
})

const emit = defineEmits<{
  navigate: [path: string]
}>()

const pathSegments = computed(() => {
  const path = props.currentPath.replace(/^[/\\]+|[/\\]+$/g, '')
  if (!path) return []
  
  const separator = props.currentPath.includes('\\') ? '\\' : '/'
  return path.split(separator).filter(Boolean)
})

const getPathAtIndex = (index: number): string => {
  const separator = props.currentPath.includes('\\') ? '\\' : '/'
  return pathSegments.value.slice(0, index + 1).join(separator)
}

const handleNavigate = (path: string) => {
  emit('navigate', path)
}
</script>

<style scoped>
.breadcrumb {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  padding: 0.75rem 1rem;
  background-color: var(--color-surface);
  border-bottom: 1px solid var(--color-border);
  overflow-x: auto;
  white-space: nowrap;
}

.breadcrumb-item {
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
  padding: 0.375rem 0.75rem;
  background: none;
  border: 1px solid transparent;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.875rem;
  color: var(--color-text);
  transition: all 0.15s;
  white-space: nowrap;
}

.breadcrumb-item:hover {
  background-color: var(--color-hover);
  color: var(--color-primary);
}

.breadcrumb-item--home svg {
  width: 18px;
  height: 18px;
  fill: currentColor;
}

.breadcrumb-item--home {
  padding: 0.375rem;
}

.breadcrumb-separator {
  color: var(--color-text-secondary);
  font-size: 0.875rem;
}

.breadcrumb-item:last-child {
  font-weight: 500;
  color: var(--color-text);
}

.breadcrumb-item:last-child:hover {
  background-color: transparent;
  color: var(--color-text);
}
</style>