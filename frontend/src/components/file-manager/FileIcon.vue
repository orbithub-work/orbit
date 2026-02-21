<template>
  <span
    class="file-icon"
    :class="`file-icon--${fileType}`"
  >
    <svg
      class="icon-svg"
      viewBox="0 0 24 24"
      xmlns="http://www.w3.org/2000/svg"
      aria-hidden="true"
    >
      <use :href="symbolHref" />
    </svg>
  </span>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  fileType: string
  size?: number
}

const props = withDefaults(defineProps<Props>(), {
  size: 24
})

const symbolHref = computed(() => {
  switch (props.fileType) {
    case 'image':
      return '#i-image'
    case 'video':
      return '#i-video'
    case 'audio':
      return '#i-audio'
    case 'document':
      return '#i-file'
    case 'folder':
      return '#i-folder'
    default:
      return '#i-file'
  }
})
</script>

<style scoped>
.file-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: currentColor;
}

.icon-svg {
  width: v-bind('props.size + "px"');
  height: v-bind('props.size + "px"');
  color: currentColor;
}

.file-icon--image {
  color: var(--color-success);
}

.file-icon--video {
  color: var(--color-primary);
}

.file-icon--audio {
  color: var(--color-warning);
}

.file-icon--document {
  color: var(--color-text-secondary);
}

.file-icon--folder {
  color: var(--color-warning);
}
</style>
