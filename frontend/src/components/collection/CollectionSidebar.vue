<template>
  <div class="collection-sidebar">
    <div class="sidebar-header">
      <h3>收藏夹</h3>
      <button
        class="create-btn"
        @click="handleCreate"
      >
        <Icon name="plus" size="sm" />
        新建
      </button>
    </div>

    <div
      v-if="loading"
      class="loading"
    >
      <span class="loading-spinner"></span>
      <span>加载中...</span>
    </div>

    <div
      v-else-if="collections.length === 0"
      class="empty"
    >
      <p>暂无收藏夹</p>
      <p class="hint">
        点击"新建"创建第一个收藏夹
      </p>
    </div>

    <div
      v-else
      class="collection-list"
    >
      <div
        v-for="collection in collections"
        :key="collection.id"
        class="collection-item"
        :class="{ 'active': currentCollectionId === collection.id }"
        @click="handleSelect(collection)"
        @contextmenu.prevent="handleContextMenu(collection, $event)"
      >
        <div class="item-icon">
          <Icon name="folder" size="sm" />
        </div>
        <div class="item-info">
          <div class="item-name">
            {{ collection.name }}
          </div>
          <div class="item-count">
            {{ collection.itemCount }} 项
          </div>
        </div>
      </div>
    </div>

    <!-- 创建收藏夹对话框 -->
    <Teleport to="body">
      <div
        v-if="showCreateDialog"
        class="dialog-overlay"
        @click="handleCreateCancel"
      >
        <div
          class="dialog"
          @click.stop
        >
          <h4>创建新收藏夹</h4>
          <div class="form-group">
            <label>名称</label>
            <input
              v-model="newCollectionName"
              type="text"
              placeholder="请输入收藏夹名称"
              @keydown.enter="handleCreateConfirm"
            />
          </div>
          <div class="form-group">
            <label>描述（可选）</label>
            <input
              v-model="newCollectionDescription"
              type="text"
              placeholder="请输入收藏夹描述"
            />
          </div>
          <div class="dialog-actions">
            <button
              class="btn-cancel"
              @click="handleCreateCancel"
            >
              取消
            </button>
            <button
              class="btn-confirm"
              @click="handleCreateConfirm"
            >
              创建
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- 右键菜单 -->
    <Teleport to="body">
      <div
        v-if="showContextMenu"
        class="context-menu"
        :style="{ left: menuX + 'px', top: menuY + 'px' }"
        @click.stop
      >
        <div
          class="menu-item"
          @click="handleRename"
        >
          <span class="menu-icon">✏️</span>
          <span>重命名</span>
        </div>
        <div
          class="menu-item"
          @click="handleDelete"
        >
          <span class="menu-icon">🗑️</span>
          <span>删除</span>
        </div>
      </div>
    </Teleport>

    <!-- 重命名对话框 -->
    <Teleport to="body">
      <div
        v-if="showRenameDialog"
        class="dialog-overlay"
        @click="handleRenameCancel"
      >
        <div
          class="dialog"
          @click.stop
        >
          <h4>重命名收藏夹</h4>
          <div class="form-group">
            <label>新名称</label>
            <input
              v-model="renameCollectionName"
              type="text"
              placeholder="请输入新名称"
              @keydown.enter="handleRenameConfirm"
            />
          </div>
          <div class="dialog-actions">
            <button
              class="btn-cancel"
              @click="handleRenameCancel"
            >
              取消
            </button>
            <button
              class="btn-confirm"
              @click="handleRenameConfirm"
            >
              确认
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- 删除确认对话框 -->
    <Teleport to="body">
      <div
        v-if="showDeleteDialog"
        class="dialog-overlay"
        @click="handleDeleteCancel"
      >
        <div
          class="dialog"
          @click.stop
        >
          <h4>删除收藏夹</h4>
          <p>确定要删除收藏夹 "{{ selectedCollection?.name }}" 吗？</p>
          <p class="warning">
            此操作将移除该收藏夹中的所有文件关联，但不会删除文件本身。
          </p>
          <div class="dialog-actions">
            <button
              class="btn-cancel"
              @click="handleDeleteCancel"
            >
              取消
            </button>
            <button
              class="btn-danger"
              @click="handleDeleteConfirm"
            >
              删除
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import Icon from '@/components/common/Icon.vue'
import { useCollectionStore } from '@/stores/collectionStore'

const collectionStore = useCollectionStore()

interface Props {
  currentCollectionId?: string
}

const props = withDefaults(defineProps<Props>(), {
  currentCollectionId: ''
})

const emit = defineEmits<{
  select: [collection: any]
  'add-files': [collectionId: string, fileIds: string[]]
}>()

const loading = ref(false)
const collections = ref<any[]>([])
const showCreateDialog = ref(false)
const newCollectionName = ref('')
const newCollectionDescription = ref('')
const showContextMenu = ref(false)
const menuX = ref(0)
const menuY = ref(0)
const selectedCollection = ref<any>(null)
const showRenameDialog = ref(false)
const renameCollectionName = ref('')
const showDeleteDialog = ref(false)

const loadCollections = async () => {
  loading.value = true
  await collectionStore.loadCollections()
  collections.value = collectionStore.collections
  loading.value = false
}

watch(() => collectionStore.collections, (newVal) => {
  collections.value = newVal
})

// 初始化加载
loadCollections()

