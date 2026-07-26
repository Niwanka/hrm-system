import { defineStore } from 'pinia'
import axios from 'axios'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: JSON.parse(localStorage.getItem('hrm_user') || 'null'),
    token: localStorage.getItem('hrm_token') || '',
    isLoading: false,
    error: null,
  }),

  getters: {
    isAuthenticated: (state) => !!state.user,
    accessLevel: (state) => state.user?.access_level || 0,
    userRole: (state) => state.user?.role_name || 'Guest',
    fullName: (state) => state.user ? `${state.user.first_name} ${state.user.last_name}` : '',
    
    // RBAC Permission Getters
    isAdmin: (state) => (state.user?.access_level || 0) >= 100,
    isHR: (state) => (state.user?.access_level || 0) >= 80,
    isManager: (state) => (state.user?.access_level || 0) >= 50,
    hasAccessLevel: (state) => (minLevel) => (state.user?.access_level || 0) >= minLevel,
  },

  actions: {
    async login(credentials) {
      this.isLoading = true
      this.error = null
      try {
        const response = await axios.post('/api/login', credentials, {
          withCredentials: true,
        })

        const { user, token } = response.data
        this.user = user
        this.token = token

        localStorage.setItem('hrm_user', JSON.stringify(user))
        localStorage.setItem('hrm_token', token)
        
        // Attach default Auth header for axios
        axios.defaults.headers.common['Authorization'] = `Bearer ${token}`

        return { success: true }
      } catch (err) {
        this.error = err.response?.data?.error || 'Invalid credentials or server error'
        return { success: false, error: this.error }
      } finally {
        this.isLoading = false
      }
    },

    async fetchCurrentUser() {
      if (!this.token && !this.user) return
      try {
        const response = await axios.get('/api/me', {
          headers: { Authorization: `Bearer ${this.token}` },
          withCredentials: true,
        })
        this.user = response.data.user
        localStorage.setItem('hrm_user', JSON.stringify(this.user))
      } catch (err) {
        console.warn('Session check failed or expired token', err)
        this.clearAuth()
      }
    },

    async logout() {
      try {
        await axios.post('/api/logout', {}, { withCredentials: true })
      } catch (err) {
        console.error('Logout error:', err)
      } finally {
        this.clearAuth()
      }
    },

    clearAuth() {
      this.user = null
      this.token = ''
      localStorage.removeItem('hrm_user')
      localStorage.removeItem('hrm_token')
      delete axios.defaults.headers.common['Authorization']
    },
  },
})
