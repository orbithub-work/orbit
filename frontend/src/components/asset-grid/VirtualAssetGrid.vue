<template>
  <div class="virtual-grid-container" ref="containerRef">
    <!-- Grid View (宫格视图) -->
    <div
      v-if="viewMode === 'grid'"
      class="grid-view"
      ref="gridRef"
    >
      <div
        v-for="item in props.items"
        :key="item.id"
        class="grid-card"
        :class="{ 'grid-card--selected': selectedIds.has(item.id) }"
        :style="{
          width: minItemWidth + 'px',
        }"
        @click.stop="handleClick(item.id, $event)"
        @dblclick="handleDoubleClick(item.id)"
        @contextmenu.prevent="handleContextMenu(item.id, $event)"
      >
        <div class="grid-thumbnail" :style="{ height: minItemWidth + 'px' }">
          <img
            v-if="item.thumbnailUrl"
            :src="item.thumbnailUrl"
            :alt="item.name"
            loading="lazy"
          />
          <div v-else class="grid-placeholder">
            <span class="file-ext">{{ getFileExtension(item.name) }}</span>
          </div>
        </div>
        <div class="grid-meta">
          <div class="grid-title" :title="item.name">{{ item.name }}</div>
        </div>
      </div>
    </div>

    <RecycleScroller
      v-else-if="viewMode === 'list'"
      class="scroller"
      :items="listRowData"
      :item-size="48"
      key-field="id"
      v-slot="{ item }"
      :buffer="400"
    >
      <div
        class="list-row"
        :class="{ 'list-row--selected': selectedIds.has(item.id) }"
        @click.stop="handleClick(item.id, $event)"
        @dblclick="handleDoubleClick(item.id)"
        @contextmenu.prevent="handleContextMenu(item.id, $event)"
      >
        <div class="list-icon">
          <img v-if="item.thumbnailUrl" :src="item.thumbnailUrl" :alt="item.name" />
          <div v-else class="icon-placeholder">
            <span>{{ getFileExtension(item.name) }}</span>
          </div>
        </div>
        <div class="list-name">{{ item.name }}</div>
        <div class="list-size">{{ formatSize(item.size) }}</div>
        <div class="list-date">{{ formatDate(item.modified_at) }}</div>
      </div>
    </RecycleScroller>

    <div v-else-if="viewMode === 'gallery'" class="gallery-view">
      <div class="gallery-sidebar">
        <RecycleScroller
          ref="sidebarScrollerRef"
          class="sidebar-scroller"
          :items="props.items"
          :item-size="80"
          key-field="id"
          v-slot="{ item, index }"
          :buffer="200"
        >
          <div
            class="sidebar-thumb"
            :class="{ 'sidebar-thumb--active': index === currentIndex }"
            @click="selectGalleryItem(index)"
          >
            <img
              v-if="item.thumbnailUrl"
              :src="item.thumbnailUrl"
              :alt="item.name"
            />
            <div v-else class="thumb-placeholder">
              <span>{{ getFileExtension(item.name) }}</span>
            </div>
          </div>
        </RecycleScroller>
      </div>
      <div class="gallery-main">
        <div class="gallery-preview">
          <img
            v-if="currentGalleryItem?.thumbnailUrl"
            :src="currentGalleryItem.thumbnailUrl"
            :alt="currentGalleryItem?.name"
            class="preview-image"
          />
          <div v-else class="preview-placeholder">
            <span class="file-ext-large">{{ getFileExtension(currentGalleryItem?.name || '') }}</span>
          </div>
        </div>
        <div class="gallery-info-bar">
          <div class="info-name">{{ currentGalleryItem?.name }}</div>
          <div class="info-meta">
            <span>{{ getFileExtension(currentGalleryItem?.name || '') }}</span>
            <span>{{ formatSize(currentGalleryItem?.size || 0) }}</span>
            <span>{{ currentIndex + 1 }} / {{ props.items.length }}</span>
          </div>
        </div>
        <button class="nav-btn nav-prev" @click="prevItem" :disabled="currentIndex <= 0">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="15 18 9 12 15 6"></polyline>
          </svg>
        </button>
        <button class="nav-btn nav-next" @click="nextItem" :disabled="currentIndex >= props.items.length - 1">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="9 18 15 12 9 6"></polyline>
          </svg>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, shallowRef, onMounted, onUnmounted, computed, watch, nextTick } from 'vue'
