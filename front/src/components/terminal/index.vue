<template>
  <div class="terminal-container">
    <div class="terminal-header">
      <div class="terminal-info">
        <Icon icon="lucide:terminal" class="terminal-icon" />
        <span class="terminal-title">{{ $t('server.terminal') }}</span>
        <span class="server-name">- {{ server.hostname }}</span>
      </div>
      <div class="terminal-actions">
        <button v-if="!connected" class="action-btn reconnect-btn" @click="connect">
          <Icon icon="lucide:refresh-cw" />
          {{ $t('server.reconnect') }}
        </button>
        <button class="action-btn close-btn" @click="handleClose">
          <Icon icon="lucide:x" />
          {{ $t('server.close') }}
        </button>
      </div>
    </div>
    <div class="terminal-body">
      <div v-if="connecting" class="connecting-state">
        <Icon icon="lucide:loader-2" class="spinner" />
        <p>{{ $t('server.connecting') }}</p>
      </div>
      <div v-else-if="connectionError" class="error-state">
        <Icon icon="lucide:alert-circle" class="error-icon" />
        <p>{{ $t('server.connectionFailed') }}</p>
        <button class="retry-btn" @click="connect">
          {{ $t('server.reconnect') }}
        </button>
      </div>
      <div v-else ref="terminalRef" class="xterm-container"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
import { Icon } from '@iconify/vue'
import { getTerminalWebSocketUrl } from '@/api/server'
import type { Server } from '@/types'

const props = defineProps<{
  server: Server
}>()

const router = useRouter()

const terminalRef = ref<HTMLElement>()
const connecting = ref(true)
const connected = ref(false)
const connectionError = ref(false)
let ws: WebSocket | null = null

const connect = () => {
  connecting.value = true
  connectionError.value = false
  
  if (ws) {
    ws.close()
  }

  try {
    const url = getTerminalWebSocketUrl(props.server.id)
    ws = new WebSocket(url)

    ws.onopen = () => {
      connecting.value = false
      connected.value = true
      connectionError.value = false
    }

    ws.onmessage = (event) => {
      if (terminalRef.value) {
        terminalRef.value.textContent += event.data
      }
    }

    ws.onerror = () => {
      connecting.value = false
      connectionError.value = true
      connected.value = false
    }

    ws.onclose = () => {
      connecting.value = false
      connected.value = false
    }
  } catch (error) {
    console.error('Failed to connect:', error)
    connecting.value = false
    connectionError.value = true
  }
}

const handleClose = () => {
  if (ws) {
    ws.close()
  }
  router.push('/servers')
}

onMounted(() => {
  connect()
})

onBeforeUnmount(() => {
  if (ws) {
    ws.close()
  }
})
</script>

<style scoped>
.terminal-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: #1e1e1e;
}

.terminal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  background: #2d2d2d;
  border-bottom: 1px solid #404040;
}

.terminal-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.terminal-icon {
  width: 18px;
  height: 18px;
  color: #4fc3f7;
}

.terminal-title {
  font-size: 14px;
  font-weight: 600;
  color: #ffffff;
}

.server-name {
  font-size: 13px;
  color: #a0a0a0;
}

.terminal-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  border: none;
}

.reconnect-btn {
  background: #4fc3f7;
  color: #1e1e1e;
}

.reconnect-btn:hover {
  background: #29b6f6;
}

.close-btn {
  background: #dc2626;
  color: #ffffff;
}

.close-btn:hover {
  background: #b91c1c;
}

.terminal-body {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}

.connecting-state,
.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  color: #a0a0a0;
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

.error-icon {
  width: 48px;
  height: 48px;
  color: #dc2626;
}

.retry-btn {
  padding: 8px 16px;
  background: #4fc3f7;
  color: #1e1e1e;
  border: none;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.retry-btn:hover {
  background: #29b6f6;
}

.xterm-container {
  width: 100%;
  height: 100%;
  padding: 16px;
  font-family: 'SF Mono', Monaco, Consolas, monospace;
  font-size: 13px;
  line-height: 1.6;
  color: #f0f0f0;
  white-space: pre-wrap;
  overflow: auto;
  background: #1e1e1e;
}
</style>
