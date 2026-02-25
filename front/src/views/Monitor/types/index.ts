// Monitor 页面类型定义
import type { BaseMonitorRecord, DiskIORecord, NetworkIORecord, TimeRange } from '@/types/monitor'

export type { BaseMonitorRecord, DiskIORecord, NetworkIORecord, TimeRange }

// 监控状态
export interface MonitorState {
  loading: boolean
  error: string | null
  selectedServer: number | null
  timeRange: TimeRange
  baseData: BaseMonitorRecord[]
  diskData: DiskIORecord[]
  netData: NetworkIORecord[]
  diskDevices: string[]
  netInterfaces: string[]
}
