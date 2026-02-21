<template>
  <div class="project-view">
    <!-- 骨架屏 -->
    <ProjectViewSkeleton v-if="loading" />
    
    <template v-else>
      <!-- 项目头部概览 -->
      <div class="project-header-section">
        <div class="header-main">
          <div class="project-info">
            <h2 class="project-title">
              {{ currentProject.name }}
            </h2>
            <div class="project-status-group">
              <span
                class="status-badge"
                :class="currentProject.status"
              >
                {{ getStatusText(currentProject.status) }}
              </span>
              <span class="project-type-tag">{{ getProjectTypeText(currentProject.project_type) }}</span>
            </div>
          </div>
          
          <div class="header-actions">
            <button
              class="action-btn sync-btn"
              title="同步监控文件夹"
              @click="handleSync"
            >
              <svg
                class="icon icon--16"
                viewBox="0 0 24 24"
                aria-hidden="true"
              >
                <use href="#i-sync" />
              </svg>
              同步
            </button>
            <button
              class="action-btn settings-btn"
              title="项目设置"
              @click="openSettings"
            >
              <svg
                class="icon icon--16"
                viewBox="0 0 24 24"
                aria-hidden="true"
              >
                <use href="#i-settings" />
              </svg>
            </button>
          </div>
        </div>

        <div class="header-stats">
          <div class="stat-card">
            <span class="stat-label">素材数量</span>
            <span class="stat-value">{{ currentProject.stats?.file_count || 0 }}</span>
          </div>
          <div class="stat-card">
            <span class="stat-label">迭代版本</span>
            <span class="stat-value">{{ currentProject.iterations?.length || 0 }}</span>
          </div>
          <div class="stat-card">
            <span class="stat-label">截止日期</span>
            <span class="stat-value">{{ currentProject.deadline ? formatDate(currentProject.deadline) : '未设置' }}</span>
          </div>
          <div class="stat-card progress-card">
            <div class="stat-label">
              完成进度
            </div>
            <div class="progress-bar-container">
              <div
                class="progress-bar-fill"
                :style="{ width: progressPercent + '%' }"
              ></div>
              <span class="progress-text">{{ progressPercent }}%</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 导航标签页 -->
      <div class="project-tabs">
        <button 
          v-for="tab in tabs" 
          :key="tab.key"
          :class="['tab-btn', { active: activeTab === tab.key }]"
          @click="activeTab = tab.key"
        >
          <svg
            class="icon icon--16"
            viewBox="0 0 24 24"
            aria-hidden="true"
          >
            <use :href="'#' + tab.icon" />
          </svg>
          {{ tab.label }}
        </button>
      </div>

      <!-- 内容区域 -->
      <div class="project-content">
        <!-- 概览页面 -->
        <div
          v-if="activeTab === 'overview'"
          class="tab-pane overview-pane"
        >
          <div class="pane-grid">
            <div class="grid-main">
              <!-- 项目描述 -->
              <section class="content-section info-section">
                <h3 class="section-title">
                  项目简介
                </h3>
                <p class="project-desc">
                  {{ currentProject.description || '暂无项目描述' }}
                </p>
              </section>

              <!-- 最近迭代 -->
              <section class="content-section iterations-section">
                <div class="section-header">
                  <h3 class="section-title">
                    最近迭代
                  </h3>
                  <button
                    class="text-btn"
                    @click="activeTab = 'workflow'"
                  >
                    查看全部
                  </button>
                </div>
                <div
                  v-if="currentProject.iterations?.length"
                  class="iteration-list"
                >
                  <div
                    v-for="it in currentProject.iterations.slice(0, 3)"
                    :key="it.id"
                    class="iteration-card"
                  >
                    <div class="it-version">
                      {{ it.version }}
                    </div>
                    <div class="it-info">
                      <div class="it-title">
                        {{ it.title }}
                      </div>
                      <div class="it-date">
                        {{ formatDate(it.created_at) }}
                      </div>
                    </div>
                  </div>
                </div>
                <div
                  v-else
                  class="empty-state"
                >
                  暂无迭代记录
                </div>
              </section>

              <!-- 监控文件夹 -->
              <section
                v-if="currentProject.folder_path"
                class="content-section folder-section"
              >
                <h3 class="section-title">
                  监控文件夹
                </h3>
                <div class="folder-path-box">
                  <svg
                    class="icon icon--20"
                    viewBox="0 0 24 24"
                    aria-hidden="true"
                  >
                    <use href="#i-folder" />
                  </svg>
                  <span class="path-text">{{ currentProject.folder_path }}</span>
                  <span
                    v-if="currentProject.auto_sync"
                    class="sync-badge"
                  >自动同步中</span>
                </div>
              </section>
            </div>

            <div class="grid-sidebar">
              <!-- 里程碑概览 -->
              <section class="content-section milestone-section">
                <h3 class="section-title">
                  关键里程碑
                </h3>
                <div class="milestone-timeline">
                  <div 
                    v-for="ms in currentProject.milestones" 
                    :key="ms.id" 
                    class="milestone-item"
                    :class="{ completed: ms.completed }"
                  >
                    <div class="ms-dot"></div>
                    <div class="ms-content">
                      <div class="ms-title">
                        {{ ms.title }}
                      </div>
                      <div class="ms-date">
                        {{ ms.due_date ? formatDate(ms.due_date) : '待定' }}
                      </div>
                    </div>
                  </div>
                </div>
              </section>

              <!-- 成品展示 -->
              <section class="content-section deliverable-section">
                <h3 class="section-title">
                  最新成品
                </h3>
                <div
                  v-if="currentProject.deliverables?.length"
                  class="deliverable-mini-grid"
                >
                  <div
                    v-for="dl in currentProject.deliverables.slice(0, 2)"
                    :key="dl.id"
                    class="dl-mini-card"
                  >
                    <div class="dl-icon">
                      <svg
                        class="icon icon--24"
                        viewBox="0 0 24 24"
                        aria-hidden="true"
                      >
                        <use :href="getFileSymbol(dl.file_type)" />
                      </svg>
                    </div>
                    <div class="dl-name">
                      {{ dl.name }}
                    </div>
                  </div>
                </div>
                <div
                  v-else
                  class="empty-state"
                >
                  暂无成品
                </div>
              </section>
            </div>
          </div>
        </div>

        <!-- 素材管理页面 -->
        <div
          v-if="activeTab === 'materials'"
          class="tab-pane materials-pane"
        >
          <div class="materials-toolbar">
            <div class="search-box">
              <input
                v-model="materialSearchQuery"
                type="text"
                placeholder="搜索项目素材..."
              />
            </div>
            <div class="filter-group">
              <button class="filter-btn active">
                全部
              </button>
              <button class="filter-btn">
                视频
              </button>
              <button class="filter-btn">
                音频
              </button>
              <button class="filter-btn">
                图片
              </button>
            </div>
          </div>
          
          <div class="material-tree-container">
            <div 
              v-for="folder in filteredProjectFolders" 
              :key="folder.id"
              class="folder-group"
            >
              <div
                class="folder-header"
                @click="toggleFolder(folder)"
              >
                <svg 
                  class="arrow-icon" 
                  :class="{ expanded: folder.expanded }"
                  viewBox="0 0 24 24"
                  aria-hidden="true"
                >
                  <use href="#i-chevron-right" />
                </svg>
                <span class="folder-name">{{ folder.name }}</span>
                <span class="item-count">{{ folder.items.length }}</span>
              </div>
              
              <div
                v-if="folder.expanded"
                class="folder-items"
              >
                <div 
                  v-for="item in folder.items" 
                  :key="item.id"
                  class="material-item"
                >
                  <div class="item-icon">
                    <svg
                      class="icon icon--20"
                      viewBox="0 0 24 24"
                      aria-hidden="true"
                    >
                      <use :href="getFileSymbol(item.type)" />
                    </svg>
                  </div>
                  <div class="item-main">
                    <div class="item-name">
                      {{ item.name }}
                    </div>
                    <div class="item-meta">
                      {{ item.type }} · {{ item.size || '未知大小' }}
                    </div>
                  </div>
                  <div class="item-actions">
                    <button
                      class="icon-btn"
                      title="查看详情"
                    >
                      <svg
                        class="icon icon--16"
                        viewBox="0 0 24 24"
                        aria-hidden="true"
                      >
                        <use href="#i-eye" />
                      </svg>
                    </button>
                    <button 
                      class="icon-btn" 
                      :class="{ 'active': isDeliverable(item) }"
                      title="标记为成品"
                      @click="toggleDeliverable(item)"
                    >
                      <svg
                        class="icon icon--16"
                        viewBox="0 0 24 24"
                        aria-hidden="true"
                      >
                        <use :href="isDeliverable(item) ? '#i-star' : '#i-star-outline'" />
                      </svg>
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 迭代 & 里程碑页面 -->
        <div
          v-if="activeTab === 'workflow'"
          class="tab-pane workflow-pane"
        >
          <div class="workflow-grid">
            <!-- 迭代记录 -->
            <div class="workflow-column">
              <div class="column-header">
                <h3>迭代记录 (Versions)</h3>
                <button
                  class="add-btn-small"
                  @click="addIteration"
                >
                  + 新迭代
                </button>
              </div>
              <div class="iteration-timeline">
                <div
                  v-for="it in currentProject.iterations"
                  :key="it.id"
                  class="it-full-card"
                >
                  <div class="it-header">
                    <span class="it-v-badge">{{ it.version }}</span>
                    <span class="it-full-title">{{ it.title }}</span>
                    <span class="it-full-date">{{ formatDate(it.created_at) }}</span>
                  </div>
                  <p class="it-full-desc">
                    {{ it.description }}
                  </p>
                  <div
                    v-if="it.file_ids?.length"
                    class="it-files"
                  >
                    <span
                      v-for="fid in it.file_ids"
                      :key="fid"
                      class="file-link"
                    >📎 相关文件</span>
                  </div>
                </div>
              </div>
            </div>

            <!-- 里程碑详情 -->
            <div class="workflow-column">
              <div class="column-header">
                <h3>项目里程碑 (Milestones)</h3>
                <button
                  class="add-btn-small"
                  @click="addMilestone"
                >
                  + 新里程碑
                </button>
              </div>
              <div class="milestone-list-full">
                <div 
                  v-for="ms in currentProject.milestones" 
                  :key="ms.id" 
                  class="ms-full-item"
                  :class="{ completed: ms.completed }"
                >
                  <div
                    class="ms-check"
                    @click="toggleMilestone(ms)"
                  >
                    <div class="check-box"></div>
                  </div>
                  <div class="ms-full-content">
                    <div class="ms-full-header">
                      <span class="ms-full-title">{{ ms.title }}</span>
                      <span class="ms-full-date">{{ ms.due_date ? formatDate(ms.due_date) : '未设日期' }}</span>
                    </div>
                    <p class="ms-full-desc">
                      {{ ms.description }}
                    </p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 成品管理页面 -->
        <div
          v-if="activeTab === 'deliverables'"
          class="tab-pane deliverables-pane"
        >
          <div class="pane-header-actions">
            <h3>项目成品 (Final Deliverables)</h3>
            <button class="btn btn-primary">
              添加成品
            </button>
          </div>
          
          <div
            v-if="currentProject.deliverables?.length"
            class="deliverables-grid"
          >
            <div
              v-for="dl in currentProject.deliverables"
              :key="dl.id"
              class="dl-card"
            >
              <div class="dl-preview">
                <!-- 占位符预览 -->
                <div class="preview-placeholder">
                  {{ getFileIcon(dl.file_type) }}
                </div>
              </div>
              <div class="dl-info">
                <div
                  class="dl-name"
                  :title="dl.name"
                >
                  {{ dl.name }}
                </div>
                <div class="dl-meta">
                  <span class="dl-v">{{ dl.version || 'Final' }}</span>
                  <span class="dl-date">{{ formatDate(dl.created_at) }}</span>
                </div>
                <div class="dl-actions">
                  <button class="btn-outline">
                    复制路径
                  </button>
                  <button class="btn-outline">
                    打开位置
                  </button>
                </div>
              </div>
            </div>
          </div>
          <div
            v-else
            class="empty-state-large"
          >
            <div class="empty-icon">
              🏆
            </div>
            <p>还没有标记为成品的资产</p>
            <button class="btn btn-outline">
              从素材中选择成品
            </button>
          </div>
        </div>

        <!-- 操作日志页面 -->
        <div
          v-if="activeTab === 'activity'"
          class="tab-pane activity-pane"
        >
          <div class="activity-timeline">
            <div
              v-for="log in sortedWorkLogs"
              :key="log.id"
              class="activity-item"
            >
              <div class="activity-time">
                {{ formatDateTime(log.created_at) }}
              </div>
              <div class="activity-marker"></div>
              <div class="activity-content">
                <span class="activity-text">{{ log.content }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 设置/编辑项目对话框 (复用 AddProjectDialog) -->
      <AddProjectDialog
        :visible="showSettingsDialog"
        :edit-mode="true"
        :initial-data="currentProject"
        @confirm="handleSettingsConfirm"
        @cancel="showSettingsDialog = false"
      />
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import ProjectViewSkeleton from './ProjectViewSkeleton.vue'
import AddProjectDialog from './project/AddProjectDialog.vue'
import { useProjectStore, type Project, type Milestone } from '@/stores/projectStore'

const projectStore = useProjectStore()

const props = defineProps<{
  projectId: string
}>()

const loading = ref(true)
const activeTab = ref('overview')
const showSettingsDialog = ref(false)
const materialSearchQuery = ref('')

const tabs = [
  { key: 'overview', label: '概览', icon: 'i-home' },
  { key: 'materials', label: '素材管理', icon: 'i-folder' },
  { key: 'workflow', label: '迭代 & 里程碑', icon: 'i-rocket' },
  { key: 'deliverables', label: '成品管理', icon: 'i-trophy' },
  { key: 'activity', label: '操作日志', icon: 'i-scroll' }
]

// 当前项目数据
const currentProject = computed(() => projectStore.currentProject as Project || {} as Project)

const progressPercent = computed(() => {
  if (!currentProject.value.milestones?.length) return 0
  const completed = currentProject.value.milestones.filter(m => m.completed).length
  return Math.round((completed / currentProject.value.milestones.length) * 100)
})

const sortedWorkLogs = computed(() => {
  if (!currentProject.value.work_logs) return []
  return [...currentProject.value.work_logs].sort((a, b) => 
    new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
  )
})

// 项目文件夹模拟数据 (实际应从 store 获取)
const projectFolders = ref<any[]>([])

const filteredProjectFolders = computed(() => {
  if (!materialSearchQuery.value) return projectFolders.value
  const query = materialSearchQuery.value.toLowerCase()
  return projectFolders.value.map(folder => ({
    ...folder,
    items: folder.items.filter((item: any) => item.name.toLowerCase().includes(query))
  })).filter(folder => folder.items.length > 0)
})

// 工具函数
const formatDate = (dateStr: string) => {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleDateString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit' })
}

const formatDateTime = (dateStr: string) => {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleString('zh-CN', { 
    month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' 
  })
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    'active': '进行中',
    'completed': '已完成',
    'archived': '已归档',
    'paused': '已暂停',
    'planning': '规划中'
  }
  return map[status] || status
}

