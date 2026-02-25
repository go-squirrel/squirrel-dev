<template>
  <div class="chart-card">
    <div class="chart-header">
      <h3 class="chart-title">
        <Icon icon="lucide:cpu" />
        {{ $t('monitor.cpuUsage') }}
      </h3>
      <span class="chart-value" :class="getUsageClass(latestValue)">
        {{ formatPercent(latestValue) }}
      </span>
    </div>
    <div class="chart-container" ref="chartContainer"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted, nextTick } from 'vue'
import { Icon } from '@iconify/vue'
import * as echarts from 'echarts'
import type { BaseMonitorRecord } from '@/types/monitor'
import { formatPercent } from '@/utils/format'

const props = defineProps<{
  data: BaseMonitorRecord[]
}>()

const chartContainer = ref<HTMLDivElement>()
let chartInstance: echarts.ECharts | null = null

const sortedData = computed(() => {
  return [...props.data].sort((a, b) =>
    new Date(a.collect_time).getTime() - new Date(b.collect_time).getTime()
  )
})

const latestValue = computed(() => {
  if (sortedData.value.length === 0) return 0
  return sortedData.value[sortedData.value.length - 1].cpu_usage
})

const getUsageClass = (value: number) => {
  if (value >= 90) return 'danger'
  if (value >= 70) return 'warning'
  return 'normal'
}

const getChartOption = () => {
  const times = sortedData.value.map(d =>
    new Date(d.collect_time).toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  )
  const values = sortedData.value.map(d => d.cpu_usage)

  return {
    tooltip: {
      trigger: 'axis',
      formatter: (params: any) => {
        return `<strong>${params[0].axisValue}</strong><br/>
                CPU: ${params[0].value.toFixed(2)}%`
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
      type: 'line',
      smooth: true,
      symbol: 'none',
      lineStyle: { width: 2, color: '#4fc3f7' },
      areaStyle: {
        color: {
          type: 'linear',
          x: 0, y: 0, x2: 0, y2: 1,
          colorStops: [
            { offset: 0, color: 'rgba(79, 195, 247, 0.3)' },
            { offset: 1, color: 'rgba(79, 195, 247, 0)' }
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

watch(() => props.data, () => nextTick(() => updateChart()), { deep: true })
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

.chart-value {
  font-size: 18px;
  font-weight: 600;
}

.chart-value.normal { color: #67c23a; }
.chart-value.warning { color: #e6a23c; }
.chart-value.danger { color: #f56c6c; }

.chart-container {
  height: 200px;
  width: 100%;
}
</style>
