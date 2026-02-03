// 应用相关 API
import { get } from '@/utils/request'
import type { Application } from '@/types'

/**
 * 获取应用列表
 */
export function fetchApplications(): Promise<Application[]> {
  return get('/application')
}
