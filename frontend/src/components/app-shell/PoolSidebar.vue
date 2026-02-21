<template>
  <div class="pool-sidebar-container" ref="containerRef">
    <section class="sidebar-group" :style="groupStyle(0)">
      <div class="group-title">目录</div>
      <div class="group-content">
        <TreeNode
          label="全部素材"
          icon="image"
          :count="totalCount"
          :depth="0"
          :is-active="activeItem === 'all' && !selectedPath"
          @click="handleSelectAll"
        />
        
        <!-- 绑定的素材目录树 -->
        <DirectoryTree
          :active-path="selectedPath"
          @select="handleSelectPath"
        />
      </div>
    </section>

    <div class="sidebar-resizer" @mousedown="startDrag(0, $event)">
      <span class="resizer-grip"></span>
    </div>

    <section class="sidebar-group" :style="groupStyle(1)">
      <div class="group-title-row">
        <div class="group-title">标签</div>
        <button class="group-add-btn" title="添加标签" @click="startAddTag()">
          <Icon name="plus" size="sm" />
        </button>
      </div>
      <div class="group-content">
        <!-- 空状态提示 -->
        <div v-if="tagRows.length === 0 && !editingTagId" class="empty-hint">
          <EmptyState
            icon="🏷️"
            title="暂无标签"
            compact
          />
        </div>
        
        <!-- 标签列表 -->
        <div
          v-for="tag in tagRows"
          :key="tag.id"
          class="tag-item"
          :class="{ active: activeTagId === tag.id }"
          :style="{ paddingLeft: `${16 + tag.depth * 20}px` }"
        >
          <!-- 编辑模式 -->
          <input
            v-if="editingTagId === tag.id"
            ref="editInputRef"
            v-model.trim="editingTagName"
            class="tag-edit-input"
            @keydown.enter.prevent="submitEditTag(tag)"
            @keydown.esc.prevent="cancelEdit"
            @blur="submitEditTag(tag)"
          />
          
          <!-- 显示模式 -->
          <template v-else>
            <Icon name="tag" size="sm" class="tag-icon" />
            <span class="tag-label" @click="selectTag(tag)" @dblclick="startEditTag(tag)">
              {{ tag.name }}
            </span>
            <span class="tag-count">{{ tag.file_count }}</span>
            <div class="tag-actions">
              <button class="tag-action-btn" title="添加子标签" @click="startAddTag(tag.id)">
                <Icon name="plus" size="xs" />
              </button>
              <button class="tag-action-btn" title="删除" @click="deleteTagConfirm(tag)">
                <Icon name="trash" size="xs" />
              </button>
            </div>
          </template>
        </div>
        
        <!-- 新增标签输入框 -->
        <div
          v-if="editingTagId === 'new'"
          class="tag-item"
          :style="{ paddingLeft: `${16 + newTagDepth * 20}px` }"
        >
          <input
            ref="newTagInputRef"
            v-model.trim="newTagName"
            class="tag-edit-input"
            placeholder="标签名称"
            @keydown.enter.prevent="submitAddTag"
            @keydown.esc.prevent="cancelEdit"
            @blur="submitAddTag"
          />
        </div>
      </div>
    </section>

    <div class="sidebar-resizer" @mousedown="startDrag(1, $event)">
      <span class="resizer-grip"></span>
    </div>

    <section class="sidebar-group" :style="groupStyle(2)">
      <div class="group-title">归档目录</div>
      <div class="group-content"></div>
    </section>

    <!-- Plugin Extension Point -->
    <ExtensionSlot 
      slot="Pool.Sidebar.Section"
      :context="{ view: 'pool', selectedPath, activeTagId }"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
import Icon from '../common/Icon.vue'
import EmptyState from '../common/EmptyState.vue'
import TreeNode from '../common/TreeNode.vue'
import DirectoryTree from '../directory/DirectoryTree.vue'
import ExtensionSlot from '../common/ExtensionSlot.vue'
import { useTagStore } from '@/stores/tagStore'

