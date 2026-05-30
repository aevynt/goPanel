import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/Login.vue'),
    },
    {
      path: '/editor',
      name: 'editor',
      component: () => import('@/views/Editor.vue'),
    },
    {
      path: '/',
      component: () => import('@/views/Layout.vue'),
      redirect: '/dashboard',
      children: [
        {
          path: 'dashboard',
          name: 'dashboard',
          component: () => import('@/views/Dashboard.vue'),
        },
        {
          path: 'services',
          name: 'services',
          component: () => import('@/views/Services.vue'),
        },
        {
          path: 'binaries',
          name: 'binaries',
          component: () => import('@/views/Binaries.vue'),
        },
        {
          path: 'files/:path(.*)?',
          name: 'files',
          component: () => import('@/views/FileManager.vue'),
          props: true,
        },
        {
          path: 'ports',
          name: 'ports',
          component: () => import('@/views/Ports.vue'),
        },
        {
          path: 'sites',
          name: 'sites',
          component: () => import('@/views/Sites.vue'),
        },
        {
          path: 'users',
          name: 'users',
          component: () => import('@/views/Users.vue'),
        },
        {
          path: 'public',
          name: 'public',
          component: () => import('@/views/Public.vue'),
        },
        {
          path: 'settings',
          name: 'settings',
          component: () => import('@/views/Settings.vue'),
        },
      ],
    },
  ],
})

router.beforeEach((to, _from, next) => {
  const normalized = to.path.replace(/\/+/g, '/')
  if (to.path !== normalized) {
    next({ path: normalized, query: to.query, hash: to.hash })
    return
  }
  const auth = useAuthStore()
  if (to.name !== 'login' && !auth.isAuthenticated) {
    next('/login')
  } else if (to.name === 'login' && auth.isAuthenticated) {
    next('/dashboard')
  } else {
    next()
  }
})

export default router
