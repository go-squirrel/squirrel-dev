# Phase 4: 前端文件管理页面

## 前置条件

已完成 Phase 1-3：
- Agent 端文件操作功能完整
- APIServer 文件存储服务已实现

## 开发目标

实现前端文件管理页面：
- 文件列表（APIServer 文件库）
- 文件上传（支持分片上传）
- 文件删除、下载
- 触发分发到 Agent 节点

## 1. 目录结构

```
front/src/
├── views/
│   └── FileManager/
│       ├── index.vue              # 主页面
│       └── components/
│           ├── FileList.vue       # 文件列表
│           ├── FileUpload.vue     # 上传组件
│           ├── ChunkUpload.vue    # 分片上传组件
│           └── DispatchDialog.vue # 分发弹窗
├── api/
│   └── files.ts                   # API 接口
├── types/
│   └── files.ts                   # 类型定义
└── composables/
    └── useFileUpload.ts           # 分片上传逻辑
```

## 2. 类型定义

```typescript
// front/src/types/files.ts

// 文件记录
interface FileRecord {
  id: number
  uuid: string
  name: string
  size: number
  path: string
  md5: string
  mime: string
  status: 'uploaded' | 'distributing' | 'distributed'
  created_at: string
  updated_at: string
}

// 分片上传初始化响应
interface ChunkInitResponse {
  upload_id: string
  chunk_size: number
  chunk_count: number
}

// 分发任务
interface DispatchTask {
  id: number
  file_id: number
  server_id: number
  server_name?: string
  target_path: string
  status: 'pending' | 'running' | 'success' | 'failed'
  progress: number
  error_msg?: string
  created_at: string
  finished_at?: string
}

// 服务器信息
interface Server {
  id: number
  name: string
  ip_address: string
  status: 'online' | 'offline'
}
```

## 3. API 接口

```typescript
// front/src/api/files.ts

import request from '@/utils/request'
import type { FileRecord, DispatchTask, Server, ChunkInitResponse } from '@/types/files'

// ==================== 文件管理 ====================

// 文件列表
export function listFiles(params: { page: number; page_size: number; name?: string }) {
  return request.get<{ list: FileRecord[]; total: number }>('/api/v1/files', { params })
}

// 文件详情
export function getFile(id: number) {
  return request.get<FileRecord>(`/api/v1/files/${id}`)
}

// 删除文件
export function deleteFile(id: number) {
  return request.delete(`/api/v1/files/${id}`)
}

// 下载文件 URL
export function getDownloadUrl(id: number) {
  return `/api/v1/files/${id}/download`
}

// 普通上传
export function uploadFile(file: File) {
  const formData = new FormData()
  formData.append('file', file)
  return request.post<FileRecord>('/api/v1/files/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}

// ==================== 分片上传 ====================

// 初始化分片上传
export function initChunkUpload(data: { file_name: string; file_size: number; file_md5?: string }) {
  return request.post<ChunkInitResponse>('/api/v1/files/chunk/init', data)
}

// 上传分片
export function uploadChunk(uploadId: string, chunkIndex: number, chunk: Blob) {
  return request.post(`/api/v1/files/chunk/upload?upload_id=${uploadId}&chunk_index=${chunkIndex}`, chunk, {
    headers: { 'Content-Type': 'application/octet-stream' }
  })
}

// 查询分片上传状态
export function getChunkStatus(uploadId: string) {
  return request.get('/api/v1/files/chunk/status', { params: { upload_id: uploadId } })
}

// 合并分片
export function completeChunkUpload(uploadId: string) {
  return request.post<FileRecord>('/api/v1/files/chunk/complete', { upload_id: uploadId })
}

// ==================== 分发任务 ====================

// 创建分发任务
export function createDispatch(data: { file_id: number; server_ids: number[]; target_path: string }) {
  return request.post('/api/v1/dispatch', data)
}

// 分发任务列表
export function listDispatchTasks(params?: { file_id?: number; status?: string }) {
  return request.get<DispatchTask[]>('/api/v1/dispatch/tasks', { params })
}

// 分发任务详情
export function getDispatchTask(id: number) {
  return request.get<DispatchTask>(`/api/v1/dispatch/tasks/${id}`)
}

// ==================== 服务器 ====================

// 服务器列表
export function listServers() {
  return request.get<Server[]>('/api/v1/servers')
}
```

