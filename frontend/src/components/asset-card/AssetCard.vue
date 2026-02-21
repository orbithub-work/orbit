<template>
  <div
    class="asset-card"
    :class="[
      `asset-card--${status}`,
      `asset-card--${assetType}`,
      { 'asset-card--selected': selected }
    ]"
    @mouseenter="hovered = true"
    @mouseleave="hovered = false"
    @click="handleClick"
    @dblclick="handleDoubleClick"
    @contextmenu.prevent="handleContextMenu"
  >
    <div class="asset-card__thumbnail">
      <slot name="thumbnail">
        <div class="thumbnail-placeholder">
          <component :is="defaultIcon" class="placeholder-icon" />
        </div>
      </slot>

      <div v-if="status !== 'ready'" class="asset-card__overlay">
        <slot name="status-overlay">
          <template v-if="status === 'loading'">
            <div class="status-spinner" />
            <span class="status-text">加载中...</span>
          </template>
          
          <template v-else-if="status === 'processing'">
            <div class="status-progress">
              <div class="progress-bar" :style="{ width: `${progress}%` }" />
            </div>
            <span class="status-text">{{ progressText }}</span>
          </template>
          
          <template v-else-if="status === 'error'">
            <svg class="status-icon status-icon--error" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10" />
              <line x1="15" y1="9" x2="9" y2="15" />
              <line x1="9" y1="9" x2="15" y2="15" />
            </svg>
            <span class="status-text">{{ errorMessage }}</span>
          </template>
          
          <template v-else-if="status === 'missing'">
            <svg class="status-icon status-icon--warning" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z" />
              <line x1="12" y1="9" x2="12" y2="13" />
              <line x1="12" y1="17" x2="12.01" y2="17" />
            </svg>
            <span class="status-text">文件缺失</span>
          </template>
        </slot>
      </div>

      <div v-if="showTypeBadge" class="asset-card__badge">
        {{ fileTypeLabel }}
      </div>

      <div v-if="selected" class="asset-card__check">
        <svg class="check-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
          <polyline points="20 6 9 17 4 12" />
        </svg>
      </div>

      <Transition name="fade">
        <div v-if="hovered && status === 'ready'" class="asset-card__hover-actions">
          <slot name="hover-actions">
            <button class="hover-btn" title="预览" @click.stop="$emit('preview')">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" />
                <circle cx="12" cy="12" r="3" />
              </svg>
            </button>
            <button class="hover-btn" title="收藏" @click.stop="$emit('favorite')">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2" />
              </svg>
            </button>
          </slot>
        </div>
      </Transition>
    </div>

    <div class="asset-card__info">
      <div class="info-name" :title="asset.name">
        {{ asset.name }}
      </div>
      <div class="info-meta">
        <span class="format-tag">{{ fileTypeLabel }}</span>
        <span class="file-size">{{ formatFileSize(asset.size) }}</span>
      </div>
    </div>

    <div v-if="$slots.extra" class="asset-card__extra">
      <slot name="extra" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, defineComponent, h } from 'vue'

export type AssetType = 'image' | 'video' | 'audio' | 'font' | 'document' | 'project' | 'archive' | 'other'
export type AssetStatus = 'loading' | 'processing' | 'ready' | 'error' | 'missing'

export interface Asset {
  id: string
  name: string
  path: string
  size: number
  file_type?: string
  thumbnail_url?: string
}

interface Props {
  asset: Asset
  status?: AssetStatus
  progress?: number
  progressText?: string
  errorMessage?: string
  selected?: boolean
  showTypeBadge?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  status: 'ready',
  progress: 0,
  progressText: '处理中...',
  errorMessage: '加载失败',
  selected: false,
  showTypeBadge: true
})

const emit = defineEmits<{
  click: [asset: Asset, event: MouseEvent]
  doubleClick: [asset: Asset]
  contextMenu: [asset: Asset, event: MouseEvent]
  preview: []
  favorite: []
}>()

