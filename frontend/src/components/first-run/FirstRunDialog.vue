<template>
  <Teleport to="body">
    <Transition name="fade">
      <div
        v-if="visible"
        class="first-run-overlay"
        @click.self="handleOverlayClick"
      >
        <div class="first-run-dialog">
          <!-- 标题区域 -->
          <div class="dialog-header">
            <div class="logo">
              <svg
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
              >
                <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z" />
              </svg>
            </div>
            <h2>是否将系统常用目录加入素材库？</h2>
            <p class="subtitle">
              智归档 OS 可以帮您管理下载、图片、视频等文件夹。我们将仅建立索引，不会移动或修改您的任何文件。
            </p>
          </div>

          <!-- 步骤1: 选择目录 -->
          <div
            v-if="currentStep === 'select'"
            class="dialog-content"
          >
            <p class="description">
              请选择要加入素材库的目录
            </p>

            <div class="directories-list">
              <div
                v-for="dir in detectedDirectories"
                :key="dir.path"
                class="directory-item"
                :class="{ selected: dir.selected, disabled: !dir.exists }"
                @click="dir.exists && toggleDirectory(dir)"
              >
                <div class="checkbox">
                  <svg
                    v-if="dir.selected"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="3"
                  >
                    <polyline points="20 6 9 17 4 12" />
                  </svg>
                </div>
                <div
                  class="icon"
                  :style="{ color: dir.color }"
                >
                  <svg
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                  >
                    <path
                      v-if="dir.type === 'documents'"
                      d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"
                    />
                    <path
                      v-if="dir.type === 'documents'"
                      d="M14 2v6h6"
                    />
                    <rect
                      v-if="dir.type === 'pictures'"
                      x="3"
                      y="3"
                      width="18"
                      height="18"
                      rx="2"
                      ry="2"
                    />
                    <circle
                      v-if="dir.type === 'pictures'"
                      cx="8.5"
                      cy="8.5"
                      r="1.5"
                    />
                    <polyline
                      v-if="dir.type === 'pictures'"
                      points="21 15 16 10 5 21"
                    />
                    <path
                      v-if="dir.type === 'music'"
                      d="M9 18V5l12-2v13"
                    />
                    <circle
                      v-if="dir.type === 'music'"
                      cx="6"
                      cy="18"
                      r="3"
                    />
                    <circle
                      v-if="dir.type === 'music'"
                      cx="18"
                      cy="16"
                      r="3"
                    />
                    <rect
                      v-if="dir.type === 'videos'"
                      x="2"
                      y="2"
                      width="20"
                      height="20"
                      rx="2.18"
                      ry="2.18"
                    />
                    <line
                      v-if="dir.type === 'videos'"
                      x1="7"
                      y1="2"
                      x2="7"
                      y2="22"
                    />
                    <line
                      v-if="dir.type === 'videos'"
                      x1="17"
                      y1="2"
                      x2="17"
                      y2="22"
                    />
                    <line
                      v-if="dir.type === 'videos'"
                      x1="2"
                      y1="12"
                      x2="22"
                      y2="12"
                    />
                    <line
                      v-if="dir.type === 'videos'"
                      x1="2"
                      y1="7"
                      x2="7"
                      y2="7"
                    />
                    <line
                      v-if="dir.type === 'videos'"
                      x1="2"
                      y1="17"
                      x2="7"
                      y2="17"
                    />
                    <line
                      v-if="dir.type === 'videos'"
                      x1="17"
                      y1="17"
                      x2="22"
                      y2="17"
                    />
                    <line
                      v-if="dir.type === 'videos'"
                      x1="17"
                      y1="7"
                      x2="22"
                      y2="7"
                    />
                  </svg>
                </div>
                <div class="info">
                  <span class="name">{{ dir.name }}</span>
                  <span class="path">{{ dir.path }}</span>
                </div>
                <div
                  v-if="!dir.exists"
                  class="status"
                >
                  未找到
                </div>
                <div
                  v-else-if="dir.fileCount !== undefined"
                  class="file-count"
                >
                  约 {{ formatNumber(dir.fileCount) }} 个文件
                </div>
              </div>
            </div>

            <div class="actions">
              <button
                class="btn-secondary"
                @click="skip"
              >
                跳过，稍后手动设置
              </button>
              <button 
                class="btn-primary" 
                :disabled="!hasSelectedDirectories"
                @click="confirmSelection"
              >
                开始索引 ({{ selectedCount }})
              </button>
            </div>
          </div>

          <!-- 步骤2: 导入进度 -->
          <div
            v-else-if="currentStep === 'importing'"
            class="dialog-content"
          >
            <div class="import-status">
              <div
                class="status-icon"
                :class="{ spinning: isImporting }"
              >
                <svg
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                >
                  <path d="M21 12a9 9 0 1 1-6.219-8.56" />
                </svg>
              </div>
              <h3>{{ importStatus }}</h3>
              <p class="status-detail">
                {{ importDetail }}
              </p>
            </div>

            <div class="progress-section">
              <div class="progress-bar">
                <div
                  class="progress-fill"
                  :style="{ width: totalProgress + '%' }"
                ></div>
              </div>
              <div class="progress-stats">
                <span>已处理: {{ processedFiles }} / {{ totalFiles }}</span>
                <span>{{ totalProgress }}%</span>
              </div>
              <div
                v-if="pendingCount > 0"
                class="background-stats"
              >
                <span class="processing-icon">⚡</span>
                <span>后台正在分析 {{ pendingCount }} 个文件 (生成指纹/缩略图)...</span>
              </div>
            </div>

            <div class="directory-progress">
              <div
                v-for="dir in importingDirectories"
                :key="dir.path"
                class="dir-progress-item"
              >
                <span class="dir-name">{{ dir.name }}</span>
                <div class="dir-progress-bar">
                  <div
                    class="dir-progress-fill"
                    :style="{ width: dir.progress + '%' }"
                  ></div>
                </div>
                <span class="dir-status">{{ dir.status }}</span>
              </div>
            </div>

            <div class="actions">
              <button 
                v-if="!isImporting"
                class="btn-primary"
                @click="finish"
              >
                完成
              </button>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { apiCall } from '@/services/api'

