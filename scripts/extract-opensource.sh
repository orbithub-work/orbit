#!/bin/bash

# OrbitHub 开源版本抽离脚本
# 从闭源仓库抽离核心功能，创建干净的开源仓库

set -e  # 遇到错误立即退出

echo "🚀 OrbitHub 开源版本抽离脚本"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# 配置
SOURCE_DIR="/Users/a/Projects/smart-archive-os"
TARGET_DIR="/Users/a/Projects/orbit"
GITHUB_REPO="git@github.com:orbithub-work/orbit.git"

# 1. 创建目标目录
echo ""
echo "1️⃣  创建目标目录..."
if [ -d "$TARGET_DIR" ]; then
    echo "⚠️  目标目录已存在，是否删除并重新创建？(y/n)"
    read -r response
    if [ "$response" = "y" ]; then
        rm -rf "$TARGET_DIR"
    else
        echo "❌ 取消操作"
        exit 1
    fi
fi

mkdir -p "$TARGET_DIR"
cd "$TARGET_DIR"

# 2. 初始化Git仓库
echo ""
echo "2️⃣  初始化Git仓库..."
git init
git config user.name "zengchangwt"
git config user.email "zengchang42@gmail.com"
git config user.signingkey ~/.ssh/github_ed25519.pub
git config commit.gpgsign true
git remote add origin "$GITHUB_REPO"

# 3. 复制核心文件
echo ""
echo "3️⃣  复制核心文件..."

# 复制目录
cp -r "$SOURCE_DIR/cmd" ./
cp -r "$SOURCE_DIR/internal" ./
cp -r "$SOURCE_DIR/frontend" ./
cp -r "$SOURCE_DIR/docs" ./
cp -r "$SOURCE_DIR/scripts" ./

# 复制根文件
cp "$SOURCE_DIR/go.mod" ./
cp "$SOURCE_DIR/go.sum" ./
cp "$SOURCE_DIR/Makefile" ./
cp "$SOURCE_DIR/.gitignore" ./

# 4. 删除Pro功能代码
echo ""
echo "4️⃣  删除Pro功能代码..."

# 删除Pro服务
rm -f internal/services/workflow_service.go
rm -f internal/services/publish_metrics_service.go
rm -f internal/services/lineage_candidate_service.go

# 删除Pro API处理器
rm -f internal/httpapi/handler_workflow.go
rm -f internal/httpapi/handler_publish_metrics.go

# 删除Pro数据模型
rm -f internal/models/workflow_template.go
rm -f internal/models/project_workflow.go
rm -f internal/models/project_workflow_step.go
rm -f internal/models/project_roadmap_item.go
rm -f internal/models/project_note.go
rm -f internal/models/publish_*.go
rm -f internal/models/metrics_*.go
rm -f internal/models/lineage_candidate.go

# 删除Pro仓库
rm -f internal/repos/workflow_*.go
rm -f internal/repos/publish_*.go
rm -f internal/repos/metrics_*.go
rm -f internal/repos/lineage_candidate_repo.go
rm -f internal/repos/project_workflow*.go
rm -f internal/repos/project_roadmap_repo.go
rm -f internal/repos/project_note_repo.go

# 删除Pro目录
rm -rf pro/

# 删除旧的插件示例（会重新创建）
rm -rf plugins/frontend/copyright-ui/
rm -rf plugins/satellite/dock/

echo "✅ Pro功能代码已删除"

# 5. 创建开源插件示例
echo ""
echo "5️⃣  创建开源插件示例..."

# 创建示例frontend插件
mkdir -p plugins/frontend/quick-tools/dist
cat > plugins/frontend/quick-tools/manifest.json << 'EOF'
{
  "id": "com.orbithub.quick-tools",
  "name": "快捷工具",
  "version": "1.0.0",
  "type": "frontend",
  "description": "常用的快捷操作工具",
  "author": "OrbitHub Team",
  "license": "MIT",
  "tier": "free",
  "entry": "./dist/index.js",
  "mounts": [
    {
      "slot": "Pool.Sidebar.Section",
      "component": "QuickTools",
      "title": "快捷工具",
      "icon": "⚡",
      "order": 100
    }
  ],
  "permissions": [
    "assets:read",
    "tags:write",
    "ui:notification"
  ]
}
EOF

echo "✅ 开源插件示例已创建"

# 6. 创建README
echo ""
echo "6️⃣  创建README..."
cat > README.md << 'EOF'
# OrbitHub

