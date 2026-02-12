<template>
  <div class="app-store-page">
    <PageHeader :title="$t('appStore.title')" :subtitle="$t('appStore.listTitle')">
      <div class="header-actions">
        <div class="search-box">
          <Icon icon="lucide:search" class="search-icon" />
          <input
            v-model="searchKeyword"
            type="text"
            class="search-input"
            :placeholder="$t('appStore.searchPlaceholder')"
          />
          <button v-if="searchKeyword" class="clear-btn" @click="searchKeyword = ''">
            <Icon icon="lucide:x" />
          </button>
        </div>
        <CategoryFilter
          v-model:category="selectedCategory"
          v-model:type="selectedType"
        />
        <Button type="primary" icon="lucide:plus" @click="handleAdd">
          {{ $t('appStore.addApp') }}
        </Button>
      </div>
    </PageHeader>

    <Loading v-if="loading" :text="$t('common.loading')" />
    <AppTable
      v-else
      :apps="filteredApps"
      @detail="handleDetail"
      @edit="handleEdit"
      @delete="handleDelete"
      @import="handleImport"
    />

    <!-- 表单弹窗 -->
    <AppForm
      v-if="showForm"
      :app="editingApp"
      @submit="handleFormSubmit"
      @cancel="closeForm"
    />

    <!-- 详情弹窗 -->
    <AppDetail
      v-if="showDetail"
      :app="detailApp"
      @close="closeDetail"
      @download="handleDownload"
      @import="handleImport"
    />

    <!-- 删除确认弹窗 -->
    <DeleteConfirm
      v-if="showDeleteConfirm"
      :app="deletingApp"
      @confirm="confirmDelete"
      @cancel="closeDeleteConfirm"
    />

    <!-- Toast 提示 -->
    <Toast
      v-if="toastVisible"
      v-model:visible="toastVisible"
      :message="toastMessage"
      :type="toastType"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import PageHeader from '@/components/PageHeader/index.vue'
import Button from '@/components/Button/index.vue'
import Loading from '@/components/Loading/index.vue'
import Toast from '@/components/Toast/index.vue'
import AppTable from './components/AppTable.vue'
import AppForm from './components/AppForm.vue'
import AppDetail from './components/AppDetail.vue'
import DeleteConfirm from './components/DeleteConfirm.vue'
import CategoryFilter from './components/CategoryFilter.vue'
import { useAppStore } from './composables/useAppStore'
import { createApplication } from '@/api/application'
import type { AppStore, CreateAppRequest, UpdateAppRequest } from '@/types'

const { t } = useI18n()
const router = useRouter()

// 搜索和筛选
const searchKeyword = ref('')
const selectedCategory = ref('')
const selectedType = ref('')

const {
  loading,
  filteredApps,
  loadApps,
  addApp,
  editApp,
  removeApp
} = useAppStore(searchKeyword, selectedCategory, selectedType)

// 表单状态
const showForm = ref(false)
const editingApp = ref<AppStore | null>(null)

// 详情状态
const showDetail = ref(false)
const detailApp = ref<AppStore | null>(null)

// 删除确认状态
const showDeleteConfirm = ref(false)
const deletingApp = ref<AppStore | null>(null)

// Toast 状态
const toastVisible = ref(false)
const toastMessage = ref('')
const toastType = ref<'success' | 'error'>('success')

// 显示 Toast
const showToast = (message: string, type: 'success' | 'error' = 'success') => {
  toastMessage.value = message
  toastType.value = type
  toastVisible.value = true
  setTimeout(() => {
    toastVisible.value = false
  }, 3000)
}

// 打开添加表单
const handleAdd = () => {
  editingApp.value = null
  showForm.value = true
}

// 打开编辑表单
const handleEdit = (app: AppStore) => {
  editingApp.value = app
  showForm.value = true
}

// 关闭表单
const closeForm = () => {
  showForm.value = false
  editingApp.value = null
}

// 表单提交
const handleFormSubmit = async (data: CreateAppRequest | UpdateAppRequest) => {
  try {
    if (editingApp.value) {
      await editApp(editingApp.value.id, data as UpdateAppRequest)
      showToast(t('appStore.updateSuccess'))
    } else {
      await addApp(data as CreateAppRequest)
      showToast(t('appStore.createSuccess'))
    }
    closeForm()
    await loadApps()
  } catch (error) {
    showToast(t('appStore.operationFailed'), 'error')
  }
}

// 查看详情
const handleDetail = (app: AppStore) => {
  detailApp.value = app
  showDetail.value = true
}

// 关闭详情
const closeDetail = () => {
  showDetail.value = false
  detailApp.value = null
}

// 下载应用
const handleDownload = (app: AppStore | null) => {
  if (!app) return
  // 创建下载
  const blob = new Blob([app.content], { type: 'text/yaml' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `${app.name}-${app.version}.yaml`
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
}

// 导入应用到应用管理
const handleImport = async (app: AppStore | null) => {
  if (!app) return
  try {
    // 转换类型：AppStore.type -> ApplicationInstance.type
    const applicationData = {
      name: app.name,
      description: app.description,
      type: app.type,
      content: app.content,
      version: app.version
    }
    await createApplication(applicationData)
    showToast(t('appStore.importSuccess'), 'success')
    // 跳转到应用管理页面
    router.push('/applications')
  } catch (error) {
    console.error('Failed to import application:', error)
    showToast(t('appStore.operationFailed'), 'error')
  }
}

// 打开删除确认
const handleDelete = (app: AppStore) => {
  deletingApp.value = app
  showDeleteConfirm.value = true
}

// 关闭删除确认
const closeDeleteConfirm = () => {
  showDeleteConfirm.value = false
  deletingApp.value = null
}

// 确认删除
const confirmDelete = async () => {
  if (!deletingApp.value) return
  try {
    await removeApp(deletingApp.value.id)
    showToast(t('appStore.deleteSuccess'))
    closeDeleteConfirm()
    await loadApps()
  } catch (error) {
    showToast(t('appStore.operationFailed'), 'error')
  }
}

// 初始化加载
onMounted(() => {
  loadApps()
})
</script>

<style scoped>
.app-store-page {
  padding: 20px;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.search-box {
  position: relative;
  display: flex;
  align-items: center;
}

.search-icon {
  position: absolute;
  left: 12px;
  width: 16px;
  height: 16px;
  color: #94a3b8;
  pointer-events: none;
}

.search-input {
  width: 280px;
  padding: 8px 12px 8px 36px;
  border: 2px solid #e2e8f0;
  border-radius: 6px;
  font-size: 13px;
  color: #1e3a5f;
  background: #f8fafc;
  transition: all 0.2s ease;
}

.search-input:focus {
  outline: none;
  border-color: #4fc3f7;
  background: #ffffff;
}

.clear-btn {
  position: absolute;
  right: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 20px;
  height: 20px;
  border: none;
  background: transparent;
  color: #94a3b8;
  cursor: pointer;
  border-radius: 4px;
  transition: all 0.2s ease;
}

.clear-btn:hover {
  background: #f1f5f9;
  color: #64748b;
}
</style>
