# Squirrel Dev

一个轻量级、现代化的运维管理平台，支持服务器管理、应用部署、监控告警和脚本执行等功能。

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.25%2B-brightgreen.svg)](https://golang.org)
[![Vue Version](https://img.shields.io/badge/Vue-3.4%2B-green.svg)](https://vuejs.org)

## 功能特性

- **服务器管理** - 统一管理多台服务器，支持 Web 终端连接
- **应用部署** - 支持 Docker Compose 应用一键部署、启动、停止
- **应用商店** - 内置常用应用模板，快速部署常用服务
- **监控告警** - 实时采集服务器 CPU、内存、磁盘、网络等资源使用情况
- **脚本管理** - 定时任务脚本管理，支持 Cron 表达式调度
- **配置中心** - 集中管理应用配置，支持动态更新
- **命令行工具** - 提供 squctl CLI 工具，方便运维操作

## 系统架构

Squirrel 采用有代理（Agent）架构设计：

```
┌─────────────┐      HTTP       ┌─────────────┐      HTTP       ┌─────────────┐
│   Frontend  │ ◄─────────────► │  API Server │ ◄─────────────► │    Agent    │
│  (Vue3/TS)  │                 │  (控制台)    │                 │  (客户端)    │
└─────────────┘                 └─────────────┘                 └─────────────┘
                                                                      │
                                                                      ▼
                                                               ┌─────────────┐
                                                               │   Servers   │
                                                               └─────────────┘
```

### 核心组件

| 组件 | 描述 | 默认端口 |
|------|------|----------|
| `squ-apiserver` | 控制台服务端，提供 API 接口和前端界面 | 10700 |
| `squ-agent` | 客户端代理，部署在目标服务器上执行具体操作 | 10750 |
| `squctl` | 命令行工具，用于与 apiserver 交互 | - |

## 技术栈

### 后端
- **Go 1.25+** - 主要开发语言
- **Gin** - HTTP Web 框架
- **GORM** - ORM 数据库框架
- **JWT** - 用户认证
- **WebSocket** - 实时终端连接
- **Cobra** - CLI 命令行框架
- **Viper** - 配置管理
- **Zap** - 日志记录

### 前端
- **Vue 3** - 渐进式 JavaScript 框架
- **TypeScript** - 类型安全的 JavaScript
- **Vite** - 前端构建工具
- **Pinia** - 状态管理
- **Vue Router** - 路由管理
- **Vue I18n** - 国际化支持
- **Sass** - CSS 预处理器

### 数据库
- **SQLite** - 默认嵌入式数据库（零配置）
- **MySQL** - 可选外置数据库

## 快速开始

### 环境要求

- Go 1.25 或更高版本
- Node.js 18+ 和 npm
- Docker 和 Docker Compose（可选，用于容器化部署）

### 构建安装

```bash
# 克隆项目
git clone https://github.com/yourusername/squirrel-dev.git
cd squirrel-dev

# 一键构建（前端 + 后端）
make all

# 打包输出
make package

# 构建 Docker 镜像
make image
```

构建完成后，二进制文件和配置将位于 `squirrel/` 目录下。

### 手动启动

```bash
# 1. 启动 API Server
./squirrel/squ-apiserver --config ./squirrel/config/apiserver.yaml

# 2. 在目标服务器上启动 Agent
./squirrel/squ-agent --config ./squirrel/config/agent.yaml

# 3. 使用 squctl 连接
./squirrel/squctl login http://localhost:10700
```

### Docker 部署

```bash
# 构建镜像
make image

# 使用 Docker Compose 启动
docker-compose up -d
```

## 配置说明

### API Server 配置 (`config/apiserver.yaml`)

```yaml
server:
  bind: "0.0.0.0"
  port: 10700
  mode: debug

db:
  type: sqlite  # 可选: mysql 或 sqlite
  sqlite:
    filePath: ./db/apiserver.db

log:
  path: ./log/apiserver
  level: info
```

### Agent 配置 (`config/agent.yaml`)

```yaml
server:
  bind: "0.0.0.0"
  port: 10750

apiserver:
  http:
    scheme: http
    server: 127.0.0.1:10700
    baseUri: /api/v1
```

## 项目结构

```
squirrel-dev/
├── cmd/                    # 程序入口
│   ├── squ-apiserver/      # API 服务端入口
│   ├── squ-agent/          # Agent 代理入口
│   └── squctl/             # CLI 工具入口
├── internal/               # 内部实现
│   ├── squ-apiserver/      # API Server 业务逻辑
│   ├── squ-agent/          # Agent 业务逻辑
│   └── squctl/             # CLI 业务逻辑
├── pkg/                    # 公共包
│   └── collector/          # 监控数据采集
├── front/                  # 前端源码 (Vue3)
├── config/                 # 配置文件模板
├── dockerfiles/            # Docker 构建文件
└── api-rest/               # API 测试请求
```

## API 接口

项目提供完整的 RESTful API 接口，详见 `api-rest/` 目录：

- `server.http` - 服务器管理接口
- `monitor.http` - 监控数据接口
- `application.http` - 应用管理接口
- `deployment.http` - 部署管理接口
- `script.http` - 脚本管理接口

## 开发计划

- [x] 服务器管理和 Web 终端
- [x] Docker Compose 应用部署
- [x] 服务器资源监控
- [x] 定时脚本任务
- [ ] Kubernetes 集群支持
- [ ] 告警通知（邮件/钉钉/企业微信）
- [ ] 日志收集与分析
- [ ] 多租户权限管理

## 贡献指南

欢迎提交 Issue 和 Pull Request！

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

## 许可证

本项目采用 [Apache 2.0](LICENSE) 许可证开源。

```
Copyright 2026 agocan

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0
```

## 致谢

感谢所有为本项目做出贡献的开发者。

---

如果本项目对您有帮助，请给个 Star ⭐️ 支持一下！
