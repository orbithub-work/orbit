<template>
  <Teleport to="body">
    <Transition name="dialog-fade">
      <dialog
        v-if="visible"
        class="modal modal-open"
        @click.self="close"
      >
        <div class="modal-box">
          <h3 class="text-lg font-bold">{{ title }}</h3>
          <p class="py-4">{{ message }}</p>
          <div class="modal-action">
            <button
              class="btn btn-ghost"
              @click="close"
            >
              {{ cancelText }}
            </button>
            <button
              class="btn"
              :class="confirmTypeClass"
              :disabled="loading"
              @click="handleConfirm"
            >
              <span v-if="loading" class="loading loading-spinner loading-sm"></span>
              {{ loading ? '处理中...' : confirmText }}
            </button>
          </div>
        </div>
        <form method="dialog" class="modal-backdrop">
          <button @click="close">close</button>
        </form>
      </dialog>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  visible: boolean
  title: string
  message: string
  confirmText?: string
  cancelText?: string
  confirmType?: 'primary' | 'danger' | 'warning'
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  confirmText: '确认',
  cancelText: '取消',
  confirmType: 'primary',
  loading: false
})

const emit = defineEmits<{
  'update:visible': [value: boolean]
  confirm: []
  cancel: []
}>()

const confirmTypeClass = computed(() => {
  switch (props.confirmType) {
    case 'danger':
      return 'btn-error'
    case 'warning':
      return 'btn-warning'
    default:
      return 'btn-primary'
  }
})

function close() {
  emit('update:visible', false)
  emit('cancel')
}

function handleConfirm() {
  emit('confirm')
}
</script>

<style scoped>
.modal {
  z-index: 9999;
}

.modal-box {
  background: var(--color-bg-surface, #1e1e1e);
  border: 1px solid var(--color-border, rgba(255, 255, 255, 0.1));
}

.modal-backdrop {
  background: rgba(0, 0, 0, 0.7);
  backdrop-filter: blur(4px);
}

.dialog-fade-enter-active,
.dialog-fade-leave-active {
  transition: opacity 0.2s;
}

.dialog-fade-enter-active .modal-box,
.dialog-fade-leave-active .modal-box {
  transition: transform 0.2s, opacity 0.2s;
}

.dialog-fade-enter-from,
.dialog-fade-leave-to {
  opacity: 0;
}

.dialog-fade-enter-from .modal-box,
.dialog-fade-leave-to .modal-box {
  transform: scale(0.95);
  opacity: 0;
}
</style>
