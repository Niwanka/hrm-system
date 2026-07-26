import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/authStore'

import LoginView from '../views/LoginView.vue'
import DashboardView from '../views/DashboardView.vue'
import DirectoryView from '../views/DirectoryView.vue'
import HierarchyView from '../views/HierarchyView.vue'
import PayrollView from '../views/PayrollView.vue'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: LoginView,
    meta: { guest: true },
  },
  {
    path: '/',
    redirect: '/dashboard',
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: DashboardView,
    meta: { requiresAuth: true, minAccessLevel: 10 },
  },
  {
    path: '/directory',
    name: 'Directory',
    component: DirectoryView,
    meta: { requiresAuth: true, minAccessLevel: 10 },
  },
  {
    path: '/hierarchy',
    name: 'Hierarchy',
    component: HierarchyView,
    meta: { requiresAuth: true, minAccessLevel: 10 },
  },
  {
    path: '/payroll',
    name: 'Payroll',
    component: PayrollView,
    meta: { requiresAuth: true, minAccessLevel: 80 }, // HR/Admin only
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/dashboard',
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

// Navigation Guard
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    return next({ name: 'Login' })
  }

  if (to.meta.guest && authStore.isAuthenticated) {
    return next({ name: 'Dashboard' })
  }

  if (to.meta.minAccessLevel && authStore.accessLevel < to.meta.minAccessLevel) {
    alert(`Access Denied: Requires minimum access level ${to.meta.minAccessLevel}. Your level is ${authStore.accessLevel}.`)
    return next({ name: 'Dashboard' })
  }

  next()
})

export default router
