<script setup>
import { ref, onMounted, computed } from 'vue'
import axios from 'axios'
import Layout from '../components/Layout.vue'
import { useAuthStore } from '../stores/authStore'
import { 
  Clock, 
  Calendar, 
  CheckCircle, 
  AlertTriangle, 
  TrendingUp, 
  Users, 
  RefreshCw, 
  LogIn, 
  LogOut 
} from 'lucide-vue-next'

const authStore = useAuthStore()

const activeTab = ref('my') // 'my' or 'team'
const loading = ref(true)
const myLogs = ref([])
const teamLogs = ref([])

// Attendance Widget State
const todayState = ref({
  is_clocked_in: false,
  is_clocked_out: false,
  log: null,
})

const fetchTodayStatus = async () => {
  try {
    const res = await axios.get('/api/attendance/today', { withCredentials: true })
    todayState.value = res.data
  } catch (err) {
    console.error('Failed to fetch today status:', err)
  }
}

const fetchMyLogs = async () => {
  loading.value = true
  try {
    const res = await axios.get('/api/attendance/my-logs', { withCredentials: true })
    myLogs.value = res.data.logs || []
  } catch (err) {
    console.error('Failed to fetch attendance logs:', err)
  } finally {
    loading.value = false
  }
}

const fetchTeamLogs = async () => {
  if (authStore.accessLevel < 50) return
  try {
    const res = await axios.get('/api/attendance/team-logs', { withCredentials: true })
    teamLogs.value = res.data.team_logs || []
  } catch (err) {
    console.error('Failed to fetch team logs:', err)
  }
}

const handleClockIn = async () => {
  try {
    await axios.post('/api/attendance/clock-in', {}, { withCredentials: true })
    await fetchTodayStatus()
    await fetchMyLogs()
  } catch (err) {
    alert(err.response?.data?.error || 'Clock in failed')
  }
}

const handleClockOut = async () => {
  try {
    await axios.post('/api/attendance/clock-out', {}, { withCredentials: true })
    await fetchTodayStatus()
    await fetchMyLogs()
  } catch (err) {
    alert(err.response?.data?.error || 'Clock out failed')
  }
}

// Summary Metrics Computation
const metrics = computed(() => {
  const totalDays = myLogs.value.length
  const presentDays = myLogs.value.filter(l => l.status === 'Present').length
  const lateDays = myLogs.value.filter(l => l.status === 'Late').length
  const totalHours = myLogs.value.reduce((acc, l) => acc + (l.hours_worked || 0), 0)
  const avgHours = totalDays > 0 ? (totalHours / totalDays).toFixed(1) : 0

  return {
    totalDays,
    presentDays,
    lateDays,
    totalHours: totalHours.toFixed(1),
    avgHours,
  }
})

const getStatusBadgeClass = (status) => {
  switch (status) {
    case 'Present':
      return 'bg-emerald-100 dark:bg-emerald-950 text-emerald-700 dark:text-emerald-400 border-emerald-400/30'
    case 'Late':
      return 'bg-amber-100 dark:bg-amber-950 text-amber-700 dark:text-amber-400 border-amber-400/30'
    case 'Half Day':
      return 'bg-blue-100 dark:bg-blue-950 text-blue-700 dark:text-blue-400 border-blue-400/30'
    case 'Absent':
      return 'bg-red-100 dark:bg-red-950 text-red-700 dark:text-red-400 border-red-400/30'
    default:
      return 'bg-slate-100 dark:bg-slate-900 text-slate-700 dark:text-slate-400'
  }
}

const formatTime = (isoString) => {
  if (!isoString) return '--:--'
  return new Date(isoString).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}

onMounted(() => {
  fetchTodayStatus()
  fetchMyLogs()
  if (authStore.accessLevel >= 50) fetchTeamLogs()
})
</script>

