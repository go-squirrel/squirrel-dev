# Phase 1: 数据模型

## 开发目标

定义 RBAC 核心数据模型，包括：
- Role（角色）
- Permission（权限）
- 关联表（UserRole、RolePermission）
- 数据库迁移和初始化数据

## 1. 目录结构

```
internal/squ-apiserver/
├── model/
│   ├── user.go          # 已存在，扩展字段
│   ├── role.go          # 新增：角色模型
│   ├── permission.go    # 新增：权限模型
│   └── rbac.go          # 新增：关联表模型
└── repository/
    ├── role.go          # 新增：角色仓储
    └── permission.go    # 新增：权限仓储
```

## 2. 数据模型定义

### 2.1 扩展 User 模型

```go
// internal/squ-apiserver/model/user.go

type User struct {
    BaseModel
    Username string  `gorm:"size:50;not null;unique"`
    Password string  `gorm:"size:100;not null"`
    Email    string  `gorm:"size:100;unique"`
    Nickname string  `gorm:"size:50"`
    Avatar   string  `gorm:"size:255"`
    Status   int     `gorm:"default:1"`  // 1: 正常, 0: 禁用
    
    // 新增字段
    LastLoginAt *time.Time `gorm:"comment:最后登录时间"`
}

// 新增：用户角色列表（GORM 关联）
func (u *User) Roles() []Role {
    // 通过 UserRole 关联表获取
}
```

### 2.2 Role 模型

```go
// internal/squ-apiserver/model/role.go

// Role 角色
type Role struct {
    BaseModel
    Name        string `gorm:"size:50;not null;unique"`
    Code        string `gorm:"size:50;not null;unique;index"`  // 角色编码
    Description string `gorm:"size:255"`
    IsSystem    bool   `gorm:"default:false;comment:是否系统内置"`
    Status      int    `gorm:"default:1"`  // 1: 启用, 0: 禁用
}

// 角色编码常量
const (
    RoleSuperAdmin = "super_admin"
    RoleAdmin      = "admin"
    RoleOperator   = "operator"
    RoleDeveloper  = "developer"
    RoleViewer     = "viewer"
)
```

### 2.3 Permission 模型

```go
// internal/squ-apiserver/model/permission.go

// Permission 权限
type Permission struct {
    BaseModel
    Name        string `gorm:"size:100;not null"`
    Code        string `gorm:"size:100;not null;unique;index"`
    Resource    string `gorm:"size:50;not null;index"`
    Action      string `gorm:"size:20;not null"`
    Description string `gorm:"size:255"`
}

// 资源类型常量
const (
    ResourceServer      = "server"
    ResourceApplication = "application"
    ResourceMonitor     = "monitor"
    ResourceScript      = "script"
    ResourceFile        = "file"
    ResourceTerminal    = "terminal"
    ResourceUser        = "user"
    ResourceRole        = "role"
)

// 操作类型常量
const (
    ActionList     = "list"
    ActionGet      = "get"
    ActionCreate   = "create"
    ActionUpdate   = "update"
    ActionDelete   = "delete"
    ActionExec     = "exec"
    ActionUpload   = "upload"
    ActionDownload = "download"
)
```

### 2.4 关联表模型

```go
// internal/squ-apiserver/model/rbac.go

// UserRole 用户-角色关联
type UserRole struct {
    UserID    uint      `gorm:"primaryKey;autoIncrement:false"`
    RoleID    uint      `gorm:"primaryKey;autoIncrement:false"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (UserRole) TableName() string {
    return "user_roles"
}

// RolePermission 角色-权限关联
type RolePermission struct {
    RoleID       uint      `gorm:"primaryKey;autoIncrement:false"`
    PermissionID uint      `gorm:"primaryKey;autoIncrement:false"`
    CreatedAt    time.Time `gorm:"autoCreateTime"`
}

func (RolePermission) TableName() string {
    return "role_permissions"
}
```

## 3. 初始化数据

### 3.1 预置权限

```go
// internal/squ-apiserver/model/permission_init.go

