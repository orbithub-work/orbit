# 0005. 动态注册式插件架构 (Active Registration Plugin Architecture)

## 背景 (Context)
为了支持更灵活的插件生态（包括内网分布式插件、无需重启即插即用），我们放弃传统的“扫描本地 Manifest 文件”模式，转向 **“插件主动注册”** 模式。
这类似于微服务架构中的服务发现，或 AI Agent 体系中的 MCP (Model Context Protocol)。

## 决策 (Decision)

### 1. 核心流程
1.  **插件启动**: 插件进程（本地 EXE 或远程服务）启动。
2.  **服务发现**: 插件读取 `server.port` 或通过 mDNS 找到 Host 地址。
3.  **握手注册**: 插件向 Host 发送 `POST /api/plugins/register`，提交自身元数据、能力描述和回调地址。
4.  **心跳保活**: 插件定期发送心跳，Host 维护在线列表。
5.  **反向调用**: Host 需要执行某项能力时，根据注册的回调地址 (HTTP/WebSocket) 调用插件。

### 2. API 协议

#### 注册接口
`POST /api/plugins/register`

**Request Body**:
```json
{
  "id": "com.example.jianying",
  "name": "剪映助手",
  "version": "1.0.0",
  "description": "剪映项目解析与同步",
  "endpoint": "http://127.0.0.1:45678", // 插件监听的地址
  "capabilities": [
    {
      "name": "import_draft",
      "description": "导入剪映草稿",
      "params_schema": { ... }
    }
  ],
  "config_schema": {
    "properties": {
      "auto_sync": { "type": "boolean", "default": true }
    }
  }
}
```

**Response**:
```json
{
  "status": "registered",
  "token": "plugin_token_123", // 用于后续鉴权
  "host_version": "0.1.0"
}
```

### 3. 数据持久化
*   **内存注册表 (Registry)**: 维护当前**在线**的插件实例及其 Endpoint。
*   **数据库 (DB)**: 仅存储用户对插件的**配置值** (Config Values) 和历史授权状态。插件重启注册后，Host 自动下发之前的配置。

## 优势
*   **灵活性**: 支持任何语言编写的插件，只要能发 HTTP 请求。
*   **分布式**: 插件可以运行在局域网的另一台机器上。
*   **解耦**: Host 不关心插件的安装位置和启动方式。
