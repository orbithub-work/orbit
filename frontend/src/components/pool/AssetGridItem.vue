<template>
  <div 
    class="grid-item"
    :class="{ selected: isSelected }"
    @click="$emit('select', asset)"
    @dblclick="$emit('open', asset)"
    @contextmenu.prevent="$emit('context-menu', asset, $event)"
  >
    <div class="item-thumbnail">
      <img v-if="asset.thumbnail_url" :src="asset.thumbnail_url" :alt="asset.name" />
      <div v-else class="thumbnail-placeholder">
        <Icon :name="getFileIcon(asset.file_type)" size="lg" />
      </div>
    </div>
    <div class="item-info">
      <div class="item-name" :title="asset.name">{{ asset.name }}</div>
      <div class="item-meta">{{ formatFileSize(asset.size) }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import Icon from '@/components/common/Icon.vue'
import type { FileItem } from '@/stores/fileStore'

defineProps<{
  asset: FileItem
  isSelected: boolean
}>()

defineEmits<{
  select: [asset: FileItem]
  open: [asset: FileItem]
  'context-menu': [asset: FileItem, event: MouseEvent]
}>()

function getFileIcon(type?: string) {
  if (!type) return 'file'
  if (type.startsWith('image/')) return 'image'
  if (type.startsWith('video/')) return 'video'
  if (type.startsWith('audio/')) return 'music'
  return 'file'
}

function formatFileSize(bytes?: number) {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${(bytes / Math.pow(k, i)).toFixed(1)} ${sizes[i]}`
}
</script>

<style scoped>
.grid-item {
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.2s;
  background: rgba(255, 255, 255, 0.03);
}

.grid-item:hover {
  background: rgba(255, 255, 255, 0.06);
  transform: translateY(-2px);
}

.grid-item.selected {
  outline: 2px solid #3b82f6;
  background: rgba(59, 130, 246, 0.1);
}

.item-thumbnail {
  aspect-ratio: 1;
  background: rgba(0, 0, 0, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
}

.item-thumbnail img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.thumbnail-placeholder {
  color: #6b7280;
}

.item-info {
  padding: 8px;
}

.item-name {
  font-size: 13px;
  color: #e5e7eb;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.item-meta {
  font-size: 11px;
  color: #9ca3af;
  margin-top: 2px;
}
</style>
