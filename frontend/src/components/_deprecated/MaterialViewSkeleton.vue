<template>
  <div class="material-grid-skeleton">
    <div
      v-for="i in 12"
      :key="i"
      class="card-skeleton"
    >
      <div class="thumb-skeleton"></div>
      <div class="meta-skeleton">
        <div class="text-line title"></div>
        <div class="text-line detail"></div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.material-grid-skeleton {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: var(--spacing-md);
  padding-bottom: 40px;
}

.card-skeleton {
  background: var(--color-bg-surface);
  border-radius: var(--radius-lg);
  overflow: hidden;
  border: 1px solid transparent;
}

.thumb-skeleton {
  width: 100%;
  aspect-ratio: 4/3;
  background: var(--color-bg-surface-hover);
  animation: pulse 1.5s infinite;
}

.meta-skeleton {
  padding: 10px 12px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.text-line {
  height: 10px;
  border-radius: 2px;
  background: var(--color-bg-surface-hover);
  animation: pulse 1.5s infinite;
}

.text-line.title {
  width: 70%;
}

.text-line.detail {
  width: 40%;
}

@keyframes pulse {
  0% { opacity: 0.6; }
  50% { opacity: 0.3; }
  100% { opacity: 0.6; }
}
</style>
