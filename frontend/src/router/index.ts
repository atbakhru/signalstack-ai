import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/authStore'

const routes = [
  {
    path: '/',
    redirect: '/login',
  },
  {
    path: '/login',
    component: () => import('../views/LoginView.vue'),
  },
  {
    path: '/register',
    component: () => import('../views/RegisterView.vue'),
  },
  {
    path: '/dashboard',
    component: () => import('../views/DashboardView.vue'),
  },
  {
    path: '/documents',
    component: () => import('../views/DocumentsView.vue'),
  },
  {
    path: '/ingestion',
    component: () => import('../views/IngestionView.vue'),
  },
  {
    path: '/chat',
    component: () => import('../views/ChatView.vue'),
  },
  {
    path: '/evaluation',
    component: () => import('../views/EvaluationView.vue'),
  },
  {
    path: '/metrics',
    component: () => import('../views/MetricsView.vue'),
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach(async (to) => {
  const authStore = useAuthStore()
  await authStore.initialize()
  if (to.path !== '/login' && to.path !== '/register' && !authStore.isAuthenticated) {
    return '/login'
  }
  if ((to.path === '/login' || to.path === '/register') && authStore.isAuthenticated) {
    return '/dashboard'
  }
  return true
})

export default router
