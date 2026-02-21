<template>
  <aside class="sidebar">
    <!-- Logo -->
    <div class="sidebar-header">
      <div class="logo">
        <div class="logo-icon">
          <svg
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
          >
            <rect
              x="3"
              y="3"
              width="18"
              height="18"
              rx="2"
              ry="2"
            />
            <circle
              cx="8.5"
              cy="8.5"
              r="1.5"
            />
            <polyline points="21 15 16 10 5 21" />
          </svg>
        </div>
        <span class="logo-text">智归档</span>
      </div>
    </div>

    <!-- 主要导航 -->
    <nav class="nav-section">
      <div class="nav-item active">
        <svg
          class="nav-icon"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
        >
          <rect
            x="3"
            y="3"
            width="7"
            height="7"
          />
          <rect
            x="14"
            y="3"
            width="7"
            height="7"
          />
          <rect
            x="14"
            y="14"
            width="7"
            height="7"
          />
          <rect
            x="3"
            y="14"
            width="7"
            height="7"
          />
        </svg>
        <span class="nav-label">全部</span>
        <span class="nav-count">19</span>
      </div>
      <div class="nav-item">
        <svg
          class="nav-icon"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
        >
          <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z" />
        </svg>
        <span class="nav-label">未分类</span>
      </div>
      <div class="nav-item">
        <svg
          class="nav-icon"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
        >
          <path d="M20.59 13.41l-7.17 7.17a2 2 0 0 1-2.83 0L2 12V2h10l8.59 8.59a2 2 0 0 1 0 2.82z" />
          <line
            x1="7"
            y1="7"
            x2="7.01"
            y2="7"
          />
        </svg>
        <span class="nav-label">未标签</span>
        <span class="nav-count">19</span>
      </div>
      <div class="nav-item">
        <svg
          class="nav-icon"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
        >
          <circle
            cx="12"
            cy="12"
            r="10"
          />
          <polyline points="12 6 12 12 16 14" />
        </svg>
        <span class="nav-label">最近使用</span>
      </div>
      <div class="nav-item">
        <svg
          class="nav-icon"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
        >
          <line
            x1="4"
            y1="9"
            x2="20"
            y2="9"
          />
          <line
            x1="4"
            y1="15"
            x2="20"
            y2="15"
          />
          <line
            x1="10"
            y1="3"
            x2="8"
            y2="21"
          />
          <line
            x1="16"
            y1="3"
            x2="14"
            y2="21"
          />
        </svg>
        <span class="nav-label">随机模式</span>
      </div>
      <div class="nav-item">
        <svg
          class="nav-icon"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
        >
          <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2" />
          <circle
            cx="9"
            cy="7"
            r="4"
          />
          <path d="M23 21v-2a4 4 0 0 0-3-3.87" />
          <path d="M16 3.13a4 4 0 0 1 0 7.75" />
        </svg>
        <span class="nav-label">资源社区</span>
      </div>
      <div class="nav-item">
        <svg
          class="nav-icon"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
        >
          <path d="M20.59 13.41l-7.17 7.17a2 2 0 0 1-2.83 0L2 12V2h10l8.59 8.59a2 2 0 0 1 0 2.82z" />
          <line
            x1="7"
            y1="7"
            x2="7.01"
            y2="7"
          />
        </svg>
        <span class="nav-label">标签管理</span>
      </div>
      <div class="nav-item">
        <svg
          class="nav-icon"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
        >
          <polyline points="3 6 5 6 21 6" />
          <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" />
        </svg>
        <span class="nav-label">回收站</span>
      </div>
    </nav>

    <!-- 智能文件夹 -->
    <div class="nav-group">
      <div class="nav-group-title">
        智能文件夹 (1)
      </div>
      <div class="nav-item active">
        <svg
          class="nav-icon"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
        >
          <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z" />
        </svg>
        <span class="nav-label">我问问</span>
      </div>
    </div>

    <!-- 文件夹 -->
    <div class="nav-group nav-group-flex">
      <div class="nav-group-title">
        文件夹 (267)
      </div>
      <div class="folder-list">
        <div class="nav-item folder-item">
          <svg
            class="nav-icon"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
          >
            <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z" />
          </svg>
          <span class="nav-label">本地资源库</span>
          <span class="nav-count">19</span>
        </div>
      </div>
    </div>

    <!-- 底部工具栏 -->
    <div class="sidebar-footer">
      <button class="footer-btn">
        <svg
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
        >
          <polygon points="22 3 2 3 10 12.46 10 19 14 21 14 12.46 22 3" />
        </svg>
        <span>筛选</span>
      </button>
      <button
        class="footer-btn"
        @click="$router.push('/settings')"
      >
        <svg
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
        >
          <circle
            cx="12"
            cy="12"
            r="3"
          />
          <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z" />
        </svg>
      </button>
      <button class="footer-btn">
        <svg
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
        >
          <circle
            cx="12"
            cy="12"
            r="10"
          />
          <path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3" />
          <line
            x1="12"
            y1="17"
            x2="12.01"
            y2="17"
          />
        </svg>
      </button>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()
