<template>
  <AssetCard
    :asset="asset"
    :status="status"
    :selected="selected"
    :show-type-badge="false"
    class="font-card"
    @click="(a, e) => $emit('click', a, e)"
    @double-click="(a) => $emit('doubleClick', a)"
    @context-menu="(a, e) => $emit('contextMenu', a, e)"
    @preview="$emit('preview')"
    @favorite="$emit('favorite')"
  >
    <template #thumbnail>
      <div class="font-preview" :style="fontStyle">
        <div class="font-sample">
          {{ previewText }}
        </div>
        <div class="font-name">
          {{ asset.font_name || asset.name }}
        </div>
        <div class="font-format">
          {{ fontFormat }}
        </div>
      </div>
    </template>

    <template #meta>
      <span>{{ asset.font_style || 'Regular' }}</span>
      <span v-if="asset.font_weight">{{ asset.font_weight }}</span>
      <span>{{ formatFileSize(asset.size) }}</span>
    </template>

    <template #extra>
      <div class="font-actions">
        <button class="font-btn font-btn--primary" @click.stop="$emit('install')">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
            <polyline points="7 10 12 15 17 10" />
            <line x1="12" y1="15" x2="12" y2="3" />
          </svg>
          安装
        </button>
        <button class="font-btn" @click.stop="$emit('previewAll')">
          预览
        </button>
      </div>
    </template>

    <template #hover-actions>
      <button class="hover-btn" title="预览全部字符" @click.stop="$emit('previewAll')">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polyline points="4 7 4 4 20 4 20 7" />
          <line x1="9" y1="20" x2="15" y2="20" />
          <line x1="12" y1="4" x2="12" y2="20" />
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
import { computed } from 'vue'
import AssetCard, { type Asset, type AssetStatus } from './AssetCard.vue'

export interface FontAsset extends Asset {
  font_name?: string
  font_style?: string
  font_weight?: number
  font_family?: string
}

interface Props {
  asset: FontAsset
  status?: AssetStatus
  selected?: boolean
  previewText?: string
  loadedFontFamily?: string
}

const props = withDefaults(defineProps<Props>(), {
  status: 'ready',
  selected: false,
  previewText: 'Aa永字八法'
})

defineEmits<{
  click: [asset: FontAsset, event: MouseEvent]
  doubleClick: [asset: FontAsset]
  contextMenu: [asset: FontAsset, event: MouseEvent]
  preview: []
  previewAll: []
  install: []
  favorite: []
}>()

const fontFormat = computed(() => {
  const ext = props.asset.name.split('.').pop()?.toLowerCase()
  const formatMap: Record<string, string> = {
    'ttf': 'TrueType',
    'otf': 'OpenType',
    'woff': 'WOFF',
    'woff2': 'WOFF2',
    'eot': 'EOT'
  }
  return formatMap[ext || ''] || ext?.toUpperCase() || ''
})

const fontStyle = computed(() => {
  if (props.loadedFontFamily) {
    return { fontFamily: props.loadedFontFamily }
  }
  return {}
})

function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(1))} ${sizes[i]}`
}
</script>

<style scoped>
.font-card {
  --card-radius: 8px;
}

.font-preview {
  position: relative;
  width: 100%;
  aspect-ratio: 1;
  background: linear-gradient(180deg, #fafafa, #f0f0f0);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 16px;
  transition: background 0.2s;
}

.font-card:hover .font-preview {
  background: linear-gradient(180deg, #fff, #f5f5f5);
}

.font-sample {
  font-size: 32px;
  color: #1a1a1a;
  line-height: 1.2;
  text-align: center;
  word-break: break-all;
  transition: transform 0.2s;
}

.font-card:hover .font-sample {
  transform: scale(1.05);
}

.font-name {
  margin-top: 12px;
  font-size: 12px;
  color: #666;
  font-family: system-ui, -apple-system, sans-serif;
  text-align: center;
  max-width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.font-format {
  position: absolute;
  bottom: 8px;
  right: 8px;
  padding: 2px 6px;
  background: rgba(0, 0, 0, 0.6);
  border-radius: 4px;
  font-size: 10px;
  font-weight: 500;
  color: #fff;
  font-family: system-ui;
}

.font-actions {
  display: flex;
  gap: 8px;
  margin-top: 4px;
}

.font-btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  padding: 6px 12px;
  font-size: 12px;
  font-weight: 500;
  background: var(--color-surface-elevated, #3a3a3a);
  border: none;
  border-radius: 6px;
  color: var(--color-text-primary, #e5e7eb);
  cursor: pointer;
  transition: all 0.2s;
}

.font-btn:hover {
  background: var(--color-hover, #4a4a4a);
}

.font-btn--primary {
  background: var(--color-primary, #3b82f6);
  color: #fff;
}

.font-btn--primary:hover {
  background: var(--color-primary-hover, #2563eb);
}

.font-btn svg {
  width: 14px;
  height: 14px;
}
</style>
