<script setup>
import { ref, onMounted, computed } from 'vue'
import axios from 'axios'
import Layout from '../components/Layout.vue'
import { useAuthStore } from '../stores/authStore'
import { 
  Palmtree, 
  Calendar, 
  Clock, 
  CheckCircle, 
  XCircle, 
  PlusCircle, 
  RefreshCw, 
  X, 
  AlertCircle,
  FileText,
  UserCheck
} from 'lucide-vue-next'

const authStore = useAuthStore()

const loading = ref(true)
const balances = ref([])
const requests = ref([])
const leaveTypes = ref([])
const pendingRequests = ref([])

// Request Modal State
const showModal = ref(false)
const submitting = ref(false)
const formError = ref(null)
const formSuccess = ref(null)

const newRequest = ref({
  leave_type_id: null,
  start_date: '',
  end_date: '',
  reason: '',
})

// Calculated Days
const calculatedDays = computed(() => {
  if (!newRequest.value.start_date || !newRequest.value.end_date) return 0
  const start = new Date(newRequest.value.start_date)
  const end = new Date(newRequest.value.end_date)
  if (end < start) return 0
  const diffTime = Math.abs(end - start)
  return Math.ceil(diffTime / (1000 * 60 * 60 * 24)) + 1
})

const fetchMyLeaveData = async () => {
  loading.value = true
  try {
    const res = await axios.get('/api/leave/my-data', { withCredentials: true })
    balances.value = res.data.balances || []
    requests.value = res.data.requests || []
    leaveTypes.value = res.data.leave_types || []

    if (leaveTypes.value.length > 0 && !newRequest.value.leave_type_id) {
      newRequest.value.leave_type_id = leaveTypes.value[0].id
    }
  } catch (err) {
    console.error('Failed to fetch leave data:', err)
  } finally {
    loading.value = false
  }
}

const fetchPendingRequests = async () => {
  if (authStore.accessLevel < 50) return
  try {
    const res = await axios.get('/api/leave/pending', { withCredentials: true })
    pendingRequests.value = res.data.pending_requests || []
  } catch (err) {
    console.error('Failed to fetch pending requests:', err)
  }
}

const openModal = () => {
  formError.value = null
  formSuccess.value = null
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
}

const handleLeaveSubmit = async () => {
  submitting.value = true
  formError.value = null
  formSuccess.value = null

  try {
    const payload = {
      leave_type_id: Number(newRequest.value.leave_type_id),
      start_date: newRequest.value.start_date,
      end_date: newRequest.value.end_date,
      reason: newRequest.value.reason,
    }

    const res = await axios.post('/api/leave/request', payload, { withCredentials: true })
    formSuccess.value = res.data.message || 'Leave request submitted!'
    
    await fetchMyLeaveData()
    if (authStore.accessLevel >= 50) fetchPendingRequests()

    setTimeout(() => {
      closeModal()
      newRequest.value = {
        leave_type_id: leaveTypes.value.length > 0 ? leaveTypes.value[0].id : null,
        start_date: '',
        end_date: '',
        reason: '',
      }
    }, 1000)

  } catch (err) {
    formError.value = err.response?.data?.error || 'Failed to submit leave request.'
  } finally {
    submitting.value = false
  }
}

const handleApprove = async (id) => {
  try {
    await axios.put(`/api/leave/${id}/approve`, {}, { withCredentials: true })
    await fetchPendingRequests()
    await fetchMyLeaveData()
  } catch (err) {
    alert(err.response?.data?.error || 'Failed to approve request')
  }
}

const handleReject = async (id) => {
  try {
    await axios.put(`/api/leave/${id}/reject`, {}, { withCredentials: true })
    await fetchPendingRequests()
    await fetchMyLeaveData()
  } catch (err) {
    alert(err.response?.data?.error || 'Failed to reject request')
  }
}

const getStatusBadgeClass = (status) => {
  switch (status) {
    case 'Approved':
      return 'bg-emerald-100 dark:bg-emerald-950 text-emerald-700 dark:text-emerald-400 border-emerald-400/30'
    case 'Pending':
      return 'bg-amber-100 dark:bg-amber-950 text-amber-700 dark:text-amber-400 border-amber-400/30'
    case 'Rejected':
      return 'bg-red-100 dark:bg-red-950 text-red-700 dark:text-red-400 border-red-400/30'
    default:
      return 'bg-slate-100 dark:bg-slate-900 text-slate-700 dark:text-slate-400'
  }
}