const getProjectTypeText = (type: string) => {
  const map: Record<string, string> = {
    'photoshoot': '摄影/摄像',
    'document_edit': '文档编辑',
    'creative': '创意项目',
    'research': '研究项目',
    'archive_org': '归档整理',
    'custom': '自定义'
  }
  return map[type] || type
}

const getFileSymbol = (type: string) => {
  const map: Record<string, string> = {
    video: '#i-video',
    audio: '#i-audio',
    image: '#i-image',
    document: '#i-file',
    folder: '#i-folder',
    '视频': '#i-video',
    '音频': '#i-audio',
    '图片': '#i-image',
    '文档': '#i-file',
    '压缩包': '#i-box'
  }
  return map[type] || '#i-file'
}

const toggleFolder = (folder: any) => {
  folder.expanded = !folder.expanded
}

const openSettings = () => {
  showSettingsDialog.value = true
}

const handleSync = async () => {
  if (props.projectId) {
    try {
      const result = await projectStore.syncProjectFiles(props.projectId)
      console.log('Sync result:', result)
      // TODO: Show success notification
      await loadProjectFiles()
    } catch (error) {
      console.error('Sync failed:', error)
    }
  }
}

const toggleMilestone = async (ms: Milestone) => {
  if (props.projectId) {
    await projectStore.toggleMilestone(props.projectId, ms.id)
  }
}

