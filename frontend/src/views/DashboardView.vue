<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import Layout from '../components/Layout.vue'
import { useAuthStore } from '../stores/authStore'
import { 
  Users, 
  GitFork, 
  ShieldCheck, 
  Key, 
  Layers, 
  CheckCircle2,
  Clock,
  LogIn,
  LogOut,
  CheckCircle,
  AlertCircle
} from 'lucide-vue-next'

const authStore = useAuthStore()

// Attendance Widget State
const attendanceState = ref({
  loading: true,
  is_clocked_in: false,
  is_clocked_out: false,
  log: null,
  message: null,
  error: null,
})

const currentTime = ref(new Date().toLocaleTimeString())
setInterval(() => {
  currentTime.value = new Date().toLocaleTimeString()
}, 1000)

const fetchTodayAttendance = async () => {
  attendanceState.value.loading = true
  try {
    const res = await axios.get('/api/attendance/today', { withCredentials: true })
    attendanceState.value.is_clocked_in = res.data.is_clocked_in
    attendanceState.value.is_clocked_out = res.data.is_clocked_out
    attendanceState.value.log = res.data.log
  } catch (err) {
    console.error('Failed to fetch today attendance:', err)
  } finally {
    attendanceState.value.loading = false
  }
}

const handleClockIn = async () => {
  attendanceState.value.message = null
  attendanceState.value.error = null
  try {
    const res = await axios.post('/api/attendance/clock-in', {}, { withCredentials: true })
    attendanceState.value.message = res.data.message
    await fetchTodayAttendance()
  } catch (err) {
    attendanceState.value.error = err.response?.data?.error || 'Clock in failed.'
  }
}

const handleClockOut = async () => {
  attendanceState.value.message = null
  attendanceState.value.error = null
  try {
    const res = await axios.post('/api/attendance/clock-out', {}, { withCredentials: true })
    attendanceState.value.message = res.data.message
    await fetchTodayAttendance()
  } catch (err) {
    attendanceState.value.error = err.response?.data?.error || 'Clock out failed.'
  }
}

onMounted(() => {
  fetchTodayAttendance()
})
</script>

