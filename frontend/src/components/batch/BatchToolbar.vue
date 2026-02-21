<template>
  <Transition name="toolbar-slide">
    <div
      v-if="selectedCount > 0"
      class="batch-toolbar fixed bottom-6 left-1/2 -translate-x-1/2 z-[1000] flex items-center gap-6 px-5 py-3 bg-base-300/95 border border-base-content/10 rounded-box shadow-xl backdrop-blur-xl"
    >
      <div class="batch-toolbar__info flex items-center gap-3 pr-6 border-r border-base-content/10">
        <svg class="w-5 h-5 text-primary">
          <use href="#icon-check-circle" />
        </svg>
        <span class="text-sm font-semibold">已选中 {{ selectedCount }} 项</span>
        <button
          class="btn btn-link btn-xs text-base-content/60"
          @click="clearSelection"
        >
          清除
        </button>
      </div>

      <div class="batch-toolbar__actions flex items-center gap-2">
        <button
          class="btn btn-sm btn-ghost gap-2"
          :title="`为 ${selectedCount} 个文件添加标签`"
          @click="showTagDialog = true"
        >
          <svg class="w-4 h-4">
            <use href="#icon-tag" />
          </svg>
          添加标签
        </button>

        <button
          class="btn btn-sm btn-ghost gap-2"
          :title="`移动 ${selectedCount} 个文件`"
          @click="showMoveDialog = true"
        >
          <svg class="w-4 h-4">
            <use href="#icon-folder-move" />
          </svg>
          移动
        </button>

        <button
          class="btn btn-sm btn-ghost gap-2"
          :title="`归档 ${selectedCount} 个文件到项目`"
          @click="showArchiveDialog = true"
        >
          <svg class="w-4 h-4">
            <use href="#icon-archive" />
          </svg>
          归档
        </button>

        <button
          class="btn btn-sm btn-ghost gap-2"
          :title="`导出 ${selectedCount} 个文件`"
          @click="handleExport"
        >
          <svg class="w-4 h-4">
            <use href="#icon-download" />
          </svg>
          导出
        </button>

        <button
          class="btn btn-sm btn-ghost btn-error gap-2"
          :title="`删除 ${selectedCount} 个文件`"
          @click="showDeleteConfirm = true"
        >
          <svg class="w-4 h-4">
            <use href="#icon-trash" />
          </svg>
          删除
        </button>

        <div class="dropdown dropdown-top dropdown-end">
          <button
            class="btn btn-sm btn-ghost btn-circle"
            @click="showMoreMenu = !showMoreMenu"
          >
            <svg class="w-4 h-4">
              <use href="#icon-more-vertical" />
            </svg>
          </button>

          <Transition name="dropdown-fade">
            <ul
              v-if="showMoreMenu"
              class="dropdown-content z-[1] menu p-2 shadow-xl bg-base-200 rounded-box w-52"
            >
              <li>
                <button @click="handleCopyPaths">
                  <svg class="w-4 h-4">
                    <use href="#icon-copy" />
                  </svg>
                  复制路径
                </button>
              </li>
              <li>
                <button @click="handleOpenInFolder">
                  <svg class="w-4 h-4">
                    <use href="#icon-folder-open" />
                  </svg>
                  在文件夹中显示
                </button>
              </li>
              <li class="menu-title">
                <span></span>
              </li>
              <li>
                <button @click="handleCreateCollection">
                  <svg class="w-4 h-4">
                    <use href="#icon-collection" />
                  </svg>
                  创建智能集合
                </button>
              </li>
              <li>
                <button @click="handleGenerateThumbnails">
                  <svg class="w-4 h-4">
                    <use href="#icon-image" />
                  </svg>
                  重新生成缩略图
                </button>
              </li>
            </ul>
          </Transition>
        </div>
      </div>

      <TagBatchDialog
        v-model:visible="showTagDialog"
        :file-ids="selectedIds"
        @success="handleSuccess"
      />

      <MoveBatchDialog
        v-model:visible="showMoveDialog"
        :file-ids="selectedIds"
        @success="handleSuccess"
      />

      <ArchiveBatchDialog
        v-model:visible="showArchiveDialog"
        :file-ids="selectedIds"
        @success="handleSuccess"
      />

      <ConfirmDialog
        v-model:visible="showDeleteConfirm"
        title="确认删除"
        :message="`确定要删除选中的 ${selectedCount} 个文件吗？此操作不可撤销。`"
        confirm-text="删除"
        confirm-type="danger"
        :loading="deleting"
        @confirm="handleDelete"
      />
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import TagBatchDialog from './TagBatchDialog.vue'
import MoveBatchDialog from './MoveBatchDialog.vue'
import ArchiveBatchDialog from './ArchiveBatchDialog.vue'
import ConfirmDialog from '../common/ConfirmDialog.vue'

interface Props {
  selectedIds: string[]
  selectedAssets: Array<{
    id: string
    name: string
    path: string
    size: number
  }>
}

const props = defineProps<Props>()

const emit = defineEmits<{
  clearSelection: []
  success: []
  delete: [ids: string[]]
  export: [ids: string[]]
}>()

const showTagDialog = ref(false)
const showMoveDialog = ref(false)
const showArchiveDialog = ref(false)
const showDeleteConfirm = ref(false)
const showMoreMenu = ref(false)
const deleting = ref(false)

const selectedCount = computed(() => props.selectedIds.length)

function clearSelection() {
  emit('clearSelection')
}

function handleSuccess() {
  emit('success')
  showTagDialog.value = false
  showMoveDialog.value = false
  showArchiveDialog.value = false
}

async function handleDelete() {
  deleting.value = true
  try {
    emit('delete', props.selectedIds)
    showDeleteConfirm.value = false
  } finally {
    deleting.value = false
  }
}

function handleExport() {
  emit('export', props.selectedIds)
}

async function handleCopyPaths() {
  const paths = props.selectedAssets.map(a => a.path).join('\n')
  try {
    await navigator.clipboard.writeText(paths)
    showMoreMenu.value = false
  } catch (err) {
    console.error('Failed to copy paths:', err)
  }
}

function handleOpenInFolder() {
  if (props.selectedAssets.length > 0) {
    const firstAsset = props.selectedAssets[0]
    window.electron?.openInFolder(firstAsset.path)
    showMoreMenu.value = false
  }
}

function handleCreateCollection() {
  showMoreMenu.value = false
}

async function handleGenerateThumbnails() {
  showMoreMenu.value = false
}
</script>

<style scoped>
.toolbar-slide-enter-active,
.toolbar-slide-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.toolbar-slide-enter-from {
  opacity: 0;
  transform: translate(-50%, 20px);
}

.toolbar-slide-leave-to {
  opacity: 0;
  transform: translate(-50%, 20px);
}

.dropdown-fade-enter-active,
.dropdown-fade-leave-active {
  transition: all 0.2s;
}

.dropdown-fade-enter-from,
.dropdown-fade-leave-to {
  opacity: 0;
  transform: translateY(4px);
}
</style>
