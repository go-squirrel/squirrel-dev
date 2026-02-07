<template>
  <div class="modal-overlay" @click.self="$emit('cancel')">
    <div class="modal">
      <div class="modal-body">
        <div class="icon-wrapper">
          <Icon icon="lucide:alert-triangle" class="warning-icon" />
        </div>
        <h3>{{ $t('configs.confirmDelete') }}</h3>
        <p class="warning-text">
          {{ $t('configs.deleteWarning', { key: config.key }) }}
        </p>
        <div class="modal-footer">
          <button class="btn btn-cancel" @click="$emit('cancel')">
            {{ $t('configs.cancel') }}
          </button>
          <button class="btn btn-danger" @click="handleConfirm" :disabled="deleting">
            {{ deleting ? $t('common.loading') : $t('configs.confirm') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Icon } from '@iconify/vue'
import { deleteConfig } from '@/api/config'
import type { Config } from '@/types'

const props = defineProps<{
  config: Config
}>()

const emit = defineEmits<{
  confirm: []
  cancel: []
}>()

const deleting = ref(false)

const handleConfirm = async () => {
  deleting.value = true
  try {
    await deleteConfig(props.config.id)
    emit('confirm')
  } catch (error) {
    console.error('Failed to delete config:', error)
  } finally {
    deleting.value = false
  }
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  padding: 20px;
}

.modal {
  background: #ffffff;
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.12);
  max-width: 400px;
  width: 100%;
}

.modal-body {
  padding: 32px 24px;
  text-align: center;
}

.icon-wrapper {
  display: flex;
  justify-content: center;
  margin-bottom: 16px;
}

.warning-icon {
  width: 48px;
  height: 48px;
  color: #f59e0b;
}

.modal-body h3 {
  font-size: 18px;
  font-weight: 600;
  color: #1e3a5f;
  margin-bottom: 12px;
}

.warning-text {
  font-size: 14px;
  color: #64748b;
  line-height: 1.6;
  margin-bottom: 24px;
}

.modal-footer {
  display: flex;
  justify-content: center;
  gap: 12px;
}

.btn {
  padding: 10px 24px;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  border: none;
}

.btn-cancel {
  background: #f5f7fa;
  color: #64748b;
}

.btn-cancel:hover {
  background: #e2e8f0;
  color: #1e3a5f;
}

.btn-danger {
  background: linear-gradient(135deg, #dc2626 0%, #b91c1c 100%);
  color: #ffffff;
}

.btn-danger:hover {
  box-shadow: 0 4px 12px rgba(220, 38, 38, 0.4);
  transform: translateY(-1px);
}
</style>
