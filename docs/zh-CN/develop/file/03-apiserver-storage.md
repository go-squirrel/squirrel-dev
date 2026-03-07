# Phase 3: APIServer 文件存储

## 前置条件

已完成 Phase 1-2（Agent 端功能），Agent 具备完整的文件操作能力：
- 目录浏览、文件读取、文件下载
- 文件写入、上传、创建、删除、重命名
- 接收推送文件

## 开发目标

实现 APIServer 文件存储服务：
- 文件上传（支持断点续传）
- 文件库管理
- 分发任务管理
- 推送文件到 Agent

## 1. 目录结构

```
internal/squ-apiserver/
├── handler/
│   └── files/
│       ├── handler.go       # Handler 入口
│       ├── upload.go        # 文件上传
│       ├── chunk.go         # 分片上传
│       ├── file.go          # 文件库管理
│       └── dispatch.go      # 分发任务
├── model/
│   └── files.go             # 数据模型
├── repository/
│   └── files.go             # 数据访问
└── service/
    └── dispatcher.go        # 分发服务

front/src/
├── api/
│   └── files.ts             # API 接口
└── types/
    └── files.ts             # 类型定义
```

## 2. 数据模型

```go
// internal/squ-apiserver/model/files.go

package model

import "time"

// FileRecord 文件记录（APIServer 文件库）
type FileRecord struct {
    ID         uint      `gorm:"primaryKey" json:"id"`
    UUID       string    `gorm:"uniqueIndex;size:36" json:"uuid"`       // 文件唯一标识
    Name       string    `gorm:"size:255;not null" json:"name"`         // 文件名
    Size       int64     `gorm:"not null" json:"size"`                  // 文件大小
    Path       string    `gorm:"size:512" json:"path"`                  // 存储路径（相对路径）
    MD5        string    `gorm:"size:32" json:"md5"`                    // 文件哈希
    Mime       string    `gorm:"size:128" json:"mime"`                  // MIME 类型
    Status     string    `gorm:"size:20;default:uploaded" json:"status"` // uploaded/distributing/distributed
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
}

// TableName 指定表名
func (FileRecord) TableName() string {
    return "file_records"
}

// ChunkUpload 分片上传记录
type ChunkUpload struct {
    ID         uint   `gorm:"primaryKey" json:"id"`
    UploadID   string `gorm:"uniqueIndex;size:36" json:"upload_id"`  // 上传会话 ID
    FileName   string `gorm:"size:255;not null" json:"file_name"`
    FileSize   int64  `gorm:"not null" json:"file_size"`
    FileMD5    string `gorm:"size:32" json:"file_md5"`               // 整体文件 MD5
    ChunkSize  int64  `gorm:"not null" json:"chunk_size"`            // 分片大小
    ChunkCount int    `gorm:"not null" json:"chunk_count"`           // 分片总数
    Chunks     string `gorm:"type:text" json:"chunks"`               // 已上传分片 JSON
    Status     string `gorm:"size:20;default:uploading" json:"status"` // uploading/completed/expired
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
    ExpiresAt  *time.Time `json:"expires_at"`
}

// TableName 指定表名
func (ChunkUpload) TableName() string {
    return "chunk_uploads"
}

// DispatchTask 分发任务
type DispatchTask struct {
    ID          uint       `gorm:"primaryKey" json:"id"`
    FileID      uint       `gorm:"index;not null" json:"file_id"`
    ServerID    uint       `gorm:"index;not null" json:"server_id"`
    TargetPath  string     `gorm:"size:512;not null" json:"target_path"` // 目标路径
    Status      string     `gorm:"size:20;default:pending" json:"status"` // pending/running/success/failed
    Progress    int        `gorm:"default:0" json:"progress"`            // 进度 0-100
    ErrorMsg    string     `gorm:"type:text" json:"error_msg"`
    CreatedAt   time.Time  `json:"created_at"`
    StartedAt   *time.Time `json:"started_at"`
    FinishedAt  *time.Time `json:"finished_at"`
}

// TableName 指定表名
func (DispatchTask) TableName() string {
    return "dispatch_tasks"
}
```

