# 评分功能使用指南

## 概述

前端已经实现了统一的资源评分功能，支持：
- 查看资源评分（1-5星）
- 设置资源评分
- 清除资源评分

## 后端接口

### 1. 获取资源列表（包含评分）

```http
GET /api/assets?projectId=xxx&limit=50
```

**返回示例：**
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "asset-123",
        "path": "/path/to/file.jpg",
        "size": 1024000,
        "mtime": 1640995200,
        "user_rating": 5,
        "width": 1920,
        "height": 1080
      }
    ],
    "total": 100,
    "nextCursor": "cursor-token",
    "hasMore": true
  }
}
```

### 2. 更新评分

```http
POST /api/assets/update
Content-Type: application/json

{
  "id": "asset-123",
  "user_rating": 5
}
```

### 3. 清除评分

```http
POST /api/assets/update
Content-Type: application/json

{
  "id": "asset-123",
  "clear_user_rating": true
}
```

## 前端使用

### 1. 在 Store 中使用

```typescript
import { useFileStore } from '@/stores/fileStore'

const fileStore = useFileStore()

// 加载资源（自动包含评分）
await fileStore.loadFiles()

// 更新评分
await fileStore.updateRating('asset-123', 5)

// 清除评分
await fileStore.updateRating('asset-123', null)
```

### 2. 在组件中使用

```vue
<template>
  <div class="asset-card">
    <img :src="asset.thumbnail_path" :alt="asset.name" />
    <div class="asset-info">
      <h3>{{ asset.name }}</h3>
      <RatingControl
        :asset-id="asset.id"
        :rating="asset.rating"
        @update="handleRatingUpdate"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useFileStore } from '@/stores/fileStore'
import RatingControl from '@/components/common/RatingControl.vue'

const fileStore = useFileStore()

const handleRatingUpdate = async (assetId: string, rating: number | null) => {
  try {
    await fileStore.updateRating(assetId, rating)
    console.log('评分更新成功')
  } catch (error) {
    console.error('评分更新失败:', error)
  }
}
</script>
```

### 3. 按评分筛选

```typescript
import { apiCall } from '@/services/api'

// 获取 4-5 星的资源
const assets = await apiCall('list_assets', {
  projectId: 'proj-123',
  ratingMin: 4,
  ratingMax: 5,
  limit: 50
})

// 获取未评分的资源
const unratedAssets = await apiCall('list_assets', {
  projectId: 'proj-123',
  ratingMin: 0,
  ratingMax: 0,
  limit: 50
})

// 获取特定评分的资源
const fiveStarAssets = await apiCall('list_assets', {
  projectId: 'proj-123',
  rating: 5,
  limit: 50
})
```

## RatingControl 组件

### Props

- `assetId: string` - 资源ID（必需）
- `rating?: number` - 当前评分（0-5）

### Events

- `@update(assetId: string, rating: number | null)` - 评分更新事件

### 使用示例

```vue
<RatingControl
  :asset-id="asset.id"
  :rating="asset.rating"
  @update="handleRatingUpdate"
/>
```

## 数据流

```
用户点击星星
    ↓
RatingControl 触发 @update 事件
    ↓
父组件调用 fileStore.updateRating()
    ↓
发送 POST /api/assets/update
    ↓
后端更新数据库
    ↓
前端更新本地状态
    ↓
UI 自动更新
```

## 注意事项

1. **评分范围**：1-5 星，0 表示未评分
2. **清除评分**：传 `null` 而不是 `0`
3. **本地更新**：`updateRating` 会自动更新本地状态，无需重新加载列表
4. **错误处理**：更新失败时会抛出异常，需要在调用处捕获

## 完整示例

```vue
<template>
  <div class="asset-list">
    <div
      v-for="asset in fileStore.files"
      :key="asset.id"
      class="asset-item"
    >
      <img :src="getThumbnailUrl(asset)" :alt="asset.name" />
      <div class="asset-details">
        <h4>{{ asset.name }}</h4>
        <p>{{ formatSize(asset.size) }}</p>
        <RatingControl
          :asset-id="asset.id"
          :rating="asset.rating"
          @update="handleRatingUpdate"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useFileStore } from '@/stores/fileStore'
import RatingControl from '@/components/common/RatingControl.vue'

const fileStore = useFileStore()

onMounted(async () => {
  fileStore.setProjectId('proj-123')
  await fileStore.loadFiles()
})

const handleRatingUpdate = async (assetId: string, rating: number | null) => {
  try {
    await fileStore.updateRating(assetId, rating)
  } catch (error) {
    console.error('更新评分失败:', error)
    alert('更新评分失败，请重试')
  }
}

const getThumbnailUrl = (asset: any) => {
  return asset.thumbnail_path || '/placeholder.png'
}

const formatSize = (bytes: number) => {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
}
</script>

<style scoped>
.asset-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 1rem;
  padding: 1rem;
}

.asset-item {
  border: 1px solid var(--color-border);
  border-radius: 0.5rem;
  overflow: hidden;
  transition: transform 0.2s;
}

.asset-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.asset-item img {
  width: 100%;
  height: 150px;
  object-fit: cover;
}

.asset-details {
  padding: 0.75rem;
}

.asset-details h4 {
  margin: 0 0 0.5rem 0;
  font-size: 0.875rem;
  font-weight: 500;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.asset-details p {
  margin: 0 0 0.5rem 0;
  font-size: 0.75rem;
  color: var(--color-text-secondary);
}
</style>
```

## API 参数完整列表

### GET /api/assets

| 参数 | 类型 | 说明 | 示例 |
|------|------|------|------|
| projectId | string | 项目ID | `proj-123` |
| search | string | 搜索关键词 | `vacation` |
| tagIds | string | 标签ID（逗号分隔） | `tag-1,tag-2` |
| types | string | 文件类型（逗号分隔） | `image,video` |
| shapes | string | 形状（逗号分隔） | `landscape,portrait` |
| rating | int | 精确评分 | `5` |
| ratingMin | int | 最低评分 | `4` |
| ratingMax | int | 最高评分 | `5` |
| sizeMin | int64 | 最小文件大小（字节） | `1048576` |
| sizeMax | int64 | 最大文件大小（字节） | `10485760` |
| widthMin | int | 最小宽度 | `1920` |
| widthMax | int | 最大宽度 | `3840` |
| heightMin | int | 最小高度 | `1080` |
| heightMax | int | 最大高度 | `2160` |
| mtimeFrom | int64 | 开始时间（Unix时间戳） | `1640995200` |
| mtimeTo | int64 | 结束时间（Unix时间戳） | `1672531200` |
| sortBy | string | 排序字段 | `mtime`, `size`, `name` |
| sortOrder | string | 排序方向 | `asc`, `desc` |
| cursor | string | 分页游标 | `cursor-token` |
| limit | int | 每页数量 | `50` |

### POST /api/assets/update

```typescript
interface UpdateAssetRequest {
  id: string                    // 资源ID（必需）
  user_rating?: number          // 用户评分 1-5
  clear_user_rating?: boolean   // 清除评分
  media_meta?: string           // 媒体元数据（JSON字符串）
}
```
