<template>
  <div class="artifact-sidebar-container">
    <div class="sidebar-group">
      <div class="group-title">状态</div>
      <div
        class="nav-item"
        :class="{ active: activeFilter === 'all' }"
        @click="setFilter('all')"
      >
        <Icon name="archive" size="sm" class="nav-icon" />
        <span class="label">全部制品</span>
        <span class="count">{{ totalCount }}</span>
      </div>
      <div
        class="nav-item"
        :class="{ active: activeFilter === 'draft' }"
        @click="setFilter('draft')"
      >
        <Icon name="edit-box" size="sm" class="nav-icon" />
        <span class="label">草稿</span>
        <span class="count">{{ draftCount }}</span>
      </div>
      <div
        class="nav-item"
        :class="{ active: activeFilter === 'published' }"
        @click="setFilter('published')"
      >
        <Icon name="check" size="sm" class="nav-icon" />
        <span class="label">已发布</span>
        <span class="count">{{ publishedCount }}</span>
      </div>
      <div
        class="nav-item"
        :class="{ active: activeFilter === 'archived' }"
        @click="setFilter('archived')"
      >
        <Icon name="folder" size="sm" class="nav-icon" />
        <span class="label">已归档</span>
        <span class="count">{{ archivedCount }}</span>
      </div>
    </div>

    <div class="sidebar-group">
      <div class="group-title">发布平台</div>
      <div
        v-for="platform in platforms"
        :key="platform.id"
        class="nav-item"
        :class="{ active: activeFilter === `platform:${platform.id}` }"
        @click="setFilter(`platform:${platform.id}`)"
      >
        <Icon name="globe" size="sm" class="nav-icon" />
        <span class="label">{{ platform.name }}</span>
        <span class="count">{{ getPlatformCount(platform.id) }}</span>
      </div>
    </div>

    <div class="sidebar-group flex-grow">
      <div class="group-title">所属项目</div>
      <div
        v-for="project in projects"
        :key="project.id"
        class="nav-item"
        :class="{ active: activeFilter === `project:${project.id}` }"
        @click="setFilter(`project:${project.id}`)"
      >
        <Icon name="folder" size="sm" class="nav-icon" />
        <span class="label">{{ project.name }}</span>
        <span class="count">{{ getProjectCount(project.id) }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import Icon from '@/components/common/Icon.vue'
import { PLATFORMS } from '@/types/artifact'
import { useArtifactStore } from '@/stores/artifactStore'

interface Project {
  id: string
  name: string
}

const props = defineProps<{
  projects: Project[]
}>()

const emit = defineEmits<{
  'update:activeFilter': [value: string]
}>()

const artifactStore = useArtifactStore()
const activeFilter = ref('all')

const platforms = PLATFORMS

const totalCount = computed(() => artifactStore.artifacts.length)
const draftCount = computed(() => artifactStore.draftArtifacts.length)
const publishedCount = computed(() => artifactStore.publishedArtifacts.length)
const archivedCount = computed(() => artifactStore.archivedArtifacts.length)

function getPlatformCount(platformId: string) {
  return artifactStore.getArtifactsByPlatform(platformId).length
}

function getProjectCount(projectId: string) {
  return artifactStore.getArtifactsByProject(projectId).length
}

function setFilter(filter: string) {
  activeFilter.value = filter
  emit('update:activeFilter', filter)
}
</script>

<style scoped>
.artifact-sidebar-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
}

.sidebar-group {
  padding: 12px 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  flex-shrink: 0;
}

.sidebar-group:last-child {
  border-bottom: none;
}

.sidebar-group.flex-grow {
  flex: 1;
  overflow-y: auto;
  min-height: 0;
}

.group-title {
  padding: 0 16px 8px;
  font-size: 11px;
  font-weight: 600;
  color: #6b7280;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.nav-item {
  display: flex;
  align-items: center;
  padding: 8px 16px;
  cursor: pointer;
  transition: background 0.15s;
}

.nav-item:hover {
  background: rgba(255, 255, 255, 0.06);
}

.nav-item.active {
  background: rgba(59, 130, 246, 0.15);
}

.nav-item .icon {
  width: 20px;
  font-size: 14px;
  text-align: center;
  margin-right: 10px;
}

.nav-item .label {
  flex: 1;
  font-size: 12px;
  color: #d1d5db;
}

.nav-item.active .label {
  color: #60a5fa;
}

.nav-item .count {
  font-size: 11px;
  color: #6b7280;
  background: rgba(255, 255, 255, 0.06);
  padding: 2px 6px;
  border-radius: 4px;
}
</style>