## 3. API 设计

```
# 文件上传
POST   /api/v1/files/upload           # 普通上传（小文件）
POST   /api/v1/files/chunk/init       # 初始化分片上传
POST   /api/v1/files/chunk/upload     # 上传分片
POST   /api/v1/files/chunk/complete   # 合并分片
GET    /api/v1/files/chunk/status     # 查询分片上传状态

# 文件库管理
GET    /api/v1/files                  # 文件列表
GET    /api/v1/files/:id              # 文件详情
DELETE /api/v1/files/:id              # 删除文件
GET    /api/v1/files/:id/download     # 下载文件

# 分发任务
POST   /api/v1/dispatch               # 创建分发任务
GET    /api/v1/dispatch/tasks         # 任务列表
GET    /api/v1/dispatch/tasks/:id     # 任务详情
```

## 4. 文件上传实现

### 4.1 Handler 入口

```go
// internal/squ-apiserver/handler/files/handler.go

package files

import (
    "github.com/gin-gonic/gin"
    
    "squirrel-dev/internal/squ-apiserver/repository"
    "squirrel-dev/internal/squ-apiserver/service"
)

// Handler 文件管理 Handler
type Handler struct {
    fileRepo   repository.FileRepository
    chunkRepo  repository.ChunkRepository
    dispatchSvc *service.DispatchService
    storagePath string // 文件存储路径
}

// NewHandler 创建 Handler
func NewHandler(
    fileRepo repository.FileRepository,
    chunkRepo repository.ChunkRepository,
    dispatchSvc *service.DispatchService,
    storagePath string,
) *Handler {
    return &Handler{
        fileRepo:    fileRepo,
        chunkRepo:   chunkRepo,
        dispatchSvc: dispatchSvc,
        storagePath: storagePath,
    }
}

// RegisterRoutes 注册路由
func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
    files := r.Group("/files")
    {
        // 上传
        files.POST("/upload", h.Upload)
        files.POST("/chunk/init", h.ChunkInit)
        files.POST("/chunk/upload", h.ChunkUpload)
        files.POST("/chunk/complete", h.ChunkComplete)
        files.GET("/chunk/status", h.ChunkStatus)
        
        // 文件库
        files.GET("", h.List)
        files.GET("/:id", h.Get)
        files.DELETE("/:id", h.Delete)
        files.GET("/:id/download", h.Download)
    }
    
    // 分发任务
    dispatch := r.Group("/dispatch")
    {
        dispatch.POST("", h.CreateDispatch)
        dispatch.GET("/tasks", h.ListTasks)
        dispatch.GET("/tasks/:id", h.GetTask)
    }
}
```

### 4.2 普通上传