## 4. 分片上传逻辑

```typescript
// front/src/composables/useFileUpload.ts

import { ref, computed } from 'vue'
import { initChunkUpload, uploadChunk, completeChunkUpload } from '@/api/files'
import type { FileRecord } from '@/types/files'

// 分片大小阈值（50MB）
const CHUNK_THRESHOLD = 50 * 1024 * 1024

export function useFileUpload() {
  const uploading = ref(false)
  const progress = ref(0)
  const currentFile = ref<File | null>(null)

  // 是否使用分片上传
  const shouldUseChunk = computed(() => {
    return currentFile.value && currentFile.value.size > CHUNK_THRESHOLD
  })

  // 上传文件
  async function upload(file: File): Promise<FileRecord> {
    currentFile.value = file
    uploading.value = true
    progress.value = 0

    try {
      if (file.size > CHUNK_THRESHOLD) {
        // 分片上传
        return await chunkUpload(file)
      } else {
        // 普通上传
        const result = await normalUpload(file)
        progress.value = 100
        return result
      }
    } finally {
      uploading.value = false
    }
  }

  // 普通上传
  async function normalUpload(file: File): Promise<FileRecord> {
    const formData = new FormData()
    formData.append('file', file)
    
    // 使用 axios 的 onUploadProgress
    const response = await fetch('/api/v1/files/upload', {
      method: 'POST',
      body: formData,
    })
    
    return response.json()
  }

  // 分片上传
  async function chunkUpload(file: File): Promise<FileRecord> {
    // 1. 初始化
    const initRes = await initChunkUpload({
      file_name: file.name,
      file_size: file.size,
    })

    const { upload_id, chunk_size, chunk_count } = initRes.data
    const uploadedChunks = new Set<number>()

    // 2. 上传分片（并发控制）
    const concurrent = 3 // 并发数
    const chunks: Promise<void>[] = []

    for (let i = 0; i < chunk_count; i++) {
      const start = i * chunk_size
      const end = Math.min(start + chunk_size, file.size)
      const chunk = file.slice(start, end)

      // 等待并发槽位
      while (chunks.length >= concurrent) {
        await Promise.race(chunks)
      }

      const promise = uploadChunk(upload_id, i, chunk)
        .then(() => {
          uploadedChunks.add(i)
          progress.value = Math.round((uploadedChunks.size / chunk_count) * 100)
        })
        .finally(() => {
          const idx = chunks.indexOf(promise)
          if (idx > -1) chunks.splice(idx, 1)
        })

      chunks.push(promise)
    }

    // 等待所有分片完成
    await Promise.all(chunks)

    // 3. 合并分片
    const result = await completeChunkUpload(upload_id)
    progress.value = 100

    return result.data
  }

  return {
    upload,
    uploading,
    progress,
    shouldUseChunk,
  }
}
```

## 5. 页面布局

```
┌──────────────────────────────────────────────────────────────────────────┐
│ 📁 文件管理                                          [上传文件] [刷新]   │
├──────────────────────────────────────────────────────────────────────────┤
│ 搜索: [___________________]  状态: [全部 ▼]                              │
├──────────────────────────────────────────────────────────────────────────┤
│ 名称              │ 大小    │ 状态     │ 上传时间      │ 操作            │
│ 📄 app.tar.gz     │ 150MB   │ 已上传   │ 2024-01-15   │ [分发][下载][删] │
│ 📄 config.yaml    │ 2KB     │ 分发中   │ 2024-01-14   │ [分发][下载][删] │
│ 📄 nginx.conf     │ 5KB     │ 已分发   │ 2024-01-13   │ [分发][下载][删] │
├──────────────────────────────────────────────────────────────────────────┤
│ 共 3 条记录                              < 1 2 3 >                        │
└──────────────────────────────────────────────────────────────────────────┘

[分发弹窗]
┌──────────────────────────────────────────────────────────────────┐
│ 分发文件: app.tar.gz                                              │
│ ──────────────────────────────────────────────────────────────── │
│ 目标服务器:  [✓] server-1 (在线)  [✓] server-2 (在线)            │
│             [ ] server-3 (离线)                                   │
│ 目标路径:    [/opt/app/_________________________]                 │
│                                                                   │
│                              [取消] [开始分发]                    │
└──────────────────────────────────────────────────────────────────┘
```

