# Phase 4: 管理接口

## 开发目标

实现 RBAC 管理相关的 REST API：
- 角色管理接口
- 用户角色分配接口
- 权限查询接口

## 1. 目录结构

```
internal/squ-apiserver/
├── handler/
│   ├── role/
│   │   ├── handler.go       # 路由入口
│   │   ├── service.go       # 业务逻辑
│   │   └── types.go         # 请求/响应类型
│   └── user/
│       ├── handler.go       # 已存在，扩展角色管理
│       └── service.go
└── router/
    ├── role.go              # 新增
    └── user.go              # 修改
```

## 2. 请求/响应类型

```go
// internal/squ-apiserver/handler/role/types.go

// ========== 角色 ==========

// RoleResponse 角色信息响应
type RoleResponse struct {
    ID          uint   `json:"id"`
    Name        string `json:"name"`
    Code        string `json:"code"`
    Description string `json:"description"`
    IsSystem    bool   `json:"is_system"`
    Status      int    `json:"status"`
}

// RoleDetailResponse 角色详情响应（含权限）
type RoleDetailResponse struct {
    RoleResponse
    Permissions []PermissionResponse `json:"permissions"`
}

// CreateRoleRequest 创建角色请求
type CreateRoleRequest struct {
    Name          string `json:"name" binding:"required"`
    Code          string `json:"code" binding:"required"`
    Description   string `json:"description"`
}

// UpdateRoleRequest 更新角色请求
type UpdateRoleRequest struct {
    Name        string `json:"name"`
    Description string `json:"description"`
    Status      *int   `json:"status"`
}

// UpdateRolePermissionsRequest 更新角色权限请求
type UpdateRolePermissionsRequest struct {
    PermissionIDs []uint `json:"permission_ids" binding:"required"`
}

// ========== 权限 ==========

// PermissionResponse 权限信息
type PermissionResponse struct {
    ID          uint   `json:"id"`
    Name        string `json:"name"`
    Code        string `json:"code"`
    Resource    string `json:"resource"`
    Action      string `json:"action"`
    Description string `json:"description"`
}

// PermissionGroupResponse 权限分组响应
type PermissionGroupResponse struct {
    Resource    string                `json:"resource"`
    Permissions []PermissionResponse `json:"permissions"`
}

// ========== 用户角色 ==========

// UserRoleResponse 用户角色响应
type UserRoleResponse struct {
    UserID uint            `json:"user_id"`
    Roles  []RoleResponse `json:"roles"`
}

// UpdateUserRolesRequest 更新用户角色请求
type UpdateUserRolesRequest struct {
    RoleIDs []uint `json:"role_ids" binding:"required"`
}

// ========== 当前用户 ==========

// UserPermissionsResponse 当前用户权限响应
type UserPermissionsResponse struct {
    Roles       []string `json:"roles"`
    Permissions []string `json:"permissions"`
}
```

## 3. 角色管理接口

### 3.1 Handler

