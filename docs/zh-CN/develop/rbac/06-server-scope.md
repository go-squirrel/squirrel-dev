# Phase 6: 服务器级权限隔离（可选）

## 开发目标

实现更精细的服务器级权限控制：
- 用户只能访问被授权的服务器
- 支持按服务器分配角色
- 前端服务器列表过滤

## 1. 使用场景

| 场景 | 说明 |
|------|------|
| 运维分区 | A 组运维只能管理测试环境服务器，B 组管理生产环境 |
| 客户隔离 | 多租户场景下，不同客户只能访问自己的服务器 |
| 权限委托 | 某用户只在特定服务器上拥有管理员权限 |

## 2. 数据模型扩展

### 2.1 用户服务器范围模型

```go
// internal/squ-apiserver/model/rbac.go

// UserServerScope 用户服务器权限范围
type UserServerScope struct {
    UserID   uint      `gorm:"primaryKey;autoIncrement:false;index"`
    ServerID uint      `gorm:"primaryKey;autoIncrement:false;index"`
    RoleID   uint      `gorm:"primaryKey;autoIncrement:false;index"`  // 该服务器上的角色
    CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (UserServerScope) TableName() string {
    return "user_server_scopes"
}
```

### 2.2 数据库迁移

```go
// internal/pkg/migration/migrations/xxxx_add_user_server_scope.go

func (m *Migration) Up() error {
    return m.db.CreateTable(&model.UserServerScope{}).Error
}

func (m *Migration) Down() error {
    return m.db.DropTable(&model.UserServerScope{}).Error
}
```

## 3. 权限服务扩展

### 3.1 服务器范围检查

```go
// internal/squ-apiserver/handler/rbac/service.go

// CheckServerScope 检查用户是否有权访问指定服务器
func (s *RBACService) CheckServerScope(userID, serverID uint) (bool, error) {
    // 1. 获取用户角色
    roles, err := s.roleRepo.GetByUserID(userID)
    if err != nil {
        return false, err
    }
    
    // 2. 超级管理员不受限制
    for _, role := range roles {
        if role.Code == "super_admin" {
            return true, nil
        }
    }
    
    // 3. 检查是否有服务器范围限制
    var count int64
    s.db.Model(&model.UserServerScope{}).
        Where("user_id = ?", userID).
        Count(&count)
    
    // 未配置范围限制，默认允许访问所有服务器
    if count == 0 {
        return true, nil
    }
    
    // 4. 检查是否在范围内
    var scope model.UserServerScope
    err = s.db.Where("user_id = ? AND server_id = ?", userID, serverID).First(&scope).Error
    return err == nil, nil
}

// GetUserServerIDs 获取用户可访问的服务器ID列表
func (s *RBACService) GetUserServerIDs(userID uint) ([]uint, error) {
    // 1. 检查是否为超级管理员
    roles, _ := s.roleRepo.GetByUserID(userID)
    for _, role := range roles {
        if role.Code == "super_admin" {
            // 返回所有服务器
            var serverIDs []uint
            s.db.Model(&model.Server{}).Pluck("id", &serverIDs)
            return serverIDs, nil
        }
    }
    
    // 2. 检查是否有范围限制
    var count int64
    s.db.Model(&model.UserServerScope{}).Where("user_id = ?", userID).Count(&count)
    
    if count == 0 {
        // 无限制，返回所有服务器
        var serverIDs []uint
        s.db.Model(&model.Server{}).Pluck("id", &serverIDs)
        return serverIDs, nil
    }
    
    // 3. 返回范围内的服务器
    var serverIDs []uint
    s.db.Model(&model.UserServerScope{}).
        Where("user_id = ?", userID).
        Pluck("server_id", &serverIDs)
    return serverIDs, nil
}
```

### 3.2 带服务器范围的权限检查

```go
// internal/squ-apiserver/handler/rbac/service.go

// CheckPermissionWithServer 检查用户在指定服务器上的权限
func (s *RBACService) CheckPermissionWithServer(userID, serverID uint, resource, action string) (bool, error) {
    // 1. 检查基础权限
    hasPermission, err := s.CheckPermission(userID, resource, action)
    if err != nil || !hasPermission {
        return false, err
    }
    
    // 2. 检查服务器范围
    hasScope, err := s.CheckServerScope(userID, serverID)
    if err != nil || !hasScope {
        return false, err
    }
    
    return true, nil
}
```

## 4. 中间件扩展

