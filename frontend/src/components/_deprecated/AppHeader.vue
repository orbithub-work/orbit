<template>
  <header class="app-header">
    <!-- 左侧：Logo + 搜索 -->
    <div class="header-left">
      <div class="logo">
        <div class="logo-icon">
          <svg
            class="icon icon--16"
            viewBox="0 0 24 24"
            aria-hidden="true"
          >
            <use href="#i-box" />
          </svg>
        </div>
        <span class="logo-text">智归档OS</span>
      </div>
      
      <div class="search-box">
        <svg
          class="search-icon icon icon--16"
          viewBox="0 0 24 24"
          aria-hidden="true"
        >
          <use href="#i-search" />
        </svg>
        <input 
          v-model="searchQuery" 
          type="text"
          class="search-input"
          placeholder="搜索素材/项目/记录"
          @input="handleSearch"
        />
      </div>
    </div>

    <!-- 右侧：设置 + 用户 + 窗口控制 -->
    <div class="header-right">
      <button
        class="header-btn"
        title="显示/隐藏 Dock"
        @click="$emit('toggle-dock')"
      >
        <svg
          class="icon icon--16"
          viewBox="0 0 24 24"
          aria-hidden="true"
        >
          <use href="#i-box" />
        </svg>
      </button>
      <button
        class="header-btn"
        title="设置"
        @click="$emit('open-settings')"
      >
        <svg
          class="icon icon--16"
          viewBox="0 0 24 24"
          aria-hidden="true"
        >
          <use href="#i-settings" />
        </svg>
      </button>
      
      <!-- 暂时隐藏“我的”页面相关入口 -->
      <!-- <div
        class="user-avatar"
        title="用户"
      >
        <svg
          class="icon icon--16"
          viewBox="0 0 24 24"
          aria-hidden="true"
        >
          <use href="#i-user" />
        </svg>
      </div> -->

      <!-- 窗口控制按钮 -->
      <div class="window-controls">
        <button
          class="window-btn close"
          title="关闭"
          @click="closeWindow"
        >
          <svg
            class="icon icon--16"
            viewBox="0 0 24 24"
            aria-hidden="true"
          >
            <use href="#i-window-close" />
          </svg>
        </button>
        <button
          class="window-btn minimize"
          title="最小化"
          @click="minimizeWindow"
        >
          <svg
            class="icon icon--16"
            viewBox="0 0 24 24"
            aria-hidden="true"
          >
            <use href="#i-window-minimize" />
          </svg>
        </button>
        <button
          class="window-btn maximize"
          title="最大化"
          @click="maximizeWindow"
        >
          <svg
            class="icon icon--16"
            viewBox="0 0 24 24"
            aria-hidden="true"
          >
            <use href="#i-window-maximize" />
          </svg>
        </button>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const searchQuery = ref('')

defineEmits<{
  'open-settings': [],
  'toggle-dock': []
}>()

// 处理搜索输入变化
const handleSearch = () => {
  console.log('搜索:', searchQuery.value)
}

// 最小化窗口
const minimizeWindow = () => {
  const handler = (window as any)?.mediaAssistant?.window?.minimize
  if (typeof handler === 'function') {
    handler()
  }
}

// 最大化/还原窗口
const maximizeWindow = () => {
  const handler = (window as any)?.mediaAssistant?.window?.maximize
  if (typeof handler === 'function') {
    handler()
  }
}

// 关闭窗口
const closeWindow = () => {
  window.close()
}
</script>

<style scoped>
.app-header {
  height: var(--header-height);
  background: var(--color-bg-header);
  border-bottom: 1px solid var(--color-border);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 16px;
  -webkit-app-region: drag;
  color: var(--color-text-primary);
  position: relative;
  z-index: 100;
}

.header-left,
.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
  -webkit-app-region: no-drag;
}

/* Logo */
.logo {
  display: flex;
  align-items: center;
  gap: 8px;
}

.logo-icon {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  background: var(--color-primary);
  border-radius: 4px;
}

.logo-icon svg {
  width: 18px;
  height: 18px;
}

.logo-text {
  font-size: 16px;
  font-weight: 800;
  color: var(--color-text-primary);
  letter-spacing: -0.01em;
}

/* 搜索框 */
.search-box {
  position: relative;
  width: 320px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.search-icon {
  position: absolute;
  left: 12px;
  top: 50%;
  transform: translateY(-50%);
  width: 14px;
  height: 14px;
  color: var(--color-text-tertiary);
  pointer-events: none;
  transition: color 0.2s;
  z-index: 2;
}

.search-input {
  width: 100%;
  height: 36px;
  padding: 0 12px 0 36px;
  background: var(--color-bg-surface);
  border: 1px solid var(--glass-border);
  border-radius: var(--radius-lg);
  color: var(--color-text-primary);
  font-size: 13px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: inset 0 2px 4px rgba(0, 0, 0, 0.1);
}

.search-input:hover {
  background: rgba(255, 255, 255, 0.08);
  border-color: var(--glass-border-bright);
}

.search-input:focus {
  outline: none;
  background: var(--color-bg-input);
  border-color: var(--color-primary);
}

.search-input:focus + .search-icon {
  color: var(--color-primary);
}

/* 右侧按钮 */
.header-btn {
  width: 34px;
  height: 34px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: 1px solid transparent;
  border-radius: var(--radius-md);
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: all 0.2s;
}

.header-btn:hover {
  background: var(--color-hover);
  color: var(--color-text-primary);
}

.header-btn svg {
  width: var(--icon-size-16);
  height: var(--icon-size-16);
}

.user-avatar {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: var(--color-bg-surface);
  border: 1px solid var(--color-border);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-text-secondary);
  cursor: pointer;
}

.user-avatar svg {
  width: var(--icon-size-16);
  height: var(--icon-size-16);
}

/* 窗口控制按钮 */
.window-controls {
  display: flex;
  align-items: center;
  gap: 8px;
  padding-left: var(--spacing-md);
  border-left: 1px solid var(--color-border);
}

.window-btn {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0;
  position: relative;
  overflow: hidden;
}

.window-btn svg {
  width: 8px;
  height: 8px;
  opacity: 0;
  color: rgba(0, 0, 0, 0.5);
  transition: opacity 0.2s;
}

.window-btn:hover svg {
  opacity: 1;
}

.window-btn.minimize { background: #febc2e; }
.window-btn.maximize { background: #28c840; }
.window-btn.close { background: #ff5f57; }

</style>