```go
// internal/squ-apiserver/handler/files/upload.go

package files

import (
    "crypto/md5"
    "encoding/hex"
    "io"
    "os"
    "path/filepath"
    "time"
    
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    
    "squirrel-dev/internal/squ-apiserver/model"
    "squirrel-dev/internal/pkg/response"
)

// Upload 普通文件上传（适合小文件 < 50MB）
func (h *Handler) Upload(c *gin.Context) {
    // 获取上传的文件
    file, header, err := c.Request.FormFile("file")
    if err != nil {
        response.Error(c, 400, "获取上传文件失败: "+err.Error())
        return
    }
    defer file.Close()
    
    // 生成 UUID
    fileUUID := uuid.New().String()
    
    // 创建存储目录
    storageDir := filepath.Join(h.storagePath, time.Now().Format("2006/01/02"))
    if err := os.MkdirAll(storageDir, 0755); err != nil {
        response.Error(c, 500, "创建存储目录失败: "+err.Error())
        return
    }
    
    // 目标文件路径
    targetPath := filepath.Join(storageDir, fileUUID+"_"+header.Filename)
    
    // 创建目标文件
    dst, err := os.Create(targetPath)
    if err != nil {
        response.Error(c, 500, "创建文件失败: "+err.Error())
        return
    }
    defer dst.Close()
    
    // 计算 MD5 并写入文件
    hash := md5.New()
    writer := io.MultiWriter(dst, hash)
    
    written, err := io.Copy(writer, file)
    if err != nil {
        response.Error(c, 500, "写入文件失败: "+err.Error())
        return
    }
    
    // 保存文件记录
    record := &model.FileRecord{
        UUID:  fileUUID,
        Name:  header.Filename,
        Size:  written,
        Path:  filepath.Join(time.Now().Format("2006/01/02"), fileUUID+"_"+header.Filename),
        MD5:   hex.EncodeToString(hash.Sum(nil)),
        Mime:  header.Header.Get("Content-Type"),
    }
    
    if err := h.fileRepo.Create(record); err != nil {
        os.Remove(targetPath)
        response.Error(c, 500, "保存文件记录失败: "+err.Error())
        return
    }
    
    response.Success(c, record)
}
```

## 5. 分片上传实现

### 5.1 初始化分片上传

```go
// internal/squ-apiserver/handler/files/chunk.go

package files

import (
    "time"
    
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    
    "squirrel-dev/internal/squ-apiserver/model"
    "squirrel-dev/internal/pkg/response"
)

// ChunkInitRequest 初始化请求
type ChunkInitRequest struct {
    FileName string `json:"file_name" binding:"required"`
    FileSize int64  `json:"file_size" binding:"required"`
    FileMD5  string `json:"file_md5"`  // 可选，用于校验
}

// ChunkInitResponse 初始化响应
type ChunkInitResponse struct {
    UploadID   string `json:"upload_id"`
    ChunkSize  int64  `json:"chunk_size"`  // 分片大小（字节）
    ChunkCount int    `json:"chunk_count"` // 分片总数
}

// ChunkInit 初始化分片上传
func (h *Handler) ChunkInit(c *gin.Context) {
    var req ChunkInitRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, 400, "参数错误: "+err.Error())
        return
    }
    
    // 计算分片大小（默认 5MB）
    chunkSize := int64(5 * 1024 * 1024)
    
    // 大文件自动增大分片
    if req.FileSize > 100*1024*1024 { // > 100MB
        chunkSize = 10 * 1024 * 1024 // 10MB
    }
    if req.FileSize > 1024*1024*1024 { // > 1GB
        chunkSize = 20 * 1024 * 1024 // 20MB
    }
    
    // 计算分片数量
    chunkCount := int(req.FileSize / chunkSize)
    if req.FileSize%chunkSize != 0 {
        chunkCount++
    }
    
    // 创建上传记录
    upload := &model.ChunkUpload{
        UploadID:   uuid.New().String(),
        FileName:   req.FileName,
        FileSize:   req.FileSize,
        FileMD5:    req.FileMD5,
        ChunkSize:  chunkSize,
        ChunkCount: chunkCount,
        Chunks:     "[]", // 空数组
        Status:     "uploading",
    }
    
    // 设置过期时间（24小时后）
    expiresAt := time.Now().Add(24 * time.Hour)
    upload.ExpiresAt = &expiresAt
    
    if err := h.chunkRepo.Create(upload); err != nil {
        response.Error(c, 500, "创建上传记录失败: "+err.Error())
        return
    }
    
    response.Success(c, ChunkInitResponse{
        UploadID:   upload.UploadID,
        ChunkSize:  chunkSize,
        ChunkCount: chunkCount,
    })
}
```

### 5.2 上传分片

