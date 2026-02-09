// 部署相关 API
import { get, post, del } from '@/utils/request'
import type { Deployment, CreateDeploymentRequest, Server } from '@/types'

/**
 * 获取部署列表
 */
export function fetchDeployments(serverId?: number): Promise<Deployment[]> {
  const url = serverId ? `/deployment?server_id=${serverId}` : '/deployment'
  return get(url)
}

/**
 * 创建部署
 */
export function createDeployment(applicationId: number, data: CreateDeploymentRequest): Promise<string> {
  return post(`/deployment/${applicationId}`, data)
}

/**
 * 启动应用
 */
export function startDeployment(deploymentId: number): Promise<string> {
  return post(`/deployment/start/${deploymentId}`)
}

/**
 * 停止应用
 */
export function stopDeployment(deploymentId: number): Promise<string> {
  return post(`/deployment/stop/${deploymentId}`)
}

/**
 * 卸载部署
 */
export function undeployDeployment(deploymentId: number): Promise<string> {
  return del(`/deployment/${deploymentId}`)
}

/**
 * 获取应用部署的服务器列表
 */
export function fetchDeploymentServers(applicationId: number): Promise<Server[]> {
  return get(`/deployment/${applicationId}/servers`)
}
