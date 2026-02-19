<template>
  <Modal :model-value="visible" :title="$t('deployment.editDeployment')" width="600px" show-footer @close="handleCancel">
    <div class="form-content">
      <div class="form-item">
        <label class="form-label">{{ $t('deployment.configContent') }} <span class="required">*</span></label>
        <textarea
          v-model="formData.content"
          class="form-textarea"
          :placeholder="$t('deployment.configContentPlaceholder')"
          rows="15"
        />
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
import type { Deployment, UpdateDeploymentRequest } from '../types'
import Modal from '@/components/Modal/index.vue'
import Button from '@/components/Button/index.vue'

interface Props {
  visible: boolean
  deployment: Deployment | null
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:visible': [value: boolean]
  submit: [id: number, data: UpdateDeploymentRequest]
}>()

const loading = ref(false)
const formData = ref<UpdateDeploymentRequest>({
  content: ''
})

watch(() => props.visible, (val) => {
  if (val && props.deployment) {
    formData.value.content = props.deployment.content || ''
  } else if (!val) {
    formData.value = { content: '' }
    loading.value = false
  }
})

const handleSubmit = async () => {
  if (!formData.value.content.trim()) {
    return
  }

  if (!props.deployment) return

  loading.value = true
  try {
    emit('submit', props.deployment.id, { ...formData.value })
    emit('update:visible', false)
  } finally {
    loading.value = false
  }
}

const handleCancel = () => {
  emit('update:visible', false)
}
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

.form-textarea {
  width: 100%;
  padding: 12px;
  border: 2px solid #e2e8f0;
  border-radius: 6px;
  font-size: 13px;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  color: #1e3a5f;
  background: #f8fafc;
  resize: vertical;
  min-height: 200px;
  transition: all 0.2s ease;
}

.form-textarea:focus {
  outline: none;
  border-color: #4fc3f7;
  background: #ffffff;
}

.form-textarea::placeholder {
  color: #94a3b8;
}
</style>
