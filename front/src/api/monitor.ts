// 监控相关 API
import { get } from '@/utils/request'
import type { MonitorData } from '@/types'
import type { PageData, BaseMonitorRecord, DiskIORecord, NetworkIORecord } from '@/types/monitor'

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
 * 获取基础监控历史数据
 * @param serverId 服务器ID
 * @param page 页码
 * @param count 每页数量
 */
export function fetchBaseMonitorHistory(
  serverId: number,
  page: number = 1,
  count: number = 100
): Promise<PageData<BaseMonitorRecord>> {
  return get(`/monitor/base/${serverId}/${page}/${count}`)
}

/**
 * 获取磁盘IO历史数据
 * @param serverId 服务器ID
 * @param page 页码
 * @param count 每页数量
 */
export function fetchDiskIOHistory(
  serverId: number,
  page: number = 1,
  count: number = 100
): Promise<PageData<DiskIORecord>> {
  return get(`/monitor/disk/${serverId}/${page}/${count}`)
}

/**
 * 获取网络IO历史数据
 * @param serverId 服务器ID
 * @param page 页码
 * @param count 每页数量
 */
export function fetchNetIOHistory(
  serverId: number,
  page: number = 1,
  count: number = 100
): Promise<PageData<NetworkIORecord>> {
  return get(`/monitor/net/${serverId}/${page}/${count}`)
}
