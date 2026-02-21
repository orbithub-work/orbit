<template>
  <div class="project-view">
    <div class="project-hero">
      <div class="project-hero-title">{{ project?.name }}</div>
      <div class="project-hero-meta">
        {{ project?.client }} · {{ project?.deadline }}
      </div>
    </div>

    <div class="project-panels">
      <div class="project-panel">
        <div class="panel-title">项目摘要</div>
        <div class="panel-row">
          <span>内容类型</span>
          <span>{{ project?.type }}</span>
        </div>
        <div class="panel-row">
          <span>负责人</span>
          <span>{{ project?.owner }}</span>
        </div>
        <div class="panel-row">
          <span>状态</span>
          <span>{{ project?.status }}</span>
        </div>
      </div>
      <div class="project-panel">
        <div class="panel-title">素材构成</div>
        <div class="panel-row">
          <span>参考素材</span>
          <span>{{ project?.references }}</span>
        </div>
        <div class="panel-row">
          <span>交付文件</span>
          <span>{{ project?.deliverables }}</span>
        </div>
        <div class="panel-row">
          <span>版本</span>
          <span>{{ project?.versions }}</span>
        </div>
      </div>
    </div>

    <div class="project-section">
      <div class="section-title">最近素材</div>
      <EmptyState
        v-if="!project"
        icon="📋"
        title="暂无项目"
        description="请先创建或选择一个项目"
      />
      <div v-else class="project-asset-grid">
        <div v-for="i in 8" :key="i" class="project-asset">
          <div class="asset-thumb"></div>
          <div class="asset-name">素材 {{ i }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import EmptyState from '@/components/common/EmptyState.vue'

interface Project {
  id: string | number
  name: string
  client: string
  owner: string
  deadline: string
  status: string
  type: string
  assets: number
  references: string
  deliverables: string
  versions: string
}

defineProps<{
  project: Project | undefined
}>()
</script>

<style scoped>
.project-view {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
  background: #1e1e1e;
}

.project-hero {
  margin-bottom: 24px;
}

.project-hero-title {
  font-size: 28px;
  font-weight: 700;
  color: #f3f4f6;
  margin-bottom: 8px;
}

.project-hero-meta {
  font-size: 14px;
  color: #6b7280;
}

.project-panels {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
  margin-bottom: 24px;
}

.project-panel {
  background: #252526;
  border-radius: 12px;
  padding: 16px;
}

.panel-title {
  font-size: 13px;
  font-weight: 600;
  color: #9ca3af;
  margin-bottom: 12px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.panel-row {
  display: flex;
  justify-content: space-between;
  padding: 8px 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  font-size: 12px;
}

.panel-row:last-child {
  border-bottom: none;
}

.panel-row span:first-child {
  color: #6b7280;
}

.panel-row span:last-child {
  color: #d1d5db;
}

.project-section {
  margin-bottom: 24px;
}

.section-title {
  font-size: 14px;
  font-weight: 600;
  color: #9ca3af;
  margin-bottom: 12px;
}

.project-asset-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
  gap: 12px;
}

.project-asset {
  background: #252526;
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  transition: transform 0.15s;
}

.project-asset:hover {
  transform: translateY(-2px);
}

.asset-thumb {
  width: 100%;
  aspect-ratio: 1;
  background: #2a2a2a;
}

.asset-name {
  padding: 8px;
  font-size: 11px;
  color: #9ca3af;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
