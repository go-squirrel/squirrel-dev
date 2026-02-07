<template>
  <div class="configs-page">
    <PageHeader :title="$t('configs.listTitle')">
      <div class="header-actions">
        <div class="search-wrapper">
          <Icon icon="lucide:search" class="search-icon" />
          <input
            v-model="searchKeyword"
            :placeholder="$t('configs.searchPlaceholder')"
            class="search-input"
          />
          <button v-if="searchKeyword" class="clear-btn" @click="searchKeyword = ''">
            <Icon icon="lucide:x" />
          </button>
        </div>
        <Button type="primary" @click="handleAdd">
          <Icon icon="lucide:plus" />
          {{ $t('configs.addConfig') }}
        </Button>
      </div>
    </PageHeader>

    <Loading v-if="loading" :text="$t('common.loading')" />

    <Empty v-else-if="filteredConfigs.length === 0" :description="$t('common.noData')" icon="lucide:settings">
      <template #action>
        <Button type="primary" @click="handleAdd">
          {{ $t('configs.addConfig') }}
        </Button>
      </template>
    </Empty>

    <ConfigTable
      v-else
      :configs="filteredConfigs"
      :sort-by="sortBy"
      :sort-order="sortOrder"
      @edit="handleEdit"
      @delete="handleDelete"
      @copy="handleCopy"
      @sort="handleSort"
    />

    <ConfigForm
      v-if="showForm"
      :config="editingConfig"
      @submit="handleFormSubmit"
      @cancel="showForm = false"
    />

    <DeleteConfirm
      v-if="showDeleteConfirm && deletingConfig"
      :config="deletingConfig"
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
import { Icon } from '@iconify/vue'
import { useI18n } from 'vue-i18n'
import { fetchConfigs, createConfig, updateConfig, deleteConfig } from '@/api/config'
import type { Config, CreateConfigRequest, UpdateConfigRequest } from '@/types'
import PageHeader from '@/components/PageHeader/index.vue'
import Button from '@/components/Button/index.vue'
import Loading from '@/components/Loading/index.vue'
import Empty from '@/components/Empty/index.vue'
import ConfigTable from './components/ConfigTable.vue'
import ConfigForm from './components/ConfigForm.vue'
import DeleteConfirm from './components/DeleteConfirm.vue'
import Toast from '@/components/Toast/index.vue'
import { useLoading } from '@/composables/useLoading'

const { t } = useI18n()
const { loading, withLoading } = useLoading()

const configs = ref<Config[]>([])
const showForm = ref(false)
const showDeleteConfirm = ref(false)
const editingConfig = ref<Config | null>(null)
const deletingConfig = ref<Config | null>(null)
const searchKeyword = ref('')
const sortBy = ref<string | null>(null)
const sortOrder = ref<'asc' | 'desc'>('asc')
const toastVisible = ref(false)
const toastMessage = ref('')

const filteredConfigs = computed(() => {
  let result = configs.value

  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    result = result.filter(config =>
      config.key.toLowerCase().includes(keyword) ||
      config.value.toLowerCase().includes(keyword)
    )
  }

  if (sortBy.value) {
    result = [...result].sort((a, b) => {
      const aValue = a[sortBy.value as keyof Config]
      const bValue = b[sortBy.value as keyof Config]
      const comparison = String(aValue).localeCompare(String(bValue))
      return sortOrder.value === 'asc' ? comparison : -comparison
    })
  }

  return result
})

const loadConfigs = async () => {
  await withLoading(async () => {
    configs.value = await fetchConfigs()
  })
}

const handleAdd = () => {
  editingConfig.value = null
  showForm.value = true
}

const handleEdit = (config: Config) => {
  editingConfig.value = config
  showForm.value = true
}

const handleDelete = (config: Config) => {
  deletingConfig.value = config
  showDeleteConfirm.value = true
}

const handleCopy = async (value: string) => {
  try {
    await navigator.clipboard.writeText(value)
    toastMessage.value = t('configs.copySuccess')
    toastVisible.value = true
    setTimeout(() => {
      toastVisible.value = false
    }, 2000)
  } catch (err) {
    console.error('Failed to copy:', err)
  }
}

const handleSort = (field: string) => {
  if (sortBy.value === field) {
    sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc'
  } else {
    sortBy.value = field
    sortOrder.value = 'asc'
  }
}

const handleFormSubmit = async (data: CreateConfigRequest | UpdateConfigRequest) => {
  let success = false
  if (editingConfig.value) {
    try {
      await updateConfig(editingConfig.value.id, data as UpdateConfigRequest)
      success = true
      toastMessage.value = t('configs.updateSuccess')
    } catch (error) {
      console.error('Failed to update config:', error)
    }
  } else {
    try {
      await createConfig(data as CreateConfigRequest)
      success = true
      toastMessage.value = t('configs.createSuccess')
    } catch (error) {
      console.error('Failed to create config:', error)
    }
  }

  if (success) {
    showForm.value = false
    await loadConfigs()
    toastVisible.value = true
    setTimeout(() => {
      toastVisible.value = false
    }, 2000)
  }
}

const confirmDelete = async () => {
  if (!deletingConfig.value) return

  try {
    await deleteConfig(deletingConfig.value.id)
    showDeleteConfirm.value = false
    await loadConfigs()
    toastMessage.value = t('configs.deleteSuccess')
    toastVisible.value = true
    setTimeout(() => {
      toastVisible.value = false
    }, 2000)
  } catch (error) {
    console.error('Failed to delete config:', error)
  }
}

onMounted(() => {
  loadConfigs()
})
</script>

<style scoped>
.configs-page {
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