interface FolderItem {
  id: string | number
  name: string
  path?: string
  count: number
}

interface TagItem {
  id: string
  name: string
  color: string | null
  icon: string | null
  parent_id: string | null
  file_count: number
}

interface TagRow extends TagItem {
  depth: number
}

const props = defineProps<{
  activeItem: string
  folders: FolderItem[]
  tags: TagItem[]
  currentProjectId?: string
}>()

const emit = defineEmits<{
  'update:activeItem': [value: string]
  'select-folder': [folder: FolderItem]
  'select-tag': [tag: TagItem]
  'select-path': [path: string]
}>()

const activeFolderId = ref<string | number | null>(null)
const activeTagId = ref<string | null>(null)
const selectedPath = ref<string>('')
const currentProjectId = ref<string>(props.currentProjectId || '')
const containerRef = ref<HTMLElement | null>(null)

// 监听 projectId 变化
watch(() => props.currentProjectId, (newId) => {
  currentProjectId.value = newId || ''
  // 切换项目时重置选中路径
  selectedPath.value = ''
})
const tagStore = useTagStore()
const editingTagId = ref<string | null>(null)
const editingTagName = ref('')
const newTagName = ref('')
const newTagParentId = ref<string | null>(null)
const newTagDepth = ref(0)
const creatingTag = ref(false)
const newTagInputRef = ref<HTMLInputElement | null>(null)
const editInputRef = ref<HTMLInputElement | null>(null)

const RESIZER_SIZE = 6
const MIN_HEIGHTS = [120, 120, 100]
const groupHeights = ref<number[]>([220, 260, 180])

const totalCount = computed(() => {
  return props.folders.reduce((sum, f) => sum + f.count, 0)
})

const tagRows = computed<TagRow[]>(() => {
  const rows: TagRow[] = []
  const byParent = new Map<string | null, TagItem[]>()
  for (const tag of props.tags) {
    const key = tag.parent_id ?? null
    if (!byParent.has(key)) byParent.set(key, [])
    byParent.get(key)!.push(tag)
  }

  for (const list of byParent.values()) {
    list.sort((a, b) => a.name.localeCompare(b.name))
  }

  const walk = (parentId: string | null, depth: number) => {
    const children = byParent.get(parentId) || []
    for (const child of children) {
      rows.push({ ...child, depth })
      walk(child.id, depth + 1)
    }
  }

  walk(null, 0)

  if (rows.length < props.tags.length) {
    const exists = new Set(rows.map(r => r.id))
    for (const tag of props.tags) {
      if (!exists.has(tag.id)) rows.push({ ...tag, depth: 0 })
    }
  }

  return rows
})

function groupStyle(index: number) {
  return {
    height: `${groupHeights.value[index]}px`,
    minHeight: `${MIN_HEIGHTS[index]}px`,
  }
}

function normalizeHeights(totalHeight: number) {
  const minSum = MIN_HEIGHTS.reduce((sum, h) => sum + h, 0)
  const target = Math.max(totalHeight, minSum)
  const next = [...groupHeights.value]
  let sum = next.reduce((s, h) => s + h, 0)

  if (sum <= 0) {
    const avg = Math.floor(target / next.length)
    for (let i = 0; i < next.length; i++) next[i] = avg
    next[next.length - 1] += target - avg * next.length
    groupHeights.value = next
    return
  }

  for (let i = 0; i < next.length; i++) {
    next[i] = Math.max(MIN_HEIGHTS[i], Math.floor((next[i] / sum) * target))
  }

  sum = next.reduce((s, h) => s + h, 0)
  if (sum < target) {
    next[next.length - 1] += target - sum
  } else if (sum > target) {
    let overflow = sum - target
    for (let i = next.length - 1; i >= 0 && overflow > 0; i--) {
      const shrinkable = Math.max(0, next[i] - MIN_HEIGHTS[i])
      const delta = Math.min(shrinkable, overflow)
      next[i] -= delta
      overflow -= delta
    }
  }

  groupHeights.value = next
}