```go
// ChunkUploadRequest 分片上传请求
type ChunkUploadRequest struct {
    UploadID   string `form:"upload_id" binding:"required"`
    ChunkIndex int    `form:"chunk_index" binding:"required"`
}

// ChunkUpload 上传分片
func (h *Handler) ChunkUpload(c *gin.Context) {
    var req ChunkUploadRequest
    if err := c.ShouldBindQuery(&req); err != nil {
        response.Error(c, 400, "参数错误: "+err.Error())
        return
    }
    
    // 获取上传记录
    upload, err := h.chunkRepo.GetByUploadID(req.UploadID)
    if err != nil {
        response.Error(c, 404, "上传记录不存在")
        return
    }
    
    // 检查是否过期
    if upload.ExpiresAt != nil && time.Now().After(*upload.ExpiresAt) {
        response.Error(c, 400, "上传已过期，请重新初始化")
        return
    }
    
    // 检查分片索引
    if req.ChunkIndex < 0 || req.ChunkIndex >= upload.ChunkCount {
        response.Error(c, 400, "分片索引超出范围")
        return
    }
    
    // 创建分片存储目录
    chunkDir := filepath.Join(h.storagePath, "chunks", req.UploadID)
    if err := os.MkdirAll(chunkDir, 0755); err != nil {
        response.Error(c, 500, "创建分片目录失败: "+err.Error())
        return
    }
    
    // 保存分片
    chunkPath := filepath.Join(chunkDir, fmt.Sprintf("%d", req.ChunkIndex))
    dst, err := os.Create(chunkPath)
    if err != nil {
        response.Error(c, 500, "创建分片文件失败: "+err.Error())
        return
    }
    defer dst.Close()
    
    written, err := io.Copy(dst, c.Request.Body)
    if err != nil {
        response.Error(c, 500, "写入分片失败: "+err.Error())
        return
    }
    
    // 更新已上传分片列表
    if err := h.chunkRepo.MarkChunkUploaded(req.UploadID, req.ChunkIndex); err != nil {
        response.Error(c, 500, "更新分片状态失败: "+err.Error())
        return
    }
    
    response.Success(c, gin.H{
        "upload_id":   req.UploadID,
        "chunk_index": req.ChunkIndex,
        "size":        written,
    })
}
```

### 5.3 合并分片

```go
// ChunkCompleteRequest 合并请求
type ChunkCompleteRequest struct {
    UploadID string `json:"upload_id" binding:"required"`
}

// ChunkComplete 合并分片
func (h *Handler) ChunkComplete(c *gin.Context) {
    var req ChunkCompleteRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, 400, "参数错误: "+err.Error())
        return
    }
    
    // 获取上传记录
    upload, err := h.chunkRepo.GetByUploadID(req.UploadID)
    if err != nil {
        response.Error(c, 404, "上传记录不存在")
        return
    }
    
    // 检查是否所有分片都已上传
    var uploadedChunks []bool
    if err := json.Unmarshal([]byte(upload.Chunks), &uploadedChunks); err != nil {
        response.Error(c, 500, "解析分片状态失败")
        return
    }
    
    for i, uploaded := range uploadedChunks {
        if !uploaded {
            response.Error(c, 400, fmt.Sprintf("分片 %d 未上传", i))
            return
        }
    }
    
    // 合并分片
    fileUUID := uuid.New().String()
    storageDir := filepath.Join(h.storagePath, time.Now().Format("2006/01/02"))
    if err := os.MkdirAll(storageDir, 0755); err != nil {
        response.Error(c, 500, "创建存储目录失败")
        return
    }
    
    targetPath := filepath.Join(storageDir, fileUUID+"_"+upload.FileName)
    dst, err := os.Create(targetPath)
    if err != nil {
        response.Error(c, 500, "创建目标文件失败")
        return
    }
    defer dst.Close()
    
    // 计算 MD5
    hash := md5.New()
    writer := io.MultiWriter(dst, hash)
    
    // 按顺序合并分片
    chunkDir := filepath.Join(h.storagePath, "chunks", req.UploadID)
    for i := 0; i < upload.ChunkCount; i++ {
        chunkPath := filepath.Join(chunkDir, fmt.Sprintf("%d", i))
        chunkFile, err := os.Open(chunkPath)
        if err != nil {
            response.Error(c, 500, fmt.Sprintf("打开分片 %d 失败", i))
            return
        }
        io.Copy(writer, chunkFile)
        chunkFile.Close()
    }
    
    // 清理分片
    os.RemoveAll(chunkDir)
    
    // 更新上传记录状态
    h.chunkRepo.UpdateStatus(req.UploadID, "completed")
    
    // 创建文件记录
    record := &model.FileRecord{
        UUID:  fileUUID,
        Name:  upload.FileName,
        Size:  upload.FileSize,
        Path:  filepath.Join(time.Now().Format("2006/01/02"), fileUUID+"_"+upload.FileName),
        MD5:   hex.EncodeToString(hash.Sum(nil)),
    }
    
    if err := h.fileRepo.Create(record); err != nil {
        os.Remove(targetPath)
        response.Error(c, 500, "保存文件记录失败")
        return
    }
    
    response.Success(c, record)
}
```

