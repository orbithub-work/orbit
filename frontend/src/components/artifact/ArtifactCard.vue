<template>
  <div 
    class="artifact-card"
    :class="{ selected: isSelected }"
    @click="$emit('select', artifact)"
    @dblclick="$emit('open', artifact)"
  >
    <div class="card-thumbnail">
      <img v-if="artifact.thumbnail" :src="artifact.thumbnail" :alt="artifact.name" />
      <div v-else class="thumbnail-placeholder">
        <Icon name="file" size="lg" />
      </div>
      <div class="card-status" :class="`status--${artifact.status}`">
        {{ statusText }}
      </div>
    </div>
    <div class="card-content">
      <div class="card-title">{{ artifact.name }}</div>
      <div class="card-meta">
        <span>{{ artifact.project_name }}</span>
        <span>{{ formatDate(artifact.created_at) }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import Icon from '@/components/common/Icon.vue'

interface Artifact {
  id: string
  name: string
  thumbnail?: string
  status: string
  project_name: string
  created_at: string
}

const props = defineProps<{
  artifact: Artifact
  isSelected: boolean
}>()

defineEmits<{
  select: [artifact: Artifact]
  open: [artifact: Artifact]
}>()

const statusText = computed(() => {
  const map: Record<string, string> = {
    draft: '草稿',
    review: '审核中',
    published: '已发布',
    archived: '已归档'
  }
  return map[props.artifact.status] || props.artifact.status
})

function formatDate(date: string) {
  return new Date(date).toLocaleDateString('zh-CN')
}
</script>

<script lang="ts">
import { computed } from 'vue'
export default { name: 'ArtifactCard' }
</script>

<style scoped>
.artifact-card {
  border-radius: 8px;
  overflow: hidden;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.06);
  cursor: pointer;
  transition: all 0.2s;
}

.artifact-card:hover {
  background: rgba(255, 255, 255, 0.06);
  transform: translateY(-2px);
}

.artifact-card.selected {
  outline: 2px solid #3b82f6;
  background: rgba(59, 130, 246, 0.1);
}

.card-thumbnail {
  position: relative;
  aspect-ratio: 16/9;
  background: rgba(0, 0, 0, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
}

.card-thumbnail img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.thumbnail-placeholder {
  color: #6b7280;
}

.card-status {
  position: absolute;
  top: 8px;
  right: 8px;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 500;
}

.status--draft { background: rgba(107, 114, 128, 0.9); color: white; }
.status--review { background: rgba(251, 191, 36, 0.9); color: white; }
.status--published { background: rgba(34, 197, 94, 0.9); color: white; }
.status--archived { background: rgba(156, 163, 175, 0.9); color: white; }

.card-content {
  padding: 12px;
}

.card-title {
  font-size: 14px;
  font-weight: 500;
  color: #e5e7eb;
  margin-bottom: 6px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.card-meta {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: #9ca3af;
}
</style>
