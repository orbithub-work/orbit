<template>
  <div class="project-view">
    <!-- 模板引导页 -->
    <ProjectTemplateGuide
      v-if="!project && showTemplateGuide"
      @confirm="handleTemplateConfirm"
      @skip="handleTemplateSkip"
    />

    <!-- 项目详情页 -->
    <template v-else-if="project">
      <!-- 项目头部 -->
    <header class="project-header">
      <div class="header-left">
        <h1 class="project-title">{{ project?.name || '未选择项目' }}</h1>
        <div class="project-meta">
          <span v-if="project?.client">{{ project.client }}</span>
          <span v-if="project?.deadline">截止：{{ project.deadline }}</span>
        </div>
      </div>
      <div class="header-right">
        <div class="status-badge" :class="`status--${project?.status}`">
          {{ getStatusText(project?.status) }}
        </div>
        <button class="action-btn">
          <Icon name="more-vertical" size="md" />
        </button>
      </div>
    </header>

    <!-- 统计卡片 -->
    <div class="stats-cards">
      <div class="stat-card">
        <div class="stat-icon">📁</div>
        <div class="stat-content">
          <div class="stat-value">{{ project?.assets || 0 }}</div>
          <div class="stat-label">素材</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon">🎬</div>
        <div class="stat-content">
          <div class="stat-value">{{ projectEngineFiles.length }}</div>
          <div class="stat-label">工程</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon">✅</div>
        <div class="stat-content">
          <div class="stat-value">{{ projectDeliverables.length }}</div>
          <div class="stat-label">成品</div>
        </div>
      </div>
    </div>

    <!-- 主内容区 -->
    <div class="project-content">
      <!-- 素材区 -->
      <section class="content-section">
        <div class="section-header">
          <h3 class="section-title">
            <Icon name="image" size="md" />
            素材 ({{ projectAssets.length }})
          </h3>
          <button class="add-btn">
            <Icon name="plus" size="sm" />
            添加素材
          </button>
        </div>
        <div class="section-body">
          <EmptyState
            v-if="projectAssets.length === 0"
            icon="📁"
            title="暂无素材"
            description="点击「添加素材」开始"
            compact
          />
          <div v-else class="asset-grid">
            <div v-for="asset in projectAssets" :key="asset.id" class="asset-card">
              <div class="asset-thumb">
                <span class="file-icon">{{ getFileIcon(asset.type) }}</span>
              </div>
              <div class="asset-name">{{ asset.name }}</div>
            </div>
          </div>
        </div>
      </section>

      <!-- 工程文件区 -->
      <section class="content-section">
        <div class="section-header">
          <h3 class="section-title">
            <Icon name="file" size="md" />
            工程文件 ({{ projectEngineFiles.length }})
          </h3>
          <button class="add-btn">
            <Icon name="plus" size="sm" />
            添加工程
          </button>
        </div>
        <div class="section-body">
          <EmptyState
            v-if="projectEngineFiles.length === 0"
            icon="🎬"
            title="暂无工程文件"
            description="添加 PR/AE/剪映等工程文件"
            compact
          />
          <div v-else class="file-list">
            <div v-for="file in projectEngineFiles" :key="file.id" class="file-item">
              <div class="file-icon-large">{{ getEngineIcon(file.type) }}</div>
              <div class="file-info">
                <div class="file-name">{{ file.name }}</div>
                <div class="file-meta">{{ file.version }} · {{ formatDate(file.updated) }}</div>
              </div>
              <button class="open-btn">
                <Icon name="external-link" size="sm" />
                打开
              </button>
            </div>
          </div>
        </div>
      </section>

      <!-- 成品区 -->
      <section class="content-section">
        <div class="section-header">
          <h3 class="section-title">
            <Icon name="check-circle" size="md" />
            成品 ({{ projectDeliverables.length }})
          </h3>
          <button class="add-btn">
            <Icon name="plus" size="sm" />
            添加成品
          </button>
        </div>
        <div class="section-body">
          <EmptyState
            v-if="projectDeliverables.length === 0"
            icon="✅"
            title="暂无成品"
            description="导出的最终文件将显示在这里"
            compact
          />
          <div v-else class="file-list">
            <div v-for="file in projectDeliverables" :key="file.id" class="file-item">
              <div class="file-icon-large">{{ getFileIcon(file.type) }}</div>
              <div class="file-info">
                <div class="file-name">{{ file.name }}</div>
                <div class="file-meta">{{ file.version }} · {{ formatSize(file.size) }}</div>
              </div>
              <button class="open-btn">
                <Icon name="folder" size="sm" />
                定位
              </button>
            </div>
          </div>
        </div>
      </section>

      <!-- 项目日志 -->
      <ProjectLogTimeline :project-id="project?.id" />
    </div>
    </template>

    <!-- 空状态 -->
    <EmptyState
      v-else
      icon="📋"
      title="暂无项目"
      description="从左侧选择一个项目，或创建新项目"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import Icon from '@/components/common/Icon.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import ProjectTemplateGuide from '@/components/project/ProjectTemplateGuide.vue'
