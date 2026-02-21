<!-- 资产浏览视图 - 使用新组件的示例 -->
<template>
  <div class="assets-view">
    <!-- 头部工具栏 -->
    <div class="assets-view__header">
      <div class="header-left">
        <h1 class="view-title">
          资产库
        </h1>
        <span class="assets-count">{{ assets.length }} 个文件</span>
      </div>

      <div class="header-right">
        <!-- 视图切换 -->
        <div class="view-toggle">
          <button
            class="toggle-btn"
            :class="{ active: viewMode === 'grid' }"
            title="网格视图"
            @click="viewMode = 'grid'"
          >
            <svg class="icon">
              <use href="#icon-grid" />
            </svg>
          </button>
          <button
            class="toggle-btn"
            :class="{ active: viewMode === 'list' }"
            title="列表视图"
            @click="viewMode = 'list'"
          >
            <svg class="icon">
              <use href="#icon-list" />
            </svg>
          </button>
        </div>

        <!-- 排序 -->
        <select
          v-model="sortBy"
          class="sort-select"
        >
          <option value="name">名称</option>
          <option value="date">修改时间</option>
          <option value="size">大小</option>
          <option value="type">类型</option>
        </select>

        <!-- 刷新 -->
        <button
          class="btn btn--icon"
          title="刷新"
          @click="loadAssets"
        >
          <svg class="icon">
            <use href="#icon-refresh" />
          </svg>
        </button>
      </div>
    </div>

    <!-- 主内容区 -->
    <div class="assets-view__content">
      <!-- 虚拟滚动网格 -->
      <VirtualAssetGrid
        v-if="viewMode === 'grid'"
        ref="gridRef"
        :items="sortedAssets"
        :loading="loading"
        :item-size="160"
        :gap="16"
        :show-metadata="true"
        empty-message="暂无文件，试试导入一些资产"
        @select="handleSelect"
        @open="handleOpen"
        @context-menu="handleContextMenu"
      />

      <!-- 列表视图（占位） -->
      <div
        v-else
        class="list-view-placeholder"
      >
        <p>列表视图开发中...</p>
      </div>
    </div>

    <!-- 批量操作工具栏 -->
    <BatchToolbar
      :selected-ids="selectedIds"
      :selected-assets="selectedAssets"
      @clear-selection="clearSelection"
      @success="handleBatchSuccess"
      @delete="handleBatchDelete"
      @export="handleBatchExport"
    />

    <!-- 快速预览 -->
    <QuickPreview
      v-model:visible="previewVisible"
      v-model:current-index="previewIndex"
      :assets="sortedAssets"
      @open="handleOpen"
      @open-in-folder="handleOpenInFolder"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import VirtualAssetGrid from '@/components/asset-grid/VirtualAssetGrid.vue'
import BatchToolbar from '@/components/batch/BatchToolbar.vue'
import QuickPreview from '@/components/preview/QuickPreview.vue'

interface Asset {
  id: string
  name: string
  path: string
  size: number
  is_directory?: boolean
  file_type?: string
  thumbnailUrl?: string
  width?: number
  height?: number
  modified_at?: number
}

const gridRef = ref<InstanceType<typeof VirtualAssetGrid> | null>(null)
const assets = ref<Asset[]>([])
const loading = ref(false)
const viewMode = ref<'grid' | 'list'>('grid')
const sortBy = ref<'name' | 'date' | 'size' | 'type'>('name')

// 选中状态
const selectedIds = ref<string[]>([])
const selectedAssets = computed(() =>
  assets.value.filter(asset => selectedIds.value.includes(asset.id))
)

// 预览状态
const previewVisible = ref(false)
const previewIndex = ref(0)

// 排序后的资产
const sortedAssets = computed(() => {
  const sorted = [...assets.value]

  switch (sortBy.value) {
    case 'name':
      return sorted.sort((a, b) => a.name.localeCompare(b.name))
    case 'date':
      return sorted.sort((a, b) => (b.modified_at || 0) - (a.modified_at || 0))
    case 'size':
      return sorted.sort((a, b) => b.size - a.size)
    case 'type':
      return sorted.sort((a, b) => {
        const extA = a.name.split('.').pop() || ''
        const extB = b.name.split('.').pop() || ''
        return extA.localeCompare(extB)
      })
    default:
      return sorted
  }
})

