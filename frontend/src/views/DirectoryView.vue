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
  Edit3, 
  Trash2, 
  X, 
  CheckCircle, 
  AlertCircle,
  Building,
  Briefcase
} from 'lucide-vue-next'

const authStore = useAuthStore()

const employees = ref([])
const roles = ref([])
const departments = ref([])
const loading = ref(true)
const searchQuery = ref('')
const selectedRoleFilter = ref('ALL')
const selectedDeptFilter = ref('ALL')

// Add Modal State
const showAddModal = ref(false)
const submittingAdd = ref(false)
const addFormError = ref(null)
const addFormSuccess = ref(null)

const newEmp = ref({
  first_name: '',
  last_name: '',
  email: '',
  password: 'password123',
  role_id: null,
  department_id: null,
  status: 'Active',
  manager_id: null,
})

// Edit Modal State
const showEditModal = ref(false)
const submittingEdit = ref(false)
const editFormError = ref(null)
const editFormSuccess = ref(null)

const editEmp = ref({
  id: null,
  first_name: '',
  last_name: '',
  email: '',
  role_id: null,
  department_id: null,
  status: 'Active',
  manager_id: null,
})

// Delete Confirmation Modal State
const showDeleteModal = ref(false)
const submittingDelete = ref(false)
const deletingEmp = ref(null)

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

const fetchRolesAndDepartments = async () => {
  try {
    const [rolesRes, deptsRes] = await Promise.all([
      axios.get('/api/roles', { withCredentials: true }),
      axios.get('/api/departments', { withCredentials: true })
    ])
    roles.value = rolesRes.data.roles || []
    departments.value = deptsRes.data.departments || []

    if (roles.value.length > 0 && !newEmp.value.role_id) {
      const empRole = roles.value.find(r => r.name === 'Employee') || roles.value[0]
      newEmp.value.role_id = empRole.id
    }
  } catch (err) {
    console.error('Failed to fetch roles or departments:', err)
  }
}

const openAddModal = () => {
  addFormError.value = null
  addFormSuccess.value = null
  showAddModal.value = true
  fetchRolesAndDepartments()
}

const closeAddModal = () => {
  showAddModal.value = false
}

const openEditModal = (emp) => {
  editFormError.value = null
  editFormSuccess.value = null
  editEmp.value = {
    id: emp.id,
    first_name: emp.first_name,
    last_name: emp.last_name,
    email: emp.email,
    role_id: emp.role_id,
    department_id: emp.department_id || null,
    status: emp.status || 'Active',
    manager_id: emp.manager_id || null,
  }
  showEditModal.value = true
  fetchRolesAndDepartments()
}

const closeEditModal = () => {
  showEditModal.value = false
}

const openDeleteModal = (emp) => {
  deletingEmp.value = emp
  showDeleteModal.value = true
}

const closeDeleteModal = () => {
  showDeleteModal.value = false
  deletingEmp.value = null
}

const handleCreateEmployee = async () => {
  submittingAdd.value = true
  addFormError.value = null
  addFormSuccess.value = null

  try {
    const payload = {
      first_name: newEmp.value.first_name,
      last_name: newEmp.value.last_name,
      email: newEmp.value.email,
      password: newEmp.value.password,
      role_id: Number(newEmp.value.role_id),
      department_id: newEmp.value.department_id ? Number(newEmp.value.department_id) : null,
      status: newEmp.value.status || 'Active',
      manager_id: newEmp.value.manager_id ? Number(newEmp.value.manager_id) : null,
    }

    const res = await axios.post('/api/employees', payload, { withCredentials: true })
    addFormSuccess.value = res.data.message || 'Employee inserted successfully!'
    
    await fetchEmployees()
    setTimeout(() => {
      closeAddModal()
      newEmp.value = {
        first_name: '',
        last_name: '',
        email: '',
        password: 'password123',
        role_id: roles.value.length > 0 ? roles.value[0].id : null,
        department_id: null,
        status: 'Active',
        manager_id: null,
      }
    }, 1000)

  } catch (err) {
    addFormError.value = err.response?.data?.error || 'Failed to create employee.'
  } finally {
    submittingAdd.value = false
  }
}

