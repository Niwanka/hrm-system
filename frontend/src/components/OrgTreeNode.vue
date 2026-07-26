<script setup>
import { ref } from 'vue'
import { ChevronDown, ChevronRight, User, Users } from 'lucide-vue-next'

const props = defineProps({
  employee: {
    type: Object,
    required: true,
  },
})

const isExpanded = ref(true)

const toggleExpand = () => {
  isExpanded.value = !isExpanded.value
}
</script>

<template>
  <div class="space-y-4">
    <!-- Manager Node Card -->
    <div class="inline-flex items-center space-x-3 p-4 rounded-xl bg-white dark:bg-[#121215] border border-slate-200 dark:border-emerald-500/40 shadow-sm dark:shadow-lg hover:border-emerald-500 transition-all duration-200">
      <button 
        v-if="employee.direct_reports && employee.direct_reports.length > 0"
        @click="toggleExpand"
        class="w-6 h-6 rounded bg-emerald-100 dark:bg-emerald-950 border border-emerald-400/40 text-emerald-600 dark:text-emerald-400 flex items-center justify-center cursor-pointer hover:bg-emerald-200 dark:hover:bg-emerald-900"
      >
        <ChevronDown v-if="isExpanded" class="w-4 h-4" />
        <ChevronRight v-else class="w-4 h-4" />
      </button>
      <div v-else class="w-6 h-6 flex items-center justify-center">
        <User class="w-4 h-4 text-slate-400" />
      </div>

      <div class="w-10 h-10 rounded-lg bg-emerald-100 dark:bg-emerald-950 text-emerald-700 dark:text-emerald-400 flex items-center justify-center font-bold text-sm border border-emerald-400/30">
        {{ employee.first_name?.[0] }}{{ employee.last_name?.[0] }}
      </div>

      <div>
        <h4 class="font-bold text-slate-900 dark:text-white text-sm flex items-center gap-2">
          <span>{{ employee.first_name }} {{ employee.last_name }}</span>
          <span class="px-2 py-0.5 rounded bg-emerald-100 dark:bg-emerald-950 border border-emerald-400/30 text-emerald-700 dark:text-emerald-400 text-[10px] font-mono">
            {{ employee.role?.name || 'Manager' }} (LVL {{ employee.role?.access_level }})
          </span>
        </h4>
        <p class="text-xs text-slate-500 dark:text-slate-400">{{ employee.email }}</p>
      </div>

      <div v-if="employee.direct_reports && employee.direct_reports.length > 0" class="pl-2 border-l border-slate-200 dark:border-[#27272a] text-xs text-emerald-600 dark:text-emerald-400 font-mono flex items-center space-x-1">
        <Users class="w-3.5 h-3.5" />
        <span>{{ employee.direct_reports.length }} Reports</span>
      </div>
    </div>

    <!-- Recursive Direct Reports Children -->
    <div 
      v-if="isExpanded && employee.direct_reports && employee.direct_reports.length > 0"
      class="pl-8 border-l-2 border-emerald-400/30 dark:border-emerald-500/20 space-y-4 ml-4"
    >
      <OrgTreeNode 
        v-for="child in employee.direct_reports" 
        :key="child.id"
        :employee="child"
      />
    </div>
  </div>
</template>
