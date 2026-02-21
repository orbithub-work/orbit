<template>
  <div class="app-shell">
    <header class="app-titlebar" :class="`app-titlebar--${platform}`">
      <div class="titlebar-left">
        <div v-if="isMac" class="traffic-lights">
          <button class="light red tooltip tooltip-top" data-tip="关闭" @click.stop="closeWindow"></button>
          <button class="light yellow tooltip tooltip-top" data-tip="最小化" @click.stop="minimizeWindow"></button>
          <button class="light green tooltip tooltip-top" data-tip="缩放" @click.stop="maximizeWindow"></button>
        </div>
        <div class="app-icon">🎬</div>
        <span class="app-name">智归档OS</span>
      </div>
      
      <div class="titlebar-center">
        <div class="tab-wrapper">
          <div class="tab-slider" :style="tabSliderStyle"></div>
          <button
            v-for="(tab, index) in tabs"
            :key="tab.key"
            class="tab-btn tooltip tooltip-top"
            :class="{ active: currentRoute === tab.key }"
            @click="router.push(`/${tab.key}`)"
            :data-tip="tab.tooltip"
            :ref="el => tabRefs[index] = el"
          >
            {{ tab.label }}
          </button>
        </div>
      </div>
      
      <div class="titlebar-right">
        <button v-if="!isMac" class="window-btn window-btn--windows tooltip tooltip-top" data-tip="最小化" @click.stop="minimizeWindow">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="5" y1="12" x2="19" y2="12"></line>
          </svg>
        </button>
        <button v-if="!isMac" class="window-btn window-btn--windows tooltip tooltip-top" data-tip="最大化" @click.stop="maximizeWindow">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="4" y="4" width="16" height="16" rx="2"></rect>
          </svg>
        </button>
        <button v-if="!isMac" class="window-btn window-btn--windows window-btn--close tooltip tooltip-top" data-tip="关闭" @click.stop="closeWindow">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="18" y1="6" x2="6" y2="18"></line>
            <line x1="6" y1="6" x2="18" y2="18"></line>
          </svg>
        </button>
      </div>
    </header>

    <main class="app-content">
      <aside class="eagle-sidebar">
        <router-view name="sidebar" 
          v-model:active-item="activeSidebarItem"
          :folders="poolFolders"
          :tags="poolTags"
          :current-project-id="fileStore.projectId"
          :projects="projectList"
          :active-project-id="activeProjectId"
          :active-filter="activeArtifactFilter"
          @select-folder="handleSelectFolder"
          @select-tag="handleSelectTag"
          @select-path="handleSelectPath"
          @create-project="handleCreateProject"
          @update:active-filter="activeArtifactFilter = $event"
          @update:active-project-id="activeProjectId = $event"
        />
      </aside>

      <div class="layout-main">
        <router-view name="main"
          :assets="fileStore.files"
          :loading="fileStore.loading"
          :selected-count="selectedAssetIds.size"
          :current-path="fileStore.currentPath"
          :project-id="fileStore.projectId"
          :project="activeProject"
          :active-filter="activeArtifactFilter"
          @select="handleAssetSelect"
          @open="handleAssetOpen"
          @context-menu="handleAssetContextMenu"
          @navigate="handleNavigate"
        />
      </div>

      <Inspector
        v-if="showInspector && selectedAsset"
        :mode="currentRoute"
        :selected-id="lastSelectedId"
        :project="activeProject"
        :selected-asset="selectedAsset"
        @close="showInspector = false"
      />
    </main>

    <TaskPanel 
      :visible="showTaskPanel"
      @close="showTaskPanel = false"
    />

    <footer class="app-statusbar">
      <StatusBar 
        :stats="statusBarStats"
        :ws-connected="wsConnected"
        @toggle-tasks="showTaskPanel = !showTaskPanel"
      />
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch, nextTick, type ComponentPublicInstance } from 'vue'
import { useRouter } from 'vue-router'
import { useProjectStore } from '@/stores/projectStore'
import { useFileStore, type FileItem } from '@/stores/fileStore'
import { useTagStore } from '@/stores/tagStore'
import { useArtifactStore } from '@/stores/artifactStore'
import { apiCall, connectWebSocket, subscribeToEvent } from '@/services/api'
import type { Project } from '@/stores/projectStore'
import type { TagItem } from '@/stores/tagStore'

