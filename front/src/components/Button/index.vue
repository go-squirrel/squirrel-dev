<template>
  <button
    :class="buttonClass"
    :disabled="disabled || loading"
    @click="handleClick"
  >
    <Icon v-if="loading" icon="lucide:loader-2" class="loading-icon" />
    <Icon v-else-if="icon" :icon="icon" class="icon" />
    <span v-if="$slots.default" class="text">
      <slot></slot>
    </span>
  </button>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface ButtonProps {
  type?: 'primary' | 'secondary' | 'danger' | 'ghost'
  size?: 'small' | 'medium' | 'large'
  disabled?: boolean
  loading?: boolean
  icon?: string
  block?: boolean
  plain?: boolean
  round?: boolean
  circle?: boolean
}

const props = withDefaults(defineProps<ButtonProps>(), {
  type: 'primary',
  size: 'medium',
  disabled: false,
  loading: false,
  block: false,
  plain: false,
  round: false,
  circle: false
})

const emit = defineEmits<{
  click: [event: MouseEvent]
}>()

const buttonClass = computed(() => {
  return [
    's-button',
    `s-button--${props.type}`,
    `s-button--${props.size}`,
    {
      's-button--disabled': props.disabled || props.loading,
      's-button--loading': props.loading,
      's-button--block': props.block,
      's-button--plain': props.plain,
      's-button--round': props.round,
      's-button--circle': props.circle
    }
  ]
})

const handleClick = (event: MouseEvent) => {
  if (!props.disabled && !props.loading) {
    emit('click', event)
  }
}
</script>

<style scoped>
.s-button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 8px 16px;
  border: none;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  white-space: nowrap;
  outline: none;
  user-select: none;
}

.s-button--small {
  padding: 6px 12px;
  font-size: 12px;
}

.s-button--medium {
  padding: 8px 16px;
  font-size: 13px;
}

.s-button--large {
  padding: 12px 24px;
  font-size: 14px;
}

.s-button--primary {
  background: linear-gradient(135deg, #4fc3f7 0%, #29b6f6 100%);
  color: #ffffff;
}

.s-button--primary:not(.s-button--disabled):not(.s-button--loading):hover {
  box-shadow: 0 4px 12px rgba(79, 195, 247, 0.4);
  transform: translateY(-1px);
}

.s-button--secondary {
  background: #f1f5f9;
  color: #475569;
}

.s-button--secondary:not(.s-button--disabled):not(.s-button--loading):hover {
  background: #e2e8f0;
}

.s-button--danger {
  background: linear-gradient(135deg, #f87171 0%, #ef4444 100%);
  color: #ffffff;
}

.s-button--danger:not(.s-button--disabled):not(.s-button--loading):hover {
  box-shadow: 0 4px 12px rgba(239, 68, 68, 0.4);
  transform: translateY(-1px);
}

.s-button--ghost {
  background: transparent;
  color: #475569;
}

.s-button--ghost:not(.s-button--disabled):not(.s-button--loading):hover {
  background: #f1f5f9;
}

.s-button--plain.s-button--primary {
  background: transparent;
  color: #29b6f6;
  border: 1px solid #29b6f6;
}

.s-button--plain.s-button--primary:not(.s-button--disabled):not(.s-button--loading):hover {
  background: rgba(41, 182, 246, 0.1);
}

.s-button--disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.s-button--loading {
  cursor: not-allowed;
}

.s-button--block {
  width: 100%;
}

.s-button--round {
  border-radius: 20px;
}

.s-button--circle {
  border-radius: 50%;
  padding: 8px;
}

.icon {
  width: 16px;
  height: 16px;
}

.loading-icon {
  width: 16px;
  height: 16px;
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
