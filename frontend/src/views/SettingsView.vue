<template>
  <div class="settings-page min-h-screen bg-base-100">
    <div class="settings-header h-12 flex items-center gap-3 px-4 border-b border-base-300 shrink-0">
      <button class="btn btn-sm btn-circle btn-ghost" @click="goBack">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M19 12H5M12 19l-7-7 7-7"/>
        </svg>
      </button>
      <h1 class="text-base font-semibold">设置</h1>
    </div>

    <div class="settings-content flex-1 overflow-y-auto p-4">
      <section class="mb-6">
        <h2 class="text-xs font-semibold text-base-content/50 uppercase tracking-wide mb-3 pb-2 border-b border-base-300">外观</h2>
        
        <div class="form-control">
          <label class="label cursor-pointer justify-between py-3 border-b border-base-200">
            <div>
              <span class="label-text text-sm">主题</span>
              <p class="text-xs text-base-content/50">选择应用的外观主题</p>
            </div>
            <select v-model="settings.theme" class="select select-sm select-bordered w-36">
              <option value="system">跟随系统</option>
              <option value="light">亮色</option>
              <option value="dark">暗色</option>
            </select>
          </label>
        </div>

        <div class="form-control">
          <label class="label cursor-pointer justify-between py-3 border-b border-base-200">
            <div>
              <span class="label-text text-sm">语言</span>
              <p class="text-xs text-base-content/50">界面显示语言</p>
            </div>
            <select v-model="settings.language" class="select select-sm select-bordered w-36">
              <option value="zh-CN">简体中文</option>
              <option value="en-US">English</option>
            </select>
          </label>
        </div>

        <div class="form-control">
          <label class="label cursor-pointer justify-between py-3 border-b border-base-200">
            <div>
              <span class="label-text text-sm">界面缩放</span>
              <p class="text-xs text-base-content/50">调整界面元素大小</p>
            </div>
            <select v-model="settings.scale" class="select select-sm select-bordered w-36">
              <option value="0.8">80%</option>
              <option value="0.9">90%</option>
              <option value="1">100%</option>
              <option value="1.1">110%</option>
              <option value="1.2">120%</option>
            </select>
          </label>
        </div>
      </section>

      <section class="mb-6">
        <h2 class="text-xs font-semibold text-base-content/50 uppercase tracking-wide mb-3 pb-2 border-b border-base-300">工作区</h2>
        
        <div class="form-control">
          <label class="label cursor-pointer justify-between py-3 border-b border-base-200">
            <div>
              <span class="label-text text-sm">监控目录深度</span>
              <p class="text-xs text-base-content/50">递归监控子目录的层级深度</p>
            </div>
            <select v-model="settings.watchDepth" class="select select-sm select-bordered w-36">
              <option value="1">1 级</option>
              <option value="2">2 级</option>
              <option value="3">3 级（免费版上限）</option>
              <option value="-1">无限制（Pro）</option>
            </select>
          </label>
        </div>

        <div class="form-control">
          <label class="label cursor-pointer justify-between py-3 border-b border-base-200">
            <div>
              <span class="label-text text-sm">文件监控并发</span>
              <p class="text-xs text-base-content/50">同时处理的文件变更数量</p>
            </div>
            <select v-model="settings.watchConcurrency" class="select select-sm select-bordered w-36">
              <option value="10">10 个</option>
              <option value="25">25 个</option>
              <option value="50">50 个（免费版上限）</option>
              <option value="100">100 个（Pro）</option>
              <option value="-1">无限制（Pro）</option>
            </select>
          </label>
        </div>

        <div class="form-control">
          <label class="label cursor-pointer justify-between py-3 border-b border-base-200">
            <div>
              <span class="label-text text-sm">索引文件数量上限</span>
              <p class="text-xs text-base-content/50">数据库中保存的最大文件记录数</p>
            </div>
            <select v-model="settings.maxIndexCount" class="select select-sm select-bordered w-36">
              <option value="1000">1,000 个（免费版上限）</option>
              <option value="5000">5,000 个（Pro）</option>
              <option value="10000">10,000 个（Pro）</option>
              <option value="-1">无限制（Pro）</option>
            </select>
          </label>
        </div>
      </section>

      <section class="mb-6">
        <h2 class="text-xs font-semibold text-base-content/50 uppercase tracking-wide mb-3 pb-2 border-b border-base-300">备份</h2>
        
        <div class="form-control">
          <label class="label cursor-pointer justify-between py-3 border-b border-base-200">
            <div>
              <span class="label-text text-sm">自动备份</span>
              <p class="text-xs text-base-content/50">定时自动备份工作区数据</p>
            </div>
            <input type="checkbox" v-model="settings.autoBackup" class="toggle toggle-primary">
          </label>
        </div>

        <div class="form-control" v-if="settings.autoBackup">
          <label class="label cursor-pointer justify-between py-3 border-b border-base-200">
            <div>
              <span class="label-text text-sm">备份频率</span>
              <p class="text-xs text-base-content/50">自动备份的时间间隔</p>
            </div>
            <select v-model="settings.backupFrequency" class="select select-sm select-bordered w-36">
              <option value="daily">每天</option>
              <option value="weekly">每周</option>
              <option value="monthly">每月</option>
            </select>
          </label>
        </div>

        <div class="form-control" v-if="settings.autoBackup">
          <label class="label cursor-pointer justify-between py-3 border-b border-base-200">
            <div>
              <span class="label-text text-sm">保留备份数量</span>
              <p class="text-xs text-base-content/50">自动清理旧备份，保留最近N个</p>
            </div>
            <select v-model="settings.backupRetention" class="select select-sm select-bordered w-36">
              <option value="3">3 个</option>
              <option value="5">5 个</option>
              <option value="10">10 个</option>
              <option value="-1">无限制</option>
            </select>
          </label>
        </div>
      </section>

      <section class="mb-6">
        <h2 class="text-xs font-semibold text-base-content/50 uppercase tracking-wide mb-3 pb-2 border-b border-base-300">性能</h2>
        
        <div class="form-control">
          <label class="label cursor-pointer justify-between py-3 border-b border-base-200">
            <div>
              <span class="label-text text-sm">缩略图缓存</span>
              <p class="text-xs text-base-content/50">缩略图缓存占用磁盘空间上限</p>
            </div>
            <select v-model="settings.thumbnailCacheSize" class="select select-sm select-bordered w-36">
              <option value="500">500 MB</option>
              <option value="1000">1 GB</option>
              <option value="2000">2 GB</option>
              <option value="5000">5 GB</option>
            </select>
          </label>
        </div>

        <div class="form-control">
          <label class="label cursor-pointer justify-between py-3 border-b border-base-200">
            <div>
              <span class="label-text text-sm">元数据提取并发</span>
              <p class="text-xs text-base-content/50">同时提取元数据的文件数量</p>
            </div>
            <select v-model="settings.metadataConcurrency" class="select select-sm select-bordered w-36">
              <option value="1">1 个</option>
              <option value="2">2 个（免费版上限）</option>
              <option value="4">4 个（Pro）</option>
              <option value="8">8 个（Pro）</option>
            </select>
          </label>
        </div>
      </section>

      <section class="mb-6">
        <h2 class="text-xs font-semibold text-base-content/50 uppercase tracking-wide mb-3 pb-2 border-b border-base-300">关于</h2>
        
        <div class="text-center py-6">
          <div class="text-xl font-bold mb-2">智归档OS</div>
          <div class="text-xs text-base-content/50 mb-1">版本 1.0.0</div>
          <div class="badge badge-primary badge-sm">免费版</div>
        </div>

        <div class="flex flex-col gap-2 px-4">
          <button class="btn btn-sm btn-ghost justify-start" @click="checkUpdate">
            检查更新
          </button>
          <button class="btn btn-sm btn-ghost justify-start" @click="openLicense">
            查看许可证
          </button>
          <button class="btn btn-sm btn-ghost justify-start" @click="clearCache">
            清除缓存
          </button>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const settings = reactive({
  theme: 'system',
  language: 'zh-CN',
  scale: '1',
  watchDepth: '3',
  watchConcurrency: '50',
  maxIndexCount: '1000',
  autoBackup: true,
  backupFrequency: 'daily',
  backupRetention: '5',
  thumbnailCacheSize: '1000',
  metadataConcurrency: '2'
})

function goBack() {
  router.push('/')
}

function checkUpdate() {
  alert('已是最新版本')
}

function openLicense() {
  alert('智归档OS v1.0\n\n本软件为免费版本，仅供个人使用。')
}

function clearCache() {
  if (confirm('确定要清除所有缓存吗？这将删除缩略图缓存和临时文件。')) {
    alert('缓存已清除')
  }
}
</script>

<style scoped>
.settings-page {
  display: flex;
  flex-direction: column;
  height: 100vh;
}
</style>
