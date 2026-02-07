// 脚本相关 API
import { get, post, del } from '@/utils/request'
import type { Script, CreateScriptRequest, UpdateScriptRequest, ExecuteScriptRequest, ScriptResult } from '@/types'

/**
 * 获取脚本列表
 */
export function fetchScripts(): Promise<Script[]> {
  return get('/scripts')
}

/**
 * 获取脚本详情
 */
export function fetchScriptDetail(scriptId: number): Promise<Script> {
  return get(`/scripts/${scriptId}`)
}

/**
 * 创建脚本
 */
export function createScript(data: CreateScriptRequest): Promise<string> {
  return post('/scripts', data)
}

/**
 * 更新脚本
 */
export function updateScript(scriptId: number, data: UpdateScriptRequest): Promise<string> {
  return post(`/scripts/${scriptId}`, data)
}

/**
 * 删除脚本
 */
export function deleteScript(scriptId: number): Promise<string> {
  return del(`/scripts/${scriptId}`)
}

/**
 * 执行脚本
 */
export function executeScript(data: ExecuteScriptRequest): Promise<string> {
  return post('/scripts/execute', data)
}

/**
 * 获取脚本执行结果
 */
export function fetchScriptResults(scriptId: number): Promise<ScriptResult[]> {
  return get(`/scripts/${scriptId}/results`)
}
