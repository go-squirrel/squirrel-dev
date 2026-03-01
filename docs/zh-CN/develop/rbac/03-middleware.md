# Phase 3: 权限中间件

## 开发目标

实现权限检查中间件，集成到现有路由中：
- RBACMiddleware 中间件
- 现有路由改造
- 权限检查失败处理

## 1. 目录结构

```
internal/squ-apiserver/
├── middleware/
│   ├── auth.go           # 已存在：JWT 认证
│   └── rbac.go           # 新增：RBAC 权限检查
└── router/
    ├── router.go         # 路由注册入口
    ├── server.go         # 服务器路由
    ├── application.go    # 应用路由
    ├── monitor.go        # 监控路由
    ├── script.go         # 脚本路由
    ├── file.go           # 文件路由
    └── ...
```

## 2. 权限中间件实现

### 2.1 基础中间件

```go
// internal/squ-apiserver/middleware/rbac.go

// RBACMiddleware 权限检查中间件
// resource: 资源类型，如 "server"、"file"
// action: 操作类型，如 "list"、"create"
func RBACMiddleware(resource, action string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 获取用户 ID（由 JWT 中间件设置）
        userID, exists := c.Get("user_id")
        if !exists {
            c.AbortWithStatusJSON(401, gin.H{
                "code": 401,
                "msg":  "未登录",
            })
            return
        }
        
        // 2. 获取权限服务
        rbacService := rbac.GetRBACService()
        if rbacService == nil {
            c.AbortWithStatusJSON(500, gin.H{
                "code": 500,
                "msg":  "权限服务未初始化",
            })
            return
        }
        
        // 3. 检查权限
        hasPermission, err := rbacService.CheckPermission(userID.(uint), resource, action)
        if err != nil {
            c.AbortWithStatusJSON(500, gin.H{
                "code": 500,
                "msg":  "权限检查失败",
            })
            return
        }
        
        if !hasPermission {
            c.AbortWithStatusJSON(403, gin.H{
                "code": 403,
                "msg":  "权限不足",
            })
            return
        }
        
        c.Next()
    }
}
```

### 2.2 带服务器 ID 的权限中间件

```go
// internal/squ-apiserver/middleware/rbac.go

// RBACMiddlewareWithServer 带服务器范围的权限检查
// 用于需要检查用户是否有权访问特定服务器的场景
func RBACMiddlewareWithServer(resource, action string, serverIDParam string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID, exists := c.Get("user_id")
        if !exists {
            c.AbortWithStatusJSON(401, gin.H{"code": 401, "msg": "未登录"})
            return
        }
        
        // 1. 先检查基础权限
        rbacService := rbac.GetRBACService()
        hasPermission, err := rbacService.CheckPermission(userID.(uint), resource, action)
        if err != nil || !hasPermission {
            c.AbortWithStatusJSON(403, gin.H{"code": 403, "msg": "权限不足"})
            return
        }
        
        // 2. 检查服务器范围（如果启用了服务器级权限）
        serverIDStr := c.Param(serverIDParam)
        if serverIDStr != "" {
            serverID, err := strconv.ParseUint(serverIDStr, 10, 64)
            if err == nil {
                // Phase 6 实现：检查用户是否有权访问该服务器
                // hasScope := rbacService.CheckServerScope(userID.(uint), uint(serverID))
                // if !hasScope {
                //     c.AbortWithStatusJSON(403, gin.H{"code": 403, "msg": "无权访问该服务器"})
                //     return
                // }
            }
        }
        
        c.Next()
    }
}
```

### 2.3 多权限检查中间件

```go
// internal/squ-apiserver/middleware/rbac.go

// RBACAnyMiddleware 检查用户是否有任意一个权限
func RBACAnyMiddleware(permissionCodes []string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID, exists := c.Get("user_id")
        if !exists {
            c.AbortWithStatusJSON(401, gin.H{"code": 401, "msg": "未登录"})
            return
        }
        
        rbacService := rbac.GetRBACService()
        hasPermission, err := rbacService.CheckAnyPermission(userID.(uint), permissionCodes)
        if err != nil || !hasPermission {
            c.AbortWithStatusJSON(403, gin.H{"code": 403, "msg": "权限不足"})
            return
        }
        
        c.Next()
    }
}
```

## 3. 现有路由改造

### 3.1 服务器路由

