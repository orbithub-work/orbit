import { defineStore } from 'pinia'
import { ref } from 'vue'
import { apiCall } from '@/services/api'

export interface FolderNode {
  id: string
  name: string
  path: string
  children?: FolderNode[]
  isExpanded?: boolean
}

export interface CommonDirectory {
  name: string
  path: string
  icon: string
}

export const useFolderStore = defineStore('folder', () => {
  const folderTree = ref<FolderNode[]>([])
  const commonDirectories = ref<CommonDirectory[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)
  const selectedFolder = ref<FolderNode | null>(null)
  const selectedCommonDir = ref<CommonDirectory | null>(null)
  const isFirstLaunch = ref(false)
  const folderApiAvailable = ref(false)

  const loadFolderTree = async () => {
    loading.value = true
    error.value = null

    try {
      const result = await apiCall<CommonDirectory[]>('get_common_directories')
      commonDirectories.value = result || []
      folderTree.value = commonDirectories.value.map(dir => ({
        id: dir.path,
        name: dir.name,
        path: dir.path
      }))
    } catch (e: any) {
      const msg = e?.message || 'Unknown error'
      error.value = msg
      console.error('Failed to load folder tree (sanitized):', msg)
    } finally {
      loading.value = false
    }
  }

  const selectFolder = (folder: FolderNode | null) => {
    selectedFolder.value = folder
  }

  const toggleFolder = (folder: FolderNode) => {
    folder.isExpanded = !folder.isExpanded
  }

  const refreshFolderTree = async () => {
    await loadFolderTree()
  }

  const fetchFolders = async () => {
    await loadFolderTree()
  }

  const loadCommonDirectories = async () => {
    try {
      const result = await apiCall<CommonDirectory[]>('get_common_directories')
      commonDirectories.value = result || []
    } catch (e: any) {
      console.error('Failed to load common directories (sanitized):', e?.message)
    }
  }

  const selectCommonDir = (dir: CommonDirectory | null) => {
    selectedCommonDir.value = dir
    if (dir) {
      selectedFolder.value = {
        id: dir.path,
        name: dir.name,
        path: dir.path
      }
    }
  }

  const checkFirstLaunch = async () => {
    try {
      const result = await apiCall<boolean>('is_first_launch')
      isFirstLaunch.value = result || false
      return isFirstLaunch.value
    } catch (e) {
      console.error('Failed to check first launch:', e)
      return false
    }
  }

  const completeFirstLaunch = async () => {
    try {
      await apiCall('complete_first_launch')
      isFirstLaunch.value = false
    } catch (e) {
      console.error('Failed to complete first launch:', e)
    }
  }

  return {
    folderTree,
    commonDirectories,
    loading,
    error,
    selectedFolder,
    selectedCommonDir,
    isFirstLaunch,
    folderApiAvailable,
    loadFolderTree,
    fetchFolders,
    selectFolder,
    toggleFolder,
    refreshFolderTree,
    loadCommonDirectories,
    selectCommonDir,
    checkFirstLaunch,
    completeFirstLaunch
  }
})
