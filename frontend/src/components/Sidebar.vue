<script setup>
import { computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '../stores/authStore'
import { useThemeStore } from '../stores/themeStore'
import { 
  LayoutDashboard, 
  Users, 
  GitFork, 
  DollarSign, 
  LogOut, 
  ShieldCheck, 
  Sun, 
  Moon 
} from 'lucide-vue-next'

const authStore = useAuthStore()
const themeStore = useThemeStore()
const router = useRouter()
const route = useRoute()

const navItems = computed(() => [
  {
    name: 'Dashboard',
    path: '/dashboard',
    icon: LayoutDashboard,
    minLevel: 10,
  },
  {
    name: 'Employee Directory',
    path: '/directory',
    icon: Users,
    minLevel: 10,
  },
  {
    name: 'Org Hierarchy',
    path: '/hierarchy',
    icon: GitFork,
    minLevel: 10,
  },
  {
    name: 'Company Payroll',
    path: '/payroll',
    icon: DollarSign,
    minLevel: 80,
    badge: 'HR / Admin Only',
  },
])

const visibleNavItems = computed(() => {
  return navItems.value.filter(item => authStore.accessLevel >= item.minLevel)
})

const handleLogout = async () => {
  await authStore.logout()
  router.push('/login')
}
</script>

<template>
  <aside class="w-64 bg-white dark:bg-[#121215] border-r border-slate-200 dark:border-[#27272a] flex flex-col justify-between h-screen sticky top-0 transition-colors duration-300">
    <!-- Top Branding -->
    <div>
      <div class="p-6 border-b border-slate-200 dark:border-[#27272a] flex items-center justify-between">
        <div class="flex items-center space-x-3">
          <div class="w-10 h-10 rounded-xl bg-gradient-to-tr from-emerald-600 to-emerald-400 flex items-center justify-center shadow-[0_0_15px_rgba(16,185,129,0.4)]">
            <ShieldCheck class="w-6 h-6 text-black font-bold" />
          </div>
          <div>
            <h1 class="font-bold text-lg tracking-wider text-slate-900 dark:text-white flex items-center gap-1.5">
              HRM<span class="text-emerald-500 dark:text-emerald-400 font-extrabold">CORE</span>
            </h1>
            <p class="text-[11px] text-emerald-600 dark:text-emerald-400/80 font-mono tracking-tight uppercase">RBAC Engine</p>
          </div>
        </div>

        <!-- Theme Toggle Button -->
        <button
          @click="themeStore.toggleTheme()"
          title="Toggle Light / Dark Mode"
          class="p-2 rounded-lg bg-slate-100 dark:bg-[#18181c] border border-slate-200 dark:border-[#27272a] text-slate-600 dark:text-emerald-400 hover:border-emerald-500/50 transition-colors cursor-pointer"
        >
          <Sun v-if="themeStore.isDark" class="w-4 h-4 text-emerald-400" />
          <Moon v-else class="w-4 h-4 text-slate-700" />
        </button>
      </div>

      <!-- User Info Badge -->
      <div class="mx-4 my-4 p-3 rounded-lg bg-slate-100 dark:bg-[#18181c] border border-slate-200 dark:border-emerald-500/20 flex items-center space-x-3">
        <div class="w-9 h-9 rounded-full bg-emerald-100 dark:bg-emerald-950 border border-emerald-400/40 text-emerald-600 dark:text-emerald-400 flex items-center justify-center font-bold text-sm">
          {{ authStore.user?.first_name?.[0] || 'U' }}{{ authStore.user?.last_name?.[0] || '' }}
        </div>
        <div class="flex-1 min-w-0">
          <p class="text-xs font-semibold text-slate-900 dark:text-white truncate">{{ authStore.fullName }}</p>
          <div class="flex items-center space-x-1.5 mt-0.5">
            <span class="inline-block w-1.5 h-1.5 rounded-full bg-emerald-500 dark:bg-emerald-400 animate-pulse"></span>
            <span class="text-[11px] text-emerald-600 dark:text-emerald-400 font-mono font-medium">
              LVL {{ authStore.accessLevel }} ({{ authStore.userRole }})
            </span>
          </div>
        </div>
      </div>

      <!-- Navigation Links -->
      <nav class="px-3 space-y-1">
        <p class="px-3 text-[10px] font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider mb-2">Navigation</p>
        <router-link
          v-for="item in visibleNavItems"
          :key="item.path"
          :to="item.path"
          :class="[
            route.path === item.path
              ? 'bg-emerald-50 dark:bg-emerald-950/50 text-emerald-600 dark:text-emerald-400 border-l-2 border-emerald-500 dark:border-emerald-400 font-medium'
              : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-slate-200 hover:bg-slate-100 dark:hover:bg-[#1c1c21]',
            'group flex items-center px-3 py-2.5 text-sm rounded-r-lg transition-all duration-200'
          ]"
        >
          <component
            :is="item.icon"
            :class="[
              route.path === item.path ? 'text-emerald-600 dark:text-emerald-400' : 'text-slate-400 group-hover:text-emerald-500 dark:group-hover:text-emerald-400',
              'mr-3 h-4 w-4 shrink-0 transition-colors'
            ]"
          />
          <span class="flex-1">{{ item.name }}</span>
          <span 
            v-if="item.badge"
            class="text-[9px] px-1.5 py-0.5 rounded bg-emerald-100 dark:bg-emerald-950 border border-emerald-400/30 text-emerald-700 dark:text-emerald-400 font-mono"
          >
            {{ item.badge }}
          </span>
        </router-link>
      </nav>
    </div>

    <!-- Bottom Logout -->
    <div class="p-4 border-t border-slate-200 dark:border-[#27272a]">
      <button
        @click="handleLogout"
        class="w-full flex items-center justify-center space-x-2 px-4 py-2.5 text-sm rounded-lg bg-red-50 dark:bg-red-950/30 text-red-600 dark:text-red-400 hover:bg-red-100 dark:hover:bg-red-950/60 border border-red-200 dark:border-red-900/40 transition-colors duration-200 font-medium cursor-pointer"
      >
        <LogOut class="w-4 h-4" />
        <span>Sign Out</span>
      </button>
    </div>
  </aside>
</template>
