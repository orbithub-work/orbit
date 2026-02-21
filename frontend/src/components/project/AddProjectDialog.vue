<template>
  <div
    v-if="visible"
    class="dialog-overlay"
    @click="handleCancel"
  >
    <div
      class="dialog"
      @click.stop
    >
      <div class="dialog-header">
        <h3>{{ title }}</h3>
        <button
          class="close-btn"
          @click="handleCancel"
        >
          <Icon name="close" size="sm" />
        </button>
      </div>

      <div class="dialog-body">
        <div class="form-section">
          <label>项目名称</label>
          <input
            v-model="projectName"
            type="text"
            placeholder="请输入项目名称"
            class="form-input"
            @keydown.enter="handleConfirm"
          />
        </div>

        <div class="form-section">
          <label>选择文件夹</label>
          <div class="folder-selector">
            <input
              v-model="folderPath"
              type="text"
              placeholder="点击右侧按钮选择文件夹"
              class="form-input"
              readonly
            />
            <button
              class="browse-btn"
              @click="browseFolder"
            >
              <Icon name="folder" size="sm" />
              浏览
            </button>
          </div>
        </div>

        <div class="form-section">
          <label class="checkbox-label">
            <input
              v-model="autoSync"
              type="checkbox"
            />
            <span>自动同步文件夹内容</span>
            <span class="hint">当文件夹中有新文件时自动添加到项目</span>
          </label>
        </div>

        <div
          v-if="selectedFiles.length > 0"
          class="preview-section"
        >
          <div class="preview-header">
            <span class="preview-title">将导入的文件 ({{ selectedFiles.length }} 个)</span>
          </div>
          <div class="preview-list">
            <div
              v-for="file in selectedFiles.slice(0, 5)"
              :key="file.path"
              class="preview-item"
            >
              <Icon :name="getFileIcon(file.name)" size="sm" class="file-icon" />
              <span class="file-name">{{ file.name }}</span>
            </div>
            <div
              v-if="selectedFiles.length > 5"
              class="preview-more"
            >
              还有 {{ selectedFiles.length - 5 }} 个文件...
            </div>
          </div>
        </div>

        <div
          v-if="loading"
          class="loading"
        >
          <span class="loading-spinner"></span>
          <span>{{ loadingText }}</span>
        </div>

        <div
          v-if="errorMessage"
          class="error-message"
        >
          {{ errorMessage }}
        </div>
      </div>

      <div class="dialog-footer">
        <button
          class="btn-cancel"
          @click="handleCancel"
        >
          取消
        </button>
        <button
          class="btn-confirm"
          :disabled="!canConfirm || loading"
          @click="handleConfirm"
        >
          创建项目
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import Icon from '@/components/common/Icon.vue'
import { useProjectStore } from '@/stores/projectStore'

interface FileInfo {
  name: string
  path: string
  size: number
  type: string
}

const props = withDefaults(defineProps<{
  visible: boolean
  title?: string
}>(), {
  title: '从文件夹创建项目'
})

const emit = defineEmits<{
  confirm: [projectData: { name: string; folderPath: string; autoSync: boolean }]
  cancel: []
}>()

const projectStore = useProjectStore()

const projectName = ref('')
const folderPath = ref('')
const autoSync = ref(true)
const loading = ref(false)
const loadingText = ref('')
const errorMessage = ref('')
const selectedFiles = ref<FileInfo[]>([])

const canConfirm = computed(() => {
  return projectName.value.trim().length > 0 && folderPath.value.trim().length > 0
})

watch(folderPath, async (newPath) => {
  if (newPath && newPath.length > 0) {
    await previewFiles(newPath)
  } else {
    selectedFiles.value = []
  }
})

const browseFolder = async () => {
  try {
    const result = await projectStore.selectProjectFolder()

    if (result) {
      folderPath.value = result as string

      if (!projectName.value.trim()) {
        const folderName = folderPath.value.split(/[/\\]/).pop() || ''
        projectName.value = folderName
      }
    }
  } catch (error) {
    console.error('Failed to browse folder:', error)
    errorMessage.value = '无法打开文件夹选择器'
  }
}

const previewFiles = async (path: string) => {
  loading.value = true
  loadingText.value = '正在扫描文件夹...'
  errorMessage.value = ''

  try {
    // 暂时移除 readDir 预览功能，因为需要后端配合实现对应的 Wails 方法
    // 如果需要此功能，后续应在 Go 后端添加相应的方法并通过 apiCall 调用
    selectedFiles.value = []
  } catch (error) {
    console.error('Failed to preview files:', error)
    errorMessage.value = '无法预览文件夹内容'
  } finally {
    loading.value = false
  }
}

const getFileType = (filename: string): string => {
  const ext = filename.split('.').pop()?.toLowerCase() || ''
  const imageExts = ['jpg', 'jpeg', 'png', 'gif', 'bmp', 'webp', 'svg', 'ico']
  const videoExts = ['mp4', 'avi', 'mkv', 'mov', 'wmv', 'flv', 'webm']
  const audioExts = ['mp3', 'wav', 'flac', 'aac', 'ogg', 'wma']

  if (imageExts.includes(ext)) return 'image'
  if (videoExts.includes(ext)) return 'video'
  if (audioExts.includes(ext)) return 'audio'
  return 'file'
}

const getFileIcon = (filename: string): string => {
  const type = getFileType(filename)
  const icons: Record<string, string> = {
    image: 'image',
    video: 'video',
    audio: 'audio',
    file: 'file'
  }
  return icons[type] || icons.file
}

