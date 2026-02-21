<template>
  <div class="pool-view">
    <PoolToolbar
      :total-count="totalCount"
      :search-query="searchQuery"
      :view-mode="viewMode"
      @refresh="loadMoreAssets"
      @update:search-query="searchQuery = $event"
      @update:view-mode="viewMode = $event"
    />

    <div class="pool-content" ref="containerRef" @scroll="handleScroll">
      <EmptyState
        v-if="!loading && displayItems.length === 0"
        icon="image"
        title="暂无素材"
        description="开始导入素材文件"
      />

      <div v-else class="asset-grid" :style="gridStyle">
        <AssetGridItem
          v-for="asset in displayItems"
          :key="asset.id"
          :asset="asset"
          :is-selected="selectedIds.has(asset.id)"
          @select="handleSelect(asset)"
          @open="$emit('open', asset)"
          @context-menu="$emit('context-menu', asset, $event)"
        />
      </div>

      <div v-if="loading" class="loading-indicator">
        <Icon name="loader" size="lg" class="spin" />
        <span>加载中...</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { PoolToolbar, AssetGridItem } from '@/components/pool'
import EmptyState from '@/components/common/EmptyState.vue'
import Icon from '@/components/common/Icon.vue'
import type { FileItem } from '@/stores/fileStore'

const props = defineProps<{
  assets: FileItem[]
  loading: boolean
  selectedCount: number
  currentPath: string
  projectId?: string
}>()

const emit = defineEmits<{
  select: [items: FileItem[]]
  open: [asset: FileItem]
  'context-menu': [asset: FileItem, event: MouseEvent]
  navigate: [path: string]
}>()

const searchQuery = ref('')
const viewMode = ref<'list' | 'grid'>('grid')
const containerRef = ref<HTMLElement>()
const selectedIds = ref<Set<string>>(new Set())

const filteredAssets = computed(() => {
  if (!searchQuery.value) return props.assets
  const query = searchQuery.value.toLowerCase()
  return props.assets.filter(a => a.name.toLowerCase().includes(query))
})

const displayItems = computed(() => filteredAssets.value)
const totalCount = computed(() => filteredAssets.value.length)

const gridStyle = computed(() => ({
  display: 'grid',
  gridTemplateColumns: viewMode.value === 'grid' 
    ? 'repeat(auto-fill, minmax(180px, 1fr))' 
    : '1fr',
  gap: '16px',
  padding: '16px'
}))

function handleSelect(asset: FileItem) {
  if (selectedIds.value.has(asset.id)) {
    selectedIds.value.delete(asset.id)
  } else {
    selectedIds.value.add(asset.id)
  }
  emit('select', Array.from(selectedIds.value).map(id => 
    props.assets.find(a => a.id === id)!
  ).filter(Boolean))
}

function handleScroll() {
  const el = containerRef.value
  if (!el) return
  const { scrollTop, scrollHeight, clientHeight } = el
  if (scrollHeight - scrollTop - clientHeight < 200) {
    loadMoreAssets()
  }
}

async function loadMoreAssets() {
  // 触发父组件加载更多
}

watch(() => props.assets, () => {
  selectedIds.value.clear()
})
</script>

<style scoped>
.pool-view {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: #1a1a1a;
}

.pool-content {
  flex: 1;
  overflow-y: auto;
}

.asset-grid {
  min-height: 100%;
}

.loading-indicator {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 24px;
  color: #9ca3af;
}

.spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
</style>
