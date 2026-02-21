export interface PlatformPublish {
  platform: string
  url?: string
  published_at: string
  note?: string
}

export interface Artifact {
  id: string
  name: string
  project_id: string
  file_path: string
  thumbnail?: string
  status: 'draft' | 'published' | 'archived'
  platforms: PlatformPublish[]
  version: number
  tags: string[]
  note?: string
  created_at: string
  updated_at: string
}

export const PLATFORMS = [
  { id: 'xiaohongshu', name: '小红书', icon: '📕' },
  { id: 'douyin', name: '抖音', icon: '🎵' },
  { id: 'wechat', name: '公众号', icon: '💬' },
  { id: 'bilibili', name: 'B站', icon: '📺' },
  { id: 'weibo', name: '微博', icon: '📝' },
  { id: 'zhihu', name: '知乎', icon: '💡' },
  { id: 'youtube', name: 'YouTube', icon: '▶️' },
  { id: 'other', name: '其他', icon: '🔗' },
] as const

export type PlatformId = typeof PLATFORMS[number]['id']
