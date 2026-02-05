<template>
  <div class="server-terminal-page">
    <Terminal v-if="server" :server="server" />
    <div v-else class="loading">
      {{ $t('server.loading') }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import Terminal from '@/components/terminal/index.vue'
import { fetchServerDetail } from '@/api/server'
import type { Server } from '@/types'

const route = useRoute()

const server = ref<Server | null>(null)

onMounted(async () => {
  const serverId = Number(route.params.id)
  if (serverId) {
    try {
      server.value = await fetchServerDetail(serverId)
    } catch (error) {
      console.error('Failed to load server:', error)
    }
  }
})
</script>

<style scoped>
.server-terminal-page {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.loading {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #64748b;
  font-size: 14px;
}
</style>
