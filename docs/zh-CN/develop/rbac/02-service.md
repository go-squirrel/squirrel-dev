# Phase 2: 权限服务层

## 开发目标

实现 RBAC 核心业务逻辑：
- RBACService 权限检查服务
- 用户角色管理服务
- 权限缓存机制

## 1. 目录结构

```
internal/squ-apiserver/
├── handler/
│   └── rbac/
│       ├── service.go          # 权限检查核心服务
│       ├── role_service.go     # 角色管理服务
│       └── user_role_service.go # 用户角色管理服务
└── repository/
    ├── role.go                 # 已创建
    └── permission.go           # 已创建
```

## 2. 权限检查服务

### 2.1 Service 定义

```go
// internal/squ-apiserver/handler/rbac/service.go

// RBACService 权限检查服务
type RBACService struct {
    roleRepo       *repository.RoleRepository
    permissionRepo *repository.PermissionRepository
    cache          *PermissionCache  // 可选：权限缓存
}

func NewRBACService(
    roleRepo *repository.RoleRepository,
    permissionRepo *repository.PermissionRepository,
) *RBACService {
    return &RBACService{
        roleRepo:       roleRepo,
        permissionRepo: permissionRepo,
    }
}
```

### 2.2 核心权限检查方法

```go
// internal/squ-apiserver/handler/rbac/service.go

// CheckPermission 检查用户是否有指定权限
// resource: 资源类型，如 "server"、"file"
// action: 操作类型，如 "list"、"create"
// 返回: 是否有权限
func (s *RBACService) CheckPermission(userID uint, resource, action string) (bool, error) {
    // 1. 获取用户所有角色
    roles, err := s.roleRepo.GetByUserID(userID)
    if err != nil {
        return false, err
    }
    
    if len(roles) == 0 {
        return false, nil  // 无角色，无权限
    }
    
    // 2. 检查是否有 super_admin 角色
    for _, role := range roles {
        if role.Code == "super_admin" {
            return true, nil  // 超级管理员拥有所有权限
        }
    }
    
    // 3. 构造权限码
    permissionCode := fmt.Sprintf("%s:%s", resource, action)
    wildcardCode := fmt.Sprintf("%s:*", resource)
    
    // 4. 检查每个角色的权限
    for _, role := range roles {
        permCodes, err := s.roleRepo.GetPermissionCodes(role.ID)
        if err != nil {
            continue
        }
        
        for _, code := range permCodes {
            // 匹配：精确匹配 或 通配符匹配
            if code == permissionCode || code == wildcardCode || code == "*" {
                return true, nil
            }
        }
    }
    
    return false, nil
}

// CheckAnyPermission 检查用户是否有任意一个指定权限
func (s *RBACService) CheckAnyPermission(userID uint, permissionCodes []string) (bool, error) {
    for _, code := range permissionCodes {
        parts := strings.Split(code, ":")
        if len(parts) != 2 {
            continue
        }
        has, err := s.CheckPermission(userID, parts[0], parts[1])
        if err != nil {
            return false, err
        }
        if has {
            return true, nil
        }
    }
    return false, nil
}
```

### 2.3 获取用户权限列表

```go
// internal/squ-apiserver/handler/rbac/service.go

// GetUserPermissions 获取用户所有权限码
func (s *RBACService) GetUserPermissions(userID uint) ([]string, error) {
    roles, err := s.roleRepo.GetByUserID(userID)
    if err != nil {
        return nil, err
    }
    
    permissionSet := make(map[string]bool)
    
    for _, role := range roles {
        // super_admin 特殊处理
        if role.Code == "super_admin" {
            return []string{"*"}, nil
        }
        
        permCodes, err := s.roleRepo.GetPermissionCodes(role.ID)
        if err != nil {
            continue
        }
        
        for _, code := range permCodes {
            permissionSet[code] = true
        }
    }
    
    result := make([]string, 0, len(permissionSet))
    for code := range permissionSet {
        result = append(result, code)
    }
    return result, nil
}

// GetUserRoleCodes 获取用户角色编码列表
func (s *RBACService) GetUserRoleCodes(userID uint) ([]string, error) {
    roles, err := s.roleRepo.GetByUserID(userID)
    if err != nil {
        return nil, err
    }
    
    codes := make([]string, len(roles))
    for i, role := range roles {
        codes[i] = role.Code
    }
    return codes, nil
}
```

