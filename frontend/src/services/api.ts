import axios from 'axios'

// API 服务器基础 URL
const getApiBaseUrl = async () => {
  try {
    const cachedPort = localStorage.getItem('api_server_port')
    if (cachedPort) {
      return `http://localhost:${cachedPort}`
    }
    return 'http://localhost:32000'
  } catch {
    return 'http://localhost:32000'
  }
}

const getBridgeBaseUrl = (): string | null => {
  const w = window as any
  if (w?.mediaAssistant?.coreBaseUrl && typeof w.mediaAssistant.coreBaseUrl === 'string') {
    return w.mediaAssistant.coreBaseUrl
  }
  return null
}

// 创建 axios 实例
const apiClient = axios.create({
  baseURL: 'http://localhost:32000',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 端口发现 - 尝试找到正在运行的 API 服务器端口
export const discoverApiPort = async (): Promise<string | null> => {
  const ports = [32000, 32001, 32002, 32003, 32004]
  
  // 并行检查所有端口，使用更短的超时
  const checks = ports.map(port => 
    axios.get(`http://localhost:${port}/api/v1/health`, { timeout: 300 })
      .then(res => res.status === 200 ? port : null)
      .catch(() => null)
  )
  
  const results = await Promise.all(checks)
  const availablePort = results.find(port => port !== null)
  
  if (availablePort) {
    localStorage.setItem('api_server_port', availablePort.toString())
    return `http://localhost:${availablePort}`
  }
  
  return null
}

let apiInitialized = false
let apiDiscoveryFailed = false

const initializeApiClient = async () => {
  if (apiInitialized) return
  if (apiDiscoveryFailed) return

  // 优先使用 Electron bridge 提供的 URL
  const bridgeBaseUrl = getBridgeBaseUrl()
  if (bridgeBaseUrl) {
    apiClient.defaults.baseURL = bridgeBaseUrl
    apiInitialized = true
    console.log('[API] Using bridge URL:', bridgeBaseUrl)
    return
  }

  // 尝试使用缓存的端口
  const baseUrl = await getApiBaseUrl()
  if (baseUrl) {
    apiClient.defaults.baseURL = baseUrl
  }
  
  // 尝试发现可用端口
  const discoveredUrl = await discoverApiPort()
  if (discoveredUrl) {
    apiClient.defaults.baseURL = discoveredUrl
    apiInitialized = true
    console.log('[API] Connected to:', discoveredUrl)
  } else {
    // 即使发现失败，也标记为已初始化，使用默认 URL
    apiInitialized = true
    console.warn('[API] Failed to discover API server, using default URL:', apiClient.defaults.baseURL)
  }
}

// 统一的 API 调用函数
export const apiCall = async <T = any>(
  command: string,
  args?: Record<string, any>
): Promise<T> => {
  await initializeApiClient()

  const normalizedArgs = normalizeArgs(command, args)
  const endpoint = mapCommandToEndpoint(command, normalizedArgs)
  const method = getHttpMethod(command)

  try {
    let response
    if (method === 'GET') {
      response = await apiClient.get(endpoint, {
        params: normalizedArgs
      })
    } else if (method === 'POST') {
      response = await apiClient.post(endpoint, normalizedArgs)
    } else if (method === 'DELETE') {
      response = await apiClient.delete(endpoint, {
        params: normalizedArgs,
        data: normalizedArgs
      })
    } else {
      throw new Error(`Unsupported HTTP method: ${method}`)
    }

    if (response.data && typeof response.data === 'object' && 'data' in response.data) {
      return response.data.data as T
    }
    return response.data as T
  } catch (error) {
    console.error(`[API] ${command} failed:`, error)
    throw error
  }
}

const normalizeArgs = (command: string, args?: Record<string, any>): Record<string, any> | undefined => {
  if (!args || typeof args !== 'object') {
    return args
  }

  // 兼容旧的 get_project_files 接口
  if (command === 'get_project_files' && args.projectId) {
    return { 
      projectId: args.projectId, 
      limit: args.limit || 50 
    }
  }

  if (command === 'update_project' && args.project && typeof args.project === 'object') {
    return args.project
  }

  if (command === 'create_project_from_folder') {
    const folderPath = typeof args.folderPath === 'string' ? args.folderPath : ''
    const fileName = folderPath.split(/[\\/]/).filter(Boolean).pop() || 'New Project'
    return {
      name: args.name || fileName,
      project_type: args.project_type || 'custom',
      path: folderPath
    }
  }

  if (command === 'open_file' || command === 'open_in_folder') {
    if (!args.id && args.asset_id) {
      return { ...args, id: args.asset_id }
    }
  }

  if ((command === 'scan_project' || command === 'sync_project_files') && !args.id && args.projectId) {
    return { ...args, id: args.projectId }
  }

  return args
}

// 项目目录树 API
export interface DirectoryNode {
  path: string
  name: string
  source_type?: string
  is_root: boolean
  exists: boolean
  has_children: boolean
}

export const getProjectBoundDirectories = async (projectId: string): Promise<DirectoryNode[]> => {
  const baseUrl = await getApiBaseUrl()
  const response = await axios.get(`${baseUrl}/api/v1/projects/directories/bound`, {
    params: { projectId }
  })
  return response.data.data || []
}

export const getDirectoryChildren = async (path: string): Promise<DirectoryNode[]> => {
  const baseUrl = await getApiBaseUrl()
  const response = await axios.get(`${baseUrl}/api/v1/library/directories/children`, {
    params: { path }
  })
  return response.data.data || []
}

const mapCommandToEndpoint = (command: string, args?: Record<string, any>): string => {
  if (command.startsWith('/')) {
    return command
  }

  const endpointMap: Record<string, any> = {
    'list_projects': '/api/v1/projects',
    'get_project': '/api/v1/projects/get',
    'create_project': '/api/v1/projects',
    'update_project': '/api/v1/projects/update',
    'delete_project': '/api/v1/projects/delete',
    'create_project_from_folder': '/api/v1/projects',
    'update_project_status': '/api/v1/projects/update',
    'update_project_deadline': '/api/v1/projects/update',
    'update_project_path': '/api/v1/projects/update-path',
    'scan_project': '/api/v1/projects/scan',
    'sync_project_files': '/api/v1/projects/scan',
    'get_project_stats': '/api/v1/projects/stats',
    'add_file_to_project': '/api/v1/projects/assets/add',
    'remove_file_from_project': '/api/v1/projects/assets/remove',

    'list_folders': '/api/v1/get_common_directories',
    'get_folder_tree': '/api/v1/folders/tree',
    'create_folder': '/api/v1/folders',
    'delete_folder': (args: any) => `/api/v1/folders/${args?.folderPath}`,
    'rename_folder': '/api/v1/folders/rename',
    'move_folder': '/api/v1/folders/move',

    'list_files': '/api/v1/files',
    'get_file_detail': '/api/v1/assets/get',
    'delete_file': '/api/v1/assets/delete',
    'rename_file': (args: any) => `/api/v1/files/${args?.id}/rename`,
    'move_file': '/api/v1/files/move',
    'copy_file': '/api/v1/files/copy',
    'archive_files': '/api/v1/assets/archive',

    'search_files': '/api/v1/search/files',
    'search_collections': '/api/v1/search/collections',
    'search_tags': '/api/v1/tags/search',

    'create_tag': '/api/v1/tags/create',
    'update_tag': '/api/v1/tags/update',
    'delete_tag': (args: any) => `/api/v1/tags/delete?id=${encodeURIComponent(args?.tagId || '')}`,
    'list_tags': '/api/v1/tags',
    'list_tag_tree': '/api/v1/tags/tree',
    'get_file_tags': (args: any) => `/api/v1/tags/file?file_id=${encodeURIComponent(args?.fileId || '')}`,
    'add_tags_to_files': '/api/v1/tags/batch-add',
    'remove_tags_from_files': '/api/v1/tags/batch-remove',

    'get_common_directories': '/api/v1/get_common_directories',
    'run_onboarding': '/api/v1/run_onboarding',
    'is_first_launch': '/api/v1/system/first-launch',
    'complete_first_launch': '/api/v1/system/first-launch/complete',

    'list_library_sources': '/api/v1/library/sources',
    'add_library_source': '/api/v1/library/sources/add',
    'remove_library_source': '/api/v1/library/sources/remove',

    'get_thumbnail': (args: any) => `/api/v1/thumbnails/${args?.id}`,
    'generate_thumbnail': '/api/v1/thumbnails/generate',

    'list_assets': '/api/v1/assets',
    'get_project_files': '/api/v1/assets',
    'list_activity_logs': '/api/v1/activity/list',
    'get_directory_warnings': '/api/v1/projects/directories/warnings',
    'get_dock_items': '/api/v1/projects',
    'get_directory_contents': '/api/v1/files',
    'update_asset': '/api/v1/assets/update',
    'get_import_progress': '/api/v1/get_import_progress',
    'get_active_tasks': '/api/v1/tasks/active',
    'get_system_info': '/api/v1/system/info',
    'get_default_dirs': '/api/v1/system/default-dirs',
  }

  const mapped = endpointMap[command]
  if (typeof mapped === 'function') {
    return mapped(args)
  }
  return mapped || `/api/v1/${command}`
}

const getHttpMethod = (command: string): 'GET' | 'POST' | 'PUT' | 'DELETE' => {
  if (command.startsWith('/')) {
    return 'POST'
  }

  const getMethods = [
    'list_projects', 'get_project', 'list_folders', 'get_folder_tree',
    'list_files', 'get_file_detail', 'search_files',
    'search_collections', 'search_tags', 'list_tags', 'list_tag_tree', 'get_file_tags',
    'get_thumbnail', 'get_common_directories', 'get_dock_items',
    'get_directory_contents', 'is_first_launch', 'get_active_tasks',
    'get_system_info', 'get_default_dirs', 'get_project_stats', 'get_import_progress',
    'get_project_files', 'list_assets', 'list_activity_logs', 'get_directory_warnings',
    'list_library_sources'
  ]

  const deleteMethods = ['delete_folder']

  if (getMethods.includes(command)) return 'GET'
  if (deleteMethods.includes(command)) return 'DELETE'
  return 'POST'
}

// WebSocket 事件处理
export type EventCallback = (data: any) => void
const eventListeners: Record<string, EventCallback[]> = {}

export const subscribeToEvent = (type: string, callback: EventCallback) => {
  if (!eventListeners[type]) {
    eventListeners[type] = []
  }
  eventListeners[type].push(callback)
}

let ws: WebSocket | null = null
let wsReconnectAttempts = 0
const WS_MAX_RECONNECT_ATTEMPTS = 5

export const connectWebSocket = async () => {
  if (ws) return
  if (wsReconnectAttempts >= WS_MAX_RECONNECT_ATTEMPTS) {
    console.warn('[WS] Max reconnect attempts reached, stopping')
    return
  }

  try {
    await initializeApiClient()
    const baseUrl = getBridgeBaseUrl() || apiClient.defaults.baseURL || await getApiBaseUrl()
    const wsUrl = baseUrl.replace('http', 'ws') + '/events'

    ws = new WebSocket(wsUrl)
    console.log('[WS] Connecting to:', wsUrl)

    ws.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data)
        console.log('[WS] Received message:', msg)
        if (msg.type && eventListeners[msg.type]) {
          eventListeners[msg.type].forEach(cb => cb(msg))
        }
      } catch (err) {
        console.error('[WS] Failed to parse message:', err)
      }
    }

    ws.onclose = () => {
      console.warn('[WS] Closed:', wsUrl)
      ws = null
      wsReconnectAttempts++
      if (wsReconnectAttempts < WS_MAX_RECONNECT_ATTEMPTS) {
        setTimeout(connectWebSocket, 3000)
      } else {
        console.warn('[WS] Stopping reconnection attempts')
      }
    }

    ws.onerror = (err) => {
      console.error('[WS] Error:', err)
      ws?.close()
    }
  } catch (err) {
    console.error('[WS] Failed to connect:', err)
  }
}

