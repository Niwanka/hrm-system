<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import Layout from '../components/Layout.vue'
import { DollarSign, ShieldAlert, CheckCircle } from 'lucide-vue-next'

const payrollRecords = ref([])
const message = ref('')
const loading = ref(true)
const error = ref(null)

const fetchPayroll = async () => {
  loading.value = true
  error.value = null
  try {
    const res = await axios.get('/api/payroll', { withCredentials: true })
    payrollRecords.value = res.data.payroll_records || []
    message.value = res.data.message
  } catch (err) {
    error.value = err.response?.data?.error || 'Access Forbidden: Insufficient clearance level'
  } finally {
    loading.value = false
  }
}

onMounted(fetchPayroll)
</script>

<template>
  <Layout>
    <div class="space-y-6">
      <!-- Header -->
      <div class="flex items-center justify-between">
        <div>
          <h2 class="text-2xl font-bold text-slate-900 dark:text-white flex items-center gap-2">
            <DollarSign class="w-6 h-6 text-emerald-500 dark:text-emerald-400" />
            <span>Company Payroll Management</span>
          </h2>
          <p class="text-slate-500 dark:text-slate-400 text-sm mt-1">
            Restricted RBAC Protected Route (Minimum Access Level 80 required)
          </p>
        </div>
        <span class="px-3 py-1 rounded-full bg-emerald-100 dark:bg-emerald-950 text-emerald-700 dark:text-emerald-400 border border-emerald-400/30 text-xs font-mono font-bold">
          RESTRICTED CLEARANCE
        </span>
      </div>

      <!-- Access Denied Alert -->
      <div v-if="error" class="p-6 rounded-2xl bg-red-50 dark:bg-red-950/40 border border-red-200 dark:border-red-900/60 text-red-600 dark:text-red-300 flex items-start space-x-3">
        <ShieldAlert class="w-6 h-6 text-red-500 shrink-0 mt-0.5" />
        <div>
          <h3 class="font-bold text-slate-900 dark:text-white">403 Forbidden Access</h3>
          <p class="text-sm mt-1 text-red-600 dark:text-red-300">{{ error }}</p>
        </div>
      </div>

      <!-- Content Table -->
      <div v-else-if="!loading" class="p-6 rounded-2xl bg-white dark:bg-[#121215] border border-slate-200 dark:border-[#27272a] space-y-6 shadow-sm">
        <div class="p-4 rounded-xl bg-emerald-50 dark:bg-emerald-950/40 border border-emerald-400/30 text-emerald-800 dark:text-emerald-300 text-sm flex items-center space-x-2">
          <CheckCircle class="w-4 h-4 text-emerald-600 dark:text-emerald-400" />
          <span>{{ message }}</span>
        </div>

        <div class="overflow-x-auto">
          <table class="w-full text-left text-sm text-slate-700 dark:text-slate-300">
            <thead class="text-xs uppercase bg-slate-100 dark:bg-[#18181c] text-emerald-700 dark:text-emerald-400 font-mono border-b border-slate-200 dark:border-[#27272a]">
              <tr>
                <th class="px-4 py-3">Payroll Period</th>
                <th class="px-4 py-3">Total Disbursed</th>
                <th class="px-4 py-3">Status</th>
                <th class="px-4 py-3 text-right">Action</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-200 dark:divide-[#27272a]">
              <tr v-for="rec in payrollRecords" :key="rec.period" class="hover:bg-slate-50 dark:hover:bg-[#18181c]/50">
                <td class="px-4 py-3.5 font-bold text-slate-900 dark:text-white font-mono">{{ rec.period }}</td>
                <td class="px-4 py-3.5 text-emerald-600 dark:text-emerald-400 font-mono font-semibold">{{ rec.total_disbursed }}</td>
                <td class="px-4 py-3.5">
                  <span class="px-2 py-1 rounded text-xs font-mono bg-emerald-100 dark:bg-emerald-950 text-emerald-700 dark:text-emerald-400 border border-emerald-400/30">
                    {{ rec.status }}
                  </span>
                </td>
                <td class="px-4 py-3.5 text-right">
                  <button class="text-xs text-slate-500 hover:text-emerald-600 dark:hover:text-emerald-400 underline font-mono cursor-pointer">Download PDF</button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </Layout>
</template>
