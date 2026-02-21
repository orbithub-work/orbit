# Media Assistant 优化路线图

> 基于竞品分析，本文档列出了缩小与 Eagle/Billfish 差距的优化方案和实施计划。

---

## 📊 当前差距诊断

### 与 Eagle 对比
| 维度 | Eagle | Media Assistant | 差距 |
|------|-------|-----------------|------|
| 格式支持 | ✅ 全面 | ⚠️ 基础 | **需补齐** |
| 预览体验 | ✅ 优秀 | ⚠️ 基础 | **需优化** |
| 标签系统 | ✅ 完整 | ✅ 完整 | **已对等** |
| 性能表现 | ✅ 流畅 | ✅ 更优 | **领先** |
| 浏览器采集 | ✅ 有 | ❌ 无 | **需实现** |
| 资产血缘 | ❌ 无 | ✅ 有 | **独家优势** |
| 工程解析 | ❌ 无 | ⚠️ 规划中 | **差异化** |

### 与 Billfish 对比
| 维度 | Billfish | Media Assistant | 差距 |
|------|----------|-----------------|------|
| 性能 (10w+) | ❌ 卡顿 | ✅ 流畅 | **已超越** |
| 免费模式 | ✅ 免费 | ✅ 免费 | **对等** |
| 数据主权 | ❌ 私有 | ✅ 开放 | **领先** |
| AI 功能 | ❌ 无 | 🔄 规划中 | **潜在优势** |

---

## 🎯 3 个月优化计划

### 📅 Month 1: 补齐基本盘（60分 → 80分）

#### Week 1-2: 前端核心组件 ⭐⭐⭐ [已完成 ✅]
**优先级**: P0 - 用户第一眼体验

**已完成内容**:
- ✅ 虚拟滚动网格（支持 10 万+ 文件）
- ✅ 快速预览面板（Space 键触发）
- ✅ 批量操作工具栏
- ✅ 标签批量管理对话框
- ✅ 多选/框选交互
- ✅ 键盘快捷键支持

**交付成果**:
- 6 个核心组件（2000+ 行代码）
- 3 份完整文档
- 1 个示例页面

---

#### Week 3-4: 格式支持大补全 ⭐⭐⭐
**优先级**: P0 - 素材库的生命线  
**目标**: 支持创作者 90% 常用格式

##### 1. 视频格式增强（3 天）
**任务**:
- [ ] 集成 FFmpeg Go binding
  ```go
  import "github.com/u2takey/ffmpeg-go"
  
  // 视频缩略图生成
  func GenerateVideoThumbnail(videoPath string, timestamp float64) (string, error) {
      outputPath := getThumbnailPath(videoPath)
      err := ffmpeg.Input(videoPath, ffmpeg.KwArgs{
          "ss": timestamp,
      }).Output(outputPath, ffmpeg.KwArgs{
          "vframes": 1,
          "vf": "scale=320:-1",
      }).Run()
      return outputPath, err
  }
  ```

- [ ] 支持格式: MP4, MOV, AVI, MKV, WEBM, FLV
- [ ] 提取视频元数据（时长、分辨率、编码格式）
- [ ] 后台任务队列批量生成

**验收标准**:
- 视频缩略图生成时间 < 2 秒
- 支持至少 6 种主流视频格式

##### 2. 设计源文件（3 天）
**任务**:
- [ ] PSD 缩略图提取
  ```bash
  # 使用 ImageMagick 或 libvips
  convert "file.psd[0]" -thumbnail 320x320 output.jpg
  ```

- [ ] Figma/Sketch 文件识别（显示图标）
- [ ] AI 文件基础支持（Adobe Illustrator）

**技术选型**:
- PSD: ImageMagick 或 psd-tools
- AI: Ghostscript 或 Inkscape
- Sketch: 仅识别，暂不预览

##### 3. 专业格式（2 天）
**任务**:
- [ ] RAW 图片支持（CR2, NEF, ARW, DNG）
  ```bash
  # 使用 libraw
  dcraw -c -w -T input.cr2 | convert - -thumbnail 320x320 output.jpg
  ```

