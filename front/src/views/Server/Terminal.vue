<template>
  <div class="server-terminal-page">
    <Terminal v-if="server" :server="server" />
    <div v-else class="loading">
      <Icon icon="lucide:loader-2" class="spinner" />
      <span>{{ $t('server.loading') }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { Icon } from '@iconify/vue'
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
  width: 100%;
  height: 100vh;
  overflow: hidden;
}

.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #64748b;
  font-size: 14px;
  gap: 12px;
}

.spinner {
  width: 32px;
  height: 32px;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}
</style>
