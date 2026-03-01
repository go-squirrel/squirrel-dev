# Phase 3: APIServer 代理转发

## 前置条件

已完成 Phase 1-2（Agent 端功能），Agent 具备完整的文件操作能力：
- 目录浏览、文件读取、文件下载
- 文件写入、上传、创建、删除、重命名

## 开发目标

实现 APIServer 代理转发：
- 请求路由到对应 Agent
- 流式代理（上传/下载）
- 统一响应格式

## 1. 目录结构

```
internal/squ-apiserver/
└── handler/
    └── fs/
        ├── handler.go       # Handler 入口
        ├── proxy.go         # 请求代理
        └── stream.go        # 流式传输
```

## 2. API 设计

APIServer 端所有接口通过 `serverId` 路由到对应 Agent：

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/fs/{serverId}/list` | 目录列表 |
| GET | `/api/v1/fs/{serverId}/read` | 文件读取 |
| GET | `/api/v1/fs/{serverId}/download` | 文件下载 |
| POST | `/api/v1/fs/{serverId}/write` | 文件写入 |
| POST | `/api/v1/fs/{serverId}/upload` | 文件上传 |
| POST | `/api/v1/fs/{serverId}/mkdir` | 创建目录 |
| POST | `/api/v1/fs/{serverId}/create` | 创建文件 |
| POST | `/api/v1/fs/{serverId}/delete` | 删除 |
| POST | `/api/v1/fs/{serverId}/rename` | 重命名 |

## 3. Handler 入口

```go
// internal/squ-apiserver/handler/fs/handler.go

package fs

import (
    "github.com/gin-gonic/gin"
    
    "squirrel-dev/internal/squ-apiserver/model"
    "squirrel-dev/internal/squ-apiserver/repository"
)

// Handler 文件管理 Handler
type Handler struct {
    serverRepo  repository.ServerRepository
    agentClient *AgentClient
}

// NewHandler 创建 Handler
func NewHandler(serverRepo repository.ServerRepository) *Handler {
    return &Handler{
        serverRepo:  serverRepo,
        agentClient: NewAgentClient(),
    }
}

// RegisterRoutes 注册路由
func RegisterRoutes(r *gin.RouterGroup, serverRepo repository.ServerRepository) {
    h := NewHandler(serverRepo)
    
    fs := r.Group("/fs/:serverId")
    {
        // 读取操作
        fs.GET("/list", h.List)
        fs.GET("/read", h.Read)
        fs.GET("/download", h.Download)
        
        // 写入操作
        fs.POST("/write", h.Write)
        fs.POST("/upload", h.Upload)
        fs.POST("/mkdir", h.Mkdir)
        fs.POST("/create", h.CreateFile)
        fs.POST("/delete", h.Delete)
        fs.POST("/rename", h.Rename)
    }
}
```

## 4. Agent 客户端

```go
// internal/squ-apiserver/handler/fs/proxy.go

package fs

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
    
    "squirrel-dev/internal/squ-apiserver/model"
    "squirrel-dev/internal/squ-apiserver/repository"
)

// AgentClient Agent HTTP 客户端
type AgentClient struct {
    httpClient *http.Client
}

// NewAgentClient 创建客户端
func NewAgentClient() *AgentClient {
    return &AgentClient{
        httpClient: &http.Client{
            Timeout: 30 * time.Second,
            // 可以配置 TLS 用于 mTLS
        },
    }
}

// getAgentURL 获取 Agent 地址
func (h *Handler) getAgentURL(serverId string) (string, error) {
    server, err := h.serverRepo.GetByID(parseUint(serverId))
    if err != nil {
        return "", fmt.Errorf("服务器不存在")
    }
    
    if server.Status != model.ServerStatusOnline {
        return "", fmt.Errorf("服务器离线")
    }
    
    return fmt.Sprintf("http://%s:%d/api/agent/v1", server.IpAddress, server.AgentPort), nil
}

// proxyGet 代理 GET 请求
func (h *Handler) proxyGet(c *gin.Context, endpoint string) {
    serverId := c.Param("serverId")
    
    // 获取 Agent 地址
    baseURL, err := h.getAgentURL(serverId)
    if err != nil {
        c.JSON(400, gin.H{"code": 400, "msg": err.Error()})
        return
    }
    
    // 构建请求 URL
    targetURL := fmt.Sprintf("%s/fs/%s?%s", baseURL, endpoint, c.Request.URL.RawQuery)
    
    // 发送请求
    resp, err := h.agentClient.httpClient.Get(targetURL)
    if err != nil {
        c.JSON(500, gin.H{"code": 500, "msg": "请求 Agent 失败: " + err.Error()})
        return
    }
    defer resp.Body.Close()
    
    // 转发响应
    c.Status(resp.StatusCode)
    for k, v := range resp.Header {
        c.Header(k, v[0])
    }
    io.Copy(c.Writer, resp.Body)
}

