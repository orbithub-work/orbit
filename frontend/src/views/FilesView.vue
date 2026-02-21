<template>
  <div class="files-view">
    <!-- 顶部工具栏 -->
    <header class="toolbar">
      <div class="toolbar-left">
        <!-- 导航按钮 -->
        <div class="nav-buttons">
          <button
            class="nav-btn"
            :disabled="!canGoBack"
            @click="goBack"
          >
            <svg
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
            >
              <polyline points="15 18 9 12 15 6" />
            </svg>
          </button>
          <button
            class="nav-btn"
            :disabled="!canGoForward"
            @click="goForward"
          >
            <svg
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
            >
              <polyline points="9 18 15 12 9 6" />
            </svg>
          </button>
        </div>

        <!-- 路径显示 -->
        <div class="breadcrumb">
          <span class="breadcrumb-text">我问问</span>
        </div>
      </div>

      <div class="toolbar-center">
        <!-- 缩放滑块 -->
        <div class="zoom-control">
          <button
            class="zoom-btn"
            @click="decreaseZoom"
          >
            <svg
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
            >
              <circle
                cx="12"
                cy="12"
                r="10"
              />
              <line
                x1="8"
                y1="12"
                x2="16"
                y2="12"
              />
            </svg>
          </button>
          <input
            v-model="zoomLevel"
            type="range"
            min="100"
            max="400"
            class="zoom-slider"
          />
          <button
            class="zoom-btn"
            @click="increaseZoom"
          >
            <svg
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
            >
              <circle
                cx="12"
                cy="12"
                r="10"
              />
              <line
                x1="12"
                y1="8"
                x2="12"
                y2="16"
              />
              <line
                x1="8"
                y1="12"
                x2="16"
                y2="12"
              />
            </svg>
          </button>
        </div>
      </div>

      <div class="toolbar-right">
        <!-- 工具按钮 -->
        <button
          class="tool-btn"
          title="添加"
        >
          <svg
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
          >
            <line
              x1="12"
              y1="5"
              x2="12"
              y2="19"
            />
            <line
              x1="5"
              y1="12"
              x2="19"
              y2="12"
            />
          </svg>
        </button>
        <button
          class="tool-btn"
          title="文件夹"
        >
          <svg
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
          >
            <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z" />
          </svg>
        </button>
        <button
          class="tool-btn"
          title="筛选"
        >
          <svg
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
          >
            <polygon points="22 3 2 3 10 12.46 10 19 14 21 14 12.46 22 3" />
          </svg>
        </button>
        <button
          class="tool-btn"
          title="搜索"
        >
          <svg
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
          >
            <circle
              cx="11"
              cy="11"
              r="8"
            />
            <line
              x1="21"
              y1="21"
              x2="16.65"
              y2="16.65"
            />
          </svg>
        </button>

        <!-- 视图切换 -->
        <div class="view-toggle">
          <button
            :class="['view-btn', { active: viewMode === 'grid' }]"
            @click="viewMode = 'grid'"
          >
            <svg
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
            >
              <rect
                x="3"
                y="3"
                width="7"
                height="7"
              />
              <rect
                x="14"
                y="3"
                width="7"
                height="7"
              />
              <rect
                x="14"
                y="14"
                width="7"
                height="7"
              />
              <rect
                x="3"
                y="14"
                width="7"
                height="7"
              />
            </svg>
          </button>
          <button
            :class="['view-btn', { active: viewMode === 'list' }]"
            @click="viewMode = 'list'"
          >
            <svg
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
            >
              <line
                x1="8"
                y1="6"
                x2="21"
                y2="6"
              />
              <line
                x1="8"
                y1="12"
                x2="21"
                y2="12"
              />
              <line
                x1="8"
                y1="18"
                x2="21"
                y2="18"
              />
              <line
                x1="3"
                y1="6"
                x2="3.01"
                y2="6"
              />
              <line
                x1="3"
                y1="12"
                x2="3.01"
                y2="12"
              />
              <line
                x1="3"
                y1="18"
                x2="3.01"
                y2="18"
              />
            </svg>
          </button>
        </div>
      </div>
    </header>

    <!-- 主内容区 -->
    <div class="content-area">
      <!-- 空状态 -->
      <div
        v-if="files.length === 0"
        class="empty-state"
      >
        <div class="empty-illustration">
          <svg
            viewBox="0 0 200 160"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
          >
            <!-- 文件夹背景 -->
            <rect
              x="40"
              y="40"
              width="120"
              height="100"
              rx="8"
              fill="#2a2a2a"
              stroke="#3a3a3a"
              stroke-width="2"
            />
            <path
              d="M40 55 L40 48 C40 43.5817 43.5817 40 48 40 L70 40 L80 55 L40 55Z"
              fill="#333"
            />
            <!-- 放大镜 -->
            <circle
              cx="100"
              cy="85"
              r="25"
              stroke="#555"
              stroke-width="3"
              fill="none"
            />
            <line
              x1="118"
              y1="103"
              x2="135"
              y2="120"
              stroke="#555"
              stroke-width="3"
              stroke-linecap="round"
            />
            <!-- 文件图标 -->
            <rect
              x="75"
              y="70"
              width="20"
              height="25"
              rx="2"
              fill="#3a3a3a"
              stroke="#555"
              stroke-width="1.5"
            />
            <line
              x1="80"
              y1="78"
              x2="90"
              y2="78"
              stroke="#555"
              stroke-width="1.5"
            />
            <line
              x1="80"
              y1="83"
              x2="90"
              y2="83"
              stroke="#555"
              stroke-width="1.5"
            />
            <line
              x1="80"
              y1="88"
              x2="85"
              y2="88"
              stroke="#555"
              stroke-width="1.5"
            />
          </svg>
        </div>
        <h3 class="empty-title">
          没有找到相关文件
        </h3>
        <p class="empty-desc">
          也许换个搜索规则就能找到它！
        </p>
      </div>

      <!-- 文件网格 -->
      <div
        v-else
        class="file-grid-container"
        :class="{ 'with-detail': selectedFile }"
      >
        <div
          v-if="viewMode === 'grid'"
          class="file-grid"
          :style="gridStyle"
        >
          <div
            v-for="file in files"
            :key="file.id"
            :class="['file-card', { selected: selectedFile?.id === file.id, pending: file.status === 'PENDING' }]"
            @click="selectFile(file)"
            @dblclick="openFile(file)"
          >
            <div class="file-thumbnail">
              <img
                v-if="file.thumbnail || file.preview"
                :src="file.thumbnail || file.preview"
                :alt="file.name"
              />
              <div
                v-else
                class="file-placeholder"
                :class="file.type"
              >
                <span class="file-ext">{{ file.extension?.toUpperCase() }}</span>
              </div>
              <div class="file-overlay">
                <button
                  class="overlay-btn"
                  title="收藏"
                >
                  <svg
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                  >
                    <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z" />
                  </svg>
                </button>
              </div>
            </div>
            <div class="file-info">
              <div class="file-name">
                {{ file.name }}
              </div>
              <div class="file-meta">
                <span v-if="file.width && file.height">{{ file.width }} × {{ file.height }}</span>
                <span>{{ formatSize(file.size, file.status) }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 列表视图 -->
        <div
          v-else
          class="file-list"
        >
          <div class="list-header">
            <div class="list-col name">
              名称
            </div>
            <div class="list-col type">
              类型
            </div>
            <div class="list-col size">
              大小
            </div>
            <div class="list-col date">
              修改日期
            </div>
          </div>
          <div
            v-for="file in files"
            :key="file.id"
            :class="['list-row', { selected: selectedFile?.id === file.id, pending: file.status === 'PENDING' }]"
            @click="selectFile(file)"
            @dblclick="openFile(file)"
          >
            <div class="list-col name">
              <div
                class="file-icon"
                :class="file.type"
              >
                {{ file.extension?.toUpperCase() }}
              </div>
              <span>{{ file.name }}</span>
            </div>
            <div class="list-col type">
              {{ file.type || '未知' }}
            </div>
            <div class="list-col size">
              {{ formatSize(file.size, file.status) }}
            </div>
            <div class="list-col date">
              {{ formatDate(file.modifiedAt) }}
            </div>
          </div>
        </div>
      </div>

      <!-- 右侧详情面板 -->
      <aside
        v-if="selectedFile"
        class="detail-panel"
      >
        <!-- 预览区 -->
        <div class="detail-preview">
          <img
            v-if="selectedFile.thumbnail || selectedFile.preview"
            :src="selectedFile.thumbnail || selectedFile.preview"
            :alt="selectedFile.name"
          />
          <div
            v-else
            class="detail-placeholder"
            :class="selectedFile.type"
          >
            <span>{{ selectedFile.extension?.toUpperCase() }}</span>
          </div>
        </div>

        <!-- 文件信息 -->
        <div class="detail-content">
          <!-- 文件名 -->
          <div class="detail-name">
            {{ selectedFile.name }}
          </div>

          <!-- 标签区 -->
          <div class="detail-tags">
            <div
              class="tag-input"
              @click="addTag"
            >
              <svg
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
              >
                <line
                  x1="12"
                  y1="5"
                  x2="12"
                  y2="19"
                />
                <line
                  x1="5"
                  y1="12"
                  x2="19"
                  y2="12"
                />
              </svg>
              <span>添加标签</span>
            </div>
          </div>

          <!-- 信息列表 -->
          <div class="detail-info">
            <div class="info-row">
              <span class="info-label">尺寸</span>
              <span class="info-value">{{ selectedFile.width && selectedFile.height ? `${selectedFile.width} × ${selectedFile.height}` : '-' }}</span>
            </div>
            <div class="info-row">
              <span class="info-label">大小</span>
              <span class="info-value">{{ formatSize(selectedFile.size) }}</span>
            </div>
            <div class="info-row">
              <span class="info-label">类型</span>
              <span class="info-value">{{ selectedFile.type || selectedFile.extension?.toUpperCase() || '未知' }}</span>
            </div>
            <div class="info-row">
              <span class="info-label">创建时间</span>
              <span class="info-value">{{ formatDate(selectedFile.modifiedAt) }}</span>
            </div>
          </div>

          <!-- 注释 -->
          <div class="detail-notes">
            <div class="notes-header">
              <svg
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
              >
                <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7" />
                <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z" />
              </svg>
              <span>添加注释...</span>
            </div>
          </div>
        </div>

        <!-- 底部操作栏 -->
        <div class="detail-footer">
          <button
            class="footer-action"
            title="在文件夹中显示"
          >
            <svg
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
            >
              <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z" />
            </svg>
          </button>
          <button
            class="footer-action"
            title="收藏"
          >
            <svg
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
            >
              <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z" />
            </svg>
          </button>
          <button
            class="footer-action"
            title="更多"
          >
            <svg
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
            >
              <circle
                cx="12"
                cy="12"
                r="1"
              />
              <circle
                cx="19"
                cy="12"
                r="1"
              />
              <circle
                cx="5"
                cy="12"
                r="1"
              />
            </svg>
          </button>
        </div>
      </aside>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { subscribeToEvent, connectWebSocket, apiCall } from '../services/api'

interface FileItem {
	id: string
	name: string
	path: string
	size: number
	type?: string
	extension?: string
	thumbnail?: string
	preview?: string
	width?: number
	height?: number
	modifiedAt: string
	tags?: string[]
	status?: 'PENDING' | 'READY'
}

const viewMode = ref<'grid' | 'list'>('grid')
const zoomLevel = ref(200)
const selectedFile = ref<FileItem | null>(null)
const canGoBack = ref(false)
const canGoForward = ref(false)

const files = ref<FileItem[]>([])

onMounted(async () => {
	// 1. 加载现有文件
	try {
		const res = await apiCall<any[]>('get_project_files', { projectId: 'default' })
		if (res) {
			files.value = res.map(detail => {
				const file: FileItem = {
					id: detail.id,
					name: detail.path.split(/[/\\]/).pop() || '',
					path: detail.path,
					size: detail.size,
					status: detail.status,
					modifiedAt: new Date(detail.updated_at * 1000).toISOString(),
				}

				if (detail.media_meta) {
					try {
						const meta = JSON.parse(detail.media_meta)
						file.width = meta.width
						file.height = meta.height
						file.thumbnail = meta.thumbnail_path
					} catch (e) {
						console.error('Failed to parse media_meta:', e)
					}
				}
				return file
			})
		}
	} catch (err) {
		console.error('Failed to load files:', err)
	}

	// 2. 连接 WebSocket 监听实时更新
	connectWebSocket()

	// 监听发现新文件 (PENDING)
	subscribeToEvent('asset_pending', (data: any) => {
		const existing = files.value.find(f => f.id === data.id)
		if (!existing) {
			files.value.unshift({
				id: data.id,
				name: data.name,
				path: data.path,
				size: 0,
				status: 'PENDING',
				modifiedAt: new Date().toISOString()
			})
		}
	})

	// 监听处理完成 (READY)
	subscribeToEvent('asset_ready', async (data: any) => {
		const index = files.value.findIndex(f => f.id === data.asset_id)
		if (index !== -1) {
			try {
				const detail = await apiCall<any>('get_file_detail', { id: data.asset_id })
				if (detail) {
					// 处理后端返回的原始 Asset 对象
					const file: FileItem = {
						id: detail.id,
						name: detail.path.split(/[/\\]/).pop() || '',
						path: detail.path,
						size: detail.size,
						status: detail.status,
						modifiedAt: new Date(detail.updated_at * 1000).toISOString(),
					}

					if (detail.media_meta) {
						try {
							const meta = JSON.parse(detail.media_meta)
							file.width = meta.width
							file.height = meta.height
							file.thumbnail = meta.thumbnail_path
						} catch (e) {
							console.error('Failed to parse media_meta:', e)
						}
					}
					
					files.value[index] = file
				} else {
					files.value[index].status = 'READY'
				}
			} catch (err) {
				console.error('Failed to fetch file detail after ready:', err)
				files.value[index].status = 'READY'
			}
		}
	})
})

const gridStyle = computed(() => {
	const size = zoomLevel.value
	return { gridTemplateColumns: `repeat(auto-fill, minmax(${size}px, 1fr))` }
})

const selectFile = (file: FileItem) => { selectedFile.value = file }
const openFile = (file: FileItem) => { console.log('Opening:', file.name) }
const goBack = () => { console.log('Go back') }
const goForward = () => { console.log('Go forward') }
const increaseZoom = () => { zoomLevel.value = Math.min(400, zoomLevel.value + 50) }
const decreaseZoom = () => { zoomLevel.value = Math.max(100, zoomLevel.value - 50) }
const addTag = () => { console.log('Add tag') }

const formatSize = (bytes: number, status?: string): string => {
	if (status === 'PENDING') return '处理中...'
	if (bytes === 0) return '0 B'
	const k = 1024
	const sizes = ['B', 'KB', 'MB', 'GB']
	const i = Math.floor(Math.log(bytes) / Math.log(k))
	return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

const formatDate = (dateStr: string): string => {
	const date = new Date(dateStr)
	return date.toLocaleDateString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit' })
}
</script>

<style scoped>
.files-view {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: #141414;
}

/* 工具栏 */
.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 16px;
  background: #1a1a1a;
  border-bottom: 1px solid var(--color-border);
  height: 48px;
}

.toolbar-left, .toolbar-center, .toolbar-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.nav-buttons {
  display: flex;
  gap: 4px;
}

.nav-btn {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  border-radius: 4px;
  color: #666;
  cursor: pointer;
}

.nav-btn:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.05);
  color: #999;
}

