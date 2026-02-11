<template>
  <div class="app-table-container">
    <table class="app-table">
      <thead>
        <tr>
          <th>{{ $t('appStore.appName') }}</th>
          <th>{{ $t('appStore.type') }}</th>
          <th>{{ $t('appStore.category') }}</th>
          <th>{{ $t('appStore.version') }}</th>
          <th>{{ $t('appStore.downloads') }}</th>
          <th>{{ $t('appStore.status') }}</th>
          <th>{{ $t('appStore.operation') }}</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="app in apps" :key="app.id" class="app-row">
          <td class="name-cell">
            <div class="name-wrapper">
              <img v-if="app.icon" :src="app.icon" class="app-icon" />
              <Icon v-else icon="lucide:package" class="app-icon-default" />
              <div class="name-info">
                <div class="name-text">
                  {{ app.name }}
                  <span v-if="app.is_official" class="official-badge">
                    {{ $t('appStore.isOfficial') }}
                  </span>
                </div>
                <div class="description-text">{{ app.description }}</div>
              </div>
            </div>
          </td>
          <td class="type-cell">
            <span class="type-tag" :class="`type-${app.type}`">
              {{ getTypeLabel(app.type) }}
            </span>
          </td>
          <td class="category-cell">
            <span class="category-tag">
              {{ getCategoryLabel(app.category) }}
            </span>
          </td>
          <td class="version-cell">
            <span class="version-text">{{ app.version }}</span>
          </td>
          <td class="downloads-cell">
            <Icon icon="lucide:download" class="downloads-icon" />
            <span>{{ app.downloads }}</span>
          </td>
          <td class="status-cell">
            <span class="status-badge" :class="`status-${app.status}`">
              {{ getStatusLabel(app.status) }}
            </span>
          </td>
          <td class="action-cell">
            <button class="action-btn import-btn" :title="$t('appStore.importToApplication')" @click="$emit('import', app)">
              <Icon icon="lucide:download-cloud" />
            </button>
            <button class="action-btn detail-btn" :title="$t('appStore.viewDetail')" @click="$emit('detail', app)">
              <Icon icon="lucide:eye" />
            </button>
            <button class="action-btn edit-btn" :title="$t('appStore.editApp')" @click="$emit('edit', app)">
              <Icon icon="lucide:edit-2" />
            </button>
            <button class="action-btn delete-btn" :title="$t('appStore.deleteApp')" @click="$emit('delete', app)">
              <Icon icon="lucide:trash-2" />
            </button>
          </td>
        </tr>
      </tbody>
    </table>
    <div v-if="apps.length === 0" class="empty-state">
      <Icon icon="lucide:inbox" class="empty-icon" />
      <p class="empty-text">{{ $t('common.noData') }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import type { AppStore } from '@/types'

defineProps<{
  apps: AppStore[]
}>()

defineEmits<{
  detail: [app: AppStore]
  edit: [app: AppStore]
  delete: [app: AppStore]
  import: [app: AppStore]
}>()

const { t } = useI18n()

const getTypeLabel = (type: string) => {
  const labels: Record<string, string> = {
    compose: t('appStore.typeCompose'),
    k8s_manifest: t('appStore.typeK8s'),
    helm_chart: t('appStore.typeHelm')
  }
  return labels[type] || type
}

const getCategoryLabel = (category: string) => {
  const labels: Record<string, string> = {
    web: t('appStore.categoryWeb'),
    database: t('appStore.categoryDatabase'),
    middleware: t('appStore.categoryMiddleware'),
    devops: t('appStore.categoryDevops')
  }
  return labels[category] || category
}

const getStatusLabel = (status: string) => {
  const labels: Record<string, string> = {
    active: t('appStore.statusActive'),
    deprecated: t('appStore.statusDeprecated')
  }
  return labels[status] || status
}
</script>

<style scoped>
.app-table-container {
  background: #ffffff;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  overflow: hidden;
}

.app-table {
  width: 100%;
  border-collapse: collapse;
}

.app-table thead {
  background: #f8fafc;
}

.app-table th {
  padding: 12px 16px;
  text-align: left;
  font-size: 12px;
  font-weight: 600;
  color: #64748b;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.app-row {
  border-bottom: 1px solid #f1f5f9;
  transition: background 0.2s ease;
}

.app-row:hover {
  background: #f8fafc;
}

.app-row:last-child {
  border-bottom: none;
}

.app-table td {
  padding: 14px 16px;
  font-size: 13px;
  color: #1e3a5f;
}

.name-cell {
  min-width: 280px;
}

.name-wrapper {
  display: flex;
  align-items: center;
  gap: 12px;
}

.app-icon {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  object-fit: cover;
  flex-shrink: 0;
}

.app-icon-default {
  width: 40px;
  height: 40px;
  color: #94a3b8;
  flex-shrink: 0;
}

.name-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.name-text {
  font-weight: 600;
  color: #1e3a5f;
  display: flex;
  align-items: center;
  gap: 8px;
}

.official-badge {
  display: inline-flex;
  align-items: center;
  padding: 2px 8px;
  background: linear-gradient(135deg, #fbbf24 0%, #f59e0b 100%);
  color: #ffffff;
  font-size: 10px;
  font-weight: 600;
  border-radius: 4px;
  text-transform: uppercase;
}

.description-text {
  font-size: 12px;
  color: #64748b;
  max-width: 300px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.type-cell,
.category-cell {
  min-width: 120px;
}

.type-tag {
  display: inline-flex;
  align-items: center;
  padding: 4px 10px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 500;
}

.type-compose {
  background: #dbeafe;
  color: #1d4ed8;
}

.type-k8s_manifest {
  background: #dcfce7;
  color: #15803d;
}

.type-helm_chart {
  background: #f3e8ff;
  color: #7c3aed;
}

.category-tag {
  display: inline-flex;
  align-items: center;
  padding: 4px 10px;
  background: #f1f5f9;
  color: #475569;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 500;
}

.version-cell {
  min-width: 80px;
}

.version-text {
  font-family: 'SF Mono', Monaco, Consolas, monospace;
  font-size: 12px;
  color: #64748b;
}

.downloads-cell {
  min-width: 100px;
  display: flex;
  align-items: center;
  gap: 6px;
  color: #64748b;
  font-size: 12px;
}

.downloads-icon {
  width: 14px;
  height: 14px;
  color: #94a3b8;
}

.status-cell {
  min-width: 80px;
}

.status-badge {
  display: inline-flex;
  align-items: center;
  padding: 4px 10px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 500;
}

.status-active {
  background: #dcfce7;
  color: #15803d;
}

.status-deprecated {
  background: #f1f5f9;
  color: #64748b;
}

.action-cell {
  min-width: 180px;
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

.detail-btn {
  background: #eff6ff;
  color: #0284c7;
}

.detail-btn:hover {
  background: #dbeafe;
  transform: scale(1.1);
}

.import-btn {
  background: #e0f2fe;
  color: #0284c7;
}

.import-btn:hover {
  background: #bae6fd;
  transform: scale(1.1);
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
