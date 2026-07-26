<script setup>
import Sidebar from './Sidebar.vue'
import { useAuthStore } from '../stores/authStore'
import { Terminal } from 'lucide-vue-next'

const authStore = useAuthStore()
</script>

<template>
  <div class="flex min-h-screen bg-slate-50 dark:bg-[#09090b] transition-colors duration-300">
    <!-- Sidebar -->
    <Sidebar />

    <!-- Main Content Area -->
    <div class="flex-1 flex flex-col min-w-0 overflow-x-hidden">
      <!-- Top Navbar -->
      <header class="h-16 bg-white/80 dark:bg-[#121215]/80 backdrop-blur-md border-b border-slate-200 dark:border-[#27272a] px-8 flex items-center justify-between sticky top-0 z-10 transition-colors duration-300">
        <div class="flex items-center space-x-3 text-xs text-slate-500 dark:text-slate-400 font-mono">
          <Terminal class="w-4 h-4 text-emerald-500 dark:text-emerald-400" />
          <span>RBAC LEVEL: <strong class="text-emerald-600 dark:text-emerald-400">{{ authStore.accessLevel }}</strong></span>
          <span class="text-slate-300 dark:text-slate-600">|</span>
          <span>ROLE: <strong class="text-slate-900 dark:text-white">{{ authStore.userRole }}</strong></span>
        </div>

        <div class="flex items-center space-x-4">
          <div class="px-3 py-1 rounded-full bg-emerald-100 dark:bg-emerald-950/60 border border-emerald-400/40 text-emerald-700 dark:text-emerald-400 text-xs font-mono flex items-center space-x-2">
            <span class="w-2 h-2 rounded-full bg-emerald-500 dark:bg-emerald-400 animate-ping"></span>
            <span>API SESSION ACTIVE</span>
          </div>
        </div>
      </header>

      <!-- Page Content -->
      <main class="flex-1 p-8">
        <slot />
      </main>
    </div>
  </div>
</template>
