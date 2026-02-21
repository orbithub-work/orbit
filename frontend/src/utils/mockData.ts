import type { Asset, ImageAsset, VideoAsset, AudioAsset, FontAsset, ProjectAsset, AssetType } from '@/components/asset-card/index.vue'

const fileExtensions: Record<AssetType, string[]> = {
  image: ['jpg', 'jpeg', 'png', 'gif', 'webp', 'svg', 'heic', 'raw'],
  video: ['mp4', 'mov', 'avi', 'mkv', 'webm', 'flv', 'prores'],
  audio: ['mp3', 'wav', 'flac', 'aac', 'ogg', 'm4a', 'aiff'],
  font: ['ttf', 'otf', 'woff', 'woff2'],
  project: ['draft', 'prproj', 'aep', 'fcpxml', 'drp', 'cap', 'veg'],
  document: ['pdf', 'doc', 'docx', 'txt', 'md'],
  archive: ['zip', 'rar', '7z', 'tar', 'gz'],
  other: ['bin', 'dat', 'tmp']
}

const imageNames = [
  'sunset_beach', 'mountain_view', 'city_skyline', 'forest_path', 'ocean_waves',
  'desert_dunes', 'autumn_leaves', 'spring_flowers', 'winter_snow', 'summer_field',
  'portrait_studio', 'product_shot', 'food_photography', 'architecture', 'street_scene',
  'night_city', 'aurora_borealis', 'waterfall', 'clouds_sky', 'abstract_art',
  'AI_generated_portrait', 'AI_landscape', 'AI_abstract', 'AI_character', 'AI_concept'
]

const videoNames = [
  'intro_animation', 'product_demo', 'interview_clip', 'nature_documentary', 'music_video',
  'commercial_ad', 'tutorial_series', 'event_highlights', 'travel_vlog', 'behind_scenes',
  'motion_graphics', 'title_sequence', 'social_media_clip', 'corporate_video', 'wedding_film'
]

const audioNames = [
  'background_music', 'podcast_episode', 'voice_over', 'sound_effect', 'ambient_noise',
  'interview_audio', 'music_track', 'audio_book', 'meditation_guide', 'live_recording',
  'podcast_intro', 'notification_sound', 'ui_feedback', 'nature_sounds', 'electronic_beat'
]

const fontNames = [
  'Helvetica', 'Arial', 'Times_New_Roman', 'Georgia', 'Verdana',
  'Roboto', 'Open_Sans', 'Montserrat', 'Lato', 'Oswald',
  'Raleway', 'Poppins', 'Playfair_Display', 'Merriweather', 'Source_Sans'
]

const projectNames = [
  '品牌宣传片_v3', '产品发布会', '年度总结视频', '社交媒体内容包', '电商详情页设计',
  'APP界面设计', '网站重设计', '品牌VI系统', '营销活动素材', '产品包装设计',
  '纪录片项目', '音乐MV制作', '企业宣传片', '活动记录', '教学视频系列'
]

const artists = [
  '张三', '李四', '王五', '赵六', '钱七',
  'Studio_A', 'Design_Lab', 'Creative_House', 'Media_Pro', 'Art_Collective'
]

const clients = [
  '阿里巴巴', '腾讯', '字节跳动', '美团', '京东',
  '小米', '华为', 'OPPO', 'vivo', '网易'
]

function randomInt(min: number, max: number): number {
  return Math.floor(Math.random() * (max - min + 1)) + min
}

function randomFloat(min: number, max: number, decimals: number = 1): number {
  return parseFloat((Math.random() * (max - min) + min).toFixed(decimals))
}

function randomElement<T>(arr: T[]): T {
  return arr[randomInt(0, arr.length - 1)]
}

function randomType(): AssetType {
  const types: AssetType[] = ['image', 'image', 'image', 'video', 'audio', 'font', 'project', 'document', 'archive']
  return randomElement(types)
}

function generateId(): string {
  return `${Date.now()}-${Math.random().toString(36).substr(2, 9)}`
}