onMounted(() => {
  fetchMyLeaveData()
  if (authStore.accessLevel >= 50) {
    fetchPendingRequests()
  }
})
</script>

<template>
  <Layout>
    <div class="space-y-6">
      <!-- Header -->
      <div class="flex items-center justify-between">
        <div>
          <h2 class="text-2xl font-bold text-slate-900 dark:text-white flex items-center gap-2">
            <Palmtree class="w-6 h-6 text-emerald-500 dark:text-emerald-400" />
            <span>Leave & Absence Engine</span>
          </h2>
          <p class="text-slate-500 dark:text-slate-400 text-sm mt-1">Leave Balances, Absence Requests & Manager Approval Queue</p>
        </div>

        <div class="flex items-center space-x-3">
          <button 
            @click="fetchMyLeaveData" 
            class="inline-flex items-center space-x-2 px-4 py-2 rounded-xl bg-white dark:bg-[#18181c] border border-slate-200 dark:border-[#27272a] hover:border-emerald-500/50 text-emerald-600 dark:text-emerald-400 text-sm font-medium transition-colors cursor-pointer shadow-sm"
          >
            <RefreshCw class="w-4 h-4" :class="{ 'animate-spin': loading }" />
            <span>Refresh</span>
          </button>

          <button 
            @click="openModal" 
            class="inline-flex items-center space-x-2 px-4 py-2 rounded-xl bg-gradient-to-r from-emerald-500 to-emerald-600 hover:from-emerald-400 hover:to-emerald-500 text-black font-bold text-sm shadow-[0_0_15px_rgba(16,185,129,0.3)] transition-all cursor-pointer"
          >
            <PlusCircle class="w-4 h-4" />
            <span>+ Request Leave</span>
          </button>
        </div>
      </div>

      <!-- Leave Balances Grid -->
      <div class="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-4 gap-4">
        <div 
          v-for="b in balances" 
          :key="b.id"
          class="p-5 rounded-2xl bg-white dark:bg-[#121215] border border-slate-200 dark:border-[#27272a] shadow-sm flex flex-col justify-between"
        >
          <div>
            <div class="flex items-center justify-between">
              <span class="text-xs font-bold text-slate-500 dark:text-slate-400 uppercase tracking-wider">{{ b.leave_type?.name }}</span>
              <div class="w-8 h-8 rounded-lg bg-emerald-100 dark:bg-emerald-950 text-emerald-600 dark:text-emerald-400 flex items-center justify-center font-bold text-xs">
                {{ b.allocated_days }}d
              </div>
            </div>
            <div class="mt-4">
              <div class="text-3xl font-extrabold text-slate-900 dark:text-white font-mono">
                {{ b.remaining_days }} <span class="text-xs font-normal text-slate-400">days left</span>
              </div>
            </div>
          </div>

          <div class="mt-4 pt-3 border-t border-slate-100 dark:border-[#27272a] flex justify-between text-xs text-slate-500 font-mono">
            <span>Used: <strong class="text-slate-900 dark:text-slate-200">{{ b.used_days }}d</strong></span>
            <span>Allocated: <strong class="text-slate-900 dark:text-slate-200">{{ b.allocated_days }}d</strong></span>
          </div>
        </div>
      </div>

      <!-- Manager Approval Queue (Access Level 50+) -->
      <div v-if="authStore.accessLevel >= 50 && pendingRequests.length > 0" class="p-6 rounded-2xl bg-amber-500/5 dark:bg-amber-500/10 border border-amber-500/30">
        <div class="flex items-center space-x-2 mb-4">
          <Clock class="w-5 h-5 text-amber-500" />
          <h3 class="text-lg font-bold text-slate-900 dark:text-white">Leave Approvals Queue ({{ pendingRequests.length }})</h3>
        </div>

        <div class="space-y-3">
          <div 
            v-for="req in pendingRequests" 
            :key="req.id"
            class="p-4 rounded-xl bg-white dark:bg-[#18181c] border border-slate-200 dark:border-[#27272a] flex flex-col md:flex-row md:items-center justify-between gap-4 shadow-sm"
          >
            <div>
              <div class="flex items-center space-x-2">
                <span class="font-bold text-slate-900 dark:text-white">{{ req.employee?.first_name }} {{ req.employee?.last_name }}</span>
                <span class="text-xs text-slate-500">({{ req.employee?.department?.name || 'Staff' }})</span>
                <span class="px-2 py-0.5 rounded bg-emerald-100 dark:bg-emerald-950 text-emerald-600 dark:text-emerald-400 text-xs font-mono font-semibold">
                  {{ req.leave_type?.name }}
                </span>
              </div>
              <p class="text-xs text-slate-500 mt-1 font-mono">
                Duration: <strong>{{ req.start_date.split('T')[0] }}</strong> to <strong>{{ req.end_date.split('T')[0] }}</strong> ({{ req.days_count }} days)
              </p>
              <p v-if="req.reason" class="text-xs text-slate-600 dark:text-slate-300 mt-1 italic">
                "{{ req.reason }}"
              </p>
            </div>

            <div class="flex items-center space-x-2">
              <button 
                @click="handleApprove(req.id)"
                class="px-3 py-1.5 rounded-lg bg-emerald-500 hover:bg-emerald-400 text-black font-bold text-xs flex items-center space-x-1 cursor-pointer transition-colors"
              >
                <CheckCircle class="w-3.5 h-3.5" />
                <span>Approve</span>
              </button>
              <button 
                @click="handleReject(req.id)"
                class="px-3 py-1.5 rounded-lg bg-red-600 hover:bg-red-500 text-white font-bold text-xs flex items-center space-x-1 cursor-pointer transition-colors"
              >
                <XCircle class="w-3.5 h-3.5" />
                <span>Reject</span>
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- My Leave Requests History Table -->
      <div class="p-6 rounded-2xl bg-white dark:bg-[#121215] border border-slate-200 dark:border-[#27272a] shadow-sm">
        <h3 class="text-lg font-bold text-slate-900 dark:text-white mb-4 flex items-center gap-2">
          <FileText class="w-5 h-5 text-emerald-500 dark:text-emerald-400" />
          <span>My Leave Requests History</span>
        </h3>

        <div v-if="loading" class="text-center py-8 text-slate-400 font-mono">Loading history...</div>
        <div v-else-if="requests.length === 0" class="text-center py-8 text-slate-400">No leave requests submitted yet.</div>

        <div v-else class="overflow-x-auto">
          <table class="w-full text-left text-sm text-slate-700 dark:text-slate-300">
            <thead class="text-xs uppercase bg-slate-100 dark:bg-[#18181c] text-emerald-700 dark:text-emerald-400 font-mono border-b border-slate-200 dark:border-[#27272a]">
              <tr>
                <th class="px-4 py-3">Submitted</th>
                <th class="px-4 py-3">Leave Type</th>
                <th class="px-4 py-3">Duration</th>
                <th class="px-4 py-3">Days</th>
                <th class="px-4 py-3">Reason</th>
                <th class="px-4 py-3">Status</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-200 dark:divide-[#27272a]">
              <tr v-for="req in requests" :key="req.id" class="hover:bg-slate-50 dark:hover:bg-[#18181c]/50 text-xs font-mono">
                <td class="px-4 py-3 text-slate-500">{{ req.created_at.split('T')[0] }}</td>
                <td class="px-4 py-3 font-sans font-semibold text-slate-900 dark:text-white">{{ req.leave_type?.name }}</td>
                <td class="px-4 py-3 text-slate-600 dark:text-slate-300">{{ req.start_date.split('T')[0] }} ➔ {{ req.end_date.split('T')[0] }}</td>
                <td class="px-4 py-3 font-bold text-slate-900 dark:text-white">{{ req.days_count }} days</td>
                <td class="px-4 py-3 font-sans text-slate-500 max-w-xs truncate">{{ req.reason || 'N/A' }}</td>
                <td class="px-4 py-3">
                  <span :class="getStatusBadgeClass(req.status)" class="px-2 py-0.5 rounded text-[10px] font-bold border">
                    {{ req.status }}
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- Modal for Submitting Leave Request -->
    <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4">
      <div class="w-full max-w-md bg-white dark:bg-[#121215] border border-slate-200 dark:border-[#27272a] rounded-2xl shadow-2xl p-6 relative">
        <button @click="closeModal" class="absolute top-5 right-5 text-slate-400 hover:text-slate-600 dark:hover:text-white cursor-pointer">
          <X class="w-5 h-5" />
        </button>

        <div class="flex items-center space-x-3 mb-6">
          <div class="w-10 h-10 rounded-xl bg-emerald-100 dark:bg-emerald-950 border border-emerald-400/40 text-emerald-600 dark:text-emerald-400 flex items-center justify-center">
            <Palmtree class="w-5 h-5" />
          </div>
          <div>
            <h3 class="text-xl font-bold text-slate-900 dark:text-white">Apply for Leave</h3>
            <p class="text-xs text-slate-500">Submit absence request to your manager</p>
          </div>
        </div>

        <div v-if="formError" class="mb-4 p-3 rounded-lg bg-red-50 dark:bg-red-950/40 text-red-600 text-xs flex items-center space-x-2">
          <AlertCircle class="w-4 h-4 shrink-0" />
          <span>{{ formError }}</span>
        </div>

        <div v-if="formSuccess" class="mb-4 p-3 rounded-lg bg-emerald-50 dark:bg-emerald-950/40 text-emerald-600 text-xs flex items-center space-x-2">
          <CheckCircle class="w-4 h-4 shrink-0" />
          <span>{{ formSuccess }}</span>
        </div>

        <form @submit.prevent="handleLeaveSubmit" class="space-y-4">
          <div>
            <label class="block text-xs font-semibold text-slate-600 dark:text-slate-300 uppercase tracking-wider mb-1">Leave Type *</label>
            <select v-model="newRequest.leave_type_id" required class="w-full px-3 py-2 rounded-lg bg-slate-50 dark:bg-[#18181c] border border-slate-200 dark:border-[#27272a] text-slate-900 dark:text-white text-sm focus:outline-none focus:border-emerald-500">
              <option v-for="lt in leaveTypes" :key="lt.id" :value="lt.id">{{ lt.name }} (Max: {{ lt.max_days_per_year }} days/yr)</option>
            </select>
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-xs font-semibold text-slate-600 dark:text-slate-300 uppercase tracking-wider mb-1">Start Date *</label>
              <input v-model="newRequest.start_date" type="date" required class="w-full px-3 py-2 rounded-lg bg-slate-50 dark:bg-[#18181c] border border-slate-200 dark:border-[#27272a] text-slate-900 dark:text-white text-sm focus:outline-none focus:border-emerald-500" />
            </div>

            <div>
              <label class="block text-xs font-semibold text-slate-600 dark:text-slate-300 uppercase tracking-wider mb-1">End Date *</label>
              <input v-model="newRequest.end_date" type="date" required class="w-full px-3 py-2 rounded-lg bg-slate-50 dark:bg-[#18181c] border border-slate-200 dark:border-[#27272a] text-slate-900 dark:text-white text-sm focus:outline-none focus:border-emerald-500" />
            </div>
          </div>

          <div v-if="calculatedDays > 0" class="p-3 rounded-lg bg-slate-100 dark:bg-[#18181c] text-xs font-mono flex justify-between items-center text-emerald-600 dark:text-emerald-400">
            <span>Total Duration:</span>
            <span class="font-bold text-sm">{{ calculatedDays }} day(s)</span>
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-600 dark:text-slate-300 uppercase tracking-wider mb-1">Reason for Leave</label>
            <textarea v-model="newRequest.reason" rows="3" placeholder="Brief explanation for absence..." class="w-full px-3 py-2 rounded-lg bg-slate-50 dark:bg-[#18181c] border border-slate-200 dark:border-[#27272a] text-slate-900 dark:text-white text-sm focus:outline-none focus:border-emerald-500"></textarea>
          </div>

          <div class="pt-4 flex justify-end space-x-3">
            <button type="button" @click="closeModal" class="px-4 py-2 rounded-xl border border-slate-200 dark:border-[#27272a] text-slate-600 text-sm font-medium hover:bg-slate-100 cursor-pointer">Cancel</button>
            <button type="submit" :disabled="submitting" class="px-5 py-2 rounded-xl bg-gradient-to-r from-emerald-500 to-emerald-600 hover:from-emerald-400 hover:to-emerald-500 text-black font-bold text-sm shadow-[0_0_15px_rgba(16,185,129,0.3)] disabled:opacity-50 cursor-pointer">
              {{ submitting ? 'Submitting...' : 'Submit Request' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </Layout>
</template>