## 6. 主页面实现

```vue
<!-- FileManager/index.vue -->

<template>
  <div class="file-manager">
    <!-- 工具栏 -->
    <div class="toolbar">
      <h2>文件管理</h2>
      <div class="actions">
        <el-button type="primary" @click="openUpload">
          <el-icon><Upload /></el-icon>
          上传文件
        </el-button>
        <el-button @click="loadFiles">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </div>
    </div>

    <!-- 搜索栏 -->
    <div class="search-bar">
      <el-input
        v-model="searchName"
        placeholder="搜索文件名"
        clearable
        @clear="loadFiles"
        @keyup.enter="loadFiles"
        style="width: 200px"
      />
      <el-select v-model="filterStatus" placeholder="状态筛选" clearable style="width: 120px">
        <el-option label="已上传" value="uploaded" />
        <el-option label="分发中" value="distributing" />
        <el-option label="已分发" value="distributed" />
      </el-select>
    </div>

    <!-- 文件列表 -->
    <FileList
      :files="fileList"
      :loading="loading"
      @dispatch="openDispatch"
      @download="handleDownload"
      @delete="handleDelete"
    />

    <!-- 分页 -->
    <el-pagination
      v-model:current-page="page"
      v-model:page-size="pageSize"
      :total="total"
      :page-sizes="[10, 20, 50]"
      layout="total, sizes, prev, pager, next"
      @change="loadFiles"
    />

    <!-- 上传弹窗 -->
    <FileUpload v-model:visible="uploadVisible" @success="handleUploadSuccess" />

    <!-- 分发弹窗 -->
    <DispatchDialog
      v-model:visible="dispatchVisible"
      :file="dispatchFile"
      @success="handleDispatchSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listFiles, deleteFile, getDownloadUrl } from '@/api/files'
import type { FileRecord } from '@/types/files'
import FileList from './components/FileList.vue'
import FileUpload from './components/FileUpload.vue'
import DispatchDialog from './components/DispatchDialog.vue'

// 状态
const fileList = ref<FileRecord[]>([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const searchName = ref('')
const filterStatus = ref('')

const uploadVisible = ref(false)
const dispatchVisible = ref(false)
const dispatchFile = ref<FileRecord | null>(null)

// 加载文件列表
async function loadFiles() {
  loading.value = true
  try {
    const res = await listFiles({
      page: page.value,
      page_size: pageSize.value,
      name: searchName.value || undefined,
    })
    fileList.value = res.data.list
    total.value = res.data.total
  } finally {
    loading.value = false
  }
}

// 打开上传弹窗
function openUpload() {
  uploadVisible.value = true
}

// 上传成功
function handleUploadSuccess() {
  ElMessage.success('上传成功')
  loadFiles()
}

// 打开发发弹窗
function openDispatch(file: FileRecord) {
  dispatchFile.value = file
  dispatchVisible.value = true
}

// 分发成功
function handleDispatchSuccess() {
  ElMessage.success('分发任务已创建')
  loadFiles()
}

// 下载文件
function handleDownload(file: FileRecord) {
  window.open(getDownloadUrl(file.id))
}

// 删除文件
async function handleDelete(file: FileRecord) {
  await ElMessageBox.confirm(`确定删除文件 "${file.name}"？`, '确认删除', {
    type: 'warning',
  })
  await deleteFile(file.id)
  ElMessage.success('删除成功')
  loadFiles()
}

// 状态筛选变化
watch(filterStatus, () => {
  page.value = 1
  loadFiles()
})

onMounted(() => {
  loadFiles()
})
</script>

<style scoped>
.file-manager {
  padding: 20px;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.search-bar {
  display: flex;
  gap: 10px;
  margin-bottom: 16px;
}
</style>
```

## 7. 文件列表组件

