<template>
  <div class="project-view-skeleton">
    <!-- 项目概览骨架 -->
    <div class="project-overview-skeleton">
      <div class="overview-header-skeleton">
        <Skeleton
          variant="text"
          class="project-title-skeleton"
        />
        <div class="status-badge-wrapper">
          <Skeleton
            variant="rounded"
            class="status-badge-skeleton"
          />
        </div>
      </div>
      
      <div class="overview-stats-skeleton">
        <div
          v-for="i in 3"
          :key="i"
          class="stat-item-skeleton"
        >
          <Skeleton
            variant="text"
            class="stat-value-skeleton"
          />
          <Skeleton
            variant="text"
            class="stat-label-skeleton"
          />
        </div>
      </div>

      <!-- 进度条骨架 -->
      <div class="progress-section-skeleton">
        <div class="progress-header-skeleton">
          <Skeleton
            variant="text"
            class="progress-text-skeleton"
          />
          <Skeleton
            variant="text"
            class="progress-percent-skeleton"
          />
        </div>
        <div class="progress-bar-wrapper">
          <Skeleton
            variant="rounded"
            class="progress-bar-skeleton"
          />
        </div>
      </div>
    </div>

    <!-- 子 Tab 骨架 -->
    <div class="sub-tabs-skeleton">
      <div
        v-for="i in 2"
        :key="i"
        class="sub-tab-wrapper"
      >
        <Skeleton
          variant="rounded"
          class="sub-tab-skeleton"
        />
      </div>
    </div>

    <!-- 内容区域骨架 -->
    <div class="tab-content-skeleton">
      <div
        v-for="i in 4"
        :key="i"
        class="folder-item-skeleton"
      >
        <div class="folder-header-skeleton">
          <div class="folder-icon-wrapper">
            <Skeleton
              variant="rounded"
              class="folder-icon-skeleton"
            />
          </div>
          <Skeleton
            variant="text"
            class="folder-name-skeleton"
          />
          <Skeleton
            variant="text"
            class="folder-count-skeleton"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import Skeleton from './Skeleton.vue'
</script>

<style scoped>
.project-view-skeleton {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

/* 项目概览骨架 */
.project-overview-skeleton {
  padding: 20px;
  background: rgba(255, 255, 255, 0.015);
  border-bottom: 1px solid rgba(255, 255, 255, 0.04);
}

.overview-header-skeleton {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 20px;
}

.project-title-skeleton {
  width: 140px;
  height: 22px;
}

.status-badge-wrapper {
  padding: 2px;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.02);
}

.status-badge-skeleton {
  width: 56px;
  height: 22px;
}

.overview-stats-skeleton {
  display: flex;
  gap: 32px;
  margin-bottom: 20px;
}

.stat-item-skeleton {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.stat-value-skeleton {
  width: 36px;
  height: 18px;
}

.stat-label-skeleton {
  width: 52px;
  height: 11px;
}

/* 进度条骨架 */
.progress-section-skeleton {
  margin-top: 16px;
}

.progress-header-skeleton {
  display: flex;
  justify-content: space-between;
  margin-bottom: 10px;
}

.progress-text-skeleton {
  width: 60px;
  height: 12px;
}

.progress-percent-skeleton {
  width: 32px;
  height: 12px;
}

.progress-bar-wrapper {
  padding: 3px;
  border-radius: 10px;
  background: rgba(255, 255, 255, 0.02);
}

.progress-bar-skeleton {
  width: 100%;
  height: 8px;
}

/* 子 Tab 骨架 */
.sub-tabs-skeleton {
  display: flex;
  gap: 8px;
  padding: 12px 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.04);
  background: rgba(255, 255, 255, 0.01);
}

.sub-tab-wrapper {
  padding: 2px;
  border-radius: 10px;
  background: rgba(255, 255, 255, 0.02);
}

.sub-tab-skeleton {
  width: 80px;
  height: 30px;
}

/* 内容区域骨架 */
.tab-content-skeleton {
  flex: 1;
  padding: 16px 20px;
  overflow-y: auto;
}

.folder-item-skeleton {
  margin-bottom: 10px;
  padding: 3px;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.015);
  border: 1px solid rgba(255, 255, 255, 0.03);
}

.folder-header-skeleton {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 14px;
}

.folder-icon-wrapper {
  padding: 2px;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.02);
}

.folder-icon-skeleton {
  width: 20px;
  height: 20px;
}

.folder-name-skeleton {
  flex: 1;
  width: 100px;
  height: 14px;
}

.folder-count-skeleton {
  width: 28px;
  height: 12px;
}
</style>