import Inspector from '@/components/app-shell/Inspector.vue'
import StatusBar from '@/components/app-shell/StatusBar.vue'
import TaskPanel from '@/components/app-shell/TaskPanel.vue'

const projectStore = useProjectStore()
const fileStore = useFileStore()
const tagStore = useTagStore()
const artifactStore = useArtifactStore()
const router = useRouter()

const selectedAssetIds = ref<Set<string>>(new Set())
const selectedAsset = ref<FileItem | null>(null)

const activeSidebarItem = ref('all')
const activeProjectId = ref<string | number>('')
const activeArtifactFilter = ref('all')
const showInspector = ref(true)
const showTaskPanel = ref(false)
const wsConnected = ref(false)
const platform = ref<'mac' | 'windows' | 'linux'>('windows')
const isMac = computed(() => platform.value === 'mac')

const currentRoute = computed(() => router.currentRoute.value.path.split('/')[1] || 'pool')

const tabs = [
  { key: 'pool', label: '素材库', tooltip: '管理所有素材文件' },
  { key: 'project', label: '项目库', tooltip: '管理和追踪创作项目' },
  { key: 'artifact', label: '交付库', tooltip: '已完成的交付物和发布内容' },
  { key: 'analytics', label: '数据看板', tooltip: '数据统计、数据分析和对比' },
]

const activeTabIndex = computed(() => tabs.findIndex(t => t.key === currentRoute.value))
const tabRefs = ref<(HTMLElement | Element | ComponentPublicInstance | null)[]>([])

const tabSliderStyle = computed(() => {
  const activeTab = tabRefs.value[activeTabIndex.value]
  if (!activeTab || typeof activeTab === 'object' && 'ctx' in activeTab) {
    return {
      transform: 'translateX(0)',
      width: '0'
    }
  }

  const domElement = activeTab as HTMLElement

  // 由于我们已经将 margin 应用到了 tab 按钮上，
  // offsetLeft 已经包含了所有必要的间距，所以不需要再调整
  return {
    transform: `translateX(${domElement.offsetLeft}px)`,
    width: `${domElement.offsetWidth}px`
  }
})

const normalizeProjectId = (id: string | number | undefined | null): string =>
  id === undefined || id === null ? '' : String(id)

const poolFolders = computed(() => {
  return projectStore.projects.map((p: Project) => ({
    id: p.id,
    name: p.name,
    path: p.path,
    count: p.stats?.file_count || 0
  }))
})

// 状态栏统计数据
const statusBarStats = computed(() => {
  const activeProject = projectStore.projects.find(p => String(p.id) === String(activeProjectId.value))
  
  return {
    projectName: activeProject?.name,
    totalAssets: fileStore.files.length,
    selectedCount: selectedAssetIds.value.size,
    storageSize: activeProject?.stats?.total_size,
  }
})

const poolTags = computed(() => {
  return tagStore.tags
})

const projectList = computed(() => {
  if (projectStore.projects.length === 0) return []
  
  return projectStore.projects.map(p => ({
    id: p.id,
    name: p.name,
    client: p.metadata?.client || '未指定客户',
    owner: p.metadata?.owner || '未指定负责人',
    deadline: p.deadline || '无截止日期',
    status: p.status === 'active' ? '进行中' : (p.status === 'completed' ? '已交付' : '归档'),
    type: p.project_type === 'default' ? '通用项目' : p.project_type,
    assets: p.stats?.file_count || 0,
    references: (p.metadata?.references_count || 0) + ' 项',
    deliverables: (p.deliverables?.length || 0) + ' 项',
    versions: (p.iterations?.length || 0) + ' 个版本',
  }))
})

const activeProject = computed(() => {
  return projectList.value.find(item => item.id === activeProjectId.value) ?? projectList.value[0]
})

const lastSelectedId = computed(() => {
  const ids = Array.from(selectedAssetIds.value)
  return ids.length > 0 ? parseInt(ids[ids.length - 1]) || null : null
})

watch(() => projectStore.projects, (newProjects) => {
  if (newProjects.length > 0 && !activeProjectId.value && currentRoute.value === 'project') {
    activeProjectId.value = newProjects[0].id
  }
}, { immediate: true })

// 监听路由切换
watch(currentRoute, async (newRoute) => {
  if (newRoute === 'pool') {
    // 切换到素材库，清空 projectId，显示全部素材
    fileStore.setProjectId('')
    await fileStore.loadFiles('')
  }
})

