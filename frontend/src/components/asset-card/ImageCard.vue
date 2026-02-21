<template>
  <AssetCard
    :asset="asset"
    :status="status"
    :selected="selected"
    :show-type-badge="showTypeBadge"
    class="image-card"
    @click="(a, e) => $emit('click', a, e)"
    @double-click="(a) => $emit('doubleClick', a)"
    @context-menu="(a, e) => $emit('contextMenu', a, e)"
    @preview="$emit('preview')"
    @favorite="$emit('favorite')"
  >
    <template #thumbnail>
      <div class="image-thumbnail">
        <img
          v-if="thumbnailUrl && !thumbnailError"
          :src="thumbnailUrl"
          :alt="asset.name"
          class="thumbnail-image"
          loading="lazy"
          @error="handleThumbnailError"
        />
        <div v-else class="image-placeholder">
          <svg class="placeholder-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="3" y="3" width="18" height="18" rx="2" ry="2" />
            <circle cx="8.5" cy="8.5" r="1.5" />
            <polyline points="21 15 16 10 5 21" />
          </svg>
        </div>
        
        <div v-if="asset.is_ai_generated" class="ai-badge">
          <svg class="ai-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2" />
          </svg>
          <span>AI</span>
        </div>

        <div v-if="asset.rating && asset.rating > 0" class="rating-badge">
          <span v-for="i in asset.rating" :key="i" class="star">★</span>
        </div>
      </div>
    </template>

    <template #meta>
      <span>{{ formatFileSize(asset.size) }}</span>
      <span v-if="asset.width && asset.height" class="dimensions">
        {{ asset.width }} × {{ asset.height }}
      </span>
    </template>

    <template #hover-actions>
      <button class="hover-btn" title="快速预览 (Space)" @click.stop="$emit('preview')">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" />
          <circle cx="12" cy="12" r="3" />
        </svg>
      </button>
      <button 
        v-if="asset.is_ai_generated && asset.prompt_id" 
        class="hover-btn hover-btn--ai" 
        title="查看 Prompt" 
        @click.stop="$emit('viewPrompt')"
      >
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
        </svg>
      </button>
      <button class="hover-btn" title="收藏" @click.stop="$emit('favorite')">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2" />
        </svg>
      </button>
      <button class="hover-btn" title="更多" @click.stop="$emit('more')">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="1" />
          <circle cx="19" cy="12" r="1" />
          <circle cx="5" cy="12" r="1" />
        </svg>
      </button>
    </template>
  </AssetCard>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import AssetCard, { type Asset, type AssetStatus } from './AssetCard.vue'

export interface ImageAsset extends Asset {
  width?: number
  height?: number
  is_ai_generated?: boolean
  prompt_id?: string
  rating?: number
}

interface Props {
  asset: ImageAsset
  status?: AssetStatus
  selected?: boolean
  thumbnailUrl?: string
  showTypeBadge?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  status: 'ready',
  selected: false,
  showTypeBadge: true
})

defineEmits<{
  click: [asset: ImageAsset, event: MouseEvent]
  doubleClick: [asset: ImageAsset]
  contextMenu: [asset: ImageAsset, event: MouseEvent]
  preview: []
  viewPrompt: []
  favorite: []
  more: []
}>()

const thumbnailError = ref(false)

function handleThumbnailError() {
  thumbnailError.value = true
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
.image-card {
  --card-radius: 8px;
}

.image-thumbnail {
  position: relative;
  width: 100%;
  aspect-ratio: 1;
  overflow: hidden;
  background: var(--color-surface, #2a2a2a);
}

.thumbnail-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.image-card:hover .thumbnail-image {
  transform: scale(1.05);
}

.image-placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #2a2a2a, #1a1a1a);
  color: var(--color-text-tertiary, #6b7280);
}

.placeholder-icon {
  width: 40px;
  height: 40px;
  opacity: 0.5;
}

.ai-badge {
  position: absolute;
  top: 8px;
  left: 8px;
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 3px 8px;
  background: linear-gradient(135deg, #8b5cf6, #6366f1);
  border-radius: 12px;
  font-size: 10px;
  font-weight: 600;
  color: #fff;
  box-shadow: 0 2px 8px rgba(139, 92, 246, 0.4);
}

.ai-icon {
  width: 10px;
  height: 10px;
}

.rating-badge {
  position: absolute;
  bottom: 8px;
  left: 8px;
  display: flex;
  gap: 1px;
}

.star {
  font-size: 10px;
  color: #fbbf24;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.5);
}

.dimensions {
  color: var(--color-text-tertiary, #6b7280);
}

.hover-btn--ai {
  background: linear-gradient(135deg, #8b5cf6, #6366f1) !important;
}

.hover-btn--ai:hover {
  background: linear-gradient(135deg, #7c3aed, #4f46e5) !important;
}
</style>
