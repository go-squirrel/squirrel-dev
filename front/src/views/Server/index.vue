<template>
  <div class="server-page">
    <PageHeader :title="$t('server.listTitle')">
      <div class="header-actions">
        <div class="search-wrapper">
          <Icon icon="lucide:search" class="search-icon" />
          <input
            v-model="searchKeyword"
            :placeholder="$t('server.searchPlaceholder')"
            class="search-input"
          />
          <button v-if="searchKeyword" class="clear-btn" @click="searchKeyword = ''">
            <Icon icon="lucide:x" />
          </button>
        </div>
        <Button type="primary" @click="handleAdd">
          <Icon icon="lucide:plus" />
          {{ $t('server.addServer') }}
        </Button>
      </div>
    </PageHeader>

    <Loading v-if="loading" :text="$t('server.loading')" />

    <Empty v-else-if="filteredServers.length === 0" :description="$t('server.noServers')" icon="lucide:server">
      <template #action>
        <Button type="primary" @click="handleAdd">
          {{ $t('server.addFirstServer') }}
        </Button>
      </template>
    </Empty>

    <ServerTable
      v-else
      :servers="filteredServers"
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
import { ref, computed, onMounted } from 'vue'
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
const searchKeyword = ref('')

const filteredServers = computed(() => {
  if (!searchKeyword.value) {
    return servers.value
  }
  const keyword = searchKeyword.value.toLowerCase()
  return servers.value.filter(server =>
    server.hostname.toLowerCase().includes(keyword) ||
    server.ip_address.toLowerCase().includes(keyword) ||
    (server.server_alias && server.server_alias.toLowerCase().includes(keyword))
  )
})

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
