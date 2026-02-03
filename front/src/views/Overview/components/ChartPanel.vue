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
    <div class="chart-container" ref="chartContainer">
      <canvas ref="chartCanvas"></canvas>
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref, watch, nextTick, onMounted, onUnmounted } from 'vue'
import { Icon } from '@iconify/vue'
import { formatBytes, formatSpeed } from '@/utils/format'
import type { ChartDataPoint } from '@/types'

const props = defineProps<{
  chartData: ChartDataPoint[]
  currentNetStats: { bytesSent: number; bytesRecv: number }
  currentIOStats: { readBytes: number; writeBytes: number }
}>()

const chartType = defineModel<'net' | 'io'>('chartType', { default: 'net' })
const chartTarget = defineModel<string>('chartTarget', { default: 'all' })

const chartTargetList = ref<string[]>([])
const chartContainer = ref<HTMLDivElement>()
const chartCanvas = ref<HTMLCanvasElement>()

// 绘制图表
const drawChart = () => {
  if (!chartCanvas.value || !chartContainer.value) return

  const canvas = chartCanvas.value
  const container = chartContainer.value
  const ctx = canvas.getContext('2d')
  if (!ctx) return

  const rect = container.getBoundingClientRect()
  canvas.width = rect.width
  canvas.height = rect.height

  const width = canvas.width
  const height = canvas.height
  const padding = 40

  ctx.clearRect(0, 0, width, height)

  if (props.chartData.length < 2) return

  const allValues = props.chartData.flatMap(d => [d.value1, d.value2])
  const maxValue = Math.max(...allValues, 1)
  const minValue = 0

  // 绘制网格线
  ctx.strokeStyle = '#e2e8f0'
  ctx.lineWidth = 1
  for (let i = 0; i <= 5; i++) {
    const y = padding + (height - 2 * padding) * (1 - i / 5)
    ctx.beginPath()
    ctx.moveTo(padding, y)
    ctx.lineTo(width - padding, y)
    ctx.stroke()
  }

  // 绘制数据线1
  ctx.strokeStyle = '#4fc3f7'
  ctx.lineWidth = 2
  ctx.beginPath()
  props.chartData.forEach((point, index) => {
    const x = padding + (width - 2 * padding) * (index / (props.chartData.length - 1))
    const y = padding + (height - 2 * padding) * (1 - (point.value1 - minValue) / (maxValue - minValue))
    if (index === 0) {
      ctx.moveTo(x, y)
    } else {
      ctx.lineTo(x, y)
    }
  })
  ctx.stroke()

  // 绘制数据线2
  ctx.strokeStyle = '#94a3b8'
  ctx.lineWidth = 2
  ctx.beginPath()
  props.chartData.forEach((point, index) => {
    const x = padding + (width - 2 * padding) * (index / (props.chartData.length - 1))
    const y = padding + (height - 2 * padding) * (1 - (point.value2 - minValue) / (maxValue - minValue))
    if (index === 0) {
      ctx.moveTo(x, y)
    } else {
      ctx.lineTo(x, y)
    }
  })
  ctx.stroke()

  // 绘制 X 轴标签
  ctx.fillStyle = '#64748b'
  ctx.font = '11px sans-serif'
  ctx.textAlign = 'center'
  const step = Math.ceil(props.chartData.length / 6)
  props.chartData.forEach((point, index) => {
    if (index % step === 0) {
      const x = padding + (width - 2 * padding) * (index / (props.chartData.length - 1))
      ctx.fillText(point.time, x, height - 10)
    }
  })
}

// 监听数据变化
watch(() => props.chartData, () => {
  nextTick(() => drawChart())
}, { deep: true })

// 窗口大小变化时重绘
const handleResize = () => {
  nextTick(() => drawChart())
}

onMounted(() => {
  window.addEventListener('resize', handleResize)
  nextTick(() => drawChart())
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
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
}

.chart-container canvas {
  width: 100%;
  height: 100%;
}
</style>
