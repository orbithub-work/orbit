<template>
  <div class="file-grid">
    <div
      v-if="loading"
      class="file-grid--loading"
    >
      <span class="loading-spinner"></span>
      <span>加载中...</span>
    </div>
    
    <div
      v-else-if="files.length === 0"
      class="file-grid--empty"
    >
      <FileIcon
        file-type="folder"
        :size="48"
      />
      <p>{{ emptyMessage }}</p>
    </div>
    
    <div
      v-else
      class="grid-items"
    >
      <div 
        v-for="file in sortedFiles" 
        :key="file.id"
        class="grid-item"
        :class="{ 
          'grid-item--selected': selectedIds.has(file.id),
          'grid-item--directory': file.is_directory
        }"
        @click="handleClick(file)"
        @dblclick="handleDoubleClick(file)"
        @contextmenu.prevent="handleContextMenu(file, $event)"
      >
        <div class="item-thumbnail">
          <FileThumbnail 
            v-if="showThumbnails"
            :file-id="file.id"
            :file-name="file.name"
            :file-type="file.file_type || getFileType(file)"
            :thumbnail-size="thumbnailSize"
            :existing-thumbnail="file.thumbnail_path"
            @thumbnail-error="handleThumbnailError"
          />
          <FileIcon 
            v-else
            :file-type="file.is_directory ? 'folder' : (file.file_type || getFileType(file))" 
            :size="thumbnailSize"
          />
        </div>
        <div
          class="item-name"
          :title="file.name"
        >
          {{ file.name }}
        </div>
        <div class="item-meta">
          {{ formatFileSize(file.size) }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import FileIcon from './FileIcon.vue'
import FileThumbnail from './FileThumbnail.vue'
import { useFileStore } from '@/stores/fileStore'

const fileStore = useFileStore()

interface FileItem {
  id: string
  name: string
  path: string
  size: number
  file_type?: string
  is_directory?: boolean
  thumbnail_path?: string | null
}

interface Props {
  files: FileItem[]
  loading?: boolean
  emptyMessage?: string
  showThumbnails?: boolean
  thumbnailSize?: number
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  emptyMessage: '此文件夹为空',
  showThumbnails: true,
  thumbnailSize: 120
})

const emit = defineEmits<{
  select: [files: FileItem[]]
  open: [file: FileItem]
  'context-menu': [file: FileItem, event: MouseEvent]
}>()

const selectedIds = ref<Set<string>>(new Set())

const sortedFiles = computed(() => {
  return fileStore.filteredFiles
})

const getFileType = (file: FileItem): string => {
  if (file.is_directory) return 'folder'
  
  const ext = file.name.split('.').pop()?.toLowerCase() || ''
  const imageExts = ['jpg', 'jpeg', 'png', 'gif', 'webp', 'bmp']
  const videoExts = ['mp4', 'mov', 'avi', 'mkv', 'webm']
  const audioExts = ['mp3', 'wav', 'flac', 'aac', 'ogg', 'm4a']
  const docExts = ['pdf', 'doc', 'docx', 'xls', 'xlsx', 'txt']
  
  if (imageExts.includes(ext)) return 'image'
  if (videoExts.includes(ext)) return 'video'
  if (audioExts.includes(ext)) return 'audio'
  if (docExts.includes(ext)) return 'document'
  
  return 'other'
}

const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const handleClick = (file: FileItem, event?: MouseEvent) => {
  if (event && (event.target as HTMLElement).tagName === 'INPUT') return
  
  if (selectedIds.value.has(file.id)) {
    selectedIds.value.delete(file.id)
  } else {
    selectedIds.value.add(file.id)
  }
  
  emitSelectedFiles()
}

const handleDoubleClick = (file: FileItem) => {
  emit('open', file)
}

const handleContextMenu = (file: FileItem, event: MouseEvent) => {
  if (!selectedIds.value.has(file.id)) {
    selectedIds.value.clear()
    selectedIds.value.add(file.id)
  }
  emit('context-menu', file, event)
}

const handleThumbnailError = () => {
  // Thumbnail loading failed, will show icon instead
}

const emitSelectedFiles = () => {
  const selectedFiles = props.files.filter(f => selectedIds.value.has(f.id))
  emit('select', selectedFiles)
}
</script>

<style scoped>
.file-grid {
  width: 100%;
  height: 100%;
  overflow: auto;
}

.file-grid--loading,
.file-grid--empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  gap: 1rem;
  color: var(--color-text-secondary);
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 3px solid var(--color-border);
  border-top-color: var(--color-primary);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.grid-items {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
  gap: 1rem;
  padding: 1rem;
}

.grid-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 1rem;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.15s;
  text-align: center;
  background-color: var(--color-surface);
  border: 2px solid transparent;
}

.grid-item:hover {
  background-color: var(--color-hover);
}

.grid-item--selected {
  background-color: var(--color-selected);
  border-color: var(--color-primary);
}

.grid-item--directory .item-name {
  font-weight: 500;
}

.item-thumbnail {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 0.75rem;
}

.item-name {
  font-size: 0.875rem;
  margin-bottom: 0.25rem;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  width: 100%;
}

.item-meta {
  font-size: 0.75rem;
  color: var(--color-text-secondary);
}
</style>