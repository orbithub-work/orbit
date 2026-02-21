<template>
  <AssetCard
    :asset="asset"
    :status="status"
    :selected="selected"
    :show-type-badge="showTypeBadge"
    class="audio-card"
    @click="(a, e) => $emit('click', a, e)"
    @double-click="(a) => $emit('doubleClick', a)"
    @context-menu="(a, e) => $emit('contextMenu', a, e)"
    @preview="$emit('preview')"
    @favorite="$emit('favorite')"
  >
    <template #thumbnail>
      <div class="audio-thumbnail" :class="{ 'is-playing': isPlaying }">
        <div class="audio-visual">
          <div v-if="waveformData && waveformData.length > 0" class="waveform">
            <div
              v-for="(bar, index) in waveformData"
              :key="index"
              class="waveform-bar"
              :style="{ height: `${bar}%` }"
            />
          </div>
          <div v-else class="audio-icon-wrapper">
            <svg class="music-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M9 18V5l12-2v13" />
              <circle cx="6" cy="18" r="3" />
              <circle cx="18" cy="16" r="3" />
            </svg>
          </div>
        </div>
        
        <div class="audio-info-overlay">
          <span class="duration">{{ formatDuration(asset.duration) }}</span>
          <span v-if="asset.bitrate" class="bitrate">{{ asset.bitrate }}kbps</span>
        </div>

        <div v-if="isPlaying" class="playing-indicator">
          <div class="sound-wave">
            <span /><span /><span /><span /><span />
          </div>
        </div>

        <div v-if="asset.artist" class="artist-badge">
          {{ truncateArtist(asset.artist) }}
        </div>
      </div>
    </template>

    <template #meta>
      <span>{{ formatFileSize(asset.size) }}</span>
      <span>{{ formatDuration(asset.duration) }}</span>
      <span v-if="asset.artist" class="artist">{{ truncateArtist(asset.artist) }}</span>
    </template>

    <template #hover-actions>
      <button class="hover-btn" title="播放" @click.stop="togglePlay">
        <svg v-if="!isPlaying" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polygon points="5 3 19 12 5 21 5 3" />
        </svg>
        <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <rect x="6" y="4" width="4" height="16" />
          <rect x="14" y="4" width="4" height="16" />
        </svg>
      </button>
      <button class="hover-btn" title="查看详情" @click.stop="$emit('details')">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="10" />
          <line x1="12" y1="16" x2="12" y2="12" />
          <line x1="12" y1="8" x2="12.01" y2="8" />
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
import { ref } from 'vue'
import AssetCard, { type Asset, type AssetStatus } from './AssetCard.vue'

export interface AudioAsset extends Asset {
  duration?: number
  bitrate?: number
  artist?: string
  album?: string
  title?: string
  sample_rate?: number
  channels?: number
}

interface Props {
  asset: AudioAsset
  status?: AssetStatus
  selected?: boolean
  showTypeBadge?: boolean
  waveformData?: number[]
  isPlaying?: boolean
}

withDefaults(defineProps<Props>(), {
  status: 'ready',
  selected: false,
  showTypeBadge: true,
  isPlaying: false
})

const emit = defineEmits<{
  click: [asset: AudioAsset, event: MouseEvent]
  doubleClick: [asset: AudioAsset]
  contextMenu: [asset: AudioAsset, event: MouseEvent]
  preview: []
  play: []
  pause: []
  details: []
  favorite: []
}>()

function togglePlay() {
  emit('play')
}

function formatDuration(seconds?: number): string {
  if (!seconds) return '0:00'
  const mins = Math.floor(seconds / 60)
  const secs = Math.floor(seconds % 60)
  return `${mins}:${secs.toString().padStart(2, '0')}`
}

function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(1))} ${sizes[i]}`
}

function truncateArtist(artist: string): string {
  return artist.length > 15 ? artist.slice(0, 15) + '...' : artist
}
</script>

<style scoped>
.audio-card {
  --card-radius: 8px;
}

.audio-thumbnail {
  position: relative;
  width: 100%;
  aspect-ratio: 1;
  overflow: hidden;
  background: linear-gradient(135deg, #1a1a2e, #16213e);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.audio-visual {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
}

.audio-icon-wrapper {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea, #764ba2);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 20px rgba(102, 126, 234, 0.4);
  transition: transform 0.3s;
}

.audio-card:hover .audio-icon-wrapper {
  transform: scale(1.1);
}

.music-icon {
  width: 32px;
  height: 32px;
  color: #fff;
}

.waveform {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 2px;
  height: 60%;
  width: 80%;
}

.waveform-bar {
  width: 3px;
  min-height: 4px;
  background: linear-gradient(180deg, #667eea, #764ba2);
  border-radius: 2px;
  transition: height 0.1s;
}

.audio-thumbnail.is-playing .waveform-bar {
  animation: waveform 0.5s ease-in-out infinite alternate;
}

.audio-thumbnail.is-playing .waveform-bar:nth-child(odd) {
  animation-delay: 0.1s;
}

@keyframes waveform {
  from { transform: scaleY(0.5); }
  to { transform: scaleY(1); }
}

.audio-info-overlay {
  position: absolute;
  bottom: 8px;
  left: 8px;
  display: flex;
  gap: 8px;
  font-size: 11px;
  color: rgba(255, 255, 255, 0.7);
}

.duration {
  background: rgba(0, 0, 0, 0.6);
  padding: 2px 6px;
  border-radius: 4px;
}

.bitrate {
  background: rgba(102, 126, 234, 0.6);
  padding: 2px 6px;
  border-radius: 4px;
}

.playing-indicator {
  position: absolute;
  top: 8px;
  right: 8px;
}

.sound-wave {
  display: flex;
  align-items: center;
  gap: 3px;
  height: 24px;
  padding: 4px 8px;
  background: rgba(59, 130, 246, 0.8);
  border-radius: 12px;
}

.sound-wave span {
  width: 4px;
  background: #fff;
  border-radius: 2px;
  animation: sound-wave 0.5s ease-in-out infinite alternate;
}

.sound-wave span:nth-child(1) { animation-delay: 0s; height: 8px; }
.sound-wave span:nth-child(2) { animation-delay: 0.1s; height: 16px; }
.sound-wave span:nth-child(3) { animation-delay: 0.2s; height: 24px; }
.sound-wave span:nth-child(4) { animation-delay: 0.3s; height: 16px; }
.sound-wave span:nth-child(5) { animation-delay: 0.4s; height: 8px; }

@keyframes sound-wave {
  from { transform: scaleY(0.5); }
  to { transform: scaleY(1); }
}

.artist-badge {
  position: absolute;
  top: 8px;
  left: 8px;
  padding: 2px 8px;
  background: rgba(0, 0, 0, 0.6);
  border-radius: 4px;
  font-size: 10px;
  color: rgba(255, 255, 255, 0.8);
  max-width: 100px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.artist {
  color: var(--color-text-tertiary, #6b7280);
}
</style>
