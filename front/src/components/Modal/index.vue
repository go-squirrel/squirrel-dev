<template>
  <Teleport to="body">
    <Transition name="modal-fade">
      <div v-if="modelValue" class="s-modal-mask" @click="handleMaskClick"></div>
    </Transition>
    <Transition name="modal-slide">
      <div v-if="modelValue" :class="modalClass" :style="modalStyle">
        <div v-if="showHeader || $slots.header" class="s-modal-header">
          <slot name="header">
            <span class="s-modal-title">{{ title }}</span>
          </slot>
          <Icon v-if="showClose" icon="lucide:x" class="s-modal-close" @click="handleClose" />
        </div>
        <div class="s-modal-body">
          <slot></slot>
        </div>
        <div v-if="showFooter || $slots.footer" class="s-modal-footer">
          <slot name="footer"></slot>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { computed, watch, onMounted, onUnmounted } from 'vue'

interface ModalProps {
  modelValue?: boolean
  title?: string
  width?: string | number
  fullscreen?: boolean
  closeOnClickModal?: boolean
  closeOnPressEscape?: boolean
  showClose?: boolean
  showHeader?: boolean
  showFooter?: boolean
  beforeClose?: () => boolean | Promise<boolean>
}

const props = withDefaults(defineProps<ModalProps>(), {
  modelValue: false,
  title: '',
  width: '520px',
  fullscreen: false,
  closeOnClickModal: true,
  closeOnPressEscape: true,
  showClose: true,
  showHeader: true,
  showFooter: false
})

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  close: []
  open: []
  opened: []
  closed: []
}>()

const modalClass = computed(() => {
  return [
    's-modal',
    {
      's-modal--fullscreen': props.fullscreen
    }
  ]
})

const modalStyle = computed(() => {
  if (props.fullscreen) return {}
  return {
    width: typeof props.width === 'number' ? `${props.width}px` : props.width
  }
})

const handleClose = async () => {
  if (props.beforeClose) {
    const canClose = await props.beforeClose()
    if (!canClose) return
  }
  emit('update:modelValue', false)
  emit('close')
}

const handleMaskClick = () => {
  if (props.closeOnClickModal) {
    handleClose()
  }
}

const handleEscape = (e: KeyboardEvent) => {
  if (e.key === 'Escape' && props.closeOnPressEscape) {
    handleClose()
  }
}

watch(() => props.modelValue, (val) => {
  if (val) {
    emit('open')
  }
})

onMounted(() => {
  document.addEventListener('keydown', handleEscape)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleEscape)
})
</script>

<style scoped>
.s-modal-mask {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 2000;
}

.s-modal {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background: #ffffff;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  z-index: 2001;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
}

.s-modal--fullscreen {
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  transform: none;
  width: 100% !important;
  height: 100%;
  max-height: 100vh;
  border-radius: 0;
}

.s-modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px;
  border-bottom: 1px solid #e5e7eb;
}

.s-modal-title {
  font-size: 16px;
  font-weight: 600;
  color: #1e3a5f;
}

.s-modal-close {
  width: 20px;
  height: 20px;
  color: #94a3b8;
  cursor: pointer;
  transition: color 0.3s ease;
}

.s-modal-close:hover {
  color: #64748b;
}

.s-modal-body {
  padding: 20px;
  overflow-y: auto;
  flex: 1;
}

.s-modal-footer {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 20px;
  border-top: 1px solid #e5e7eb;
}

.modal-fade-enter-active,
.modal-fade-leave-active {
  transition: opacity 0.3s ease;
}

.modal-fade-enter-from,
.modal-fade-leave-to {
  opacity: 0;
}

.modal-slide-enter-active,
.modal-slide-leave-active {
  transition: all 0.3s ease;
}

.modal-slide-enter-from,
.modal-slide-leave-to {
  opacity: 0;
  transform: translate(-50%, -60%);
}
</style>
