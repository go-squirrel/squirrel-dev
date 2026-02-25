<template>
  <div class="chart-card">
    <div class="chart-header">
      <h3 class="chart-title">
        <Icon icon="lucide:hard-drive" />
        {{ $t('monitor.diskUsage') }}
      </h3>
      <select v-model="selectedMount" class="device-select">
        <option v-for="mount in mountPoints" :key="mount" :value="mount">
          {{ mount }}
        </option>
      </select>
    </div>
    <div class="chart-info">
      {{ formatBytes(latestUsed) }} / {{ formatBytes(latestTotal) }}
    </div>
    <div class="chart-container" ref="chartContainer"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted, nextTick } from 'vue'
import { Icon } from '@iconify/vue'
import * as echarts from 'echarts'
import type { DiskUsageRecord } from '@/types/monitor'
import { formatBytes } from '@/utils/format'

const props = defineProps<{
  data: DiskUsageRecord[]
  mountPoints: string[]
}>()

const chartContainer = ref<HTMLDivElement>()
let chartInstance: echarts.ECharts | null = null

const selectedMount = ref('/')

const filteredData = computed(() => {
  return props.data.filter(d => d.mount_point === selectedMount.value)
})

const latestUsed = computed(() => {
  if (filteredData.value.length === 0) return 0
  return filteredData.value[filteredData.value.length - 1].used
})

const latestTotal = computed(() => {
  if (filteredData.value.length === 0) return 0
  return filteredData.value[filteredData.value.length - 1].total
})

const getChartOption = () => {
  const times = filteredData.value.map(d =>
    new Date(d.collect_time).toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  )
  const values = filteredData.value.map(d => d.usage)

  return {
    tooltip: {
      trigger: 'axis',
      formatter: (params: any) => {
        return `<strong>${params[0].axisValue}</strong><br/>
                ${params[0].seriesName}: ${params[0].value.toFixed(2)}%`
      }
    },
    grid: {
      left: '2%',
      right: '2%',
      bottom: '15%',
      top: '10%',
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
      min: 0,
      max: 100,
      axisLine: { show: false },
      axisTick: { show: false },
      splitLine: { lineStyle: { color: '#e2e8f0', type: 'dashed' } },
      axisLabel: {
        color: '#64748b',
        fontSize: 11,
        formatter: '{value}%'
      }
    },
    series: [{
      name: '磁盘使用率',
      type: 'line',
      smooth: true,
      symbol: 'none',
      lineStyle: { width: 2, color: '#e6a23c' },
      areaStyle: {
        color: {
          type: 'linear',
          x: 0, y: 0, x2: 0, y2: 1,
          colorStops: [
            { offset: 0, color: 'rgba(230, 162, 60, 0.3)' },
            { offset: 1, color: 'rgba(230, 162, 60, 0)' }
          ]
        }
      },
      data: values
    }]
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

watch([() => props.data, selectedMount], () => nextTick(() => updateChart()), { deep: true })
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
  margin-bottom: 8px;
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

.chart-info {
  font-size: 12px;
  color: #64748b;
  margin-bottom: 12px;
}

.chart-container {
  height: 200px;
  width: 100%;
}
</style>