watch(() => activeProjectId.value, async (nextId) => {
  // 只在项目库标签页才设置 projectId
  if (currentRoute.value === 'project') {
    const projectId = normalizeProjectId(nextId)
    if (!projectId) return
    fileStore.setProjectId(projectId)
    await fileStore.loadFiles(fileStore.currentPath)
  }
})

function handleCreateProject() {
  // 清空当前项目，触发模板引导页
  activeProjectId.value = ''
}

onMounted(async () => {
  const rawPlatform = (window as any)?.mediaAssistant?.platform
  if (rawPlatform === 'darwin') platform.value = 'mac'
  else if (rawPlatform === 'win32') platform.value = 'windows'
  else if (rawPlatform === 'linux') platform.value = 'linux'
  else if (/Mac|iPhone|iPad|iPod/i.test(navigator.platform)) platform.value = 'mac'
  else if (/Win/i.test(navigator.platform)) platform.value = 'windows'

  await Promise.all([
    projectStore.fetchProjects(),
    tagStore.loadTags()
  ])
  artifactStore.loadMockData()

  // 默认选择第一个项目
  const initialProjectId = normalizeProjectId(projectStore.projects[0]?.id)
  if (initialProjectId) {
    fileStore.setProjectId(initialProjectId)
    await fileStore.loadFiles('')
  }

  // 连接 WebSocket 并监听权限告警
  const warned = new Set<string>()
  try {
    await connectWebSocket()
    wsConnected.value = true
  } catch (e) {
    console.error('WebSocket connection failed:', e)
    wsConnected.value = false
  }

  subscribeToEvent('system_warning', (msg) => {
    const d = msg?.data
    if (d?.code !== 'watch_permission_denied') return
    const key = `${d.project_id}:${d.path}`
    if (warned.has(key)) return
    warned.add(key)
    console.warn(`[权限告警] 无权限监听目录: ${d.path}`)
  })

  // 启动补偿：拉最近日志
  try {
    const logs = await apiCall<any[]>('list_activity_logs', { limit: 100 })
    logs
      ?.filter(l => l.level === 'WARN' && String(l.message).startsWith('monitor permission denied:'))
      .forEach(l => {
        console.warn(`[历史告警] ${l.message}`)
      })
  } catch (err) {
    console.error('[活动日志] 拉取失败:', err)
  }

  try {
    const isFirstLaunch = await apiCall<boolean>('is_first_launch')
    if (isFirstLaunch) {
      router.replace('/onboarding')
    }
  } catch (e) {
    console.warn('Failed to check first launch status', e)
  }

  // 组件挂载后初始化滑块位置
  nextTick(() => {
    // 强制重新计算滑块样式
  })
})

function handleAssetSelect(items: FileItem[]) {
  selectedAssetIds.value = new Set(items.map(i => i.id))
  selectedAsset.value = items.length > 0 ? items[items.length - 1] : null
  if (items.length > 0) {
    showInspector.value = true
  }
}

function handleAssetOpen(item: FileItem) {
  if (item.is_directory) {
    fileStore.loadFiles(item.path)
  } else {
    fileStore.openFile(item)
  }
}

function handleAssetContextMenu(item: FileItem, event: MouseEvent) {
  console.log('Context menu:', item, event)
}

function handleNavigate(path: string) {
  fileStore.loadFiles(path)
}

function handleSelectFolder(folder: any) {
  console.log('Selected folder:', folder)
  // TODO: 根据文件夹筛选素材
}

function handleSelectTag(tag: any) {
  console.log('Selected tag:', tag)
  // TODO: 根据标签筛选素材
}

function handleSelectPath(path: string) {
  console.log('Selected path:', path)
  // 根据选中的目录路径筛选素材
  fileStore.loadFiles(path)
}

function minimizeWindow() {
  const handler = (window as any)?.mediaAssistant?.window?.minimize
  if (typeof handler === 'function') handler()
}

function maximizeWindow() {
  const handler = (window as any)?.mediaAssistant?.window?.maximize
  if (typeof handler === 'function') handler()
}

function closeWindow() {
  const handler = (window as any)?.mediaAssistant?.window?.close
  if (typeof handler === 'function') handler()
  else window.close()
}
</script>

<style>
* {
  box-sizing: border-box;
}

html, body, #app {
  height: 100%;
  width: 100%;
  margin: 0;
  padding: 0;
  overflow: hidden;
}

