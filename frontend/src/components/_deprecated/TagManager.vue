<!-- 标签管理组件 -->
<template>
  <div class="tag-manager">
    <!-- 头部 -->
    <div class="tag-manager-header">
      <h2 class="tag-manager-title">
        标签管理
      </h2>
      <div class="tag-manager-actions">
        <button
          class="btn btn--primary"
          @click="showCreateDialog = true"
        >
          <svg class="icon icon--16"><use href="#i-plus" /></svg>
          新建标签
        </button>
        <button
          class="btn btn--ghost"
          :disabled="loading"
          @click="loadTags"
        >
          <svg class="icon icon--16"><use href="#i-refresh" /></svg>
        </button>
      </div>
    </div>

    <!-- 搜索框 -->
    <div class="tag-manager-search">
      <div class="search-input-wrapper">
        <svg class="icon icon--16"><use href="#i-search" /></svg>
        <input
          v-model="searchQuery"
          type="text"
          class="search-input"
          placeholder="搜索标签..."
          @input="handleSearch"
        />
      </div>
    </div>

    <!-- 标签列表 -->
    <div class="tag-manager-content">
      <!-- 加载状态 -->
      <div
        v-if="loading"
        class="loading-state"
      >
        <div class="loading-spinner"></div>
        <span>加载中...</span>
      </div>

      <!-- 空状态 -->
      <div
        v-else-if="displayedTags.length === 0"
        class="empty-state"
      >
        <svg class="icon icon--48 empty-icon"><use href="#i-label" /></svg>
        <p class="empty-text">
          {{ searchQuery ? '没有找到匹配的标签' : '还没有创建标签' }}
        </p>
        <button
          v-if="!searchQuery"
          class="btn btn--primary"
          @click="showCreateDialog = true"
        >
          创建第一个标签
        </button>
      </div>

      <!-- 标签列表 -->
      <div
        v-else
        class="tag-list"
      >
        <div
          v-for="tag in displayedTags"
          :key="tag.id"
          class="tag-item"
          :class="{ 'tag-item--editing': editingTagId === tag.id }"
        >
          <!-- 标签颜色指示器 -->
          <div
            class="tag-color"
            :style="{ backgroundColor: tag.color || '#6b7280' }"
          ></div>

          <!-- 标签信息 -->
          <div class="tag-info">
            <div class="tag-name">
              {{ tag.name }}
            </div>
            <div class="tag-meta">
              <span class="tag-count">{{ tag.file_count }} 个文件</span>
              <span class="tag-date">{{ formatDate(tag.created_at) }}</span>
            </div>
          </div>

          <!-- 标签操作 -->
          <div class="tag-actions">
            <button
              class="btn btn--icon btn--ghost"
              title="编辑"
              @click="startEdit(tag)"
            >
              <svg class="icon icon--16"><use href="#i-edit" /></svg>
            </button>
            <button
              class="btn btn--icon btn--ghost"
              title="删除"
              @click="confirmDelete(tag)"
            >
              <svg class="icon icon--16"><use href="#i-trash" /></svg>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 创建/编辑对话框 -->
    <div
      v-if="showCreateDialog || showEditDialog"
      class="dialog-overlay"
      @click.self="closeDialogs"
    >
      <div class="dialog">
        <div class="dialog-header">
          <h3 class="dialog-title">
            {{ showCreateDialog ? '创建标签' : '编辑标签' }}
          </h3>
          <button
            class="btn btn--icon btn--ghost"
            @click="closeDialogs"
          >
            <svg class="icon icon--16"><use href="#i-window-close" /></svg>
          </button>
        </div>

        <div class="dialog-body">
          <div class="form-group">
            <label class="form-label">标签名称</label>
            <input
              v-model="tagForm.name"
              type="text"
              class="form-input"
              placeholder="输入标签名称..."
            />
          </div>

          <div class="form-group">
            <label class="form-label">标签颜色</label>
            <div class="color-picker">
              <div
                v-for="color in presetColors"
                :key="color"
                class="color-option"
                :class="{ active: tagForm.color === color }"
                :style="{ backgroundColor: color }"
                @click="tagForm.color = color"
              ></div>
              <input
                v-model="tagForm.color"
                type="color"
                class="color-input"
                title="自定义颜色"
              />
            </div>
          </div>

          <div class="form-group">
            <label class="form-label">图标（可选）</label>
            <div class="icon-picker">
              <div
                v-for="icon in presetIcons"
                :key="icon"
                class="icon-option"
                :class="{ active: tagForm.icon === icon }"
                @click="tagForm.icon = icon"
              >
                <svg class="icon icon--20"><use :href="`#${icon}`" /></svg>
              </div>
              <button
                class="btn btn--icon btn--ghost"
                title="清除图标"
                @click="tagForm.icon = null"
              >
                <svg class="icon icon--16"><use href="#i-clear" /></svg>
              </button>
            </div>
          </div>
        </div>

        <div class="dialog-footer">
          <button
            class="btn btn--ghost"
            @click="closeDialogs"
          >
            取消
          </button>
          <button
            class="btn btn--primary"
            :disabled="!tagForm.name.trim() || saving"
            @click="saveTag"
          >
            {{ saving ? '保存中...' : '保存' }}
          </button>
        </div>
      </div>
    </div>

    <!-- 删除确认对话框 -->
    <div
      v-if="showDeleteConfirm"
      class="dialog-overlay"
      @click.self="showDeleteConfirm = false"
    >
      <div class="dialog dialog--small">
        <div class="dialog-header">
          <h3 class="dialog-title">
            确认删除
          </h3>
        </div>

        <div class="dialog-body">
          <p>确定要删除标签 "{{ tagToDelete?.name }}" 吗？</p>
          <p
            v-if="tagToDelete && tagToDelete.file_count > 0"
            class="warning-text"
          >
            该标签关联了 {{ tagToDelete.file_count }} 个文件，删除后这些文件将不再显示该标签。
          </p>
        </div>

        <div class="dialog-footer">
          <button
            class="btn btn--ghost"
            @click="showDeleteConfirm = false"
          >
            取消
          </button>
          <button
            class="btn btn--danger"
            :disabled="deleting"
            @click="deleteTag"
          >
            {{ deleting ? '删除中...' : '删除' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, reactive } from 'vue'
import { useTagStore, TagItem } from '../stores/tagStore'

// Store
const tagStore = useTagStore()

// 状态
const loading = ref(false)
const searchQuery = ref('')
const displayedTags = ref<TagItem[]>([])
const showCreateDialog = ref(false)
const showEditDialog = ref(false)
const showDeleteConfirm = ref(false)
const editingTagId = ref<string | null>(null)
const tagToDelete = ref<TagItem | null>(null)
const saving = ref(false)
const deleting = ref(false)

// 标签表单
const tagForm = reactive({
  name: '',
  color: '#6b7280',
  icon: null as string | null
})

// 预设颜色
const presetColors = [
  '#6b7280', // 灰色
  '#ef4444', // 红色
  '#f97316', // 橙色
  '#eab308', // 黄色
  '#22c55e', // 绿色
  '#06b6d4', // 青色
  '#3b82f6', // 蓝色
  '#8b5cf6', // 紫色
  '#ec4899', // 粉色
]

// 预设图标
const presetIcons = [
  'i-label',
  'i-star',
  'i-heart',
  'i-bookmark',
  'i-flag',
  'i-folder',
  'i-image',
  'i-video',
  'i-music',
]

// 加载标签
async function loadTags() {
  loading.value = true
  try {
    await tagStore.loadTags()
    displayedTags.value = tagStore.sortedTags
  } finally {
    loading.value = false
  }
}

// 搜索处理
async function handleSearch() {
  if (!searchQuery.value.trim()) {
    displayedTags.value = tagStore.sortedTags
    return
  }

  loading.value = true
  try {
    const results = await tagStore.searchTags(searchQuery.value, 50)
    displayedTags.value = results
  } finally {
    loading.value = false
  }
}

// 开始编辑
function startEdit(tag: TagItem) {
  editingTagId.value = tag.id
  tagForm.name = tag.name
  tagForm.color = tag.color || presetColors[0]
  tagForm.icon = tag.icon
  showEditDialog.value = true
}

// 确认删除
function confirmDelete(tag: TagItem) {
  tagToDelete.value = tag
  showDeleteConfirm.value = true
}

// 保存标签
async function saveTag() {
  saving.value = true
  try {
    if (showCreateDialog.value) {
      await tagStore.createTag({
        name: tagForm.name.trim(),
        color: tagForm.color,
        icon: tagForm.icon,
      })
    } else if (showEditDialog.value && editingTagId.value) {
      await tagStore.updateTag({
        id: editingTagId.value,
        name: tagForm.name.trim(),
        color: tagForm.color,
        icon: tagForm.icon,
      })
    }

    await loadTags()
    closeDialogs()
  } finally {
    saving.value = false
  }
}

// 删除标签
async function deleteTag() {
  if (!tagToDelete.value) return

  deleting.value = true
  try {
    await tagStore.deleteTag(tagToDelete.value.id)
    displayedTags.value = displayedTags.value.filter(t => t.id !== tagToDelete.value!.id)
    showDeleteConfirm.value = false
    tagToDelete.value = null
  } finally {
    deleting.value = false
  }
}

// 关闭对话框
function closeDialogs() {
  showCreateDialog.value = false
  showEditDialog.value = false
  showDeleteConfirm.value = false
  editingTagId.value = null
  tagToDelete.value = null

  // 重置表单
  tagForm.name = ''
  tagForm.color = presetColors[0]
  tagForm.icon = null
}

// 格式化日期
function formatDate(timestamp: number) {
  const date = new Date(timestamp)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))

  if (days === 0) {
    return '今天'
  } else if (days === 1) {
    return '昨天'
  } else if (days < 7) {
    return `${days} 天前`
  } else {
    return date.toLocaleDateString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit' })
  }
}

