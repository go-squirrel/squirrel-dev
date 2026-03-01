# Phase 2: Agent 写入功能

## 前置条件

已完成 Phase 1（01-agent-core.md），具备：
- 目录浏览功能
- 文件读取功能
- 安全校验模块
- Service 结构体

## 开发目标

实现 Agent 端文件写入操作：
- 文件写入（保存内容）
- 文件上传（接收上传）
- 创建文件/目录
- 删除文件/目录
- 重命名文件/目录

## 1. 新增目录结构

```
internal/squ-agent/
└── handler/
    └── fs/
        ├── write.go       # 文件写入（新增）
        ├── upload.go      # 文件上传（新增）
        ├── create.go      # 创建文件/目录（新增）
        ├── delete.go      # 删除操作（新增）
        ├── rename.go      # 重命名操作（新增）
        └── types.go       # 新增请求/响应结构
```

## 2. 新增数据模型

```go
// internal/squ-agent/handler/fs/types.go 新增

// WriteRequest 文件写入请求
type WriteRequest struct {
    Path    string `json:"path" binding:"required"`
    Content string `json:"content"`       // 文件内容
    Mode    string `json:"mode"`          // 权限模式，如 "644"
    Backup  bool   `json:"backup"`        // 是否备份原文件
}

// WriteResponse 文件写入响应
type WriteResponse struct {
    Path     string `json:"path"`      // 文件路径
    Size     int64  `json:"size"`      // 写入字节数
    Backup   string `json:"backup"`    // 备份文件路径（如果有）
}

// MkdirRequest 创建目录请求
type MkdirRequest struct {
    Path string `json:"path" binding:"required"`
    Mode string `json:"mode"`  // 权限模式
}

// CreateFileRequest 创建文件请求
type CreateFileRequest struct {
    Path    string `json:"path" binding:"required"`
    Content string `json:"content"`  // 初始内容
    Mode    string `json:"mode"`     // 权限模式
}

// DeleteRequest 删除请求
type DeleteRequest struct {
    Path    string `json:"path" binding:"required"`
    Recursive bool `json:"recursive"`  // 递归删除（目录）
}

// RenameRequest 重命名请求
type RenameRequest struct {
    OldPath string `json:"old_path" binding:"required"`
    NewPath string `json:"new_path" binding:"required"`
}

// UploadResponse 上传响应
type UploadResponse struct {
    Path string `json:"path"`  // 保存路径
    Size int64  `json:"size"`  // 文件大小
    Name string `json:"name"`  // 文件名
}
```

## 3. 文件写入实现

```go
// internal/squ-agent/handler/fs/write.go

package fs

import (
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
    "time"
    
    "github.com/gin-gonic/gin"
    
    "squirrel-dev/internal/pkg/response"
)

// Write 文件写入
func (s *Service) Write(c *gin.Context) {
    var req WriteRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, 400, "参数错误: "+err.Error())
        return
    }
    
    // 安全校验
    if err := s.security.ValidatePath(req.Path); err != nil {
        response.Error(c, 403, err.Error())
        return
    }
    
    // 检查父目录是否存在
    parentDir := filepath.Dir(req.Path)
    if _, err := os.Stat(parentDir); os.IsNotExist(err) {
        response.Error(c, 400, "父目录不存在")
        return
    }
    
    // 备份原文件
    var backupPath string
    if req.Backup {
        if _, err := os.Stat(req.Path); err == nil {
            backupPath = fmt.Sprintf("%s.bak.%d", req.Path, time.Now().Unix())
            if err := copyFile(req.Path, backupPath); err != nil {
                response.Error(c, 500, "备份文件失败: "+err.Error())
                return
            }
        }
    }
    
    // 写入文件
    content := []byte(req.Content)
    if err := ioutil.WriteFile(req.Path, content, 0644); err != nil {
        response.Error(c, 500, "写入文件失败: "+err.Error())
        return
    }
    
    // 设置权限
    if req.Mode != "" {
        mode := parseFileMode(req.Mode)
        if err := os.Chmod(req.Path, mode); err != nil {
            // 权限设置失败不影响写入结果，记录日志即可
        }
    }
    
    response.Success(c, WriteResponse{
        Path:   req.Path,
        Size:   int64(len(content)),
        Backup: backupPath,
    })
}

// copyFile 复制文件
func copyFile(src, dst string) error {
    input, err := ioutil.ReadFile(src)
    if err != nil {
        return err
    }
    return ioutil.WriteFile(dst, input, 0644)
}

// parseFileMode 解析权限模式字符串
func parseFileMode(mode string) os.FileMode {
    var m uint32
    fmt.Sscanf(mode, "%o", &m)
    return os.FileMode(m)
}
```