import ProjectLogTimeline from '@/components/project/ProjectLogTimeline.vue'

interface Project {
  id: string | number
  name: string
  client: string
  owner: string
  deadline: string
  status: string
  type: string
  assets: number
  references: string
  deliverables: string
  versions: string
}

const props = defineProps<{
  project: Project | undefined
}>()

const showTemplateGuide = computed(() => !props.project)

// Mock 数据 - 后续接入真实 API
const projectAssets = computed(() => [])
const projectEngineFiles = computed(() => [])
const projectDeliverables = computed(() => [])

function handleTemplateConfirm(templateId: string) {
  console.log('Selected template:', templateId)
  // TODO: 创建项目并应用模板
  // 这里应该调用 projectStore.createProject() 并传入模板 ID
}

function handleTemplateSkip() {
  console.log('Skip template, create blank project')
  // TODO: 创建空白项目
}

function getStatusText(status?: string) {
  const map: Record<string, string> = {
    active: '进行中',
    completed: '已完成',
    archived: '已归档'
  }
  return map[status || ''] || '未知'
}

function getFileIcon(type: string) {
  const icons: Record<string, string> = {
    image: '🖼️',
    video: '🎬',
    audio: '🎵',
    document: '📄'
  }
  return icons[type] || '📄'
}

function getEngineIcon(type: string) {
  const icons: Record<string, string> = {
    premiere: '🎬',
    aftereffects: '🎨',
    photoshop: '🖼️',
    davinci: '🎞️'
  }
  return icons[type] || '📄'
}

function formatDate(date: string) {
  return new Date(date).toLocaleDateString('zh-CN')
}

function formatSize(bytes: number) {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  if (bytes < 1024 * 1024 * 1024) return (bytes / 1024 / 1024).toFixed(1) + ' MB'
  return (bytes / 1024 / 1024 / 1024).toFixed(1) + ' GB'
}
</script>

<style scoped>
.project-view {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: #1b1c1f;
  overflow: hidden;
}

.project-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 24px 24px 16px;
  border-bottom: 1px solid #2b2b2f;
}

.header-left {
  flex: 1;
  min-width: 0;
}

.project-title {
  font-size: 20px;
  font-weight: 600;
  color: #e5e7eb;
  margin: 0 0 4px 0;
}

.project-meta {
  display: flex;
  gap: 12px;
  font-size: 13px;
  color: #6b7280;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.status-badge {
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.status--active {
  background: rgba(34, 197, 94, 0.15);
  color: #22c55e;
}

.status--completed {
  background: rgba(59, 130, 246, 0.15);
  color: #60a5fa;
}

.status--archived {
  background: rgba(107, 114, 128, 0.15);
  color: #9ca3af;
}

.action-btn {
  width: 32px;
  height: 32px;
  padding: 0;
  background: transparent;
  border: none;
  border-radius: 6px;
  color: #9ca3af;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.action-btn:hover {
  background: rgba(255, 255, 255, 0.08);
  color: #e5e7eb;
}

.stats-cards {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
  padding: 16px 24px;
  border-bottom: 1px solid #2b2b2f;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: #252526;
  border-radius: 12px;
}

.stat-icon {
  font-size: 32px;
  opacity: 0.8;
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 24px;
  font-weight: 600;
  color: #e5e7eb;
  line-height: 1;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 12px;
  color: #6b7280;
}

.project-content {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.content-section {
  background: #252526;
  border-radius: 12px;
  overflow: hidden;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.section-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 15px;
  font-weight: 600;
  color: #e5e7eb;
  margin: 0;
}

.add-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: rgba(59, 130, 246, 0.15);
  border: none;
  border-radius: 6px;
  color: #60a5fa;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
}

.add-btn:hover {
  background: rgba(59, 130, 246, 0.25);
}

.section-body {
  padding: 20px;
}

.asset-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
  gap: 12px;
}

.asset-card {
  background: #1e1e1e;
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  transition: transform 0.15s;
}

.asset-card:hover {
  transform: translateY(-2px);
}

.asset-thumb {
  width: 100%;
  aspect-ratio: 1;
  background: #2a2a2a;
  display: flex;
  align-items: center;
  justify-content: center;
}

.file-icon {
  font-size: 32px;
}

.asset-name {
  padding: 8px;
  font-size: 11px;
  color: #9ca3af;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.file-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.file-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: #1e1e1e;
  border-radius: 8px;
  transition: background 0.15s;
}

.file-item:hover {
  background: #2a2a2a;
}

.file-icon-large {
  font-size: 32px;
  flex-shrink: 0;
}

.file-info {
  flex: 1;
  min-width: 0;
}

.file-name {
  font-size: 13px;
  color: #e5e7eb;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-bottom: 2px;
}

.file-meta {
  font-size: 11px;
  color: #6b7280;
}

.open-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: transparent;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 6px;
  color: #9ca3af;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
  flex-shrink: 0;
}

.open-btn:hover {
  background: rgba(255, 255, 255, 0.08);
  border-color: rgba(255, 255, 255, 0.2);
  color: #e5e7eb;
}
</style>
