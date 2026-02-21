<template>
  <div class="collections-view">
    <CollectionSidebar
      :current-collection-id="currentCollectionId"
      @select="handleSelectCollection"
    />

    <div class="content-area">
      <div
        v-if="loading"
        class="loading"
      >
        <span class="loading-spinner"></span>
        <span>加载中...</span>
      </div>

      <div
        v-else-if="!currentCollection"
        class="empty-content"
      >
        <div class="empty-icon">
          <Icon name="folder" size="xl" />
        </div>
        <p>请选择一个收藏夹</p>
        <p class="hint">
          点击左侧收藏夹查看内容
        </p>
      </div>

      <div
        v-else-if="collectionFiles.length === 0"
        class="empty-content"
      >
        <div class="empty-icon">
          <Icon name="document" size="xl" />
        </div>
        <p>{{ currentCollection.name }} 收藏夹为空</p>
        <p class="hint">
          将文件添加到收藏夹以快速访问
        </p>
      </div>

      <div
        v-else
        class="files-container"
      >
        <div class="container-header">
          <h2>{{ currentCollection.name }}</h2>
          <p class="collection-description">
            {{ currentCollection.description }}
          </p>
        </div>

        <FileGrid
          :files="collectionFiles"
          :loading="loading"
          @open="handleOpenFile"
          @context-menu="handleContextMenu"
        />
      </div>
    </div>

    <!-- 文件操作菜单 -->
    <Teleport to="body">
      <div
        v-if="showFileContextMenu"
        class="file-context-menu"
        :style="{ left: menuX + 'px', top: menuY + 'px' }"
        @click.stop
      >
        <div
          class="menu-item"
          @click="handleRemoveFromCollection"
        >
          <span class="menu-icon">➖</span>
          <span>从收藏夹移除</span>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import Icon from '@/components/common/Icon.vue'
import CollectionSidebar from '@/components/collection/CollectionSidebar.vue'
import FileGrid from '@/components/file-manager/FileGrid.vue'
import { useCollectionStore } from '@/stores/collectionStore'

const collectionStore = useCollectionStore()

const currentCollectionId = ref<string>('')
const currentCollection = ref<any>(null)
const collectionFiles = ref<any[]>([])
const loading = ref(false)
const showFileContextMenu = ref(false)
const menuX = ref(0)
const menuY = ref(0)
const selectedFile = ref<any>(null)

const handleSelectCollection = async (collection: any) => {
  currentCollectionId.value = collection.id
  currentCollection.value = collection
  loading.value = true
  await collectionStore.loadCollectionFiles(collection.id)
  collectionFiles.value = collectionStore.collectionFiles
  loading.value = false
}

const handleOpenFile = (file: any) => {
  console.log('Opening file:', file)
  // TODO: 实现打开文件功能
}

const handleContextMenu = (file: any, event: MouseEvent) => {
  showFileContextMenu.value = true
  menuX.value = event.clientX
  menuY.value = event.clientY
  selectedFile.value = file
}

const handleRemoveFromCollection = async () => {
  if (!selectedFile.value || !currentCollectionId.value) return

  try {
    await collectionStore.removeFileFromCollection(
      currentCollectionId.value,
      selectedFile.value.id
    )
    collectionFiles.value = collectionFiles.value.filter(f => f.id !== selectedFile.value.id)
    showFileContextMenu.value = false
  } catch (error) {
    console.error('Failed to remove from collection:', error)
  }
}

// 点击其他地方关闭右键菜单
const handleClickOutside = (event: MouseEvent) => {
  if (!(event.target as HTMLElement).closest('.file-context-menu')) {
    showFileContextMenu.value = false
  }
}

window.addEventListener('click', handleClickOutside)
</script>

<style scoped>
.collections-view {
  display: flex;
  height: 100%;
  width: 100%;
}

.content-area {
  flex: 1;
  height: 100%;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: var(--color-text-secondary);
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 3px solid var(--color-border);
  border-top-color: var(--color-primary);
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 1rem;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.empty-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: var(--color-text-secondary);
  text-align: center;
}

.empty-icon {
  font-size: 4rem;
  margin-bottom: 1rem;
  opacity: 0.5;
}

.empty-content p {
  margin: 0;
  font-size: 1rem;
}

.empty-content .hint {
  margin-top: 0.5rem;
  font-size: 0.875rem;
  color: var(--color-text-tertiary);
}

.files-container {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.container-header {
  padding: 1rem 1.5rem;
  border-bottom: 1px solid var(--color-border);
  background-color: var(--color-surface);
}

.container-header h2 {
  margin: 0 0 0.5rem 0;
  font-size: 1.5rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.collection-description {
  margin: 0;
  font-size: 0.875rem;
  color: var(--color-text-secondary);
  line-height: 1.5;
}

.file-context-menu {
  position: fixed;
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
  min-width: 150px;
  z-index: 1001;
  overflow: hidden;
}

.menu-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1rem;
  cursor: pointer;
  transition: background-color 0.2s;
  font-size: 0.875rem;
  color: var(--color-text-primary);
}

.menu-item:hover {
  background-color: var(--color-hover);
}

.menu-icon {
  font-size: 1rem;
}
</style>