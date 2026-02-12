<template>
  <div class="server-table-container">
    <table class="server-table">
      <thead>
        <tr>
          <th>{{ $t('server.hostname') }}</th>
          <th>{{ $t('server.ipAddress') }}</th>
          <th>{{ $t('server.sshPort') }}</th>
          <th>{{ $t('server.username') }}</th>
          <th>{{ $t('server.authType') }}</th>
          <th>{{ $t('server.status') }}</th>
          <th>{{ $t('server.operation') }}</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="server in servers" :key="server.id" class="server-row">
          <td class="hostname-cell">
            <div class="hostname-wrapper">
              <Icon icon="lucide:server" class="server-icon" />
              <div class="hostname-text">
                <div class="hostname">{{ server.hostname }}</div>
                <div v-if="server.server_alias" class="alias">{{ server.server_alias }}</div>
              </div>
            </div>
          </td>
          <td class="ip-cell">{{ server.ip_address }}</td>
          <td class="port-cell">{{ server.ssh_port }}</td>
          <td class="username-cell">{{ server.ssh_username }}</td>
          <td class="auth-cell">
            <span class="auth-badge" :class="server.auth_type">
              {{ server.auth_type === 'password' ? $t('server.password') : $t('server.key') }}
            </span>
          </td>
          <td class="status-cell">
            <span class="status-badge" :class="server.status">
              <span class="status-dot"></span>
              {{ getStatusText(server.status) }}
            </span>
          </td>
          <td class="action-cell">
            <button class="action-btn terminal-btn" :title="$t('server.connectTerminal')" @click="$emit('terminal', server)">
              <Icon icon="lucide:terminal" />
            </button>
            <button class="action-btn detail-btn" :title="$t('server.viewDetail')" @click="$emit('detail', server)">
              <Icon icon="lucide:info" />
            </button>
            <button class="action-btn edit-btn" :title="$t('server.editServer')" @click="$emit('edit', server)">
              <Icon icon="lucide:edit-2" />
            </button>
            <button class="action-btn delete-btn" :title="$t('server.deleteServer')" @click="$emit('delete', server)">
              <Icon icon="lucide:trash-2" />
            </button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import type { Server } from '@/types'

defineProps<{
  servers: Server[]
}>()

defineEmits<{
  terminal: [server: Server]
  detail: [server: Server]
  edit: [server: Server]
  delete: [server: Server]
}>()

const { t } = useI18n()

const getStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    online: t('server.online'),
    offline: t('server.offline'),
    unknown: t('server.unknown'),
    active: t('server.active'),
    inactive: t('server.inactive')
  }
  return statusMap[status] || status
}
</script>

<style scoped>
.server-table-container {
  background: #ffffff;
  border-radius: 8px;
  overflow: hidden;
}

.server-table {
  width: 100%;
  border-collapse: collapse;
}

.server-table thead {
  background: #f8fafc;
}

.server-table th {
  padding: 12px 16px;
  text-align: left;
  font-weight: 600;
  font-size: 13px;
  color: #475569;
  border-bottom: 2px solid #e2e8f0;
}

.server-row {
  border-bottom: 1px solid #f1f5f9;
  transition: background 0.2s ease;
}

.server-row:hover {
  background: #f8fafc;
}

.server-row:last-child {
  border-bottom: none;
}

.server-table td {
  padding: 12px 16px;
  font-size: 13px;
  color: #1e3a5f;
}

.hostname-cell {
  min-width: 180px;
}

.hostname-wrapper {
  display: flex;
  align-items: center;
  gap: 8px;
}

.server-icon {
  width: 16px;
  height: 16px;
  color: #64748b;
}

.hostname-text {
  display: flex;
  flex-direction: column;
}

.hostname {
  font-weight: 500;
}

.alias {
  font-size: 11px;
  color: #94a3b8;
  margin-top: 2px;
}

.ip-cell {
  font-family: 'SF Mono', Monaco, Consolas, monospace;
  color: #64748b;
}

.port-cell {
  color: #64748b;
}

.username-cell {
  color: #64748b;
}

.auth-cell {
  min-width: 80px;
}

.auth-badge {
  display: inline-flex;
  align-items: center;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.auth-badge.password {
  background: #dbeafe;
  color: #1e40af;
}

.auth-badge.key {
  background: #fef3c7;
  color: #92400e;
}

.status-cell {
  min-width: 80px;
}

.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 11px;
  font-weight: 500;
}

.status-badge.online,
.status-badge.active {
  background: #dcfce7;
  color: #16a34a;
}

.status-badge.offline,
.status-badge.inactive {
  background: #fee2e2;
  color: #dc2626;
}

.status-badge.unknown {
  background: #f1f5f9;
  color: #64748b;
}

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: currentColor;
}

.action-cell {
  width: 140px;
  display: flex;
  align-items: center;
  gap: 4px;
}

.action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border: none;
  background: transparent;
  color: #64748b;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.action-btn:hover {
  background: #f1f5f9;
  color: #475569;
}

.action-btn.terminal-btn:hover {
  background: #dcfce7;
  color: #16a34a;
}

.action-btn.detail-btn:hover {
  background: #e0f2fe;
  color: #0284c7;
}

.action-btn.edit-btn:hover {
  background: #fef3c7;
  color: #92400e;
}

.action-btn.delete-btn:hover {
  background: #fee2e2;
  color: #dc2626;
}
</style>
