// 监控数据缓存管理
import type { ChartDataPoint, MonitorData } from '../types'

interface ServerCache {
  monitorData: MonitorData
  chartData: ChartDataPoint[]
  currentNetStats: { bytesSent: number; bytesRecv: number }
  currentIOStats: { readBytes: number; writeBytes: number }
  lastNetStats: { bytesSent: number; bytesRecv: number; timestamp: number }
  lastIOStats: { readBytes: number; writeBytes: number; timestamp: number }
  chartTargetList: string[]
  timestamp: number
}

// 缓存存储 - 使用 Map 存储每个服务器的数据
const cacheStore = new Map<number, ServerCache>()

// 缓存有效期（毫秒）- 5分钟
const CACHE_EXPIRY = 5 * 60 * 1000

export function useMonitorCache() {
  // 获取缓存数据
  const getCache = (serverId: number): ServerCache | null => {
    const cache = cacheStore.get(serverId)
    if (!cache) return null

    // 检查缓存是否过期
    if (Date.now() - cache.timestamp > CACHE_EXPIRY) {
      cacheStore.delete(serverId)
      return null
    }

    return cache
  }

  // 设置缓存数据
  const setCache = (serverId: number, data: Partial<ServerCache>) => {
    const existing = cacheStore.get(serverId)
    const newCache: ServerCache = {
      monitorData: data.monitorData ?? existing?.monitorData ?? {},
      chartData: data.chartData ?? existing?.chartData ?? [],
      currentNetStats: data.currentNetStats ?? existing?.currentNetStats ?? { bytesSent: 0, bytesRecv: 0 },
      currentIOStats: data.currentIOStats ?? existing?.currentIOStats ?? { readBytes: 0, writeBytes: 0 },
      lastNetStats: data.lastNetStats ?? existing?.lastNetStats ?? { bytesSent: 0, bytesRecv: 0, timestamp: 0 },
      lastIOStats: data.lastIOStats ?? existing?.lastIOStats ?? { readBytes: 0, writeBytes: 0, timestamp: 0 },
      chartTargetList: data.chartTargetList ?? existing?.chartTargetList ?? [],
      timestamp: Date.now()
    }
    cacheStore.set(serverId, newCache)
  }

  // 清除指定服务器的缓存
  const clearCache = (serverId: number) => {
    cacheStore.delete(serverId)
  }

  // 清除所有缓存
  const clearAllCache = () => {
    cacheStore.clear()
  }

  // 检查是否有有效缓存
  const hasValidCache = (serverId: number): boolean => {
    return getCache(serverId) !== null
  }

  return {
    getCache,
    setCache,
    clearCache,
    clearAllCache,
    hasValidCache
  }
}
