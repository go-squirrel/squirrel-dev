<template>
  <div class="login-page">
    <div class="login-container">
      <div class="login-left">
        <LoginBrand />
      </div>
      <div class="login-right">
        <div class="login-form-wrapper">
          <h2 class="login-title">{{ $t('login.welcomeBack') }}</h2>
          <p class="login-subtitle">{{ $t('login.loginAccount') }}</p>
          <LoginForm :loading="loading" @submit="handleLogin" />
          <p class="login-tip">{{ $t('login.tip') }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useUserStore } from '@/store'
import { login } from '@/api'
import LoginForm from './components/LoginForm.vue'
import LoginBrand from './components/LoginBrand.vue'
import { useLoading } from '@/composables/useLoading'

const router = useRouter()
const userStore = useUserStore()
const { loading, withLoading } = useLoading()

const handleLogin = async (formData: { username: string; password: string; remember: boolean }) => {
  await withLoading(async () => {
    const data = await login({
      username: formData.username,
      password: formData.password
    })

    if (data.token) {
      userStore.setToken(data.token)
      if (data.user) {
        userStore.setUser({
          id: data.user.id,
          username: data.user.username,
          role: 'admin'
        })
      }
      router.push('/')
    }
  })
}
</script>

<style scoped>
.login-page {
  width: 100%;
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #1e3a5f 0%, #0f1f33 100%);
  padding: 20px;
}

.login-container {
  display: flex;
  width: 100%;
  max-width: 1000px;
  height: 600px;
  background: #ffffff;
  border-radius: 24px;
  overflow: hidden;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
}

.login-left {
  flex: 1;
  background: linear-gradient(135deg, #1e3a5f 0%, #0f1f33 100%);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px;
  position: relative;
  overflow: hidden;
}

.login-left::before {
  content: '';
  position: absolute;
  width: 300px;
  height: 300px;
  background: rgba(79, 195, 247, 0.1);
  border-radius: 50%;
  top: -100px;
  right: -100px;
}

.login-left::after {
  content: '';
  position: absolute;
  width: 200px;
  height: 200px;
  background: rgba(79, 195, 247, 0.05);
  border-radius: 50%;
  bottom: -50px;
  left: -50px;
}

.login-right {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px;
}

.login-form-wrapper {
  width: 100%;
  max-width: 360px;
}

.login-title {
  font-size: 28px;
  font-weight: 700;
  color: #1e3a5f;
  margin-bottom: 8px;
}

.login-subtitle {
  font-size: 14px;
  color: #64748b;
  margin-bottom: 32px;
}

.login-tip {
  text-align: center;
  font-size: 12px;
  color: #94a3b8;
  margin-top: 20px;
}
</style>
