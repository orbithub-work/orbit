<template>
  <div class="main-layout">
    <header class="header">
      <slot name="header"></slot>
    </header>
    
    <div class="content">
      <slot></slot>
    </div>
    
    <footer class="footer">
      <slot name="footer"></slot>
    </footer>
  </div>
</template>

<script setup lang="ts">
// Main layout component that wraps the app content
// Includes header, main content area, and footer
</script>

<style scoped>
.main-layout {
  display: flex;
  flex-direction: column;
  height: 100%;
  min-height: 0;
}

.header {
  height: 60px;
  background-color: var(--color-surface);
  border-bottom: 1px solid var(--color-border);
  display: flex;
  align-items: center;
  padding: 0 1rem;
}

.content {
  flex: 1;
  display: flex;
  min-height: 0;
  overflow: hidden;
}

.footer {
  height: 60px;
  background-color: var(--color-surface);
  border-top: 1px solid var(--color-border);
  display: flex;
  align-items: center;
  padding: 0 1rem;
}
</style>
