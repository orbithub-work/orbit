<template>
  <div class="rating-control">
    <button
      v-for="star in 5"
      :key="star"
      class="star-btn"
      :class="{ active: star <= (currentRating || 0) }"
      @click="handleRatingClick(star)"
      :title="`${star} 星`"
    >
      <Icon :name="star <= (currentRating || 0) ? 'star-filled' : 'star'" size="sm" />
    </button>
    <button
      v-if="currentRating"
      class="clear-btn"
      @click="handleClearRating"
      title="清除评分"
    >
      <Icon name="x" size="sm" />
    </button>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import Icon from '@/components/common/Icon.vue'

interface Props {
  assetId: string
  rating?: number
}

interface Emits {
  (e: 'update', assetId: string, rating: number | null): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const currentRating = ref(props.rating || 0)

watch(() => props.rating, (newRating) => {
  currentRating.value = newRating || 0
})

const handleRatingClick = (star: number) => {
  currentRating.value = star
  emit('update', props.assetId, star)
}

const handleClearRating = () => {
  currentRating.value = 0
  emit('update', props.assetId, null)
}
</script>

<style scoped>
.rating-control {
  display: flex;
  align-items: center;
  gap: 0.25rem;
}

.star-btn,
.clear-btn {
  background: none;
  border: none;
  padding: 0.25rem;
  cursor: pointer;
  color: var(--color-text-secondary);
  transition: color 0.2s, transform 0.1s;
}

.star-btn:hover,
.clear-btn:hover {
  color: var(--color-primary);
  transform: scale(1.1);
}

.star-btn.active {
  color: var(--color-primary);
}

.clear-btn {
  margin-left: 0.25rem;
  opacity: 0.6;
}

.clear-btn:hover {
  opacity: 1;
}
</style>
