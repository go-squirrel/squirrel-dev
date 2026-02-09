<template>
  <section class="monitor-section">
    <h3 class="section-title">
      <Icon icon="lucide:activity" />
      <span>{{ $t('overview.status') }}</span>
    </h3>
    <div class="metrics-grid">
      <!-- 负载 - 显示原始值，无单位 -->
      <MetricCircle
        :value="loadMetric.usage"
        :label="$t('overview.load')"
        :sub-label="loadMetric.status || ''"
        :show-tooltip="activeTooltip === 'load'"
        @show-detail="$emit('showTooltip', 'load')"
        @hide-detail="$emit('hideTooltip')"
      >
        <LoadDetail :data="monitorData.loadAverage" />
      </MetricCircle>

      <!-- CPU - 显示百分比 -->
      <MetricCircle
        :value="cpuMetric.usage"
        :label="$t('overview.cpu')"
        :sub-label="`( ${cpuMetric.used} / ${cpuMetric.total} ) ${$t('overview.cores')}`"
        unit="%"
        :show-tooltip="activeTooltip === 'cpu'"
        @show-detail="$emit('showTooltip', 'cpu')"
        @hide-detail="$emit('hideTooltip')"
      >
        <CpuDetail 
          :cpu="monitorData.cpu" 
          @view-process="$emit('showProcess', 'cpu')" 
        />
      </MetricCircle>

      <!-- 内存 - 显示百分比 -->
      <MetricCircle
        :value="memoryMetric.usage"
        :label="$t('overview.memory')"
        :sub-label="`${memoryMetric.used} / ${memoryMetric.total}`"
        unit="%"
        :show-tooltip="activeTooltip === 'memory'"
        @show-detail="$emit('showTooltip', 'memory')"
        @hide-detail="$emit('hideTooltip')"
      >
        <MemoryDetail 
          :memory="monitorData.memory" 
          @view-process="$emit('showProcess', 'memory')" 
        />
      </MetricCircle>

      <!-- 磁盘 - 显示百分比 -->
      <MetricCircle
        :value="diskMetric.usage"
        :label="$t('overview.disk')"
        :sub-label="`${diskMetric.used} / ${diskMetric.total}`"
        unit="%"
        :show-tooltip="activeTooltip === 'disk'"
        @show-detail="$emit('showTooltip', 'disk')"
        @hide-detail="$emit('hideTooltip')"
      >
        <DiskDetail :disk="monitorData.disk" />
      </MetricCircle>
    </div>
  </section>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue'
import MetricCircle from './MetricCircle.vue'
import LoadDetail from './details/LoadDetail.vue'
import CpuDetail from './details/CpuDetail.vue'
import MemoryDetail from './details/MemoryDetail.vue'
import DiskDetail from './details/DiskDetail.vue'
import type { MonitorData } from '@/types'

interface MetricData {
  usage: number
  status?: string
  used?: string | number
  total?: string | number
}

defineProps<{
  monitorData: MonitorData
  loadMetric: MetricData
  cpuMetric: MetricData
  memoryMetric: MetricData
  diskMetric: MetricData
  activeTooltip: string | null
}>()

defineEmits<{
  showTooltip: [type: string]
  hideTooltip: []
  showProcess: [type: 'cpu' | 'memory']
}>()
</script>

<style scoped>
.monitor-section {
  background: #ffffff;
  border-radius: 8px;
  padding: 10px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.section-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  font-weight: 600;
  color: #1e3a5f;
  margin-bottom: 16px;
}

.metrics-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 24px;
}
</style>
