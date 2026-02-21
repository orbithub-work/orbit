<template>
  <div class="tree-node" :style="{ paddingLeft: `${depth * 16}px` }">
    <div
      class="tree-node-content"
      :class="{ 
        active: isActive,
        'has-children': hasChildren 
      }"
      @click="handleClick"
    >
      <button
        v-if="hasChildren"
        class="tree-toggle"
        @click.stop="toggleExpand"
      >
        <Icon :name="isExpanded ? 'chevron-down' : 'chevron-right'" size="sm" />
      </button>
      <div v-else class="tree-spacer"></div>
      
      <Icon :name="icon" size="sm" class="tree-icon" />
      <span class="tree-label">{{ label }}</span>
      <span v-if="hasWarning" class="tree-warning" title="监控权限不足">⚠️</span>
      <span v-if="count !== undefined" class="tree-count">{{ count }}</span>
    </div>
    
    <div v-if="hasChildren && isExpanded" class="tree-children">
      <slot></slot>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import Icon from '@/components/common/Icon.vue'

const props = defineProps<{
  label: string
  icon?: string
  count?: number
  depth?: number
  isActive?: boolean
  hasChildren?: boolean
  defaultExpanded?: boolean
  hasWarning?: boolean
}>()

const emit = defineEmits<{
  click: []
}>()

const isExpanded = ref(props.defaultExpanded ?? true)

function toggleExpand() {
  isExpanded.value = !isExpanded.value
}

function handleClick() {
  emit('click')
}
</script>

<style scoped>
.tree-node {
  user-select: none;
}

.tree-node-content {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 8px;
  border-radius: 6px;
  cursor: pointer;
  transition: background 0.15s;
}

.tree-node-content:hover {
  background: rgba(255, 255, 255, 0.06);
}

.tree-node-content.active {
  background: rgba(59, 130, 246, 0.15);
  color: #60a5fa;
}

.tree-toggle {
  width: 16px;
  height: 16px;
  padding: 0;
  background: transparent;
  border: none;
  color: #9ca3af;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 3px;
  transition: all 0.15s;
  flex-shrink: 0;
}

.tree-toggle:hover {
  background: rgba(255, 255, 255, 0.08);
  color: #e5e7eb;
}

.tree-spacer {
  width: 16px;
  flex-shrink: 0;
}

.tree-icon {
  color: #9ca3af;
  flex-shrink: 0;
}

.tree-node-content.active .tree-icon {
  color: #60a5fa;
}

.tree-warning {
  font-size: 12px;
  margin-left: 4px;
  opacity: 0.8;
}

.tree-label {
  flex: 1;
  font-size: 13px;
  color: #e5e7eb;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.tree-count {
  font-size: 11px;
  color: #6b7280;
  background: rgba(255, 255, 255, 0.06);
  padding: 2px 6px;
  border-radius: 10px;
  flex-shrink: 0;
}

.tree-node-content.active .tree-count {
  background: rgba(59, 130, 246, 0.2);
  color: #60a5fa;
}

.tree-children {
  margin-left: 0;
}
</style>