```vue
<!-- FileManager/components/FileList.vue -->

<template>
  <el-table :data="files" v-loading="loading" stripe>
    <!-- 文件名 -->
    <el-table-column label="文件名" prop="name" min-width="200">
      <template #default="{ row }">
        <div class="file-name">
          <el-icon :size="20"><Document /></el-icon>
          <span>{{ row.name }}</span>
        </div>
      </template>
    </el-table-column>

    <!-- 大小 -->
    <el-table-column label="大小" width="100">
      <template #default="{ row }">
        {{ formatSize(row.size) }}
      </template>
    </el-table-column>

    <!-- 状态 -->
    <el-table-column label="状态" width="100">
      <template #default="{ row }">
        <el-tag :type="getStatusType(row.status)">
          {{ getStatusLabel(row.status) }}
        </el-tag>
      </template>
    </el-table-column>

    <!-- 上传时间 -->
    <el-table-column label="上传时间" width="180">
      <template #default="{ row }">
        {{ formatDate(row.created_at) }}
      </template>
    </el-table-column>

    <!-- 操作 -->
    <el-table-column label="操作" width="200" fixed="right">
      <template #default="{ row }">
        <el-button size="small" type="primary" link @click="emit('dispatch', row)">
          分发
        </el-button>
        <el-button size="small" link @click="emit('download', row)">
          下载
        </el-button>
        <el-button size="small" type="danger" link @click="emit('delete', row)">
          删除
        </el-button>
      </template>
    </el-table-column>
  </el-table>
</template>

<script setup lang="ts">
import type { FileRecord } from '@/types/files'

defineProps<{
  files: FileRecord[]
  loading: boolean
}>()

const emit = defineEmits<{
  dispatch: [file: FileRecord]
  download: [file: FileRecord]
  delete: [file: FileRecord]
}>()

// 格式化文件大小
function formatSize(bytes: number): string {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  if (bytes < 1024 * 1024 * 1024) return (bytes / 1024 / 1024).toFixed(1) + ' MB'
  return (bytes / 1024 / 1024 / 1024).toFixed(2) + ' GB'
}

// 状态类型
function getStatusType(status: string): string {
  const map: Record<string, string> = {
    uploaded: 'info',
    distributing: 'warning',
    distributed: 'success',
  }
  return map[status] || 'info'
}

// 状态标签
function getStatusLabel(status: string): string {
  const map: Record<string, string> = {
    uploaded: '已上传',
    distributing: '分发中',
    distributed: '已分发',
  }
  return map[status] || status
}

// 格式化日期
function formatDate(date: string): string {
  return new Date(date).toLocaleString()
}
</script>

<style scoped>
.file-name {
  display: flex;
  align-items: center;
  gap: 8px;
}
</style>
```

## 8. 上传组件

```vue
<!-- FileManager/components/FileUpload.vue -->

<template>
  <el-dialog v-model="visible" title="上传文件" width="500px">
    <!-- 拖拽上传区域 -->
    <el-upload
      drag
      :auto-upload="false"
      :show-file-list="false"
      :on-change="handleFileChange"
    >
      <el-icon class="el-icon--upload"><UploadFilled /></el-icon>
      <div class="el-upload__text">
        拖拽文件到此处或 <em>点击选择</em>
      </div>
      <template #tip>
        <div class="el-upload__tip">
          支持任意类型文件，大文件自动启用分片上传
        </div>
      </template>
    </el-upload>

    <!-- 上传进度 -->
    <div v-if="uploading" class="upload-progress">
      <p>{{ currentFile?.name }}</p>
      <el-progress :percentage="progress" :status="progress === 100 ? 'success' : ''" />
    </div>

    <template #footer>
      <el-button @click="close">取消</el-button>
      <el-button type="primary" :loading="uploading" :disabled="!currentFile" @click="doUpload">
        上传
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useFileUpload } from '@/composables/useFileUpload'

const props = defineProps<{
  visible: boolean
}>()

const emit = defineEmits<{
  'update:visible': [value: boolean]
  success: []
}>()

const { upload, uploading, progress } = useFileUpload()
const currentFile = ref<File | null>(null)

// 选择文件
function handleFileChange(file: any) {
  currentFile.value = file.raw
}

// 执行上传
async function doUpload() {
  if (!currentFile.value) return
  
  try {
    await upload(currentFile.value)
    emit('success')
    close()
  } catch (error: any) {
    console.error('上传失败', error)
  }
}

// 关闭弹窗
function close() {
  currentFile.value = null
  emit('update:visible', false)
}

// 重置状态
watch(() => props.visible, (val) => {
  if (!val) {
    currentFile.value = null
  }
})
</script>

<style scoped>
.upload-progress {
  margin-top: 16px;
  padding: 12px;
  background: #f5f7fa;
  border-radius: 4px;
}
</style>
```

