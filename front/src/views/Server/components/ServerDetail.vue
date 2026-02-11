<template>
  <div class="modal-overlay" @click.self="$emit('close')">
    <div class="modal">
      <div class="modal-header">
        <h3>{{ $t('server.serverDetail') }}</h3>
        <button class="close-btn" @click="$emit('close')">
          <Icon icon="lucide:x" />
        </button>
      </div>

      <div class="modal-body">
        <div v-if="loading" class="loading">
          <Icon icon="lucide:loader-2" class="spinner" />
          <p>{{ $t('server.loading') }}</p>
        </div>

        <template v-else-if="serverDetail">
          <div class="section">
            <h4>{{ $t('server.basicInfo') }}</h4>
            <div class="info-grid">
              <div class="info-item">
                <span class="label">{{ $t('server.hostname') }}</span>
                <span class="value">{{ serverDetail.hostname }}</span>
              </div>
              <div class="info-item">
                <span class="label">{{ $t('server.ipAddress') }}</span>
                <span class="value">{{ serverDetail.ip_address }}</span>
              </div>
              <div class="info-item">
                <span class="label">{{ $t('server.sshPort') }}</span>
                <span class="value">{{ serverDetail.ssh_port }}</span>
              </div>
              <div class="info-item">
                <span class="label">{{ $t('server.username') }}</span>
                <span class="value">{{ serverDetail.ssh_username }}</span>
              </div>
              <div class="info-item">
                <span class="label">{{ $t('server.authType') }}</span>
                <span class="value">
                  {{ serverDetail.auth_type === 'password' ? $t('server.password') : $t('server.key') }}
                </span>
              </div>
              <div class="info-item">
                <span class="label">{{ $t('server.status') }}</span>
                <span class="value">
                  <span class="status-badge" :class="serverDetail.status">
                    <span class="status-dot"></span>
                    {{ getStatusText(serverDetail.status) }}
                  </span>
                </span>
              </div>
              <div v-if="serverDetail.server_alias" class="info-item">
                <span class="label">{{ $t('server.serverAlias') }}</span>
                <span class="value">{{ serverDetail.server_alias }}</span>
              </div>
            </div>
          </div>

          <div v-if="serverDetail.server_info" class="section">
            <h4>{{ $t('server.systemInfo') }}</h4>
            <div class="info-grid">
              <div class="info-item">
                <span class="label">{{ $t('server.os') }}</span>
                <span class="value">{{ serverDetail.server_info.os }}</span>
              </div>
              <div class="info-item">
                <span class="label">{{ $t('server.kernel') }}</span>
                <span class="value">{{ serverDetail.server_info.kernelVersion }}</span>
              </div>
              <div class="info-item">
                <span class="label">{{ $t('server.architecture') }}</span>
                <span class="value">{{ serverDetail.server_info.architecture }}</span>
              </div>
              <div class="info-item">
                <span class="label">{{ $t('server.platform') }}</span>
                <span class="value">{{ serverDetail.server_info.platform }}</span>
              </div>
              <div class="info-item">
                <span class="label">{{ $t('server.platformVersion') }}</span>
                <span class="value">{{ serverDetail.server_info.platformVersion }}</span>
              </div>
              <div class="info-item">
                <span class="label">{{ $t('server.uptime') }}</span>
                <span class="value">{{ serverDetail.server_info.uptimeStr }}</span>
              </div>
            </div>

            <div v-if="serverDetail.server_info.ipAddresses" class="ip-list">
              <h5>{{ $t('server.ipAddresses') }}</h5>
              <div v-for="addr in serverDetail.server_info.ipAddresses" :key="addr.name" class="ip-group">
                <span class="ip-name">{{ addr.name }}</span>
                <div class="ip-values">
                  <span v-if="addr.ipv4.length > 0" class="ip-tag ipv4">
                    IPv4: {{ addr.ipv4.join(', ') }}
                  </span>
                  <span v-if="addr.ipv6.length > 0" class="ip-tag ipv6">
                    IPv6: {{ addr.ipv6.join(', ') }}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { fetchServerDetail } from '@/api/server'
import type { Server } from '@/types'

const props = defineProps<{
  server: Server
}>()

defineEmits<{
  close: []
}>()

const { t } = useI18n()

const loading = ref(true)
const serverDetail = ref<Server | null>(null)

const loadServerDetail = async () => {
  loading.value = true
  try {
    serverDetail.value = await fetchServerDetail(props.server.id)
  } catch (error) {
    console.error('Failed to load server detail:', error)
  } finally {
    loading.value = false
  }
}

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

onMounted(() => {
  loadServerDetail()
})
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  padding: 20px;
}

.modal {
  background: #ffffff;
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.12);
  max-width: 600px;
  width: 100%;
  max-height: 90vh;
  overflow-y: auto;
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid #f1f5f9;
}

.modal-header h3 {
  font-size: 16px;
  font-weight: 600;
  color: #1e3a5f;
}

.close-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border-radius: 6px;
  background: #f5f7fa;
  color: #64748b;
  border: none;
  cursor: pointer;
  transition: all 0.2s ease;
}

.close-btn:hover {
  background: #fee2e2;
  color: #dc2626;
}

.modal-body {
  padding: 24px;
}

.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px;
  color: #64748b;
}

.spinner {
  width: 32px;
  height: 32px;
  animation: spin 1s linear infinite;
  margin-bottom: 12px;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

.section {
  margin-bottom: 24px;
}

.section:last-child {
  margin-bottom: 0;
}

.section h4 {
  font-size: 14px;
  font-weight: 600;
  color: #1e3a5f;
  margin-bottom: 16px;
  padding-bottom: 8px;
  border-bottom: 2px solid #f1f5f9;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.info-item .label {
  font-size: 11px;
  font-weight: 500;
  color: #64748b;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.info-item .value {
  font-size: 13px;
  font-weight: 500;
  color: #1e3a5f;
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

.ip-list {
  margin-top: 20px;
}

.ip-list h5 {
  font-size: 12px;
  font-weight: 600;
  color: #64748b;
  margin-bottom: 12px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.ip-group {
  margin-bottom: 12px;
}

.ip-group:last-child {
  margin-bottom: 0;
}

.ip-name {
  display: block;
  font-size: 12px;
  font-weight: 500;
  color: #1e3a5f;
  margin-bottom: 6px;
}

.ip-values {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.ip-tag {
  display: inline-block;
  padding: 6px 12px;
  border-radius: 6px;
  font-size: 12px;
  font-family: 'SF Mono', Monaco, Consolas, monospace;
}

.ip-tag.ipv4 {
  background: #e0f2fe;
  color: #0284c7;
}

.ip-tag.ipv6 {
  background: #fef3c7;
  color: #d97706;
}
</style>
