<!-- 标签选择器组件 -->
<template>
  <div class="tag-selector">
    <div
      v-if="label"
      class="tag-selector-label"
    >
      {{ label }}
    </div>

    <div
      class="tag-selector-container"
      :class="{
        'tag-selector--open': isOpen,
        'tag-selector--disabled': disabled,
        'tag-selector--error': error
      }"
    >
      <!-- 选中的标签 -->
      <div
        class="tag-selector-display"
        @click="openDropdown"
      >
        <div class="tag-selector-tags">
          <div
            v-for="tag in selectedTags"
            :key="tag.id"
            class="tag-chip"
          >
            <div
              class="tag-chip-color"
              :style="{ backgroundColor: tag.color || '#6b7280' }"
            ></div>
            <span class="tag-chip-text">{{ tag.name }}</span>
            <button
              v-if="!disabled"
              class="tag-chip-remove"
              @click.stop="removeTag(tag.id)"
            >
              <svg class="icon icon--12"><use href="#i-close" /></svg>
            </button>
          </div>
          <span
            v-if="selectedTags.length === 0"
            class="tag-selector-placeholder"
          >
            {{ placeholder }}
          </span>
        </div>
        <svg class="icon icon--16 tag-selector-arrow">
          <use :href="isOpen ? '#i-chevron-up' : '#i-chevron-down'" />
        </svg>
      </div>

      <!-- 下拉菜单 -->
      <div
        v-if="isOpen"
        class="tag-dropdown"
      >
        <!-- 搜索框 -->
        <div class="tag-dropdown-search">
          <svg class="icon icon--14"><use href="#i-search" /></svg>
          <input
            ref="searchInputRef"
            v-model="searchQuery"
            type="text"
            class="tag-dropdown-input"
            :placeholder="searchPlaceholder"
            @input="handleSearch"
          />
        </div>

        <!-- 创建新标签按钮 -->
        <button
          v-if="allowCreate && searchQuery.trim() && !tagExists(searchQuery.trim())"
          class="tag-dropdown-create"
          @click="createNewTag"
        >
          <svg class="icon icon--14"><use href="#i-plus" /></svg>
          创建 "{{ searchQuery.trim() }}" 标签
        </button>

        <!-- 标签列表 -->
        <div class="tag-dropdown-list">
          <!-- 加载状态 -->
          <div
            v-if="loading"
            class="tag-dropdown-loading"
          >
            <div class="loading-spinner"></div>
          </div>

          <!-- 空状态 -->
          <div
            v-else-if="filteredTags.length === 0"
            class="tag-dropdown-empty"
          >
            <svg class="icon icon--32"><use href="#i-label" /></svg>
            <span class="empty-text">
              {{ searchQuery ? '没有找到匹配的标签' : '还没有标签' }}
            </span>
          </div>

          <!-- 标签选项 -->
          <div
            v-else
            class="tag-dropdown-options"
          >
            <div
              v-for="tag in filteredTags"
              :key="tag.id"
              class="tag-dropdown-option"
              :class="{
                'tag-dropdown-option--selected': isTagSelected(tag.id),
                'tag-dropdown-option--focused': focusedTagId === tag.id
              }"
              @click="toggleTag(tag)"
              @mouseenter="focusedTagId = tag.id"
            >
              <div
                class="tag-option-color"
                :style="{ backgroundColor: tag.color || '#6b7280' }"
              ></div>
              <span class="tag-option-name">{{ tag.name }}</span>
              <span class="tag-option-count">{{ tag.file_count }}</span>
              <svg
                v-if="isTagSelected(tag.id)"
                class="icon icon--16 tag-option-check"
              >
                <use href="#i-check" />
              </svg>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 错误信息 -->
    <div
      v-if="error"
      class="tag-selector-error"
    >
      {{ error }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick, onMounted, onUnmounted } from 'vue'
import { useTagStore, TagItem } from '../stores/tagStore'

// Props
interface Props {
  modelValue: string[] // 选中的标签ID列表
  label?: string
  placeholder?: string
  searchPlaceholder?: string
  disabled?: boolean
  error?: string
  multiple?: boolean
  allowCreate?: boolean
  max?: number
}

const props = withDefaults(defineProps<Props>(), {
  label: '',
  placeholder: '选择标签...',
  searchPlaceholder: '搜索或创建标签...',
  disabled: false,
  error: '',
  multiple: true,
  allowCreate: true,
  max: undefined
})

// Emits
const emit = defineEmits<{
  'update:modelValue': [value: string[]]
  'change': [tags: TagItem[]]
}>()

// Store
const tagStore = useTagStore()

// 状态
const isOpen = ref(false)
const loading = ref(false)
const searchQuery = ref('')
const searchInputRef = ref<HTMLInputElement | null>(null)
const focusedTagId = ref<string | null>(null)
const availableTags = ref<TagItem[]>([])

