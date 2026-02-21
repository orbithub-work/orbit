<template>
  <div
    class="dock-container"
    :class="{ 'dock-collapsed': isCollapsed }"
  >
    <!-- Drag handle for moving the dock -->
    <div
      class="dock-handle"
      @mousedown="startDrag"
      @dblclick="toggleCollapse"
    >
      <div class="dock-dots">
        <span></span>
        <span></span>
        <span></span>
      </div>
      <div class="dock-title">
        智归档 Dock
      </div>
      <div class="dock-controls">
        <button
          class="dock-btn"
          :title="isCollapsed ? '展开' : '收起'"
          @click="toggleCollapse"
        >
          <svg
            v-if="isCollapsed"
            width="12"
            height="12"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
          >
            <polyline points="6 9 12 15 18 9" />
          </svg>
          <svg
            v-else
            width="12"
            height="12"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
          >
            <polyline points="18 15 12 9 6 15" />
          </svg>
        </button>
        <button
          class="dock-btn"
          title="关闭"
          @click="closeDock"
        >
          <svg
            width="12"
            height="12"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
          >
            <line
              x1="18"
              y1="6"
              x2="6"
              y2="18"
            />
            <line
              x1="6"
              y1="6"
              x2="18"
              y2="18"
            />
          </svg>
        </button>
      </div>
    </div>

    <!-- Main content area -->
    <div
      v-show="!isCollapsed"
      class="dock-content"
    >
      <!-- Quick filters -->
      <div class="dock-filters">
        <button
          v-for="filter in filters"
          :key="filter.id"
          class="filter-btn"
          :class="{ active: activeFilter === filter.id }"
          @click="setFilter(filter.id)"
        >
          {{ filter.label }}
        </button>
      </div>

      <div
        v-if="isDrawerOpen"
        class="dock-drawer"
      >
        <div class="drawer-header">
          <div class="drawer-title">
            {{ activeEntry?.name || '抽屉' }}
          </div>
          <div class="drawer-controls">
            <input
              v-model="searchQuery"
              class="drawer-search"
              placeholder="搜索当前抽屉"
            />
            <button
              class="drawer-close"
              @click="closeDrawer"
            >
              <svg
                width="12"
                height="12"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
              >
                <line
                  x1="18"
                  y1="6"
                  x2="6"
                  y2="18"
                />
                <line
                  x1="6"
                  y1="6"
                  x2="18"
                  y2="18"
                />
              </svg>
            </button>
          </div>
        </div>
        <div class="drawer-body">
          <div
            v-if="drawerLoading"
            class="drawer-loading"
          >
            加载中...
          </div>
          <div
            v-else-if="filteredDrawerItems.length === 0"
            class="drawer-empty"
          >
            暂无内容
          </div>
          <div
            v-else
            class="drawer-list"
          >
            <div
              v-for="item in filteredDrawerItems"
              :key="item.id"
              class="drawer-item"
              draggable="true"
              @click="openFileItem(item)"
              @dragstart="handleFileDragStart($event, item)"
            >
              <div
                class="drawer-icon"
                :class="getFileIconClass(item.type)"
              >
                <svg
                  width="18"
                  height="18"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="1.5"
                >
                  <rect
                    x="3"
                    y="3"
                    width="18"
                    height="18"
                    rx="2"
                    ry="2"
                  />
                  <circle
                    v-if="item.type === 'video'"
                    cx="12"
                    cy="12"
                    r="5"
                  />
                  <polygon
                    v-if="item.type === 'video'"
                    points="10 8 16 12 10 16 10 8"
                  />
                  <path
                    v-if="item.type === 'image'"
                    d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"
                  />
                  <polyline
                    v-if="item.type === 'image'"
                    points="7 10 12 15 17 10"
                  />
                  <line
                    v-if="item.type === 'image'"
                    x1="12"
                    y1="15"
                    x2="12"
                    y2="3"
                  />
                </svg>
              </div>
              <div class="drawer-info">
                <div class="drawer-name">
                  {{ item.name }}
                </div>
                <div class="drawer-meta">
                  {{ formatDate(item.modifiedAt) }}
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Scrollable asset strip -->
      <div
        ref="stripRef"
        class="dock-strip"
        @wheel="handleWheel"
      >
        <div class="dock-items">
          <div
            v-for="item in filteredDockEntries"
            :key="item.id"
            class="dock-item"
            :title="item.name"
            :data-project-id="item.kind === 'project' ? item.id : ''"
            draggable="true"
            @click="openEntry(item)"
            @dragstart="handleEntryDragStart($event, item)"
          >
            <div class="item-thumbnail">
              <div
                class="item-icon"
                :class="getFileIconClass(item.type)"
              >
                <svg
                  width="24"
                  height="24"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="1.5"
                >
                  <rect
                    x="3"
                    y="3"
                    width="18"
                    height="18"
                    rx="2"
                    ry="2"
                  />
                  <circle
                    v-if="item.type === 'video'"
                    cx="12"
                    cy="12"
                    r="5"
                  />
                  <polygon
                    v-if="item.type === 'video'"
                    points="10 8 16 12 10 16 10 8"
                  />
                  <path
                    v-if="item.type === 'image'"
                    d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"
                  />
                  <polyline
                    v-if="item.type === 'image'"
                    points="7 10 12 15 17 10"
                  />
                  <line
                    v-if="item.type === 'image'"
                    x1="12"
                    y1="15"
                    x2="12"
                    y2="3"
                  />
                </svg>
              </div>
            </div>
            <div class="item-info">
              <div class="item-name">
                {{ truncateName(item.name) }}
              </div>
              <div class="item-meta">
                {{ item.kind === 'project' ? '项目' : '目录' }}
              </div>
            </div>
          </div>

          <!-- Empty state -->
          <div
            v-if="filteredDockEntries.length === 0"
            class="dock-empty"
          >
            <p>暂无 Dock 入口</p>
            <small>请先配置常用目录或项目</small>
          </div>
        </div>

        <!-- Scroll indicators -->
        <button
          v-if="canScrollLeft"
          class="scroll-btn left"
          @click="scroll(-1)"
        >
          <svg
            width="16"
            height="16"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
          >
            <polyline points="15 18 9 12 15 6" />
          </svg>
        </button>
        <button
          v-if="canScrollRight"
          class="scroll-btn right"
          @click="scroll(1)"
        >
          <svg
            width="16"
            height="16"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
          >
            <polyline points="9 18 15 12 9 6" />
          </svg>
        </button>
      </div>

      <!-- Quick actions -->
      <div class="dock-actions">
        <button
          class="action-btn"
          title="打开主窗口"
          @click="openMainWindow"
        >
          <svg
            width="16"
            height="16"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
          >
            <rect
              x="2"
              y="3"
              width="20"
              height="14"
              rx="2"
              ry="2"
            />
            <line
              x1="8"
              y1="21"
              x2="16"
              y2="21"
            />
            <line
              x1="12"
              y1="17"
              x2="12"
              y2="21"
            />
          </svg>
        </button>
        <button
          class="action-btn"
          title="刷新"
          @click="refreshItems"
        >
          <svg
            width="16"
            height="16"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
          >
            <polyline points="23 4 23 10 17 10" />
            <polyline points="1 20 1 14 7 14" />
            <path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15" />
          </svg>
        </button>
        <button
          class="action-btn"
          title="搜索"
          @click="showSearch"
        >
          <svg
            width="16"
            height="16"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
          >
            <circle
              cx="11"
              cy="11"
              r="8"
            />
            <line
              x1="21"
              y1="21"
              x2="16.65"
              y2="16.65"
            />
          </svg>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { apiCall } from '@/services/api'

