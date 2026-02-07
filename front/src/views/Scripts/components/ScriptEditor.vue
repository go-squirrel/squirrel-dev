<template>
  <div class="script-editor-overlay" @click.self="$emit('cancel')">
    <div class="script-editor-container">
      <div class="editor-header">
        <h3 class="editor-title">
          {{ isEdit ? $t('scripts.editScript') : $t('scripts.addScript') }}
        </h3>
        <button class="close-btn" @click="$emit('cancel')">
          <Icon icon="lucide:x" />
        </button>
      </div>

      <div class="editor-body">
        <div class="form-group">
          <label class="form-label">
            {{ $t('scripts.scriptName') }}
            <span class="required">*</span>
          </label>
          <input
            v-model="formData.name"
            type="text"
            class="form-input"
            :placeholder="$t('scripts.nameRequired')"
          />
          <span v-if="errors.name" class="error-message">{{ errors.name }}</span>
        </div>

        <div class="form-group">
          <label class="form-label">
            {{ $t('scripts.scriptContent') }}
            <span class="required">*</span>
          </label>
          <textarea
            v-model="formData.content"
            class="form-textarea"
            :placeholder="$t('scripts.shebangRequired')"
            spellcheck="false"
          />
          <span v-if="errors.content" class="error-message">{{ errors.content }}</span>
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
import { Icon } from '@iconify/vue'
import { useI18n } from 'vue-i18n'
import Button from '@/components/Button/index.vue'
import type { Script, CreateScriptRequest, UpdateScriptRequest } from '@/types'

const props = defineProps<{
  script: Script | null
}>()

const emit = defineEmits<{
  submit: [data: CreateScriptRequest | UpdateScriptRequest]
  cancel: []
}>()

const { t } = useI18n()

const isEdit = computed(() => props.script !== null)
const submitting = ref(false)

const formData = ref<CreateScriptRequest>({
  name: '',
  content: ''
})

const errors = ref({
  name: '',
  content: ''
})

watch(() => props.script, (newScript) => {
  if (newScript) {
    formData.value = {
      name: newScript.name,
      content: newScript.content
    }
  } else {
    formData.value = {
      name: '',
      content: ''
    }
  }
  errors.value = { name: '', content: '' }
}, { immediate: true })

const validate = (): boolean => {
  errors.value = { name: '', content: '' }
  let isValid = true

  if (!formData.value.name.trim()) {
    errors.value.name = t('scripts.nameRequired')
    isValid = false
  }

  if (!formData.value.content.trim()) {
    errors.value.content = t('scripts.contentRequired')
    isValid = false
  } else if (!formData.value.content.trim().startsWith('#!')) {
    errors.value.content = t('scripts.shebangRequired')
    isValid = false
  }

  return isValid
}

const handleSubmit = async () => {
  if (!validate()) return

  submitting.value = true
  try {
    const data: CreateScriptRequest | UpdateScriptRequest = isEdit.value
      ? { id: props.script!.id, ...formData.value }
      : formData.value
    emit('submit', data)
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.script-editor-overlay {
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

.script-editor-container {
  background: #ffffff;
  border-radius: 16px;
  width: 90%;
  max-width: 800px;
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
  min-height: 300px;
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