import { RecycleScroller } from 'vue-virtual-scroller'
import 'vue-virtual-scroller/dist/vue-virtual-scroller.css'

interface AssetItem {
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

interface RowData {
  id: number
  items: AssetItem[]
}

interface Props {
  items: AssetItem[]
  loading?: boolean
  minItemWidth?: number
  gap?: number
  viewMode?: 'list' | 'grid' | 'gallery'
  hasMore?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  minItemWidth: 160,
  gap: 8,
  viewMode: 'grid'
})

const emit = defineEmits<{
  select: [items: AssetItem[]]
  open: [item: AssetItem]
  contextMenu: [item: AssetItem, event: MouseEvent]
  'load-more': []
}>()

const containerRef = ref<HTMLElement | null>(null)
const gridRef = ref<HTMLElement | null>(null)
const selectedIds = shallowRef<Set<string>>(new Set())
const containerWidth = ref(800)
const sidebarScrollerRef = ref<any>(null)

const currentIndex = ref(0)

// Masonry layout state
const masonryColumnCount = ref(4)
const masonryGap = computed(() => props.gap)
const masonryPadding = 16
const masonryMetaHeight = 36

// Virtual scroll state for masonry
const masonryScrollTop = ref(0)
const masonryContainerHeight = ref(800)
const masonryBuffer = 800 // 缓冲区高度（像素）

// 懒计算缓存
const masonryPositionCache = ref<Map<string, { left: number; top: number; width: number; height: number }>>(new Map())
const masonryColumnHeights = ref<number[]>([0, 0, 0, 0]) // 每列当前高度

// Track loaded images for masonry
const loadedImages = ref<Set<string>>(new Set())

// Masonry scroll handler
function onMasonryScroll(e: Event) {
  const target = e.target as HTMLElement
  masonryScrollTop.value = target.scrollTop
  masonryContainerHeight.value = target.clientHeight
  
  // 检查是否滚动到底部
  const scrollBottom = target.scrollTop + target.clientHeight
  const scrollHeight = target.scrollHeight
  const threshold = 200 // 距离底部200px时触发
  
  if (scrollBottom >= scrollHeight - threshold && props.hasMore && !props.loading) {
    emit('load-more')
  }
}

// 重置瀑布流缓存
function resetMasonryCache() {
  masonryPositionCache.value.clear()
  masonryColumnHeights.value = new Array(masonryColumnCount.value).fill(0)
}

const currentGalleryItem = computed(() => {
  return props.items[currentIndex.value] || null
})

function scrollToActive(index: number) {
  nextTick(() => {
    if (!sidebarScrollerRef.value) return
    sidebarScrollerRef.value.scrollToItem(index)
  })
}

function selectGalleryItem(index: number) {
  currentIndex.value = index
  const item = props.items[index]
  if (item) {
    selectedIds.value = new Set([item.id])
    emit('select', [item])
  }
  scrollToActive(index)
}

function prevItem() {
  if (currentIndex.value > 0) {
    selectGalleryItem(currentIndex.value - 1)
  }
}

function nextItem() {
  if (currentIndex.value < props.items.length - 1) {
    selectGalleryItem(currentIndex.value + 1)
  }
}

function handleKeydown(e: KeyboardEvent) {
  if (props.viewMode !== 'gallery') return
  if (e.key === 'ArrowLeft') {
    prevItem()
  } else if (e.key === 'ArrowRight') {
    nextItem()
  }
}

