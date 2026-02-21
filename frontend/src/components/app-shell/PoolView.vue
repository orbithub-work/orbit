<template>
  <div class="pool-view">
    <header class="eagle-header">
      <div class="header-top-row">
        <div class="header-left">
          <button class="nav-btn" @click="$emit('navigate', '')" title="根目录">
            <Icon name="home" size="lg" />
          </button>
          <button class="nav-btn" @click="goUp" :disabled="!currentPath" title="上级目录">
            <Icon name="chevron-left" size="lg" />
          </button>
          <span class="breadcrumb">{{ breadcrumbText }}</span>
        </div>
        <div class="header-center">
          <div class="search-box">
            <Icon name="search" size="sm" class="search-icon" />
            <input 
              class="search-input" 
              type="text" 
              placeholder="搜索素材" 
              v-model="searchQuery"
              @keyup.enter="handleSearch"
            />
            <!-- ✅ 简洁的快捷筛选 -->
            <div class="search-quick-filters">
              <button 
                v-for="filter in quickFiltersSimple" 
                :key="filter.value"
                :class="['quick-btn', { active: quickFilter === filter.value }]"
                :title="filter.label"
                @click="toggleQuickFilter(filter.value)"
              >
                {{ filter.icon }}
              </button>
            </div>
          </div>
        </div>
        <div class="header-right">
          <input v-if="viewMode === 'grid'" type="range" min="50" max="300" class="zoom-slider" v-model.number="zoomLevel" />
          <div class="view-buttons">
            <button class="view-btn" :class="{ active: viewMode === 'list' }" title="列表视图" @click="viewMode = 'list'">
              <Icon name="list" size="md" />
            </button>
            <button class="view-btn" :class="{ active: viewMode === 'grid' }" title="宫格视图" @click="viewMode = 'grid'">
              <Icon name="grid" size="md" />
            </button>
          </div>
        </div>
      </div>
      <div class="header-bottom-row">
        <div class="filter-row" ref="filterRowRef">
          <div class="filter-item">
            <Icon name="play" size="md" />
          </div>
          <button 
            v-for="filter in filterButtons" 
            :key="filter.id"
            class="filter-btn" 
            :class="{ 'filter-btn--active': isFilterActive(filter.id), 'filter-btn--open': activeFilter === filter.id }" 
            :ref="el => setButtonRef(filter.id, el as HTMLElement)"
            @click="toggleFilter(filter.id)"
          >
            {{ filter.label }}
            <Icon name="chevron-down" size="sm" class="filter-arrow" :class="{ 'filter-arrow--open': activeFilter === filter.id }" />
          </button>
        </div>
      </div>
    </header>

    <FilterDropdown 
      :visible="activeFilter === 'time'" 
      :target-rect="getButtonRect('time')"
      @close="activeFilter = null"
    >
      <TimeFilter
        v-model="selectedTimeRange"
        v-model:start-date="customDateStart"
        v-model:end-date="customDateEnd"
        @apply="activeFilter = null"
      />
    </FilterDropdown>

    <FilterDropdown 
      :visible="activeFilter === 'type'" 
      :target-rect="getButtonRect('type')"
      @close="activeFilter = null"
    >
      <TypeFilter
        v-model="selectedFileTypes"
        :types="fileTypes"
        @apply="activeFilter = null"
      />
    </FilterDropdown>

    <FilterDropdown 
      :visible="activeFilter === 'shape'" 
      :target-rect="getButtonRect('shape')"
      @close="activeFilter = null"
    >
      <ShapeFilter
        v-model="selectedShapes"
        @apply="activeFilter = null"
      />
    </FilterDropdown>

    <FilterDropdown 
      :visible="activeFilter === 'size'" 
      :target-rect="getButtonRect('size')"
      @close="activeFilter = null"
    >
      <SizeFilter
        v-model="selectedSizePreset"
        v-model:width-min="customWidthMin"
        v-model:width-max="customWidthMax"
        v-model:height-min="customHeightMin"
        v-model:height-max="customHeightMax"
        @apply="activeFilter = null"
      />
    </FilterDropdown>

    <FilterDropdown 
      :visible="activeFilter === 'rating'" 
      :target-rect="getButtonRect('rating')"
      @close="activeFilter = null"
    >
      <RatingFilter
        v-model="selectedRating"
        @apply="activeFilter = null"
      />
    </FilterDropdown>

    <FilterDropdown 
      :visible="activeFilter === 'fileSize'" 
      :target-rect="getButtonRect('fileSize')"
      @close="activeFilter = null"
    >
      <FileSizeFilter
        v-model="selectedFileSizes"
        @apply="activeFilter = null"
      />
    </FilterDropdown>

    <FilterDropdown 
      :visible="activeFilter === 'sort'" 
      :target-rect="getButtonRect('sort')"
      @close="activeFilter = null"
    >
      <SortFilter
        v-model="selectedSort"
        v-model:direction="sortDirection"
      />
    </FilterDropdown>

    <div class="assets-count-bar">
      <span class="count-text">素材 ({{ displayItems.length }})</span>
      <div v-if="hasActiveFilters" class="active-filter-tags">
        <span v-if="selectedFileTypes.length" class="filter-tag">
          类型: {{ selectedFileTypes.length }}项
          <button class="tag-remove" @click="selectedFileTypes = []">×</button>
        </span>
        <span v-if="selectedTimeRange" class="filter-tag">
          时间: {{ getTimeRangeLabel(selectedTimeRange) }}
          <button class="tag-remove" @click="selectedTimeRange = ''">×</button>
        </span>
        <span v-if="selectedShapes.length" class="filter-tag">
          形状: {{ selectedShapes.length }}项
          <button class="tag-remove" @click="selectedShapes = []">×</button>
        </span>
        <span v-if="selectedSizePreset" class="filter-tag">
          尺寸: {{ getSizeLabel(selectedSizePreset) }}
          <button class="tag-remove" @click="selectedSizePreset = ''">×</button>
        </span>
        <span v-if="selectedRating >= 0" class="filter-tag">
          评分: {{ getRatingLabel(selectedRating) }}
          <button class="tag-remove" @click="selectedRating = -1">×</button>
        </span>
        <span v-if="selectedFileSizes.length" class="filter-tag">
          大小: {{ selectedFileSizes.length }}项
          <button class="tag-remove" @click="selectedFileSizes = []">×</button>
        </span>
      </div>
      <Icon name="chevron-down" size="sm" class="count-arrow" />
    </div>

    <div class="eagle-grid-virtual">
      <!-- 骨架屏：首次加载 -->
      <SkeletonGrid 
        v-if="loading && displayItems.length === 0"
        :count="12"
        :item-size="zoomLevel"
      />
      
      <!-- 空状态：加载完成但无数据 -->
      <EmptyState
        v-else-if="!loading && displayItems.length === 0"
        icon="📂"
        title="暂无素材"
        description="当前文件夹为空，请导入素材或切换到其他文件夹"
      />
      
      <!-- 正常显示 -->
      <VirtualAssetGrid
        v-else
        :key="viewMode"
        :items="displayItems"
        :loading="isLoadingMore"
        :min-item-width="zoomLevel"
        :gap="8"
        :view-mode="viewMode"
        :has-more="hasMore"
        @select="handleSelect"
        @open="handleOpen"
        @context-menu="handleContextMenu"
        @load-more="loadMoreAssets"
      />

      <div v-if="displayItems.length > 0" class="grid-stats">
        <span>共 {{ displayItems.length }} 项</span>
        <span v-if="selectedCount > 0">已选 {{ selectedCount }} 个</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick, onUnmounted } from 'vue'
