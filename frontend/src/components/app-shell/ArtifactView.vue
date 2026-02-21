<template>
  <div class="artifact-view">
    <header class="artifact-header">
      <div class="header-top-row">
        <div class="header-left">
          <span class="header-title">制品</span>
        </div>
        <div class="header-center">
          <div class="search-box">
            <svg class="search-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="11" cy="11" r="8"></circle>
              <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
            </svg>
            <input class="search-input" type="text" placeholder="搜索制品" v-model="searchQuery" />
          </div>
        </div>
        <div class="header-right">
          <button class="view-btn" :class="{ active: viewMode === 'grid' }" title="网格视图" @click="viewMode = 'grid'">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="3" y="3" width="7" height="7"></rect>
              <rect x="14" y="3" width="7" height="7"></rect>
              <rect x="14" y="14" width="7" height="7"></rect>
              <rect x="3" y="14" width="7" height="7"></rect>
            </svg>
          </button>
          <button class="view-btn" :class="{ active: viewMode === 'list' }" title="列表视图" @click="viewMode = 'list'">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="8" y1="6" x2="21" y2="6"></line>
              <line x1="8" y1="12" x2="21" y2="12"></line>
              <line x1="8" y1="18" x2="21" y2="18"></line>
              <line x1="3" y1="6" x2="3.01" y2="6"></line>
              <line x1="3" y1="12" x2="3.01" y2="12"></line>
              <line x1="3" y1="18" x2="3.01" y2="18"></line>
            </svg>
          </button>
        </div>
      </div>
    </header>

    <div class="artifact-count-bar">
      <span class="count-text">制品 ({{ filteredArtifacts.length }})</span>
    </div>

    <div class="artifact-grid-container">
      <EmptyState
        v-if="filteredArtifacts.length === 0"
        icon="📦"
        title="暂无制品"
        description="还没有任何交付制品，完成项目后可在此查看"
      />
      <div v-else-if="viewMode === 'grid'" class="artifact-grid">
        <div
          v-for="artifact in filteredArtifacts"
          :key="artifact.id"
          class="artifact-card"
          :class="{ 'artifact-card--selected': selectedIds.has(artifact.id) }"
          @click="handleSelect(artifact, $event)"
          @dblclick="handleOpen(artifact)"
        >
          <div class="card-thumbnail">
            <div class="thumbnail-placeholder">
              <span class="file-icon">{{ getFileIcon(artifact) }}</span>
            </div>
            <div class="card-status" :class="`status--${artifact.status}`">
              {{ getStatusText(artifact.status) }}
            </div>
          </div>
          <div class="card-content">
            <div class="card-title">{{ artifact.name }}</div>
            <div class="card-meta">
              <span class="card-version">v{{ artifact.version }}</span>
              <div class="card-platforms">
                <span
                  v-for="p in artifact.platforms"
                  :key="p.platform"
                  class="platform-badge"
                  :title="getPlatformName(p.platform)"
                >
                  {{ getPlatformIcon(p.platform) }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div v-else class="artifact-list">
        <div
          v-for="artifact in filteredArtifacts"
          :key="artifact.id"
          class="artifact-list-item"
          :class="{ 'artifact-list-item--selected': selectedIds.has(artifact.id) }"
          @click="handleSelect(artifact, $event)"
          @dblclick="handleOpen(artifact)"
        >
          <div class="list-icon">{{ getFileIcon(artifact) }}</div>
          <div class="list-name">{{ artifact.name }}</div>
          <div class="list-status" :class="`status--${artifact.status}`">
            {{ getStatusText(artifact.status) }}
          </div>
          <div class="list-version">v{{ artifact.version }}</div>
          <div class="list-platforms">
            <span
              v-for="p in artifact.platforms"
              :key="p.platform"
              class="platform-badge"
            >
              {{ getPlatformIcon(p.platform) }}
            </span>
          </div>
          <div class="list-date">{{ formatDate(artifact.updated_at) }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useArtifactStore } from '@/stores/artifactStore'
import { PLATFORMS } from '@/types/artifact'
import type { Artifact } from '@/types/artifact'
import EmptyState from '@/components/common/EmptyState.vue'

const artifactStore = useArtifactStore()

const searchQuery = ref('')
const viewMode = ref<'grid' | 'list'>('grid')
const selectedIds = ref<Set<string>>(new Set())

const props = defineProps<{
  activeFilter?: string
}>()

const filteredArtifacts = computed(() => {
  let result = artifactStore.artifacts

  if (props.activeFilter && props.activeFilter !== 'all') {
    if (props.activeFilter.startsWith('status:')) {
      const status = props.activeFilter.replace('status:', '') as Artifact['status']
      result = result.filter(a => a.status === status)
    } else if (props.activeFilter.startsWith('platform:')) {
      const platformId = props.activeFilter.replace('platform:', '')
      result = result.filter(a => a.platforms.some(p => p.platform === platformId))
    } else if (props.activeFilter.startsWith('project:')) {
      const projectId = props.activeFilter.replace('project:', '')
      result = result.filter(a => a.project_id === projectId)
    } else if (['draft', 'published', 'archived'].includes(props.activeFilter)) {
      result = result.filter(a => a.status === props.activeFilter)
    }
  }

  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(a => 
      a.name.toLowerCase().includes(query) ||
      a.tags.some(t => t.toLowerCase().includes(query))
    )
  }

  return result
})

function getFileIcon(artifact: Artifact): string {
  const ext = artifact.file_path.split('.').pop()?.toLowerCase() || ''
  if (['jpg', 'jpeg', 'png', 'gif', 'webp'].includes(ext)) return '🖼️'
  if (['mp4', 'mov', 'avi', 'mkv'].includes(ext)) return '🎬'
  if (['mp3', 'wav', 'flac'].includes(ext)) return '🎵'
  if (['pdf'].includes(ext)) return '📄'
  return '📄'
}