const gridColumnCount = computed(() => {
  const availableWidth = containerWidth.value - 20
  if (availableWidth <= 0) return 1
  
  let cols = 1
  while (true) {
    const requiredWidth = cols * props.minItemWidth + (cols - 1) * props.gap
    if (requiredWidth > availableWidth) {
      cols--
      break
    }
    cols++
  }
  return Math.max(1, cols)
})

const gridRowHeight = computed(() => {
  const availableWidth = containerWidth.value - 20
  const cardWidth = (availableWidth - (gridColumnCount.value - 1) * props.gap) / gridColumnCount.value
  const imageHeight = cardWidth * 0.75
  const textHeight = 12 + 17 + 4 + 15 + 12
  return Math.ceil(imageHeight + textHeight + props.gap)
})

const gridRowData = computed<RowData[]>(() => {
  const rows: RowData[] = []
  const cols = gridColumnCount.value
  
  for (let i = 0; i < props.items.length; i += cols) {
    rows.push({
      id: Math.floor(i / cols),
      items: props.items.slice(i, i + cols)
    })
  }
  
  return rows
})

const listRowData = computed(() => {
  return props.items.map((item) => ({
    ...item
  }))
})

// 瀑布流列宽
const masonryColumnWidth = computed(() => {
  const availableWidth = containerWidth.value - masonryPadding * 2
  const gaps = (masonryColumnCount.value - 1) * masonryGap.value
  return Math.floor((availableWidth - gaps) / masonryColumnCount.value)
})

// 计算单个项目的高度
function calculateItemHeight(item: AssetItem): number {
  let imageHeight = masonryColumnWidth.value
  if (item.width && item.height) {
    imageHeight = Math.round((masonryColumnWidth.value * item.height) / item.width)
  }
  imageHeight = Math.max(imageHeight, masonryColumnWidth.value * 0.5)
  imageHeight = Math.min(imageHeight, masonryColumnWidth.value * 2)
  return imageHeight + masonryMetaHeight
}

// 懒计算：获取或计算项目位置
function getItemPosition(item: AssetItem, index: number) {
  const cacheKey = item.id
  
  // 检查缓存
  if (masonryPositionCache.value.has(cacheKey)) {
    return masonryPositionCache.value.get(cacheKey)!
  }
  
  // 计算到当前索引的所有项目位置
  for (let i = masonryPositionCache.value.size; i <= index; i++) {
    const currentItem = props.items[i]
    if (!currentItem) continue
    
    // 找到最短的列
    const colHeights = masonryColumnHeights.value
    const shortestCol = colHeights.indexOf(Math.min(...colHeights))
    
    const itemHeight = calculateItemHeight(currentItem)
    const left = masonryPadding + shortestCol * (masonryColumnWidth.value + masonryGap.value)
    const top = colHeights[shortestCol]
    
    // 更新列高度
    colHeights[shortestCol] += itemHeight + masonryGap.value
    
    // 缓存位置
    masonryPositionCache.value.set(currentItem.id, {
      left,
      top,
      width: masonryColumnWidth.value,
      height: itemHeight
    })
  }
  
  return masonryPositionCache.value.get(cacheKey)
}

// 估算总高度（用于滚动条）
const estimatedTotalHeight = computed(() => {
  if (props.items.length === 0) return 0
  
  // 基于已计算的样本估算
  if (masonryPositionCache.value.size > 0) {
    const totalCalculated = masonryColumnHeights.value.reduce((a, b) => a + b, 0)
    const avgHeight = totalCalculated / masonryPositionCache.value.size
    const remaining = props.items.length - masonryPositionCache.value.size
    return Math.max(...masonryColumnHeights.value) + (remaining * avgHeight / masonryColumnCount.value)
  }
  
  // 粗略估算：假设平均高度250px
  return (props.items.length / masonryColumnCount.value) * 250
})