const isDeliverable = (item: any) => {
  return currentProject.value.deliverables?.some(d => d.path === item.path)
}

const toggleDeliverable = async (item: any) => {
  if (props.projectId) {
    await projectStore.toggleDeliverable(props.projectId, item)
  }
}

const addMilestone = async () => {
  const title = prompt('请输入里程碑名称:')
  if (title && props.projectId) {
    const dueDate = prompt('请输入截止日期 (YYYY-MM-DD, 可选):')
    await projectStore.addMilestone(props.projectId, title, dueDate || undefined)
  }
}

const addIteration = async () => {
  const version = prompt('请输入版本号 (如 v1.0):')
  if (version && props.projectId) {
    const title = prompt('请输入迭代标题:')
    const description = prompt('请输入迭代描述 (可选):')
    await projectStore.addIteration(props.projectId, version, title || '', description || undefined)
  }
}

const handleSettingsConfirm = async (data: any) => {
  try {
    await projectStore.updateProject({ ...currentProject.value, ...data })
    showSettingsDialog.value = false
  } catch (error) {
    console.error('Update failed:', error)
  }
}

// 加载项目文件逻辑 (复用并优化)
const loadProjectFiles = async () => {
  try {
    const files = await projectStore.getProjectFiles(props.projectId)
    projectFolders.value = convertFilesToTree(files)
  } catch (error) {
    console.error('Failed to load project files:', error)
  }
}

