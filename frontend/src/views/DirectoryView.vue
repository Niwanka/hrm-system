<script setup>
import { ref, onMounted, computed } from 'vue'
import axios from 'axios'
import Layout from '../components/Layout.vue'
import { useAuthStore } from '../stores/authStore'
import { 
  Users, 
  Search, 
  Mail, 
  RefreshCw, 
  UserPlus, 
  X, 
  CheckCircle, 
  AlertCircle,
  User,
  Lock,
  Shield,
  UserCheck
} from 'lucide-vue-next'

const authStore = useAuthStore()

const employees = ref([])
const roles = ref([])
const loading = ref(true)
const searchQuery = ref('')
const selectedRoleFilter = ref('ALL')

// Modal state
const showModal = ref(false)
const submitting = ref(false)
const formError = ref(null)
const formSuccess = ref(null)

const newEmp = ref({
  first_name: '',
  last_name: '',
  email: '',
  password: 'password123',
  role_id: null,
  manager_id: null,
})

const fetchEmployees = async () => {
  loading.value = true
  try {
    const res = await axios.get('/api/employees', { withCredentials: true })
    employees.value = res.data.employees || []
  } catch (err) {
    console.error('Failed to load employee directory:', err)
  } finally {
    loading.value = false
  }
}

const fetchRoles = async () => {
  try {
    const res = await axios.get('/api/roles', { withCredentials: true })
    roles.value = res.data.roles || []
    if (roles.value.length > 0 && !newEmp.value.role_id) {
      // Default to Employee role or last role
      const empRole = roles.value.find(r => r.name === 'Employee') || roles.value[0]
      newEmp.value.role_id = empRole.id
    }
  } catch (err) {
    console.error('Failed to fetch roles:', err)
  }
}

const openModal = () => {
  formError.value = null
  formSuccess.value = null
  showModal.value = true
  fetchRoles()
}

const closeModal = () => {
  showModal.value = false
}

const handleCreateEmployee = async () => {
  submitting.value = true
  formError.value = null
  formSuccess.value = null

  try {
    const payload = {
      first_name: newEmp.value.first_name,
      last_name: newEmp.value.last_name,
      email: newEmp.value.email,
      password: newEmp.value.password,
      role_id: Number(newEmp.value.role_id),
      manager_id: newEmp.value.manager_id ? Number(newEmp.value.manager_id) : null,
    }

    const res = await axios.post('/api/employees', payload, { withCredentials: true })
    formSuccess.value = res.data.message || 'Employee inserted successfully!'
    
    // Refresh list & reset form
    await fetchEmployees()
    setTimeout(() => {
      closeModal()
      newEmp.value = {
        first_name: '',
        last_name: '',
        email: '',
        password: 'password123',
        role_id: roles.value.length > 0 ? roles.value[0].id : null,
        manager_id: null,
      }
    }, 1200)

  } catch (err) {
    formError.value = err.response?.data?.error || 'Failed to create employee. Please check inputs.'
  } finally {
    submitting.value = false
  }
}

const filteredEmployees = computed(() => {
  return employees.value.filter(emp => {
    const matchesSearch = 
      `${emp.first_name} ${emp.last_name}`.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      emp.email.toLowerCase().includes(searchQuery.value.toLowerCase())
    
    const matchesRole = selectedRoleFilter.value === 'ALL' || emp.role?.name === selectedRoleFilter.value
    return matchesSearch && matchesRole
  })
})

onMounted(() => {
  fetchEmployees()
  fetchRoles()
})
</script>

