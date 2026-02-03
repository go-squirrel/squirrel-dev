import { createI18n } from 'vue-i18n'
import zhCN from './zh-CN'

const messages = {
  'zh-CN': zhCN
}

export const i18n = createI18n({
  legacy: false,
  locale: 'zh-CN',
  fallbackLocale: 'zh-CN',
  messages
})
