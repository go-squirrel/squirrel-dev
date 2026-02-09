<template>
  <div class="deployment-table-container">
    <table class="deployment-table">
      <thead>
        <tr>
          <th>{{ $t('deployment.deployId') }}</th>
          <th>{{ $t('deployment.applicationName') }}</th>
          <th>{{ $t('deployment.applicationType') }}</th>
          <th>{{ $t('deployment.applicationVersion') }}</th>
          <th>{{ $t('deployment.serverAddress') }}</th>
          <th>{{ $t('deployment.status') }}</th>
          <th>{{ $t('deployment.deployedAt') }}</th>
          <th>{{ $t('common.operation') }}</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="deployment in deployments" :key="deployment.id" class="deployment-row">
          <td class="deploy-id-cell">
            <span class="deploy-id-text">{{ deployment.deploy_id }}</span>
          </td>
          <td class="name-cell">
            <div class="name-wrapper">
              <Icon :icon="getTypeIcon(deployment.application.type)" class="type-icon" />
              <span class="name-text">{{ deployment.application.name }}</span>
            </div>
          </td>
          <td class="type-cell">
            <span class="type-badge" :class="`type-badge--${deployment.application.type}`">
              {{ getTypeLabel(deployment.application.type) }}
            </span>
          </td>
          <td class="version-cell">
            <span class="version-text">{{ deployment.application.version }}</span>
          </td>
          <td class="server-cell">
            <div class="server-wrapper">
              <Icon icon="lucide:server" class="server-icon" />
              <div class="server-info">
                <span class="server-ip">{{ deployment.server.ip_address }}</span>
                <span class="agent-port">:{{ deployment.server.agent_port || 8080 }}</span>
              </div>
            </div>
          </td>
          <td class="status-cell">
            <span class="status-badge" :class="`status-badge--${deployment.status}`">
              {{ getStatusLabel(deployment.status) }}
            </span>
          </td>
          <td class="deployed-at-cell">
            <span class="deployed-at-text">{{ deployment.deployed_at }}</span>
          </td>
          <td class="action-cell">
            <button class="action-btn view-btn" :title="$t('deployment.viewDetail')" @click="$emit('view', deployment)">
              <Icon icon="lucide:eye" />
            </button>
            <button
              v-if="deployment.status === 'running'"
              class="action-btn stop-btn"
              :title="$t('deployment.stop')"
              @click="$emit('stop', deployment)"
            >
              <Icon icon="lucide:square" />
            </button>
            <button
              v-else
              class="action-btn start-btn"
              :title="$t('deployment.start')"
              @click="$emit('start', deployment)"
            >
              <Icon icon="lucide:play" />
            </button>
            <button class="action-btn undeploy-btn" :title="$t('deployment.undeploy')" @click="$emit('undeploy', deployment)">
              <Icon icon="lucide:trash-2" />
            </button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue'
import type { Deployment } from '../types'

defineProps<{
  deployments: Deployment[]
}>()

defineEmits<{
  view: [deployment: Deployment]
  start: [deployment: Deployment]
  stop: [deployment: Deployment]
  undeploy: [deployment: Deployment]
}>()

const getTypeIcon = (type: string) => {
  const icons: Record<string, string> = {
    compose: 'lucide:box',
    script: 'lucide:terminal',
    binary: 'lucide:file-code'
  }
  return icons[type] || 'lucide:package'
}

const getTypeLabel = (type: string) => {
  const labels: Record<string, string> = {
    compose: 'Compose',
    script: '脚本',
    binary: '二进制'
  }
  return labels[type] || type
}

const getStatusLabel = (status: string) => {
  const labels: Record<string, string> = {
    running: '运行中',
    stopped: '已停止',
    not_deployed: '未部署',
    error: '错误'
  }
  return labels[status] || status
}
</script>

<style scoped>
.deployment-table-container {
  background: #ffffff;
  border-radius: 8px;
  overflow: hidden;
}

.deployment-table {
  width: 100%;
  border-collapse: collapse;
}

.deployment-table thead {
  background: #f8fafc;
}

.deployment-table th {
  padding: 12px 16px;
  text-align: left;
  font-weight: 600;
  font-size: 13px;
  color: #475569;
  border-bottom: 2px solid #e2e8f0;
}

.deployment-row {
  border-bottom: 1px solid #f1f5f9;
  transition: background 0.2s ease;
}

.deployment-row:hover {
  background: #f8fafc;
}

.deployment-row:last-child {
  border-bottom: none;
}

.deployment-table td {
  padding: 12px 16px;
  font-size: 13px;
  color: #1e3a5f;
}

.deploy-id-cell {
  width: 100px;
}

.deploy-id-text {
  font-family: 'Monaco', 'Menlo', monospace;
  font-size: 12px;
  color: #64748b;
}

.name-cell {
  width: 180px;
}

.name-wrapper {
  display: flex;
  align-items: center;
  gap: 8px;
}

.type-icon {
  width: 16px;
  height: 16px;
  color: #64748b;
}

.name-text {
  font-weight: 500;
}

.type-cell {
  width: 100px;
}

.type-badge {
  display: inline-block;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.type-badge--compose {
  background: #dbeafe;
  color: #1e40af;
}

.type-badge--script {
  background: #dcfce7;
  color: #166534;
}

.type-badge--binary {
  background: #fef3c7;
  color: #92400e;
}

.version-cell {
  width: 100px;
}

.version-text {
  font-family: 'Monaco', 'Menlo', monospace;
  font-size: 12px;
  color: #64748b;
}

.server-cell {
  width: 180px;
}

.server-wrapper {
  display: flex;
  align-items: center;
  gap: 6px;
}

.server-icon {
  width: 14px;
  height: 14px;
  color: #64748b;
}

.server-info {
  display: flex;
  align-items: center;
}

.server-ip {
  font-family: 'Monaco', 'Menlo', monospace;
  font-size: 12px;
  color: #1e3a5f;
}

.agent-port {
  font-family: 'Monaco', 'Menlo', monospace;
  font-size: 12px;
  color: #64748b;
}

.status-cell {
  width: 100px;
}

.status-badge {
  display: inline-block;
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.status-badge--running {
  background: #dcfce7;
  color: #166534;
}

.status-badge--stopped {
  background: #f1f5f9;
  color: #64748b;
}

.status-badge--not_deployed {
  background: #fef3c7;
  color: #92400e;
}

.status-badge--error {
  background: #fee2e2;
  color: #dc2626;
}

.deployed-at-cell {
  width: 140px;
}

.deployed-at-text {
  font-size: 12px;
  color: #64748b;
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

.action-btn.view-btn:hover {
  background: #e0f2fe;
  color: #0284c7;
}

.action-btn.start-btn:hover {
  background: #dcfce7;
  color: #16a34a;
}

.action-btn.stop-btn:hover {
  background: #fef3c7;
  color: #ca8a04;
}

.action-btn.undeploy-btn:hover {
  background: #fee2e2;
  color: #dc2626;
}
</style>
