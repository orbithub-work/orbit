# 📦 应用打包与启动流程模拟 (Application Packaging & Launch Simulation)

## 1. 概述

Media Assistant 采用 **双进程架构 (Split Architecture)**，在发布时需要将 Go 后端和 Electron 前端捆绑在一起。

本文档描述了未来的发布形态、目录结构以及用户启动的完整流程。

## 2. 目录结构 (Release Structure)

打包后的根目录（即用户解压后的文件夹）将包含以下核心文件：

```text
MediaAssistant/
├── media-assistant-core.exe  # [入口] Go 后端服务，负责系统托盘、API、数据库
├── media-assistant-ui.exe    # [UI] Electron 前端容器，由 Core 自动唤起
├── Start-MediaAssistant.bat  # [快捷方式] 模拟用户双击启动的脚本
├── resources/                # Electron 资源包 (app.asar)
├── locales/                  # 语言包
├── swiftshader/              # 图形渲染库
└── ... (其他 DLL 依赖)
```

## 3. 打包脚本 (Packaging Script)

为了实现上述结构，我们提供了 `scripts/simulate_packaging.ps1` 脚本。

### 脚本逻辑
1.  **构建 Go Core**: 编译 `cmd/core` 生成 `media-assistant-core.exe`。
2.  **构建 Electron UI**:
    *   `pnpm build`: 生成 Vue 静态资源。
    *   `electron-builder`: 将 Electron 打包为可执行文件 (Unpacked 模式)。
3.  **合并**: 将 Electron 的输出文件全部复制到发布目录，并将 UI 可执行文件重命名为 `media-assistant-ui.exe`。
4.  **配置**: 确保 Go Core 和 UI 在同一级目录，满足 Core 的自动查找逻辑。

## 4. 用户启动流程 (Launch Sequence)

### 步骤 1: 用户操作
用户双击 `media-assistant-core.exe` (或其快捷方式)。

### 步骤 2: Go Core 初始化
1.  **启动**: Core 进程启动，加载 `media_assistant.db` 数据库。
2.  **系统托盘**: 在任务栏右下角显示应用图标。
3.  **API 服务**: 在随机端口 (如 `32000`) 启动 HTTP/WebSocket 服务。
4.  **环境检测**:
    *   检测是否为首次运行 (First Launch)。
    *   检测是否已有实例运行 (Single Instance Lock)。

### 步骤 3: 唤起 UI (Spawn UI)
如果 Core 判断需要显示界面（首次启动或用户点击托盘），它会执行以下逻辑：
1.  **查找 UI**: 在当前目录下查找 `media-assistant-ui.exe`。
2.  **注入环境变量**:
    *   `MA_CORE_BASE_URL=http://localhost:32000` (告知 UI 后端地址)
    *   `MA_CORE_MODE=external` (告知 UI 它是被外部唤起的)
3.  **执行**: `exec.Command("media-assistant-ui.exe")`。

### 步骤 4: Electron UI 启动
1.  **Main Process**: Electron 主进程启动，读取 `MA_CORE_BASE_URL`。
2.  **创建窗口**:
    *   创建 `BrowserWindow`。
    *   加载 Vue 页面 (`index.html`)。
3.  **连接**: 前端通过 API Client 连接到 `http://localhost:32000`。
4.  **渲染**: 用户看到管理界面。

## 5. 如何测试模拟

1.  **前置条件**: 确保已安装 Node.js, pnpm, Go。
2.  **运行脚本**:
    ```powershell
    ./scripts/simulate_packaging.ps1
    ```
3.  **验证结果**:
    *   查看 `release/` 目录。
    *   双击运行 `release/media-assistant-core.exe`。
    *   观察系统托盘是否出现图标。
    *   点击托盘 -> "显示界面"，验证 Electron 窗口是否成功弹出。
