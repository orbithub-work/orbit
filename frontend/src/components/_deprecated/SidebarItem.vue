<template>
  <div class="sidebar-item-container">
    <div
      class="sidebar-item"
      :class="{ 
        'active': isActive,
        'is-expanded': isExpanded,
        'has-children': hasChildren
      }"
      :style="{ paddingLeft: paddingLeft }"
      @click="handleClick"
    >
      <!-- 展开/折叠箭头 -->
      <button 
        v-if="hasChildren"
        class="toggle-btn"
        @click.stop="toggleExpand"
      >
        <svg 
          class="icon icon--12 arrow-icon" 
          :class="{ 'rotated': isExpanded }"
          viewBox="0 0 24 24"
        >
          <path
            d="M9 18l6-6-6-6"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          />
        </svg>
      </button>
      <span
        v-else
        class="toggle-placeholder"
      ></span>

      <!-- 图标 -->
      <svg
        v-if="icon"
        class="item-icon icon icon--16"
        viewBox="0 0 24 24"
        aria-hidden="true"
      >
        <use :href="icon" />
      </svg>
      <slot
        v-else
        name="icon"
      ></slot>

      <!-- 文本 -->
      <span class="item-label truncate">{{ label }}</span>

      <!-- 计数 -->
      <span
        v-if="count !== undefined"
        class="item-count"
      >{{ count }}</span>
      
      <!-- 更多操作 (Hover显示) -->
      <div class="item-actions">
        <slot name="actions"></slot>
      </div>
    </div>

    <!-- 子级递归 -->
    <div
      v-if="isExpanded && hasChildren"
      class="sidebar-sub-menu"
    >
      <SidebarItem
        v-for="child in children"
        :id="child.id"
        :key="child.id"
        :label="child.label"
        :icon="child.icon"
        :count="child.count"
        :children="child.children"
        :level="level + 1"
        :active-id="activeId"
        @select="$emit('select', $event)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'

interface Item {
  id: string
  label: string
  icon?: string
  count?: number
  children?: Item[]
}

const props = withDefaults(defineProps<{
  id: string
  label: string
  icon?: string
  count?: number
  children?: Item[]
  level?: number
  activeId?: string
  defaultExpanded?: boolean
}>(), {
  level: 0,
  defaultExpanded: false,
  children: () => []
})

const emit = defineEmits<{
  select: [id: string]
}>()

const isExpanded = ref(props.defaultExpanded)

const hasChildren = computed(() => props.children && props.children.length > 0)
const isActive = computed(() => props.activeId === props.id)
const paddingLeft = computed(() => `${(props.level * 12) + 8}px`) // Base padding 8px + level indent

const toggleExpand = () => {
  isExpanded.value = !isExpanded.value
}

const handleClick = () => {
  emit('select', props.id)
}
</script>

<style scoped>
.sidebar-item-container {
  display: flex;
  flex-direction: column;
}

.sidebar-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 12px;
  margin: 1px 8px;
  border-radius: 4px;
  cursor: pointer;
  color: var(--color-text-secondary);
  transition: all 0.1s ease-out;
  position: relative;
  overflow: hidden;
  font-size: 13px;
  font-weight: 500;
  border: 1px solid transparent;
}

.sidebar-item:hover {
  background: var(--color-hover);
  color: var(--color-text-primary);
}

.sidebar-item.active {
  background: var(--color-active-bg);
  color: var(--color-primary);
}

/* Eagle style: Active indicator bar or icon color */
.sidebar-item.active .item-icon {
  color: var(--color-primary);
}

.sidebar-item.active .item-count {
  background: var(--color-hover);
  color: var(--color-text-primary);
}

.item-icon {
  width: 16px;
  height: 16px;
  flex-shrink: 0;
  opacity: 0.8;
  transition: color 0.2s;
}

.sidebar-item:hover .item-icon {
  opacity: 1;
}

.item-label {
  flex: 1;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  line-height: 1.5;
}

.item-count {
  font-size: 11px;
  background: var(--color-hover);
  padding: 1px 6px;
  border-radius: 10px;
  color: var(--color-text-tertiary);
  transition: all 0.2s;
}

.toggle-btn {
  background: none;
  border: none;
  padding: 0;
  cursor: pointer;
  color: var(--color-text-tertiary);
  display: flex;
  align-items: center;
  justify-content: center;
  width: 16px;
  height: 16px;
  margin-right: -4px;
  transition: color 0.2s;
}

.toggle-btn:hover {
  color: var(--color-text-primary);
}

.arrow-icon {
  width: 10px;
  height: 10px;
  transition: transform 0.2s;
}

.arrow-icon.rotated {
  transform: rotate(90deg);
}

.toggle-placeholder {
  width: 12px;
  display: inline-block;
}

.sidebar-sub-menu {
  display: flex;
  flex-direction: column;
}
</style>
