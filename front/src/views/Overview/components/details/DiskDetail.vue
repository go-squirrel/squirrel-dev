<template>
  <div class="tooltip-header">{{ $t('overview.diskDetail') }}</div>
  <div class="tooltip-content">
    <div class="disk-list">
      <div
        v-for="partition in disk?.partitions || []"
        :key="partition.device"
        class="disk-item"
      >
        <div class="disk-header">
          <span>{{ partition.device }}</span>
          <span>{{ partition.mountPoint }}</span>
        </div>
        <div class="disk-bar">
          <div class="disk-fill" :style="{ width: partition.usage + '%' }"></div>
        </div>
        <div class="disk-info">
          <span>{{ formatBytes(partition.used) }} / {{ formatBytes(partition.total) }}</span>
          <span>{{ partition.usage.toFixed(1) }}%</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { formatBytes } from '@/utils/format'
import type { DiskPartition } from '@/types'

defineProps<{
  disk?: {
    available: number
    partitions: DiskPartition[]
    total: number
    usage: number
    used: number
  }
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

.disk-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.disk-item {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.disk-header {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  font-weight: 500;
  color: #1e3a5f;
}

.disk-bar {
  height: 6px;
  background: #f1f5f9;
  border-radius: 3px;
  overflow: hidden;
}

.disk-fill {
  height: 100%;
  background: linear-gradient(90deg, #4fc3f7 0%, #29b6f6 100%);
  border-radius: 3px;
  transition: width 0.3s ease;
}

.disk-info {
  display: flex;
  justify-content: space-between;
  font-size: 11px;
  color: #64748b;
}
</style>
