import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Artifact, PlatformPublish } from '@/types/artifact'

export const useArtifactStore = defineStore('artifact', () => {
  const artifacts = ref<Artifact[]>([])
  const loading = ref(false)
  const currentArtifact = ref<Artifact | null>(null)

  const draftArtifacts = computed(() => 
    artifacts.value.filter(a => a.status === 'draft')
  )

  const publishedArtifacts = computed(() => 
    artifacts.value.filter(a => a.status === 'published')
  )

  const archivedArtifacts = computed(() => 
    artifacts.value.filter(a => a.status === 'archived')
  )

  function loadMockData() {
    artifacts.value = [
      {
        id: '1',
        name: '春季穿搭分享',
        project_id: '1',
        file_path: '/artifacts/spring-fashion.jpg',
        thumbnail: '',
        status: 'published',
        platforms: [
          { platform: 'xiaohongshu', url: 'https://xiaohongshu.com/xxx', published_at: '2024-01-15' }
        ],
        version: 2,
        tags: ['穿搭', '春季'],
        created_at: '2024-01-10',
        updated_at: '2024-01-15'
      },
      {
        id: '2',
        name: '产品测评视频',
        project_id: '2',
        file_path: '/artifacts/review.mp4',
        thumbnail: '',
        status: 'draft',
        platforms: [],
        version: 1,
        tags: ['测评', '数码'],
        created_at: '2024-01-12',
        updated_at: '2024-01-12'
      },
      {
        id: '3',
        name: '旅行Vlog',
        project_id: '3',
        file_path: '/artifacts/travel.mp4',
        thumbnail: '',
        status: 'published',
        platforms: [
          { platform: 'douyin', url: 'https://douyin.com/xxx', published_at: '2024-01-14' },
          { platform: 'bilibili', url: 'https://bilibili.com/xxx', published_at: '2024-01-14' }
        ],
        version: 3,
        tags: ['旅行', 'Vlog'],
        created_at: '2024-01-08',
        updated_at: '2024-01-14'
      },
      {
        id: '4',
        name: '美食探店',
        project_id: '1',
        file_path: '/artifacts/food.jpg',
        thumbnail: '',
        status: 'archived',
        platforms: [
          { platform: 'xiaohongshu', url: 'https://xiaohongshu.com/yyy', published_at: '2023-12-20' }
        ],
        version: 1,
        tags: ['美食', '探店'],
        created_at: '2023-12-18',
        updated_at: '2023-12-20'
      },
    ]
  }

  function getArtifactsByProject(projectId: string) {
    return artifacts.value.filter(a => a.project_id === projectId)
  }

  function getArtifactsByStatus(status: Artifact['status']) {
    return artifacts.value.filter(a => a.status === status)
  }

  function getArtifactsByPlatform(platformId: string) {
    return artifacts.value.filter(a => 
      a.platforms.some(p => p.platform === platformId)
    )
  }

  function addArtifact(artifact: Omit<Artifact, 'id' | 'created_at' | 'updated_at'>) {
    const newArtifact: Artifact = {
      ...artifact,
      id: Date.now().toString(),
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString()
    }
    artifacts.value.push(newArtifact)
    return newArtifact
  }

  function updateArtifact(id: string, updates: Partial<Artifact>) {
    const index = artifacts.value.findIndex(a => a.id === id)
    if (index !== -1) {
      artifacts.value[index] = {
        ...artifacts.value[index],
        ...updates,
        updated_at: new Date().toISOString()
      }
    }
  }

  function deleteArtifact(id: string) {
    const index = artifacts.value.findIndex(a => a.id === id)
    if (index !== -1) {
      artifacts.value.splice(index, 1)
    }
  }

  function addPlatform(id: string, platform: PlatformPublish) {
    const artifact = artifacts.value.find(a => a.id === id)
    if (artifact) {
      artifact.platforms.push(platform)
      artifact.updated_at = new Date().toISOString()
    }
  }

  return {
    artifacts,
    loading,
    currentArtifact,
    draftArtifacts,
    publishedArtifacts,
    archivedArtifacts,
    loadMockData,
    getArtifactsByProject,
    getArtifactsByStatus,
    getArtifactsByPlatform,
    addArtifact,
    updateArtifact,
    deleteArtifact,
    addPlatform
  }
})