const convertFilesToTree = (files: any[]) => {
  const folders: any[] = []
  const filesByType: Record<string, any[]> = {
    '视频素材': [],
    '音频素材': [],
    '设计文件': [],
    '其他文件': []
  }

  files.forEach(file => {
    const ext = file.path.split('.').pop()?.toLowerCase() || ''
    const type = getFileTypeByExtension(ext)
    const category = type === '视频' ? '视频素材' : 
                     type === '音频' ? '音频素材' : 
                     (type === '图片' || type === '文档') ? '设计文件' : '其他文件'
    filesByType[category].push({ ...file, type })
  })

  Object.entries(filesByType).forEach(([name, items]) => {
    if (items.length > 0) {
      folders.push({ id: name, name, expanded: true, items })
    }
  })
  return folders
}

const getFileTypeByExtension = (ext: string) => {
  if (['mp4', 'mov', 'avi'].includes(ext)) return '视频'
  if (['mp3', 'wav', 'flac'].includes(ext)) return '音频'
  if (['jpg', 'png', 'gif', 'webp'].includes(ext)) return '图片'
  if (['pdf', 'doc', 'docx', 'txt'].includes(ext)) return '文档'
  return '其他'
}

onMounted(async () => {
  if (props.projectId) {
    await projectStore.getProject(props.projectId)
    await loadProjectFiles()
  }
  loading.value = false
})
</script>