## 6. 文件库管理

```go
// internal/squ-apiserver/handler/files/file.go

package files

import (
    "fmt"
    "os"
    "path/filepath"
    
    "github.com/gin-gonic/gin"
    
    "squirrel-dev/internal/pkg/response"
)

// ListRequest 列表请求
type ListRequest struct {
    Page     int    `form:"page" binding:"required,min=1"`
    PageSize int    `form:"page_size" binding:"required,min=1,max=100"`
    Name     string `form:"name"` // 文件名搜索
}

// List 文件列表
func (h *Handler) List(c *gin.Context) {
    var req ListRequest
    if err := c.ShouldBindQuery(&req); err != nil {
        response.Error(c, 400, "参数错误: "+err.Error())
        return
    }
    
    records, total, err := h.fileRepo.List(req.Page, req.PageSize, req.Name)
    if err != nil {
        response.Error(c, 500, "查询失败: "+err.Error())
        return
    }
    
    response.Success(c, gin.H{
        "list":  records,
        "total": total,
    })
}

// Get 文件详情
func (h *Handler) Get(c *gin.Context) {
    id := c.Param("id")
    
    record, err := h.fileRepo.GetByID(parseUint(id))
    if err != nil {
        response.Error(c, 404, "文件不存在")
        return
    }
    
    response.Success(c, record)
}

// Delete 删除文件
func (h *Handler) Delete(c *gin.Context) {
    id := c.Param("id")
    
    record, err := h.fileRepo.GetByID(parseUint(id))
    if err != nil {
        response.Error(c, 404, "文件不存在")
        return
    }
    
    // 删除物理文件
    filePath := filepath.Join(h.storagePath, record.Path)
    os.Remove(filePath)
    
    // 删除记录
    if err := h.fileRepo.Delete(parseUint(id)); err != nil {
        response.Error(c, 500, "删除失败: "+err.Error())
        return
    }
    
    response.Success(c, gin.H{"id": id})
}

// Download 下载文件
func (h *Handler) Download(c *gin.Context) {
    id := c.Param("id")
    
    record, err := h.fileRepo.GetByID(parseUint(id))
    if err != nil {
        response.Error(c, 404, "文件不存在")
        return
    }
    
    filePath := filepath.Join(h.storagePath, record.Path)
    
    c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", record.Name))
    c.Header("Content-Type", "application/octet-stream")
    c.File(filePath)
}
```

## 7. 分发任务实现

### 7.1 创建分发任务