const handleUpdateEmployee = async () => {
  submittingEdit.value = true
  editFormError.value = null
  editFormSuccess.value = null

  try {
    const payload = {
      first_name: editEmp.value.first_name,
      last_name: editEmp.value.last_name,
      email: editEmp.value.email,
      role_id: Number(editEmp.value.role_id),
      department_id: editEmp.value.department_id ? Number(editEmp.value.department_id) : null,
      status: editEmp.value.status,
      manager_id: editEmp.value.manager_id ? Number(editEmp.value.manager_id) : null,
    }

    const res = await axios.put(`/api/employees/${editEmp.value.id}`, payload, { withCredentials: true })
    editFormSuccess.value = res.data.message || 'Employee updated successfully!'

    await fetchEmployees()
    setTimeout(() => {
      closeEditModal()
    }, 1000)

  } catch (err) {
    editFormError.value = err.response?.data?.error || 'Failed to update employee.'
  } finally {
    submittingEdit.value = false
  }
}

const handleDeleteEmployee = async () => {
  if (!deletingEmp.value) return
  submittingDelete.value = true

  try {
    await axios.delete(`/api/employees/${deletingEmp.value.id}`, { withCredentials: true })
    await fetchEmployees()
    closeDeleteModal()
  } catch (err) {
    alert(err.response?.data?.error || 'Failed to delete employee')
  } finally {
    submittingDelete.value = false
  }
}

const getStatusBadgeClass = (status) => {
  switch (status) {
    case 'Active':
      return 'bg-emerald-100 dark:bg-emerald-950 text-emerald-700 dark:text-emerald-400 border-emerald-400/30'
    case 'Onboarding':
      return 'bg-blue-100 dark:bg-blue-950 text-blue-700 dark:text-blue-400 border-blue-400/30'
    case 'On Leave':
      return 'bg-amber-100 dark:bg-amber-950 text-amber-700 dark:text-amber-400 border-amber-400/30'
    case 'Terminated':
      return 'bg-red-100 dark:bg-red-950 text-red-700 dark:text-red-400 border-red-400/30'
    default:
      return 'bg-slate-100 dark:bg-slate-900 text-slate-700 dark:text-slate-400 border-slate-300'
  }
}

const filteredEmployees = computed(() => {
  return employees.value.filter(emp => {
    const matchesSearch = 
      `${emp.first_name} ${emp.last_name}`.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      emp.email.toLowerCase().includes(searchQuery.value.toLowerCase())
    
    const matchesRole = selectedRoleFilter.value === 'ALL' || emp.role?.name === selectedRoleFilter.value
    const matchesDept = selectedDeptFilter.value === 'ALL' || emp.department?.name === selectedDeptFilter.value

    return matchesSearch && matchesRole && matchesDept
  })
})

