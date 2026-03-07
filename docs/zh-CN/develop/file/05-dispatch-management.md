# Phase 5: 分发管理页面（可选）

## 开发目标

实现分发管理页面，用于查看：
- 所有 Agent 节点状态
- 各节点文件同步情况
- 分发任务进度和历史

## 1. 页面定位

|| 页面 | 路由 | 数据来源 | 用途 |
||------|------|----------|------|
|| 文件管理 | `/files` | APIServer 文件库 | 上传、管理源文件 |
|| 分发管理 | `/dispatch` | Agent 节点 | 查看目标节点状态 |

## 2. 整体架构

```
┌─────────────────────────────────────────────────────────────────┐
│                     Frontend (/dispatch)                         │
├─────────────────────────────────────────────────────────────────┤
│ APIServer                                                        │
│ ├── 节点状态聚合                                                  │
│ ├── 分发任务查询                                                  │
│ └── 文件一致性检查                                                │
├─────────────────────────────────────────────────────────────────┤
│ Agent 1  ──── 返回文件列表                                        │
│ Agent 2  ──── 返回文件列表                                        │
│ Agent N  ──── 返回文件列表                                        │
└─────────────────────────────────────────────────────────────────┘
```

## 3. 数据模型扩展

### 3.1 服务器文件状态

```go
// internal/squ-apiserver/handler/dispatch/types.go

// ServerFileStatus 服务器文件状态
type ServerFileStatus struct {
    ServerID    uint   `json:"server_id"`
    ServerName  string `json:"server_name"`
    ServerIP    string `json:"server_ip"`
    Status      string `json:"status"`       // online/offline
    
    // 文件同步状态
    FileStatus  []FileSyncStatus `json:"file_status"`
}

// FileSyncStatus 文件同步状态
type FileSyncStatus struct {
    FileID      uint   `json:"file_id"`
    FileName    string `json:"file_name"`
    FileMD5     string `json:"file_md5"`
    
    // 节点上的状态
    Exists      bool   `json:"exists"`       // 文件是否存在
    Synced      bool   `json:"synced"`       // MD5 是否一致
    LocalMD5    string `json:"local_md5"`    // 节点上的 MD5
    TargetPath  string `json:"target_path"`  // 目标路径
}
```

### 3.2 APIServer 扩展 API

```
# 分发管理 API
GET  /api/v1/servers                    # 服务器列表（含文件状态概览）
GET  /api/v1/servers/:id/files          # 指定服务器的文件列表
GET  /api/v1/servers/files              # 所有服务器文件状态聚合

# 分发任务管理
GET  /api/v1/dispatch/tasks             # 分发任务列表（已实现）
GET  /api/v1/dispatch/tasks/:id         # 任务详情（已实现）
POST /api/v1/dispatch/tasks/:id/retry   # 重试失败任务
```

## 4. 后端实现

### 4.1 Handler

```go
// internal/squ-apiserver/handler/dispatch/handler.go

package dispatch

import (
    "github.com/gin-gonic/gin"
    
    "squirrel-dev/internal/squ-apiserver/repository"
    "squirrel-dev/internal/squ-apiserver/service"
)

// Handler 分发管理 Handler
type Handler struct {
    serverRepo  repository.ServerRepository
    fileRepo    repository.FileRepository
    taskRepo    repository.DispatchRepository
    agentClient *service.AgentClient
}

// NewHandler 创建 Handler
func NewHandler(
    serverRepo repository.ServerRepository,
    fileRepo repository.FileRepository,
    taskRepo repository.DispatchRepository,
) *Handler {
    return &Handler{
        serverRepo:  serverRepo,
        fileRepo:    fileRepo,
        taskRepo:    taskRepo,
        agentClient: service.NewAgentClient(),
    }
}

// RegisterRoutes 注册路由
func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
    servers := r.Group("/servers")
    {
        servers.GET("", h.ListServers)
        servers.GET("/files", h.ListServerFiles)
        servers.GET("/:id/files", h.GetServerFiles)
    }
    
    dispatch := r.Group("/dispatch")
    {
        dispatch.GET("/tasks", h.ListTasks)
        dispatch.GET("/tasks/:id", h.GetTask)
        dispatch.POST("/tasks/:id/retry", h.RetryTask)
    }
}
```

### 4.2 服务器列表

