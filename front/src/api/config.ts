import { get, post, del } from '../utils/request'
import type { Config, CreateConfigRequest, UpdateConfigRequest } from '../types'

export function fetchConfigs(): Promise<Config[]> {
  return get('/config')
}

export function fetchConfigDetail(configId: number): Promise<Config> {
  return get(`/config/${configId}`)
}

export function createConfig(data: CreateConfigRequest): Promise<string> {
  return post('/config', data)
}

export function updateConfig(configId: number, data: UpdateConfigRequest): Promise<string> {
  return post(`/config/${configId}`, data)
}

export function deleteConfig(configId: number): Promise<string> {
  return del(`/config/${configId}`)
}
