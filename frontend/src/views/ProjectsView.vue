<template>
  <div class="projects-view">
    <!-- Header -->
    <div class="page-header">
      <h1>项目管理</h1>
      <div class="header-actions">
        <button
          class="btn btn-primary"
          @click="showCreateDialog = true"
        >
          <Icon name="plus" size="sm" />
          新建项目
        </button>
      </div>
    </div>

    <!-- Projects List -->
    <div class="projects-container">
      <div
        v-if="loading"
        class="loading-container"
      >
        <div class="spinner"></div>
        <p>加载项目...</p>
      </div>

      <div
        v-else-if="projects.length === 0"
        class="empty-state"
      >
        <div class="empty-icon">
          <Icon name="folder" size="xl" />
        </div>
        <h3>暂无项目</h3>
        <p>点击"新建项目"按钮创建您的第一个项目</p>
        <button
          class="btn btn-primary"
          @click="showCreateDialog = true"
        >
          新建项目
        </button>
      </div>

      <div
        v-else
        class="projects-grid"
      >
        <div
          v-for="project in projects"
          :key="project.id"
          class="project-card"
          :class="{
            'status-active': project.status === 'active',
            'status-completed': project.status === 'completed',
            'status-archived': project.status === 'archived',
            'status-paused': project.status === 'paused'
          }"
        >
          <div class="project-header">
            <div class="project-icon">
              <Icon :name="getProjectIcon(project.project_type)" size="lg" />
            </div>
            <div class="project-title">
              <h3>{{ project.name }}</h3>
              <p class="project-type">
                {{ getProjectTypeText(project.project_type) }}
              </p>
            </div>
            <div class="project-actions">
              <button
                class="btn-icon"
                @click="editProject(project)"
              >
                <Icon name="pencil" size="sm" />
              </button>
              <button
                class="btn-icon"
                @click="deleteProject(project)"
              >
                <Icon name="trash" size="sm" />
              </button>
            </div>
          </div>

          <div class="project-details">
            <p
              v-if="project.description"
              class="project-description"
            >
              {{ project.description }}
            </p>

            <div class="project-meta">
              <div class="meta-item">
                <Icon name="calendar" size="sm" />
                <span v-if="project.deadline">截止: {{ formatDate(project.deadline) }}</span>
                <span v-else>无截止日期</span>
              </div>

              <div class="meta-item">
                <Icon name="document" size="sm" />
                {{ project.stats.file_count }} 个文件
              </div>

              <div class="meta-item">
                <Icon name="chart-bar" size="sm" />
                {{ formatSize(project.stats.total_size) }}
              </div>

              <div class="meta-item">
                <Icon name="archive" size="sm" />
                归档进度: {{ project.stats.archive_progress.toFixed(0) }}%
              </div>

              <div
                v-if="project.path"
                class="meta-item project-path"
              >
                <Icon name="folder" size="sm" />
                <span
                  class="path-text"
                  :title="project.path"
                >{{ project.path }}</span>
              </div>
            </div>

            <div class="project-status">
              <span
                class="status-badge"
                :class="`status-${project.status}`"
              >
                {{ getStatusText(project.status) }}
              </span>
            </div>
          </div>

          <div class="project-footer">
            <div class="footer-buttons">
              <button
                class="btn btn-small"
                @click="viewProject(project)"
              >
                查看详情
              </button>
              <button
                v-if="!project.path"
                class="btn btn-small btn-secondary"
                @click="handleSelectFolder(project)"
              >
                设置目录
              </button>
              <template v-else>
                <button
                  class="btn btn-small btn-secondary"
                  :disabled="scanningIds.has(project.id)"
                  @click="handleScanProject(project)"
                >
                  {{ scanningIds.has(project.id) ? '扫描中...' : '扫描项目' }}
                </button>
              </template>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Create Project Dialog -->
    <Teleport to="body">
      <div
        v-if="showCreateDialog"
        class="modal-overlay"
        @click="showCreateDialog = false"
      >
        <div
          class="modal"
          @click.stop
        >
          <div class="modal-header">
            <h2>新建项目</h2>
            <button
              class="btn-close"
              @click="showCreateDialog = false"
            >
              ×
            </button>
          </div>

          <div class="modal-body">
            <form @submit.prevent="handleCreateProject">
              <div class="form-group">
                <label for="projectName">项目名称</label>
                <input
                  id="projectName"
                  v-model="newProject.name"
                  type="text"
                  class="form-control"
                  placeholder="请输入项目名称"
                  required
                />
              </div>

              <div class="form-group">
                <label for="projectType">项目类型</label>
                <select
                  id="projectType"
                  v-model="newProject.project_type"
                  class="form-control"
                  required
                >
                  <option value="photoshoot">
                    摄影/摄像
                  </option>
                  <option value="document_edit">
                    文档编辑
                  </option>
                  <option value="creative">
                    创意项目
                  </option>
                  <option value="research">
                    研究项目
                  </option>
                  <option value="archive_org">
                    归档整理
                  </option>
                  <option value="custom">
                    自定义
                  </option>
                </select>
              </div>

              <div class="form-group">
                <label for="projectDescription">项目描述</label>
                <textarea
                  id="projectDescription"
                  v-model="newProject.description"
                  class="form-control"
                  placeholder="请输入项目描述"
                  rows="3"
                ></textarea>
              </div>

              <div class="form-group">
                <label for="projectDeadline">截止日期</label>
                <input
                  id="projectDeadline"
                  v-model="newProject.deadline"
                  type="date"
                  class="form-control"
                />
              </div>

              <div class="modal-footer">
                <button
                  type="button"
                  class="btn btn-secondary"
                  @click="showCreateDialog = false"
                >
                  取消
                </button>
                <button
                  type="submit"
                  class="btn btn-primary"
                  :disabled="createLoading"
                >
                  {{ createLoading ? '创建中...' : '创建' }}
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Edit Project Dialog -->
    <Teleport to="body">
      <div
        v-if="editingProject"
        class="modal-overlay"
        @click="editingProject = null"
      >
        <div
          class="modal"
          @click.stop
        >
          <div class="modal-header">
            <h2>编辑项目</h2>
            <button
              class="btn-close"
              @click="editingProject = null"
            >
              ×
            </button>
          </div>

          <div class="modal-body">
            <form @submit.prevent="handleUpdateProject">
              <div class="form-group">
                <label for="editProjectName">项目名称</label>
                <input
                  id="editProjectName"
                  v-model="editingProject.name"
                  type="text"
                  class="form-control"
                  placeholder="请输入项目名称"
                  required
                />
              </div>

              <div class="form-group">
                <label for="editProjectType">项目类型</label>
                <select
                  id="editProjectType"
                  v-model="editingProject.project_type"
                  class="form-control"
                  required
                >
                  <option value="photoshoot">
                    摄影/摄像
                  </option>
                  <option value="document_edit">
                    文档编辑
                  </option>
                  <option value="creative">
                    创意项目
                  </option>
                  <option value="research">
                    研究项目
                  </option>
                  <option value="archive_org">
                    归档整理
                  </option>
                  <option value="custom">
                    自定义
                  </option>
                </select>
              </div>

              <div class="form-group">
                <label for="editProjectDescription">项目描述</label>
                <textarea
                  id="editProjectDescription"
                  v-model="editingProject.description"
                  class="form-control"
                  placeholder="请输入项目描述"
                  rows="3"
                ></textarea>
              </div>

              <div class="form-group">
                <label for="editProjectDeadline">截止日期</label>
                <input
                  id="editProjectDeadline"
                  v-model="editingProject.deadline"
                  type="date"
                  class="form-control"
                />
              </div>

              <div class="form-group">
                <label for="editProjectStatus">项目状态</label>
                <select
                  id="editProjectStatus"
                  v-model="editingProject.status"
                  class="form-control"
                  required
                >
                  <option value="active">
                    进行中
                  </option>
                  <option value="paused">
                    暂停
                  </option>
                  <option value="completed">
                    已完成
                  </option>
                  <option value="archived">
                    已归档
                  </option>
                </select>
              </div>

              <div class="form-group">
                <label>项目目录</label>
                <div class="input-group">
                  <input
                    v-model="editingProject.path"
                    type="text"
                    class="form-control"
                    placeholder="未设置目录"
                    readonly
                  />
                  <button
                    type="button"
                    class="btn btn-secondary btn-append"
                    @click="handleEditSelectFolder"
                  >
                    选择
                  </button>
                </div>
              </div>

              <div class="modal-footer">
                <button
                  type="button"
                  class="btn btn-secondary"
                  @click="editingProject = null"
                >
                  取消
                </button>
                <button
                  type="submit"
                  class="btn btn-primary"
                  :disabled="updateLoading"
                >
                  {{ updateLoading ? '保存中...' : '保存' }}
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Icon from '@/components/common/Icon.vue'
import { useProjectStore } from '@/stores/projectStore'
import { showNotification } from '@/utils/notifications'

