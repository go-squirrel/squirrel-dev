<template>
  <div class="application-table-container">
    <table class="application-table">
      <thead>
        <tr>
          <th class="sortable" @click="$emit('sort', 'name')">
            {{ $t('application.applicationName') }}
            <Icon :icon="getSortIcon('name')" class="sort-icon" />
          </th>
          <th class="sortable" @click="$emit('sort', 'type')">
            {{ $t('application.applicationType') }}
            <Icon :icon="getSortIcon('type')" class="sort-icon" />
          </th>
          <th class="sortable" @click="$emit('sort', 'version')">
            {{ $t('application.applicationVersion') }}
            <Icon :icon="getSortIcon('version')" class="sort-icon" />
          </th>
          <th>{{ $t('application.applicationDescription') }}</th>
          <th>{{ $t('common.operation') }}</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="app in applications" :key="app.id" class="application-row">
          <td class="name-cell">
            <div class="name-wrapper">
              <Icon :icon="getTypeIcon(app.type)" class="type-icon" />
              <span class="name-text">{{ app.name }}</span>
            </div>
          </td>
          <td class="type-cell">
            <span class="type-badge" :class="`type-badge--${app.type}`">
              {{ getTypeLabel(app.type) }}
            </span>
          </td>
          <td class="version-cell">
            <span class="version-text">{{ app.version }}</span>
          </td>
          <td class="description-cell">
            <div class="description-wrapper">
              <span class="description-text" :title="app.description">{{ app.description }}</span>
            </div>
          </td>
          <td class="action-cell">
            <button class="action-btn view-btn" :title="$t('application.viewDetail')" @click="$emit('view', app)">
              <Icon icon="lucide:eye" />
            </button>
            <button class="action-btn edit-btn" :title="$t('application.editApplication')" @click="$emit('edit', app)">
              <Icon icon="lucide:edit-2" />
            </button>
            <button class="action-btn delete-btn" :title="$t('application.deleteApplication')" @click="$emit('delete', app)">
              <Icon icon="lucide:trash-2" />
            </button>
            <button class="action-btn deploy-btn" :title="$t('application.deploy')" @click="$emit('deploy', app)">
              <Icon icon="lucide:rocket" />
            </button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue'
import type { ApplicationInstance } from '@/types'

const props = defineProps<{
  applications: ApplicationInstance[]
  sortBy?: string | null
  sortOrder?: 'asc' | 'desc'
}>()

defineEmits<{
  edit: [application: ApplicationInstance]
  delete: [application: ApplicationInstance]
  view: [application: ApplicationInstance]
  deploy: [application: ApplicationInstance]
  sort: [field: string]
}>()

const getSortIcon = (field: string) => {
  if (field !== props.sortBy) return 'lucide:chevrons-up-down'
  return props.sortOrder === 'asc' ? 'lucide:chevron-up' : 'lucide:chevron-down'
}

const getTypeIcon = (type: string) => {
  const icons: Record<string, string> = {
    compose: 'lucide:box',
    script: 'lucide:terminal',
    binary: 'lucide:file-code'
  }
  return icons[type] || 'lucide:package'
}

const getTypeLabel = (type: string) => {
  const labels: Record<string, string> = {
    compose: 'Compose',
    script: '脚本',
    binary: '二进制'
  }
  return labels[type] || type
}
</script>

<style scoped>
.application-table-container {
  background: #ffffff;
  border-radius: 8px;
  overflow: hidden;
}

.application-table {
  width: 100%;
  border-collapse: collapse;
}

.application-table thead {
  background: #f8fafc;
}

.application-table th {
  padding: 12px 16px;
  text-align: left;
  font-weight: 600;
  font-size: 13px;
  color: #475569;
  border-bottom: 2px solid #e2e8f0;
}

.application-table th.sortable {
  cursor: pointer;
  user-select: none;
  transition: background 0.2s ease;
}

.application-table th.sortable:hover {
  background: #f1f5f9;
}

.sort-icon {
  width: 14px;
  height: 14px;
  margin-left: 4px;
  color: #94a3b8;
}

.application-row {
  border-bottom: 1px solid #f1f5f9;
  transition: background 0.2s ease;
}

.application-row:hover {
  background: #f8fafc;
}

.application-row:last-child {
  border-bottom: none;
}

.application-table td {
  padding: 12px 16px;
  font-size: 13px;
  color: #1e3a5f;
}

.name-cell {
  width: 200px;
}

.name-wrapper {
  display: flex;
  align-items: center;
  gap: 8px;
}

.type-icon {
  width: 16px;
  height: 16px;
  color: #64748b;
}

.name-text {
  font-weight: 500;
}

.type-cell {
  width: 100px;
}

.type-badge {
  display: inline-block;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.type-badge--compose {
  background: #dbeafe;
  color: #1e40af;
}

.type-badge--script {
  background: #dcfce7;
  color: #166534;
}

.type-badge--binary {
  background: #fef3c7;
  color: #92400e;
}

.version-cell {
  width: 100px;
}

.version-text {
  font-family: 'Monaco', 'Menlo', monospace;
  font-size: 12px;
  color: #64748b;
}

.description-cell {
  max-width: 300px;
}

.description-wrapper {
  display: flex;
  align-items: center;
}

.description-text {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: #64748b;
}

.action-cell {
  width: 160px;
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

.action-btn.view-btn:hover {
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

.action-btn.deploy-btn:hover {
  background: #dbeafe;
  color: #0284c7;
}
</style>
