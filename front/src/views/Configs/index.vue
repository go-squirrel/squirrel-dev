<template>
  <div class="configs-page">
    <PageHeader :title="$t('configs.listTitle')">
      <Button type="primary" @click="handleAdd">
        <Icon icon="lucide:plus" />
        {{ $t('configs.addConfig') }}
      </Button>
    </PageHeader>

    <Loading v-if="loading" :text="$t('common.loading')" />

    <Empty v-else-if="configs.length === 0" :description="$t('common.noData')" icon="lucide:settings">
      <template #action>
        <Button type="primary" @click="handleAdd">
          {{ $t('configs.addConfig') }}
        </Button>
      </template>
    </Empty>

    <ConfigTable
      v-else
      :configs="configs"
      @edit="handleEdit"
      @delete="handleDelete"
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
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { fetchConfigs } from '@/api/config'
import type { Config } from '@/types'
import PageHeader from '@/components/PageHeader/index.vue'
import Button from '@/components/Button/index.vue'
import Loading from '@/components/Loading/index.vue'
import Empty from '@/components/Empty/index.vue'
import ConfigTable from './components/ConfigTable.vue'
import ConfigForm from './components/ConfigForm.vue'
import DeleteConfirm from './components/DeleteConfirm.vue'
import { useLoading } from '@/composables/useLoading'

const { loading, withLoading } = useLoading()

const configs = ref<Config[]>([])
const showForm = ref(false)
const showDeleteConfirm = ref(false)
const editingConfig = ref<Config | null>(null)
const deletingConfig = ref<Config | null>(null)

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

const handleFormSubmit = async () => {
  showForm.value = false
  await loadConfigs()
}

const confirmDelete = async () => {
  showDeleteConfirm.value = false
  await loadConfigs()
}

onMounted(() => {
  loadConfigs()
})
</script>

<style scoped>
.configs-page {
  padding: 20px;
}
</style>