```go
// internal/squ-apiserver/handler/role/handler.go

var svc *Service

func Init(db *gorm.DB) {
    svc = NewService(db)
}

// List 获取角色列表
// GET /api/v1/roles
func List(c *gin.Context) {
    roles, err := svc.List()
    if err != nil {
        response.Error(c, 500, "获取角色列表失败")
        return
    }
    
    list := make([]RoleResponse, len(roles))
    for i, role := range roles {
        list[i] = RoleResponse{
            ID:          role.ID,
            Name:        role.Name,
            Code:        role.Code,
            Description: role.Description,
            IsSystem:    role.IsSystem,
            Status:      role.Status,
        }
    }
    
    response.Success(c, list)
}

// Get 获取角色详情
// GET /api/v1/roles/:id
func Get(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        response.Error(c, 400, "无效的角色ID")
        return
    }
    
    detail, err := svc.GetDetail(uint(id))
    if err != nil {
        response.Error(c, 404, "角色不存在")
        return
    }
    
    response.Success(c, detail)
}

// Create 创建角色
// POST /api/v1/roles
func Create(c *gin.Context) {
    var req CreateRoleRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, 400, "参数错误: "+err.Error())
        return
    }
    
    role, err := svc.Create(req)
    if err != nil {
        response.Error(c, 500, "创建角色失败: "+err.Error())
        return
    }
    
    response.Success(c, RoleResponse{
        ID:          role.ID,
        Name:        role.Name,
        Code:        role.Code,
        Description: role.Description,
        IsSystem:    role.IsSystem,
        Status:      role.Status,
    })
}

// Update 更新角色
// PUT /api/v1/roles/:id
func Update(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    
    var req UpdateRoleRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, 400, "参数错误")
        return
    }
    
    err := svc.Update(uint(id), req)
    if err != nil {
        response.Error(c, 500, "更新角色失败: "+err.Error())
        return
    }
    
    response.Success(c, nil)
}

// Delete 删除角色
// DELETE /api/v1/roles/:id
func Delete(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    
    err := svc.Delete(uint(id))
    if err != nil {
        response.Error(c, 500, "删除角色失败: "+err.Error())
        return
    }
    
    response.Success(c, nil)
}

// GetPermissions 获取角色权限
// GET /api/v1/roles/:id/permissions
func GetPermissions(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    
    permissions, err := svc.GetPermissions(uint(id))
    if err != nil {
        response.Error(c, 500, "获取权限失败")
        return
    }
    
    response.Success(c, permissions)
}

// UpdatePermissions 更新角色权限
// PUT /api/v1/roles/:id/permissions
func UpdatePermissions(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    
    var req UpdateRolePermissionsRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, 400, "参数错误")
        return
    }
    
    err := svc.UpdatePermissions(uint(id), req.PermissionIDs)
    if err != nil {
        response.Error(c, 500, "更新权限失败: "+err.Error())
        return
    }
    
    response.Success(c, nil)
}
```

### 3.2 Service

```go
// internal/squ-apiserver/handler/role/service.go

type Service struct {
    db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
    return &Service{db: db}
}

func (s *Service) List() ([]model.Role, error) {
    var roles []model.Role
    err := s.db.Where("status = 1").Order("id").Find(&roles).Error
    return roles, err
}

func (s *Service) GetDetail(id uint) (*RoleDetailResponse, error) {
    var role model.Role
    err := s.db.First(&role, id).Error
    if err != nil {
        return nil, err
    }
    
    // 获取权限
    var permissions []model.Permission
    s.db.Table("permissions").
        Select("permissions.*").
        Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
        Where("role_permissions.role_id = ?", id).
        Find(&permissions)
    
    permList := make([]PermissionResponse, len(permissions))
    for i, p := range permissions {
        permList[i] = PermissionResponse{
            ID:          p.ID,
            Name:        p.Name,
            Code:        p.Code,
            Resource:    p.Resource,
            Action:      p.Action,
            Description: p.Description,
        }
    }
    
    return &RoleDetailResponse{
        RoleResponse: RoleResponse{
            ID:          role.ID,
            Name:        role.Name,
            Code:        role.Code,
            Description: role.Description,
            IsSystem:    role.IsSystem,
            Status:      role.Status,
        },
        Permissions: permList,
    }, nil
}

func (s *Service) Create(req CreateRoleRequest) (*model.Role, error) {
    role := model.Role{
        Name:        req.Name,
        Code:        req.Code,
        Description: req.Description,
        IsSystem:    false,
        Status:      1,
    }
    err := s.db.Create(&role).Error
    return &role, err
}

func (s *Service) Update(id uint, req UpdateRoleRequest) error {
    updates := map[string]interface{}{}
    if req.Name != "" {
        updates["name"] = req.Name
    }
    if req.Description != "" {
        updates["description"] = req.Description
    }
    if req.Status != nil {
        updates["status"] = *req.Status
    }
    
    return s.db.Model(&model.Role{}).Where("id = ?", id).Updates(updates).Error
}

func (s *Service) Delete(id uint) error {
    // 检查是否为系统角色
    var role model.Role
    if err := s.db.First(&role, id).Error; err != nil {
        return err
    }
    if role.IsSystem {
        return fmt.Errorf("系统内置角色不能删除")
    }
    
    // 删除角色关联
    return s.db.Transaction(func(tx *gorm.DB) error {
        // 删除用户-角色关联
        tx.Where("role_id = ?", id).Delete(&model.UserRole{})
        // 删除角色-权限关联
        tx.Where("role_id = ?", id).Delete(&model.RolePermission{})
        // 删除角色
        return tx.Delete(&model.Role{}, id).Error
    })
}

func (s *Service) GetPermissions(roleID uint) ([]PermissionResponse, error) {
    var permissions []model.Permission
    err := s.db.Table("permissions").
        Select("permissions.*").
        Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
        Where("role_permissions.role_id = ?", roleID).
        Find(&permissions).Error
    
    if err != nil {
        return nil, err
    }
    
    result := make([]PermissionResponse, len(permissions))
    for i, p := range permissions {
        result[i] = PermissionResponse{
            ID:          p.ID,
            Name:        p.Name,
            Code:        p.Code,
            Resource:    p.Resource,
            Action:      p.Action,
            Description: p.Description,
        }
    }
    return result, nil
}

func (s *Service) UpdatePermissions(roleID uint, permissionIDs []uint) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // 删除旧关联
        tx.Where("role_id = ?", roleID).Delete(&model.RolePermission{})
        // 创建新关联
        for _, permID := range permissionIDs {
            rp := model.RolePermission{RoleID: roleID, PermissionID: permID}
            if err := tx.Create(&rp).Error; err != nil {
                return err
            }
        }
        return nil
    })
}
```

