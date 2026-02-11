<template>
  <div class="deploy-modal-overlay" @click.self="$emit('cancel')">
    <div class="deploy-modal-container">
      <div class="modal-header">
        <h3 class="modal-title">{{ $t('application.deployApplication') }}</h3>
        <button class="close-btn" @click="$emit('cancel')">
          <Icon icon="lucide:x" />
        </button>
      </div>

      <div class="modal-body">
        <div class="form-group">
          <label class="form-label">
            {{ $t('application.applicationName') }}
          </label>
          <div class="application-info">
            <Icon :icon="getTypeIcon(application?.type)" class="type-icon" />
            <div class="info-text">
              <div class="name">{{ application?.name }}</div>
              <div class="version">{{ application?.version }}</div>
            </div>
          </div>
        </div>

        <div class="form-group">
          <label class="form-label">
            {{ $t('application.selectServer') }}
            <span class="required">*</span>
          </label>
          <select
            v-model="selectedServerId"
            class="form-select"
            :disabled="loading || servers.length === 0"
          >
            <option value="">
              {{ $t('application.selectServerPlaceholder') }}
            </option>
            <option v-for="server in servers" :key="server.id" :value="server.id">
              {{ server.server_alias || server.hostname }} ({{ server.ip_address }})
            </option>
          </select>
          <span v-if="errors.serverId" class="error-message">{{ errors.serverId }}</span>
        </div>

        <div v-if="servers.length === 0" class="no-servers">
          <Icon icon="lucide:server-off" class="no-servers-icon" />
          <p>{{ $t('server.noServers') }}</p>
        </div>
      </div>

      <div class="modal-footer">
        <Button type="secondary" @click="$emit('cancel')">
          {{ $t('common.cancel') }}
        </Button>
        <Button
          type="primary"
          :loading="loading"
          :disabled="!selectedServerId || servers.length === 0"
          @click="handleDeploy"
        >
          <Icon icon="lucide:rocket" />
          {{ $t('application.deploy') }}
        </Button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import Button from '@/components/Button/index.vue'
import { fetchServers } from '@/api/server'
import { post } from '@/utils/request'
import type { Server, ApplicationInstance } from '@/types'

const props = defineProps<{
  application: ApplicationInstance | null
}>()

const emit = defineEmits<{
  cancel: []
  success: []
  error: [message: string]
}>()

const { t } = useI18n()

const servers = ref<Server[]>([])
const selectedServerId = ref<number | null>(null)
const loading = ref(false)
const errors = ref({
  serverId: ''
})

const getTypeIcon = (type?: string) => {
  if (!type) return 'lucide:package'
  const icons: Record<string, string> = {
    compose: 'lucide:box',
    script: 'lucide:terminal',
    binary: 'lucide:file-code'
  }
  return icons[type] || 'lucide:package'
}

const loadServers = async () => {
  try {
    servers.value = await fetchServers()
  } catch (error) {
    console.error('Failed to load servers:', error)
  }
}

const validate = (): boolean => {
  errors.value = { serverId: '' }
  let isValid = true

  if (!selectedServerId.value) {
    errors.value.serverId = t('application.required')
    isValid = false
  }

  return isValid
}

const handleDeploy = async () => {
  if (!validate()) return
  if (!selectedServerId.value || !props.application) return

  loading.value = true
  try {
    await post(`/deployment/deploy/${props.application.id}`, {
      server_id: selectedServerId.value
    })
    // post 函数在 code !== 0 时会抛出异常，所以这里直接表示成功
    emit('success')
  } catch (error: any) {
    console.error('Failed to deploy application:', error)
    const errorMessage = error.message || t('application.deployFailed')
    emit('error', errorMessage)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadServers()
})
</script>

<style scoped>
.deploy-modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.deploy-modal-container {
  background: #ffffff;
  border-radius: 16px;
  width: 90%;
  max-width: 500px;
  max-height: 90vh;
  overflow: hidden;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  display: flex;
  flex-direction: column;
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid #e2e8f0;
  flex-shrink: 0;
}

.modal-title {
  font-size: 18px;
  font-weight: 600;
  color: #1e3a5f;
}

.close-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: none;
  background: transparent;
  color: #64748b;
  cursor: pointer;
  border-radius: 6px;
  transition: all 0.2s ease;
  flex-shrink: 0;
}

.close-btn:hover {
  background: #f1f5f9;
  color: #1e3a5f;
}

.modal-body {
  padding: 24px;
  overflow-y: auto;
  flex: 1;
}

.deploy-hint {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: #e0f2fe;
  border-radius: 8px;
  margin-bottom: 20px;
}

.hint-icon {
  width: 20px;
  height: 20px;
  color: #0284c7;
  flex-shrink: 0;
}

.deploy-hint p {
  margin: 0;
  font-size: 14px;
  color: #0284c7;
}

.form-group {
  margin-bottom: 20px;
}

.form-label {
  display: block;
  margin-bottom: 8px;
  font-size: 14px;
  font-weight: 500;
  color: #1e3a5f;
}

.form-label .required {
  color: #ef4444;
  margin-left: 4px;
}

.application-info {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: #f8fafc;
  border-radius: 8px;
}

.type-icon {
  width: 40px;
  height: 40px;
  color: #4fc3f7;
  flex-shrink: 0;
}

.info-text {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.name {
  font-weight: 600;
  color: #1e3a5f;
  font-size: 14px;
}

.version {
  font-size: 12px;
  color: #64748b;
  font-family: 'Monaco', 'Menlo', monospace;
}

.form-select {
  width: 100%;
  padding: 10px 14px;
  border: 2px solid #e2e8f0;
  border-radius: 8px;
  font-size: 14px;
  color: #1e3a5f;
  background: #f8fafc;
  transition: all 0.2s ease;
  cursor: pointer;
}

.form-select:focus {
  outline: none;
  border-color: #4fc3f7;
  background: #ffffff;
}

.form-select:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.error-message {
  margin-top: 6px;
  font-size: 12px;
  color: #ef4444;
}

.no-servers {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  color: #94a3b8;
}

.no-servers-icon {
  width: 48px;
  height: 48px;
  margin-bottom: 12px;
  color: #cbd5e1;
}

.no-servers p {
  font-size: 14px;
  color: #94a3b8;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 24px;
  border-top: 1px solid #e2e8f0;
  background: #f8fafc;
  flex-shrink: 0;
}
</style>
