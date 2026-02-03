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

    <!-- 左侧服务器列表 -->
    <aside class="server-sidebar">
      <div class="sidebar-header">
        <h3>服务器列表</h3>
      </div>
      <div class="server-list">
        <div
          v-for="server in serverList"
          :key="server.id"
          class="server-item"
          :class="{ active: currentServerId === server.id }"
          @click="switchServer(server.id)"
        >
          <div class="server-icon">
            <Icon icon="lucide:server" />
          </div>
          <div class="server-info">
            <span class="server-name">{{ server.hostname }}</span>
            <span class="server-ip">{{ server.ip_address }}</span>
          </div>
          <div class="server-status" :class="getStatusClass(server.status)"></div>
        </div>
      </div>
    </aside>

    <!-- 主内容区 -->
    <main class="main-area">
      <!-- 基础数据展示 -->
      <section class="stats-section">
        <div class="stats-grid">
          <div class="stat-card">
            <div class="stat-icon blue">
              <Icon icon="lucide:globe" />
            </div>
            <div class="stat-content">
              <span class="stat-value">{{ baseStats.website || 0 }}</span>
              <span class="stat-label">网站</span>
            </div>
          </div>
          <div class="stat-card">
            <div class="stat-icon purple">
              <Icon icon="lucide:database" />
            </div>
            <div class="stat-content">
              <span class="stat-value">{{ baseStats.database || 0 }}</span>
              <span class="stat-label">数据库 · 所有</span>
            </div>
          </div>
          <div class="stat-card">
            <div class="stat-icon orange">
              <Icon icon="lucide:clock" />
            </div>
            <div class="stat-content">
              <span class="stat-value">{{ baseStats.cron || 0 }}</span>
              <span class="stat-label">计划任务</span>
            </div>
          </div>
          <div class="stat-card">
            <div class="stat-icon green">
              <Icon icon="lucide:package" />
            </div>
            <div class="stat-content">
              <span class="stat-value">{{ baseStats.installedApps || 0 }}</span>
              <span class="stat-label">已安装应用</span>
            </div>
          </div>
        </div>
      </section>

      <!-- 监控指标展示 -->
      <section class="monitor-section">
        <h3 class="section-title">
          <Icon icon="lucide:activity" />
          <span>状态</span>
        </h3>
        <div class="metrics-grid">
          <!-- 负载 -->
          <div class="metric-item" @mouseenter="showTooltip('load')" @mouseleave="hideTooltip">
            <div class="metric-circle">
              <svg viewBox="0 0 100 100">
                <circle class="circle-bg" cx="50" cy="50" r="42" />
                <circle
                  class="circle-progress"
                  cx="50"
                  cy="50"
                  r="42"
                  :stroke-dasharray="getCircleProgress(loadMetric.usage)"
                />
              </svg>
              <div class="metric-value">
                <span class="value">{{ loadMetric.usage.toFixed(2) }}</span>
              </div>
            </div>
            <span class="metric-label">负载</span>
            <span class="metric-sub">{{ loadMetric.status }}</span>

            <!-- 负载悬停详情 -->
            <div v-if="activeTooltip === 'load'" class="tooltip-panel">
              <div class="tooltip-header">负载详情</div>
              <div class="tooltip-content">
                <div class="load-item">
                  <span>1分钟负载</span>
                  <span class="load-value">{{ monitorData.loadAverage?.load1 || 0 }}</span>
                </div>
                <div class="load-item">
                  <span>5分钟负载</span>
                  <span class="load-value">{{ monitorData.loadAverage?.load5 || 0 }}</span>
                </div>
                <div class="load-item">
                  <span>15分钟负载</span>
                  <span class="load-value">{{ monitorData.loadAverage?.load15 || 0 }}</span>
                </div>
              </div>
            </div>
          </div>

          <!-- CPU -->
          <div class="metric-item" @mouseenter="showTooltip('cpu')" @mouseleave="hideTooltip">
            <div class="metric-circle">
              <svg viewBox="0 0 100 100">
                <circle class="circle-bg" cx="50" cy="50" r="42" />
                <circle
                  class="circle-progress"
                  cx="50"
                  cy="50"
                  r="42"
                  :stroke-dasharray="getCircleProgress(cpuMetric.usage)"
                />
              </svg>
              <div class="metric-value">
                <span class="value">{{ cpuMetric.usage.toFixed(2) }}</span>
                <span class="unit">%</span>
              </div>
            </div>
            <span class="metric-label">CPU</span>
            <span class="metric-sub">( {{ cpuMetric.used }} / {{ cpuMetric.total }} ) 核</span>

            <!-- CPU悬停详情 -->
            <div v-if="activeTooltip === 'cpu'" class="tooltip-panel">
              <div class="tooltip-header">CPU 详情</div>
              <div class="tooltip-content">
                <div class="cpu-info">
                  <span>型号: {{ monitorData.cpu?.model || '-' }}</span>
                  <span>频率: {{ ((monitorData.cpu?.frequency || 0) / 1000).toFixed(2) }} GHz</span>
                  <span>核心数: {{ monitorData.cpu?.cores || 0 }} 核</span>
                </div>
                <div class="cpu-cores">
                  <div
                    v-for="(usage, index) in monitorData.cpu?.perCoreUsage || []"
                    :key="index"
                    class="core-item"
                  >
                    <span>核心 {{ index + 1 }}</span>
                    <div class="core-bar">
                      <div class="core-fill" :style="{ width: usage + '%' }"></div>
                    </div>
                    <span>{{ usage.toFixed(1) }}%</span>
                  </div>
                </div>
                <div class="tooltip-actions">
                  <button @click="showProcessList('cpu')">查看使用率前五进程</button>
                </div>
              </div>
            </div>
          </div>

          <!-- 内存 -->
          <div class="metric-item" @mouseenter="showTooltip('memory')" @mouseleave="hideTooltip">
            <div class="metric-circle">
              <svg viewBox="0 0 100 100">
                <circle class="circle-bg" cx="50" cy="50" r="42" />
                <circle
                  class="circle-progress"
                  cx="50"
                  cy="50"
                  r="42"
                  :stroke-dasharray="getCircleProgress(memoryMetric.usage)"
                />
              </svg>
              <div class="metric-value">
                <span class="value">{{ memoryMetric.usage.toFixed(2) }}</span>
                <span class="unit">%</span>
              </div>
            </div>
            <span class="metric-label">内存</span>
            <span class="metric-sub">{{ memoryMetric.used }} / {{ memoryMetric.total }}</span>

            <!-- 内存悬停详情 -->
            <div v-if="activeTooltip === 'memory'" class="tooltip-panel">
              <div class="tooltip-header">内存详情</div>
              <div class="tooltip-content">
                <div class="memory-info">
                  <div class="mem-item">
                    <span>总计</span>
                    <span>{{ formatBytes(monitorData.memory?.total || 0) }}</span>
                  </div>
                  <div class="mem-item">
                    <span>已用</span>
                    <span>{{ formatBytes(monitorData.memory?.used || 0) }}</span>
                  </div>
                  <div class="mem-item">
                    <span>可用</span>
                    <span>{{ formatBytes(monitorData.memory?.available || 0) }}</span>
                  </div>
                  <div class="mem-item">
                    <span>Swap 总计</span>
                    <span>{{ formatBytes(monitorData.memory?.swapTotal || 0) }}</span>
                  </div>
                  <div class="mem-item">
                    <span>Swap 已用</span>
                    <span>{{ formatBytes(monitorData.memory?.swapUsed || 0) }}</span>
                  </div>
                </div>
                <div class="tooltip-actions">
                  <button @click="showProcessList('memory')">查看使用率前五进程</button>
                </div>
              </div>
            </div>
          </div>

          <!-- 磁盘 -->
          <div class="metric-item" @mouseenter="showTooltip('disk')" @mouseleave="hideTooltip">
            <div class="metric-circle">
              <svg viewBox="0 0 100 100">
                <circle class="circle-bg" cx="50" cy="50" r="42" />
                <circle
                  class="circle-progress"
                  cx="50"
                  cy="50"
                  r="42"
                  :stroke-dasharray="getCircleProgress(diskMetric.usage)"
                />
              </svg>
              <div class="metric-value">
                <span class="value">{{ diskMetric.usage.toFixed(2) }}</span>
                <span class="unit">%</span>
              </div>
            </div>
            <span class="metric-label">磁盘</span>
            <span class="metric-sub">{{ diskMetric.used }} / {{ diskMetric.total }}</span>

            <!-- 磁盘悬停详情 -->
            <div v-if="activeTooltip === 'disk'" class="tooltip-panel">
              <div class="tooltip-header">磁盘详情</div>
              <div class="tooltip-content">
                <div class="disk-list">
                  <div
                    v-for="partition in monitorData.disk?.partitions || []"
                    :key="partition.device"
                    class="disk-item"
                  >
                    <div class="disk-header">
                      <span>{{ partition.device }}</span>
                      <span>{{ partition.mountPoint }}</span>
                    </div>
                    <div class="disk-bar">
                      <div class="disk-fill" :style="{ width: partition.usage + '%' }"></div>
                    </div>
                    <div class="disk-info">
                      <span>{{ formatBytes(partition.used) }} / {{ formatBytes(partition.total) }}</span>
                      <span>{{ partition.usage.toFixed(1) }}%</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- 监控图表 -->
      <section class="chart-section">
        <div class="chart-header">
          <h3 class="section-title">
            <Icon icon="lucide:bar-chart-2" />
            <span>监控</span>
          </h3>
          <div class="chart-controls">
            <select v-model="chartType" class="control-select">
              <option value="net">网卡流量</option>
              <option value="io">磁盘 IO</option>
            </select>
            <select v-model="chartTarget" class="control-select">
              <option value="all">所有</option>
              <option v-for="item in chartTargetList" :key="item" :value="item">{{ item }}</option>
            </select>
          </div>
        </div>
        <div class="chart-stats">
          <div class="chart-stat-item" v-if="chartType === 'net'">
            <span class="stat-name">上行</span>
            <span class="stat-value">{{ formatSpeed(currentNetStats.bytesSent) }}</span>
          </div>
          <div class="chart-stat-item" v-if="chartType === 'net'">
            <span class="stat-name">下行</span>
            <span class="stat-value">{{ formatSpeed(currentNetStats.bytesRecv) }}</span>
          </div>
          <div class="chart-stat-item" v-if="chartType === 'io'">
            <span class="stat-name">读取</span>
            <span class="stat-value">{{ formatBytes(currentIOStats.readBytes) }}</span>
          </div>
          <div class="chart-stat-item" v-if="chartType === 'io'">
            <span class="stat-name">写入</span>
            <span class="stat-value">{{ formatBytes(currentIOStats.writeBytes) }}</span>
          </div>
        </div>
        <div class="chart-container" ref="chartContainer">
          <canvas ref="chartCanvas"></canvas>
        </div>
      </section>
    </main>

    <!-- 右侧信息区 -->
    <aside class="info-sidebar">
      <!-- 系统信息 -->
      <div class="info-card">
        <div class="info-header">
          <h3>系统信息</h3>
          <div class="info-actions">
            <button class="icon-btn"><Icon icon="lucide:refresh-cw" /></button>
            <button class="icon-btn"><Icon icon="lucide:maximize-2" /></button>
          </div>
        </div>
        <div class="info-content">
          <div class="info-item">
            <span class="info-label">主机名称</span>
            <span class="info-value">{{ systemInfo.hostname }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">发行版本</span>
            <span class="info-value">{{ systemInfo.os }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">内核版本</span>
            <span class="info-value">{{ systemInfo.kernel }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">系统类型</span>
            <span class="info-value">{{ systemInfo.arch }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">主机地址</span>
            <span class="info-value">{{ systemInfo.ip }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">启动时间</span>
            <span class="info-value">{{ systemInfo.bootTime }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">运行时间</span>
            <span class="info-value">{{ systemInfo.uptime }}</span>
          </div>
        </div>
      </div>

      <!-- 应用列表 -->
      <div class="info-card">
        <div class="info-header">
          <h3>应用</h3>
          <button class="icon-btn"><Icon icon="lucide:settings" /></button>
        </div>
        <div class="app-list">
          <div v-for="app in appList" :key="app.id" class="app-item">
            <div class="app-icon" :style="{ backgroundColor: app.color }">
              <Icon :icon="app.icon" />
            </div>
            <div class="app-info">
              <div class="app-name-row">
                <span class="app-name">{{ app.name }}</span>
                <Icon icon="lucide:chevron-down" class="app-expand" />
              </div>
              <span class="app-version">版本: {{ app.version }}</span>
              <div class="app-actions">
                <button class="app-action-btn">关闭</button>
                <button class="app-action-btn">重启</button>
                <button class="app-action-btn">目录</button>
              </div>
            </div>
            <div class="app-status">
              <span class="status-badge" :class="app.status">{{ getAppStatusText(app.status) }}</span>
            </div>
          </div>
        </div>
      </div>
    </aside>

    <!-- 进程列表弹窗 -->
    <div v-if="showProcessModal" class="modal-overlay" @click="closeProcessModal">
      <div class="modal" @click.stop>
        <div class="modal-header">
          <h3>{{ processModalTitle }}</h3>
          <button class="btn-close" @click="closeProcessModal">
            <Icon icon="lucide:x" />
          </button>
        </div>
        <div class="modal-body">
          <table class="process-table">
            <thead>
              <tr>
                <th>PID</th>
                <th>名称</th>
                <th>CPU%</th>
                <th>内存(MB)</th>
                <th>状态</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="proc in processList" :key="proc.pid">
                <td>{{ proc.pid }}</td>
                <td>{{ proc.name }}</td>
                <td>{{ proc.cpuPercent.toFixed(2) }}%</td>
                <td>{{ proc.memoryMB.toFixed(2) }}</td>
                <td>{{ proc.status }}</td>
                <td>
                  <button class="btn-icon danger" @click="killProcess(proc.pid)">
                    <Icon icon="lucide:trash-2" />
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { Icon } from '@iconify/vue'

// 服务器列表
interface Server {
  id: number
  hostname: string
  ip_address: string
  status: string
}

const serverList = ref<Server[]>([])
const currentServerId = ref<number>(0)

// 基础统计数据
const baseStats = ref({
  website: 1,
  database: 2,
  cron: 1,
  installedApps: 6
})

// 监控数据
interface MonitorData {
  hostname?: string
  cpu?: {
    cores: number
    frequency: number
    model: string
    perCoreUsage: number[]
    usage: number
  }
  memory?: {
    available: number
    swapTotal: number
    swapUsed: number
    total: number
    usage: number
    used: number
  }
  disk?: {
    available: number
    partitions: Array<{
      available: number
      device: string
      fsType: string
      mountPoint: string
      total: number
      usage: number
      used: number
    }>
    total: number
    usage: number
    used: number
  }
  loadAverage?: {
    load1: number
    load5: number
    load15: number
  }
  topCPU?: ProcessInfo[]
  topMemory?: ProcessInfo[]
}

interface ProcessInfo {
  pid: number
  name: string
  cpuPercent: number
  memoryMB: number
  memoryPercent: number
  status: string
  createTime: number
}

const monitorData = ref<MonitorData>({})

// 计算监控指标
const loadMetric = computed(() => {
  const usage = monitorData.value.loadAverage?.load1 || 0
  let status = '运行流畅'
  if (usage > 80) status = '负载过高'
  else if (usage > 50) status = '负载中等'
  return { usage, status }
})

const cpuMetric = computed(() => {
  // API 返回的 usage 是小数（如 0.4975 表示 49.75%）
  const rawUsage = monitorData.value.cpu?.usage || 0
  // 如果 rawUsage > 1，说明 API 返回的是百分比形式（如 49.75）
  // 如果 rawUsage <= 1，说明 API 返回的是小数形式（如 0.4975）
  const usage = rawUsage > 1 ? rawUsage : rawUsage 
  const cores = monitorData.value.cpu?.cores || 0
  // 计算使用的核心数
  const used = ((usage / 100) * cores).toFixed(2)
  return { usage, used, total: cores }
})

const memoryMetric = computed(() => {
  // API 返回的 usage 可能是小数（如 0.219）或百分比（如 21.9）
  const rawUsage = monitorData.value.memory?.usage || 0
  const usage = rawUsage > 1 ? rawUsage : rawUsage
  const used = formatBytes(monitorData.value.memory?.used || 0)
  const total = formatBytes(monitorData.value.memory?.total || 0)
  return { usage, used, total }
})

const diskMetric = computed(() => {
  // API 返回的 usage 可能是小数（如 0.1148）或百分比（如 11.48）
  const rawUsage = monitorData.value.disk?.usage || 0
  const usage = rawUsage > 1 ? rawUsage : rawUsage 
  const used = formatBytes(monitorData.value.disk?.used || 0)
  const total = formatBytes(monitorData.value.disk?.total || 0)
  return { usage, used, total }
})

// 圆形进度条计算
const getCircleProgress = (percentage: number) => {
  const radius = 42
  const circumference = 2 * Math.PI * radius
  const progress = Math.min(percentage, 100) / 100
  const dasharray = `${circumference * progress} ${circumference}`
  return dasharray
}

// 工具函数
function formatBytes(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

function formatSpeed(bytes: number): string {
  return formatBytes(bytes) + '/s'
}

function getStatusClass(status: string): string {
  switch (status) {
    case 'online': return 'online'
    case 'offline': return 'offline'
    default: return 'unknown'
  }
}

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

const showProcessList = (type: 'cpu' | 'memory') => {
  if (type === 'cpu') {
    processModalTitle.value = 'CPU 使用率前五进程'
    processList.value = monitorData.value.topCPU || []
  } else {
    processModalTitle.value = '内存使用率前五进程'
    processList.value = monitorData.value.topMemory || []
  }
  showProcessModal.value = true
}

const closeProcessModal = () => {
  showProcessModal.value = false
}

const killProcess = async (pid: number) => {
  // TODO: 实现杀死进程 API
  console.log('Kill process:', pid)
}

// 图表相关
const chartType = ref<'net' | 'io'>('net')
const chartTarget = ref('all')
const chartTargetList = ref<string[]>([])
const chartContainer = ref<HTMLDivElement>()
const chartCanvas = ref<HTMLCanvasElement>()

interface ChartDataPoint {
  time: string
  value1: number
  value2: number
}

const chartData = ref<ChartDataPoint[]>([])
const currentNetStats = ref({ bytesSent: 0, bytesRecv: 0 })
const currentIOStats = ref({ readBytes: 0, writeBytes: 0 })

// 上一次的数据（用于计算速度）
let lastNetStats = { bytesSent: 0, bytesRecv: 0, timestamp: 0 }
let lastIOStats = { readBytes: 0, writeBytes: 0, timestamp: 0 }

// 系统信息
const systemInfo = ref({
  hostname: '-',
  os: '-',
  kernel: '-',
  arch: '-',
  ip: '-',
  bootTime: '-',
  uptime: '-'
})

// 应用列表
interface App {
  id: number
  name: string
  version: string
  status: string
  icon: string
  color: string
}

const appList = ref<App[]>([])

function getAppStatusText(status: string): string {
  const map: Record<string, string> = {
    running: '运行中',
    stopped: '已停止',
    error: '错误'
  }
  return map[status] || status
}

// API 请求函数
const API_BASE = import.meta.env.DEV ? 'http://192.168.37.20:10700/api/v1' : '/api/v1'

async function fetchServers() {
  try {
    const response = await fetch(`${API_BASE}/server`)
    const result = await response.json()
    if (result.code === 0) {
      serverList.value = result.data
      if (serverList.value.length > 0 && currentServerId.value === 0) {
        currentServerId.value = serverList.value[0].id
        // 获取到服务器后立即加载监控数据和系统信息
        await fetchServerDetail()
        await fetchMonitorStats()
        await fetchChartData()
      }
    }
  } catch (error) {
    console.error('Failed to fetch servers:', error)
  }
}

async function fetchMonitorStats() {
  if (!currentServerId.value) return
  try {
    const response = await fetch(`${API_BASE}/monitor/stats/${currentServerId.value}`)
    const result = await response.json()
    if (result.code === 0) {
      monitorData.value = result.data
    }
  } catch (error) {
    console.error('Failed to fetch monitor stats:', error)
  }
}

// 获取服务器详细信息
async function fetchServerDetail() {
  if (!currentServerId.value) return
  try {
    const response = await fetch(`${API_BASE}/server/${currentServerId.value}`)
    const result = await response.json()
    if (result.code === 0 && result.data.server_info) {
      const info = result.data.server_info
      systemInfo.value.hostname = info.hostname || '-'
      systemInfo.value.os = `${info.platform || '-'} ${info.platformVersion || ''}`.trim() || '-'
      systemInfo.value.kernel = info.kernelVersion || '-'
      systemInfo.value.arch = info.architecture || '-'
      systemInfo.value.uptime = info.uptimeStr || '-'

      // 获取第一个 IPv4 地址作为主 IP
      if (info.ipAddresses && info.ipAddresses.length > 0) {
        const firstInterface = info.ipAddresses[0]
        if (firstInterface.ipv4 && firstInterface.ipv4.length > 0) {
          systemInfo.value.ip = firstInterface.ipv4[0]
        } else {
          systemInfo.value.ip = '-'
        }
      } else {
        systemInfo.value.ip = '-'
      }
    }
  } catch (error) {
    console.error('Failed to fetch server detail:', error)
  }
}

async function fetchApplications() {
  try {
    const response = await fetch(`${API_BASE}/application`)
    const result = await response.json()
    if (result.code === 0) {
      appList.value = result.data.map((app: any) => ({
        id: app.id,
        name: app.name,
        version: app.version || '1.0.0',
        status: app.status || 'stopped',
        icon: 'lucide:box',
        color: '#4fc3f7'
      }))
      baseStats.value.installedApps = appList.value.length
    }
  } catch (error) {
    console.error('Failed to fetch applications:', error)
  }
}

async function fetchChartData() {
  if (!currentServerId.value) return
  try {
    if (chartType.value === 'net') {
      const target = chartTarget.value
      const response = await fetch(`${API_BASE}/monitor/stats/net/${currentServerId.value}/${target}`)
      const result = await response.json()
      if (result.code === 0) {
        const now = Date.now()
        const { bytesSent, bytesRecv } = result.data

        // 计算速度（字节/秒）
        if (lastNetStats.timestamp > 0) {
          const timeDiff = (now - lastNetStats.timestamp) / 1000 // 秒
          const sentSpeed = timeDiff > 0 ? Math.max(0, (bytesSent - lastNetStats.bytesSent) / timeDiff) : 0
          const recvSpeed = timeDiff > 0 ? Math.max(0, (bytesRecv - lastNetStats.bytesRecv) / timeDiff) : 0

          currentNetStats.value = { bytesSent: sentSpeed, bytesRecv: recvSpeed }
          updateChartData(sentSpeed, recvSpeed)
        }

        // 保存当前数据
        lastNetStats = { bytesSent, bytesRecv, timestamp: now }
      }
    } else {
      const target = chartTarget.value
      const response = await fetch(`${API_BASE}/monitor/stats/io/${currentServerId.value}/${target}`)
      const result = await response.json()
      if (result.code === 0) {
        const now = Date.now()
        const { readBytes, writeBytes } = result.data

        // 计算速度（字节/秒）
        if (lastIOStats.timestamp > 0) {
          const timeDiff = (now - lastIOStats.timestamp) / 1000 // 秒
          const readSpeed = timeDiff > 0 ? Math.max(0, (readBytes - lastIOStats.readBytes) / timeDiff) : 0
          const writeSpeed = timeDiff > 0 ? Math.max(0, (writeBytes - lastIOStats.writeBytes) / timeDiff) : 0

          currentIOStats.value = { readBytes: readSpeed, writeBytes: writeSpeed }
          updateChartData(readSpeed, writeSpeed)
        }

        // 保存当前数据
        lastIOStats = { readBytes, writeBytes, timestamp: now }
      }
    }
  } catch (error) {
    console.error('Failed to fetch chart data:', error)
  }
}

function updateChartData(value1: number, value2: number) {
  const now = new Date()
  const time = `${now.getHours().toString().padStart(2, '0')}:${now.getMinutes().toString().padStart(2, '0')}:${now.getSeconds().toString().padStart(2, '0')}`

  chartData.value.push({
    time,
    value1,
    value2
  })

  // 保留最近 30 个数据点
  if (chartData.value.length > 30) {
    chartData.value.shift()
  }

  drawChart()
}

function drawChart() {
  if (!chartCanvas.value || !chartContainer.value) return

  const canvas = chartCanvas.value
  const container = chartContainer.value
  const ctx = canvas.getContext('2d')
  if (!ctx) return

  // 设置 canvas 尺寸
  const rect = container.getBoundingClientRect()
  canvas.width = rect.width
  canvas.height = rect.height

  const width = canvas.width
  const height = canvas.height
  const padding = 40

  // 清空画布
  ctx.clearRect(0, 0, width, height)

  if (chartData.value.length < 2) return

  // 计算数据范围
  const allValues = chartData.value.flatMap(d => [d.value1, d.value2])
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
  chartData.value.forEach((point, index) => {
    const x = padding + (width - 2 * padding) * (index / (chartData.value.length - 1))
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
  chartData.value.forEach((point, index) => {
    const x = padding + (width - 2 * padding) * (index / (chartData.value.length - 1))
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
  const step = Math.ceil(chartData.value.length / 6)
  chartData.value.forEach((point, index) => {
    if (index % step === 0) {
      const x = padding + (width - 2 * padding) * (index / (chartData.value.length - 1))
      ctx.fillText(point.time, x, height - 10)
    }
  })
}

// 切换服务器
async function switchServer(serverId: number) {
  currentServerId.value = serverId
  chartData.value = []
  // 重置统计数据
  lastNetStats = { bytesSent: 0, bytesRecv: 0, timestamp: 0 }
  lastIOStats = { readBytes: 0, writeBytes: 0, timestamp: 0 }
  currentNetStats.value = { bytesSent: 0, bytesRecv: 0 }
  currentIOStats.value = { readBytes: 0, writeBytes: 0 }
  await fetchServerDetail()
  fetchMonitorStats()
  fetchChartData()
}

// 定时器
let monitorTimer: ReturnType<typeof setInterval> | null = null
let chartTimer: ReturnType<typeof setInterval> | null = null

// 监听图表类型变化
watch(chartType, () => {
  chartTarget.value = 'all'
  chartData.value = []
  fetchChartData()
})

watch(chartTarget, () => {
  chartData.value = []
  fetchChartData()
})

// 窗口大小变化时重绘图表
const handleResize = () => {
  nextTick(() => {
    drawChart()
  })
}

onMounted(() => {
  fetchServers()
  fetchApplications()

  // 5秒刷新一次监控数据
  monitorTimer = setInterval(() => {
    fetchMonitorStats()
  }, 5000)

  // 10秒刷新一次图表数据
  chartTimer = setInterval(() => {
    fetchChartData()
  }, 10000)

  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  if (monitorTimer) clearInterval(monitorTimer)
  if (chartTimer) clearInterval(chartTimer)
  window.removeEventListener('resize', handleResize)
})
</script>

<style scoped>
.overview-page {
  display: flex;
  height: 100%;
  gap: 16px;
}

/* 左侧服务器列表 */
.server-sidebar {
  width: 220px;
  background: #ffffff;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.sidebar-header {
  padding: 16px;
  border-bottom: 1px solid #f1f5f9;
}

.sidebar-header h3 {
  font-size: 14px;
  font-weight: 600;
  color: #1e3a5f;
}

.server-list {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
}

.server-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.3s ease;
  margin-bottom: 4px;
}

.server-item:hover {
  background: #f8fafc;
}

.server-item.active {
  background: linear-gradient(135deg, #4fc3f7 0%, #29b6f6 100%);
}

.server-item.active .server-name,
.server-item.active .server-ip {
  color: #ffffff;
}

.server-icon {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f1f5f9;
  border-radius: 6px;
  color: #4fc3f7;
  font-size: 16px;
}

.server-item.active .server-icon {
  background: rgba(255, 255, 255, 0.2);
  color: #ffffff;
}

.server-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.server-name {
  font-size: 13px;
  font-weight: 500;
  color: #1e3a5f;
}

.server-ip {
  font-size: 11px;
  color: #94a3b8;
}

.server-status {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #94a3b8;
}

.server-status.online {
  background: #16a34a;
}

.server-status.offline {
  background: #dc2626;
}

/* 主内容区 */
.main-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 16px;
  overflow-y: auto;
}

/* 基础统计数据 */
.stats-section {
  background: #ffffff;
  border-radius: 8px;
  padding: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: #f8fafc;
  border-radius: 8px;
  transition: all 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.stat-icon {
  width: 44px;
  height: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  font-size: 20px;
}

.stat-icon.blue {
  background: #e0f2fe;
  color: #0284c7;
}

.stat-icon.purple {
  background: #f3e8ff;
  color: #9333ea;
}

.stat-icon.orange {
  background: #ffedd5;
  color: #ea580c;
}

.stat-icon.green {
  background: #dcfce7;
  color: #16a34a;
}

.stat-content {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: #1e3a5f;
}

.stat-label {
  font-size: 12px;
  color: #64748b;
}

/* 监控指标 */
.monitor-section {
  background: #ffffff;
  border-radius: 8px;
  padding: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.section-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  font-weight: 600;
  color: #1e3a5f;
  margin-bottom: 16px;
}

.metrics-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 24px;
}

.metric-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  position: relative;
  padding: 16px;
}

.metric-circle {
  position: relative;
  width: 100px;
  height: 100px;
}

.metric-circle svg {
  width: 100%;
  height: 100%;
  transform: rotate(-90deg);
}

.circle-bg {
  fill: none;
  stroke: #f1f5f9;
  stroke-width: 4;
}

.circle-progress {
  fill: none;
  stroke: url(#gradient);
  stroke-width: 4;
  stroke-linecap: round;
  transition: stroke-dasharray 0.5s ease;
}

.metric-value {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  display: flex;
  align-items: baseline;
  gap: 2px;
}

.metric-value .value {
  font-size: 20px;
  font-weight: 700;
  color: #1e3a5f;
}

.metric-value .unit {
  font-size: 12px;
  color: #64748b;
}

.metric-label {
  font-size: 13px;
  font-weight: 500;
  color: #1e3a5f;
}

.metric-sub {
  font-size: 11px;
  color: #94a3b8;
}

/* 悬停提示面板 */
.tooltip-panel {
  position: absolute;
  top: 100%;
  left: 50%;
  transform: translateX(-50%);
  width: 280px;
  background: #ffffff;
  border-radius: 8px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
  z-index: 100;
  margin-top: 8px;
}

.tooltip-header {
  padding: 12px 16px;
  border-bottom: 1px solid #f1f5f9;
  font-size: 13px;
  font-weight: 600;
  color: #1e3a5f;
}

.tooltip-content {
  padding: 12px 16px;
}

/* 负载详情 */
.load-item {
  display: flex;
  justify-content: space-between;
  padding: 6px 0;
  font-size: 12px;
}

.load-item span:first-child {
  color: #64748b;
}

.load-value {
  font-weight: 500;
  color: #1e3a5f;
}

/* CPU 详情 */
.cpu-info {
  display: flex;
  flex-direction: column;
  gap: 6px;
  margin-bottom: 12px;
  font-size: 12px;
  color: #64748b;
}

.cpu-cores {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.core-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 11px;
}

.core-item span:first-child {
  width: 50px;
  color: #64748b;
}

.core-bar {
  flex: 1;
  height: 4px;
  background: #f1f5f9;
  border-radius: 2px;
  overflow: hidden;
}

.core-fill {
  height: 100%;
  background: linear-gradient(90deg, #4fc3f7 0%, #29b6f6 100%);
  border-radius: 2px;
  transition: width 0.3s ease;
}

.core-item span:last-child {
  width: 40px;
  text-align: right;
  color: #1e3a5f;
}

/* 内存详情 */
.memory-info {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.mem-item {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
}

.mem-item span:first-child {
  color: #64748b;
}

.mem-item span:last-child {
  color: #1e3a5f;
  font-weight: 500;
}

/* 磁盘详情 */
.disk-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.disk-item {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.disk-header {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  font-weight: 500;
  color: #1e3a5f;
}

.disk-bar {
  height: 6px;
  background: #f1f5f9;
  border-radius: 3px;
  overflow: hidden;
}

.disk-fill {
  height: 100%;
  background: linear-gradient(90deg, #4fc3f7 0%, #29b6f6 100%);
  border-radius: 3px;
  transition: width 0.3s ease;
}

.disk-info {
  display: flex;
  justify-content: space-between;
  font-size: 11px;
  color: #64748b;
}

/* 提示面板操作按钮 */
.tooltip-actions {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid #f1f5f9;
}

.tooltip-actions button {
  width: 100%;
  padding: 8px;
  font-size: 12px;
  font-weight: 500;
  color: #4fc3f7;
  background: #e0f2fe;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.tooltip-actions button:hover {
  background: #4fc3f7;
  color: #ffffff;
}

/* 图表区域 */
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

.chart-header .section-title {
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

/* 右侧信息区 */
.info-sidebar {
  width: 320px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.info-card {
  background: #ffffff;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  overflow: hidden;
}

.info-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  border-bottom: 1px solid #f1f5f9;
}

.info-header h3 {
  font-size: 14px;
  font-weight: 600;
  color: #1e3a5f;
}

.info-actions {
  display: flex;
  gap: 6px;
}

.icon-btn {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  background: #f8fafc;
  border-radius: 6px;
  color: #64748b;
  cursor: pointer;
  transition: all 0.3s ease;
}

.icon-btn:hover {
  background: #f1f5f9;
  color: #4fc3f7;
}

.info-content {
  padding: 12px 16px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  padding: 8px 0;
  font-size: 12px;
  border-bottom: 1px solid #f8fafc;
}

.info-item:last-child {
  border-bottom: none;
}

.info-label {
  color: #64748b;
}

.info-value {
  color: #1e3a5f;
  font-weight: 500;
  max-width: 150px;
  text-align: right;
  word-break: break-all;
}

/* 应用列表 */
.app-list {
  padding: 12px 16px;
  max-height: 400px;
  overflow-y: auto;
}

.app-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 12px 0;
  border-bottom: 1px solid #f8fafc;
}

.app-item:last-child {
  border-bottom: none;
}

.app-icon {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  color: #ffffff;
  font-size: 20px;
  flex-shrink: 0;
}

.app-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.app-name-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.app-name {
  font-size: 13px;
  font-weight: 600;
  color: #1e3a5f;
}

.app-expand {
  width: 14px;
  height: 14px;
  color: #94a3b8;
  cursor: pointer;
}

.app-version {
  font-size: 11px;
  color: #94a3b8;
}

.app-actions {
  display: flex;
  gap: 6px;
  margin-top: 4px;
}

.app-action-btn {
  padding: 3px 8px;
  font-size: 11px;
  color: #64748b;
  background: #f8fafc;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.app-action-btn:hover {
  background: #e0f2fe;
  color: #0284c7;
}

.app-status {
  flex-shrink: 0;
}

.status-badge {
  padding: 3px 8px;
  border-radius: 10px;
  font-size: 10px;
  font-weight: 500;
}

.status-badge.running {
  background: #dcfce7;
  color: #16a34a;
}

.status-badge.stopped {
  background: #f1f5f9;
  color: #64748b;
}

.status-badge.error {
  background: #fee2e2;
  color: #dc2626;
}

/* 进程弹窗 */
.modal-overlay {
  position: fixed;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.5);
  z-index: 1000;
  padding: 16px;
}

.modal {
  width: 100%;
  max-width: 600px;
  background: #ffffff;
  border-radius: 8px;
  overflow: hidden;
  max-height: 80vh;
  display: flex;
  flex-direction: column;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid #f1f5f9;
}

.modal-header h3 {
  font-size: 14px;
  font-weight: 600;
  color: #1e3a5f;
}

.btn-close {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border-radius: 6px;
  background: #f8fafc;
  color: #64748b;
  border: none;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn-close:hover {
  background: #f1f5f9;
  color: #1e3a5f;
}

.modal-body {
  padding: 16px;
  overflow-y: auto;
}

.process-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 12px;
}

.process-table th,
.process-table td {
  padding: 10px 12px;
  text-align: left;
  border-bottom: 1px solid #f1f5f9;
}

.process-table th {
  font-weight: 600;
  color: #1e3a5f;
  background: #f8fafc;
}

.process-table td {
  color: #64748b;
}

.btn-icon.danger {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  background: #f8fafc;
  border-radius: 6px;
  color: #64748b;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn-icon.danger:hover {
  background: #fee2e2;
  color: #dc2626;
}
</style>
