<template>
  <div class="modal-overlay" @click.self="$emit('cancel')">
    <div class="modal">
      <div class="modal-header">
        <h3>{{ isEdit ? $t('configs.editConfig') : $t('configs.addConfig') }}</h3>
        <button class="close-btn" @click="$emit('cancel')">
          <Icon icon="lucide:x" />
        </button>
      </div>

      <form @submit.prevent="handleSubmit" class="modal-body">
        <div class="form-section">
          <h4>{{ $t('configs.basicInfo') }}</h4>
          <div class="form-group">
            <label>{{ $t('configs.configKey') }} *</label>
            <input
              v-model="formData.key"
              type="text"
              :placeholder="$t('configs.required')"
              :class="{ error: errors.key }"
            />
            <span v-if="errors.key" class="error-text">{{ errors.key }}</span>
          </div>

          <div class="form-group">
            <label>{{ $t('configs.configValue') }} *</label>
            <textarea
              v-model="formData.value"
              rows="6"
              :placeholder="$t('configs.required')"
              :class="{ error: errors.value }"
            ></textarea>
            <span v-if="errors.value" class="error-text">{{ errors.value }}</span>
          </div>
        </div>

        <div class="modal-footer">
          <button type="button" class="btn btn-cancel" @click="$emit('cancel')">
            {{ $t('configs.cancel') }}
          </button>
          <button type="submit" class="btn btn-primary" :disabled="submitting">
            {{ submitting ? $t('common.loading') : $t('configs.save') }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { Icon } from '@iconify/vue'
import { createConfig, updateConfig } from '@/api/config'
import type { Config, CreateConfigRequest, UpdateConfigRequest } from '@/types'

const props = defineProps<{
  config?: Config | null
}>()

const emit = defineEmits<{
  submit: []
  cancel: []
}>()

const { t } = useI18n()

const isEdit = computed(() => !!props.config)
const submitting = ref(false)

const formData = reactive<CreateConfigRequest>({
  key: '',
  value: ''
})

const errors = reactive<Record<string, string>>({})

watch(() => props.config, (config) => {
  if (config) {
    formData.key = config.key
    formData.value = config.value
  } else {
    resetForm()
  }
}, { immediate: true })

const resetForm = () => {
  formData.key = ''
  formData.value = ''
  Object.keys(errors).forEach(key => delete errors[key])
}

const validate = () => {
  Object.keys(errors).forEach(key => delete errors[key])

  if (!formData.key) {
    errors.key = t('configs.required')
  }

  if (!formData.value) {
    errors.value = t('configs.required')
  }

  return Object.keys(errors).length === 0
}

const handleSubmit = async () => {
  if (!validate()) return

  submitting.value = true
  try {
    if (isEdit.value && props.config) {
      const updateData: UpdateConfigRequest = {
        id: props.config.id,
        key: formData.key,
        value: formData.value
      }
      await updateConfig(props.config.id, updateData)
    } else {
      await createConfig(formData)
    }
    emit('submit')
  } catch (error) {
    console.error('Failed to save config:', error)
  } finally {
    submitting.value = false
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
  max-width: 500px;
  width: 100%;
  max-height: 90vh;
  overflow-y: auto;
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid #f1f5f9;
}

.modal-header h3 {
  font-size: 16px;
  font-weight: 600;
  color: #1e3a5f;
}

.close-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border-radius: 6px;
  background: #f5f7fa;
  color: #64748b;
  border: none;
  cursor: pointer;
  transition: all 0.2s ease;
}

.close-btn:hover {
  background: #fee2e2;
  color: #dc2626;
}

.modal-body {
  padding: 24px;
}

.form-section {
  margin-bottom: 24px;
}

.form-section:last-of-type {
  margin-bottom: 0;
}

.form-section h4 {
  font-size: 13px;
  font-weight: 600;
  color: #1e3a5f;
  margin-bottom: 16px;
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  font-size: 12px;
  font-weight: 500;
  color: #64748b;
  margin-bottom: 6px;
}

.form-group input,
.form-group select,
.form-group textarea {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  font-size: 13px;
  color: #1e3a5f;
  transition: all 0.2s ease;
  background: #ffffff;
}

.form-group input:focus,
.form-group select:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #4fc3f7;
  box-shadow: 0 0 0 3px rgba(79, 195, 247, 0.1);
}

.form-group input.error,
.form-group select.error,
.form-group textarea.error {
  border-color: #dc2626;
}

.form-group textarea {
  resize: vertical;
  font-family: 'SF Mono', Monaco, Consolas, monospace;
  font-size: 12px;
}

.error-text {
  display: block;
  font-size: 11px;
  color: #dc2626;
  margin-top: 4px;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding-top: 20px;
  border-top: 1px solid #f1f5f9;
}

.btn {
  padding: 10px 20px;
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

.btn-primary {
  background: linear-gradient(135deg, #4fc3f7 0%, #29b6f6 100%);
  color: #ffffff;
}

.btn-primary:hover:not(:disabled) {
  box-shadow: 0 4px 12px rgba(79, 195, 247, 0.4);
  transform: translateY(-1px);
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>
