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
          <span class="icon">×</span>
        </button>
      </div>

      <div class="dialog-body">
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
            创建第一个收藏夹来管理您的文件
          </p>
        </div>

        <div
          v-else
          class="collection-list"
        >
          <div class="list-header">
            <span>收藏夹</span>
            <span class="item-count">文件数量</span>
          </div>

          <div
            v-for="collection in collections"
            :key="collection.id"
            class="collection-item"
            :class="{ 'selected': selectedCollectionIds.has(collection.id) }"
            @click="handleCollectionClick(collection)"
          >
            <div class="item-info">
              <div class="item-name">
                {{ collection.name }}
              </div>
              <div class="item-description">
                {{ collection.description }}
              </div>
            </div>
            <div class="item-count">
              {{ collection.itemCount }}
            </div>
          </div>
        </div>

        <div class="create-collection-section">
          <button
            class="create-btn"
            @click="handleCreateCollection"
          >
            <span class="icon">+</span>
            新建收藏夹
          </button>
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
          :disabled="selectedCollectionIds.size === 0"
          @click="handleConfirm"
        >
          添加到收藏夹
        </button>
      </div>
    </div>
  </div>

  <!-- 新建收藏夹对话框 -->
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
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useCollectionStore } from '@/stores/collectionStore'

const collectionStore = useCollectionStore()

interface Props {
  visible: boolean
  fileIds: string[]
  multiple?: boolean
  title?: string
}

const props = withDefaults(defineProps<Props>(), {
  multiple: true,
  title: '添加到收藏夹'
})

const emit = defineEmits<{
  confirm: [collectionIds: string[]]
  cancel: []
}>()

const loading = ref(false)
const collections = ref<any[]>([])
const selectedCollectionIds = ref<Set<string>>(new Set())
const showCreateDialog = ref(false)
const newCollectionName = ref('')
const newCollectionDescription = ref('')

const loadCollections = async () => {
  loading.value = true
  await collectionStore.loadCollections()
  collections.value = collectionStore.collections
  loading.value = false
}

watch(() => props.visible, (newVal) => {
  if (newVal) {
    loadCollections()
    selectedCollectionIds.value.clear()
  }
})

const handleCollectionClick = (collection: any) => {
  if (selectedCollectionIds.value.has(collection.id)) {
    selectedCollectionIds.value.delete(collection.id)
  } else {
    if (!props.multiple) {
      selectedCollectionIds.value.clear()
    }
    selectedCollectionIds.value.add(collection.id)
  }
}

const handleCreateCollection = () => {
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
    const newCollection = await collectionStore.createCollection(
      newCollectionName.value.trim(),
      newCollectionDescription.value.trim() || undefined
    )
    collections.value.push(newCollection)
    selectedCollectionIds.value.add(newCollection.id)
    showCreateDialog.value = false
  } catch (error) {
    console.error('Failed to create collection:', error)
  }
}

const handleConfirm = () => {
  const collectionIds = Array.from(selectedCollectionIds.value)
  if (collectionIds.length > 0) {
    emit('confirm', collectionIds)
  }
}

const handleCancel = () => {
  emit('cancel')
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
  background-color: white;
  border-radius: 8px;
  width: 500px;
  max-width: 90%;
  max-height: 80%;
  display: flex;
  flex-direction: column;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
}

.dialog-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1rem 1.5rem;
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
  width: 30px;
  height: 30px;
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
  text-align: center;
  padding: 2rem;
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
  margin-bottom: 1rem;
}

.list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.5rem 0;
  font-size: 0.75rem;
  font-weight: 500;
  color: var(--color-text-secondary);
  text-transform: uppercase;
  border-bottom: 1px solid var(--color-border);
}

.collection-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75rem;
  border-radius: 6px;
  cursor: pointer;
  transition: background-color 0.2s;
  margin-top: 0.25rem;
}

.collection-item:hover {
  background-color: var(--color-hover);
}

.collection-item.selected {
  background-color: var(--color-primary-light);
  color: var(--color-primary);
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

.item-description {
  font-size: 0.75rem;
  color: var(--color-text-secondary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.item-count {
  font-size: 0.75rem;
  color: var(--color-text-secondary);
  margin-left: 1rem;
}

.create-collection-section {
  padding-top: 1rem;
  border-top: 1px solid var(--color-border);
}

.create-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 1rem;
  background-color: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 6px;
  font-size: 0.875rem;
  cursor: pointer;
  transition: background-color 0.2s;
  width: 100%;
}

.create-btn:hover {
  background-color: var(--color-hover);
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

.btn-confirm:hover:not(:disabled) {
  background-color: var(--color-primary-hover);
}

.btn-confirm:disabled {
  background-color: var(--color-border);
  cursor: not-allowed;
}

/* 新建收藏夹对话框样式 */
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
</style>
