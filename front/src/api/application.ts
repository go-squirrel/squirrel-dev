// 应用相关 API
import { get, post, del } from '@/utils/request'
import type { ApplicationInstance, CreateApplicationRequest, UpdateApplicationRequest } from '@/types'

/**
 * 获取应用列表
 */
export function fetchApplications(): Promise<ApplicationInstance[]> {
  return get('/application')
}

/**
 * 获取应用详情
 */
export function fetchApplicationDetail(applicationId: number): Promise<ApplicationInstance> {
  return get(`/application/${applicationId}`)
}

/**
 * 创建应用
 */
export function createApplication(data: CreateApplicationRequest): Promise<string> {
  return post('/application', data)
}

/**
 * 更新应用
 */
export function updateApplication(applicationId: number, data: UpdateApplicationRequest): Promise<string> {
  return post(`/application/${applicationId}`, data)
}

/**
 * 删除应用
 */
export function deleteApplication(applicationId: number): Promise<string> {
  return del(`/application/${applicationId}`)
}
