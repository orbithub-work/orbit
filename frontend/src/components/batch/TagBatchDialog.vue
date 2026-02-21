<!-- 批量添加标签对话框 -->
<template>
  <Teleport to="body">
    <Transition name="dialog-fade">
      <div
        v-if="visible"
        class="dialog-overlay"
        @click.self="close"
      >
        <div class="dialog dialog--medium">
          <!-- 对话框头部 -->
          <div class="dialog-header">
            <h3 class="dialog-title">
              为 {{ fileIds.length }} 个文件添加标签
            </h3>
            <button
              class="dialog-close"
              @click="close"
            >
              <svg class="icon">
                <use href="#icon-close" />
              </svg>
            </button>
          </div>

          <!-- 对话框内容 -->
          <div class="dialog-body">
            <!-- 搜索标签 -->
            <div class="search-box">
              <svg class="icon icon--search">
                <use href="#icon-search" />
              </svg>
              <input
                v-model="searchQuery"
                type="text"
                placeholder="搜索或创建标签..."
                class="search-input"
                @keydown.enter="handleCreateTag"
              />
            </div>

            <!-- 标签列表 -->
            <div class="tags-section">
              <h4 class="section-title">
                选择标签
              </h4>
              <div
                v-if="loading"
                class="loading-state"
              >
                <div class="loading-spinner" />
                <span>加载中...</span>
              </div>

              <div
                v-else-if="filteredTags.length === 0"
                class="empty-state"
              >
                <p>{{ searchQuery ? '没有找到匹配的标签' : '暂无标签' }}</p>
                <button
                  v-if="searchQuery"
                  class="btn btn--primary btn--sm"
                  @click="handleCreateTag"
                >
                  创建 "{{ searchQuery }}"
                </button>
              </div>

              <div
                v-else
                class="tags-grid"
              >
                <div
                  v-for="tag in filteredTags"
                  :key="tag.id"
                  class="tag-card"
                  :class="{ 'tag-card--selected': selectedTagIds.has(tag.id) }"
                  @click="toggleTag(tag.id)"
                >
                  <div
                    class="tag-color"
                    :style="{ backgroundColor: tag.color || '#6b7280' }"
                  />
                  <div class="tag-content">
                    <div class="tag-name">
                      {{ tag.name }}
                    </div>
                    <div class="tag-count">
                      {{ tag.file_count }} 个文件
                    </div>
                  </div>
                  <div
                    v-if="selectedTagIds.has(tag.id)"
                    class="tag-check"
                  >
                    <svg class="icon">
                      <use href="#icon-check" />
                    </svg>
                  </div>
                </div>
              </div>
            </div>

            <!-- 已选标签 -->
            <div
              v-if="selectedTagIds.size > 0"
              class="selected-tags"
            >
              <h4 class="section-title">
                已选 {{ selectedTagIds.size }} 个标签
              </h4>
              <div class="selected-tags-list">
                <span
                  v-for="tagId in Array.from(selectedTagIds)"
                  :key="tagId"
                  class="selected-tag"
                  :style="{ '--tag-color': getTagById(tagId)?.color || '#6b7280' }"
                >
                  {{ getTagById(tagId)?.name }}
                  <button
                    class="tag-remove"
                    @click="selectedTagIds.delete(tagId)"
                  >
                    <svg class="icon">
                      <use href="#icon-close" />
                    </svg>
                  </button>
                </span>
              </div>
            </div>
          </div>

          <!-- 对话框底部 -->
          <div class="dialog-footer">
            <button
              class="btn btn--secondary"
              @click="close"
            >
              取消
            </button>
            <button
              class="btn btn--primary"
              :disabled="selectedTagIds.size === 0 || saving"
              @click="handleSave"
            >
              {{ saving ? '保存中...' : `添加到 ${fileIds.length} 个文件` }}
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useTagStore } from '@/stores/tagStore'
import type { TagItem } from '@/stores/tagStore'

interface Props {
  visible: boolean
  fileIds: string[]
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:visible': [value: boolean]
  success: []
}>()

const tagStore = useTagStore()
const searchQuery = ref('')
const selectedTagIds = ref<Set<string>>(new Set())
const loading = ref(false)
const saving = ref(false)

// 过滤标签
const filteredTags = computed(() => {
  if (!searchQuery.value) return tagStore.sortedTags

  const query = searchQuery.value.toLowerCase()
  return tagStore.sortedTags.filter(tag =>
    tag.name.toLowerCase().includes(query)
  )
})

// 根据ID获取标签
function getTagById(id: string): TagItem | undefined {
  return tagStore.getTagById(id)
}

// 切换标签选择
function toggleTag(tagId: string) {
  if (selectedTagIds.value.has(tagId)) {
    selectedTagIds.value.delete(tagId)
  } else {
    selectedTagIds.value.add(tagId)
  }
}

// 创建新标签
async function handleCreateTag() {
  if (!searchQuery.value.trim()) return

  try {
    const newTag = await tagStore.createTag({
      name: searchQuery.value.trim(),
      color: '#3b82f6'
    })

    selectedTagIds.value.add(newTag.id)
    searchQuery.value = ''
  } catch (err) {
    console.error('Failed to create tag:', err)
  }
}

// 保存
async function handleSave() {
  if (selectedTagIds.value.size === 0) return

  saving.value = true
  try {
    await tagStore.addTagsToFiles(props.fileIds, Array.from(selectedTagIds.value))
    emit('success')
    close()
  } catch (err) {
    console.error('Failed to add tags:', err)
  } finally {
    saving.value = false
  }
}