- [ ] EXR/TGA/HDR 支持（OpenImageIO）
- [ ] 字体文件预览（TTF, OTF, WOFF）

**API 设计**:
```go
// internal/services/format_service.go
type FormatService struct {
    converters map[string]FormatConverter
}

type FormatConverter interface {
    GenerateThumbnail(inputPath string, size int) (string, error)
    ExtractMetadata(inputPath string) (map[string]interface{}, error)
}
```

---

### 📅 Month 2: 降维打击免费市场（80分 → 90分）

#### Week 5-6: 性能与智能化 ⭐⭐
**优先级**: P1 - 碾压 Billfish 的关键  
**目标**: "10 万文件秒开，内存占用 < 500MB"

##### 1. 超大库性能优化（4 天）
**任务**:
- [ ] 实现分页加载 API
  ```go
  // GET /api/assets?page=1&limit=200&sort=name
  func (h *Handler) handleListAssets(w http.ResponseWriter, r *http.Request) {
      page := getQueryInt(r, "page", 1)
      limit := getQueryInt(r, "limit", 200)
      
      assets, total, err := h.assetService.ListAssetsPaginated(
          r.Context(), page, limit,
      )
      
      json.NewEncoder(w).Encode(map[string]interface{}{
          "assets": assets,
          "total": total,
          "page": page,
          "limit": limit,
      })
  }
  ```

- [ ] 缩略图缓存策略
  - LRU 内存缓存（最近 1000 张）
  - 磁盘缓存（WebP 格式压缩）
  - 预加载策略（滚动方向预测）

- [ ] 数据库查询优化
  ```sql
  -- 添加复合索引
  CREATE INDEX idx_assets_project_status ON assets(project_id, status);
  CREATE INDEX idx_assets_type_created ON assets(file_type, created_at DESC);
  
  -- 使用 EXPLAIN QUERY PLAN 分析慢查询
  ```

**性能目标**:
- 10 万文件初始加载: < 500ms
- 滚动加载 200 个文件: < 100ms
- 缩略图缓存命中率: > 80%

##### 2. 智能集合实装（3 天）
**任务**:
- [ ] 动态查询构建器
  ```go
  type SmartCollectionRule struct {
      Field    string      `json:"field"`     // "size", "type", "color", "date"
      Operator string      `json:"operator"`  // "gt", "lt", "eq", "contains"
      Value    interface{} `json:"value"`
  }
  
  type SmartCollection struct {
      ID    string                 `json:"id"`
      Name  string                 `json:"name"`
      Rules []SmartCollectionRule  `json:"rules"`
      Logic string                 `json:"logic"` // "AND", "OR"
  }
  ```

- [ ] 预设智能规则
  - "大于 100MB 的视频"
  - "最近 7 天修改的图片"
  - "红色主色调的设计稿"
  - "未打标签的文件"

- [ ] 自动更新机制（文件变更时刷新集合）

**UI 设计**:
```
智能集合创建对话框
├── 规则编辑器
│   ├── 字段选择 (下拉菜单)
│   ├── 操作符 (>, <, =, 包含)
│   └── 值输入 (根据字段类型动态)
├── 逻辑关系 (AND / OR)
└── 预览结果 (实时显示匹配数量)
```

---

#### Week 7-8: 颜色搜索 + 视觉检索 ⭐⭐
**优先级**: P1 - Eagle 的杀手锏  
**目标**: 实现基础的"按颜色找图"

##### 1. 颜色提取（2 天）
**任务**:
- [ ] 实现 K-means 聚类算法提取主色调
  ```go
  // internal/services/color_service.go
  func ExtractDominantColors(imagePath string, k int) ([]Color, error) {
      img, err := loadImage(imagePath)
      if err != nil {
          return nil, err
      }
      
      // 降采样加速
      img = resize.Thumbnail(100, 100, img, resize.Lanczos3)
      
      // K-means 聚类
      pixels := extractPixels(img)
      colors := kmeans(pixels, k, 20) // k个颜色，最多20次迭代
      
      return colors, nil
  }
  ```

