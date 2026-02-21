<template>
  <div class="search-panel">
    <div class="search-header">
      <h3>高级搜索</h3>
      <button
        class="clear-btn"
        @click="clearAll"
      >
        清空
      </button>
    </div>

    <div class="search-form">
      <!-- 文件名搜索 -->
      <div class="form-group">
        <label for="name">文件名</label>
        <input
          id="name"
          v-model="searchFilters.namePattern"
          type="text"
          placeholder="搜索文件名..."
          class="form-input"
        />
      </div>

      <!-- 文件类型筛选 -->
      <div class="form-group">
        <label>文件类型</label>
        <div class="checkbox-group">
          <label
            v-for="type in fileTypes"
            :key="type.value"
            class="checkbox-label"
          >
            <input
              v-model="searchFilters.fileTypes"
              :value="type.value"
              type="checkbox"
            />
            {{ type.label }}
          </label>
        </div>
      </div>

      <!-- 文件大小筛选 -->
      <div class="form-group">
        <label>文件大小</label>
        <div class="size-range">
          <input
            v-model.number="searchFilters.minSize"
            type="number"
            placeholder="最小"
            class="size-input"
          />
          <span class="size-separator">-</span>
          <input
            v-model.number="searchFilters.maxSize"
            type="number"
            placeholder="最大"
            class="size-input"
          />
          <select
            v-model="searchFilters.sizeUnit"
            class="size-unit"
          >
            <option value="B">
              B
            </option>
            <option value="KB">
              KB
            </option>
            <option value="MB">
              MB
            </option>
            <option value="GB">
              GB
            </option>
          </select>
        </div>
      </div>

      <!-- 日期范围筛选 -->
      <div class="form-group">
        <label>修改日期</label>
        <div class="date-range">
          <input
            v-model="searchFilters.dateFrom"
            type="date"
            class="date-input"
          />
          <span class="date-separator">-</span>
          <input
            v-model="searchFilters.dateTo"
            type="date"
            class="date-input"
          />
        </div>
      </div>

      <!-- 标签筛选 -->
      <div class="form-group">
        <label for="tags">标签</label>
        <input
          id="tags"
          v-model="searchFilters.tags"
          type="text"
          placeholder="输入标签（逗号分隔）"
          class="form-input"
        />
      </div>

      <!-- 搜索按钮 -->
      <div class="form-actions">
        <button
          class="btn-primary"
          @click="performSearch"
        >
          搜索
        </button>
        <button
          class="btn-secondary"
          @click="cancelSearch"
        >
          取消
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useFileStore } from '@/stores/fileStore'

const fileStore = useFileStore()

const fileTypes = [
  { value: 'image', label: '图片' },
  { value: 'video', label: '视频' },
  { value: 'audio', label: '音频' },
  { value: 'document', label: '文档' },
  { value: 'other', label: '其他' },
]

interface SearchFilters {
  namePattern: string
  fileTypes: string[]
  minSize: number | null
  maxSize: number | null
  sizeUnit: 'B' | 'KB' | 'MB' | 'GB'
  dateFrom: string
  dateTo: string
  tags: string
}

const searchFilters = ref<SearchFilters>({
  namePattern: '',
  fileTypes: [],
  minSize: null,
  maxSize: null,
  sizeUnit: 'KB',
  dateFrom: '',
  dateTo: '',
  tags: '',
})

const clearAll = () => {
  searchFilters.value = {
    namePattern: '',
    fileTypes: [],
    minSize: null,
    maxSize: null,
    sizeUnit: 'KB',
    dateFrom: '',
    dateTo: '',
    tags: '',
  }
}

const performSearch = () => {
  // 转换文件大小到字节
  let minSizeBytes: number | null = null
  let maxSizeBytes: number | null = null

  if (searchFilters.value.minSize !== null) {
    minSizeBytes = convertToBytes(searchFilters.value.minSize, searchFilters.value.sizeUnit)
  }

  if (searchFilters.value.maxSize !== null) {
    maxSizeBytes = convertToBytes(searchFilters.value.maxSize, searchFilters.value.sizeUnit)
  }

  // 解析标签
  const tags = searchFilters.value.tags
    .split(',')
    .map(tag => tag.trim())
    .filter(tag => tag.length > 0)

  // 执行搜索
  fileStore.searchFilesAdvanced({
    namePattern: searchFilters.value.namePattern || undefined,
    fileTypes: searchFilters.value.fileTypes.length > 0 ? searchFilters.value.fileTypes : undefined,
    minSize: minSizeBytes,
    maxSize: maxSizeBytes,
    dateFrom: searchFilters.value.dateFrom || undefined,
    dateTo: searchFilters.value.dateTo || undefined,
    tags,
  })
}

const cancelSearch = () => {
  // 清除搜索
  fileStore.searchFiles('')
}

const convertToBytes = (size: number, unit: string): number => {
  const units = {
    B: 1,
    KB: 1024,
    MB: 1024 * 1024,
    GB: 1024 * 1024 * 1024,
  }

  return size * units[unit]
}
</script>

<style scoped>
.search-panel {
  background-color: var(--color-surface);
  border-radius: 8px;
  padding: 1rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.search-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 1rem;
}

.search-header h3 {
  margin: 0;
  font-size: 1.125rem;
  font-weight: 600;
}

.clear-btn {
  background: none;
  border: none;
  color: var(--color-primary);
  cursor: pointer;
  font-size: 0.875rem;
}

.clear-btn:hover {
  text-decoration: underline;
}

.search-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-group label {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--color-text-secondary);
}

.form-input,
.size-input,
.date-input,
.size-unit {
  padding: 0.5rem;
  border: 1px solid var(--color-border);
  border-radius: 4px;
  background-color: var(--color-surface);
  color: var(--color-text);
  font-size: 0.875rem;
}

.form-input:focus,
.size-input:focus,
.date-input:focus,
.size-unit:focus {
  outline: none;
  border-color: var(--color-primary);
}

.checkbox-group {
  display: flex;
  flex-wrap: wrap;
  gap: 0.75rem;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  cursor: pointer;
  font-size: 0.875rem;
}

.size-range,
.date-range {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.size-input {
  width: 100px;
}

.size-separator,
.date-separator {
  color: var(--color-text-secondary);
}

.size-unit {
  width: 80px;
}

.form-actions {
  display: flex;
  gap: 0.75rem;
  justify-content: flex-end;
  margin-top: 0.5rem;
}

.btn-primary,
.btn-secondary {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.875rem;
  transition: background-color 0.15s;
}

.btn-primary {
  background-color: var(--color-primary);
  color: white;
}

.btn-primary:hover {
  background-color: var(--color-primary-dark);
}

.btn-secondary {
  background-color: var(--color-surface);
  color: var(--color-text);
  border: 1px solid var(--color-border);
}

.btn-secondary:hover {
  background-color: var(--color-hover);
}

@media (max-width: 768px) {
  .size-range,
  .date-range {
    flex-direction: column;
    align-items: stretch;
  }

  .size-input {
    width: 100%;
  }

  .size-unit {
    width: 100%;
  }

  .form-actions {
    flex-direction: column;
  }

  .btn-primary,
  .btn-secondary {
    width: 100%;
  }
}
</style>