## 3. 角色管理服务

```go
// internal/squ-apiserver/handler/rbac/role_service.go

// RoleService 角色管理服务
type RoleService struct {
    roleRepo       *repository.RoleRepository
    permissionRepo *repository.PermissionRepository
    db             *gorm.DB
}

// List 获取角色列表
func (s *RoleService) List() ([]model.Role, error) {
    var roles []model.Role
    err := s.db.Where("status = 1").Order("id").Find(&roles).Error
    return roles, err
}

// GetByID 获取角色详情
func (s *RoleService) GetByID(id uint) (*model.Role, error) {
    var role model.Role
    err := s.db.First(&role, id).Error
    return &role, err
}

// GetRolePermissions 获取角色的权限列表
func (s *RoleService) GetRolePermissions(roleID uint) ([]model.Permission, error) {
    var permissions []model.Permission
    err := s.db.Table("permissions").
        Select("permissions.*").
        Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
        Where("role_permissions.role_id = ?", roleID).
        Find(&permissions).Error
    return permissions, err
}

// UpdateRolePermissions 更新角色权限
func (s *RoleService) UpdateRolePermissions(roleID uint, permissionIDs []uint) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // 1. 删除旧关联
        if err := tx.Where("role_id = ?", roleID).Delete(&model.RolePermission{}).Error; err != nil {
            return err
        }
        
        // 2. 创建新关联
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

## 4. 用户角色管理服务

```go
// internal/squ-apiserver/handler/rbac/user_role_service.go

// UserRoleService 用户角色管理服务
type UserRoleService struct {
    db       *gorm.DB
    roleRepo *repository.RoleRepository
}

// GetUserRoles 获取用户的角色列表
func (s *UserRoleService) GetUserRoles(userID uint) ([]model.Role, error) {
    return s.roleRepo.GetByUserID(userID)
}

// UpdateUserRoles 更新用户角色
func (s *UserRoleService) UpdateUserRoles(userID uint, roleIDs []uint) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // 1. 删除旧关联
        if err := tx.Where("user_id = ?", userID).Delete(&model.UserRole{}).Error; err != nil {
            return err
        }
        
        // 2. 创建新关联
        for _, roleID := range roleIDs {
            ur := model.UserRole{UserID: userID, RoleID: roleID}
            if err := tx.Create(&ur).Error; err != nil {
                return err
            }
        }
        
        return nil
    })
}

// AssignRole 为用户分配角色（追加）
func (s *UserRoleService) AssignRole(userID, roleID uint) error {
    ur := model.UserRole{UserID: userID, RoleID: roleID}
    return s.db.Create(&ur).Error
}

// RemoveRole 移除用户角色
func (s *UserRoleService) RemoveRole(userID, roleID uint) error {
    return s.db.Where("user_id = ? AND role_id = ?", userID, roleID).
        Delete(&model.UserRole{}).Error
}
```

## 5. 权限缓存（可选优化）

```go
// internal/squ-apiserver/handler/rbac/cache.go

// PermissionCache 权限缓存
type PermissionCache struct {
    store sync.Map  // key: userID, value: *UserPermissionCache
    ttl   time.Duration
}

type UserPermissionCache struct {
    Permissions []string
    Roles       []string
    ExpiredAt   time.Time
}

func NewPermissionCache(ttl time.Duration) *PermissionCache {
    return &PermissionCache{ttl: ttl}
}

// Get 从缓存获取用户权限
func (c *PermissionCache) Get(userID uint) (*UserPermissionCache, bool) {
    val, ok := c.store.Load(userID)
    if !ok {
        return nil, false
    }
    
    cache := val.(*UserPermissionCache)
    if time.Now().After(cache.ExpiredAt) {
        c.store.Delete(userID)
        return nil, false
    }
    
    return cache, true
}

