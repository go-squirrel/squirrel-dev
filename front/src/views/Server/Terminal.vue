<template>
  <div class="terminal-page">
    <!-- 左侧服务器列表侧边栏 -->
    <div class="server-sidebar">
      <div class="sidebar-header">
        <Icon icon="lucide:server" class="sidebar-icon" />
        <span class="sidebar-title">{{ $t('server.serverList') }}</span>
      </div>
      <div class="server-search">
        <Icon icon="lucide:search" class="search-icon" />
        <input
          v-model="searchKeyword"
          :placeholder="$t('server.searchPlaceholder')"
          class="search-input"
        />
      </div>
      <div class="server-list">
        <Loading v-if="loading" :text="$t('server.loading')" />
        <Empty v-else-if="filteredServers.length === 0" :description="$t('server.noServers')" icon="lucide:server" />
        <template v-else>
          <div
            v-for="srv in filteredServers"
            :key="srv.id"
            class="server-item"
            :class="{ active: selectedServer?.id === srv.id }"
            @click="selectServer(srv)"
          >
            <div class="server-info">
              <Icon icon="lucide:monitor" class="server-icon" />
              <div class="server-detail">
                <span class="server-hostname">{{ srv.hostname || srv.ip_address }}</span>
                <span class="server-ip">{{ srv.ip_address }}</span>
              </div>
            </div>
            <Icon
              v-if="selectedServer?.id === srv.id"
              icon="lucide:check-circle"
              class="check-icon"
            />
          </div>
        </template>
      </div>
    </div>

    <!-- 右侧 Terminal 区域 -->
    <div class="terminal-area">
      <div v-if="!selectedServer" class="empty-state">
        <Icon icon="lucide:terminal" class="empty-icon" />
        <p>{{ $t('server.selectServer') }}</p>
      </div>
      <TerminalComponent
        v-else
        :key="selectedServer.id"
        :server="selectedServer"
        @close="handleClose"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { fetchServers } from '@/api/server'
import type { Server } from '@/types'
import TerminalComponent from '@/components/terminal/index.vue'
import Loading from '@/components/Loading/index.vue'
import Empty from '@/components/Empty/index.vue'

const route = useRoute()
const router = useRouter()

const servers = ref<Server[]>([])
const loading = ref(false)
const searchKeyword = ref('')
const selectedServer = ref<Server | null>(null)

const filteredServers = computed(() => {
  if (!searchKeyword.value) {
    return servers.value
  }
  const keyword = searchKeyword.value.toLowerCase()
  return servers.value.filter(srv =>
    srv.hostname.toLowerCase().includes(keyword) ||
    srv.ip_address.toLowerCase().includes(keyword) ||
    (srv.server_alias && srv.server_alias.toLowerCase().includes(keyword))
  )
})

const loadServers = async () => {
  loading.value = true
  try {
    servers.value = await fetchServers()
    // 如果 URL 中有 id 参数，自动选中对应服务器
    const serverId = Number(route.params.id) || Number(route.query.id)
    if (serverId) {
      const server = servers.value.find(s => s.id === serverId)
      if (server) {
        selectedServer.value = server
      }
    }
  } catch (error) {
    console.error('Failed to load servers:', error)
  } finally {
    loading.value = false
  }
}

const selectServer = (server: Server) => {
  selectedServer.value = server
  // 更新 URL，但不刷新页面
  router.replace({
    path: `/servers/${server.id}/terminal`,
    replace: true
  })
}

const handleClose = () => {
  router.push('/servers')
}

onMounted(() => {
  loadServers()
})
</script>

<style scoped>
.terminal-page {
  display: flex;
  width: 100%;
  height: 100%;
  overflow: hidden;
}

/* 侧边栏样式 */
.server-sidebar {
  width: 280px;
  min-width: 280px;
  background: #ffffff;
  border-right: 1px solid #e2e8f0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.sidebar-header {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 16px;
  border-bottom: 1px solid #e2e8f0;
}

.sidebar-icon {
  width: 20px;
  height: 20px;
  color: #4fc3f7;
}

.sidebar-title {
  font-size: 15px;
  font-weight: 600;
  color: #1e3a5f;
}

.server-search {
  display: flex;
  align-items: center;
  position: relative;
  padding: 12px 16px;
  border-bottom: 1px solid #e2e8f0;
}

.server-search .search-icon {
  position: absolute;
  left: 28px;
  width: 16px;
  height: 16px;
  color: #94a3b8;
  pointer-events: none;
}

.server-search .search-input {
  width: 100%;
  padding: 8px 12px 8px 36px;
  border: 2px solid #e2e8f0;
  border-radius: 6px;
  font-size: 13px;
  color: #1e3a5f;
  background: #f8fafc;
  transition: all 0.2s ease;
}

.server-search .search-input:focus {
  outline: none;
  border-color: #4fc3f7;
  background: #ffffff;
}

.server-list {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
}

.server-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  margin-bottom: 4px;
}

.server-item:hover {
  background: #f1f5f9;
}

.server-item.active {
  background: #e0f2fe;
}

.server-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.server-icon {
  width: 20px;
  height: 20px;
  color: #64748b;
}

.server-item.active .server-icon {
  color: #0284c7;
}

.server-detail {
  display: flex;
  flex-direction: column;
}

.server-hostname {
  font-size: 14px;
  font-weight: 500;
  color: #1e3a5f;
}

.server-ip {
  font-size: 12px;
  color: #94a3b8;
}

.check-icon {
  width: 18px;
  height: 18px;
  color: #10b981;
}

/* Terminal 区域样式 */
.terminal-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background: #1e1e1e;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #a0a0a0;
  gap: 16px;
}

.empty-icon {
  width: 64px;
  height: 64px;
  color: #4fc3f7;
}

.empty-state p {
  font-size: 14px;
}

:deep(.terminal-fullscreen) {
  position: static !important;
  height: 100% !important;
}
</style>
