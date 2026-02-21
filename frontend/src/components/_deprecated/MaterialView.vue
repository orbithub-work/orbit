<template>
  <div class="material-view grid-12">
    <!-- Main Content Area -->
    <div 
      class="view-content"
      :class="selectedItem ? 'col-span-9' : 'col-span-12'"
    >
      <!-- Filter Bar -->
      <div class="filter-bar">
        <div class="filter-group">
          <button
            v-for="filter in typeFilters"
            :key="filter.key"
            :class="['pill', { active: activeTypeFilter === filter.key }]"
            @click="activeTypeFilter = filter.key"
          >
            {{ filter.label }}
          </button>
        </div>
        <div class="view-controls">
          <div class="slider-container">
            <svg class="icon icon--16"><use href="#i-image" /></svg>
            <input
              v-model="cardSize"
              type="range"
              min="100"
              max="400"
              class="size-slider"
            />
          </div>
        </div>
      </div>

      <!-- Scrollable Grid Area -->
      <div class="scroll-container">
        <!-- Debug Info -->
        <div style="padding: 10px; background: rgba(255,0,0,0.1); margin-bottom: 10px; border-radius: 4px; font-size: 12px; flex-shrink: 0;">
          Loading: {{ loading }} | Materials: {{ materials.length }} | Filtered: {{ filteredMaterials.length }} | Filter: {{ activeTypeFilter }}
        </div>

        <!-- Debug Warning -->
        <div style="background: red; color: white; padding: 10px; text-align: center; font-weight: bold; margin-bottom: 10px; border-radius: 4px;">
          DEBUG MODE: FILE LOADING DISABLED
        </div>

        <!-- Skeleton -->
        <MaterialViewSkeleton v-if="loading" />

        <!-- Material Grid -->
        <div
          v-else-if="filteredMaterials.length > 0"
          class="material-grid"
          :style="{ '--card-size': cardSize + 'px' }"
        >
          <div
            v-for="item in filteredMaterials"
            :key="item.id"
            class="material-card card"
            :class="{ selected: selectedItem?.id === item.id }"
            @click="selectMaterial(item)"
          >
            <!-- Thumbnail Wrapper -->
            <div class="card-thumbnail">
              <img
                v-if="item.thumbnail"
                :src="item.thumbnail"
                :alt="item.name"
                loading="lazy"
              />
              <div
                v-else
                class="thumbnail-placeholder"
                :class="item.type"
              >
                <svg class="icon icon--32">
                  <use :href="getTypeIcon(item.type)" />
                </svg>
              </div>

              <!-- Hover Overlay Actions -->
              <div class="card-overlay">
                <button
                  class="action-btn"
                  title="预览"
                >
                  <svg class="icon icon--16"><use href="#i-eye" /></svg>
                </button>
                <button
                  class="action-btn"
                  title="收藏"
                >
                  <svg class="icon icon--16"><use href="#i-star" /></svg>
                </button>
              </div>

              <!-- Duration Badge (Video) -->
              <span
                v-if="item.type === 'video'"
                class="duration-badge"
              >02:14</span>
            </div>

            <!-- Card Meta -->
            <div class="card-meta">
              <div
                class="meta-title"
                :title="item.name"
              >
                {{ item.name }}
              </div>
              <div class="meta-details">
                <span class="meta-ext">{{ getExtension(item.name) }}</span>
                <span class="meta-size">{{ item.size }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Empty State -->
        <div
          v-else
          class="empty-state"
        >
          <div class="empty-icon">
            <svg class="icon icon--64"><use href="#i-folder-open" /></svg>
          </div>
          <p class="empty-text">
            {{ error ? '加载失败' : '暂无素材' }}
          </p>
          <p class="empty-subtext">
            {{ error ? error : '拖入文件或点击添加文件夹' }}
          </p>
        </div>
      </div>
    </div>

    <!-- Inspector Panel -->
    <div
      v-if="selectedItem"
      class="view-inspector col-span-3"
    >
      <Inspector 
        :item="selectedItem" 
        @close="selectedItem = null"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue'
import MaterialViewSkeleton from './MaterialViewSkeleton.vue'
import Inspector from './Inspector.vue'
import { useFolderStore } from '@/stores/folderStore'
import { useProjectStore } from '@/stores/projectStore'
import { apiCall } from '@/services/api'

const folderStore = useFolderStore()
const projectStore = useProjectStore()
const activeTypeFilter = ref('all')
const loading = ref(true)
const error = ref<string | null>(null)
const selectedItem = ref<any>(null)
const cardSize = ref(200)

// Mock Data
const materials = ref<any[]>([])

// Computed property for filtered materials
const filteredMaterials = computed(() => {
  if (activeTypeFilter.value === 'all') {
    return materials.value
  }
  return materials.value.filter(item => item.type === activeTypeFilter.value)
})

// Watch for folder selection change
watch(() => folderStore.selectedFolder, (newFolder) => {
  if (newFolder) {
    refreshMaterials(newFolder.path)
  }
})

const refreshMaterials = async (path: string) => {
  // DEBUG: Prevent crash by delaying and catching
  console.log('MaterialView: refreshMaterials start')
  loading.value = true
  error.value = null
  
  // Artificial delay to see if initial render works
  await new Promise(resolve => setTimeout(resolve, 1000))

  try {
    // console.log('MaterialView: Fetching projects...')
    // // Get projects first to find the default project
    // const projects = await apiCall<any[]>('list_projects')
    // console.log('MaterialView: Projects loaded:', projects)
    
    // if (!projects || !Array.isArray(projects)) {
    //      console.warn('MaterialView: projects is not an array', projects)
    //      materials.value = []
    //      return
    // }

    // const defaultProject = projects.find((p: any) => p.name === 'Default Library')

    // if (!defaultProject) {
    //   error.value = 'Default Library project not found'
    //   console.warn('MaterialView: Default Library project not found')
    //   materials.value = []
    //   return
    // }
    console.log('MaterialView: Projects fetching DISABLED for debugging')
    materials.value = []
    loading.value = false
    return

    // console.log('MaterialView: Fetching files...')
    // // Get files from the default project
    // const files = await apiCall<any[]>('get_project_files', {
    //   projectId: defaultProject.id
    // })
    // console.log('MaterialView: Files loaded in refresh (count):', files?.length)

    // if (files && files.length > 0) {
    //   materials.value = files.map((file: any) => ({
    //     id: file.id,
    //     name: file.name,
    //     type: getFileType(file.file_type || ''),
    //     size: formatFileSize(file.size),
    //     path: file.path,
    //     thumbnail: file.thumbnail_path,
    //     status: 'online'
    //   }))
    //   // console.log('MaterialView: Materials mapped in refresh:', materials.value)
    // } else {
    //   materials.value = []
    // }
    materials.value = []
    console.log('MaterialView: Files fetching DISABLED for debugging')
  } catch (e: any) {
    const msg = e?.message || 'Unknown error'
    error.value = msg
    console.error('MaterialView: Failed to load materials (sanitized):', msg)
    materials.value = []
  } finally {
    loading.value = false
  }
}

const typeFilters = [
  { key: 'all', label: '全部' },
  { key: 'video', label: '视频' },
  { key: 'image', label: '图片' },
  { key: 'audio', label: '音频' },
  { key: 'project', label: '工程' }
]

const getFileType = (mimeType: string): string => {
  const mime_lower = mimeType.toLowerCase()
  if (mime_lower.startsWith('image/')) {
    return 'image'
  }
  if (mime_lower.startsWith('video/')) {
    return 'video'
  }
  if (mime_lower.startsWith('audio/')) {
    return 'audio'
  }
  return 'project'
}

const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i]
}