// ... existing types ...

// Types
type DockEntryType = 'system' | 'project'
type FileType = 'image' | 'video' | 'audio' | 'document' | 'other'

interface DockEntry {
  id: string
  name: string
  path: string
  type: 'other'
  kind: DockEntryType
  projectId?: string
}

interface DrawerItem {
  id: string
  name: string
  path: string
  type: FileType
  modifiedAt: Date
  size?: number
}

interface FilterOption {
  id: string
  label: string
}

// State
const isCollapsed = ref(false)
const activeFilter = ref('all')
const stripRef = ref<HTMLElement | null>(null)
const canScrollLeft = ref(false)
const canScrollRight = ref(false)
const dockEntries = ref<DockEntry[]>([])
const activeEntry = ref<DockEntry | null>(null)
const drawerItems = ref<DrawerItem[]>([])
const isDrawerOpen = ref(false)
const drawerLoading = ref(false)
const searchQuery = ref('')

const filters: FilterOption[] = [
  { id: 'all', label: '全部' },
  { id: 'system', label: '系统目录' },
  { id: 'project', label: '项目' }
]

// Computed
const filteredDockEntries = computed(() => {
  if (activeFilter.value === 'all') return dockEntries.value
  return dockEntries.value.filter(entry => entry.kind === activeFilter.value)
})

