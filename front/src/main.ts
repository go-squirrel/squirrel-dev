import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { pinia } from './store'
import i18n from './lang'
import { setupIcons } from './utils/icons'

const app = createApp(App)

app.use(pinia)
app.use(router)
app.use(i18n)

// 预加载图标数据并全局注册 Icon 组件（离线模式）
setupIcons(app)

app.mount('#app')
