/**
 * Smart Archive OS Plugin SDK
 * 
 * Provides API access for plugins to interact with the host application
 */

export interface PluginAPI {
  assets: AssetAPI
  tags: TagAPI
  ui: UIAPI
  http: HttpAPI
  context: ContextAPI
}

export interface AssetAPI {
  list(query?: AssetQuery): Promise<Asset[]>
  get(id: number): Promise<Asset>
  update(id: number, data: Partial<Asset>): Promise<void>
  delete(id: number): Promise<void>
  getSelected(): Asset[]
}

export interface TagAPI {
  list(): Promise<Tag[]>
  create(tag: CreateTagInput): Promise<Tag>
  update(id: number, data: Partial<Tag>): Promise<void>
  delete(id: number): Promise<void>
  batchAdd(assetIds: number[], tagIds: number[]): Promise<void>
  batchRemove(assetIds: number[], tagIds: number[]): Promise<void>
}

export interface UIAPI {
  showNotification(message: string, type?: 'success' | 'error' | 'warning' | 'info', duration?: number): Promise<void>
  showDialog(options: DialogOptions): Promise<DialogResult>
  showConfirm(message: string, title?: string): Promise<boolean>
}

export interface HttpAPI {
  get<T = any>(url: string, options?: RequestOptions): Promise<T>
  post<T = any>(url: string, data?: any, options?: RequestOptions): Promise<T>
  put<T = any>(url: string, data?: any, options?: RequestOptions): Promise<T>
  delete<T = any>(url: string, options?: RequestOptions): Promise<T>
}

export interface ContextAPI {
  getSelectedAssets(): Asset[]
  getSelectedTags(): Tag[]
  getCurrentView(): string
  getCurrentPath(): string
  setContext(context: UIContext): Promise<void>
}

// Types
export interface Asset {
  id: number
  name: string
  path: string
  size: number
  mime_type: string
  created_at: string
  updated_at: string
  thumbnail_path?: string
  rating?: number
  tags?: Tag[]
}

export interface Tag {
  id: number
  name: string
  color?: string
  parent_id?: number
}

export interface AssetQuery {
  search?: string
  tags?: number[]
  mime_type?: string
  rating?: number
  limit?: number
  offset?: number
}

export interface CreateTagInput {
  name: string
  color?: string
  parent_id?: number
}

export interface DialogOptions {
  title: string
  message: string
  type?: 'alert' | 'confirm' | 'prompt'
  buttons?: string[]
  data?: Record<string, any>
}

export interface DialogResult {
  confirmed: boolean
  value?: string
  data?: Record<string, any>
}

export interface RequestOptions {
  headers?: Record<string, string>
  timeout?: number
}

export interface UIContext {
  asset_ids?: number[]
  tag_ids?: number[]
  path?: string
  view?: string
  action?: string
  data?: Record<string, any>
}

/**
 * Create plugin API instance
 */
export function createPluginAPI(baseURL: string = 'http://localhost:32000'): PluginAPI {
  const api: PluginAPI = {
    assets: {
      async list(query) {
        const params = new URLSearchParams()
        if (query?.search) params.set('search', query.search)
        if (query?.tags) params.set('tags', query.tags.join(','))
        if (query?.mime_type) params.set('mime_type', query.mime_type)
        if (query?.rating) params.set('rating', query.rating.toString())
        if (query?.limit) params.set('limit', query.limit.toString())
        if (query?.offset) params.set('offset', query.offset.toString())
        
        const response = await fetch(`${baseURL}/api/assets?${params}`)
        const result = await response.json()
        return result.data || []
      },

      async get(id) {
        const response = await fetch(`${baseURL}/api/assets/get?id=${id}`)
        const result = await response.json()
        return result.data
      },

      async update(id, data) {
        await fetch(`${baseURL}/api/assets/update`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ id, ...data })
        })
      },

      async delete(id) {
        await fetch(`${baseURL}/api/assets/delete`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ id })
        })
      },

      getSelected() {
        return (window as any).__PLUGIN_CONTEXT__?.selectedAssets || []
      }
    },

    tags: {
      async list() {
        const response = await fetch(`${baseURL}/api/tags`)
        const result = await response.json()
        return result.data || []
      },

      async create(tag) {
        const response = await fetch(`${baseURL}/api/tags/create`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(tag)
        })
        const result = await response.json()
        return result.data
      },

      async update(id, data) {
        await fetch(`${baseURL}/api/tags/update`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ id, ...data })
        })
      },

      async delete(id) {
        await fetch(`${baseURL}/api/tags/delete`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ id })
        })
      },

      async batchAdd(assetIds, tagIds) {
        await fetch(`${baseURL}/api/tags/batch-add`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ asset_ids: assetIds, tag_ids: tagIds })
        })
      },

      async batchRemove(assetIds, tagIds) {
        await fetch(`${baseURL}/api/tags/batch-remove`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ asset_ids: assetIds, tag_ids: tagIds })
        })
      }
    },

    ui: {
      async showNotification(message, type = 'info', duration = 0) {
        await fetch(`${baseURL}/api/ui/notification`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ message, type, duration })
        })
      },

      async showDialog(options) {
        const response = await fetch(`${baseURL}/api/ui/dialog`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(options)
        })
        const result = await response.json()
        return result.data
      },

      async showConfirm(message, title = '确认') {
        const result = await this.showDialog({
          title,
          message,
          type: 'confirm',
          buttons: ['取消', '确认']
        })
        return result.confirmed
      }
    },

    http: {
      async get(url, options) {
        const response = await fetch(url, {
          method: 'GET',
          headers: options?.headers,
        })
        return response.json()
      },

      async post(url, data, options) {
        const response = await fetch(url, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            ...options?.headers
          },
          body: JSON.stringify(data)
        })
        return response.json()
      },

      async put(url, data, options) {
        const response = await fetch(url, {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json',
            ...options?.headers
          },
          body: JSON.stringify(data)
        })
        return response.json()
      },

      async delete(url, options) {
        const response = await fetch(url, {
          method: 'DELETE',
          headers: options?.headers,
        })
        return response.json()
      }
    },

    context: {
      getSelectedAssets() {
        return (window as any).__PLUGIN_CONTEXT__?.selectedAssets || []
      },

      getSelectedTags() {
        return (window as any).__PLUGIN_CONTEXT__?.selectedTags || []
      },

      getCurrentView() {
        return (window as any).__PLUGIN_CONTEXT__?.currentView || 'pool'
      },

      getCurrentPath() {
        return (window as any).__PLUGIN_CONTEXT__?.currentPath || ''
      },

      async setContext(context) {
        await fetch(`${baseURL}/api/ui/context`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(context)
        })
      }
    }
  }

  return api
}

// Export for use in plugins
export default createPluginAPI