## 4. 文件上传实现

```go
// internal/squ-agent/handler/fs/upload.go

package fs

import (
    "io"
    "os"
    "path/filepath"
    
    "github.com/gin-gonic/gin"
    
    "squirrel-dev/internal/pkg/response"
)

// Upload 文件上传
func (s *Service) Upload(c *gin.Context) {
    // 获取目标路径
    targetPath := c.PostForm("path")
    if targetPath == "" {
        response.Error(c, 400, "path 参数必填")
        return
    }
    
    // 安全校验
    if err := s.security.ValidatePath(targetPath); err != nil {
        response.Error(c, 403, err.Error())
        return
    }
    
    // 获取上传的文件
    file, header, err := c.Request.FormFile("file")
    if err != nil {
        response.Error(c, 400, "获取上传文件失败: "+err.Error())
        return
    }
    defer file.Close()
    
    // 如果目标是目录，使用原文件名
    info, err := os.Stat(targetPath)
    if err == nil && info.IsDir() {
        targetPath = filepath.Join(targetPath, header.Filename)
    }
    
    // 再次校验最终路径
    if err := s.security.ValidatePath(targetPath); err != nil {
        response.Error(c, 403, err.Error())
        return
    }
    
    // 创建目标文件
    dst, err := os.Create(targetPath)
    if err != nil {
        response.Error(c, 500, "创建文件失败: "+err.Error())
        return
    }
    defer dst.Close()
    
    // 流式写入
    written, err := io.Copy(dst, file)
    if err != nil {
        response.Error(c, 500, "写入文件失败: "+err.Error())
        return
    }
    
    response.Success(c, UploadResponse{
        Path: targetPath,
        Size: written,
        Name: filepath.Base(targetPath),
    })
}
```

## 5. 创建操作实现

```go
// internal/squ-agent/handler/fs/create.go

package fs

import (
    "os"
    
    "github.com/gin-gonic/gin"
    
    "squirrel-dev/internal/pkg/response"
)

// Mkdir 创建目录
func (s *Service) Mkdir(c *gin.Context) {
    var req MkdirRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, 400, "参数错误: "+err.Error())
        return
    }
    
    // 安全校验
    if err := s.security.ValidatePath(req.Path); err != nil {
        response.Error(c, 403, err.Error())
        return
    }
    
    // 创建目录（包括父目录）
    mode := os.FileMode(0755)
    if req.Mode != "" {
        mode = parseFileMode(req.Mode)
    }
    
    if err := os.MkdirAll(req.Path, mode); err != nil {
        response.Error(c, 500, "创建目录失败: "+err.Error())
        return
    }
    
    response.Success(c, gin.H{
        "path": req.Path,
    })
}

// CreateFile 创建空文件
func (s *Service) CreateFile(c *gin.Context) {
    var req CreateFileRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, 400, "参数错误: "+err.Error())
        return
    }
    
    // 安全校验
    if err := s.security.ValidatePath(req.Path); err != nil {
        response.Error(c, 403, err.Error())
        return
    }
    
    // 检查文件是否已存在
    if _, err := os.Stat(req.Path); err == nil {
        response.Error(c, 400, "文件已存在")
        return
    }
    
    // 创建文件
    mode := os.FileMode(0644)
    if req.Mode != "" {
        mode = parseFileMode(req.Mode)
    }
    
    file, err := os.OpenFile(req.Path, os.O_CREATE|os.O_WRONLY, mode)
    if err != nil {
        response.Error(c, 500, "创建文件失败: "+err.Error())
        return
    }
    defer file.Close()
    
    // 写入初始内容
    if req.Content != "" {
        if _, err := file.WriteString(req.Content); err != nil {
            response.Error(c, 500, "写入内容失败: "+err.Error())
            return
        }
    }
    
    response.Success(c, gin.H{
        "path": req.Path,
        "size": int64(len(req.Content)),
    })
}
```

## 6. 删除操作实现

```go
// internal/squ-agent/handler/fs/delete.go

package fs

import (
    "os"
    
    "github.com/gin-gonic/gin"
    
    "squirrel-dev/internal/pkg/response"
)

// Delete 删除文件/目录
func (s *Service) Delete(c *gin.Context) {
    var req DeleteRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, 400, "参数错误: "+err.Error())
        return
    }
    
    // 安全校验
    if err := s.security.ValidatePath(req.Path); err != nil {
        response.Error(c, 403, err.Error())
        return
    }
    
    // 检查路径是否存在
    info, err := os.Stat(req.Path)
    if err != nil {
        response.Error(c, 404, "路径不存在")
        return
    }
    
    // 删除操作
    if info.IsDir() {
        if req.Recursive {
            err = os.RemoveAll(req.Path)
        } else {
            err = os.Remove(req.Path)
        }
    } else {
        err = os.Remove(req.Path)
    }
    
    if err != nil {
        response.Error(c, 500, "删除失败: "+err.Error())
        return
    }
    
    response.Success(c, gin.H{
        "path": req.Path,
    })
}
```