import Icon from '@/components/common/Icon.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import SkeletonGrid from '@/components/common/SkeletonGrid.vue'
import VirtualAssetGrid from '@/components/asset-grid/VirtualAssetGrid.vue'
import FilterDropdown from '@/components/filters/FilterDropdown.vue'
import TimeFilter from '@/components/filters/TimeFilter.vue'
import TypeFilter from '@/components/filters/TypeFilter.vue'
import ShapeFilter from '@/components/filters/ShapeFilter.vue'
import SizeFilter from '@/components/filters/SizeFilter.vue'
import RatingFilter from '@/components/filters/RatingFilter.vue'
import FileSizeFilter from '@/components/filters/FileSizeFilter.vue'
import SortFilter from '@/components/filters/SortFilter.vue'
import type { FileItem } from '@/stores/fileStore'
import { apiCall } from '@/services/api'

interface GridItem {
  id: string
  name: string
  path: string
  size: number
  is_directory?: boolean
  file_type?: string
  thumbnailUrl?: string
  modified_at?: number
  width?: number
  height?: number
}

const props = defineProps<{
  assets: FileItem[]
  loading: boolean
  selectedCount: number
  currentPath: string
  projectId?: string
}>()

const emit = defineEmits<{
  select: [items: FileItem[]]
  open: [item: FileItem]
  'context-menu': [item: FileItem, event: MouseEvent]
  navigate: [path: string]
}>()

