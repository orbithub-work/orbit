<template>
  <div class="settings-view">
    <div class="settings-header">
      <h2 class="settings-title">
        设置
      </h2>
      <button
        class="close-btn"
        @click="$emit('close')"
      >
        <svg class="icon icon--20"><use href="#i-window-close" /></svg>
      </button>
    </div>

    <div class="settings-content">
      <div class="settings-section">
        <h3 class="section-label">
          界面主题
        </h3>
        <div class="theme-grid">
          <div 
            v-for="theme in themeStore.themes" 
            :key="theme.id"
            class="theme-card"
            :class="{ active: themeStore.currentThemeId === theme.id }"
            @click="themeStore.setTheme(theme.id)"
          >
            <div
              class="theme-preview"
              :style="{ background: theme.variables['--color-bg-base'] }"
            >
              <div
                class="preview-glow"
                :style="{ background: theme.glowColors[0] }"
              ></div>
              <div
                class="preview-accent"
                :style="{ background: theme.variables['--grad-primary'] }"
              ></div>
            </div>
            <div class="theme-info">
              <div class="theme-name">
                {{ theme.name }}
              </div>
              <div class="theme-desc">
                {{ theme.description }}
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="settings-section">
        <div class="setting-item">
          <div class="setting-info">
            <div class="setting-name">
              启动时自动扫描
            </div>
            <div class="setting-desc">
              应用启动时自动检查监控文件夹的变化
            </div>
          </div>
          <div class="setting-control">
            <input
              type="checkbox"
              checked
            />
          </div>
        </div>
      </div>

      <div class="settings-section">
        <h3 class="section-label">
          素材库
        </h3>
        <div
          class="setting-item clickable"
          @click="showTagManager = true"
        >
          <div class="setting-info">
            <div class="setting-name">
              标签管理
            </div>
            <div class="setting-desc">
              管理素材标签，创建、编辑或删除标签
            </div>
          </div>
          <div class="setting-control">
            <svg class="icon icon--16"><use href="#i-chevron-right" /></svg>
          </div>
        </div>
        <div class="setting-item">
          <div class="setting-info">
            <div class="setting-name">
              缩略图缓存
            </div>
            <div class="setting-desc">
              缓存生成的缩略图以提高浏览速度
            </div>
          </div>
          <div class="setting-control">
            <button class="btn btn-secondary">
              清理缓存
            </button>
          </div>
        </div>
      </div>

      <!-- 标签管理弹窗 -->
      <div
        v-if="showTagManager"
        class="modal-overlay"
        @click.self="showTagManager = false"
      >
        <div class="modal-content">
          <div class="modal-header">
            <h3>标签管理</h3>
            <button
              class="btn btn--icon btn--ghost"
              @click="showTagManager = false"
            >
              <svg class="icon icon--16"><use href="#i-window-close" /></svg>
            </button>
          </div>
          <div class="modal-body">
            <TagManager />
          </div>
        </div>
      </div>

      <div class="settings-section">
        <h3 class="section-label">
          关于
        </h3>
        <div class="setting-item">
          <div class="setting-info">
            <div class="setting-name">
              版本信息
            </div>
            <div class="setting-desc">
              智归档OS v0.1.0-alpha
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useThemeStore } from '@/stores/themeStore';
import TagManager from './TagManager.vue';

const themeStore = useThemeStore();
const showTagManager = ref(false);

defineEmits<{
  close: []
}>()
</script>

<style scoped>
.settings-view {
  height: 100%;
  background: var(--color-bg-sidebar);
  backdrop-filter: blur(var(--glass-blur));
  -webkit-backdrop-filter: blur(var(--glass-blur));
  display: flex;
  flex-direction: column;
}

.settings-header {
  padding: var(--spacing-xl);
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid var(--glass-border);
}

.settings-title {
  font-size: 24px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.close-btn {
  color: var(--color-text-tertiary);
  cursor: pointer;
  transition: color 0.2s;
}

.close-btn:hover {
  color: var(--color-text-primary);
}

.settings-content {
  flex: 1;
  padding: var(--spacing-xl);
  overflow-y: auto;
  max-width: 800px;
  margin: 0 auto;
  width: 100%;
}

.settings-section {
  margin-bottom: var(--spacing-2xl);
}

.section-label {
  font-size: 14px;
  font-weight: 600;
  color: var(--color-primary);
  margin-bottom: var(--spacing-lg);
  text-transform: uppercase;
  letter-spacing: 1px;
}

/* 主题网格 */
.theme-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
  gap: var(--spacing-md);
}

.theme-card {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid var(--glass-border);
  border-radius: var(--radius-lg);
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.theme-card:hover {
  background: rgba(255, 255, 255, 0.08);
  border-color: var(--glass-border-bright);
  transform: translateY(-2px);
}

.theme-card.active {
  border-color: var(--color-primary);
  box-shadow: 0 0 15px var(--color-primary-bg);
  background: rgba(255, 255, 255, 0.05);
}

.theme-preview {
  height: 80px;
  position: relative;
  overflow: hidden;
  background: #000;
}

.preview-glow {
  position: absolute;
  width: 100px;
  height: 100px;
  filter: blur(20px);
  border-radius: 50%;
  top: -20px;
  right: -20px;
  opacity: 0.5;
}

.preview-accent {
  position: absolute;
  bottom: 10px;
  left: 10px;
  width: 40px;
  height: 4px;
  border-radius: 2px;
}

.theme-info {
  padding: 12px;
}

.theme-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: 4px;
}

.theme-desc {
  font-size: 11px;
  color: var(--color-text-tertiary);
}

.setting-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--spacing-lg);
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid var(--glass-border);
  border-radius: var(--radius-md);
  margin-bottom: var(--spacing-md);
}

.setting-info {
  flex: 1;
}

.setting-name {
  font-size: 15px;
  font-weight: 500;
  color: var(--color-text-primary);
  margin-bottom: 2px;
}

.setting-desc {
  font-size: 13px;
  color: var(--color-text-tertiary);
}

.setting-control {
  margin-left: var(--spacing-xl);
}

.select-input {
  background: rgba(0, 0, 0, 0.2);
  border: 1px solid var(--glass-border);
  color: var(--color-text-primary);
  padding: 6px 12px;
  border-radius: var(--radius-sm);
  outline: none;
}

input[type="checkbox"] {
  width: 18px;
  height: 18px;
  cursor: pointer;
}

.clickable {
  cursor: pointer;
  transition: background 0.2s;
}

.clickable:hover {
  background: rgba(255, 255, 255, 0.05);
}

/* 弹窗样式 */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
  -webkit-backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  width: 90vw;
  max-width: 800px;
  max-height: 90vh;
  background: var(--color-bg-sidebar);
  border: 1px solid var(--glass-border);
  border-radius: var(--radius-lg);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.4);
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid var(--glass-border);
}

.modal-header h3 {
  font-size: 18px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
}

.modal-body {
  flex: 1;
  overflow: hidden;
  padding: 0;
}
</style>