const getTypeIcon = (type: string) => {
  const map: Record<string, string> = {
    video: '#i-video',
    image: '#i-image',
    audio: '#i-audio',
    project: '#i-file'
  }
  return map[type] || '#i-file'
}

const getExtension = (name: string) => {
  if (!name) return ''
  return name.includes('.') ? name.split('.').pop()?.toUpperCase() : ''
}

const selectMaterial = (item: any) => {
  selectedItem.value = item
}

onMounted(async () => {
  // Load materials from default project
  console.log('MaterialView: onMounted called')
  await refreshMaterials('/')
})
</script>

<style scoped>
.material-view {
  height: 100%;
  overflow: hidden;
  background: var(--color-bg-base);
}

.view-content {
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
  transition: all var(--motion-duration) var(--motion-easing);
}

/* Filter Bar */
.filter-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--spacing-md) var(--spacing-lg);
  background: rgba(18, 18, 20, 0.2);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid var(--glass-border);
  flex-shrink: 0;
}

.pill {
  padding: 6px 16px;
  border-radius: 99px;
  background: rgba(255, 255, 255, 0.03);
  color: var(--color-text-secondary);
  font-size: 13px;
  cursor: pointer;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid var(--glass-border);
}

.pill:hover {
  color: var(--color-text-primary);
  background: var(--color-bg-surface-hover);
  border-color: var(--glass-border-bright);
  transform: translateY(-1px);
}

