# Phase 5: 文件传输中心（可选）

## 开发目标

实现重量级文件传输功能：
- 大文件上传（断点续传）
- 文件库管理
- 多节点分发
- 传输记录追踪

## 1. 定位对比

| 维度 | 文件管理 | 文件传输中心 |
|------|---------|-------------|
| 用途 | 实时操作 | 批量分发 |
| 文件大小 | < 50MB | 无限制 |
| 目标节点 | 单节点 | 多节点 |
| 断点续传 | ❌ | ✅ |
| 版本管理 | ❌ | ✅ |

## 2. 整体架构

```
┌─────────────────────────────────────────────────────────────────┐
│                     Frontend (/transfer)                         │
├─────────────────────────────────────────────────────────────────┤
│ APIServer                                                        │
│ ├── 文件存储服务                                                  │
│ ├── 传输任务队列                                                  │
│ └── 任务调度器                                                    │
├─────────────────────────────────────────────────────────────────┤
│ Agent 1  ←──┐                                                    │
│ Agent 2  ←──┼── 推送/拉取文件                                     │
│ Agent N  ←──┘                                                    │
└─────────────────────────────────────────────────────────────────┘
```

## 3. 数据模型

```go
// internal/squ-apiserver/model/transfer.go

// FileRecord 文件记录
type FileRecord struct {
    ID         uint      `json:"id"`
    UUID       string    `json:"uuid"`       // 文件唯一标识
    Name       string    `json:"name"`       // 文件名
    Size       int64     `json:"size"`       // 文件大小
    Path       string    `json:"path"`       // 存储路径
    Mime       string    `json:"mime"`       // MIME 类型
    MD5        string    `json:"md5"`        // 文件哈希
    Version    int       `json:"version"`    // 版本号
    Uploader   string    `json:"uploader"`   // 上传者
    ExpiresAt  *time.Time `json:"expires_at"` // 过期时间
    CreatedAt  time.Time `json:"created_at"`
}

// TransferTask 传输任务
type TransferTask struct {
    ID          uint      `json:"id"`
    FileID      uint      `json:"file_id"`      // 关联文件
    ServerID    uint      `json:"server_id"`    // 目标服务器
    TargetPath  string    `json:"target_path"`  // 目标路径
    Status      string    `json:"status"`       // pending/running/success/failed
    Progress    int       `json:"progress"`     // 进度 0-100
    ErrorMsg    string    `json:"error_msg"`    // 错误信息
    StartedAt   *time.Time `json:"started_at"`
    FinishedAt  *time.Time `json:"finished_at"`
    CreatedAt   time.Time `json:"created_at"`
}
```

## 4. API 设计

### APIServer 端

```
# 文件库管理
POST   /api/v1/transfer/upload      # 上传文件（支持分片）
GET    /api/v1/transfer/files       # 文件列表
DELETE /api/v1/transfer/files/{id}  # 删除文件

# 分发任务
POST   /api/v1/transfer/dispatch    # 创建分发任务
GET    /api/v1/transfer/tasks       # 任务列表
GET    /api/v1/transfer/tasks/{id}  # 任务详情
DELETE /api/v1/transfer/tasks/{id}  # 取消任务

# 分片上传
POST   /api/v1/transfer/chunk/init      # 初始化分片上传
POST   /api/v1/transfer/chunk/upload    # 上传分片
POST   /api/v1/transfer/chunk/complete  # 合并分片
```

### Agent 端

```
# 接收文件
POST /api/agent/v1/fs/receive   # 接收推送的文件
GET  /api/agent/v1/fs/pull      # 拉取文件
```

## 5. 分片上传流程

```
1. 初始化
   POST /chunk/init { name, size, md5 }
   ← { upload_id, chunk_size, chunk_count }

2. 上传分片（并行）
   POST /chunk/upload { upload_id, chunk_index, chunk_data }
   ← { received: true }

3. 合并分片
   POST /chunk/complete { upload_id }
   ← { file_id, path }
```

