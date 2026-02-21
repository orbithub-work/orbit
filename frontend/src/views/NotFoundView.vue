<template>
  <div class="not-found-view">
    <h1>404 - 页面未找到</h1>
    <p>抱歉，您访问的页面不存在。</p>
    <router-link
      to="/"
      class="home-link"
    >
      返回首页
    </router-link>
  </div>
</template>

<script setup lang="ts">
// Not found view component
</script>

<style scoped>
.not-found-view {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100vh;
  text-align: center;
  padding: 1rem;
}

.home-link {
  color: var(--color-primary);
  text-decoration: none;
  margin-top: 1rem;
  padding: 0.5rem 1rem;
  border: 1px solid var(--color-border);
  border-radius: 4px;
}

.home-link:hover {
  background-color: var(--color-primary-100);
}
</style>