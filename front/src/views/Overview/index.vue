<template>
  <div class="overview-page">
    <!-- SVG 渐变定义 -->
    <svg width="0" height="0" style="position: absolute;">
      <defs>
        <linearGradient id="gradient" x1="0%" y1="0%" x2="100%" y2="0%">
          <stop offset="0%" style="stop-color:#4fc3f7;stop-opacity:1" />
          <stop offset="100%" style="stop-color:#29b6f6;stop-opacity:1" />
        </linearGradient>
      </defs>
    </svg>

    <!-- 加载中状态 -->
    <div v-if="loading" class="loading-container">
      <Icon icon="lucide:loader-2" class="spinner" />
      <p>{{ $t('overview.loading') }}</p>
    </div>

    <template v-else>
      <!-- 左侧服务器列表 -->
      <ServerList 
        :servers="serverList" 
        :current-server-id="currentServerId"
        @switch="handleServerSwitch"
      />

      <!-- 主内容区 -->
      <main class="main-area">
        <!-- 基础数据展示 -->
        <StatsCard :stats="baseStats" />

        <!-- 监控指标展示 -->
        <MetricsPanel
          :monitor-data="monitorData"
          :load-metric="loadMetric"
          :cpu-metric="cpuMetric"
          :memory-metric="memoryMetric"
          :disk-metric="diskMetric"
          :active-tooltip="activeTooltip"
          @show-tooltip="showTooltip"
          @hide-tooltip="hideTooltip"
          @show-process="showProcessList"
        />

        <!-- 监控图表 -->
        <ChartPanel
          v-model:chart-type="chartType"
          v-model:chart-target="chartTarget"
          :chart-data="chartData"
          :current-net-stats="currentNetStats"
          :current-i-o-stats="currentIOStats"
          :chart-target-list="chartTargetList"
        />
      </main>

      <!-- 右侧信息区 -->
      <aside class="info-sidebar">
        <SystemInfo :info="systemInfo" />
        <AppList :apps="appList" />
      </aside>
    </template>

    <!-- 进程列表弹窗 -->
    <ProcessModal
      :visible="showProcessModal"
      :title="processModalTitle"
      :processes="processList"
      @close="closeProcessModal"
      @kill="killProcess"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { Icon } from '@iconify/vue'
import ServerList from './components/ServerList.vue'
import StatsCard from './components/StatsCard.vue'
import MetricsPanel from './components/MetricsPanel.vue'
import ChartPanel from './components/ChartPanel.vue'
import SystemInfo from './components/SystemInfo.vue'
import AppList from './components/AppList.vue'
import ProcessModal from './components/ProcessModal.vue'
import { useServers, useMonitor, useSystemInfo, useApplications } from './composables'
import type { ProcessInfo } from '@/types'

// 基础统计数据
const baseStats = ref({
  website: 1,
  database: 2,
  cron: 1,
  installedApps: 6
})

// 页面加载状态
const loading = ref(true)

// 使用 composables
const { serverList, currentServerId, loadServers, switchServer } = useServers()
const { 
  monitorData, 
  chartData, 
  chartType, 
  chartTarget,
  chartTargetList,
  currentNetStats,
  currentIOStats,
  loadMetric, 
  cpuMetric, 
  memoryMetric, 
  diskMetric,
  resetData,
  stopTimers,
  startTimers,
  loadMonitorStats,
  loadChartData
} = useMonitor(currentServerId)
const { systemInfo, loadSystemInfo } = useSystemInfo(currentServerId)
const { appList } = useApplications()

// 悬停提示
const activeTooltip = ref<string | null>(null)
let tooltipTimer: ReturnType<typeof setTimeout> | null = null

const showTooltip = (type: string) => {
  if (tooltipTimer) clearTimeout(tooltipTimer)
  activeTooltip.value = type
}

const hideTooltip = () => {
  tooltipTimer = setTimeout(() => {
    activeTooltip.value = null
  }, 300)
}

// 进程列表弹窗
const showProcessModal = ref(false)
const processModalTitle = ref('')
const processList = ref<ProcessInfo[]>([])

const { t } = useI18n()

const showProcessList = (type: 'cpu' | 'memory') => {
  if (type === 'cpu') {
    processModalTitle.value = t('overview.topCpuProcesses')
    processList.value = monitorData.value.topCPU || []
  } else {
    processModalTitle.value = t('overview.topMemoryProcesses')
    processList.value = monitorData.value.topMemory || []
  }
  showProcessModal.value = true
}

const closeProcessModal = () => {
  showProcessModal.value = false
}

const killProcess = async (pid: number) => {
  console.log('Kill process:', pid)
}

// 切换服务器
const handleServerSwitch = async (serverId: number) => {
  switchServer(serverId)
  resetData()
  await loadSystemInfo()
  loadMonitorStats()
  loadChartData()
}

// 初始化数据
const initData = async () => {
  loading.value = true
  try {
    // 先加载服务器列表
    await loadServers()
    // 如果有服务器，加载监控数据
    if (currentServerId.value) {
      await loadSystemInfo()
      await loadMonitorStats()
      await loadChartData()
      startTimers()
    }
  } catch (error) {
    console.error('Failed to init data:', error)
  } finally {
    loading.value = false
  }
}

// 监听服务器变化
watch(currentServerId, async (newId, oldId) => {
  if (newId && newId !== oldId) {
    stopTimers()
    await loadSystemInfo()
    await loadMonitorStats()
    await loadChartData()
    startTimers()
  }
})

onMounted(() => {
  initData()
})
</script>

<style scoped>
.overview-page {
  display: flex;
  height: 100%;
  gap: 16px;
}

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  color: #64748b;
}

.spinner {
  width: 48px;
  height: 48px;
  animation: spin 1s linear infinite;
  margin-bottom: 16px;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

.main-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 16px;
  overflow-y: auto;
}

.info-sidebar {
  width: 320px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}
</style>
