<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import Layout from '../components/Layout.vue'
import { ShieldAlert, Activity, RefreshCw, KeyRound, Globe, User } from 'lucide-vue-next'

const logs = ref([])
const loading = ref(true)

const fetchLogs = async () => {
  loading.value = true
  try {
    const res = await axios.get('/api/audit-logs', { withCredentials: true })
    logs.value = res.data.logs || []
  } catch (err) {
    console.error('Failed to fetch audit logs:', err)
  } finally {
    loading.value = false
  }
}

const formatDate = (isoString) => {
  if (!isoString) return ''
  const date = new Date(isoString)
  return date.toLocaleString()
}

const getActionBadgeClass = (action) => {
  switch (action) {
    case 'LOGIN_SUCCESS':
      return 'bg-emerald-100 dark:bg-emerald-950 text-emerald-700 dark:text-emerald-400 border-emerald-400/30'
    case 'TOKEN_REFRESH_SUCCESS':
      return 'bg-cyan-100 dark:bg-cyan-950 text-cyan-700 dark:text-cyan-400 border-cyan-400/30'
    case 'LOGIN_FAILED':
    case 'ACCESS_DENIED':
    case 'RATE_LIMIT_EXCEEDED':
      return 'bg-red-100 dark:bg-red-950 text-red-700 dark:text-red-400 border-red-400/30'
    case 'LOGOUT':
      return 'bg-slate-100 dark:bg-slate-900 text-slate-700 dark:text-slate-400 border-slate-300 dark:border-slate-800'
    default:
      return 'bg-purple-100 dark:bg-purple-950 text-purple-700 dark:text-purple-400 border-purple-400/30'
  }
}

onMounted(fetchLogs)
</script>

<template>
  <Layout>
    <div class="space-y-6">
      <!-- Header -->
      <div class="flex items-center justify-between">
        <div>
          <h2 class="text-2xl font-bold text-slate-900 dark:text-white flex items-center gap-2">
            <Activity class="w-6 h-6 text-emerald-500 dark:text-emerald-400" />
            <span>Security Audit Trail Logs</span>
          </h2>
          <p class="text-slate-500 dark:text-slate-400 text-sm mt-1">
            Real-time security events, session refreshes, and access control audit logs
          </p>
        </div>

        <button 
          @click="fetchLogs" 
          class="inline-flex items-center space-x-2 px-4 py-2 rounded-xl bg-white dark:bg-[#18181c] border border-slate-200 dark:border-[#27272a] hover:border-emerald-500/50 text-emerald-600 dark:text-emerald-400 text-sm font-medium transition-colors cursor-pointer shadow-sm"
        >
          <RefreshCw class="w-4 h-4" :class="{ 'animate-spin': loading }" />
          <span>Refresh Logs</span>
        </button>
      </div>

      <!-- Logs Table -->
      <div class="p-6 rounded-2xl bg-white dark:bg-[#121215] border border-slate-200 dark:border-[#27272a] shadow-sm">
        <div v-if="loading" class="text-center py-12 text-slate-400 font-mono">
          Loading security audit logs...
        </div>

        <div v-else-if="logs.length === 0" class="text-center py-12 text-slate-400">
          No audit logs recorded yet.
        </div>

        <div v-else class="overflow-x-auto">
          <table class="w-full text-left text-sm text-slate-700 dark:text-slate-300">
            <thead class="text-xs uppercase bg-slate-100 dark:bg-[#18181c] text-emerald-700 dark:text-emerald-400 font-mono border-b border-slate-200 dark:border-[#27272a]">
              <tr>
                <th class="px-4 py-3">Timestamp</th>
                <th class="px-4 py-3">Action Event</th>
                <th class="px-4 py-3">User / Identity</th>
                <th class="px-4 py-3">IP Address</th>
                <th class="px-4 py-3">Audit Details</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-200 dark:divide-[#27272a]">
              <tr v-for="log in logs" :key="log.id" class="hover:bg-slate-50 dark:hover:bg-[#18181c]/50 font-mono text-xs">
                <td class="px-4 py-3.5 text-slate-500 dark:text-slate-400 whitespace-nowrap">
                  {{ formatDate(log.created_at) }}
                </td>
                <td class="px-4 py-3.5">
                  <span :class="getActionBadgeClass(log.action)" class="px-2.5 py-1 rounded text-[11px] font-bold border">
                    {{ log.action }}
                  </span>
                </td>
                <td class="px-4 py-3.5 font-sans font-medium text-slate-900 dark:text-white">
                  {{ log.employee ? `${log.employee.first_name} ${log.employee.last_name}` : 'Anonymous / System' }}
                </td>
                <td class="px-4 py-3.5 text-slate-500 dark:text-slate-400">
                  {{ log.ip_address || '127.0.0.1' }}
                </td>
                <td class="px-4 py-3.5 font-sans text-slate-600 dark:text-slate-300">
                  {{ log.details }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </Layout>
</template>
