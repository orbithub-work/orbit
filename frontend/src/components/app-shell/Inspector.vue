<template>
  <aside class="inspector-drawer">
    <div class="drawer-overlay" @click="$emit('close')"></div>
    <div class="drawer-content">
      <button class="close-btn-overlay" @click="$emit('close')" title="关闭">
        <Icon name="close" size="sm" />
      </button>

      <div v-if="mode === 'pool' && selectedAsset" class="drawer-body">
        <!-- 预览区域 -->
        <div class="preview-section">
          <div class="preview-box" :class="`preview-${fileCategory}`">
            <template v-if="fileCategory === 'image'">
              <img
                v-if="!imageError"
                :src="getFileUrl(selectedAsset)"
                :alt="selectedAsset.name"
                class="preview-img"
                @error="handleImageError"
              />
              <div v-else class="preview-error">
                <span class="error-icon">🖼️</span>
                <span class="error-text">预览加载失败</span>
              </div>
            </template>

            <template v-else-if="fileCategory === 'video'">
              <video
                v-if="!videoError"
                ref="videoRef"
                :src="getFileUrl(selectedAsset)"
                class="preview-video"
                controls
                preload="metadata"
                @error="handleVideoError"
              ></video>
              <div v-else class="preview-error">
                <span class="error-icon">🎬</span>
                <span class="error-text">视频加载失败</span>
              </div>
            </template>

            <template v-else-if="fileCategory === 'audio'">
              <div class="audio-preview">
                <div class="audio-bg">
                  <div class="audio-waves">
                    <span v-for="i in 20" :key="i" class="wave-bar" :style="{ height: `${Math.random() * 60 + 20}%` }"></span>
                  </div>
                </div>
                <div class="audio-info">
                  <span class="audio-icon">🎵</span>
                  <span class="audio-name">{{ selectedAsset.name }}</span>
                </div>
                <audio
                  ref="audioRef"
                  :src="getFileUrl(selectedAsset)"
                  class="audio-player"
                  controls
                  preload="metadata"
                  @error="handleAudioError"
                ></audio>
              </div>
            </template>

            <template v-else-if="fileCategory === 'document'">
              <div class="document-preview">
                <Icon :name="getFileIcon(selectedAsset)" size="xl" class="doc-icon" />
                <span class="doc-ext">{{ getFileExtension(selectedAsset).toUpperCase() }}</span>
              </div>
            </template>

            <template v-else>
              <div class="file-preview">
                <Icon :name="getFileIcon(selectedAsset)" size="xl" class="file-icon" />
                <span class="file-ext">{{ getFileExtension(selectedAsset).toUpperCase() }}</span>
              </div>
            </template>
          </div>

          <!-- 颜色提取 -->
          <div v-if="colors.length > 0" class="color-palette">
            <div
              v-for="(color, index) in colors"
              :key="index"
              class="color-swatch"
              :style="{ backgroundColor: color }"
              :title="color"
              @click="copyColor(color)"
            ></div>
          </div>
        </div>

        <!-- 文件名 -->
        <div class="filename-section">
          <div class="filename">{{ selectedAsset.name }}</div>
        </div>

        <!-- 标签区域 -->
        <div class="tags-section">
          <div class="tags-header">
            <span class="section-label">标签</span>
            <button class="add-btn" @click="showTagInput = !showTagInput">
              <Icon name="plus" size="sm" />
            </button>
          </div>
          <div v-if="tags.length > 0" class="tags-list">
            <span v-for="tag in tags" :key="tag" class="tag-item">
              {{ tag }}
              <button class="tag-remove" @click="removeTag(tag)">×</button>
            </span>
          </div>
          <div v-if="showTagInput" class="tag-input-wrapper">
            <input
              v-model="newTag"
              type="text"
              class="tag-input"
              placeholder="输入标签按回车"
              @keyup.enter="addTag"
              @blur="showTagInput = false"
            />
          </div>
        </div>

        <!-- 备注 -->
        <div class="note-section">
          <div class="section-label">备注</div>
          <textarea
            v-model="note"
            class="note-input"
            placeholder="添加备注..."
            rows="2"
            @blur="saveNote"
          ></textarea>
        </div>

        <!-- 链接 -->
        <div class="link-section">
          <div class="section-label">链接</div>
          <div class="link-input-wrapper">
            <Icon name="link" size="sm" class="link-icon" />
            <input
              v-model="linkUrl"
              type="text"
              class="link-input"
              placeholder="https://"
              @blur="saveLink"
            />
          </div>
        </div>

        <!-- 文件夹位置 -->
        <div class="folder-section">
          <div class="folder-header">
            <span class="section-label">文件夹</span>
            <span class="folder-name" :title="parentFolder">{{ parentFolder }}</span>
          </div>
        </div>

        <!-- 素材信息 -->
        <div class="info-section">
          <div class="section-title">素材信息</div>

          <!-- 评分 -->
          <div class="prop-row">
            <span class="prop-label">评分</span>
            <div class="rating-stars">
              <span
                v-for="i in 5"
                :key="i"
                class="star"
                :class="{ active: i <= rating }"
                @click="setRating(i)"
              >★</span>
            </div>
          </div>

          <!-- 类型 -->
          <div class="prop-row">
            <span class="prop-label">类型</span>
            <span class="prop-value">{{ getFileExtension(selectedAsset).toLowerCase() }}</span>
          </div>

          <!-- 尺寸 -->
          <div v-if="mediaInfo?.width && mediaInfo?.height" class="prop-row">
            <span class="prop-label">尺寸</span>
            <span class="prop-value">{{ mediaInfo.width }} × {{ mediaInfo.height }}</span>
          </div>

          <!-- 文件大小 -->
          <div class="prop-row">
            <span class="prop-label">文件大小</span>
            <span class="prop-value">{{ formatSize(selectedAsset.size) }}</span>
          </div>

          <!-- 时长 -->
          <div v-if="mediaInfo?.duration" class="prop-row">
            <span class="prop-label">时长</span>
            <span class="prop-value">{{ formatDuration(mediaInfo.duration) }}</span>
          </div>

          <!-- 导入时间 -->
          <div v-if="selectedAsset.created_at" class="prop-row">
            <span class="prop-label">导入时间</span>
            <span class="prop-value">{{ formatDate(selectedAsset.created_at) }}</span>
          </div>

          <!-- 创建时间 -->
          <div v-if="selectedAsset.created_at" class="prop-row">
            <span class="prop-label">创建时间</span>
            <span class="prop-value">{{ formatDate(selectedAsset.created_at) }}</span>
          </div>

          <!-- 修改时间 -->
          <div class="prop-row">
            <span class="prop-label">修改时间</span>
            <span class="prop-value">{{ formatDate(selectedAsset.modified_at) }}</span>
          </div>
        </div>

        <!-- 操作按钮 -->
        <div class="actions-section">
          <button class="action-btn primary" @click="exportFile">
            <Icon name="download" size="sm" />
            导出
          </button>
          <button class="action-btn" @click="openFile">
            <Icon name="external-link" size="sm" />
            打开文件
          </button>
          <button class="action-btn" @click="openInFolder">
            <Icon name="folder-open" size="sm" />
            打开目录
          </button>
        </div>
      </div>

      <div v-else-if="mode === 'project' && project" class="drawer-body">
        <div class="project-cover">{{ project.name }}</div>

        <div class="info-section">
          <div class="section-title">项目信息</div>

          <div class="prop-row">
            <span class="prop-label">客户</span>
            <span class="prop-value">{{ project.client }}</span>
          </div>
          <div class="prop-row">
            <span class="prop-label">负责人</span>
            <span class="prop-value">{{ project.owner }}</span>
          </div>
          <div class="prop-row">
            <span class="prop-label">截止</span>
            <span class="prop-value">{{ project.deadline }}</span>
          </div>
          <div class="prop-row">
            <span class="prop-label">素材量</span>
            <span class="prop-value">{{ project.assets }} 个</span>
          </div>
          <div class="prop-row">
            <span class="prop-label">交付项</span>
            <span class="prop-value">{{ project.deliverables }}</span>
          </div>
        </div>
      </div>

      <div v-else class="drawer-empty">
        <span class="empty-icon">📋</span>
        <span class="empty-text">请选择一个{{ mode === 'pool' ? '文件' : '项目' }}</span>
      </div>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import Icon from '@/components/common/Icon.vue'
