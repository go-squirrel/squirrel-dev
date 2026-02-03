// 系统信息相关逻辑
import { ref, onMounted } from 'vue'
import { fetchServerDetail } from '@/api'
import type { SystemInfo } from '@/types'

export function useSystemInfo(serverId: Ref<number>) {
  const systemInfo = ref<SystemInfo>({
    hostname: '-',
    os: '-',
    kernel: '-',
    arch: '-',
    ip: '-',
    bootTime: '-',
    uptime: '-'
  })

  const loadSystemInfo = async () => {
    if (!serverId.value) return
    try {
      const data = await fetchServerDetail(serverId.value)
      if (data.server_info) {
        const info = data.server_info
        systemInfo.value.hostname = info.hostname || '-'
        systemInfo.value.os = `${info.platform || '-'} ${info.platformVersion || ''}`.trim() || '-'
        systemInfo.value.kernel = info.kernelVersion || '-'
        systemInfo.value.arch = info.architecture || '-'
        systemInfo.value.uptime = info.uptimeStr || '-'

        if (info.ipAddresses && info.ipAddresses.length > 0) {
          const firstInterface = info.ipAddresses[0]
          if (firstInterface.ipv4 && firstInterface.ipv4.length > 0) {
            systemInfo.value.ip = firstInterface.ipv4[0]
          } else {
            systemInfo.value.ip = '-'
          }
        } else {
          systemInfo.value.ip = '-'
        }
      }
    } catch (error) {
      console.error('Failed to fetch server detail:', error)
    }
  }

  onMounted(() => {
    if (serverId.value) {
      loadSystemInfo()
    }
  })

  return {
    systemInfo,
    loadSystemInfo
  }
}

import type { Ref } from 'vue'