// 选中的标签
const selectedTags = computed(() => {
  return props.modelValue.map(id => tagStore.getTagById(id)).filter((t): t is TagItem => t !== undefined)
})

// 过滤后的标签
const filteredTags = computed(() => {
  if (!searchQuery.value.trim()) {
    return availableTags.value
  }

  const query = searchQuery.value.toLowerCase()
  return availableTags.value.filter(tag =>
    tag.name.toLowerCase().includes(query)
  )
})

// 检查标签是否已存在
function tagExists(name: string) {
  return availableTags.value.some(tag =>
    tag.name.toLowerCase() === name.toLowerCase()
  )
}

// 检查标签是否被选中
function isTagSelected(tagId: string) {
  return props.modelValue.includes(tagId)
}

// 打开下拉菜单
function openDropdown() {
  if (props.disabled) return

  isOpen.value = true

  // 加载标签
  if (availableTags.value.length === 0) {
    loadTags()
  }

  // 聚焦搜索框
  nextTick(() => {
    searchInputRef.value?.focus()
  })
}

// 关闭下拉菜单
function closeDropdown() {
  isOpen.value = false
  searchQuery.value = ''
  focusedTagId.value = null
}

// 加载标签
async function loadTags() {
  loading.value = true
  try {
    await tagStore.loadTags()
    availableTags.value = tagStore.sortedTags
  } finally {
    loading.value = false
  }
}

// 搜索处理
async function handleSearch() {
  if (!searchQuery.value.trim()) {
    availableTags.value = tagStore.sortedTags
    return
  }

  loading.value = true
  try {
    const results = await tagStore.searchTags(searchQuery.value, 50)
    availableTags.value = results
  } finally {
    loading.value = false
  }
}

// 切换标签选中状态
function toggleTag(tag: TagItem) {
  const isSelected = isTagSelected(tag.id)

  if (props.multiple) {
    const newValue = isSelected
      ? props.modelValue.filter(id => id !== tag.id)
      : (props.max && props.modelValue.length >= props.max
          ? props.modelValue
          : [...props.modelValue, tag.id])

    emit('update:modelValue', newValue)
    emit('change', newValue.map(id => tagStore.getTagById(id)!).filter(t => t))
  } else {
    const newValue = isSelected ? [] : [tag.id]
    emit('update:modelValue', newValue)
    emit('change', newValue.map(id => tagStore.getTagById(id)!).filter(t => t))
    closeDropdown()
  }
}

// 移除标签
function removeTag(tagId: string) {
  if (props.disabled) return

  const newValue = props.modelValue.filter(id => id !== tagId)
  emit('update:modelValue', newValue)
  emit('change', newValue.map(id => tagStore.getTagById(id)!).filter(t => t))
}

// 创建新标签
async function createNewTag() {
  const name = searchQuery.value.trim()
  if (!name) return

  try {
    const newTag = await tagStore.createTag({
      name,
      color: '#6b7280',
    })

    availableTags.value.push(newTag)

    if (props.multiple) {
      const newValue = [...props.modelValue, newTag.id]
      emit('update:modelValue', newValue)
      emit('change', newValue.map(id => tagStore.getTagById(id)!).filter(t => t))
    } else {
      emit('update:modelValue', [newTag.id])
      emit('change', [newTag])
      closeDropdown()
    }

    searchQuery.value = ''
  } catch (e) {
    console.error('Failed to create tag:', e)
  }
}

// 键盘导航
function handleKeydown(event: KeyboardEvent) {
  if (!isOpen.value) return

  switch (event.key) {
    case 'Escape':
      event.preventDefault()
      closeDropdown()
      break
    case 'ArrowDown':
      event.preventDefault()
      focusNextTag()
      break
    case 'ArrowUp':
      event.preventDefault()
      focusPreviousTag()
      break
    case 'Enter':
      event.preventDefault()
      if (focusedTagId.value) {
        const tag = filteredTags.value.find(t => t.id === focusedTagId.value)
        if (tag) toggleTag(tag)
      }
      break
    case 'Tab':
      closeDropdown()
      break
  }
}

function focusNextTag() {
  if (filteredTags.value.length === 0) return

  const currentIndex = filteredTags.value.findIndex(t => t.id === focusedTagId.value)
  const nextIndex = currentIndex < filteredTags.value.length - 1 ? currentIndex + 1 : 0
  focusedTagId.value = filteredTags.value[nextIndex].id
}

function focusPreviousTag() {
  if (filteredTags.value.length === 0) return

  const currentIndex = filteredTags.value.findIndex(t => t.id === focusedTagId.value)
  const prevIndex = currentIndex > 0 ? currentIndex - 1 : filteredTags.value.length - 1
  focusedTagId.value = filteredTags.value[prevIndex].id
}

// 点击外部关闭
function handleClickOutside(event: MouseEvent) {
  const target = event.target as HTMLElement
  if (!target.closest('.tag-selector')) {
    closeDropdown()
  }
}