> 摆脱平台引力，进入你的自主轨道

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![GitHub release](https://img.shields.io/github/release/orbithub-work/orbit.svg)](https://github.com/orbithub-work/orbit/releases)

[English](README_EN.md) | 简体中文

## 💡 核心理念

**像开发软件一样做内容创作**

- 📦 **素材池** = 代码库（管理你的创作资产）
- 🎬 **工作台** = 开发项目（组织创作流程）
- 🎁 **成品库** = 发布版本（管理发布作品）
- 📊 **数据看板** = 监控面板（追踪创作数据）

## ✨ 完整的创作工作流

### 📦 素材池 (Pool)
导入、分类、搜索你的所有创作素材

- 🏷️ 多级标签系统
- 🔍 强大的搜索筛选
- ⭐ 评分管理
- 👁️ 缩略图预览
- 🔄 自动扫描监控

### 🎬 工作台 (Workspace)
像管理代码项目一样管理创作项目

- 📝 项目管理（策划 → 制作 → 完成）
- 🔗 关联素材
- 📋 项目笔记（Markdown）
- 📅 路线图
- 📊 进度追踪

### 🎁 成品库 (Artifact)
管理你发布的所有作品

- 📤 发布记录（B站、抖音、YouTube...）
- 🔗 关联源项目
- 📈 基础数据统计
- 🏷️ 成品分类

### 📊 数据看板 (Analytics)
了解你的创作数据

- 📊 素材统计
- 📈 项目进度
- 🎯 发布数据
- 📉 趋势分析

## 🚀 快速开始

### 下载安装

[下载最新版本](https://github.com/orbithub-work/orbit/releases)

### 从源码编译

```bash
# 克隆仓库
git clone https://github.com/orbithub-work/orbit.git
cd orbit

# 编译后端
go build -o bin/core cmd/core/main.go

# 编译前端
cd frontend
npm install
npm run build

# 启动
./bin/core
```

## 🔌 插件生态

OrbitHub支持插件扩展，打造专属工作流：

- 🎬 剪映草稿导入
- 📝 Notion同步
- 🤖 AI智能标签
- 🚀 一键发布

[浏览插件市场](https://orbithub.work/plugins)

## 💎 Pro版本

开源版已经提供完整工作流！Pro版提供更深度的功能：

- 🔒 **版权确权** - 数字签名保护你的创作
- 🔗 **血缘链追溯** - 追踪素材的来源和使用
- 📋 **自定义工作流** - 打造专属创作流程
- 🚀 **一键发布** - 同时发布到多个平台
- 📊 **高级数据分析** - AI驱动的数据洞察
- 🤝 **团队协作** - 多人协同创作

[了解Pro版本](https://orbithub.work/pricing) | 价格：¥299/年

## 🎯 适用人群

- 📹 视频创作者
- 📷 摄影师
- 🎨 设计师
- 🎬 剪辑师
- 📝 自媒体作者

## 📖 文档

- [架构设计](docs/ARCHITECTURE.md)
- [插件开发](docs/PLUGIN_DEVELOPMENT.md)
- [贡献指南](CONTRIBUTING.md)

## 🤝 贡献

欢迎贡献代码、开发插件、完善文档！

查看 [贡献指南](CONTRIBUTING.md)

## 📄 开源协议

[MIT License](LICENSE)

## 🌟 Star History

如果这个项目对你有帮助，请给个Star ⭐️

---

**让创作更有条理，让数据更有价值** ✨

[OrbitHub.work](https://orbithub.work)
EOF

echo "✅ README已创建"

# 7. 创建LICENSE
echo ""
echo "7️⃣  创建LICENSE..."
cat > LICENSE << 'EOF'
MIT License

Copyright (c) 2026 OrbitHub

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
EOF

echo "✅ LICENSE已创建"

# 8. 创建.gitignore
echo ""
echo "8️⃣  更新.gitignore..."
cat >> .gitignore << 'EOF'

# OrbitHub specific
/data/
/bin/
*.db
*.db-shm
*.db-wal
server.port

# Pro features (should not exist in open source)
/pro/
EOF

echo "✅ .gitignore已更新"

# 9. 第一次提交
echo ""
echo "9️⃣  第一次提交..."
git add .
git commit -m "🎉 Initial commit: OrbitHub open source release

- Complete asset management workflow
- Project workspace with basic workflow
- Artifact management
- Analytics dashboard
- Plugin system framework
- MIT License

OrbitHub.work - Break free from platform gravity"

# 10. 推送到GitHub
echo ""
echo "🔟 推送到GitHub..."
echo "⚠️  即将推送到 $GITHUB_REPO"
echo "   是否继续？(y/n)"
read -r response
if [ "$response" = "y" ]; then
    git branch -M main
    git push -u origin main
    echo "✅ 推送成功！"
else
    echo "⏸️  跳过推送，你可以稍后手动推送："
    echo "   cd $TARGET_DIR"
    echo "   git push -u origin main"
fi

# 完成
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ OrbitHub 开源版本抽离完成！"
echo ""
echo "📁 开源仓库位置: $TARGET_DIR"
echo "🔗 GitHub仓库: https://github.com/orbithub-work/orbit"
echo "🌐 官网: https://orbithub.work"
echo ""
echo "📝 下一步："
echo "1. 访问 https://github.com/orbithub-work/orbit"
echo "2. 完善仓库描述和Topics"
echo "3. 创建第一个Release"
echo "4. 部署官网到 orbithub-work.github.io"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
