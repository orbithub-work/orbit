<template>
  <div class="file-list">
    <div
      v-if="loading"
      class="file-list--loading"
    >
      <span class="loading-spinner"></span>
      <span>加载中...</span>
    </div>
    
    <div
      v-else-if="files.length === 0"
      class="file-list--empty"
    >
      <FileIcon
        file-type="folder"
        :size="48"
      />
      <p>{{ emptyMessage }}</p>
    </div>
    
    <table
      v-else
      class="file-list-table"
    >
      <thead>
        <tr>
          <th class="col-checkbox">
            <input 
              type="checkbox" 
              :checked="allSelected" 
              @change="toggleSelectAll"
            />
          </th>
          <th class="col-name">
            名称
            <button
              class="sort-btn"
              @click="sortBy('name')"
            >
              <span v-if="fileStore.sortField === 'name'">{{ fileStore.sortOrder === 'asc' ? '↑' : '↓' }}</span>
            </button>
          </th>
          <th class="col-size">
            大小
            <button
              class="sort-btn"
              @click="sortBy('size')"
            >
              <span v-if="fileStore.sortField === 'size'">{{ fileStore.sortOrder === 'asc' ? '↑' : '↓' }}</span>
            </button>
          </th>
          <th class="col-type">
            类型
            <button
              class="sort-btn"
              @click="sortBy('type')"
            >
              <span v-if="fileStore.sortField === 'type'">{{ fileStore.sortOrder === 'asc' ? '↑' : '↓' }}</span>
            </button>
          </th>
          <th class="col-modified">
            修改时间
            <button
              class="sort-btn"
              @click="sortBy('modified_at')"
            >
              <span v-if="fileStore.sortField === 'modified_at'">{{ fileStore.sortOrder === 'asc' ? '↓' : '↑' }}</span>
            </button>
          </th>
        </tr>
      </thead>
      <tbody>
        <tr 
          v-for="file in sortedFiles" 
          :key="file.id"
          class="file-row"
          :class="{ 
            'file-row--selected': selectedIds.has(file.id),
            'file-row--directory': file.is_directory
          }"
          @click="handleRowClick(file, $event)"
          @dblclick="handleRowDoubleClick(file)"
          @contextmenu.prevent="handleContextMenu(file, $event)"
        >
          <td class="col-checkbox">
            <input 
              type="checkbox" 
              :checked="selectedIds.has(file.id)"
              @change="toggleSelect(file.id, $event)"
            />
          </td>
          <td class="col-name">
            <FileIcon
              :file-type="file.file_type || getFileType(file)"
              :size="20"
            />
            <span class="file-name">{{ file.name }}</span>
          </td>
          <td class="col-size">
            {{ formatFileSize(file.size) }}
          </td>
          <td class="col-type">
            {{ file.file_type || getFileType(file) }}
          </td>
          <td class="col-modified">
            {{ formatDate(file.modified_at) }}
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import FileIcon from './FileIcon.vue'
import { useFileStore } from '@/stores/fileStore'

const fileStore = useFileStore()

interface FileItem {
  id: string
  name: string
  path: string
  size: number
  file_type?: string
  is_directory?: boolean
  modified_at: string | Date
}

interface Props {
  files: FileItem[]
  loading?: boolean
  emptyMessage?: string
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  emptyMessage: '此文件夹为空'
})

const emit = defineEmits<{
  select: [files: FileItem[]]
  open: [file: FileItem]
  'context-menu': [file: FileItem, event: MouseEvent]
  sort: [field: string, order: 'asc' | 'desc']
}>()

const selectedIds = ref<Set<string>>(new Set())

const sortedFiles = computed(() => {
  return fileStore.filteredFiles
})

const allSelected = computed(() => {
  return props.files.length > 0 && selectedIds.value.size === props.files.length
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

const formatDate = (date: string | Date): string => {
  const d = typeof date === 'string' ? new Date(date) : date
  return new Intl.DateTimeFormat('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  }).format(d)
}

const handleRowClick = (file: FileItem, event: MouseEvent) => {
  if ((event.target as HTMLElement).tagName === 'INPUT') return
  
  toggleSelect(file.id)
}

const handleRowDoubleClick = (file: FileItem) => {
  emit('open', file)
}

const handleContextMenu = (file: FileItem, event: MouseEvent) => {
  if (!selectedIds.value.has(file.id)) {
    selectedIds.value.clear()
    selectedIds.value.add(file.id)
  }
  emit('context-menu', file, event)
}

const toggleSelect = (id: string, event?: Event) => {
  if (event) event.stopPropagation()
  
  if (selectedIds.value.has(id)) {
    selectedIds.value.delete(id)
  } else {
    selectedIds.value.add(id)
  }
  
  emitSelectedFiles()
}

const toggleSelectAll = () => {
  if (allSelected.value) {
    selectedIds.value.clear()
  } else {
    selectedIds.value = new Set(props.files.map(f => f.id))
  }
  emitSelectedFiles()
}

const sortBy = (field: string) => {
  fileStore.toggleSort(field)
  emit('sort', field, fileStore.sortOrder)
}

const emitSelectedFiles = () => {
  const selectedFiles = props.files.filter(f => selectedIds.value.has(f.id))
  emit('select', selectedFiles)
}
</script>

<style scoped>
.file-list {
  width: 100%;
  height: 100%;
  overflow: auto;
}

.file-list--loading,
.file-list--empty {
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

.file-list-table {
  width: 100%;
  border-collapse: collapse;
}

.file-list-table thead {
  position: sticky;
  top: 0;
  background-color: var(--color-surface);
  z-index: 1;
}

.file-list-table th {
  padding: 0.75rem 1rem;
  text-align: left;
  font-weight: 600;
  border-bottom: 2px solid var(--color-border);
  user-select: none;
}

.file-list-table th .sort-btn {
  background: none;
  border: none;
  cursor: pointer;
  padding: 0 0.25rem;
  margin-left: 0.25rem;
  color: var(--color-text-secondary);
}

.file-list-table th .sort-btn:hover {
  color: var(--color-primary);
}

.file-row {
  border-bottom: 1px solid var(--color-border);
  cursor: pointer;
  transition: background-color 0.15s;
}

.file-row:hover {
  background-color: var(--color-hover);
}

.file-row--selected {
  background-color: var(--color-selected);
}

.file-row--directory {
  font-weight: 500;
}

.file-row td {
  padding: 0.75rem 1rem;
}

.col-checkbox {
  width: 40px;
  text-align: center;
}

.col-name {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.file-name {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.col-size {
  width: 100px;
  text-align: right;
}

.col-type {
  width: 100px;
}

.col-modified {
  width: 160px;
}
</style>