// DefaultPermissions 预置权限列表
var DefaultPermissions = []Permission{
    // 服务器管理
    {Name: "查看服务器列表", Code: "server:list", Resource: "server", Action: "list"},
    {Name: "查看服务器详情", Code: "server:get", Resource: "server", Action: "get"},
    {Name: "添加服务器", Code: "server:create", Resource: "server", Action: "create"},
    {Name: "编辑服务器", Code: "server:update", Resource: "server", Action: "update"},
    {Name: "删除服务器", Code: "server:delete", Resource: "server", Action: "delete"},
    
    // 应用管理
    {Name: "查看应用列表", Code: "application:list", Resource: "application", Action: "list"},
    {Name: "查看应用详情", Code: "application:get", Resource: "application", Action: "get"},
    {Name: "创建应用", Code: "application:create", Resource: "application", Action: "create"},
    {Name: "编辑应用", Code: "application:update", Resource: "application", Action: "update"},
    {Name: "删除应用", Code: "application:delete", Resource: "application", Action: "delete"},
    {Name: "部署应用", Code: "application:deploy", Resource: "application", Action: "exec"},
    
    // 监控
    {Name: "查看监控数据", Code: "monitor:get", Resource: "monitor", Action: "get"},
    {Name: "查看监控历史", Code: "monitor:list", Resource: "monitor", Action: "list"},
    
    // 脚本管理
    {Name: "查看脚本列表", Code: "script:list", Resource: "script", Action: "list"},
    {Name: "查看脚本详情", Code: "script:get", Resource: "script", Action: "get"},
    {Name: "创建脚本", Code: "script:create", Resource: "script", Action: "create"},
    {Name: "编辑脚本", Code: "script:update", Resource: "script", Action: "update"},
    {Name: "删除脚本", Code: "script:delete", Resource: "script", Action: "delete"},
    {Name: "执行脚本", Code: "script:exec", Resource: "script", Action: "exec"},
    
    // 文件管理
    {Name: "浏览文件", Code: "file:list", Resource: "file", Action: "list"},
    {Name: "查看文件内容", Code: "file:get", Resource: "file", Action: "get"},
    {Name: "编辑文件", Code: "file:update", Resource: "file", Action: "update"},
    {Name: "删除文件", Code: "file:delete", Resource: "file", Action: "delete"},
    {Name: "上传文件", Code: "file:upload", Resource: "file", Action: "upload"},
    {Name: "下载文件", Code: "file:download", Resource: "file", Action: "download"},
    
    // Web终端
    {Name: "使用Web终端", Code: "terminal:exec", Resource: "terminal", Action: "exec"},
    
    // 用户管理
    {Name: "查看用户列表", Code: "user:list", Resource: "user", Action: "list"},
    {Name: "查看用户详情", Code: "user:get", Resource: "user", Action: "get"},
    {Name: "创建用户", Code: "user:create", Resource: "user", Action: "create"},
    {Name: "编辑用户", Code: "user:update", Resource: "user", Action: "update"},
    {Name: "删除用户", Code: "user:delete", Resource: "user", Action: "delete"},
    
    // 角色管理
    {Name: "查看角色列表", Code: "role:list", Resource: "role", Action: "list"},
    {Name: "查看角色详情", Code: "role:get", Resource: "role", Action: "get"},
}
```

### 3.2 预置角色

```go
// internal/squ-apiserver/model/role_init.go

// DefaultRoles 预置角色
var DefaultRoles = []Role{
    {
        Name:        "超级管理员",
        Code:        "super_admin",
        Description: "系统最高权限，拥有所有功能访问权限",
        IsSystem:    true,
        Status:      1,
    },
    {
        Name:        "运维管理员",
        Code:        "admin",
        Description: "拥有大部分管理权限，无法管理用户",
        IsSystem:    true,
        Status:      1,
    },
    {
        Name:        "运维工程师",
        Code:        "operator",
        Description: "日常运维操作，可管理服务器、应用、脚本、文件",
        IsSystem:    true,
        Status:      1,
    },
    {
        Name:        "开发人员",
        Code:        "developer",
        Description: "可查看服务器、使用终端、部署应用",
        IsSystem:    true,
        Status:      1,
    },
    {
        Name:        "只读用户",
        Code:        "viewer",
        Description: "仅可查看服务器状态和监控数据",
        IsSystem:    true,
        Status:      1,
    },
}
```

### 3.3 角色权限映射

```go
// internal/squ-apiserver/model/rbac_init.go