// Set 设置用户权限缓存
func (c *PermissionCache) Set(userID uint, permissions, roles []string) {
    cache := &UserPermissionCache{
        Permissions: permissions,
        Roles:       roles,
        ExpiredAt:   time.Now().Add(c.ttl),
    }
    c.store.Store(userID, cache)
}

// Delete 删除用户权限缓存
func (c *PermissionCache) Delete(userID uint) {
    c.store.Delete(userID)
}

// Clear 清空所有缓存
func (c *PermissionCache) Clear() {
    c.store.Range(func(key, value interface{}) bool {
        c.store.Delete(key)
        return true
    })
}
```

### 5.1 带缓存的权限检查

```go
// internal/squ-apiserver/handler/rbac/service.go

// CheckPermissionWithCache 带缓存的权限检查
func (s *RBACService) CheckPermissionWithCache(userID uint, resource, action string) (bool, error) {
    // 1. 尝试从缓存获取
    if s.cache != nil {
        if cache, ok := s.cache.Get(userID); ok {
            // 从缓存数据检查权限
            return s.checkPermissionFromCache(cache, resource, action), nil
        }
    }
    
    // 2. 缓存未命中，从数据库获取
    permissions, err := s.GetUserPermissions(userID)
    if err != nil {
        return false, err
    }
    
    roles, err := s.GetUserRoleCodes(userID)
    if err != nil {
        return false, err
    }
    
    // 3. 设置缓存
    if s.cache != nil {
        s.cache.Set(userID, permissions, roles)
    }
    
    // 4. 检查权限
    return s.checkPermissionFromCodes(permissions, resource, action), nil
}

func (s *RBACService) checkPermissionFromCodes(codes []string, resource, action string) bool {
    permissionCode := fmt.Sprintf("%s:%s", resource, action)
    wildcardCode := fmt.Sprintf("%s:*", resource)
    
    for _, code := range codes {
        if code == "*" || code == permissionCode || code == wildcardCode {
            return true
        }
    }
    return false
}
```

## 6. 服务初始化

```go
// internal/squ-apiserver/handler/rbac/init.go

var (
    rbacService      *RBACService
    roleService      *RoleService
    userRoleService  *UserRoleService
    permissionCache  *PermissionCache
)

// Init 初始化 RBAC 服务
func Init(db *gorm.DB) {
    // 初始化 Repository
    roleRepo := repository.NewRoleRepository(db)
    permissionRepo := repository.NewPermissionRepository(db)
    
    // 初始化缓存（可选，TTL 5分钟）
    permissionCache = NewPermissionCache(5 * time.Minute)
    
    // 初始化服务
    rbacService = NewRBACService(roleRepo, permissionRepo)
    rbacService.cache = permissionCache  // 注入缓存
    
    roleService = NewRoleService(roleRepo, permissionRepo, db)
    userRoleService = NewUserRoleService(db, roleRepo)
}

// GetRBACService 获取 RBAC 服务实例
func GetRBACService() *RBACService {
    return rbacService
}

// GetRoleService 获取角色服务实例
func GetRoleService() *RoleService {
    return roleService
}

// GetUserRoleService 获取用户角色服务实例
func GetUserRoleService() *UserRoleService {
    return userRoleService
}

// ClearUserPermissionCache 清除用户权限缓存
// 当用户角色变更时调用
func ClearUserPermissionCache(userID uint) {
    if permissionCache != nil {
        permissionCache.Delete(userID)
    }
}
```

## 7. 完成检查

- [ ] RBACService 结构体定义
- [ ] CheckPermission 方法实现
- [ ] CheckAnyPermission 方法实现
- [ ] GetUserPermissions 方法实现
- [ ] RoleService 实现
- [ ] UserRoleService 实现
- [ ] PermissionCache 实现（可选）
- [ ] 服务初始化函数

## 8. 测试要点

1. **权限检查测试**
   - super_admin 拥有所有权限
   - 通配符权限匹配
   - 多角色权限合并
   - 无角色用户无权限

2. **缓存测试**
   - 缓存命中
   - 缓存过期
   - 缓存清除

## 9. 下一阶段

完成本阶段后，继续开发 [03-middleware.md](./03-middleware.md) 实现权限中间件。
