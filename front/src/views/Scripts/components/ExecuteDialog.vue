<template>
  <div class="execute-dialog-overlay" @click.self="$emit('cancel')">
    <div class="execute-dialog-container">
      <div class="dialog-header">
        <h3 class="dialog-title">{{ $t('scripts.executeScript') }}</h3>
        <button class="close-btn" @click="$emit('cancel')">
          <Icon icon="lucide:x" />
        </button>
      </div>

      <div class="dialog-body">
        <div class="script-preview">
          <div class="script-preview-title">{{ $t('scripts.scriptDetail') }}</div>
          <div class="script-preview-name">{{ script?.name }}</div>
          <pre class="script-preview-content">{{ script?.content }}</pre>
        </div>

        <div class="server-select-wrapper">
          <label class="server-select-label">{{ $t('scripts.selectServer') }}</label>
          <select v-model="selectedServerId" class="server-select">
            <option value="">{{ $t('scripts.selectServer') }}</option>
            <option v-for="server in servers" :key="server.id" :value="server.id">
              {{ server.server_alias || server.hostname }} ({{ server.ip_address }})
            </option>
          </select>
          <span v-if="error" class="error-message">{{ error }}</span>
        </div>
      </div>

      <div class="dialog-footer">
        <Button type="secondary" @click="$emit('cancel')">
          {{ $t('common.cancel') }}
        </Button>
        <Button type="primary" :loading="executing" :disabled="!selectedServerId" @click="handleExecute">
          {{ $t('scripts.execute') }}
        </Button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useI18n } from 'vue-i18n'
import Button from '@/components/Button/index.vue'
import { fetchServers } from '@/api/server'
import type { Script, Server } from '@/types'

defineProps<{
  script: Script | null
}>()

const emit = defineEmits<{
  execute: [serverId: number]
  cancel: []
}>()

const { t } = useI18n()

const servers = ref<Server[]>([])
const selectedServerId = ref<number | ''>('')
const executing = ref(false)
const error = ref('')

onMounted(async () => {
  try {
    servers.value = await fetchServers()
  } catch (err) {
    console.error('Failed to fetch servers:', err)
  }
})

const handleExecute = async () => {
  if (!selectedServerId.value) {
    error.value = t('scripts.selectServer')
    return
  }
  error.value = ''
  executing.value = true
  try {
    emit('execute', Number(selectedServerId.value))
  } finally {
    executing.value = false
  }
}
</script>

<style scoped>
.execute-dialog-overlay {
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

.execute-dialog-container {
  background: #ffffff;
  border-radius: 16px;
  width: 90%;
  max-width: 500px;
  overflow: hidden;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.dialog-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid #e2e8f0;
}

.dialog-title {
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
}

.close-btn:hover {
  background: #f1f5f9;
  color: #1e3a5f;
}

.dialog-body {
  padding: 24px;
}

.script-preview {
  background: #f8fafc;
  border-radius: 8px;
  padding: 16px;
  margin-bottom: 20px;
}

.script-preview-title {
  font-size: 12px;
  font-weight: 600;
  color: #64748b;
  text-transform: uppercase;
  margin-bottom: 8px;
}

.script-preview-name {
  font-size: 14px;
  font-weight: 500;
  color: #1e3a5f;
  margin-bottom: 8px;
}

.script-preview-content {
  font-size: 12px;
  font-family: 'Fira Code', 'Consolas', monospace;
  color: #94a3b8;
  background: #1e3a5f;
  padding: 12px;
  border-radius: 6px;
  max-height: 120px;
  overflow-y: auto;
  white-space: pre-wrap;
  word-break: break-all;
  margin: 0;
}

.server-select-wrapper {
  margin-bottom: 8px;
}

.server-select-label {
  display: block;
  margin-bottom: 8px;
  font-size: 14px;
  font-weight: 500;
  color: #1e3a5f;
}

.server-select {
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

.server-select:focus {
  outline: none;
  border-color: #4fc3f7;
  background: #ffffff;
}

.error-message {
  margin-top: 6px;
  font-size: 12px;
  color: #ef4444;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 24px;
  border-top: 1px solid #e2e8f0;
  background: #f8fafc;
}
</style>
