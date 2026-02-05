<template>
  <div class="main-layout">
    <Sidebar :is-collapsed="isSidebarCollapsed" />
    <div class="main-wrapper">
      <Header 
        :user-name="userName" 
        @toggle-sidebar="toggleSidebar"
        @logout="handleLogout"
      />
      <MainContent>
        <slot />
      </MainContent>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore, useLayoutStore } from '@/store'
import Sidebar from './components/Sidebar/index.vue'
import Header from './components/Header/index.vue'
import MainContent from './components/MainContent/index.vue'

const router = useRouter()
const userStore = useUserStore()
const layoutStore = useLayoutStore()

const isSidebarCollapsed = computed(() => layoutStore.isSidebarCollapsed)
const userName = computed(() => userStore.currentUser?.username || 'Admin')

const toggleSidebar = () => {
  layoutStore.toggleSidebar()
}

const handleLogout = () => {
  userStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.main-layout {
  display: flex;
  width: 100%;
  height: 100vh;
  background: #f5f7fa;
}

.main-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
</style>
