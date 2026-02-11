<template>
  <div class="config-editor-overlay" @click.self="$emit('cancel')">
    <div class="config-editor-container">
      <div class="editor-header">
        <h3 class="editor-title">
          {{ isEdit ? $t('configs.editConfig') : $t('configs.addConfig') }}
        </h3>
        <button class="close-btn" @click="$emit('cancel')">
          <Icon icon="lucide:x" />
        </button>
      </div>

      <div class="editor-body">
        <div class="form-group">
          <label class="form-label">
            {{ $t('configs.configKey') }}
            <span class="required">*</span>
          </label>
          <input
            v-model="formData.key"
            type="text"
            class="form-input"
            :placeholder="$t('configs.required')"
          />
          <span v-if="errors.key" class="error-message">{{ errors.key }}</span>
        </div>

        <div class="form-group">
          <label class="form-label">
            {{ $t('configs.configValue') }}
            <span class="required">*</span>
          </label>
          <textarea
            v-model="formData.value"
            class="form-textarea"
            :placeholder="$t('configs.required')"
            rows="8"
          />
          <span v-if="errors.value" class="error-message">{{ errors.value }}</span>
        </div>
      </div>

      <div class="editor-footer">
        <Button type="secondary" @click="$emit('cancel')">
          {{ $t('common.cancel') }}
        </Button>
        <Button type="primary" :loading="submitting" @click="handleSubmit">
          {{ $t('common.save') }}
        </Button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import Button from '@/components/Button/index.vue'
import type { Config, CreateConfigRequest, UpdateConfigRequest } from '@/types'

const props = defineProps<{
  config: Config | null
}>()

const emit = defineEmits<{
  submit: [data: CreateConfigRequest | UpdateConfigRequest]
  cancel: []
}>()

const { t } = useI18n()

const isEdit = computed(() => props.config !== null)
const submitting = ref(false)

const formData = ref<CreateConfigRequest>({
  key: '',
  value: ''
})

const errors = ref({
  key: '',
  value: ''
})

watch(() => props.config, (newConfig) => {
  if (newConfig) {
    formData.value = {
      key: newConfig.key,
      value: newConfig.value
    }
  } else {
    formData.value = {
      key: '',
      value: ''
    }
  }
  errors.value = { key: '', value: '' }
}, { immediate: true })

const validate = (): boolean => {
  errors.value = { key: '', value: '' }
  let isValid = true

  if (!formData.value.key.trim()) {
    errors.value.key = t('configs.required')
    isValid = false
  }

  if (!formData.value.value.trim()) {
    errors.value.value = t('configs.required')
    isValid = false
  }

  return isValid
}

const handleSubmit = async () => {
  if (!validate()) return

  submitting.value = true
  try {
    const data: CreateConfigRequest | UpdateConfigRequest = isEdit.value
      ? { id: props.config!.id, ...formData.value }
      : formData.value
    emit('submit', data)
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.config-editor-overlay {
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

.config-editor-container {
  background: #ffffff;
  border-radius: 16px;
  width: 90%;
  max-width: 600px;
  max-height: 90vh;
  overflow: hidden;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.editor-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid #e2e8f0;
}

.editor-title {
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

.editor-body {
  padding: 24px;
  overflow-y: auto;
  max-height: calc(90vh - 140px);
}

.form-group {
  margin-bottom: 20px;
}

.form-label {
  display: block;
  margin-bottom: 8px;
  font-size: 14px;
  font-weight: 500;
  color: #1e3a5f;
}

.form-label .required {
  color: #ef4444;
  margin-left: 4px;
}

.form-input {
  width: 100%;
  padding: 10px 14px;
  border: 2px solid #e2e8f0;
  border-radius: 8px;
  font-size: 14px;
  color: #1e3a5f;
  background: #f8fafc;
  transition: all 0.2s ease;
}

.form-input:focus {
  outline: none;
  border-color: #4fc3f7;
  background: #ffffff;
}

.form-textarea {
  width: 100%;
  min-height: 120px;
  padding: 12px 14px;
  border: 2px solid #e2e8f0;
  border-radius: 8px;
  font-size: 13px;
  font-family: 'Fira Code', 'Consolas', monospace;
  color: #1e3a5f;
  background: #f8fafc;
  transition: all 0.2s ease;
  resize: vertical;
  line-height: 1.6;
}

.form-textarea:focus {
  outline: none;
  border-color: #4fc3f7;
  background: #ffffff;
}

.error-message {
  margin-top: 6px;
  font-size: 12px;
  color: #ef4444;
}

.editor-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 24px;
  border-top: 1px solid #e2e8f0;
  background: #f8fafc;
}
</style>