const filteredDrawerItems = computed(() => {
  const query = searchQuery.value.trim().toLowerCase()
  if (!query) return drawerItems.value
  return drawerItems.value.filter(item => item.name.toLowerCase().includes(query))
})

// Methods
const handleDomDragOver = (event: DragEvent) => {
  event.preventDefault()
}

const handleDomDrop = (event: DragEvent) => {
  event.preventDefault()
  const files = Array.from(event.dataTransfer?.files ?? [])
  if (files.length === 0) return

  // Check if we are in Electron environment (files have absolute paths)
  const isElectron = files.every(f => (f as any).path && (f as any).path !== '' && !(f as any).path.startsWith('/fake/'))
  
  if (isElectron) {
    const paths = files.map(file => (file as any).path)
    handleExternalDrop(event.clientX, event.clientY, paths)
  } else {
    console.warn('Web drag and drop is not supported for indexing. Please use the Electron app.')
  }
}

function setupFileDrop() {
  document.addEventListener('dragover', handleDomDragOver)
  document.addEventListener('drop', handleDomDrop)
}

function cleanupFileDrop() {
  document.removeEventListener('dragover', handleDomDragOver)
  document.removeEventListener('drop', handleDomDrop)
}

async function handleExternalDrop(x: number, y: number, paths: string[]) {
  // Find target project under coordinates
  const element = document.elementFromPoint(x, y)
  const dockItem = element?.closest('.dock-item') as HTMLElement
  const projectId = dockItem?.dataset.projectId

  if (!projectId) {
    console.warn('No project target for dropped files')
    // TODO: Drop into a default project or show selection
    return
  }

  try {
    await apiCall('archive_files', {
      project_id: projectId,
      paths: paths
    })
    console.log('Successfully archived files to project:', projectId)
    // Refresh if drawer is open for this project
    if (isDrawerOpen.value && activeEntry.value?.id === projectId) {
      loadDrawerItems(activeEntry.value)
    }
  } catch (err) {
    console.error('Failed to archive files:', err)
  }
}

function toggleCollapse() {
  isCollapsed.value = !isCollapsed.value
}

function closeDock() {
  isCollapsed.value = true
}

function setFilter(filterId: string) {
  activeFilter.value = filterId
  updateScrollIndicators()
}

async function openEntry(item: DockEntry) {
  activeEntry.value = item
  isDrawerOpen.value = true
  await loadDrawerItems(item)
}

function handleEntryDragStart(event: DragEvent, item: DockEntry) {
  if (event.dataTransfer) {
    event.dataTransfer.setData('text/plain', item.path)
    event.dataTransfer.effectAllowed = 'copy'
  }
}

async function openFileItem(item: DrawerItem) {
  if (item.id) {
    try {
      await apiCall('open_file', { id: item.id })
      return
    } catch {
      await apiCall('open_in_folder', { path: item.path })
      return
    }
  }
  await apiCall('open_in_folder', { path: item.path })
}

