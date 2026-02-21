# 前端接入智能搜索 1.0

## 后端已实现的功能

### 1. 快捷筛选参数
- `quickFilter`: "recent" | "unrated" | "large" | "vertical" | "horizontal"
- `datePreset`: "today" | "thisWeek" | "thisMonth" | "lastMonth"

### 2. 搜索历史
- `GET /api/search/history?limit=10`
- `DELETE /api/search/history/clear`

### 3. 增强的 ListAssets API
```
GET /api/assets
  ?quickFilter=recent
  &datePreset=thisWeek
  &search=xxx
  &projectId=xxx
  &limit=50
  &cursor=xxx
```

---

## 前端接入方案

### 阶段 1：API 层（1小时）

#### 1.1 更新 `api.ts` 的类型定义

**文件：** `frontend/src/services/api.ts`

```typescript
// 在现有的 apiCall 函数之后添加

export interface ListAssetsParams {
  projectId?: string
  directory?: string
  search?: string
  tagIds?: string[]
  types?: string[]
  shapes?: string[]
  
  // ✅ 新增：智能搜索参数
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

// 封装的 API 调用
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
  
  return apiCall<ListAssetsResult>(`/api/assets?${queryParams.toString()}`)
}

export const getSearchHistory = async (limit = 10): Promise<SearchHistoryItem[]> => {
  return apiCall<SearchHistoryItem[]>(`/api/search/history?limit=${limit}`)
}

export const clearSearchHistory = async (): Promise<void> => {
  return apiCall('/api/search/history/clear', {}, 'DELETE')
}
```

---

### 阶段 2：UI 组件（2-3小时）

#### 2.1 快捷筛选按钮组

**文件：** `frontend/src/components/common/QuickFilters.vue`（新建）

```vue
<template>
  <div class="quick-filters">
    <button
      v-for="filter in filters"
      :key="filter.value"
      :class="['filter-btn', { active: modelValue === filter.value }]"
      @click="$emit('update:modelValue', modelValue === filter.value ? '' : filter.value)"
    >
      <Icon :name="filter.icon" size="sm" />
      <span>{{ filter.label }}</span>
    </button>
  </div>
</template>

<script setup lang="ts">
import Icon from './Icon.vue'

defineProps<{
  modelValue?: string
}>()

defineEmits<{
  'update:modelValue': [value: string]
}>()

const filters = [
  { value: 'recent', label: '最近', icon: 'clock' },
  { value: 'unrated', label: '未评分', icon: 'star' },
  { value: 'large', label: '大文件', icon: 'archive-box' },
  { value: 'vertical', label: '竖屏', icon: 'arrow-up' },
  { value: 'horizontal', label: '横屏', icon: 'arrow-right' },
]
</script>

<style scoped>
.quick-filters {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.filter-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 6px;
  color: #9ca3af;
  font-size: 13px;
  line-height: 1;
  cursor: pointer;
  transition: all 0.15s;
}

.filter-btn:hover {
  background: rgba(255, 255, 255, 0.08);
  border-color: rgba(255, 255, 255, 0.2);
  color: #e5e7eb;
}

.filter-btn.active {
  background: rgba(59, 130, 246, 0.15);
  border-color: #3b82f6;
  color: #60a5fa;
}
</style>
```

#### 2.2 日期预设选择器

**文件：** `frontend/src/components/common/DatePresets.vue`（新建）