```go
// internal/squ-apiserver/handler/files/dispatch.go

package files

import (
    "github.com/gin-gonic/gin"
    
    "squirrel-dev/internal/squ-apiserver/model"
    "squirrel-dev/internal/pkg/response"
)

// CreateDispatchRequest 创建分发请求
type CreateDispatchRequest struct {
    FileID      uint   `json:"file_id" binding:"required"`
    ServerIDs   []uint `json:"server_ids" binding:"required"`  // 目标服务器列表
    TargetPath  string `json:"target_path" binding:"required"` // 目标路径
}

// CreateDispatch 创建分发任务
func (h *Handler) CreateDispatch(c *gin.Context) {
    var req CreateDispatchRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, 400, "参数错误: "+err.Error())
        return
    }
    
    // 获取文件记录
    file, err := h.fileRepo.GetByID(req.FileID)
    if err != nil {
        response.Error(c, 404, "文件不存在")
        return
    }
    
    // 创建分发任务
    var tasks []*model.DispatchTask
    for _, serverID := range req.ServerIDs {
        task := &model.DispatchTask{
            FileID:     req.FileID,
            ServerID:   serverID,
            TargetPath: req.TargetPath,
            Status:     "pending",
        }
        tasks = append(tasks, task)
    }
    
    if err := h.dispatchSvc.CreateTasks(tasks); err != nil {
        response.Error(c, 500, "创建任务失败: "+err.Error())
        return
    }
    
    // 更新文件状态
    h.fileRepo.UpdateStatus(req.FileID, "distributing")
    
    // 启动异步分发
    go h.dispatchSvc.ProcessPendingTasks()
    
    response.Success(c, gin.H{
        "task_count": len(tasks),
        "file":       file,
    })
}

// ListTasks 任务列表
func (h *Handler) ListTasks(c *gin.Context) {
    fileID := c.Query("file_id")
    status := c.Query("status")
    
    tasks, err := h.dispatchSvc.ListTasks(parseUint(fileID), status)
    if err != nil {
        response.Error(c, 500, "查询失败: "+err.Error())
        return
    }
    
    response.Success(c, tasks)
}

// GetTask 任务详情
func (h *Handler) GetTask(c *gin.Context) {
    id := c.Param("id")
    
    task, err := h.dispatchSvc.GetTask(parseUint(id))
    if err != nil {
        response.Error(c, 404, "任务不存在")
        return
    }
    
    response.Success(c, task)
}
```

### 7.2 分发服务

