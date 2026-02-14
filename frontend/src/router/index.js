import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue')
  },
  {
    path: '/',
    name: 'Dashboard',
    component: () => import('../views/Dashboard.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/contracts',
    name: 'Contracts',
    component: () => import('../views/Contracts.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/contracts/:id',
    name: 'ContractDetail',
    component: () => import('../views/ContractDetail.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/obligations',
    name: 'Obligations',
    component: () => import('../views/Obligations.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/financial',
    name: 'Financial',
    component: () => import('../views/Financial.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/audit',
    name: 'Audit',
    component: () => import('../views/Audit.vue'),
    meta: { requiresAuth: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next('/login')
  } else if (to.path === '/login' && authStore.isAuthenticated) {
    next('/')
  } else {
    next()
  }
})

export default router