const hovered = ref(false)

const assetType = computed<AssetType>(() => {
  const ext = getFileExtension(props.asset.name).toLowerCase()
  
  if (['jpg', 'jpeg', 'png', 'gif', 'webp', 'bmp', 'svg', 'heic', 'raw', 'cr2', 'nef', 'arw', 'dng'].includes(ext)) return 'image'
  if (['mp4', 'mov', 'avi', 'mkv', 'webm', 'flv', 'wmv', 'm4v', 'prores'].includes(ext)) return 'video'
  if (['mp3', 'wav', 'flac', 'aac', 'ogg', 'm4a', 'wma', 'aiff'].includes(ext)) return 'audio'
  if (['ttf', 'otf', 'woff', 'woff2', 'eot'].includes(ext)) return 'font'
  if (['pdf', 'doc', 'docx', 'xls', 'xlsx', 'ppt', 'pptx', 'txt', 'md'].includes(ext)) return 'document'
  if (['draft', 'prproj', 'aep', 'fcpxml', 'drp', 'vpj', 'cap', 'veg'].includes(ext)) return 'project'
  if (['zip', 'rar', '7z', 'tar', 'gz', 'bz2'].includes(ext)) return 'archive'
  return 'other'
})

const fileTypeLabel = computed(() => {
  return getFileExtension(props.asset.name).toUpperCase()
})

const defaultIcon = computed(() => {
  const icons: Record<AssetType, any> = {
    image: defineComponent({
      render: () => h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '2' }, [
        h('rect', { x: 3, y: 3, width: 18, height: 18, rx: 2, ry: 2 }),
        h('circle', { cx: 8.5, cy: 8.5, r: 1.5 }),
        h('polyline', { points: '21 15 16 10 5 21' })
      ])
    }),
    video: defineComponent({
      render: () => h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '2' }, [
        h('polygon', { points: '23 7 16 12 23 17 23 7' }),
        h('rect', { x: 1, y: 5, width: 15, height: 14, rx: 2, ry: 2 })
      ])
    }),
    audio: defineComponent({
      render: () => h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '2' }, [
        h('path', { d: 'M9 18V5l12-2v13' }),
        h('circle', { cx: 6, cy: 18, r: 3 }),
        h('circle', { cx: 18, cy: 16, r: 3 })
      ])
    }),
    font: defineComponent({
      render: () => h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '2' }, [
        h('polyline', { points: '4 7 4 4 20 4 20 7' }),
        h('line', { x1: 9, y1: 20, x2: 15, y2: 20 }),
        h('line', { x1: 12, y1: 4, x2: 12, y2: 20 })
      ])
    }),
    document: defineComponent({
      render: () => h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '2' }, [
        h('path', { d: 'M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z' }),
        h('polyline', { points: '14 2 14 8 20 8' }),
        h('line', { x1: 16, y1: 13, x2: 8, y2: 13 }),
        h('line', { x1: 16, y1: 17, x2: 8, y2: 17 }),
        h('polyline', { points: '10 9 9 9 8 9' })
      ])
    }),
    project: defineComponent({
      render: () => h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '2' }, [
        h('path', { d: 'M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z' })
      ])
    }),
    archive: defineComponent({
      render: () => h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '2' }, [
        h('polyline', { points: '21 8 21 21 3 21 3 8' }),
        h('rect', { x: 1, y: 3, width: 22, height: 5 }),
        h('line', { x1: 10, y1: 12, x2: 14, y2: 12 })
      ])
    }),
    other: defineComponent({
      render: () => h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '2' }, [
        h('path', { d: 'M13 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z' }),
        h('polyline', { points: '13 2 13 9 20 9' })
      ])
    })
  }
  return icons[assetType.value]
})

function handleClick(event: MouseEvent) {
  emit('click', props.asset, event)
}

function handleDoubleClick() {
  emit('doubleClick', props.asset)
}

