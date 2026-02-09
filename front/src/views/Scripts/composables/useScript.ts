import { ref } from 'vue'
import { fetchScripts, fetchScriptDetail, createScript, updateScript, deleteScript, executeScript, fetchScriptResults } from '@/api/script'
import type { Script, CreateScriptRequest, UpdateScriptRequest, ExecuteScriptRequest, ScriptResult } from '@/types'

export function useScript() {
  const scripts = ref<Script[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  const loadScripts = async () => {
    loading.value = true
    error.value = null
    try {
      scripts.value = await fetchScripts()
    } catch (err) {
      error.value = err instanceof Error ? err.message : '加载失败'
    } finally {
      loading.value = false
    }
  }

  const getScriptDetail = async (scriptId: number): Promise<Script | null> => {
    try {
      return await fetchScriptDetail(scriptId)
    } catch (err) {
      error.value = err instanceof Error ? err.message : '获取详情失败'
      return null
    }
  }

  const addScript = async (data: CreateScriptRequest): Promise<boolean> => {
    try {
      await createScript(data)
      return true
    } catch (err) {
      error.value = err instanceof Error ? err.message : '创建失败'
      return false
    }
  }

  const editScript = async (scriptId: number, data: UpdateScriptRequest): Promise<boolean> => {
    try {
      await updateScript(scriptId, data)
      return true
    } catch (err) {
      error.value = err instanceof Error ? err.message : '更新失败'
      return false
    }
  }

  const removeScript = async (scriptId: number): Promise<boolean> => {
    try {
      await deleteScript(scriptId)
      return true
    } catch (err) {
      error.value = err instanceof Error ? err.message : '删除失败'
      return false
    }
  }

  const runScript = async (data: ExecuteScriptRequest): Promise<boolean> => {
    try {
      await executeScript(data)
      return true
    } catch (err) {
      error.value = err instanceof Error ? err.message : '执行失败'
      return false
    }
  }

  const getScriptResults = async (scriptId: number): Promise<ScriptResult[]> => {
    try {
      return await fetchScriptResults(scriptId)
    } catch (err) {
      error.value = err instanceof Error ? err.message : '获取结果失败'
      return []
    }
  }

  return {
    scripts,
    loading,
    error,
    loadScripts,
    getScriptDetail,
    addScript,
    editScript,
    removeScript,
    runScript,
    getScriptResults
  }
}
