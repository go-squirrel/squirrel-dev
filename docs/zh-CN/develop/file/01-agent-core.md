# Phase 1: Agent 核心功能

## 开发目标

实现 Agent 端文件管理的核心功能：
- 目录浏览（列出目录内容）
- 文件读取（文本预览）
- 文件下载（流式传输）
- 安全校验（路径沙箱限制）

## 1. 目录结构

```
internal/squ-agent/
├── handler/
│   └── fs/
│       ├── handler.go        # Handler 入口，路由注册
│       ├── service.go        # 业务逻辑结构体
│       ├── security.go       # 安全校验
│       ├── list.go           # 目录列表
│       ├── read.go           # 文件读取
│       └── download.go       # 文件下载
├── model/
│   └── fs_config.go          # 文件管理配置
└── router/
    └── fs.go                 # 路由注册（或在现有router中添加）
```

## 2. 配置设计

### 2.1 配置结构

```go
// internal/squ-agent/model/fs_config.go

package model

// FSConfig 文件管理配置
type FSConfig struct {
    Enabled      bool     `yaml:"enabled" json:"enabled"`
    AllowedPaths []string `yaml:"allowed_paths" json:"allowed_paths"`
    DeniedPaths  []string `yaml:"denied_paths" json:"denied_paths"`
    DeniedExts   []string `yaml:"denied_exts" json:"denied_exts"`
    MaxReadSize  int64    `yaml:"max_read_size" json:"max_read_size"` // MB
}

// DefaultFSConfig 默认配置
func DefaultFSConfig() FSConfig {
    return FSConfig{
        Enabled: true,
        AllowedPaths: []string{
            "/home",
            "/opt",
            "/var/log",
            "/tmp",
        },
        DeniedPaths: []string{
            "/etc/shadow",
            "/etc/passwd",
            "/root/.ssh",
        },
        DeniedExts: []string{
            ".key",
            ".pem",
        },
        MaxReadSize: 10, // 10MB
    }
}
```

### 2.2 配置文件

```yaml
# config/agent.yaml 新增
filesystem:
  enabled: true
  allowed_paths:
    - /home
    - /opt
    - /var/log
  denied_paths:
    - /etc/shadow
    - /root/.ssh
  denied_exts:
    - .key
    - .pem
  max_read_size: 10  # MB
```

## 3. 数据模型

```go
// internal/squ-agent/handler/fs/types.go

package fs

import "time"

// FileInfo 文件/目录信息
type FileInfo struct {
    Name      string    `json:"name"`       // 文件名
    Path      string    `json:"path"`       // 完整路径
    IsDir     bool      `json:"is_dir"`     // 是否目录
    Size      int64     `json:"size"`       // 大小（字节）
    Mode      string    `json:"mode"`       // 权限（如 755）
    ModTime   time.Time `json:"mod_time"`   // 修改时间
    Extension string    `json:"extension"`  // 扩展名
    Mime      string    `json:"mime"`       // MIME 类型
}

// ListRequest 目录列表请求
type ListRequest struct {
    Path       string `form:"path" binding:"required"`
    ShowHidden bool   `form:"show_hidden"`  // 是否显示隐藏文件
    SortBy     string `form:"sort_by"`      // name, size, time
    Order      string `form:"order"`        // asc, desc
}

// ListResponse 目录列表响应
type ListResponse struct {
    Path    string     `json:"path"`     // 当前路径
    Parent  string     `json:"parent"`   // 父目录路径
    Entries []FileInfo `json:"entries"`  // 文件列表
}

// ReadRequest 文件读取请求
type ReadRequest struct {
    Path   string `form:"path" binding:"required"`
    Offset int64  `form:"offset"`  // 读取偏移量（可选）
    Limit  int64  `form:"limit"`   // 读取字节数（可选）
}

// ReadResponse 文件读取响应
type ReadResponse struct {
    Path    string `json:"path"`      // 文件路径
    Content string `json:"content"`   // 文件内容（Base64 或原文）
    Size    int64  `json:"size"`      // 文件总大小
    Mime    string `json:"mime"`      // MIME 类型
    Encoding string `json:"encoding"` // 编码（utf-8, gbk 等）
}
```

## 4. 安全模块实现

