import type { ApplicationInstance, Server } from '@/types'

// 部署状态
export type DeploymentStatus = 'running' | 'stopped' | 'not_deployed' | 'error'

// 部署信息
export interface Deployment {
  id: number
  deploy_id: number
  application: ApplicationInstance
  server: Server
  status: DeploymentStatus
  deployed_at: string
}

// 创建部署请求
export interface CreateDeploymentRequest {
  application_id: number
  server_id: number
}

// 表单状态
export interface FormState {
  visible: boolean
  data: CreateDeploymentRequest | null
}

// 卸载确认状态
export interface UndeployState {
  visible: boolean
  deployment: Deployment | null
}

// 详情状态
export interface DetailState {
  visible: boolean
  deployment: Deployment | null
}
