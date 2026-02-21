<template>
  <aside class="inspector">
    <div class="inspector-header">
      <h3 class="title">
        详细信息
      </h3>
      <button
        class="close-btn"
        @click="$emit('close')"
      >
        <svg
          class="icon icon--16"
          viewBox="0 0 24 24"
          aria-hidden="true"
        >
          <use href="#i-window-close" />
        </svg>
      </button>
    </div>

    <div
      v-if="item"
      class="inspector-content"
    >
      <!-- 预览图 -->
      <div class="preview-section">
        <div class="preview-container">
          <img 
            v-if="item.thumbnail" 
            :src="item.thumbnail" 
            class="preview-image" 
            :alt="item.name"
          />
          <div
            v-else
            class="preview-placeholder"
            :class="item.type"
          >
            <svg
              class="icon icon--32"
              viewBox="0 0 24 24"
              aria-hidden="true"
            >
              <use :href="getTypeIcon(item.type)" />
            </svg>
          </div>
        </div>
      </div>

      <!-- 基本信息 -->
      <div class="info-section">
        <div class="info-row name-row">
          <span class="label">名称</span>
          <div
            class="value editable"
            contenteditable="true"
          >
            {{ item.name }}
          </div>
          <span class="ext">{{ getExtension(item.name) }}</span>
        </div>
        
        <div class="info-grid">
          <div class="info-item">
            <span class="label">格式</span>
            <span class="value badge">{{ item.type.toUpperCase() }}</span>
          </div>
          <div class="info-item">
            <span class="label">大小</span>
            <span class="value">{{ item.size }}</span>
          </div>
          <div class="info-item">
            <span class="label">尺寸</span>
            <span class="value">1920×1080</span>
          </div>
          <div class="info-item">
            <span class="label">创建</span>
            <span class="value">2023/10/24</span>
          </div>
        </div>
      </div>

      <!-- 色板 (Eagle 特色) -->
      <div
        v-if="item.type === 'image'"
        class="palette-section"
      >
        <div class="section-title">
          <svg
            class="icon icon--16"
            viewBox="0 0 24 24"
            aria-hidden="true"
          ><use href="#i-box" /></svg>
          色彩分析
        </div>
        <div class="color-palette">
          <div
            class="color-swatch"
            style="background: #2d3748"
            title="#2D3748"
          ></div>
          <div
            class="color-swatch"
            style="background: #4a5568"
            title="#4A5568"
          ></div>
          <div
            class="color-swatch"
            style="background: #718096"
            title="#718096"
          ></div>
          <div
            class="color-swatch"
            style="background: #63b3ed"
            title="#63B3ED"
          ></div>
          <div
            class="color-swatch"
            style="background: #f6e05e"
            title="#F6E05E"
          ></div>
        </div>
      </div>

      <!-- 标签 -->
      <div class="tags-section">
        <div class="section-title">
          <svg
            class="icon icon--16"
            viewBox="0 0 24 24"
            aria-hidden="true"
          ><use href="#i-label" /></svg>
          标签
        </div>
        <div class="tags-content">
          <!-- 标签选择器 -->
          <TagSelector
            v-model="fileTags"
            :multiple="true"
            :allow-create="true"
            placeholder="添加标签..."
            @change="handleTagsChange"
          />
          <!-- 标签数量显示 -->
          <div class="tags-count">
            {{ fileTags.length }} 个标签
          </div>
        </div>
      </div>

      <!-- 评分 -->
      <div class="rating-section">
        <div class="section-title">
          评分
        </div>
        <div class="rating-stars">
          <svg
            v-for="i in 5"
            :key="i"
            class="icon icon--16 star"
            :class="{ filled: i <= 4 }"
            viewBox="0 0 24 24"
            aria-hidden="true"
          >
            <use :href="i <= 4 ? '#i-star' : '#i-star-outline'" />
          </svg>
        </div>
      </div>

      <!-- 备注 -->
      <div class="notes-section">
        <div class="section-title">
          <svg
            class="icon icon--16"
            viewBox="0 0 24 24"
            aria-hidden="true"
          ><use href="#i-edit" /></svg>
          备注
        </div>
        <textarea
          class="notes-input"
          placeholder="添加备注信息..."
        ></textarea>
      </div>
    </div>
    
    <div
      v-else
      class="empty-inspector"
    >
      <p>选择项目查看详情</p>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import TagSelector from './TagSelector.vue'
import { useTagStore, TagItem } from '@/stores/tagStore'

// Props
interface Props {
  item: any
}

const props = defineProps<Props>()

// Emits
defineEmits<{
  close: []
}>()

// Store
const tagStore = useTagStore()

// 状态
const fileTags = ref<string[]>([])

// 获取类型图标
const getTypeIcon = (type: string) => {
  const map: Record<string, string> = {
    video: '#i-video',
    image: '#i-image',
    audio: '#i-audio',
    project: '#i-file'
  }
  return map[type] || '#i-file'
}

// 获取文件扩展名
const getExtension = (name: string) => {
  return name.includes('.') ? '.' + name.split('.').pop() : ''
}

