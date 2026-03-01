# Phase 4: 前端文件管理页面

## 前置条件

已完成 Phase 1-3：
- Agent 端文件操作功能完整
- APIServer 代理转发已实现

## 开发目标

实现前端文件管理页面：
- 文件浏览器（文件列表）
- 文件编辑器（在线编辑）
- 文件上传/下载
- 基础操作（创建、删除、重命名）

## 1. 目录结构

```
front/src/
├── views/
│   └── FileManager/
│       ├── index.vue              # 主页面
│       └── components/
│           ├── FileList.vue       # 文件列表
│           ├── FileEditor.vue     # 在线编辑器
│           ├── FileUpload.vue     # 上传组件
│           └── Breadcrumb.vue     # 路径导航
├── api/
│   └── filesystem.ts              # API 接口
├── types/
│   └── filesystem.ts              # 类型定义
└── composables/
    └── useFileManager.ts          # 文件操作逻辑
```

## 2. 类型定义

```typescript
// front/src/types/filesystem.ts

// 文件信息
interface FileInfo {
  name: string
  path: string
  is_dir: boolean
  size: number
  mode: string
  mod_time: string
  extension: string
  mime: string
}

// 目录列表响应
interface ListResponse {
  path: string
  parent: string
  entries: FileInfo[]
}

// 文件读取响应
interface ReadResponse {
  path: string
  content: string
  size: number
  mime: string
  encoding: string
}
```

## 3. API 接口

```typescript
// front/src/api/filesystem.ts

// 目录列表
listDir(serverId, path, params?) → ListResponse

// 文件读取
readFile(serverId, path) → ReadResponse

// 文件下载 URL
getDownloadUrl(serverId, path) → string

// 文件写入
writeFile(serverId, { path, content, backup? }) → { path, size }

// 文件上传
uploadFile(serverId, path, file) → { path, size, name }

// 创建目录
mkdir(serverId, { path }) → { path }

// 创建文件
createFile(serverId, path, content?) → { path }

// 删除
deleteFile(serverId, { path, recursive? }) → { path }

// 重命名
rename(serverId, { old_path, new_path }) → { old_path, new_path }
```

## 4. 页面布局

```
┌──────────────────────────────────────────────────────────────────┐
│ 📁 文件管理              [服务器 ▼]      [上传] [新建目录] [新建文件] │
├──────────────────────────────────────────────────────────────────┤
│ 📍 /home/app/config > logs                                       │
├──────────────────────────────────────────────────────────────────┤
│ 名称              │ 大小   │ 修改时间      │ 权限  │ 操作        │
│ 📂 logs/          │ -      │ 2024-01-15    │ 755  │ [删除]      │
│ 📄 config.yaml    │ 2KB    │ 2024-01-14    │ 644  │ [编辑][下载] │
│ 📄 app.log        │ 5MB    │ 2024-01-15    │ 644  │ [查看][下载] │
└──────────────────────────────────────────────────────────────────┘
```

## 5. 主页面逻辑

```vue
<!-- FileManager/index.vue 伪代码 -->

<template>
  <div class="file-manager">
    <!-- 工具栏：服务器选择 + 操作按钮 -->
    <Toolbar 
      :servers="serverList"
      v-model:selected-server="selectedServerId"
      @upload="openUpload"
      @mkdir="handleMkdir"
      @create="handleCreateFile"
    />
    
    <!-- 路径导航 -->
    <Breadcrumb :path="currentPath" @navigate="navigate" />
    
    <!-- 文件列表 -->
    <FileList 
      :entries="fileList"
      :loading="loading"
      @open="handleOpen"
      @edit="openEditor"
      @download="handleDownload"
      @delete="handleDelete"
      @rename="handleRename"
    />
    
    <!-- 编辑器弹窗 -->
    <FileEditor v-model:visible="editorVisible" :file="editFile" />
    
    <!-- 上传弹窗 -->
    <FileUpload v-model:visible="uploadVisible" :path="currentPath" />
  </div>
</template>

<script setup>
// 状态
const serverList = ref([])
const selectedServerId = ref(null)
const currentPath = ref('/')
const fileList = ref([])
const loading = ref(false)

// 加载服务器列表后自动加载根目录
onMounted(async () => {
  serverList.value = await getServerList()
  if (serverList.value.length) {
    selectedServerId.value = serverList.value[0].id
    loadFileList()
  }
})

// 加载文件列表
async function loadFileList() {
  loading.value = true
  fileList.value = await listDir(selectedServerId.value, currentPath.value)
  loading.value = false
}

// 打开文件/目录
function handleOpen(file) {
  if (file.is_dir) {
    currentPath.value = file.path
    loadFileList()
  } else {
    openEditor(file)
  }
}

// 下载文件
function handleDownload(file) {
  window.open(getDownloadUrl(selectedServerId.value, file.path))
}

// 删除文件（带确认）
async function handleDelete(file) {
  await confirm(`确定删除 ${file.name}？`)
  await deleteFile(selectedServerId.value, { path: file.path, recursive: file.is_dir })
  loadFileList()
}
</script>
```

