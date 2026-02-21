<template>
  <div
    class="file-thumbnail"
    :class="{ 'file-thumbnail--has-image': hasImage }"
  >
    <img
      v-if="hasImage && imageUrl"
      :src="imageUrl"
      :alt="fileName"
      class="thumbnail-image"
      @error="handleImageError"
    />
    <FileIcon
      v-else
      :file-type="fileType"
      :size="thumbnailSize"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import FileIcon from './FileIcon.vue'
import { apiCall } from '@/services/api'

interface Props {
  fileId: string
  fileName: string
  fileType: string
  thumbnailSize?: number
  existingThumbnail?: string | null
}

const props = withDefaults(defineProps<Props>(), {
  thumbnailSize: 150
})

const emit = defineEmits<{
  thumbnailLoaded: [url: string]
  thumbnailError: []
}>()

const imageUrl = ref<string | null>(props.existingThumbnail || null)
const loading = ref(false)
const hasImage = computed(() => ['image', 'video'].includes(props.fileType))

let abortController: AbortController | null = null

const loadThumbnail = async () => {
  if (!hasImage.value) return

  loading.value = true
  try {
    abortController = new AbortController()
    const thumbnailData = await apiCall<string>('get_thumbnail', {
      id: props.fileId,
      size: props.thumbnailSize
    })

    if (thumbnailData) {
      imageUrl.value = thumbnailData
      emit('thumbnailLoaded', thumbnailData)
    }
  } catch (error) {
    console.error('Failed to load thumbnail:', error)
    emit('thumbnailError')
  } finally {
    loading.value = false
    abortController = null
  }
}

const handleImageError = () => {
  imageUrl.value = null
  emit('thumbnailError')
}

onMounted(() => {
  if (hasImage.value && !imageUrl.value) {
    loadThumbnail()
  }
})

onUnmounted(() => {
  if (abortController) {
    abortController.abort()
  }
})
</script>

<style scoped>
.file-thumbnail {
  display: flex;
  align-items: center;
  justify-content: center;
  width: v-bind('props.thumbnailSize + "px"');
  height: v-bind('props.thumbnailSize + "px"');
  background-color: var(--color-surface);
  border-radius: 4px;
  overflow: hidden;
}

.file-thumbnail--has-image {
  background-color: #f5f5f5;
}

.thumbnail-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
</style>