// RolePermissionMapping 角色默认权限映射
// key: 角色Code, value: 权限Code列表
var RolePermissionMapping = map[string][]string{
    "super_admin": {"*"},  // 特殊处理：所有权限
    
    "admin": {
        "server:*", "application:*", "monitor:*",
        "script:*", "file:*", "terminal:exec",
        "role:list", "role:get",
    },
    
    "operator": {
        "server:list", "server:get",
        "application:*", "monitor:*",
        "script:*", "file:*", "terminal:exec",
    },
    
    "developer": {
        "server:list", "server:get",
        "application:list", "application:get", "application:deploy",
        "monitor:*", "terminal:exec",
        "file:list", "file:get", "file:download",
    },
    
    "viewer": {
        "server:list", "server:get",
        "application:list", "application:get",
        "monitor:*", "script:list",
    },
}
```

## 4. 数据库迁移

### 4.1 迁移文件

```go
// internal/pkg/migration/migrations/xxxx_add_rbac_tables.go

// Up 创建 RBAC 相关表
func (m *Migration) Up() error {
    // 创建 roles 表
    // 创建 permissions 表
    // 创建 user_roles 关联表
    // 创建 role_permissions 关联表
    // 为 users 表添加 last_login_at 字段
    return nil
}

// Down 回滚
func (m *Migration) Down() error {
    // 删除关联表
    // 删除 permissions 表
    // 删除 roles 表
    // 移除 users 表新增字段
    return nil
}
```

### 4.2 初始化数据迁移

```go
// internal/pkg/migration/migrations/xxxx_init_rbac_data.go

func (m *Migration) Up() error {
    // 1. 插入预置权限
    for _, perm := range DefaultPermissions {
        db.FirstOrCreate(&perm, Permission{Code: perm.Code})
    }
    
    // 2. 插入预置角色
    for _, role := range DefaultRoles {
        db.FirstOrCreate(&role, Role{Code: role.Code})
    }
    
    // 3. 建立角色-权限关联
    for roleCode, permCodes := range RolePermissionMapping {
        var role Role
        db.Where("code = ?", roleCode).First(&role)
        
        for _, permCode := range permCodes {
            if permCode == "*" {
                continue  // super_admin 特殊处理
            }
            var perm Permission
            db.Where("code = ?", permCode).First(&perm)
            db.Create(&RolePermission{RoleID: role.ID, PermissionID: perm.ID})
        }
    }
    
    // 4. 为现有管理员用户分配 super_admin 角色
    //    根据实际需求确定如何识别管理员
    
    return nil
}
```

## 5. Repository 层

### 5.1 Role Repository

```go
// internal/squ-apiserver/repository/role.go

type RoleRepository struct {
    db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
    return &RoleRepository{db: db}
}

// 获取用户的所有角色
func (r *RoleRepository) GetByUserID(userID uint) ([]model.Role, error) {
    var roles []model.Role
    err := r.db.Table("roles").
        Select("roles.*").
        Joins("JOIN user_roles ON user_roles.role_id = roles.id").
        Where("user_roles.user_id = ? AND roles.status = 1", userID).
        Find(&roles).Error
    return roles, err
}

// 获取角色的所有权限 Code
func (r *RoleRepository) GetPermissionCodes(roleID uint) ([]string, error) {
    var codes []string
    err := r.db.Table("permissions").
        Select("permissions.code").
        Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
        Where("role_permissions.role_id = ?", roleID).
        Pluck("code", &codes).Error
    return codes, err
}
```

### 5.2 Permission Repository

```go
// internal/squ-apiserver/repository/permission.go

type PermissionRepository struct {
    db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
    return &PermissionRepository{db: db}
}

// 获取所有权限（按资源分组）
func (r *PermissionRepository) GetAllGrouped() (map[string][]model.Permission, error) {
    var permissions []model.Permission
    err := r.db.Find(&permissions).Error
    if err != nil {
        return nil, err
    }
    
    result := make(map[string][]model.Permission)
    for _, p := range permissions {
        result[p.Resource] = append(result[p.Resource], p)
    }
    return result, nil
}
```

## 6. 完成检查

- [ ] Role 模型定义
- [ ] Permission 模型定义
- [ ] UserRole 关联表定义
- [ ] RolePermission 关联表定义
- [ ] 预置权限数据
- [ ] 预置角色数据
- [ ] 角色权限映射数据
- [ ] 数据库迁移文件
- [ ] 初始化数据迁移
- [ ] RoleRepository 实现
- [ ] PermissionRepository 实现

## 7. 测试要点

1. **模型测试**
   - 角色创建、查询、更新、删除
   - 权限创建、查询
   - 关联关系正确建立

2. **迁移测试**
   - 迁移执行成功
   - 回滚执行成功
   - 初始化数据正确

## 8. 下一阶段

完成本阶段后，继续开发 [02-service.md](./02-service.md) 实现权限服务层。
