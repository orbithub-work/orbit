<template>
  <AssetCard
    :asset="asset"
    :status="status"
    :selected="selected"
    :show-type-badge="false"
    class="project-card"
    @click="(a, e) => $emit('click', a, e)"
    @double-click="(a) => $emit('doubleClick', a)"
    @context-menu="(a, e) => $emit('contextMenu', a, e)"
    @preview="$emit('preview')"
    @favorite="$emit('favorite')"
  >
    <template #thumbnail>
      <div class="project-thumbnail">
        <div class="project-icon" :class="`project-icon--${projectTypeKey}`">
          <component :is="projectIcon" />
        </div>
        
        <div class="project-meta-overlay">
          <span class="project-type">{{ projectTypeLabel }}</span>
          <span v-if="asset.material_count" class="material-count">
            {{ asset.material_count }} 个素材
          </span>
        </div>
        
        <div v-if="asset.broken_links && asset.broken_links > 0" class="broken-links-warning">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z" />
            <line x1="12" y1="9" x2="12" y2="13" />
            <line x1="12" y1="17" x2="12.01" y2="17" />
          </svg>
          <span>{{ asset.broken_links }} 个断链</span>
        </div>

        <div v-if="asset.duration" class="project-duration">
          {{ formatDuration(asset.duration) }}
        </div>
      </div>
    </template>

    <template #meta>
      <span>{{ projectTypeLabel }}</span>
      <span v-if="asset.modified_at">{{ formatDate(asset.modified_at) }}</span>
      <span>{{ formatFileSize(asset.size) }}</span>
    </template>

    <template #hover-actions>
      <button class="hover-btn" title="打开工程" @click.stop="$emit('open')">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6" />
          <polyline points="15 3 21 3 21 9" />
          <line x1="10" y1="14" x2="21" y2="3" />
        </svg>
      </button>
      <button class="hover-btn" title="检查断链" @click.stop="$emit('checkLinks')">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71" />
          <path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71" />
        </svg>
      </button>
      <button class="hover-btn" title="查看素材" @click.stop="$emit('viewMaterials')">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polygon points="12 2 2 7 12 12 22 7 12 2" />
          <polyline points="2 17 12 22 22 17" />
          <polyline points="2 12 12 17 22 12" />
        </svg>
      </button>
      <button class="hover-btn" title="收藏" @click.stop="$emit('favorite')">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2" />
        </svg>
      </button>
    </template>
  </AssetCard>
</template>

<script setup lang="ts">
import { computed, defineComponent, h } from 'vue'
import AssetCard, { type Asset, type AssetStatus } from './AssetCard.vue'

export interface ProjectAsset extends Asset {
  material_count?: number
  broken_links?: number
  duration?: number
  modified_at?: string
  software_version?: string
}

interface Props {
  asset: ProjectAsset
  status?: AssetStatus
  selected?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  status: 'ready',
  selected: false
})

defineEmits<{
  click: [asset: ProjectAsset, event: MouseEvent]
  doubleClick: [asset: ProjectAsset]
  contextMenu: [asset: ProjectAsset, event: MouseEvent]
  preview: []
  open: []
  checkLinks: []
  viewMaterials: []
  favorite: []
}>()

interface ProjectTypeConfig {
  key: string
  label: string
  icon: any
  color: string
}

const projectTypeMap: Record<string, ProjectTypeConfig> = {
  'draft': {
    key: 'jianying',
    label: '剪映',
    icon: defineComponent({
      render: () => h('svg', { viewBox: '0 0 24 24', fill: 'currentColor' }, [
        h('path', { d: 'M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z' })
      ])
    }),
    color: '#000'
  },
  'prproj': {
    key: 'premiere',
    label: 'Premiere',
    icon: defineComponent({
      render: () => h('svg', { viewBox: '0 0 24 24', fill: 'currentColor' }, [
        h('path', { d: 'M4 4h16v16H4V4zm2 2v12h12V6H6zm3 2h2v8H9V8zm4 0h2c1.1 0 2 .9 2 2v4c0 1.1-.9 2-2 2h-2V8zm2 2v4h0V10z' })
      ])
    }),
    color: '#9999ff'
  },
  'aep': {
    key: 'aftereffects',
    label: 'After Effects',
    icon: defineComponent({
      render: () => h('svg', { viewBox: '0 0 24 24', fill: 'currentColor' }, [
        h('path', { d: 'M4 4h16v16H4V4zm2 2v12h12V6H6zm3 2h2l2 8h-2l-.5-2h-1.5l-.5 2H8l2-8zm1 2l-.5 2h1l-.5-2z' })
      ])
    }),
    color: '#9999ff'
  },
  'fcpxml': {
    key: 'fcpx',
    label: 'Final Cut',
    icon: defineComponent({
      render: () => h('svg', { viewBox: '0 0 24 24', fill: 'currentColor' }, [
        h('path', { d: 'M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5' })
      ])
    }),
    color: '#fff'
  },
  'drp': {
    key: 'davinci',
    label: 'DaVinci',
    icon: defineComponent({
      render: () => h('svg', { viewBox: '0 0 24 24', fill: 'currentColor' }, [
        h('path', { d: 'M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 18c-4.41 0-8-3.59-8-8s3.59-8 8-8 8 3.59 8 8-3.59 8-8 8zm-1-13h2v6h-2zm0 8h2v2h-2z' })
      ])
    }),
    color: '#ff6b35'
  },
  'cap': {
    key: 'capcut',
    label: 'CapCut',
    icon: defineComponent({
      render: () => h('svg', { viewBox: '0 0 24 24', fill: 'currentColor' }, [
        h('path', { d: 'M8 5v14l11-7z' })
      ])
    }),
    color: '#00f2ea'
  },
  'veg': {
    key: 'vegas',
    label: 'Vegas',
    icon: defineComponent({
      render: () => h('svg', { viewBox: '0 0 24 24', fill: 'currentColor' }, [
        h('path', { d: 'M4 4h16v16H4V4zm2 2v12h12V6H6z' })
      ])
    }),
    color: '#1a1a1a'
  }
}

