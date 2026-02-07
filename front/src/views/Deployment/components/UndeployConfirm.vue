<template>
  <Modal :model-value="visible" :title="$t('deployment.confirmUndeploy')" width="400px" show-footer @close="handleCancel">
    <div class="confirm-content">
      <div class="warning-icon">
        <Icon icon="lucide:alert-triangle" />
      </div>
      <p class="confirm-message">
        {{ $t('deployment.undeployWarning', { name: deployment?.application?.name || '' }) }}
      </p>
    </div>

    <template #footer>
      <Button type="secondary" @click="handleCancel">{{ $t('common.cancel') }}</Button>
      <Button type="danger" :loading="loading" @click="handleConfirm">{{ $t('common.confirm') }}</Button>
    </template>
  </Modal>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Icon } from '@iconify/vue'
import type { Deployment } from '../types'
import Modal from '@/components/Modal/index.vue'
import Button from '@/components/Button/index.vue'

defineProps<{
  visible: boolean
  deployment: Deployment | null
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

const handleCancel = () => {
  emit('cancel')
}
</script>

<style scoped>
.confirm-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px 0;
  text-align: center;
}

.warning-icon {
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #fef3c7;
  border-radius: 50%;
  margin-bottom: 16px;
  color: #f59e0b;
  font-size: 24px;
}

.confirm-message {
  font-size: 14px;
  color: #475569;
  line-height: 1.6;
}
</style>
