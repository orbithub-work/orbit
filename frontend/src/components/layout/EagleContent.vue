<template>
  <div class="eagle-grid">
    <div class="grid-container">
      <div 
        v-for="i in 20" 
        :key="i" 
        class="grid-item"
        :class="{ selected: selectedIds.includes(i) }"
        @click="selectItem(i, $event)"
      >
        <div class="thumbnail-wrapper">
          <!-- Placeholder for image -->
          <div
            class="thumbnail-placeholder"
            :style="{ backgroundColor: getRandomColor(i) }"
          >
            <span class="file-type">JPG</span>
          </div>
          <div class="hover-overlay">
            <span class="checkbox"></span>
          </div>
        </div>
        <div class="item-info">
          <div class="item-name">
            Image_Asset_{{ i }}.jpg
          </div>
          <div class="item-meta">
            <span class="meta-tag">2.4 MB</span>
            <span class="meta-tag">4K</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const selectedIds = ref<number[]>([])

const selectItem = (id: number, event: MouseEvent) => {
  if (event.ctrlKey || event.metaKey) {
    if (selectedIds.value.includes(id)) {
      selectedIds.value = selectedIds.value.filter(i => i !== id)
    } else {
      selectedIds.value.push(id)
    }
  } else {
    selectedIds.value = [id]
  }
}

const getRandomColor = (i: number) => {
  const colors = ['#2a2a2a', '#333333', '#3a3a3a', '#252526']
  return colors[i % colors.length]
}
</script>

<style scoped>
.eagle-grid {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  background: var(--color-bg-base, #1e1e1e);
}

.grid-container {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
  gap: 16px;
}

.grid-item {
  display: flex;
  flex-direction: column;
  border-radius: 6px;
  overflow: hidden;
  cursor: pointer;
  transition: transform 0.1s;
  border: 2px solid transparent;
}

.grid-item:hover {
  transform: translateY(-2px);
}

.grid-item.selected {
  border-color: var(--color-primary, #007acc);
  background: var(--color-bg-surface-hover);
}

.thumbnail-wrapper {
  aspect-ratio: 4/3;
  background: #000;
  position: relative;
  border-radius: 4px;
  overflow: hidden;
}

.thumbnail-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #555;
  font-weight: bold;
}

.item-info {
  padding: 8px 4px;
}

.item-name {
  font-size: 13px;
  color: var(--color-text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-bottom: 4px;
}

.item-meta {
  display: flex;
  gap: 6px;
}

.meta-tag {
  font-size: 10px;
  color: var(--color-text-tertiary);
  background: rgba(255,255,255,0.05);
  padding: 1px 4px;
  border-radius: 2px;
}

/* Scrollbar */
.eagle-grid::-webkit-scrollbar {
  width: 8px;
}
.eagle-grid::-webkit-scrollbar-track {
  background: transparent;
}
.eagle-grid::-webkit-scrollbar-thumb {
  background: #3e3e42;
  border-radius: 4px;
}
</style>