```go
// internal/squ-apiserver/handler/dispatch/server.go

package dispatch

import (
    "github.com/gin-gonic/gin"
    
    "squirrel-dev/internal/pkg/response"
)

// ServerWithStatus 带状态的服务器信息
type ServerWithStatus struct {
    ID          uint   `json:"id"`
    Name        string `json:"name"`
    IPAddress   string `json:"ip_address"`
    Status      string `json:"status"`
    AgentPort   int    `json:"agent_port"`
    
    // 分发统计
    TotalFiles  int    `json:"total_files"`   // 分发到此节点的文件数
    SyncedFiles int    `json:"synced_files"`  // 已同步文件数
    FailedFiles int    `json:"failed_files"`  // 失败文件数
}

// ListServers 服务器列表
func (h *Handler) ListServers(c *gin.Context) {
    // 获取服务器列表
    servers, err := h.serverRepo.List()
    if err != nil {
        response.Error(c, 500, "获取服务器列表失败: "+err.Error())
        return
    }
    
    // 组装状态信息
    var result []ServerWithStatus
    for _, s := range servers {
        // 统计分发任务
        tasks, _ := h.taskRepo.GetByServerID(s.ID)
        
        total := len(tasks)
        synced := 0
        failed := 0
        
        for _, t := range tasks {
            if t.Status == "success" {
                synced++
            } else if t.Status == "failed" {
                failed++
            }
        }
        
        result = append(result, ServerWithStatus{
            ID:          s.ID,
            Name:        s.Name,
            IPAddress:   s.IPAddress,
            Status:      s.Status,
            AgentPort:   s.AgentPort,
            TotalFiles:  total,
            SyncedFiles: synced,
            FailedFiles: failed,
        })
    }
    
    response.Success(c, result)
}
```

### 4.3 服务器文件状态

```go
// internal/squ-apiserver/handler/dispatch/files.go

package dispatch

import (
    "fmt"
    "io/ioutil"
    "net/http"
    
    "github.com/gin-gonic/gin"
    
    "squirrel-dev/internal/pkg/response"
)

// ListServerFiles 所有服务器文件状态
func (h *Handler) ListServerFiles(c *gin.Context) {
    // 获取所有服务器
    servers, _ := h.serverRepo.List()
    
    // 获取所有文件
    files, _ := h.fileRepo.List(1, 1000, "")
    
    // 获取所有分发任务
    allTasks := make(map[uint]map[uint]*DispatchTask) // fileID -> serverID -> task
    
    var result []ServerFileStatus
    
    for _, server := range servers {
        status := ServerFileStatus{
            ServerID:   server.ID,
            ServerName: server.Name,
            ServerIP:   server.IPAddress,
            Status:     server.Status,
        }
        
        // 获取该服务器的分发任务
        tasks, _ := h.taskRepo.GetByServerID(server.ID)
        taskMap := make(map[uint]*DispatchTask)
        for _, t := range tasks {
            taskMap[t.FileID] = t
        }
        
        // 遍历文件
        for _, file := range files {
            task, hasTask := taskMap[file.ID]
            if !hasTask {
                continue // 该文件未分发到此服务器
            }
            
            fileStatus := FileSyncStatus{
                FileID:     file.ID,
                FileName:   file.Name,
                FileMD5:    file.MD5,
                TargetPath: task.TargetPath,
                Exists:     task.Status == "success",
                Synced:     task.Status == "success",
            }
            
            // 如果服务器在线，可以检查实际文件状态
            if server.Status == "online" && task.Status == "success" {
                localMD5, err := h.getAgentFileMD5(server, task.TargetPath)
                if err == nil {
                    fileStatus.LocalMD5 = localMD5
                    fileStatus.Synced = localMD5 == file.MD5
                }
            }
            
            status.FileStatus = append(status.FileStatus, fileStatus)
        }
        
        result = append(result, status)
    }
    
    response.Success(c, result)
}

// GetServerFiles 指定服务器的文件列表
func (h *Handler) GetServerFiles(c *gin.Context) {
    id := c.Param("id")
    
    server, err := h.serverRepo.GetByID(parseUint(id))
    if err != nil {
        response.Error(c, 404, "服务器不存在")
        return
    }
    
    // 获取该服务器的分发任务
    tasks, _ := h.taskRepo.GetByServerID(server.ID)
    
    var result []FileSyncStatus
    
    for _, task := range tasks {
        file, _ := h.fileRepo.GetByID(task.FileID)
        if file == nil {
            continue
        }
        
        fileStatus := FileSyncStatus{
            FileID:     file.ID,
            FileName:   file.Name,
            FileMD5:    file.MD5,
            TargetPath: task.TargetPath,
            Exists:     task.Status == "success",
            Synced:     task.Status == "success",
        }
        
        result = append(result, fileStatus)
    }
    
    response.Success(c, result)
}

// getAgentFileMD5 获取 Agent 上文件的 MD5
func (h *Handler) getAgentFileMD5(server *model.Server, path string) (string, error) {
    url := fmt.Sprintf("http://%s:%d/api/agent/v1/fs/read?path=%s", 
        server.IPAddress, server.AgentPort, path)
    
    resp, err := h.agentClient.Get(url)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != 200 {
        return "", fmt.Errorf("agent returned %d", resp.StatusCode)
    }
    
    // 这里简化处理，实际应该调用 Agent 的 md5 接口
    // 或者在 Agent 端增加一个专门的文件校验接口
    return "", nil
}
```