const searchQuery = ref('')
const zoomLevel = ref(160)
const viewMode = ref<'list' | 'grid'>('grid')
const activeFilter = ref<string | null>(null)

// ✅ 智能搜索状态
const quickFilter = ref('')
const datePreset = ref('')

// 简化的快捷筛选
const quickFiltersSimple = [
  { value: 'recent', label: '最近', icon: '🕐' },
  { value: 'unrated', label: '未评分', icon: '⭐' },
  { value: 'vertical', label: '竖屏', icon: '📱' },
]

// 分页加载状态
const currentPage = ref(1)
const pageSize = ref(50)
const hasMore = ref(true)
const isLoadingMore = ref(false)
const allAssets = ref<FileItem[]>([])
const nextCursor = ref<string | null>(null)

const selectedFileTypes = ref<string[]>([])
const selectedTimeRange = ref('')
const selectedShapes = ref<string[]>([])
const selectedSizePreset = ref('')
const selectedRating = ref(-1)
const selectedFileSizes = ref<string[]>([])
const selectedSort = ref('name')
const sortDirection = ref<'asc' | 'desc'>('desc')

const customDateStart = ref('')
const customDateEnd = ref('')
const customWidthMin = ref<number | null>(null)
const customWidthMax = ref<number | null>(null)
const customHeightMin = ref<number | null>(null)
const customHeightMax = ref<number | null>(null)

const previousViewMode = ref<'list' | 'grid'>('grid')

const filterRowRef = ref<HTMLElement | null>(null)
const buttonRefs: Record<string, HTMLElement | null> = {}

const filterButtons = [
  { id: 'time', label: '导入时间' },
  { id: 'type', label: '类型' },
  { id: 'shape', label: '形状' },
  { id: 'size', label: '尺寸' },
  { id: 'rating', label: '评分' },
  { id: 'fileSize', label: '文件大小' },
  { id: 'sort', label: '排序' },
]

function setButtonRef(id: string, el: HTMLElement | null) {
  buttonRefs[id] = el
}

function getButtonRect(id: string): DOMRect | null {
  return buttonRefs[id]?.getBoundingClientRect() || null
}

function isFilterActive(id: string): boolean {
  switch (id) {
    case 'time': return !!selectedTimeRange.value || !!(customDateStart.value || customDateEnd.value)
    case 'type': return selectedFileTypes.value.length > 0
    case 'shape': return selectedShapes.value.length > 0
    case 'size': return !!selectedSizePreset.value || !!(customWidthMin.value || customWidthMax.value || customHeightMin.value || customHeightMax.value)
    case 'rating': return selectedRating.value >= 0
    case 'fileSize': return selectedFileSizes.value.length > 0
    case 'sort': return selectedSort.value !== 'name' || sortDirection.value !== 'desc'
    default: return false
  }
}

const hasActiveFilters = computed(() => {
  return filterButtons.some(f => isFilterActive(f.id))
})

watch(searchQuery, (newQuery) => {
  if (newQuery && viewMode.value !== 'grid' && viewMode.value !== 'masonry') {
    previousViewMode.value = viewMode.value
    viewMode.value = 'grid'
  }
})

const sourceAssets = computed(() => (props.projectId ? allAssets.value : props.assets))

const fileTypes = computed(() => {
  const types: Record<string, number> = {}
  sourceAssets.value.forEach(asset => {
    const ext = asset.name.split('.').pop()?.toLowerCase() || ''
    types[ext] = (types[ext] || 0) + 1
  })
  
  const sortedTypes = Object.entries(types)
    .map(([ext, count]) => ({ ext, count, label: ext.toUpperCase() }))
    .sort((a, b) => b.count - a.count)
  
  return sortedTypes
})

