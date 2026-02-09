<template>
  <Transition name="toast">
    <div v-if="visible" class="toast" :class="`toast--${type}`">
      <Icon :icon="getIcon()" class="toast-icon" />
      <span class="toast-message">{{ message }}</span>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue'

interface Props {
  visible?: boolean
  message?: string
  type?: 'success' | 'error' | 'warning' | 'info'
}

const props = withDefaults(defineProps<Props>(), {
  visible: false,
  message: '',
  type: 'success'
})

const getIcon = () => {
  const icons: Record<string, string> = {
    success: 'lucide:check-circle',
    error: 'lucide:x-circle',
    warning: 'lucide:alert-triangle',
    info: 'lucide:info'
  }
  return icons[props.type]
}
</script>

<style scoped>
.toast {
  position: fixed;
  top: 20px;
  right: 20px;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 16px;
  border-radius: 6px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  z-index: 9999;
  min-width: 200px;
}

.toast--success {
  background: #dcfce7;
  color: #16a34a;
}

.toast--error {
  background: #fee2e2;
  color: #dc2626;
}

.toast--warning {
  background: #fef3c7;
  color: #d97706;
}

.toast--info {
  background: #e0f2fe;
  color: #0284c7;
}

.toast-icon {
  width: 20px;
  height: 20px;
  flex-shrink: 0;
}

.toast-message {
  font-size: 13px;
  font-weight: 500;
}

.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}

.toast-enter-from {
  opacity: 0;
  transform: translateX(100%);
}

.toast-leave-to {
  opacity: 0;
  transform: translateX(100%);
}
</style>