<style scoped>
.project-view {
  display: flex;
  flex-direction: column;
  height: 100%;
  background-color: var(--color-bg-base);
  color: var(--color-text-primary);
  overflow: hidden;
}

/* Header Section */
.project-header-section {
  padding: var(--spacing-lg) var(--spacing-xl);
  background-color: var(--color-bg-surface);
  border-bottom: 1px solid var(--color-border);
}

.header-main {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: var(--spacing-lg);
}

.project-title {
  font-size: 24px;
  font-weight: 700;
  margin: 0 0 8px 0;
  letter-spacing: -0.02em;
}

.project-status-group {
  display: flex;
  gap: var(--spacing-sm);
  align-items: center;
}

.status-badge {
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
}

.status-badge.active { background: rgba(16, 185, 129, 0.1); color: var(--color-success); }
.status-badge.planning { background: rgba(59, 130, 246, 0.1); color: var(--color-primary); }
.status-badge.paused { background: rgba(245, 158, 11, 0.1); color: var(--color-warning); }

.project-type-tag {
  font-size: 13px;
  color: var(--color-text-secondary);
}

.header-actions {
  display: flex;
  gap: var(--spacing-sm);
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  border-radius: var(--radius-md);
  border: 1px solid var(--color-border);
  background: var(--color-bg-base);
  cursor: pointer;
  transition: all 0.2s;
  font-size: 13px;
  color: var(--color-text-secondary);
}