<template>
  <Layout>
    <div class="space-y-6">
      <!-- Header -->
      <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
        <div>
          <h2 class="text-2xl font-bold text-slate-900 dark:text-white flex items-center gap-2">
            <Clock class="w-6 h-6 text-emerald-500 dark:text-emerald-400" />
            <span>Attendance & Timesheets</span>
          </h2>
          <p class="text-slate-500 dark:text-slate-400 text-sm mt-1">Track daily clock-in/out hours, monthly timesheets & team attendance</p>
        </div>

        <div class="flex items-center space-x-3">
          <button 
            @click="() => { fetchTodayStatus(); fetchMyLogs(); fetchTeamLogs(); }" 
            class="inline-flex items-center space-x-2 px-4 py-2 rounded-xl bg-white dark:bg-[#18181c] border border-slate-200 dark:border-[#27272a] hover:border-emerald-500/50 text-emerald-600 dark:text-emerald-400 text-sm font-medium transition-colors cursor-pointer shadow-sm"
          >
            <RefreshCw class="w-4 h-4" :class="{ 'animate-spin': loading }" />
            <span>Refresh</span>
          </button>

          <!-- 1-Click Clock In/Out Quick Action -->
          <button 
            v-if="!todayState.is_clocked_in && !todayState.is_clocked_out"
            @click="handleClockIn"
            class="inline-flex items-center space-x-2 px-4 py-2 rounded-xl bg-gradient-to-r from-emerald-500 to-emerald-600 hover:from-emerald-400 text-black font-extrabold text-sm shadow-[0_0_15px_rgba(16,185,129,0.3)] cursor-pointer"
          >
            <LogIn class="w-4 h-4" />
            <span>Clock In Now</span>
          </button>

          <button 
            v-else-if="todayState.is_clocked_in"
            @click="handleClockOut"
            class="inline-flex items-center space-x-2 px-4 py-2 rounded-xl bg-red-600 hover:bg-red-500 text-white font-extrabold text-sm shadow-md cursor-pointer"
          >
            <LogOut class="w-4 h-4" />
            <span>Clock Out</span>
          </button>
        </div>
      </div>

      <!-- Summary Metrics Grid -->
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div class="p-5 rounded-2xl bg-white dark:bg-[#121215] border border-slate-200 dark:border-[#27272a] shadow-sm">
          <div class="flex items-center justify-between text-slate-500 dark:text-slate-400 text-xs uppercase font-bold tracking-wider">
            <span>Days Present</span>
            <CheckCircle class="w-4 h-4 text-emerald-500" />
          </div>
          <div class="text-3xl font-extrabold text-slate-900 dark:text-white mt-3 font-mono">
            {{ metrics.presentDays }} <span class="text-xs text-slate-400 font-normal">/ {{ metrics.totalDays }} days</span>
          </div>
        </div>

        <div class="p-5 rounded-2xl bg-white dark:bg-[#121215] border border-slate-200 dark:border-[#27272a] shadow-sm">
          <div class="flex items-center justify-between text-slate-500 dark:text-slate-400 text-xs uppercase font-bold tracking-wider">
            <span>Late Arrivals</span>
            <AlertTriangle class="w-4 h-4 text-amber-500" />
          </div>
          <div class="text-3xl font-extrabold text-slate-900 dark:text-white mt-3 font-mono">
            {{ metrics.lateDays }} <span class="text-xs text-slate-400 font-normal">times</span>
          </div>
        </div>

        <div class="p-5 rounded-2xl bg-white dark:bg-[#121215] border border-slate-200 dark:border-[#27272a] shadow-sm">
          <div class="flex items-center justify-between text-slate-500 dark:text-slate-400 text-xs uppercase font-bold tracking-wider">
            <span>Total Hours Worked</span>
            <Clock class="w-4 h-4 text-emerald-500" />
          </div>
          <div class="text-3xl font-extrabold text-slate-900 dark:text-white mt-3 font-mono">
            {{ metrics.totalHours }} <span class="text-xs text-slate-400 font-normal">hrs</span>
          </div>
        </div>

        <div class="p-5 rounded-2xl bg-white dark:bg-[#121215] border border-slate-200 dark:border-[#27272a] shadow-sm">
          <div class="flex items-center justify-between text-slate-500 dark:text-slate-400 text-xs uppercase font-bold tracking-wider">
            <span>Avg Daily Hours</span>
            <TrendingUp class="w-4 h-4 text-emerald-500" />
          </div>
          <div class="text-3xl font-extrabold text-slate-900 dark:text-white mt-3 font-mono">
            {{ metrics.avgHours }} <span class="text-xs text-slate-400 font-normal">hrs/day</span>
          </div>
        </div>
      </div>

      <!-- Navigation Tabs (My Timesheet vs Team Overview) -->
      <div v-if="authStore.accessLevel >= 50" class="flex border-b border-slate-200 dark:border-[#27272a]">
        <button 
          @click="activeTab = 'my'" 
          :class="activeTab === 'my' ? 'border-emerald-500 text-emerald-600 dark:text-emerald-400 font-bold' : 'border-transparent text-slate-500 hover:text-slate-700 dark:hover:text-slate-300'"
          class="px-4 py-2.5 text-sm border-b-2 font-medium transition-colors cursor-pointer"
        >
          My Timesheet History
        </button>
        <button 
          @click="activeTab = 'team'" 
          :class="activeTab === 'team' ? 'border-emerald-500 text-emerald-600 dark:text-emerald-400 font-bold' : 'border-transparent text-slate-500 hover:text-slate-700 dark:hover:text-slate-300'"
          class="px-4 py-2.5 text-sm border-b-2 font-medium transition-colors cursor-pointer flex items-center space-x-2"
        >
          <Users class="w-4 h-4" />
          <span>Team Attendance Overview</span>
        </button>
      </div>

      <!-- My Timesheet Table -->
      <div v-if="activeTab === 'my'" class="p-6 rounded-2xl bg-white dark:bg-[#121215] border border-slate-200 dark:border-[#27272a] shadow-sm">
        <h3 class="text-lg font-bold text-slate-900 dark:text-white mb-4">My Monthly Timesheet Logs</h3>

        <div v-if="loading" class="text-center py-8 text-slate-400 font-mono">Loading timesheet records...</div>
        <div v-else-if="myLogs.length === 0" class="text-center py-8 text-slate-400">No clock-in records found. Click "Clock In Now" to log today's shift!</div>

        <div v-else class="overflow-x-auto">
          <table class="w-full text-left text-sm text-slate-700 dark:text-slate-300">
            <thead class="text-xs uppercase bg-slate-100 dark:bg-[#18181c] text-emerald-700 dark:text-emerald-400 font-mono border-b border-slate-200 dark:border-[#27272a]">
              <tr>
                <th class="px-4 py-3">Date</th>
                <th class="px-4 py-3">Clock In</th>
                <th class="px-4 py-3">Clock Out</th>
                <th class="px-4 py-3">Hours Worked</th>
                <th class="px-4 py-3">Status</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-200 dark:divide-[#27272a]">
              <tr v-for="log in myLogs" :key="log.id" class="hover:bg-slate-50 dark:hover:bg-[#18181c]/50 text-xs font-mono">
                <td class="px-4 py-3.5 font-bold text-slate-900 dark:text-white">{{ log.date.split('T')[0] }}</td>
                <td class="px-4 py-3.5 text-emerald-600 dark:text-emerald-400 font-semibold">{{ formatTime(log.clock_in) }}</td>
                <td class="px-4 py-3.5 text-slate-600 dark:text-slate-300">{{ formatTime(log.clock_out) }}</td>
                <td class="px-4 py-3.5 font-bold text-slate-900 dark:text-white">{{ log.hours_worked || 0 }} hrs</td>
                <td class="px-4 py-3.5">
                  <span :class="getStatusBadgeClass(log.status)" class="px-2 py-0.5 rounded text-[10px] font-bold border">
                    {{ log.status }}
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Team Overview Table (Manager/HR/Admin Level 50+) -->
      <div v-else-if="activeTab === 'team'" class="p-6 rounded-2xl bg-white dark:bg-[#121215] border border-slate-200 dark:border-[#27272a] shadow-sm">
        <h3 class="text-lg font-bold text-slate-900 dark:text-white mb-4">Team Attendance Logs Overview</h3>

        <div v-if="teamLogs.length === 0" class="text-center py-8 text-slate-400">No team clock-in logs found.</div>

        <div v-else class="overflow-x-auto">
          <table class="w-full text-left text-sm text-slate-700 dark:text-slate-300">
            <thead class="text-xs uppercase bg-slate-100 dark:bg-[#18181c] text-emerald-700 dark:text-emerald-400 font-mono border-b border-slate-200 dark:border-[#27272a]">
              <tr>
                <th class="px-4 py-3">Employee</th>
                <th class="px-4 py-3">Department</th>
                <th class="px-4 py-3">Date</th>
                <th class="px-4 py-3">Clock In</th>
                <th class="px-4 py-3">Clock Out</th>
                <th class="px-4 py-3">Hours</th>
                <th class="px-4 py-3">Status</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-200 dark:divide-[#27272a]">
              <tr v-for="log in teamLogs" :key="log.id" class="hover:bg-slate-50 dark:hover:bg-[#18181c]/50 text-xs font-mono">
                <td class="px-4 py-3.5 font-sans font-bold text-slate-900 dark:text-white">
                  {{ log.employee ? `${log.employee.first_name} ${log.employee.last_name}` : 'Unknown' }}
                </td>
                <td class="px-4 py-3.5 text-slate-500 font-sans">{{ log.employee?.department?.name || 'Staff' }}</td>
                <td class="px-4 py-3.5 text-slate-500">{{ log.date.split('T')[0] }}</td>
                <td class="px-4 py-3.5 text-emerald-600 dark:text-emerald-400 font-semibold">{{ formatTime(log.clock_in) }}</td>
                <td class="px-4 py-3.5 text-slate-600 dark:text-slate-300">{{ formatTime(log.clock_out) }}</td>
                <td class="px-4 py-3.5 font-bold text-slate-900 dark:text-white">{{ log.hours_worked || 0 }} hrs</td>
                <td class="px-4 py-3.5">
                  <span :class="getStatusBadgeClass(log.status)" class="px-2 py-0.5 rounded text-[10px] font-bold border">
                    {{ log.status }}
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </Layout>
</template>
