<template>
  <div class="overview-grid">
    <StatsCard :stats="baseStats" />
  </div>

  <div class="overview-page">
    <ServerList 
      :servers="serverList" 
      :current-server-id="currentServerId"
      @switch="handleServerSwitch"
    />

    <main class="main-area">
      <Loading v-if="loading" :text="$t('overview.loading')" />
      <template v-else>
        
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

        <ChartPanel
          v-model:chart-type="chartType"
          v-model:chart-target="chartTarget"
          :chart-data="chartData"
          :current-net-stats="currentNetStats"
          :current-i-o-stats="currentIOStats"
          :chart-target-list="chartTargetList"
        />
      </template>
    </main>

    <aside v-if="!loading" class="info-sidebar">
      <SystemInfo :info="systemInfo" />
      <AppList :apps="appList" />
    </aside>

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
import { ref, watch, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import Loading from '@/components/Loading/index.vue'
import ServerList from './components/ServerList.vue'
import StatsCard from './components/StatsCard.vue'
import MetricsPanel from './components/MetricsPanel.vue'
import ChartPanel from './components/ChartPanel.vue'
import SystemInfo from './components/SystemInfo.vue'
import AppList from './components/AppList.vue'
import ProcessModal from './components/ProcessModal.vue'
import { useServers, useMonitor, useSystemInfo, useApplications } from './composables'
import { useLoading } from '@/composables/useLoading'
import type { ProcessInfo } from '@/types'

const { loading, withLoading } = useLoading()

// 基础统计数据
const baseStats = ref({
  script: 1,
  deployment: 2
})

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
  loadChartData,
  restoreFromCache,
  hasValidCache
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
  if (hasValidCache(serverId)) {
    restoreFromCache(serverId)
  }
  await loadSystemInfo()
  loadMonitorStats()
  loadChartData()
}

const initData = async () => {
  await withLoading(async () => {
    await loadServers()
    if (currentServerId.value) {
      if (hasValidCache(currentServerId.value)) {
        restoreFromCache(currentServerId.value)
      }
      await loadSystemInfo()
      await loadMonitorStats()
      await loadChartData()
      startTimers()
    }
  })
}

watch(currentServerId, async (newId, oldId) => {
  if (newId && newId !== oldId) {
    stopTimers()
    if (hasValidCache(newId)) {
      restoreFromCache(newId)
    }
    await loadSystemInfo()
    await loadMonitorStats()
    await loadChartData()
    startTimers()
  }
})

onMounted(() => {
  initData()
})

onUnmounted(() => {
  stopTimers()
})
</script>

<style scoped>
.overview-grid {
  margin-bottom: 14px;
}

.overview-page {
  display: flex;
  height: 100%;
  gap: 16px;
}

.main-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 16px;
  overflow-y: auto;
  position: relative;
}

.info-sidebar {
  width: 320px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.main-area :deep(.s-loading) {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}
</style>
