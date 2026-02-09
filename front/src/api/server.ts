// 服务器相关 API
import { get, post, del } from '@/utils/request'
import type { Server, CreateServerRequest, UpdateServerRequest } from '@/types'

/**
 * 获取服务器列表
 */
export function fetchServers(): Promise<Server[]> {
  return get('/server')
}

/**
 * 获取服务器详情
 */
export function fetchServerDetail(serverId: number): Promise<Server> {
  return get(`/server/${serverId}`)
}

/**
 * 创建服务器
 */
export function createServer(data: CreateServerRequest): Promise<string> {
  return post('/server', data)
}

/**
 * 更新服务器
 */
export function updateServer(serverId: number, data: UpdateServerRequest): Promise<string> {
  return post(`/server/${serverId}`, data)
}

/**
 * 删除服务器
 */
export function deleteServer(serverId: number): Promise<string> {
  return del(`/server/${serverId}`)
}

/**
 * 获取终端 WebSocket URL
 */
export function getTerminalWebSocketUrl(serverId: number): string {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const host = window.location.host
  const token = localStorage.getItem('token') || ''
  return `${protocol}//${host}/api/v1/ws/server/${serverId}?token=${encodeURIComponent(token)}`
}
