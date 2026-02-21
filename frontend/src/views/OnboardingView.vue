<template>
  <div class="onboarding-view">
    <div class="onboarding-container">
      <div class="logo-section">
        <div class="logo">🎬</div>
        <h1>欢迎使用智归档 OS</h1>
        <p class="tagline">专业创作者的素材管理工具</p>
      </div>

      <div class="options-section">
        <h2>选择素材库位置</h2>
        
        <div class="option-cards">
          <div 
            class="option-card recommended"
            :class="{ active: selectedOption === 'create' }"
            @click="selectedOption = 'create'"
          >
            <div class="card-header">
              <Icon name="folder-plus" size="xl" />
              <span class="badge">推荐</span>
            </div>
            <h3>新建统一素材库</h3>
            <p class="card-desc">自动创建标准目录结构，适合新手快速开始</p>
            <div v-if="selectedOption === 'create'" class="card-detail">
              <div class="path-preview">
                <Icon name="folder" size="sm" />
                <span>{{ suggestedPath }}</span>
              </div>
              <button class="btn-link" @click.stop="handleChangePath">更改位置</button>
            </div>
          </div>

          <div 
            class="option-card"
            :class="{ active: selectedOption === 'existing' }"
            @click="selectedOption = 'existing'"
          >
            <div class="card-header">
              <Icon name="folder-open" size="xl" />
            </div>
            <h3>选择已有目录</h3>
            <p class="card-desc">使用现有的素材文件夹，保持原有结构</p>
            <div v-if="selectedOption === 'existing'" class="card-detail">
              <button class="btn-select" @click.stop="handleSelectFolder">
                <Icon name="folder" size="sm" />
                {{ selectedPath || '点击选择文件夹' }}
              </button>
            </div>
          </div>

          <div 
            class="option-card"
            :class="{ active: selectedOption === 'skip' }"
            @click="selectedOption = 'skip'"
          >
            <div class="card-header">
              <Icon name="arrow-right" size="xl" />
            </div>
            <h3>稍后设置</h3>
            <p class="card-desc">先体验功能，之后再配置素材库</p>
          </div>
        </div>
      </div>

      <div class="actions">
        <button 
          class="btn-primary" 
          :disabled="submitting || !canProceed"
          @click="handleStart"
        >
          <span v-if="submitting">初始化中...</span>
          <span v-else>{{ actionButtonText }}</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { apiCall } from '@/services/api'
import Icon from '@/components/common/Icon.vue'

const router = useRouter()
const submitting = ref(false)
const selectedOption = ref<'create' | 'existing' | 'skip'>('create')
const selectedPath = ref('')
const suggestedPath = ref('')

onMounted(() => {
  suggestedPath.value = '~/Documents/我的创作素材库'
})

const canProceed = computed(() => {
  if (selectedOption.value === 'create') return true
  if (selectedOption.value === 'existing') return !!selectedPath.value
  if (selectedOption.value === 'skip') return true
  return false
})

const actionButtonText = computed(() => {
  if (selectedOption.value === 'create') return '创建并开始'
  if (selectedOption.value === 'existing') return '确认并开始'
  return '进入应用'
})

async function handleChangePath() {
  // 模拟文件选择
  const mockPath = prompt('输入路径（仅用于预览）:', suggestedPath.value)
  if (mockPath) suggestedPath.value = mockPath
}

async function handleSelectFolder() {
  // 模拟文件选择
  const mockPath = prompt('输入路径（仅用于预览）:', '~/Documents/MyAssets')
  if (mockPath) selectedPath.value = mockPath
}

async function handleStart() {
  submitting.value = true
  
  try {
    if (selectedOption.value === 'skip') {
      // 直接标记完成并跳转
      await apiCall('complete_first_launch')
      router.replace('/pool')
      return
    }

    const path = selectedOption.value === 'create' ? suggestedPath.value : selectedPath.value
    
    // 调用后端初始化接口
    await apiCall('run_onboarding', {
      mode: selectedOption.value,
      path: path
    })
    
    // 标记首次启动完成
    await apiCall('complete_first_launch')
    
    // 跳转到主页面
    router.replace('/pool')
  } catch (e) {
    console.error('Onboarding failed:', e)
    alert('初始化失败: ' + (e instanceof Error ? e.message : String(e)))
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.onboarding-view {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #1a1a1a 0%, #2d2d2d 100%);
  padding: 24px;
}

.onboarding-container {
  max-width: 800px;
  width: 100%;
}

.logo-section {
  text-align: center;
  margin-bottom: 48px;
}

.logo {
  font-size: 64px;
  margin-bottom: 16px;
}

.logo-section h1 {
  font-size: 32px;
  font-weight: 600;
  color: #e5e7eb;
  margin: 0 0 8px;
}

.tagline {
  font-size: 16px;
  color: #9ca3af;
  margin: 0;
}

.options-section h2 {
  font-size: 20px;
  font-weight: 600;
  color: #e5e7eb;
  text-align: center;
  margin-bottom: 32px;
}

.option-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 16px;
  margin-bottom: 32px;
}

.option-card {
  background: rgba(255, 255, 255, 0.05);
  border: 2px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  padding: 24px;
  cursor: pointer;
  transition: all 0.2s;
  position: relative;
}

.option-card:hover {
  background: rgba(255, 255, 255, 0.08);
  border-color: rgba(255, 255, 255, 0.2);
  transform: translateY(-2px);
}

.option-card.active {
  background: rgba(59, 130, 246, 0.1);
  border-color: #3b82f6;
}

.option-card.recommended::before {
  content: '';
  position: absolute;
  top: -2px;
  left: -2px;
  right: -2px;
  bottom: -2px;
  background: linear-gradient(135deg, #3b82f6, #8b5cf6);
  border-radius: 12px;
  z-index: -1;
  opacity: 0.3;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
  color: #60a5fa;
}

.badge {
  font-size: 11px;
  padding: 4px 8px;
  background: rgba(59, 130, 246, 0.2);
  border: 1px solid #3b82f6;
  border-radius: 12px;
  color: #60a5fa;
  font-weight: 500;
}

.option-card h3 {
  font-size: 18px;
  font-weight: 600;
  color: #e5e7eb;
  margin: 0 0 8px;
}

.card-desc {
  font-size: 14px;
  color: #9ca3af;
  line-height: 1.5;
  margin: 0;
}

.card-detail {
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.path-preview {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background: rgba(0, 0, 0, 0.2);
  border-radius: 6px;
  font-size: 13px;
  color: #d1d5db;
  margin-bottom: 8px;
  word-break: break-all;
}

.btn-link {
  background: none;
  border: none;
  color: #60a5fa;
  font-size: 13px;
  line-height: 1;
  cursor: pointer;
  padding: 4px 0;
  text-decoration: underline;
}

.btn-link:hover {
  color: #93c5fd;
}

.btn-select {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 10px 16px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 6px;
  color: #e5e7eb;
  font-size: 13px;
  line-height: 1;
  cursor: pointer;
  transition: all 0.15s;
}

.btn-select:hover {
  background: rgba(255, 255, 255, 0.08);
}

.actions {
  display: flex;
  justify-content: center;
}

.btn-primary {
  padding: 14px 48px;
  background: linear-gradient(135deg, #3b82f6, #2563eb);
  border: none;
  border-radius: 8px;
  color: white;
  font-size: 16px;
  font-weight: 500;
  line-height: 1;
  cursor: pointer;
  transition: all 0.2s;
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(59, 130, 246, 0.4);
}

.btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none;
}
</style>
