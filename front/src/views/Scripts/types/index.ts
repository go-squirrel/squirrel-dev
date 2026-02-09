import type { Script, CreateScriptRequest, UpdateScriptRequest, ExecuteScriptRequest, ScriptResult } from '@/types'

export type { Script, CreateScriptRequest, UpdateScriptRequest, ExecuteScriptRequest, ScriptResult }

// 表单状态
export interface FormState {
  visible: boolean
  mode: 'create' | 'edit'
  data: CreateScriptRequest | UpdateScriptRequest | null
}

// 执行对话框状态
export interface ExecuteState {
  visible: boolean
  script: Script | null
}

// 结果列表状态
export interface ResultState {
  visible: boolean
  script: Script | null
}