body {
  background: #1e1e1e;
  color: #ccc;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
  font-size: 12px;
  -webkit-font-smoothing: antialiased;
}

.app-shell {
  display: flex;
  flex-direction: column;
  height: 100vh;
  height: 100dvh;
  min-height: 0;
  width: 100%;
  background: #1e1e1e;
  overflow: hidden;
}

.app-titlebar {
  height: 44px;
  display: grid;
  grid-template-columns: 1fr auto 1fr;
  align-items: center;
  background: #1b1c1f;
  border-bottom: 1px solid #2b2b2f;
  flex-shrink: 0;
  padding: 0 12px;
  user-select: none;
  -webkit-app-region: drag;
}

.titlebar-left {
  display: flex;
  align-items: center;
  gap: 12px;
  justify-self: start;
  -webkit-app-region: no-drag;
}

.traffic-lights {
  display: flex;
  gap: 8px;
}

.light {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  border: none;
  cursor: pointer;
  padding: 0;
}

.light.red {
  background: #f9615b;
}

.light.yellow {
  background: #fbd875;
}

.light.green {
  background: #40c764;
}

.app-icon {
  font-size: 16px;
}

.app-name {
  font-size: 13px;
  font-weight: 500;
  color: #e5e7eb;
}

.titlebar-center {
  display: flex;
  align-items: center;
  justify-self: center;
  -webkit-app-region: no-drag;
}

.tab-wrapper {
  position: relative;
  display: flex;
  padding: 4px;
  background: rgba(0, 0, 0, 0.3);
  border-radius: 10px;
  border: 1px solid rgba(255, 255, 255, 0.06);
}

.tab-slider {
  position: absolute;
  top: 4px;
  left: 4px;
  width: 68px;
  height: calc(100% - 8px);
  background: linear-gradient(135deg, #2d3748 0%, #1a202c 100%);
  border-radius: 8px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  box-shadow:
    0 2px 8px rgba(0, 0, 0, 0.4),
    inset 0 1px 0 rgba(255, 255, 255, 0.05);
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  z-index: 0;
  pointer-events: none;
}

.tab-btn {
  position: relative;
  z-index: 1;
  width: 68px;
  padding: 8px 0;
  border: none;
  background: transparent;
  color: rgba(255, 255, 255, 0.4);
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  border-radius: 6px;
  transition: color 0.25s ease;
  text-align: center;
  white-space: nowrap;
}

.tab-btn:hover {
  color: rgba(255, 255, 255, 0.7);
}

.tab-btn.active {
  color: #fff;
}

.titlebar-right {
  display: flex;
  align-items: center;
  justify-self: end;
  gap: 8px;
  -webkit-app-region: no-drag;
}

.window-btn {
  width: 28px;
  height: 28px;
  border: none;
  background: transparent;
  color: #9ca3af;
  cursor: pointer;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.window-btn:hover {
  background: rgba(255, 255, 255, 0.08);
  color: #e5e7eb;
}

.app-titlebar--windows,
.app-titlebar--linux {
  padding-right: 0;
}

.app-titlebar--windows .titlebar-right,
.app-titlebar--linux .titlebar-right {
  gap: 0;
}

.window-btn--windows {
  width: 46px;
  height: 32px;
  border-radius: 0;
  color: #a7acb8;
}

.window-btn--windows:hover {
  background: rgba(255, 255, 255, 0.08);
}

.window-btn--close:hover {
  background: #e81123;
  color: #ffffff;
}

.app-content {
  flex: 1;
  display: flex;
  overflow: hidden;
  min-height: 0;
}

.eagle-sidebar {
  width: 220px;
  background: #1b1c1f;
  border-right: 1px solid #2b2b2f;
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
  min-height: 0;
  overflow: hidden;
}

.app-content:has(.analytics-view) .eagle-sidebar {
  display: none;
}

.app-content:has(.analytics-view) .layout-main {
  width: 100%;
}

.layout-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
  min-height: 0;
  overflow: hidden;
}

.app-statusbar {
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: #1b1c1f;
  border-top: 1px solid #2b2b2f;
  flex-shrink: 0;
  padding: 0 12px;
  font-size: 11px;
  color: #6b7280;
}

.statusbar-left,
.statusbar-center,
.statusbar-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.status-item {
  white-space: nowrap;
}
</style>