```vue
<template>
  <div class="date-presets">
    <button
      v-for="preset in presets"
      :key="preset.value"
      :class="['preset-btn', { active: modelValue === preset.value }]"
      @click="$emit('update:modelValue', modelValue === preset.value ? '' : preset.value)"
    >
      {{ preset.label }}
    </button>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  modelValue?: string
}>()

defineEmits<{
  'update:modelValue': [value: string]
}>()

const presets = [
  { value: 'today', label: '今天' },
  { value: 'thisWeek', label: '本周' },
  { value: 'thisMonth', label: '本月' },
  { value: 'lastMonth', label: '上月' },
]
</script>

<style scoped>
.date-presets {
  display: flex;
  gap: 4px;
}

.preset-btn {
  padding: 4px 10px;
  background: transparent;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 4px;
  color: #9ca3af;
  font-size: 12px;
  line-height: 1;
  cursor: pointer;
  transition: all 0.15s;
}

.preset-btn:hover {
  border-color: rgba(255, 255, 255, 0.2);
  color: #e5e7eb;
}

.preset-btn.active {
  background: rgba(59, 130, 246, 0.15);
  border-color: #3b82f6;
  color: #60a5fa;
}
</style>
```

#### 2.3 搜索历史下拉

**文件：** `frontend/src/components/common/SearchHistory.vue`（新建）

```vue
<template>
  <div v-if="show && history.length > 0" class="search-history">
    <div class="history-header">
      <span>最近搜索</span>
      <button class="clear-btn" @click="handleClear">清空</button>
    </div>
    <div class="history-list">
      <button
        v-for="item in history"
        :key="item.id"
        class="history-item"
        @click="$emit('select', item.query)"
      >
        <Icon name="clock" size="sm" />
        <span class="query">{{ item.query }}</span>
        <span class="count">{{ item.count }}次</span>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import Icon from './Icon.vue'
import { getSearchHistory, clearSearchHistory, type SearchHistoryItem } from '@/services/api'

const props = defineProps<{
  show: boolean
}>()

const emit = defineEmits<{
  select: [query: string]
  close: []
}>()

const history = ref<SearchHistoryItem[]>([])

watch(() => props.show, async (show) => {
  if (show) {
    try {
      history.value = await getSearchHistory(10)
    } catch (e) {
      console.error('Failed to load search history:', e)
    }
  }
})

const handleClear = async () => {
  try {
    await clearSearchHistory()
    history.value = []
    emit('close')
  } catch (e) {
    console.error('Failed to clear search history:', e)
  }
}
</script>

<style scoped>
.search-history {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  margin-top: 4px;
  background: #2d2d2d;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 6px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
  z-index: 100;
  max-height: 300px;
  overflow-y: auto;
}

.history-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  font-size: 12px;
  color: #9ca3af;
}

.clear-btn {
  background: none;
  border: none;
  color: #60a5fa;
  font-size: 12px;
  cursor: pointer;
  padding: 0;
}

.clear-btn:hover {
  color: #93c5fd;
}

.history-list {
  padding: 4px;
}

.history-item {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
  padding: 8px;
  background: transparent;
  border: none;
  border-radius: 4px;
  color: #e5e7eb;
  font-size: 13px;
  text-align: left;
  cursor: pointer;
  transition: background 0.15s;
}

.history-item:hover {
  background: rgba(255, 255, 255, 0.05);
}

.query {
  flex: 1;
}

.count {
  font-size: 11px;
  color: #6b7280;
}
</style>
```

---

### 阶段 3：集成到现有页面（1-2小时）

#### 3.1 修改 PoolView.vue

**文件：** `frontend/src/components/app-shell/PoolView.vue`