function handleContextMenu(event: MouseEvent) {
  emit('contextMenu', props.asset, event)
}

function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(1))} ${sizes[i]}`
}

function getFileExtension(filename: string): string {
  const parts = filename.split('.')
  return parts.length > 1 ? parts[parts.length - 1] : ''
}

defineExpose({
  assetType,
  formatFileSize,
  getFileExtension
})
</script>

<style scoped>
.asset-card {
  position: relative;
  display: flex;
  flex-direction: column;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  background: var(--color-surface, #2a2a2a);
  border: 2px solid transparent;
  overflow: hidden;
}

.asset-card:hover {
  background: var(--color-hover, #333);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
}

.asset-card--selected {
  background: var(--color-selected, #1e3a5f) !important;
  border-color: var(--color-primary, #3b82f6);
}

.asset-card--loading .asset-card__thumbnail,
.asset-card--processing .asset-card__thumbnail,
.asset-card--error .asset-card__thumbnail,
.asset-card--missing .asset-card__thumbnail {
  opacity: 0.5;
}

.asset-card__thumbnail {
  position: relative;
  width: 100%;
  aspect-ratio: 1;
  border-radius: 6px 6px 0 0;
  overflow: hidden;
  background: var(--color-bg-elevated, #1e1e1e);
  display: flex;
  align-items: center;
  justify-content: center;
}

.thumbnail-placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  color: var(--color-text-tertiary, #6b7280);
}

.placeholder-icon {
  width: 48px;
  height: 48px;
  opacity: 0.5;
}

.asset-card__overlay {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
}

.status-spinner {
  width: 32px;
  height: 32px;
  border: 3px solid rgba(255, 255, 255, 0.2);
  border-top-color: var(--color-primary, #3b82f6);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.status-progress {
  width: 80%;
  height: 4px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 2px;
  overflow: hidden;
}

.progress-bar {
  height: 100%;
  background: var(--color-primary, #3b82f6);
  transition: width 0.3s;
}

.status-icon {
  width: 32px;
  height: 32px;
}

.status-icon--error {
  color: #ef4444;
}

.status-icon--warning {
  color: #f59e0b;
}

.status-text {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.9);
}

.asset-card__badge {
  position: absolute;
  bottom: 6px;
  right: 6px;
  padding: 2px 6px;
  font-size: 10px;
  font-weight: 600;
  text-transform: uppercase;
  background: rgba(0, 0, 0, 0.7);
  color: #fff;
  border-radius: 4px;
  backdrop-filter: blur(4px);
}

.asset-card__check {
  position: absolute;
  top: 6px;
  right: 6px;
  width: 24px;
  height: 24px;
  background: var(--color-primary, #3b82f6);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.check-icon {
  width: 14px;
  height: 14px;
}

.asset-card__hover-actions {
  position: absolute;
  top: 6px;
  left: 6px;
  display: flex;
  gap: 4px;
}

.hover-btn {
  width: 28px;
  height: 28px;
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.7);
  border: none;
  border-radius: 6px;
  color: #fff;
  cursor: pointer;
  transition: all 0.2s;
  backdrop-filter: blur(4px);
}

.hover-btn:hover {
  background: var(--color-primary, #3b82f6);
  transform: scale(1.1);
}

.hover-btn svg {
  width: 14px;
  height: 14px;
}

.asset-card__info {
  padding: 8px 10px;
}

.info-name {
  font-size: 13px;
  font-weight: 500;
  color: var(--color-text-primary, #e5e7eb);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  line-height: 1.4;
  margin-bottom: 6px;
}

.info-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 11px;
}

.format-tag {
  padding: 2px 6px;
  font-size: 10px;
  font-weight: 600;
  text-transform: uppercase;
  background: rgba(255, 255, 255, 0.15);
  color: #ffffff;
  border-radius: 3px;
}

.file-size {
  color: var(--color-text-secondary, #9ca3af);
}

.asset-card__extra {
  padding: 0 10px 10px;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.15s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