export { initializeApiClient }

// ============================================
// 智能搜索 API
// ============================================

export interface ListAssetsParams {
  projectId?: string
  directory?: string
  search?: string
  tagIds?: string[]
  types?: string[]
  shapes?: string[]
  
  quickFilter?: 'recent' | 'unrated' | 'large' | 'vertical' | 'horizontal'
  datePreset?: 'today' | 'thisWeek' | 'thisMonth' | 'lastMonth'
  
  sizeMin?: number
  sizeMax?: number
  ratingMin?: number
  ratingMax?: number
  mtimeFrom?: number
  mtimeTo?: number
  
  sortBy?: string
  sortOrder?: 'asc' | 'desc'
  limit?: number
  cursor?: string
}

export interface AssetListItem {
  id: string
  name: string
  path: string
  size: number
  mtime: number
  modified_at: number
  file_type?: string
  status?: string
  shape: string
  suggested_rating?: number
  user_rating?: number
  thumbnail_path?: string
  width?: number
  height?: number
  created_at: number
}

export interface ListAssetsResult {
  items: AssetListItem[]
  nextCursor?: string
  hasMore: boolean
  total: number
}

export interface SearchHistoryItem {
  id: string
  query: string
  filters: Record<string, any>
  count: number
  created_at: number
  updated_at: number
}

