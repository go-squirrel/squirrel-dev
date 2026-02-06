<template>
  <div :class="loadingClass">
    <Icon icon="lucide:loader-2" class="spinner" />
    <p v-if="text" class="text">{{ text }}</p>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Icon } from '@iconify/vue'

interface LoadingProps {
  text?: string
  size?: 'small' | 'medium' | 'large'
  fullscreen?: boolean
}

const props = withDefaults(defineProps<LoadingProps>(), {
  text: '',
  size: 'medium',
  fullscreen: false
})

const loadingClass = computed(() => {
  return [
    's-loading',
    `s-loading--${props.size}`,
    {
      's-loading--fullscreen': props.fullscreen
    }
  ]
})
</script>

<style scoped>
.s-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #64748b;
}

.s-loading--small .spinner {
  width: 32px;
  height: 32px;
}

.s-loading--small .text {
  font-size: 12px;
}

.s-loading--medium .spinner {
  width: 48px;
  height: 48px;
}

.s-loading--medium .text {
  font-size: 14px;
}

.s-loading--large .spinner {
  width: 64px;
  height: 64px;
}

.s-loading--large .text {
  font-size: 16px;
}

.spinner {
  animation: spin 1s linear infinite;
}

.text {
  margin-top: 16px;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

.s-loading--fullscreen {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255, 255, 255, 0.9);
  z-index: 9999;
}
</style>
