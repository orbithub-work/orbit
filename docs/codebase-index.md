# 智归档 OS 代码索引

最后更新：2026-02-17

## 1. 项目概览

- 仓库名：`smart-archive-os`
- 主要技术栈：后端 Go 1.24（SQLite + Bun ORM）；前端 Vue 3 + TypeScript + Pinia + Vite + Electron
- 代码规模（排除 `frontend/node_modules`）：约 217 个文件

## 2. 顶层目录索引

- `cmd/`：Go 程序入口
- `internal/`：后端业务核心（服务、仓储、HTTP API、处理器）
- `frontend/`：Vue + Electron 前端工程
- `docs/`：产品、架构、开发文档与 ADR
- `scripts/`：辅助脚本
- `data/`：运行期数据目录（开发环境）

## 3. 运行入口

### 后端入口

- `cmd/core/main.go`
- 功能：初始化日志、单实例检查、启动 `core.System`、挂载 HTTP API、处理 IPC 与系统能力（打开文件/目录、项目管理、索引、任务、插件、活动日志等）

### 后端系统装配

- `internal/core/bootstrap.go`
- 功能：解析数据目录、打开与迁移数据库、初始化 Repo/Service、注册媒体解析器、启动队列与扫描器、启动文件监听器

### 前端入口

- `frontend/src/main.ts`
- 功能：创建 Vue App，注册 Pinia、Router、虚拟滚动组件

### 路由入口

- `frontend/src/router/index.ts`
- 路由模式：`createWebHashHistory()`
- 主要页面：`Dashboard`、`Files`、`Collections`、`Projects`、`Settings`、`Onboarding`、`Dock`、`TrayMenu`

## 4. 后端模块索引（`internal/`）

- `internal/httpapi/`：HTTP 接口层
- 关键文件：`server.go`、`deps.go`、`middleware.go`、`handler_*.go`
- `internal/services/`：业务服务层
- 覆盖：资产、扫描、项目、任务、标签、插件、活动、设置、文件监听、去重、缓存等
- `internal/repos/`：数据访问层
- 覆盖：`project`、`asset`、`project_asset`、`asset_lineage`、`tag`、`activity`、`task`、`artifact`、`event_log`
- `internal/processor/`：媒体解析器框架与实现
- 子模块：`image`、`video`、`psd`、`raw`、`fallback`
- `internal/db/`：数据库初始化与迁移
- `internal/models/`：领域模型定义
- `internal/infra/`：基础设施能力（数据目录、单实例、端口文件等）
- `internal/utils/`：跨模块工具函数
- `internal/pkg/logger/`：日志模块

## 5. 前端模块索引（`frontend/src/`）

- `views/`：页面级视图（路由组件）
- `layouts/`：布局壳层（如 `MainLayout.vue`）
- `components/`：复用组件
- `asset-grid/`：资产网格（含 `VirtualAssetGrid.vue`）
- `app-shell/`：主应用壳组件（含 `PoolView.vue`）
- `file-manager/`：文件管理组件集
- `asset-card/`：资产卡片渲染组件
- `collection/`、`batch/`、`preview/`、`project/`：功能分区组件
- `stores/`：Pinia 状态管理（project、collection、file、tag、theme、plugin 等）
- `composables/`：组合式逻辑（如 `useApi.ts`）
- `services/`：API 请求封装（`api.ts`）
- `styles/`：全局样式与主题变量

## 6. 构建与开发命令

### 后端（根目录）

- `make`（查看/执行工程定义目标，见 `Makefile`）
- 可执行入口：`cmd/core/main.go`

### 前端（`frontend/`）

- `npm run dev`：Vite 开发服务（端口 `5176`）
- `npm run build`：构建前端
- `npm run electron:dev`：启动 Electron
- `npm run lint`：ESLint（带 `--fix`）

## 7. 当前高频改动区域（按近期工作语义）

- `frontend/src/components/asset-grid/VirtualAssetGrid.vue`
- `frontend/src/components/app-shell/PoolView.vue`
- `frontend/src/AppSafe.vue`
- `frontend/src/layouts/MainLayout.vue`
- `frontend/src/composables/useApi.ts`

## 8. 建议阅读顺序（新加入开发者）

1. `docs/README.md`
2. `docs/development-guide.md`
3. `cmd/core/main.go`
4. `internal/core/bootstrap.go`
5. `internal/httpapi/server.go` 与 `internal/httpapi/deps.go`
6. `frontend/src/main.ts` 与 `frontend/src/router/index.ts`
7. `frontend/src/stores/` + `frontend/src/services/api.ts`