## 9. 分发弹窗组件

```vue
<!-- FileManager/components/DispatchDialog.vue -->

<template>
  <el-dialog v-model="visible" title="分发文件" width="500px">
    <div v-if="file" class="dispatch-form">
      <!-- 文件信息 -->
      <div class="file-info">
        <p><strong>文件名:</strong> {{ file.name }}</p>
        <p><strong>大小:</strong> {{ formatSize(file.size) }}</p>
      </div>

      <el-divider />

      <!-- 目标服务器 -->
      <el-form-item label="目标服务器">
        <div class="server-list">
          <el-checkbox
            v-for="server in servers"
            :key="server.id"
            v-model="selectedServers"
            :label="server.id"
            :disabled="server.status !== 'online'"
          >
            {{ server.name }} ({{ server.ip_address }})
            <el-tag size="small" :type="server.status === 'online' ? 'success' : 'info'">
              {{ server.status === 'online' ? '在线' : '离线' }}
            </el-tag>
          </el-checkbox>
        </div>
      </el-form-item>

      <!-- 目标路径 -->
      <el-form-item label="目标路径">
        <el-input v-model="targetPath" placeholder="/opt/app/" />
      </el-form-item>
    </div>

    <template #footer>
      <el-button @click="close">取消</el-button>
      <el-button type="primary" :loading="submitting" :disabled="selectedServers.length === 0" @click="doDispatch">
        开始分发
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { listServers, createDispatch } from '@/api/files'
import type { FileRecord, Server } from '@/types/files'

const props = defineProps<{
  visible: boolean
  file: FileRecord | null
}>()

const emit = defineEmits<{
  'update:visible': [value: boolean]
  success: []
}>()

const servers = ref<Server[]>([])
const selectedServers = ref<number[]>([])
const targetPath = ref('/opt/app/')
const submitting = ref(false)

// 加载服务器列表
async function loadServers() {
  const res = await listServers()
  servers.value = res.data
}

// 执行分发
async function doDispatch() {
  if (!props.file) return
  
  submitting.value = true
  try {
    await createDispatch({
      file_id: props.file.id,
      server_ids: selectedServers.value,
      target_path: targetPath.value,
    })
    emit('success')
    close()
  } catch (error: any) {
    console.error('分发失败', error)
  } finally {
    submitting.value = false
  }
}

// 关闭弹窗
function close() {
  selectedServers.value = []
  emit('update:visible', false)
}

// 格式化文件大小
function formatSize(bytes: number): string {
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  if (bytes < 1024 * 1024 * 1024) return (bytes / 1024 / 1024).toFixed(1) + ' MB'
  return (bytes / 1024 / 1024 / 1024).toFixed(2) + ' GB'
}

// 重置状态
watch(() => props.visible, (val) => {
  if (val) {
    selectedServers.value = []
    targetPath.value = '/opt/app/'
  }
})

onMounted(() => {
  loadServers()
})
</script>

<style scoped>
.dispatch-form {
  padding: 0 20px;
}

.file-info {
  background: #f5f7fa;
  padding: 12px;
  border-radius: 4px;
}

.file-info p {
  margin: 4px 0;
}

.server-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
</style>
```

## 10. 路由配置

```typescript
// front/src/router/index.ts

{
  path: '/files',
  name: 'FileManager',
  component: () => import('@/views/FileManager/index.vue'),
  meta: { title: '文件管理' }
}
```

## 11. 完成检查

- [ ] 类型定义完成
- [ ] API 接口封装完成
- [ ] 分片上传逻辑实现
- [ ] 主页面实现完成
- [ ] 文件列表组件完成
- [ ] 上传组件完成
- [ ] 分发弹窗完成
- [ ] 路由配置完成
- [ ] 功能测试通过

## 12. 下一阶段

完成本阶段后，文件管理核心功能已完整。如需查看分发状态、各节点文件状态，继续开发 [05-dispatch-management.md](./05-dispatch-management.md)。
