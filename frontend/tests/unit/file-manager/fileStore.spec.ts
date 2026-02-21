import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useFileStore } from '@/stores/fileStore'
import { apiCall } from '@/services/api'

vi.mock('@/services/api', () => ({
  apiCall: vi.fn()
}))

describe('fileStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  describe('initial state', () => {
    it('should have default values', () => {
      const store = useFileStore()
      
      expect(store.currentPath).toBe('/')
      expect(store.files).toEqual([])
      expect(store.selectedFiles).toEqual([])
      expect(store.folderTree).toEqual([])
      expect(store.viewMode).toBe('list')
      expect(store.searchQuery).toBe('')
      expect(store.filterType).toBe(null)
      expect(store.loading).toBe(false)
      expect(store.error).toBe(null)
    })
  })

  describe('loadFiles', () => {
    it('should load files successfully', async () => {
      const mockFiles = [
        { id: '1', name: 'test.jpg', path: '/test.jpg', size: 1024, file_type: 'image', modified_at: '2024-01-01' },
        { id: '2', name: 'test.pdf', path: '/test.pdf', size: 2048, file_type: 'document', modified_at: '2024-01-02' }
      ]
      
      vi.mocked(apiCall).mockResolvedValue(mockFiles)
      
      const store = useFileStore()
      await store.loadFiles('/test')
      
      expect(store.files).toEqual(mockFiles)
      expect(store.currentPath).toBe('/test')
      expect(store.loading).toBe(false)
      expect(store.error).toBe(null)
    })

    it('should handle loading state', async () => {
      vi.mocked(apiCall).mockImplementation(() => new Promise(resolve => setTimeout(() => resolve([]), 100)))
      
      const store = useFileStore()
      const loadPromise = store.loadFiles('/test')
      
      expect(store.loading).toBe(true)
      await loadPromise
      expect(store.loading).toBe(false)
    })

    it('should handle errors', async () => {
      vi.mocked(apiCall).mockRejectedValue(new Error('Failed to load'))
      
      const store = useFileStore()
      await store.loadFiles('/test')
      
      expect(store.error).toBe('Failed to load')
      expect(store.loading).toBe(false)
    })
  })

  describe('searchFiles', () => {
    it('should search files with query', async () => {
      const mockResults = [
        { id: '1', name: 'test.jpg', path: '/test.jpg', size: 1024, file_type: 'image', modified_at: '2024-01-01' }
      ]
      
      vi.mocked(apiCall).mockResolvedValue(mockResults)
      
      const store = useFileStore()
      await store.searchFiles('test')
      
      expect(store.searchQuery).toBe('test')
      expect(store.files).toEqual(mockResults)
    })

    it('should clear search when query is empty', async () => {
      vi.mocked(apiCall).mockResolvedValue([])
      
      const store = useFileStore()
      store.searchQuery = 'test'
      await store.searchFiles('')
      
      expect(store.searchQuery).toBe('')
    })
  })

  describe('filteredFiles', () => {
    it('should filter by file type', () => {
      const store = useFileStore()
      store.files = [
        { id: '1', name: 'test.jpg', path: '/test.jpg', size: 1024, file_type: 'image', modified_at: '2024-01-01' },
        { id: '2', name: 'test.pdf', path: '/test.pdf', size: 2048, file_type: 'document', modified_at: '2024-01-02' }
      ]
      
      store.setFilterType('image')
      
      expect(store.filteredFiles).toHaveLength(1)
      expect(store.filteredFiles[0].file_type).toBe('image')
    })

    it('should filter by search query', () => {
      const store = useFileStore()
      store.files = [
        { id: '1', name: 'test.jpg', path: '/test.jpg', size: 1024, file_type: 'image', modified_at: '2024-01-01' },
        { id: '2', name: 'demo.pdf', path: '/demo.pdf', size: 2048, file_type: 'document', modified_at: '2024-01-02' }
      ]
      
      store.searchQuery = 'test'
      
      expect(store.filteredFiles).toHaveLength(1)
      expect(store.filteredFiles[0].name).toBe('test.jpg')
    })

    it('should combine filters', () => {
      const store = useFileStore()
      store.files = [
        { id: '1', name: 'test.jpg', path: '/test.jpg', size: 1024, file_type: 'image', modified_at: '2024-01-01' },
        { id: '2', name: 'test.pdf', path: '/test.pdf', size: 2048, file_type: 'document', modified_at: '2024-01-02' }
      ]
      
      store.searchQuery = 'test'
      store.setFilterType('image')
      
      expect(store.filteredFiles).toHaveLength(1)
      expect(store.filteredFiles[0].file_type).toBe('image')
    })
  })

  describe('navigation', () => {
    it('should navigate to path', async () => {
      vi.mocked(apiCall).mockResolvedValue([])
      
      const store = useFileStore()
      await store.navigateToPath('/test')
      
      expect(store.currentPath).toBe('/test')
    })

    it('should track navigation history', async () => {
      vi.mocked(invoke).mockResolvedValue([])
      
      const store = useFileStore()
      await store.navigateToPath('/test1')
      await store.navigateToPath('/test2')
      
      expect(store.canGoBack).toBe(true)
      expect(store.canGoForward).toBe(false)
    })

    it('should go back in history', async () => {
      vi.mocked(invoke).mockResolvedValue([])
      
      const store = useFileStore()
      await store.navigateToPath('/test1')
      await store.navigateToPath('/test2')
      await store.goBack()
      
      expect(store.currentPath).toBe('/test1')
    })

    it('should go forward in history', async () => {
      vi.mocked(invoke).mockResolvedValue([])
      
      const store = useFileStore()
      await store.navigateToPath('/test1')
      await store.navigateToPath('/test2')
      await store.goBack()
      await store.goForward()
      
      expect(store.currentPath).toBe('/test2')
    })
  })

  describe('view mode', () => {
    it('should set view mode', () => {
      const store = useFileStore()
      store.setViewMode('grid')
      
      expect(store.viewMode).toBe('grid')
    })
  })

  describe('deleteFiles', () => {
    it('should delete files successfully', async () => {
      vi.mocked(invoke).mockResolvedValue(undefined)
      vi.mocked(invoke).mockResolvedValue([])
      
      const store = useFileStore()
      await store.deleteFiles(['1', '2'])
      
      expect(invoke).toHaveBeenCalledWith('delete_file', { id: '1' })
      expect(invoke).toHaveBeenCalledWith('delete_file', { id: '2' })
      expect(store.selectedFiles).toEqual([])
    })
  })

  describe('renameFile', () => {
    it('should rename file successfully', async () => {
      vi.mocked(invoke).mockResolvedValue(undefined)
      vi.mocked(invoke).mockResolvedValue([])
      
      const store = useFileStore()
      await store.renameFile('1', 'new-name.jpg')
      
      expect(invoke).toHaveBeenCalledWith('rename_file', { id: '1', newName: 'new-name.jpg' })
    })
  })
})