export function generateMockAsset(index: number): Asset {
  const type = randomType()
  const ext = randomElement(fileExtensions[type])
  const id = generateId()
  
  const baseAsset = {
    id,
    path: `/mock/path/${id}.${ext}`,
    size: randomInt(100000, 500000000),
    file_type: type
  }
  
  switch (type) {
    case 'image':
      return generateImageAsset(baseAsset, ext, index)
    case 'video':
      return generateVideoAsset(baseAsset, ext, index)
    case 'audio':
      return generateAudioAsset(baseAsset, ext, index)
    case 'font':
      return generateFontAsset(baseAsset, ext, index)
    case 'project':
      return generateProjectAsset(baseAsset, ext, index)
    default:
      return {
        ...baseAsset,
        name: `Asset_${index}.${ext}`
      }
  }
}

function generateImageAsset(base: any, ext: string, index: number): ImageAsset {
  const name = randomElement(imageNames)
  const isAI = name.startsWith('AI_')
  
  return {
    ...base,
    name: `${name}_${index}.${ext}`,
    width: randomInt(800, 8000),
    height: randomInt(600, 6000),
    is_ai_generated: isAI,
    prompt_id: isAI ? generateId() : undefined,
    rating: Math.random() > 0.7 ? randomInt(1, 5) : undefined,
    thumbnailUrl: `https://picsum.photos/seed/${index}/400/400`
  }
}

function generateVideoAsset(base: any, ext: string, index: number): VideoAsset {
  const name = randomElement(videoNames)
  const codecs = ['H.264', 'H.265', 'ProRes', 'DNxHD', 'AV1']
  
  return {
    ...base,
    name: `${name}_${index}.${ext}`,
    duration: randomFloat(10, 3600, 1),
    codec: randomElement(codecs),
    fps: randomElement([24, 25, 30, 60, 120]),
    width: randomElement([1920, 2560, 3840, 4096]),
    height: randomElement([1080, 1440, 2160, 2160]),
    thumbnailUrl: `https://picsum.photos/seed/${index + 1000}/400/225`
  }
}

function generateAudioAsset(base: any, ext: string, index: number): AudioAsset {
  const name = randomElement(audioNames)
  
  return {
    ...base,
    name: `${name}_${index}.${ext}`,
    duration: randomFloat(30, 600, 1),
    bitrate: randomElement([128, 192, 256, 320, 1411]),
    artist: randomElement(artists),
    album: `Album_${randomInt(1, 20)}`,
    sample_rate: randomElement([44100, 48000, 96000]),
    channels: randomElement([1, 2])
  }
}

function generateFontAsset(base: any, ext: string, index: number): FontAsset {
  const fontName = randomElement(fontNames)
  const styles = ['Regular', 'Bold', 'Italic', 'Light', 'Medium', 'Black']
  const weights = [100, 200, 300, 400, 500, 600, 700, 800, 900]
  
  return {
    ...base,
    name: `${fontName}-${randomElement(styles)}.${ext}`,
    font_name: fontName,
    font_style: randomElement(styles),
    font_weight: randomElement(weights),
    size: randomInt(20000, 5000000)
  }
}

function generateProjectAsset(base: any, ext: string, index: number): ProjectAsset {
  const name = randomElement(projectNames)
  
  return {
    ...base,
    name: `${name}.${ext}`,
    material_count: randomInt(5, 100),
    broken_links: Math.random() > 0.8 ? randomInt(1, 10) : 0,
    duration: randomFloat(60, 7200, 1),
    modified_at: new Date(Date.now() - randomInt(0, 30) * 24 * 60 * 60 * 1000).toISOString(),
    size: randomInt(1000000, 5000000000)
  }
}

export function generateMockAssets(count: number = 5000): Asset[] {
  const assets: Asset[] = []
  
  for (let i = 0; i < count; i++) {
    assets.push(generateMockAsset(i))
  }
  
  return assets
}

export function getAssetTypeDistribution(assets: Asset[]): Record<AssetType, number> {
  const distribution: Record<AssetType, number> = {
    image: 0,
    video: 0,
    audio: 0,
    font: 0,
    project: 0,
    document: 0,
    archive: 0,
    other: 0
  }
  
  assets.forEach(asset => {
    const type = getAssetTypeFromName(asset.name)
    distribution[type]++
  })
  
  return distribution
}

function getAssetTypeFromName(filename: string): AssetType {
  const ext = filename.split('.').pop()?.toLowerCase() || ''
  
  for (const [type, extensions] of Object.entries(fileExtensions)) {
    if (extensions.includes(ext)) {
      return type as AssetType
    }
  }
  
  return 'other'
}