onMounted(() => {
  fetchEmployees()
  fetchRolesAndDepartments()
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
          <p class="text-slate-500 dark:text-slate-400 text-sm mt-1">Browse all staff members, roles, departments, and manager assignments</p>
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
            @click="openAddModal" 
            class="inline-flex items-center space-x-2 px-4 py-2 rounded-xl bg-gradient-to-r from-emerald-500 to-emerald-600 hover:from-emerald-400 hover:to-emerald-500 text-black font-bold text-sm shadow-[0_0_15px_rgba(16,185,129,0.3)] transition-all cursor-pointer"
          >
            <UserPlus class="w-4 h-4" />
            <span>+ Add Employee</span>
          </button>
        </div>
      </div>

      <!-- Filters & Search Bar -->
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

        <div class="flex flex-wrap items-center gap-3">
          <div class="flex items-center space-x-2">
            <label class="text-xs text-slate-500 dark:text-slate-400 font-semibold uppercase">Role:</label>
            <select 
              v-model="selectedRoleFilter"
              class="px-3 py-2 rounded-lg bg-slate-50 dark:bg-[#18181c] border border-slate-200 dark:border-[#27272a] text-slate-900 dark:text-white text-sm focus:outline-none focus:border-emerald-500"
            >
              <option value="ALL">All Roles</option>
              <option value="Admin">Admin (LVL 100)</option>
              <option value="HR">HR (LVL 80)</option>
              <option value="Manager">Manager (LVL 50)</option>
              <option value="Employee">Employee (LVL 10)</option>
            </select>
          </div>

          <div class="flex items-center space-x-2">
            <label class="text-xs text-slate-500 dark:text-slate-400 font-semibold uppercase">Department:</label>
            <select 
              v-model="selectedDeptFilter"
              class="px-3 py-2 rounded-lg bg-slate-50 dark:bg-[#18181c] border border-slate-200 dark:border-[#27272a] text-slate-900 dark:text-white text-sm focus:outline-none focus:border-emerald-500"
            >
              <option value="ALL">All Departments</option>
              <option v-for="d in departments" :key="d.id" :value="d.name">
                {{ d.name }}
              </option>
            </select>
          </div>
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
          class="p-6 rounded-2xl bg-white dark:bg-[#121215] border border-slate-200 dark:border-[#27272a] hover:border-emerald-500/40 shadow-sm hover:shadow-md transition-all duration-200 group relative flex flex-col justify-between"
        >
          <div>
            <div class="flex items-start justify-between">
              <div class="w-12 h-12 rounded-xl bg-emerald-100 dark:bg-emerald-950/80 border border-emerald-400/40 text-emerald-700 dark:text-emerald-400 flex items-center justify-center font-bold text-lg">
                {{ emp.first_name[0] }}{{ emp.last_name[0] }}
              </div>
              
              <div class="flex flex-col items-end space-y-1">
                <span class="px-2.5 py-0.5 rounded-full bg-emerald-100 dark:bg-emerald-950 text-emerald-700 dark:text-emerald-400 text-xs font-mono font-semibold border border-emerald-400/30">
                  {{ emp.role?.name || 'User' }} (LVL {{ emp.role?.access_level }})
                </span>
                <span :class="getStatusBadgeClass(emp.status || 'Active')" class="px-2 py-0.5 rounded text-[10px] font-mono font-bold border">
                  {{ emp.status || 'Active' }}
                </span>
              </div>
            </div>

            <div class="mt-4">
              <h3 class="text-lg font-bold text-slate-900 dark:text-white group-hover:text-emerald-600 dark:group-hover:text-emerald-400 transition-colors">
                {{ emp.first_name }} {{ emp.last_name }}
              </h3>
              <div class="flex items-center space-x-2 text-xs text-slate-500 dark:text-slate-400 mt-1">
                <Mail class="w-3.5 h-3.5 text-emerald-500 dark:text-emerald-400" />
                <span>{{ emp.email }}</span>
              </div>
              <div v-if="emp.department" class="flex items-center space-x-2 text-xs text-slate-500 dark:text-slate-400 mt-1">
                <Building class="w-3.5 h-3.5 text-emerald-500 dark:text-emerald-400" />
                <span>{{ emp.department.name }}</span>
              </div>
            </div>
          </div>

          <div class="mt-4 pt-4 border-t border-slate-100 dark:border-[#27272a]/60 flex items-center justify-between">
            <div class="text-xs text-slate-500 dark:text-slate-400 font-mono">
              <span class="block text-[10px] text-slate-400">Reports To:</span>
              <span class="text-slate-900 dark:text-white font-sans font-medium">
                {{ emp.manager ? `${emp.manager.first_name} ${emp.manager.last_name}` : 'None (Top Level)' }}
              </span>
            </div>

            <!-- Edit & Delete Buttons for Level 50+ Manager/HR/Admin -->
            <div v-if="authStore.accessLevel >= 50" class="flex items-center space-x-2">
              <button 
                @click="openEditModal(emp)" 
                title="Edit Employee Profile"
                class="p-1.5 rounded-lg bg-slate-100 dark:bg-[#18181c] text-slate-600 dark:text-slate-300 hover:text-emerald-500 hover:border-emerald-500/50 border border-slate-200 dark:border-[#27272a] transition-colors cursor-pointer"
              >
                <Edit3 class="w-4 h-4" />
              </button>

              <button 
                v-if="authStore.accessLevel >= 80"
                @click="openDeleteModal(emp)" 
                title="Soft Delete / Terminate Employee"
                class="p-1.5 rounded-lg bg-red-50 dark:bg-red-950/40 text-red-600 dark:text-red-400 hover:bg-red-100 dark:hover:bg-red-900/60 border border-red-200 dark:border-red-900/40 transition-colors cursor-pointer"
              >
                <Trash2 class="w-4 h-4" />
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Modal for Adding New Employee -->
    <div v-if="showAddModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4">
      <div class="w-full max-w-lg bg-white dark:bg-[#121215] border border-slate-200 dark:border-[#27272a] rounded-2xl shadow-2xl p-6 relative">
        <button @click="closeAddModal" class="absolute top-5 right-5 text-slate-400 hover:text-slate-600 dark:hover:text-white cursor-pointer">
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

        <div v-if="addFormError" class="mb-4 p-3 rounded-lg bg-red-50 dark:bg-red-950/40 text-red-600 dark:text-red-300 text-xs flex items-center space-x-2">
          <AlertCircle class="w-4 h-4 shrink-0" />
          <span>{{ addFormError }}</span>
        </div>

        <div v-if="addFormSuccess" class="mb-4 p-3 rounded-lg bg-emerald-50 dark:bg-emerald-950/40 text-emerald-700 dark:text-emerald-300 text-xs flex items-center space-x-2">
          <CheckCircle class="w-4 h-4 shrink-0" />
          <span>{{ addFormSuccess }}</span>
        </div>

        <form @submit.prevent="handleCreateEmployee" class="space-y-4">
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-xs font-semibold text-slate-600 dark:text-slate-300 uppercase tracking-wider mb-1">First Name *</label>
              <input v-model="newEmp.first_name" type="text" required placeholder="John" class="w-full px-3 py-2 rounded-lg bg-slate-50 dark:bg-[#18181c] border border-slate-200 dark:border-[#27272a] text-slate-900 dark:text-white text-sm focus:outline-none focus:border-emerald-500" />
            </div>
            <div>
              <label class="block text-xs font-semibold text-slate-600 dark:text-slate-300 uppercase tracking-wider mb-1">Last Name *</label>
              <input v-model="newEmp.last_name" type="text" required placeholder="Doe" class="w-full px-3 py-2 rounded-lg bg-slate-50 dark:bg-[#18181c] border border-slate-200 dark:border-[#27272a] text-slate-900 dark:text-white text-sm focus:outline-none focus:border-emerald-500" />
            </div>
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-600 dark:text-slate-300 uppercase tracking-wider mb-1">Email Address *</label>
            <input v-model="newEmp.email" type="email" required placeholder="john.doe@company.com" class="w-full px-3 py-2 rounded-lg bg-slate-50 dark:bg-[#18181c] border border-slate-200 dark:border-[#27272a] text-slate-900 dark:text-white text-sm focus:outline-none focus:border-emerald-500" />
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-600 dark:text-slate-300 uppercase tracking-wider mb-1">Password *</label>
            <input v-model="newEmp.password" type="password" required placeholder="••••••••" class="w-full px-3 py-2 rounded-lg bg-slate-50 dark:bg-[#18181c] border border-slate-200 dark:border-[#27272a] text-slate-900 dark:text-white text-sm focus:outline-none focus:border-emerald-500" />
          </div>

          <div class="grid grid-cols-3 gap-4">
            <div>
              <label class="block text-xs font-semibold text-slate-600 dark:text-slate-300 uppercase tracking-wider mb-1">Role *</label>
              <select v-model="newEmp.role_id" required class="w-full px-3 py-2 rounded-lg bg-slate-50 dark:bg-[#18181c] border border-slate-200 dark:border-[#27272a] text-slate-900 dark:text-white text-sm focus:outline-none focus:border-emerald-500">
                <option v-for="r in roles" :key="r.id" :value="r.id">{{ r.name }}</option>
              </select>
            </div>

            <div>
              <label class="block text-xs font-semibold text-slate-600 dark:text-slate-300 uppercase tracking-wider mb-1">Department</label>
              <select v-model="newEmp.department_id" class="w-full px-3 py-2 rounded-lg bg-slate-50 dark:bg-[#18181c] border border-slate-200 dark:border-[#27272a] text-slate-900 dark:text-white text-sm focus:outline-none focus:border-emerald-500">
                <option :value="null">None</option>
                <option v-for="d in departments" :key="d.id" :value="d.id">{{ d.name }}</option>
              </select>
            </div>

            <div>
              <label class="block text-xs font-semibold text-slate-600 dark:text-slate-300 uppercase tracking-wider mb-1">Status</label>
              <select v-model="newEmp.status" class="w-full px-3 py-2 rounded-lg bg-slate-50 dark:bg-[#18181c] border border-slate-200 dark:border-[#27272a] text-slate-900 dark:text-white text-sm focus:outline-none focus:border-emerald-500">
                <option value="Active">Active</option>
                <option value="Onboarding">Onboarding</option>
                <option value="On Leave">On Leave</option>
                <option value="Terminated">Terminated</option>
              </select>
            </div>
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-600 dark:text-slate-300 uppercase tracking-wider mb-1">Reporting Manager</label>
            <select v-model="newEmp.manager_id" class="w-full px-3 py-2 rounded-lg bg-slate-50 dark:bg-[#18181c] border border-slate-200 dark:border-[#27272a] text-slate-900 dark:text-white text-sm focus:outline-none focus:border-emerald-500">
              <option :value="null">None (Top Level)</option>
              <option v-for="mgr in employees" :key="mgr.id" :value="mgr.id">{{ mgr.first_name }} {{ mgr.last_name }} ({{ mgr.role?.name }})</option>
            </select>
          </div>

          <div class="pt-4 flex justify-end space-x-3">
            <button type="button" @click="closeAddModal" class="px-4 py-2 rounded-xl border border-slate-200 dark:border-[#27272a] text-slate-600 dark:text-slate-300 text-sm font-medium hover:bg-slate-100 dark:hover:bg-[#18181c] cursor-pointer">Cancel</button>
            <button type="submit" :disabled="submittingAdd" class="px-5 py-2 rounded-xl bg-gradient-to-r from-emerald-500 to-emerald-600 hover:from-emerald-400 hover:to-emerald-500 text-black font-bold text-sm shadow-[0_0_15px_rgba(16,185,129,0.3)] disabled:opacity-50 cursor-pointer">
              {{ submittingAdd ? 'Inserting...' : 'Insert Employee' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Modal for Editing Employee -->
    <div v-if="showEditModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4">
      <div class="w-full max-w-lg bg-white dark:bg-[#121215] border border-slate-200 dark:border-[#27272a] rounded-2xl shadow-2xl p-6 relative">
        <button @click="closeEditModal" class="absolute top-5 right-5 text-slate-400 hover:text-slate-600 dark:hover:text-white cursor-pointer">
          <X class="w-5 h-5" />
        </button>

        <div class="flex items-center space-x-3 mb-6">
          <div class="w-10 h-10 rounded-xl bg-emerald-100 dark:bg-emerald-950 border border-emerald-400/40 text-emerald-600 dark:text-emerald-400 flex items-center justify-center">
            <Edit3 class="w-5 h-5" />
          </div>
          <div>
            <h3 class="text-xl font-bold text-slate-900 dark:text-white">Edit Employee Profile</h3>
            <p class="text-xs text-slate-500 dark:text-slate-400">Update employee details and manager assignments</p>
          </div>
        </div>

        <div v-if="editFormError" class="mb-4 p-3 rounded-lg bg-red-50 dark:bg-red-950/40 text-red-600 dark:text-red-300 text-xs flex items-center space-x-2">
          <AlertCircle class="w-4 h-4 shrink-0" />
          <span>{{ editFormError }}</span>
        </div>

        <div v-if="editFormSuccess" class="mb-4 p-3 rounded-lg bg-emerald-50 dark:bg-emerald-950/40 text-emerald-700 dark:text-emerald-300 text-xs flex items-center space-x-2">
          <CheckCircle class="w-4 h-4 shrink-0" />
          <span>{{ editFormSuccess }}</span>
        </div>

        <form @submit.prevent="handleUpdateEmployee" class="space-y-4">
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-xs font-semibold text-slate-600 dark:text-slate-300 uppercase tracking-wider mb-1">First Name *</label>
              <input v-model="editEmp.first_name" type="text" required class="w-full px-3 py-2 rounded-lg bg-slate-50 dark:bg-[#18181c] border border-slate-200 dark:border-[#27272a] text-slate-900 dark:text-white text-sm focus:outline-none focus:border-emerald-500" />
            </div>
            <div>
              <label class="block text-xs font-semibold text-slate-600 dark:text-slate-300 uppercase tracking-wider mb-1">Last Name *</label>
              <input v-model="editEmp.last_name" type="text" required class="w-full px-3 py-2 rounded-lg bg-slate-50 dark:bg-[#18181c] border border-slate-200 dark:border-[#27272a] text-slate-900 dark:text-white text-sm focus:outline-none focus:border-emerald-500" />
            </div>
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-600 dark:text-slate-300 uppercase tracking-wider mb-1">Email Address *</label>
            <input v-model="editEmp.email" type="email" required class="w-full px-3 py-2 rounded-lg bg-slate-50 dark:bg-[#18181c] border border-slate-200 dark:border-[#27272a] text-slate-900 dark:text-white text-sm focus:outline-none focus:border-emerald-500" />
          </div>

          <div class="grid grid-cols-3 gap-4">
            <div>
              <label class="block text-xs font-semibold text-slate-600 dark:text-slate-300 uppercase tracking-wider mb-1">Role *</label>
              <select v-model="editEmp.role_id" required class="w-full px-3 py-2 rounded-lg bg-slate-50 dark:bg-[#18181c] border border-slate-200 dark:border-[#27272a] text-slate-900 dark:text-white text-sm focus:outline-none focus:border-emerald-500">
                <option v-for="r in roles" :key="r.id" :value="r.id">{{ r.name }}</option>
              </select>
            </div>

            <div>
              <label class="block text-xs font-semibold text-slate-600 dark:text-slate-300 uppercase tracking-wider mb-1">Department</label>
              <select v-model="editEmp.department_id" class="w-full px-3 py-2 rounded-lg bg-slate-50 dark:bg-[#18181c] border border-slate-200 dark:border-[#27272a] text-slate-900 dark:text-white text-sm focus:outline-none focus:border-emerald-500">
                <option :value="null">None</option>
                <option v-for="d in departments" :key="d.id" :value="d.id">{{ d.name }}</option>
              </select>
            </div>

            <div>
              <label class="block text-xs font-semibold text-slate-600 dark:text-slate-300 uppercase tracking-wider mb-1">Status</label>
              <select v-model="editEmp.status" class="w-full px-3 py-2 rounded-lg bg-slate-50 dark:bg-[#18181c] border border-slate-200 dark:border-[#27272a] text-slate-900 dark:text-white text-sm focus:outline-none focus:border-emerald-500">
                <option value="Active">Active</option>
                <option value="Onboarding">Onboarding</option>
                <option value="On Leave">On Leave</option>
                <option value="Terminated">Terminated</option>
              </select>
            </div>
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-600 dark:text-slate-300 uppercase tracking-wider mb-1">Reporting Manager</label>
            <select v-model="editEmp.manager_id" class="w-full px-3 py-2 rounded-lg bg-slate-50 dark:bg-[#18181c] border border-slate-200 dark:border-[#27272a] text-slate-900 dark:text-white text-sm focus:outline-none focus:border-emerald-500">
              <option :value="null">None (Top Level)</option>
              <option v-for="mgr in employees" :key="mgr.id" :value="mgr.id" v-show="mgr.id !== editEmp.id">
                {{ mgr.first_name }} {{ mgr.last_name }} ({{ mgr.role?.name }})
              </option>
            </select>
          </div>

          <div class="pt-4 flex justify-end space-x-3">
            <button type="button" @click="closeEditModal" class="px-4 py-2 rounded-xl border border-slate-200 dark:border-[#27272a] text-slate-600 dark:text-slate-300 text-sm font-medium hover:bg-slate-100 dark:hover:bg-[#18181c] cursor-pointer">Cancel</button>
            <button type="submit" :disabled="submittingEdit" class="px-5 py-2 rounded-xl bg-gradient-to-r from-emerald-500 to-emerald-600 hover:from-emerald-400 hover:to-emerald-500 text-black font-bold text-sm shadow-[0_0_15px_rgba(16,185,129,0.3)] disabled:opacity-50 cursor-pointer">
              {{ submittingEdit ? 'Updating...' : 'Save Changes' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Modal for Soft Delete Confirmation -->
    <div v-if="showDeleteModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4">
      <div class="w-full max-w-md bg-white dark:bg-[#121215] border border-red-200 dark:border-red-900/60 rounded-2xl shadow-2xl p-6 relative">
        <div class="flex items-center space-x-3 mb-4">
          <div class="w-10 h-10 rounded-xl bg-red-100 dark:bg-red-950 text-red-600 dark:text-red-400 flex items-center justify-center">
            <Trash2 class="w-5 h-5" />
          </div>
          <div>
            <h3 class="text-lg font-bold text-slate-900 dark:text-white">Confirm Soft Delete</h3>
            <p class="text-xs text-slate-500 dark:text-slate-400">Preserves audit trail & history</p>
          </div>
        </div>

        <p class="text-sm text-slate-600 dark:text-slate-300 mb-6">
          Are you sure you want to soft-delete employee <strong class="text-slate-900 dark:text-white">{{ deletingEmp?.first_name }} {{ deletingEmp?.last_name }}</strong> (`{{ deletingEmp?.email }}`)?
        </p>

        <div class="flex justify-end space-x-3">
          <button @click="closeDeleteModal" type="button" class="px-4 py-2 rounded-xl border border-slate-200 dark:border-[#27272a] text-slate-600 dark:text-slate-300 text-sm font-medium hover:bg-slate-100 cursor-pointer">Cancel</button>
          <button @click="handleDeleteEmployee" :disabled="submittingDelete" type="button" class="px-5 py-2 rounded-xl bg-red-600 hover:bg-red-500 text-white font-bold text-sm shadow-md disabled:opacity-50 cursor-pointer">
            {{ submittingDelete ? 'Deleting...' : 'Soft Delete' }}
          </button>
        </div>
      </div>
    </div>
  </Layout>
</template>