// 目录类型定义
interface DirectoryInfo {
  type: 'documents' | 'pictures' | 'music' | 'videos'
  name: string
  path: string
  exists: boolean
  selected: boolean
  color: string
  fileCount?: number
}

// 导入进度
interface ImportProgress {
  name: string
  path: string
  progress: number
  status: string
}

// Props
const props = defineProps<{
  visible: boolean
}>()

// Emits
const emit = defineEmits<{
  (e: 'complete', directories: string[]): void
  (e: 'skip'): void
}>()

// 状态
const currentStep = ref<'select' | 'importing'>('select')
const detectedDirectories = ref<DirectoryInfo[]>([])
const isImporting = ref(false)
const importStatus = ref('正在准备...')
const importDetail = ref('')
const processedFiles = ref(0)
const totalFiles = ref(0)
const pendingCount = ref(0)
const importingDirectories = ref<ImportProgress[]>([])

// 计算属性
const selectedCount = computed(() => 
  detectedDirectories.value.filter(d => d.selected).length
)

const hasSelectedDirectories = computed(() => selectedCount.value > 0)

const selectedDirectories = computed(() => 
  detectedDirectories.value.filter(d => d.selected)
)

const totalProgress = computed(() => {
  if (totalFiles.value === 0) return 0
  return Math.round((processedFiles.value / totalFiles.value) * 100)
})

// 方法
function toggleDirectory(dir: DirectoryInfo) {
  dir.selected = !dir.selected
}

function formatNumber(num: number): string {
  if (num >= 10000) {
    return (num / 10000).toFixed(1) + '万'
  }
  return num.toLocaleString()
}

function handleOverlayClick() {
  // 点击遮罩不关闭，防止误操作
}

async function detectDirectories() {
  try {
    const dirs = await apiCall<DirectoryInfo[]>('get_common_directories')
    const preferred = new Set(['下载', '图片', '视频', 'Downloads', 'Pictures', 'Videos'])
    detectedDirectories.value = dirs.map(dir => ({
      ...dir,
      selected: dir.exists && preferred.has(dir.name)
    }))
  } catch (error) {
    console.error('检测目录失败:', error)
    // 使用默认目录列表
    detectedDirectories.value = getDefaultDirectories()
  }
}