- [ ] 存储到 `media_meta` JSON 字段
  ```json
  {
    "dominant_colors": [
      {"hex": "#FF5733", "percentage": 0.45},
      {"hex": "#33FF57", "percentage": 0.30},
      {"hex": "#3357FF", "percentage": 0.25}
    ]
  }
  ```

- [ ] 后台任务批量处理现有图片

##### 2. 颜色搜索 UI（3 天）
**任务**:
- [ ] 色盘选择器组件
  ```vue
  <ColorPicker
    v-model="selectedColor"
    :preset-colors="presetColors"
    @change="handleColorSearch"
  />
  ```

- [ ] 颜色相似度算法（HSV 空间距离）
  ```go
  func ColorDistance(c1, c2 Color) float64 {
      // 转换到 HSV 色彩空间
      h1, s1, v1 := RGBtoHSV(c1)
      h2, s2, v2 := RGBtoHSV(c2)
      
      // 计算距离（色相环形距离）
      dh := min(abs(h1-h2), 360-abs(h1-h2))
      ds := abs(s1 - s2)
      dv := abs(v1 - v2)
      
      return sqrt(dh*dh + ds*ds + dv*dv)
  }
  ```

- [ ] 搜索结果排序（相似度降序）

**API 设计**:
```
GET /api/assets/search-by-color?color=#FF5733&threshold=30
Response: [
  {
    "id": "asset-1",
    "dominant_colors": [...],
    "similarity": 0.95
  }
]
```

##### 3. 以图搜图（简化版）（2 天）
**任务**:
- [ ] 感知哈希（pHash）生成
  ```go
  func GeneratePerceptualHash(imagePath string) (string, error) {
      // 1. 缩放到 32x32
      // 2. 转灰度
      // 3. DCT 变换
      // 4. 提取低频系数
      // 5. 二值化生成 64 位哈希
      return hash, nil
  }
  ```

- [ ] 汉明距离相似度匹配
  ```go
  func HammingDistance(hash1, hash2 string) int {
      distance := 0
      for i := 0; i < len(hash1); i++ {
          if hash1[i] != hash2[i] {
              distance++
          }
      }
      return distance
  }
  ```

- [ ] 拖拽图片搜索界面

**性能优化**:
- pHash 存储到数据库字段
- 使用 BK-Tree 或 VP-Tree 加速相似搜索

---

### 📅 Month 3: 打出差异化王炸（90分 → 120分）

#### Week 9-10: 剪映/PR 工程解析 ⭐⭐⭐
**优先级**: P0 - 护城河功能  
**目标**: 创作者再也不怕工程红名

##### 1. 剪映工程解析（5 天）
**任务**:
- [ ] `.draft` 文件 JSON 解析
  ```go
  // internal/services/project_parser_service.go
  type JianyingProject struct {
      ID         string            `json:"id"`
      Name       string            `json:"name"`
      Duration   float64           `json:"duration"`
      Materials  []MaterialRef     `json:"materials"`
  }
  
  type MaterialRef struct {
      ID           string  `json:"id"`
      Type         string  `json:"type"`        // "video", "image", "audio"
      Path         string  `json:"path"`        // 原始路径
      InPoint      float64 `json:"in_point"`
      OutPoint     float64 `json:"out_point"`
  }
  
  func ParseJianyingProject(draftPath string) (*JianyingProject, error) {
      data, err := os.ReadFile(draftPath)
      if err != nil {
          return nil, err
      }
      
      var project JianyingProject
      err = json.Unmarshal(data, &project)
      return &project, err
  }
  ```