// 初始化
onMounted(() => {
  loadTags()
})
</script>

<style scoped>
.tag-manager {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: var(--color-bg-base);
}

.tag-manager-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid var(--color-border-light);
}

.tag-manager-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
}

.tag-manager-actions {
  display: flex;
  gap: 8px;
}

.tag-manager-search {
  padding: 12px 20px;
  border-bottom: 1px solid var(--color-border-light);
}

.search-input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.search-input-wrapper .icon {
  position: absolute;
  left: 12px;
  color: var(--color-text-tertiary);
}

.search-input {
  width: 100%;
  padding: 8px 12px 8px 36px;
  font-size: 14px;
  border: 1px solid var(--color-border-light);
  border-radius: 6px;
  background: var(--color-bg-elevated);
  color: var(--color-text-primary);
  outline: none;
  transition: border-color 0.2s;
}

.search-input:focus {
  border-color: var(--color-primary);
}

.tag-manager-content {
  flex: 1;
  overflow-y: auto;
  padding: 8px 0;
}

.loading-state,
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 48px 20px;
  color: var(--color-text-secondary);
}

.empty-icon {
  margin-bottom: 16px;
  color: var(--color-text-tertiary);
}

.empty-text {
  margin: 0 0 16px 0;
  font-size: 14px;
}

