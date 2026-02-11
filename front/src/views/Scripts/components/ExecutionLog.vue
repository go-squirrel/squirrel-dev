<template>
  <div class="execution-log-overlay" @click.self="$emit('close')">
    <div class="execution-log-container">
      <div class="log-header">
        <div class="log-header-left">
          <h3 class="log-title">{{ $t('scripts.executionLog') }}</h3>
          <div class="log-meta">
            <span class="meta-item">
              <Icon icon="lucide:file-code" class="meta-icon" />
              {{ script?.name }}
            </span>
            <span class="meta-item">
              <Icon icon="lucide:server" class="meta-icon" />
              {{ result?.server_ip }}:{{ result?.agent_port }}
            </span>
            <span class="meta-item">
              <Icon icon="lucide:clock" class="meta-icon" />
              {{ result?.created_at }}
            </span>
            <span class="status-badge" :class="`status-${result?.status}`">
              {{ getStatusText(result?.status) }}
            </span>
          </div>
        </div>
        <div class="log-header-right">
          <button class="action-btn" :title="$t('common.copy')" @click="handleCopy">
            <Icon icon="lucide:copy" />
          </button>
          <button class="action-btn" :title="$t('common.refresh')" @click="$emit('refresh')">
            <Icon icon="lucide:refresh-cw" />
          </button>
          <button class="close-btn" @click="$emit('close')">
            <Icon icon="lucide:x" />
          </button>
        </div>
      </div>

      <div class="log-body">
        <div v-if="result?.error_message" class="error-banner">
          <Icon icon="lucide:alert-circle" class="error-icon" />
          <span class="error-text">{{ result.error_message }}</span>
        </div>

        <div class="log-content-wrapper">
          <div v-if="!result?.output" class="empty-log">
            <Icon icon="lucide:file-text" class="empty-icon" />
            <p>{{ $t('scripts.noOutput') }}</p>
          </div>
          <pre v-else class="log-content">{{ result.output }}</pre>
        </div>
      </div>

      <div class="log-footer">
        <div class="log-stats">
          <span class="stat-item">{{ $t('scripts.taskId') }}: {{ result?.task_id }}</span>
        </div>
        <Button type="secondary" @click="$emit('close')">
          {{ $t('common.close') }}
        </Button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import Button from '@/components/Button/index.vue'
import type { Script, ScriptResult } from '@/types'

const props = defineProps<{
  script: Script | null
  result: ScriptResult | null
}>()

defineEmits<{
  close: []
  refresh: []
}>()

const { t } = useI18n()

const getStatusText = (status: string | undefined): string => {
  if (!status) return ''
  const statusMap: Record<string, string> = {
    running: t('scripts.statusRunning'),
    success: t('scripts.statusSuccess'),
    failed: t('scripts.statusFailed')
  }
  return statusMap[status] || status
}

const handleCopy = async () => {
  if (!props.result?.output) return
  try {
    await navigator.clipboard.writeText(props.result.output)
    // 可以添加复制成功的提示
  } catch (err) {
    console.error('Failed to copy:', err)
  }
}
</script>

<style scoped>
.execution-log-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1100;
  padding: 20px;
}

.execution-log-container {
  background: #ffffff;
  border-radius: 16px;
  width: 100%;
  max-width: 1200px;
  height: 90vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  box-shadow: 0 25px 80px rgba(0, 0, 0, 0.4);
}

.log-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid #e2e8f0;
  background: #f8fafc;
}

.log-header-left {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.log-title {
  font-size: 18px;
  font-weight: 600;
  color: #1e3a5f;
  margin: 0;
}

.log-meta {
  display: flex;
  align-items: center;
  gap: 16px;
  flex-wrap: wrap;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: #64748b;
}

.meta-icon {
  width: 14px;
  height: 14px;
}

.status-badge {
  display: inline-flex;
  align-items: center;
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 500;
}

.status-running {
  background: rgba(245, 158, 11, 0.1);
  color: #f59e0b;
}

.status-success {
  background: rgba(34, 197, 94, 0.1);
  color: #22c55e;
}

.status-failed {
  background: rgba(239, 68, 68, 0.1);
  color: #ef4444;
}

.log-header-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border: none;
  background: transparent;
  color: #64748b;
  cursor: pointer;
  border-radius: 8px;
  transition: all 0.2s ease;
}

.action-btn:hover {
  background: #e2e8f0;
  color: #1e3a5f;
}

.close-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border: none;
  background: transparent;
  color: #64748b;
  cursor: pointer;
  border-radius: 8px;
  transition: all 0.2s ease;
  margin-left: 8px;
}

.close-btn:hover {
  background: #fee2e2;
  color: #ef4444;
}

.log-body {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.error-banner {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 24px;
  background: rgba(239, 68, 68, 0.1);
  border-bottom: 1px solid rgba(239, 68, 68, 0.2);
}

.error-icon {
  width: 20px;
  height: 20px;
  color: #ef4444;
  flex-shrink: 0;
}

.error-text {
  font-size: 14px;
  color: #ef4444;
  font-weight: 500;
}

.log-content-wrapper {
  flex: 1;
  overflow: auto;
  background: #0f172a;
  padding: 20px 24px;
}

.log-content {
  margin: 0;
  font-family: 'Fira Code', 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.7;
  color: #e2e8f0;
  white-space: pre-wrap;
  word-break: break-all;
}

.empty-log {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #64748b;
}

.empty-icon {
  width: 48px;
  height: 48px;
  margin-bottom: 16px;
  color: #475569;
}

.log-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 24px;
  border-top: 1px solid #e2e8f0;
  background: #f8fafc;
}

.log-stats {
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-item {
  font-size: 13px;
  color: #64748b;
}
</style>
