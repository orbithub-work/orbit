# 素材库迁移策略设计

> 创建日期: 2026-02-17
> 状态: 设计阶段
> 目标: 支持从 Eagle、Billfish 等主流素材管理工具无缝迁移

## 1. 背景分析

### 1.1 用户痛点

| 痛点 | 描述 |
|------|------|
| 迁移成本高 | 用户已在其他工具积累了大量素材和标签体系 |
| 数据丢失风险 | 标签、分类、评分等元数据无法保留 |
| 学习成本 | 用户需要重新适应新的管理方式 |
| 时间成本 | 重新整理素材需要大量时间 |

### 1.2 目标工具分析

| 工具 | 市场份额 | 存储方式 | 迁移难度 |
|------|----------|----------|----------|
| **Eagle** | 高（付费 199 元） | 侵占式，随机ID目录，metadata.json | 中等 |
| **Billfish** | 高（免费） | 三种模式（复制/剪切/索引），.bfmeta | 低 |
| **Pixave** | 中（Mac） | 数据库存储 | 高 |
| **Inboard** | 低 | 数据库存储 | 高 |

> **Billfish 导入模式说明**：
> - **复制模式（推荐）**：素材复制到素材库目录，保留原位置文件
> - **剪切模式**：素材移动到素材库目录，原位置文件删除
> - **索引模式**：不移动文件，仅建立索引引用

## 2. Eagle 数据结构分析

### 2.1 库目录结构

```
MyLibrary.library/
├── images/
│   ├── FLK2N823JH23/              # 随机 ID 命名
│   │   ├── 海报设计.jpg           # 实际文件
│   │   └── metadata.json          # 元数据
│   ├── KJH234HJ23K/
│   │   ├── 背景.png
│   │   └── metadata.json
│   └── ...
├── folders.json                   # 文件夹结构
├── tags.json                      # 标签列表
└── library.json                   # 库配置
```

### 2.2 metadata.json 结构

```json
{
  "id": "FLK2N823JH23",
  "name": "海报设计.jpg",
  "size": 1024000,
  "btime": 1609459200000,
  "mtime": 1609459200000,
  "tags": ["海报", "设计", "蓝色"],
  "folders": ["FOLDER_ID_1", "FOLDER_ID_2"],
  "annotation": "这是一个很好的海报设计参考",
  "url": "https://example.com/source",
  "star": 5,
  "isDeleted": false,
  "lastModified": 1609459200000
}
```

### 2.3 folders.json 结构

```json
[
  {
    "id": "FOLDER_ID_1",
    "name": "设计素材",
    "children": [
      {
        "id": "FOLDER_ID_2",
        "name": "海报"
      }
    ]
  }
]
```

### 2.4 tags.json 结构

```json
[
  {
    "id": "TAG_ID_1",
    "name": "海报",
    "color": "#3498db"
  },
  {
    "id": "TAG_ID_2",
    "name": "设计",
    "color": "#e74c3c"
  }
]
```

## 3. Billfish 数据结构分析

### 3.1 导入模式对比

| 模式 | 文件位置 | 适用场景 | 磁盘占用 |
|------|----------|----------|----------|
| **复制** | 复制到素材库目录 | 素材备份，原位置不变 | 双倍 |
| **剪切** | 移动到素材库目录 | 整理零散素材 | 单份 |
| **索引** | 原位置不变 | 跨磁盘大量素材 | 无额外占用 |

### 3.2 库目录结构

**复制/剪切模式：**
```
MyLibrary/                        # 素材库目录
├── 海报/
│   ├── 海报设计.jpg
│   ├── 海报设计.jpg.bfmeta      # 元数据文件
│   ├── 背景.png
│   └── 背景.png.bfmeta
├── 图标/
│   └── ...
└── .billfish/                    # 库配置目录
    ├── library.json
    └── ...
```

**索引模式：**
```
# 素材保持原位置，.bfmeta 文件也在原位置
D:/素材/海报设计.jpg
D:/素材/海报设计.jpg.bfmeta
E:/图片/背景.png
E:/图片/背景.png.bfmeta
```

### 3.3 .bfmeta 结构（推测）

