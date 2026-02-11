<template>
  <form class="login-form" @submit.prevent="handleSubmit">
    <div class="form-group">
      <label class="form-label">{{ $t('login.username') }}</label>
      <div class="input-wrapper">
        <Icon icon="lucide:user" class="input-icon" />
        <input
          v-model="formData.username"
          type="text"
          class="form-input"
          :placeholder="$t('login.usernamePlaceholder')"
          required
        />
      </div>
    </div>
    <div class="form-group">
      <label class="form-label">{{ $t('login.password') }}</label>
      <div class="input-wrapper">
        <Icon icon="lucide:lock" class="input-icon" />
        <input
          v-model="formData.password"
          :type="showPassword ? 'text' : 'password'"
          class="form-input"
          :placeholder="$t('login.passwordPlaceholder')"
          required
        />
        <button
          type="button"
          class="toggle-password"
          @click="showPassword = !showPassword"
        >
          <Icon :icon="showPassword ? 'lucide:eye-off' : 'lucide:eye'" />
        </button>
      </div>
    </div>
    <div class="form-options">
      <label class="checkbox-label">
        <input v-model="formData.remember" type="checkbox" />
        <span>{{ $t('login.rememberMe') }}</span>
      </label>
    </div>
    <Button type="primary" size="large" block :loading="loading" @click="handleSubmit">
      {{ $t('login.login') }}
    </Button>
  </form>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import Button from '@/components/Button/index.vue'

defineProps<{
  loading: boolean
}>()

const emit = defineEmits<{
  submit: [data: { username: string; password: string; remember: boolean }]
}>()

const formData = ref({
  username: '',
  password: '',
  remember: false
})

const showPassword = ref(false)

const handleSubmit = () => {
  emit('submit', { ...formData.value })
}
</script>

<style scoped>
.login-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-label {
  font-size: 14px;
  font-weight: 500;
  color: #1e3a5f;
}

.input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.input-icon {
  position: absolute;
  left: 14px;
  width: 20px;
  height: 20px;
  color: #94a3b8;
  pointer-events: none;
}

.form-input {
  width: 100%;
  padding: 14px 14px 14px 44px;
  font-size: 14px;
  border: 2px solid #e2e8f0;
  border-radius: 12px;
  background: #f8fafc;
  color: #1e3a5f;
  transition: all 0.3s ease;
}

.form-input:focus {
  outline: none;
  border-color: #4fc3f7;
  background: #ffffff;
  box-shadow: 0 0 0 4px rgba(79, 195, 247, 0.1);
}

.toggle-password {
  position: absolute;
  right: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  background: none;
  border: none;
  color: #94a3b8;
  cursor: pointer;
  transition: color 0.3s ease;
}

.toggle-password:hover {
  color: #4fc3f7;
}

.form-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #64748b;
  cursor: pointer;
}

.checkbox-label input {
  width: 16px;
  height: 16px;
  accent-color: #4fc3f7;
}
</style>