```go
// internal/squ-apiserver/router/server.go

func Server(group *gin.RouterGroup) {
    servers := group.Group("/servers")
    
    // 列表：需要 server:list 权限
    servers.GET("", 
        middleware.RBACMiddleware(model.ResourceServer, model.ActionList),
        server.List,
    )
    
    // 详情：需要 server:get 权限
    servers.GET("/:id",
        middleware.RBACMiddleware(model.ResourceServer, model.ActionGet),
        server.Get,
    )
    
    // 创建：需要 server:create 权限
    servers.POST("",
        middleware.RBACMiddleware(model.ResourceServer, model.ActionCreate),
        server.Create,
    )
    
    // 更新：需要 server:update 权限
    servers.PUT("/:id",
        middleware.RBACMiddleware(model.ResourceServer, model.ActionUpdate),
        server.Update,
    )
    
    // 删除：需要 server:delete 权限
    servers.DELETE("/:id",
        middleware.RBACMiddleware(model.ResourceServer, model.ActionDelete),
        server.Delete,
    )
    
    // 终端：需要 terminal:exec 权限
    servers.GET("/:id/terminal",
        middleware.RBACMiddleware(model.ResourceTerminal, model.ActionExec),
        terminal.Handler,
    )
    
    // 终端 WebSocket
    servers.GET("/:id/terminal/ws",
        middleware.RBACMiddleware(model.ResourceTerminal, model.ActionExec),
        terminal.WSHandler,
    )
}
```

### 3.2 应用路由

```go
// internal/squ-apiserver/router/application.go

func Application(group *gin.RouterGroup) {
    apps := group.Group("/applications")
    
    apps.GET("",
        middleware.RBACMiddleware(model.ResourceApplication, model.ActionList),
        application.List,
    )
    
    apps.GET("/:id",
        middleware.RBACMiddleware(model.ResourceApplication, model.ActionGet),
        application.Get,
    )
    
    apps.POST("",
        middleware.RBACMiddleware(model.ResourceApplication, model.ActionCreate),
        application.Create,
    )
    
    apps.PUT("/:id",
        middleware.RBACMiddleware(model.ResourceApplication, model.ActionUpdate),
        application.Update,
    )
    
    apps.DELETE("/:id",
        middleware.RBACMiddleware(model.ResourceApplication, model.ActionDelete),
        application.Delete,
    )
    
    // 部署：需要 application:deploy 权限
    apps.POST("/:id/deploy",
        middleware.RBACMiddleware(model.ResourceApplication, model.ActionExec),
        application.Deploy,
    )
    
    apps.POST("/:id/stop",
        middleware.RBACMiddleware(model.ResourceApplication, model.ActionExec),
        application.Stop,
    )
    
    apps.POST("/:id/restart",
        middleware.RBACMiddleware(model.ResourceApplication, model.ActionExec),
        application.Restart,
    )
}
```

### 3.3 监控路由

```go
// internal/squ-apiserver/router/monitor.go

func Monitor(group *gin.RouterGroup) {
    monitor := group.Group("/monitor")
    
    // 实时数据
    monitor.GET("/stats/:serverId",
        middleware.RBACMiddleware(model.ResourceMonitor, model.ActionGet),
        monitor.GetStats,
    )
    
    // 历史数据
    monitor.GET("/base/history/:serverId",
        middleware.RBACMiddleware(model.ResourceMonitor, model.ActionList),
        monitor.GetBaseHistory,
    )
    
    monitor.GET("/diskio/history/:serverId",
        middleware.RBACMiddleware(model.ResourceMonitor, model.ActionList),
        monitor.GetDiskIOHistory,
    )
    
    monitor.GET("/netio/history/:serverId",
        middleware.RBACMiddleware(model.ResourceMonitor, model.ActionList),
        monitor.GetNetIOHistory,
    )
}
```

### 3.4 脚本路由

```go
// internal/squ-apiserver/router/script.go

func Script(group *gin.RouterGroup) {
    scripts := group.Group("/scripts")
    
    scripts.GET("",
        middleware.RBACMiddleware(model.ResourceScript, model.ActionList),
        script.List,
    )
    
    scripts.GET("/:id",
        middleware.RBACMiddleware(model.ResourceScript, model.ActionGet),
        script.Get,
    )
    
    scripts.POST("",
        middleware.RBACMiddleware(model.ResourceScript, model.ActionCreate),
        script.Create,
    )
    
    scripts.PUT("/:id",
        middleware.RBACMiddleware(model.ResourceScript, model.ActionUpdate),
        script.Update,
    )
    
    scripts.DELETE("/:id",
        middleware.RBACMiddleware(model.ResourceScript, model.ActionDelete),
        script.Delete,
    )
    
    // 执行脚本
    scripts.POST("/:id/exec",
        middleware.RBACMiddleware(model.ResourceScript, model.ActionExec),
        script.Exec,
    )
}
```

### 3.5 文件路由