```json
{
  "tags": ["海报", "设计"],
  "rating": 5,
  "annotation": "备注信息",
  "source": "https://example.com",
  "createTime": 1609459200000,
  "modifyTime": 1609459200000
}
```

### 3.4 特点

- **多种导入模式**：复制、剪切、索引三种方式
- **元数据独立**：每个文件有对应的 .bfmeta 文件
- **文件夹结构保留**：复制/剪切模式下保持用户创建的文件夹结构
- **兼容性好**：即使删除软件，文件仍可正常使用

## 4. 迁移策略设计

### 4.1 自动检测机制

```go
type LibraryType string

const (
    LibraryTypeEagle    LibraryType = "eagle"
    LibraryTypeBillfish LibraryType = "billfish"
    LibraryTypeUnknown  LibraryType = "unknown"
)

func DetectLibraryType(path string) LibraryType {
    // Eagle: 检测 *.library 目录
    if strings.HasSuffix(path, ".library") {
        if _, err := os.Stat(filepath.Join(path, "images")); err == nil {
            if _, err := os.Stat(filepath.Join(path, "folders.json")); err == nil {
                return LibraryTypeEagle
            }
        }
    }
    
    // Billfish: 检测 .billfish 目录
    if _, err := os.Stat(filepath.Join(path, ".billfish")); err == nil {
        return LibraryTypeBillfish
    }
    
    // 检测 .bfmeta 文件
    if hasBfmetaFiles(path) {
        return LibraryTypeBillfish
    }
    
    return LibraryTypeUnknown
}
```

### 4.2 迁移流程

```
用户添加监听目录
        ↓
  自动检测库类型
        ↓
  ┌─────┴─────┐
  ↓           ↓
 Eagle      Billfish
  ↓           ↓
解析元数据   解析元数据
  ↓           ↓
  └─────┬─────┘
        ↓
  映射到本地数据模型
        ↓
  导入标签/文件夹/评分
        ↓
  建立素材索引
        ↓
    完成
```

### 4.3 数据映射关系

| Eagle/Billfish | Smart Archive | 说明 |
|----------------|---------------|------|
| tags | tags 表 | 标签直接导入 |
| folders | collections 表 | 文件夹映射为收藏集 |
| star/rating | rating 字段 | 评分保留 |
| annotation | notes 字段 | 备注保留 |
| url | source_url 字段 | 来源链接保留 |
| btime/mtime | created_at/updated_at | 时间保留 |

### 4.4 Eagle 迁移实现

```go
type EagleImporter struct {
    libraryPath string
    assetRepo   *repos.AssetRepo
    tagRepo     *repos.TagRepo
    collectionRepo *repos.CollectionRepo
}

type EagleMetadata struct {
    ID         string   `json:"id"`
    Name       string   `json:"name"`
    Size       int64    `json:"size"`
    Btime      int64    `json:"btime"`
    Mtime      int64    `json:"mtime"`
    Tags       []string `json:"tags"`
    Folders    []string `json:"folders"`
    Annotation string   `json:"annotation"`
    URL        string   `json:"url"`
    Star       int      `json:"star"`
}

func (i *EagleImporter) Import(ctx context.Context) (*ImportResult, error) {
    // 1. 解析文件夹结构
    folders, err := i.parseFolders()
    if err != nil {
        return nil, err
    }
    
    // 2. 解析标签
    tags, err := i.parseTags()
    if err != nil {
        return nil, err
    }
    
    // 3. 创建收藏集（文件夹映射）
    collectionMap := i.createCollections(ctx, folders)
    
    // 4. 创建标签
    tagMap := i.createTags(ctx, tags)
    
    // 5. 导入素材
    imagesDir := filepath.Join(i.libraryPath, "images")
    entries, _ := os.ReadDir(imagesDir)
    
    var imported, skipped int
    for _, entry := range entries {
        if !entry.IsDir() {
            continue
        }
        
        metaPath := filepath.Join(imagesDir, entry.Name(), "metadata.json")
        meta, err := i.parseMetadata(metaPath)
        if err != nil {
            skipped++
            continue
        }
        
        // 查找实际文件
        filePath := i.findActualFile(filepath.Join(imagesDir, entry.Name()), meta.Name)
        if filePath == "" {
            skipped++
            continue
        }
        
        // 创建资产
        asset, err := i.createAsset(ctx, filePath, meta)
        if err != nil {
            skipped++
            continue
        }
        
        // 关联标签
        for _, tagName := range meta.Tags {
            if tagID, ok := tagMap[tagName]; ok {
                i.tagRepo.LinkTag(ctx, asset.ID, tagID)
            }
        }
        
        // 关联收藏集
        for _, folderID := range meta.Folders {
            if collectionID, ok := collectionMap[folderID]; ok {
                i.collectionRepo.LinkAsset(ctx, collectionID, asset.ID)
            }
        }
        
        imported++
    }
    
    return &ImportResult{
        Imported: imported,
        Skipped:  skipped,
        Tags:     len(tagMap),
        Folders:  len(collectionMap),
    }, nil
}
```

