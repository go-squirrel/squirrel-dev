import type { ApplicationInstance, CreateApplicationRequest, UpdateApplicationRequest } from '@/types'

export type { ApplicationInstance, CreateApplicationRequest, UpdateApplicationRequest }

export interface FormState {
  visible: boolean
  mode: 'create' | 'edit'
  data: CreateApplicationRequest | UpdateApplicationRequest | null
}

export interface DeleteState {
  visible: boolean
  application: ApplicationInstance | null
}

export interface DetailState {
  visible: boolean
  application: ApplicationInstance | null
}
