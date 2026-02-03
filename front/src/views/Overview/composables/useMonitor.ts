// 监控数据相关逻辑
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { fetchMonitorStats, fetchNetStats, fetchIOStats } from '@/api'
import { formatBytes } from '@/utils/format'
import type { MonitorData, ChartDataPoint, ChartType } from '../types'

export function useMonitor(serverId: Ref<number>) {
  const monitorData = ref<MonitorData>({})
  const chartData = ref<ChartDataPoint[]>([])
  const chartType = ref<ChartType>('net')
  const chartTarget = ref('all')
  const currentNetStats = ref({ bytesSent: 0, bytesRecv: 0 })
  const currentIOStats = ref({ readBytes: 0, writeBytes: 0 })
  
  let lastNetStats = { bytesSent: 0, bytesRecv: 0, timestamp: 0 }
  let lastIOStats = { readBytes: 0, writeBytes: 0, timestamp: 0 }
  let monitorTimer: ReturnType<typeof setInterval> | null = null
  let chartTimer: ReturnType<typeof setInterval> | null = null

  // 计算指标
  const loadMetric = computed(() => {
    const usage = monitorData.value.loadAverage?.load1 || 0
    let status = '运行流畅'
    if (usage > 80) status = '负载过高'
    else if (usage > 50) status = '负载中等'
    return { usage, status }
  })

  const cpuMetric = computed(() => {
    const rawUsage = monitorData.value.cpu?.usage || 0
    const usage = rawUsage > 1 ? rawUsage : rawUsage
    const cores = monitorData.value.cpu?.cores || 0
    const used = ((usage / 100) * cores).toFixed(2)
    return { usage, used, total: cores }
  })

  const memoryMetric = computed(() => {
    const rawUsage = monitorData.value.memory?.usage || 0
    const usage = rawUsage > 1 ? rawUsage : rawUsage
    const used = formatBytes(monitorData.value.memory?.used || 0)
    const total = formatBytes(monitorData.value.memory?.total || 0)
    return { usage, used, total }
  })

  const diskMetric = computed(() => {
    const rawUsage = monitorData.value.disk?.usage || 0
    const usage = rawUsage > 1 ? rawUsage : rawUsage
    const used = formatBytes(monitorData.value.disk?.used || 0)
    const total = formatBytes(monitorData.value.disk?.total || 0)
    return { usage, used, total }
  })

  // 获取监控数据
  const loadMonitorStats = async () => {
    if (!serverId.value) return
    try {
      const data = await fetchMonitorStats(serverId.value)
      monitorData.value = data
    } catch (error) {
      console.error('Failed to fetch monitor stats:', error)
    }
  }

  // 获取图表数据
  const loadChartData = async () => {
    if (!serverId.value) return
    try {
      if (chartType.value === 'net') {
        const data = await fetchNetStats(serverId.value, chartTarget.value)
        const now = Date.now()
        const { bytesSent, bytesRecv } = data

        if (lastNetStats.timestamp > 0) {
          const timeDiff = (now - lastNetStats.timestamp) / 1000
          const sentSpeed = timeDiff > 0 ? Math.max(0, (bytesSent - lastNetStats.bytesSent) / timeDiff) : 0
          const recvSpeed = timeDiff > 0 ? Math.max(0, (bytesRecv - lastNetStats.bytesRecv) / timeDiff) : 0

          currentNetStats.value = { bytesSent: sentSpeed, bytesRecv: recvSpeed }
          updateChartData(sentSpeed, recvSpeed)
        }
        lastNetStats = { bytesSent, bytesRecv, timestamp: now }
      } else {
        const data = await fetchIOStats(serverId.value, chartTarget.value)
        const now = Date.now()
        const { readBytes, writeBytes } = data

        if (lastIOStats.timestamp > 0) {
          const timeDiff = (now - lastIOStats.timestamp) / 1000
          const readSpeed = timeDiff > 0 ? Math.max(0, (readBytes - lastIOStats.readBytes) / timeDiff) : 0
          const writeSpeed = timeDiff > 0 ? Math.max(0, (writeBytes - lastIOStats.writeBytes) / timeDiff) : 0

          currentIOStats.value = { readBytes: readSpeed, writeBytes: writeSpeed }
          updateChartData(readSpeed, writeSpeed)
        }
        lastIOStats = { readBytes, writeBytes, timestamp: now }
      }
    } catch (error) {
      console.error('Failed to fetch chart data:', error)
    }
  }

  // 更新图表数据
  const updateChartData = (value1: number, value2: number) => {
    const now = new Date()
    const time = `${now.getHours().toString().padStart(2, '0')}:${now.getMinutes().toString().padStart(2, '0')}:${now.getSeconds().toString().padStart(2, '0')}`

    chartData.value.push({ time, value1, value2 })
    if (chartData.value.length > 30) {
      chartData.value.shift()
    }
  }

  // 重置数据
  const resetData = () => {
    chartData.value = []
    lastNetStats = { bytesSent: 0, bytesRecv: 0, timestamp: 0 }
    lastIOStats = { readBytes: 0, writeBytes: 0, timestamp: 0 }
    currentNetStats.value = { bytesSent: 0, bytesRecv: 0 }
    currentIOStats.value = { readBytes: 0, writeBytes: 0 }
  }

  // 启动定时器
  const startTimers = () => {
    monitorTimer = setInterval(loadMonitorStats, 5000)
    chartTimer = setInterval(loadChartData, 10000)
  }

  // 停止定时器
  const stopTimers = () => {
    if (monitorTimer) clearInterval(monitorTimer)
    if (chartTimer) clearInterval(chartTimer)
  }

  onMounted(() => {
    if (serverId.value) {
      loadMonitorStats()
      loadChartData()
      startTimers()
    }
  })

  onUnmounted(() => {
    stopTimers()
  })

  return {
    monitorData,
    chartData,
    chartType,
    chartTarget,
    currentNetStats,
    currentIOStats,
    loadMetric,
    cpuMetric,
    memoryMetric,
    diskMetric,
    loadMonitorStats,
    loadChartData,
    resetData,
    startTimers,
    stopTimers
  }
}

import type { Ref } from 'vue'