.nav-btn:disabled {
  opacity: 0.3;
  cursor: not-allowed;
}

.nav-btn svg {
  width: 16px;
  height: 16px;
}

.breadcrumb-text {
  font-size: 14px;
  color: #e0e0e0;
}

/* 缩放控制 */
.zoom-control {
  display: flex;
  align-items: center;
  gap: 8px;
}

.zoom-btn {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  color: #666;
  cursor: pointer;
}

.zoom-btn:hover {
  color: #999;
}

.zoom-btn svg {
  width: 16px;
  height: 16px;
}

.zoom-slider {
  width: 80px;
  height: 3px;
  -webkit-appearance: none;
  background: #333;
  border-radius: 2px;
  outline: none;
}

.zoom-slider::-webkit-slider-thumb {
  -webkit-appearance: none;
  width: 10px;
  height: 10px;
  background: #666;
  border-radius: 50%;
  cursor: pointer;
}

.zoom-slider::-webkit-slider-thumb:hover {
  background: #888;
}

/* 工具按钮 */
.tool-btn {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  border-radius: 4px;
  color: #666;
  cursor: pointer;
}

.tool-btn:hover {
  background: rgba(255, 255, 255, 0.05);
  color: #999;
}

.tool-btn svg {
  width: 16px;
  height: 16px;
}

