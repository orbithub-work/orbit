# API 优化完成总结

## 已完成的工作

### 1. 后端接口（已存在，无需修改）

✅ **GET /api/assets** - 统一的资源列表接口
- 支持所有筛选参数（projectId, search, tagIds, types, shapes, rating, size, dimensions, time）
- 支持排序和分页
- 返回统一结构：`{ items, nextCursor, hasMore, total }`

✅ **POST /api/assets/update** - 资源更新接口
- 支持更新评分：`{ id, user_rating }` 或 `{ id, clear_user_rating: true }`
- 支持更新元数据：`{ id, media_meta }`

### 2. 前端优化

#### 2.1 API 服务层 (`frontend/src/services/api.ts`)

✅ 统一接口映射
```typescript
'list_assets': '/api/assets',
'get_project_files': '/api/assets',  // 重定向到统一接口
'update_asset': '/api/assets/update',
```

✅ 参数兼容性转换
```typescript
// 兼容旧的 get_project_files 接口
if (command === 'get_project_files' && args.projectId) {
  return { projectId: args.projectId, limit: args.limit || 50 }
}
```

#### 2.2 状态管理 (`frontend/src/stores/fileStore.ts`)

✅ 更新 `loadFiles` 方法使用统一接口
```typescript
const result = await apiCall<{ items: AssetDto[], total: number }>('list_assets', {
  projectId: projectId.value,
  limit: 1000
})
```

✅ 添加 `rating` 字段到 `FileItem` 接口
```typescript
export interface FileItem {
  // ... 其他字段
  rating?: number
}
```

✅ 新增 `updateRating` 方法
```typescript
async function updateRating(assetId: string, rating: number | null) {
  await apiCall('update_asset', {
    id: assetId,
    ...(rating !== null ? { user_rating: rating } : { clear_user_rating: true })
  })
  
  // 更新本地状态
  const file = files.value.find(f => f.id === assetId)
  if (file) {
    file.rating = rating ?? undefined
  }
}
```

#### 2.3 UI 组件

✅ 创建 `RatingControl.vue` 组件
- 5星评分控件
- 支持清除评分
- 响应式更新

### 3. 文档

✅ 创建 `docs/api-simplification.md` - API 简化方案
✅ 创建 `docs/rating-feature-guide.md` - 评分功能使用指南

## 使用示例

### 基础用法

```typescript
import { useFileStore } from '@/stores/fileStore'

const fileStore = useFileStore()

// 1. 加载资源列表（包含评分）
fileStore.setProjectId('proj-123')
await fileStore.loadFiles()

// 2. 更新评分
await fileStore.updateRating('asset-123', 5)

// 3. 清除评分
await fileStore.updateRating('asset-123', null)
```

### 在组件中使用

```vue
<template>
  <RatingControl
    :asset-id="asset.id"
    :rating="asset.rating"
    @update="handleRatingUpdate"
  />
</template>

<script setup lang="ts">
import { useFileStore } from '@/stores/fileStore'
import RatingControl from '@/components/common/RatingControl.vue'

const fileStore = useFileStore()

const handleRatingUpdate = async (assetId: string, rating: number | null) => {
  await fileStore.updateRating(assetId, rating)
}
</script>
```

### 高级筛选

```typescript
import { apiCall } from '@/services/api'

// 获取 4-5 星的图片
const assets = await apiCall('list_assets', {
  projectId: 'proj-123',
  types: ['image'],
  ratingMin: 4,
  ratingMax: 5,
  sortBy: 'mtime',
  sortOrder: 'desc',
  limit: 50
})
```

## 优势

1. **统一接口** - 一个接口完成所有资源列表查询
2. **简化代码** - 减少重复代码，提高可维护性
3. **类型安全** - 统一的返回结构，更好的 TypeScript 支持
4. **性能优化** - 后端统一处理，可以更好地优化查询
5. **扩展性强** - 新增筛选条件只需添加参数，无需新接口
6. **向后兼容** - 保留旧接口名称，渐进式迁移

## 接口对比

### 优化前

```typescript
// 多个不同的接口
await apiCall('get_project_files', { projectId: 'xxx' })
await apiCall('search_files', { query: 'xxx' })
await apiCall('advanced_search', { ... })
```

### 优化后

```typescript
// 统一使用 list_assets
await apiCall('list_assets', { projectId: 'xxx' })
await apiCall('list_assets', { search: 'xxx' })
await apiCall('list_assets', { projectId: 'xxx', types: ['image'], ratingMin: 4 })
```

## 下一步建议

### 1. 逐步迁移其他接口

可以考虑将以下接口也统一到 `list_assets`：
- `search_files` → `list_assets` with `search` param
- `advanced_search` → `list_assets` with multiple filters

### 2. 添加更多筛选条件

后端可以考虑添加：
- `hasTag` - 是否有标签
- `isArchived` - 是否已归档
- `colorTags` - 颜色标签
- `collections` - 所属收藏夹

### 3. 优化分页体验

- 实现虚拟滚动 + 无限加载
- 使用 `cursor` 实现高效分页
- 缓存已加载的数据

### 4. 添加批量操作

```typescript
// 批量更新评分
POST /api/assets/batch-update
{
  "updates": [
    { "id": "asset-1", "user_rating": 5 },
    { "id": "asset-2", "user_rating": 4 }
  ]
}
```

## 测试清单

- [ ] 加载项目资源列表
- [ ] 搜索资源
- [ ] 按类型筛选
- [ ] 按评分筛选
- [ ] 按尺寸筛选
- [ ] 按时间筛选
- [ ] 设置评分（1-5星）
- [ ] 清除评分
- [ ] 排序功能
- [ ] 分页加载
- [ ] 本地状态更新
- [ ] 错误处理

## 相关文件

### 后端
- `internal/httpapi/handler_asset.go` - 资源接口处理
- `internal/services/asset_service.go` - 资源服务

### 前端
- `frontend/src/services/api.ts` - API 服务
- `frontend/src/stores/fileStore.ts` - 文件状态管理
- `frontend/src/components/common/RatingControl.vue` - 评分控件

### 文档
- `docs/api-simplification.md` - API 简化方案
- `docs/rating-feature-guide.md` - 评分功能使用指南
- `docs/api/openapi.yaml` - API 规范

## 总结

通过这次优化，我们实现了：

1. ✅ **统一的资源列表接口** - 一个接口替代多个接口
2. ✅ **完整的评分功能** - 读取、设置、清除评分
3. ✅ **向后兼容** - 保留旧接口名称，平滑迁移
4. ✅ **类型安全** - TypeScript 类型定义完善
5. ✅ **文档完善** - 详细的使用指南和示例

前端现在可以：
- 用一个接口完成所有资源列表查询
- 用一个更新接口完成评分修改
- 享受更好的代码可维护性和扩展性
