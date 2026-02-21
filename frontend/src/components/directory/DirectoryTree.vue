<template>
  <div class="directory-tree">
    <DirectoryTreeNode
      v-for="node in rootNodes"
      :key="node.path"
      :node="node"
      :active-path="activePath"
      :depth="0"
      @select="handleSelect"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import DirectoryTreeNode from './DirectoryTreeNode.vue'
import { apiCall, type DirectoryNode } from '@/services/api'

const props = defineProps<{
  activePath?: string
}>()

const emit = defineEmits<{
  select: [path: string]
}>()

const rootNodes = ref<DirectoryNode[]>([])

onMounted(async () => {
  await loadLibrarySources()
})

async function loadLibrarySources() {
  try {
    const sources = await apiCall<any[]>('list_library_sources')
    rootNodes.value = (sources || []).map(s => ({
      path: s.root_path || s.rootPath,
      name: s.root_path?.split('/').pop() || s.rootPath?.split('/').pop() || 'Unknown',
      has_children: true
    }))
  } catch (error) {
    console.error('Failed to load library sources:', error)
  }
}

function handleSelect(path: string) {
  emit('select', path)
}
</script>

<style scoped>
.directory-tree {
  display: flex;
  flex-direction: column;
}
</style>
