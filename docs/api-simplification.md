# API 简化方案

## 现状分析

### 后端接口（已实现）

✅ **GET /api/assets** - 统一的资源列表接口
- 支持所有筛选参数：projectId, search, tagIds, types, shapes, rating, sizeMin/Max, widthMin/Max, heightMin/Max, mtimeFrom/To
- 支持排序：sortBy, sortOrder
- 支持分页：cursor, limit
- 返回统一结构：`{ items: [], nextCursor: string, hasMore: bool, total: int }`

✅ **POST /api/assets/update** - 资源更新接口
- 支持更新评分：`{ id: string, user_rating: int }` 或 `{ id: string, clear_user_rating: true }`
- 支持更新元数据：`{ id: string, media_meta: string }`

### 前端现状

❌ **问题：** 前端使用了多个不同的接口获取资源列表
- `get_project_files` - 获取项目文件
- `search_files` - 搜索文件
- `advanced_search` - 高级搜索
- `list_files` - 列出文件

这些接口实际上都可以用 **GET /api/assets** 统一替代。

## 优化方案

### 1. 统一使用 `list_assets` 接口

**前端调用示例：**

```typescript
// 基础列表（替代 get_project_files）
const assets = await apiCall('list_assets', {
  projectId: 'proj-123',
  limit: 50
})

// 搜索（替代 search_files）
const assets = await apiCall('list_assets', {
  search: 'vacation',
  types: ['image', 'video'],
  limit: 100
})

// 高级筛选（替代 advanced_search）
const assets = await apiCall('list_assets', {
  projectId: 'proj-123',
  types: ['image'],
  shapes: ['landscape'],
  ratingMin: 4,
  sizeMin: 1024 * 1024, // 1MB
  mtimeFrom: 1640995200, // 2022-01-01
  sortBy: 'mtime',
  sortOrder: 'desc',
  cursor: 'next-page-token',
  limit: 50
})

// 按标签筛选
const assets = await apiCall('list_assets', {
  tagIds: ['tag-1', 'tag-2'],
  limit: 50
})
```

### 2. 更新评分

```typescript
// 设置评分
await apiCall('update_asset', {
  id: 'asset-123',
  user_rating: 5
})

// 清除评分
await apiCall('update_asset', {
  id: 'asset-123',
  clear_user_rating: true
})
```

## 实施步骤

### Step 1: 更新 API 映射

修改 `frontend/src/services/api.ts`：

```typescript
const endpointMap: Record<string, any> = {
  // ... 其他映射
  
  // 统一使用 list_assets
  'list_assets': '/api/assets',
  'get_project_files': '/api/assets',  // 重定向到统一接口
  'search_files': '/api/assets',       // 重定向到统一接口
  'advanced_search': '/api/assets',    // 重定向到统一接口
  
  // 更新资源
  'update_asset': '/api/assets/update',
}

const getMethods = [
  // ... 其他方法
  'list_assets',
  'get_project_files',
  'search_files',
  'advanced_search',
]
```

### Step 2: 重构 fileStore

修改 `frontend/src/stores/fileStore.ts`：