const projectTypeKey = computed(() => {
  const ext = props.asset.name.split('.').pop()?.toLowerCase() || ''
  return projectTypeMap[ext]?.key || 'unknown'
})

const projectTypeLabel = computed(() => {
  const ext = props.asset.name.split('.').pop()?.toLowerCase() || ''
  return projectTypeMap[ext]?.label || '工程文件'
})

const projectIcon = computed(() => {
  const ext = props.asset.name.split('.').pop()?.toLowerCase() || ''
  return projectTypeMap[ext]?.icon || defineComponent({
    render: () => h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '2' }, [
      h('path', { d: 'M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z' })
    ])
  })
})

function formatDuration(seconds?: number): string {
  if (!seconds) return ''
  const mins = Math.floor(seconds / 60)
  const secs = Math.floor(seconds % 60)
  if (mins >= 60) {
    const hours = Math.floor(mins / 60)
    const remainMins = mins % 60
    return `${hours}:${remainMins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
  }
  return `${mins}:${secs.toString().padStart(2, '0')}`
}

function formatDate(dateStr?: string): string {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))
  
  if (days === 0) return '今天'
  if (days === 1) return '昨天'
  if (days < 7) return `${days}天前`
  if (days < 30) return `${Math.floor(days / 7)}周前`
  if (days < 365) return `${Math.floor(days / 30)}个月前`
  return `${Math.floor(days / 365)}年前`
}

function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(1))} ${sizes[i]}`
}
</script>

<style scoped>
.project-card {
  --card-radius: 8px;
}

.project-thumbnail {
  position: relative;
  aspect-ratio: 16/9;
  background: linear-gradient(135deg, #1e1e2e, #2d2d3f);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  overflow: hidden;
}

.project-icon {
  width: 48px;
  height: 48px;
  padding: 10px;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: transform 0.2s;
}

.project-card:hover .project-icon {
  transform: scale(1.1);
}

.project-icon svg {
  width: 28px;
  height: 28px;
}

.project-icon--jianying {
  background: linear-gradient(135deg, #000, #333);
  color: #fff;
}

.project-icon--premiere {
  background: linear-gradient(135deg, #9999ff, #00005b);
  color: #fff;
}

.project-icon--aftereffects {
  background: linear-gradient(135deg, #9999ff, #00005b);
  color: #fff;
}

.project-icon--fcpx {
  background: linear-gradient(135deg, #fff, #ccc);
  color: #000;
}

.project-icon--davinci {
  background: linear-gradient(135deg, #ff6b35, #f7931e);
  color: #fff;
}

.project-icon--capcut {
  background: linear-gradient(135deg, #00f2ea, #ff0050);
  color: #fff;
}

.project-icon--vegas {
  background: linear-gradient(135deg, #1a1a1a, #4a4a4a);
  color: #fff;
}

.project-meta-overlay {
  position: absolute;
  bottom: 8px;
  left: 8px;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.project-type {
  font-size: 11px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.9);
}

.material-count {
  font-size: 10px;
  color: rgba(255, 255, 255, 0.6);
}

.broken-links-warning {
  position: absolute;
  top: 8px;
  right: 8px;
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
  background: rgba(239, 68, 68, 0.9);
  border-radius: 4px;
  font-size: 11px;
  color: #fff;
}

.broken-links-warning svg {
  width: 12px;
  height: 12px;
}

.project-duration {
  position: absolute;
  bottom: 8px;
  right: 8px;
  padding: 2px 6px;
  background: rgba(0, 0, 0, 0.7);
  border-radius: 4px;
  font-size: 11px;
  color: #fff;
}
</style>
