<template>
  <aside class="sidebar" :class="{ collapsed: isCollapsed }">
    <div class="logo">
      <Icon icon="lucide:squirrel" class="logo-icon" />
      <span v-if="!isCollapsed" class="logo-text">Squirrel</span>
    </div>
    <nav class="nav">
      <router-link
        v-for="item in navItems"
        :key="item.path"
        :to="item.path"
        class="nav-item"
        :class="{ active: isActive(item.path) }"
        :title="isCollapsed ? item.label : ''"
      >
        <Icon :icon="item.icon" class="nav-icon" />
        <span v-if="!isCollapsed">{{ item.label }}</span>
      </router-link>
    </nav>
  </aside>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router'
import { Icon } from '@iconify/vue'
import type { NavItem } from '../../types'

defineProps<{
  isCollapsed: boolean
}>()

const route = useRoute()

const navItems: NavItem[] = [
  { path: '/', label: '概览', icon: 'lucide:layout-dashboard' },
  { path: '/servers', label: '服务器', icon: 'lucide:server' },
  { path: '/applications', label: '应用', icon: 'lucide:box' },
  { path: '/monitor', label: '监控', icon: 'lucide:activity' },
  { path: '/scripts', label: '脚本', icon: 'lucide:terminal' },
  { path: '/configs', label: '配置', icon: 'lucide:settings' },
  { path: '/deployments', label: '部署', icon: 'lucide:rocket' },
  { path: '/appstore', label: '应用商店', icon: 'lucide:shopping-bag' }
]

const isActive = (path: string) => {
  return route.path === path || route.path.startsWith(path + '/')
}
</script>

<style scoped>
.sidebar {
  width: 200px;
  background: linear-gradient(180deg, #1e3a5f 0%, #0f1f33 100%);
  display: flex;
  flex-direction: column;
  padding: 16px 0;
  box-shadow: 4px 0 12px rgba(0, 0, 0, 0.1);
  transition: width 0.3s ease;
}

.sidebar.collapsed {
  width: 64px;
}

.logo {
  display: flex;
  align-items: center;
  padding: 0 16px;
  margin-bottom: 20px;
}

.logo-icon {
  width: 24px;
  height: 24px;
  color: #4fc3f7;
  flex-shrink: 0;
}

.logo-text {
  margin-left: 8px;
  font-size: 16px;
  font-weight: 700;
  color: #ffffff;
  letter-spacing: 0.5px;
  white-space: nowrap;
}

.nav {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 0 8px;
}

.nav-item {
  display: flex;
  align-items: center;
  padding: 10px 12px;
  margin-bottom: 2px;
  border-radius: 4px;
  color: rgba(255, 255, 255, 0.6);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  font-size: 13px;
  font-weight: 500;
}

.nav-item:hover {
  background: rgba(79, 195, 247, 0.15);
  color: #4fc3f7;
  transform: translateX(2px);
}

.nav-item.active {
  background: linear-gradient(135deg, #4fc3f7 0%, #29b6f6 100%);
  color: #ffffff;
  box-shadow: 0 2px 8px rgba(79, 195, 247, 0.4);
}

.nav-icon {
  width: 16px;
  height: 16px;
  margin-right: 8px;
  flex-shrink: 0;
}

.sidebar.collapsed .nav-icon {
  margin-right: 0;
}
</style>
