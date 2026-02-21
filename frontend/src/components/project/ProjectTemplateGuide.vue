<template>
  <div class="project-template-guide">
    <div class="guide-header">
      <h2 class="guide-title">选择项目模板</h2>
      <p class="guide-subtitle">根据你的创作类型，选择合适的模板快速开始</p>
    </div>

    <div class="template-grid">
      <div
        v-for="template in templates"
        :key="template.id"
        class="template-card"
        :class="{ selected: selectedTemplate === template.id }"
        @click="selectedTemplate = template.id"
      >
        <div class="template-icon">{{ template.icon }}</div>
        <div class="template-content">
          <h3 class="template-name">{{ template.name }}</h3>
          <p class="template-desc">{{ template.description }}</p>
          <div class="template-features">
            <span v-for="feature in template.features" :key="feature" class="feature-tag">
              {{ feature }}
            </span>
          </div>
        </div>
        <div class="template-check">
          <Icon v-if="selectedTemplate === template.id" name="check-circle" size="lg" />
        </div>
      </div>
    </div>

    <div class="guide-actions">
      <button class="btn-secondary" @click="$emit('skip')">
        跳过，使用空白项目
      </button>
      <button class="btn-primary" :disabled="!selectedTemplate" @click="handleConfirm">
        使用此模板
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import Icon from '@/components/common/Icon.vue'

const emit = defineEmits<{
  confirm: [templateId: string]
  skip: []
}>()

const selectedTemplate = ref<string>('')

const templates = [
  {
    id: 'short-video',
    name: '短视频创作',
    icon: '🎬',
    description: '适合抖音、B站、小红书等短视频创作',
    features: ['素材管理', '剪辑工程', '多平台发布', '数据追踪']
  },
  {
    id: 'design',
    name: '设计项目',
    icon: '🎨',
    description: '适合平面设计、UI设计、品牌设计等',
    features: ['设计稿管理', '版本迭代', '客户交付', '素材复用']
  },
  {
    id: 'photography',
    name: '商单摄影',
    icon: '📷',
    description: '适合摄影师、修图师的商业拍摄项目',
    features: ['原片管理', '精修流程', '客户选片', '交付归档']
  },
  {
    id: 'game-dev',
    name: '游戏开发',
    icon: '🎮',
    description: '适合独立游戏开发、素材管理',
    features: ['资源管理', '版本控制', '构建发布', '迭代追踪']
  },
  {
    id: 'content-creation',
    name: '内容创作',
    icon: '✍️',
    description: '适合文章、教程、课程等内容创作',
    features: ['素材收集', '内容编写', '发布管理', '数据分析']
  },
  {
    id: 'blank',
    name: '空白项目',
    icon: '📋',
    description: '从零开始，自定义你的工作流程',
    features: ['完全自定义', '灵活配置', '适合任何场景']
  }
]

function handleConfirm() {
  if (selectedTemplate.value) {
    emit('confirm', selectedTemplate.value)
  }
}
</script>

<style scoped>
.project-template-guide {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: #1b1c1f;
  padding: 48px;
  overflow-y: auto;
}

.guide-header {
  text-align: center;
  margin-bottom: 48px;
}

.guide-title {
  font-size: 28px;
  font-weight: 600;
  color: #e5e7eb;
  margin: 0 0 12px 0;
}

.guide-subtitle {
  font-size: 15px;
  color: #9ca3af;
  margin: 0;
}

.template-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 20px;
  margin-bottom: 48px;
}

.template-card {
  position: relative;
  display: flex;
  gap: 16px;
  padding: 24px;
  background: #252526;
  border: 2px solid transparent;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.template-card:hover {
  background: #2a2a2a;
  border-color: rgba(59, 130, 246, 0.3);
  transform: translateY(-2px);
}

.template-card.selected {
  background: rgba(59, 130, 246, 0.1);
  border-color: #3b82f6;
}

.template-icon {
  font-size: 48px;
  flex-shrink: 0;
}

.template-content {
  flex: 1;
  min-width: 0;
}

.template-name {
  font-size: 16px;
  font-weight: 600;
  color: #e5e7eb;
  margin: 0 0 8px 0;
}

.template-desc {
  font-size: 13px;
  color: #9ca3af;
  margin: 0 0 12px 0;
  line-height: 1.5;
}

.template-features {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.feature-tag {
  padding: 3px 8px;
  background: rgba(255, 255, 255, 0.06);
  border-radius: 4px;
  font-size: 11px;
  color: #9ca3af;
}

.template-card.selected .feature-tag {
  background: rgba(59, 130, 246, 0.2);
  color: #60a5fa;
}

.template-check {
  position: absolute;
  top: 16px;
  right: 16px;
  color: #3b82f6;
}

.guide-actions {
  display: flex;
  justify-content: center;
  gap: 16px;
  padding-top: 24px;
  border-top: 1px solid #2b2b2f;
}

.btn-secondary,
.btn-primary {
  padding: 12px 32px;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-secondary {
  background: transparent;
  border: 1px solid rgba(255, 255, 255, 0.1);
  color: #9ca3af;
}

.btn-secondary:hover {
  background: rgba(255, 255, 255, 0.05);
  border-color: rgba(255, 255, 255, 0.2);
  color: #e5e7eb;
}

.btn-primary {
  background: #3b82f6;
  color: #fff;
}

.btn-primary:hover:not(:disabled) {
  background: #2563eb;
}

.btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
