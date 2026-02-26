<template>
  <div class="modal-overlay" @click.self="$emit('cancel')">
    <div class="modal">
      <div class="modal-header">
        <h3>{{ isEdit ? $t('server.editServer') : $t('server.addServer') }}</h3>
        <button class="close-btn" @click="$emit('cancel')">
          <Icon icon="lucide:x" />
        </button>
      </div>

      <form @submit.prevent="handleSubmit" class="modal-body">
        <div class="form-section">
          <h4>{{ $t('server.basicInfo') }}</h4>
          <div class="form-group">
            <label>{{ $t('server.ipAddress') }} *</label>
            <input
              v-model="formData.ip_address"
              type="text"
              :placeholder="$t('server.invalidIp')"
              :class="{ error: errors.ip_address }"
            />
            <span v-if="errors.ip_address" class="error-text">{{ errors.ip_address }}</span>
          </div>

          <div class="form-group">
            <label>{{ $t('server.agentPort') }} *</label>
            <input
              v-model.number="formData.port"
              type="number"
              min="1"
              max="65535"
              :class="{ error: errors.port }"
            />
            <span v-if="errors.port" class="error-text">{{ errors.port }}</span>
          </div>

          <div class="form-group">
            <label>{{ $t('server.sshPort') }} *</label>
            <input
              v-model.number="formData.ssh_port"
              type="number"
              min="1"
              max="65535"
              :class="{ error: errors.ssh_port }"
            />
            <span v-if="errors.ssh_port" class="error-text">{{ errors.ssh_port }}</span>
          </div>

          <div class="form-group">
            <label>{{ $t('server.username') }} *</label>
            <input
              v-model="formData.ssh_username"
              type="text"
              :class="{ error: errors.ssh_username }"
            />
            <span v-if="errors.ssh_username" class="error-text">{{ errors.ssh_username }}</span>
          </div>

          <div class="form-group">
            <label>{{ $t('server.serverAlias') }} ({{ $t('server.optional') }})</label>
            <input
              v-model="formData.server_alias"
              type="text"
              :placeholder="$t('server.serverAliasPlaceholder')"
            />
          </div>

          <div class="form-group">
            <label>{{ $t('server.status') }} *</label>
            <select v-model="formData.status">
              <option value="active">{{ $t('server.active') }}</option>
              <option value="inactive">{{ $t('server.inactive') }}</option>
            </select>
          </div>
        </div>

        <div class="form-section">
          <h4>{{ $t('server.authInfo') }}</h4>
          <div class="form-group">
            <label>{{ $t('server.authType') }} *</label>
            <select v-model="formData.auth_type">
              <option value="password">{{ $t('server.authPassword') }}</option>
              <option value="key">{{ $t('server.authKey') }}</option>
            </select>
          </div>

          <div v-if="formData.auth_type === 'password'" class="form-group">
            <label>{{ $t('server.sshPassword') }} ({{ $t('server.optional') }})</label>
            <div class="password-input-wrapper">
              <input
                v-model="formData.ssh_password"
                :type="showPassword ? 'text' : 'password'"
                :placeholder="$t('server.optional')"
              />
              <button type="button" class="toggle-password-btn" @click="showPassword = !showPassword">
                <Icon :icon="showPassword ? 'lucide:eye-off' : 'lucide:eye'" />
              </button>
            </div>
          </div>

          <div v-if="formData.auth_type === 'key'" class="form-group">
            <label>{{ $t('server.sshPrivateKey') }} ({{ $t('server.optional') }})</label>
            <textarea
              v-model="formData.ssh_private_key"
              rows="6"
              :placeholder="$t('server.optional')"
            ></textarea>
          </div>
        </div>

        <div class="modal-footer">
          <button type="button" class="btn btn-cancel" @click="$emit('cancel')">
            {{ $t('server.cancel') }}
          </button>
          <button type="submit" class="btn btn-primary" :disabled="submitting">
            {{ submitting ? $t('server.loading') : $t('server.save') }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { createServer, updateServer } from '@/api/server'
import type { Server, CreateServerRequest, UpdateServerRequest } from '@/types'

const props = defineProps<{
  server?: Server | null
}>()

const emit = defineEmits<{
  submit: []
  cancel: []
}>()

const { t } = useI18n()

const isEdit = computed(() => !!props.server)
const submitting = ref(false)
const showPassword = ref(false)

const formData = reactive<CreateServerRequest>({
  ip_address: '',
  port: 10750,
  ssh_username: '',
  ssh_port: 22,
  ssh_password: '',
  ssh_private_key: '',
  auth_type: 'password',
  status: 'active',
  server_alias: ''
})

const errors = reactive<Record<string, string>>({})

const resetForm = () => {
  formData.ip_address = ''
  formData.port = 10750
  formData.ssh_username = ''
  formData.ssh_port = 22
  formData.ssh_password = ''
  formData.ssh_private_key = ''
  formData.auth_type = 'password'
  formData.status = 'active'
  formData.server_alias = ''
  Object.keys(errors).forEach(key => delete errors[key])
}

watch(() => props.server, (server) => {
  if (server) {
    formData.ip_address = server.ip_address
    formData.port = server.port
    formData.ssh_username = server.ssh_username
    formData.ssh_port = server.ssh_port
    formData.auth_type = server.auth_type
    // 后端返回的 status 可能是 'online'/'offline'，需要映射到 'active'/'inactive'
    formData.status = (server.status === 'online' || server.status === 'active') ? 'active' : 'inactive'
    formData.server_alias = server.server_alias || ''
    // 填充密码（如果后端返回了）
    formData.ssh_password = (server as any).ssh_password || ''
    formData.ssh_private_key = (server as any).ssh_private_key || ''
  } else {
    resetForm()
  }
}, { immediate: true })

const validate = () => {
  Object.keys(errors).forEach(key => delete errors[key])

  if (!formData.ip_address) {
    errors.ip_address = '必填项'
  } else if (!/^(\d{1,3}\.){3}\d{1,3}$/.test(formData.ip_address)) {
    errors.ip_address = '请输入有效的IP地址'
  }

  if (!formData.ssh_username) {
    errors.ssh_username = '必填项'
  }

  if (!formData.ssh_port || formData.ssh_port < 1 || formData.ssh_port > 65535) {
    errors.ssh_port = '端口范围 1-65535'
  }

  if (!formData.port || formData.port < 1 || formData.port > 65535) {
    errors.port = '端口范围 1-65535'
  }

  return Object.keys(errors).length === 0
}

const handleSubmit = async () => {
  if (!validate()) return

  submitting.value = true
  try {
    if (isEdit.value && props.server) {
      const updateData: UpdateServerRequest = {
        ip_address: formData.ip_address,
        port: formData.port,
        ssh_username: formData.ssh_username,
        ssh_port: formData.ssh_port,
        auth_type: formData.auth_type,
        status: formData.status,
        server_alias: formData.server_alias || undefined
      }
      if (formData.ssh_password) updateData.ssh_password = formData.ssh_password
      if (formData.ssh_private_key) updateData.ssh_private_key = formData.ssh_private_key
      await updateServer(props.server.id, updateData)
    } else {
      await createServer(formData)
    }
    emit('submit')
  } catch (error) {
    console.error('Failed to save server:', error)
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

.password-input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.password-input-wrapper input {
  padding-right: 40px;
}

.toggle-password-btn {
  position: absolute;
  right: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border-radius: 4px;
  background: transparent;
  color: #64748b;
  border: none;
  cursor: pointer;
  transition: all 0.2s ease;
}

.toggle-password-btn:hover {
  color: #4fc3f7;
  background: #f5f7fa;
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
