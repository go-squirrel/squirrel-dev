<template>
  <div class="tooltip-header">{{ $t('overview.memoryDetail') }}</div>
  <div class="tooltip-content">
    <div class="memory-info">
      <div class="mem-item">
        <span>{{ $t('overview.total') }}</span>
        <span>{{ formatBytes(memory?.total || 0) }}</span>
      </div>
      <div class="mem-item">
        <span>{{ $t('overview.used') }}</span>
        <span>{{ formatBytes(memory?.used || 0) }}</span>
      </div>
      <div class="mem-item">
        <span>{{ $t('overview.available') }}</span>
        <span>{{ formatBytes(memory?.available || 0) }}</span>
      </div>
      <div class="mem-item">
        <span>{{ $t('overview.swapTotal') }}</span>
        <span>{{ formatBytes(memory?.swapTotal || 0) }}</span>
      </div>
      <div class="mem-item">
        <span>{{ $t('overview.swapUsed') }}</span>
        <span>{{ formatBytes(memory?.swapUsed || 0) }}</span>
      </div>
    </div>
    <div class="tooltip-actions">
      <button @click="$emit('viewProcess')">{{ $t('overview.viewTopProcess') }}</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { formatBytes } from '@/utils/format'

defineProps<{
  memory?: {
    available: number
    swapTotal: number
    swapUsed: number
    total: number
    usage: number
    used: number
  }
}>()

defineEmits<{
  viewProcess: []
}>()
</script>

<style scoped>
.tooltip-header {
  padding: 12px 16px;
  border-bottom: 1px solid #f1f5f9;
  font-size: 13px;
  font-weight: 600;
  color: #1e3a5f;
}

.tooltip-content {
  padding: 12px 16px;
}

.memory-info {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.mem-item {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
}

.mem-item span:first-child {
  color: #64748b;
}

.mem-item span:last-child {
  color: #1e3a5f;
  font-weight: 500;
}

.tooltip-actions {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid #f1f5f9;
}

.tooltip-actions button {
  width: 100%;
  padding: 8px;
  font-size: 12px;
  font-weight: 500;
  color: #4fc3f7;
  background: #e0f2fe;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.tooltip-actions button:hover {
  background: #4fc3f7;
  color: #ffffff;
}
</style>