### 4.4 重试任务

```go
// internal/squ-apiserver/handler/dispatch/task.go

// RetryTask 重试失败的任务
func (h *Handler) RetryTask(c *gin.Context) {
    id := c.Param("id")
    
    task, err := h.taskRepo.GetByID(parseUint(id))
    if err != nil {
        response.Error(c, 404, "任务不存在")
        return
    }
    
    if task.Status != "failed" {
        response.Error(c, 400, "只能重试失败的任务")
        return
    }
    
    // 重置状态
    task.Status = "pending"
    task.ErrorMsg = ""
    task.StartedAt = nil
    task.FinishedAt = nil
    
    if err := h.taskRepo.Update(task); err != nil {
        response.Error(c, 500, "重试失败: "+err.Error())
        return
    }
    
    // 触发分发
    go h.dispatchSvc.ProcessPendingTasks()
    
    response.Success(c, task)
}
```

## 5. 前端实现

### 5.1 目录结构

```
front/src/
└── views/
    └── DispatchManagement/
        ├── index.vue           # 主页面
        └── components/
            ├── ServerList.vue  # 服务器列表
            ├── TaskList.vue    # 任务列表
            └── FileSync.vue    # 文件同步状态
```

### 5.2 API 接口

```typescript
// front/src/api/dispatch.ts

import request from '@/utils/request'

// 服务器列表
export function listServers() {
  return request.get<ServerWithStatus[]>('/api/v1/servers')
}

// 所有服务器文件状态
export function listServerFiles() {
  return request.get<ServerFileStatus[]>('/api/v1/servers/files')
}

// 指定服务器文件列表
export function getServerFiles(serverId: number) {
  return request.get<FileSyncStatus[]>(`/api/v1/servers/${serverId}/files`)
}

// 重试任务
export function retryTask(taskId: number) {
  return request.post(`/api/v1/dispatch/tasks/${taskId}/retry`)
}
```

### 5.3 页面布局

```
┌──────────────────────────────────────────────────────────────────────────┐
│ 📡 分发管理                                            [刷新]            │
├──────────────────────────────────────────────────────────────────────────┤
│ 📋 服务器状态                                                             │
│ ┌────────────────┬───────────┬───────────┬───────────┬─────────────────┐ │
│ │ 服务器          │ 状态      │ 总文件    │ 已同步    │ 失败            │ │
│ │ server-1       │ 🟢 在线   │ 5         │ 5         │ 0               │ │
│ │ server-2       │ 🟢 在线   │ 5         │ 4         │ 1               │ │
│ │ server-3       │ 🔴 离线   │ 3         │ 3         │ 0               │ │
│ └────────────────┴───────────┴───────────┴───────────┴─────────────────┘ │
├──────────────────────────────────────────────────────────────────────────┤
│ 📦 分发任务                                                               │
│ ┌────────────────┬──────────┬──────────┬───────────┬────────────────────┐│
│ │ 文件名          │ 目标服务器│ 目标路径  │ 状态      │ 操作              ││
│ │ app.tar.gz     │ server-2 │ /opt/app │ ❌ 失败   │ [重试] [详情]     ││
│ │ config.yaml    │ server-1 │ /etc     │ ✅ 成功   │ [详情]            ││
│ │ nginx.conf     │ server-2 │ /etc     │ ⏳ 进行中  │ [详情]            ││
│ └────────────────┴──────────┴──────────┴───────────┴────────────────────┘│
├──────────────────────────────────────────────────────────────────────────┤
│ 🔄 文件同步状态                                                           │
│ ┌────────────────┬──────────────┬──────────────┬──────────┬─────────────┐│
│ │ 文件名          │ server-1     │ server-2     │ server-3 │ 备注        ││
│ │ app.tar.gz     │ ✅ md5:abc.. │ ❌ 不存在     │ ✅ md5:.. │             ││
│ │ config.yaml    │ ✅ md5:def.. │ ✅ md5:def..  │ ✅ md5..  │             ││
│ └────────────────┴──────────────┴──────────────┴──────────┴─────────────┘│
└──────────────────────────────────────────────────────────────────────────┘
```