- [ ] 提取素材路径列表
- [ ] 检测断链（素材不存在）
  ```go
  type LinkStatus struct {
      MaterialID  string `json:"material_id"`
      Path        string `json:"path"`
      Exists      bool   `json:"exists"`
      AssetID     string `json:"asset_id,omitempty"` // 如果在库中
  }
  
  func CheckProjectLinks(project *JianyingProject) []LinkStatus {
      statuses := make([]LinkStatus, 0)
      for _, mat := range project.Materials {
          exists := fileExists(mat.Path)
          assetID, _ := findAssetByPath(mat.Path)
          
          statuses = append(statuses, LinkStatus{
              MaterialID: mat.ID,
              Path:       mat.Path,
              Exists:     exists,
              AssetID:    assetID,
          })
      }
      return statuses
  }
  ```

- [ ] 自动修复建议（基于指纹匹配）
  ```go
  func SuggestLinkFix(brokenPath string) []string {
      // 1. 提取文件名
      filename := filepath.Base(brokenPath)
      
      // 2. 按文件名搜索
      candidates := searchAssetsByName(filename)
      
      // 3. 如果有原文件指纹，按指纹匹配
      if fp := getFingerprint(brokenPath); fp != "" {
          candidates = filterByFingerprint(candidates, fp)
      }
      
      return extractPaths(candidates)
  }
  ```

##### 2. 血缘关系自动建立（2 天）
**任务**:
- [ ] 工程文件 → 素材自动关联
  ```go
  func LinkProjectToAssets(projectPath string) error {
      project, err := ParseJianyingProject(projectPath)
      if err != nil {
          return err
      }
      
      // 1. 创建或获取 Project 记录
      proj, err := projectService.CreateOrGet(project.Name, projectPath)
      
      // 2. 关联素材
      for _, mat := range project.Materials {
          asset, err := assetService.IndexFile(ctx, IndexFileRequest{
              Path:      mat.Path,
              ProjectID: proj.ID,
          })
          
          if err == nil {
              // 3. 建立血缘关系
              lineageService.Create(ctx, projectPath, mat.Path, "USES")
          }
      }
      
      return nil
  }
  ```

- [ ] UI 展示"这个素材被哪些工程用过"
  ```
  资产详情面板
  ├── 基本信息
  ├── 标签
  └── 📽️ 使用情况
      ├── 🎬 我的旅行 Vlog.draft (2023-12-01)
      ├── 🎬 年度总结视频.draft (2023-11-15)
      └── 🎬 产品宣传片_v2.draft (2023-10-20)
  ```

**API 设计**:
```
GET /api/assets/:id/usage
Response: {
  "projects": [
    {
      "id": "proj-1",
      "name": "我的旅行 Vlog",
      "path": "/path/to/project.draft",
      "last_used": "2023-12-01T10:00:00Z"
    }
  ],
  "total": 3
}
```

---

#### Week 11-12: 非侵入式优势强化 ⭐⭐
**优先级**: P1 - 强化核心卖点  
**目标**: "文件在哪儿，就在哪儿索引"

##### 1. 实时文件监听（3 天）
**任务**:
- [ ] 集成 fsnotify（已有依赖）
  ```go
  // internal/services/watcher_service.go
  func (w *WatcherService) WatchDirectory(path string) error {
      watcher, err := fsnotify.NewWatcher()
      if err != nil {
          return err
      }
      
      err = watcher.Add(path)
      if err != nil {
          return err
      }
      
      go func() {
          for {
              select {
              case event, ok := <-watcher.Events:
                  if !ok {
                      return
                  }
                  w.handleFileEvent(event)
                  
              case err, ok := <-watcher.Errors:
                  if !ok {
                      return
                  }
                  log.Error("Watcher error:", err)
              }
          }
      }()
      
      return nil
  }
  
  func (w *WatcherService) handleFileEvent(event fsnotify.Event) {
      switch {
      case event.Op&fsnotify.Create == fsnotify.Create:
          w.onFileCreated(event.Name)
      case event.Op&fsnotify.Write == fsnotify.Write:
          w.onFileModified(event.Name)
      case event.Op&fsnotify.Remove == fsnotify.Remove:
          w.onFileDeleted(event.Name)
      case event.Op&fsnotify.Rename == fsnotify.Rename:
          w.onFileRenamed(event.Name)
      }
  }
  ```

- [ ] 自动检测文件变更/移动/删除
- [ ] 增量更新索引（不全量扫描）

