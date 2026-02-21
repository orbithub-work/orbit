# 智归档OS 插件清单规范（Manifest Specification）

## 核心原则

**一切组件必须通过 manifest.json 向 Go 内核注册，Go 内核是唯一的调度中心。**

---

## Manifest 结构

### 通用字段（所有类型必需）

```json
{
  "id": "com.smartarchive.component-name",
  "name": "组件显示名称",
  "version": "1.0.0",
  "type": "frontend | backend | satellite",
  "description": "组件描述",
  "author": "作者名称",
  "license": "MIT | proprietary",
  "tier": "free | pro | enterprise"
}
```

---

## 三种组件类型

### 1. Frontend（前端组件 - UI皮肤）

**类比**：浏览器插件、VS Code主题

**特点**：
- 纯静态资源（JS/CSS/Vue组件）
- 由Go内核通过HTTP提供给Electron加载
- 无独立进程，无网络通信

**Manifest示例**：
```json
{
  "id": "com.smartarchive.copyright-ui",
  "name": "版权状态面板",
  "version": "1.0.0",
  "type": "frontend",
  "tier": "pro",
  "entry": "./dist/index.js",
  "mounts": [
    {
      "slot": "Pool.Sidebar.Section",
      "component": "CopyrightStatus",
      "title": "版权状态",
      "icon": "🔒",
      "order": 100
    }
  ],
  "permissions": [
    "assets:read",
    "ui:notification"
  ]
}
```

**加载方式**：
```go
// Go内核启动时扫描 plugins/frontend/ 目录
// 通过 /plugins/{id}/index.js 提供静态资源
```

---

### 2. Backend（后端扩展 - 进程调用）

**类比**：显卡驱动、FFmpeg、ImageMagick

**特点**：
- 独立可执行文件（.exe/.bin）
- 由Go内核通过 os/exec 启动和管理
- 通过 stdin/stdout 或 HTTP 与Go通信
- Go负责生命周期管理

**Manifest示例**：
```json
{
  "id": "com.smartarchive.ffmpeg-proxy",
  "name": "FFmpeg视频处理",
  "version": "1.0.0",
  "type": "backend",
  "tier": "free",
  "executable": "./bin/ffmpeg-proxy",
  "capabilities": [
    {
      "name": "generate_thumbnail",
      "input": ["video_path"],
      "output": ["thumbnail_path"]
    },
    {
      "name": "extract_metadata",
      "input": ["video_path"],
      "output": ["duration", "resolution", "codec"]
    }
  ],
  "permissions": [
    "fs:read",
    "fs:write"
  ]
}
```

**启动方式**：
```go
// Go内核按需启动
cmd := exec.Command(plugin.Executable)
cmd.Stdin = jsonInput
cmd.Stdout = jsonOutput
cmd.Start()
```

---

### 3. Satellite（卫星应用 - 独立内网通信）

**类比**：蓝牙设备、打印机、Dock应用

**特点**：
- 完全独立的应用程序（可以单独启动）
- 通过HTTP/WebSocket与Go内核通信
- 启动后主动向Go注册
- Go不负责启动，只负责路由和鉴权

**Manifest示例**：
```json
{
  "id": "com.smartarchive.dock",
  "name": "剪辑助手 Dock",
  "version": "1.0.0",
  "type": "satellite",
  "tier": "pro",
  "entry": "http://127.0.0.1:9090",
  "capabilities": [
    {
      "name": "generate_pr_draft",
      "description": "生成Premiere Pro草稿",
      "endpoint": "/api/generate-draft",
      "method": "POST"
    },
    {
      "name": "sync_timeline",
      "description": "同步时间线",
      "endpoint": "/api/sync-timeline",
      "method": "POST"
    }
  ],
  "permissions": [
    "assets:read",
    "lineage:read",
    "artifacts:write"
  ],
  "heartbeat": {
    "interval": 5,
    "endpoint": "/health"
  }
}
```

**注册流程**：
```
1. Dock应用启动
2. 读取自己的manifest.json
3. POST http://127.0.0.1:8848/api/plugins/register
4. Go内核验证License和权限
5. Go内核返回Token
6. Dock每5秒发送心跳
```

---

## 权限系统

### 开源版权限（Free Tier）
```
assets:read          - 读取素材
assets:write         - 修改素材元数据
tags:read            - 读取标签
tags:write           - 修改标签
fs:read              - 读取文件系统
ui:notification      - 显示通知
ui:dialog            - 显示对话框
```

