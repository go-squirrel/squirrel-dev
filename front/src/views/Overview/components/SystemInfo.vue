<template>
  <div class="info-card">
    <div class="info-header">
      <h3>{{ $t('overview.systemInfo') }}</h3>
      <div class="info-actions">
        <button class="icon-btn"><Icon icon="lucide:refresh-cw" /></button>
        <button class="icon-btn"><Icon icon="lucide:maximize-2" /></button>
      </div>
    </div>
    <div class="info-content">
      <div class="info-item">
        <span class="info-label">{{ $t('overview.hostname') }}</span>
        <span class="info-value">{{ info.hostname }}</span>
      </div>
      <div class="info-item">
        <span class="info-label">{{ $t('overview.os') }}</span>
        <span class="info-value">{{ info.os }}</span>
      </div>
      <div class="info-item">
        <span class="info-label">{{ $t('overview.kernel') }}</span>
        <span class="info-value">{{ info.kernel }}</span>
      </div>
      <div class="info-item">
        <span class="info-label">{{ $t('overview.arch') }}</span>
        <span class="info-value">{{ info.arch }}</span>
      </div>
      <div class="info-item">
        <span class="info-label">{{ $t('overview.ip') }}</span>
        <span class="info-value">{{ info.ip }}</span>
      </div>
      <div class="info-item">
        <span class="info-label">{{ $t('overview.bootTime') }}</span>
        <span class="info-value">{{ info.bootTime }}</span>
      </div>
      <div class="info-item">
        <span class="info-label">{{ $t('overview.uptime') }}</span>
        <span class="info-value">{{ info.uptime }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue'
import type { SystemInfo } from '@/types'

defineProps<{
  info: SystemInfo
}>()
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

.info-actions {
  display: flex;
  gap: 6px;
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

.info-content {
  padding: 12px 16px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  padding: 8px 0;
  font-size: 12px;
  border-bottom: 1px solid #f8fafc;
}

.info-item:last-child {
  border-bottom: none;
}

.info-label {
  color: #64748b;
}

.info-value {
  color: #1e3a5f;
  font-weight: 500;
  max-width: 150px;
  text-align: right;
  word-break: break-all;
}
</style>