```go
// internal/squ-apiserver/service/dispatcher.go

package service

import (
    "fmt"
    "io"
    "net/http"
    "os"
    "path/filepath"
    "sync"
    "time"
    
    "squirrel-dev/internal/squ-apiserver/model"
    "squirrel-dev/internal/squ-apiserver/repository"
)

// DispatchService 分发服务
type DispatchService struct {
    taskRepo    repository.DispatchRepository
    fileRepo    repository.FileRepository
    serverRepo  repository.ServerRepository
    storagePath string
    client      *http.Client
    mu          sync.Mutex
}

// NewDispatchService 创建分发服务
func NewDispatchService(
    taskRepo repository.DispatchRepository,
    fileRepo repository.FileRepository,
    serverRepo repository.ServerRepository,
    storagePath string,
) *DispatchService {
    return &DispatchService{
        taskRepo:    taskRepo,
        fileRepo:    fileRepo,
        serverRepo:  serverRepo,
        storagePath: storagePath,
        client: &http.Client{
            Timeout: 30 * time.Minute, // 大文件传输可能需要较长时间
        },
    }
}

// ProcessPendingTasks 处理待执行任务
func (s *DispatchService) ProcessPendingTasks() {
    tasks, err := s.taskRepo.GetPendingTasks()
    if err != nil {
        return
    }
    
    var wg sync.WaitGroup
    for _, task := range tasks {
        wg.Add(1)
        go func(t *model.DispatchTask) {
            defer wg.Done()
            s.processTask(t)
        }(task)
    }
    wg.Wait()
}

// processTask 处理单个任务
func (s *DispatchService) processTask(task *model.DispatchTask) {
    s.mu.Lock()
    // 更新状态为运行中
    now := time.Now()
    task.Status = "running"
    task.StartedAt = &now
    s.taskRepo.Update(task)
    s.mu.Unlock()
    
    // 获取文件信息
    file, err := s.fileRepo.GetByID(task.FileID)
    if err != nil {
        s.failTask(task, "文件不存在")
        return
    }
    
    // 获取服务器信息
    server, err := s.serverRepo.GetByID(task.ServerID)
    if err != nil {
        s.failTask(task, "服务器不存在")
        return
    }
    
    if server.Status != model.ServerStatusOnline {
        s.failTask(task, "服务器离线")
        return
    }
    
    // 打开文件
    filePath := filepath.Join(s.storagePath, file.Path)
    f, err := os.Open(filePath)
    if err != nil {
        s.failTask(task, "打开文件失败: "+err.Error())
        return
    }
    defer f.Close()
    
    // 构建请求
    url := fmt.Sprintf("http://%s:%d/api/agent/v1/fs/receive?target_path=%s&overwrite=true",
        server.IpAddress, server.AgentPort, task.TargetPath)
    
    req, err := http.NewRequest("POST", url, f)
    if err != nil {
        s.failTask(task, "创建请求失败: "+err.Error())
        return
    }
    req.Header.Set("Content-Type", "application/octet-stream")
    
    // 发送请求
    resp, err := s.client.Do(req)
    if err != nil {
        s.failTask(task, "发送请求失败: "+err.Error())
        return
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != 200 {
        body, _ := io.ReadAll(resp.Body)
        s.failTask(task, fmt.Sprintf("Agent 返回错误: %s", string(body)))
        return
    }
    
    // 成功
    finishedAt := time.Now()
    task.Status = "success"
    task.Progress = 100
    task.FinishedAt = &finishedAt
    s.taskRepo.Update(task)
    
    // 更新文件状态
    s.updateFileStatus(task.FileID)
}

// failTask 标记任务失败
func (s *DispatchService) failTask(task *model.DispatchTask, errMsg string) {
    finishedAt := time.Now()
    task.Status = "failed"
    task.ErrorMsg = errMsg
    task.FinishedAt = &finishedAt
    s.taskRepo.Update(task)
}

// updateFileStatus 更新文件状态
func (s *DispatchService) updateFileStatus(fileID uint) {
    tasks, _ := s.taskRepo.GetByFileID(fileID)
    
    allSuccess := true
    anyRunning := false
    
    for _, t := range tasks {
        if t.Status == "running" || t.Status == "pending" {
            anyRunning = true
        }
        if t.Status != "success" {
            allSuccess = false
        }
    }
    
    if anyRunning {
        s.fileRepo.UpdateStatus(fileID, "distributing")
    } else if allSuccess {
        s.fileRepo.UpdateStatus(fileID, "distributed")
    }
}
```

## 8. 路由注册

在 `internal/squ-apiserver/router/router.go` 中添加：

```go
import (
    fileshandler "squirrel-dev/internal/squ-apiserver/handler/files"
    "squirrel-dev/internal/squ-apiserver/repository"
    "squirrel-dev/internal/squ-apiserver/service"
)

func InitRouter(r *gin.Engine, conf *config.Config, db database.DB) {
    // ... 现有代码 ...
    
    // 文件存储配置
    storagePath := conf.Storage.Path // 从配置读取
    
    // 初始化依赖
    fileRepo := repository.NewFileRepository(db)
    chunkRepo := repository.NewChunkRepository(db)
    taskRepo := repository.NewDispatchRepository(db)
    serverRepo := repository.NewServerRepository(db)
    dispatchSvc := service.NewDispatchService(taskRepo, fileRepo, serverRepo, storagePath)
    
    // 创建 Handler
    filesHandler := fileshandler.NewHandler(fileRepo, chunkRepo, dispatchSvc, storagePath)
    
    // 注册路由
    api := r.Group("/api/v1")
    api.Use(auth.Middleware()) // 认证中间件
    {
        fileshandler.RegisterRoutes(api, filesHandler)
    }
}
```

