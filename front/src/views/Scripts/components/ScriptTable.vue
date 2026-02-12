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
  border-radius: 8px;
  overflow: hidden;
}

.script-table {
  width: 100%;
  border-collapse: collapse;
}

.script-table thead {
  background: #f8fafc;
}

.script-table th {
  padding: 12px 16px;
  text-align: left;
  font-weight: 600;
  font-size: 13px;
  color: #475569;
  border-bottom: 2px solid #e2e8f0;
}

.script-row {
  border-bottom: 1px solid #f1f5f9;
  transition: background 0.2s ease;
}

.script-row:hover {
  background: #f8fafc;
}

.script-row:last-child {
  border-bottom: none;
}

.script-table td {
  padding: 12px 16px;
  font-size: 13px;
  color: #1e3a5f;
}

.name-cell {
  width: 60%;
}

.name-wrapper {
  display: flex;
  align-items: center;
  gap: 8px;
}

.name-icon {
  width: 16px;
  height: 16px;
  color: #64748b;
}

.name-text {
  font-size: 13px;
  font-weight: 500;
  color: #1e3a5f;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.action-cell {
  width: 40%;
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

.action-btn.edit-btn:hover {
  background: #fef3c7;
  color: #92400e;
}

.action-btn.execute-btn:hover {
  background: #dcfce7;
  color: #16a34a;
}

.action-btn.result-btn:hover {
  background: #e0f2fe;
  color: #0284c7;
}

.action-btn.delete-btn:hover {
  background: #fee2e2;
  color: #dc2626;
}
</style>
