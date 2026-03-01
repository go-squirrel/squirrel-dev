# Phase 5: 前端集成

## 开发目标

在前端实现权限控制：
- 权限 Store
- 路由守卫
- 按钮级权限控制
- 用户角色管理页面

## 1. 目录结构

```
front/src/
├── api/
│   ├── role.ts           # 新增：角色 API
│   └── permission.ts     # 新增：权限 API
├── stores/
│   ├── user.ts           # 已存在，扩展
│   └── permission.ts     # 新增：权限 Store
├── types/
│   └── rbac.ts           # 新增：RBAC 类型定义
├── composables/
│   └── usePermission.ts  # 新增：权限组合函数
├── views/
│   └── System/
│       ├── Role/
│       │   ├── index.vue        # 角色列表
│       │   └── components/
│       │       └── PermissionEditor.vue  # 权限编辑器
│       └── User/
│           └── components/
│               └── RoleAssign.vue  # 角色分配组件
└── router/
    └── index.ts          # 修改：添加路由守卫
```

## 2. 类型定义

```typescript
// front/src/types/rbac.ts

// 角色
export interface Role {
  id: number
  name: string
  code: string
  description: string
  is_system: boolean
  status: number
}

// 角色详情（含权限）
export interface RoleDetail extends Role {
  permissions: Permission[]
}

// 权限
export interface Permission {
  id: number
  name: string
  code: string
  resource: string
  action: string
  description: string
}

// 权限分组
export interface PermissionGroup {
  resource: string
  permissions: Permission[]
}

// 用户权限响应
export interface UserPermissions {
  roles: string[]
  permissions: string[]
}
```

## 3. API 实现

```typescript
// front/src/api/role.ts

import { get, post, put, del } from '@/utils/request'
import type { Role, RoleDetail } from '@/types/rbac'

// 获取角色列表
export function getRoleList(): Promise<Role[]> {
  return get('/roles')
}

// 获取角色详情
export function getRoleDetail(id: number): Promise<RoleDetail> {
  return get(`/roles/${id}`)
}

// 创建角色
export function createRole(data: Partial<Role>): Promise<Role> {
  return post('/roles', data)
}

// 更新角色
export function updateRole(id: number, data: Partial<Role>): Promise<void> {
  return put(`/roles/${id}`, data)
}

// 删除角色
export function deleteRole(id: number): Promise<void> {
  return del(`/roles/${id}`)
}

// 获取角色权限
export function getRolePermissions(id: number): Promise<Permission[]> {
  return get(`/roles/${id}/permissions`)
}

// 更新角色权限
export function updateRolePermissions(id: number, permissionIds: number[]): Promise<void> {
  return put(`/roles/${id}/permissions`, { permission_ids: permissionIds })
}

// 获取用户角色
export function getUserRoles(userId: number): Promise<{ user_id: number; roles: Role[] }> {
  return get(`/users/${userId}/roles`)
}

// 更新用户角色
export function updateUserRoles(userId: number, roleIds: number[]): Promise<void> {
  return put(`/users/${userId}/roles`, { role_ids: roleIds })
}
```

```typescript
// front/src/api/permission.ts

import { get } from '@/utils/request'
import type { UserPermissions, PermissionGroup } from '@/types/rbac'

// 获取当前用户权限
export function getCurrentUserPermissions(): Promise<UserPermissions> {
  return get('/user/permissions')
}

// 获取所有权限（分组）
export function getAllPermissions(): Promise<PermissionGroup[]> {
  return get('/permissions')
}
```

## 4. 权限 Store