### 4.5 Billfish 迁移实现

```go
type BillfishImporter struct {
    libraryPath string
    assetRepo   *repos.AssetRepo
    tagRepo     *repos.TagRepo
    collectionRepo *repos.CollectionRepo
}

type BillfishMetadata struct {
    Tags       []string `json:"tags"`
    Rating     int      `json:"rating"`
    Annotation string   `json:"annotation"`
    Source     string   `json:"source"`
    CreateTime int64    `json:"createTime"`
    ModifyTime int64    `json:"modifyTime"`
}

func (i *BillfishImporter) Import(ctx context.Context) (*ImportResult, error) {
    // 检测导入模式
    mode := i.detectImportMode()
    
    var imported, skipped int
    tagMap := make(map[string]string)
    collectionMap := make(map[string]string)
    
    switch mode {
    case "copy", "cut":
        // 复制/剪切模式：素材在库目录内
        err := filepath.Walk(i.libraryPath, func(path string, info os.FileInfo, err error) error {
            if err != nil || info.IsDir() {
                return nil
            }
            
            // 跳过 .bfmeta 和 .billfish 目录
            if strings.HasSuffix(path, ".bfmeta") || strings.Contains(path, ".billfish") {
                return nil
            }
            
            // 读取元数据
            meta := i.readMetadata(path + ".bfmeta")
            
            // 创建资产
            asset, err := i.createAsset(ctx, path, meta)
            if err != nil {
                skipped++
                return nil
            }
            
            // 关联标签和收藏集
            i.linkMetadata(ctx, asset, meta, tagMap, collectionMap, filepath.Dir(path))
            
            imported++
            return nil
        })
        return &ImportResult{Imported: imported, Skipped: skipped, Tags: len(tagMap)}, err
        
    case "index":
        // 索引模式：素材分散各处，需要从 .bfmeta 获取路径
        // 或从 .billfish 配置中读取索引信息
        err := filepath.Walk(i.libraryPath, func(path string, info os.FileInfo, err error) error {
            if err != nil || info.IsDir() || !strings.HasSuffix(path, ".bfmeta") {
                return nil
            }
            
            // .bfmeta 文件名去掉后缀就是素材路径
            assetPath := strings.TrimSuffix(path, ".bfmeta")
            if _, err := os.Stat(assetPath); os.IsNotExist(err) {
                skipped++
                return nil
            }
            
            meta := i.readMetadata(path)
            asset, err := i.createAsset(ctx, assetPath, meta)
            if err != nil {
                skipped++
                return nil
            }
            
            i.linkMetadata(ctx, asset, meta, tagMap, collectionMap, "")
            imported++
            return nil
        })
        return &ImportResult{Imported: imported, Skipped: skipped, Tags: len(tagMap)}, err
    }
    
    return nil, errors.New("unknown import mode")
}

func (i *BillfishImporter) detectImportMode() string {
    // 检查 .billfish 目录是否存在配置文件
    // 判断是否为索引模式
    cfgPath := filepath.Join(i.libraryPath, ".billfish", "library.json")
    if data, err := os.ReadFile(cfgPath); err == nil {
        var cfg struct {
            Mode string `json:"mode"`
        }
        if json.Unmarshal(data, &cfg) == nil {
            return cfg.Mode
        }
    }
    
    // 默认当作复制/剪切模式
    return "copy"
}
```

## 5. 用户界面设计

### 5.1 检测提示

