// HTTP 请求工具
import { getErrorMessage, isAuthError as checkAuthError, type ApiError } from './errorHandler'

const API_BASE = '/api/v1'

/**
 * 处理认证失败：清除token并跳转到登录页
 */
function handleAuthError(): void {
  localStorage.removeItem('token')
  if (window.location.pathname !== '/login') {
    window.location.href = '/login'
  }
}

/**
 * 处理请求响应
 */
async function handleResponse<T>(response: Response): Promise<T> {
  // HTTP 401 未授权
  if (response.status === 401) {
    handleAuthError()
    throw new Error(getErrorMessage(66004))
  }

  if (!response.ok) {
    throw new Error(getErrorMessage({ code: response.status, message: `HTTP error! status: ${response.status}` }))
  }

  const result = await response.json()
  if (result.code !== 0) {
    const apiError: ApiError = { code: result.code, message: result.message }
    // 后端返回的认证错误码
    if (checkAuthError(apiError)) {
      handleAuthError()
    }
    throw new Error(getErrorMessage(apiError))
  }

  return result.data
}

/**
 * 发送 GET 请求
 */
export async function get<T>(url: string): Promise<T> {
  const response = await fetch(`${API_BASE}${url}`, {
    headers: {
      'Authorization': `Bearer ${localStorage.getItem('token') || ''}`
    }
  })

  return handleResponse<T>(response)
}

/**
 * 发送 POST 请求
 */
export async function post<T>(url: string, data?: any): Promise<T> {
  const response = await fetch(`${API_BASE}${url}`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${localStorage.getItem('token') || ''}`
    },
    body: data ? JSON.stringify(data) : undefined
  })

  return handleResponse<T>(response)
}

/**
 * 发送 DELETE 请求
 */
export async function del<T>(url: string): Promise<T> {
  const response = await fetch(`${API_BASE}${url}`, {
    method: 'DELETE',
    headers: {
      'Authorization': `Bearer ${localStorage.getItem('token') || ''}`
    }
  })

  return handleResponse<T>(response)
}

/**
 * 发送 PUT 请求
 */
export async function put<T>(url: string, data?: any): Promise<T> {
  const response = await fetch(`${API_BASE}${url}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${localStorage.getItem('token') || ''}`
    },
    body: data ? JSON.stringify(data) : undefined
  })

  return handleResponse<T>(response)
}

/**
 * 发送 PATCH 请求
 */
export async function patchRequest<T>(url: string, data?: any): Promise<T> {
  const response = await fetch(`${API_BASE}${url}`, {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${localStorage.getItem('token') || ''}`
    },
    body: data ? JSON.stringify(data) : undefined
  })

  return handleResponse<T>(response)
}
