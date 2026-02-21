# 智归档OS 开发指南

> 快速上手指南，从环境搭建到本地调试

---

## 1. 环境准备

### 1.1 必备工具

| 工具 | 版本要求 | 用途 |
|------|---------|------|
| Go | 1.22+ | 后端开发 |
| Node.js | 18+ | 前端开发 |
| pnpm | 8+ | 前端包管理 |
| Git | 任意 | 版本控制 |

### 1.2 安装步骤

**安装 Go**:
前往 [golang.org](https://golang.org/dl/) 下载并安装。

**安装 Node.js**:
建议使用 nvm 安装 Node.js 18 或更高版本。

**安装 pnpm**:
```bash
npm install -g pnpm
```

### 1.3 验证安装

```bash
go version
node --version
pnpm --version
```

---

## 2. 项目结构

```
media-assistant-rust/
├── cmd/
│   └── core/                    # 常驻服务入口（托盘）
├── internal/                    # 核心业务逻辑
│   ├── db/                      # 数据库逻辑
│   ├── repos/                   # 数据持久化层
│   ├── services/                # 核心服务层
│   └── httpapi/                 # 内部 HTTP API 服务
├── models/                      # 数据模型
├── frontend/                    # 前端与 Electron
│   ├── src/                     # Vue 3 代码
│   ├── electron/                # Electron 主进程
│   └── package.json             # 前端依赖
├── build/                       # 构建资源
└── docs/                        # 文档库
```

---

## 3. 编码规范 (必须遵守)

### 3.1 注释要求
*   **强制中文注释**：所有核心逻辑、复杂算法、函数入口及关键变量定义，**必须**使用清晰的中文注释。
*   **文档注释**：公共函数必须使用 `///` 进行文档化说明。

### 3.2 代码风格
*   **Go**：遵循 `gofmt` 格式化。
*   **Vue/TS**：遵循 ESLint 规则，使用 Composition API 风格。

### 3.3 架构原则
*   **非侵入式**：严禁修改用户原始素材文件。
*   **内容实例分离**：资产管理必须基于指纹 (Hash) 与路径 (Path) 的映射。
*   **异步处理**：所有磁盘 I/O 和计算密集型任务必须在异步任务中执行。

---

## 4. AI 编码指南 (Prompt 专用)
> 当你使用 AI 辅助编码时，请务必在 Prompt 中包含以下内容：

**“请严格遵守项目根目录下 `docs/development-guide.md` 中的开发规范。特别注意：所有核心模块的关键部分必须使用【中文注释】；遵循‘内容与实例分离’的架构原则；保持对用户素材的非侵入性。”**

---

## 5. 开发流程

### 5.1 首次启动

**克隆仓库**:
```bash
git clone <repository-url>
cd media-assistant-rust
```

**安装依赖**:
```bash
go mod tidy
cd frontend
pnpm install
```

**启动后端常驻服务**:
```bash
go run ./cmd/core
```

**启动前端开发服务器**:
```bash
cd frontend
pnpm dev
```

**启动 Electron UI**:
```bash
cd frontend
pnpm electron:dev
```

如需指定后端地址，可在启动前设置：
```bash
set MA_CORE_BASE_URL=http://127.0.0.1:32000
```

### 5.2 日常开发

**只跑前端**:
```bash
cd frontend
pnpm dev
```

**只跑后端**:
```bash
go run ./cmd/core
```

**前后端联调**:
```bash
go run ./cmd/core
cd frontend
pnpm dev
pnpm electron:dev
```

---

## 6. 调试技巧

### 6.1 后端调试
```bash
go test ./...
```

### 6.2 前端调试
```bash
pnpm lint
pnpm exec vue-tsc --noEmit
```

---

## 7. 学习资源

### Go
- [Go 文档](https://go.dev/doc/)

### Electron
- [Electron Documentation](https://www.electronjs.org/docs/latest)

### Vue
- [Vue 3 Documentation](https://vuejs.org/guide/)