const handleConfirm = async () => {
  if (!canConfirm.value || loading.value) return

  loading.value = true
  errorMessage.value = ''
  loadingText.value = '正在创建项目...'

  try {
    await projectStore.createProjectFromFolder(
      folderPath.value,
      autoSync.value
    )
    emit('confirm', {
      name: projectName.value.trim(),
      folderPath: folderPath.value,
      autoSync: autoSync.value
    })
    resetForm()
  } catch (error: any) {
    console.error('Failed to create project:', error)
    errorMessage.value = error.message || '创建项目失败'
  } finally {
    loading.value = false
    loadingText.value = ''
  }
}

const handleCancel = () => {
  if (!loading.value) {
    resetForm()
    emit('cancel')
  }
}

const resetForm = () => {
  projectName.value = ''
  folderPath.value = ''
  autoSync.value = true
  selectedFiles.value = []
  errorMessage.value = ''
  loading.value = false
  loadingText.value = ''
}
</script>

<style scoped>
.dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.dialog {
  background-color: var(--color-surface);
  border-radius: 12px;
  width: 560px;
  max-width: 90%;
  max-height: 85%;
  display: flex;
  flex-direction: column;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2);
  border: 1px solid var(--color-border);
}

.dialog-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1.25rem 1.5rem;
  border-bottom: 1px solid var(--color-border);
}

.dialog-header h3 {
  margin: 0;
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.close-btn {
  background: none;
  border: none;
  font-size: 1.5rem;
  color: var(--color-text-secondary);
  cursor: pointer;
  padding: 0;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  transition: background-color 0.2s;
}

.close-btn:hover {
  background-color: var(--color-hover);
}

.dialog-body {
  flex: 1;
  overflow: auto;
  padding: 1.5rem;
}

.form-section {
  margin-bottom: 1.25rem;
}

.form-section > label {
  display: block;
  margin-bottom: 0.5rem;
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--color-text-primary);
}

.form-input {
  width: 100%;
  padding: 0.625rem 0.875rem;
  border: 1px solid var(--color-border);
  border-radius: 8px;
  font-size: 0.875rem;
  color: var(--color-text-primary);
  background-color: var(--color-surface-light);
  box-sizing: border-box;
  transition: all 0.2s;
}

.form-input:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px var(--color-primary-light);
  background-color: var(--color-surface);
}

.form-input::placeholder {
  color: var(--color-text-tertiary);
}

.folder-selector {
  display: flex;
  gap: 0.5rem;
}

.folder-selector .form-input {
  flex: 1;
}

.browse-btn {
  display: flex;
  align-items: center;
  gap: 0.375rem;
  padding: 0.625rem 1rem;
  background-color: var(--color-surface-light);
  border: 1px solid var(--color-border);
  border-radius: 8px;
  font-size: 0.875rem;
  color: var(--color-text-primary);
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}

.browse-btn:hover {
  background-color: var(--color-hover);
  border-color: var(--color-primary);
}

.checkbox-label {
  display: flex;
  align-items: flex-start;
  gap: 0.5rem;
  cursor: pointer;
}

.checkbox-label input[type="checkbox"] {
  width: 18px;
  height: 18px;
  margin: 0;
  cursor: pointer;
  accent-color: var(--color-primary);
}

.checkbox-label span {
  font-size: 0.875rem;
  color: var(--color-text-primary);
}

.checkbox-label .hint {
  display: block;
  font-size: 0.75rem;
  color: var(--color-text-secondary);
  margin-top: 0.125rem;
}

.preview-section {
  margin-top: 1.5rem;
  padding: 1rem;
  background-color: var(--color-surface-light);
  border-radius: 8px;
  border: 1px solid var(--color-border);
}

.preview-header {
  margin-bottom: 0.75rem;
}

.preview-title {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--color-text-primary);
}

.preview-list {
  display: flex;
  flex-direction: column;
  gap: 0.375rem;
}

.preview-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.375rem 0;
}

.file-icon {
  font-size: 1rem;
}

.file-name {
  font-size: 0.8125rem;
  color: var(--color-text-secondary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.preview-more {
  font-size: 0.75rem;
  color: var(--color-text-tertiary);
  font-style: italic;
}

.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 1.5rem;
  color: var(--color-text-secondary);
}

.loading-spinner {
  width: 28px;
  height: 28px;
  border: 3px solid var(--color-border);
  border-top-color: var(--color-primary);
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 0.75rem;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.error-message {
  margin-top: 1rem;
  padding: 0.75rem 1rem;
  background-color: var(--color-error-light);
  border: 1px solid var(--color-error);
  border-radius: 8px;
  color: var(--color-error);
  font-size: 0.875rem;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  padding: 1rem 1.5rem;
  border-top: 1px solid var(--color-border);
  background-color: var(--color-surface-light);
}

.btn-cancel,
.btn-confirm {
  padding: 0.625rem 1.25rem;
  border: none;
  border-radius: 8px;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-cancel {
  background-color: var(--color-surface);
  color: var(--color-text-primary);
  border: 1px solid var(--color-border);
}

.btn-cancel:hover {
  background-color: var(--color-hover);
}

.btn-confirm {
  background-color: var(--color-primary);
  color: white;
}

.btn-confirm:hover:not(:disabled) {
  background-color: var(--color-primary-hover);
  transform: translateY(-1px);
}

.btn-confirm:disabled {
  background-color: var(--color-border);
  color: var(--color-text-tertiary);
  cursor: not-allowed;
  transform: none;
}
</style>