## 6. 分发流程

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   上传文件   │ ──► │  创建任务    │ ──► │  任务队列    │
└─────────────┘     └─────────────┘     └─────────────┘
                                               │
                                               ▼
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   完成      │ ◄── │  更新进度    │ ◄── │  推送到Agent │
└─────────────┘     └─────────────┘     └─────────────┘
```

## 7. 前端页面

```
┌──────────────────────────────────────────────────────────────────┐
│ 📦 文件传输中心                                [上传文件]        │
├──────────────────────────────────────────────────────────────────┤
│ 📋 传输任务                                                       │
│ ──────────────────────────────────────────────────────────────── │
│ app-v2.0.tar.gz  150MB   ━━━━━━━░░  80%   3/5节点   [取消]       │
│ config.zip       2MB     ━━━━━━━━━━  100%  5/5节点   ✓ 完成      │
├──────────────────────────────────────────────────────────────────┤
│ 📂 文件库                              [清理过期]               │
│ ──────────────────────────────────────────────────────────────── │
│ app-v2.0.tar.gz  150MB  2024-01-15  [分发] [下载] [删除]         │
│ nginx.conf       5KB    2024-01-14  [分发] [下载] [删除]         │
└──────────────────────────────────────────────────────────────────┘

[分发弹窗]
┌──────────────────────────────────────────────────────────────────┐
│ 分发文件: app-v2.0.tar.gz                                        │
│ ──────────────────────────────────────────────────────────────── │
│ 目标服务器:  [✓] server-1  [✓] server-2  [ ] server-3            │
│ 目标路径:    [/home/app/_________________________]               │
│                                                                  │
│                              [取消] [开始分发]                   │
└──────────────────────────────────────────────────────────────────┘
```

## 8. 目录结构

```
internal/squ-apiserver/
├── handler/
│   └── transfer/
│       ├── upload.go       # 文件上传
│       ├── chunk.go        # 分片上传
│       ├── file.go         # 文件库管理
│       └── dispatch.go     # 分发任务
├── model/
│   └── transfer.go         # 数据模型
├── repository/
│   └── transfer.go         # 数据访问
└── service/
    └── dispatcher.go       # 任务调度器

front/src/
└── views/
    └── TransferCenter/
        ├── index.vue       # 主页面
        └── components/
            ├── FileLibrary.vue    # 文件库
            ├── TaskList.vue       # 任务列表
            ├── DispatchDialog.vue # 分发弹窗
            └── ChunkUpload.vue    # 分片上传
```

## 9. 配置设计

```yaml
# config/apiserver.yaml
transfer:
  enabled: true
  storage:
    path: /data/squirrel/uploads    # 存储路径
    max_size: 10GB                  # 最大存储
    expire_days: 7                  # 过期天数
  upload:
    max_file_size: 1024             # MB, 0=无限制
    chunk_size: 5                   # MB, 分片大小
    concurrent_chunks: 3            # 并行上传数
  dispatch:
    concurrent_tasks: 5             # 并行任务数
    retry_times: 3                  # 失败重试次数
```

## 10. 开发步骤

| 步骤 | 内容 | 预计代码量 |
|------|------|-----------|
| 1 | 数据模型 + 数据库迁移 | ~100 行 |
| 2 | 文件上传（普通） | ~150 行 |
| 3 | 分片上传 | ~200 行 |
| 4 | 文件库管理 | ~100 行 |
| 5 | 分发任务 + 调度器 | ~250 行 |
| 6 | Agent 接收接口 | ~80 行 |
| 7 | 前端页面 | ~400 行 |
| **总计** | | **~1280 行** |

## 11. 注意事项

1. **存储清理**：定时任务清理过期文件
2. **任务去重**：相同文件 + 相同目标避免重复分发
3. **断点续传**：记录已上传分片，支持中断恢复
4. **进度同步**：WebSocket 推送实时进度
5. **错误处理**：失败重试 + 告警通知

## 12. 完成检查

- [ ] 数据模型设计完成
- [ ] 文件上传功能实现
- [ ] 分片上传实现
- [ ] 文件库管理实现
- [ ] 分发任务实现
- [ ] 任务调度器实现
- [ ] Agent 接收接口
- [ ] 前端页面实现
- [ ] 功能测试通过