/* 视图切换 */
.view-toggle {
  display: flex;
  background: rgba(255, 255, 255, 0.03);
  border-radius: 4px;
  padding: 2px;
  margin-left: 8px;
}

.view-btn {
  width: 26px;
  height: 26px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  border-radius: 3px;
  color: #666;
  cursor: pointer;
}

.view-btn:hover {
  color: #999;
}

.view-btn.active {
  background: rgba(255, 255, 255, 0.08);
  color: #e0e0e0;
}

.view-btn svg {
  width: 14px;
  height: 14px;
}

/* 内容区 */
.content-area {
  display: flex;
  flex: 1;
  overflow: hidden;
}

/* 空状态 */
.empty-state {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px;
}

.empty-illustration {
  width: 200px;
  height: 160px;
  margin-bottom: 24px;
}

.empty-illustration svg {
  width: 100%;
  height: 100%;
}

.empty-title {
  font-size: 16px;
  font-weight: 500;
  color: #e0e0e0;
  margin-bottom: 8px;
}

.empty-desc {
  font-size: 13px;
  color: #666;
}

/* 文件网格 */
.file-grid-container {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}

.file-grid {
  display: grid;
  gap: 12px;
}

.file-card {
  cursor: pointer;
}

.file-card:hover .file-thumbnail {
  box-shadow: 0 0 0 1px #3a3a3a;
}

