<template>
  <div class="project-sidebar-container" ref="containerRef">
    <section class="sidebar-group" :style="groupStyle(0)">
      <div class="group-title-row">
        <div class="group-title">项目列表</div>
        <button class="group-add-btn" title="新建项目" @click="$emit('create-project')">
          <Icon name="plus" size="sm" />
        </button>
      </div>
      <div class="group-content">
        <!-- 始终显示结构，避免布局变形 -->
        <div v-if="projects.length === 0" class="empty-hint">
          <EmptyState
            icon="📁"
            title="暂无项目"
            compact
          />
        </div>
        <TreeNode
          v-for="project in projects"
          :key="project.id"
          :label="project.name"
          icon="archive"
          :count="project.assets"
          :depth="0"
          :is-active="activeProjectId === project.id"
          @click="$emit('update:activeProjectId', project.id)"
        />
      </div>
    </section>

    <div class="sidebar-resizer" @mousedown="startDrag(0, $event)">
      <span class="resizer-grip"></span>
    </div>

    <section class="sidebar-group" :style="groupStyle(1)">
      <div class="group-title">项目状态</div>
      <div class="group-content">
        <TreeNode label="进行中" icon="rocket" :count="4" :depth="0" />
        <TreeNode label="已交付" icon="check" :count="12" :depth="0" />
        <TreeNode label="归档" icon="cube" :count="6" :depth="0" />
      </div>
    </section>

    <div class="sidebar-resizer" @mousedown="startDrag(1, $event)">
      <span class="resizer-grip"></span>
    </div>

    <section class="sidebar-group" :style="groupStyle(2)">
      <div class="group-title">成员</div>
      <div class="group-content">
        <div class="folder-list">
          <div class="nav-item">
            <Icon name="user" size="sm" class="nav-icon" />
            <span class="label">叶柯</span>
          </div>
          <div class="nav-item">
            <Icon name="user" size="sm" class="nav-icon" />
            <span class="label">文豪</span>
          </div>
          <div class="nav-item">
            <span class="icon">👤</span>
            <span class="label">Rina</span>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import Icon from '@/components/common/Icon.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import TreeNode from '@/components/common/TreeNode.vue'

interface Project {
  id: string | number
  name: string
  assets: number
}

defineProps<{
  projects: Project[]
  activeProjectId: string | number
}>()

defineEmits<{
  'update:activeProjectId': [value: string | number]
  'create-project': []
}>()

const containerRef = ref<HTMLElement | null>(null)
const RESIZER_SIZE = 6
const MIN_HEIGHTS = [120, 100, 100]
const groupHeights = ref<number[]>([260, 140, 220])

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
</script>

<style scoped>
.project-sidebar-container {
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
  min-height: 60px;
  padding-bottom: 12px;
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

.nav-item .icon {
  width: 20px;
  font-size: 14px;
  text-align: center;
  margin-right: 10px;
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
</style>
