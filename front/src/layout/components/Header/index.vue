<template>
  <header class="header">
    <div class="header-left">
      <button v-if="!isFullLayout" class="toggle-btn" @click="$emit('toggleSidebar')">
        <Icon icon="lucide:menu" />
      </button>
      <h1 class="page-title">{{ pageTitle }}</h1>
    </div>
    <div class="header-right">
      <!-- 语言切换 -->
      <div class="lang-switcher">
        <button class="lang-btn" @click="showLangMenu = !showLangMenu">
          <Icon icon="lucide:globe" />
          <span>{{ currentLocaleName }}</span>
          <Icon icon="lucide:chevron-down" class="chevron" :class="{ open: showLangMenu }" />
        </button>
        <div v-if="showLangMenu" class="lang-menu">
          <button
            v-for="loc in availableLocales"
            :key="loc.code"
            class="lang-option"
            :class="{ active: currentLocale === loc.code }"
            @click="switchLocale(loc.code)"
          >
            {{ loc.name }}
          </button>
        </div>
      </div>
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
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import { Icon } from '@iconify/vue'
import { useI18n } from 'vue-i18n'
import { availableLocales } from '@/lang'

defineProps<{
  userName: string
}>()

defineEmits<{
  toggleSidebar: []
  logout: []
}>()

const route = useRoute()
const { locale, t } = useI18n()

const showLangMenu = ref(false)

const isFullLayout = computed(() => route.meta.layout === 'full')
const currentLocale = computed(() => locale.value)
const currentLocaleName = computed(() => {
  const loc = availableLocales.find(l => l.code === locale.value)
  return loc?.name || locale.value
})

const pageTitle = computed(() => {
  const titles: Record<string, string> = {
    '/': t('layout.overview'),
    '/servers': t('server.title'),
    '/servers/:id/terminal': t('server.terminal'),
    '/applications': t('layout.applications'),
    '/monitor': t('layout.monitor'),
    '/scripts': t('layout.scripts'),
    '/configs': t('layout.configs'),
    '/deployments': t('layout.deployments'),
    '/appstore': t('layout.appstore')
  }
  
  for (const [path, title] of Object.entries(titles)) {
    if (route.path === path || route.path.startsWith(path + '/')) {
      return title
    }
  }
  return 'Squirrel'
})

const switchLocale = (code: string) => {
  locale.value = code
  localStorage.setItem('locale', code)
  showLangMenu.value = false
}
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

/* 语言切换器 */
.lang-switcher {
  position: relative;
}

.lang-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: #f5f7fa;
  border-radius: 6px;
  color: #64748b;
  font-size: 13px;
  border: none;
  cursor: pointer;
  transition: all 0.3s ease;
}

.lang-btn:hover {
  background: #e0f2fe;
  color: #4fc3f7;
}

.lang-btn .chevron {
  width: 14px;
  height: 14px;
  transition: transform 0.3s ease;
}

.lang-btn .chevron.open {
  transform: rotate(180deg);
}

.lang-menu {
  position: absolute;
  top: 100%;
  right: 0;
  margin-top: 4px;
  background: #ffffff;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  padding: 4px;
  min-width: 120px;
  z-index: 100;
}

.lang-option {
  display: block;
  width: 100%;
  padding: 8px 12px;
  text-align: left;
  font-size: 13px;
  color: #64748b;
  background: none;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.lang-option:hover {
  background: #f5f7fa;
  color: #1e3a5f;
}

.lang-option.active {
  background: #e0f2fe;
  color: #0284c7;
  font-weight: 500;
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
