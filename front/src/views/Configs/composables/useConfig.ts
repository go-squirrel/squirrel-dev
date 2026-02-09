import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { fetchConfigs, fetchConfigDetail, createConfig, updateConfig, deleteConfig } from '@/api/config'
import type { Config, CreateConfigRequest, UpdateConfigRequest } from '@/types'

export function useConfig() {
  const { t } = useI18n()
  const configs = ref<Config[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  const loadConfigs = async () => {
    loading.value = true
    error.value = null
    try {
      configs.value = await fetchConfigs()
    } catch (err) {
      error.value = t('configs.operationFailed')
      console.error('Failed to load configs:', err)
    } finally {
      loading.value = false
    }
  }

  const loadConfigDetail = async (id: number): Promise<Config | null> => {
    loading.value = true
    error.value = null
    try {
      return await fetchConfigDetail(id)
    } catch (err) {
      error.value = t('configs.operationFailed')
      console.error('Failed to load config detail:', err)
      return null
    } finally {
      loading.value = false
    }
  }

  const addConfig = async (data: CreateConfigRequest): Promise<boolean> => {
    loading.value = true
    error.value = null
    try {
      await createConfig(data)
      await loadConfigs()
      return true
    } catch (err) {
      error.value = t('configs.operationFailed')
      console.error('Failed to create config:', err)
      return false
    } finally {
      loading.value = false
    }
  }

  const editConfig = async (id: number, data: UpdateConfigRequest): Promise<boolean> => {
    loading.value = true
    error.value = null
    try {
      await updateConfig(id, data)
      await loadConfigs()
      return true
    } catch (err) {
      error.value = t('configs.operationFailed')
      console.error('Failed to update config:', err)
      return false
    } finally {
      loading.value = false
    }
  }

  const removeConfig = async (id: number): Promise<boolean> => {
    loading.value = true
    error.value = null
    try {
      await deleteConfig(id)
      await loadConfigs()
      return true
    } catch (err) {
      error.value = t('configs.operationFailed')
      console.error('Failed to delete config:', err)
      return false
    } finally {
      loading.value = false
    }
  }

  return {
    configs,
    loading,
    error,
    loadConfigs,
    loadConfigDetail,
    addConfig,
    editConfig,
    removeConfig
  }
}