const projectStore = useProjectStore()
const router = useRouter()

// State
const projects = ref<any[]>([])
const loading = ref(false)
const showCreateDialog = ref(false)
const editingProject = ref<any>(null)
const createLoading = ref(false)
const updateLoading = ref(false)
const scanningIds = ref<Set<string>>(new Set())

const newProject = ref({
  name: '',
  description: '',
  project_type: 'custom',
  deadline: ''
})

// Lifecycle
onMounted(() => {
  loadProjects()
})

// Methods
async function loadProjects() {
  loading.value = true
  try {
    await projectStore.fetchProjects()
    projects.value = projectStore.projects
  } catch (error) {
    console.error('Failed to load projects:', error)
    showNotification('加载项目失败', 'error')
  } finally {
    loading.value = false
  }
}

async function handleSelectFolder(project: any) {
  const path = await projectStore.selectProjectFolder()
  if (path) {
    try {
      await projectStore.updateProjectPath(project.id, path)
      showNotification('目录设置成功', 'success')
      loadProjects()
    } catch (error) {
      showNotification('设置目录失败', 'error')
    }
  }
}

async function handleEditSelectFolder() {
  const path = await projectStore.selectProjectFolder()
  if (path && editingProject.value) {
    editingProject.value.path = path
  }
}