function getTimeRangeLabel(range: string): string {
  const labels: Record<string, string> = {
    today: '今天',
    week: '本周',
    month: '本月',
    year: '今年',
  }
  return labels[range] || range
}

function getSizeLabel(preset: string): string {
  const labels: Record<string, string> = {
    '4k': '4K',
    '2k': '2K',
    '1080p': '1080P',
    '720p': '720P',
    'small': '小图',
  }
  return labels[preset] || preset
}

function getRatingLabel(rating: number): string {
  const labels = ['未评分', '1星', '2星', '3星', '4星', '5星']
  return labels[rating] || ''
}

const breadcrumbText = computed(() => {
  if (!props.currentPath) return '根目录'
  return props.currentPath.split('/').filter(Boolean).join(' / ') || '根目录'
})

// 加载更多数据
async function loadMoreAssets() {
  if (isLoadingMore.value || !hasMore.value) return
  
  isLoadingMore.value = true
  try {
    if (props.projectId) {
      const params: Record<string, any> = {
        projectId: props.projectId,
        limit: pageSize.value,
        sortBy: selectedSort.value,
        sortOrder: sortDirection.value,
        // ✅ 添加智能搜索参数
        quickFilter: quickFilter.value || undefined,
        datePreset: datePreset.value || undefined,
      }
      if (nextCursor.value) {
        params.cursor = nextCursor.value
      }

      const result = await apiCall<{
        items: any[]
        nextCursor: string | null
        hasMore: boolean
        total: number
      }>('list_assets', params)
      
      if (result && result.items) {
        const newAssets = result.items.map((item): FileItem => {
          const fileName = item.name || item.path?.split('/').pop() || ''
          const ext = fileName.split('.').pop()?.toLowerCase() || ''
          
          let fileType = 'file'
          if (['jpg', 'jpeg', 'png', 'gif', 'webp', 'bmp', 'svg'].includes(ext)) fileType = 'image'
          else if (['mp4', 'avi', 'mov', 'mkv', 'webm'].includes(ext)) fileType = 'video'
          else if (['mp3', 'wav', 'flac', 'aac', 'ogg'].includes(ext)) fileType = 'audio'
          
          const modified = item.mtime ?? item.modified_at

          return {
            id: item.id,
            name: fileName,
            path: item.path || '',
            size: item.size || 0,
            file_type: fileType,
            is_directory: false,
            modified_at: typeof modified === 'number' ? new Date(modified * 1000) : undefined,
            thumbnail_path: `assets://thumbnails/${item.id}`,
            width: item.width,
            height: item.height
          }
        })
        
        allAssets.value.push(...newAssets)
        nextCursor.value = result.nextCursor
        hasMore.value = result.hasMore !== false
      } else {
        hasMore.value = false
      }
    } else {
      const startIndex = (currentPage.value - 1) * pageSize.value
      const nextSlice = props.assets.slice(startIndex, startIndex + pageSize.value)
      if (nextSlice.length === 0) {
        hasMore.value = false
        return
      }
      allAssets.value.push(...nextSlice)
      currentPage.value += 1
      hasMore.value = props.assets.length > startIndex + nextSlice.length
    }
  } catch (e) {
    console.error('Failed to load more:', e)
    hasMore.value = false
  } finally {
    isLoadingMore.value = false
  }
}

// 重置分页
function resetPagination() {
  currentPage.value = 1
  hasMore.value = true
  allAssets.value = []
  nextCursor.value = null
}

