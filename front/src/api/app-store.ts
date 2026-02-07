import { get, post, del } from '@/utils/request'
import type { AppStore, CreateAppRequest, UpdateAppRequest } from '@/types'

/**
 * 获取应用列表
 */
export function fetchApps(): Promise<AppStore[]> {
  return get('/app-store')
}

/**
 * 获取应用详情
 */
export function fetchAppDetail(appId: number): Promise<AppStore> {
  return get(`/app-store/${appId}`)
}

/**
 * 创建应用
 */
export function createApp(data: CreateAppRequest): Promise<string> {
  return post('/app-store', data)
}

/**
 * 更新应用
 */
export function updateApp(appId: number, data: UpdateAppRequest): Promise<string> {
  return post(`/app-store/${appId}`, data)
}

/**
 * 删除应用
 */
export function deleteApp(appId: number): Promise<string> {
  return del(`/app-store/${appId}`)
}