.loading-spinner {
  width: 32px;
  height: 32px;
  border: 2px solid var(--color-border-light);
  border-top-color: var(--color-primary);
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.tag-list {
  padding: 0 8px;
}

.tag-item {
  display: flex;
  align-items: center;
  padding: 12px 12px;
  border-radius: 6px;
  transition: background-color 0.2s;
}

.tag-item:hover {
  background-color: var(--color-bg-hover);
}

.tag-item--editing {
  background-color: var(--color-bg-selected);
}

.tag-color {
  width: 24px;
  height: 24px;
  border-radius: 4px;
  flex-shrink: 0;
  margin-right: 12px;
}

.tag-info {
  flex: 1;
  min-width: 0;
}

.tag-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.tag-meta {
  display: flex;
  gap: 12px;
  margin-top: 2px;
  font-size: 12px;
  color: var(--color-text-tertiary);
}

.tag-actions {
  display: flex;
  gap: 4px;
  opacity: 0;
  transition: opacity 0.2s;
}

.tag-item:hover .tag-actions {
  opacity: 1;
}

/* 对话框样式 */
.dialog-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.dialog {
  background: var(--color-bg-elevated);
  border-radius: 8px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
  max-width: 480px;
  width: 90%;
  max-height: 80vh;
  display: flex;
  flex-direction: column;
}

.dialog--small {
  max-width: 400px;
}

.dialog-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid var(--color-border-light);
}

.dialog-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
}

.dialog-body {
  padding: 20px;
  overflow-y: auto;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  padding: 16px 20px;
  border-top: 1px solid var(--color-border-light);
}

.form-group {
  margin-bottom: 16px;
}

.form-group:last-child {
  margin-bottom: 0;
}

.form-label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text-primary);
  margin-bottom: 8px;
}

.form-input {
  width: 100%;
  padding: 8px 12px;
  font-size: 14px;
  border: 1px solid var(--color-border-light);
  border-radius: 6px;
  background: var(--color-bg-base);
  color: var(--color-text-primary);
  outline: none;
  transition: border-color 0.2s;
}

.form-input:focus {
  border-color: var(--color-primary);
}

.color-picker {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.color-option {
  width: 32px;
  height: 32px;
  border-radius: 6px;
  cursor: pointer;
  border: 2px solid transparent;
  transition: transform 0.2s, border-color 0.2s;
}

.color-option:hover {
  transform: scale(1.1);
}

.color-option.active {
  border-color: var(--color-primary);
}

.color-input {
  width: 32px;
  height: 32px;
  padding: 0;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  background: none;
}

.icon-picker {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.icon-option {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid var(--color-border-light);
  border-radius: 6px;
  cursor: pointer;
  color: var(--color-text-secondary);
  transition: all 0.2s;
}

.icon-option:hover {
  background: var(--color-bg-hover);
}

.icon-option.active {
  border-color: var(--color-primary);
  color: var(--color-primary);
  background: var(--color-bg-selected);
}

.warning-text {
  color: var(--color-warning);
  font-size: 14px;
  margin: 12px 0 0 0;
}
</style>