// 监听模型值变化
watch(() => props.modelValue, (newValue) => {
  // 如果有新的标签ID，确保标签已加载
  if (newValue.length > 0 && availableTags.value.length === 0) {
    loadTags()
  }
}, { immediate: true })

// 生命周期
onMounted(() => {
  document.addEventListener('keydown', handleKeydown)
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeydown)
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
.tag-selector {
  position: relative;
}

.tag-selector-label {
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text-primary);
  margin-bottom: 8px;
}

.tag-selector-container {
  position: relative;
}

.tag-selector-display {
  display: flex;
  align-items: center;
  min-height: 36px;
  padding: 4px 8px 4px 12px;
  background: var(--color-bg-elevated);
  border: 1px solid var(--color-border-light);
  border-radius: 6px;
  cursor: pointer;
  transition: border-color 0.2s;
}

.tag-selector-display:hover {
  border-color: var(--color-border-medium);
}

.tag-selector--open .tag-selector-display {
  border-color: var(--color-primary);
}

.tag-selector--disabled .tag-selector-display {
  cursor: not-allowed;
  opacity: 0.6;
}

.tag-selector--error .tag-selector-display {
  border-color: var(--color-error);
}

.tag-selector-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  flex: 1;
  align-items: center;
}

.tag-chip {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 2px 6px 2px 8px;
  background: var(--color-bg-hover);
  border-radius: 4px;
  font-size: 13px;
}

.tag-chip-color {
  width: 10px;
  height: 10px;
  border-radius: 2px;
  flex-shrink: 0;
}

.tag-chip-text {
  color: var(--color-text-primary);
}

.tag-chip-remove {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 14px;
  height: 14px;
  padding: 0;
  background: none;
  border: none;
  border-radius: 2px;
  cursor: pointer;
  color: var(--color-text-tertiary);
  transition: all 0.2s;
}

.tag-chip-remove:hover {
  background: var(--color-bg-selected);
  color: var(--color-text-primary);
}

.tag-selector-placeholder {
  color: var(--color-text-tertiary);
  font-size: 14px;
}

.tag-selector-arrow {
  color: var(--color-text-tertiary);
  flex-shrink: 0;
  margin-left: 8px;
}

.tag-dropdown {
  position: absolute;
  top: calc(100% + 4px);
  left: 0;
  right: 0;
  background: var(--color-bg-elevated);
  border: 1px solid var(--color-border-light);
  border-radius: 6px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.2);
  z-index: 100;
  max-height: 320px;
  overflow: hidden;
}

.tag-dropdown-search {
  display: flex;
  align-items: center;
  padding: 8px 12px;
  border-bottom: 1px solid var(--color-border-light);
}

.tag-dropdown-search .icon {
  color: var(--color-text-tertiary);
  margin-right: 8px;
}

.tag-dropdown-input {
  flex: 1;
  padding: 6px 0;
  font-size: 14px;
  background: none;
  border: none;
  color: var(--color-text-primary);
  outline: none;
}

.tag-dropdown-input::placeholder {
  color: var(--color-text-tertiary);
}

.tag-dropdown-create {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
  padding: 10px 12px;
  background: none;
  border: none;
  border-bottom: 1px solid var(--color-border-light);
  color: var(--color-primary);
  font-size: 14px;
  cursor: pointer;
  transition: background 0.2s;
}

.tag-dropdown-create:hover {
  background: var(--color-bg-hover);
}

.tag-dropdown-list {
  max-height: 240px;
  overflow-y: auto;
}

.tag-dropdown-loading {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
}

.loading-spinner {
  width: 24px;
  height: 24px;
  border: 2px solid var(--color-border-light);
  border-top-color: var(--color-primary);
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.tag-dropdown-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 32px 20px;
  color: var(--color-text-tertiary);
}

.tag-dropdown-empty .icon {
  margin-bottom: 8px;
}

.empty-text {
  font-size: 13px;
}

.tag-dropdown-options {
  padding: 4px 0;
}

.tag-dropdown-option {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 12px;
  cursor: pointer;
  transition: background 0.2s;
}

.tag-dropdown-option:hover,
.tag-dropdown-option--focused {
  background: var(--color-bg-hover);
}

.tag-dropdown-option--selected {
  background: var(--color-bg-selected);
}

.tag-option-color {
  width: 14px;
  height: 14px;
  border-radius: 3px;
  flex-shrink: 0;
}

.tag-option-name {
  flex: 1;
  font-size: 14px;
  color: var(--color-text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.tag-option-count {
  font-size: 12px;
  color: var(--color-text-tertiary);
  margin-right: 8px;
}

.tag-option-check {
  color: var(--color-primary);
  flex-shrink: 0;
}

.tag-selector-error {
  margin-top: 6px;
  font-size: 12px;
  color: var(--color-error);
}
</style>
