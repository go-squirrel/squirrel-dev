<template>
  <Modal :model-value="visible" :title="$t('deployment.deploymentDetail')" width="600px" show-footer @close="handleClose">
    <div v-if="deployment" class="detail-content">
      <div class="detail-section">
        <h3 class="section-title">{{ $t('deployment.basicInfo') }}</h3>
        <div class="info-grid">
          <div class="info-item">
            <span class="info-label">{{ $t('deployment.deployId') }}</span>
            <span class="info-value">{{ deployment.deploy_id }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">{{ $t('deployment.status') }}</span>
            <span class="info-value">
              <span class="status-badge" :class="`status-badge--${deployment.status}`">
                {{ getStatusLabel(deployment.status) }}
              </span>
            </span>
          </div>
          <div class="info-item">
            <span class="info-label">{{ $t('deployment.deployedAt') }}</span>
            <span class="info-value">{{ deployment.deployed_at }}</span>
          </div>
        </div>
      </div>

      <div class="detail-section">
        <h3 class="section-title">{{ $t('deployment.applicationInfo') }}</h3>
        <div class="info-grid">
          <div class="info-item">
            <span class="info-label">{{ $t('deployment.applicationName') }}</span>
            <span class="info-value">{{ deployment.application.name }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">{{ $t('deployment.applicationType') }}</span>
            <span class="info-value">
              <span class="type-badge" :class="`type-badge--${deployment.application.type}`">
                {{ getTypeLabel(deployment.application.type) }}
              </span>
            </span>
          </div>
          <div class="info-item">
            <span class="info-label">{{ $t('deployment.applicationVersion') }}</span>
            <span class="info-value">{{ deployment.application.version }}</span>
          </div>
          <div class="info-item info-item--full">
            <span class="info-label">{{ $t('deployment.applicationDescription') }}</span>
            <span class="info-value description">{{ deployment.application.description || '-' }}</span>
          </div>
        </div>
      </div>

      <div class="detail-section">
        <h3 class="section-title">{{ $t('deployment.serverInfo') }}</h3>
        <div class="info-grid">
          <div class="info-item">
            <span class="info-label">{{ $t('deployment.serverAddress') }}</span>
            <span class="info-value">{{ deployment.server.ip_address }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">{{ $t('deployment.agentPort') }}</span>
            <span class="info-value">{{ deployment.server.agent_port || 8080 }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">{{ $t('deployment.serverHostname') }}</span>
            <span class="info-value">{{ deployment.server.hostname || '-' }}</span>
          </div>
        </div>
      </div>

      <div class="detail-section">
        <h3 class="section-title">{{ $t('deployment.configContent') }}</h3>
        <div class="content-wrapper">
          <pre class="content-code">{{ deployment.content || '-' }}</pre>
        </div>
      </div>
    </div>

    <template #footer>
      <Button type="secondary" @click="handleClose">{{ $t('common.close') }}</Button>
      <Button
        v-if="deployment?.status === 'running'"
        type="primary"
        @click="$emit('stop', deployment)"
      >
        <Icon icon="lucide:square" />
        {{ $t('deployment.stop') }}
      </Button>
      <Button
        v-else-if="deployment"
        type="primary"
        @click="$emit('start', deployment)"
      >
        <Icon icon="lucide:play" />
        {{ $t('deployment.start') }}
      </Button>
      <Button v-if="deployment" type="danger" @click="$emit('undeploy', deployment)">
        <Icon icon="lucide:trash-2" />
        {{ $t('deployment.undeploy') }}
      </Button>
    </template>
  </Modal>
</template>

<script setup lang="ts">
import type { Deployment } from '../types'
import Modal from '@/components/Modal/index.vue'
import Button from '@/components/Button/index.vue'

defineProps<{
  visible: boolean
  deployment: Deployment | null
}>()

const emit = defineEmits<{
  close: []
  start: [deployment: Deployment]
  stop: [deployment: Deployment]
  undeploy: [deployment: Deployment]
}>()

const handleClose = () => {
  emit('close')
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

const getTypeLabel = (type: string) => {
  const labels: Record<string, string> = {
    compose: 'Compose',
    script: '脚本',
    binary: '二进制'
  }
  return labels[type] || type
}
</script>

<style scoped>
.detail-content {
  padding: 20px 0;
}

.detail-section {
  margin-bottom: 24px;
}

.detail-section:last-child {
  margin-bottom: 0;
}

.section-title {
  font-size: 14px;
  font-weight: 600;
  color: #1e3a5f;
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid #e2e8f0;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.info-item--full {
  grid-column: span 2;
}

.info-label {
  font-size: 12px;
  color: #64748b;
}

.info-value {
  font-size: 13px;
  font-weight: 500;
  color: #1e3a5f;
}

.info-value.description {
  font-weight: 400;
  line-height: 1.5;
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

.content-wrapper {
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  padding: 12px;
  max-height: 300px;
  overflow: auto;
}

.content-code {
  margin: 0;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 12px;
  line-height: 1.6;
  color: #334155;
  white-space: pre-wrap;
  word-break: break-all;
}
</style>
