# 文件管理模块设计概述

## 1. 模块定位

文件管理模块采用统一入口设计，所有文件先上传到 APIServer 文件库，然后可选分发到 Agent 节点：

|| 页面 | 路由 | 定位 | 说明 |
||------|------|------|------|
|| **文件管理** | `/files` | APIServer 文件库 | 上传、管理 APIServer 上的文件，触发分发 |
|| **分发管理** | `/dispatch` | Agent 节点视图 | 查看各节点文件状态、同步情况 |

## 2. 整体架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          Frontend (Vue3)                                     │
│  ├── /files      - 文件管理页面（APIServer 文件库）                          │
│  └── /dispatch   - 分发管理页面（Agent 节点视图）                            │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │ HTTP
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                     squ-apiserver (核心节点)                                  │
│  ├── 文件存储服务        - 统一文件库，支持断点续传                          │
│  ├── 分发任务管理        - 创建分发任务，跟踪进度                            │
│  └── Agent 通信          - 推送文件到目标节点                                │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │ mTLS / HTTP
              ┌─────────────────────┼─────────────────────┐
              ▼                     ▼                     ▼
┌──────────────────────┐ ┌──────────────────────┐ ┌──────────────────────┐
│   squ-agent (主机1)   │ │   squ-agent (主机2)   │ │   squ-agent (主机N)   │
├──────────────────────┤ ├──────────────────────┤ ├──────────────────────┤
│ • 文件系统操作        │ │ • 文件系统操作        │ │ • 文件系统操作        │
│ • 沙箱路径限制        │ │ • 沙箱路径限制        │ │ • 沙箱路径限制        │
│ • 接收推送文件        │ │ • 接收推送文件        │ │ • 接收推送文件        │
└──────────────────────┘ └──────────────────────┘ └──────────────────────┘
```

## 3. 核心流程

### 3.1 文件上传流程

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   前端      │ ──► │  APIServer  │ ──► │   文件库    │
└─────────────┘     └─────────────┘     └─────────────┘
                          │
                    （可选分发）
                          ▼
                    ┌─────────────┐
                    │   Agents    │
                    └─────────────┘
```

1. 用户在**文件管理**页面上传文件到 APIServer
2. 文件存储在 APIServer 文件库中
3. 可选择分发到指定的 Agent 节点

### 3.2 分发流程

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  选择文件   │ ──► │  选择节点   │ ──► │  创建任务   │
└─────────────┘     └─────────────┘     └─────────────┘
                                              │
                                              ▼
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   完成      │ ◄── │  更新进度   │ ◄── │  推送文件   │
└─────────────┘     └─────────────┘     └─────────────┘
```

## 4. 功能范围

### 4.1 文件管理页面

|| 功能 | Phase | 说明 |
||------|-------|------|
|| 文件上传 | 1 | 上传文件到 APIServer（支持断点续传） |
|| 文件列表 | 1 | 查看 APIServer 文件库 |
|| 文件删除 | 1 | 删除文件库中的文件 |
|| 触发分发 | 2 | 选择文件分发到 Agent 节点 |
|| 文件预览 | 2 | 预览文本文件内容 |
|| 下载文件 | 2 | 从 APIServer 下载文件 |

### 4.2 分发管理页面（Phase 3）

|| 功能 | 说明 |
||------|------|
|| 服务器列表 | 查看所有在线 Agent 节点 |
|| 文件状态 | 查看各节点文件同步状态 |
|| 分发任务 | 查看分发任务进度和历史 |
|| 一致性检查 | 检查文件在各节点的一致性 |

## 5. 开发阶段规划

```
Phase 1: Agent 核心功能
├── 01-agent-core.md    # 目录浏览 + 文件读取 + 安全校验
└── 完成后可独立测试

Phase 2: Agent 写入功能
├── 02-agent-write.md   # 文件写入 + 上传 + 接收推送
└── 完成后 Agent 端功能完整

Phase 3: APIServer 文件存储
├── 03-apiserver-storage.md  # 文件存储 + 分片上传 + 分发任务
└── 完成后 APIServer 文件库可用

Phase 4: 前端页面
├── 04-frontend-filemanager.md  # 文件管理页面
└── 完成后用户可用

Phase 5: 分发管理（可选）
├── 05-dispatch-management.md  # 分发管理页面
└── 按需开发
```

## 6. 数据模型概览

### 6.1 APIServer 数据模型

```go
// FileRecord 文件记录（APIServer 文件库）
type FileRecord struct {
    ID         uint      `json:"id"`
    UUID       string    `json:"uuid"`       // 文件唯一标识
    Name       string    `json:"name"`       // 文件名
    Size       int64     `json:"size"`       // 文件大小
    Path       string    `json:"path"`       // 存储路径
    MD5        string    `json:"md5"`        // 文件哈希
    Status     string    `json:"status"`     // uploaded/distributing/distributed
    CreatedAt  time.Time `json:"created_at"`
}

