<template>
  <header class="eagle-header">
    <div class="header-left">
      <div class="nav-controls">
        <button
          class="icon-btn"
          title="Back"
        >
          <Icon name="arrow-left" size="sm" />
        </button>
        <button
          class="icon-btn"
          title="Forward"
          disabled
        >
          <Icon name="arrow-right" size="sm" />
        </button>
      </div>
      <div class="breadcrumb">
        Library / All
      </div>
    </div>

    <div class="header-center">
      <div class="search-bar">
        <Icon name="search" size="sm" class="search-icon" />
        <input
          type="text"
          placeholder="Search..."
        />
      </div>
    </div>

    <div class="header-right">
      <div class="view-controls">
        <button
          class="icon-btn active"
          title="Grid View"
        >
          <Icon name="grid" size="sm" />
        </button>
        <button
          class="icon-btn"
          title="Waterfall View"
        >
          <Icon name="table" size="sm" />
        </button>
        <button
          class="icon-btn"
          title="List View"
        >
          <Icon name="list" size="sm" />
        </button>
      </div>
      <div class="zoom-slider">
        <Icon name="minus" size="sm" />
        <input
          type="range"
          min="50"
          max="300"
          value="150"
        />
        <Icon name="plus" size="sm" />
      </div>
      <button
        class="icon-btn primary"
        title="Filter"
      >
        <Icon name="filter" size="sm" /> Filter
      </button>
    </div>
  </header>
</template>

<script setup lang="ts">
import Icon from '@/components/common/Icon.vue'
</script>

<style scoped>
.eagle-header {
  height: var(--header-height, 56px);
  background: var(--color-bg-header, #2d2d2d);
  border-bottom: 1px solid var(--color-border, #3e3e42);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 16px;
  -webkit-app-region: drag; /* Electron drag */
}

.header-left, .header-center, .header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-left {
  flex: 1;
}

.header-center {
  flex: 2;
  justify-content: center;
}

.header-right {
  flex: 1;
  justify-content: flex-end;
  -webkit-app-region: no-drag;
}

.nav-controls {
  display: flex;
  gap: 4px;
  -webkit-app-region: no-drag;
}

.breadcrumb {
  font-size: 13px;
  color: var(--color-text-secondary);
  font-weight: 500;
}

.search-bar {
  background: var(--color-bg-input, #3c3c3c);
  border-radius: 6px;
  padding: 6px 12px;
  display: flex;
  align-items: center;
  width: 100%;
  max-width: 400px;
  -webkit-app-region: no-drag;
  border: 1px solid transparent;
  transition: all 0.2s;
}

.search-bar:focus-within {
  border-color: var(--color-primary);
  background: var(--color-bg-base);
}

.search-bar input {
  background: transparent;
  border: none;
  color: var(--color-text-primary);
  margin-left: 8px;
  width: 100%;
  outline: none;
  font-size: 13px;
}

.icon-btn {
  background: transparent;
  border: none;
  color: var(--color-text-secondary);
  padding: 6px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.icon-btn:hover:not(:disabled) {
  background: var(--color-bg-surface-hover);
  color: var(--color-text-primary);
}

.icon-btn.active {
  background: var(--color-bg-surface-hover);
  color: var(--color-primary);
}

.icon-btn:disabled {
  opacity: 0.3;
  cursor: not-allowed;
}

.zoom-slider {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 10px;
  color: var(--color-text-tertiary);
  -webkit-app-region: no-drag;
}

input[type=range] {
  width: 80px;
  height: 4px;
  border-radius: 2px;
  background: #4d4d4d;
  outline: none;
  appearance: none;
}
</style>