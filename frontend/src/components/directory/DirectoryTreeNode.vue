<template>
  <div class="tree-node-wrapper">
    <TreeNode
      :label="node.name"
      :icon="node.is_root ? 'folder' : 'folder'"
      :depth="depth"
      :is-active="isActive"
      :has-children="node.has_children"
      :default-expanded="false"
      @click="handleClick"
    >
      <!-- 懒加载子节点 -->
      <template v-if="isExpanded && node.has_children">
        <div v-if="loading" class="loading-placeholder">
          <span class="loading-text">加载中...</span>
        </div>
        <DirectoryTreeNode
          v-for="child in children"
          :key="child.path"
          :node="child"
          :active-path="activePath"
          :depth="depth + 1"
          @select="$emit('select', $event)"
          @load-children="$emit('load-children', $event)"
        />
      </template>
    </TreeNode>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import TreeNode from '@/components/common/TreeNode.vue'
import { getDirectoryChildren, type DirectoryNode } from '@/services/api'

const props = defineProps<{
  node: DirectoryNode
  activePath?: string
  depth: number
}>()

const emit = defineEmits<{
  select: [path: string]
  'load-children': [node: DirectoryNode, children: DirectoryNode[]]
}>()

const isExpanded = ref(false)
const children = ref<DirectoryNode[]>([])
const loading = ref(false)

const isActive = computed(() => props.activePath === props.node.path)

async function handleClick() {
  emit('select', props.node.path)
  
  // 如果有子节点且未加载，则加载
  if (props.node.has_children && !isExpanded.value && children.value.length === 0) {
    await loadChildren()
  }
  
  isExpanded.value = !isExpanded.value
}

async function loadChildren() {
  if (loading.value) return
  
  loading.value = true
  try {
    const result = await getDirectoryChildren(props.node.path)
    children.value = result
    emit('load-children', props.node, result)
  } catch (error) {
    console.error('Failed to load children:', error)
  } finally {
    loading.value = false
  }
}

// 监听展开状态变化
watch(isExpanded, (expanded) => {
  if (expanded && props.node.has_children && children.value.length === 0) {
    loadChildren()
  }
})
</script>

<style scoped>
.tree-node-wrapper {
  user-select: none;
}

.loading-placeholder {
  padding: 6px 8px 6px 32px;
  font-size: 12px;
  color: #6b7280;
}

.loading-text {
  font-style: italic;
}
</style>
