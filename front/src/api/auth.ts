// 认证相关 API
import { post } from '@/utils/request'

export interface LoginParams {
  username: string
  password: string
}

export interface LoginResult {
  token: string
}

/**
 * 用户登录
 */
export function login(params: LoginParams): Promise<LoginResult> {
  return post('/login', params)
}
