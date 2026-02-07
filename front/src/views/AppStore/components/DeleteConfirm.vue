<template>
  <div class="delete-confirm-overlay" @click.self="$emit('cancel')">
    <div class="delete-confirm-container">
      <div class="confirm-header">
        <div class="warning-icon">
          <Icon icon="lucide:alert-triangle" />
        </div>
        <h3 class="confirm-title">{{ $t('appStore.confirmDelete') }}</h3>
      </div>

      <div class="confirm-body">
        <p class="confirm-message">
          {{ $t('appStore.deleteWarning', { name: app?.name }) }}
        </p>
      </div>

      <div class="confirm-footer">
        <Button type="secondary" @click="$emit('cancel')">
          {{ $t('common.cancel') }}
        </Button>
        <Button type="danger" :loading="deleting" @click="handleConfirm">
          {{ $t('common.delete') }}
        </Button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Icon } from '@iconify/vue'
import Button from '@/components/Button/index.vue'
import type { AppStore } from '@/types'

defineProps<{
  app: AppStore | null
}>()

const emit = defineEmits<{
  confirm: []
  cancel: []
}>()

const deleting = ref(false)

const handleConfirm = async () => {
  deleting.value = true
  try {
    emit('confirm')
  } finally {
    deleting.value = false
  }
}
</script>

<style scoped>
.delete-confirm-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.delete-confirm-container {
  background: #ffffff;
  border-radius: 16px;
  width: 90%;
  max-width: 420px;
  overflow: hidden;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.confirm-header {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 32px 24px 16px;
}

.warning-icon {
  width: 64px;
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #fef3c7;
  border-radius: 50%;
  color: #f59e0b;
  font-size: 32px;
  margin-bottom: 16px;
}

.confirm-title {
  font-size: 18px;
  font-weight: 600;
  color: #1e3a5f;
}

.confirm-body {
  padding: 0 24px 24px;
  text-align: center;
}

.confirm-message {
  font-size: 14px;
  color: #64748b;
  line-height: 1.6;
}

.confirm-footer {
  display: flex;
  justify-content: center;
  gap: 12px;
  padding: 16px 24px 24px;
}
</style>