const displayItems = computed<GridItem[]>(() => {
  let items = allAssets.value.length > 0 ? allAssets.value : props.assets.slice(0, pageSize.value)
  
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    items = items.filter(f => f.name.toLowerCase().includes(query))
  }
  
  if (selectedFileTypes.value.length > 0) {
    items = items.filter(f => {
      const ext = f.name.split('.').pop()?.toLowerCase() || ''
      return selectedFileTypes.value.includes(ext)
    })
  }
  
  return items.map(f => {
    const ext = f.name.split('.').pop()?.toLowerCase() || ''
    const isImage = ['jpg', 'jpeg', 'png', 'gif', 'webp', 'bmp', 'svg'].includes(ext)
    const isVideo = ['mp4', 'avi', 'mov', 'mkv', 'webm'].includes(ext)

    return {
      id: f.id,
      name: f.name,
      path: f.path,
      size: f.size,
      is_directory: f.is_directory,
      file_type: isImage ? 'image' : isVideo ? 'video' : f.file_type,
      thumbnailUrl: f.thumbnail_path,
      modified_at: typeof f.modified_at === 'string'
        ? new Date(f.modified_at).getTime()
        : f.modified_at instanceof Date
          ? f.modified_at.getTime()
          : undefined,
      width: f.width,
      height: f.height
    }
  })
})

watch(() => props.assets, async (newAssets) => {
  if (props.projectId) return
  resetPagination()

  // 兼容模式：使用本地数据 slice
  if (newAssets.length > 0) {
    const firstPage = newAssets.slice(0, pageSize.value)
    allAssets.value = [...firstPage]
    hasMore.value = newAssets.length > pageSize.value
    currentPage.value = 2
    loadThumbnails()
  } else {
    allAssets.value = []
    hasMore.value = false
  }
}, { immediate: true })

watch(() => props.projectId, async (newProjectId) => {
  resetPagination()
  if (newProjectId) {
    await loadMoreAssets()
  } else {
    const firstPage = props.assets.slice(0, pageSize.value)
    allAssets.value = [...firstPage]
    hasMore.value = props.assets.length > pageSize.value
    currentPage.value = firstPage.length > 0 ? 2 : 1
    loadThumbnails()
  }
}, { immediate: true })

function goUp() {
  if (!props.currentPath) return
  const parts = props.currentPath.split('/').filter(Boolean)
  parts.pop()
  emit('navigate', '/' + parts.join('/'))
}

// ✅ 智能搜索处理
function toggleQuickFilter(value: string) {
  quickFilter.value = quickFilter.value === value ? '' : value
  handleSmartFilterChange()
}

function handleSmartFilterChange() {
  resetPagination()
  if (props.projectId) {
    loadMoreAssets()
  }
}

function handleSearch() {
  handleSmartFilterChange()
}

function handleSelect(items: GridItem[]) {
  const fileItems = items
    .map(item => sourceAssets.value.find(f => f.id === item.id))
    .filter((f): f is FileItem => f !== undefined)
  emit('select', fileItems)
}

function handleOpen(item: GridItem) {
  const original = sourceAssets.value.find(f => f.id === item.id)
  if (original) {
    emit('open', original)
  }
}

function handleContextMenu(item: GridItem, event: MouseEvent) {
  const original = sourceAssets.value.find(f => f.id === item.id)
  if (original) {
    emit('context-menu', original, event)
  }
}

function toggleFilter(filterType: string) {
  if (activeFilter.value === filterType) {
    activeFilter.value = null
  } else {
    activeFilter.value = filterType
  }
}
</script>

<style scoped>
.pool-view {
  flex: 1;
  display: flex;
  flex-direction: column;
  height: 100%;
  min-height: 0;
  background: #1b1c1f;
  overflow: hidden;
}

.eagle-header {
  display: flex;
  flex-direction: column;
  background: #1b1c1f;
  border-bottom: 1px solid #2b2b2f;
  flex-shrink: 0;
  position: relative;
}

