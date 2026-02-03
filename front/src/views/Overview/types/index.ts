// Overview 页面类型定义
import type { Server, Application, MonitorData, SystemInfo, ProcessInfo, ChartDataPoint } from '@/types'

export type {
  Server,
  Application,
  MonitorData,
  SystemInfo,
  ProcessInfo,
  ChartDataPoint
}

// 基础统计数据
export interface BaseStats {
  website: number
  database: number
  cron: number
  installedApps: number
}

// 监控指标
export interface MetricData {
  usage: number
  status?: string
  used?: string | number
  total?: string | number
}

// 图表配置
export type ChartType = 'net' | 'io'

export interface ChartStats {
  bytesSent?: number
  bytesRecv?: number
  readBytes?: number
  writeBytes?: number
}
