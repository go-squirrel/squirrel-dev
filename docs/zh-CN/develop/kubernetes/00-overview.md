# Kubernetes 支持设计概览

## 1. 背景

当前系统主要支持 Docker Compose 应用部署，越来越多场景需要管理 Kubernetes 集群：
- 容器化应用部署到 K8s 集群
- 查看和管理 K8s 资源
- Pod 日志和终端访问

## 2. 设计目标

- **集群管理**：支持多 K8s 集群接入
- **Workload 管理**：Deployment、StatefulSet、DaemonSet、Job、CronJob
- **资源查看**：Pod、Service、ConfigMap、Secret、PVC 等
- **运维操作**：Pod 日志、终端、重启、扩缩容

## 3. 整体架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         Kubernetes 支持架构                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   ┌────────────────────────────────────────────────────────────────────┐    │
│   │                          Frontend (Vue)                             │    │
│   │  ├── 集群管理页面                                                   │    │
│   │  ├── Workload 管理页面                                              │    │
│   │  ├── Pod 详情页面 (日志/终端)                                        │    │
│   │  └── 资源管理页面 (Service/ConfigMap/Secret)                        │    │
│   └────────────────────────────────────────────────────────────────────┘    │
│                                    │                                         │
│                                    ▼                                         │
│   ┌────────────────────────────────────────────────────────────────────┐    │
│   │                        squ-apiserver                                │    │
│   │  ├── K8s Client 管理 (多集群)                                       │    │
│   │  ├── Workload 操作 API                                              │    │
│   │  └── WebSocket 代理 (日志/终端)                                     │    │
│   └────────────────────────────────────────────────────────────────────┘    │
│                                    │                                         │
│              ┌─────────────────────┼─────────────────────┐                  │
│              ▼                     ▼                     ▼                  │
│   ┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐          │
│   │  K8s Cluster 1   │  │  K8s Cluster 2   │  │  K8s Cluster N   │          │
│   │  (kubeconfig)    │  │  (kubeconfig)    │  │  (kubeconfig)    │          │
│   └──────────────────┘  └──────────────────┘  └──────────────────┘          │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

## 4. 核心功能

| 功能模块 | 说明 |
|---------|------|
| **集群管理** | 接入 K8s 集群，管理 kubeconfig |
| **命名空间** | 命名空间查看和切换 |
| **Workload** | Deployment/StatefulSet/DaemonSet 的查看、创建、更新、删除 |
| **Pod 管理** | Pod 列表、详情、日志、终端、重启 |
| **配置管理** | ConfigMap、Secret 的查看和管理 |
| **网络管理** | Service、Ingress 的查看和管理 |
| **存储管理** | PVC、PV、StorageClass 的查看 |

## 5. 数据模型概要

```go
// K8s 集群配置
type K8sCluster struct {
    ID          uint
    Name        string    // 集群名称
    Alias       string    // 显示名称
    Kubeconfig  string    // kubeconfig 内容 (加密存储)
    APIServer   string    // API Server 地址
    Enabled     bool
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

// 集群访问凭证由 APIServer 统一管理，通过 RBAC 控制访问权限
```

## 6. 开发阶段规划

```
Phase 1: 集群接入
├── 集群模型和 CRUD
├── K8s Client 初始化
├── 健康检查和连接状态
└── 命名空间列表

Phase 2: Workload 管理
├── Deployment 列表和详情
├── Deployment 创建/更新/删除
├── 扩缩容操作
├── 重启操作
└── YAML 编辑

Phase 3: Pod 操作
├── Pod 列表和详情
├── Pod 日志 (WebSocket)
├── Pod 终端 (WebSocket)
├── Pod 删除
└── 容器切换

Phase 4: 资源管理
├── Service 管理
├── ConfigMap 管理
├── Secret 管理
└── PVC 管理

Phase 5: 前端集成
├── 集群管理页面
├── Workload 管理页面
├── Pod 详情页面
└── YAML 编辑器
```

## 7. 与现有模块的关系

```
Server (服务器管理) ──┬──► Docker Compose 应用 (现有)
                     │
                     └──► K8s Cluster (新增)
                               │
                               ▼
RBAC (权限控制) ◄─── K8s 资源操作权限
                               │
                               ▼
Terminal (终端) ◄─── Pod Exec (复用 WebSocket)
```

## 8. 技术选型

| 组件 | 选择 | 原因 |
|------|------|------|
| K8s Client | client-go | 官方 SDK，功能完整 |
| YAML 编辑 | Monaco Editor | 复用文件管理模块编辑器 |
| WebSocket | gin-gonic/websocket | 复用终端模块实现 |

## 9. 设计要点

### 9.1 多集群支持
- 每个集群存储独立的 kubeconfig
- 支持多个 kubeconfig context
- 集群访问凭证加密存储

### 9.2 权限控制
- 集成 RBAC，控制集群访问权限
- 可限制用户只能访问特定命名空间
- 敏感操作（删除、更新）需要更高权限

### 9.3 前端设计
- 集群切换选择器
- 命名空间选择器
- 资源类型导航
- YAML 编辑器 + 表单编辑器双模式

## 10. 待详细设计的问题

- [ ] kubeconfig 如何安全存储？（加密方案）
- [ ] 是否支持 ServiceAccount 方式接入？
- [ ] 多集群资源如何聚合展示？
- [ ] 是否需要 Helm Chart 支持？
- [ ] 日志是否需要持久化存储？
- [ ] 是否支持跨集群应用分发？
