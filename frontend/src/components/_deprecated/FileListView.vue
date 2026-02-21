<template>
  <div class="file-list-view">
    <!-- 工具栏 -->
    <div class="toolbar">
      <div class="view-controls">
        <button 
          v-for="mode in viewModes" 
          :key="mode.key"
          :class="['view-mode-btn', { active: currentViewMode === mode.key }]"
          :title="mode.title"
          @click="switchViewMode(mode.key)"
        >
          {{ mode.icon }}
        </button>
      </div>
      
      <div class="sort-controls">
        <select
          v-model="sortBy.field"
          class="sort-field"
          @change="applySorting"
        >
          <option value="name">
            名称
          </option>
          <option value="size">
            大小
          </option>
          <option value="modifiedAt">
            修改时间
          </option>
          <option value="createdAt">
            创建时间
          </option>
        </select>
        
        <button 
          class="sort-direction-btn" 
          :title="sortDirection === 'asc' ? '升序' : '降序'"
          @click="toggleSortDirection"
        >
          {{ sortDirection === 'asc' ? '↑' : '↓' }}
        </button>
      </div>
    </div>
    
    <!-- 文件列表容器 -->
    <div
      class="file-list-container"
      :class="`view-mode-${currentViewMode}`"
    >
      <!-- 列表视图 -->
      <div
        v-if="currentViewMode === 'list'"
        class="list-view"
      >
        <div class="list-header">
          <div class="col col-name">
            名称
          </div>
          <div class="col col-size">
            大小
          </div>
          <div class="col col-type">
            类型
          </div>
          <div class="col col-modified">
            修改时间
          </div>
        </div>
        <div 
          v-for="file in sortedFiles" 
          :key="file.id" 
          class="list-row"
          :class="{ selected: file.isSelected }"
          @click="selectFile(file)"
        >
          <div class="col col-name">
            <span class="file-icon">{{ getFileIcon(file) }}</span>
            <span class="file-name">{{ file.name }}</span>
          </div>
          <div class="col col-size">
            {{ formatFileSize(file.size) }}
          </div>
          <div class="col col-type">
            {{ file.mimeType || getFileType(file.name) }}
          </div>
          <div class="col col-modified">
            {{ formatDate(file.modifiedAt) }}
          </div>
        </div>
      </div>
      
      <!-- 网格视图 -->
      <div
        v-else-if="currentViewMode === 'grid'"
        class="grid-view"
      >
        <div 
          v-for="file in sortedFiles" 
          :key="file.id" 
          class="grid-item"
          :class="{ selected: file.isSelected }"
          @click="selectFile(file)"
        >
          <div class="item-icon">
            {{ getFileIcon(file) }}
          </div>
          <div class="item-name">
            {{ file.name }}
          </div>
          <div class="item-meta">
            {{ formatFileSize(file.size) }}
          </div>
        </div>
      </div>
      
      <!-- 缩略图视图 -->
      <div
        v-else-if="currentViewMode === 'thumbnail'"
        class="thumbnail-view"
      >
        <div 
          v-for="file in sortedFiles" 
          :key="file.id" 
          class="thumbnail-item"
          :class="{ selected: file.isSelected }"
          @click="selectFile(file)"
        >
          <div
            v-if="file.thumbnail"
            class="thumbnail-img"
          >
            <img
              :src="file.thumbnail"
              :alt="file.name"
            />
          </div>
          <div
            v-else
            class="thumbnail-placeholder"
          >
            {{ getFileIcon(file) }}
          </div>
          <div class="item-name">
            {{ file.name }}
          </div>
          <div class="item-meta">
            {{ formatFileSize(file.size) }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, reactive } from 'vue'

// 定义文件项接口
interface FileItem {
  id: string;
  name: string;
  path: string;
  size: number;
  type: 'file' | 'directory';
  mimeType?: string;
  thumbnail?: string;
  createdAt: Date;
  modifiedAt: Date;
  isSelected?: boolean;
}

// 定义视图模式
enum ViewMode {
  List = 'list',
  Grid = 'grid',
  Thumbnail = 'thumbnail'
}

// 定义排序字段
enum SortField {
  Name = 'name',
  Size = 'size',
  ModifiedAt = 'modifiedAt',
  CreatedAt = 'createdAt'
}

// 定义排序方向
enum SortDirection {
  Asc = 'asc',
  Desc = 'desc'
}

