import { defineStore } from 'pinia'

export const useThemeStore = defineStore('theme', {
  state: () => ({
    isDark: localStorage.getItem('hrm_theme') !== 'light',
  }),

  actions: {
    initTheme() {
      if (this.isDark) {
        document.documentElement.classList.add('dark')
      } else {
        document.documentElement.classList.remove('dark')
      }
    },

    toggleTheme() {
      this.isDark = !this.isDark
      if (this.isDark) {
        document.documentElement.classList.add('dark')
        localStorage.setItem('hrm_theme', 'dark')
      } else {
        document.documentElement.classList.remove('dark')
        localStorage.setItem('hrm_theme', 'light')
      }
    },
  },
})