function handleFileDragStart(event: DragEvent, item: DrawerItem) {
  if (event.dataTransfer) {
    event.dataTransfer.setData('text/plain', item.path)
    event.dataTransfer.setData('text/uri-list', `file://${item.path}`)
    event.dataTransfer.effectAllowed = 'copy'
  }
}

function getFileIconClass(type: string): string {
  return `icon-${type}`
}

function truncateName(name: string, maxLength: number = 12): string {
  if (name.length <= maxLength) return name
  const ext = name.lastIndexOf('.')
  if (ext > 0) {
    const extName = name.slice(ext)
    const baseName = name.slice(0, ext)
    if (baseName.length > maxLength - 3) {
      return baseName.slice(0, maxLength - 3) + '...' + extName
    }
  }
  return name.slice(0, maxLength) + '...'
}

function formatDate(date: Date): string {
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const minutes = Math.floor(diff / 60000)
  const hours = Math.floor(diff / 3600000)
  const days = Math.floor(diff / 86400000)

  if (minutes < 60) return `${minutes}分钟前`
  if (hours < 24) return `${hours}小时前`
  if (days < 7) return `${days}天前`
  return date.toLocaleDateString('zh-CN')
}

function handleWheel(event: WheelEvent) {
  if (stripRef.value) {
    event.preventDefault()
    stripRef.value.scrollLeft += event.deltaY
    updateScrollIndicators()
  }
}

function scroll(direction: number) {
  if (stripRef.value) {
    const scrollAmount = 200
    stripRef.value.scrollBy({ left: direction * scrollAmount, behavior: 'smooth' })
    setTimeout(updateScrollIndicators, 300)
  }
}

function updateScrollIndicators() {
  if (!stripRef.value) return
  const { scrollLeft, scrollWidth, clientWidth } = stripRef.value
  canScrollLeft.value = scrollLeft > 0
  canScrollRight.value = scrollLeft < scrollWidth - clientWidth - 10
}

function openMainWindow() {
  if ((window as any).go?.main?.DockApp?.OpenMain) {
    (window as any).go.main.DockApp.OpenMain()
  }
}

function refreshItems() {
  loadDockEntries()
  if (activeEntry.value) {
    loadDrawerItems(activeEntry.value)
  }
}

function showSearch() {
  isDrawerOpen.value = true
  if (activeEntry.value) {
    loadDrawerItems(activeEntry.value)
  }
}

function closeDrawer() {
  isDrawerOpen.value = false
}

async function loadDockEntries() {
  const [commonDirs, projects] = await Promise.all([
    apiCall<any[]>('get_common_directories').catch(() => []),
    apiCall<any[]>('list_projects').catch(() => [])
  ])
  const systemEntries: DockEntry[] = (commonDirs || []).map((dir: any) => ({
    id: `system-${dir.path}`,
    name: dir.name || dir.path,
    path: dir.path,
    type: 'other',
    kind: 'system'
  }))
  const projectEntries: DockEntry[] = (projects || []).map((project: any) => ({
    id: `project-${project.id}`,
    name: project.name || '未命名项目',
    path: project.path || '',
    type: 'other',
    kind: 'project',
    projectId: project.id
  }))
  dockEntries.value = [...systemEntries, ...projectEntries]
  updateScrollIndicators()
}

async function loadDrawerItems(entry: DockEntry) {
  drawerLoading.value = true
  try {
    let items: any[] = []
    if (entry.kind === 'project' && entry.projectId) {
      items = await apiCall<any[]>('get_project_files', { projectId: entry.projectId })
    } else if (entry.path) {
      items = await apiCall<any[]>('list_files', { path: entry.path })
    }
    drawerItems.value = normalizeDrawerItems(items)
  } finally {
    drawerLoading.value = false
  }
}

function normalizeDrawerItems(items: any[]): DrawerItem[] {
  return (items || [])
    .filter(item => !item.is_directory)
    .map(item => {
      const name = item.name || getNameFromPath(item.path || '')
      const modifiedAt = parseDate(item.modified_at || item.modifiedAt)
      return {
        id: item.id || item.path || name,
        name,
        path: item.path || '',
        type: getFileType(name, item.file_type),
        modifiedAt,
        size: item.size
      }
    })
    .sort((a, b) => b.modifiedAt.getTime() - a.modifiedAt.getTime())
}