## 4. 用户角色管理接口

```go
// internal/squ-apiserver/handler/user/handler.go

// GetUserRoles 获取用户角色
// GET /api/v1/users/:id/roles
func GetUserRoles(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    
    roles, err := userRoleService.GetUserRoles(uint(id))
    if err != nil {
        response.Error(c, 500, "获取用户角色失败")
        return
    }
    
    list := make([]RoleResponse, len(roles))
    for i, role := range roles {
        list[i] = RoleResponse{
            ID:          role.ID,
            Name:        role.Name,
            Code:        role.Code,
            Description: role.Description,
            IsSystem:    role.IsSystem,
            Status:      role.Status,
        }
    }
    
    response.Success(c, UserRoleResponse{
        UserID: uint(id),
        Roles:  list,
    })
}

// UpdateUserRoles 更新用户角色
// PUT /api/v1/users/:id/roles
func UpdateUserRoles(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    
    var req UpdateUserRolesRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, 400, "参数错误")
        return
    }
    
    err := userRoleService.UpdateUserRoles(uint(id), req.RoleIDs)
    if err != nil {
        response.Error(c, 500, "更新用户角色失败: "+err.Error())
        return
    }
    
    // 清除用户权限缓存
    rbac.ClearUserPermissionCache(uint(id))
    
    response.Success(c, nil)
}
```

## 5. 权限查询接口

```go
// internal/squ-apiserver/handler/user/handler.go

// GetCurrentUserPermissions 获取当前用户权限
// GET /api/v1/user/permissions
func GetCurrentUserPermissions(c *gin.Context) {
    userID, _ := c.Get("user_id")
    
    rbacService := rbac.GetRBACService()
    
    roles, err := rbacService.GetUserRoleCodes(userID.(uint))
    if err != nil {
        response.Error(c, 500, "获取用户角色失败")
        return
    }
    
    permissions, err := rbacService.GetUserPermissions(userID.(uint))
    if err != nil {
        response.Error(c, 500, "获取用户权限失败")
        return
    }
    
    response.Success(c, UserPermissionsResponse{
        Roles:       roles,
        Permissions: permissions,
    })
}

// GetAllPermissions 获取所有权限列表（用于前端权限配置）
// GET /api/v1/permissions
func GetAllPermissions(c *gin.Context) {
    permissionRepo := repository.NewPermissionRepository(database.GetDB())
    
    grouped, err := permissionRepo.GetAllGrouped()
    if err != nil {
        response.Error(c, 500, "获取权限列表失败")
        return
    }
    
    result := make([]PermissionGroupResponse, 0)
    for resource, perms := range grouped {
        permList := make([]PermissionResponse, len(perms))
        for i, p := range perms {
            permList[i] = PermissionResponse{
                ID:          p.ID,
                Name:        p.Name,
                Code:        p.Code,
                Resource:    p.Resource,
                Action:      p.Action,
                Description: p.Description,
            }
        }
        result = append(result, PermissionGroupResponse{
            Resource:    resource,
            Permissions: permList,
        })
    }
    
    response.Success(c, result)
}
```

## 6. 路由注册