```go
// internal/squ-apiserver/middleware/rbac.go

// RBACMiddlewareWithServerScope 带服务器范围检查的权限中间件
func RBACMiddlewareWithServerScope(resource, action string, serverIDParam string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID, exists := c.Get("user_id")
        if !exists {
            c.AbortWithStatusJSON(401, gin.H{"code": 401, "msg": "未登录"})
            return
        }
        
        // 1. 获取服务器 ID
        serverIDStr := c.Param(serverIDParam)
        if serverIDStr == "" {
            // 从查询参数获取
            serverIDStr = c.Query("serverId")
        }
        
        if serverIDStr == "" {
            c.AbortWithStatusJSON(400, gin.H{"code": 400, "msg": "缺少服务器ID"})
            return
        }
        
        serverID, err := strconv.ParseUint(serverIDStr, 10, 64)
        if err != nil {
            c.AbortWithStatusJSON(400, gin.H{"code": 400, "msg": "无效的服务器ID"})
            return
        }
        
        // 2. 检查权限（包含服务器范围）
        rbacService := rbac.GetRBACService()
        hasPermission, err := rbacService.CheckPermissionWithServer(
            userID.(uint), uint(serverID), resource, action,
        )
        
        if err != nil || !hasPermission {
            c.AbortWithStatusJSON(403, gin.H{"code": 403, "msg": "权限不足"})
            return
        }
        
        c.Next()
    }
}
```

## 5. API 接口

### 5.1 类型定义

```go
// internal/squ-apiserver/handler/user/types.go

// UserServerScopeRequest 用户服务器范围请求
type UserServerScopeRequest struct {
    ServerIDs []uint `json:"server_ids" binding:"required"`
}

// UserServerScopeResponse 用户服务器范围响应
type UserServerScopeResponse struct {
    UserID     uint             `json:"user_id"`
    Servers    []ServerScopeItem `json:"servers"`
    Unrestricted bool           `json:"unrestricted"`  // 是否无限制
}

type ServerScopeItem struct {
    ServerID   uint   `json:"server_id"`
    ServerName string `json:"server_name"`
    RoleID     uint   `json:"role_id"`
    RoleName   string `json:"role_name"`
}
```

### 5.2 Handler

```go
// internal/squ-apiserver/handler/user/handler.go

// GetUserServerScopes 获取用户服务器范围
// GET /api/v1/users/:id/servers
func GetUserServerScopes(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    
    scopes, unrestricted, err := userService.GetServerScopes(uint(id))
    if err != nil {
        response.Error(c, 500, "获取服务器范围失败")
        return
    }
    
    response.Success(c, UserServerScopeResponse{
        UserID:       uint(id),
        Servers:      scopes,
        Unrestricted: unrestricted,
    })
}

// UpdateUserServerScopes 更新用户服务器范围
// PUT /api/v1/users/:id/servers
func UpdateUserServerScopes(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    
    var req UserServerScopeRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, 400, "参数错误")
        return
    }
    
    err := userService.UpdateServerScopes(uint(id), req.ServerIDs)
    if err != nil {
        response.Error(c, 500, "更新服务器范围失败")
        return
    }
    
    // 清除缓存
    rbac.ClearUserPermissionCache(uint(id))
    
    response.Success(c, nil)
}
```

### 5.3 Service

```go
// internal/squ-apiserver/handler/user/service.go

func (s *Service) GetServerScopes(userID uint) ([]ServerScopeItem, bool, error) {
    // 检查是否无限制
    var count int64
    s.db.Model(&model.UserServerScope{}).Where("user_id = ?", userID).Count(&count)
    
    if count == 0 {
        return nil, true, nil  // 无限制
    }
    
    // 获取范围列表
    var scopes []model.UserServerScope
    s.db.Where("user_id = ?", userID).Find(&scopes)
    
    var result []ServerScopeItem
    for _, scope := range scopes {
        var server model.Server
        var role model.Role
        s.db.First(&server, scope.ServerID)
        s.db.First(&role, scope.RoleID)
        
        result = append(result, ServerScopeItem{
            ServerID:   scope.ServerID,
            ServerName: server.Hostname,
            RoleID:     scope.RoleID,
            RoleName:   role.Name,
        })
    }
    
    return result, false, nil
}

func (s *Service) UpdateServerScopes(userID uint, serverIDs []uint) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // 删除旧范围
        tx.Where("user_id = ?", userID).Delete(&model.UserServerScope{})
        
        // 获取用户默认角色
        var userRole model.UserRole
        if err := tx.Where("user_id = ?", userID).First(&userRole).Error; err != nil {
            return err
        }
        
        // 创建新范围
        for _, serverID := range serverIDs {
            scope := model.UserServerScope{
                UserID:   userID,
                ServerID: serverID,
                RoleID:   userRole.RoleID,  // 使用默认角色
            }
            if err := tx.Create(&scope).Error; err != nil {
                return err
            }
        }
        
        return nil
    })
}
```

### 5.4 路由注册

```go
// internal/squ-apiserver/router/user.go

func User(group *gin.RouterGroup) {
    users := group.Group("/users")
    
    // ... 现有路由 ...
    
    // 服务器范围管理
    users.GET("/:id/servers",
        middleware.RBACMiddleware(model.ResourceUser, model.ActionGet),
        user.GetUserServerScopes,
    )
    users.PUT("/:id/servers",
        middleware.RBACMiddleware(model.ResourceUser, model.ActionUpdate),
        user.UpdateUserServerScopes,
    )
}
```

## 6. 前端集成

### 6.1 服务器列表过滤

