<template>
  <div class="menu-item-wrapper">
    <component
      :is="item.path ? 'router-link' : 'div'"
      :to="item.path || ''"
      :class="[
        'menu-item',
        { 'active': isActive(item.path), 'has-children': item.children && item.children.length > 0 }
      ]"
      @click="handleItemClick"
    >
      <span
        v-if="item.icon"
        class="menu-icon"
      >
        <!-- In a real implementation, we would use actual icon components -->
        <span class="icon-placeholder">{{ getIcon(item.icon) }}</span>
      </span>
      <span
        v-if="!collapsed"
        class="menu-label"
      >{{ item.label }}</span>
      <span
        v-if="item.children && item.children.length > 0 && !collapsed"
        class="expand-icon"
        @click.stop="toggleExpanded"
      >
        {{ expanded ? '▼' : '▶' }}
      </span>
    </component>
    
    <div
      v-if="expanded && item.children && !collapsed"
      class="submenu"
    >
      <MenuItem
        v-for="child in item.children"
        :key="child.id"
        :item="child"
        :level="level + 1"
        :collapsed="collapsed"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'

interface MenuItem {
  id: string;
  label: string;
  icon?: string;
  path?: string;
  children?: MenuItem[];
  meta?: {
    title?: string;
    requiresAuth?: boolean;
  };
}

interface Props {
  item: MenuItem;
  level: number;
  collapsed: boolean;
}

const props = defineProps<Props>()
const route = useRoute()
const expanded = ref(false)

const isActive = (path?: string) => {
  return path && route.path.startsWith(path)
}

const toggleExpanded = () => {
  expanded.value = !expanded.value
}

const handleItemClick = () => {
  if (props.item.children && props.item.children.length > 0) {
    toggleExpanded()
  }
}

const getIcon = (iconName: string) => {
  // Return a character representing the icon
  // In a real implementation, we would use actual icon components
  switch(iconName) {
    case 'dashboard-icon':
      return '📊'
    case 'folder-icon':
      return '📁'
    case 'collection-icon':
      return '📚'
    case 'image-icon':
      return '🖼️'
    case 'video-icon':
      return '🎬'
    case 'audio-icon':
      return '🎵'
    case 'settings-icon':
      return '⚙️'
    case 'project-icon':
      return '📋'
    default:
      return '📦'
  }
}
</script>

<style scoped>
.menu-item-wrapper {
  margin-left: 0.5rem;
}

.menu-item {
  display: flex;
  align-items: center;
  padding: 0.75rem 1rem;
  margin: 0.25rem 0.5rem;
  border-radius: 6px;
  cursor: pointer;
  text-decoration: none;
  color: var(--color-text);
  transition: all 0.2s ease;
}

.menu-item:hover {
  background-color: var(--color-primary-100);
}

.menu-item.active {
  background-color: var(--color-primary);
  color: white;
}

.menu-item.has-children {
  justify-content: space-between;
}

.menu-icon {
  margin-right: 0.75rem;
  width: 20px;
  display: flex;
  align-items: center;
}

.icon-placeholder {
  font-size: 1.2rem;
}

.menu-label {
  flex: 1;
  white-space: nowrap;
}

.expand-icon {
  cursor: pointer;
  padding: 0.25rem;
  border-radius: 4px;
}

.expand-icon:hover {
  background-color: rgba(0, 0, 0, 0.1);
}

.submenu {
  margin-left: 2rem;
  border-left: 1px solid var(--color-border);
  padding-left: 0.5rem;
}

/* Adjust styles for nested levels */
.menu-item.level-1 { padding-left: 2rem; }
.menu-item.level-2 { padding-left: 3rem; }
.menu-item.level-3 { padding-left: 4rem; }
</style>