## 7. 重命名操作实现

```go
// internal/squ-agent/handler/fs/rename.go

package fs

import (
    "os"
    
    "github.com/gin-gonic/gin"
    
    "squirrel-dev/internal/pkg/response"
)

// Rename 重命名文件/目录
func (s *Service) Rename(c *gin.Context) {
    var req RenameRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, 400, "参数错误: "+err.Error())
        return
    }
    
    // 安全校验 - 旧路径
    if err := s.security.ValidatePath(req.OldPath); err != nil {
        response.Error(c, 403, "旧路径: "+err.Error())
        return
    }
    
    // 安全校验 - 新路径
    if err := s.security.ValidatePath(req.NewPath); err != nil {
        response.Error(c, 403, "新路径: "+err.Error())
        return
    }
    
    // 检查旧路径是否存在
    if _, err := os.Stat(req.OldPath); os.IsNotExist(err) {
        response.Error(c, 404, "原路径不存在")
        return
    }
    
    // 检查新路径是否已存在
    if _, err := os.Stat(req.NewPath); err == nil {
        response.Error(c, 400, "目标路径已存在")
        return
    }
    
    // 重命名
    if err := os.Rename(req.OldPath, req.NewPath); err != nil {
        response.Error(c, 500, "重命名失败: "+err.Error())
        return
    }
    
    response.Success(c, gin.H{
        "old_path": req.OldPath,
        "new_path": req.NewPath,
    })
}
```

## 8. 更新路由注册

```go
// internal/squ-agent/handler/fs/handler.go 更新

// RegisterRoutes 注册路由
func RegisterRoutes(r *gin.RouterGroup) {
    fs := r.Group("/fs")
    {
        // Phase 1: 读取操作
        fs.GET("/list", svc.List)         // 目录列表
        fs.GET("/read", svc.Read)         // 文件读取
        fs.GET("/download", svc.Download) // 文件下载
        
        // Phase 2: 写入操作
        fs.POST("/write", svc.Write)           // 文件写入
        fs.POST("/upload", svc.Upload)         // 文件上传
        fs.POST("/mkdir", svc.Mkdir)           // 创建目录
        fs.POST("/create", svc.CreateFile)     // 创建文件
        fs.POST("/delete", svc.Delete)         // 删除
        fs.POST("/rename", svc.Rename)         // 重命名
    }
}
```

## 9. 测试用例

### 测试文件写入

```bash
# 写入文件
curl -X POST "http://localhost:10750/api/agent/v1/fs/write" \
  -H "Content-Type: application/json" \
  -d '{"path": "/home/user/test.txt", "content": "Hello World", "backup": true}'
```

### 测试文件上传

```bash
# 上传文件
curl -X POST "http://localhost:10750/api/agent/v1/fs/upload" \
  -F "path=/home/user/" \
  -F "file=@./local-file.txt"
```

### 测试创建目录

```bash
# 创建目录
curl -X POST "http://localhost:10750/api/agent/v1/fs/mkdir" \
  -H "Content-Type: application/json" \
  -d '{"path": "/home/user/newdir"}'
```

### 测试删除

```bash
# 删除文件
curl -X POST "http://localhost:10750/api/agent/v1/fs/delete" \
  -H "Content-Type: application/json" \
  -d '{"path": "/home/user/test.txt"}'

# 递归删除目录
curl -X POST "http://localhost:10750/api/agent/v1/fs/delete" \
  -H "Content-Type: application/json" \
  -d '{"path": "/home/user/newdir", "recursive": true}'
```

### 测试重命名

```bash
# 重命名
curl -X POST "http://localhost:10750/api/agent/v1/fs/rename" \
  -H "Content-Type: application/json" \
  -d '{"old_path": "/home/user/test.txt", "new_path": "/home/user/renamed.txt"}'
```

## 10. 完成检查

- [ ] 文件写入功能实现
- [ ] 文件上传功能实现
- [ ] 创建目录功能实现
- [ ] 创建文件功能实现
- [ ] 删除功能实现
- [ ] 重命名功能实现
- [ ] 路由更新完成
- [ ] 单元测试通过
- [ ] 集成测试通过

## 11. 下一阶段

完成本阶段后，Agent 端文件管理功能已完整。继续开发 [03-apiserver-proxy.md](./03-apiserver-proxy.md) 实现 APIServer 代理转发。
