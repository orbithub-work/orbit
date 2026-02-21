<!-- 快速预览组件 - 按Space键触发 -->
<template>
  <Teleport to="body">
    <Transition name="preview-fade">
      <div
        v-if="visible"
        class="quick-preview"
        @click.self="close"
        @keydown.esc="close"
        @keydown.left="navigatePrev"
        @keydown.right="navigateNext"
        @keydown.space.prevent="close"
      >
        <!-- 背景遮罩 -->
        <div class="quick-preview__backdrop" />

        <!-- 主内容区 -->
        <div class="quick-preview__content">
          <!-- 关闭按钮 -->
          <button
            class="quick-preview__close"
            @click="close"
          >
            <svg class="icon">
              <use href="#icon-close" />
            </svg>
          </button>

          <!-- 导航按钮 -->
          <button
            v-if="hasPrev"
            class="quick-preview__nav quick-preview__nav--prev"
            @click="navigatePrev"
          >
            <svg class="icon">
              <use href="#icon-chevron-left" />
            </svg>
          </button>

          <button
            v-if="hasNext"
            class="quick-preview__nav quick-preview__nav--next"
            @click="navigateNext"
          >
            <svg class="icon">
              <use href="#icon-chevron-right" />
            </svg>
          </button>

          <!-- 媒体预览区 -->
          <div class="quick-preview__media">
            <!-- 图片预览 -->
            <img
              v-if="currentAsset && isImage(currentAsset)"
              :src="getPreviewUrl(currentAsset)"
              :alt="currentAsset.name"
              class="preview-image"
              @load="handleMediaLoad"
              @error="handleMediaError"
            />

            <!-- 视频预览 -->
            <video
              v-else-if="currentAsset && isVideo(currentAsset)"
              :src="getPreviewUrl(currentAsset)"
              class="preview-video"
              controls
              autoplay
              @loadedmetadata="handleMediaLoad"
              @error="handleMediaError"
            />

            <!-- 音频预览 -->
            <div
              v-else-if="currentAsset && isAudio(currentAsset)"
              class="preview-audio"
            >
              <div class="audio-icon">
                <svg class="icon">
                  <use href="#icon-music" />
                </svg>
              </div>
              <audio
                :src="getPreviewUrl(currentAsset)"
                controls
                autoplay
                @loadedmetadata="handleMediaLoad"
                @error="handleMediaError"
              />
            </div>

            <!-- PDF预览 -->
            <iframe
              v-else-if="currentAsset && isPDF(currentAsset)"
              :src="getPreviewUrl(currentAsset)"
              class="preview-pdf"
              @load="handleMediaLoad"
            />

            <!-- 不支持预览 -->
            <div
              v-else
              class="preview-unsupported"
            >
              <svg class="icon icon--large">
                <use href="#icon-file" />
              </svg>
              <p>暂不支持预览此文件类型</p>
              <button
                class="btn btn--primary"
                @click="openInDefault"
              >
                使用默认应用打开
              </button>
            </div>

            <!-- 加载状态 -->
            <div
              v-if="loading"
              class="preview-loading"
            >
              <div class="loading-spinner" />
              <span>加载中...</span>
            </div>

            <!-- 错误状态 -->
            <div
              v-if="error"
              class="preview-error"
            >
              <svg class="icon icon--large">
                <use href="#icon-alert-circle" />
              </svg>
              <p>{{ error }}</p>
            </div>
          </div>

          <!-- 信息侧边栏 -->
          <div class="quick-preview__sidebar">
            <div
              v-if="currentAsset"
              class="preview-info"
            >
              <!-- 文件名 -->
              <h3 class="preview-info__title">
                {{ currentAsset.name }}
              </h3>

              <!-- 元数据 -->
              <div class="preview-info__meta">
                <div class="meta-item">
                  <span class="meta-label">大小</span>
                  <span class="meta-value">{{ formatFileSize(currentAsset.size) }}</span>
                </div>

                <div
                  v-if="currentAsset.width && currentAsset.height"
                  class="meta-item"
                >
                  <span class="meta-label">尺寸</span>
                  <span class="meta-value">{{ currentAsset.width }} × {{ currentAsset.height }}</span>
                </div>

                <div
                  v-if="currentAsset.modified_at"
                  class="meta-item"
                >
                  <span class="meta-label">修改时间</span>
                  <span class="meta-value">{{ formatDate(currentAsset.modified_at) }}</span>
                </div>

                <div class="meta-item">
                  <span class="meta-label">路径</span>
                  <span class="meta-value meta-value--path">{{ currentAsset.path }}</span>
                </div>
              </div>

              <!-- 标签 -->
              <div
                v-if="currentAsset.tags && currentAsset.tags.length > 0"
                class="preview-info__tags"
              >
                <h4 class="section-title">
                  标签
                </h4>
                <div class="tags-list">
                  <span
                    v-for="tag in currentAsset.tags"
                    :key="tag.id"
                    class="tag"
                    :style="{ '--tag-color': tag.color }"
                  >
                    {{ tag.name }}
                  </span>
                </div>
              </div>

              <!-- 操作按钮 -->
              <div class="preview-info__actions">
                <button
                  class="btn btn--secondary btn--block"
                  @click="openInFolder"
                >
                  <svg class="icon">
                    <use href="#icon-folder-open" />
                  </svg>
                  在文件夹中显示
                </button>

                <button
                  class="btn btn--secondary btn--block"
                  @click="copyPath"
                >
                  <svg class="icon">
                    <use href="#icon-copy" />
                  </svg>
                  复制路径
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- 底部工具栏 -->
        <div class="quick-preview__toolbar">
          <div class="toolbar-left">
            <span class="toolbar-text">
              {{ currentIndex + 1 }} / {{ total }}
            </span>
          </div>

          <div class="toolbar-right">
            <span class="toolbar-hint">
              <kbd>←</kbd> <kbd>→</kbd> 切换 | <kbd>Space</kbd> 或 <kbd>Esc</kbd> 关闭
            </span>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'

