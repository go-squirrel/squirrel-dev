<template>
  <div class="terminal-fullscreen">
    <div class="terminal-header">
      <div class="terminal-info">
        <Icon icon="lucide:terminal" class="terminal-icon" />
        <span class="terminal-title">{{ $t('server.terminal') }}</span>
        <span class="server-name" v-if="server">- {{ server.hostname || server.ip_address }}</span>
      </div>
      <div class="terminal-actions">
        <button v-if="!connected && !connecting" class="action-btn reconnect-btn" @click="connect">
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
      <div ref="terminalRef" class="xterm-container"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
import { Icon } from '@iconify/vue'
import { Terminal } from '@xterm/xterm'
import { FitAddon } from '@xterm/addon-fit'
import '@xterm/xterm/css/xterm.css'
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
let term: Terminal | null = null
let fitAddon: FitAddon | null = null
let resizeObserver: ResizeObserver | null = null

const connect = () => {
  connecting.value = true
  connectionError.value = false
  connected.value = false

  if (ws) {
    ws.close()
    ws = null
  }

  // 清除之前的终端
  if (terminalRef.value) {
    terminalRef.value.innerHTML = ''
  }

  try {
    const url = getTerminalWebSocketUrl(props.server.id)
    ws = new WebSocket(url)

    ws.onopen = () => {
      connecting.value = false
      connected.value = true
      connectionError.value = false

      // 初始化 xterm
      initTerminal()
    }

    ws.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data)
        if (message.type === 'stdout' && message.data && term) {
          term.write(message.data)
        } else if (message.type === 'stderr' && message.data && term) {
          term.write(message.data)
        }
      } catch {
        // 如果不是 JSON 格式，直接显示原始数据
        if (term) {
          term.write(event.data)
        }
      }
    }

    ws.onerror = () => {
      connecting.value = false
      connectionError.value = true
      connected.value = false
      if (term) {
        term.writeln('\r\n\x1b[31m连接发生错误\x1b[0m')
      }
    }

    ws.onclose = () => {
      connecting.value = false
      connected.value = false
      if (term && !connectionError.value) {
        term.writeln('\r\n\x1b[33m连接已关闭\x1b[0m')
      }
    }
  } catch (error) {
    console.error('Failed to connect:', error)
    connecting.value = false
    connectionError.value = true
  }
}

const initTerminal = () => {
  if (!terminalRef.value || !ws) return

  // 创建终端实例
  term = new Terminal({
    fontSize: 14,
    fontFamily: 'SF Mono, Monaco, Consolas, "Liberation Mono", "Courier New", monospace',
    cursorBlink: true,
    cursorStyle: 'block',
    theme: {
      background: '#1e1e1e',
      foreground: '#f0f0f0',
      cursor: '#f0f0f0',
      selectionBackground: '#264f78',
      black: '#000000',
      red: '#cd3131',
      green: '#0dbc79',
      yellow: '#e5e510',
      blue: '#2472c8',
      magenta: '#bc3fbc',
      cyan: '#11a8cd',
      white: '#e5e5e5',
      brightBlack: '#666666',
      brightRed: '#f14c4c',
      brightGreen: '#23d18b',
      brightYellow: '#f5f543',
      brightBlue: '#3b8eea',
      brightMagenta: '#d670d6',
      brightCyan: '#29b8db',
      brightWhite: '#e5e5e5'
    },
    scrollback: 10000,
    allowProposedApi: true
  })

  // 添加自适应插件
  fitAddon = new FitAddon()
  term.loadAddon(fitAddon)

  // 打开终端
  term.open(terminalRef.value)

  // 自适应大小
  fitAddon.fit()

  // 发送初始大小到后端
  const dims = fitAddon.proposeDimensions()
  if (dims && ws && ws.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify({
      type: 'resize',
      cols: dims.cols,
      rows: dims.rows
    }))
  }

  // 处理输入
  term.onData((data) => {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({ type: 'stdin', data }))
    }
  })

  // 监听大小变化
  resizeObserver = new ResizeObserver(() => {
    if (fitAddon && term && ws && ws.readyState === WebSocket.OPEN) {
      fitAddon.fit()
      const dims = fitAddon.proposeDimensions()
      if (dims) {
        ws.send(JSON.stringify({
          type: 'resize',
          cols: dims.cols,
          rows: dims.rows
        }))
      }
    }
  })

  if (terminalRef.value) {
    resizeObserver.observe(terminalRef.value)
  }

  // 聚焦终端
  term.focus()
}

const handleClose = () => {
  if (resizeObserver && terminalRef.value) {
    resizeObserver.unobserve(terminalRef.value)
    resizeObserver.disconnect()
  }
  if (term) {
    term.dispose()
    term = null
  }
  if (ws) {
    ws.close()
    ws = null
  }
  router.push('/servers')
}

onMounted(() => {
  connect()
})

onBeforeUnmount(() => {
  if (resizeObserver && terminalRef.value) {
    resizeObserver.unobserve(terminalRef.value)
    resizeObserver.disconnect()
  }
  if (term) {
    term.dispose()
    term = null
  }
  if (ws) {
    ws.close()
    ws = null
  }
})
</script>

<style scoped>
.terminal-fullscreen {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  flex-direction: column;
  background: #1e1e1e;
  z-index: 1000;
}

.terminal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  background: #2d2d2d;
  border-bottom: 1px solid #404040;
  flex-shrink: 0;
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
  overflow: hidden;
  position: relative;
}

.xterm-container {
  width: 100%;
  height: 100%;
  padding: 8px;
}

:deep(.xterm) {
  height: 100% !important;
}

:deep(.xterm-screen) {
  height: 100% !important;
}

.connecting-state,
.error-state {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  color: #a0a0a0;
  z-index: 10;
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
</style>