import type { FileItem } from '@/stores/fileStore'
import { apiCall } from '@/services/api'

interface Project {
  id: string | number
  name: string
  client: string
  owner: string
  deadline: string
  assets: number
  deliverables: string
}

const props = defineProps<{
  mode: 'pool' | 'project' | 'artifact'
  selectedId: number | null
  project?: Project
  selectedAsset?: FileItem | null
}>()

const emit = defineEmits<{
  close: []
}>()

const videoRef = ref<HTMLVideoElement | null>(null)
const audioRef = ref<HTMLAudioElement | null>(null)
const videoError = ref(false)
const audioError = ref(false)
const imageError = ref(false)

// 标签功能
const tags = ref<string[]>([])
const newTag = ref('')
const showTagInput = ref(false)

// 备注功能
const note = ref('')

// 链接功能
const linkUrl = ref('')

// 评分功能
const rating = ref(0)

// 颜色提取（模拟）
const colors = ref<string[]>(['#E8B89D', '#2D2D2D', '#8B7355', '#D4A574', '#A67B5B', '#4A4A4A', '#F5DEB3'])

const fileCategory = computed(() => {
  if (!props.selectedAsset) return 'other'
  return getFileCategory(props.selectedAsset)
})

interface MediaInfo {
  width?: number
  height?: number
  duration?: number
  format?: string
}

