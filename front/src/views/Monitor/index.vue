<template>
  <div class="monitor-page">
    <div class="page-header">
      <h1 class="page-title">{{ $t('monitor.title') }}</h1>
      <ServerSelector
        v-model="selectedServer"
        :servers="servers"
        @change="loadMonitorData"
      />
    </div>

    <Loading v-if="loading && !baseData.length" :text="$t('common.loading')" />

    <template v-else-if="selectedServer">
      <TimeRangeSelector v-model="timeRange" />
      <div class="charts-row">
        <CPUMonitorChart :data="baseData" :time-range="timeRange" />
        <MemoryMonitorChart :data="baseData" :time-range="timeRange" />
        <DiskUsageChart :data="diskUsageData" :mount-points="mountPoints" :time-range="timeRange" />
      </div>
      <div class="charts-row io-row">
        <DiskIOChart :data="diskData" :devices="diskDevices" :time-range="timeRange" />
        <NetIOChart :data="netData" :interfaces="netInterfaces" :time-range="timeRange" />
      </div>
    </template>

    <Empty
      v-else
      :description="$t('monitor.noServerSelected')"
      icon="lucide:monitor"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Loading from '@/components/Loading/index.vue'
import Empty from '@/components/Empty/index.vue'
import ServerSelector from './components/ServerSelector.vue'
import CPUMonitorChart from './components/CPUMonitorChart.vue'
import MemoryMonitorChart from './components/MemoryMonitorChart.vue'
import DiskUsageChart from './components/DiskUsageChart.vue'
import DiskIOChart from './components/DiskIOChart.vue'
import NetIOChart from './components/NetIOChart.vue'
import TimeRangeSelector from './components/TimeRangeSelector.vue'
import {
  fetchBaseMonitorHistory,
  fetchDiskIOHistory,
  fetchNetIOHistory,
  fetchDiskUsageHistory
} from '@/api/monitor'
import { fetchServers } from '@/api/server'
import type {
  BaseMonitorRecord,
  DiskIORecord,
  NetworkIORecord,
  DiskUsageRecord,
  TimeRange
} from '@/types/monitor'
import type { Server } from '@/types'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const selectedServer = ref<number | null>(null)
const timeRange = ref<TimeRange>('1h')
const servers = ref<Server[]>([])
const baseData = ref<BaseMonitorRecord[]>([])
const diskUsageData = ref<DiskUsageRecord[]>([])
const diskData = ref<DiskIORecord[]>([])
const netData = ref<NetworkIORecord[]>([])
const mountPoints = ref<string[]>([])
const diskDevices = ref<string[]>([])
const netInterfaces = ref<string[]>([])

const loadServers = async () => {
  try {
    servers.value = await fetchServers()
    if (servers.value.length > 0) {
      const queryServerId = Number(route.query.serverId)
      if (queryServerId && servers.value.some(s => s.id === queryServerId)) {
        selectedServer.value = queryServerId
      } else if (!selectedServer.value) {
        selectedServer.value = servers.value[0].id
      }
    }
  } catch (error) {
    console.error('Failed to load servers:', error)
  }
}

const loadMonitorData = async () => {
  if (!selectedServer.value) return

  loading.value = true
  try {
    const [base, disk, net, diskUsage] = await Promise.all([
      fetchBaseMonitorHistory(selectedServer.value, timeRange.value),
      fetchDiskIOHistory(selectedServer.value, timeRange.value),
      fetchNetIOHistory(selectedServer.value, timeRange.value),
      fetchDiskUsageHistory(selectedServer.value, timeRange.value)
    ])

    baseData.value = base || []
    diskData.value = disk || []
    netData.value = net || []
    diskUsageData.value = diskUsage || []

    mountPoints.value = [...new Set(diskUsageData.value.map((d: DiskUsageRecord) => d.mount_point))]
    diskDevices.value = [...new Set(diskData.value.map(d => d.disk_name))]
    netInterfaces.value = [...new Set(netData.value.map(d => d.interface_name))]
  } catch (error) {
    console.error('Failed to load monitor data:', error)
  } finally {
    loading.value = false
  }
}

watch(selectedServer, (newVal) => {
  if (newVal) {
    router.replace({ query: { serverId: String(newVal) } })
    loadMonitorData()
  }
})

watch(timeRange, () => {
  loadMonitorData()
})

onMounted(() => {
  loadServers()
})
</script>

<style scoped lang="scss">
.monitor-page {
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  height: 100%;
  overflow-y: auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.page-title {
  font-size: 18px;
  font-weight: 600;
  color: #1e3a5f;
  margin: 0;
}

.charts-row {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
}

.charts-row.io-row {
  grid-template-columns: repeat(2, 1fr);
}

@media (max-width: 1200px) {
  .charts-row {
    grid-template-columns: 1fr;
  }
}
</style>