```go
// internal/squ-apiserver/router/file.go

func File(group *gin.RouterGroup) {
    files := group.Group("/servers/:id/files")
    
    // 浏览目录
    files.GET("",
        middleware.RBACMiddleware(model.ResourceFile, model.ActionList),
        file.List,
    )
    
    // 读取文件
    files.GET("/read",
        middleware.RBACMiddleware(model.ResourceFile, model.ActionGet),
        file.Read,
    )
    
    // 编辑文件
    files.PUT("/write",
        middleware.RBACMiddleware(model.ResourceFile, model.ActionUpdate),
        file.Write,
    )
    
    // 上传文件
    files.POST("/upload",
        middleware.RBACMiddleware(model.ResourceFile, model.ActionUpload),
        file.Upload,
    )
    
    // 下载文件
    files.GET("/download",
        middleware.RBACMiddleware(model.ResourceFile, model.ActionDownload),
        file.Download,
    )
    
    // 删除文件
    files.DELETE("",
        middleware.RBACMiddleware(model.ResourceFile, model.ActionDelete),
        file.Delete,
    )
    
    // 创建目录
    files.POST("/mkdir",
        middleware.RBACMiddleware(model.ResourceFile, model.ActionCreate),
        file.Mkdir,
    )
}
```

## 4. 权限资源常量映射

```go
// internal/squ-apiserver/model/permission.go

// 资源与路由对应关系（文档参考）
var ResourceRouteMapping = map[string][]string{
    ResourceServer:      {"/servers", "/servers/:id"},
    ResourceApplication: {"/applications", "/applications/:id"},
    ResourceMonitor:     {"/monitor/*"},
    ResourceScript:      {"/scripts", "/scripts/:id"},
    ResourceFile:        {"/servers/:id/files/*"},
    ResourceTerminal:    {"/servers/:id/terminal"},
    ResourceUser:        {"/users", "/users/:id"},
    ResourceRole:        {"/roles", "/roles/:id"},
}
```

## 5. 中间件注册顺序

```go
// internal/squ-apiserver/router/router.go

func Init(r *gin.Engine, conf *config.Config, db database.DB) {
    // 公开路由（无需认证）
    public := r.Group("/api/v1")
    {
        public.POST("/auth/login", auth.Login)
        public.GET("/health", health.Check)
    }
    
    // 需要认证的路由
    api := r.Group("/api/v1")
    api.Use(middleware.JWTMiddleware(conf))  // 先执行 JWT 认证
    {
        // 用户信息（认证后即可访问）
        api.GET("/user/info", user.Info)
        api.GET("/user/permissions", user.Permissions)  // 获取当前用户权限
        
        // 需要 RBAC 权限的路由
        Server(api)
        Application(api)
        Monitor(api)
        Script(api)
        File(api)
        
        // 角色管理（管理员权限）
        Role(api)
        
        // 用户管理（管理员权限）
        User(api)
    }
}
```

## 6. 权限检查失败处理

### 6.1 统一错误响应

```go
// internal/squ-apiserver/middleware/rbac.go

// 权限错误码
const (
    ErrCodeUnauthorized = 401
    ErrCodeForbidden    = 403
)

// 权限检查失败的统一响应
func permissionDenied(c *gin.Context, code int, msg string) {
    // 记录日志
    zap.L().Warn("权限检查失败",
        zap.Int("code", code),
        zap.String("msg", msg),
        zap.String("path", c.Request.URL.Path),
        zap.String("method", c.Request.Method),
        zap.Any("user_id", c.Get("user_id")),
    )
    
    c.AbortWithStatusJSON(code, gin.H{
        "code": code,
        "msg":  msg,
    })
}
```

### 6.2 审计日志

```go
// internal/squ-apiserver/middleware/rbac.go

// RBACMiddlewareWithAudit 带审计日志的权限中间件
func RBACMiddlewareWithAudit(resource, action string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID, _ := c.Get("user_id")
        
        rbacService := rbac.GetRBACService()
        hasPermission, err := rbacService.CheckPermission(userID.(uint), resource, action)
        
        // 记录审计日志
        auditLog := model.AuditLog{
            UserID:     userID.(uint),
            Resource:   resource,
            Action:     action,
            Path:       c.Request.URL.Path,
            Method:     c.Request.Method,
            Allowed:    hasPermission,
            RequestAt:  time.Now(),
        }
        
        if err != nil || !hasPermission {
            auditLog.Reason = "权限不足"
            // 保存审计日志...
            permissionDenied(c, 403, "权限不足")
            return
        }
        
        // 保存审计日志...
        c.Next()
    }
}
```

## 7. 完成检查

- [ ] RBACMiddleware 实现
- [ ] RBACMiddlewareWithServer 实现
- [ ] RBACAnyMiddleware 实现
- [ ] 服务器路由改造
- [ ] 应用路由改造
- [ ] 监控路由改造
- [ ] 脚本路由改造
- [ ] 文件路由改造
- [ ] 中间件注册顺序正确
- [ ] 权限检查失败日志

## 8. 测试要点

1. **基础权限检查**
   - 有权限用户可访问
   - 无权限用户返回 403
   - 未登录用户返回 401

2. **通配符权限**
   - `server:*` 可访问所有服务器操作
   - `*` 可访问所有接口

3. **多角色权限**
   - 用户有多个角色时，权限合并

## 9. 下一阶段

完成本阶段后，继续开发 [04-api.md](./04-api.md) 实现角色管理接口。