const mediaInfo = computed((): MediaInfo | null => {
  // TODO: 从后端获取真实的媒体信息
  if (!props.selectedAsset?.path) return null
  // 模拟数据
  if (fileCategory.value === 'image') {
    return { width: 1920, height: 1080, format: 'JPEG' }
  }
  if (fileCategory.value === 'video') {
    return { width: 720, height: 1280, duration: 65, format: 'MP4' }
  }
  return null
})

// 获取父文件夹名称
const parentFolder = computed(() => {
  if (!props.selectedAsset?.path) return ''
  const parts = props.selectedAsset.path.split('\\')
  return parts[parts.length - 2] || ''
})

watch(() => props.selectedAsset, (newAsset) => {
  videoError.value = false
  audioError.value = false
  imageError.value = false

  if (videoRef.value) {
    videoRef.value.pause()
    videoRef.value.currentTime = 0
  }
  if (audioRef.value) {
    audioRef.value.pause()
    audioRef.value.currentTime = 0
  }

  // 重置编辑状态
  if (newAsset) {
    // TODO: 从后端加载实际的标签、备注、链接、评分数据
    tags.value = []
    note.value = ''
    linkUrl.value = ''
    rating.value = 0
  }
})

function getFileCategory(file: FileItem): string {
  if (file.is_directory) return 'folder'

  const ext = getFileExtension(file).toLowerCase()

  if (['jpg', 'jpeg', 'png', 'gif', 'webp', 'bmp', 'svg', 'ico', 'tiff', 'tif'].includes(ext)) return 'image'
  if (['mp4', 'avi', 'mov', 'mkv', 'webm', 'wmv', 'flv', 'm4v', '3gp'].includes(ext)) return 'video'
  if (['mp3', 'wav', 'flac', 'aac', 'ogg', 'wma', 'm4a', 'aiff'].includes(ext)) return 'audio'
  if (['pdf', 'doc', 'docx', 'xls', 'xlsx', 'ppt', 'pptx', 'txt', 'rtf'].includes(ext)) return 'document'
  if (['psd', 'ai', 'eps', 'sketch', 'fig', 'xd', 'indd'].includes(ext)) return 'design'
  if (['zip', 'rar', '7z', 'tar', 'gz', 'bz2'].includes(ext)) return 'archive'
  if (['ttf', 'otf', 'woff', 'woff2', 'eot'].includes(ext)) return 'font'

  return 'other'
}

