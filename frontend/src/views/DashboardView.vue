<script setup>
import Layout from '../components/Layout.vue'
import { useAuthStore } from '../stores/authStore'
import { 
  Users, 
  GitFork, 
  ShieldCheck, 
  Key, 
  Layers, 
  CheckCircle2 
} from 'lucide-vue-next'

const authStore = useAuthStore()
</script>

<template>
  <Layout>
    <div class="space-y-8">
      <!-- Welcome Header -->
      <div class="p-8 rounded-2xl bg-white dark:bg-gradient-to-r dark:from-[#121215] dark:via-[#18181c] dark:to-[#121215] border border-slate-200 dark:border-emerald-500/30 relative overflow-hidden shadow-xl dark:shadow-2xl transition-colors duration-300">
        <div class="absolute right-0 top-0 w-96 h-96 bg-emerald-500/10 rounded-full blur-3xl pointer-events-none"></div>
        <div class="relative z-10">
          <div class="inline-flex items-center space-x-2 px-3 py-1 rounded-full bg-emerald-100 dark:bg-emerald-950/80 border border-emerald-400/40 text-emerald-700 dark:text-emerald-400 text-xs font-mono mb-4">
            <ShieldCheck class="w-3.5 h-3.5" />
            <span>SESSION VALIDATED • GORM POSTGRES BACKEND</span>
          </div>
          <h2 class="text-3xl font-extrabold text-slate-900 dark:text-white">
            Welcome back, <span class="text-emerald-600 dark:text-emerald-400">{{ authStore.fullName }}</span>
          </h2>
          <p class="text-slate-600 dark:text-slate-400 text-sm mt-2 max-w-2xl">
            You are logged into the Next-Gen HRM Core Portal with role <strong class="text-slate-900 dark:text-white">{{ authStore.userRole }}</strong> and Access Level <strong class="text-emerald-600 dark:text-emerald-400">LVL {{ authStore.accessLevel }}</strong>.
          </p>
        </div>
      </div>

      <!-- Quick Metrics Cards -->
      <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
        <div class="p-6 rounded-xl bg-white dark:bg-[#121215] border border-slate-200 dark:border-[#27272a] hover:border-emerald-500/40 shadow-sm transition-all">
          <div class="flex items-center justify-between">
            <span class="text-xs font-semibold text-slate-500 dark:text-slate-400 uppercase tracking-wider">Your Access Level</span>
            <Key class="w-5 h-5 text-emerald-500 dark:text-emerald-400" />
          </div>
          <div class="text-3xl font-bold text-slate-900 dark:text-white mt-3 font-mono">LVL {{ authStore.accessLevel }}</div>
          <div class="text-xs text-emerald-600 dark:text-emerald-400/80 mt-1 font-mono">RBAC Permission Threshold</div>
        </div>

        <div class="p-6 rounded-xl bg-white dark:bg-[#121215] border border-slate-200 dark:border-[#27272a] hover:border-emerald-500/40 shadow-sm transition-all">
          <div class="flex items-center justify-between">
            <span class="text-xs font-semibold text-slate-500 dark:text-slate-400 uppercase tracking-wider">User Role</span>
            <ShieldCheck class="w-5 h-5 text-emerald-500 dark:text-emerald-400" />
          </div>
          <div class="text-3xl font-bold text-slate-900 dark:text-white mt-3 font-mono">{{ authStore.userRole }}</div>
          <div class="text-xs text-slate-500 dark:text-slate-400 mt-1">Assigned Role ID: {{ authStore.user?.role_id }}</div>
        </div>

        <div class="p-6 rounded-xl bg-white dark:bg-[#121215] border border-slate-200 dark:border-[#27272a] hover:border-emerald-500/40 shadow-sm transition-all">
          <div class="flex items-center justify-between">
            <span class="text-xs font-semibold text-slate-500 dark:text-slate-400 uppercase tracking-wider">Employee Directory</span>
            <Users class="w-5 h-5 text-emerald-500 dark:text-emerald-400" />
          </div>
          <div class="text-3xl font-bold text-slate-900 dark:text-white mt-3 font-mono">Accessible</div>
          <div class="text-xs text-slate-500 dark:text-slate-400 mt-1">Level 10+ Required</div>
        </div>

        <div class="p-6 rounded-xl bg-white dark:bg-[#121215] border border-slate-200 dark:border-[#27272a] hover:border-emerald-500/40 shadow-sm transition-all">
          <div class="flex items-center justify-between">
            <span class="text-xs font-semibold text-slate-500 dark:text-slate-400 uppercase tracking-wider">Company Payroll</span>
            <Layers class="w-5 h-5 text-emerald-500 dark:text-emerald-400" />
          </div>
          <div class="text-3xl font-bold mt-3 font-mono" :class="authStore.isHR ? 'text-emerald-600 dark:text-emerald-400' : 'text-slate-400 dark:text-slate-600'">
            {{ authStore.isHR ? 'Granted' : 'Locked' }}
          </div>
          <div class="text-xs text-slate-500 dark:text-slate-400 mt-1">Requires Level 80 (HR/Admin)</div>
        </div>
      </div>

      <!-- Feature Capabilities Matrix -->
      <div class="p-6 rounded-xl bg-white dark:bg-[#121215] border border-slate-200 dark:border-[#27272a] shadow-sm transition-colors duration-300">
        <h3 class="text-lg font-bold text-slate-900 dark:text-white mb-4 flex items-center gap-2">
          <CheckCircle2 class="w-5 h-5 text-emerald-500 dark:text-emerald-400" />
          <span>Active Session Capabilities Matrix</span>
        </h3>
        <div class="space-y-3">
          <div class="p-3 rounded-lg bg-slate-50 dark:bg-[#18181c] flex items-center justify-between text-sm">
            <span class="text-slate-700 dark:text-slate-300">View Employee Directory & Contact Profiles</span>
            <span class="px-2.5 py-1 rounded bg-emerald-100 dark:bg-emerald-950 text-emerald-700 dark:text-emerald-400 text-xs font-mono border border-emerald-400/30">PASS (LVL 10+)</span>
          </div>
          <div class="p-3 rounded-lg bg-slate-50 dark:bg-[#18181c] flex items-center justify-between text-sm">
            <span class="text-slate-700 dark:text-slate-300">View Organizational Tree & Recursive Hierarchy</span>
            <span class="px-2.5 py-1 rounded bg-emerald-100 dark:bg-emerald-950 text-emerald-700 dark:text-emerald-400 text-xs font-mono border border-emerald-400/30">PASS (LVL 10+)</span>
          </div>
          <div class="p-3 rounded-lg bg-slate-50 dark:bg-[#18181c] flex items-center justify-between text-sm">
            <span class="text-slate-700 dark:text-slate-300">Managerial Oversight & Direct Reports Access</span>
            <span :class="authStore.isManager ? 'bg-emerald-100 dark:bg-emerald-950 text-emerald-700 dark:text-emerald-400 border-emerald-400/30' : 'bg-slate-200 dark:bg-slate-900 text-slate-500 border-slate-300 dark:border-slate-800'" class="px-2.5 py-1 rounded text-xs font-mono border">
              {{ authStore.isManager ? 'PASS (LVL 50+)' : 'RESTRICTED (Requires LVL 50)' }}
            </span>
          </div>
          <div class="p-3 rounded-lg bg-slate-50 dark:bg-[#18181c] flex items-center justify-between text-sm">
            <span class="text-slate-700 dark:text-slate-300">Company Payroll & Financial Disbursal Reports</span>
            <span :class="authStore.isHR ? 'bg-emerald-100 dark:bg-emerald-950 text-emerald-700 dark:text-emerald-400 border-emerald-400/30' : 'bg-slate-200 dark:bg-slate-900 text-slate-500 border-slate-300 dark:border-slate-800'" class="px-2.5 py-1 rounded text-xs font-mono border">
              {{ authStore.isHR ? 'PASS (LVL 80+)' : 'RESTRICTED (Requires LVL 80)' }}
            </span>
          </div>
        </div>
      </div>
    </div>
  </Layout>
</template>