function getDefaultDirectories(): DirectoryInfo[] {
  const homeDir = ''
  return [
    {
      type: 'documents',
      name: '下载',
      path: `${homeDir}/Downloads`,
      exists: false,
      selected: true,
      color: '#8b5cf6'
    },
    {
      type: 'pictures',
      name: '图片',
      path: `${homeDir}/Pictures`,
      exists: false,
      selected: true,
      color: '#3b82f6'
    },
    {
      type: 'videos',
      name: '视频',
      path: `${homeDir}/Videos`,
      exists: false,
      selected: true,
      color: '#ef4444'
    },
    {
      type: 'music',
      name: '音乐',
      path: `${homeDir}/Music`,
      exists: false,
      selected: false,
      color: '#f59e0b'
    },
    {
      type: 'documents',
      name: '桌面',
      path: `${homeDir}/Desktop`,
      exists: false,
      selected: false,
      color: '#6b7280'
    }
  ]
}

async function confirmSelection() {
  if (!hasSelectedDirectories.value) return

  currentStep.value = 'importing'
  isImporting.value = true

  // 初始化导入进度
  importingDirectories.value = selectedDirectories.value.map(dir => ({
    name: dir.name,
    path: dir.path,
    progress: 0,
    status: '等待中'
  }))

  importStatus.value = '正在快速扫描文件...'

  try {
    // Call the new onboarding API
    await apiCall('run_onboarding', {
      import_downloads: selectedDirectories.value.some(d => d.name === 'Downloads' || d.name === '下载'),
      import_pictures: selectedDirectories.value.some(d => d.name === 'Pictures' || d.name === '图片'),
      import_videos: selectedDirectories.value.some(d => d.name === 'Videos' || d.name === '视频'),
      import_music: selectedDirectories.value.some(d => d.name === 'Music' || d.name === '音乐'),
      import_desktop: selectedDirectories.value.some(d => d.name === 'Desktop' || d.name === '桌面')
    })

    // 开始监听进度
    startProgressListener()
  } catch (error) {
    console.error('启动导入失败:', error)
    importStatus.value = '导入启动失败'
    importDetail.value = String(error)
    isImporting.value = false
  }
}

function startProgressListener() {
  // 使用轮询或事件监听进度
  const checkProgress = setInterval(async () => {
    try {
      const progress = await apiCall<{
        processed: number
        total: number
        directories: ImportProgress[]
        status: string
        detail: string
        complete: boolean
        pending_count?: number
      }>('get_import_progress')

      processedFiles.value = progress.processed
      totalFiles.value = progress.total
      pendingCount.value = progress.pending_count || 0
      importingDirectories.value = progress.directories
      importStatus.value = progress.status
      importDetail.value = progress.detail

      if (progress.complete) {
        clearInterval(checkProgress)
        isImporting.value = false
        importStatus.value = '导入完成！'
        importDetail.value = `共处理 ${progress.processed} 个文件`
        
        // 延迟后自动完成
        setTimeout(() => {
          finish()
        }, 1500)
      }
    } catch (error) {
      console.error('获取进度失败:', error)
    }
  }, 500)
}

function skip() {
  emit('skip')
}

function finish() {
  const paths = selectedDirectories.value.map(d => d.path)
  emit('complete', paths)
}

// 生命周期
onMounted(() => {
  if (props.visible) {
    detectDirectories()
  }
})
</script>

<style scoped>
.first-run-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.8);
  backdrop-filter: blur(8px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
}

.first-run-dialog {
  background: var(--color-surface);
  border-radius: 16px;
  width: 90%;
  max-width: 600px;
  max-height: 90vh;
  overflow: hidden;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
}

.dialog-header {
  text-align: center;
  padding: 2rem 2rem 1rem;
  background: linear-gradient(135deg, var(--color-primary) 0%, var(--color-primary-dark) 100%);
  color: white;
}

.logo {
  width: 64px;
  height: 64px;
  margin: 0 auto 1rem;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo svg {
  width: 32px;
  height: 32px;
  color: white;
}

.dialog-header h2 {
  margin: 0 0 0.5rem;
  font-size: 1.5rem;
  font-weight: 600;
}

.subtitle {
  margin: 0;
  opacity: 0.9;
  font-size: 0.95rem;
}

.dialog-content {
  padding: 1.5rem 2rem 2rem;
}

.description {
  margin: 0 0 1.5rem;
  color: var(--color-text-secondary);
  font-size: 0.9rem;
  line-height: 1.6;
}

.directories-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  margin-bottom: 1.5rem;
  max-height: 300px;
  overflow-y: auto;
}