// 加载资产
async function loadAssets() {
  loading.value = true
  try {
    // 这里应该调用实际的 API
    // const response = await apiCall('list_assets')
    // assets.value = response

    // 模拟数据（演示用）
    assets.value = generateMockAssets(100)
  } finally {
    loading.value = false
  }
}

// 生成模拟数据
function generateMockAssets(count: number): Asset[] {
  const types = ['jpg', 'png', 'mp4', 'pdf', 'psd']
  const mockAssets: Asset[] = []

  for (let i = 0; i < count; i++) {
    const type = types[Math.floor(Math.random() * types.length)]
    mockAssets.push({
      id: `asset-${i}`,
      name: `file-${i}.${type}`,
      path: `/path/to/file-${i}.${type}`,
      size: Math.floor(Math.random() * 10000000),
      file_type: type,
      width: type === 'jpg' || type === 'png' ? 1920 : undefined,
      height: type === 'jpg' || type === 'png' ? 1080 : undefined,
      modified_at: Date.now() / 1000 - Math.random() * 86400 * 30,
      thumbnailUrl: type === 'jpg' || type === 'png' ? `https://picsum.photos/200/200?random=${i}` : undefined
    })
  }

  return mockAssets
}

// 处理选择
function handleSelect(items: Asset[]) {
  selectedIds.value = items.map(item => item.id)
}

// 清除选择
function clearSelection() {
  selectedIds.value = []
  gridRef.value?.clearSelection()
}

// 处理打开
function handleOpen(asset: Asset) {
  if (asset.is_directory) {
    // 打开文件夹
    console.log('Open folder:', asset.path)
  } else {
    // 打开文件预览
    previewIndex.value = sortedAssets.value.findIndex(a => a.id === asset.id)
    previewVisible.value = true
  }
}

// 处理右键菜单
function handleContextMenu(asset: Asset, event: MouseEvent) {
  console.log('Context menu:', asset, event)
  // 这里应该显示右键菜单
}

// 在文件夹中显示
function handleOpenInFolder(asset: Asset) {
  console.log('Open in folder:', asset.path)
  // 调用系统 API 打开文件夹
}

// 批量操作成功
function handleBatchSuccess() {
  clearSelection()
  loadAssets()
}

// 批量删除
async function handleBatchDelete(ids: string[]) {
  try {
    // 调用删除 API
    console.log('Delete assets:', ids)
    await loadAssets()
    clearSelection()
  } catch (err) {
    console.error('Failed to delete assets:', err)
  }
}

// 批量导出
function handleBatchExport(ids: string[]) {
  console.log('Export assets:', ids)
  // 实现导出逻辑
}

// 键盘快捷键
function handleKeyDown(event: KeyboardEvent) {
  // Space 键打开预览
  if (event.key === ' ' && selectedIds.value.length === 1 && !previewVisible.value) {
    event.preventDefault()
    const asset = selectedAssets.value[0]
    handleOpen(asset)
  }
}

onMounted(() => {
  loadAssets()
  document.addEventListener('keydown', handleKeyDown)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeyDown)
})
</script>

<style scoped>
.assets-view {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: var(--color-bg-base, #1a1a1a);
}

.assets-view__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 24px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  background: var(--color-bg-elevated, #1e1e1e);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.view-title {
  font-size: 20px;
  font-weight: 600;
  color: #fff;
  margin: 0;
}

.assets-count {
  padding: 4px 12px;
  font-size: 13px;
  background: rgba(255, 255, 255, 0.05);
  color: #9ca3af;
  border-radius: 12px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.view-toggle {
  display: flex;
  gap: 4px;
  padding: 4px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 8px;
}

.toggle-btn {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: none;
  border: none;
  border-radius: 6px;
  color: #9ca3af;
  cursor: pointer;
  transition: all 0.2s;
}

.toggle-btn:hover {
  background: rgba(255, 255, 255, 0.05);
  color: #fff;
}

.toggle-btn.active {
  background: rgba(59, 130, 246, 0.2);
  color: #3b82f6;
}

.sort-select {
  padding: 8px 12px;
  font-size: 13px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 6px;
  color: #fff;
  cursor: pointer;
  outline: none;
}

.btn {
  padding: 8px 12px;
  font-size: 14px;
  font-weight: 500;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
}

.btn--icon {
  width: 36px;
  height: 36px;
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.05);
  color: #9ca3af;
}

.btn--icon:hover {
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
}

.assets-view__content {
  flex: 1;
  overflow: hidden;
}

.list-view-placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #6b7280;
}

.icon {
  width: 18px;
  height: 18px;
}
</style>
