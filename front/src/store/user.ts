// 用户状态管理
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User } from '@/types'

export const useUserStore = defineStore('user', () => {
  // 状态
  const user = ref<User | null>(null)
  const token = ref<string>(localStorage.getItem('token') || '')
  
  // 计算属性
  const isLoggedIn = computed(() => !!token.value)
  const currentUser = computed(() => user.value)
  
  // 方法
  function setUser(userData: User) {
    user.value = userData
  }
  
  function setToken(newToken: string) {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }
  
  function clearUser() {
    user.value = null
    token.value = ''
    localStorage.removeItem('token')
  }
  
  function logout() {
    clearUser()
  }
  
  return {
    user,
    token,
    isLoggedIn,
    currentUser,
    setUser,
    setToken,
    clearUser,
    logout
  }
})
