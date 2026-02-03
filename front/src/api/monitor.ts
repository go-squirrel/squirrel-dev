// 监控相关 API
import { get } from '@/utils/request'
import type { MonitorData } from '@/types'

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
 * 获取磁盘 IO 统计
 */
export function fetchIOStats(serverId: number, target: string): Promise<any> {
  return get(`/monitor/stats/io/${serverId}/${target}`)
}
