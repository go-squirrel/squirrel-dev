<template>
  <section class="chart-section">
    <div class="chart-header">
      <h3 class="section-title">
        <Icon icon="lucide:bar-chart-2" />
        <span>{{ $t('overview.monitor') }}</span>
      </h3>
      <div class="chart-controls">
        <select v-model="chartType" class="control-select">
          <option value="net">{{ $t('overview.netTraffic') }}</option>
          <option value="io">{{ $t('overview.diskIO') }}</option>
        </select>
        <select v-model="chartTarget" class="control-select">
          <option value="all">{{ $t('overview.all') }}</option>
          <option v-for="item in chartTargetList" :key="item" :value="item">{{ item }}</option>
        </select>
      </div>
    </div>
    <div class="chart-stats">
      <div class="chart-stat-item" v-if="chartType === 'net'">
        <span class="stat-name">{{ $t('overview.upload') }}</span>
        <span class="stat-value">{{ formatSpeed(currentNetStats.bytesSent) }}</span>
      </div>
      <div class="chart-stat-item" v-if="chartType === 'net'">
        <span class="stat-name">{{ $t('overview.download') }}</span>
        <span class="stat-value">{{ formatSpeed(currentNetStats.bytesRecv) }}</span>
      </div>
      <div class="chart-stat-item" v-if="chartType === 'io'">
        <span class="stat-name">{{ $t('overview.read') }}</span>
        <span class="stat-value">{{ formatBytes(currentIOStats.readBytes) }}</span>
      </div>
      <div class="chart-stat-item" v-if="chartType === 'io'">
        <span class="stat-name">{{ $t('overview.write') }}</span>
        <span class="stat-value">{{ formatBytes(currentIOStats.writeBytes) }}</span>
      </div>
    </div>
    <div class="chart-container" ref="chartContainer"></div>
  </section>
</template>

<script setup lang="ts">
import { ref, watch, nextTick, onMounted, onUnmounted } from 'vue'
import { formatBytes, formatSpeed } from '@/utils/format'
import type { ChartDataPoint } from '@/types'
import * as echarts from 'echarts'

const props = defineProps<{
  chartData: ChartDataPoint[]
  currentNetStats: { bytesSent: number; bytesRecv: number }
  currentIOStats: { readBytes: number; writeBytes: number }
  chartTargetList: string[]
}>()

const chartType = defineModel<'net' | 'io'>('chartType', { default: 'net' })
const chartTarget = defineModel<string>('chartTarget', { default: 'all' })

const chartContainer = ref<HTMLDivElement>()
let chartInstance: echarts.ECharts | null = null

// 配置 ECharts 选项
const getChartOption = () => {
  const times = props.chartData.map(d => d.time)
  const values1 = props.chartData.map(d => d.value1)
  const values2 = props.chartData.map(d => d.value2)

  const isNet = chartType.value === 'net'
  const color1 = '#4fc3f7'
  const color2 = '#94a3b8'
  const name1 = isNet ? '上传' : '读取'
  const name2 = isNet ? '下载' : '写入'

  return {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross'
      },
      formatter: (params: any) => {
        if (!params || params.length === 0) return ''
        let result = `<strong>${params[0].axisValue}</strong><br/>`
        params.forEach((param: any) => {
          const value = isNet ? formatSpeed(param.value) : formatBytes(param.value)
          result += `${param.marker} ${param.seriesName}: ${value}<br/>`
        })
        return result
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
      axisLine: {
        lineStyle: {
          color: '#e2e8f0'
        }
      },
      axisLabel: {
        color: '#64748b',
        fontSize: 11,
        rotate: 30
      },
      axisTick: {
        show: false
      }
    },
    yAxis: {
      type: 'value',
      splitLine: {
        lineStyle: {
          color: '#e2e8f0',
          type: 'dashed'
        }
      },
      axisLine: {
        show: false
      },
      axisLabel: {
        color: '#64748b',
        fontSize: 11,
        formatter: (value: number) => isNet ? formatSpeed(value) : formatBytes(value)
      },
      axisTick: {
        show: false
      }
    },
    series: [
      {
        name: name1,
        type: 'line',
        smooth: true,
        symbol: 'none',
        lineStyle: {
          width: 2,
          color: color1
        },
        areaStyle: {
          color: {
            type: 'linear',
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: [
              {
                offset: 0,
                color: 'rgba(79, 195, 247, 0.3)'
              },
              {
                offset: 1,
                color: 'rgba(79, 195, 247, 0)'
              }
            ]
          }
        },
        data: values1
      },
      {
        name: name2,
        type: 'line',
        smooth: true,
        symbol: 'none',
        lineStyle: {
          width: 2,
          color: color2
        },
        areaStyle: {
          color: {
            type: 'linear',
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: [
              {
                offset: 0,
                color: 'rgba(148, 163, 184, 0.3)'
              },
              {
                offset: 1,
                color: 'rgba(148, 163, 184, 0)'
              }
            ]
          }
        },
        data: values2
      }
    ]
  }
}

// 初始化图表
const initChart = () => {
  if (!chartContainer.value) return

  chartInstance = echarts.init(chartContainer.value)
  updateChart()
}

// 更新图表
const updateChart = () => {
  if (!chartInstance) return

  const option = getChartOption()
  chartInstance.setOption(option)
}

// 监听数据变化
watch(() => props.chartData, () => {
  nextTick(() => updateChart())
}, { deep: true })

// 监听图表类型变化
watch(() => chartType.value, () => {
  nextTick(() => updateChart())
})

// 窗口大小变化时重绘
const handleResize = () => {
  if (chartInstance) {
    chartInstance.resize()
  }
}

onMounted(() => {
  window.addEventListener('resize', handleResize)
  nextTick(() => initChart())
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  if (chartInstance) {
    chartInstance.dispose()
  }
})
</script>

<style scoped>
.chart-section {
  background: #ffffff;
  border-radius: 8px;
  padding: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  flex: 1;
  min-height: 300px;
}

.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  font-weight: 600;
  color: #1e3a5f;
  margin-bottom: 0;
}

.chart-controls {
  display: flex;
  gap: 8px;
}

.control-select {
  padding: 6px 12px;
  font-size: 12px;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  background: #f8fafc;
  color: #1e3a5f;
  cursor: pointer;
}

.chart-stats {
  display: flex;
  gap: 16px;
  margin-bottom: 16px;
}

.chart-stat-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: #f8fafc;
  border-radius: 6px;
  font-size: 12px;
}

.chart-stat-item .stat-name {
  color: #64748b;
}

.chart-stat-item .stat-value {
  font-weight: 500;
  color: #1e3a5f;
}

.chart-container {
  height: 200px;
  position: relative;
  width: 100%;
}
</style>
