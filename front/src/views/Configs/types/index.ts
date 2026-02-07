import type { Config, CreateConfigRequest, UpdateConfigRequest } from '../../types'

export type { Config, CreateConfigRequest, UpdateConfigRequest }

export interface FormState {
  visible: boolean
  mode: 'create' | 'edit'
  data: CreateConfigRequest | UpdateConfigRequest | null
}

export interface DeleteState {
  visible: boolean
  config: Config | null
}