```typescript
// front/src/stores/permission.ts

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { getCurrentUserPermissions } from '@/api/permission'
import type { UserPermissions } from '@/types/rbac'

export const usePermissionStore = defineStore('permission', () => {
  // State
  const roles = ref<string[]>([])
  const permissions = ref<string[]>([])
  const loaded = ref(false)

  // Getters
  const isSuperAdmin = computed(() => roles.value.includes('super_admin'))
  
  const isAdmin = computed(() => 
    roles.value.includes('super_admin') || roles.value.includes('admin')
  )

  // Actions
  async function fetchPermissions() {
    if (loaded.value) return
    
    try {
      const data = await getCurrentUserPermissions()
      roles.value = data.roles
      permissions.value = data.permissions
      loaded.value = true
    } catch (error) {
      console.error('获取权限失败:', error)
    }
  }

  // 检查单个权限
  function hasPermission(code: string): boolean {
    // 超级管理员拥有所有权限
    if (isSuperAdmin.value) return true
    
    // 精确匹配
    if (permissions.value.includes(code)) return true
    
    // 通配符匹配：server:* 匹配 server:list, server:get 等
    const [resource] = code.split(':')
    if (permissions.value.includes(`${resource}:*`)) return true
    
    return false
  }

  // 检查多个权限（满足任一即可）
  function hasAnyPermission(codes: string[]): boolean {
    if (isSuperAdmin.value) return true
    return codes.some(code => hasPermission(code))
  }

  // 检查多个权限（需全部满足）
  function hasAllPermissions(codes: string[]): boolean {
    if (isSuperAdmin.value) return true
    return codes.every(code => hasPermission(code))
  }

  // 清除权限（退出登录时）
  function clearPermissions() {
    roles.value = []
    permissions.value = []
    loaded.value = false
  }

  return {
    roles,
    permissions,
    loaded,
    isSuperAdmin,
    isAdmin,
    fetchPermissions,
    hasPermission,
    hasAnyPermission,
    hasAllPermissions,
    clearPermissions,
  }
})
```

## 5. 权限组合函数

```typescript
// front/src/composables/usePermission.ts

import { usePermissionStore } from '@/stores/permission'

export function usePermission() {
  const permissionStore = usePermissionStore()

  // 权限检查函数
  const check = (code: string): boolean => {
    return permissionStore.hasPermission(code)
  }

  // 多权限检查（任一）
  const checkAny = (codes: string[]): boolean => {
    return permissionStore.hasAnyPermission(codes)
  }

  // 权限指令值
  const can = (code: string | string[]): boolean => {
    if (Array.isArray(code)) {
      return checkAny(code)
    }
    return check(code)
  }

  return {
    check,
    checkAny,
    can,
    isSuperAdmin: permissionStore.isSuperAdmin,
    isAdmin: permissionStore.isAdmin,
  }
}
```

## 6. 路由守卫

```typescript
// front/src/router/index.ts

import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { usePermissionStore } from '@/stores/permission'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    // ... 现有路由配置
  ],
})

// 全局前置守卫
router.beforeEach(async (to, from) => {
  const userStore = useUserStore()
  const permissionStore = usePermissionStore()

  // 不需要登录的页面
  const publicPages = ['/login']
  if (publicPages.includes(to.path)) {
    return
  }

  // 检查登录状态
  if (!userStore.isLoggedIn) {
    return { path: '/login', query: { redirect: to.fullPath } }
  }

  // 获取权限（仅首次）
  if (!permissionStore.loaded) {
    await permissionStore.fetchPermissions()
  }

  // 检查路由权限
  const requiredPermissions = to.meta.permissions as string[] | undefined
  if (requiredPermissions && requiredPermissions.length > 0) {
    if (!permissionStore.hasAnyPermission(requiredPermissions)) {
      // 无权限，跳转到 403 页面
      return { path: '/403' }
    }
  }
})

export default router
```

### 6.1 路由权限配置

```typescript
// front/src/router/index.ts

const routes = [
  {
    path: '/users',
    name: 'Users',
    component: () => import('@/views/System/User/index.vue'),
    meta: {
      permissions: ['user:list'],  // 需要的权限
    },
  },
  {
    path: '/roles',
    name: 'Roles',
    component: () => import('@/views/System/Role/index.vue'),
    meta: {
      permissions: ['role:list'],
    },
  },
  {
    path: '/servers',
    name: 'Servers',
    component: () => import('@/views/Server/index.vue'),
    meta: {
      permissions: ['server:list'],
    },
  },
  {
    path: '/servers/:id/terminal',
    name: 'Terminal',
    component: () => import('@/views/Terminal/index.vue'),
    meta: {
      permissions: ['terminal:exec'],
    },
  },
  // ...
]
```

## 7. 权限指令

