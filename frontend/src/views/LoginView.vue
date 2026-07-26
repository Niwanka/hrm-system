<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/authStore'
import { useThemeStore } from '../stores/themeStore'
import { ShieldCheck, Lock, Mail, ArrowRight, Sparkles, Sun, Moon } from 'lucide-vue-next'

const email = ref('admin@company.com')
const password = ref('password123')
const authStore = useAuthStore()
const themeStore = useThemeStore()
const router = useRouter()

const handleLogin = async () => {
  const result = await authStore.login({
    email: email.value,
    password: password.value,
  })

  if (result.success) {
    router.push('/dashboard')
  }
}

const setDemoAccount = (demoEmail) => {
  email.value = demoEmail
  password.value = 'password123'
}
</script>

<template>
  <div class="min-h-screen bg-slate-50 dark:bg-[#09090b] flex items-center justify-center p-6 relative overflow-hidden transition-colors duration-300">
    <!-- Ambient Glow -->
    <div class="absolute -top-40 -left-40 w-96 h-96 bg-emerald-500/10 dark:bg-emerald-600/15 rounded-full blur-3xl pointer-events-none"></div>
    <div class="absolute -bottom-40 -right-40 w-96 h-96 bg-emerald-500/10 rounded-full blur-3xl pointer-events-none"></div>

    <!-- Top Right Theme Switcher -->
    <button
      @click="themeStore.toggleTheme()"
      title="Toggle Light / Dark Mode"
      class="absolute top-6 right-6 p-2.5 rounded-xl bg-white dark:bg-[#121215] border border-slate-200 dark:border-[#27272a] shadow-md text-slate-700 dark:text-emerald-400 hover:border-emerald-500/50 transition-colors cursor-pointer"
    >
      <Sun v-if="themeStore.isDark" class="w-5 h-5 text-emerald-400" />
      <Moon v-else class="w-5 h-5 text-slate-700" />
    </button>

    <div class="w-full max-w-md">
      <!-- Logo Header -->
      <div class="text-center mb-8">
        <div class="w-16 h-16 rounded-2xl bg-gradient-to-tr from-emerald-600 to-emerald-400 mx-auto flex items-center justify-center shadow-[0_0_30px_rgba(16,185,129,0.4)] mb-4">
          <ShieldCheck class="w-9 h-9 text-black font-extrabold" />
        </div>
        <h1 class="text-3xl font-extrabold tracking-tight text-slate-900 dark:text-white">
          HRM<span class="text-emerald-500 dark:text-emerald-400">CORE</span>
        </h1>
        <p class="text-sm text-slate-500 dark:text-slate-400 mt-1">Enterprise Core Authentication & RBAC Directory</p>
      </div>

      <!-- Login Card -->
      <div class="glass-card rounded-2xl p-8 border border-slate-200 dark:border-[#27272a] shadow-2xl relative bg-white dark:bg-[#121215]">
        <div v-if="authStore.error" class="mb-6 p-4 rounded-xl bg-red-50 dark:bg-red-950/40 border border-red-200 dark:border-red-900/60 text-red-600 dark:text-red-300 text-sm flex items-start space-x-2">
          <span>⚠️ {{ authStore.error }}</span>
        </div>

        <form @submit.prevent="handleLogin" class="space-y-5">
          <div>
            <label class="block text-xs font-semibold text-slate-600 dark:text-slate-300 uppercase tracking-wider mb-2">
              Work Email Address
            </label>
            <div class="relative">
              <Mail class="w-5 h-5 text-slate-400 absolute left-3.5 top-3" />
              <input
                v-model="email"
                type="email"
                required
                placeholder="name@company.com"
                class="w-full pl-11 pr-4 py-2.5 rounded-xl bg-slate-50 dark:bg-[#18181c] border border-slate-200 dark:border-[#27272a] text-slate-900 dark:text-white text-sm focus:outline-none focus:border-emerald-500 dark:focus:border-emerald-400 focus:ring-1 focus:ring-emerald-500"
              />
            </div>
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-600 dark:text-slate-300 uppercase tracking-wider mb-2">
              Password
            </label>
            <div class="relative">
              <Lock class="w-5 h-5 text-slate-400 absolute left-3.5 top-3" />
              <input
                v-model="password"
                type="password"
                required
                placeholder="••••••••"
                class="w-full pl-11 pr-4 py-2.5 rounded-xl bg-slate-50 dark:bg-[#18181c] border border-slate-200 dark:border-[#27272a] text-slate-900 dark:text-white text-sm focus:outline-none focus:border-emerald-500 dark:focus:border-emerald-400 focus:ring-1 focus:ring-emerald-500"
              />
            </div>
          </div>

          <button
            type="submit"
            :disabled="authStore.isLoading"
            class="w-full py-3 px-4 rounded-xl bg-gradient-to-r from-emerald-500 to-emerald-600 hover:from-emerald-400 hover:to-emerald-500 text-black font-bold text-sm flex items-center justify-center space-x-2 shadow-[0_0_20px_rgba(16,185,129,0.35)] transition-all duration-200 disabled:opacity-50 cursor-pointer"
          >
            <span>{{ authStore.isLoading ? 'Authenticating...' : 'Sign In to Portal' }}</span>
            <ArrowRight v-if="!authStore.isLoading" class="w-4 h-4" />
          </button>
        </form>

        <!-- Quick Demo Switcher -->
        <div class="mt-8 pt-6 border-t border-slate-200 dark:border-[#27272a]">
          <div class="flex items-center space-x-1 text-xs text-slate-500 dark:text-slate-400 mb-3 font-semibold">
            <Sparkles class="w-3.5 h-3.5 text-emerald-500 dark:text-emerald-400" />
            <span>Quick Select Demo Accounts:</span>
          </div>
          <div class="grid grid-cols-2 gap-2 text-xs font-mono">
            <button
              @click="setDemoAccount('admin@company.com')"
              type="button"
              class="p-2 rounded-lg bg-slate-50 dark:bg-[#18181c] hover:bg-emerald-50 dark:hover:bg-emerald-950/60 border border-slate-200 dark:border-[#27272a] hover:border-emerald-500/50 text-left transition-colors cursor-pointer"
            >
              <div class="text-emerald-600 dark:text-emerald-400 font-bold">Admin (LVL 100)</div>
              <div class="text-[10px] text-slate-500 dark:text-slate-400 truncate">admin@company.com</div>
            </button>
            <button
              @click="setDemoAccount('hr@company.com')"
              type="button"
              class="p-2 rounded-lg bg-slate-50 dark:bg-[#18181c] hover:bg-emerald-50 dark:hover:bg-emerald-950/60 border border-slate-200 dark:border-[#27272a] hover:border-emerald-500/50 text-left transition-colors cursor-pointer"
            >
              <div class="text-emerald-600 dark:text-emerald-400 font-bold">HR (LVL 80)</div>
              <div class="text-[10px] text-slate-500 dark:text-slate-400 truncate">hr@company.com</div>
            </button>
            <button
              @click="setDemoAccount('manager@company.com')"
              type="button"
              class="p-2 rounded-lg bg-slate-50 dark:bg-[#18181c] hover:bg-emerald-50 dark:hover:bg-emerald-950/60 border border-slate-200 dark:border-[#27272a] hover:border-emerald-500/50 text-left transition-colors cursor-pointer"
            >
              <div class="text-emerald-600 dark:text-emerald-400 font-bold">Manager (LVL 50)</div>
              <div class="text-[10px] text-slate-500 dark:text-slate-400 truncate">manager@company.com</div>
            </button>
            <button
              @click="setDemoAccount('employee@company.com')"
              type="button"
              class="p-2 rounded-lg bg-slate-50 dark:bg-[#18181c] hover:bg-emerald-50 dark:hover:bg-emerald-950/60 border border-slate-200 dark:border-[#27272a] hover:border-emerald-500/50 text-left transition-colors cursor-pointer"
            >
              <div class="text-emerald-600 dark:text-emerald-400 font-bold">Employee (LVL 10)</div>
              <div class="text-[10px] text-slate-500 dark:text-slate-400 truncate">employee@company.com</div>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
