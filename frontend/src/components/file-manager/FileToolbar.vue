<template>
  <div class="file-toolbar">
    <div class="toolbar-left">
      <button
        class="toolbar-btn toolbar-btn--icon"
        :disabled="!canGoBack"
        :title="'返回上级'"
        @click="$emit('go-back')"
      >
        <svg
          viewBox="0 0 24 24"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path d="M20 11H7.83l5.59-5.59L12 4l-8 8 8 8 1.41-1.41L7.83 13H20v-2z" />
        </svg>
      </button>

      <button
        class="toolbar-btn toolbar-btn--icon"
        :disabled="!canGoForward"
        :title="'前进'"
        @click="$emit('go-forward')"
      >
        <svg
          viewBox="0 0 24 24"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path d="M12 4l-1.41 1.41L16.17 11H4v2h12.17l-5.58 5.59L12 20l8-8z" />
        </svg>
      </button>

      <button
        class="toolbar-btn toolbar-btn--icon"
        :title="'刷新'"
        @click="$emit('refresh')"
      >
        <svg
          viewBox="0 0 24 24"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path d="M17.65 6.35C16.2 4.9 14.21 4 12 4c-4.42 0-7.99 3.58-7.99 8s3.57 8 7.99 8c3.73 0 6.84-2.55 7.73-6h-2.08c-.82 2.33-3.04 4-5.65 4-3.31 0-6-2.69-6-6s2.69-6 6-6c1.66 0 3.14.69 4.22 1.78L13 11h7V4l-2.35 2.35z" />
        </svg>
      </button>

      <button
        class="toolbar-btn toolbar-btn--icon"
        :title="'新建文件夹'"
        @click="$emit('create-folder')"
      >
        <svg
          viewBox="0 0 24 24"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path d="M10 4H4c-1.1 0-1.99.9-1.99 2L2 18c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V8c0-1.1-.9-2-2-2h-8l-2-2z" />
        </svg>
      </button>
    </div>
    
    <div class="toolbar-center">
      <div class="search-box">
        <svg
          class="search-icon"
          viewBox="0 0 24 24"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path d="M15.5 14h-.79l-.28-.27C15.41 12.59 16 11.11 16 9.5 16 5.91 13.09 3 9.5 3S3 5.91 3 9.5 5.91 16 9.5 16c1.61 0 3.09-.59 4.23-1.57l.27.28v.79l5 4.99L20.49 19l-4.99-5zm-6 0C7.01 14 5 11.99 5 9.5S7.01 5 9.5 5 14 7.01 14 9.5 11.99 14 9.5 14z" />
        </svg>
        <input 
          v-model="searchQuery" 
          type="text"
          placeholder="搜索文件..."
          class="search-input"
          @input="handleSearchInput"
          @keyup.enter="handleSearch"
        />
        <button 
          v-if="searchQuery"
          class="search-clear"
          @click="clearSearch"
        >
          <svg
            viewBox="0 0 24 24"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z" />
          </svg>
        </button>
      </div>
    </div>
    
    <div class="toolbar-right">
      <button
        class="toolbar-btn toolbar-btn--icon"
        :title="'高级搜索'"
        @click="$emit('advanced-search')"
      >
        <svg
          viewBox="0 0 24 24"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path d="M15.5 14h-.79l-.28-.27C15.41 12.59 16 11.11 16 9.5 16 5.91 13.09 3 9.5 3S3 5.91 3 9.5 5.91 16 9.5 16c1.61 0 3.09-.59 4.23-1.57l.27.28v.79l5 4.99L20.49 19l-4.99-5zm-6 0C7.01 14 5 11.99 5 9.5S7.01 5 9.5 5 14 7.01 14 9.5 11.99 14 9.5 14z" />
        </svg>
      </button>

      <div class="filter-dropdown">
        <select
          v-model="filterType"
          class="filter-select"
          @change="handleFilterChange"
        >
          <option value="">
            所有类型
          </option>
          <option value="image">
            图片
          </option>
          <option value="video">
            视频
          </option>
          <option value="audio">
            音频
          </option>
          <option value="document">
            文档
          </option>
          <option value="folder">
            文件夹
          </option>
        </select>
      </div>

      <div class="view-toggle">
        <button
          v-for="mode in viewModes"
          :key="mode.value"
          class="view-btn"
          :class="{ 'view-btn--active': currentViewMode === mode.value }"
          :title="mode.label"
          @click="$emit('view-change', mode.value)"
        >
          <svg
            viewBox="0 0 24 24"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path :d="mode.icon" />
          </svg>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

interface Props {
  canGoBack?: boolean
  canGoForward?: boolean
  currentViewMode?: 'list' | 'grid'
}

