<template>
  <div class="info-card">
    <div class="info-header">
      <h3>{{ $t('overview.applications') }}</h3>
      <button class="icon-btn"><Icon icon="lucide:settings" /></button>
    </div>
    <div class="app-list">
      <div v-for="app in apps" :key="app.id" class="app-item">
        <div class="app-icon" :style="{ backgroundColor: app.color }">
          <Icon :icon="app.icon" />
        </div>
        <div class="app-info">
          <div class="app-name-row">
            <span class="app-name">{{ app.name }}</span>
            <Icon icon="lucide:chevron-down" class="app-expand" />
          </div>
          <span class="app-version">{{ $t('overview.version') }}: {{ app.version }}</span>
          <div class="app-actions">
            <button class="app-action-btn">{{ $t('overview.stop') }}</button>
            <button class="app-action-btn">{{ $t('overview.restart') }}</button>
            <button class="app-action-btn">{{ $t('overview.directory') }}</button>
          </div>
        </div>
        <div class="app-status">
          <span class="status-badge" :class="app.status">{{ getAppStatusText(app.status) }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue'
import type { Application } from '@/types'

defineProps<{
  apps: Application[]
}>()

const getAppStatusText = (status: string): string => {
  const map: Record<string, string> = {
    running: '运行中',
    stopped: '已停止',
    error: '错误'
  }
  return map[status] || status
}
</script>

<style scoped>
.info-card {
  background: #ffffff;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  overflow: hidden;
}

.info-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  border-bottom: 1px solid #f1f5f9;
}

.info-header h3 {
  font-size: 14px;
  font-weight: 600;
  color: #1e3a5f;
}

.icon-btn {
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

.icon-btn:hover {
  background: #f1f5f9;
  color: #4fc3f7;
}

.app-list {
  padding: 12px 16px;
  max-height: 400px;
  overflow-y: auto;
}

.app-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 12px 0;
  border-bottom: 1px solid #f8fafc;
}

.app-item:last-child {
  border-bottom: none;
}

.app-icon {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  color: #ffffff;
  font-size: 20px;
  flex-shrink: 0;
}

.app-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.app-name-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.app-name {
  font-size: 13px;
  font-weight: 600;
  color: #1e3a5f;
}

.app-expand {
  width: 14px;
  height: 14px;
  color: #94a3b8;
  cursor: pointer;
}

.app-version {
  font-size: 11px;
  color: #94a3b8;
}

.app-actions {
  display: flex;
  gap: 6px;
  margin-top: 4px;
}

.app-action-btn {
  padding: 3px 8px;
  font-size: 11px;
  color: #64748b;
  background: #f8fafc;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.app-action-btn:hover {
  background: #e0f2fe;
  color: #0284c7;
}

.app-status {
  flex-shrink: 0;
}

.status-badge {
  padding: 3px 8px;
  border-radius: 10px;
  font-size: 10px;
  font-weight: 500;
}

.status-badge.running {
  background: #dcfce7;
  color: #16a34a;
}

.status-badge.stopped {
  background: #f1f5f9;
  color: #64748b;
}

.status-badge.error {
  background: #fee2e2;
  color: #dc2626;
}
</style>