.header-top-row {
  height: 48px;
  display: grid;
  grid-template-columns: 1fr auto 1fr;
  align-items: center;
  padding: 0 16px;
  gap: 8px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.nav-btn {
  background: transparent;
  border: none;
  cursor: pointer;
  padding: 6px;
  border-radius: 6px;
  color: #9ca3af;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.nav-btn:hover {
  background: rgba(255, 255, 255, 0.08);
  color: #e5e7eb;
}

.nav-btn:disabled {
  opacity: 0.3;
  cursor: not-allowed;
}

.breadcrumb {
  font-size: 13px;
  color: #9ca3af;
  margin-left: 4px;
}

.header-center {
  display: flex;
  justify-content: center;
}

.search-box {
  position: relative;
  display: flex;
  align-items: center;
  width: 400px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  padding: 0 12px;
  transition: all 0.2s;
}

.search-box:focus-within {
  background: rgba(255, 255, 255, 0.08);
  border-color: rgba(59, 130, 246, 0.5);
}

.search-icon {
  color: #6b7280;
  margin-right: 8px;
}

.search-input {
  flex: 1;
  background: transparent;
  border: none;
  outline: none;
  color: #e5e7eb;
  font-size: 14px;
  padding: 10px 0;
}

.search-input::placeholder {
  color: #6b7280;
}

.search-quick-filters {
  display: flex;
  gap: 4px;
  margin-left: 8px;
  padding-left: 8px;
  border-left: 1px solid rgba(255, 255, 255, 0.1);
}

.quick-btn {
  background: transparent;
  border: none;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 16px;
  line-height: 1;
  cursor: pointer;
  opacity: 0.5;
  transition: all 0.15s;
}

.quick-btn:hover {
  opacity: 0.8;
  background: rgba(255, 255, 255, 0.05);
}

.quick-btn.active {
  opacity: 1;
  background: rgba(59, 130, 246, 0.15);
}

.header-right {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 12px;
}

.zoom-slider {
  width: 80px;
  height: 4px;
  -webkit-appearance: none;
  appearance: none;
  background: #3a3a3a;
  border-radius: 2px;
  cursor: pointer;
}

.zoom-slider::-webkit-slider-thumb {
  -webkit-appearance: none;
  appearance: none;
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: #fff;
  cursor: pointer;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
}

.view-buttons {
  display: flex;
  gap: 2px;
  background: #2a2a2a;
  border-radius: 6px;
  padding: 2px;
}

.view-btn {
  background: transparent;
  border: none;
  cursor: pointer;
  padding: 6px;
  border-radius: 4px;
  color: #9ca3af;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.view-btn:hover {
  color: #e5e7eb;
}

.view-btn.active {
  background: #3a3a3a;
  color: #e5e7eb;
}

.header-bottom-row {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.filter-row {
  padding: 0 16px 8px 16px;
}

.filter-row {
  display: flex;
  align-items: center;
  gap: 4px;
}

.filter-item {
  padding: 6px 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #9ca3af;
}

.filter-btn {
  background: transparent;
  border: none;
  cursor: pointer;
  padding: 6px 12px;
  border-radius: 6px;
  color: #9ca3af;
  font-size: 13px;
  font-weight: 500;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  transition: all 0.2s;
  white-space: nowrap;
}

.filter-btn:hover {
  background: rgba(255, 255, 255, 0.06);
  color: #e5e7eb;
}

.filter-btn--active {
  background: rgba(249, 115, 22, 0.15);
  color: #f97316;
}

.filter-arrow {
  opacity: 0.7;
  transition: transform 0.2s;
  flex-shrink: 0;
}

.filter-arrow--open {
  transform: rotate(180deg);
}

.assets-count-bar {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  background: #1b1c1f;
  color: #9ca3af;
  font-size: 13px;
  border-bottom: 1px solid #2b2b2f;
}

.count-text {
  color: #e5e7eb;
}

.active-filter-tags {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.filter-tag {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 2px 8px;
  background: rgba(249, 115, 22, 0.15);
  border-radius: 4px;
  font-size: 12px;
  color: #f97316;
}

.tag-remove {
  background: transparent;
  border: none;
  color: #f97316;
  cursor: pointer;
  font-size: 14px;
  padding: 0;
  line-height: 1;
}

.tag-remove:hover {
  color: #ef4444;
}

.count-arrow {
  opacity: 0.7;
  margin-left: auto;
}

.eagle-grid-virtual {
  flex: 1;
  min-height: 0;
  position: relative;
  padding: 8px 8px 0 8px;
  background: #1b1c1f;
  overflow: hidden;
}

.eagle-grid-virtual :deep(.virtual-grid-container) {
  position: absolute;
  top: 8px;
  left: 8px;
  right: 8px;
  bottom: 8px;
  background: #1b1c1f;
}

.grid-stats {
  position: absolute;
  bottom: 16px;
  left: 24px;
  display: flex;
  gap: 16px;
  font-size: 11px;
  color: #6b7280;
  background: rgba(30, 30, 30, 0.9);
  padding: 6px 12px;
  border-radius: 4px;
  pointer-events: none;
}
</style>
