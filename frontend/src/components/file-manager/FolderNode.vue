<template>
  <div class="folder-node">
    <div 
      :class="['folder-item', { active: isSelected }]"
      :style="{ paddingLeft: `${12 + level * 16}px` }"
      @click="handleSelect"
    >
      <!-- 展开/折叠图标 -->
      <span 
        v-if="hasChildren"
        class="toggle-icon"
        @click.stop="handleToggle"
      >
        <svg 
          viewBox="0 0 24 24" 
          :class="{ expanded: isExpanded }"
          class="arrow-icon"
        >
          <path d="M8.59 16.59L13.17 12 8.59 7.41 10 6l6 6-6 6-1.41-1.41z" />
        </svg>
      </span>
      <span
        v-else
        class="toggle-icon placeholder"
      ></span>
      
      <!-- 文件夹图标 -->
      <span class="folder-icon">
        <svg
          viewBox="0 0 24 24"
          fill="currentColor"
        >
          <path d="M10 4H4c-1.1 0-1.99.9-1.99 2L2 18c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V8c0-1.1-.9-2-2-2h-8l-2-2z" />
        </svg>
      </span>
      
      <!-- 文件夹名称 -->
      <span class="folder-name">{{ node.name }}</span>
      
      <!-- 文件数量 -->
      <span
        v-if="node.children?.length"
        class="folder-count"
      >
        {{ node.children.length }}
      </span>
    </div>
    
    <!-- 子文件夹 -->
    <div
      v-if="hasChildren && isExpanded"
      class="folder-children"
    >
      <FolderNode
        v-for="child in node.children"
        :key="child.id"
        :node="child"
        :level="level + 1"
        :selected-path="selectedPath"
        @select="$emit('select', $event)"
        @toggle="$emit('toggle', $event)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

interface FolderNodeData {
  id: string
  name: string
  path: string
  children?: FolderNodeData[]
  isExpanded?: boolean
}

interface Props {
  node: FolderNodeData
  level?: number
  selectedPath?: string
}

const props = withDefaults(defineProps<Props>(), {
  level: 0,
  selectedPath: ''
})

const emit = defineEmits<{
  select: [node: FolderNodeData]
  toggle: [node: FolderNodeData]
}>()

const isExpanded = ref(props.node.isExpanded ?? false)

const hasChildren = computed(() => {
  return props.node.children && props.node.children.length > 0
})

const isSelected = computed(() => {
  return props.selectedPath === props.node.path
})

const handleSelect = () => {
  emit('select', props.node)
}

const handleToggle = () => {
  isExpanded.value = !isExpanded.value
  emit('toggle', props.node)
}
</script>

<style scoped>
.folder-node {
  user-select: none;
}

.folder-item {
  display: flex;
  align-items: center;
  padding: 6px 8px;
  margin: 2px 0;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
  color: #b0b0b0;
  font-size: 13px;
}

.folder-item:hover {
  background: rgba(255, 255, 255, 0.05);
  color: #fff;
}

.folder-item.active {
  background: var(--color-active-bg);
  color: var(--color-primary);
}

.toggle-icon {
  width: 16px;
  height: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 4px;
  cursor: pointer;
  border-radius: 3px;
}

.toggle-icon:hover {
  background: rgba(255, 255, 255, 0.1);
}

.toggle-icon.placeholder {
  cursor: default;
}

.toggle-icon.placeholder:hover {
  background: transparent;
}

.arrow-icon {
  width: 16px;
  height: 16px;
  fill: currentColor;
  opacity: 0.5;
  transition: transform 0.2s;
}

.arrow-icon.expanded {
  transform: rotate(90deg);
}

.folder-icon {
  width: 18px;
  height: 18px;
  margin-right: 8px;
  color: var(--color-primary);
  opacity: 0.8;
}

.folder-item.active .folder-icon {
  opacity: 1;
}

.folder-icon svg {
  width: 100%;
  height: 100%;
}

.folder-name {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.folder-count {
  font-size: 10px;
  color: #666;
  background: rgba(255, 255, 255, 0.05);
  padding: 1px 6px;
  border-radius: 8px;
  margin-left: 4px;
}

.folder-item.active .folder-count {
  color: var(--color-primary);
  background: var(--color-active-bg);
}

.folder-children {
  animation: slideDown 0.2s ease;
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-4px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