**性能优化**:
- 批量处理事件（防止短时间大量触发）
- 忽略临时文件（.tmp, .swp, ~）
- 可配置的监听规则

##### 2. XMP Sidecar 支持（2 天）
**任务**:
- [ ] 标签/评分写入 `.xmp` 文件
  ```go
  // internal/services/xmp_service.go
  func WriteXMP(assetPath string, metadata Metadata) error {
      xmpPath := assetPath + ".xmp"
      
      xmpContent := fmt.Sprintf(`<?xml version="1.0"?>
  <x:xmpmeta xmlns:x="adobe:ns:meta/">
    <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
      <rdf:Description rdf:about="">
        <dc:subject>
          <rdf:Bag>
            %s
          </rdf:Bag>
        </dc:subject>
        <xmp:Rating>%d</xmp:Rating>
      </rdf:Description>
    </rdf:RDF>
  </x:xmpmeta>`, generateTagsXML(metadata.Tags), metadata.Rating)
      
      return os.WriteFile(xmpPath, []byte(xmpContent), 0644)
  }
  ```

- [ ] 读取 XMP 元数据（导入时）
- [ ] 跨软件兼容（Bridge, Lightroom, Capture One）

**用户价值**:
- 数据不被软件锁死
- 随时可迁移到其他工具

##### 3. 项目统计仪表板（2 天）
**任务**:
- [ ] 文件数/总大小/类型分布统计
  ```go
  type ProjectStats struct {
      TotalFiles    int              `json:"total_files"`
      TotalSize     int64            `json:"total_size"`
      FileTypeStats map[string]int   `json:"file_type_stats"`
      TagStats      map[string]int   `json:"tag_stats"`
      TimelineStats []TimelinePoint  `json:"timeline_stats"`
  }
  
  func (s *ProjectService) GetProjectStats(ctx context.Context, projectID string) (*ProjectStats, error) {
      // SQL 聚合查询
      stats := &ProjectStats{
          FileTypeStats: make(map[string]int),
          TagStats:      make(map[string]int),
      }
      
      // 统计总数和大小
      row := s.db.QueryRow(`
          SELECT COUNT(*), SUM(size)
          FROM assets
          WHERE project_id = ?
      `, projectID)
      row.Scan(&stats.TotalFiles, &stats.TotalSize)
      
      // 按类型统计
      rows, _ := s.db.Query(`
          SELECT file_type, COUNT(*)
          FROM assets
          WHERE project_id = ?
          GROUP BY file_type
      `, projectID)
      // ...
      
      return stats, nil
  }
  ```

- [ ] 血缘关系可视化（依赖树）
  ```vue
  <template>
    <div class="lineage-graph">
      <vue-flow
        :nodes="lineageNodes"
        :edges="lineageEdges"
        @node-click="handleNodeClick"
      />
    </div>
  </template>
  ```

- [ ] ECharts 图表展示

**UI 设计**:
```
项目仪表板
├── 概览卡片
│   ├── 📊 2,345 个文件
│   ├── 💾 48.3 GB
│   └── 🏷️ 128 个标签
├── 文件类型分布 (饼图)
├── 时间线 (折线图)
└── 血缘关系图 (力导向图)
```

---

## 🚀 快速见效的"小胜利"（Quick Wins）

以下功能可以在 **1 周内完成**，且用户感知明显：

### 1️⃣ 批量操作增强（2 天）
**任务**:
- [ ] 批量归档到项目
  ```go
  func (s *AssetService) BatchArchive(ctx context.Context, projectID string, assetIDs []string) error {
      project, err := s.projects.Get(ctx, projectID)
      if err != nil {
          return err
      }
      
      for _, assetID := range assetIDs {
          asset, err := s.assets.GetByID(ctx, assetID)
          if err != nil {
              continue
          }
          
          // 复制文件到项目目录
          targetDir := filepath.Join(project.Path, "collected")
          os.MkdirAll(targetDir, 0755)
          
          targetPath := filepath.Join(targetDir, filepath.Base(asset.Path))
          utils.CopyFile(asset.Path, targetPath)
          
          // 建立关联
          s.projectAssets.Link(ctx, projectID, assetID)
      }
      
      return nil
  }
  ```