async function handleScanProject(project: any) {
  if (!project.path) {
    showNotification('请先设置项目目录', 'warning')
    return
  }

  scanningIds.value.add(project.id)
  try {
    await projectStore.scanProject(project.id, project.path)
    showNotification('项目扫描已启动', 'success')
  } catch (error) {
    console.error('Failed to scan project:', error)
    showNotification('扫描项目失败', 'error')
  } finally {
    // 简单起见，这里直接移除。实际可能需要监听进度事件。
    setTimeout(() => {
      scanningIds.value.delete(project.id)
    }, 2000)
  }
}

async function handleCreateProject() {
  createLoading.value = true
  try {
    const projectData = {
      ...newProject.value,
      deadline: newProject.value.deadline || null
    }

    await projectStore.createProject(
      projectData.name,
      projectData.project_type
    )

    // Reset form
    newProject.value = {
      name: '',
      description: '',
      project_type: 'custom',
      deadline: ''
    }

    showCreateDialog.value = false
    showNotification('项目创建成功', 'success')
    loadProjects()
  } catch (error) {
    console.error('Failed to create project:', error)
    showNotification('创建项目失败', 'error')
  } finally {
    createLoading.value = false
  }
}

function editProject(project: any) {
  editingProject.value = JSON.parse(JSON.stringify(project))
}

async function handleUpdateProject() {
  updateLoading.value = true
  try {
    await projectStore.updateProject(editingProject.value)
    showNotification('项目更新成功', 'success')
    editingProject.value = null
    loadProjects()
  } catch (error) {
    console.error('Failed to update project:', error)
    showNotification('更新项目失败', 'error')
  } finally {
    updateLoading.value = false
  }
}

async function deleteProject(project: any) {
  if (!confirm(`确定要删除项目"${project.name}"吗？`)) {
    return
  }

  try {
    await projectStore.deleteProject(project.id)
    showNotification('项目删除成功', 'success')
    loadProjects()
  } catch (error) {
    console.error('Failed to delete project:', error)
    showNotification('删除项目失败', 'error')
  }
}

function viewProject(project: any) {
  router.push({ name: 'ProjectDetail', params: { id: project.id } })
}

function getProjectIcon(type: string) {
  const icons: Record<string, string> = {
    photoshoot: 'camera',
    document_edit: 'document',
    creative: 'sparkles',
    research: 'beaker',
    archive_org: 'archive',
    custom: 'folder'
  }
  return icons[type] || 'folder'
}

function getProjectTypeText(type: string) {
  const texts: Record<string, string> = {
    photoshoot: '摄影/摄像',
    document_edit: '文档编辑',
    creative: '创意项目',
    research: '研究项目',
    archive_org: '归档整理',
    custom: '自定义'
  }
  return texts[type] || '自定义'
}

