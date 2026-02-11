<template>
  <div class="application-form-overlay" @click.self="$emit('cancel')">
    <div class="application-form-container">
      <div class="form-header">
        <h3 class="form-title">
          {{ isEdit ? $t('application.editApplication') : $t('application.addApplication') }}
        </h3>
        <button class="close-btn" @click="$emit('cancel')">
          <Icon icon="lucide:x" />
        </button>
      </div>

      <div class="form-body">
        <div class="form-group">
          <label class="form-label">
            {{ $t('application.applicationName') }}
            <span class="required">*</span>
          </label>
          <input
            v-model="formData.name"
            type="text"
            class="form-input"
            :placeholder="$t('application.required')"
          />
          <span v-if="errors.name" class="error-message">{{ errors.name }}</span>
        </div>

        <div class="form-group">
          <label class="form-label">
            {{ $t('application.applicationType') }}
            <span class="required">*</span>
          </label>
          <select v-model="formData.type" class="form-select">
            <option value="compose">{{ $t('application.typeCompose') }}</option>
            <option value="script">{{ $t('application.typeScript') }}</option>
            <option value="binary">{{ $t('application.typeBinary') }}</option>
          </select>
          <span v-if="errors.type" class="error-message">{{ errors.type }}</span>
        </div>

        <div class="form-group">
          <label class="form-label">
            {{ $t('application.applicationVersion') }}
            <span class="required">*</span>
          </label>
          <input
            v-model="formData.version"
            type="text"
            class="form-input"
            :placeholder="$t('application.required')"
          />
          <span v-if="errors.version" class="error-message">{{ errors.version }}</span>
        </div>

        <div class="form-group">
          <label class="form-label">
            {{ $t('application.applicationDescription') }}
          </label>
          <textarea
            v-model="formData.description"
            class="form-textarea"
            :placeholder="$t('application.applicationDescription')"
            rows="3"
          />
          <span v-if="errors.description" class="error-message">{{ errors.description }}</span>
        </div>

        <div class="form-group">
          <label class="form-label">
            {{ $t('application.applicationContent') }}
            <span class="required">*</span>
          </label>
          <textarea
            v-model="formData.content"
            class="form-textarea code-editor"
            :placeholder="$t('application.applicationContent')"
            rows="12"
          />
          <span v-if="errors.content" class="error-message">{{ errors.content }}</span>
        </div>
      </div>

      <div class="form-footer">
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
import type { ApplicationInstance, CreateApplicationRequest, UpdateApplicationRequest } from '@/types'

const props = defineProps<{
  application: ApplicationInstance | null
}>()

const emit = defineEmits<{
  submit: [data: CreateApplicationRequest | UpdateApplicationRequest]
  cancel: []
}>()

const { t } = useI18n()

const isEdit = computed(() => props.application !== null)
const submitting = ref(false)

const formData = ref<CreateApplicationRequest>({
  name: '',
  description: '',
  type: 'compose',
  content: '',
  version: ''
})

const errors = ref({
  name: '',
  type: '',
  version: '',
  description: '',
  content: ''
})

watch(() => props.application, (newApplication) => {
  if (newApplication) {
    formData.value = {
      name: newApplication.name,
      description: newApplication.description,
      type: newApplication.type,
      content: newApplication.content,
      version: newApplication.version
    }
  } else {
    formData.value = {
      name: '',
      description: '',
      type: 'compose',
      content: '',
      version: ''
    }
  }
  errors.value = { name: '', type: '', version: '', description: '', content: '' }
}, { immediate: true })

const validate = (): boolean => {
  errors.value = { name: '', type: '', version: '', description: '', content: '' }
  let isValid = true

  if (!formData.value.name.trim()) {
    errors.value.name = t('application.required')
    isValid = false
  }

  if (!formData.value.type) {
    errors.value.type = t('application.required')
    isValid = false
  }

  if (!formData.value.version.trim()) {
    errors.value.version = t('application.required')
    isValid = false
  }

  if (!formData.value.content.trim()) {
    errors.value.content = t('application.required')
    isValid = false
  }

  return isValid
}

const handleSubmit = async () => {
  if (!validate()) return

  submitting.value = true
  try {
    const data: CreateApplicationRequest | UpdateApplicationRequest = isEdit.value
      ? { id: props.application!.id, ...formData.value }
      : formData.value
    emit('submit', data)
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.application-form-overlay {
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

.application-form-container {
  background: #ffffff;
  border-radius: 16px;
  width: 90%;
  max-width: 800px;
  max-height: 90vh;
  overflow: hidden;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.form-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid #e2e8f0;
}

.form-title {
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

.form-body {
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

.form-select {
  width: 100%;
  padding: 10px 14px;
  border: 2px solid #e2e8f0;
  border-radius: 8px;
  font-size: 14px;
  color: #1e3a5f;
  background: #f8fafc;
  transition: all 0.2s ease;
  cursor: pointer;
}

.form-select:focus {
  outline: none;
  border-color: #4fc3f7;
  background: #ffffff;
}

.form-textarea {
  width: 100%;
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

.form-textarea.code-editor {
  font-size: 12px;
  background: #1e293b;
  color: #e2e8f0;
}

.form-textarea.code-editor:focus {
  background: #0f172a;
  color: #f8fafc;
}

.error-message {
  margin-top: 6px;
  font-size: 12px;
  color: #ef4444;
}

.form-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 24px;
  border-top: 1px solid #e2e8f0;
  background: #f8fafc;
}
</style>
