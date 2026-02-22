# 目录结构

## 1. 概述

本项目采用标准的 Go 项目布局，遵循社区最佳实践，将代码按职责分层组织。

## 2. 顶层目录结构

```
squirrel-dev/
├── cmd/                    # 应用程序入口
├── internal/               # 私有应用代码
├── pkg/                    # 公共库代码
├── api-rest/               # REST API 定义文件
├── config/                 # 配置文件
├── docs/                   # 文档
├── front/                  # 前端代码
├── dockerfiles/            # Docker 构建文件
├── compose/                # Docker Compose 配置
├── db/                     # 本地数据库文件
├── log/                    # 日志目录
├── squirrel/               # 编译产物目录
├── go.mod                  # Go 模块定义
├── go.sum                  # Go 依赖校验
├── Makefile                # 构建脚本
└── README.md               # 项目说明
```

## 3. cmd 目录 - 应用入口

`cmd` 目录存放每个可执行程序的入口代码，每个子目录对应一个独立的程序。

### 3.1 结构示例

```
cmd/
├── squ-agent/              # Agent 程序入口
│   ├── agent.go            # main 函数
│   └── app/                # 应用初始化
│       ├── server.go
│       └── options/
│           └── options.go
├── squ-apiserver/          # API Server 程序入口
│   ├── apiserver.go        # main 函数
│   └── app/
│       ├── server.go
│       └── options/
│           └── options.go
└── squctl/                 # CLI 工具入口
    ├── squctl.go           # main 函数
    └── app/
        ├── server.go
        ├── certs.go
        └── options/
            └── options.go
```

### 3.2 入口文件规范

入口文件（如 `agent.go`、`apiserver.go`）应保持简洁，仅负责：

1. 调用应用初始化函数
2. 执行命令
3. 处理错误

```go
// cmd/squ-apiserver/apiserver.go
package main

import (
    "log"

    "squirrel-dev/cmd/squ-apiserver/app"
)

func main() {
    cmd := app.NewServerCommand()
    if err := cmd.Execute(); err != nil {
        log.Fatalln(err)
    }
}
```

### 3.3 app 子目录

`app` 目录负责：

- `server.go`: 定义命令行入口和初始化流程
- `options/`: 命令行参数和配置选项定义

## 4. internal 目录 - 业务逻辑

`internal` 目录存放私有应用代码，Go 编译器会阻止其他项目导入此目录下的代码。

### 4.1 结构示例

```
internal/
├── pkg/                    # 内部共享包
│   ├── database/           # 数据库连接
│   ├── middleware/         # HTTP 中间件
│   ├── migration/          # 数据库迁移
│   └── response/           # 统一响应格式
├── squ-agent/              # Agent 业务代码
├── squ-apiserver/          # API Server 业务代码
└── squctl/                 # CLI 工具业务代码
```

### 4.2 业务模块分层

每个业务模块（如 `squ-apiserver`）采用标准的分层架构：

```
internal/squ-apiserver/
├── config/                 # 配置定义与加载
├── model/                  # 数据模型定义
├── repository/             # 数据访问层
├── handler/                # 业务处理层
├── router/                 # 路由注册
├── server/                 # 服务器启动逻辑
├── cron/                   # 定时任务
└── terminal/               # 终端功能
```

### 4.3 各层职责

| 层级 | 目录 | 职责 |
|------|------|------|
| 入口层 | `cmd/*/app/` | 命令行解析、应用初始化 |
| 配置层 | `config/` | 配置结构定义、配置加载 |
| 路由层 | `router/` | HTTP 路由注册、路由分组 |
| 处理层 | `handler/` | 请求处理、业务逻辑编排 |
| 模型层 | `model/` | 数据结构定义、数据库映射 |
| 存储层 | `repository/` | 数据库操作、CRUD 封装 |
| 服务层 | `server/` | 服务启动、依赖注入 |

### 4.4 handler 目录组织

当 handler 较多时，按资源/功能划分子目录：

```
handler/
├── app_store/
│   ├── handler.go          # Handler 结构体定义
│   ├── service.go          # 业务逻辑
│   └── common.go           # 共用函数
├── application/
│   ├── handler.go
│   └── service.go
├── auth/
├── deployment/
├── monitor/
├── script/
└── server/
```

## 5. pkg 目录 - 公共库

`pkg` 目录存放可被外部项目引用的公共代码。

```
pkg/
├── cert/                   # 证书工具
├── collector/              # 系统指标采集
├── compose/                # Docker Compose 操作
├── docker/                 # Docker 客户端封装
├── email/                  # 邮件发送
├── execute/                # 命令执行
├── file/                   # 文件操作
├── hash/                   # 密码哈希
├── httpclient/             # HTTP 客户端
├── jwt/                    # JWT 工具
├── k8s/                    # Kubernetes 客户端
├── ssh/                    # SSH 客户端
└── utils/                  # 通用工具函数
```

## 6. 目录选择原则

### 6.1 internal vs pkg

| 场景 | 选择 |
|------|------|
| 仅项目内部使用的代码 | `internal/` |
| 可能被其他项目复用的代码 | `pkg/` |
| 业务特定逻辑 | `internal/{module}/` |
| 通用工具函数 | `pkg/` |

### 6.2 cmd vs internal

| 场景 | 选择 |
|------|------|
| 程序入口（main 函数） | `cmd/{program}/` |
| 入口初始化逻辑 | `cmd/{program}/app/` |
| 核心业务逻辑 | `internal/{module}/` |

## 7. 命名规范

1. **目录名**: 使用小写，单词间用 `-` 连接（如 `squ-agent`）
2. **文件名**: 使用小写，单词间用 `_` 连接（如 `server_types.go`）
3. **包名**: 使用小写单词，简洁有意义（如 `handler`、`model`）
4. **入口文件**: 与目录名相关，如 `squ-agent/agent.go`
