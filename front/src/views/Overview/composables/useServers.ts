import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { fetchServers } from '@/api'
import type { Server } from '@/types'

export function useServers() {
  const route = useRoute()
  const router = useRouter()
  const serverList = ref<Server[]>([])
  const currentServerId = ref<number>(0)
  const loading = ref(false)

  const loadServers = async () => {
    loading.value = true
    try {
      const data = await fetchServers()
      serverList.value = data
      if (data.length > 0) {
        const queryServerId = Number(route.query.serverId)
        if (queryServerId && data.some(s => s.id === queryServerId)) {
          currentServerId.value = queryServerId
        } else if (currentServerId.value === 0) {
          currentServerId.value = data[0].id
        }
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
    router.replace({ query: { serverId: String(serverId) } })
  }

  return {
    serverList,
    currentServerId,
    loading,
    loadServers,
    switchServer
  }
}
