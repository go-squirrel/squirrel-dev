<template>
  <div class="app-form-overlay" @click.self="$emit('cancel')">
    <div class="app-form-container">
      <div class="form-header">
        <h3 class="form-title">
          {{ isEdit ? $t('appStore.editApp') : $t('appStore.addApp') }}
        </h3>
        <button class="close-btn" @click="$emit('cancel')">
          <Icon icon="lucide:x" />
        </button>
      </div>

      <div class="form-body">
        <!-- 基本信息 -->
        <div class="form-section">
          <h4 class="section-title">{{ $t('appStore.basicInfo') }}</h4>
          
          <div class="form-row">
            <div class="form-group">
              <label class="form-label">
                {{ $t('appStore.appName') }}
                <span class="required">*</span>
              </label>
              <input
                v-model="formData.name"
                type="text"
                class="form-input"
                :placeholder="$t('appStore.required')"
              />
              <span v-if="errors.name" class="error-message">{{ errors.name }}</span>
            </div>

            <div class="form-group">
              <label class="form-label">
                {{ $t('appStore.version') }}
                <span class="required">*</span>
              </label>
              <input
                v-model="formData.version"
                type="text"
                class="form-input"
                :placeholder="$t('appStore.required')"
              />
              <span v-if="errors.version" class="error-message">{{ errors.version }}</span>
            </div>
          </div>

          <div class="form-group">
            <label class="form-label">
              {{ $t('appStore.description') }}
              <span class="required">*</span>
            </label>
            <textarea
              v-model="formData.description"
              class="form-textarea"
              :placeholder="$t('appStore.required')"
              rows="3"
            />
            <span v-if="errors.description" class="error-message">{{ errors.description }}</span>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label class="form-label">
                {{ $t('appStore.type') }}
                <span class="required">*</span>
              </label>
              <select v-model="formData.type" class="form-select">
                <option value="compose">{{ $t('appStore.typeCompose') }}</option>
                <option value="k8s_manifest">{{ $t('appStore.typeK8s') }}</option>
                <option value="helm_chart">{{ $t('appStore.typeHelm') }}</option>
              </select>
              <span v-if="errors.type" class="error-message">{{ errors.type }}</span>
            </div>

            <div class="form-group">
              <label class="form-label">
                {{ $t('appStore.category') }}
                <span class="required">*</span>
              </label>
              <select v-model="formData.category" class="form-select">
                <option value="web">{{ $t('appStore.categoryWeb') }}</option>
                <option value="database">{{ $t('appStore.categoryDatabase') }}</option>
                <option value="middleware">{{ $t('appStore.categoryMiddleware') }}</option>
                <option value="devops">{{ $t('appStore.categoryDevops') }}</option>
              </select>
              <span v-if="errors.category" class="error-message">{{ errors.category }}</span>
            </div>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label class="form-label">{{ $t('appStore.icon') }}</label>
              <input
                v-model="formData.icon"
                type="text"
                class="form-input"
                :placeholder="'URL'"
              />
            </div>

            <div class="form-group">
              <label class="form-label">
                {{ $t('appStore.author') }}
                <span class="required">*</span>
              </label>
              <input
                v-model="formData.author"
                type="text"
                class="form-input"
                :placeholder="$t('appStore.required')"
              />
              <span v-if="errors.author" class="error-message">{{ errors.author }}</span>
            </div>
          </div>

          <div class="form-group">
            <label class="form-label">{{ $t('appStore.tags') }}</label>
            <input
              v-model="formData.tags"
              type="text"
              class="form-input"
              :placeholder="'tag1, tag2, tag3'"
            />
          </div>
        </div>

        <!-- 链接信息 -->
        <div class="form-section">
          <h4 class="section-title">{{ $t('appStore.repoUrl') }}</h4>
          
          <div class="form-group">
            <label class="form-label">{{ $t('appStore.repoUrl') }}</label>
            <input
              v-model="formData.repo_url"
              type="text"
              class="form-input"
              placeholder="https://github.com/..."
            />
          </div>

          <div class="form-group">
            <label class="form-label">{{ $t('appStore.homepageUrl') }}</label>
            <input
              v-model="formData.homepage_url"
              type="text"
              class="form-input"
              placeholder="https://..."
            />
          </div>

          <div class="form-group checkbox-group">
            <label class="checkbox-label">
              <input v-model="formData.is_official" type="checkbox" class="checkbox-input" />
              <span class="checkbox-text">{{ $t('appStore.isOfficial') }}</span>
            </label>
          </div>
        </div>

        <!-- 模板内容 -->
        <div class="form-section">
          <h4 class="section-title">{{ $t('appStore.templateContent') }}</h4>
          
          <div class="form-group">
            <label class="form-label">
              {{ $t('appStore.content') }}
              <span class="required">*</span>
            </label>
            <textarea
              v-model="formData.content"
              class="form-textarea code-textarea"
              :placeholder="$t('appStore.required')"
              rows="12"
            />
            <span v-if="errors.content" class="error-message">{{ errors.content }}</span>
          </div>
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
import { Icon } from '@iconify/vue'
import Button from '@/components/Button/index.vue'
import type { AppStore, CreateAppRequest, UpdateAppRequest } from '@/types'