<template>
  <Layout>
    <div class="space-y-6">
      <!-- Header -->
      <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
        <div>
          <h2 class="text-2xl font-bold text-slate-900 dark:text-white flex items-center gap-2">
            <Users class="w-6 h-6 text-emerald-500 dark:text-emerald-400" />
            <span>Employee Directory</span>
          </h2>
          <p class="text-slate-500 dark:text-slate-400 text-sm mt-1">Browse all staff members, roles, and manager assignments</p>
        </div>

        <div class="flex items-center space-x-3">
          <button 
            @click="fetchEmployees" 
            class="inline-flex items-center space-x-2 px-4 py-2 rounded-xl bg-white dark:bg-[#18181c] border border-slate-200 dark:border-[#27272a] hover:border-emerald-500/50 text-emerald-600 dark:text-emerald-400 text-sm font-medium transition-colors cursor-pointer shadow-sm"
          >
            <RefreshCw class="w-4 h-4" :class="{ 'animate-spin': loading }" />
            <span>Refresh</span>
          </button>

          <!-- Add Employee Button (Requires Level 50+ Manager/HR/Admin) -->
          <button 
            v-if="authStore.accessLevel >= 50"
            @click="openModal" 
            class="inline-flex items-center space-x-2 px-4 py-2 rounded-xl bg-gradient-to-r from-emerald-500 to-emerald-600 hover:from-emerald-400 hover:to-emerald-500 text-black font-bold text-sm shadow-[0_0_15px_rgba(16,185,129,0.3)] transition-all cursor-pointer"
          >
            <UserPlus class="w-4 h-4" />
            <span>+ Add Employee</span>
          </button>
        </div>
      </div>

      <!-- Filters & Search -->
      <div class="p-4 rounded-xl bg-white dark:bg-[#121215] border border-slate-200 dark:border-[#27272a] flex flex-col md:flex-row gap-4 shadow-sm">
        <div class="relative flex-1">
          <Search class="w-4 h-4 text-slate-400 absolute left-3.5 top-3" />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Search employee by name or email..."
            class="w-full pl-10 pr-4 py-2 rounded-lg bg-slate-50 dark:bg-[#18181c] border border-slate-200 dark:border-[#27272a] text-slate-900 dark:text-white text-sm focus:outline-none focus:border-emerald-500 dark:focus:border-emerald-400"
          />
        </div>

        <div class="flex items-center space-x-2">
          <label class="text-xs text-slate-500 dark:text-slate-400 font-semibold uppercase">Role Filter:</label>
          <select 
            v-model="selectedRoleFilter"
            class="px-3 py-2 rounded-lg bg-slate-50 dark:bg-[#18181c] border border-slate-200 dark:border-[#27272a] text-slate-900 dark:text-white text-sm focus:outline-none focus:border-emerald-500 dark:focus:border-emerald-400"
          >
            <option value="ALL">All Roles</option>
            <option value="Admin">Admin (LVL 100)</option>
            <option value="HR">HR (LVL 80)</option>
            <option value="Manager">Manager (LVL 50)</option>
            <option value="Employee">Employee (LVL 10)</option>
          </select>
        </div>
      </div>

      <!-- Directory Cards Grid -->
      <div v-if="loading" class="text-center py-12 text-slate-400 font-mono">
        Loading employee directory...
      </div>

      <div v-else-if="filteredEmployees.length === 0" class="text-center py-12 text-slate-400">
        No employees found matching filter criteria.
      </div>

      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div 
          v-for="emp in filteredEmployees" 
          :key="emp.id"
          class="p-6 rounded-2xl bg-white dark:bg-[#121215] border border-slate-200 dark:border-[#27272a] hover:border-emerald-500/40 shadow-sm hover:shadow-md transition-all duration-200 group relative"
        >
          <div class="flex items-start justify-between">
            <div class="w-12 h-12 rounded-xl bg-emerald-100 dark:bg-emerald-950/80 border border-emerald-400/40 text-emerald-700 dark:text-emerald-400 flex items-center justify-center font-bold text-lg">
              {{ emp.first_name[0] }}{{ emp.last_name[0] }}
            </div>
            <span class="px-2.5 py-1 rounded-full bg-emerald-100 dark:bg-emerald-950 text-emerald-700 dark:text-emerald-400 text-xs font-mono font-semibold border border-emerald-400/30">
              {{ emp.role?.name || 'User' }} (LVL {{ emp.role?.access_level }})
            </span>
          </div>

          <div class="mt-4">
            <h3 class="text-lg font-bold text-slate-900 dark:text-white group-hover:text-emerald-600 dark:group-hover:text-emerald-400 transition-colors">
              {{ emp.first_name }} {{ emp.last_name }}
            </h3>
            <div class="flex items-center space-x-2 text-xs text-slate-500 dark:text-slate-400 mt-1">
              <Mail class="w-3.5 h-3.5 text-emerald-500 dark:text-emerald-400" />
              <span>{{ emp.email }}</span>
            </div>
          </div>

          <div class="mt-4 pt-4 border-t border-slate-100 dark:border-[#27272a]/60 text-xs flex justify-between items-center text-slate-500 dark:text-slate-400 font-mono">
            <span>Reporting Manager:</span>
            <span class="text-slate-900 dark:text-white font-sans font-medium">
              {{ emp.manager ? `${emp.manager.first_name} ${emp.manager.last_name}` : 'None (Top Level)' }}
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- Modal for Adding New Employee -->
    <div 
      v-if="showModal" 
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4"
    >
      <div class="w-full max-w-lg bg-white dark:bg-[#121215] border border-slate-200 dark:border-[#27272a] rounded-2xl shadow-2xl p-6 relative">
        <!-- Close Button -->
        <button 
          @click="closeModal" 
          class="absolute top-5 right-5 text-slate-400 hover:text-slate-600 dark:hover:text-white cursor-pointer"
        >
          <X class="w-5 h-5" />
        </button>

        <div class="flex items-center space-x-3 mb-6">
          <div class="w-10 h-10 rounded-xl bg-emerald-100 dark:bg-emerald-950 border border-emerald-400/40 text-emerald-600 dark:text-emerald-400 flex items-center justify-center">
            <UserPlus class="w-5 h-5" />
          </div>
          <div>
            <h3 class="text-xl font-bold text-slate-900 dark:text-white">Insert New Employee</h3>
            <p class="text-xs text-slate-500 dark:text-slate-400">Add a new staff member into the database</p>
          </div>
        </div>

        <!-- Success/Error Alert -->
        <div v-if="formError" class="mb-4 p-3 rounded-lg bg-red-50 dark:bg-red-950/40 border border-red-200 dark:border-red-900/60 text-red-600 dark:text-red-300 text-xs flex items-center space-x-2">
          <AlertCircle class="w-4 h-4 shrink-0" />
          <span>{{ formError }}</span>
        </div>

        <div v-if="formSuccess" class="mb-4 p-3 rounded-lg bg-emerald-50 dark:bg-emerald-950/40 border border-emerald-400/30 text-emerald-700 dark:text-emerald-300 text-xs flex items-center space-x-2">
          <CheckCircle class="w-4 h-4 shrink-0" />
          <span>{{ formSuccess }}</span>
        </div>

        <form @submit.prevent="handleCreateEmployee" class="space-y-4">
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-xs font-semibold text-slate-600 dark:text-slate-300 uppercase tracking-wider mb-1">
                First Name *
              </label>
              <input
                v-model="newEmp.first_name"
                type="text"
                required
                placeholder="John"
                class="w-full px-3 py-2 rounded-lg bg-slate-50 dark:bg-[#18181c] border border-slate-200 dark:border-[#27272a] text-slate-900 dark:text-white text-sm focus:outline-none focus:border-emerald-500"
              />
            </div>
            <div>
              <label class="block text-xs font-semibold text-slate-600 dark:text-slate-300 uppercase tracking-wider mb-1">
                Last Name *
              </label>
              <input
                v-model="newEmp.last_name"
                type="text"
                required
                placeholder="Doe"
                class="w-full px-3 py-2 rounded-lg bg-slate-50 dark:bg-[#18181c] border border-slate-200 dark:border-[#27272a] text-slate-900 dark:text-white text-sm focus:outline-none focus:border-emerald-500"
              />
            </div>
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-600 dark:text-slate-300 uppercase tracking-wider mb-1">
              Email Address *
            </label>
            <input
              v-model="newEmp.email"
              type="email"
              required
              placeholder="john.doe@company.com"
              class="w-full px-3 py-2 rounded-lg bg-slate-50 dark:bg-[#18181c] border border-slate-200 dark:border-[#27272a] text-slate-900 dark:text-white text-sm focus:outline-none focus:border-emerald-500"
            />
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-600 dark:text-slate-300 uppercase tracking-wider mb-1">
              Password *
            </label>
            <input
              v-model="newEmp.password"
              type="password"
              required
              placeholder="••••••••"
              class="w-full px-3 py-2 rounded-lg bg-slate-50 dark:bg-[#18181c] border border-slate-200 dark:border-[#27272a] text-slate-900 dark:text-white text-sm focus:outline-none focus:border-emerald-500"
            />
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-xs font-semibold text-slate-600 dark:text-slate-300 uppercase tracking-wider mb-1">
                Assigned Role *
              </label>
              <select 
                v-model="newEmp.role_id"
                required
                class="w-full px-3 py-2 rounded-lg bg-slate-50 dark:bg-[#18181c] border border-slate-200 dark:border-[#27272a] text-slate-900 dark:text-white text-sm focus:outline-none focus:border-emerald-500"
              >
                <option v-for="r in roles" :key="r.id" :value="r.id">
                  {{ r.name }} (LVL {{ r.access_level }})
                </option>
              </select>
            </div>

            <div>
              <label class="block text-xs font-semibold text-slate-600 dark:text-slate-300 uppercase tracking-wider mb-1">
                Reporting Manager
              </label>
              <select 
                v-model="newEmp.manager_id"
                class="w-full px-3 py-2 rounded-lg bg-slate-50 dark:bg-[#18181c] border border-slate-200 dark:border-[#27272a] text-slate-900 dark:text-white text-sm focus:outline-none focus:border-emerald-500"
              >
                <option :value="null">None (Top Level)</option>
                <option v-for="mgr in employees" :key="mgr.id" :value="mgr.id">
                  {{ mgr.first_name }} {{ mgr.last_name }} ({{ mgr.role?.name }})
                </option>
              </select>
            </div>
          </div>

          <div class="pt-4 flex justify-end space-x-3">
            <button
              type="button"
              @click="closeModal"
              class="px-4 py-2 rounded-xl border border-slate-200 dark:border-[#27272a] text-slate-600 dark:text-slate-300 text-sm font-medium hover:bg-slate-100 dark:hover:bg-[#18181c] cursor-pointer"
            >
              Cancel
            </button>
            <button
              type="submit"
              :disabled="submitting"
              class="px-5 py-2 rounded-xl bg-gradient-to-r from-emerald-500 to-emerald-600 hover:from-emerald-400 hover:to-emerald-500 text-black font-bold text-sm shadow-[0_0_15px_rgba(16,185,129,0.3)] disabled:opacity-50 cursor-pointer"
            >
              {{ submitting ? 'Inserting...' : 'Insert Employee' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </Layout>
</template>