const props = withDefaults(defineProps<Props>(), {
  canGoBack: false,
  canGoForward: false,
  currentViewMode: 'list'
})

const emit = defineEmits<{
  search: [query: string]
  'filter-change': [type: string]
  'view-change': [mode: 'list' | 'grid']
  'go-back': []
  'go-forward': []
  refresh: []
  'create-folder': []
  'advanced-search': []
}>()

const searchQuery = ref('')
const filterType = ref('')
let searchDebounceTimer: number | null = null

const viewModes = [
  { value: 'list', label: '列表视图', icon: 'M3 13h2v-2H3v2zm0 4h2v-2H3v2zm0-8h2V7H3v2zm4 4h14v-2H7v2zm0 4h14v-2H7v2zM7 7v2h14V7H7z' },
  { value: 'grid', label: '网格视图', icon: 'M3 3v8h8V3H3zm2 6V5h4v4H5zm8-6v8h8V3h-8zm2 6V5h4v4h-4zM3 13v8h8v-8H3zm2 6v-4h4v4H5zm8-6v8h8v-8h-8zm2 6v-4h4v4h-4z' }
]

const handleSearchInput = () => {
  if (searchDebounceTimer) {
    clearTimeout(searchDebounceTimer)
  }
  
  searchDebounceTimer = window.setTimeout(() => {
    handleSearch()
  }, 300)
}

const handleSearch = () => {
  emit('search', searchQuery.value.trim())
}

const clearSearch = () => {
  searchQuery.value = ''
  emit('search', '')
}

const handleFilterChange = () => {
  emit('filter-change', filterType.value)
}

onUnmounted(() => {
  if (searchDebounceTimer) {
    clearTimeout(searchDebounceTimer)
  }
})
</script>

<style scoped>
.file-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  padding: 0.75rem 1rem;
  background-color: var(--color-surface);
  border-bottom: 1px solid var(--color-border);
}

.toolbar-left,
.toolbar-center,
.toolbar-right {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.toolbar-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0.5rem;
  background: none;
  border: 1px solid transparent;
  border-radius: 4px;
  cursor: pointer;
  color: var(--color-text);
  transition: all 0.15s;
}

.toolbar-btn:hover:not(:disabled) {
  background-color: var(--color-hover);
  color: var(--color-primary);
}

.toolbar-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.toolbar-btn--icon svg {
  width: 20px;
  height: 20px;
  fill: currentColor;
}

.search-box {
  display: flex;
  align-items: center;
  position: relative;
  width: 300px;
}

.search-icon {
  position: absolute;
  left: 0.75rem;
  width: 18px;
  height: 18px;
  fill: var(--color-text-secondary);
  pointer-events: none;
}

.search-input {
  width: 100%;
  padding: 0.5rem 2.5rem 0.5rem 2.25rem;
  border: 1px solid var(--color-border);
  border-radius: 4px;
  background-color: var(--color-surface);
  color: var(--color-text);
  font-size: 0.875rem;
}

.search-input:focus {
  outline: none;
  border-color: var(--color-primary);
}

.search-clear {
  position: absolute;
  right: 0.5rem;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0.25rem;
  background: none;
  border: none;
  cursor: pointer;
  color: var(--color-text-secondary);
}

.search-clear:hover {
  color: var(--color-text);
}

.search-clear svg {
  width: 16px;
  height: 16px;
  fill: currentColor;
}

.filter-select {
  padding: 0.5rem 2rem 0.5rem 0.75rem;
  border: 1px solid var(--color-border);
  border-radius: 4px;
  background-color: var(--color-surface);
  color: var(--color-text);
  font-size: 0.875rem;
  cursor: pointer;
  appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath d='M7 10l5 5 5-5z' fill='%23666'/%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 0.5rem center;
  background-size: 16px;
}

.filter-select:focus {
  outline: none;
  border-color: var(--color-primary);
}

.view-toggle {
  display: flex;
  border: 1px solid var(--color-border);
  border-radius: 4px;
  overflow: hidden;
}

.view-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0.5rem 0.75rem;
  background: none;
  border: none;
  cursor: pointer;
  color: var(--color-text-secondary);
  transition: all 0.15s;
}

.view-btn:hover {
  background-color: var(--color-hover);
  color: var(--color-text);
}

.view-btn--active {
  background-color: var(--color-primary);
  color: white;
}

.view-btn svg {
  width: 18px;
  height: 18px;
  fill: currentColor;
}

@media (max-width: 768px) {
  .file-toolbar {
    flex-wrap: wrap;
  }
  
  .search-box {
    width: 200px;
  }
  
  .filter-select {
    font-size: 0.8rem;
  }
}
</style>