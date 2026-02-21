<template>
  <header class="title-bar">
    <div class="title-left">
      <div class="app-mark">智归档OS</div>
    </div>
    <div class="title-tabs">
      <button
        class="btn btn-sm btn-ghost title-tab"
        :class="{ 'btn-active': modelValue === 'pool' }"
        @click="$emit('update:modelValue', 'pool')"
      >
        素材池
      </button>
      <button
        class="btn btn-sm btn-ghost title-tab"
        :class="{ 'btn-active': modelValue === 'project' }"
        @click="$emit('update:modelValue', 'project')"
      >
        项目
      </button>
    </div>
    <div class="title-right">
      <button
        class="btn btn-sm btn-circle btn-ghost"
        :class="{ 'btn-active': showInspector }"
        @click="$emit('toggle-inspector')"
        title="文件信息"
      >
        <Icon name="info" size="sm" />
      </button>
      <button
        class="btn btn-sm btn-circle btn-ghost"
        @click="goToSettings"
        title="设置"
      >
        <Icon name="settings" size="sm" />
      </button>
      <div class="window-controls">
        <button class="window-btn minimize" @click="minimizeWindow">
          <Icon name="minus" size="sm" />
        </button>
        <button class="window-btn maximize" @click="maximizeWindow">
          <Icon name="square" size="sm" />
        </button>
        <button class="window-btn close" @click="closeWindow">
          <Icon name="close" size="sm" />
        </button>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import Icon from '@/components/common/Icon.vue'

defineProps<{
  modelValue: 'pool' | 'project'
  showInspector: boolean
}>()

defineEmits<{
  'update:modelValue': [value: 'pool' | 'project']
  'toggle-inspector': []
}>()

const router = useRouter()

function goToSettings() {
  router.push('/settings')
}

function minimizeWindow() {
  const handler = (window as any)?.mediaAssistant?.window?.minimize
  if (typeof handler === 'function') handler()
}

function maximizeWindow() {
  const handler = (window as any)?.mediaAssistant?.window?.maximize
  if (typeof handler === 'function') handler()
}

function closeWindow() {
  const handler = (window as any)?.mediaAssistant?.window?.close
  if (typeof handler === 'function') handler()
  else window.close()
}
</script>

<style scoped>
.title-bar {
  height: 44px;
  display: grid;
  grid-template-columns: 1fr auto 1fr;
  align-items: center;
  padding: 0 14px;
  background: var(--color-bg-sidebar, #1b1c1f);
  border-bottom: 1px solid var(--color-border, #2b2b2f);
  -webkit-app-region: drag;
  z-index: 50;
}

.title-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.app-mark {
  font-size: 12px;
  font-weight: 600;
  color: var(--color-text-primary, #d1d5db);
}

.title-tabs {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 3px;
  background: var(--color-bg-surface, #202124);
  border-radius: 8px;
  border: 1px solid var(--color-border, #2b2b2f);
}

.title-tabs .title-tab {
  -webkit-app-region: no-drag;
  min-height: 28px;
  padding: 4px 16px;
  font-size: 13px;
  font-weight: 600;
}

.title-right {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 4px;
}

.title-right .btn {
  -webkit-app-region: no-drag;
}

.window-controls {
  display: flex;
  gap: 8px;
  -webkit-app-region: no-drag;
  margin-left: 8px;
}

.window-btn {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  border: none;
  cursor: pointer;
  font-size: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: opacity 0.15s;
}

.window-btn.minimize {
  background: #fbd875;
  color: rgba(0, 0, 0, 0.5);
}

.window-btn.maximize {
  background: #40c764;
  color: rgba(0, 0, 0, 0.5);
}

.window-btn.close {
  background: #f9615b;
  color: rgba(0, 0, 0, 0.5);
}

.window-btn:hover {
  opacity: 0.8;
}
</style>