### 5.4 主页面实现

```vue
<!-- DispatchManagement/index.vue -->

<template>
  <div class="dispatch-management">
    <!-- 工具栏 -->
    <div class="toolbar">
      <h2>分发管理</h2>
      <el-button @click="loadAll">
        <el-icon><Refresh /></el-icon>
        刷新
      </el-button>
    </div>

    <!-- 服务器状态 -->
    <el-card header="服务器状态" class="section">
      <ServerList :servers="servers" :loading="serversLoading" @select="selectServer" />
    </el-card>

    <!-- 分发任务 -->
    <el-card header="分发任务" class="section">
      <TaskList :tasks="tasks" :loading="tasksLoading" @retry="handleRetry" />
    </el-card>

    <!-- 文件同步状态 -->
    <el-card header="文件同步状态" class="section">
      <FileSync :data="fileStatus" :loading="fileStatusLoading" />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { listServers, listServerFiles, listDispatchTasks, retryTask } from '@/api/dispatch'
import ServerList from './components/ServerList.vue'
import TaskList from './components/TaskList.vue'
import FileSync from './components/FileSync.vue'

// 状态
const servers = ref([])
const serversLoading = ref(false)
const tasks = ref([])
const tasksLoading = ref(false)
const fileStatus = ref([])
const fileStatusLoading = ref(false)

// 加载服务器列表
async function loadServers() {
  serversLoading.value = true
  try {
    const res = await listServers()
    servers.value = res.data
  } finally {
    serversLoading.value = false
  }
}

// 加载任务列表
async function loadTasks() {
  tasksLoading.value = true
  try {
    const res = await listDispatchTasks()
    tasks.value = res.data
  } finally {
    tasksLoading.value = false
  }
}

// 加载文件同步状态
async function loadFileStatus() {
  fileStatusLoading.value = true
  try {
    const res = await listServerFiles()
    fileStatus.value = res.data
  } finally {
    fileStatusLoading.value = false
  }
}

// 加载全部
function loadAll() {
  loadServers()
  loadTasks()
  loadFileStatus()
}

// 重试任务
async function handleRetry(taskId: number) {
  await retryTask(taskId)
  ElMessage.success('任务已重新加入队列')
  loadTasks()
}

// 选择服务器
function selectServer(serverId: number) {
  // 可以展开显示该服务器的详细信息
}

onMounted(() => {
  loadAll()
})
</script>

<style scoped>
.dispatch-management {
  padding: 20px;
}

.section {
  margin-bottom: 20px;
}
</style>
```

## 6. 定时刷新

分发任务进行中时，可以使用定时器自动刷新：

```typescript
// 在 onMounted 中启用定时刷新
let refreshTimer: number | null = null

onMounted(() => {
  loadAll()
  
  // 每 5 秒刷新一次
  refreshTimer = window.setInterval(() => {
    loadTasks()
  }, 5000)
})

// 离开页面时清除定时器
onUnmounted(() => {
  if (refreshTimer) {
    clearInterval(refreshTimer)
  }
})
```

## 7. 路由配置

```typescript
// front/src/router/index.ts

{
  path: '/dispatch',
  name: 'DispatchManagement',
  component: () => import('@/views/DispatchManagement/index.vue'),
  meta: { title: '分发管理' }
}
```

## 8. 开发步骤

|| 步骤 | 内容 | 预计代码量 |
||------|------|-----------|
|| 1 | API 接口扩展 | ~50 行 |
|| 2 | 服务器状态 Handler | ~100 行 |
|| 3 | 文件状态 Handler | ~100 行 |
|| 4 | 前端页面 | ~300 行 |
|| **总计** | | **~550 行** |

## 9. 注意事项

1. **性能优化**：大量服务器时，聚合请求可以并行
2. **缓存策略**：文件状态可以缓存，避免频繁查询 Agent
3. **错误处理**：Agent 离线时显示离线状态，不影响其他节点
4. **实时更新**：可以使用 WebSocket 推送任务进度

## 10. 完成检查

- [ ] API 接口扩展完成
- [ ] 服务器状态查询实现
- [ ] 文件同步状态查询实现
- [ ] 前端页面实现
- [ ] 定时刷新功能
- [ ] 功能测试通过
