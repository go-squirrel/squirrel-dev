<template>
  <div class="tooltip-header">{{ $t('overview.loadDetail') }}</div>
  <div class="tooltip-content">
    <div class="load-item">
      <span>{{ $t('overview.load1') }}</span>
      <span class="load-value">{{ data?.load1 || 0 }}</span>
    </div>
    <div class="load-item">
      <span>{{ $t('overview.load5') }}</span>
      <span class="load-value">{{ data?.load5 || 0 }}</span>
    </div>
    <div class="load-item">
      <span>{{ $t('overview.load15') }}</span>
      <span class="load-value">{{ data?.load15 || 0 }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  data?: {
    load1: number
    load5: number
    load15: number
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

.load-item {
  display: flex;
  justify-content: space-between;
  padding: 6px 0;
  font-size: 12px;
}

.load-item span:first-child {
  color: #64748b;
}

.load-value {
  font-weight: 500;
  color: #1e3a5f;
}
</style>
