<template>
  <div class="config-table-container">
    <table class="config-table">
      <thead>
        <tr>
          <th class="sortable" @click="$emit('sort', 'key')">
            {{ $t('configs.configKey') }}
            <Icon :icon="getSortIcon('key')" class="sort-icon" />
          </th>
          <th class="sortable" @click="$emit('sort', 'value')">
            {{ $t('configs.configValue') }}
            <Icon :icon="getSortIcon('value')" class="sort-icon" />
          </th>
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
            <button class="action-btn copy-btn" :title="$t('configs.copy')" @click="$emit('copy', config.value)">
              <Icon icon="lucide:copy" />
            </button>
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
import type { Config } from '@/types'

const props = defineProps<{
  configs: Config[]
  sortBy?: string | null
  sortOrder?: 'asc' | 'desc'
}>()

defineEmits<{
  edit: [config: Config]
  delete: [config: Config]
  copy: [value: string]
  sort: [field: string]
}>()

const getSortIcon = (field: string) => {
  if (field !== props.sortBy) return 'lucide:chevrons-up-down'
  return props.sortOrder === 'asc' ? 'lucide:chevron-up' : 'lucide:chevron-down'
}
</script>

<style scoped>
.config-table-container {
  background: #ffffff;
  border-radius: 8px;
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
  font-weight: 600;
  font-size: 13px;
  color: #475569;
  border-bottom: 2px solid #e2e8f0;
}

.config-table th.sortable {
  cursor: pointer;
  user-select: none;
  transition: background 0.2s ease;
}

.config-table th.sortable:hover {
  background: #f1f5f9;
}

.sort-icon {
  width: 14px;
  height: 14px;
  margin-left: 4px;
  color: #94a3b8;
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
  padding: 12px 16px;
  font-size: 13px;
  color: #1e3a5f;
}

.key-cell {
  min-width: 200px;
}

.key-wrapper {
  display: flex;
  align-items: center;
  gap: 8px;
}

.key-icon {
  width: 16px;
  height: 16px;
  color: #64748b;
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
  width: 140px;
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

.action-btn.copy-btn:hover {
  background: #e0f2fe;
  color: #0284c7;
}

.action-btn.edit-btn:hover {
  background: #fef3c7;
  color: #92400e;
}

.action-btn.delete-btn:hover {
  background: #fee2e2;
  color: #dc2626;
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