- [ ] 批量打标签（前端已完成）
- [ ] 批量移动文件

### 2️⃣ 搜索体验优化（1 天）
**任务**:
- [ ] 实时搜索建议（输入时展示）
  ```go
  // GET /api/assets/search-suggest?q=vac
  func (h *Handler) handleSearchSuggest(w http.ResponseWriter, r *http.Request) {
      query := r.URL.Query().Get("q")
      
      // 前缀匹配
      suggestions := []string{}
      
      // 1. 文件名匹配
      fileNames := h.assetService.SearchByPrefix(ctx, query, 5)
      suggestions = append(suggestions, fileNames...)
      
      // 2. 标签匹配
      tags := h.tagService.SearchByPrefix(ctx, query, 5)
      suggestions = append(suggestions, tags...)
      
      json.NewEncoder(w).Encode(suggestions)
  }
  ```

- [ ] 搜索历史记录
- [ ] 支持通配符（`*.psd`, `vacation_*`）

### 3️⃣ 拖拽采集增强（1 天）
**任务**:
- [ ] 从浏览器拖拽图片到应用
- [ ] 自动下载并索引
- [ ] 记录来源 URL（血缘关系）
  ```go
  // 前端发送
  {
    "url": "https://example.com/image.jpg",
    "source_page": "https://example.com/gallery",
    "project_id": "project-uuid"
  }
  
  // 后端处理
  func (h *Handler) handleImportFromURL(w http.ResponseWriter, r *http.Request) {
      var req ImportURLRequest
      json.NewDecoder(r.Body).Decode(&req)
      
      // 1. 下载文件
      localPath, err := downloadFile(req.URL)
      
      // 2. 索引文件
      asset, err := h.assetService.IndexFile(ctx, IndexFileRequest{
          Path:      localPath,
          ProjectID: req.ProjectID,
      })
      
      // 3. 记录来源（血缘关系）
      h.lineageService.Create(ctx, req.SourcePage, asset.ID, "DOWNLOADED_FROM")
      
      json.NewEncoder(w).Encode(asset)
  }
  ```

---

## 💡 关键技术选型

### 格式支持
| 需求 | 技术方案 | 依赖库 |
|------|---------|--------|
| 视频缩略图 | FFmpeg | github.com/u2takey/ffmpeg-go |
| 图片处理 | Go 原生 | github.com/disintegration/imaging |
| PSD 预览 | ImageMagick | 系统依赖 |
| RAW 图片 | dcraw | 系统依赖 |
| 字体预览 | FreeType | github.com/golang/freetype |

### 前端组件
| 需求 | 技术方案 | 依赖库 |
|------|---------|--------|
| 虚拟列表 | 自定义实现 | ✅ 已完成 |
| 颜色选择器 | 第三方组件 | @ckpack/vue-color |
| 文件拖拽 | Vue Draggable | vue-draggable |
| 图表 | ECharts | echarts (已有) |
| 血缘图 | 力导向图 | @vue-flow/core |

### 性能优化
| 需求 | 技术方案 | 优势 |
|------|---------|------|
| 缩略图缓存 | LRU + 磁盘 | 减少重复计算 |
| 数据库索引 | 复合索引 | 加速查询 |
| 并发处理 | Goroutine Pool | 充分利用多核 |
| 增量扫描 | 文件监听 | 减少全量扫描 |

---

## 📊 成功指标（KPI）

### 性能指标
| 指标 | 当前 | 目标 | 说明 |
|------|------|------|------|
| 10K 文件加载 | - | < 500ms | 初始启动速度 |
| 100K 文件滚动 | - | 60 FPS | 用户体验流畅度 |
| 内存占用 | - | < 500MB | 大库场景 |
| 缩略图生成 | - | < 2s | 视频/PSD |