```vue
<template>
  <div class="pool-view">
    <!-- 搜索栏 -->
    <div class="search-bar">
      <div class="search-input-wrapper">
        <input
          v-model="searchQuery"
          type="text"
          placeholder="搜索素材..."
          @focus="showHistory = true"
          @blur="handleSearchBlur"
          @keyup.enter="handleSearch"
        />
        <SearchHistory
          :show="showHistory"
          @select="handleHistorySelect"
          @close="showHistory = false"
        />
      </div>
      <button class="search-btn" @click="handleSearch">
        <Icon name="magnifying-glass" />
      </button>
    </div>

    <!-- ✅ 新增：快捷筛选 -->
    <div class="filters-section">
      <QuickFilters v-model="quickFilter" @update:modelValue="handleFilterChange" />
      <DatePresets v-model="datePreset" @update:modelValue="handleFilterChange" />
    </div>

    <!-- 素材列表 -->
    <div class="assets-grid">
      <!-- 现有的素材展示逻辑 -->
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import Icon from '@/components/common/Icon.vue'
import QuickFilters from '@/components/common/QuickFilters.vue'
import DatePresets from '@/components/common/DatePresets.vue'
import SearchHistory from '@/components/common/SearchHistory.vue'
import { listAssets, type ListAssetsParams } from '@/services/api'

const searchQuery = ref('')
const quickFilter = ref('')
const datePreset = ref('')
const showHistory = ref(false)

const handleSearch = async () => {
  const params: ListAssetsParams = {
    search: searchQuery.value,
    quickFilter: quickFilter.value || undefined,
    datePreset: datePreset.value || undefined,
    limit: 50,
  }
  
  try {
    const result = await listAssets(params)
    // 更新素材列表
    console.log('Assets:', result)
  } catch (e) {
    console.error('Search failed:', e)
  }
}

const handleFilterChange = () => {
  // 筛选条件变化时自动搜索
  handleSearch()
}

const handleHistorySelect = (query: string) => {
  searchQuery.value = query
  showHistory.value = false
  handleSearch()
}

const handleSearchBlur = () => {
  // 延迟关闭，让点击历史项有时间触发
  setTimeout(() => {
    showHistory.value = false
  }, 200)
}
</script>

<style scoped>
.search-bar {
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
}

.search-input-wrapper {
  position: relative;
  flex: 1;
}

.search-input-wrapper input {
  width: 100%;
  padding: 10px 16px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 6px;
  color: #e5e7eb;
  font-size: 14px;
}

.search-btn {
  padding: 10px 16px;
  background: #3b82f6;
  border: none;
  border-radius: 6px;
  color: white;
  cursor: pointer;
}

.filters-section {
  display: flex;
  gap: 16px;
  margin-bottom: 16px;
  padding-bottom: 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}
</style>
```

---

## 实现步骤

### Day 1 上午（2小时）
1. ✅ 更新 `api.ts` 类型定义和封装函数
2. ✅ 创建 `QuickFilters.vue` 组件
3. ✅ 创建 `DatePresets.vue` 组件

### Day 1 下午（2小时）
4. ✅ 创建 `SearchHistory.vue` 组件
5. ✅ 集成到 `PoolView.vue`
6. ✅ 测试基本功能

### Day 2（2小时）
7. ✅ 集成到 `ProjectView.vue`（如果需要）
8. ✅ 样式优化
9. ✅ 边界情况处理

---

## 测试清单

```bash
# 1. 快捷筛选
- [ ] 点击"最近"按钮，显示7天内素材
- [ ] 点击"未评分"按钮，显示未评分素材
- [ ] 点击"竖屏"按钮，显示竖屏素材
- [ ] 再次点击取消筛选

# 2. 日期预设
- [ ] 选择"今天"，显示今天的素材
- [ ] 选择"本周"，显示本周的素材
- [ ] 选择"本月"，显示本月的素材

# 3. 组合筛选
- [ ] 同时选择"竖屏"+"本周"
- [ ] 搜索关键词+"最近"

# 4. 搜索历史
- [ ] 搜索后，历史记录自动保存
- [ ] 点击历史记录，自动填充搜索框
- [ ] 清空历史记录
```

---

## 关键点

1. **不破坏现有功能**
   - 所有新参数都是可选的
   - 不传参数时行为与之前完全一致

2. **渐进式增强**
   - 先在 PoolView 实现
   - 验证后再推广到其他页面

3. **用户体验**
   - 快捷筛选一键切换
   - 搜索历史自动记录
   - 组合筛选自动生效

4. **性能考虑**
   - 搜索历史异步加载
   - 防抖处理（如果需要）
   - 分页保持稳定