// proxyPost 代理 POST 请求
func (h *Handler) proxyPost(c *gin.Context, endpoint string) {
    serverId := c.Param("serverId")
    
    // 获取 Agent 地址
    baseURL, err := h.getAgentURL(serverId)
    if err != nil {
        c.JSON(400, gin.H{"code": 400, "msg": err.Error()})
        return
    }
    
    // 构建请求 URL
    targetURL := fmt.Sprintf("%s/fs/%s", baseURL, endpoint)
    
    // 读取请求体
    body, err := io.ReadAll(c.Request.Body)
    if err != nil {
        c.JSON(500, gin.H{"code": 500, "msg": "读取请求体失败"})
        return
    }
    
    // 创建请求
    req, err := http.NewRequest("POST", targetURL, bytes.NewReader(body))
    if err != nil {
        c.JSON(500, gin.H{"code": 500, "msg": "创建请求失败"})
        return
    }
    req.Header.Set("Content-Type", c.ContentType())
    
    // 发送请求
    resp, err := h.agentClient.httpClient.Do(req)
    if err != nil {
        c.JSON(500, gin.H{"code": 500, "msg": "请求 Agent 失败: " + err.Error()})
        return
    }
    defer resp.Body.Close()
    
    // 转发响应
    c.Status(resp.StatusCode)
    for k, v := range resp.Header {
        c.Header(k, v[0])
    }
    io.Copy(c.Writer, resp.Body)
}
```

## 5. 各接口实现

```go
// internal/squ-apiserver/handler/fs/handler.go 续

// List 目录列表
func (h *Handler) List(c *gin.Context) {
    h.proxyGet(c, "list")
}

// Read 文件读取
func (h *Handler) Read(c *gin.Context) {
    h.proxyGet(c, "read")
}

// Download 文件下载
func (h *Handler) Download(c *gin.Context) {
    serverId := c.Param("serverId")
    
    baseURL, err := h.getAgentURL(serverId)
    if err != nil {
        c.JSON(400, gin.H{"code": 400, "msg": err.Error()})
        return
    }
    
    targetURL := fmt.Sprintf("%s/fs/download?%s", baseURL, c.Request.URL.RawQuery)
    
    // 流式下载
    resp, err := h.agentClient.httpClient.Get(targetURL)
    if err != nil {
        c.JSON(500, gin.H{"code": 500, "msg": "下载失败: " + err.Error()})
        return
    }
    defer resp.Body.Close()
    
    // 设置响应头
    c.Header("Content-Disposition", resp.Header.Get("Content-Disposition"))
    c.Header("Content-Length", resp.Header.Get("Content-Length"))
    c.Header("Content-Type", resp.Header.Get("Content-Type"))
    
    // 流式传输
    io.Copy(c.Writer, resp.Body)
}

// Write 文件写入
func (h *Handler) Write(c *gin.Context) {
    h.proxyPost(c, "write")
}

// Upload 文件上传
func (h *Handler) Upload(c *gin.Context) {
    serverId := c.Param("serverId")
    
    baseURL, err := h.getAgentURL(serverId)
    if err != nil {
        c.JSON(400, gin.H{"code": 400, "msg": err.Error()})
        return
    }
    
    targetURL := fmt.Sprintf("%s/fs/upload?%s", baseURL, c.Request.URL.RawQuery)
    
    // 创建 multipart 请求
    c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 50*1024*1024) // 50MB 限制
    
    // 转发请求
    resp, err := http.Post(targetURL, c.ContentType(), c.Request.Body)
    if err != nil {
        c.JSON(500, gin.H{"code": 500, "msg": "上传失败: " + err.Error()})
        return
    }
    defer resp.Body.Close()
    
    c.Status(resp.StatusCode)
    io.Copy(c.Writer, resp.Body)
}

// Mkdir 创建目录
func (h *Handler) Mkdir(c *gin.Context) {
    h.proxyPost(c, "mkdir")
}

// CreateFile 创建文件
func (h *Handler) CreateFile(c *gin.Context) {
    h.proxyPost(c, "create")
}

// Delete 删除
func (h *Handler) Delete(c *gin.Context) {
    h.proxyPost(c, "delete")
}

// Rename 重命名
func (h *Handler) Rename(c *gin.Context) {
    h.proxyPost(c, "rename")
}
```

## 6. 路由注册

在 `internal/squ-apiserver/router/router.go` 中添加：

```go
import (
    fshandler "squirrel-dev/internal/squ-apiserver/handler/fs"
)

func InitRouter(r *gin.Engine, conf *config.Config, db database.DB) {
    // ... 现有代码 ...
    
    // 文件管理路由
    api := r.Group("/api/v1")
    api.Use(auth.Middleware()) // 认证中间件
    {
        // 文件管理
        fshandler.RegisterRoutes(api, repository.NewServerRepository(db))
    }
}
```

## 7. 测试用例

### 测试目录列表

```bash
# 获取服务器 ID 为 1 的目录列表
curl "http://localhost:10700/api/v1/fs/1/list?path=/home"
```

### 测试文件下载

```bash
# 下载文件
curl -O "http://localhost:10700/api/v1/fs/1/download?path=/home/user/config.yaml"
```

### 测试文件上传

```bash
# 上传文件到服务器 ID 为 1
curl -X POST "http://localhost:10700/api/v1/fs/1/upload?path=/home/user/" \
  -F "file=@./local-file.txt"
```

### 测试文件写入

```bash
# 写入文件
curl -X POST "http://localhost:10700/api/v1/fs/1/write" \
  -H "Content-Type: application/json" \
  -d '{"path": "/home/user/test.txt", "content": "Hello World"}'
```

## 8. 完成检查

- [ ] Handler 结构体创建
- [ ] Agent 客户端实现
- [ ] GET 请求代理实现
- [ ] POST 请求代理实现
- [ ] 流式下载实现
- [ ] 流式上传实现
- [ ] 路由注册完成
- [ ] 集成测试通过

## 9. 下一阶段

完成本阶段后，可通过 APIServer 访问 Agent 文件管理功能。继续开发 [04-frontend-filemanager.md](./04-frontend-filemanager.md) 实现前端页面。