## 9. 配置设计

```yaml
# config/apiserver.yaml
storage:
  path: /data/squirrel/files       # 文件存储路径
  max_size: 10GB                   # 最大存储空间
  chunk_size: 5MB                  # 默认分片大小
  chunk_expire: 24h                # 分片上传过期时间

dispatch:
  concurrent: 5                    # 并发分发数
  retry: 3                         # 失败重试次数
  timeout: 30m                     # 传输超时时间
```

## 10. 数据库迁移

```sql
-- 文件记录表
CREATE TABLE file_records (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    uuid VARCHAR(36) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    size BIGINT NOT NULL,
    path VARCHAR(512),
    md5 VARCHAR(32),
    mime VARCHAR(128),
    status VARCHAR(20) DEFAULT 'uploaded',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_status (status),
    INDEX idx_created (created_at)
);

-- 分片上传表
CREATE TABLE chunk_uploads (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    upload_id VARCHAR(36) UNIQUE NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_size BIGINT NOT NULL,
    file_md5 VARCHAR(32),
    chunk_size BIGINT NOT NULL,
    chunk_count INT NOT NULL,
    chunks TEXT,
    status VARCHAR(20) DEFAULT 'uploading',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    expires_at TIMESTAMP,
    INDEX idx_upload_id (upload_id),
    INDEX idx_status (status)
);

-- 分发任务表
CREATE TABLE dispatch_tasks (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    file_id BIGINT UNSIGNED NOT NULL,
    server_id BIGINT UNSIGNED NOT NULL,
    target_path VARCHAR(512) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    progress INT DEFAULT 0,
    error_msg TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    started_at TIMESTAMP NULL,
    finished_at TIMESTAMP NULL,
    INDEX idx_file_id (file_id),
    INDEX idx_server_id (server_id),
    INDEX idx_status (status),
    FOREIGN KEY (file_id) REFERENCES file_records(id)
);
```

## 11. 测试用例

### 测试文件上传

```bash
# 普通上传
curl -X POST "http://localhost:10700/api/v1/files/upload" \
  -F "file=@./test-file.txt"
```

### 测试分片上传

```bash
# 1. 初始化
curl -X POST "http://localhost:10700/api/v1/files/chunk/init" \
  -H "Content-Type: application/json" \
  -d '{"file_name": "large-file.tar.gz", "file_size": 104857600}'

# 返回: {"upload_id": "xxx", "chunk_size": 5242880, "chunk_count": 20}

# 2. 上传分片
curl -X POST "http://localhost:10700/api/v1/files/chunk/upload?upload_id=xxx&chunk_index=0" \
  --data-binary "@./chunk_0"

# 3. 合并分片
curl -X POST "http://localhost:10700/api/v1/files/chunk/complete" \
  -H "Content-Type: application/json" \
  -d '{"upload_id": "xxx"}'
```

### 测试分发

```bash
# 创建分发任务
curl -X POST "http://localhost:10700/api/v1/dispatch" \
  -H "Content-Type: application/json" \
  -d '{"file_id": 1, "server_ids": [1, 2, 3], "target_path": "/opt/app/app.tar.gz"}'

# 查看任务列表
curl "http://localhost:10700/api/v1/dispatch/tasks"
```

## 12. 完成检查

- [ ] 数据模型定义完成
- [ ] 数据库迁移执行
- [ ] 普通文件上传实现
- [ ] 分片上传实现
- [ ] 文件库管理实现
- [ ] 分发任务创建实现
- [ ] 分发服务实现
- [ ] 路由注册完成
- [ ] 单元测试通过
- [ ] 集成测试通过

## 13. 下一阶段

完成本阶段后，APIServer 文件存储功能已完整。继续开发 [04-frontend-filemanager.md](./04-frontend-filemanager.md) 实现前端页面。
