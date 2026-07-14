import { createRouter, createWebHistory } from 'vue-router/auto'
import { routes } from 'vue-router/auto-routes'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to: { meta: { requiresAuth: boolean | undefined } }) => {
  const authStore = useAuthStore()
  const requiresAuth = to.meta.requiresAuth as boolean | undefined

  if (requiresAuth && !authStore.isAuthenticated) {
    return { path: '/Login' }
  }
})

export default router