.file-card.selected .file-thumbnail {
  box-shadow: 0 0 0 2px #4a9eff;
}

.file-card.pending {
  opacity: 0.7;
  position: relative;
}

.file-card.pending::after {
  content: '处理中...';
  position: absolute;
  top: 8px;
  right: 8px;
  font-size: 10px;
  background: rgba(0, 0, 0, 0.6);
  padding: 2px 6px;
  border-radius: 10px;
  color: #fff;
  z-index: 1;
}

.file-thumbnail {
  position: relative;
  aspect-ratio: 1;
  background: #1f1f1f;
  border-radius: 6px;
  overflow: hidden;
  transition: box-shadow 0.15s;
}

.file-thumbnail img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.file-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #2a2a2a 0%, #1f1f1f 100%);
}

.file-ext {
  font-size: 20px;
  font-weight: 600;
  color: #555;
}

.file-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.4);
  display: flex;
  align-items: flex-start;
  justify-content: flex-end;
  padding: 8px;
  opacity: 0;
  transition: opacity 0.15s;
}

.file-card:hover .file-overlay {
  opacity: 1;
}

.overlay-btn {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.9);
  border: none;
  border-radius: 50%;
  color: #333;
  cursor: pointer;
}

.overlay-btn svg {
  width: 12px;
  height: 12px;
  fill: currentColor;
}