```typescript
// front/src/directives/permission.ts

import type { Directive } from 'vue'
import { usePermissionStore } from '@/stores/permission'

export const vPermission: Directive<HTMLElement, string | string[]> = {
  mounted(el, binding) {
    const permissionStore = usePermissionStore()
    const value = binding.value
    
    const hasPermission = Array.isArray(value)
      ? permissionStore.hasAnyPermission(value)
      : permissionStore.hasPermission(value)
    
    if (!hasPermission) {
      el.style.display = 'none'
      // 或者直接移除元素
      // el.remove()
    }
  },
}

// 注册指令
// main.ts
// app.directive('permission', vPermission)
```

### 7.1 使用示例

```vue
<template>
  <!-- 单个权限 -->
  <button v-permission="'server:create'">添加服务器</button>
  
  <!-- 多个权限（满足任一即可） -->
  <button v-permission="['user:update', 'role:update']">编辑</button>
  
  <!-- 使用 v-if -->
  <button v-if="can('server:delete')">删除</button>
</template>

<script setup lang="ts">
import { usePermission } from '@/composables/usePermission'

const { can } = usePermission()
</script>
```

## 8. 角色管理页面

### 8.1 角色列表

```vue
<!-- front/src/views/System/Role/index.vue -->

<template>
  <div class="role-page">
    <div class="page-header">
      <h2>角色管理</h2>
      <el-button v-permission="'user:update'" type="primary" @click="showCreate">
        创建角色
      </el-button>
    </div>

    <el-table :data="roles" v-loading="loading">
      <el-table-column prop="name" label="角色名称" />
      <el-table-column prop="code" label="角色编码" />
      <el-table-column prop="description" label="描述" />
      <el-table-column prop="is_system" label="系统角色">
        <template #default="{ row }">
          <el-tag v-if="row.is_system" type="info">系统</el-tag>
          <el-tag v-else>自定义</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200">
        <template #default="{ row }">
          <el-button link @click="showPermissions(row)">权限</el-button>
          <el-button 
            v-if="!row.is_system" 
            link 
            v-permission="'user:update'"
            @click="showEdit(row)"
          >
            编辑
          </el-button>
          <el-button 
            v-if="!row.is_system" 
            link 
            type="danger"
            v-permission="'user:update'"
            @click="handleDelete(row)"
          >
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 权限编辑对话框 -->
    <PermissionEditor
      v-model:visible="permissionDialogVisible"
      :role="currentRole"
      @success="loadRoles"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getRoleList, deleteRole } from '@/api/role'
import type { Role } from '@/types/rbac'
import PermissionEditor from './components/PermissionEditor.vue'

const roles = ref<Role[]>([])
const loading = ref(false)
const permissionDialogVisible = ref(false)
const currentRole = ref<Role | null>(null)

onMounted(() => {
  loadRoles()
})

async function loadRoles() {
  loading.value = true
  try {
    roles.value = await getRoleList()
  } finally {
    loading.value = false
  }
}

function showPermissions(role: Role) {
  currentRole.value = role
  permissionDialogVisible.value = true
}

async function handleDelete(role: Role) {
  if (!confirm(`确定删除角色 "${role.name}"？`)) return
  
  await deleteRole(role.id)
  loadRoles()
}
</script>
```

### 8.2 权限编辑器

