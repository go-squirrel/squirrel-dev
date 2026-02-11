import i18n from '@/lang'

// 错误码到模块映射
const ERROR_CODE_MAP: Record<number, string> = {
  // Server (60000-60059)
  60001: 'server',
  60002: 'server',
  60003: 'server',
  60004: 'server',
  60005: 'server',
  60006: 'server',
  60021: 'server',
  60022: 'server',
  60023: 'server',
  60024: 'server',
  60041: 'server',
  60042: 'server',
  60043: 'server',

  // Config (65000-65019)
  65001: 'config',
  65002: 'config',
  65003: 'config',
  65004: 'config',
  65005: 'config',
  65006: 'config',

  // Auth (66000-66019)
  66001: 'auth',
  66002: 'auth',
  66003: 'auth',
  66004: 'auth',
  66005: 'auth',

  // Application (71000-71019)
  71001: 'application',
  71002: 'application',
  71003: 'application',
  71004: 'application',
  71005: 'application',
  71006: 'application',
  71007: 'application',

  // Deployment (72000-72039)
  72001: 'deployment',
  72002: 'deployment',
  72003: 'deployment',
  72004: 'deployment',
  72005: 'deployment',
  72006: 'deployment',
  72007: 'deployment',
  72008: 'deployment',
  72009: 'deployment',
  72010: 'deployment',
  72021: 'deployment',
  72022: 'deployment',
  72023: 'deployment',
  72024: 'deployment',
  72025: 'deployment',
  72026: 'deployment',
  72027: 'deployment',

  // AppStore (73000-73019)
  73001: 'appStore',
  73002: 'appStore',
  73003: 'appStore',
  73004: 'appStore',
  73005: 'appStore',
  73006: 'appStore',
  73007: 'appStore',

  // Script (80000-80039)
  80001: 'script',
  80002: 'script',
  80003: 'script',
  80004: 'script',
  80005: 'script',
  80021: 'script',
  80022: 'script',
  80023: 'script',

  // Monitor (81000-81019)
  81001: 'monitor',
  81002: 'monitor',
  81003: 'monitor',
  81004: 'monitor',
}

// API 错误类型
export interface ApiError {
  code: number
  message?: string
  data?: unknown
}

/**
 * 获取错误消息
 */
export function getErrorMessage(error: ApiError | number): string {
  const { t } = i18n.global

  let code: number
  let message: string | undefined

  if (typeof error === 'number') {
    code = error
  } else {
    code = error.code
    message = error.message
  }

  const module = ERROR_CODE_MAP[code]

  if (module) {
    const errorKey = `error.${module}.${code}`
    const translated = t(errorKey)
    // 如果翻译结果与键相同，说明没有找到翻译，返回原始消息或通用错误
    if (translated !== errorKey) {
      return translated
    }
  }

  // 如果没有对应的模块或翻译，返回原始消息或通用错误
  return message || t('error.common.unknownError')
}

/**
 * 判断是否为网络错误
 */
export function isNetworkError(error: unknown): boolean {
  if (error instanceof TypeError) {
    return error.message.includes('fetch') || error.message.includes('network')
  }
  return false
}

/**
 * 判断是否为认证错误
 */
export function isAuthError(error: ApiError): boolean {
  return error.code >= 66000 && error.code <= 66019
}

/**
 * 判断是否为权限错误
 */
export function isPermissionError(error: ApiError): boolean {
  return error.code >= 66000 && error.code <= 66019
}
