import { ref } from 'vue'

export function useLoading() {
  const loading = ref(false)

  const withLoading = async <T = any>(fn: () => Promise<T>): Promise<T | undefined> => {
    loading.value = true
    try {
      return await fn()
    } finally {
      loading.value = false
    }
  }

  const startLoading = () => {
    loading.value = true
  }

  const stopLoading = () => {
    loading.value = false
  }

  return {
    loading,
    withLoading,
    startLoading,
    stopLoading
  }
}