<template>
  <Layout>
    <div class="space-y-8">
      <!-- Welcome Header -->
      <div class="p-8 rounded-2xl bg-white dark:bg-gradient-to-r dark:from-[#121215] dark:via-[#18181c] dark:to-[#121215] border border-slate-200 dark:border-emerald-500/30 relative overflow-hidden shadow-xl dark:shadow-2xl transition-colors duration-300 flex flex-col lg:flex-row lg:items-center justify-between gap-6">
        <div class="absolute right-0 top-0 w-96 h-96 bg-emerald-500/10 rounded-full blur-3xl pointer-events-none"></div>
        <div class="relative z-10">
          <div class="inline-flex items-center space-x-2 px-3 py-1 rounded-full bg-emerald-100 dark:bg-emerald-950/80 border border-emerald-400/40 text-emerald-700 dark:text-emerald-400 text-xs font-mono mb-4">
            <ShieldCheck class="w-3.5 h-3.5" />
            <span>SESSION VALIDATED • GORM POSTGRES BACKEND</span>
          </div>
          <h2 class="text-3xl font-extrabold text-slate-900 dark:text-white">
            Welcome back, <span class="text-emerald-600 dark:text-emerald-400">{{ authStore.fullName }}</span>
          </h2>
          <p class="text-slate-600 dark:text-slate-400 text-sm mt-2 max-w-xl">
            You are logged into the Next-Gen HRM Core Portal with role <strong class="text-slate-900 dark:text-white">{{ authStore.userRole }}</strong> and Access Level <strong class="text-emerald-600 dark:text-emerald-400">LVL {{ authStore.accessLevel }}</strong>.
          </p>
        </div>

        <!-- 1-Click Attendance Widget -->
        <div class="relative z-10 p-5 rounded-2xl bg-slate-50 dark:bg-[#18181c] border border-slate-200 dark:border-[#27272a] shadow-inner min-w-[300px]">
          <div class="flex items-center justify-between mb-3">
            <span class="text-xs font-bold text-slate-500 dark:text-slate-400 uppercase tracking-wider flex items-center gap-1.5">
              <Clock class="w-4 h-4 text-emerald-500" />
              <span>Daily Attendance</span>
            </span>
            <span class="text-xs font-mono text-emerald-600 dark:text-emerald-400 font-bold">{{ currentTime }}</span>
          </div>

          <div v-if="attendanceState.message" class="mb-3 p-2 rounded bg-emerald-100 dark:bg-emerald-950 text-emerald-700 dark:text-emerald-300 text-[11px]">
            {{ attendanceState.message }}
          </div>

          <div v-if="attendanceState.error" class="mb-3 p-2 rounded bg-red-100 dark:bg-red-950 text-red-700 dark:text-red-300 text-[11px]">
            {{ attendanceState.error }}
          </div>

          <div v-if="attendanceState.loading" class="text-xs text-slate-400 font-mono py-2 text-center">
            Checking attendance status...
          </div>

          <div v-else class="space-y-3">
            <div v-if="!attendanceState.is_clocked_in && !attendanceState.is_clocked_out" class="space-y-2">
              <button 
                @click="handleClockIn"
                class="w-full py-2.5 px-4 rounded-xl bg-gradient-to-r from-emerald-500 to-emerald-600 hover:from-emerald-400 hover:to-emerald-500 text-black font-extrabold text-sm shadow-[0_0_15px_rgba(16,185,129,0.3)] transition-all flex items-center justify-center space-x-2 cursor-pointer"
              >
                <LogIn class="w-4 h-4" />
                <span>Clock In Now</span>
              </button>
            </div>

            <div v-else-if="attendanceState.is_clocked_in" class="space-y-2">
              <div class="p-2 rounded bg-emerald-50 dark:bg-emerald-950/60 text-emerald-700 dark:text-emerald-400 text-xs font-mono flex justify-between items-center border border-emerald-400/30">
                <span>Clocked In:</span>
                <strong class="font-sans font-bold">{{ attendanceState.log?.clock_in ? new Date(attendanceState.log.clock_in).toLocaleTimeString() : '' }}</strong>
              </div>
              <button 
                @click="handleClockOut"
                class="w-full py-2 px-4 rounded-xl bg-red-600 hover:bg-red-500 text-white font-bold text-sm shadow-md transition-all flex items-center justify-center space-x-2 cursor-pointer"
              >
                <LogOut class="w-4 h-4" />
                <span>Clock Out</span>
              </button>
            </div>

            <div v-else-if="attendanceState.is_clocked_out" class="p-3 rounded-xl bg-emerald-100 dark:bg-emerald-950/80 border border-emerald-400/40 text-emerald-700 dark:text-emerald-400 text-xs font-mono text-center">
              <div class="font-bold">Shift Completed Today 🎉</div>
              <div class="mt-1 text-[11px]">Total: <strong>{{ attendanceState.log?.hours_worked }} hrs</strong></div>
            </div>
          </div>
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
            <span class="text-slate-700 dark:text-slate-300">Leave & Absence Request Management</span>
            <span class="px-2.5 py-1 rounded bg-emerald-100 dark:bg-emerald-950 text-emerald-700 dark:text-emerald-400 text-xs font-mono border border-emerald-400/30">PASS (LVL 10+)</span>
          </div>
          <div class="p-3 rounded-lg bg-slate-50 dark:bg-[#18181c] flex items-center justify-between text-sm">
            <span class="text-slate-700 dark:text-slate-300">Managerial Oversight & Leave Approvals Queue</span>
            <span :class="authStore.isManager ? 'bg-emerald-100 dark:bg-emerald-950 text-emerald-700 dark:text-emerald-400 border-emerald-400/30' : 'bg-slate-200 dark:bg-slate-900 text-slate-500 border-slate-300 dark:border-slate-800'" class="px-2.5 py-1 rounded text-xs font-mono border">
              {{ authStore.isManager ? 'PASS (LVL 50+)' : 'RESTRICTED (Requires LVL 50)' }}
            </span>
          </div>
        </div>
      </div>
    </div>
  </Layout>
</template>