function getStatusText(status: string) {
  const texts: Record<string, string> = {
    active: '进行中',
    paused: '暂停',
    completed: '已完成',
    archived: '已归档'
  }
  return texts[status] || '未知'
}

function formatDate(dateStr: string) {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN')
}

function formatSize(bytes: number) {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}
</script>

<style scoped>
.projects-view {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}

.page-header h1 {
  font-size: 24px;
  color: #1a1a1a;
  margin: 0;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.projects-container {
  min-height: 400px;
}

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid #f3f3f3;
  border-top: 3px solid #4a90d9;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 16px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: #666;
}

.empty-icon {
  font-size: 64px;
  margin-bottom: 16px;
  opacity: 0.5;
}

.empty-state h3 {
  font-size: 20px;
  margin: 16px 0 8px;
  color: #333;
}

.projects-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 20px;
}

.project-card {
  background: white;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: transform 0.2s, box-shadow 0.2s;
  display: flex;
  flex-direction: column;
}

.project-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
}

.project-card.status-active {
  border-left: 4px solid #4CAF50;
}

.project-card.status-completed {
  border-left: 4px solid #2196F3;
}

.project-card.status-archived {
  border-left: 4px solid #9E9E9E;
  opacity: 0.8;
}

.project-card.status-paused {
  border-left: 4px solid #FFC107;
}

.project-header {
  display: flex;
  align-items: flex-start;
  margin-bottom: 16px;
}

.project-icon {
  font-size: 48px;
  margin-right: 16px;
  flex-shrink: 0;
}

.project-title {
  flex: 1;
  min-width: 0;
}

.project-title h3 {
  font-size: 18px;
  margin: 0 0 4px 0;
  color: #1a1a1a;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.project-type {
  font-size: 12px;
  color: #666;
  margin: 0;
}

.project-actions {
  display: flex;
  gap: 4px;
}

.btn-icon {
  background: none;
  border: none;
  padding: 4px;
  cursor: pointer;
  color: #999;
  font-size: 16px;
  border-radius: 4px;
  transition: background-color 0.2s, color 0.2s;
}

.btn-icon:hover {
  background-color: #f0f0f0;
  color: #333;
}

.project-details {
  flex: 1;
  margin-bottom: 16px;
}

.project-description {
  font-size: 14px;
  color: #666;
  margin: 0 0 16px 0;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.project-meta {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 16px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  color: #666;
}

.meta-item .icon {
  font-size: 14px;
}

.project-status {
  margin-bottom: 16px;
}

.status-badge {
  display: inline-block;
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.status-active {
  background-color: #E8F5E9;
  color: #2E7D32;
}

.status-completed {
  background-color: #E3F2FD;
  color: #1565C0;
}

.status-archived {
  background-color: #F5F5F5;
  color: #424242;
}

.status-paused {
  background-color: #FFF8E1;
  color: #F57C00;
}

.project-footer {
  text-align: center;
}

.footer-buttons {
  display: flex;
  justify-content: center;
  gap: 8px;
}

.btn-small {
  padding: 6px 16px;
  font-size: 12px;
}

.btn-secondary {
  background-color: #f5f5f5;
  color: #333;
  border: 1px solid #ddd;
}

.btn-secondary:hover {
  background-color: #e0e0e0;
}

.project-path {
  color: #4a90d9;
  font-style: italic;
}

.path-text {
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.input-group {
  display: flex;
  gap: 8px;
}

.btn-append {
  flex-shrink: 0;
}

/* Modal styles */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
}

.modal {
  background: white;
  border-radius: 8px;
  max-width: 500px;
  width: 100%;
  max-height: 90vh;
  overflow-y: auto;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #e0e0e0;
}

.modal-header h2 {
  font-size: 18px;
  margin: 0;
  color: #1a1a1a;
}

.btn-close {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #999;
  padding: 0;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  transition: background-color 0.2s, color 0.2s;
}

.btn-close:hover {
  background-color: #f0f0f0;
  color: #333;
}

.modal-body {
  padding: 20px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
  color: #333;
}

.form-control {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  box-sizing: border-box;
}

.form-control:focus {
  outline: none;
  border-color: #4a90d9;
  box-shadow: 0 0 0 2px rgba(74, 144, 217, 0.2);
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  padding: 20px;
  border-top: 1px solid #e0e0e0;
}

/* Responsive */
@media (max-width: 768px) {
  .projects-grid {
    grid-template-columns: 1fr;
  }

  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }
}
</style>