function parseDate(value: string | Date | undefined) {
  if (!value) return new Date(0)
  if (value instanceof Date) return value
  const parsed = new Date(value)
  if (Number.isNaN(parsed.getTime())) {
    return new Date(0)
  }
  return parsed
}

function getNameFromPath(path: string) {
  const parts = path.split(/[/\\]/)
  return parts[parts.length - 1] || path
}

function getFileType(name: string, rawType?: string): FileType {
  const type = rawType?.toLowerCase()
  if (type && ['image', 'video', 'audio', 'document'].includes(type)) {
    return type as FileType
  }
  const ext = name.split('.').pop()?.toLowerCase() || ''
  const imageExts = ['jpg', 'jpeg', 'png', 'gif', 'bmp', 'webp', 'svg', 'ico']
  const videoExts = ['mp4', 'avi', 'mkv', 'mov', 'wmv', 'flv', 'webm']
  const audioExts = ['mp3', 'wav', 'flac', 'aac', 'ogg', 'wma']
  const docExts = ['pdf', 'doc', 'docx', 'ppt', 'pptx', 'xls', 'xlsx', 'txt', 'md']
  if (imageExts.includes(ext)) return 'image'
  if (videoExts.includes(ext)) return 'video'
  if (audioExts.includes(ext)) return 'audio'
  if (docExts.includes(ext)) return 'document'
  return 'other'
}

function startDrag() {}

// Lifecycle
onMounted(() => {
  loadDockEntries()
  setupFileDrop()
  if (stripRef.value) {
    stripRef.value.addEventListener('scroll', updateScrollIndicators)
  }
})

onUnmounted(() => {
  cleanupFileDrop()
  if (stripRef.value) {
    stripRef.value.removeEventListener('scroll', updateScrollIndicators)
  }
})
</script>

<style scoped>
.dock-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: rgba(30, 30, 35, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
  user-select: none;
}

.dock-collapsed {
  height: auto;
}

.dock-handle {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  background: rgba(255, 255, 255, 0.05);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  cursor: grab;
}

.dock-handle:active {
  cursor: grabbing;
}

.dock-dots {
  display: flex;
  gap: 4px;
}

.dock-dots span {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.3);
}

.dock-dots span:first-child {
  background: #ff5f57;
}

.dock-dots span:nth-child(2) {
  background: #febc2e;
}

.dock-dots span:last-child {
  background: #28c840;
}

.dock-title {
  font-size: 12px;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.7);
}

.dock-controls {
  display: flex;
  gap: 4px;
}

.dock-btn {
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.1);
  border: none;
  border-radius: 4px;
  color: rgba(255, 255, 255, 0.7);
  cursor: pointer;
  transition: all 0.2s;
}

.dock-btn:hover {
  background: rgba(255, 255, 255, 0.2);
  color: white;
}

.dock-content {
  display: flex;
  flex-direction: column;
  flex: 1;
  overflow: hidden;
}

.dock-filters {
  display: flex;
  gap: 4px;
  padding: 8px 12px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  overflow-x: auto;
  scrollbar-width: none;
}

.dock-filters::-webkit-scrollbar {
  display: none;
}

.dock-drawer {
  margin: 8px 12px;
  padding: 8px;
  background: rgba(0, 0, 0, 0.35);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 10px;
}

.drawer-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.drawer-title {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.85);
}

.drawer-controls {
  display: flex;
  align-items: center;
  gap: 6px;
}

.drawer-search {
  width: 160px;
  height: 24px;
  padding: 0 8px;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.9);
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.12);
  border-radius: 6px;
  outline: none;
}

.drawer-search::placeholder {
  color: rgba(255, 255, 255, 0.4);
}

.drawer-close {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.08);
  border: none;
  border-radius: 6px;
  color: rgba(255, 255, 255, 0.6);
  cursor: pointer;
  transition: all 0.2s;
}

