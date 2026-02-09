// HTTP 请求工具

const API_BASE = '/api/v1'

// 认证错误码
const AUTH_ERROR_CODES = [41003, 41004, 41005] // JWT相关错误码

/**
 * 处理认证失败：清除token并跳转到登录页
 */
function handleAuthError(): void {
  localStorage.removeItem('token')
  window.location.href = '/login'
}

/**
 * 检查响应是否为认证错误
 */
function isAuthError(code: number, status: number): boolean {
  return status === 401 || AUTH_ERROR_CODES.includes(code)
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
  
  // HTTP 401 未授权
  if (response.status === 401) {
    handleAuthError()
    throw new Error('认证失败，请重新登录')
  }
  
  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`)
  }
  
  const result = await response.json()
  if (result.code !== 0) {
    // 后端返回的认证错误码
    if (isAuthError(result.code, 200)) {
      handleAuthError()
    }
    throw new Error(result.message || '请求失败')
  }
  
  return result.data
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
  
  // HTTP 401 未授权
  if (response.status === 401) {
    handleAuthError()
    throw new Error('认证失败，请重新登录')
  }
  
  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`)
  }
  
  const result = await response.json()
  if (result.code !== 0) {
    // 后端返回的认证错误码
    if (isAuthError(result.code, 200)) {
      handleAuthError()
    }
    throw new Error(result.message || '请求失败')
  }
  
  return result.data
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
  
  // HTTP 401 未授权
  if (response.status === 401) {
    handleAuthError()
    throw new Error('认证失败，请重新登录')
  }
  
  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`)
  }
  
  const result = await response.json()
  if (result.code !== 0) {
    // 后端返回的认证错误码
    if (isAuthError(result.code, 200)) {
      handleAuthError()
    }
    throw new Error(result.message || '请求失败')
  }
  
  return result.data
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
  
  // HTTP 401 未授权
  if (response.status === 401) {
    handleAuthError()
    throw new Error('认证失败，请重新登录')
  }
  
  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`)
  }
  
  const result = await response.json()
  if (result.code !== 0) {
    if (isAuthError(result.code, 200)) {
      handleAuthError()
    }
    throw new Error(result.message || '请求失败')
  }
  
  return result.data
}

/**
 * 发送 PATCH 请求
 */
export async function patch<T>(url: string, data?: any): Promise<T> {
  const response = await fetch(`${API_BASE}${url}`, {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${localStorage.getItem('token') || ''}`
    },
    body: data ? JSON.stringify(data) : undefined
  })
  
  // HTTP 401 未授权
  if (response.status === 401) {
    handleAuthError()
    throw new Error('认证失败，请重新登录')
  }
  
  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`)
  }
  
  const result = await response.json()
  if (result.code !== 0) {
    if (isAuthError(result.code, 200)) {
      handleAuthError()
    }
    throw new Error(result.message || '请求失败')
  }
  
  return result.data
}