### 功能覆盖
| 功能 | Eagle | Billfish | 我们 | 状态 |
|------|-------|----------|------|------|
| 格式支持 | 95% | 85% | 60% → **90%** | 📈 提升中 |
| 标签系统 | 100% | 100% | 100% | ✅ 已对等 |
| 批量操作 | 100% | 80% | 100% | ✅ 已完成 |
| 智能集合 | 80% | 50% | 0% → **80%** | 📈 待开发 |
| 工程解析 | 0% | 0% | 0% → **100%** | 🚀 独家 |

### 用户满意度
- **目标**: 从 60 分提升到 85 分
- **关键指标**:
  - 导入速度（10K 文件 < 5 分钟）
  - 搜索响应（< 100ms）
  - 崩溃率（< 0.1%）

---

## 🔄 迭代策略

### Phase 1: 基础对标（Month 1）
**目标**: 达到 Billfish 80% 功能覆盖  
**重点**: 格式支持、UI 体验

### Phase 2: 性能超越（Month 2）
**目标**: 性能指标全面超越 Billfish  
**重点**: 虚拟滚动、智能集合、颜色搜索

### Phase 3: 差异化突破（Month 3）
**目标**: 推出独家功能（工程解析、血缘追踪）  
**重点**: 剪映/PR 解析、实时监听、XMP 支持

---

## 📝 开发规范

### 代码质量
- ✅ 所有 API 必须有单元测试（覆盖率 > 70%）
- ✅ 所有错误必须有合理的错误处理
- ✅ 关键功能必须有性能测试
- ✅ 新功能必须更新文档

### Git 工作流
```bash
# 功能分支命名
feature/format-support-video
feature/color-search
feature/jianying-parser

# 提交消息格式
feat: 添加视频缩略图生成功能
fix: 修复大文件库卡顿问题
perf: 优化颜色搜索性能
docs: 更新 API 文档
```

### 性能测试
```bash
# 每周进行性能测试
cd tests
go test -bench=. -benchmem ./...

# 生成性能报告
go test -cpuprofile=cpu.prof -memprofile=mem.prof
go tool pprof cpu.prof
```

---

## 🐛 已知风险与对策

### 风险 1: FFmpeg 集成复杂度
**影响**: 视频缩略图生成可能延期  
**对策**: 
- 使用成熟的 Go binding 库
- 先支持 MP4/MOV，逐步扩展
- 提供降级方案（只显示图标）

### 风险 2: 性能优化效果不达标
**影响**: 大文件库仍然卡顿  
**对策**:
- 分阶段优化（先数据库，后缓存，最后 UI）
- 提供性能监控面板
- 用户可自定义缓存大小

### 风险 3: 工程解析兼容性
**影响**: 部分剪映版本无法解析  
**对策**:
- 支持多个剪映版本的格式
- 提供手动导入备选方案
- 记录解析失败的案例用于改进

---

## 📞 支持与反馈

### 开发团队联系
- 技术问题：查看各组件文档
- Bug 反馈：提交 GitHub Issue
- 功能建议：讨论区

### 文档索引
- 📄 [竞品策略](./COMPETITIVE_STRATEGY.md)
- 📄 [项目路线图](./ROADMAP.md)
- 📄 [P0 UI 组件](../frontend/P0_UI_COMPONENTS.md)
- 📄 [快速启动](../frontend/QUICK_START.md)

---

## 🎉 总结

### 当前进度
- ✅ **Month 1 Week 1-2**: P0 UI 组件（已完成）
- 🔄 **Month 1 Week 3-4**: 格式支持（进行中）
- 📅 **Month 2**: 性能与智能化（待开始）
- 📅 **Month 3**: 差异化功能（待开始）

### 关键里程碑
1. **Week 4**: 视频/PSD 支持上线
2. **Week 8**: 智能集合 + 颜色搜索上线
3. **Week 12**: 剪映工程解析上线

### 终极目标
**从 60 分基础产品 → 85 分专业工具 → 独家功能领先者**

🚀 **让我们一起缩小与竞品的差距，打造创作者的最佳资产管理工具！**