.drawer-close:hover {
  background: rgba(255, 255, 255, 0.16);
  color: rgba(255, 255, 255, 0.9);
}

.drawer-body {
  margin-top: 8px;
  max-height: 180px;
  overflow: auto;
}

.drawer-loading,
.drawer-empty {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.45);
  padding: 12px 4px;
}

.drawer-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.drawer-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 8px;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.08);
  cursor: pointer;
  transition: all 0.2s;
}

.drawer-item:hover {
  background: rgba(255, 255, 255, 0.1);
  border-color: rgba(255, 255, 255, 0.18);
}

.drawer-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border-radius: 6px;
  background: rgba(0, 0, 0, 0.25);
  color: rgba(255, 255, 255, 0.6);
}

.drawer-info {
  flex: 1;
  min-width: 0;
}

.drawer-name {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.85);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.drawer-meta {
  font-size: 10px;
  color: rgba(255, 255, 255, 0.4);
  margin-top: 2px;
}

.filter-btn {
  padding: 4px 12px;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.6);
  background: transparent;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}

.filter-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  color: rgba(255, 255, 255, 0.9);
}

.filter-btn.active {
  background: rgba(59, 130, 246, 0.3);
  border-color: rgba(59, 130, 246, 0.5);
  color: #60a5fa;
}

.dock-strip {
  position: relative;
  flex: 1;
  display: flex;
  align-items: center;
  overflow: hidden;
}

.dock-items {
  display: flex;
  gap: 8px;
  padding: 8px 12px;
  overflow-x: auto;
  scroll-behavior: smooth;
  scrollbar-width: none;
  -webkit-overflow-scrolling: touch;
}

.dock-items::-webkit-scrollbar {
  display: none;
}

.dock-item {
  flex-shrink: 0;
  width: 80px;
  padding: 8px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.dock-item:hover {
  background: rgba(255, 255, 255, 0.1);
  border-color: rgba(255, 255, 255, 0.2);
  transform: translateY(-2px);
}

.dock-item:active {
  transform: translateY(0);
}

.item-thumbnail {
  width: 100%;
  height: 50px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.3);
  border-radius: 4px;
  overflow: hidden;
  margin-bottom: 6px;
}

.item-thumbnail img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.item-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  color: rgba(255, 255, 255, 0.5);
}

.icon-image {
  color: #10b981;
}

.icon-video {
  color: #f59e0b;
}

.icon-audio {
  color: #8b5cf6;
}

.icon-document {
  color: #3b82f6;
}

.item-info {
  text-align: center;
}

.item-name {
  font-size: 10px;
  color: rgba(255, 255, 255, 0.9);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.item-meta {
  font-size: 9px;
  color: rgba(255, 255, 255, 0.4);
  margin-top: 2px;
}

.dock-empty {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 20px;
  color: rgba(255, 255, 255, 0.4);
}

.dock-empty p {
  font-size: 12px;
  margin: 0;
}

.dock-empty small {
  font-size: 10px;
  margin-top: 4px;
}

.scroll-btn {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  width: 24px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.5);
  border: none;
  color: rgba(255, 255, 255, 0.7);
  cursor: pointer;
  transition: all 0.2s;
  z-index: 10;
}

.scroll-btn:hover {
  background: rgba(0, 0, 0, 0.7);
  color: white;
}

.scroll-btn.left {
  left: 0;
  border-radius: 0 4px 4px 0;
}

.scroll-btn.right {
  right: 0;
  border-radius: 4px 0 0 4px;
}

.dock-actions {
  display: flex;
  gap: 8px;
  padding: 8px 12px;
  border-top: 1px solid rgba(255, 255, 255, 0.05);
  justify-content: flex-end;
}

.action-btn {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 6px;
  color: rgba(255, 255, 255, 0.6);
  cursor: pointer;
  transition: all 0.2s;
}

.action-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  border-color: rgba(255, 255, 255, 0.2);
  color: rgba(255, 255, 255, 0.9);
}
</style>