function getFileUrl(file: FileItem): string {
  const baseUrl = localStorage.getItem('api_server_port')
    ? `http://localhost:${localStorage.getItem('api_server_port')}`
    : 'http://localhost:32000'
  return `${baseUrl}/api/v1/assets/file?path=${encodeURIComponent(file.path)}`
}

function getFileExtension(file: FileItem): string {
  return file.name.split('.').pop() || ''
}

function getFileIcon(file: FileItem): string {
  if (file.is_directory) return 'folder'

  const category = getFileCategory(file)
  const ext = getFileExtension(file).toLowerCase()

  switch (category) {
    case 'image': return 'image'
    case 'video': return 'video'
    case 'audio': return 'audio'
    case 'document':
      if (ext === 'pdf') return 'document'
      if (['doc', 'docx'].includes(ext)) return 'document'
      if (['xls', 'xlsx'].includes(ext)) return 'chart-bar'
      if (['ppt', 'pptx'].includes(ext)) return 'chart'
      return 'document'
    case 'design': return 'sparkles'
    case 'archive': return 'archive'
    case 'font': return 'code'
    default: return 'file'
  }
}

function handleImageError() {
  imageError.value = true
}

function handleVideoError() {
  videoError.value = true
}

function handleAudioError() {
  audioError.value = true
}

