<template>
  <div class="script-table-container">
    <table class="script-table">
      <thead>
        <tr>
          <th>{{ $t('scripts.scriptName') }}</th>
          <th>{{ $t('scripts.operation') }}</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="script in scripts" :key="script.id" class="script-row">
          <td class="name-cell">
            <div class="name-wrapper">
              <Icon icon="lucide:file-code" class="name-icon" />
              <span class="name-text" :title="script.name">{{ script.name }}</span>
            </div>
          </td>
          <td class="action-cell">
            <button class="action-btn edit-btn" :title="$t('scripts.editScript')" @click="$emit('edit', script)">
              <Icon icon="lucide:edit-2" />
            </button>
            <button class="action-btn execute-btn" :title="$t('scripts.executeScript')" @click="$emit('execute', script)">
              <Icon icon="lucide:play" />
            </button>
            <button class="action-btn result-btn" :title="$t('scripts.viewResults')" @click="$emit('viewResults', script)">
              <Icon icon="lucide:list" />
            </button>
            <button class="action-btn delete-btn" :title="$t('scripts.deleteScript')" @click="$emit('delete', script)">
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
import type { Script } from '@/types'

defineProps<{
  scripts: Script[]
}>()

defineEmits<{
  edit: [script: Script]
  delete: [script: Script]
  execute: [script: Script]
  viewResults: [script: Script]
}>()
</script>

<style scoped>
.script-table-container {
  background: #ffffff;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.script-table {
  width: 100%;
  border-collapse: collapse;
}

.script-table th {
  padding: 14px 20px;
  text-align: left;
  font-weight: 600;
  font-size: 13px;
  color: #64748b;
  background: #f8fafc;
  border-bottom: 1px solid #e2e8f0;
}

.script-table td {
  padding: 16px 20px;
  border-bottom: 1px solid #f1f5f9;
}

.script-row:hover {
  background: #f8fafc;
}

.name-cell {
  width: 60%;
}

.name-wrapper {
  display: flex;
  align-items: center;
  gap: 10px;
}

.name-icon {
  width: 20px;
  height: 20px;
  color: #4fc3f7;
  flex-shrink: 0;
}

.name-text {
  font-size: 14px;
  font-weight: 500;
  color: #1e3a5f;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.action-cell {
  width: 40%;
}

.action-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  margin-right: 8px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
  background: transparent;
}

.action-btn:hover {
  transform: translateY(-1px);
}

.edit-btn {
  color: #4fc3f7;
}

.edit-btn:hover {
  background: rgba(79, 195, 247, 0.1);
}

.execute-btn {
  color: #22c55e;
}

.execute-btn:hover {
  background: rgba(34, 197, 94, 0.1);
}

.result-btn {
  color: #8b5cf6;
}

.result-btn:hover {
  background: rgba(139, 92, 246, 0.1);
}

.delete-btn {
  color: #ef4444;
}

.delete-btn:hover {
  background: rgba(239, 68, 68, 0.1);
}
</style>
