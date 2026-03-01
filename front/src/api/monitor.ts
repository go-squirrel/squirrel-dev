// 监控相关 API
import { get } from '@/utils/request'
import type { MonitorData } from '@/types'
import type { TimeRange, BaseMonitorRecord, DiskIORecord, NetworkIORecord, DiskUsageRecord } from '@/types/monitor'

/**
 * 获取监控统计
 */
export function fetchMonitorStats(serverId: number): Promise<MonitorData> {
  return get(`/monitor/stats/${serverId}`)
}

/**
 * 获取网络统计
 */
export function fetchNetStats(serverId: number, target: string): Promise<any> {
  return get(`/monitor/stats/net/${serverId}/${target}`)
}

/**
 * 获取磁盘 Io 统计
 */
export function fetchIOStats(serverId: number, target: string): Promise<any> {
  return get(`/monitor/stats/io/${serverId}/${target}`)
}

/**
 * 获取基础监控历史数据（按时间范围）
 * @param serverId 服务器ID
 * @param range 时间范围 (1h, 6h, 24h, 7d)
 */
export function fetchBaseMonitorHistory(
  serverId: number,
  range: TimeRange = '1h'
): Promise<BaseMonitorRecord[]> {
  return get(`/monitor/base/${serverId}?range=${range}`)
}

/**
 * 获取磁盘IO历史数据（按时间范围）
 * @param serverId 服务器ID
 * @param range 时间范围 (1h, 6h, 24h, 7d)
 */
export function fetchDiskIOHistory(
  serverId: number,
  range: TimeRange = '1h'
): Promise<DiskIORecord[]> {
  return get(`/monitor/disk/${serverId}?range=${range}`)
}

/**
 * 获取网络IO历史数据（按时间范围）
 * @param serverId 服务器ID
 * @param range 时间范围 (1h, 6h, 24h, 7d)
 */
export function fetchNetIOHistory(
  serverId: number,
  range: TimeRange = '1h'
): Promise<NetworkIORecord[]> {
  return get(`/monitor/net/${serverId}?range=${range}`)
}

/**
 * 获取磁盘使用量历史数据（按时间范围）
 * @param serverId 服务器ID
 * @param range 时间范围 (1h, 6h, 24h, 7d)
 */
export function fetchDiskUsageHistory(
  serverId: number,
  range: TimeRange = '1h'
): Promise<DiskUsageRecord[]> {
  return get(`/monitor/disk-usage/${serverId}?range=${range}`)
}
