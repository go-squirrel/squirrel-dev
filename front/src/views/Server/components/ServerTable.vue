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
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
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
  font-size: 12px;
  font-weight: 600;
  color: #64748b;
  text-transform: uppercase;
  letter-spacing: 0.5px;
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
  padding: 14px 16px;
  font-size: 13px;
  color: #1e3a5f;
}

.hostname-cell {
  min-width: 180px;
}

.hostname-wrapper {
  display: flex;
  align-items: center;
  gap: 10px;
}

.server-icon {
  width: 18px;
  height: 18px;
  color: #4fc3f7;
  flex-shrink: 0;
}

.hostname-text {
  display: flex;
  flex-direction: column;
}

.hostname {
  font-weight: 500;
  color: #1e3a5f;
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
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 11px;
  font-weight: 500;
}

.auth-badge.password {
  background: #e0f2fe;
  color: #0284c7;
}

.auth-badge.key {
  background: #fef3c7;
  color: #d97706;
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
  min-width: 140px;
}

.action-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 6px;
  border: none;
  cursor: pointer;
  transition: all 0.2s ease;
  margin-right: 4px;
}

.action-btn:last-child {
  margin-right: 0;
}

.terminal-btn {
  background: #f0fdf4;
  color: #16a34a;
}

.terminal-btn:hover {
  background: #dcfce7;
  transform: scale(1.1);
}

.detail-btn {
  background: #eff6ff;
  color: #0284c7;
}

.detail-btn:hover {
  background: #dbeafe;
  transform: scale(1.1);
}

.edit-btn {
  background: #fef9c3;
  color: #ca8a04;
}

.edit-btn:hover {
  background: #fef08a;
  transform: scale(1.1);
}

.delete-btn {
  background: #fef2f2;
  color: #dc2626;
}

.delete-btn:hover {
  background: #fee2e2;
  transform: scale(1.1);
}
</style>