interface Asset {
  id: string
  name: string
  path: string
  size: number
  file_type?: string
  width?: number
  height?: number
  modified_at?: number
  tags?: Array<{ id: string; name: string; color: string }>
}

interface Props {
  visible: boolean
  assets: Asset[]
  currentIndex: number
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:visible': [value: boolean]
  'update:currentIndex': [value: number]
  open: [asset: Asset]
  openInFolder: [asset: Asset]
}>()

const loading = ref(false)
const error = ref<string | null>(null)

const currentAsset = computed(() => props.assets[props.currentIndex])
const total = computed(() => props.assets.length)
const hasPrev = computed(() => props.currentIndex > 0)
const hasNext = computed(() => props.currentIndex < props.assets.length - 1)

// 关闭预览
function close() {
  emit('update:visible', false)
}

// 上一个
function navigatePrev() {
  if (hasPrev.value) {
    emit('update:currentIndex', props.currentIndex - 1)
    resetLoadingState()
  }
}

// 下一个
function navigateNext() {
  if (hasNext.value) {
    emit('update:currentIndex', props.currentIndex + 1)
    resetLoadingState()
  }
}

// 重置加载状态
function resetLoadingState() {
  loading.value = true
  error.value = null
}

// 媒体加载完成
function handleMediaLoad() {
  loading.value = false
  error.value = null
}

// 媒体加载错误
function handleMediaError() {
  loading.value = false
  error.value = '无法加载文件预览'
}

// 判断文件类型
function isImage(asset: Asset): boolean {
  const ext = getExtension(asset.name).toLowerCase()
  return ['jpg', 'jpeg', 'png', 'gif', 'webp', 'bmp', 'svg'].includes(ext)
}

function isVideo(asset: Asset): boolean {
  const ext = getExtension(asset.name).toLowerCase()
  return ['mp4', 'mov', 'avi', 'mkv', 'webm', 'flv'].includes(ext)
}

function isAudio(asset: Asset): boolean {
  const ext = getExtension(asset.name).toLowerCase()
  return ['mp3', 'wav', 'flac', 'aac', 'ogg', 'm4a'].includes(ext)
}

function isPDF(asset: Asset): boolean {
  return getExtension(asset.name).toLowerCase() === 'pdf'
}

function getExtension(filename: string): string {
  const parts = filename.split('.')
  return parts.length > 1 ? parts[parts.length - 1] : ''
}

// 获取预览URL
function getPreviewUrl(asset: Asset): string {
  // 这里应该返回实际的预览URL
  // 可能需要通过API获取，或者直接使用文件路径
  return `/api/assets/${asset.id}/preview`
}

// 在默认应用中打开
function openInDefault() {
  if (currentAsset.value) {
    emit('open', currentAsset.value)
  }
}

// 在文件夹中显示
function openInFolder() {
  if (currentAsset.value) {
    emit('openInFolder', currentAsset.value)
  }
}

// 复制路径
async function copyPath() {
  if (!currentAsset.value) return

  try {
    await navigator.clipboard.writeText(currentAsset.value.path)
    // 这里应该显示一个成功提示
  } catch (err) {
    console.error('Failed to copy path:', err)
  }
}

