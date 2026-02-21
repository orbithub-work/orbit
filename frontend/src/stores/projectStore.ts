// 项目 Store - 项目管理状态
//
// 提供项目的增删改查、文件关联、统计信息等功能

import { defineStore } from 'pinia'
import { ref } from 'vue'
import { apiCall } from '@/services/api'

// 项目统计信息
export interface ProjectStats {
  file_count: number
  total_size: number
  archive_progress: number
}

// 工作日志
export interface WorkLog {
  id: string
  date: string
  content: string
  created_at: string
  updated_at: string
}

// 里程碑
export interface Milestone {
  id: string
  title: string
  description?: string
  due_date?: string
  completed: boolean
  completed_at?: string
}

// 迭代版本
export interface Iteration {
  id: string
  version: string
  title: string
  description?: string
  file_ids: string[]
  created_at: string
}

// 交付物
export interface Deliverable {
  id: string
  name: string
  path: string
  file_type: string
  version?: string
  created_at: string
}

// 项目实体
export interface Project {
  id: string
  name: string
  description?: string
  project_type: string
  status: string
  color?: string
  icon?: string
  custom_icon_path?: string
  path?: string
  auto_sync: boolean
  file_ids: string[]
  collection_ids: string[]
  tags: string[]
  start_date?: string
  end_date?: string
  deadline?: string
  created_at: string
  updated_at: string
  work_logs: WorkLog[]
  milestones: Milestone[]
  iterations: Iteration[]
  deliverables: Deliverable[]
  metadata: Record<string, any>
  stats: ProjectStats
}