function formatSize(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

function formatDate(date: string | Date): string {
  const d = typeof date === 'string' ? new Date(date) : date
  return new Intl.DateTimeFormat('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  }).format(d)
}

function formatDuration(seconds: number): string {
  const h = Math.floor(seconds / 3600)
  const m = Math.floor((seconds % 3600) / 60)
  const s = Math.floor(seconds % 60)

  if (h > 0) {
    return `${h}:${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`
  }
  return `${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`
}

// 标签操作
function addTag() {
  if (newTag.value.trim() && !tags.value.includes(newTag.value.trim())) {
    tags.value.push(newTag.value.trim())
    newTag.value = ''
    showTagInput.value = false
    // TODO: 保存到后端
  }
}

function removeTag(tag: string) {
  tags.value = tags.value.filter(t => t !== tag)
  // TODO: 保存到后端
}

// 评分操作
function setRating(value: number) {
  rating.value = value
  // TODO: 保存到后端
}

// 保存备注
function saveNote() {
  // TODO: 保存到后端
}

// 保存链接
function saveLink() {
  // TODO: 保存到后端
}

// 复制颜色
function copyColor(color: string) {
  navigator.clipboard.writeText(color)
}

// 文件操作
async function openFile() {
  if (!props.selectedAsset) return
  try {
    await apiCall('open_file', { path: props.selectedAsset.path, asset_id: props.selectedAsset.id })
  } catch (e) {
    console.error('Failed to open file:', e)
  }
}

async function openInFolder() {
  if (!props.selectedAsset) return
  try {
    await apiCall('open_in_folder', { path: props.selectedAsset.path, asset_id: props.selectedAsset.id })
  } catch (e) {
    console.error('Failed to open in folder:', e)
  }
}

async function exportFile() {
  if (!props.selectedAsset) return
  try {
    await apiCall('export_file', { path: props.selectedAsset.path, asset_id: props.selectedAsset.id })
  } catch (e) {
    console.error('Failed to export file:', e)
  }
}
</script>

<style scoped>
.inspector-drawer {
  position: fixed;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;
  z-index: 100;
  display: flex;
  justify-content: flex-end;
}

.drawer-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.4);
  backdrop-filter: blur(4px);
  -webkit-backdrop-filter: blur(4px);
  cursor: pointer;
  animation: fadeIn 0.2s ease;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.drawer-content {
  position: relative;
  width: 300px;
  background: #1e1e1e;
  border-left: 1px solid #2d2d2d;
  display: flex;
  flex-direction: column;
  box-shadow: -4px 0 24px rgba(0, 0, 0, 0.4);
  animation: slideIn 0.25s ease;
}

.close-btn-overlay {
  position: absolute;
  top: 50%;
  left: -16px;
  transform: translateY(-50%);
  width: 32px;
  height: 32px;
  border-radius: 50%;
  border: none;
  background: rgba(40, 40, 40, 0.95);
  color: #9ca3af;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
  z-index: 10;
}

.close-btn-overlay:hover {
  background: rgba(60, 60, 60, 0.95);
  color: #e5e7eb;
}

@keyframes slideIn {
  from { transform: translateX(100%); }
  to { transform: translateX(0); }
}

.drawer-body {
  flex: 1;
  overflow-y: auto;
  padding: 0;
}

.drawer-empty {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  color: #6b7280;
  padding: 40px 20px;
}

.empty-icon {
  font-size: 40px;
  opacity: 0.5;
}

.empty-text {
  font-size: 12px;
}

/* Preview Section */
.preview-section {
  padding: 16px;
  background: #252526;
  border-bottom: 1px solid #2d2d2d;
}

.preview-box {
  width: 100%;
  aspect-ratio: 1;
  background: #1e1e1e;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  position: relative;
  margin-bottom: 12px;
}

.preview-img {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.preview-video {
  width: 100%;
  height: 100%;
  object-fit: contain;
  background: #000;
}

.preview-error {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
}

.preview-error .error-icon {
  font-size: 32px;
  opacity: 0.5;
}

.preview-error .error-text {
  font-size: 11px;
  color: #6b7280;
}

/* Color Palette */
.color-palette {
  display: flex;
  gap: 4px;
  justify-content: center;
}

.color-swatch {
  width: 28px;
  height: 28px;
  border-radius: 4px;
  cursor: pointer;
  border: 1px solid rgba(255, 255, 255, 0.1);
  transition: transform 0.15s;
}

.color-swatch:hover {
  transform: scale(1.1);
}

/* Audio Preview */
.audio-preview {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%);
}

.audio-bg {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
}

.audio-waves {
  display: flex;
  align-items: center;
  gap: 2px;
  height: 50px;
}

.wave-bar {
  width: 3px;
  background: linear-gradient(to top, #3b82f6, #8b5cf6);
  border-radius: 2px;
  opacity: 0.6;
  animation: wave 1s ease-in-out infinite;
}

.wave-bar:nth-child(odd) {
  animation-delay: 0.2s;
}

@keyframes wave {
  0%, 100% { transform: scaleY(0.5); }
  50% { transform: scaleY(1); }
}

.audio-info {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 8px;
  background: rgba(0, 0, 0, 0.3);
}

.audio-icon {
  font-size: 24px;
  margin-bottom: 4px;
}

.audio-name {
  font-size: 10px;
  color: #d1d5db;
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.audio-player {
  width: 100%;
  height: 32px;
}

/* Document & File Preview */
.document-preview,
.file-preview {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.doc-icon,
.file-icon {
  font-size: 48px;
  opacity: 0.6;
}

.doc-ext,
.file-ext {
  font-size: 12px;
  font-weight: 500;
  color: #9ca3af;
  background: #3a3a3a;
  padding: 2px 8px;
  border-radius: 3px;
}

/* Filename Section */
.filename-section {
  padding: 12px 16px;
  border-bottom: 1px solid #2d2d2d;
}

.filename {
  font-size: 13px;
  font-weight: 500;
  color: #e5e7eb;
  word-break: break-all;
  line-height: 1.4;
}

/* Common Section Styles */
.section-label {
  font-size: 11px;
  color: #6b7280;
  font-weight: 500;
}

/* Tags Section */
.tags-section {
  padding: 12px 16px;
  border-bottom: 1px solid #2d2d2d;
}

.tags-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
}

.add-btn {
  width: 20px;
  height: 20px;
  border-radius: 4px;
  border: none;
  background: #3a3a3a;
  color: #9ca3af;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s;
}

.add-btn:hover {
  background: #4a4a4a;
  color: #e5e7eb;
}

.tags-list {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.tag-item {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
  background: #3a3a3a;
  border-radius: 4px;
  font-size: 11px;
  color: #d1d5db;
}

.tag-remove {
  width: 14px;
  height: 14px;
  border: none;
  background: none;
  color: #6b7280;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  line-height: 1;
  padding: 0;
}

.tag-remove:hover {
  color: #ef4444;
}

.tag-input-wrapper {
  margin-top: 8px;
}

.tag-input {
  width: 100%;
  padding: 6px 10px;
  background: #2a2a2a;
  border: 1px solid #3a3a3a;
  border-radius: 4px;
  font-size: 11px;
  color: #e5e7eb;
  outline: none;
}

.tag-input:focus {
  border-color: #007acc;
}

/* Note Section */
.note-section {
  padding: 12px 16px;
  border-bottom: 1px solid #2d2d2d;
}

.note-section .section-label {
  display: block;
  margin-bottom: 8px;
}

.note-input {
  width: 100%;
  padding: 8px 10px;
  background: #2a2a2a;
  border: 1px solid #3a3a3a;
  border-radius: 4px;
  font-size: 11px;
  color: #e5e7eb;
  resize: none;
  outline: none;
  font-family: inherit;
}

.note-input:focus {
  border-color: #007acc;
}

.note-input::placeholder {
  color: #6b7280;
}

/* Link Section */
.link-section {
  padding: 12px 16px;
  border-bottom: 1px solid #2d2d2d;
}

.link-section .section-label {
  display: block;
  margin-bottom: 8px;
}

.link-input-wrapper {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 10px;
  background: #2a2a2a;
  border: 1px solid #3a3a3a;
  border-radius: 4px;
}

.link-icon {
  color: #6b7280;
  flex-shrink: 0;
}

.link-input {
  flex: 1;
  border: none;
  background: none;
  font-size: 11px;
  color: #e5e7eb;
  outline: none;
}

.link-input::placeholder {
  color: #6b7280;
}

/* Folder Section */
.folder-section {
  padding: 12px 16px;
  border-bottom: 1px solid #2d2d2d;
}

.folder-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.folder-name {
  font-size: 11px;
  color: #9ca3af;
  max-width: 180px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* Info Section */
.info-section {
  padding: 12px 16px;
  border-bottom: 1px solid #2d2d2d;
}

.section-title {
  font-size: 11px;
  font-weight: 600;
  color: #e5e7eb;
  margin-bottom: 12px;
}

.prop-row {
  display: flex;
  align-items: center;
  margin-bottom: 8px;
}

.prop-row:last-child {
  margin-bottom: 0;
}

.prop-label {
  width: 60px;
  flex-shrink: 0;
  font-size: 11px;
  color: #6b7280;
}

.prop-value {
  flex: 1;
  font-size: 11px;
  color: #d1d5db;
}

/* Rating Stars */
.rating-stars {
  display: flex;
  gap: 4px;
}

.star {
  font-size: 14px;
  color: #3a3a3a;
  cursor: pointer;
  transition: color 0.15s;
}

.star.active {
  color: #f59e0b;
}

.star:hover {
  color: #fbbf24;
}

/* Actions Section */
.actions-section {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 10px;
  background: #2a2a2a;
  border: 1px solid #3a3a3a;
  border-radius: 6px;
  color: #d1d5db;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.15s;
}

.action-btn:hover {
  background: #3a3a3a;
  border-color: #4a4a4a;
}

.action-btn.primary {
  background: #007acc;
  border-color: #007acc;
  color: white;
}

.action-btn.primary:hover {
  background: #0066aa;
}

/* Project Cover */
.project-cover {
  margin: 16px;
  width: calc(100% - 32px);
  height: 100px;
  background: linear-gradient(135deg, #3b82f6 0%, #8b5cf6 100%);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 600;
  color: white;
}
</style>
