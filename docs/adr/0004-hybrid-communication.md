# 0004. 混合通信架构 (Hybrid Communication Architecture)

## 背景 (Context)
Media Assistant 需要同时支持两类客户端：
1.  **主界面 (Main UI)**: 运行在 Electron 渲染进程中，通过预加载桥接访问主进程能力。
2.  **插件/外部工具 (Plugins)**: 运行在外部进程（如 PR 插件、独立脚本），通过网络协议通信，权限受限。

## 决策 (Decision)
采用 **混合通信模式 (Hybrid Communication Mode)**，根据功能特性将 API 分为两类。

### 1. 原生 API (Native API)
*   **通信通道**: Electron 主进程 + 预加载桥接 (IPC)。
*   **适用场景**:
    *   窗口管理 (Window Management): 显示/隐藏 Dock、置顶、调整大小。
    *   系统交互 (System Integration): 打开文件选择框、系统通知。
    *   应用生命周期 (Lifecycle): 退出、重启、更新检查。
*   **实现位置**: Electron 主进程与预加载桥接。

### 2. 数据 API (Data API)
*   **通信通道**: 主界面与插件统一使用 HTTP/WebSocket。
*   **适用场景**:
    *   资产管理 (Assets): 搜索、索引、元数据提取。
    *   项目管理 (Projects): 创建、列表、状态更新。
    *   日志与统计 (Logs/Stats)。
*   **实现架构**:
    *   核心逻辑封装在 `internal/services` 层 (如 `AssetService`)。
    *   `httpapi` Handler 负责对外 HTTP/WS 调用，统一入口。

## 接口分类清单 (API Classification)

| 类别 | 功能 | 方法/路径 | 通道 | 备注 |
| :--- | :--- | :--- | :--- | :--- |
| **Native** | Dock 控制 | `ShowDock`, `HideDock` | Electron Only | 必须在主进程执行 |
| **Native** | 文件选择 | `SelectProjectFolder` | Electron Only | 阻塞 UI |
| **Data** | 搜索 | `SearchFiles` / `/api/search` | Hybrid | 插件需频繁调用 |
| **Data** | 项目列表 | `ListProjects` / `/api/projects` | Hybrid | 插件需知道项目上下文 |

## 前端策略 (Frontend Strategy)
前端 `api.ts` 统一通过 HTTP API 调用后端服务：
*   Electron 渲染进程通过 `window.mediaAssistant` 获得后端基址。
*   浏览器模式或插件使用相同的 HTTP 入口。
