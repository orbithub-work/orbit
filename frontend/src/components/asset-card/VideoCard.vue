<template>
  <AssetCard
    :asset="asset"
    :status="status"
    :selected="selected"
    :show-type-badge="showTypeBadge"
    class="video-card"
    @click="(a, e) => $emit('click', a, e)"
    @double-click="(a) => $emit('doubleClick', a)"
    @context-menu="(a, e) => $emit('contextMenu', a, e)"
    @preview="$emit('preview')"
    @favorite="$emit('favorite')"
  >
    <template #thumbnail>
      <div class="video-thumbnail">
        <img
          v-if="thumbnailUrl && !thumbnailError"
          :src="thumbnailUrl"
          :alt="asset.name"
          class="thumbnail-image"
          loading="lazy"
          @error="handleThumbnailError"
        />
        <div v-else class="video-placeholder">
          <svg class="placeholder-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polygon points="23 7 16 12 23 17 23 7" />
            <rect x="1" y="5" width="15" height="14" rx="2" ry="2" />
          </svg>
        </div>
        
        <div class="duration-badge">
          {{ formatDuration(asset.duration) }}
        </div>
        
        <div class="play-overlay" @click.stop="$emit('play')">
          <button class="play-btn">
            <svg viewBox="0 0 24 24" fill="currentColor">
              <polygon points="5 3 19 12 5 21 5 3" />
            </svg>
          </button>
        </div>

        <div v-if="asset.codec" class="codec-badge">
          {{ asset.codec }}
        </div>
      </div>
    </template>

    <template #meta>
      <span>{{ formatFileSize(asset.size) }}</span>
      <span>{{ formatDuration(asset.duration) }}</span>
      <span v-if="asset.fps">{{ asset.fps }}fps</span>
    </template>

    <template #hover-actions>
      <button class="hover-btn" title="播放" @click.stop="$emit('play')">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polygon points="5 3 19 12 5 21 5 3" />
        </svg>
      </button>
      <button class="hover-btn" title="提取帧" @click.stop="$emit('extractFrame')">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <rect x="3" y="3" width="18" height="18" rx="2" ry="2" />
          <circle cx="8.5" cy="8.5" r="1.5" />
          <polyline points="21 15 16 10 5 21" />
        </svg>
      </button>
      <button class="hover-btn" title="收藏" @click.stop="$emit('favorite')">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2" />
        </svg>
      </button>
    </template>
  </AssetCard>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import AssetCard, { type Asset, type AssetStatus } from './AssetCard.vue'

export interface VideoAsset extends Asset {
  duration?: number
  codec?: string
  fps?: number
  width?: number
  height?: number
}

interface Props {
  asset: VideoAsset
  status?: AssetStatus
  selected?: boolean
  thumbnailUrl?: string
  showTypeBadge?: boolean
}

withDefaults(defineProps<Props>(), {
  status: 'ready',
  selected: false,
  showTypeBadge: true
})

defineEmits<{
  click: [asset: VideoAsset, event: MouseEvent]
  doubleClick: [asset: VideoAsset]
  contextMenu: [asset: VideoAsset, event: MouseEvent]
  preview: []
  play: []
  extractFrame: []
  favorite: []
}>()

const thumbnailError = ref(false)

function handleThumbnailError() {
  thumbnailError.value = true
}

function formatDuration(seconds?: number): string {
  if (!seconds) return '0:00'
  const mins = Math.floor(seconds / 60)
  const secs = Math.floor(seconds % 60)
  if (mins >= 60) {
    const hours = Math.floor(mins / 60)
    const remainMins = mins % 60
    return `${hours}:${remainMins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
  }
  return `${mins}:${secs.toString().padStart(2, '0')}`
}

function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(1))} ${sizes[i]}`
}
</script>

<style scoped>
.video-card {
  --card-radius: 8px;
}

.video-thumbnail {
  position: relative;
  width: 100%;
  aspect-ratio: 16/9;
  overflow: hidden;
  background: var(--color-surface, #2a2a2a);
}

.thumbnail-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.video-card:hover .thumbnail-image {
  transform: scale(1.05);
}

.video-placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #1a1a2e, #16213e);
  color: var(--color-text-tertiary, #6b7280);
}

.placeholder-icon {
  width: 40px;
  height: 40px;
  opacity: 0.5;
}

.duration-badge {
  position: absolute;
  bottom: 8px;
  right: 8px;
  padding: 2px 8px;
  background: rgba(0, 0, 0, 0.8);
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
  color: #fff;
  backdrop-filter: blur(4px);
}

.codec-badge {
  position: absolute;
  top: 8px;
  left: 8px;
  padding: 2px 6px;
  background: rgba(59, 130, 246, 0.8);
  border-radius: 4px;
  font-size: 10px;
  font-weight: 600;
  color: #fff;
  text-transform: uppercase;
}

.play-overlay {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.3);
  opacity: 0;
  transition: opacity 0.2s;
  cursor: pointer;
}

.video-card:hover .play-overlay {
  opacity: 1;
}

.play-btn {
  width: 48px;
  height: 48px;
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.95);
  border: none;
  border-radius: 50%;
  color: #000;
  cursor: pointer;
  transition: all 0.2s;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
}

.play-btn:hover {
  transform: scale(1.1);
  background: #fff;
}

.play-btn svg {
  width: 20px;
  height: 20px;
  margin-left: 2px;
}
</style>
