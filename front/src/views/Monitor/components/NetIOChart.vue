<template>
  <div class="chart-card">
    <div class="chart-header">
      <h3 class="chart-title">
        <Icon icon="lucide:activity" />
        {{ $t('monitor.networkTraffic') }}
      </h3>
      <select v-model="selectedInterface" class="device-select">
        <option value="all">{{ $t('monitor.all') }}</option>
        <option v-for="iface in interfaces" :key="iface" :value="iface">
          {{ iface }}
        </option>
      </select>
    </div>
    <div class="chart-stats">
      <div class="stat-item">
        <span class="stat-label">{{ $t('monitor.upload') }}</span>
        <span class="stat-value">{{ formatBytes(latestUploadSpeed) }}/s</span>
      </div>
      <div class="stat-item">
        <span class="stat-label">{{ $t('monitor.download') }}</span>
        <span class="stat-value">{{ formatBytes(latestDownloadSpeed) }}/s</span>
      </div>
    </div>
    <div class="chart-container" ref="chartContainer"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import { Icon } from '@iconify/vue'
import * as echarts from 'echarts'
import type { NetworkIORecord } from '@/types/monitor'
import { formatBytes } from '@/utils/format'
import { calculateNetworkIOSpeed, groupAndCalculateSpeed, type NetworkIOSpeedRecord } from '@/utils/monitor'

const props = defineProps<{
  data: NetworkIORecord[]
  interfaces: string[]
}>()

const { t } = useI18n()
const selectedInterface = ref('all')
const chartContainer = ref<HTMLDivElement>()
let chartInstance: echarts.ECharts | null = null

// 计算速率数据
const speedData = computed(() => {
  if (props.data.length < 2) return []

  // 按设备分组并计算速率
  const speedByInterface = groupAndCalculateSpeed(
    props.data,
    'interface_name',
    calculateNetworkIOSpeed
  )

  if (selectedInterface.value === 'all') {
    // 聚合所有网卡的速率
    const timeMap = new Map<string, { upload: number; download: number }>()

    speedByInterface.forEach((records) => {
      records.forEach((record: NetworkIOSpeedRecord) => {
        const time = record.collect_time
        if (!timeMap.has(time)) {
          timeMap.set(time, { upload: 0, download: 0 })
        }
        const current = timeMap.get(time)!
        current.upload += record.upload_speed
        current.download += record.download_speed
      })
    })

    return Array.from(timeMap.entries())
      .map(([collect_time, speeds]) => ({
        collect_time,
        upload_speed: speeds.upload,
        download_speed: speeds.download
      }))
      .sort((a, b) => new Date(a.collect_time).getTime() - new Date(b.collect_time).getTime())
  } else {
    // 返回特定网卡的速率
    return speedByInterface.get(selectedInterface.value) || []
  }
})

const latestUploadSpeed = computed(() => {
  if (speedData.value.length === 0) return 0
  return speedData.value[speedData.value.length - 1].upload_speed
})

const latestDownloadSpeed = computed(() => {
  if (speedData.value.length === 0) return 0
  return speedData.value[speedData.value.length - 1].download_speed
})

const getChartOption = () => {
  const times = speedData.value.map(d =>
    new Date(d.collect_time).toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  )
  const uploadValues = speedData.value.map(d => d.upload_speed)
  const downloadValues = speedData.value.map(d => d.download_speed)

  return {
    tooltip: {
      trigger: 'axis',
      formatter: (params: any) => {
        let result = `<strong>${params[0].axisValue}</strong><br/>`
        params.forEach((param: any) => {
          result += `${param.marker} ${param.seriesName}: ${formatBytes(param.value)}/s<br/>`
        })
        return result
      }
    },
    legend: {
      data: [t('monitor.upload'), t('monitor.download')],
      right: '2%',
      top: 0,
      textStyle: { color: '#64748b', fontSize: 12 }
    },
    grid: {
      left: '2%',
      right: '2%',
      bottom: '15%',
      top: '15%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: times,
      axisLine: { lineStyle: { color: '#e2e8f0' } },
      axisLabel: { color: '#64748b', fontSize: 11, rotate: 30 },
      axisTick: { show: false }
    },
    yAxis: {
      type: 'value',
      axisLine: { show: false },
      axisTick: { show: false },
      splitLine: { lineStyle: { color: '#e2e8f0', type: 'dashed' } },
      axisLabel: {
        color: '#64748b',
        fontSize: 11,
        formatter: (value: number) => formatBytes(value) + '/s'
      }
    },
    series: [
      {
        name: t('monitor.upload'),
        type: 'line',
        smooth: true,
        symbol: 'none',
        lineStyle: { width: 2, color: '#4fc3f7' },
        data: uploadValues
      },
      {
        name: t('monitor.download'),
        type: 'line',
        smooth: true,
        symbol: 'none',
        lineStyle: { width: 2, color: '#94a3b8' },
        data: downloadValues
      }
    ]
  }
}

const initChart = () => {
  if (!chartContainer.value) return
  chartInstance = echarts.init(chartContainer.value)
  updateChart()
}

const updateChart = () => {
  if (!chartInstance) return
  chartInstance.setOption(getChartOption())
}

const handleResize = () => chartInstance?.resize()

onMounted(() => {
  window.addEventListener('resize', handleResize)
  nextTick(() => initChart())
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  chartInstance?.dispose()
})

watch([() => props.data, selectedInterface], () => nextTick(() => updateChart()), { deep: true })
</script>

<style scoped lang="scss">
.chart-card {
  background: #ffffff;
  border-radius: 8px;
  padding: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.chart-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  font-weight: 600;
  color: #1e3a5f;
  margin: 0;
}

.device-select {
  padding: 4px 8px;
  font-size: 12px;
  border: 1px solid #e2e8f0;
  border-radius: 4px;
  background: #f8fafc;
  color: #1e3a5f;
  cursor: pointer;
}

.chart-stats {
  display: flex;
  gap: 16px;
  margin-bottom: 12px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 8px;
  background: #f8fafc;
  border-radius: 4px;
  font-size: 12px;
}

.stat-label {
  color: #64748b;
}

.stat-value {
  font-weight: 500;
  color: #1e3a5f;
}

.chart-container {
  height: 200px;
  width: 100%;
}
</style>
