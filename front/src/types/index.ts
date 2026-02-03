// 全局类型定义

// 服务器类型
export interface Server {
  id: number
  hostname: string
  ip_address: string
  status: 'online' | 'offline' | 'unknown'
}

// 应用类型
export interface Application {
  id: number
  name: string
  version: string
  status: 'running' | 'stopped' | 'error'
  icon: string
  color: string
}

// 监控数据类型
export interface MonitorData {
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
    partitions: DiskPartition[]
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

export interface DiskPartition {
  available: number
  device: string
  fsType: string
  mountPoint: string
  total: number
  usage: number
  used: number
}

export interface ProcessInfo {
  pid: number
  name: string
  cpuPercent: number
  memoryMB: number
  memoryPercent: number
  status: string
  createTime: number
}

// 系统信息类型
export interface SystemInfo {
  hostname: string
  os: string
  kernel: string
  arch: string
  ip: string
  bootTime: string
  uptime: string
}

// 用户类型
export interface User {
  id: number
  username: string
  avatar?: string
  role: string
}

// API 响应类型
export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

// 图表数据点
export interface ChartDataPoint {
  time: string
  value1: number
  value2: number
}
