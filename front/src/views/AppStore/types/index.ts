import type { AppStore, CreateAppRequest, UpdateAppRequest } from '@/types'

export type { AppStore, CreateAppRequest, UpdateAppRequest }

// 表单状态
export interface FormState {
  visible: boolean
  mode: 'create' | 'edit'
  data: CreateAppRequest | UpdateAppRequest | null
}

// 详情状态
export interface DetailState {
  visible: boolean
  app: AppStore | null
}

// 删除确认状态
export interface DeleteState {
  visible: boolean
  app: AppStore | null
}

// 筛选状态
export interface FilterState {
  category?: string
  type?: string
  keyword?: string
}
