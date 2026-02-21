<template>
  <div class="plugin-view">
    <MicroAppContainer :plugin-id="pluginId" />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import MicroAppContainer from '@/components/MicroAppContainer.vue'

const route = useRoute()
const pluginId = computed(() => route.params.id as string)
</script>

<style scoped>
.plugin-view {
  width: 100%;
  height: 100%;
  overflow: hidden;
}
</style>