// 格式化文件大小
function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(1))} ${sizes[i]}`
}

// 格式化日期
function formatDate(timestamp: number): string {
  return new Date(timestamp * 1000).toLocaleString('zh-CN')
}

// 监听可见性变化
watch(() => props.visible, (newVal) => {
  if (newVal) {
    resetLoadingState()
    document.body.style.overflow = 'hidden'
  } else {
    document.body.style.overflow = ''
  }
})

// 全局键盘事件
function handleGlobalKeyDown(event: KeyboardEvent) {
  if (!props.visible) return

  if (event.key === 'ArrowLeft') {
    event.preventDefault()
    navigatePrev()
  } else if (event.key === 'ArrowRight') {
    event.preventDefault()
    navigateNext()
  } else if (event.key === 'Escape' || event.key === ' ') {
    event.preventDefault()
    close()
  }
}

onMounted(() => {
  document.addEventListener('keydown', handleGlobalKeyDown)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleGlobalKeyDown)
  document.body.style.overflow = ''
})
</script>

<style scoped>
.quick-preview {
  position: fixed;
  inset: 0;
  z-index: 9999;
  display: flex;
  flex-direction: column;
}

.quick-preview__backdrop {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.95);
  backdrop-filter: blur(8px);
}

.quick-preview__content {
  position: relative;
  flex: 1;
  display: grid;
  grid-template-columns: 1fr 320px;
  gap: 0;
  overflow: hidden;
}

.quick-preview__close {
  position: absolute;
  top: 20px;
  right: 340px;
  z-index: 10;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.5);
  border: none;
  border-radius: 50%;
  color: #fff;
  cursor: pointer;
  transition: all 0.2s;
}

.quick-preview__close:hover {
  background: rgba(0, 0, 0, 0.8);
  transform: scale(1.1);
}

.quick-preview__nav {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  z-index: 10;
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.5);
  border: none;
  border-radius: 50%;
  color: #fff;
  cursor: pointer;
  transition: all 0.2s;
}

.quick-preview__nav:hover {
  background: rgba(0, 0, 0, 0.8);
  transform: translateY(-50%) scale(1.1);
}

.quick-preview__nav--prev {
  left: 20px;
}

.quick-preview__nav--next {
  right: 340px;
}

.quick-preview__media {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 60px;
  overflow: hidden;
}

.preview-image,
.preview-video {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
  border-radius: 8px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.5);
}

.preview-audio {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 24px;
}

.audio-icon {
  width: 120px;
  height: 120px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 50%;
  color: #fff;
}

.audio-icon .icon {
  width: 64px;
  height: 64px;
}

.preview-audio audio {
  width: 400px;
}

.preview-pdf {
  width: 100%;
  height: 100%;
  border: none;
  border-radius: 8px;
}

.preview-unsupported {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  color: #9ca3af;
}

.preview-unsupported .icon--large {
  width: 80px;
  height: 80px;
}

.preview-loading,
.preview-error {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 16px;
  color: #fff;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 3px solid rgba(255, 255, 255, 0.2);
  border-top-color: #3b82f6;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.quick-preview__sidebar {
  background: rgba(20, 20, 20, 0.95);
  border-left: 1px solid rgba(255, 255, 255, 0.1);
  overflow-y: auto;
}

.preview-info {
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.preview-info__title {
  font-size: 16px;
  font-weight: 600;
  color: #fff;
  margin: 0;
  word-break: break-word;
}

.preview-info__meta {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.meta-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.meta-label {
  font-size: 12px;
  color: #9ca3af;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.meta-value {
  font-size: 14px;
  color: #e5e7eb;
}

.meta-value--path {
  font-family: monospace;
  font-size: 12px;
  word-break: break-all;
  opacity: 0.8;
}

.section-title {
  font-size: 12px;
  font-weight: 600;
  color: #9ca3af;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin: 0 0 12px 0;
}

.tags-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.tag {
  padding: 4px 12px;
  font-size: 12px;
  background: var(--tag-color, #6b7280);
  color: #fff;
  border-radius: 12px;
}

.preview-info__actions {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.btn {
  padding: 10px 16px;
  font-size: 14px;
  font-weight: 500;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.btn--primary {
  background: #3b82f6;
  color: #fff;
}

.btn--primary:hover {
  background: #2563eb;
}

.btn--secondary {
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
}

.btn--secondary:hover {
  background: rgba(255, 255, 255, 0.15);
}

.btn--block {
  width: 100%;
}

.quick-preview__toolbar {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 24px;
  background: rgba(0, 0, 0, 0.5);
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.toolbar-text,
.toolbar-hint {
  font-size: 13px;
  color: #9ca3af;
}

.toolbar-hint {
  display: flex;
  align-items: center;
  gap: 8px;
}

kbd {
  padding: 2px 6px;
  font-size: 11px;
  font-family: monospace;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 4px;
  color: #fff;
}

/* 过渡动画 */
.preview-fade-enter-active,
.preview-fade-leave-active {
  transition: opacity 0.2s;
}

.preview-fade-enter-from,
.preview-fade-leave-to {
  opacity: 0;
}

/* 滚动条 */
.quick-preview__sidebar::-webkit-scrollbar {
  width: 8px;
}

.quick-preview__sidebar::-webkit-scrollbar-track {
  background: transparent;
}

.quick-preview__sidebar::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.2);
  border-radius: 4px;
}

.quick-preview__sidebar::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.3);
}
</style>