function getStatusText(status: Artifact['status']): string {
  const map: Record<Artifact['status'], string> = {
    draft: '草稿',
    published: '已发布',
    archived: '已归档'
  }
  return map[status]
}

function getPlatformIcon(platformId: string): string {
  const platform = PLATFORMS.find(p => p.id === platformId)
  return platform?.icon || '🔗'
}

function getPlatformName(platformId: string): string {
  const platform = PLATFORMS.find(p => p.id === platformId)
  return platform?.name || '其他'
}

function formatDate(dateStr: string): string {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN', { month: 'short', day: 'numeric' })
}

function handleSelect(artifact: Artifact, event: MouseEvent) {
  if (event.ctrlKey || event.metaKey) {
    const newSet = new Set(selectedIds.value)
    if (newSet.has(artifact.id)) {
      newSet.delete(artifact.id)
    } else {
      newSet.add(artifact.id)
    }
    selectedIds.value = newSet
  } else {
    selectedIds.value = new Set([artifact.id])
  }
}

function handleOpen(artifact: Artifact) {
  console.log('Open artifact:', artifact)
}
</script>

<style scoped>
.artifact-view {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: #1b1c1f;
  overflow: hidden;
}

.artifact-header {
  display: flex;
  flex-direction: column;
  background: #1b1c1f;
  border-bottom: 1px solid #2b2b2f;
  flex-shrink: 0;
}

.header-top-row {
  height: 48px;
  display: grid;
  grid-template-columns: 1fr auto 1fr;
  align-items: center;
  padding: 0 16px;
  gap: 8px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.header-title {
  font-size: 14px;
  font-weight: 600;
  color: #e5e7eb;
}

.header-center {
  display: flex;
  justify-content: center;
}

.search-box {
  position: relative;
  display: flex;
  align-items: center;
}

.search-icon {
  position: absolute;
  left: 12px;
  color: #6b7280;
}

.search-input {
  width: 320px;
  height: 36px;
  background: #2a2a2a;
  border: 1px solid #3a3a3a;
  border-radius: 20px;
  padding: 0 12px 0 36px;
  font-size: 13px;
  color: #e5e7eb;
  outline: none;
  transition: all 0.2s;
}

.search-input:focus {
  border-color: #3b82f6;
  background: #252525;
}

.header-right {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 4px;
}

.view-btn {
  background: transparent;
  border: none;
  cursor: pointer;
  padding: 8px;
  border-radius: 6px;
  color: #9ca3af;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.view-btn:hover {
  background: rgba(255, 255, 255, 0.08);
  color: #e5e7eb;
}

.view-btn.active {
  background: rgba(59, 130, 246, 0.15);
  color: #60a5fa;
}

.artifact-count-bar {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 8px 16px;
  background: #1b1c1f;
  color: #9ca3af;
  font-size: 13px;
  border-bottom: 1px solid #2b2b2f;
}

.count-text {
  color: #e5e7eb;
}

.artifact-grid-container {
  flex: 1;
  min-height: 0;
  overflow: auto;
  padding: 16px;
}

.artifact-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 16px;
}

.artifact-card {
  background: #252526;
  border-radius: 12px;
  overflow: hidden;
  cursor: pointer;
  border: 2px solid transparent;
  transition: all 0.2s ease-in-out;
}

.artifact-card:hover {
  background: #2a2a2a;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
}

.artifact-card--selected {
  border-color: #3b82f6;
  background: rgba(59, 130, 246, 0.1);
}

.card-thumbnail {
  width: 100%;
  aspect-ratio: 16/10;
  background: #1e1e1e;
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
}

.thumbnail-placeholder {
  color: #6b7280;
}

.file-icon {
  font-size: 32px;
}

.card-status {
  position: absolute;
  top: 8px;
  right: 8px;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 10px;
  font-weight: 600;
}

.status--draft {
  background: rgba(251, 191, 36, 0.2);
  color: #fbbf24;
}

.status--published {
  background: rgba(34, 197, 94, 0.2);
  color: #22c55e;
}

.status--archived {
  background: rgba(107, 114, 128, 0.2);
  color: #9ca3af;
}

.card-content {
  padding: 12px;
}

.card-title {
  font-size: 13px;
  font-weight: 500;
  color: #e5e7eb;
  margin-bottom: 8px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.card-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.card-version {
  font-size: 11px;
  color: #6b7280;
  background: rgba(255, 255, 255, 0.06);
  padding: 2px 6px;
  border-radius: 4px;
}

.card-platforms {
  display: flex;
  gap: 4px;
}

.platform-badge {
  font-size: 12px;
}

.artifact-list {
  display: flex;
  flex-direction: column;
  gap: 1px;
}

.artifact-list-item {
  display: grid;
  grid-template-columns: 32px 1fr 80px 60px 100px 100px;
  align-items: center;
  padding: 12px 16px;
  background: #1e1e1e;
  cursor: pointer;
  transition: background 0.15s;
}

.artifact-list-item:hover {
  background: #252526;
}

.artifact-list-item--selected {
  background: rgba(59, 130, 246, 0.15);
}

.list-icon {
  font-size: 16px;
}

.list-name {
  font-size: 13px;
  color: #e5e7eb;
}

.list-status {
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 4px;
  text-align: center;
}

.list-version {
  font-size: 12px;
  color: #6b7280;
}

.list-platforms {
  display: flex;
  gap: 4px;
}

.list-date {
  font-size: 12px;
  color: #6b7280;
}
</style>
