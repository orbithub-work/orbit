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
