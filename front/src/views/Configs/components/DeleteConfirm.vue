<template>
  <div class="delete-confirm-overlay" @click.self="$emit('cancel')">
    <div class="delete-confirm-container">
      <div class="dialog-header">
        <h3 class="dialog-title">{{ $t('configs.confirmDelete') }}</h3>
        <button class="close-btn" @click="$emit('cancel')">
          <Icon icon="lucide:x" />
        </button>
      </div>

      <div class="dialog-body">
        <div class="warning-icon-wrapper">
          <Icon icon="lucide:alert-triangle" class="warning-icon" />
        </div>
        <p class="warning-text">
          {{ $t('configs.deleteWarning', { key: config?.key }) }}
        </p>
      </div>

      <div class="dialog-footer">
        <Button type="secondary" @click="$emit('cancel')">
          {{ $t('common.cancel') }}
        </Button>
        <Button type="danger" @click="$emit('confirm')">
          {{ $t('common.delete') }}
        </Button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import Button from '@/components/Button/index.vue'
import type { Config } from '@/types'

defineProps<{
  config: Config | null
}>()

defineEmits<{
  confirm: []
  cancel: []
}>()
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
  max-width: 400px;
  overflow: hidden;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.dialog-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid #e2e8f0;
}

.dialog-title {
  font-size: 18px;
  font-weight: 600;
  color: #1e3a5f;
}

.close-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: none;
  background: transparent;
  color: #64748b;
  cursor: pointer;
  border-radius: 6px;
  transition: all 0.2s ease;
}

.close-btn:hover {
  background: #f1f5f9;
  color: #1e3a5f;
}

.dialog-body {
  padding: 32px 24px;
  text-align: center;
}

.warning-icon-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 64px;
  height: 64px;
  margin: 0 auto 20px;
  background: rgba(245, 158, 11, 0.1);
  border-radius: 50%;
}

.warning-icon {
  width: 32px;
  height: 32px;
  color: #f59e0b;
}

.warning-text {
  font-size: 14px;
  color: #64748b;
  line-height: 1.6;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 24px;
  border-top: 1px solid #e2e8f0;
  background: #f8fafc;
}
</style>
