<template>
  <div class="artifact-view">
    <FilterBar
      :filters="filters"
      :active-filter="activeFilter"
      @update:active-filter="activeFilter = $event"
    />

    <div class="artifact-content">
      <EmptyState
        v-if="filteredArtifacts.length === 0"
        icon="package"
        title="暂无交付物"
        description="创建并导出项目成品"
      />

      <div v-else class="artifact-grid">
        <ArtifactCard
          v-for="artifact in filteredArtifacts"
          :key="artifact.id"
          :artifact="artifact"
          :is-selected="selectedIds.has(artifact.id)"
          @select="handleSelect(artifact)"
          @open="handleOpen(artifact)"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ArtifactCard, FilterBar } from '@/components/artifact'
import EmptyState from '@/components/common/EmptyState.vue'

interface Artifact {
  id: string
  name: string
  thumbnail?: string
  status: string
  project_name: string
  created_at: string
}

const props = defineProps<{
  activeFilter: string
}>()

const activeFilter = ref('all')
const selectedIds = ref<Set<string>>(new Set())

const mockArtifacts = ref<Artifact[]>([
  {
    id: '1',
    name: '宣传视频_v1.mp4',
    status: 'published',
    project_name: '品牌宣传片',
    created_at: '2026-02-15'
  },
  {
    id: '2',
    name: '产品介绍_草稿.mp4',
    status: 'draft',
    project_name: '产品发布会',
    created_at: '2026-02-18'
  }
])

const filters = computed(() => [
  { key: 'all', label: '全部', count: mockArtifacts.value.length },
  { key: 'draft', label: '草稿', count: mockArtifacts.value.filter(a => a.status === 'draft').length },
  { key: 'review', label: '审核中', count: mockArtifacts.value.filter(a => a.status === 'review').length },
  { key: 'published', label: '已发布', count: mockArtifacts.value.filter(a => a.status === 'published').length },
  { key: 'archived', label: '已归档', count: mockArtifacts.value.filter(a => a.status === 'archived').length }
])

const filteredArtifacts = computed(() => {
  if (activeFilter.value === 'all') return mockArtifacts.value
  return mockArtifacts.value.filter(a => a.status === activeFilter.value)
})

function handleSelect(artifact: Artifact) {
  if (selectedIds.value.has(artifact.id)) {
    selectedIds.value.delete(artifact.id)
  } else {
    selectedIds.value.add(artifact.id)
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
  background: #1a1a1a;
}

.artifact-content {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}

.artifact-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
}
</style>