let ro: ResizeObserver | null = null
onMounted(() => {
  if (!containerRef.value) return
  const available = containerRef.value.clientHeight - RESIZER_SIZE * 2
  normalizeHeights(available)
  ro = new ResizeObserver((entries) => {
    const entry = entries[0]
    if (!entry) return
    const nextAvailable = entry.contentRect.height - RESIZER_SIZE * 2
    normalizeHeights(nextAvailable)
  })
  ro.observe(containerRef.value)
})

onUnmounted(() => {
  ro?.disconnect()
  ro = null
  stopDrag()
})

let draggingIndex: number | null = null
let startY = 0
let startA = 0
let startB = 0

function startDrag(index: number, event: MouseEvent) {
  draggingIndex = index
  startY = event.clientY
  startA = groupHeights.value[index]
  startB = groupHeights.value[index + 1]
  window.addEventListener('mousemove', onDrag)
  window.addEventListener('mouseup', stopDrag)
  document.body.style.cursor = 'row-resize'
  event.preventDefault()
}

function onDrag(event: MouseEvent) {
  if (draggingIndex === null) return

  const i = draggingIndex
  const total = startA + startB
  const delta = event.clientY - startY
  const minA = MIN_HEIGHTS[i]
  const minB = MIN_HEIGHTS[i + 1]

  let nextA = startA + delta
  nextA = Math.max(minA, Math.min(nextA, total - minB))
  const nextB = total - nextA

  const next = [...groupHeights.value]
  next[i] = nextA
  next[i + 1] = nextB
  groupHeights.value = next
}

function stopDrag() {
  if (draggingIndex === null) return
  draggingIndex = null
  window.removeEventListener('mousemove', onDrag)
  window.removeEventListener('mouseup', stopDrag)
  document.body.style.cursor = ''
}

function handleSelectAll() {
  selectedPath.value = ''
  activeFolderId.value = null
  activeTagId.value = null
  emit('update:activeItem', 'all')
}

function handleSelectPath(path: string) {
  selectedPath.value = path
  activeFolderId.value = null
  activeTagId.value = null
  emit('select-path', path)
}

function selectFolder(folder: FolderItem) {
  selectedPath.value = ''
  activeFolderId.value = folder.id
  activeTagId.value = null
  emit('select-folder', folder)
}

function selectTag(tag: TagItem) {
  selectedPath.value = ''
  activeTagId.value = tag.id
  activeFolderId.value = null
  emit('select-tag', tag)
}

// 标签编辑功能
function startAddTag(parentId: string | null = null) {
  editingTagId.value = 'new'
  newTagName.value = ''
  newTagParentId.value = parentId
  
  // 计算缩进深度
  if (parentId) {
    const parent = tagRows.value.find(t => t.id === parentId)
    newTagDepth.value = parent ? parent.depth + 1 : 0
  } else {
    newTagDepth.value = 0
  }
  
  nextTick(() => newTagInputRef.value?.focus())
}

function startEditTag(tag: TagItem) {
  editingTagId.value = tag.id
  editingTagName.value = tag.name
  nextTick(() => editInputRef.value?.focus())
}

function cancelEdit() {
  editingTagId.value = null
  editingTagName.value = ''
  newTagName.value = ''
  newTagParentId.value = null
}

async function submitAddTag() {
  const name = newTagName.value.trim()
  if (!name || creatingTag.value) {
    cancelEdit()
    return
  }

  creatingTag.value = true
  try {
    await tagStore.createTag({ 
      name,
      parent_id: newTagParentId.value 
    })
    cancelEdit()
  } catch (e) {
    console.error('Failed to create tag:', e)
  } finally {
    creatingTag.value = false
  }
}

async function submitEditTag(tag: TagItem) {
  const name = editingTagName.value.trim()
  if (!name || name === tag.name) {
    cancelEdit()
    return
  }

  try {
    await tagStore.updateTag({ 
      id: tag.id,
      name 
    })
    cancelEdit()
  } catch (e) {
    console.error('Failed to update tag:', e)
  }
}