.action-btn:hover {
  background: var(--color-bg-surface-hover);
  color: var(--color-text-primary);
  border-color: var(--color-text-secondary);
}

.action-btn svg {
  width: var(--icon-size-16);
  height: var(--icon-size-16);
}

.header-stats {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
  gap: var(--spacing-lg);
}

.stat-card {
  display: flex;
  flex-direction: column;
}

.stat-label {
  font-size: 11px;
  color: var(--color-text-tertiary);
  margin-bottom: 4px;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.stat-value {
  font-size: 18px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.progress-card {
  grid-column: span 2;
}

.progress-bar-container {
  height: 4px;
  background: var(--color-bg-base);
  border-radius: 2px;
  position: relative;
  margin-top: 8px;
  overflow: hidden;
}

.progress-bar-fill {
  height: 100%;
  background: var(--color-primary);
  border-radius: 2px;
  transition: width 0.3s ease;
}

.progress-text {
  position: absolute;
  right: 0;
  top: -16px;
  font-size: 11px;
  font-weight: 600;
  color: var(--color-primary);
}

/* Tabs */
.project-tabs {
  display: flex;
  padding: 0 var(--spacing-xl);
  background-color: var(--color-bg-surface);
  border-bottom: 1px solid var(--color-border);
  gap: var(--spacing-lg);
}

.tab-btn {
  padding: 12px 4px;
  background: none;
  border: none;
  border-bottom: 2px solid transparent;
  cursor: pointer;
  color: var(--color-text-secondary);
  font-weight: 500;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
}

.tab-btn:hover {
  color: var(--color-text-primary);
}

.tab-btn.active {
  color: var(--color-primary);
  border-bottom-color: var(--color-primary);
}

/* Content Area */
.project-content {
  flex: 1;
  overflow-y: auto;
  padding: var(--spacing-lg);
}

.pane-grid {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: var(--spacing-lg);
}

.content-section {
  background: var(--color-bg-surface);
  border-radius: var(--radius-lg);
  padding: var(--spacing-lg);
  margin-bottom: var(--spacing-lg);
  border: 1px solid var(--color-border);
  box-shadow: var(--shadow-sm);
}

.section-title {
  font-size: 15px;
  font-weight: 600;
  margin: 0 0 var(--spacing-md) 0;
  color: var(--color-text-primary);
}

/* Iteration Cards */
.iteration-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);
}

.iteration-card {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  padding: 12px;
  background: var(--color-bg-base);
  border-radius: var(--radius-md);
  border: 1px solid var(--color-border);
  transition: all 0.2s;
}

.iteration-card:hover {
  border-color: var(--color-text-secondary);
  transform: translateY(-1px);
}

.it-version {
  background: var(--color-primary);
  color: white;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 700;
}

.it-title { font-weight: 600; font-size: 13px; }
.it-date { font-size: 11px; color: var(--color-text-tertiary); }

/* Milestone Timeline */
.milestone-timeline {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
  position: relative;
  padding-left: 8px;
}

.milestone-timeline::before {
  content: '';
  position: absolute;
  left: 11px;
  top: 6px;
  bottom: 6px;
  width: 1px;
  background: var(--color-border);
}

.milestone-item {
  display: flex;
  gap: var(--spacing-md);
  position: relative;
}

.ms-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: var(--color-bg-base);
  border: 1px solid var(--color-text-secondary);
  z-index: 1;
  margin-top: 5px;
  box-sizing: content-box;
}

.milestone-item.completed .ms-dot {
  background: var(--color-success);
}

.ms-title { font-weight: 600; font-size: 0.875rem; }
.ms-date { font-size: 0.75rem; color: var(--color-text-tertiary); }