export const listAssets = async (params: ListAssetsParams): Promise<ListAssetsResult> => {
  const queryParams = new URLSearchParams()
  
  if (params.projectId) queryParams.set('projectId', params.projectId)
  if (params.directory) queryParams.set('directory', params.directory)
  if (params.search) queryParams.set('search', params.search)
  if (params.quickFilter) queryParams.set('quickFilter', params.quickFilter)
  if (params.datePreset) queryParams.set('datePreset', params.datePreset)
  if (params.tagIds?.length) queryParams.set('tagIds', params.tagIds.join(','))
  if (params.types?.length) queryParams.set('types', params.types.join(','))
  if (params.shapes?.length) queryParams.set('shapes', params.shapes.join(','))
  if (params.sortBy) queryParams.set('sortBy', params.sortBy)
  if (params.sortOrder) queryParams.set('sortOrder', params.sortOrder)
  if (params.limit) queryParams.set('limit', params.limit.toString())
  if (params.cursor) queryParams.set('cursor', params.cursor)
  if (params.sizeMin) queryParams.set('sizeMin', params.sizeMin.toString())
  if (params.sizeMax) queryParams.set('sizeMax', params.sizeMax.toString())
  if (params.ratingMin !== undefined) queryParams.set('ratingMin', params.ratingMin.toString())
  if (params.ratingMax !== undefined) queryParams.set('ratingMax', params.ratingMax.toString())
  if (params.mtimeFrom) queryParams.set('mtimeFrom', params.mtimeFrom.toString())
  if (params.mtimeTo) queryParams.set('mtimeTo', params.mtimeTo.toString())
  
  return apiCall<ListAssetsResult>(`/api/v1/assets?${queryParams.toString()}`)
}

export const getSearchHistory = async (limit = 10): Promise<SearchHistoryItem[]> => {
  const result = await apiCall<SearchHistoryItem[]>(`/api/v1/search/history?limit=${limit}`)
  return Array.isArray(result) ? result : []
}

export const clearSearchHistory = async (): Promise<void> => {
  await apiCall('/api/v1/search/history/clear', {}, 'DELETE')
}
