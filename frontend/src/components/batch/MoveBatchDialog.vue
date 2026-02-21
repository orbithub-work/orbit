<!-- 移动文件批量对话框（占位组件） -->
<template>
  <Teleport to="body">
    <Transition name="dialog-fade">
      <div
        v-if="visible"
        class="dialog-overlay"
        @click.self="close"
      >
        <div class="dialog dialog--medium">
          <div class="dialog-header">
            <h3 class="dialog-title">
              移动 {{ fileIds.length }} 个文件
            </h3>
            <button
              class="dialog-close"
              @click="close"
            >
              <svg class="icon">
                <use href="#icon-close" />
              </svg>
            </button>
          </div>

          <div class="dialog-body">
            <p>移动文件功能开发中...</p>
          </div>

          <div class="dialog-footer">
            <button
              class="btn btn--secondary"
              @click="close"
            >
              取消
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
interface Props {
  visible: boolean
  fileIds: string[]
}

defineProps<Props>()

const emit = defineEmits<{
  'update:visible': [value: boolean]
  success: []
}>()

function close() {
  emit('update:visible', false)
}
</script>

<style scoped>
/* 复用 TagBatchDialog 的样式 */
</style>
