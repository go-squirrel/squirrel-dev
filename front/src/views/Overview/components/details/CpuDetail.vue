<template>
  <div class="tooltip-header">{{ $t('overview.cpuDetail') }}</div>
  <div class="tooltip-content">
    <div class="cpu-info">
      <span>{{ $t('overview.cpuModel') }}: {{ cpu?.model || '-' }}</span>
      <span>{{ $t('overview.cpuFreq') }}: {{ ((cpu?.frequency || 0) / 1000).toFixed(2) }} GHz</span>
      <span>{{ $t('overview.cpuCores') }}: {{ cpu?.cores || 0 }} {{ $t('overview.cores') }}</span>
    </div>
    <div class="cpu-cores">
      <div
        v-for="(usage, index) in cpu?.perCoreUsage || []"
        :key="index"
        class="core-item"
      >
        <span>{{ $t('overview.core') }} {{ index + 1 }}</span>
        <div class="core-bar">
          <div class="core-fill" :style="{ width: usage + '%' }"></div>
        </div>
        <span>{{ usage.toFixed(1) }}%</span>
      </div>
    </div>
    <div class="tooltip-actions">
      <button @click="$emit('viewProcess')">{{ $t('overview.viewTopProcess') }}</button>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  cpu?: {
    cores: number
    frequency: number
    model: string
    perCoreUsage: number[]
    usage: number
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

.cpu-info {
  display: flex;
  flex-direction: column;
  gap: 6px;
  margin-bottom: 12px;
  font-size: 12px;
  color: #64748b;
}

.cpu-cores {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.core-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 11px;
}

.core-item span:first-child {
  width: 50px;
  color: #64748b;
}

.core-bar {
  flex: 1;
  height: 4px;
  background: #f1f5f9;
  border-radius: 2px;
  overflow: hidden;
}

.core-fill {
  height: 100%;
  background: linear-gradient(90deg, #4fc3f7 0%, #29b6f6 100%);
  border-radius: 2px;
  transition: width 0.3s ease;
}

.core-item span:last-child {
  width: 40px;
  text-align: right;
  color: #1e3a5f;
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