export const useProjectStore = defineStore('project', () => {
  // 状态
  const projects = ref<Project[]>([])
  const currentProject = ref<Project | null>(null)
  const loading = ref(false)

  // 加载所有项目
  async function fetchProjects() {
    loading.value = true
    try {
      const result = await apiCall<Project[]>('list_projects')
      projects.value = result || []
    } catch (error: any) {
      const msg = error?.message || 'Unknown error'
      console.error('Failed to load projects (sanitized):', msg)
      // 浏览器模式下使用模拟数据
      projects.value = [
        {
          id: '1',
          name: '示例项目',
          project_type: 'default',
          status: 'active',
          auto_sync: true,
          file_ids: [],
          collection_ids: [],
          tags: [],
          work_logs: [],
          milestones: [],
          iterations: [],
          deliverables: [],
          metadata: {},
          stats: { file_count: 0, total_size: 0, archive_progress: 0 },
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString()
        }
      ]
    } finally {
      loading.value = false
    }
  }

  // 获取单个项目
  async function getProject(id: string) {
    loading.value = true
    try {
      const result = await apiCall<Project | null>('get_project', { id })
      currentProject.value = result
      return result
    } catch (error) {
      console.error('Failed to get project:', error)
      return null
    } finally {
      loading.value = false
    }
  }

  // 创建项目
  async function createProject(name: string, project_type: string) {
    loading.value = true
    try {
      const result = await apiCall<Project>('create_project', {
        name,
        project_type
      })
      projects.value.push(result)
      return result
    } catch (error) {
      console.error('Failed to create project:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  // 更新项目
  async function updateProject(project: Project) {
    loading.value = true
    try {
      await apiCall('update_project', { project })
      const index = projects.value.findIndex(p => p.id === project.id)
      if (index !== -1) {
        projects.value[index] = project
      }
      if (currentProject.value?.id === project.id) {
        currentProject.value = project
      }
    } catch (error) {
      console.error('Failed to update project:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  // 删除项目
  async function deleteProject(id: string) {
    loading.value = true
    try {
      await apiCall('delete_project', { id })
      projects.value = projects.value.filter(p => p.id !== id)
      if (currentProject.value?.id === id) {
        currentProject.value = null
      }
    } catch (error) {
      console.error('Failed to delete project:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  // 更新项目状态
  async function updateProjectStatus(id: string, status: string) {
    loading.value = true
    try {
      await apiCall('update_project_status', { id, status })
      const project = projects.value.find(p => p.id === id)
      if (project) {
        project.status = status
        if (status === 'completed' && !project.end_date) {
          project.end_date = new Date().toISOString().split('T')[0]
        }
        project.updated_at = new Date().toISOString()
      }
      if (currentProject.value?.id === id) {
        currentProject.value.status = status
        if (status === 'completed' && !currentProject.value.end_date) {
          currentProject.value.end_date = new Date().toISOString().split('T')[0]
        }
        currentProject.value.updated_at = new Date().toISOString()
      }
    } catch (error) {
      console.error('Failed to update project status:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  // 获取项目统计信息
  async function getProjectStats(id: string) {
    loading.value = true
    try {
      const result = await apiCall<ProjectStats>('get_project_stats', { id })
      return result
    } catch (error) {
      console.error('Failed to get project stats:', error)
      return { file_count: 0, total_size: 0, archive_progress: 0 }
    } finally {
      loading.value = false
    }
  }

  // 添加文件到项目
  async function addFileToProject(projectId: string, fileId: string) {
    loading.value = true
    try {
      await apiCall('add_file_to_project', { projectId, fileId })
      const project = projects.value.find(p => p.id === projectId)
      if (project) {
        if (!project.file_ids.includes(fileId)) {
          project.file_ids.push(fileId)
        }
        const stats = await getProjectStats(projectId)
        project.stats = stats
        project.updated_at = new Date().toISOString()
      }
      if (currentProject.value?.id === projectId) {
        if (!currentProject.value.file_ids.includes(fileId)) {
          currentProject.value.file_ids.push(fileId)
        }
        const stats = await getProjectStats(projectId)
        currentProject.value.stats = stats
        currentProject.value.updated_at = new Date().toISOString()
      }
    } catch (error) {
      console.error('Failed to add file to project:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  // 从项目移除文件
  async function removeFileFromProject(projectId: string, fileId: string) {
    loading.value = true
    try {
      await apiCall('remove_file_from_project', { projectId, fileId })
      const project = projects.value.find(p => p.id === projectId)
      if (project) {
        project.file_ids = project.file_ids.filter(id => id !== fileId)
        const stats = await getProjectStats(projectId)
        project.stats = stats
        project.updated_at = new Date().toISOString()
      }
      if (currentProject.value?.id === projectId) {
        currentProject.value.file_ids = currentProject.value.file_ids.filter(id => id !== fileId)
        const stats = await getProjectStats(projectId)
        currentProject.value.stats = stats
        currentProject.value.updated_at = new Date().toISOString()
      }
    } catch (error) {
      console.error('Failed to remove file from project:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  // 更新项目截止日期
  async function updateProjectDeadline(id: string, deadline?: string) {
    loading.value = true
    try {
      await apiCall('update_project_deadline', { id, deadline })
      const project = projects.value.find(p => p.id === id)
      if (project) {
        project.deadline = deadline
        project.updated_at = new Date().toISOString()
      }
      if (currentProject.value?.id === id) {
        currentProject.value.deadline = deadline
        currentProject.value.updated_at = new Date().toISOString()
      }
    } catch (error) {
      console.error('Failed to update project deadline:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  // 里程碑管理
  // 添加里程碑
  async function addMilestone(projectId: string, title: string, dueDate?: string) {
    if (!currentProject.value || currentProject.value.id !== projectId) return
    
    const newMilestone: Milestone = {
      id: crypto.randomUUID(),
      title,
      due_date: dueDate,
      completed: false
    }
    
    const updatedProject = {
      ...currentProject.value,
      milestones: [...(currentProject.value.milestones || []), newMilestone]
    }
    
    await updateProject(updatedProject)
  }

  // 切换里程碑完成状态
  async function toggleMilestone(projectId: string, milestoneId: string) {
    if (!currentProject.value || currentProject.value.id !== projectId) return
    
    const milestones = currentProject.value.milestones.map(m => {
      if (m.id === milestoneId) {
        const completed = !m.completed
        return {
          ...m,
          completed,
          completed_at: completed ? new Date().toISOString() : undefined
        }
      }
      return m
    })
    
    await updateProject({ ...currentProject.value, milestones })
  }

  // 迭代管理
  // 添加迭代版本
  async function addIteration(projectId: string, version: string, title: string, description?: string) {
    if (!currentProject.value || currentProject.value.id !== projectId) return
    
    const newIteration: Iteration = {
      id: crypto.randomUUID(),
      version,
      title,
      description,
      file_ids: [],
      created_at: new Date().toISOString()
    }
    
    const updatedProject = {
      ...currentProject.value,
      iterations: [...(currentProject.value.iterations || []), newIteration]
    }
    
    await updateProject(updatedProject)
  }

  // 成品管理
  // 切换文件为成品状态
  async function toggleDeliverable(projectId: string, file: any) {
    if (!currentProject.value || currentProject.value.id !== projectId) return
    
    const deliverables = currentProject.value.deliverables || []
    const index = deliverables.findIndex(d => d.path === file.path)
    
    let newDeliverables
    if (index === -1) {
      const newDeliverable: Deliverable = {
        id: crypto.randomUUID(),
        name: file.name,
        path: file.path,
        file_type: file.file_type || 'unknown',
        created_at: new Date().toISOString()
      }
      newDeliverables = [...deliverables, newDeliverable]
    } else {
      newDeliverables = deliverables.filter(d => d.path !== file.path)
    }
    
    await updateProject({ ...currentProject.value, deliverables: newDeliverables })
  }

  // 从文件夹创建项目
  async function createProjectFromFolder(folderPath: string, autoSync: boolean = true) {
    loading.value = true
    try {
      const result = await apiCall<Project>('create_project_from_folder', {
        folderPath
      })
      projects.value.push(result)
      return result
    } catch (error) {
      console.error('Failed to create project from folder:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  // 获取项目文件列表
  async function getProjectFiles(projectId: string) {
    loading.value = true
    try {
      interface FileItemDto {
        id: string
        name: string
        path: string
        size: number
        file_type?: string
        is_directory?: boolean
        modified_at: string
        thumbnail_path?: string
      }
      const result = await apiCall<FileItemDto[]>('get_project_files', { projectId })
      return result
    } catch (error) {
      console.error('Failed to get project files:', error)
      return []
    } finally {
      loading.value = false
    }
  }

  // 同步项目文件
  async function syncProjectFiles(projectId: string) {
    loading.value = true
    try {
      interface SyncResult {
        added_count: number
        removed_count: number
        updated_count: number
        total_files: number
      }
      const result = await apiCall<SyncResult>('sync_project_files', { projectId })
      console.log('Sync result:', result)
      await fetchProjects()
      if (currentProject.value?.id === projectId) {
        const updated = projects.value.find(p => p.id === projectId)
        if (updated) currentProject.value = updated
      }
      return result
    } catch (error) {
      console.error('Failed to sync project files:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  // 选择项目文件夹
  async function selectProjectFolder() {
    try {
      const path = await apiCall<string>('select_project_folder')
      return path
    } catch (error) {
      console.error('Failed to select project folder:', error)
      return ''
    }
  }

  // 更新项目路径
  async function updateProjectPath(id: string, path: string) {
    loading.value = true
    try {
      await apiCall('update_project_path', { id, path })
      const project = projects.value.find(p => p.id === id)
      if (project) {
        project.path = path
        project.updated_at = new Date().toISOString()
      }
      if (currentProject.value?.id === id) {
        currentProject.value.path = path
        currentProject.value.updated_at = new Date().toISOString()
      }
    } catch (error) {
      console.error('Failed to update project path:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  // 扫描项目
  async function scanProject(id: string, path: string) {
    loading.value = true
    try {
      await apiCall('scan_project', { id, path })
    } catch (error) {
      console.error('Failed to scan project:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  return {
    projects,
    currentProject,
    loading,
    fetchProjects,
    getProject,
    createProject,
    createProjectFromFolder,
    getProjectFiles,
    syncProjectFiles,
    selectProjectFolder,
    updateProjectPath,
    scanProject,
    updateProject,
    deleteProject,
    updateProjectStatus,
    getProjectStats,
    addFileToProject,
    removeFileFromProject,
    updateProjectDeadline,
    addMilestone,
    toggleMilestone,
    addIteration,
    toggleDeliverable
  }
})
