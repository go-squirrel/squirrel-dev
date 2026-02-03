// HTTP 请求工具

const API_BASE = '/api/v1'

/**
 * 发送 GET 请求
 */
export async function get<T>(url: string): Promise<T> {
  const response = await fetch(`${API_BASE}${url}`, {
    headers: {
      'Authorization': `Bearer ${localStorage.getItem('token') || ''}`
    }
  })
  
  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`)
  }
  
  const result = await response.json()
  if (result.code !== 0) {
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
  
  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`)
  }
  
  const result = await response.json()
  if (result.code !== 0) {
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
  
  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`)
  }
  
  const result = await response.json()
  if (result.code !== 0) {
    throw new Error(result.message || '请求失败')
  }
  
  return result.data
}
