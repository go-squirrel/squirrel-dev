import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { getErrorMessage, isNetworkError, isAuthError, type ApiError } from '@/utils/errorHandler'

export function useErrorHandler() {
  const { t } = useI18n()
  const error = ref<string | null>(null)
  const errorCode = ref<number | null>(null)

  const handleError = (err: ApiError | unknown) => {
    if (isNetworkError(err)) {
      error.value = t('error.common.networkError')
      return error.value
    }

    if (err && typeof err === 'object' && 'code' in err) {
      const apiError = err as ApiError
      errorCode.value = apiError.code
      error.value = getErrorMessage(apiError)

      // 认证错误跳转
      if (isAuthError(apiError)) {
        localStorage.removeItem('token')
        window.location.href = '/login'
      }
    } else if (err instanceof Error) {
      error.value = err.message
    } else {
      error.value = t('error.common.unknownError')
    }

    return error.value
  }

  const clearError = () => {
    error.value = null
    errorCode.value = null
  }

  return {
    error,
    errorCode,
    handleError,
    clearError,
  }
}