```vue
<!-- front/src/views/System/Role/components/PermissionEditor.vue -->

<template>
  <el-dialog
    :model-value="visible"
    :title="`编辑权限 - ${role?.name || ''}`"
    width="600px"
    @update:model-value="$emit('update:visible', $event)"
  >
    <div v-loading="loading">
      <div v-for="group in permissionGroups" :key="group.resource" class="permission-group">
        <h4>{{ getResourceLabel(group.resource) }}</h4>
        <el-checkbox-group v-model="selectedPermissions">
          <el-checkbox
            v-for="perm in group.permissions"
            :key="perm.id"
            :label="perm.id"
          >
            {{ perm.name }}
          </el-checkbox>
        </el-checkbox-group>
      </div>
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
import { getAllPermissions } from '@/api/permission'
import { getRolePermissions, updateRolePermissions } from '@/api/role'
import type { PermissionGroup, Permission } from '@/types/rbac'

const props = defineProps<{
  visible: boolean
  role: { id: number; name: string } | null
}>()

const emit = defineEmits<{
  'update:visible': [value: boolean]
  'success': []
}>()

const loading = ref(false)
const saving = ref(false)
const permissionGroups = ref<PermissionGroup[]>([])
const selectedPermissions = ref<number[]>([])

// 资源名称映射
const resourceLabels: Record<string, string> = {
  server: '服务器',
  application: '应用',
  monitor: '监控',
  script: '脚本',
  file: '文件',
  terminal: '终端',
  user: '用户',
  role: '角色',
}

function getResourceLabel(resource: string): string {
  return resourceLabels[resource] || resource
}

watch(() => props.visible, async (visible) => {
  if (visible && props.role) {
    loading.value = true
    try {
      // 获取所有权限
      permissionGroups.value = await getAllPermissions()
      
      // 获取角色当前权限
      const rolePerms = await getRolePermissions(props.role.id)
      selectedPermissions.value = rolePerms.map(p => p.id)
    } finally {
      loading.value = false
    }
  }
})

async function handleSave() {
  if (!props.role) return
  
  saving.value = true
  try {
    await updateRolePermissions(props.role.id, selectedPermissions.value)
    emit('success')
    emit('update:visible', false)
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
.permission-group {
  margin-bottom: 20px;
}
.permission-group h4 {
  margin-bottom: 10px;
  color: #606266;
}
</style>
```

### 8.3 用户角色分配组件

```vue
<!-- front/src/views/System/User/components/RoleAssign.vue -->

<template>
  <el-dialog
    :model-value="visible"
    title="分配角色"
    width="400px"
    @update:model-value="$emit('update:visible', $event)"
  >
    <el-checkbox-group v-model="selectedRoles" v-loading="loading">
      <div v-for="role in roles" :key="role.id" class="role-item">
        <el-checkbox :label="role.id" :disabled="role.status !== 1">
          {{ role.name }}
          <span class="role-code">({{ role.code }})</span>
        </el-checkbox>
        <div class="role-desc">{{ role.description }}</div>
      </div>
    </el-checkbox-group>

    <template #footer>
      <el-button @click="$emit('update:visible', false)">取消</el-button>
      <el-button type="primary" @click="handleSave" :loading="saving">
        保存
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { getRoleList } from '@/api/role'
import { getUserRoles, updateUserRoles } from '@/api/role'
import type { Role } from '@/types/rbac'

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
const roles = ref<Role[]>([])
const selectedRoles = ref<number[]>([])

watch(() => props.visible, async (visible) => {
  if (visible && props.userId) {
    loading.value = true
    try {
      // 获取所有角色
      roles.value = await getRoleList()
      
      // 获取用户当前角色
      const userRoles = await getUserRoles(props.userId)
      selectedRoles.value = userRoles.roles.map(r => r.id)
    } finally {
      loading.value = false
    }
  }
})

async function handleSave() {
  if (!props.userId) return
  
  saving.value = true
  try {
    await updateUserRoles(props.userId, selectedRoles.value)
    emit('success')
    emit('update:visible', false)
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
.role-item {
  margin-bottom: 15px;
}
.role-code {
  color: #909399;
  font-size: 12px;
}
.role-desc {
  margin-left: 24px;
  color: #909399;
  font-size: 12px;
}
</style>
```

## 9. 完成检查

- [ ] RBAC 类型定义
- [ ] 角色 API
- [ ] 权限 API
- [ ] Permission Store
- [ ] usePermission 组合函数
- [ ] 路由守卫
- [ ] 权限指令 v-permission
- [ ] 角色列表页面
- [ ] 权限编辑器组件
- [ ] 用户角色分配组件

## 10. 测试要点

1. **权限检查**
   - 不同角色用户看到不同菜单
   - 无权限按钮隐藏
   - 无权限路由跳转 403

2. **角色管理**
   - 创建/编辑/删除角色
   - 权限分配保存正确

3. **用户角色分配**
   - 分配角色后权限立即生效
   - 缓存清除正常

## 11. 下一阶段

完成本阶段后，可选开发 [06-server-scope.md](./06-server-scope.md) 实现服务器级权限隔离。
