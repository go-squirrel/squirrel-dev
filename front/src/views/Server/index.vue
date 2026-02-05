<template>
  <div class="server-page">
    <div class="page-header">
      <h2 class="page-title">{{ $t('server.listTitle') }}</h2>
      <button class="add-btn" @click="handleAdd">
        <Icon icon="lucide:plus" />
        {{ $t('server.addServer') }}
      </button>
    </div>

    <div v-if="loading" class="loading">
      {{ $t('server.loading') }}
    </div>

    <div v-else-if="servers.length === 0" class="empty-state">
      <Icon icon="lucide:server" class="empty-icon" />
      <p>{{ $t('server.noServers') }}</p>
      <button class="add-btn" @click="handleAdd">
        {{ $t('server.addFirstServer') }}
      </button>
    </div>

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
import { Icon } from '@iconify/vue'
import { fetchServers } from '@/api/server'
import type { Server } from '@/types'
import ServerTable from './components/ServerTable.vue'
import ServerForm from './components/ServerForm.vue'
import ServerDetail from './components/ServerDetail.vue'
import DeleteConfirm from './components/DeleteConfirm.vue'

const router = useRouter()

const servers = ref<Server[]>([])
const loading = ref(true)
const showForm = ref(false)
const showDetail = ref(false)
const showDeleteConfirm = ref(false)
const editingServer = ref<Server | null>(null)
const selectedServer = ref<Server | null>(null)
const deletingServer = ref<Server | null>(null)

const loadServers = async () => {
  loading.value = true
  try {
    servers.value = await fetchServers()
  } catch (error) {
    console.error('Failed to load servers:', error)
  } finally {
    loading.value = false
  }
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

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
}

.page-title {
  font-size: 18px;
  font-weight: 600;
  color: #1e3a5f;
}

.add-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: linear-gradient(135deg, #4fc3f7 0%, #29b6f6 100%);
  color: #ffffff;
  border: none;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
}

.add-btn:hover {
  box-shadow: 0 4px 12px rgba(79, 195, 247, 0.4);
  transform: translateY(-1px);
}

.loading {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 60px;
  color: #64748b;
  font-size: 14px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
  color: #64748b;
}

.empty-icon {
  width: 64px;
  height: 64px;
  color: #cbd5e1;
  margin-bottom: 16px;
}

.empty-state p {
  font-size: 14px;
  margin-bottom: 20px;
}
</style>
