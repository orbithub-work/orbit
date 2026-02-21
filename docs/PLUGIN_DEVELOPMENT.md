# 智归档OS 插件系统

## 概述

智归档OS采用**动态注册式插件架构**，支持3种插件模式：

1. **frontend** - 纯前端插件（Vue组件，动态加载）
2. **local_process** - 本地进程插件（独立EXE，HTTP通信）
3. **network_service** - 网络服务插件（远程服务，HTTP通信）

## 扩展点（Extension Points）

### 页面级（Page）
- `Global.Page` - 全局页面（独占MainView）

### 面板级（Panel）
- `Pool.Sidebar.Section` - 素材库左侧边栏面板
- `Pool.Inspector.Tab` - 素材库右侧检查器标签页
- `Workspace.Sidebar.Section` - 工作台左侧边栏面板
- `Workspace.Inspector.Tab` - 工作台右侧检查器标签页
- `Artifact.Sidebar.Section` - 成品库左侧边栏面板
- `Artifact.Inspector.Tab` - 成品库右侧检查器标签页
- `Rights.Sidebar.Section` - 版权中心左侧边栏面板
- `Rights.Inspector.Tab` - 版权中心右侧检查器标签页

### 按钮级（Action）
- `Pool.Toolbar.Action` - 素材库工具栏按钮
- `Pool.ContextMenu.Item` - 素材右键菜单项
- `Workspace.Toolbar.Action` - 工作台工具栏按钮
- `Workspace.ContextMenu.Item` - 工作台右键菜单项
- `Artifact.Toolbar.Action` - 成品库工具栏按钮
- `Artifact.ContextMenu.Item` - 成品右键菜单项
- `Rights.Toolbar.Action` - 版权中心工具栏按钮

### 状态栏级（StatusBar）
- `Global.StatusBar.Left` - 状态栏左侧小部件
- `Global.StatusBar.Right` - 状态栏右侧小部件

## 插件开发

### 1. 创建插件manifest

```json
{
  "id": "com.example.myplugin",
  "name": "我的插件",
  "version": "1.0.0",
  "description": "插件描述",
  "mode": "frontend",
  "tier": "free",
  "permissions": [
    "assets:read",
    "assets:write",
    "ui:notification"
  ],
  "mounts": [
    {
      "slot": "Pool.Sidebar.Section",
      "entry": "MyComponent",
      "title": "我的面板",
      "icon": "🎨",
      "order": 100
    }
  ]
}
```

### 2. 创建Vue组件

```vue
<template>
  <div class="my-plugin-panel">
    <h3>{{ title }}</h3>
    <button @click="handleAction">执行操作</button>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { createPluginAPI } from '@/core/pluginAPI'

const api = createPluginAPI()
const title = ref('我的插件')

async function handleAction() {
  const selected = api.context.getSelectedAssets()
  await api.ui.showNotification(`选中了 ${selected.length} 个素材`, 'info')
}
</script>
```

### 3. 导出插件

```typescript
// src/index.ts
import MyComponent from './components/MyComponent.vue'

export const components = {
  MyComponent
}

export const actions = {
  async myAction(context: any) {
    console.log('Action triggered', context)
  }
}

export default {
  components,
  actions
}
```

### 4. 注册插件

```bash
# 开发模式：插件自动扫描
# 生产模式：插件需要注册

POST http://localhost:32000/api/plugins/register
Content-Type: application/json

{
  "plugin_id": "com.example.myplugin",
  "name": "我的插件",
  "version": "1.0.0",
  "mode": "frontend",
  "endpoint": "http://localhost:8001",
  "mounts": [...]
}
```

## Plugin API

### Assets API
```typescript
// 列出素材
const assets = await api.assets.list({ search: 'photo', limit: 10 })

// 获取素材详情
const asset = await api.assets.get(123)

// 更新素材
await api.assets.update(123, { rating: 5 })

// 删除素材
await api.assets.delete(123)

// 获取选中的素材
const selected = api.assets.getSelected()
```

### Tags API
```typescript
// 列出标签
const tags = await api.tags.list()

// 创建标签
const tag = await api.tags.create({ name: '重要', color: '#ff0000' })

// 批量添加标签
await api.tags.batchAdd([1, 2, 3], [10, 11])
```

### UI API
```typescript
// 显示通知
await api.ui.showNotification('操作成功', 'success')

// 显示对话框
const result = await api.ui.showDialog({
  title: '确认',
  message: '确定要删除吗？',
  type: 'confirm'
})

// 显示确认框
const confirmed = await api.ui.showConfirm('确定要继续吗？')
```

### Context API
```typescript
// 获取选中的素材
const assets = api.context.getSelectedAssets()

// 获取当前视图
const view = api.context.getCurrentView() // 'pool' | 'workspace' | 'artifact' | 'rights'

// 设置上下文
await api.context.setContext({
  asset_ids: [1, 2, 3],
  view: 'pool',
  action: 'select'
})
```

## 示例插件

### 版权中心插件（Pro）

位置：`pro/plugins/copyright/`

功能：
- 版权状态面板（侧边栏）
- 快速确权按钮（工具栏）
- 血缘链图表（检查器）
- 查看血缘链（右键菜单）
- 版权指示器（状态栏）

## 权限系统

### 开源版权限
- `assets:read` - 读取素材
- `assets:write` - 修改素材
- `tags:read` - 读取标签
- `tags:write` - 修改标签
- `ui:notification` - 显示通知
- `ui:dialog` - 显示对话框

### Pro版权限
- `lineage:read` - 读取血缘链
- `lineage:write` - 修改血缘链
- `artifacts:read` - 读取成品
- `artifacts:write` - 修改成品
- `workflow:read` - 读取工作流
- `workflow:write` - 修改工作流
- `publish:read` - 读取发布记录
- `publish:write` - 创建发布任务
- `metrics:read` - 读取数据指标
- `metrics:write` - 上报数据

## 开发工具

### 调试插件
```bash
# 启动主程序
npm run dev

# 启动插件开发服务器（如果是network_service模式）
cd pro/plugins/copyright
npm run dev
```

### 查看已注册插件
```bash
curl http://localhost:32000/api/plugins/list
```

### 查看扩展点
```bash
curl http://localhost:32000/api/extensions/slots
```

## 最佳实践

1. **插件命名**：使用反向域名格式（com.company.plugin）
2. **版本管理**：遵循语义化版本（Semantic Versioning）
3. **权限最小化**：只申请必需的权限
4. **错误处理**：妥善处理API调用失败
5. **性能优化**：避免频繁调用API，使用缓存
6. **用户体验**：提供清晰的加载状态和错误提示

## 发布插件

### 开源插件
1. 提交到GitHub
2. 发布到npm
3. 提交到插件市场

### Pro插件
1. 编译混淆
2. 加密打包
3. 上传到Pro插件服务器
