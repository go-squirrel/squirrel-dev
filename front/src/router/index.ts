import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/store'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login/index.vue'),
    meta: { requiresAuth: false, layout: 'full' }
  },
  {
    path: '/',
    name: 'Overview',
    component: () => import('@/views/Overview/index.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/servers',
    name: 'Servers',
    component: () => import('@/views/Server/index.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/scripts',
    name: 'Scripts',
    component: () => import('@/views/Scripts/index.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/configs',
    name: 'Configs',
    component: () => import('@/views/Configs/index.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/servers/:id/terminal',
    name: 'ServerTerminal',
    component: () => import('@/views/Server/Terminal.vue'),
    meta: { requiresAuth: true, layout: 'full' }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, _from, next) => {
  const userStore = useUserStore()
  const requiresAuth = to.meta.requiresAuth !== false

  if (requiresAuth && !userStore.isLoggedIn) {
    next('/login')
  } else {
    next()
  }
})

export default router
