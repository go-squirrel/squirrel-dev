<template>
  <Modal :model-value="visible" :title="$t('deployment.addDeployment')" width="500px" show-footer @close="handleCancel">
    <div class="form-content">
      <div class="form-item">
        <label class="form-label">{{ $t('deployment.selectApplication') }} <span class="required">*</span></label>
        <div class="select-wrapper">
          <select v-model="formData.application_id" class="form-select">
            <option value="0" disabled>{{ $t('deployment.selectApplication') }}</option>
            <option v-for="app in applications" :key="app.id" :value="app.id">
              {{ app.name }} ({{ app.version }})
            </option>
          </select>
          <Icon icon="lucide:chevron-down" class="select-icon" />
        </div>
      </div>

      <div class="form-item">
        <label class="form-label">{{ $t('deployment.selectServer') }} <span class="required">*</span></label>
        <div class="select-wrapper">
          <select v-model="formData.server_id" class="form-select">
            <option value="0" disabled>{{ $t('deployment.selectServer') }}</option>
            <option v-for="server in servers" :key="server.id" :value="server.id">
              {{ server.ip_address }}{{ server.server_alias ? ` (${server.server_alias})` : '' }}
            </option>
          </select>
          <Icon icon="lucide:chevron-down" class="select-icon" />
        </div>
      </div>
    </div>

    <template #footer>
      <Button type="secondary" @click="handleCancel">{{ $t('common.cancel') }}</Button>
      <Button type="primary" :loading="loading" @click="handleSubmit">{{ $t('common.confirm') }}</Button>
    </template>
  </Modal>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { Icon } from '@iconify/vue'
import type { ApplicationInstance, Server } from '@/types'
import type { CreateDeploymentRequest } from '../types'
import Modal from '@/components/Modal/index.vue'
import Button from '@/components/Button/index.vue'

interface Props {
  visible: boolean
  applications: ApplicationInstance[]
  servers: Server[]
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:visible': [value: boolean]
  submit: [data: CreateDeploymentRequest]
}>()

const loading = ref(false)
const formData = ref<CreateDeploymentRequest>({
  application_id: 0,
  server_id: 0
})

const handleSubmit = async () => {
  if (formData.value.application_id === 0) {
    return
  }
  if (formData.value.server_id === 0) {
    return
  }

  loading.value = true
  try {
    emit('submit', { ...formData.value })
    emit('update:visible', false)
  } finally {
    loading.value = false
  }
}

const handleCancel = () => {
  emit('update:visible', false)
}

watch(() => props.visible, (val) => {
  if (!val) {
    formData.value = { application_id: 0, server_id: 0 }
    loading.value = false
  }
})
</script>

<style scoped>
.form-content {
  padding: 20px 0;
}

.form-item {
  margin-bottom: 20px;
}

.form-item:last-child {
  margin-bottom: 0;
}

.form-label {
  display: block;
  margin-bottom: 8px;
  font-size: 13px;
  font-weight: 500;
  color: #1e3a5f;
}

.required {
  color: #dc2626;
}

.select-wrapper {
  position: relative;
}

.form-select {
  width: 100%;
  padding: 10px 36px 10px 12px;
  border: 2px solid #e2e8f0;
  border-radius: 6px;
  font-size: 13px;
  color: #1e3a5f;
  background: #f8fafc;
  cursor: pointer;
  appearance: none;
  transition: all 0.2s ease;
}

.form-select:focus {
  outline: none;
  border-color: #4fc3f7;
  background: #ffffff;
}

.select-icon {
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  width: 16px;
  height: 16px;
  color: #64748b;
  pointer-events: none;
}
</style>
