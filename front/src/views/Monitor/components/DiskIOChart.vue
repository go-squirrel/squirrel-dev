<template>
  <div class="chart-card">
    <div class="chart-header">
      <h3 class="chart-title">
        <Icon icon="lucide:database" />
        {{ $t('monitor.diskIO') }}
      </h3>
      <select v-model="selectedDevice" class="device-select">
        <option value="all">{{ $t('monitor.all') }}</option>
        <option v-for="device in devices" :key="device" :value="device">
          {{ device }}
        </option>
      </select>
    </div>
    <div class="chart-stats">
      <div class="stat-item">
        <span class="stat-label">{{ $t('monitor.read') }}</span>
        <span class="stat-value">{{ formatBytes(latestReadSpeed) }}/s</span>
      </div>
      <div class="stat-item">
        <span class="stat-label">{{ $t('monitor.write') }}</span>
        <span class="stat-value">{{ formatBytes(latestWriteSpeed) }}/s</span>
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
import type { DiskIORecord, TimeRange } from '@/types/monitor'
import { formatBytes } from '@/utils/format'
import { calculateDiskIOSpeed, groupAndCalculateSpeed, type DiskIOSpeedRecord } from '@/utils/monitor'

const props = defineProps<{
  data: DiskIORecord[]
  devices: string[]
  timeRange: TimeRange
}>()

const { t } = useI18n()
const selectedDevice = ref('all')
const chartContainer = ref<HTMLDivElement>()
let chartInstance: echarts.ECharts | null = null

// 根据时间范围计算开始时间
const getTimeRangeBounds = (range: TimeRange): [Date, Date] => {
  const now = new Date()
  let start: Date
  switch (range) {
    case '1h':
      start = new Date(now.getTime() - 1 * 60 * 60 * 1000)
      break
    case '6h':
      start = new Date(now.getTime() - 6 * 60 * 60 * 1000)
      break
    case '24h':
      start = new Date(now.getTime() - 24 * 60 * 60 * 1000)
      break
    case '7d':
      start = new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000)
      break
    default:
      start = new Date(now.getTime() - 1 * 60 * 60 * 1000)
  }
  return [start, now]
}

// 计算速率数据
const speedData = computed(() => {
  if (props.data.length < 2) return []

  // 按设备分组并计算速率
  const speedByDevice = groupAndCalculateSpeed(
    props.data,
    'disk_name',
    calculateDiskIOSpeed
  )

  if (selectedDevice.value === 'all') {
    // 聚合所有设备的速率
    const timeMap = new Map<string, { read: number; write: number }>()

    speedByDevice.forEach((records) => {
      records.forEach((record: DiskIOSpeedRecord) => {
        const time = record.collect_time
        if (!timeMap.has(time)) {
          timeMap.set(time, { read: 0, write: 0 })
        }
        const current = timeMap.get(time)!
        current.read += record.read_speed
        current.write += record.write_speed
      })
    })

    return Array.from(timeMap.entries())
      .map(([collect_time, speeds]) => ({
        collect_time,
        read_speed: speeds.read,
        write_speed: speeds.write
      }))
      .sort((a, b) => new Date(a.collect_time).getTime() - new Date(b.collect_time).getTime())
  } else {
    // 返回特定设备的速率
    return speedByDevice.get(selectedDevice.value) || []
  }
})

const latestReadSpeed = computed(() => {
  if (speedData.value.length === 0) return 0
  return speedData.value[speedData.value.length - 1].read_speed
})

const latestWriteSpeed = computed(() => {
  if (speedData.value.length === 0) return 0
  return speedData.value[speedData.value.length - 1].write_speed
})

const getChartOption = () => {
  const [startTime, endTime] = getTimeRangeBounds(props.timeRange)

  // 使用时间戳作为 x 轴数据
  const readDataPoints = speedData.value.map(d => [
    new Date(d.collect_time).getTime(),
    d.read_speed
  ])
  const writeDataPoints = speedData.value.map(d => [
    new Date(d.collect_time).getTime(),
    d.write_speed
  ])

  return {
    tooltip: {
      trigger: 'axis',
      formatter: (params: any) => {
        if (!params[0] || params[0].value[1] === null) return ''
        const time = new Date(params[0].value[0]).toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
        let result = `<strong>${time}</strong><br/>`
        params.forEach((param: any) => {
          result += `${param.marker} ${param.seriesName}: ${formatBytes(param.value[1])}/s<br/>`
        })
        return result
      }
    },
    legend: {
      data: [t('monitor.read'), t('monitor.write')],
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
      type: 'time',
      min: startTime.getTime(),
      max: endTime.getTime(),
      axisLine: { lineStyle: { color: '#e2e8f0' } },
      axisLabel: {
        color: '#64748b',
        fontSize: 11,
        rotate: 30,
        formatter: (value: number) => {
          return new Date(value).toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
        }
      },
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
        name: t('monitor.read'),
        type: 'line',
        smooth: true,
        symbol: 'none',
        lineStyle: { width: 2, color: '#4fc3f7' },
        data: readDataPoints
      },
      {
        name: t('monitor.write'),
        type: 'line',
        smooth: true,
        symbol: 'none',
        lineStyle: { width: 2, color: '#94a3b8' },
        data: writeDataPoints
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
  chartInstance.setOption(getChartOption(), { notMerge: true })
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

watch([() => props.data, () => props.timeRange, selectedDevice], () => nextTick(() => updateChart()), { deep: true })
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