```typescript
// 统一的加载方法
async function loadAssets(filters: {
  projectId?: string
  search?: string
  types?: string[]
  tagIds?: string[]
  shapes?: string[]
  ratingMin?: number
  ratingMax?: number
  sizeMin?: number
  sizeMax?: number
  mtimeFrom?: number
  mtimeTo?: number
  sortBy?: string
  sortOrder?: 'asc' | 'desc'
  cursor?: string
  limit?: number
} = {}) {
  loading.value = true
  error.value = null

  try {
    const result = await apiCall<{
      items: AssetDto[]
      nextCursor: string
      hasMore: boolean
      total: number
    }>('list_assets', {
      projectId: filters.projectId || projectId.value,
      search: filters.search || searchQuery.value,
      types: filters.types,
      tagIds: filters.tagIds,
      shapes: filters.shapes,
      ratingMin: filters.ratingMin,
      ratingMax: filters.ratingMax,
      sizeMin: filters.sizeMin,
      sizeMax: filters.sizeMax,
      mtimeFrom: filters.mtimeFrom,
      mtimeTo: filters.mtimeTo,
      sortBy: filters.sortBy || sortField.value,
      sortOrder: filters.sortOrder || sortOrder.value,
      cursor: filters.cursor,
      limit: filters.limit || 50
    })

    files.value = result.items.map(transformAssetToFileItem)
    nextCursor.value = result.nextCursor
    hasMore.value = result.hasMore
    totalCount.value = result.total
  } catch (e) {
    error.value = e instanceof Error ? e.message : '加载资源失败'
    console.error('Failed to load assets:', e)
  } finally {
    loading.value = false
  }
}

// 简化的方法
async function loadFiles(path?: string) {
  await loadAssets({ projectId: projectId.value })
}

async function searchFiles(query: string) {
  searchQuery.value = query
  await loadAssets({ search: query })
}

async function searchFilesAdvanced(filters: any) {
  await loadAssets({
    search: filters.namePattern,
    types: filters.fileTypes,
    tagIds: filters.tags,
    mtimeFrom: filters.dateFrom ? new Date(filters.dateFrom).getTime() / 1000 : undefined,
    mtimeTo: filters.dateTo ? new Date(filters.dateTo).getTime() / 1000 : undefined,
    sizeMin: filters.minSize,
    sizeMax: filters.maxSize,
    ratingMin: filters.rating,
  })
}

// 更新评分
async function updateRating(assetId: string, rating: number | null) {
  try {
    await apiCall('update_asset', {
      id: assetId,
      ...(rating !== null ? { user_rating: rating } : { clear_user_rating: true })
    })
    
    // 更新本地状态
    const file = files.value.find(f => f.id === assetId)
    if (file) {
      file.rating = rating
    }
  } catch (e) {
    error.value = e instanceof Error ? e.message : '更新评分失败'
    console.error('Failed to update rating:', e)
  }
}
```

### Step 3: 更新组件

在 `PoolView.vue` 或其他组件中：

```vue
<script setup lang="ts">
import { useFileStore } from '@/stores/fileStore'

const fileStore = useFileStore()

// 加载资源
onMounted(() => {
  fileStore.loadAssets({
    projectId: props.projectId,
    limit: 50
  })
})

// 搜索
const handleSearch = (query: string) => {
  fileStore.searchFiles(query)
}

// 筛选
const handleFilter = (filters: any) => {
  fileStore.loadAssets({
    types: filters.types,
    shapes: filters.shapes,
    ratingMin: filters.ratingMin,
    // ... 其他筛选条件
  })
}

// 更新评分
const handleRatingChange = (assetId: string, rating: number) => {
  fileStore.updateRating(assetId, rating)
}

// 加载更多（分页）
const handleLoadMore = () => {
  if (fileStore.hasMore) {
    fileStore.loadAssets({
      cursor: fileStore.nextCursor
    })
  }
}
</script>
```

## 优势

1. **代码简化** - 一个接口替代多个接口，减少维护成本
2. **类型安全** - 统一的返回结构，更好的 TypeScript 支持
3. **性能优化** - 后端统一处理，可以更好地优化查询
4. **扩展性** - 新增筛选条件只需添加参数，无需新接口
5. **一致性** - 所有资源列表使用相同的数据结构和分页逻辑

## 兼容性

为了保持向后兼容，可以在 `api.ts` 中保留旧的命令名称，但都映射到新的统一接口：

```typescript
const normalizeArgs = (command: string, args?: Record<string, any>): Record<string, any> | undefined => {
  // 兼容旧接口
  if (command === 'get_project_files' && args?.projectId) {
    return { projectId: args.projectId, limit: args.limit || 50 }
  }
  
  if (command === 'search_files' && args?.query) {
    return { search: args.query, types: args.file_types, limit: args.max_results || 100 }
  }
  
  if (command === 'advanced_search') {
    return {
      search: args?.query,
      types: args?.file_types,
      tagIds: args?.tags,
      mtimeFrom: args?.date_from ? new Date(args.date_from).getTime() / 1000 : undefined,
      mtimeTo: args?.date_to ? new Date(args.date_to).getTime() / 1000 : undefined,
      sizeMin: args?.min_size,
      sizeMax: args?.max_size,
      ratingMin: args?.rating,
    }
  }
  
  return args
}
```

## 总结

通过统一使用 `GET /api/assets` 接口，前端可以：
- ✅ 用一个接口完成所有资源列表查询
- ✅ 用一个更新接口完成评分修改
- ✅ 减少代码重复，提高可维护性
- ✅ 保持向后兼容，渐进式迁移