// 视图模式配置
const viewModes = [
  { key: ViewMode.List, title: '列表视图', icon: '📋' },
  { key: ViewMode.Grid, title: '网格视图', icon: ' squares' }, // 实际使用时会替换为图标
  { key: ViewMode.Thumbnail, title: '缩略图视图', icon: '🖼️' }
]

// 排序配置
interface SortConfig {
  field: SortField;
  direction: SortDirection;
}

// Props
interface Props {
  files: FileItem[];
  initialViewMode?: ViewMode;
  initialSortBy?: SortConfig;
}

const props = withDefaults(defineProps<Props>(), {
  initialViewMode: () => ViewMode.List,
  initialSortBy: () => ({
    field: SortField.ModifiedAt,
    direction: SortDirection.Desc
  })
})

// State
const currentViewMode = ref<ViewMode>(props.initialViewMode)
const sortBy = reactive<SortConfig>({ ...props.initialSortBy })

// Computed properties
const sortedFiles = computed(() => {
  const files = [...props.files] // 创建副本以避免修改原始数据
  
  return files.sort((a, b) => {
    let aValue: any
    let bValue: any
    
    switch (sortBy.field) {
      case SortField.Name:
        aValue = a.name.toLowerCase()
        bValue = b.name.toLowerCase()
        break
      case SortField.Size:
        aValue = a.size
        bValue = b.size
        break
      case SortField.ModifiedAt:
        aValue = new Date(a.modifiedAt).getTime()
        bValue = new Date(b.modifiedAt).getTime()
        break
      case SortField.CreatedAt:
        aValue = new Date(a.createdAt).getTime()
        bValue = new Date(b.createdAt).getTime()
        break
      default:
        aValue = a.name.toLowerCase()
        bValue = b.name.toLowerCase()
    }
    
    if (sortBy.direction === SortDirection.Asc) {
      return aValue > bValue ? 1 : aValue < bValue ? -1 : 0
    } else {
      return aValue < bValue ? 1 : aValue > bValue ? -1 : 0
    }
  })
})

// Methods
const switchViewMode = (mode: ViewMode) => {
  currentViewMode.value = mode
}

const applySorting = () => {
  // 排序已经在computed属性中处理
}

const toggleSortDirection = () => {
  sortBy.direction = sortBy.direction === SortDirection.Asc 
    ? SortDirection.Desc 
    : SortDirection.Asc
}

const selectFile = (file: FileItem) => {
  // 创建新对象以触发响应式更新
  const updatedFile = { ...file, isSelected: !file.isSelected }
  
  // 更新父组件的文件列表
  const index = props.files.findIndex(f => f.id === file.id)
  if (index !== -1) {
    // 注意：这里我们不能直接修改props，所以需要通过emit通知父组件
    emit('file-selected', { file: updatedFile, index })
  }
}

const getFileIcon = (file: FileItem): string => {
  if (file.type === 'directory') {
    return '📁'
  }
  
  const ext = file.name.split('.').pop()?.toLowerCase() || ''
  switch (ext) {
    case 'jpg':
    case 'jpeg':
    case 'png':
    case 'gif':
    case 'webp':
      return '🖼️'
    case 'mp4':
    case 'avi':
    case 'mov':
    case 'mkv':
      return '🎬'
    case 'mp3':
    case 'wav':
    case 'flac':
    case 'aac':
      return '🎵'
    case 'pdf':
      return '📄'
    case 'doc':
    case 'docx':
      return '📝'
    case 'xls':
    case 'xlsx':
      return '📊'
    case 'zip':
    case 'rar':
    case '7z':
      return '📦'
    default:
      return '📄'
  }
}

const getFileType = (fileName: string): string => {
  const ext = fileName.split('.').pop()?.toLowerCase() || ''
  switch (ext) {
    case 'jpg':
    case 'jpeg':
    case 'png':
    case 'gif':
    case 'webp':
      return '图片'
    case 'mp4':
    case 'avi':
    case 'mov':
    case 'mkv':
      return '视频'
    case 'mp3':
    case 'wav':
    case 'flac':
    case 'aac':
      return '音频'
    case 'pdf':
      return 'PDF文档'
    case 'doc':
    case 'docx':
      return 'Word文档'
    case 'xls':
    case 'xlsx':
      return 'Excel文档'
    case 'zip':
    case 'rar':
    case '7z':
      return '压缩文件'
    default:
      return '文件'
  }
}

