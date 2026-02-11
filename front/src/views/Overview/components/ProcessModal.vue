<template>
  <div v-if="visible" class="modal-overlay" @click="$emit('close')">
    <div class="modal" @click.stop>
      <div class="modal-header">
        <h3>{{ title }}</h3>
        <button class="btn-close" @click="$emit('close')">
          <Icon icon="lucide:x" />
        </button>
      </div>
      <div class="modal-body">
        <table class="process-table">
          <thead>
            <tr>
              <th>PID</th>
              <th>{{ $t('overview.processName') }}</th>
              <th>CPU%</th>
              <th>{{ $t('overview.memoryMB') }}</th>
              <th>{{ $t('overview.status') }}</th>
              <th>{{ $t('overview.operation') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="proc in processes" :key="proc.pid">
              <td>{{ proc.pid }}</td>
              <td>{{ proc.name }}</td>
              <td>{{ proc.cpuPercent.toFixed(2) }}%</td>
              <td>{{ proc.memoryMB.toFixed(2) }}</td>
              <td>{{ proc.status }}</td>
              <td>
                <button class="btn-icon danger" @click="$emit('kill', proc.pid)">
                  <Icon icon="lucide:trash-2" />
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { ProcessInfo } from '@/types'

defineProps<{
  visible: boolean
  title: string
  processes: ProcessInfo[]
}>()

defineEmits<{
  close: []
  kill: [pid: number]
}>()
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.5);
  z-index: 1000;
  padding: 16px;
}

.modal {
  width: 100%;
  max-width: 600px;
  background: #ffffff;
  border-radius: 8px;
  overflow: hidden;
  max-height: 80vh;
  display: flex;
  flex-direction: column;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid #f1f5f9;
}

.modal-header h3 {
  font-size: 14px;
  font-weight: 600;
  color: #1e3a5f;
}

.btn-close {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border-radius: 6px;
  background: #f8fafc;
  color: #64748b;
  border: none;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn-close:hover {
  background: #f1f5f9;
  color: #1e3a5f;
}

.modal-body {
  padding: 16px;
  overflow-y: auto;
}

.process-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 12px;
}

.process-table th,
.process-table td {
  padding: 10px 12px;
  text-align: left;
  border-bottom: 1px solid #f1f5f9;
}

.process-table th {
  font-weight: 600;
  color: #1e3a5f;
  background: #f8fafc;
}

.process-table td {
  color: #64748b;
}

.btn-icon.danger {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  background: #f8fafc;
  border-radius: 6px;
  color: #64748b;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn-icon.danger:hover {
  background: #fee2e2;
  color: #dc2626;
}
</style>
