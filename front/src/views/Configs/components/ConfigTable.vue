<template>
  <div class="config-table-container">
    <table class="config-table">
      <thead>
        <tr>
          <th>{{ $t('configs.configKey') }}</th>
          <th>{{ $t('configs.configValue') }}</th>
          <th>{{ $t('configs.operation') }}</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="config in configs" :key="config.id" class="config-row">
          <td class="key-cell">
            <div class="key-wrapper">
              <Icon icon="lucide:key" class="key-icon" />
              <span class="key-text">{{ config.key }}</span>
            </div>
          </td>
          <td class="value-cell">
            <div class="value-wrapper">
              <span class="value-text" :title="config.value">{{ config.value }}</span>
            </div>
          </td>
          <td class="action-cell">
            <button class="action-btn edit-btn" :title="$t('configs.editConfig')" @click="$emit('edit', config)">
              <Icon icon="lucide:edit-2" />
            </button>
            <button class="action-btn delete-btn" :title="$t('configs.deleteConfig')" @click="$emit('delete', config)">
              <Icon icon="lucide:trash-2" />
            </button>
          </td>
        </tr>
      </tbody>
    </table>
    <div v-if="configs.length === 0" class="empty-state">
      <Icon icon="lucide:inbox" class="empty-icon" />
      <p class="empty-text">{{ $t('common.noData') }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue'
import type { Config } from '@/types'

defineProps<{
  configs: Config[]
}>()

defineEmits<{
  edit: [config: Config]
  delete: [config: Config]
}>()
</script>

<style scoped>
.config-table-container {
  background: #ffffff;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  overflow: hidden;
}

.config-table {
  width: 100%;
  border-collapse: collapse;
}

.config-table thead {
  background: #f8fafc;
}

.config-table th {
  padding: 12px 16px;
  text-align: left;
  font-size: 12px;
  font-weight: 600;
  color: #64748b;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.config-row {
  border-bottom: 1px solid #f1f5f9;
  transition: background 0.2s ease;
}

.config-row:hover {
  background: #f8fafc;
}

.config-row:last-child {
  border-bottom: none;
}

.config-table td {
  padding: 14px 16px;
  font-size: 13px;
  color: #1e3a5f;
}

.key-cell {
  min-width: 200px;
}

.key-wrapper {
  display: flex;
  align-items: center;
  gap: 10px;
}

.key-icon {
  width: 16px;
  height: 16px;
  color: #4fc3f7;
  flex-shrink: 0;
}

.key-text {
  font-weight: 500;
  color: #1e3a5f;
  font-family: 'SF Mono', Monaco, Consolas, monospace;
}

.value-cell {
  max-width: 400px;
}

.value-wrapper {
  display: flex;
  align-items: center;
}

.value-text {
  font-family: 'SF Mono', Monaco, Consolas, monospace;
  color: #64748b;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 400px;
  display: block;
}

.action-cell {
  min-width: 100px;
}

.action-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 6px;
  border: none;
  cursor: pointer;
  transition: all 0.2s ease;
  margin-right: 4px;
}

.action-btn:last-child {
  margin-right: 0;
}

.edit-btn {
  background: #fef9c3;
  color: #ca8a04;
}

.edit-btn:hover {
  background: #fef08a;
  transform: scale(1.1);
}

.delete-btn {
  background: #fef2f2;
  color: #dc2626;
}

.delete-btn:hover {
  background: #fee2e2;
  transform: scale(1.1);
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  color: #94a3b8;
}

.empty-icon {
  width: 48px;
  height: 48px;
  margin-bottom: 12px;
  color: #cbd5e1;
}

.empty-text {
  font-size: 13px;
  color: #94a3b8;
}
</style>
