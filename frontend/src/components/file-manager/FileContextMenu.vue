<template>
  <Teleport to="body">
    <Transition name="fade">
      <div 
        v-if="visible" 
        class="context-menu-overlay"
        @click="close"
      >
        <div 
          class="context-menu"
          :style="{ top: position.y + 'px', left: position.x + 'px' }"
          @click.stop
        >
          <div 
            v-for="item in menuItems"
            :key="item.id"
            class="context-menu-item"
            :class="{ 
              'context-menu-item--divider': item.divider,
              'context-menu-item--disabled': item.disabled
            }"
            @click="handleItemClick(item)"
          >
            <span
              v-if="!item.divider"
              class="item-icon"
            >
              <svg
                viewBox="0 0 24 24"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path :d="item.icon" />
              </svg>
            </span>
            <span
              v-if="!item.divider"
              class="item-label"
            >{{ item.label }}</span>
            <span
              v-if="item.shortcut && !item.divider"
              class="item-shortcut"
            >{{ item.shortcut }}</span>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue'

interface FileItem {
  id: string
  name: string
  path: string
  is_directory?: boolean
}

interface MenuItem {
  id: string
  label: string
  icon: string
  shortcut?: string
  action: string
  disabled?: boolean
  divider?: boolean
}

interface Props {
  visible: boolean
  file: FileItem | null
  position: { x: number; y: number }
}

const props = defineProps<Props>()

const emit = defineEmits<{
  close: []
  action: [action: string, file: FileItem]
}>()

const menuItems = computed<MenuItem[]>(() => {
  if (!props.file) return []
  
  const isDirectory = props.file.is_directory
  
  const items: MenuItem[] = [
    {
      id: 'open',
      label: isDirectory ? '打开文件夹' : '打开',
      icon: 'M19 19H5V5h7l7 7v7zM12 5H5v14h14v-7l-7-7zm0 6l5.5 5.5H12V11zm0 2v5h5v-5h-5z',
      action: 'open'
    },
    {
      id: 'open-in-folder',
      label: '在文件夹中显示',
      icon: 'M20 6h-8l-2-2H4c-1.1 0-1.99.9-1.99 2L2 18c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V8c0-1.1-.9-2-2-2zm0 12H4V8h16v10z',
      action: 'open_in_folder',
      disabled: isDirectory
    },
    { divider: true },
    {
      id: 'rename',
      label: '重命名',
      icon: 'M3 17.25V21h3.75L17.81 9.94l-3.75-3.75L3 17.25zM20.71 7.04c.39-.39.39-1.02 0-1.41l-2.34-2.34c-.39-.39-1.02-.39-1.41 0l-1.83 1.83 3.75 3.75 1.83-1.83z',
      action: 'rename'
    },
    {
      id: 'delete',
      label: '删除',
      icon: 'M6 19c0 1.1.9 2 2 2h8c1.1 0 2-.9 2-2V7H6v12zM19 4h-3.5l-1-1h-5l-1 1H5v2h14V4z',
      action: 'delete'
    },
    { divider: true },
    {
      id: 'copy',
      label: '复制',
      icon: 'M16 1H4c-1.1 0-2 .9-2 2v14h2V3h12V1zm3 4H8c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h11c1.1 0 2-.9 2-2V7c0-1.1-.9-2-2-2zm0 16H8V7h11v14z',
      action: 'copy'
    },
    {
      id: 'cut',
      label: '剪切',
      icon: 'M9.64 7.64c.23-.5.36-1.05.36-1.64 0-2.21-1.79-4-4-4S2 3.79 2 6s1.79 4 4 4c.59 0 1.14-.13 1.64-.36L10 12l-2.36 2.36C7.14 14.13 6.59 14 6 14c-2.21 0-4 1.79-4 4s1.79 4 4 4 4-1.79 4-4c0-.59-.13-1.14-.36-1.64L12 14l7 7h3v-1L9.64 7.64zM6 8c-1.1 0-2-.89-2-2s.9-2 2-2 2 .89 2 2-.9 2-2 2zm0 12c-1.1 0-2-.89-2-2s.9-2 2-2 2 .89 2 2-.9 2-2 2zm6-7.5c-.28 0-.5-.22-.5-.5s.22-.5.5-.5.5.22.5.5-.22.5-.5.5zM19 3l-6 6 2 2 7-7V3z',
      action: 'cut'
    },
    { divider: true },
    {
      id: 'add-to-collection',
      label: '添加到收藏夹',
      icon: 'M17 3H7c-1.1 0-1.99.9-1.99 2L5 21l7-3 7 3V5c0-1.1-.9-2-2-2zm0 15l-5-2.18L7 18V5h10v13z',
      action: 'add_to_collection'
    },
    { divider: true },
    {
      id: 'properties',
      label: '属性',
      icon: 'M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 15h-2v-6h2v6zm0-8h-2V7h2v2z',
      action: 'properties'
    }
  ]
  
  return items
})

const close = () => {
  emit('close')
}

const handleItemClick = (item: MenuItem) => {
  if (item.disabled || item.divider) return
  
  if (props.file) {
    emit('action', item.action, props.file)
    close()
  }
}

watch(() => props.visible, (visible) => {
  if (visible) {
    nextTick(() => {
      adjustPosition()
    })
  }
})

const adjustPosition = () => {
  const menu = document.querySelector('.context-menu') as HTMLElement
  if (!menu) return
  
  const rect = menu.getBoundingClientRect()
  const windowWidth = window.innerWidth
  const windowHeight = window.innerHeight
  
  let x = props.position.x
  let y = props.position.y
  
  if (x + rect.width > windowWidth) {
    x = windowWidth - rect.width - 10
  }
  
  if (y + rect.height > windowHeight) {
    y = windowHeight - rect.height - 10
  }
  
  menu.style.left = x + 'px'
  menu.style.top = y + 'px'
}
</script>

<style scoped>
.context-menu-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 9999;
}

.context-menu {
  position: fixed;
  min-width: 200px;
  background-color: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 4px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  padding: 0.25rem 0;
  z-index: 10000;
}

.context-menu-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.5rem 0.75rem;
  cursor: pointer;
  transition: background-color 0.15s;
  user-select: none;
}

.context-menu-item:not(.context-menu-item--divider):not(.context-menu-item--disabled):hover {
  background-color: var(--color-primary);
  color: white;
}

.context-menu-item--disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.context-menu-item--divider {
  height: 1px;
  margin: 0.25rem 0;
  background-color: var(--color-border);
  padding: 0;
}

.item-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 18px;
  height: 18px;
  flex-shrink: 0;
}

.item-icon svg {
  width: 100%;
  height: 100%;
  fill: currentColor;
}

.item-label {
  flex: 1;
  font-size: 0.875rem;
}

.item-shortcut {
  font-size: 0.75rem;
  color: var(--color-text-secondary);
}

.context-menu-item:hover .item-shortcut {
  color: rgba(255, 255, 255, 0.7);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.15s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>