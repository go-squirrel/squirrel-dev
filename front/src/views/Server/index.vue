<template>
  <div class="server-page">
    <PageHeader :title="$t('server.listTitle')">
      <Button type="primary" @click="handleAdd">
        <Icon icon="lucide:plus" />
        {{ $t('server.addServer') }}
      </Button>
    </PageHeader>

    <Loading v-if="loading" :text="$t('server.loading')" />

    <Empty v-else-if="servers.length === 0" :description="$t('server.noServers')" icon="lucide:server">
      <template #action>
        <Button type="primary" @click="handleAdd">
          {{ $t('server.addFirstServer') }}
        </Button>
      </template>
    </Empty>

    <ServerTable
      v-else
      :servers="servers"
      @terminal="handleTerminal"
      @detail="handleDetail"
      @edit="handleEdit"
      @delete="handleDelete"
    />

    <ServerForm
      v-if="showForm"
      :server="editingServer"
      @submit="handleFormSubmit"
      @cancel="showForm = false"
    />

    <ServerDetail
      v-if="showDetail && selectedServer"
      :server="selectedServer"
      @close="showDetail = false"
    />

    <DeleteConfirm
      v-if="showDeleteConfirm && deletingServer"
      :server="deletingServer"
      @confirm="confirmDelete"
      @cancel="showDeleteConfirm = false"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { fetchServers } from '@/api/server'
import type { Server } from '@/types'
import PageHeader from '@/components/PageHeader/index.vue'
import Button from '@/components/Button/index.vue'
import Loading from '@/components/Loading/index.vue'
import Empty from '@/components/Empty/index.vue'
import ServerTable from './components/ServerTable.vue'
import ServerForm from './components/ServerForm.vue'
import ServerDetail from './components/ServerDetail.vue'
import DeleteConfirm from './components/DeleteConfirm.vue'
import { useLoading } from '@/composables/useLoading'

const router = useRouter()
const { loading, withLoading } = useLoading()

const servers = ref<Server[]>([])
const showForm = ref(false)
const showDetail = ref(false)
const showDeleteConfirm = ref(false)
const editingServer = ref<Server | null>(null)
const selectedServer = ref<Server | null>(null)
const deletingServer = ref<Server | null>(null)

const loadServers = async () => {
  await withLoading(async () => {
    servers.value = await fetchServers()
  })
}

const handleAdd = () => {
  editingServer.value = null
  showForm.value = true
}

const handleEdit = (server: Server) => {
  editingServer.value = server
  showForm.value = true
}

const handleDetail = (server: Server) => {
  selectedServer.value = server
  showDetail.value = true
}

const handleTerminal = (server: Server) => {
  router.push(`/servers/${server.id}/terminal`)
}

const handleDelete = (server: Server) => {
  deletingServer.value = server
  showDeleteConfirm.value = true
}

const handleFormSubmit = async () => {
  showForm.value = false
  await loadServers()
}

const confirmDelete = async () => {
  showDeleteConfirm.value = false
  await loadServers()
}

onMounted(() => {
  loadServers()
})
</script>

<style scoped>
.server-page {
  padding: 20px;
}
</style>
