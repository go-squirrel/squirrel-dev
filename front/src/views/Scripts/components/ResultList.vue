<template>
  <div class="result-list-overlay" @click.self="$emit('close')">
    <div class="result-list-container">
      <div class="dialog-header">
        <h3 class="dialog-title">{{ $t('scripts.executionResults') }} - {{ script?.name }}</h3>
        <button class="close-btn" @click="$emit('close')">
          <Icon icon="lucide:x" />
        </button>
      </div>

      <div class="result-list-body">
        <div v-if="results.length === 0" class="empty-results">
          <Icon icon="lucide:clipboard-list" class="empty-results-icon" />
          <p>{{ $t('common.noData') }}</p>
        </div>

        <div v-else class="result-table-container">
          <table class="result-table">
            <thead>
              <tr>
                <th>{{ $t('scripts.taskId') }}</th>
                <th>{{ $t('scripts.serverIP') }}</th>
                <th>{{ $t('scripts.status') }}</th>
                <th>{{ $t('scripts.createdAt') }}</th>
                <th>{{ $t('scripts.operation') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="result in results" :key="result.id" class="result-row">
                <td>{{ result.task_id }}</td>
                <td>{{ result.server_ip }}:{{ result.agent_port }}</td>
                <td>
                  <span class="status-badge" :class="`status-${result.status}`">
                    {{ getStatusText(result.status) }}
                  </span>
                </td>
                <td>{{ result.created_at }}</td>
                <td>
                  <button class="view-log-btn" @click="$emit('viewLog', result)">
                    <Icon icon="lucide:file-text" />
                    {{ $t('scripts.viewLog') }}
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div class="dialog-footer">
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

defineProps<{
  script: Script | null
  results: ScriptResult[]
}>()

defineEmits<{
  close: []
  viewLog: [result: ScriptResult]
}>()

const { t } = useI18n()

const getStatusText = (status: string): string => {
  const statusMap: Record<string, string> = {
    running: t('scripts.statusRunning'),
    success: t('scripts.statusSuccess'),
    failed: t('scripts.statusFailed')
  }
  return statusMap[status] || status
}
</script>

<style scoped>
.result-list-overlay {
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

.result-list-container {
  background: #ffffff;
  border-radius: 16px;
  width: 90%;
  max-width: 800px;
  max-height: 90vh;
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

.result-list-body {
  padding: 24px;
  overflow-y: auto;
  max-height: calc(90vh - 140px);
}

.result-table-container {
  background: #ffffff;
  border-radius: 12px;
  border: 1px solid #e2e8f0;
  overflow: hidden;
}

.result-table {
  width: 100%;
  border-collapse: collapse;
}

.result-table th {
  padding: 12px 16px;
  text-align: left;
  font-weight: 600;
  font-size: 12px;
  color: #64748b;
  background: #f8fafc;
  border-bottom: 1px solid #e2e8f0;
}

.result-table td {
  padding: 14px 16px;
  border-bottom: 1px solid #f1f5f9;
  font-size: 13px;
  color: #1e3a5f;
}

.result-row:hover {
  background: #f8fafc;
}

.status-badge {
  display: inline-flex;
  align-items: center;
  padding: 4px 10px;
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

.view-log-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  border: none;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 500;
  color: #4fc3f7;
  background: rgba(79, 195, 247, 0.1);
  cursor: pointer;
  transition: all 0.2s ease;
}

.view-log-btn:hover {
  background: rgba(79, 195, 247, 0.2);
  transform: translateY(-1px);
}

.empty-results {
  text-align: center;
  padding: 40px;
  color: #64748b;
}

.empty-results-icon {
  width: 48px;
  height: 48px;
  margin-bottom: 16px;
  color: #cbd5e1;
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