// DispatchTask 分发任务
type DispatchTask struct {
    ID          uint      `json:"id"`
    FileID      uint      `json:"file_id"`      // 关联文件
    ServerID    uint      `json:"server_id"`    // 目标服务器
    TargetPath  string    `json:"target_path"`  // 目标路径
    Status      string    `json:"status"`       // pending/running/success/failed
    Progress    int       `json:"progress"`     // 进度 0-100
    ErrorMsg    string    `json:"error_msg"`    // 错误信息
    CreatedAt   time.Time `json:"created_at"`
    FinishedAt  *time.Time `json:"finished_at"`
}

// ChunkUpload 分片上传记录
type ChunkUpload struct {
    ID         uint   `json:"id"`
    UploadID   string `json:"upload_id"`   // 上传会话 ID
    FileName   string `json:"file_name"`
    FileSize   int64  `json:"file_size"`
    FileMD5    string `json:"file_md5"`
    ChunkSize  int64  `json:"chunk_size"`
    ChunkCount int    `json:"chunk_count"`
    Chunks     string `json:"chunks"`      // JSON: [true, false, true, ...]
    Status     string `json:"status"`      // uploading/completed
}
```

### 6.2 Agent 数据模型

```go
// FileInfo 文件/目录信息
type FileInfo struct {
    Name      string    `json:"name"`       // 文件名
    Path      string    `json:"path"`       // 完整路径
    IsDir     bool      `json:"is_dir"`     // 是否目录
    Size      int64     `json:"size"`       // 大小（字节）
    Mode      string    `json:"mode"`       // 权限
    ModTime   time.Time `json:"mod_time"`   // 修改时间
}
```

## 7. API 设计概览

### 7.1 文件管理 API

```
POST   /api/v1/files/upload           # 上传文件
POST   /api/v1/files/chunk/init       # 初始化分片上传
POST   /api/v1/files/chunk/upload     # 上传分片
POST   /api/v1/files/chunk/complete   # 合并分片
GET    /api/v1/files                  # 文件列表
DELETE /api/v1/files/{id}             # 删除文件
GET    /api/v1/files/{id}/download    # 下载文件

POST   /api/v1/dispatch               # 创建分发任务
GET    /api/v1/dispatch/tasks         # 分发任务列表
GET    /api/v1/dispatch/tasks/{id}    # 任务详情
```

### 7.2 分发管理 API

```
GET    /api/v1/servers               # 服务器列表
GET    /api/v1/servers/files         # 所有服务器文件状态
GET    /api/v1/servers/{id}/files    # 指定服务器文件列表
```

### 7.3 Agent 内部 API

```
GET    /api/agent/v1/fs/list         # 目录列表
GET    /api/agent/v1/fs/read         # 文件读取
GET    /api/agent/v1/fs/download     # 文件下载
POST   /api/agent/v1/fs/write        # 文件写入
POST   /api/agent/v1/fs/upload       # 文件上传
POST   /api/agent/v1/fs/receive      # 接收推送文件
POST   /api/agent/v1/fs/mkdir        # 创建目录
POST   /api/agent/v1/fs/delete       # 删除
POST   /api/agent/v1/fs/rename       # 重命名
```

## 8. 文档索引

|| 文档 | 内容 | 开发阶段 |
||------|------|----------|
|| [01-agent-core.md](./01-agent-core.md) | Agent 核心实现 | Phase 1 |
|| [02-agent-write.md](./02-agent-write.md) | Agent 写入操作 | Phase 2 |
|| [03-apiserver-storage.md](./03-apiserver-storage.md) | APIServer 文件存储 | Phase 3 |
|| [04-frontend-filemanager.md](./04-frontend-filemanager.md) | 前端页面 | Phase 4 |
|| [05-dispatch-management.md](./05-dispatch-management.md) | 分发管理页面 | Phase 5 |

## 9. 技术栈

- **后端**：Go + Gin + GORM
- **前端**：Vue 3 + TypeScript + Monaco Editor
- **传输**：HTTP 流式传输 + 分片上传
- **安全**：路径沙箱 + mTLS

## 10. 相关参考

- 项目目录结构：`docs/zh-CN/guide/directory-structure.md`
- Agent 加入设计：`docs/zh-CN/develop/agent_join_design.md`
- 监控模块设计：`docs/zh-CN/develop/monitor/01-monitor.md`