async function deleteTagConfirm(tag: TagItem) {
  if (!confirm(`确定删除标签"${tag.name}"吗？`)) return
  
  try {
    await tagStore.deleteTag(tag.id)
    if (activeTagId.value === tag.id) {
      activeTagId.value = null
    }
  } catch (e) {
    console.error('Failed to delete tag:', e)
  }
}

function toggleAddTag() {
  startAddTag()
}

function cancelAddTag() {
  cancelEdit()
}
</script>

<style scoped>
.pool-sidebar-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
}

.sidebar-group {
  display: flex;
  flex-direction: column;
  padding: 12px 0 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  min-height: 0;
  overflow: hidden;
}

.sidebar-group:last-child {
  border-bottom: none;
}

.group-content {
  flex: 1;
  overflow-y: auto;
  min-height: 0;
  padding-bottom: 12px;
  /* 确保即使为空也保持最小高度 */
  min-height: 60px;
}

.empty-hint {
  min-height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.sidebar-resizer {
  height: 6px;
  cursor: row-resize;
  display: flex;
  align-items: center;
  justify-content: center;
  user-select: none;
}

.resizer-grip {
  width: 40px;
  height: 2px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.12);
}

.sidebar-resizer:hover .resizer-grip {
  background: rgba(96, 165, 250, 0.6);
}

.group-title {
  padding: 0 16px 8px;
  font-size: 11px;
  font-weight: 600;
  color: #6b7280;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.group-title-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.group-add-btn {
  margin-right: 12px;
  width: 20px;
  height: 20px;
  border: none;
  border-radius: 4px;
  background: transparent;
  color: #9ca3af;
  cursor: pointer;
  font-size: 16px;
  line-height: 1;
}

.group-add-btn:hover {
  background: rgba(255, 255, 255, 0.08);
  color: #e5e7eb;
}

/* 标签项 */
.tag-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 16px;
  cursor: pointer;
  transition: background 0.15s;
  position: relative;
}

.tag-item:hover {
  background: rgba(255, 255, 255, 0.05);
}

.tag-item.active {
  background: rgba(59, 130, 246, 0.15);
}

.tag-icon {
  color: #9ca3af;
  flex-shrink: 0;
}

.tag-label {
  flex: 1;
  font-size: 13px;
  color: #d1d5db;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.tag-count {
  font-size: 11px;
  color: #6b7280;
  flex-shrink: 0;
}

.tag-actions {
  display: none;
  align-items: center;
  gap: 2px;
}

.tag-item:hover .tag-actions {
  display: flex;
}

.tag-action-btn {
  width: 20px;
  height: 20px;
  border: none;
  border-radius: 3px;
  background: transparent;
  color: #9ca3af;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0;
}

.tag-action-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  color: #e5e7eb;
}

/* 标签编辑输入框 - 下划线样式 */
.tag-edit-input {
  flex: 1;
  background: transparent;
  border: none;
  border-bottom: 1px solid #3b82f6;
  color: #e5e7eb;
  font-size: 13px;
  padding: 2px 4px;
  outline: none;
  font-family: inherit;
}

.tag-edit-input::placeholder {
  color: #6b7280;
}

.nav-item {
  display: flex;
  align-items: center;
  padding: 8px 16px;
  cursor: pointer;
  transition: background 0.15s;
}

.nav-item:hover {
  background: rgba(255, 255, 255, 0.06);
}

.nav-item.active {
  background: rgba(59, 130, 246, 0.15);
}

.nav-item .nav-icon {
  width: 16px;
  height: 16px;
  margin-right: 10px;
  flex-shrink: 0;
}

.nav-item .label {
  flex: 1;
  font-size: 12px;
  color: #d1d5db;
}

.nav-item.active .label {
  color: #60a5fa;
}

.nav-item .count {
  font-size: 11px;
  color: #6b7280;
  background: rgba(255, 255, 255, 0.06);
  padding: 2px 6px;
  border-radius: 4px;
}

.nav-item--tag {
  min-height: 30px;
}
</style>
