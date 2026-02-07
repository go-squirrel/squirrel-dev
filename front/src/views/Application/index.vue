<template>
  <div class="application-page">
    <PageHeader :title="$t('application.listTitle')">
      <div class="header-actions">
        <div class="search-wrapper">
          <Icon icon="lucide:search" class="search-icon" />
          <input
            v-model="searchKeyword"
            :placeholder="$t('application.searchPlaceholder')"
            class="search-input"
          />
          <button v-if="searchKeyword" class="clear-btn" @click="searchKeyword = ''">
            <Icon icon="lucide:x" />
          </button>
        </div>
        <Button type="secondary" @click="handleImportFromStore">
          <Icon icon="lucide:shopping-bag" />
          {{ $t('application.importFromStore') }}
        </Button>
        <Button type="primary" @click="handleAdd">
          <Icon icon="lucide:plus" />
          {{ $t('application.addApplication') }}
        </Button>
      </div>
    </PageHeader>

    <Loading v-if="loading" :text="$t('common.loading')" />

    <Empty v-else-if="filteredApplications.length === 0" :description="$t('common.noData')" icon="lucide:package">
      <template #action>
        <Button type="primary" @click="handleAdd">
          {{ $t('application.addApplication') }}
        </Button>
      </template>
    </Empty>

    <ApplicationTable
      v-else
      :applications="filteredApplications"
      :sort-by="sortBy"
      :sort-order="sortOrder"
      @edit="handleEdit"
      @delete="handleDelete"
      @view="handleView"
      @sort="handleSort"
    />

    <ApplicationForm
      v-if="showForm"
      :application="editingApplication"
      @submit="handleFormSubmit"
      @cancel="showForm = false"
    />

    <ApplicationDetail
      v-if="showDetail && viewingApplication"
      :application="viewingApplication"
      @close="showDetail = false"
      @edit="handleEditFromDetail"
    />

    <DeleteConfirm
      v-if="showDeleteConfirm && deletingApplication"
      :application="deletingApplication"
      @confirm="confirmDelete"
      @cancel="showDeleteConfirm = false"
    />

    <Toast
      :visible="toastVisible"
      :message="toastMessage"
      type="success"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Icon } from '@iconify/vue'
import { useI18n } from 'vue-i18n'
import { fetchApplications, createApplication, updateApplication, deleteApplication } from '@/api/application'
import type { ApplicationInstance, CreateApplicationRequest, UpdateApplicationRequest } from '@/types'
import PageHeader from '@/components/PageHeader/index.vue'
import Button from '@/components/Button/index.vue'
import Loading from '@/components/Loading/index.vue'
import Empty from '@/components/Empty/index.vue'
import ApplicationTable from './components/ApplicationTable.vue'
import ApplicationForm from './components/ApplicationForm.vue'
import ApplicationDetail from './components/ApplicationDetail.vue'
import DeleteConfirm from './components/DeleteConfirm.vue'
import Toast from '@/components/Toast/index.vue'
import { useLoading } from '@/composables/useLoading'

const { t } = useI18n()
const router = useRouter()
const { loading, withLoading } = useLoading()

const applications = ref<ApplicationInstance[]>([])
const showForm = ref(false)
const showDeleteConfirm = ref(false)
const showDetail = ref(false)
const editingApplication = ref<ApplicationInstance | null>(null)
const deletingApplication = ref<ApplicationInstance | null>(null)
const viewingApplication = ref<ApplicationInstance | null>(null)
const searchKeyword = ref('')
const sortBy = ref<string | null>(null)
const sortOrder = ref<'asc' | 'desc'>('asc')
const toastVisible = ref(false)
const toastMessage = ref('')

const filteredApplications = computed(() => {
  let result = applications.value

  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    result = result.filter(app => 
      app.name.toLowerCase().includes(keyword) ||
      app.description.toLowerCase().includes(keyword)
    )
  }

  if (sortBy.value) {
    result = [...result].sort((a, b) => {
      const aValue = a[sortBy.value as keyof ApplicationInstance]
      const bValue = b[sortBy.value as keyof ApplicationInstance]
      const comparison = String(aValue).localeCompare(String(bValue))
      return sortOrder.value === 'asc' ? comparison : -comparison
    })
  }

  return result
})

const loadApplications = async () => {
  await withLoading(async () => {
    applications.value = await fetchApplications()
  })
}

const handleAdd = () => {
  editingApplication.value = null
  showForm.value = true
}

const handleImportFromStore = () => {
  router.push('/app-store')
}

const handleEdit = (application: ApplicationInstance) => {
  editingApplication.value = application
  showForm.value = true
}

const handleEditFromDetail = (application: ApplicationInstance) => {
  showDetail.value = false
  editingApplication.value = application
  showForm.value = true
}

const handleView = (application: ApplicationInstance) => {
  viewingApplication.value = application
  showDetail.value = true
}

const handleDelete = (application: ApplicationInstance) => {
  deletingApplication.value = application
  showDeleteConfirm.value = true
}

const handleSort = (field: string) => {
  if (sortBy.value === field) {
    sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc'
  } else {
    sortBy.value = field
    sortOrder.value = 'asc'
  }
}

const handleFormSubmit = async (data: CreateApplicationRequest | UpdateApplicationRequest) => {
  let success = false
  if (editingApplication.value) {
    try {
      await updateApplication(editingApplication.value.id, data as UpdateApplicationRequest)
      success = true
      toastMessage.value = t('application.updateSuccess')
    } catch (error) {
      console.error('Failed to update application:', error)
    }
  } else {
    try {
      await createApplication(data as CreateApplicationRequest)
      success = true
      toastMessage.value = t('application.createSuccess')
    } catch (error) {
      console.error('Failed to create application:', error)
    }
  }

  if (success) {
    showForm.value = false
    await loadApplications()
    toastVisible.value = true
    setTimeout(() => {
      toastVisible.value = false
    }, 2000)
  }
}

const confirmDelete = async () => {
  if (!deletingApplication.value) return

  try {
    await deleteApplication(deletingApplication.value.id)
    showDeleteConfirm.value = false
    await loadApplications()
    toastMessage.value = t('application.deleteSuccess')
    toastVisible.value = true
    setTimeout(() => {
      toastVisible.value = false
    }, 2000)
  } catch (error) {
    console.error('Failed to delete application:', error)
  }
}

onMounted(() => {
  loadApplications()
})
</script>

<style scoped>
.application-page {
  padding: 20px;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.search-wrapper {
  display: flex;
  align-items: center;
  position: relative;
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

.search-icon {
  position: absolute;
  left: 12px;
  width: 16px;
  height: 16px;
  color: #94a3b8;
  pointer-events: none;
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