```go
// internal/squ-agent/handler/fs/security.go

package fs

import (
    "fmt"
    "path/filepath"
    "strings"
    
    "squirrel-dev/internal/squ-agent/model"
)

// Security 安全校验器
type Security struct {
    config model.FSConfig
}

// NewSecurity 创建安全校验器
func NewSecurity(config model.FSConfig) *Security {
    return &Security{config: config}
}

// ValidatePath 验证路径是否允许访问
func (s *Security) ValidatePath(requestedPath string) error {
    // 1. 解析为绝对路径
    absPath, err := filepath.Abs(requestedPath)
    if err != nil {
        return fmt.Errorf("无效的路径: %w", err)
    }
    
    // 2. 检查是否在禁止列表
    for _, denied := range s.config.DeniedPaths {
        if strings.HasPrefix(absPath, denied) {
            return fmt.Errorf("禁止访问该路径: %s", absPath)
        }
    }
    
    // 3. 检查是否在允许列表
    if len(s.config.AllowedPaths) > 0 {
        allowed := false
        for _, path := range s.config.AllowedPaths {
            if strings.HasPrefix(absPath, path) {
                allowed = true
                break
            }
        }
        if !allowed {
            return fmt.Errorf("路径不在允许范围内: %s", absPath)
        }
    }
    
    // 4. 检查路径穿越攻击（处理符号链接）
    realPath, err := filepath.EvalSymlinks(absPath)
    if err == nil {
        // 再次验证真实路径
        for _, denied := range s.config.DeniedPaths {
            if strings.HasPrefix(realPath, denied) {
                return fmt.Errorf("禁止访问该路径: %s", realPath)
            }
        }
    }
    
    return nil
}

// ValidateExtension 验证文件扩展名是否允许
func (s *Security) ValidateExtension(path string) error {
    ext := strings.ToLower(filepath.Ext(path))
    for _, denied := range s.config.DeniedExts {
        if ext == denied {
            return fmt.Errorf("禁止访问该类型文件: %s", ext)
        }
    }
    return nil
}

// ValidateReadSize 验证文件大小是否允许读取
func (s *Security) ValidateReadSize(size int64) error {
    maxSize := s.config.MaxReadSize * 1024 * 1024
    if size > maxSize {
        return fmt.Errorf("文件过大，超过最大读取限制 %dMB", s.config.MaxReadSize)
    }
    return nil
}
```

## 5. 目录列表实现

```go
// internal/squ-agent/handler/fs/list.go

package fs

import (
    "os"
    "sort"
    "strings"
    
    "github.com/gin-gonic/gin"
    
    "squirrel-dev/internal/pkg/response"
)

// List 目录列表
func (s *Service) List(c *gin.Context) {
    var req ListRequest
    if err := c.ShouldBindQuery(&req); err != nil {
        response.Error(c, 400, "参数错误: "+err.Error())
        return
    }
    
    // 安全校验
    if err := s.security.ValidatePath(req.Path); err != nil {
        response.Error(c, 403, err.Error())
        return
    }
    
    // 读取目录
    entries, err := s.readDir(req.Path, req.ShowHidden)
    if err != nil {
        response.Error(c, 500, "读取目录失败: "+err.Error())
        return
    }
    
    // 排序
    s.sortEntries(entries, req.SortBy, req.Order)
    
    // 获取父目录
    parent := ""
    if req.Path != "/" {
        parent = filepath.Dir(req.Path)
    }
    
    response.Success(c, ListResponse{
        Path:    req.Path,
        Parent:  parent,
        Entries: entries,
    })
}

// readDir 读取目录内容
func (s *Service) readDir(path string, showHidden bool) ([]FileInfo, error) {
    dir, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer dir.Close()
    
    files, err := dir.Readdir(-1)
    if err != nil {
        return nil, err
    }
    
    var entries []FileInfo
    for _, f := range files {
        // 过滤隐藏文件
        if !showHidden && strings.HasPrefix(f.Name(), ".") {
            continue
        }
        
        fullPath := filepath.Join(path, f.Name())
        entries = append(entries, FileInfo{
            Name:      f.Name(),
            Path:      fullPath,
            IsDir:     f.IsDir(),
            Size:      f.Size(),
            Mode:      f.Mode().String(),
            ModTime:   f.ModTime(),
            Extension: strings.ToLower(filepath.Ext(f.Name())),
            Mime:      getMimeByExtension(f.Name()),
        })
    }
    
    return entries, nil
}

// sortEntries 排序文件列表
func (s *Service) sortEntries(entries []FileInfo, sortBy, order string) {
    sort.Slice(entries, func(i, j int) bool {
        // 目录始终在前
        if entries[i].IsDir != entries[j].IsDir {
            return entries[i].IsDir
        }
        
        var less bool
        switch sortBy {
        case "size":
            less = entries[i].Size < entries[j].Size
        case "time":
            less = entries[i].ModTime.Before(entries[j].ModTime)
        default: // name
            less = entries[i].Name < entries[j].Name
        }
        
        if order == "desc" {
            return !less
        }
        return less
    })
}
```

## 6. 文件读取实现