.pill.active {
  background: var(--grad-primary);
  color: #020617;
  border-color: transparent;
  font-weight: 700;
  box-shadow: 0 4px 15px rgba(125, 211, 252, 0.4);
}

.size-slider {
  width: 100px;
  accent-color: var(--color-text-secondary);
}

.scroll-container {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
  padding: var(--spacing-lg);
  min-height: 0;
  display: flex;
  flex-direction: column;
}

/* Grid Layout */
.material-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(var(--card-size), 1fr));
  gap: var(--spacing-md);
  padding-bottom: 40px;
  width: 100%;
}

.material-card {
  position: relative;
  overflow: hidden;
  border-radius: var(--radius-md);
  transition: all 0.15s ease-out;
  background: var(--color-bg-surface);
  border: 1px solid var(--color-border);
}

.material-card:hover {
  border-color: var(--color-border-hover);
  background: var(--color-bg-surface-hover);
  transform: translateY(-2px);
}

.material-card.selected {
  border-color: var(--color-primary);
  background: var(--color-primary-bg);
  box-shadow: 0 0 0 1px var(--color-primary);
}

.card-thumbnail {
  width: 100%;
  aspect-ratio: 4/3; /* Consistent aspect ratio */
  background: rgba(0, 0, 0, 0.2);
  position: relative;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
}

.card-thumbnail img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: opacity 0.2s;
}

.material-card:hover .card-thumbnail img {
  opacity: 0.9;
}

.thumbnail-placeholder {
  color: var(--color-text-tertiary);
  opacity: 0.5;
}

/* Hover Overlay */
.card-overlay {
  position: absolute;
  top: 6px;
  right: 6px;
  display: flex;
  gap: 4px;
  opacity: 0;
  transition: opacity 0.2s;
}

.material-card:hover .card-overlay {
  opacity: 1;
}

.action-btn {
  width: 24px;
  height: 24px;
  border-radius: 4px;
  background: rgba(0, 0, 0, 0.6);
  color: white;
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: background 0.2s;
}

.action-btn:hover {
  background: var(--color-primary);
}

.duration-badge {
  position: absolute;
  bottom: 6px;
  right: 6px;
  background: rgba(0, 0, 0, 0.6);
  color: white;
  font-size: 10px;
  padding: 2px 4px;
  border-radius: 2px;
  font-family: var(--font-mono);
}

/* Card Meta */
.card-meta {
  padding: 8px 10px;
}

.meta-title {
  font-size: var(--font-size-sm);
  color: var(--color-text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-bottom: 2px;
}

.meta-details {
  display: flex;
  justify-content: space-between;
  font-size: var(--font-size-xs);
  color: var(--color-text-tertiary);
}

/* Empty State */
.empty-state {
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: var(--color-text-tertiary);
  width: 100%;
}

.empty-icon {
  color: var(--color-bg-surface-active);
  margin-bottom: var(--spacing-md);
}
.view-inspector {
  background: var(--color-bg-sidebar);
  border-left: 1px solid var(--color-border);
  height: 100%;
  overflow-y: auto;
}
</style>
