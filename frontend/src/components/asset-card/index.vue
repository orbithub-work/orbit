<template>
  <component
    :is="cardComponent"
    :asset="asset"
    :status="status"
    :selected="selected"
    :thumbnail-url="thumbnailUrl"
    :show-type-badge="showTypeBadge"
    v-bind="typeSpecificProps"
    @click="(a: Asset, e: MouseEvent) => $emit('click', a, e)"
    @double-click="(a: Asset) => $emit('doubleClick', a)"
    @context-menu="(a: Asset, e: MouseEvent) => $emit('contextMenu', a, e)"
    @preview="$emit('preview')"
    @favorite="$emit('favorite')"
    v-on="typeSpecificEvents"
  />
</template>

<script setup lang="ts">
import { computed } from 'vue'
import AssetCard, { type Asset, type AssetType, type AssetStatus } from './AssetCard.vue'
import ImageCard, { type ImageAsset } from './ImageCard.vue'
import VideoCard, { type VideoAsset } from './VideoCard.vue'
import AudioCard, { type AudioAsset } from './AudioCard.vue'
import FontCard, { type FontAsset } from './FontCard.vue'
import ProjectCard, { type ProjectAsset } from './ProjectCard.vue'

export type { Asset, AssetType, AssetStatus, ImageAsset, VideoAsset, AudioAsset, FontAsset, ProjectAsset }

interface Props {
  asset: Asset
  status?: AssetStatus
  selected?: boolean
  thumbnailUrl?: string
  showTypeBadge?: boolean
  waveformData?: number[]
  isPlaying?: boolean
  previewText?: string
  loadedFontFamily?: string
}

const props = withDefaults(defineProps<Props>(), {
  status: 'ready',
  selected: false,
  showTypeBadge: true,
  isPlaying: false
})

defineEmits<{
  click: [asset: Asset, event: MouseEvent]
  doubleClick: [asset: Asset]
  contextMenu: [asset: Asset, event: MouseEvent]
  preview: []
  favorite: []
  play: []
  pause: []
  viewPrompt: []
  extractFrame: []
  details: []
  install: []
  previewAll: []
  open: []
  checkLinks: []
  viewMaterials: []
  more: []
}>()

function getFileExtension(filename: string): string {
  const parts = filename.split('.')
  return parts.length > 1 ? parts[parts.length - 1].toLowerCase() : ''
}

function getAssetType(filename: string): AssetType {
  const ext = getFileExtension(filename)
  
  if (['jpg', 'jpeg', 'png', 'gif', 'webp', 'bmp', 'svg', 'heic', 'raw', 'cr2', 'nef', 'arw', 'dng'].includes(ext)) return 'image'
  if (['mp4', 'mov', 'avi', 'mkv', 'webm', 'flv', 'wmv', 'm4v', 'prores'].includes(ext)) return 'video'
  if (['mp3', 'wav', 'flac', 'aac', 'ogg', 'm4a', 'wma', 'aiff'].includes(ext)) return 'audio'
  if (['ttf', 'otf', 'woff', 'woff2', 'eot'].includes(ext)) return 'font'
  if (['pdf', 'doc', 'docx', 'xls', 'xlsx', 'ppt', 'pptx', 'txt', 'md'].includes(ext)) return 'document'
  if (['draft', 'prproj', 'aep', 'fcpxml', 'drp', 'vpj', 'cap', 'veg'].includes(ext)) return 'project'
  if (['zip', 'rar', '7z', 'tar', 'gz', 'bz2'].includes(ext)) return 'archive'
  return 'other'
}

const assetType = computed<AssetType>(() => {
  return getAssetType(props.asset.name)
})

const cardComponent = computed(() => {
  const componentMap: Record<AssetType, any> = {
    image: ImageCard,
    video: VideoCard,
    audio: AudioCard,
    font: FontCard,
    project: ProjectCard,
    document: AssetCard,
    archive: AssetCard,
    other: AssetCard
  }
  
  return componentMap[assetType.value]
})

const typeSpecificProps = computed(() => {
  const base: Record<string, any> = {}
  
  if (assetType.value === 'audio') {
    base.waveformData = props.waveformData
    base.isPlaying = props.isPlaying
  }
  
  if (assetType.value === 'font') {
    base.previewText = props.previewText
    base.loadedFontFamily = props.loadedFontFamily
  }
  
  return base
})

const typeSpecificEvents = computed(() => {
  return {}
})
</script>
