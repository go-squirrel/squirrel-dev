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
      <div class="charts-grid">
        <CPUMonitorChart :data="baseData" />
        <MemoryMonitorChart :data="baseData" />
        <DiskIOChart :data="diskData" :devices="diskDevices" />
        <NetIOChart :data="netData" :interfaces="netInterfaces" />
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
import Loading from '@/components/Loading/index.vue'
import Empty from '@/components/Empty/index.vue'
import ServerSelector from './components/ServerSelector.vue'
import CPUMonitorChart from './components/CPUMonitorChart.vue'
import MemoryMonitorChart from './components/MemoryMonitorChart.vue'
import DiskIOChart from './components/DiskIOChart.vue'
import NetIOChart from './components/NetIOChart.vue'
import TimeRangeSelector from './components/TimeRangeSelector.vue'
import {
  fetchBaseMonitorHistory,
  fetchDiskIOHistory,
  fetchNetIOHistory
} from '@/api/monitor'
import { fetchServers } from '@/api/server'
import type {
  BaseMonitorRecord,
  DiskIORecord,
  NetworkIORecord,
  TimeRange
} from '@/types/monitor'
import type { Server } from '@/types'

const loading = ref(false)
const selectedServer = ref<number | null>(null)
const timeRange = ref<TimeRange>('1h')
const servers = ref<Server[]>([])
const baseData = ref<BaseMonitorRecord[]>([])
const diskData = ref<DiskIORecord[]>([])
const netData = ref<NetworkIORecord[]>([])
const diskDevices = ref<string[]>([])
const netInterfaces = ref<string[]>([])

const loadServers = async () => {
  try {
    servers.value = await fetchServers()
    if (servers.value.length > 0 && !selectedServer.value) {
      selectedServer.value = servers.value[0].id
      await loadMonitorData()
    }
  } catch (error) {
    console.error('Failed to load servers:', error)
  }
}

const loadMonitorData = async () => {
  if (!selectedServer.value) return

  loading.value = true
  try {
    const [base, disk, net] = await Promise.all([
      fetchBaseMonitorHistory(selectedServer.value, 1, 100),
      fetchDiskIOHistory(selectedServer.value, 1, 100),
      fetchNetIOHistory(selectedServer.value, 1, 100)
    ])

    baseData.value = base.list || []
    diskData.value = disk.list || []
    netData.value = net.list || []

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

.charts-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
  flex: 1;
}

@media (max-width: 1200px) {
  .charts-grid {
    grid-template-columns: 1fr;
  }
}
</style>
