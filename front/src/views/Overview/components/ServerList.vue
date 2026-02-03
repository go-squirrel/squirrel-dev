<template>
  <aside class="server-sidebar">
    <div class="sidebar-header">
      <h3>{{ $t('overview.serverList') }}</h3>
    </div>
    <div class="server-list">
      <div
        v-for="server in servers"
        :key="server.id"
        class="server-item"
        :class="{ active: currentServerId === server.id }"
        @click="$emit('switch', server.id)"
      >
        <div class="server-icon">
          <Icon icon="lucide:server" />
        </div>
        <div class="server-info">
          <span class="server-name">{{ server.hostname }}</span>
          <span class="server-ip">{{ server.ip_address }}</span>
        </div>
        <div class="server-status" :class="getStatusClass(server.status)"></div>
      </div>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue'
import { getStatusClass } from '@/utils/format'
import type { Server } from '@/types'

defineProps<{
  servers: Server[]
  currentServerId: number
}>()

defineEmits<{
  switch: [serverId: number]
}>()
</script>

<style scoped>
.server-sidebar {
  width: 220px;
  background: #ffffff;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.sidebar-header {
  padding: 16px;
  border-bottom: 1px solid #f1f5f9;
}

.sidebar-header h3 {
  font-size: 14px;
  font-weight: 600;
  color: #1e3a5f;
}

.server-list {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
}

.server-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.3s ease;
  margin-bottom: 4px;
}

.server-item:hover {
  background: #f8fafc;
}

.server-item.active {
  background: linear-gradient(135deg, #4fc3f7 0%, #29b6f6 100%);
}

.server-item.active .server-name,
.server-item.active .server-ip {
  color: #ffffff;
}

.server-icon {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f1f5f9;
  border-radius: 6px;
  color: #4fc3f7;
  font-size: 16px;
}

.server-item.active .server-icon {
  background: rgba(255, 255, 255, 0.2);
  color: #ffffff;
}

.server-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.server-name {
  font-size: 13px;
  font-weight: 500;
  color: #1e3a5f;
}

.server-ip {
  font-size: 11px;
  color: #94a3b8;
}

.server-status {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #94a3b8;
}

.server-status.online {
  background: #16a34a;
}

.server-status.offline {
  background: #dc2626;
}
</style>
