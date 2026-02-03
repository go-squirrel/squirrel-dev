// 布局状态管理
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useLayoutStore = defineStore('layout', () => {
  // 状态
  const sidebarCollapsed = ref(false)
  const currentLayout = ref<'default' | 'compact' | 'full'>('default')
  
  // 计算属性
  const isSidebarCollapsed = computed(() => sidebarCollapsed.value)
  const layoutMode = computed(() => currentLayout.value)
  
  // 方法
  function toggleSidebar() {
    sidebarCollapsed.value = !sidebarCollapsed.value
  }
  
  function setSidebarCollapsed(collapsed: boolean) {
    sidebarCollapsed.value = collapsed
  }
  
  function setLayoutMode(mode: 'default' | 'compact' | 'full') {
    currentLayout.value = mode
  }
  
  return {
    sidebarCollapsed,
    currentLayout,
    isSidebarCollapsed,
    layoutMode,
    toggleSidebar,
    setSidebarCollapsed,
    setLayoutMode
  }
})
