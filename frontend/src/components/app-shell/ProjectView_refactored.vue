<template>
  <div class="project-view">
    <EmptyState
      v-if="!project"
      icon="folder"
      title="未选择项目"
      description="从左侧选择一个项目查看详情"
    />

    <ProjectTemplateGuide
      v-else-if="!project.template_applied"
      @select-template="handleSelectTemplate"
    />

    <div v-else class="project-content">
      <div class="project-header">
        <h2>{{ project.name }}</h2>
        <div class="project-actions">
          <button class="action-btn">
            <Icon name="settings" size="sm" />
            设置
          </button>
        </div>
      </div>

      <div class="stats-grid">
        <StatCard icon="image" label="素材" :value="projectStats.assetCount" color="rgba(59, 130, 246, 0.2)" />
        <StatCard icon="file" label="工程文件" :value="projectStats.engineCount" color="rgba(168, 85, 247, 0.2)" />
        <StatCard icon="package" label="成品" :value="projectStats.artifactCount" color="rgba(34, 197, 94, 0.2)" />
      </div>

      <div class="sections-grid">
        <SectionCard icon="image" title="素材" :count="projectAssets.length">
          <template #actions>
            <button class="icon-btn" @click="handleAddAsset">
              <Icon name="plus" size="sm" />
            </button>
          </template>
          <EmptyState
            v-if="projectAssets.length === 0"
            icon="image"
            title="暂无素材"
            description="添加素材到项目"
            compact
          />
          <div v-else class="asset-grid">
            <div v-for="asset in projectAssets" :key="asset.id" class="asset-item">
              <img v-if="asset.thumbnail_url" :src="asset.thumbnail_url" :alt="asset.name" />
              <div v-else class="asset-placeholder">
                <Icon name="file" size="md" />
              </div>
              <span class="asset-name">{{ asset.name }}</span>
            </div>
          </div>
        </SectionCard>

        <SectionCard icon="file" title="工程文件" :count="engineFiles.length">
          <template #actions>
            <button class="icon-btn" @click="handleOpenFolder">
              <Icon name="folder-open" size="sm" />
            </button>
          </template>
          <EmptyState
            v-if="engineFiles.length === 0"
            icon="file"
            title="暂无工程文件"
            description="在绑定目录中创建工程文件"
            compact
          />
          <div v-else class="file-list">
            <div v-for="file in engineFiles" :key="file.path" class="file-item">
              <Icon name="file" size="sm" />
              <span>{{ file.name }}</span>
            </div>
          </div>
        </SectionCard>

        <SectionCard icon="package" title="成品" :count="artifacts.length">
          <template #actions>
            <button class="icon-btn" @click="handleExport">
              <Icon name="upload" size="sm" />
            </button>
          </template>
          <EmptyState
            v-if="artifacts.length === 0"
            icon="package"
            title="暂无成品"
            description="导出项目成品"
            compact
          />
          <div v-else class="file-list">
            <div v-for="artifact in artifacts" :key="artifact.id" class="file-item">
              <Icon name="package" size="sm" />
              <span>{{ artifact.name }}</span>
            </div>
          </div>
        </SectionCard>
      </div>

      <ProjectLogTimeline :project-id="project.id" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { StatCard, SectionCard, ProjectTemplateGuide, ProjectLogTimeline } from '@/components/project'
import EmptyState from '@/components/common/EmptyState.vue'
import Icon from '@/components/common/Icon.vue'
import type { Project } from '@/stores/projectStore'

const props = defineProps<{
  project?: Project
}>()

const projectAssets = ref<any[]>([])
const engineFiles = ref<any[]>([])
const artifacts = ref<any[]>([])

const projectStats = computed(() => ({
  assetCount: projectAssets.value.length,
  engineCount: engineFiles.value.length,
  artifactCount: artifacts.value.length
}))

function handleSelectTemplate(template: string) {
  console.log('Selected template:', template)
}

function handleAddAsset() {
  console.log('Add asset')
}

function handleOpenFolder() {
  console.log('Open folder')
}

function handleExport() {
  console.log('Export')
}
</script>

<style scoped>
.project-view {
  height: 100%;
  overflow-y: auto;
  background: #1a1a1a;
}

.project-content {
  padding: 24px;
}

.project-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
}

.project-header h2 {
  font-size: 24px;
  font-weight: 600;
  color: #e5e7eb;
  margin: 0;
}

.project-actions {
  display: flex;
  gap: 8px;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 6px;
  color: #e5e7eb;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.15s;
}

.action-btn:hover {
  background: rgba(255, 255, 255, 0.08);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
  margin-bottom: 24px;
}

.sections-grid {
  display: grid;
  gap: 24px;
  margin-bottom: 24px;
}

.icon-btn {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 6px;
  color: #9ca3af;
  cursor: pointer;
  transition: all 0.15s;
}

.icon-btn:hover {
  background: rgba(255, 255, 255, 0.08);
  color: #e5e7eb;
}

.asset-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
  gap: 12px;
}

.asset-item {
  display: flex;
  flex-direction: column;
  gap: 6px;
  cursor: pointer;
}

.asset-item img,
.asset-placeholder {
  aspect-ratio: 1;
  border-radius: 6px;
  background: rgba(0, 0, 0, 0.2);
  object-fit: cover;
}

.asset-placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  color: #6b7280;
}

.asset-name {
  font-size: 12px;
  color: #9ca3af;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.file-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.file-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background: rgba(255, 255, 255, 0.03);
  border-radius: 6px;
  color: #e5e7eb;
  font-size: 13px;
  cursor: pointer;
  transition: background 0.15s;
}

.file-item:hover {
  background: rgba(255, 255, 255, 0.06);
}
</style>