.file-info {
  padding: 6px 2px;
}

.file-name {
  font-size: 12px;
  color: #b0b0b0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.file-meta {
  display: flex;
  gap: 8px;
  margin-top: 2px;
  font-size: 11px;
  color: #555;
}

/* 列表视图 */
.file-list {
  display: flex;
  flex-direction: column;
}

.list-header {
  display: flex;
  padding: 8px 12px;
  border-bottom: 1px solid var(--color-border);
  font-size: 12px;
  color: #666;
}

.list-row {
  display: flex;
  align-items: center;
  padding: 10px 12px;
  border-bottom: 1px solid var(--color-border-light);
  cursor: pointer;
}

.list-row:hover {
  background: rgba(255, 255, 255, 0.02);
}

.list-row.selected {
  background: rgba(74, 158, 255, 0.08);
}

.list-row.pending {
  opacity: 0.7;
  font-style: italic;
}

.list-col {
  font-size: 13px;
  color: #a0a0a0;
}

.list-col.name {
  flex: 2;
  display: flex;
  align-items: center;
  gap: 10px;
}

.list-col.type, .list-col.size, .list-col.date {
  flex: 1;
}

.file-icon {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #2a2a2a;
  border-radius: 4px;
  font-size: 9px;
  font-weight: 600;
  color: #666;
}

