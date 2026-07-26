<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import Layout from '../components/Layout.vue'
import OrgTreeNode from '../components/OrgTreeNode.vue'
import { GitFork, RefreshCw } from 'lucide-vue-next'

const hierarchyData = ref([])
const loading = ref(true)

const fetchHierarchy = async () => {
  loading.value = true
  try {
    const res = await axios.get('/api/employees/hierarchy', { withCredentials: true })
    hierarchyData.value = res.data.hierarchy || []
  } catch (err) {
    console.error('Failed to load hierarchy data:', err)
  } finally {
    loading.value = false
  }
}

onMounted(fetchHierarchy)
</script>

<template>
  <Layout>
    <div class="space-y-6">
      <!-- Header -->
      <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
        <div>
          <h2 class="text-2xl font-bold text-slate-900 dark:text-white flex items-center gap-2">
            <GitFork class="w-6 h-6 text-emerald-500 dark:text-emerald-400" />
            <span>Organizational Hierarchy Chart</span>
          </h2>
          <p class="text-slate-500 dark:text-slate-400 text-sm mt-1">
            Fetched via GORM preloaded tree query (`/api/employees/hierarchy`)
          </p>
        </div>

        <button 
          @click="fetchHierarchy" 
          class="inline-flex items-center space-x-2 px-4 py-2 rounded-xl bg-white dark:bg-[#18181c] border border-slate-200 dark:border-[#27272a] hover:border-emerald-500/50 text-emerald-600 dark:text-emerald-400 text-sm font-medium transition-colors cursor-pointer shadow-sm"
        >
          <RefreshCw class="w-4 h-4" :class="{ 'animate-spin': loading }" />
          <span>Reload Hierarchy</span>
        </button>
      </div>

      <!-- Tree Content -->
      <div class="p-8 rounded-2xl bg-white dark:bg-[#121215] border border-slate-200 dark:border-[#27272a] shadow-sm transition-colors duration-300">
        <div v-if="loading" class="text-center py-12 text-slate-400 font-mono">
          Loading organizational hierarchy tree...
        </div>

        <div v-else-if="hierarchyData.length === 0" class="text-center py-12 text-slate-400">
          No top-level management structure found.
        </div>

        <div v-else class="space-y-8">
          <OrgTreeNode 
            v-for="manager in hierarchyData" 
            :key="manager.id" 
            :employee="manager" 
          />
        </div>
      </div>
    </div>
  </Layout>
</template>
