<script setup>
import { onMounted } from 'vue'
import { useAuthStore } from './stores/authStore'
import { useThemeStore } from './stores/themeStore'

const authStore = useAuthStore()
const themeStore = useThemeStore()

onMounted(() => {
  themeStore.initTheme()
  if (authStore.isAuthenticated) {
    authStore.fetchCurrentUser()
  }
})
</script>

<template>
  <div class="min-h-screen bg-slate-50 dark:bg-[#09090b] text-slate-900 dark:text-slate-100 selection:bg-emerald-500 selection:text-black font-sans transition-colors duration-300">
    <router-view />
  </div>
</template>
