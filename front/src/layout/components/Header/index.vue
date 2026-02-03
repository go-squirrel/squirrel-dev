<template>
  <header class="header">
    <div class="header-left">
      <button class="toggle-btn" @click="$emit('toggleSidebar')">
        <Icon icon="lucide:menu" />
      </button>
      <h1 class="page-title">{{ pageTitle }}</h1>
    </div>
    <div class="header-right">
      <div class="user-info">
        <Icon icon="lucide:user" class="user-icon" />
        <span class="user-name">{{ userName }}</span>
      </div>
      <button class="logout-btn" @click="$emit('logout')">
        <Icon icon="lucide:log-out" />
      </button>
    </div>
  </header>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { Icon } from '@iconify/vue'

defineProps<{
  userName: string
}>()

defineEmits<{
  toggleSidebar: []
  logout: []
}>()

const route = useRoute()

const pageTitle = computed(() => {
  const titles: Record<string, string> = {
    '/': '概览',
    '/servers': '服务器',
    '/applications': '应用',
    '/monitor': '监控',
    '/scripts': '脚本',
    '/configs': '配置',
    '/deployments': '部署',
    '/appstore': '应用商店'
  }
  
  for (const [path, title] of Object.entries(titles)) {
    if (route.path === path || route.path.startsWith(path + '/')) {
      return title
    }
  }
  return 'Squirrel'
})
</script>

<style scoped>
.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 10px;
  background: #ffffff;
  border-bottom: 1px solid #e8ecf1;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.toggle-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 6px;
  background: #f5f7fa;
  color: #64748b;
  border: none;
  cursor: pointer;
  transition: all 0.3s ease;
}

.toggle-btn:hover {
  background: #e0f2fe;
  color: #4fc3f7;
}

.page-title {
  font-size: 18px;
  font-weight: 600;
  color: #1e3a5f;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: #f5f7fa;
  border-radius: 12px;
}

.user-icon {
  width: 16px;
  height: 16px;
  color: #4fc3f7;
}

.user-name {
  font-size: 13px;
  font-weight: 500;
  color: #1e3a5f;
}

.logout-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 6px;
  background: #f5f7fa;
  color: #64748b;
  transition: all 0.3s ease;
  border: none;
  cursor: pointer;
}

.logout-btn:hover {
  background: #fee2e2;
  color: #ef4444;
  transform: rotate(-5deg);
}
</style>