// 可视区域项目（懒计算）
const visibleMasonryItems = computed(() => {
  const viewportStart = masonryScrollTop.value - masonryBuffer
  const viewportEnd = masonryScrollTop.value + masonryContainerHeight.value + masonryBuffer
  
  const items: Array<AssetItem & { left: number; top: number; width: number; height: number }> = []
  
  // 二分查找优化：先找到大概的起始位置
  let startIndex = 0
  let endIndex = props.items.length
  
  // 如果缓存已经有数据，估算起始索引
  if (masonryPositionCache.value.size > 0) {
    const avgItemHeight = 250 // 平均高度估算
    const itemsPerScreen = Math.ceil((viewportEnd - viewportStart) / avgItemHeight) * masonryColumnCount.value
    startIndex = Math.max(0, Math.floor(viewportStart / avgItemHeight) * masonryColumnCount.value - itemsPerScreen)
    endIndex = Math.min(props.items.length, startIndex + itemsPerScreen * 3)
  }
  
  for (let i = startIndex; i < endIndex && i < props.items.length; i++) {
    const item = props.items[i]
    const position = getItemPosition(item, i)
    
    if (!position) continue
    
    const itemBottom = position.top + position.height
    
    // 检查是否在可视区域内
    if (itemBottom >= viewportStart && position.top <= viewportEnd) {
      items.push({
        ...item,
        ...position
      })
    } else if (position.top > viewportEnd) {
      // 如果已经超过可视区域底部，提前退出
      // 但需要继续计算一些缓冲区的项目
      if (items.length > 100) break
    }
  }
  
  return items
})

function updateMasonryColumns() {
  const minWidth = props.minItemWidth
  const availableWidth = containerWidth.value - masonryPadding * 2
  masonryColumnCount.value = Math.max(1, Math.floor(availableWidth / (minWidth + masonryGap.value)))
}

function onMasonryImageLoad(id: string) {
  loadedImages.value.add(id)
}

function handleClick(id: string, event: MouseEvent) {
  const item = props.items.find(i => i.id === id)
  if (!item) return

  if (event.ctrlKey || event.metaKey) {
    const newSet = new Set(selectedIds.value)
    if (newSet.has(id)) {
      newSet.delete(id)
    } else {
      newSet.add(id)
    }
    selectedIds.value = newSet
  } else {
    selectedIds.value = new Set([id])
  }
  emit('select', props.items.filter(i => selectedIds.value.has(i.id)))
}

function handleDoubleClick(id: string) {
  const item = props.items.find(i => i.id === id)
  if (item) {
    emit('open', item)
  }
}

function handleContextMenu(id: string, event: MouseEvent) {
  const item = props.items.find(i => i.id === id)
  if (!item) return

  if (!selectedIds.value.has(id)) {
    selectedIds.value = new Set([id])
    emit('select', [item])
  }
  emit('contextMenu', item, event)
}

function getFileExtension(name: string): string {
  if (!name) return ''
  const parts = name.split('.')
  return parts.length > 1 ? parts.pop()!.toUpperCase() : ''
}

function formatSize(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i]
}

function formatDate(timestamp?: number): string {
  if (!timestamp) return '-'
  const date = new Date(timestamp)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit'
  })
}

let resizeObserver: ResizeObserver | null = null

onMounted(() => {
  if (containerRef.value) {
    containerWidth.value = containerRef.value.clientWidth
    updateMasonryColumns()
    resizeObserver = new ResizeObserver((entries) => {
      for (const entry of entries) {
        containerWidth.value = entry.contentRect.width
        updateMasonryColumns()
      }
    })
    resizeObserver.observe(containerRef.value)
  }
  document.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  if (resizeObserver) {
    resizeObserver.disconnect()
  }
  document.removeEventListener('keydown', handleKeydown)
})

watch(() => props.viewMode, (newMode) => {
  if (newMode === 'gallery' && props.items.length > 0) {
    nextTick(() => {
      scrollToActive(currentIndex.value)
    })
  }
})

defineExpose({
  clearSelection: () => {
    selectedIds.value = new Set()
    emit('select', [])
  },
  selectAll: () => {
    selectedIds.value = new Set(props.items.map(i => i.id))
    emit('select', props.items)
  }
})
</script>

