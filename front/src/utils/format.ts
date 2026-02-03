// 格式化工具函数

/**
 * 格式化字节大小
 */
export function formatBytes(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

/**
 * 格式化速度
 */
export function formatSpeed(bytes: number): string {
  return formatBytes(bytes) + '/s'
}

/**
 * 格式化日期时间
 */
export function formatDateTime(date: Date | string | number): string {
  const d = new Date(date)
  return d.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

/**
 * 格式化运行时间
 */
export function formatUptime(seconds: number): string {
  const days = Math.floor(seconds / 86400)
  const hours = Math.floor((seconds % 86400) / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)
  
  if (days > 0) {
    return `${days}天 ${hours}小时 ${minutes}分钟`
  } else if (hours > 0) {
    return `${hours}小时 ${minutes}分钟`
  } else {
    return `${minutes}分钟`
  }
}

/**
 * 获取状态样式类名
 */
export function getStatusClass(status: string): string {
  switch (status) {
    case 'online':
    case 'running':
      return 'online'
    case 'offline':
    case 'stopped':
      return 'offline'
    case 'error':
      return 'error'
    default:
      return 'unknown'
  }
}

/**
 * 获取状态显示文本
 */
export function getStatusText(status: string): string {
  const map: Record<string, string> = {
    online: '在线',
    offline: '离线',
    running: '运行中',
    stopped: '已停止',
    error: '错误',
    unknown: '未知'
  }
  return map[status] || status
}