// 关闭对话框
function close() {
  emit('update:visible', false)
  // 重置状态
  setTimeout(() => {
    searchQuery.value = ''
    selectedTagIds.value.clear()
  }, 300)
}

// 加载标签
async function loadTags() {
  loading.value = true
  try {
    await tagStore.loadTags()
  } finally {
    loading.value = false
  }
}

// 监听可见性
watch(() => props.visible, (newVal) => {
  if (newVal) {
    loadTags()
  }
})

onMounted(() => {
  if (props.visible) {
    loadTags()
  }
})
</script>

<style scoped>
.dialog-overlay {
  position: fixed;
  inset: 0;
  z-index: 9999;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.7);
  backdrop-filter: blur(4px);
}

.dialog {
  display: flex;
  flex-direction: column;
  background: var(--color-bg-elevated, #1e1e1e);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.5);
  max-height: 80vh;
}

.dialog--medium {
  width: 90%;
  max-width: 600px;
}

.dialog-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.dialog-title {
  font-size: 18px;
  font-weight: 600;
  color: #fff;
  margin: 0;
}

.dialog-close {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: none;
  border: none;
  border-radius: 6px;
  color: #9ca3af;
  cursor: pointer;
  transition: all 0.2s;
}

.dialog-close:hover {
  background: rgba(255, 255, 255, 0.05);
  color: #fff;
}

.dialog-body {
  flex: 1;
  padding: 24px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.search-box {
  position: relative;
  display: flex;
  align-items: center;
}

.icon--search {
  position: absolute;
  left: 14px;
  width: 18px;
  height: 18px;
  color: #6b7280;
  pointer-events: none;
}

.search-input {
  width: 100%;
  padding: 12px 14px 12px 44px;
  font-size: 14px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  color: #fff;
  outline: none;
  transition: all 0.2s;
}

.search-input:focus {
  border-color: #3b82f6;
  background: rgba(255, 255, 255, 0.08);
}

.tags-section {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.section-title {
  font-size: 13px;
  font-weight: 600;
  color: #9ca3af;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin: 0;
}

.loading-state,
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  gap: 16px;
  color: #6b7280;
}

.loading-spinner {
  width: 32px;
  height: 32px;
  border: 3px solid rgba(255, 255, 255, 0.1);
  border-top-color: #3b82f6;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.tags-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 12px;
  max-height: 300px;
  overflow-y: auto;
  padding: 2px;
}

.tag-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: rgba(255, 255, 255, 0.03);
  border: 2px solid transparent;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.tag-card:hover {
  background: rgba(255, 255, 255, 0.05);
  border-color: rgba(255, 255, 255, 0.1);
}

.tag-card--selected {
  background: rgba(59, 130, 246, 0.15);
  border-color: #3b82f6;
}

.tag-color {
  width: 32px;
  height: 32px;
  border-radius: 6px;
  flex-shrink: 0;
}

.tag-content {
  flex: 1;
  min-width: 0;
}

.tag-name {
  font-size: 14px;
  font-weight: 500;
  color: #fff;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.tag-count {
  font-size: 12px;
  color: #6b7280;
  margin-top: 2px;
}

.tag-check {
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #3b82f6;
  border-radius: 50%;
  color: #fff;
  flex-shrink: 0;
}

.tag-check .icon {
  width: 12px;
  height: 12px;
}

.selected-tags {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding-top: 20px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.selected-tags-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.selected-tag {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 10px 6px 12px;
  font-size: 13px;
  background: var(--tag-color);
  color: #fff;
  border-radius: 16px;
}

.tag-remove {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 16px;
  height: 16px;
  background: rgba(0, 0, 0, 0.2);
  border: none;
  border-radius: 50%;
  color: #fff;
  cursor: pointer;
  transition: background 0.2s;
}

.tag-remove:hover {
  background: rgba(0, 0, 0, 0.4);
}

.tag-remove .icon {
  width: 10px;
  height: 10px;
}

.dialog-footer {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 12px;
  padding: 20px 24px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.btn {
  padding: 10px 20px;
  font-size: 14px;
  font-weight: 500;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.btn--primary {
  background: #3b82f6;
  color: #fff;
}

.btn--primary:hover:not(:disabled) {
  background: #2563eb;
}

.btn--primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn--secondary {
  background: rgba(255, 255, 255, 0.05);
  color: #fff;
}

.btn--secondary:hover {
  background: rgba(255, 255, 255, 0.1);
}

.btn--sm {
  padding: 6px 12px;
  font-size: 13px;
}

/* 动画 */
.dialog-fade-enter-active,
.dialog-fade-leave-active {
  transition: opacity 0.2s;
}

.dialog-fade-enter-active .dialog,
.dialog-fade-leave-active .dialog {
  transition: transform 0.2s, opacity 0.2s;
}

.dialog-fade-enter-from,
.dialog-fade-leave-to {
  opacity: 0;
}

.dialog-fade-enter-from .dialog,
.dialog-fade-leave-to .dialog {
  transform: scale(0.95);
  opacity: 0;
}

/* 滚动条 */
.tags-grid::-webkit-scrollbar,
.dialog-body::-webkit-scrollbar {
  width: 8px;
}

.tags-grid::-webkit-scrollbar-track,
.dialog-body::-webkit-scrollbar-track {
  background: transparent;
}

.tags-grid::-webkit-scrollbar-thumb,
.dialog-body::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 4px;
}

.tags-grid::-webkit-scrollbar-thumb:hover,
.dialog-body::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.2);
}
</style>