### Pro版权限（Pro Tier）
```
lineage:read         - 读取血缘链
lineage:write        - 创建血缘关系
artifacts:read       - 读取成品
artifacts:write      - 创建成品
workflow:read        - 读取工作流
workflow:write       - 修改工作流
publish:read         - 读取发布记录
publish:write        - 创建发布任务
```

### 企业版权限（Enterprise Tier）
```
admin:users          - 管理用户
admin:license        - 管理授权
admin:audit          - 审计日志
```

---

## Go内核的"收费站"

### 路由鉴权中间件

```go
func (h *Handler) withPluginAuth(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // 1. 从请求头获取插件Token
        token := r.Header.Get("X-Plugin-Token")
        
        // 2. 验证Token，获取插件信息
        plugin := h.pluginService.GetPluginByToken(token)
        if plugin == nil {
            writeJSON(w, 403, APIResponse{
                Success: false,
                Error: "Invalid plugin token"
            })
            return
        }
        
        // 3. 检查插件Tier是否有权限
        if plugin.Tier == "pro" && !h.licenseService.IsProActive() {
            writeJSON(w, 403, APIResponse{
                Success: false,
                Error: "This feature requires Pro license",
                Data: map[string]interface{}{
                    "upgrade_url": "https://smartarchive.cn/pricing"
                }
            })
            return
        }
        
        // 4. 检查具体权限
        requiredPermission := getRequiredPermission(r.URL.Path)
        if !hasPermission(plugin.Permissions, requiredPermission) {
            writeJSON(w, 403, APIResponse{
                Success: false,
                Error: "Permission denied: " + requiredPermission
            })
            return
        }
        
        // 5. 通过鉴权，执行请求
        next(w, r)
    }
}
```

---

## 目录结构

```
smart-archive-os/
├── plugins/
│   ├── frontend/              # 前端组件
│   │   ├── copyright-ui/
│   │   │   ├── manifest.json
│   │   │   └── dist/
│   │   │       └── index.js
│   │   └── quick-tools/
│   │       ├── manifest.json
│   │       └── dist/
│   │           └── index.js
│   │
│   ├── backend/               # 后端扩展
│   │   ├── ffmpeg-proxy/
│   │   │   ├── manifest.json
│   │   │   └── bin/
│   │   │       └── ffmpeg-proxy
│   │   └── ai-analyzer/
│   │       ├── manifest.json
│   │       └── bin/
│   │           └── ai-analyzer
│   │
│   └── satellites/            # 卫星应用（可选，可以独立部署）
│       └── dock/
│           ├── manifest.json
│           └── SmartArchiveDock.app
│
├── cmd/core/                  # Go内核
└── frontend/                  # Electron主界面
```

---

## 通信协议

### 统一使用 JSON-RPC over HTTP

**请求格式**：
```json
{
  "jsonrpc": "2.0",
  "method": "generate_thumbnail",
  "params": {
    "asset_id": "123",
    "width": 300,
    "height": 200
  },
  "id": 1
}
```

**响应格式**：
```json
{
  "jsonrpc": "2.0",
  "result": {
    "thumbnail_path": "/path/to/thumb.jpg"
  },
  "id": 1
}
```

**错误格式**：
```json
{
  "jsonrpc": "2.0",
  "error": {
    "code": -32600,
    "message": "Invalid request"
  },
  "id": 1
}
```

---

## 生命周期

### Frontend组件
```
1. Go启动时扫描 plugins/frontend/
2. 解析manifest.json
3. 注册到PluginRegistry
4. 通过HTTP提供静态资源
5. Electron加载时动态import
```

### Backend扩展
```
1. Go按需启动进程
2. 通过stdin发送任务
3. 通过stdout接收结果
4. 任务完成后保持进程或关闭
5. 异常时自动重启
```

### Satellite应用
```
1. 用户手动启动（或开机自启）
2. 读取manifest.json
3. POST /api/plugins/register
4. 获得Token
5. 每5秒发送心跳
6. 接收Go转发的请求
7. 用户关闭时注销
```

---

## 下一步实现

1. ✅ 定义manifest.json规范（本文档）
2. ⏭️ 实现Go的PluginScanner（扫描并解析manifest）
3. ⏭️ 实现PluginRegistry（注册表）
4. ⏭️ 实现PluginRouter（路由转发）
5. ⏭️ 实现LicenseValidator（授权验证）
6. ⏭️ 创建示例插件（frontend、backend、satellite各一个）