/* Deliverables Grid */
.deliverables-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: var(--spacing-lg);
}

.dl-card {
  background: var(--color-bg-base);
  border-radius: var(--radius-lg);
  overflow: hidden;
  border: 1px solid var(--color-border);
  transition: all 0.2s;
}

.dl-card:hover { 
  transform: translateY(-2px);
  border-color: var(--color-text-secondary);
  box-shadow: var(--shadow-md);
}

.dl-preview {
  height: 120px;
  background: var(--color-bg-input);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 32px;
  color: var(--color-text-tertiary);
}

.dl-info { padding: 12px; }
.dl-name { font-weight: 600; font-size: 13px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; color: var(--color-text-primary); }
.dl-meta { display: flex; justify-content: space-between; font-size: 11px; color: var(--color-text-tertiary); margin: 4px 0 8px; }
.dl-actions { display: flex; gap: 8px; }

.btn-outline {
  flex: 1;
  padding: 6px;
  font-size: 11px;
  border: 1px solid var(--color-border);
  background: transparent;
  border-radius: 4px;
  cursor: pointer;
  color: var(--color-text-secondary);
  transition: all 0.2s;
}

.btn-outline:hover {
  background: var(--color-bg-surface-hover);
  color: var(--color-text-primary);
  border-color: var(--color-text-secondary);
}

/* Activity Timeline */
.activity-timeline {
  display: flex;
  flex-direction: column;
}

.activity-item {
  display: flex;
  gap: var(--spacing-lg);
  padding-bottom: var(--spacing-lg);
  position: relative;
}

.activity-item::before {
  content: '';
  position: absolute;
  left: 95px;
  top: 20px;
  bottom: 0;
  width: 1px;
  background: var(--color-border);
}

.activity-item:last-child::before { display: none; }

.activity-time { width: 80px; font-size: 11px; color: var(--color-text-tertiary); text-align: right; padding-top: 2px; }
.activity-marker { 
  width: 8px; 
  height: 8px; 
  border-radius: 50%; 
  background: var(--color-bg-base); 
  border: 2px solid var(--color-primary); 
  z-index: 1; 
  margin-top: 5px; 
  box-sizing: border-box;
}
.activity-content { flex: 1; font-size: 13px; line-height: 1.5; color: var(--color-text-secondary); }

/* Utils */
.empty-state { color: var(--color-text-tertiary); font-size: 13px; padding: 16px 0; text-align: center; }
.empty-state-large { text-align: center; padding: 64px 32px; }
.empty-icon { font-size: 48px; margin-bottom: 16px; opacity: 0.5; }
.text-btn { background: none; border: none; color: var(--color-primary); cursor: pointer; font-size: 13px; }
.text-btn:hover { text-decoration: underline; }

/* Materials Styles (Brief) */
.materials-toolbar { display: flex; gap: var(--spacing-md); margin-bottom: var(--spacing-lg); }
.search-box input { 
  flex: 1; 
  padding: 8px 12px; 
  border-radius: var(--radius-md); 
  border: 1px solid var(--color-border); 
  background: var(--color-bg-input); 
  width: 300px; 
  font-size: 13px;
  color: var(--color-text-primary);
}
.search-box input:focus { outline: none; border-color: var(--color-primary); }

.filter-btn { 
  padding: 8px 16px; 
  border-radius: var(--radius-md); 
  border: 1px solid var(--color-border); 
  background: var(--color-bg-base); 
  cursor: pointer; 
  font-size: 13px;
  color: var(--color-text-secondary);
}
.filter-btn:hover { background: var(--color-bg-surface-hover); color: var(--color-text-primary); }
.filter-btn.active { background: var(--color-primary); color: white; border-color: var(--color-primary); }

.material-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  padding: 8px 12px;
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all 0.1s;
}

.material-item:hover { background: var(--color-bg-surface-hover); }
.item-icon { font-size: 20px; opacity: 0.8; }
.item-name { font-weight: 500; font-size: 13px; color: var(--color-text-primary); }
.item-meta { font-size: 11px; color: var(--color-text-tertiary); }
</style>