```go
// internal/squ-apiserver/router/role.go

func Role(group *gin.RouterGroup) {
    roles := group.Group("/roles")
    
    // 需要 role:list 权限
    roles.GET("",
        middleware.RBACMiddleware(model.ResourceRole, model.ActionList),
        role.List,
    )
    
    // 需要 role:get 权限
    roles.GET("/:id",
        middleware.RBACMiddleware(model.ResourceRole, model.ActionGet),
        role.Get,
    )
    
    // 需要管理员权限（使用 user:update 权限）
    roles.POST("",
        middleware.RBACMiddleware(model.ResourceUser, model.ActionUpdate),
        role.Create,
    )
    
    roles.PUT("/:id",
        middleware.RBACMiddleware(model.ResourceUser, model.ActionUpdate),
        role.Update,
    )
    
    roles.DELETE("/:id",
        middleware.RBACMiddleware(model.ResourceUser, model.ActionUpdate),
        role.Delete,
    )
    
    roles.GET("/:id/permissions",
        middleware.RBACMiddleware(model.ResourceRole, model.ActionGet),
        role.GetPermissions,
    )
    
    roles.PUT("/:id/permissions",
        middleware.RBACMiddleware(model.ResourceUser, model.ActionUpdate),
        role.UpdatePermissions,
    )
}

// internal/squ-apiserver/router/user.go

func User(group *gin.RouterGroup) {
    users := group.Group("/users")
    
    // 用户管理
    users.GET("",
        middleware.RBACMiddleware(model.ResourceUser, model.ActionList),
        user.List,
    )
    users.GET("/:id",
        middleware.RBACMiddleware(model.ResourceUser, model.ActionGet),
        user.Get,
    )
    users.POST("",
        middleware.RBACMiddleware(model.ResourceUser, model.ActionCreate),
        user.Create,
    )
    users.PUT("/:id",
        middleware.RBACMiddleware(model.ResourceUser, model.ActionUpdate),
        user.Update,
    )
    users.DELETE("/:id",
        middleware.RBACMiddleware(model.ResourceUser, model.ActionDelete),
        user.Delete,
    )
    
    // 用户角色
    users.GET("/:id/roles",
        middleware.RBACMiddleware(model.ResourceUser, model.ActionGet),
        user.GetUserRoles,
    )
    users.PUT("/:id/roles",
        middleware.RBACMiddleware(model.ResourceUser, model.ActionUpdate),
        user.UpdateUserRoles,
    )
    
    // 权限列表
    group.GET("/permissions",
        middleware.RBACMiddleware(model.ResourceRole, model.ActionList),
        user.GetAllPermissions,
    )
}

// 在 router.go 中添加
func Init(r *gin.Engine, conf *config.Config, db database.DB) {
    // ... 现有代码 ...
    
    // 当前用户权限（认证后即可访问）
    api.GET("/user/permissions", user.GetCurrentUserPermissions)
    
    // 角色管理
    Role(api)
    
    // 用户管理（已存在，需修改）
    User(api)
}
```

## 7. API 接口汇总

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| GET | /api/v1/roles | 角色列表 | role:list |
| GET | /api/v1/roles/:id | 角色详情 | role:get |
| POST | /api/v1/roles | 创建角色 | user:update |
| PUT | /api/v1/roles/:id | 更新角色 | user:update |
| DELETE | /api/v1/roles/:id | 删除角色 | user:update |
| GET | /api/v1/roles/:id/permissions | 角色权限 | role:get |
| PUT | /api/v1/roles/:id/permissions | 更新权限 | user:update |
| GET | /api/v1/users/:id/roles | 用户角色 | user:get |
| PUT | /api/v1/users/:id/roles | 更新角色 | user:update |
| GET | /api/v1/user/permissions | 当前用户权限 | 登录即可 |
| GET | /api/v1/permissions | 所有权限列表 | role:list |

## 8. 完成检查

- [ ] 角色管理 Handler
- [ ] 角色管理 Service
- [ ] 用户角色管理 Handler
- [ ] 权限查询 Handler
- [ ] 路由注册
- [ ] 权限中间件配置

## 9. 测试要点

1. **角色管理**
   - 创建/更新/删除角色
   - 系统角色不能删除
   - 角色名称/编码唯一

2. **用户角色**
   - 分配角色成功
   - 清除缓存生效

3. **权限查询**
   - 当前用户权限正确返回
   - 权限分组正确

## 10. 下一阶段

完成本阶段后，继续开发 [05-frontend.md](./05-frontend.md) 实现前端权限集成。