/* 详情面板 - Eagle 风格 */
.detail-panel {
  width: 260px;
  background: #1a1a1a;
  border-left: 1px solid var(--color-border);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.detail-preview {
  width: 100%;
  aspect-ratio: 1;
  background: #141414;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.detail-preview img {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
}

.detail-placeholder {
  width: 80px;
  height: 80px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #2a2a2a 0%, #1f1f1f 100%);
  border-radius: 12px;
  font-size: 24px;
  font-weight: 600;
  color: #444;
}

.detail-content {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}

.detail-name {
  font-size: 14px;
  font-weight: 500;
  color: #e0e0e0;
  line-height: 1.4;
  word-break: break-all;
  margin-bottom: 12px;
}

/* 标签区 */
.detail-tags {
  margin-bottom: 16px;
}

.tag-input {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 10px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 4px;
  font-size: 12px;
  color: #666;
  cursor: pointer;
  transition: all 0.15s;
}

.tag-input:hover {
  background: rgba(255, 255, 255, 0.08);
  color: #888;
}

.tag-input svg {
  width: 12px;
  height: 12px;
}

/* 信息列表 */
.detail-info {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 16px;
}

.info-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
}

.info-row .info-label {
  color: #666;
}

.info-row .info-value {
  color: #999;
}

/* 注释 */
.detail-notes {
  padding-top: 12px;
  border-top: 1px solid #2a2a2a;
}

.notes-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  color: #666;
  cursor: pointer;
  padding: 4px 0;
  transition: color 0.15s;
}

.notes-header:hover {
  color: #888;
}

.notes-header svg {
  width: 14px;
  height: 14px;
}

/* 底部操作栏 */
.detail-footer {
  display: flex;
  justify-content: center;
  gap: 8px;
  padding: 12px 16px;
  border-top: 1px solid var(--color-border);
  background: #1a1a1a;
}

.footer-action {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  border-radius: 6px;
  color: #666;
  cursor: pointer;
  transition: all 0.15s;
}

.footer-action:hover {
  background: rgba(255, 255, 255, 0.05);
  color: #999;
}

.footer-action svg {
  width: 16px;
  height: 16px;
}
</style>