.directory-item {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem;
  background: var(--color-surface-elevated);
  border: 2px solid transparent;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.directory-item:hover:not(.disabled) {
  border-color: var(--color-primary);
  transform: translateX(4px);
}

.directory-item.selected {
  border-color: var(--color-primary);
  background: var(--color-primary-alpha);
}

.directory-item.disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.checkbox {
  width: 24px;
  height: 24px;
  border: 2px solid var(--color-border);
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: all 0.2s;
}

.directory-item.selected .checkbox {
  background: var(--color-primary);
  border-color: var(--color-primary);
}

.checkbox svg {
  width: 16px;
  height: 16px;
  color: white;
}

.icon {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.icon svg {
  width: 24px;
  height: 24px;
}

.info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  min-width: 0;
}

.name {
  font-weight: 500;
  color: var(--color-text);
}

.path {
  font-size: 0.8rem;
  color: var(--color-text-secondary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.status {
  font-size: 0.8rem;
  color: var(--color-text-tertiary);
  padding: 0.25rem 0.5rem;
  background: var(--color-surface);
  border-radius: 4px;
}

.file-count {
  font-size: 0.8rem;
  color: var(--color-primary);
  font-weight: 500;
}

.actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
}

.btn-primary,
.btn-secondary {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 8px;
  font-size: 0.95rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary {
  background: var(--color-primary);
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: var(--color-primary-dark);
  transform: translateY(-1px);
}

.btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-secondary {
  background: transparent;
  color: var(--color-text-secondary);
  border: 1px solid var(--color-border);
}

.btn-secondary:hover {
  background: var(--color-surface-elevated);
}

/* 导入进度样式 */
.import-status {
  text-align: center;
  margin-bottom: 2rem;
}

.status-icon {
  width: 64px;
  height: 64px;
  margin: 0 auto 1rem;
  color: var(--color-primary);
}

.status-icon svg {
  width: 100%;
  height: 100%;
}

.status-icon.spinning svg {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.import-status h3 {
  margin: 0 0 0.5rem;
  font-size: 1.25rem;
  color: var(--color-text);
}

.status-detail {
  margin: 0;
  color: var(--color-text-secondary);
  font-size: 0.9rem;
}

.progress-section {
  margin-bottom: 1.5rem;
}

.progress-bar {
  height: 8px;
  background: var(--color-surface-elevated);
  border-radius: 4px;
  overflow: hidden;
  margin-bottom: 0.5rem;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, var(--color-primary) 0%, var(--color-primary-light) 100%);
  border-radius: 4px;
  transition: width 0.3s ease;
}

.progress-stats {
  display: flex;
  justify-content: space-between;
  font-size: 0.85rem;
  color: var(--color-text-secondary);
}

.background-stats {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-top: 0.75rem;
  padding: 0.75rem;
  background: var(--color-surface-elevated);
  border-radius: 8px;
  font-size: 0.85rem;
  color: var(--color-text-secondary);
  animation: fadeIn 0.3s ease;
}

.processing-icon {
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0% { opacity: 0.5; }
  50% { opacity: 1; }
  100% { opacity: 0.5; }
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(-5px); }
  to { opacity: 1; transform: translateY(0); }
}

.directory-progress {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  margin-bottom: 1.5rem;
  max-height: 200px;
  overflow-y: auto;
}

.dir-progress-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  font-size: 0.85rem;
}

.dir-name {
  width: 60px;
  color: var(--color-text);
  flex-shrink: 0;
}

.dir-progress-bar {
  flex: 1;
  height: 6px;
  background: var(--color-surface-elevated);
  border-radius: 3px;
  overflow: hidden;
}

.dir-progress-fill {
  height: 100%;
  background: var(--color-primary);
  border-radius: 3px;
  transition: width 0.3s ease;
}

.dir-status {
  width: 60px;
  text-align: right;
  color: var(--color-text-secondary);
  font-size: 0.8rem;
}

/* 过渡动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