```go
// internal/squ-agent/handler/fs/read.go

package fs

import (
    "encoding/base64"
    "io/ioutil"
    "os"
    
    "github.com/gin-gonic/gin"
    
    "squirrel-dev/internal/pkg/response"
)

// Read 文件读取（文本预览）
func (s *Service) Read(c *gin.Context) {
    var req ReadRequest
    if err := c.ShouldBindQuery(&req); err != nil {
        response.Error(c, 400, "参数错误: "+err.Error())
        return
    }
    
    // 安全校验
    if err := s.security.ValidatePath(req.Path); err != nil {
        response.Error(c, 403, err.Error())
        return
    }
    if err := s.security.ValidateExtension(req.Path); err != nil {
        response.Error(c, 403, err.Error())
        return
    }
    
    // 获取文件信息
    info, err := os.Stat(req.Path)
    if err != nil {
        response.Error(c, 404, "文件不存在: "+err.Error())
        return
    }
    
    // 检查是否为目录
    if info.IsDir() {
        response.Error(c, 400, "不能读取目录")
        return
    }
    
    // 检查文件大小
    if err := s.security.ValidateReadSize(info.Size()); err != nil {
        response.Error(c, 400, err.Error())
        return
    }
    
    // 读取文件
    content, err := ioutil.ReadFile(req.Path)
    if err != nil {
        response.Error(c, 500, "读取文件失败: "+err.Error())
        return
    }
    
    // 检测编码
    encoding := detectEncoding(content)
    
    // 如果是二进制文件，返回 Base64
    isText := isTextFile(content)
    if !isText {
        content = []byte(base64.StdEncoding.EncodeToString(content))
    }
    
    response.Success(c, ReadResponse{
        Path:     req.Path,
        Content:  string(content),
        Size:     info.Size(),
        Mime:     getMimeByExtension(req.Path),
        Encoding: encoding,
    })
}
```

## 7. 文件下载实现

```go
// internal/squ-agent/handler/fs/download.go

package fs

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    
    "github.com/gin-gonic/gin"
)

// Download 文件下载（流式传输）
func (s *Service) Download(c *gin.Context) {
    path := c.Query("path")
    if path == "" {
        c.JSON(400, gin.H{"error": "path 参数必填"})
        return
    }
    
    // 安全校验
    if err := s.security.ValidatePath(path); err != nil {
        c.JSON(403, gin.H{"error": err.Error()})
        return
    }
    
    // 打开文件
    file, err := os.Open(path)
    if err != nil {
        c.JSON(404, gin.H{"error": "文件不存在"})
        return
    }
    defer file.Close()
    
    // 获取文件信息
    info, err := file.Stat()
    if err != nil {
        c.JSON(500, gin.H{"error": "获取文件信息失败"})
        return
    }
    
    // 如果是目录，拒绝下载
    if info.IsDir() {
        c.JSON(400, gin.H{"error": "不能下载目录"})
        return
    }
    
    // 设置响应头
    filename := filepath.Base(path)
    c.Header("Content-Description", "File Transfer")
    c.Header("Content-Type", "application/octet-stream")
    c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
    c.Header("Content-Transfer-Encoding", "binary")
    c.Header("Content-Length", fmt.Sprintf("%d", info.Size()))
    
    // 流式传输
    io.Copy(c.Writer, file)
}
```

## 8. Service 和 Handler 入口

```go
// internal/squ-agent/handler/fs/service.go

package fs

import "squirrel-dev/internal/squ-agent/model"

// Service 文件管理服务
type Service struct {
    security *Security
    config   model.FSConfig
}

// NewService 创建服务
func NewService(config model.FSConfig) *Service {
    return &Service{
        security: NewSecurity(config),
        config:   config,
    }
}
```

```go
// internal/squ-agent/handler/fs/handler.go

package fs

import (
    "github.com/gin-gonic/gin"
    
    "squirrel-dev/internal/squ-agent/model"
)

var svc *Service

// Init 初始化文件管理模块
func Init(config model.FSConfig) {
    svc = NewService(config)
}

// RegisterRoutes 注册路由
func RegisterRoutes(r *gin.RouterGroup) {
    fs := r.Group("/fs")
    {
        fs.GET("/list", svc.List)         // 目录列表
        fs.GET("/read", svc.Read)         // 文件读取
        fs.GET("/download", svc.Download) // 文件下载
    }
}
```

## 9. 路由注册

在 `internal/squ-agent/router/router.go` 中添加：

```go
// 在现有的路由组中添加文件管理路由
func InitRouter(r *gin.Engine, conf *config.Config) {
    // ... 现有代码 ...
    
    // 文件管理路由（需要 mTLS 认证）
    fsGroup := api.Group("/agent/v1")
    fsGroup.Use(mtls.Middleware()) // 如果有 mTLS 中间件
    fs.RegisterRoutes(fsGroup)
}
```

## 10. 测试用例

### 测试目录列表

```bash
# 列出 /home 目录
curl "http://localhost:10750/api/agent/v1/fs/list?path=/home"
```

### 测试文件读取

```bash
# 读取文件内容
curl "http://localhost:10750/api/agent/v1/fs/read?path=/home/user/config.yaml"
```

### 测试文件下载

```bash
# 下载文件
curl -O "http://localhost:10750/api/agent/v1/fs/download?path=/home/user/app.log"
```

## 11. 完成检查

- [ ] 配置结构定义完成
- [ ] 安全校验模块实现
- [ ] 目录列表功能实现
- [ ] 文件读取功能实现
- [ ] 文件下载功能实现
- [ ] 路由注册完成
- [ ] 单元测试通过
- [ ] 集成测试通过

## 12. 下一阶段

完成本阶段后，继续开发 [02-agent-write.md](./02-agent-write.md) 实现文件写入操作。