```typescript
// front/src/stores/permission.ts

// 在 permission store 中添加

const accessibleServerIDs = ref<number[]>([])
const serverScopeUnrestricted = ref(true)

async function fetchServerScope() {
  try {
    const data = await getUserServerScopes(currentUserId)
    accessibleServerIDs.value = data.servers.map(s => s.server_id)
    serverScopeUnrestricted.value = data.unrestricted
  } catch (error) {
    console.error('获取服务器范围失败:', error)
  }
}

function canAccessServer(serverId: number): boolean {
  if (isSuperAdmin.value || serverScopeUnrestricted.value) {
    return true
  }
  return accessibleServerIDs.value.includes(serverId)
}
```

### 6.2 服务器选择器过滤

```vue
<!-- 组件：服务器选择下拉框 -->

<template>
  <el-select v-model="selectedServer" placeholder="选择服务器">
    <el-option
      v-for="server in accessibleServers"
      :key="server.id"
      :label="server.hostname"
      :value="server.id"
    />
  </el-select>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { usePermissionStore } from '@/stores/permission'

const props = defineProps<{
  servers: Server[]
}>()

const permissionStore = usePermissionStore()

// 过滤出用户可访问的服务器
const accessibleServers = computed(() => {
  if (permissionStore.isSuperAdmin || permissionStore.serverScopeUnrestricted) {
    return props.servers
  }
  return props.servers.filter(s => 
    permissionStore.accessibleServerIDs.includes(s.id)
  )
})
</script>
```

### 6.3 用户服务器范围管理

```vue
<!-- front/src/views/System/User/components/ServerScope.vue -->

<template>
  <el-dialog
    :model-value="visible"
    title="服务器访问范围"
    width="600px"
    @update:model-value="$emit('update:visible', $event)"
  >
    <div v-loading="loading">
      <el-alert
        v-if="unrestricted"
        type="info"
        :closable="false"
        style="margin-bottom: 15px"
      >
        当前用户可访问所有服务器
      </el-alert>
      
      <el-checkbox
        v-model="restrictAccess"
        style="margin-bottom: 15px"
      >
        限制访问范围
      </el-checkbox>
      
      <el-transfer
        v-if="restrictAccess"
        v-model="selectedServers"
        :data="allServers"
        :titles="['可访问', '不可访问']"
        :props="{ key: 'id', label: 'hostname' }"
      />
    </div>

    <template #footer>
      <el-button @click="$emit('update:visible', false)">取消</el-button>
      <el-button type="primary" @click="handleSave" :loading="saving">
        保存
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { getServerList } from '@/api/server'
import { getUserServerScopes, updateUserServerScopes } from '@/api/user'

const props = defineProps<{
  visible: boolean
  userId: number | null
}>()

const emit = defineEmits<{
  'update:visible': [value: boolean]
  'success': []
}>()

const loading = ref(false)
const saving = ref(false)
const allServers = ref<any[]>([])
const selectedServers = ref<number[]>([])
const unrestricted = ref(true)
const restrictAccess = ref(false)

watch(() => props.visible, async (visible) => {
  if (visible && props.userId) {
    loading.value = true
    try {
      // 获取所有服务器
      const servers = await getServerList()
      allServers.value = servers
      
      // 获取用户当前范围
      const data = await getUserServerScopes(props.userId)
      unrestricted.value = data.unrestricted
      selectedServers.value = data.servers.map(s => s.server_id)
      restrictAccess.value = !data.unrestricted
    } finally {
      loading.value = false
    }
  }
})

async function handleSave() {
  if (!props.userId) return
  
  saving.value = true
  try {
    if (restrictAccess.value) {
      await updateUserServerScopes(props.userId, selectedServers.value)
    } else {
      // 清空范围限制（设置为无限制）
      await updateUserServerScopes(props.userId, [])
    }
    emit('success')
    emit('update:visible', false)
  } finally {
    saving.value = false
  }
}
</script>
```

## 7. 完成检查

- [ ] UserServerScope 模型定义
- [ ] 数据库迁移
- [ ] CheckServerScope 方法
- [ ] GetUserServerIDs 方法
- [ ] CheckPermissionWithServer 方法
- [ ] RBACMiddlewareWithServerScope 中间件
- [ ] 服务器范围管理 API
- [ ] 前端服务器过滤
- [ ] 服务器范围管理组件

## 8. 测试要点

1. **服务器范围检查**
   - 超级管理员不受限制
   - 无配置时默认允许所有
   - 配置后只能访问指定服务器

2. **前端过滤**
   - 服务器下拉框只显示可访问的
   - 直接访问无权限服务器返回 403

3. **范围管理**
   - 设置范围后立即生效
   - 清空范围恢复无限制

## 9. 注意事项

1. **性能考虑**：服务器范围检查会增加数据库查询，建议缓存结果
2. **默认行为**：未配置范围时，是允许所有还是拒绝所有，需根据安全需求确定
3. **与基础权限的关系**：服务器范围是第二层检查，基础权限是第一层