const handleCreate = () => {
  showCreateDialog.value = true
  newCollectionName.value = ''
  newCollectionDescription.value = ''
}

const handleCreateCancel = () => {
  showCreateDialog.value = false
}

const handleCreateConfirm = async () => {
  if (!newCollectionName.value.trim()) {
    return
  }

  try {
    await collectionStore.createCollection(
      newCollectionName.value.trim(),
      newCollectionDescription.value.trim() || undefined
    )
    showCreateDialog.value = false
  } catch (error) {
    console.error('Failed to create collection:', error)
  }
}

const handleSelect = (collection: any) => {
  emit('select', collection)
}

const handleContextMenu = (collection: any, event: MouseEvent) => {
  showContextMenu.value = true
  menuX.value = event.clientX
  menuY.value = event.clientY
  selectedCollection.value = collection
}

const handleRename = () => {
  if (!selectedCollection.value) return
  showContextMenu.value = false
  showRenameDialog.value = true
  renameCollectionName.value = selectedCollection.value.name
}

const handleRenameCancel = () => {
  showRenameDialog.value = false
}

const handleRenameConfirm = async () => {
  if (!selectedCollection.value || !renameCollectionName.value.trim()) {
    return
  }

  try {
    await collectionStore.renameCollection(
      selectedCollection.value.id,
      renameCollectionName.value.trim()
    )
    showRenameDialog.value = false
  } catch (error) {
    console.error('Failed to rename collection:', error)
  }
}

const handleDelete = () => {
  if (!selectedCollection.value) return
  showContextMenu.value = false
  showDeleteDialog.value = true
}

const handleDeleteCancel = () => {
  showDeleteDialog.value = false
}

const handleDeleteConfirm = async () => {
  if (!selectedCollection.value) return

  try {
    await collectionStore.deleteCollection(selectedCollection.value.id)
    showDeleteDialog.value = false
  } catch (error) {
    console.error('Failed to delete collection:', error)
  }
}

// 点击其他地方关闭右键菜单
const handleClickOutside = (event: MouseEvent) => {
  if (!(event.target as HTMLElement).closest('.collection-item') &&
      !(event.target as HTMLElement).closest('.context-menu')) {
    showContextMenu.value = false
  }
}

window.addEventListener('click', handleClickOutside)
</script>

<style scoped>
.collection-sidebar {
  width: 250px;
  height: 100%;
  background-color: var(--color-surface);
  border-right: 1px solid var(--color-border);
  display: flex;
  flex-direction: column;
}

.sidebar-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1rem;
  border-bottom: 1px solid var(--color-border);
}

.sidebar-header h3 {
  margin: 0;
  font-size: 1rem;
  font-weight: 500;
  color: var(--color-text-primary);
}

.create-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 1rem;
  background-color: var(--color-primary);
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 0.875rem;
  cursor: pointer;
  transition: background-color 0.2s;
}

.create-btn:hover {
  background-color: var(--color-primary-hover);
}

.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 2rem;
  color: var(--color-text-secondary);
}

.loading-spinner {
  width: 24px;
  height: 24px;
  border: 2px solid var(--color-border);
  border-top-color: var(--color-primary);
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 0.5rem;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.empty {
  padding: 2rem;
  text-align: center;
  color: var(--color-text-secondary);
}

.empty p {
  margin: 0;
  font-size: 0.875rem;
}

.empty .hint {
  margin-top: 0.5rem;
  font-size: 0.75rem;
  color: var(--color-text-tertiary);
}

.collection-list {
  flex: 1;
  overflow: auto;
  padding: 0.5rem;
}

.collection-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem;
  border-radius: 6px;
  cursor: pointer;
  transition: background-color 0.2s;
  margin-bottom: 0.25rem;
}

.collection-item:hover {
  background-color: var(--color-hover);
}

.collection-item.active {
  background-color: var(--color-primary-light);
  color: var(--color-primary);
}

.item-icon {
  font-size: 1.25rem;
}

.item-info {
  flex: 1;
  min-width: 0;
}

.item-name {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--color-text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.item-count {
  font-size: 0.75rem;
  color: var(--color-text-secondary);
}

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
  background-color: white;
  border-radius: 8px;
  padding: 1.5rem;
  min-width: 300px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
}

.dialog h4 {
  margin: 0 0 1rem 0;
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.form-group {
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--color-text-primary);
}

.form-group input {
  width: 100%;
  padding: 0.5rem;
  border: 1px solid var(--color-border);
  border-radius: 6px;
  font-size: 0.875rem;
  box-sizing: border-box;
}

.form-group input:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px var(--color-primary-light);
}

.dialog-actions {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  margin-top: 1.5rem;
}

.btn-cancel,
.btn-confirm,
.btn-danger {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 6px;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s;
}

.btn-cancel {
  background-color: var(--color-surface);
  color: var(--color-text-primary);
}

.btn-cancel:hover {
  background-color: var(--color-hover);
}

.btn-confirm {
  background-color: var(--color-primary);
  color: white;
}

.btn-confirm:hover {
  background-color: var(--color-primary-hover);
}

.btn-danger {
  background-color: var(--color-danger);
  color: white;
}

.btn-danger:hover {
  background-color: var(--color-danger-hover);
}

.context-menu {
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

.warning {
  color: var(--color-danger);
  font-size: 0.75rem;
  margin-top: 0.5rem;
}
</style>
