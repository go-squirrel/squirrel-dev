// 应用列表相关逻辑
import { ref, onMounted } from 'vue'
import { fetchApplications } from '@/api'
import type { Application } from '@/types'

export function useApplications() {
  const appList = ref<Application[]>([])
  const loading = ref(false)

  const loadApplications = async () => {
    loading.value = true
    try {
      const data = await fetchApplications()
      appList.value = data.map((app: any) => ({
        id: app.id,
        name: app.name,
        version: app.version || '1.0.0',
        status: app.status || 'stopped',
        icon: 'lucide:box',
        color: '#4fc3f7'
      }))
    } catch (error) {
      console.error('Failed to fetch applications:', error)
    } finally {
      loading.value = false
    }
  }

  const getAppStatusText = (status: string): string => {
    const map: Record<string, string> = {
      running: '运行中',
      stopped: '已停止',
      error: '错误'
    }
    return map[status] || status
  }

  onMounted(() => {
    loadApplications()
  })

  return {
    appList,
    loading,
    loadApplications,
    getAppStatusText
  }
}