const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 Bytes'
  
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const formatDate = (date: Date | string): string => {
  const d = typeof date === 'string' ? new Date(date) : date
  return new Intl.DateTimeFormat('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  }).format(d)
}

// Emit
const emit = defineEmits<{
  'file-selected': [payload: { file: FileItem; index: number }]
}>()
</script>

<style scoped>
.file-list-view {
  display: flex;
  flex-direction: column;
  height: 100%;
  width: 100%;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.5rem;
  background-color: var(--color-surface);
  border-bottom: 1px solid var(--color-border);
  gap: 1rem;
}

.view-controls {
  display: flex;
  gap: 0.25rem;
}

.view-mode-btn {
  background: none;
  border: 1px solid var(--color-border);
  border-radius: 4px;
  padding: 0.25rem 0.5rem;
  cursor: pointer;
  font-size: 1rem;
  color: var(--color-text);
}

.view-mode-btn.active {
  background-color: var(--color-primary);
  color: white;
  border-color: var(--color-primary);
}

.sort-controls {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.sort-field {
  padding: 0.25rem 0.5rem;
  border: 1px solid var(--color-border);
  border-radius: 4px;
  background-color: var(--color-surface);
  color: var(--color-text);
}

.sort-direction-btn {
  background: none;
  border: 1px solid var(--color-border);
  border-radius: 4px;
  padding: 0.25rem 0.5rem;
  cursor: pointer;
  font-size: 1rem;
  color: var(--color-text);
}

.file-list-container {
  flex: 1;
  overflow-y: auto;
  padding: 0.5rem;
}

/* 列表视图样式 */
.list-view {
  width: 100%;
}

.list-header {
  display: flex;
  background-color: var(--color-surface);
  border-bottom: 2px solid var(--color-border);
  font-weight: bold;
}

.list-row {
  display: flex;
  padding: 0.5rem;
  border-bottom: 1px solid var(--color-border);
  cursor: pointer;
  transition: background-color 0.2s;
}

.list-row:hover {
  background-color: var(--color-primary-100);
}

.list-row.selected {
  background-color: var(--color-primary-100);
  outline: 2px solid var(--color-primary);
}

.col {
  padding: 0.25rem;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.col-name {
  flex: 3;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.file-icon {
  font-size: 1.2rem;
}

.col-size {
  flex: 1;
  text-align: right;
}

.col-type {
  flex: 1;
}

.col-modified {
  flex: 1.5;
}

/* 网格视图样式 */
.grid-view {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: 1rem;
  padding: 1rem 0;
}

.grid-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 1rem;
  border-radius: 8px;
  cursor: pointer;
  transition: background-color 0.2s;
  text-align: center;
}

.grid-item:hover {
  background-color: var(--color-primary-100);
}

.grid-item.selected {
  background-color: var(--color-primary-100);
  outline: 2px solid var(--color-primary);
}

.item-icon {
  font-size: 2.5rem;
  margin-bottom: 0.5rem;
}

.item-name {
  font-size: 0.9rem;
  margin-bottom: 0.25rem;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  width: 100%;
}

.item-meta {
  font-size: 0.8rem;
  color: var(--color-text-secondary);
}

/* 缩略图视图样式 */
.thumbnail-view {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 1.5rem;
  padding: 1rem 0;
}

.thumbnail-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 1rem;
  border-radius: 8px;
  cursor: pointer;
  transition: background-color 0.2s;
  text-align: center;
}

.thumbnail-item:hover {
  background-color: var(--color-primary-100);
}

.thumbnail-item.selected {
  background-color: var(--color-primary-100);
  outline: 2px solid var(--color-primary);
}

.thumbnail-img {
  width: 150px;
  height: 150px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 0.5rem;
  overflow: hidden;
  border-radius: 4px;
  background-color: var(--color-surface);
}

.thumbnail-img img {
  max-width: 100%;
  max-height: 100%;
  object-fit: cover;
}

.thumbnail-placeholder {
  width: 150px;
  height: 150px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 0.5rem;
  font-size: 3rem;
  background-color: var(--color-surface);
  border-radius: 4px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .grid-view {
    grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
  }
  
  .thumbnail-view {
    grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  }
  
  .list-header, .list-row {
    flex-wrap: wrap;
  }
  
  .col {
    flex: 1 0 100px;
  }
}
</style>