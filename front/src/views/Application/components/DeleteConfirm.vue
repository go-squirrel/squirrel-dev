<template>
  <div class="delete-confirm-overlay" @click.self="$emit('cancel')">
    <div class="delete-confirm-container">
      <div class="confirm-icon">
        <Icon icon="lucide:alert-triangle" />
      </div>
      <h3 class="confirm-title">{{ $t('application.confirmDelete') }}</h3>
      <p class="confirm-message">
        {{ $t('application.deleteWarning', { name: application?.name }) }}
      </p>
      <div class="confirm-actions">
        <Button type="secondary" @click="$emit('cancel')">
          {{ $t('common.cancel') }}
        </Button>
        <Button type="danger" :loading="loading" @click="handleConfirm">
          {{ $t('common.confirm') }}
        </Button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Icon } from '@iconify/vue'
import Button from '@/components/Button/index.vue'
import type { ApplicationInstance } from '@/types'

defineProps<{
  application: ApplicationInstance | null
}>()

const emit = defineEmits<{
  confirm: []
  cancel: []
}>()

const loading = ref(false)

const handleConfirm = async () => {
  loading.value = true
  try {
    emit('confirm')
  } finally {
    loading.value = false
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
  padding: 32px;
  max-width: 400px;
  text-align: center;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.confirm-icon {
  width: 48px;
  height: 48px;
  margin: 0 auto 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #fee2e2;
  border-radius: 50%;
  color: #dc2626;
}

.confirm-icon svg {
  width: 24px;
  height: 24px;
}

.confirm-title {
  font-size: 18px;
  font-weight: 600;
  color: #1e3a5f;
  margin-bottom: 12px;
}

.confirm-message {
  font-size: 14px;
  color: #64748b;
  line-height: 1.6;
  margin-bottom: 24px;
}

.confirm-actions {
  display: flex;
  justify-content: center;
  gap: 12px;
}
</style>
