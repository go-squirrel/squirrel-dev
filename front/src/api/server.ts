// 服务器相关 API
import { get } from '@/utils/request'
import type { Server } from '@/types'

/**
 * 获取服务器列表
 */
export function fetchServers(): Promise<Server[]> {
  return get('/server')
}

/**
 * 获取服务器详情
 */
export function fetchServerDetail(serverId: number): Promise<any> {
  return get(`/server/${serverId}`)
}
