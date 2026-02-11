import { addCollection, Icon as OfflineIcon } from '@iconify/vue/offline'
import lucide from '@iconify/json/json/lucide.json'
import type { App } from 'vue'

// 预加载 lucide 图标集合
export const setupIcons = (app: App) => {
  addCollection(lucide)
  // 全局注册离线 Icon 组件
  app.component('Icon', OfflineIcon)
}
