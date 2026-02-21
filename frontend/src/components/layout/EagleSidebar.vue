<template>
  <aside class="eagle-sidebar">
    <!-- Library Section -->
    <div class="sidebar-group">
      <div class="group-title">
        LIBRARY
      </div>
      
      <div 
        class="nav-item" 
        :class="{ active: activeItem === 'all' }"
        @click="$emit('select', 'all')"
      >
        <Icon name="image" size="sm" class="nav-icon" />
        <span class="label">All</span>
        <span class="count">1,240</span>
      </div>

      <div 
        class="nav-item"
        :class="{ active: activeItem === 'uncategorized' }"
        @click="$emit('select', 'uncategorized')"
      >
        <Icon name="folder" size="sm" class="nav-icon" />
        <span class="label">Uncategorized</span>
        <span class="count">45</span>
      </div>
    </div>

    <!-- Smart Folders -->
    <div class="sidebar-group">
      <div class="group-title">
        SMART FOLDERS
      </div>
      <div class="nav-item">
        <Icon name="clock" size="sm" class="nav-icon" />
        <span class="label">Last Imported</span>
      </div>
      <div class="nav-item">
        <Icon name="tag" size="sm" class="nav-icon" />
        <span class="label">Random</span>
      </div>
    </div>

    <!-- Folders (Scrollable) -->
    <div class="sidebar-group flex-grow">
      <div class="group-title">
        FOLDERS
      </div>
      <div class="folder-list">
        <div
          v-for="i in 5"
          :key="i"
          class="nav-item"
        >
          <Icon name="folder" size="sm" class="nav-icon" />
          <span class="label">Project Alpha {{ i }}</span>
        </div>
      </div>
    </div>

    <!-- Tags -->
    <div class="sidebar-group">
      <div class="group-title">
        TAGS
      </div>
      <div class="tags-cloud">
        <span class="tag-pill">design</span>
        <span class="tag-pill">ui</span>
        <span class="tag-pill">illustration</span>
      </div>
    </div>
    
    <!-- Trash -->
    <div class="sidebar-group mt-auto">
      <div class="nav-item">
        <Icon name="trash" size="sm" class="nav-icon" />
        <span class="label">Trash</span>
      </div>
    </div>
  </aside>
</template>

<script setup lang="ts">
import Icon from '@/components/common/Icon.vue'

defineProps<{
  activeItem?: string
}>()

defineEmits<{
  (e: 'select', id: string): void
}>()
</script>

<style scoped>
.eagle-sidebar {
  width: var(--sidebar-width, 260px);
  background: var(--color-bg-sidebar, #252526);
  border-right: 1px solid var(--color-border, #3e3e42);
  display: flex;
  flex-direction: column;
  height: 100%;
  user-select: none;
}

.sidebar-group {
  padding: 12px 0;
  border-bottom: 1px solid var(--color-divider, rgba(255,255,255,0.05));
}

.sidebar-group.flex-grow {
  flex: 1;
  overflow-y: auto;
  min-height: 0;
}

.sidebar-group.mt-auto {
  margin-top: auto;
  border-top: 1px solid var(--color-border);
  border-bottom: none;
}

.group-title {
  padding: 0 16px 8px;
  font-size: 11px;
  font-weight: 600;
  color: var(--color-text-tertiary, #666);
  letter-spacing: 0.5px;
}

.nav-item {
  display: flex;
  align-items: center;
  padding: 6px 16px;
  margin: 2px 8px;
  border-radius: 6px;
  cursor: pointer;
  color: var(--color-text-secondary, #999);
  transition: all 0.1s ease;
  font-size: 13px;
}

.nav-item:hover {
  background: var(--color-bg-surface-hover, rgba(255,255,255,0.08));
  color: var(--color-text-primary, #ccc);
}

.nav-item.active {
  background: var(--color-active-bg, rgba(0, 122, 204, 0.2));
  color: var(--color-primary, #007acc);
}

.nav-item .icon {
  margin-right: 10px;
  font-size: 14px;
  width: 16px;
  text-align: center;
}

.nav-item .label {
  flex: 1;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.nav-item .count {
  font-size: 11px;
  color: var(--color-text-tertiary, #666);
}

.folder-list {
  padding-bottom: 8px;
}

.tags-cloud {
  padding: 0 16px;
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.tag-pill {
  background: var(--color-bg-surface-hover);
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 11px;
  color: var(--color-text-secondary);
  cursor: pointer;
}

.tag-pill:hover {
  color: var(--color-text-primary);
  background: var(--color-bg-input);
}
</style>