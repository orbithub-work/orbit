<template>
  <div class="folder-tree">
    <div
      v-if="loading"
      class="folder-tree-loading"
    >
      <span class="loading-spinner"></span>
      <span>加载中...</span>
    </div>
    
    <div
      v-else-if="tree.length === 0"
      class="folder-tree-empty"
    >
      <p>没有文件夹</p>
    </div>
    
    <div
      v-else
      class="tree-root"
    >
      <FolderNode
        v-for="node in tree"
        :key="node.id"
        :node="node"
        :level="0"
        @select="handleSelect"
        @toggle="handleToggle"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import FolderNode from './FolderNode.vue'

export interface FolderNodeData {
  id: string
  name: string
  path: string
  children?: FolderNodeData[]
  isExpanded?: boolean
}

interface Props {
  tree: FolderNodeData[]
  loading?: boolean
}

withDefaults(defineProps<Props>(), {
  loading: false
})

const emit = defineEmits<{
  select: [node: FolderNodeData]
  toggle: [node: FolderNodeData]
}>()

const handleSelect = (node: FolderNodeData) => {
  emit('select', node)
}

const handleToggle = (node: FolderNodeData) => {
  emit('toggle', node)
}
</script>

<style scoped>
.folder-tree {
  font-size: 13px;
}

.folder-tree-loading,
.folder-tree-empty {
  padding: 20px;
  text-align: center;
  color: #666;
  font-size: 12px;
}

.loading-spinner {
  display: inline-block;
  width: 16px;
  height: 16px;
  border: 2px solid #333;
  border-top-color: #667eea;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin-right: 8px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.tree-root {
  padding: 4px 0;
}
</style>
