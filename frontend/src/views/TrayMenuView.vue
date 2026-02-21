<template>
  <div
    class="tray-overlay"
    @click.self="handleAction('hide-tray-menu')"
  >
    <div class="tray-menu">
      <div class="menu-list">
        <div
          class="menu-item"
          @click="handleAction('show-main')"
        >
          <i class="icon">💻</i>
          <span>显示主界面</span>
        </div>
      
        <div class="menu-divider"></div>
      
        <div
          class="menu-item"
          @click="handleAction('capture')"
        >
          <i class="icon">📸</i>
          <span>屏幕截图</span>
        </div>
        <div
          class="menu-item"
          @click="handleAction('import')"
        >
          <i class="icon">📥</i>
          <span>快速导入</span>
        </div>
      
        <div class="menu-divider"></div>
      
        <div
          class="menu-item danger"
          @click="handleAction('quit')"
        >
          <i class="icon">🚪</i>
          <span>退出</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'

const handleAction = (action: string) => {
  // Use the exposed electron API to send action back to main process
  if ((window as any).mediaAssistant?.tray?.sendAction) {
    (window as any).mediaAssistant.tray.sendAction(action)
  } else {
    console.warn('Tray API not available')
  }
}

onMounted(() => {
  // Adjust window size to fit content if possible, or just ensure it looks good
  document.body.style.backgroundColor = 'transparent'
})
</script>

<style scoped>
.tray-overlay {
  width: 100%;
  height: 100%;
  background: transparent;
  padding: 8px;
  box-sizing: border-box;
}

.tray-menu {
  width: 220px;
  /* Dark theme background similar to Windows 11 context menu */
  background: #2b2b2b;
  color: #ffffff;
  display: flex;
  flex-direction: column;
  user-select: none;
  border: 1px solid #454545;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.5);
  padding: 4px;
  box-sizing: border-box;
}

.menu-list {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.menu-item {
  display: flex;
  align-items: center;
  padding: 8px 12px;
  cursor: pointer;
  transition: background-color 0.1s;
  border-radius: 4px;
  font-size: 13px;
  font-family: 'Segoe UI', sans-serif;
}

.menu-item:hover {
  background-color: #414141;
}

.menu-item .icon {
  margin-right: 10px;
  font-style: normal;
  width: 16px;
  text-align: center;
  font-size: 14px;
}

.menu-divider {
  height: 1px;
  background-color: #454545;
  margin: 4px 0;
}

.menu-item.danger:hover {
  background-color: #c42b1c;
}
</style>