// 处理标签变化
async function handleTagsChange(tags: TagItem[]) {
  if (!props.item?.id) return

  const currentTagIds = fileTags.value
  const newTagIds = tags.map(t => t.id)

  // 找出需要添加的标签
  const toAdd = newTagIds.filter(id => !currentTagIds.includes(id))
  // 找出需要移除的标签
  const toRemove = currentTagIds.filter(id => !newTagIds.includes(id))

  try {
    if (toAdd.length > 0) {
      await tagStore.addTagsToFiles([props.item.id], toAdd)
    }
    if (toRemove.length > 0) {
      await tagStore.removeTagsFromFiles([props.item.id], toRemove)
    }

    fileTags.value = newTagIds
  } catch (e) {
    console.error('Failed to update tags:', e)
    // 回滚
    fileTags.value = currentTagIds
  }
}

// 监听项目变化，加载标签
watch(() => props.item, async (newItem) => {
  if (newItem?.id) {
    try {
      const tags = await tagStore.getFileTags(newItem.id)
      fileTags.value = tags.map(t => t.id)
    } catch (e) {
      console.error('Failed to load file tags:', e)
      fileTags.value = []
    }
  } else {
    fileTags.value = []
  }
}, { immediate: true })
</script>

<style scoped>
.inspector {
  height: 100%;
  background: var(--color-bg-sidebar);
  border-left: 1px solid var(--color-border);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.inspector-header {
  height: var(--header-height);
  padding: 0 var(--spacing-md);
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid var(--color-border);
}

.title {
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.close-btn {
  background: transparent;
  border: none;
  color: var(--color-text-secondary);
  cursor: pointer;
  padding: 4px;
  border-radius: 4px;
}

.close-btn:hover {
  background: var(--color-bg-surface-hover);
  color: var(--color-text-primary);
}

.inspector-content {
  flex: 1;
  overflow-y: auto;
  padding: var(--spacing-md);
}

.preview-section {
  margin-bottom: var(--spacing-lg);
}

.preview-container {
  width: 100%;
  aspect-ratio: 16/9;
  background: var(--color-bg-base);
  border-radius: var(--radius-md);
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid var(--color-border);
}

.preview-image {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.preview-placeholder {
  color: var(--color-text-tertiary);
}

.info-section {
  margin-bottom: var(--spacing-lg);
}

.info-row {
  display: flex;
  align-items: center;
  margin-bottom: var(--spacing-sm);
}

.name-row {
  background: var(--color-bg-input);
  padding: 8px;
  border-radius: var(--radius-md);
  border: 1px solid transparent;
}

.name-row:hover {
  border-color: var(--color-border);
}

.name-row .value {
  flex: 1;
  font-weight: 500;
  color: var(--color-text-primary);
  outline: none;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.name-row .ext {
  color: var(--color-text-tertiary);
  font-size: 12px;
}

.info-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
  margin-top: 12px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.info-item .label {
  font-size: 11px;
  color: var(--color-text-tertiary);
}

.info-item .value {
  font-size: 13px;
  color: var(--color-text-secondary);
}

.info-item .badge {
  display: inline-block;
  background: var(--color-bg-surface);
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 11px;
  font-family: var(--font-mono);
}

.section-title {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  font-weight: 600;
  color: var(--color-text-tertiary);
  margin-bottom: 10px;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.palette-section, .tags-section, .rating-section, .notes-section {
  margin-bottom: var(--spacing-xl);
  padding-bottom: var(--spacing-lg);
  border-bottom: 1px solid var(--color-border-light);
}

.color-palette {
  display: flex;
  gap: 6px;
}

.color-swatch {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  border: 2px solid var(--color-bg-surface);
  box-shadow: 0 0 0 1px var(--color-border);
  cursor: pointer;
  transition: transform 0.2s;
}

.color-swatch:hover {
  transform: scale(1.1);
  z-index: 1;
}

.tags-list {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.tags-content {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.tags-count {
  font-size: 11px;
  color: var(--color-text-tertiary);
}

.tag {
  background: var(--color-primary-bg);
  color: var(--color-primary);
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.tag:hover {
  background: var(--color-primary);
  color: white;
}

.add-tag-btn {
  background: var(--color-bg-surface);
  border: 1px dashed var(--color-border);
  color: var(--color-text-secondary);
  width: 24px;
  height: 24px;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
}

.add-tag-btn:hover {
  border-color: var(--color-primary);
  color: var(--color-primary);
}

.rating-stars {
  display: flex;
  gap: 2px;
  color: var(--color-text-disabled);
}

.star.filled {
  color: var(--color-warning);
}

.notes-input {
  width: 100%;
  height: 80px;
  background: var(--color-bg-input);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  padding: 8px;
  color: var(--color-text-primary);
  font-size: 13px;
  resize: none;
  transition: all 0.2s;
}

.notes-input:focus {
  outline: none;
  border-color: var(--color-primary);
  background: var(--color-bg-base);
}

.empty-inspector {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-text-tertiary);
  font-size: 13px;
}
</style>