## 6. 文件列表组件

```vue
<!-- FileList.vue 伪代码 -->

<template>
  <el-table :data="entries" v-loading="loading">
    <!-- 名称列（带图标） -->
    <el-table-column label="名称">
      <FolderIcon v-if="row.is_dir" />
      <FileIcon v-else />
      {{ row.name }}
    </el-table-column>
    
    <!-- 大小列 -->
    <el-table-column label="大小">
      {{ row.is_dir ? '-' : formatSize(row.size) }}
    </el-table-column>
    
    <!-- 修改时间 -->
    <el-table-column label="修改时间">
      {{ formatDate(row.mod_time) }}
    </el-table-column>
    
    <!-- 操作列 -->
    <el-table-column label="操作">
      <el-button @click="emit('edit', row)" v-if="!row.is_dir">编辑</el-button>
      <el-button @click="emit('download', row)" v-if="!row.is_dir">下载</el-button>
      <el-button @click="emit('rename', row)">重命名</el-button>
      <el-button @click="emit('delete', row)" type="danger">删除</el-button>
    </el-table-column>
  </el-table>
</template>
```

## 7. 文件编辑器

```vue
<!-- FileEditor.vue 伪代码 -->

<template>
  <el-dialog :title="file?.name" width="80%">
    <!-- 使用 Monaco Editor 或 CodeMirror -->
    <MonacoEditor v-model="content" :language="detectLanguage(file)" />
    
    <template #footer>
      <el-button @click="close">取消</el-button>
      <el-button type="primary" @click="save" :loading="saving">保存</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
const content = ref('')

// 加载文件内容
watch(() => props.file, async (file) => {
  if (file) {
    const res = await readFile(serverId, file.path)
    content.value = res.content
  }
}, { immediate: true })

// 保存
async function save() {
  await writeFile(serverId, { path: props.file.path, content: content.value, backup: true })
  emit('success')
}
</script>
```

## 8. 上传组件

```vue
<!-- FileUpload.vue 伪代码 -->

<template>
  <el-dialog title="上传文件">
    <el-upload
      drag
      :action="uploadUrl"
      :on-success="handleSuccess"
      :on-error="handleError"
    >
      拖拽文件到此处或点击上传
    </el-upload>
  </el-dialog>
</template>

<script setup>
const uploadUrl = computed(() => `/api/v1/fs/${serverId}/upload?path=${path}`)

function handleSuccess() {
  emit('success')
  ElMessage.success('上传成功')
}
</script>
```

## 9. 路由配置

```typescript
// front/src/router/index.ts

{
  path: '/files',
  name: 'FileManager',
  component: () => import('@/views/FileManager/index.vue'),
  meta: { title: '文件管理' }
}
```

## 10. 完成检查

- [ ] 类型定义完成
- [ ] API 接口封装完成
- [ ] 主页面实现完成
- [ ] 文件列表组件完成
- [ ] 文件编辑器完成
- [ ] 上传组件完成
- [ ] 路由配置完成
- [ ] 功能测试通过

## 11. 下一阶段

完成本阶段后，文件管理核心功能已完整。如需大文件传输、多节点分发等功能，继续开发 [05-transfer-center.md](./05-transfer-center.md)。
