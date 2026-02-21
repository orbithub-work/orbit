<template>
  <aside
    class="sidebar"
    style="padding: 20px; color: #888;"
  >
    <h3>Sidebar (Safe Mode)</h3>
    <p>Content disabled for debugging.</p>
  </aside>
</template>

<script setup lang="ts">
defineProps<{
  currentTab: string
}>()
</script>

<style scoped>
.sidebar {
  width: 250px;
  background: #252526;
  height: 100%;
}
</style>
