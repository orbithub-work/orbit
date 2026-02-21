// 文件 Store - 文件管理状态
// 
// 提供文件列表、文件夹树、搜索、筛选、排序等功能

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { apiCall } from '@/services/api'

// 文件项
export interface FileItem {
  id: string
  name: string
  path: string
  size: number
  file_type?: string
  is_directory?: boolean
  modified_at: string | Date
  thumbnail_path?: string | null
  width?: number
  height?: number
  rating?: number
}

// 文件夹节点
export interface FolderNode {
  id: string
  name: string
  path: string
  children?: FolderNode[]
}

export const useFileStore = defineStore('file', () => {
  const currentPath = ref('')
  const files = ref<FileItem[]>([])
  const selectedFiles = ref<FileItem[]>([])
  const folderTree = ref<FolderNode[]>([])
  const viewMode = ref<'list' | 'grid'>('grid')
  const searchQuery = ref('')
  const filterType = ref<string | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)
  const projectId = ref<string>('')

  // 导航历史
  const navigationHistory = ref<string[]>(['/'])
  const historyIndex = ref(0)
  const sortField = ref<string>('modified_at')
  const sortOrder = ref<'asc' | 'desc'>('desc')

  // 过滤后的文件列表
  const filteredFiles = computed(() => {
    let result = [...files.value]

    if (filterType.value) {
      result = result.filter(file => {
        if (filterType.value === 'folder') {
          return file.is_directory
        }
        return file.file_type === filterType.value
      })
    }

    if (searchQuery.value) {
      const query = searchQuery.value.toLowerCase()
      result = result.filter(file =>
        file.name.toLowerCase().includes(query)
      )
    }

    result.sort((a, b) => {
      const aValue = a[sortField.value as keyof FileItem]
      const bValue = b[sortField.value as keyof FileItem]

      if (aValue === undefined || bValue === undefined) return 0

      if (typeof aValue === 'string' && typeof bValue === 'string') {
        return sortOrder.value === 'asc'
          ? aValue.localeCompare(bValue)
          : bValue.localeCompare(aValue)
      }

      if (typeof aValue === 'number' && typeof bValue === 'number') {
        return sortOrder.value === 'asc'
          ? aValue - bValue
          : bValue - aValue
      }

      return 0
    })

    return result
  })

  // 设置排序字段
  function setSort(field: string) {
    if (sortField.value === field) {
      sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc'
    } else {
      sortField.value = field
      sortOrder.value = 'asc'
    }
  }

  // 切换排序顺序
  function toggleSort(field: string) {
    if (sortField.value === field) {
      sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc'
    } else {
      sortField.value = field
      sortOrder.value = 'asc'
    }
  }

  // 计算属性：是否可以后退/前进
  const canGoBack = computed(() => historyIndex.value > 0)
  const canGoForward = computed(() => historyIndex.value < navigationHistory.value.length - 1)

  function setProjectId(id: string) {
    projectId.value = (id || '').trim()
  }

  async function loadFiles(path?: string) {
    loading.value = true
    error.value = null

    try {
      if (typeof path === 'string') {
        currentPath.value = path
        if (path !== navigationHistory.value[historyIndex.value]) {
          navigationHistory.value = navigationHistory.value.slice(0, historyIndex.value + 1)
          navigationHistory.value.push(path)
          historyIndex.value = navigationHistory.value.length - 1
        }
      }

      interface AssetDto {
        id: string
        path: string
        size: number
        mtime: number
        media_meta?: string
        user_rating?: number
        width?: number
        height?: number
      }
      
      // 如果没有 projectId，暂时返回空列表
      // TODO: 等后端修复 /api/v1/assets 接口后再启用
      if (!projectId.value) {
        files.value = []
        return
      }
      
      // 使用统一的 list_assets 接口
      const params: any = {
        projectId: projectId.value,
        limit: 1000
      }
      
      const result = await apiCall<{ items: AssetDto[], total: number }>('list_assets', params)

      const items = result.items || result as any as AssetDto[]

      files.value = items.map((asset): FileItem => {
        const fileName = asset.path.split('/').pop() || asset.path
        const ext = fileName.split('.').pop()?.toLowerCase() || ''
        
        let fileType = 'file'
        if (['jpg', 'jpeg', 'png', 'gif', 'webp', 'bmp', 'svg'].includes(ext)) fileType = 'image'
        else if (['mp4', 'avi', 'mov', 'mkv', 'webm'].includes(ext)) fileType = 'video'
        else if (['mp3', 'wav', 'flac', 'aac', 'ogg'].includes(ext)) fileType = 'audio'
        else if (['pdf', 'doc', 'docx', 'xls', 'xlsx', 'ppt', 'pptx'].includes(ext)) fileType = 'document'
        
        let thumbnailPath: string | undefined
        let width: number | undefined = asset.width
        let height: number | undefined = asset.height
        if (asset.media_meta) {
          try {
            const meta = JSON.parse(asset.media_meta)
            thumbnailPath = meta.thumbnail_path
            if (!width) width = meta.width
            if (!height) height = meta.height
          } catch {
            // ignore parse error
          }
        }
        
        return {
          id: asset.id,
          name: fileName,
          path: asset.path,
          size: asset.size,
          file_type: fileType,
          is_directory: false,
          modified_at: new Date(asset.mtime * 1000),
          thumbnail_path: thumbnailPath,
          width,
          height,
          rating: asset.user_rating
        }
      })
    } catch (e) {
      error.value = e instanceof Error ? e.message : '加载文件失败'
      console.error('Failed to load files:', e)
    } finally {
      loading.value = false
    }
  }

  // 加载文件夹树
  async function loadFolderTree() {
    loading.value = true
    error.value = null

    try {
      const result = await apiCall<any[]>('list_folders', {
        recursive: true
      })

      folderTree.value = (result || []).map((item: any) => ({
        id: item.id || item.path || item.name,
        name: item.name || item.path || '',
        path: item.path || '',
        children: Array.isArray(item.children) ? item.children : undefined
      }))
    } catch (e) {
      error.value = e instanceof Error ? e.message : '加载文件夹树失败'
      console.error('Failed to load folder tree:', e)
    } finally {
      loading.value = false
    }
  }

  // 获取统计信息
  async function fetchStats() {
    try {
      // 暂时模拟，实际可能需要调用后端 API
      console.log('Fetching file stats...')
    } catch (e) {
      console.error('Failed to fetch stats:', e)
    }
  }

  // 搜索文件
  async function searchFiles(query: string) {
    if (!query.trim()) {
      searchQuery.value = ''
      await loadFiles()
      return
    }

    loading.value = true
    error.value = null
    searchQuery.value = query

    try {
      interface SearchResultDto {
        files: FileItem[]
        total: number
        query: string
      }
      const result = await apiCall<SearchResultDto>('search_files', {
        query,
        file_types: filterType.value ? [filterType.value] : undefined,
        max_results: 100
      })
      files.value = result.files
    } catch (e) {
      error.value = e instanceof Error ? e.message : '搜索文件失败'
      console.error('Failed to search files:', e)
    } finally {
      loading.value = false
    }
  }

  // 高级搜索文件
  async function searchFilesAdvanced(filters: {
    namePattern?: string
    fileTypes?: string[]
    minSize?: number
    maxSize?: number
    dateFrom?: string
    dateTo?: string
    tags?: string[]
  }) {
    loading.value = true
    error.value = null
    searchQuery.value = 'advanced'

    try {
      interface SearchResultDto {
        files: FileItem[]
        total: number
        query: string
      }
      const result = await apiCall<SearchResultDto>('advanced_search', {
        query: filters.namePattern,
        file_types: filters.fileTypes,
        tags: filters.tags,
        date_from: filters.dateFrom,
        date_to: filters.dateTo,
        min_size: filters.minSize,
        max_size: filters.maxSize,
        is_archived: false,
        rating: undefined,
      })
      files.value = result.files
    } catch (e) {
      error.value = e instanceof Error ? e.message : '高级搜索失败'
      console.error('Failed to perform advanced search:', e)
    } finally {
      loading.value = false
    }
  }

  // 删除文件
  async function deleteFiles(ids: string[]) {
    loading.value = true
    error.value = null

    try {
      for (const id of ids) {
        await apiCall('delete_file', { id })
      }

      await loadFiles()
      selectedFiles.value = []
    } catch (e) {
      error.value = e instanceof Error ? e.message : '删除文件失败'
      console.error('Failed to delete files:', e)
    } finally {
      loading.value = false
    }
  }

  // 重命名文件
  async function renameFile(id: string, newName: string) {
    loading.value = true
    error.value = null

    try {
      await apiCall('rename_file', { id, newName })
      await loadFiles()
    } catch (e) {
      error.value = e instanceof Error ? e.message : '重命名文件失败'
      console.error('Failed to rename file:', e)
    } finally {
      loading.value = false
    }
  }

  // 打开文件
  async function openFile(file: FileItem) {
    try {
      await apiCall('open_file', { id: file.id })
    } catch (e) {
      error.value = e instanceof Error ? e.message : '打开文件失败'
      console.error('Failed to open file:', e)
    }
  }

  // 在文件夹中打开
  async function openInFolder(file: FileItem) {
    try {
      await apiCall('open_in_folder', { path: file.path })
    } catch (e) {
      error.value = e instanceof Error ? e.message : '在文件夹中打开失败'
      console.error('Failed to open in folder:', e)
    }
  }

  // 导航到指定路径
  function navigateToPath(path: string) {
    loadFiles(path)
  }

  // 后退
  function goBack() {
    if (canGoBack.value) {
      historyIndex.value--
      loadFiles(navigationHistory.value[historyIndex.value])
    }
  }

  // 前进
  function goForward() {
    if (canGoForward.value) {
      historyIndex.value++
      loadFiles(navigationHistory.value[historyIndex.value])
    }
  }

  // 返回上一级目录
  function goUp() {
    const parts = currentPath.value.split(/[\\/]/).filter(Boolean)
    if (parts.length > 0) {
      const parentPath = parts.slice(0, -1).join('/') || '/'
      navigateToPath(parentPath)
    }
  }

  // 设置视图模式
  function setViewMode(mode: 'list' | 'grid') {
    viewMode.value = mode
  }

  // 设置过滤类型
  function setFilterType(type: string | null) {
    filterType.value = type
  }

  // 设置选中的文件
  function setSelectedFiles(files: FileItem[]) {
    selectedFiles.value = files
  }

  // 创建文件夹
  async function createFolder(parentPath: string, folderName: string) {
    loading.value = true
    error.value = null

    try {
      await apiCall('create_folder', {
        parentPath,
        folderName
      })

      await loadFolderTree()
      if (currentPath.value === parentPath) {
        await loadFiles()
      }
    } catch (e) {
      error.value = e instanceof Error ? e.message : '创建文件夹失败'
      console.error('Failed to create folder:', e)
    } finally {
      loading.value = false
    }
  }

  // 删除文件夹
  async function deleteFolder(folderPath: string, force: boolean = false) {
    loading.value = true
    error.value = null

    try {
      await apiCall('delete_folder', {
        folderPath,
        force
      })

      await loadFolderTree()
      if (currentPath.value.startsWith(folderPath)) {
        await loadFiles('/')
      } else if (currentPath.value === folderPath) {
        await loadFiles('/')
      } else {
        await loadFiles()
      }
    } catch (e) {
      error.value = e instanceof Error ? e.message : '删除文件夹失败'
      console.error('Failed to delete folder:', e)
    } finally {
      loading.value = false
    }
  }

  // 重命名文件夹
  async function renameFolder(folderPath: string, newName: string) {
    loading.value = true
    error.value = null

    try {
      await apiCall('rename_folder', {
        folderPath,
        newName
      })

      await loadFolderTree()
      if (currentPath.value.startsWith(folderPath)) {
        const parentPath = folderPath.substring(0, folderPath.lastIndexOf('/'))
        const newPath = parentPath ? `${parentPath}/${newName}` : `/${newName}`
        await loadFiles(newPath)
      } else {
        await loadFiles()
      }
    } catch (e) {
      error.value = e instanceof Error ? e.message : '重命名文件夹失败'
      console.error('Failed to rename folder:', e)
    } finally {
      loading.value = false
    }
  }

  // 移动文件夹
  async function moveFolder(folderPath: string, targetParentPath: string) {
    loading.value = true
    error.value = null

    try {
      await apiCall('move_folder', {
        folderPath,
        targetParentPath
      })

      await loadFolderTree()
      if (currentPath.value.startsWith(folderPath)) {
        const folderName = folderPath.substring(folderPath.lastIndexOf('/') + 1)
        const newPath = targetParentPath ? `${targetParentPath}/${folderName}` : `/${folderName}`
        await loadFiles(newPath)
      } else {
        await loadFiles()
      }
    } catch (e) {
      error.value = e instanceof Error ? e.message : '移动文件夹失败'
      console.error('Failed to move folder:', e)
    } finally {
      loading.value = false
    }
  }

  // 清除错误信息
  function clearError() {
    error.value = null
  }

  // 更新资源评分
  async function updateRating(assetId: string, rating: number | null) {
    try {
      await apiCall('update_asset', {
        id: assetId,
        ...(rating !== null ? { user_rating: rating } : { clear_user_rating: true })
      })
      
      // 更新本地状态
      const file = files.value.find(f => f.id === assetId)
      if (file) {
        file.rating = rating ?? undefined
      }
    } catch (e) {
      error.value = e instanceof Error ? e.message : '更新评分失败'
      console.error('Failed to update rating:', e)
      throw e
    }
  }

  return {
    currentPath,
    files,
    selectedFiles,
    folderTree,
    viewMode,
    searchQuery,
    filterType,
    loading,
    error,
    filteredFiles,
    canGoBack,
    canGoForward,
    sortField,
    sortOrder,
    projectId,
    setProjectId,
    loadFiles,
    loadFolderTree,
    fetchStats,
    searchFiles,
    searchFilesAdvanced,
    deleteFiles,
    renameFile,
    openFile,
    openInFolder,
    navigateToPath,
    goBack,
    goForward,
    goUp,
    setViewMode,
    setFilterType,
    setSelectedFiles,
    clearError,
    createFolder,
    deleteFolder,
    renameFolder,
    moveFolder,
    setSort,
    toggleSort,
    updateRating
  }
})