```
┌─────────────────────────────────────────────────────────┐
│  🔍 检测到素材库                                        │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  检测到目录 "D:\素材库.library" 是 Eagle 素材库         │
│                                                         │
│  包含:                                                  │
│  • 3,256 个素材                                        │
│  • 89 个标签                                           │
│  • 24 个文件夹                                         │
│                                                         │
│  是否导入现有分类和标签？                               │
│                                                         │
│  [跳过，仅索引文件]  [导入全部元数据]                   │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

### 5.2 导入进度

```
┌─────────────────────────────────────────────────────────┐
│  正在导入 Eagle 素材库...                               │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  ████████████████████░░░░░░░░░░  65%                   │
│                                                         │
│  已导入: 2,116 / 3,256 个素材                          │
│  已创建: 89 个标签, 24 个收藏集                        │
│                                                         │
│  当前: 海报设计.jpg                                    │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

### 5.3 导入结果

```
┌─────────────────────────────────────────────────────────┐
│  ✅ 导入完成                                            │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  素材: 3,256 个 (跳过 12 个)                           │
│  标签: 89 个                                           │
│  收藏集: 24 个                                         │
│  评分: 156 个素材带有评分                              │
│  备注: 89 个素材带有备注                               │
│                                                         │
│  [查看素材库]                                          │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

## 6. API 设计

### 6.1 检测接口

```
POST /api/library/detect
Request:  { "path": "D:\\素材库.library" }
Response: { 
  "type": "eagle",
  "stats": {
    "assets": 3256,
    "tags": 89,
    "folders": 24
  }
}
```

### 6.2 导入接口

```
POST /api/library/import
Request:  { 
  "path": "D:\\素材库.library",
  "type": "eagle",
  "options": {
    "importTags": true,
    "importFolders": true,
    "importRatings": true,
    "importAnnotations": true
  }
}
Response: { 
  "taskId": "import_123",
  "status": "running"
}
```

### 6.3 导入状态

```
GET /api/library/import/:taskId
Response: {
  "status": "running",
  "progress": 65,
  "imported": 2116,
  "total": 3256,
  "tags": 89,
  "folders": 24
}
```

## 7. 注意事项

### 7.1 Eagle 特殊处理

| 问题 | 解决方案 |
|------|----------|
| 文件名随机 | 使用 metadata.json 中的 name 字段 |
| 文件重复引用 | 一个素材可在多个文件夹，用关联表处理 |
| 图片格式 | 支持所有 Eagle 支持的格式 |

### 7.2 Billfish 特殊处理

| 问题 | 解决方案 |
|------|----------|
| 多种导入模式 | 检测 .billfish 目录判断是否为库目录 |
| 复制/剪切模式 | 素材在库目录内，直接索引 |
| 索引模式 | 素材分散各处，需读取 .bfmeta 获取路径 |
| 元数据文件 | .bfmeta 文件可能不存在，需要容错 |
| 文件夹结构 | 保持原有目录结构，映射为收藏集 |

### 7.3 通用处理

| 问题 | 解决方案 |
|------|----------|
| 大库导入 | 异步任务，支持断点续传 |
| 重复素材 | 使用指纹去重 |
| 编码问题 | 统一使用 UTF-8 |

## 8. 实现优先级

| 优先级 | 功能 | 理由 |
|--------|------|------|
| P0 | Eagle 库检测 | 用户量大，付费用户迁移意愿强 |
| P0 | Eagle 元数据导入 | 核心迁移功能 |
| P1 | Billfish 库检测 | 免费用户量大 |
| P1 | Billfish 元数据导入 | 核心迁移功能 |
| P2 | 导入进度显示 | 用户体验 |
| P2 | 断点续传 | 大库稳定性 |
| P3 | 其他工具支持 | 扩展用户群 |

## 9. 后续扩展

1. **Pixave 支持** - Mac 用户群体
2. **Inboard 支持** - Mac 用户群体
3. **自定义导入规则** - 用户可配置映射关系
4. **增量同步** - 定期同步原库更新
5. **双向同步** - 支持导出到其他格式

---

## 变更记录

| 日期 | 变更内容 |
|------|----------|
| 2026-02-17 | 初版设计文档创建 |