const props = defineProps<{
  app: AppStore | null
}>()

const emit = defineEmits<{
  submit: [data: CreateAppRequest | UpdateAppRequest]
  cancel: []
}>()

const { t } = useI18n()

const isEdit = computed(() => props.app !== null)
const submitting = ref(false)

const formData = ref<CreateAppRequest>({
  name: '',
  description: '',
  type: 'compose',
  category: 'web',
  icon: '',
  version: '',
  content: '',
  tags: '',
  author: '',
  repo_url: '',
  homepage_url: '',
  is_official: false
})

const errors = ref({
  name: '',
  description: '',
  type: '',
  category: '',
  version: '',
  author: '',
  content: ''
})

watch(() => props.app, (newApp) => {
  if (newApp) {
    formData.value = {
      name: newApp.name,
      description: newApp.description,
      type: newApp.type,
      category: newApp.category,
      icon: newApp.icon || '',
      version: newApp.version,
      content: newApp.content,
      tags: newApp.tags,
      author: newApp.author,
      repo_url: newApp.repo_url || '',
      homepage_url: newApp.homepage_url || '',
      is_official: newApp.is_official
    }
  } else {
    formData.value = {
      name: '',
      description: '',
      type: 'compose',
      category: 'web',
      icon: '',
      version: '',
      content: '',
      tags: '',
      author: '',
      repo_url: '',
      homepage_url: '',
      is_official: false
    }
  }
  errors.value = { name: '', description: '', type: '', category: '', version: '', author: '', content: '' }
}, { immediate: true })

const validate = (): boolean => {
  errors.value = { name: '', description: '', type: '', category: '', version: '', author: '', content: '' }
  let isValid = true

  if (!formData.value.name.trim()) {
    errors.value.name = t('appStore.required')
    isValid = false
  }

  if (!formData.value.description.trim()) {
    errors.value.description = t('appStore.required')
    isValid = false
  }

  if (!formData.value.type) {
    errors.value.type = t('appStore.required')
    isValid = false
  }

  if (!formData.value.category) {
    errors.value.category = t('appStore.required')
    isValid = false
  }

  if (!formData.value.version.trim()) {
    errors.value.version = t('appStore.required')
    isValid = false
  }

  if (!formData.value.author.trim()) {
    errors.value.author = t('appStore.required')
    isValid = false
  }

  if (!formData.value.content.trim()) {
    errors.value.content = t('appStore.required')
    isValid = false
  }

  return isValid
}

const handleSubmit = async () => {
  if (!validate()) return

  submitting.value = true
  try {
    const data: CreateAppRequest | UpdateAppRequest = isEdit.value
      ? { id: props.app!.id, ...formData.value }
      : formData.value
    emit('submit', data)
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.app-form-overlay {
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

.app-form-container {
  background: #ffffff;
  border-radius: 16px;
  width: 90%;
  max-width: 800px;
  max-height: 90vh;
  overflow: hidden;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  display: flex;
  flex-direction: column;
}

.form-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid #e2e8f0;
  flex-shrink: 0;
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

.form-section {
  margin-bottom: 28px;
}

.form-section:last-child {
  margin-bottom: 0;
}

.section-title {
  font-size: 14px;
  font-weight: 600;
  color: #1e3a5f;
  margin-bottom: 16px;
  padding-bottom: 8px;
  border-bottom: 1px solid #e2e8f0;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.form-group {
  margin-bottom: 16px;
}

.form-group:last-child {
  margin-bottom: 0;
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

.form-input,
.form-select {
  width: 100%;
  padding: 10px 14px;
  border: 2px solid #e2e8f0;
  border-radius: 8px;
  font-size: 14px;
  color: #1e3a5f;
  background: #f8fafc;
  transition: all 0.2s ease;
}

.form-input:focus,
.form-select:focus {
  outline: none;
  border-color: #4fc3f7;
  background: #ffffff;
}

.form-textarea {
  width: 100%;
  min-height: 80px;
  padding: 12px 14px;
  border: 2px solid #e2e8f0;
  border-radius: 8px;
  font-size: 14px;
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

.code-textarea {
  font-family: 'Fira Code', 'Consolas', monospace;
  font-size: 13px;
  min-height: 200px;
}

.checkbox-group {
  display: flex;
  align-items: center;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.checkbox-input {
  width: 18px;
  height: 18px;
  cursor: pointer;
  accent-color: #4fc3f7;
}

.checkbox-text {
  font-size: 14px;
  color: #1e3a5f;
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
  flex-shrink: 0;
}
</style>
