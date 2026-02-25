<template>
  <div class="time-range-selector">
    <button
      v-for="range in ranges"
      :key="range.value"
      class="range-btn"
      :class="{ active: modelValue === range.value }"
      @click="$emit('update:modelValue', range.value)"
    >
      {{ range.label }}
    </button>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import type { TimeRange } from '@/types/monitor'

defineProps<{
  modelValue: TimeRange
}>()

defineEmits<{
  'update:modelValue': [value: TimeRange]
}>()

const { t } = useI18n()

const ranges = [
  { value: '1h' as TimeRange, label: t('monitor.last1Hour') },
  { value: '6h' as TimeRange, label: t('monitor.last6Hours') },
  { value: '24h' as TimeRange, label: t('monitor.last24Hours') },
  { value: '7d' as TimeRange, label: t('monitor.last7Days') },
  { value: '30d' as TimeRange, label: t('monitor.last30Days') }
]
</script>

<style scoped lang="scss">
.time-range-selector {
  display: flex;
  gap: 8px;
  padding: 16px;
  background: #ffffff;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.range-btn {
  padding: 8px 16px;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  background: #fff;
  font-size: 13px;
  color: #64748b;
  cursor: pointer;
  transition: all 0.2s;
}

.range-btn:hover {
  border-color: #4fc3f7;
  color: #4fc3f7;
}

.range-btn.active {
  background: #4fc3f7;
  border-color: #4fc3f7;
  color: #fff;
}
</style>
