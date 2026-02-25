// 监控页面类型定义

// 基础监控历史记录
export interface BaseMonitorRecord {
  id: number
  cpu_usage: number
  memory_usage: number
  memory_total: number
  memory_used: number
  disk_usage: number
  disk_total: number
  disk_used: number
  collect_time: string
}

// 磁盘IO历史记录
export interface DiskIORecord {
  id: number
  disk_name: string
  read_count: number
  write_count: number
  read_bytes: number
  write_bytes: number
  read_time: number
  write_time: number
  io_time: number
  weighted_io_time: number
  iops_in_progress: number
  collect_time: string
}

// 网络IO历史记录
export interface NetworkIORecord {
  id: number
  interface_name: string
  bytes_sent: number
  bytes_recv: number
  packets_sent: number
  packets_recv: number
  err_in: number
  err_out: number
  drop_in: number
  drop_out: number
  fifo_in: number
  fifo_out: number
  collect_time: string
}

// 分页数据
export interface PageData<T> {
  list: T[]
  total: number
  page: number
  size: number
}

// 时间范围
export type TimeRange = '1h' | '6h' | '24h' | '7d' | '30d'

// 图表数据点
export interface ChartDataPoint {
  time: string
  value1: number
  value2?: number
}

// 磁盘IO速率记录（前端计算后）
export interface DiskIOSpeedRecord {
  collect_time: string
  read_speed: number      // bytes/s
  write_speed: number     // bytes/s
}

// 网络IO速率记录（前端计算后）
export interface NetworkIOSpeedRecord {
  collect_time: string
  upload_speed: number    // bytes/s
  download_speed: number  // bytes/s
}

// 设备列表
export interface DeviceList {
  disks: string[]
  interfaces: string[]
}
