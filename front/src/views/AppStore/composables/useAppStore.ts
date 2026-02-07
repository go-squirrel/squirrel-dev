import { ref, computed, type Ref } from 'vue'
import { fetchApps, createApp, updateApp, deleteApp } from '@/api/app-store'
import type { AppStore, CreateAppRequest, UpdateAppRequest } from '@/types'

export function useAppStore(
  searchKeyword: Ref<string>,
  selectedCategory: Ref<string>,
  selectedType: Ref<string>
) {
  const apps = ref<AppStore[]>([])
  const loading = ref(false)

  // 过滤后的应用列表
  const filteredApps = computed(() => {
    let result = apps.value

    if (searchKeyword.value) {
      const keyword = searchKeyword.value.toLowerCase()
      result = result.filter(app =>
        app.name.toLowerCase().includes(keyword) ||
        app.description.toLowerCase().includes(keyword) ||
        app.tags.toLowerCase().includes(keyword)
      )
    }

    if (selectedCategory.value) {
      result = result.filter(app => app.category === selectedCategory.value)
    }

    if (selectedType.value) {
      result = result.filter(app => app.type === selectedType.value)
    }

    return result
  })

  // 加载应用列表
  const loadApps = async () => {
    loading.value = true
    try {
      apps.value = await fetchApps()
    } catch (error) {
      console.error('Failed to load apps:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  // 添加应用
  const addApp = async (data: CreateAppRequest) => {
    await createApp(data)
  }

  // 编辑应用
  const editApp = async (appId: number, data: UpdateAppRequest) => {
    await updateApp(appId, data)
  }

  // 移除应用
  const removeApp = async (appId: number) => {
    await deleteApp(appId)
  }

  return {
    apps,
    loading,
    filteredApps,
    loadApps,
    addApp,
    editApp,
    removeApp
  }
}
