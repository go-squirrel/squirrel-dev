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
  return `${protocol}//${host}/api/v1/ws/server/${serverId}`
}

/**
 * 测试 SSH 连接
 */
export function testSSHConnection(serverId: number): Promise<{ message: string; hostname: string; ip_address: string; ssh_port: number }> {
  return post(`/ssh/test/${serverId}`)
}

/**
 * 获取当前用户的 token
 */
export function getAuthToken(): string {
  return localStorage.getItem('token') || ''
}