<style scoped>
.virtual-grid-container {
  width: 100%;
  height: 100%;
  min-height: 0;
  background: var(--color-bg-base, #18181b);
}

.scroller {
  height: 100%;
}

.virtual-grid-container :deep(.vue-recycle-scroller) {
  height: 100% !important;
}

.virtual-grid-container :deep(.vue-recycle-scroller__item-wrapper) {
  contain: layout;
}

.virtual-grid-container :deep(.vue-recycle-scroller__item-view) {
  overflow: visible;
}

/* Grid View */
.grid-view {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(v-bind(minItemWidth + 'px'), 1fr));
  gap: 12px;
  padding: 16px;
  height: 100%;
  overflow-y: auto;
}

.grid-card {
  background: var(--color-bg-surface, #2a2a2a);
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  border: 2px solid transparent;
  transition: all 0.2s ease-in-out;
}

.grid-card:hover {
  background: var(--color-bg-surface-hover, #333);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
}

.grid-card--selected {
  border-color: var(--color-primary, #3b82f6);
  background: var(--color-primary-bg, rgba(59, 130, 246, 0.15));
}

.grid-thumbnail {
  width: 100%;
  background: rgba(0, 0, 0, 0.3);
  position: relative;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
}

.grid-thumbnail img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.grid-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-text-tertiary, #6b7280);
}

.grid-placeholder .file-ext {
  font-size: 14px;
  font-weight: 600;
  padding: 4px 8px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 4px;
}

.grid-meta {
  padding: 8px;
}

.grid-title {
  font-size: 12px;
  color: var(--color-text-primary, #e5e7eb);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.grid-row {
  display: grid;
  grid-template-columns: repeat(v-bind(gridColumnCount), 1fr);
  gap: v-bind(gap + 'px');
  padding: 0 10px;
}

.grid-card {
  background: var(--color-bg-surface, #2a2a2a);
  border-radius: 12px;
  overflow: hidden;
  cursor: pointer;
  border: 2px solid transparent;
  transition: all 0.2s ease-in-out;
}

.grid-card:hover {
  background: var(--color-bg-surface-hover, #333);
  transform: scale(1.01);
  box-shadow: 0 5px 26px rgba(0, 0, 0, 0.45), 0 0 12px rgba(205, 221, 246, 0.18);
}

.grid-card--selected {
  border-color: var(--color-primary, #3b82f6);
  background: var(--color-primary-bg, rgba(59, 130, 246, 0.15));
}

.card-thumbnail {
  width: 100%;
  aspect-ratio: 4/3;
  background: rgba(0, 0, 0, 0.3);
  position: relative;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
}

.card-thumbnail img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.thumbnail-placeholder {
  color: var(--color-text-tertiary, #6b7280);
  opacity: 0.5;
}

.file-ext {
  font-size: 14px;
  font-weight: 600;
  padding: 4px 8px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 4px;
}

.card-meta {
  padding: 12px;
}

.meta-title {
  font-size: 13px;
  color: var(--color-text-primary, #e5e7eb);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-bottom: 4px;
}

.meta-details {
  display: flex;
  justify-content: space-between;
  font-size: 11px;
  color: var(--color-text-tertiary, #6b7280);
}

.meta-ext {
  background: rgba(255, 255, 255, 0.08);
  padding: 2px 6px;
  border-radius: 3px;
  font-size: 10px;
}

/* List View */
.list-row {
  display: flex;
  align-items: center;
  height: 48px;
  padding: 0 16px;
  gap: 12px;
  cursor: pointer;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  transition: all 0.15s ease;
}

.list-row:hover {
  background: rgba(255, 255, 255, 0.04);
}

.list-row--selected {
  background: rgba(59, 130, 246, 0.15);
  border-left: 3px solid var(--color-primary, #3b82f6);
}

.list-icon {
  width: 32px;
  height: 32px;
  border-radius: 6px;
  overflow: hidden;
  background: rgba(0, 0, 0, 0.3);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.list-icon img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.icon-placeholder {
  font-size: 10px;
  font-weight: 600;
  color: var(--color-text-tertiary, #6b7280);
  padding: 2px 4px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 3px;
}

.list-name {
  flex: 1;
  font-size: 13px;
  color: var(--color-text-primary, #e5e7eb);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.list-size {
  width: 80px;
  font-size: 12px;
  color: var(--color-text-tertiary, #6b7280);
  text-align: right;
}

.list-date {
  width: 100px;
  font-size: 12px;
  color: var(--color-text-tertiary, #6b7280);
  text-align: right;
}

/* Gallery View */
.gallery-view {
  display: flex;
  height: 100%;
  background: #0a0a0a;
}

.gallery-main {
  flex: 1;
  position: relative;
  display: flex;
  flex-direction: column;
}

.gallery-preview {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
  min-height: 0;
}

.preview-image {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
  border-radius: 8px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.5);
}

.preview-placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
}

.file-ext-large {
  font-size: 24px;
  font-weight: 600;
  padding: 16px 32px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  color: var(--color-text-tertiary, #6b7280);
}

.gallery-info-bar {
  padding: 12px 20px;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.info-name {
  font-size: 14px;
  color: #fff;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 60%;
}

.info-meta {
  display: flex;
  gap: 16px;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.6);
}

.nav-btn {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  background: rgba(0, 0, 0, 0.5);
  border: none;
  border-radius: 50%;
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: #fff;
  transition: all 0.2s;
  z-index: 10;
}

.nav-btn:hover:not(:disabled) {
  background: rgba(0, 0, 0, 0.8);
}

.nav-btn:disabled {
  opacity: 0.3;
  cursor: not-allowed;
}

.nav-prev {
  left: 16px;
}

.nav-next {
  right: 16px;
}

.gallery-sidebar {
  width: 100px;
  background: rgba(0, 0, 0, 0.4);
  border-right: 1px solid rgba(255, 255, 255, 0.1);
}

.sidebar-scroller {
  height: 100%;
}

.sidebar-thumb {
  width: 80px;
  height: 60px;
  margin: 8px 10px;
  border-radius: 6px;
  overflow: hidden;
  cursor: pointer;
  border: 2px solid transparent;
  transition: all 0.2s;
  background: rgba(0, 0, 0, 0.3);
}

.sidebar-thumb:hover {
  border-color: rgba(255, 255, 255, 0.3);
}

.sidebar-thumb--active {
  border-color: var(--color-primary, #3b82f6);
  box-shadow: 0 0 12px rgba(59, 130, 246, 0.4);
}

.sidebar-thumb img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.thumb-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 10px;
  font-weight: 600;
  color: var(--color-text-tertiary, #6b7280);
  background: rgba(255, 255, 255, 0.05);
}

/* Masonry View */
.masonry-view {
  position: relative;
  width: 100%;
  height: 100%;
  overflow-y: auto;
  padding-top: v-bind(masonryGap + 'px');
}

.masonry-card {
  background: var(--color-bg-surface, #2a2a2a);
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  border: 2px solid transparent;
  transition: all 0.2s ease;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
}

.masonry-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4), 0 0 12px rgba(205, 221, 246, 0.15);
  z-index: 1;
}

.masonry-card--selected {
  border-color: var(--color-primary, #3b82f6);
  box-shadow: 0 0 0 2px var(--color-primary, #3b82f6), 0 8px 24px rgba(0, 0, 0, 0.4);
  z-index: 2;
}

.masonry-thumbnail {
  width: 100%;
  height: calc(100% - v-bind(masonryMetaHeight + 'px'));
  background: rgba(0, 0, 0, 0.3);
  position: relative;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
}

.masonry-thumbnail img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.masonry-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.05);
}

.masonry-meta {
  height: v-bind(masonryMetaHeight + 'px');
  padding: 8px 10px;
  display: flex;
  align-items: center;
  background: var(--color-bg-surface, #2a2a2a);
  border-top: 1px solid rgba(255, 255, 255, 0.05);
}

.masonry-title {
  font-size: 12px;
  color: var(--color-text-primary, #e5e7eb);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
