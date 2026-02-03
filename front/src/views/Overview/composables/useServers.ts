// 服务器相关逻辑
import { ref, onMounted } from 'vue'
import { fetchServers } from '@/api'
import type { Server } from '@/types'

export function useServers() {
  const serverList = ref<Server[]>([])
  const currentServerId = ref<number>(0)
  const loading = ref(false)

  const loadServers = async () => {
    loading.value = true
    try {
      const data = await fetchServers()
      serverList.value = data
      if (data.length > 0 && currentServerId.value === 0) {
        currentServerId.value = data[0].id
      }
      return data
    } catch (error) {
      console.error('Failed to fetch servers:', error)
      return []
    } finally {
      loading.value = false
    }
  }

  const switchServer = (serverId: number) => {
    currentServerId.value = serverId
  }

  onMounted(() => {
    loadServers()
  })

  return {
    serverList,
    currentServerId,
    loading,
    loadServers,
    switchServer
  }
}