</script>

<style scoped>
.sidebar {
  width: 240px;
  height: 100vh;
  background: #1a1a1a;
  color: #a0a0a0;
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
  border-right: 1px solid var(--color-border);
}

.sidebar-header {
  padding: 12px 16px;
  border-bottom: 1px solid var(--color-border);
}

.logo {
  display: flex;
  align-items: center;
  gap: 8px;
}

.logo-icon {
  width: 28px;
  height: 28px;
  color: #fff;
}

.logo-icon svg {
  width: 100%;
  height: 100%;
}

.logo-text {
  font-size: 16px;
  font-weight: 600;
  color: #fff;
}

/* 导航区域 */
.nav-section {
  padding: 8px 0;
}

.nav-group {
  padding: 8px 0;
}

.nav-group-flex {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.nav-group-title {
  padding: 4px 16px;
  font-size: 11px;
  color: #666;
  font-weight: 500;
}

.folder-list {
  flex: 1;
  overflow-y: auto;
}

.nav-item {
  display: flex;
  align-items: center;
  padding: 6px 12px;
  margin: 0 8px;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.15s;
  font-size: 13px;
  color: #a0a0a0;
}

.nav-item:hover {
  background: rgba(255, 255, 255, 0.05);
  color: #e0e0e0;
}

.nav-item.active {
  background: rgba(255, 255, 255, 0.08);
  color: #fff;
}

.nav-icon {
  width: 16px;
  height: 16px;
  margin-right: 10px;
  flex-shrink: 0;
}

.nav-icon svg {
  width: 100%;
  height: 100%;
}

.nav-label {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.nav-count {
  font-size: 11px;
  color: #666;
  margin-left: 4px;
}

.nav-item.active .nav-count {
  color: #888;
}

/* 底部工具栏 */
.sidebar-footer {
  display: flex;
  align-items: center;
  padding: 8px 12px;
  border-top: 1px solid var(--color-border);
  gap: 4px;
}

.footer-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 10px;
  background: transparent;
  border: none;
  border-radius: 4px;
  color: #888;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.15s;
}

.footer-btn:hover {
  background: rgba(255, 255, 255, 0.05);
  color: #e0e0e0;
}

.footer-btn svg {
  width: 14px;
  height: 14px;
}

.footer-btn:last-child {
  margin-left: auto;
  padding: 6px;
}

/* Scrollbar styling */
.folder-list::-webkit-scrollbar {
  width: 4px;
}

.folder-list::-webkit-scrollbar-track {
  background: transparent;
}

.folder-list::-webkit-scrollbar-thumb {
  background: #333;
  border-radius: 2px;
}

.folder-list::-webkit-scrollbar-thumb:hover {
  background: #444;
}
</style>
